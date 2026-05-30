"""
使用方式：
  立即執行一次：
    python main.py --now

  排程模式（預設每天 08:00，或用 --time 指定）：
    python main.py
    python main.py --time 14:30

  模擬模式（不實際寫入 DB，只印出抓到的價格）：
    python main.py --now --dry-run
"""

import argparse
import logging
import sys
import time

import schedule

from api_client import fetch_items, record_price
from config import BETWEEN_ITEMS_DELAY, SCHEDULE_TIME, SCHEDULE_INTERVAL_MINUTES
from scraper import get_game_window, scrape_item, verify_price_header

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s  %(levelname)s  %(message)s",
    datefmt="%H:%M:%S",
    handlers=[
        logging.StreamHandler(sys.stdout),
        logging.FileHandler("scraper.log", encoding="utf-8"),
    ],
)
logger = logging.getLogger(__name__)


def run(dry_run: bool = False) -> None:
    logger.info("=" * 50)
    logger.info(f"開始抓取{'（模擬模式）' if dry_run else ''}")
    logger.info("=" * 50)

    try:
        win = get_game_window()
    except RuntimeError as e:
        logger.error(str(e))
        return

    try:
        verify_price_header(win)
    except RuntimeError as e:
        logger.error(f"前置檢查失敗：{e}")
        return

    try:
        items = fetch_items()
    except Exception as e:
        logger.error(f"無法取得商品列表: {e}")
        return

    if not items:
        logger.warning("商品列表為空，結束")
        return

    logger.info(f"共 {len(items)} 個商品待抓取")

    ok, fail = 0, []
    total = len(items)

    for idx, item in enumerate(items, 1):
        name = item.get("item_name", "")
        item_id = item.get("item_id")

        if not name or item_id is None:
            continue

        item_type = item.get("item_type", 1)
        try:
            price = scrape_item(win, name, item_type)

            if price is None:
                logger.warning(f"▶ [{idx}/{total}] {name} → 找不到價格，跳過")
                fail.append(name)
            else:
                if dry_run:
                    logger.info(f"▶ [{idx}/{total}] {name} → {price:,} → 模擬模式，不寫入")
                    ok += 1
                else:
                    if record_price(item_id, price):
                        logger.info(f"▶ [{idx}/{total}] {name} → {price:,} → 已寫入 DB")
                        ok += 1
                    else:
                        logger.warning(f"▶ [{idx}/{total}] {name} → {price:,} → 寫入 DB 失敗")
                        fail.append(name)

        except Exception as e:
            logger.error(f"▶ [{idx}/{total}] {name} → 錯誤：{e}")
            fail.append(name)

        time.sleep(BETWEEN_ITEMS_DELAY)

    logger.info("=" * 50)
    logger.info(f"完成：{ok}/{len(items)} 筆成功")
    if fail:
        logger.warning(f"失敗項目（{len(fail)} 筆）：{', '.join(fail)}")
    logger.info("=" * 50)


def debug_ocr() -> None:
    """截圖並印出所有 OCR 辨識結果，用於診斷找不到欄位的問題。"""
    from scraper import get_game_window, _capture_full, _get_reader
    from PIL import Image

    win = get_game_window()
    logger.info("截圖中...")
    img, wx, wy = _capture_full(win)
    save_path = "debug_screenshot.png"
    Image.fromarray(img).save(save_path)
    logger.info(f"截圖已儲存：{save_path}")

    logger.info("執行 OCR...")
    results = _get_reader().readtext(img)
    logger.info(f"共辨識到 {len(results)} 個文字區塊：")
    for i, (bbox, text, conf) in enumerate(results):
        cx = int((bbox[0][0] + bbox[2][0]) / 2)
        cy = int((bbox[0][1] + bbox[2][1]) / 2)
        logger.info(f"  [{i:02d}] ({cx:4d},{cy:4d})  conf={conf:.2f}  text={repr(text)}")


def main() -> None:
    parser = argparse.ArgumentParser(description="Artale 拍賣場自動抓價腳本")
    parser.add_argument("--now", action="store_true", help="立即執行一次，不進入排程")
    parser.add_argument("--interval", action="store_true",
                        help=f"間隔模式：每 {SCHEDULE_INTERVAL_MINUTES} 分鐘自動執行一次")
    parser.add_argument("--time", default=SCHEDULE_TIME, metavar="HH:MM",
                        help=f"每日執行時間（預設 {SCHEDULE_TIME}）")
    parser.add_argument("--dry-run", action="store_true",
                        help="模擬模式：只印出價格，不寫入 DB")
    parser.add_argument("--debug-ocr", action="store_true",
                        help="截圖並印出所有 OCR 結果，用於診斷辨識問題")
    args = parser.parse_args()

    dry_run: bool = args.dry_run

    if args.debug_ocr:
        debug_ocr()
        return

    if args.now:
        run(dry_run=dry_run)
        return

    if args.interval:
        logger.info(f"間隔模式：每 {SCHEDULE_INTERVAL_MINUTES} 分鐘自動執行")
        logger.info("按 Ctrl+C 停止")
        run(dry_run=dry_run)
        schedule.every(SCHEDULE_INTERVAL_MINUTES).minutes.do(run, dry_run=dry_run)
        try:
            while True:
                schedule.run_pending()
                time.sleep(30)
        except KeyboardInterrupt:
            logger.info("已停止排程")
        return

    run_time: str = args.time
    logger.info(f"排程模式：每天 {run_time} 自動執行")
    logger.info("按 Ctrl+C 停止")

    schedule.every().day.at(run_time).do(run, dry_run=dry_run)

    try:
        while True:
            schedule.run_pending()
            time.sleep(30)
    except KeyboardInterrupt:
        logger.info("已停止排程")


if __name__ == "__main__":
    main()
