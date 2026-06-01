"""
提醒道具掃價 Bot

從後端取得啟用中的價格提醒道具清單，對每個道具到拍賣場查詢最低價並寫入 DB。
後端在寫入價格後會自動比對門檻並發送通知。

使用方式：
  python bot.py
  python bot.py --equip-region 1427,412,1713,501 --default-region 1722,415,1975,503

設定方式（config.py 或環境變數）：
  SCHEDULE_INTERVAL_MINUTES  掃描間隔（分鐘，預設 5）
"""

import argparse
import logging
import sys
import time
from datetime import datetime, timedelta

import requests
import schedule

from config import API_BASE_URL, BETWEEN_ITEMS_DELAY, SCHEDULE_INTERVAL_MINUTES
from notify import send_message, build_alert_message
from scraper import get_game_window, scrape_item, verify_price_header

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s  %(levelname)s  %(message)s",
    datefmt="%H:%M:%S",
    handlers=[
        logging.StreamHandler(sys.stdout),
        logging.FileHandler("bot.log", encoding="utf-8"),
    ],
)
logger = logging.getLogger(__name__)


class AlertBot:
    def __init__(self, equip_region=None, default_region=None) -> None:
        self._win = None
        self._equip_region = equip_region
        self._default_region = default_region

    # ------------------------------------------------------------------
    # 私有：API 操作
    # ------------------------------------------------------------------

    def _fetch_alert_items(self) -> list[dict]:
        """取得所有啟用中提醒的道具清單（不需登入）。"""
        try:
            r = requests.get(
                f"{API_BASE_URL}/api/bot/alert-items",
                timeout=10,
            )
            r.raise_for_status()
        except requests.RequestException as e:
            logger.error(f"[Bot] 取得提醒道具失敗：{e}")
            return []

        resp = r.json()
        items = resp.get("data", resp) if isinstance(resp, dict) else resp
        return items if isinstance(items, list) else []

    def _record_price(self, item_id: int, price: int) -> bool:
        """將最低價寫入後端（後端自動比對門檻並發送通知）。"""
        try:
            r = requests.post(
                f"{API_BASE_URL}/api/items/{item_id}/prices",
                json={"price": float(price)},
                headers={"X-User-ID": "alert-bot"},
                timeout=10,
            )
            return r.ok
        except requests.RequestException as e:
            logger.error(f"[Bot] 寫入價格失敗：{e}")
            return False

    # ------------------------------------------------------------------
    # 私有：單次掃描
    # ------------------------------------------------------------------

    def _run_once(self) -> None:
        logger.info("=" * 50)
        logger.info("[Bot] 開始掃描提醒道具")

        items = self._fetch_alert_items()
        if not items:
            logger.info("[Bot] 目前無啟用中的提醒，結束")
            return

        logger.info(f"[Bot] 共 {len(items)} 個道具待掃描")

        # 初始化遊戲視窗
        try:
            if self._win is None:
                self._win = get_game_window()
            verify_price_header(self._win)
        except RuntimeError as e:
            logger.error(f"[Bot] 視窗錯誤：{e}")
            self._win = None
            return

        ok, fail = 0, []

        total = len(items)
        for idx, item in enumerate(items, 1):
            name            = item.get("item_name", "")
            item_id         = item.get("item_id")
            item_type       = item.get("item_type", 1)
            search_mode     = item.get("search_mode", 1)
            english_name    = item.get("english_name", "")
            threshold_price = item.get("threshold_price", 0)
            bot_id          = item.get("bot_id")
            search_name     = english_name if search_mode == 2 and english_name else name

            if not name or item_id is None:
                continue

            if search_name != name:
                logger.info(f"  使用英文名稱查詢：{search_name}")
            try:
                price = scrape_item(self._win, search_name, item_type,
                                   equip_region=self._equip_region,
                                   default_region=self._default_region)

                if price is None:
                    logger.warning(f"▶ [{idx}/{total}] {name} → 找不到價格，跳過")
                    fail.append(name)
                else:
                    if self._record_price(item_id, price):
                        logger.info(f"▶ [{idx}/{total}] {name} → {price:,} → 已寫入 DB")
                        ok += 1
                    else:
                        logger.warning(f"▶ [{idx}/{total}] {name} → {price:,} → 寫入 DB 失敗")
                        fail.append(name)

                    if threshold_price > 0 and price <= threshold_price:
                        now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                        msg = build_alert_message(name, price, threshold_price, now)
                        if send_message(msg, bot_id=bot_id):
                            logger.info(f"  已發送價格通知：{name} → {price:,}")
                        else:
                            logger.warning(f"  價格符合門檻但通知發送失敗：{name} → {price:,}（請確認後台提醒有設定通知機器人）")

            except Exception as e:
                logger.error(f"▶ [{idx}/{total}] {name} → 錯誤：{e}")
                fail.append(name)

            time.sleep(BETWEEN_ITEMS_DELAY)

        logger.info(f"[Bot] 完成：{ok}/{len(items)} 筆成功")
        if fail:
            logger.warning(f"[Bot] 失敗項目（{len(fail)} 筆）：{', '.join(fail)}")
            msg = f"⚠️ Bot 掃描完成：{ok}/{len(items)} 筆成功\n❌ 失敗（{len(fail)} 筆）：{', '.join(fail)}"
            send_message(msg)

        next_run = datetime.now() + timedelta(minutes=SCHEDULE_INTERVAL_MINUTES)
        logger.info(f"[Bot] 下次更新時間：{next_run.strftime('%H:%M:%S')}")
        logger.info("=" * 50)

    # ------------------------------------------------------------------
    # 公開：啟動
    # ------------------------------------------------------------------

    def start(self) -> None:
        """立即執行一次，之後每 SCHEDULE_INTERVAL_MINUTES 分鐘自動觸發。"""
        logger.info(f"[Bot] 啟動，間隔：{SCHEDULE_INTERVAL_MINUTES} 分鐘")
        logger.info("[Bot] 按 Ctrl+C 停止")

        self._run_once()

        schedule.every(SCHEDULE_INTERVAL_MINUTES).minutes.do(self._run_once)

        try:
            while True:
                schedule.run_pending()
                time.sleep(30)
        except KeyboardInterrupt:
            logger.info("[Bot] 已停止")


def _parse_region(s: str) -> tuple[int, int, int, int]:
    parts = [int(v.strip()) for v in s.split(",")]
    if len(parts) != 4:
        raise ValueError(f"region 格式必須是 x1,y1,x2,y2，收到：{s}")
    return tuple(parts)  # type: ignore


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Artale 提醒道具掃價 Bot")
    parser.add_argument("--equip-region", default=None, metavar="x1,y1,x2,y2",
                        help="覆蓋裝備價格擷取區域，例如 1427,412,1713,501")
    parser.add_argument("--default-region", default=None, metavar="x1,y1,x2,y2",
                        help="覆蓋卷軸/技能書價格擷取區域，例如 1722,415,1975,503")
    args = parser.parse_args()

    equip_region = _parse_region(args.equip_region) if args.equip_region else None
    default_region = _parse_region(args.default_region) if args.default_region else None

    AlertBot(equip_region=equip_region, default_region=default_region).start()
