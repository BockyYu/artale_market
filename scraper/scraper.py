"""
遊戲視窗自動化模組。
依賴：easyocr、pyautogui、pygetwindow、mss、Pillow、pyperclip、numpy

座標系統：所有 _FIXED_* 常數都是相對於遊戲視窗左上角（win.left, win.top），
由 debug-ocr 校準。
"""

import logging
import re
import time

import easyocr
import mss
import numpy as np
import pyautogui
import pygetwindow as gw
import pyperclip
from PIL import Image

from config import (
    AFTER_SEARCH_DELAY,
    AFTER_SORT_DELAY,
    WINDOW_TITLE,
)

logger = logging.getLogger(__name__)

# OCR 引擎懶初始化
_reader: easyocr.Reader | None = None

# 搜尋輸入框中心（固定座標，相對於視窗左上角）
_SEARCH_BOX_X = 403
_SEARCH_BOX_Y = 77

# 「每個價錢」欄位快取：(center_x, center_y, col_x_min, col_x_max, row_y_min)
_price_header_cache: tuple[int, int, int, int, int] | None = None


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
    wins = gw.getWindowsWithTitle(WINDOW_TITLE)
    if not wins:
        raise RuntimeError(f"找不到遊戲視窗，請確認已開啟：{WINDOW_TITLE}")
    win = wins[0]
    if win.isMinimized:
        win.restore()
        time.sleep(0.3)
    win.activate()
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

def search_item(win, item_name: str) -> None:
    """在搜尋欄輸入商品名稱並送出（固定座標，不需 OCR）。"""
    pyautogui.click(win.left + _SEARCH_BOX_X, win.top + _SEARCH_BOX_Y)
    time.sleep(0.2)
    pyautogui.hotkey("ctrl", "a")
    pyperclip.copy(item_name)
    pyautogui.hotkey("ctrl", "v")
    time.sleep(0.1)
    pyautogui.press("enter")
    logger.debug(f"  搜尋送出: {item_name}")
    time.sleep(AFTER_SEARCH_DELAY)


def _find_price_header(win) -> tuple[int, int, int, int, int]:
    """
    OCR 自動偵測「每個價錢」欄位標題位置。
    回傳 (center_x, center_y, col_x_min, col_x_max, row_y_min)，座標相對於視窗左上角。
    結果會快取，視窗移動後重啟腳本即可重新偵測。
    """
    global _price_header_cache
    if _price_header_cache is not None:
        return _price_header_cache

    logger.info("  自動偵測「每個價錢」欄位位置...")
    img, _, _ = _capture_full(win)
    results = _get_reader().readtext(img)

    # 單一文字塊命中
    for bbox, text, _ in results:
        if "每個價錢" in text.replace(" ", ""):
            x1, x2 = int(bbox[0][0]), int(bbox[2][0])
            cy = int((bbox[0][1] + bbox[2][1]) / 2)
            cx = (x1 + x2) // 2
            pad = max(60, (x2 - x1) // 2)
            _price_header_cache = (cx, cy, max(0, x1 - pad), x2 + pad, cy + 20)
            logger.info(f"  偵測成功（單塊）：center=({cx},{cy})")
            return _price_header_cache

    # 合併同列文字塊後再找
    row_map: dict[int, list] = {}
    for bbox, text, _ in results:
        cy = int((bbox[0][1] + bbox[2][1]) / 2)
        cx = int((bbox[0][0] + bbox[2][0]) / 2)
        row_map.setdefault(round(cy / 15) * 15, []).append(
            (cx, cy, text.replace(" ", ""), int(bbox[0][0]), int(bbox[2][0]))
        )

    for frags in sorted(row_map.values(), key=lambda f: f[0][1]):
        frags.sort(key=lambda f: f[0])
        merged = "".join(f[2] for f in frags)
        if "每個價錢" in merged or ("每個" in merged and "價錢" in merged):
            relevant = [f for f in frags if any(k in f[2] for k in ["每個", "價錢"])] or frags
            cx = sum(f[0] for f in relevant) // len(relevant)
            cy = sum(f[1] for f in relevant) // len(relevant)
            x1 = min(f[3] for f in relevant)
            x2 = max(f[4] for f in relevant)
            pad = max(60, (x2 - x1) // 2)
            _price_header_cache = (cx, cy, max(0, x1 - pad), x2 + pad, cy + 20)
            logger.info(f"  偵測成功（合併）：center=({cx},{cy})")
            return _price_header_cache

    logger.warning("  無法偵測「每個價錢」位置，使用預設座標")
    _price_header_cache = (1195, 243, 1095, 1295, 260)
    return _price_header_cache


def click_sort_by_lowest(win) -> tuple[int, int, int]:
    """點擊「每個價錢」欄位標題排序，回傳 (col_x_min, col_x_max, row_y_min)。"""
    cx, cy, x_min, x_max, y_min = _find_price_header(win)
    pyautogui.click(win.left + cx, win.top + cy)
    logger.debug(f"  已點擊每個價錢排序 ({cx}, {cy})")
    time.sleep(AFTER_SORT_DELAY)
    return x_min, x_max, y_min


def _parse_price(text: str) -> int | None:
    """將 OCR 讀到的價格文字轉成整數，支援「萬」單位（例如 3,099萬 → 30990000）。"""
    text = text.strip()
    if "萬" in text or "万" in text:
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
        logger.debug(f"  OCR 列 y≈{cy}: {merged!r}")
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


# ---------------------------------------------------------------------------
# 主要入口
# ---------------------------------------------------------------------------

def scrape_item(win, item_name: str) -> int | None:
    """搜尋單一商品並回傳最低單價。"""
    search_item(win, item_name)
    x_min, x_max, y_min = click_sort_by_lowest(win)
    return read_lowest_price(win, x_min, x_max, y_min)
