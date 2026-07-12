"""
遊戲視窗自動化模組。
依賴：easyocr、pyautogui、pygetwindow、mss、Pillow、pyperclip、numpy

座標系統：所有 _FIXED_* 常數都是相對於遊戲視窗左上角（win.left, win.top），
由 debug-ocr 校準。
"""

import logging
import platform
import re
import subprocess
import time

import easyocr
import mss
import numpy as np
import pyautogui
import pyperclip
from PIL import Image

_IS_MAC = platform.system() == "Darwin"
_MOD_KEY = "command" if _IS_MAC else "ctrl"

if not _IS_MAC:
    import pygetwindow as gw


class _MacWindow:
    """macOS 視窗包裝，介面與 pygetwindow 相容。"""

    def __init__(self, app_name: str, bounds: dict):
        self._app_name = app_name
        self.left = bounds["left"]
        self.top = bounds["top"]
        self.width = bounds["width"]
        self.height = bounds["height"]
        self.isMinimized = bounds.get("minimized", False)

    def restore(self):
        subprocess.run(
            ["osascript", "-e",
             f'tell application "{self._app_name}" to set miniaturized of window 1 to false'],
            check=False,
        )

    def activate(self):
        subprocess.run(
            ["osascript", "-e", f'tell application "{self._app_name}" to activate'],
            check=False,
        )


def _mac_get_windows_with_title(title: str) -> list[_MacWindow]:
    script = f"""
set matchTitle to "{title}"
set results to {{}}
tell application "System Events"
    repeat with proc in (every process whose background only is false)
        try
            repeat with w in (every window of proc)
                if name of w contains matchTitle then
                    set b to position of w
                    set s to size of w
                    set end of results to {{appName:name of proc, x:item 1 of b, y:item 2 of b, w:item 1 of s, h:item 2 of s}}
                end if
            end repeat
        end try
    end repeat
end tell
return results
"""
    try:
        out = subprocess.check_output(["osascript", "-e", script], stderr=subprocess.DEVNULL, timeout=10).decode().strip()
    except subprocess.CalledProcessError:
        return []
    if not out:
        return []
    # osascript 回傳格式：appName:MapleStory ..., x:100, y:50, w:1280, h:720, ...
    wins = []
    for entry in out.split(", appName:"):
        entry = entry.strip().lstrip("appName:")
        try:
            parts = dict(p.strip().split(":", 1) for p in ("appName:" + entry).split(", "))
            wins.append(_MacWindow(parts["appName"], {
                "left": int(parts["x"]),
                "top": int(parts["y"]),
                "width": int(parts["w"]),
                "height": int(parts["h"]),
            }))
        except Exception:
            continue
    return wins

from config import (
    AFTER_SEARCH_DELAY,
    AFTER_SORT_DELAY,
    AUCTION_BTN_POS,
    AUCTION_EXIT_POS,
    AUCTION_WAIT_SECS,
    PRICE_REGION_DEFAULT,
    PRICE_REGION_EQUIP,
    PRICE_ROW_POS,
    PRICE_SORT_POS,
    SEARCH_BOX_POS,
    WINDOW_TITLE,
)

logger = logging.getLogger(__name__)

# OCR 引擎懶初始化
_reader: easyocr.Reader | None = None

# Retina / HiDPI 縮放比例快取（Mac 通常為 2.0，Windows 為 1.0）
_scale_factor: float | None = None


def _get_scale_factor(win) -> float:
    """偵測螢幕實體像素與邏輯座標的縮放比例，結果快取。"""
    global _scale_factor
    if _scale_factor is not None:
        return _scale_factor
    with mss.MSS() as sct:
        shot = sct.grab({"left": win.left, "top": win.top, "width": 100, "height": 100})
    _scale_factor = shot.width / 100
    return _scale_factor

# 商品類型常數（對應後端 model.ItemType）
_ITEM_TYPE_EQUIP = 6

# 搜尋輸入框快取：(center_x, center_y)，相對於視窗左上角
_search_box_cache: tuple[int, int] | None = None

# 「每個價錢」欄位快取，以 item_type 為 key
# key=-1 是 verify_price_header 啟動時偵測的通用快取，供各 item_type 初次使用
_price_header_cache: dict[int, tuple[int, int, int, int, int]] = {}
_UNIVERSAL_CACHE_KEY = -1


