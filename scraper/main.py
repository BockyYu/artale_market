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
from datetime import datetime

import schedule

from api_client import fetch_items, fetch_unrecorded_items, fetch_alert_map, fetch_latest_price, record_price
from config import BETWEEN_ITEMS_DELAY, SCHEDULE_TIME, SCHEDULE_INTERVAL_MINUTES
from notify import send_message, build_alert_message, build_price_surge_message, build_price_changes_summary
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


def _parse_region(s: str) -> tuple[int, int, int, int]:
    parts = [int(v.strip()) for v in s.split(",")]
    if len(parts) != 4:
        raise ValueError(f"region 格式必須是 x1,y1,x2,y2，收到：{s}")
    return tuple(parts)  # type: ignore


def run(dry_run: bool = False, fill_missing: bool = False, equip_region=None, default_region=None) -> None:
    mode_label = "（補漏模式）" if fill_missing else ""
    logger.info("=" * 50)
    logger.info(f"開始抓取{'（模擬模式）' if dry_run else ''}{mode_label}")
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
        if fill_missing:
            items = fetch_unrecorded_items()
        else:
            items = fetch_items()
    except Exception as e:
        logger.error(f"無法取得商品列表: {e}")
        return

    if not items:
        if fill_missing:
            logger.info("今天所有追蹤商品都已有價格記錄，無需補漏")
        else:
            logger.warning("商品列表為空，結束")
        return

    logger.info(f"共 {len(items)} 個商品待抓取")

    alert_map = fetch_alert_map()

    ok, fail = 0, []
    changed_items = []
    total = len(items)

    for idx, item in enumerate(items, 1):
        name = item.get("item_name", "")
        item_id = item.get("item_id")

        if not name or item_id is None:
            continue

        item_type = item.get("item_type", 1)
        search_mode = item.get("search_mode", 1)
        english_name = item.get("english_name", "")
        search_name = english_name if search_mode == 2 and english_name else name
        if search_name != name:
            logger.info(f"  使用英文名稱查詢：{search_name}")
        try:
            price = scrape_item(win, search_name, item_type, equip_region=equip_region, default_region=default_region)

            if price is None:
                logger.warning(f"▶ [{idx}/{total}] {name} → 找不到價格，跳過")
                fail.append(name)
            else:
                if dry_run:
                    logger.info(f"▶ [{idx}/{total}] {name} → {price:,} → 模擬模式，不寫入")
                    ok += 1
                else:
                    prev_price = fetch_latest_price(item_id)
                    if record_price(item_id, price):
                        logger.info(f"▶ [{idx}/{total}] {name} → {price:,} → 已寫入 DB")
                        ok += 1
                    else:
                        logger.warning(f"▶ [{idx}/{total}] {name} → {price:,} → 寫入 DB 失敗")
                        fail.append(name)

                    if prev_price is not None and prev_price > 0:
                        change_pct = (price - prev_price) / prev_price * 100
                        if price != prev_price:
                            changed_items.append({
                                "name": name,
                                "prev_price": prev_price,
                                "new_price": price,
                                "change_pct": round(change_pct, 1),
                            })
                        if abs(change_pct) >= 50:
                            now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                            msg = build_price_surge_message(name, prev_price, price, change_pct, now)
                            if send_message(msg):
                                logger.info(f"  已發送價格異動通知：{name} {change_pct:+.1f}%")
                            else:
                                logger.warning(f"  價格異動幅度達 {change_pct:+.1f}% 但通知發送失敗：{name}")

                alert = alert_map.get(item_id)
                if alert:
                    threshold = alert["threshold_price"]
                    bot_id = alert.get("bot_id")
                    if threshold > 0 and price <= threshold:
                        now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                        msg = build_alert_message(name, price, threshold, now)
                        if send_message(msg, bot_id=bot_id):
                            logger.info(f"  已發送價格通知：{name} → {price:,}")
                        else:
                            logger.warning(f"  價格符合門檻但通知發送失敗：{name} → {price:,}")

        except Exception as e:
            logger.error(f"▶ [{idx}/{total}] {name} → 錯誤：{e}")
            fail.append(name)

        time.sleep(BETWEEN_ITEMS_DELAY)

    logger.info("=" * 50)
    logger.info(f"完成：{ok}/{len(items)} 筆成功")
    if fail:
        logger.warning(f"失敗項目（{len(fail)} 筆）：{', '.join(fail)}")
    logger.info("=" * 50)

    if not dry_run:
        msg = f"✅ 抓取完成：{ok}/{len(items)} 筆成功"
        if fail:
            msg += f"\n❌ 失敗（{len(fail)} 筆）：{', '.join(fail)}"
        send_message(msg)

        if changed_items:
            now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
            summary = build_price_changes_summary(changed_items, now)
            if send_message(summary):
                logger.info(f"  已發送價格異動摘要（{len(changed_items)} 筆）")
            else:
                logger.warning("  價格異動摘要發送失敗")


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


def track_mouse() -> None:
    """即時顯示滑鼠相對遊戲視窗的座標，按 Ctrl+C 停止。"""
    import pyautogui
    from scraper import get_game_window

    win = get_game_window()
    print(f"視窗位置：left={win.left}, top={win.top}, width={win.width}, height={win.height}")
    print("移動滑鼠到任意位置，即時顯示視窗相對座標（按 Ctrl+C 停止）\n")
    try:
        while True:
            sx, sy = pyautogui.position()
            wx = sx - win.left
            wy = sy - win.top
            print(f"\r  螢幕({sx:5d},{sy:5d}) | 視窗相對({wx:5d},{wy:5d})   ", end="", flush=True)
            time.sleep(0.05)
    except KeyboardInterrupt:
        print("\n停止追蹤")


