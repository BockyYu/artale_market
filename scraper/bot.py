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

from api_client import fetch_latest_prices_batch
from config import API_BASE_URL, BETWEEN_ITEMS_DELAY, SCHEDULE_INTERVAL_MINUTES
from notify import send_message, build_alert_message
from scraper import get_game_window, scrape_item, verify_price_header, is_in_auction_screen, enter_auction, reenter_auction, exit_auction, preload_ocr

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
    def __init__(self, equip_region=None, default_region=None, auction_btn=None) -> None:
        self._win = None
        self._equip_region = equip_region
        self._default_region = default_region
        self._auction_btn = auction_btn

    # ------------------------------------------------------------------
    # 私有：API 操作
    # ------------------------------------------------------------------

    def _fetch_alert_items(self) -> list[dict]:
        """取得所有啟用中提醒的道具清單（不需登入），失敗自動重試一次。"""
        for attempt in range(2):
            try:
                r = requests.get(
                    f"{API_BASE_URL}/api/bot/alert-items",
                    timeout=10,
                )
                r.raise_for_status()
                resp = r.json()
                items = resp.get("data", resp) if isinstance(resp, dict) else resp
                return items if isinstance(items, list) else []
            except requests.RequestException as e:
                logger.warning(f"[Bot] 取得提醒道具失敗（第 {attempt + 1} 次）：{e}")
                if attempt == 0:
                    time.sleep(3)
        return []

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

        # 初始化遊戲視窗並進入拍賣畫面（每輪重抓，確保視窗位置最新）
        try:
            self._win = get_game_window()
            if not enter_auction(self._win, btn_pos=self._auction_btn):
                logger.error("[Bot] 無法進入拍賣畫面，請確認按鈕座標設定正確")
                return
            verify_price_header(self._win)
        except RuntimeError as e:
            logger.error(f"[Bot] 視窗錯誤：{e}")
            send_message(f"🚨 Bot 警告：找不到遊戲視窗，遊戲可能已崩潰\n{e}")
            self._win = None
            return

        ok, fail = 0, []
        consecutive_fails = 0 # 連續失敗
        total = len(items)

        # 一次批次取回所有商品的最新 DB 價格，loop 內直接查 dict
        all_ids = [i["item_id"] for i in items if i.get("item_id") is not None]
        try:
            latest_price_map = fetch_latest_prices_batch(all_ids)
        except Exception as e:
            logger.warning(f"[Bot] 批次取得最新價格失敗，本輪不做價格比對：{e}")
            latest_price_map = {}

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
                    consecutive_fails += 1
                    logger.warning(f"▶ [{idx}/{total}] {name} → 找不到價格，跳過")
                    fail.append(name)
                    if consecutive_fails >= 2:
                        logger.warning("[Bot] 連續 2 筆找不到價格，檢查是否已離開拍賣畫面...")
                        if not is_in_auction_screen(self._win):
                            logger.warning("[Bot] 已離開拍賣畫面，準備重新進入")
                            if reenter_auction(self._win, btn_pos=self._auction_btn):
                                verify_price_header(self._win)
                                consecutive_fails = 0
                            else:
                                logger.error("[Bot] 無法重新進入拍賣畫面，終止本次掃描")
                                break
                else:
                    consecutive_fails = 0
                    latest_price = latest_price_map.get(item_id)

                    verified = self._verify_price(
                        price, latest_price, name, search_name, item_type, idx, total
                    )
                    if verified is None:
                        fail.append(name)
                    else:
                        price = verified

                    if verified is not None:
                        if latest_price is not None and latest_price == price:
                            logger.info(f"▶ [{idx}/{total}] {name} → {price:,} → 與 DB 最新價格相同，略過寫入")
                            ok += 1
                        elif self._record_price(item_id, price):
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

        else:
            # for loop 正常結束（未被 break）才離開拍賣畫面
            exit_auction(self._win)
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

    def _verify_price(
        self,
        price: int,
        latest_price: int | None,
        name: str,
        search_name: str,
        item_type: int,
        idx: int,
        total: int,
    ) -> int | None:
        """
        比較當前查到的金額與 DB 最新金額，決定最終寫入價格。

        漲/跌幅 ≤ 50%         → 直接回傳原始價格
        漲/跌幅 > 50%，第二次失敗  → 回傳 None（略過本次更新）
        漲/跌幅 > 50%，兩次一致    → 確認正確，回傳該價格
        漲/跌幅 > 50%，兩次不一致  → 採用第二次結果回傳
        """
        if not latest_price or latest_price <= 0:
            return price

        change_pct = abs(price - latest_price) / latest_price
        if change_pct <= 0.5:
            return price

        direction = "漲" if price > latest_price else "跌"
        logger.warning(
            f"▶ [{idx}/{total}] {name} → 價格{direction}幅 {change_pct*100:.0f}%"
            f"（{latest_price:,} → {price:,}），重新查詢確認..."
        )

        verify_price = scrape_item(
            self._win, search_name, item_type,
            equip_region=self._equip_region,
            default_region=self._default_region,
        )

        if verify_price is None:
            logger.warning(f"▶ [{idx}/{total}] {name} → 重新查詢失敗，略過本次更新")
            return None

        if verify_price != price:
            logger.warning(
                f"▶ [{idx}/{total}] {name} → 兩次查詢結果不一致"
                f"（第一次 {price:,} / 第二次 {verify_price:,}），採用第二次結果"
            )
            return verify_price

        logger.info(f"▶ [{idx}/{total}] {name} → 重新確認一致（{price:,}），繼續寫入")
        return price

    def _ensure_positions(self) -> None:
        """
        檢查所有必要的螢幕座標是否已設定。
        若有缺漏，引導使用者互動校準並寫入 .env，同時更新 in-memory 設定。
        """
        import config as _cfg
        import scraper as _scr

        def _sync(key: str, val: tuple[int, int]) -> None:
            setattr(_cfg, key, val)
            if hasattr(_scr, key):
                setattr(_scr, key, val)

        # ── 1. 拍賣按鈕（進場前設定，不需先開拍賣）────────────────────
        if _cfg.AUCTION_BTN_POS == (0, 0):
            print("\n" + "=" * 55)
            print("⚙️  首次設定：拍賣按鈕位置（AUCTION_BTN_POS）")
            print("=" * 55)
            val = _calibrate_and_save(
                "請在遊戲中回到拍賣場外",
                "請將滑鼠移到拍賣按鈕上",
                "AUCTION_BTN_POS",
                need_game_open=False,
            )
            _sync("AUCTION_BTN_POS", val)
            self._auction_btn = val

        # ── 2. 拍賣場內位置（需要先進拍賣）──────────────────────────────
        in_auction_items = [
            ("SEARCH_BOX_POS",  "將滑鼠移到搜尋輸入框上"),
            ("PRICE_SORT_POS",  "搜尋任一道具讓列表出現，再將滑鼠移到「每個價錢」欄位標題上"),
            ("AUCTION_EXIT_POS","將滑鼠移到離開拍賣的按鈕上"),
        ]
        missing = [(k, p) for k, p in in_auction_items if getattr(_cfg, k) == (0, 0)]

        if missing:
            print("\n" + "=" * 55)
            print("⚙️  首次設定：拍賣場內位置")
            print("需要校準以下座標：")
            for k, _ in missing:
                print(f"   • {k}")
            print("=" * 55)
            print("請先【手動開啟拍賣畫面】，確認搜尋框可見後按 Enter 繼續...")
            input()

            for env_key, prompt_move in missing:
                print(f"\n── 設定 {env_key} ──")
                val = _calibrate_and_save(
                    "",          # prompt_before 不再使用（拍賣已開啟）
                    prompt_move,
                    env_key,
                    need_game_open=False,
                )
                _sync(env_key, val)

        logger.info("[Bot] 所有座標已就緒")

    def start(self) -> None:
        """立即執行一次，之後每 SCHEDULE_INTERVAL_MINUTES 分鐘自動觸發。"""
        logger.info(f"[Bot] 啟動，間隔：{SCHEDULE_INTERVAL_MINUTES} 分鐘")
        logger.info("[Bot] 按 Ctrl+C 停止")
        self._ensure_positions()
        logger.info("[Bot] 預載 OCR 引擎...")
        preload_ocr()
        logger.info("[Bot] OCR 就緒，開始執行")

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


