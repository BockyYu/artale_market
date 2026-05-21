"""快速截圖工具 — 只截遊戲視窗，存成 screenshot.png"""
import sys
import mss
import pygetwindow as gw
from PIL import Image

TITLE = "MapleStory Worlds-Artale"

wins = gw.getWindowsWithTitle(TITLE)
if not wins:
    print(f"找不到視窗：{TITLE}")
    print("目前所有視窗：")
    for t in gw.getAllTitles():
        if t.strip():
            print(f"  {repr(t)}")
    sys.exit(1)

win = wins[0]
print(f"找到視窗：{repr(win.title)}")
print(f"位置：left={win.left}, top={win.top}, width={win.width}, height={win.height}")

if win.isMinimized:
    win.restore()

win.activate()

import time; time.sleep(0.5)

with mss.mss() as sct:
    region = {
        "left": max(0, win.left),
        "top": max(0, win.top),
        "width": win.width,
        "height": win.height,
    }
    shot = sct.grab(region)

img = Image.frombytes("RGB", shot.size, shot.rgb)
img.save("screenshot.png")
print("截圖已儲存：screenshot.png")