def debug_pos(win_x: int, win_y: int) -> None:
    """移動滑鼠到指定視窗座標，截圖標記，確認點擊位置是否正確。"""
    import pyautogui
    from scraper import get_game_window, _capture_full
    from PIL import Image, ImageDraw

    win = get_game_window()
    screen_x = win.left + win_x
    screen_y = win.top + win_y
    logger.info(f"視窗座標 ({win_x}, {win_y}) → 螢幕座標 ({screen_x}, {screen_y})")
    logger.info("移動滑鼠中（3 秒後截圖）...")
    pyautogui.moveTo(screen_x, screen_y)
    time.sleep(3)

    img_arr, _, _ = _capture_full(win)
    img = Image.fromarray(img_arr)
    draw = ImageDraw.Draw(img)
    r = 12
    draw.ellipse([win_x - r, win_y - r, win_x + r, win_y + r], outline="red", width=3)
    draw.line([win_x - r - 5, win_y, win_x + r + 5, win_y], fill="red", width=2)
    draw.line([win_x, win_y - r - 5, win_x, win_y + r + 5], fill="red", width=2)
    save_path = f"debug_pos_{win_x}_{win_y}.png"
    img.save(save_path)
    logger.info(f"截圖已儲存：{save_path}（紅圈標記目標位置）")


def main() -> None:
    parser = argparse.ArgumentParser(description="Artale 拍賣場自動抓價腳本")
    parser.add_argument("--now", action="store_true", help="立即執行一次，不進入排程")
    parser.add_argument("--interval", action="store_true",
                        help=f"間隔模式：每 {SCHEDULE_INTERVAL_MINUTES} 分鐘自動執行一次")
    parser.add_argument("--time", default=SCHEDULE_TIME, metavar="HH:MM",
                        help=f"每日執行時間（預設 {SCHEDULE_TIME}）")
    parser.add_argument("--dry-run", action="store_true",
                        help="模擬模式：只印出價格，不寫入 DB")
    parser.add_argument("--fill-missing", action="store_true",
                        help="補漏模式：只抓今天還沒有價格記錄的商品")
    parser.add_argument("--test-notify", action="store_true",
                        help="發送一則測試訊息，確認通知功能是否正常")
    parser.add_argument("--debug-ocr", action="store_true",
                        help="截圖並印出所有 OCR 結果，用於診斷辨識問題")
    parser.add_argument("--debug-pos", nargs=2, type=int, metavar=("X", "Y"),
                        help="移動滑鼠到指定視窗座標並截圖標記，例如 --debug-pos 1848 391")
    parser.add_argument("--track-mouse", action="store_true",
                        help="即時顯示滑鼠相對遊戲視窗座標，用於定位 UI 元素")
    parser.add_argument("--equip-region", default=None, metavar="x1,y1,x2,y2",
                        help="覆蓋裝備價格擷取區域（視窗相對座標），例如 1427,412,1713,501")
    parser.add_argument("--default-region", default=None, metavar="x1,y1,x2,y2",
                        help="覆蓋卷軸/技能書價格擷取區域，例如 1722,415,1975,503")
    args = parser.parse_args()

    dry_run: bool = args.dry_run
    fill_missing: bool = args.fill_missing
    equip_region = _parse_region(args.equip_region) if args.equip_region else None
    default_region = _parse_region(args.default_region) if args.default_region else None

    if args.test_notify:
        from notify import _get_active_bot_ids
        bot_ids = _get_active_bot_ids()
        if not bot_ids:
            logger.error("後端無啟用中的通知 bot，請至後台新增並啟用")
            return
        logger.info(f"找到 {len(bot_ids)} 個啟用中的 bot：{bot_ids}")
        if send_message("✅ Artale Market 通知測試訊息"):
            logger.info("測試訊息發送成功")
        else:
            logger.error("測試訊息發送失敗")
        return

    if args.debug_ocr:
        debug_ocr()
        return

    if args.debug_pos:
        debug_pos(*args.debug_pos)
        return

    if args.track_mouse:
        track_mouse()
        return

    if args.now:
        run(dry_run=dry_run, fill_missing=fill_missing, equip_region=equip_region, default_region=default_region)
        return

    if args.interval:
        logger.info(f"間隔模式：每 {SCHEDULE_INTERVAL_MINUTES} 分鐘自動執行")
        logger.info("按 Ctrl+C 停止")
        run(dry_run=dry_run, fill_missing=fill_missing, equip_region=equip_region, default_region=default_region)
        schedule.every(SCHEDULE_INTERVAL_MINUTES).minutes.do(
            run, dry_run=dry_run, fill_missing=fill_missing, equip_region=equip_region, default_region=default_region
        )
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

    schedule.every().day.at(run_time).do(
        run, dry_run=dry_run, fill_missing=fill_missing, equip_region=equip_region, default_region=default_region
    )

    try:
        while True:
            schedule.run_pending()
            time.sleep(30)
    except KeyboardInterrupt:
        logger.info("已停止排程")


if __name__ == "__main__":
    main()