def _get_reader() -> easyocr.Reader:
    global _reader
    if _reader is None:
        logger.info("初始化 OCR 引擎（繁體中文 + 英文）...")
        _reader = easyocr.Reader(["ch_tra", "en"], gpu=False)
        logger.info("OCR 引擎就緒")
    return _reader


# ---------------------------------------------------------------------------
# 視窗操作
# ---------------------------------------------------------------------------

def get_game_window():
    """找到遊戲視窗並將其帶到前景。"""
    if _IS_MAC:
        wins = _mac_get_windows_with_title(WINDOW_TITLE)
    else:
        wins = gw.getWindowsWithTitle(WINDOW_TITLE)
    if not wins:
        raise RuntimeError(f"找不到遊戲視窗，請確認已開啟：{WINDOW_TITLE}")
    win = wins[0]
    if win.isMinimized:
        win.restore()
        time.sleep(0.3)
    try:
        win.activate()
    except Exception:
        pass
    time.sleep(0.4)
    return win


def _capture_column(win, x_min: int, x_max: int, y_min: int) -> np.ndarray:
    """只截取「每個價錢」欄位的資料區域，供 OCR 使用。"""
    with mss.MSS() as sct:
        region = {
            "left": max(0, win.left + x_min),
            "top": max(0, win.top + y_min),
            "width": x_max - x_min,
            "height": win.height - y_min,
        }
        shot = sct.grab(region)
    return np.array(Image.frombytes("RGB", shot.size, shot.rgb))


def _capture_full(win) -> tuple[np.ndarray, int, int]:
    """截取完整視窗，回傳 (RGB ndarray, win.left, win.top)。供 debug 使用。"""
    with mss.MSS() as sct:
        region = {
            "left": max(0, win.left),
            "top": max(0, win.top),
            "width": win.width,
            "height": win.height,
        }
        shot = sct.grab(region)
    img = np.array(Image.frombytes("RGB", shot.size, shot.rgb))
    return img, win.left, win.top


# ---------------------------------------------------------------------------
# 自動化步驟
# ---------------------------------------------------------------------------

def _find_search_box(win) -> tuple[int, int]:
    """
    回傳搜尋輸入框的視窗相對座標 (cx, cy)。
    優先使用 .env 的 SEARCH_BOX_POS，否則 OCR 偵測並快取。
    """
    global _search_box_cache
    if _search_box_cache is not None:
        return _search_box_cache

    # 優先使用 .env 固定座標（跳過 OCR）
    sx, sy = SEARCH_BOX_POS
    if sx != 0 or sy != 0:
        _search_box_cache = (sx, sy)
        logger.info(f"  使用 .env 搜尋框座標：({sx},{sy})")
        return _search_box_cache

    logger.info("  自動偵測搜尋輸入框位置...")
    img, _, _ = _capture_full(win)
    results = _get_reader().readtext(img)
    scale = _get_scale_factor(win)

    for bbox, text, _ in results:
        if any(k in text.replace(" ", "") for k in ["請輸入", "道具名稱", "輸入道具"]):
            cx = int((bbox[0][0] + bbox[2][0]) / 2 / scale)
            cy = int((bbox[0][1] + bbox[2][1]) / 2 / scale)
            _search_box_cache = (cx, cy)
            logger.info(f"  搜尋框偵測成功：center=({cx},{cy})")
            return _search_box_cache

    raise RuntimeError("無法偵測搜尋輸入框，請確認遊戲畫面正確顯示拍賣介面")


def search_item(win, item_name: str) -> None:
    """在搜尋欄輸入商品名稱並送出。"""
    cx, cy = _find_search_box(win)
    try:
        win.activate()
    except Exception:
        pass
    time.sleep(0.3)
    pyautogui.click(win.left + cx, win.top + cy)
    time.sleep(0.5)
    pyautogui.hotkey(_MOD_KEY, "a")  # 選取全部
    time.sleep(0.1)
    pyautogui.press("delete")        # 刪除選取內容，確保欄位清空
    time.sleep(0.1)
    pyperclip.copy(item_name)
    time.sleep(0.3)                  # 確保 pbcopy 完成
    pyautogui.hotkey(_MOD_KEY, "v")
    time.sleep(0.5)
    pyautogui.press("enter")
    logger.info(f"  搜尋送出: {item_name}")
    time.sleep(AFTER_SEARCH_DELAY)


