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

# ── 固定座標（根據 debug-ocr 校準）────────────────────────────────────────────
# 搜尋輸入框中心
_SEARCH_BOX_X = 403
_SEARCH_BOX_Y = 77

# 「每個價錢」欄位標題（點擊後按單價排序）
_PRICE_HEADER_X = 1195
_PRICE_HEADER_Y = 243

# 「每個價錢」欄位的 x 範圍 & 資料列起始 y（用於 OCR 裁切）
_PRICE_COL_X_MIN = 1095
_PRICE_COL_X_MAX = 1295
_PRICE_ROW_Y_MIN = 260   # 標題列以下才是資料
# ─────────────────────────────────────────────────────────────────────────────


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


def _capture_column(win) -> np.ndarray:
    """只截取「每個價錢」欄位的資料區域，供 OCR 使用。"""
    with mss.MSS() as sct:
        region = {
            "left": max(0, win.left + _PRICE_COL_X_MIN),
            "top": max(0, win.top + _PRICE_ROW_Y_MIN),
            "width": _PRICE_COL_X_MAX - _PRICE_COL_X_MIN,
            "height": win.height - _PRICE_ROW_Y_MIN,
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


def click_sort_by_lowest(win) -> None:
    """點擊「每個價錢」欄位標題，按單價由低到高排序（固定座標，不需 OCR）。"""
    pyautogui.click(win.left + _PRICE_HEADER_X, win.top + _PRICE_HEADER_Y)
    logger.debug("  已點擊每個價錢排序")
    time.sleep(AFTER_SORT_DELAY)


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


def read_lowest_price(win) -> int | None:
    """
    讀取排序後第一行的「每個價錢」。
    只對欄位資料區做 OCR（小圖），取 y 最小的有效數字。
    """
    img = _capture_column(win)
    results = _get_reader().readtext(img)

    candidates: list[tuple[int, int]] = []
    for bbox, text, conf in results:
        cy = int((bbox[0][1] + bbox[2][1]) / 2)
        # 跳過括號內的合計列（例如「(8萬 3,000)」）
        if "(" in text or "（" in text:
            continue
        price = _parse_price(text)
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
    click_sort_by_lowest(win)
    return read_lowest_price(win)