def _update_env(key: str, value: str) -> None:
    from pathlib import Path
    env_path = Path(__file__).parent / ".env"
    lines = env_path.read_text(encoding="utf-8").splitlines() if env_path.exists() else []
    for i, line in enumerate(lines):
        if line.startswith(f"{key}="):
            lines[i] = f"{key}={value}"
            break
    else:
        lines.append(f"{key}={value}")
    env_path.write_text("\n".join(lines) + "\n", encoding="utf-8")


def _calibrate_and_save(prompt_before: str, prompt_move: str, env_key: str, need_game_open: bool = True) -> tuple[int, int]:
    import threading, pyautogui as _pag, time as _time
    from scraper import get_game_window
    if need_game_open:
        print(f"{prompt_before}，確認後按任意鍵繼續...")
        input()
    win = get_game_window()
    print(f"\n{prompt_move}，準備好後按 Enter 確認\n")
    _stop = threading.Event()
    def _show():
        while not _stop.is_set():
            sx, sy = _pag.position()
            print(f"\r  視窗相對座標：({sx - win.left:5d}, {sy - win.top:5d})   ", end="", flush=True)
            _time.sleep(0.05)
    t = threading.Thread(target=_show, daemon=True)
    t.start()
    input()
    _stop.set()
    sx, sy = _pag.position()
    bx, by = sx - win.left, sy - win.top
    print(f"\n已記錄座標：({bx}, {by})")
    _update_env(env_key, f"{bx},{by}")
    logger.info(f"已將 {env_key}={bx},{by} 寫入 .env")
    return (bx, by)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Artale 提醒道具掃價 Bot")
    parser.add_argument("--equip-region", default=None, metavar="x1,y1,x2,y2",
                        help="覆蓋裝備價格擷取區域，例如 1427,412,1713,501")
    parser.add_argument("--default-region", default=None, metavar="x1,y1,x2,y2",
                        help="覆蓋卷軸/技能書價格擷取區域，例如 1722,415,1975,503")
    parser.add_argument("--auction-btn", default=None, metavar="x,y",
                        help="拍賣按鈕的視窗相對座標（覆蓋 .env 的 AUCTION_BTN_POS）")
    parser.add_argument("--set-auction-btn", action="store_true",
                        help="互動式校準並儲存拍賣按鈕座標到 .env")
    parser.add_argument("--set-search-box", action="store_true",
                        help="互動式校準並儲存拍賣搜尋輸入框座標到 .env")
    parser.add_argument("--set-price-sort", action="store_true",
                        help="互動式校準並儲存「每個價錢」排序欄位座標到 .env")
    parser.add_argument("--set-price-row", action="store_true",
                        help="互動式校準並儲存第一筆價格列點擊座標到 .env")
    parser.add_argument("--set-auction-exit", action="store_true",
                        help="互動式校準並儲存離開拍賣按鈕座標到 .env")
    args = parser.parse_args()

    equip_region = _parse_region(args.equip_region) if args.equip_region else None
    default_region = _parse_region(args.default_region) if args.default_region else None
    auction_btn = tuple(int(v) for v in args.auction_btn.split(",")) if args.auction_btn else None

    if args.set_auction_btn:
        _calibrate_and_save("請在遊戲中回到拍賣外面", "請將滑鼠移到拍賣按鈕上", "AUCTION_BTN_POS", need_game_open=False)
        sys.exit(0)

    if args.set_search_box:
        _calibrate_and_save("請先手動開啟拍賣畫面，確認搜尋輸入框可見", "請將滑鼠移到搜尋輸入框上", "SEARCH_BOX_POS")
        sys.exit(0)

    if args.set_price_sort:
        _calibrate_and_save("請先手動開啟拍賣畫面，確認「每個價錢」欄位標題可見", "請將滑鼠移到「每個價錢」欄位標題上", "PRICE_SORT_POS")
        sys.exit(0)

    if args.set_price_row:
        _calibrate_and_save("請先搜尋任一道具讓列表出現", "請將滑鼠移到第一筆價格列上", "PRICE_ROW_POS")
        sys.exit(0)

    if args.set_auction_exit:
        _calibrate_and_save("請先開啟拍賣畫面，讓離開按鈕可見", "請將滑鼠移到離開拍賣的按鈕上", "AUCTION_EXIT_POS")
        sys.exit(0)

    AlertBot(equip_region=equip_region, default_region=default_region, auction_btn=auction_btn).start()