def _find_price_header(win, item_type: int = 0) -> tuple[int, int, int, int, int]:
    """
    OCR 自動偵測「每個價錢」欄位標題位置。
    回傳 (center_x, center_y, col_x_min, col_x_max, row_y_min)，座標相對於視窗左上角。
    結果以 item_type 為 key 快取，同類型道具不需重複偵測。
    """
    global _price_header_cache
    if item_type in _price_header_cache:
        return _price_header_cache[item_type]
    # 沿用 verify_price_header 啟動時的通用偵測，避免重複 OCR（~12s）
    if _UNIVERSAL_CACHE_KEY in _price_header_cache:
        _price_header_cache[item_type] = _price_header_cache[_UNIVERSAL_CACHE_KEY]
        logger.info(f"  使用啟動快取欄位座標 (item_type={item_type})")
        return _price_header_cache[item_type]

    # 優先使用 .env 固定座標（跳過 OCR）
    px, py = PRICE_SORT_POS
    if px != 0 or py != 0:
        _price_header_cache[_UNIVERSAL_CACHE_KEY] = (px, py, max(0, px - 60), px + 60, py + 20)
        _price_header_cache[item_type] = _price_header_cache[_UNIVERSAL_CACHE_KEY]
        logger.info(f"  使用 .env 排序欄位座標：({px},{py})")
        return _price_header_cache[item_type]

    logger.info("  自動偵測「每個價錢」欄位位置...")
    img, _, _ = _capture_full(win)
    results = _get_reader().readtext(img)
    scale = _get_scale_factor(win)

    # 單一文字塊命中
    for bbox, text, _ in results:
        if "每個" in text.replace(" ", ""):
            x1, x2 = int(bbox[0][0] / scale), int(bbox[2][0] / scale)
            cy = int((bbox[0][1] + bbox[2][1]) / 2 / scale)
            cx = (x1 + x2) // 2
            pad = max(60, (x2 - x1) // 2)
            _price_header_cache[item_type] = (cx, cy, max(0, x1 - pad), x2 + pad, cy + 20)
            logger.info(f"  偵測成功（單塊）：center=({cx},{cy})")
            return _price_header_cache[item_type]

    # 合併同列文字塊後再找
    row_map: dict[int, list] = {}
    for bbox, text, _ in results:
        cy = int((bbox[0][1] + bbox[2][1]) / 2)
        cx = int((bbox[0][0] + bbox[2][0]) / 2)
        row_map.setdefault(round(cy / 15) * 15, []).append(
            (cx, cy, text.replace(" ", ""), int(bbox[0][0]), int(bbox[2][0]))
        )
    raise RuntimeError("無法偵測「每個價錢」欄位位置，請確認遊戲畫面正確顯示拍賣列表")


def click_sort_by_lowest(win, item_type: int = 0) -> tuple[int, int, int]:
    """點擊「每個價錢」欄位標題排序，回傳 (col_x_min, col_x_max, row_y_min)。"""
    cx, cy, x_min, x_max, y_min = _find_price_header(win, item_type)
    pyautogui.click(win.left + cx, win.top + cy)
    logger.debug(f"  已點擊每個價錢排序 ({cx}, {cy})")
    time.sleep(AFTER_SORT_DELAY)
    return x_min, x_max, y_min


def _parse_price(text: str) -> int | None:
    """將 OCR 讀到的價格文字轉成整數，支援「萬」單位（例如 3,099萬 → 30990000）。"""
    text = text.strip()
    if "萬" in text or "万" in text:
        logger.info(f"  _parse_price 偵測到萬：原始輸入 {repr(text)}")
        parts = re.split(r"[萬万]", text)
        wan_str = re.sub(r"[^\d]", "", parts[0])
        rem_str = re.sub(r"[^\d]", "", parts[1]) if len(parts) > 1 else ""
        if not wan_str:
            return None
        value = int(wan_str) * 10000
        if rem_str:
            value += int(rem_str)
        return value
    cleaned = re.sub(r"[^\d]", "", text)
    return int(cleaned) if len(cleaned) >= 3 else None


def read_lowest_price(win, x_min: int, x_max: int, y_min: int) -> int | None:
    """
    讀取排序後第一行的「每個價錢」。
    只對欄位資料區做 OCR（小圖），取 y 最小的有效數字。
    """
    img = _capture_column(win, x_min, x_max, y_min)
    results = _get_reader().readtext(img)

    # 把同一列的 OCR 碎片先合併（OCR 有時把「3099萬」切成「3099」+「萬」兩塊）
    row_map: dict[int, list[tuple[int, int, str]]] = {}
    for bbox, text, _conf in results:
        cy = int((bbox[0][1] + bbox[2][1]) / 2)
        cx = int((bbox[0][0] + bbox[2][0]) / 2)
        row_key = round(cy / 12) * 12  # 以 12px 為單位分列
        row_map.setdefault(row_key, []).append((cx, cy, text))

    candidates: list[tuple[int, int]] = []
    for row_key in sorted(row_map.keys()):
        fragments = sorted(row_map[row_key], key=lambda f: f[0])
        cy = int(sum(f[1] for f in fragments) / len(fragments))
        merged = "".join(f[2] for f in fragments)
        # 跳過括號內的合計列（例如「(8萬 3,000)」）
        if "(" in merged or "（" in merged:
            continue
        price = _parse_price(merged)
        if price is not None:
            candidates.append((cy, price))

    if not candidates:
        logger.warning("  無法讀取任何價格數字")
        return None

    candidates.sort(key=lambda t: t[0])
    price = candidates[0][1]
    logger.debug(f"  讀取到最低價: {price:,}")
    return price


def read_price_from_region(win, x1: int, y1: int, x2: int, y2: int) -> int | None:
    """從固定視窗相對區域截圖並 OCR 讀取最低價格。"""
    with mss.MSS() as sct:
        region = {
            "left": win.left + x1,
            "top": win.top + y1,
            "width": x2 - x1,
            "height": y2 - y1,
        }
        shot = sct.grab(region)
    img = np.array(Image.frombytes("RGB", shot.size, shot.rgb))
    results = _get_reader().readtext(img)

    row_map: dict[int, list[tuple[int, int, str]]] = {}
    for bbox, text, _conf in results:
        cy = int((bbox[0][1] + bbox[2][1]) / 2)
        cx = int((bbox[0][0] + bbox[2][0]) / 2)
        row_key = round(cy / 12) * 12
        row_map.setdefault(row_key, []).append((cx, cy, text))

    candidates: list[tuple[int, int]] = []
    for row_key in sorted(row_map.keys()):
        fragments = sorted(row_map[row_key], key=lambda f: f[0])
        cy = int(sum(f[1] for f in fragments) / len(fragments))
        merged = "".join(f[2] for f in fragments)
        if "(" in merged or "（" in merged:
            continue
        price = _parse_price(merged)
        if price is not None:
            candidates.append((cy, price))

    if not candidates:
        logger.warning("  無法讀取任何價格數字")
        return None

    candidates.sort(key=lambda t: t[0])
    price = candidates[0][1]
    logger.debug(f"  讀取到最低價: {price:,}")
    return price


# ---------------------------------------------------------------------------
# 拍賣畫面偵測與自動重進
# ---------------------------------------------------------------------------

def calibrate_auction_btn(win) -> tuple[int, int]:
    """
    互動式校準拍賣按鈕座標。
    移動滑鼠到拍賣按鈕上，按 Enter 確認，回傳視窗相對座標 (x, y)。
    """
    import threading
    print("\n請將滑鼠移到右下角的拍賣按鈕上，準備好後按 Enter 確認")
    print("（即時顯示視窗相對座標，按 Enter 鎖定）\n")

    _stop = threading.Event()

    def _show():
        while not _stop.is_set():
            sx, sy = pyautogui.position()
            wx, wy = sx - win.left, sy - win.top
            print(f"\r  視窗相對座標：({wx:5d}, {wy:5d})   ", end="", flush=True)
            time.sleep(0.05)

    t = threading.Thread(target=_show, daemon=True)
    t.start()
    input()
    _stop.set()

    sx, sy = pyautogui.position()
    bx, by = sx - win.left, sy - win.top
    print(f"\n已記錄拍賣按鈕座標：({bx}, {by})")
    return bx, by


def calibrate_search_box(win) -> tuple[int, int]:
    """
    互動式校準搜尋輸入框座標。
    移動滑鼠到搜尋框上，按 Enter 確認，回傳視窗相對座標 (x, y)。
    """
    import threading
    print("\n請將滑鼠移到拍賣搜尋輸入框上，準備好後按 Enter 確認")
    print("（即時顯示視窗相對座標，按 Enter 鎖定）\n")

    _stop = threading.Event()

    def _show():
        while not _stop.is_set():
            sx, sy = pyautogui.position()
            wx, wy = sx - win.left, sy - win.top
            print(f"\r  視窗相對座標：({wx:5d}, {wy:5d})   ", end="", flush=True)
            time.sleep(0.05)

    t = threading.Thread(target=_show, daemon=True)
    t.start()
    input()
    _stop.set()

    sx, sy = pyautogui.position()
    bx, by = sx - win.left, sy - win.top
    print(f"\n已記錄搜尋框座標：({bx}, {by})")
    return bx, by


def is_in_auction_screen(win) -> bool:
    """截圖並 OCR，判斷目前是否仍在拍賣畫面。"""
    img, _, _ = _capture_full(win)
    results = _get_reader().readtext(img)
    scale = _get_scale_factor(win)
    for _, text, _ in results:
        t = text.replace(" ", "")
        if "請輸入" in t or "每個" in t:
            return True
    return False


def reset_auction_caches() -> None:
    """清除介面元素快取，讓 verify_price_header 在重進後重新偵測。"""
    global _search_box_cache, _price_header_cache
    _search_box_cache = None
    _price_header_cache = {}


def enter_auction(win, btn_pos: tuple[int, int] | None = None) -> bool:
    """
    點擊拍賣按鈕進入拍賣畫面，回傳是否成功。
    btn_pos 優先，否則使用 AUCTION_BTN_POS（.env）。
    """
    bx, by = btn_pos if btn_pos else AUCTION_BTN_POS
    if bx == 0 and by == 0:
        logger.error("未設定拍賣按鈕座標（AUCTION_BTN_POS），請先執行 --set-auction-btn")
        return False

    try:
        win.activate()
    except Exception:
        pass
    time.sleep(0.5)

    pyautogui.click(win.left + bx, win.top + by)
    logger.info(f"  已點擊拍賣按鈕 ({bx}, {by})，等待畫面載入...")
    time.sleep(3.0)

    reset_auction_caches()

    # 座標已校準時跳過全視窗 OCR 驗證（省去 10~15 秒），直接信任點擊結果
    sx, sy = SEARCH_BOX_POS
    px, py = PRICE_SORT_POS
    if (sx != 0 or sy != 0) and (px != 0 or py != 0):
        logger.info("  座標已校準，略過 OCR 驗證，視為成功進入拍賣畫面")
        return True

    if is_in_auction_screen(win):
        logger.info("  成功進入拍賣畫面")
        return True

    logger.warning("  點擊後仍未偵測到拍賣介面，可能需要調整按鈕座標")
    return False


def reenter_auction(win, wait_secs: int = AUCTION_WAIT_SECS, btn_pos: tuple[int, int] | None = None) -> bool:
    """被踢出後等待 wait_secs 秒再重新進入。"""
    logger.info(f"  離開拍賣畫面，等待 {wait_secs} 秒後重新進入...")
    time.sleep(wait_secs)
    return enter_auction(win, btn_pos=btn_pos)


# ---------------------------------------------------------------------------
# 主要入口
# ---------------------------------------------------------------------------

def verify_price_header(win) -> None:
    """
    確認搜尋框與「每個價錢」欄位座標已就緒。
    若 .env 已設定 SEARCH_BOX_POS 與 PRICE_SORT_POS，直接使用，完全跳過 OCR。
    否則做一次 OCR 偵測並快取，後續呼叫沿用快取。
    """
    global _search_box_cache, _price_header_cache

    if _search_box_cache is not None and _UNIVERSAL_CACHE_KEY in _price_header_cache:
        logger.info("  沿用已快取的介面元素座標")
        return

    # .env 已有固定座標：直接套用，不做 OCR
    sx, sy = SEARCH_BOX_POS
    px, py = PRICE_SORT_POS
    if (sx != 0 or sy != 0) and (px != 0 or py != 0):
        _search_box_cache = (sx, sy)
        _price_header_cache[_UNIVERSAL_CACHE_KEY] = (px, py, max(0, px - 60), px + 60, py + 20)
        logger.info(f"  使用 .env 座標：搜尋框=({sx},{sy})，排序欄=({px},{py})")
        return

    # .env 只設了其中一個：先套用已設定的，剩下的才用 OCR
    if _search_box_cache is None and (sx != 0 or sy != 0):
        _search_box_cache = (sx, sy)
        logger.info(f"  使用 .env 搜尋框座標：({sx},{sy})")

    if _UNIVERSAL_CACHE_KEY not in _price_header_cache and (px != 0 or py != 0):
        _price_header_cache[_UNIVERSAL_CACHE_KEY] = (px, py, max(0, px - 60), px + 60, py + 20)
        logger.info(f"  使用 .env 排序欄位座標：({px},{py})")

    if _search_box_cache is not None and _UNIVERSAL_CACHE_KEY in _price_header_cache:
        return

    logger.info("  偵測介面元素（搜尋框 + 每個價錢欄位）...")
    img, _, _ = _capture_full(win)
    results = _get_reader().readtext(img)
    scale = _get_scale_factor(win)

    if _search_box_cache is None:
        for bbox, text, _ in results:
            if any(k in text.replace(" ", "") for k in ["請輸入", "道具名稱", "輸入道具"]):
                cx = int((bbox[0][0] + bbox[2][0]) / 2 / scale)
                cy = int((bbox[0][1] + bbox[2][1]) / 2 / scale)
                _search_box_cache = (cx, cy)
                logger.info(f"  搜尋框：center=({cx},{cy})")
                break
        if _search_box_cache is None:
            raise RuntimeError("無法偵測搜尋輸入框，請確認遊戲畫面正確顯示拍賣介面")

    if _UNIVERSAL_CACHE_KEY not in _price_header_cache:
        for bbox, text, _ in results:
            if "每個" in text.replace(" ", ""):
                x1, x2 = int(bbox[0][0] / scale), int(bbox[2][0] / scale)
                cy = int((bbox[0][1] + bbox[2][1]) / 2 / scale)
                cx = (x1 + x2) // 2
                pad = max(60, (x2 - x1) // 2)
                _price_header_cache[_UNIVERSAL_CACHE_KEY] = (cx, cy, max(0, x1 - pad), x2 + pad, cy + 20)
                logger.info(f"  每個價錢欄位：center=({cx},{cy})")
                return
        raise RuntimeError("無法偵測「每個價錢」欄位，請確認遊戲畫面正確顯示拍賣列表")


def exit_auction(win) -> None:
    """點擊離開拍賣的按鈕。"""
    ex, ey = AUCTION_EXIT_POS
    if ex == 0 and ey == 0:
        logger.warning("  未設定離開拍賣按鈕座標（AUCTION_EXIT_POS），請執行 --set-auction-exit")
        return
    try:
        win.activate()
    except Exception:
        pass
    time.sleep(0.2)
    pyautogui.click(win.left + ex, win.top + ey)
    logger.info(f"  已點擊離開拍賣按鈕 ({ex}, {ey})")
    time.sleep(1.0)


def scrape_item(
    win,
    item_name: str,
    item_type: int = 1,
    equip_region: tuple[int, int, int, int] | None = None,
    default_region: tuple[int, int, int, int] | None = None,
) -> int | None:
    """搜尋單一商品並回傳最低單價。"""
    search_item(win, item_name)
    cx, cy, _, _, _ = _find_price_header(win, item_type)
    pyautogui.click(win.left + cx, win.top + cy)
    logger.debug(f"  已點擊每個價錢排序 ({cx}, {cy})")
    time.sleep(AFTER_SORT_DELAY)

    if item_type == _ITEM_TYPE_EQUIP:
        return read_price_from_region(win, *(equip_region or PRICE_REGION_EQUIP))
    return read_price_from_region(win, *(default_region or PRICE_REGION_DEFAULT))
