import os

from pathlib import Path
from dotenv import load_dotenv
load_dotenv(Path(__file__).parent / ".env")

# 後端 API 位址
API_BASE_URL = os.getenv("API_BASE_URL", "http://localhost:8080")

# 遊戲視窗標題（完整或部分符合即可）
WINDOW_TITLE = os.getenv("WINDOW_TITLE", "MapleStory Worlds-Artale")

# 每天自動執行時間，格式 HH:MM（24小時制）
# 也可以在啟動時用 --time 14:30 覆蓋
SCHEDULE_TIME = os.getenv("SCHEDULE_TIME", "08:00")

# 間隔模式：每隔幾分鐘執行一次（--interval 模式 / bot.py 使用）
SCHEDULE_INTERVAL_MINUTES = int(os.getenv("SCHEDULE_INTERVAL_MINUTES", "1"))

# ---- Bot 帳號設定 ----
# bot.py 需要登入後台才能取得提醒列表
BOT_ADMIN_USERNAME = os.getenv("BOT_ADMIN_USERNAME", "")
BOT_ADMIN_PASSWORD = os.getenv("BOT_ADMIN_PASSWORD", "")

# ---- 通知設定 ----
# 優先：後端 Bot ID（後台「通知機器人」頁面的 ID），設定後由後端發送，不需填寫 TG 憑證
NOTIFY_BOT_ID = int(os.getenv("NOTIFY_BOT_ID", "0"))
# 備用：直接發 Telegram（不設定 NOTIFY_BOT_ID 時才使用）
NOTIFY_TG_BOT_TOKEN = os.getenv("NOTIFY_TG_BOT_TOKEN", "")
NOTIFY_TG_CHAT_ID = os.getenv("NOTIFY_TG_CHAT_ID", "")
# LINE Notify：設定 TOKEN 即可啟用（未來支援）
NOTIFY_LINE_TOKEN = os.getenv("NOTIFY_LINE_TOKEN", "")

# ---- 等待時間（秒）----
# 輸入搜尋詞後等待結果載入的時間
AFTER_SEARCH_DELAY = float(os.getenv("AFTER_SEARCH_DELAY", "2.5"))
# 點擊排序按鈕後等待畫面更新的時間
AFTER_SORT_DELAY = float(os.getenv("AFTER_SORT_DELAY", "1.5"))
# 每個商品之間的間隔，避免操作太快
BETWEEN_ITEMS_DELAY = float(os.getenv("BETWEEN_ITEMS_DELAY", "3.0"))
# 點擊裝備列表第一筆後，等待詳情面板載入的時間
AFTER_ROW_CLICK_DELAY = float(os.getenv("AFTER_ROW_CLICK_DELAY", "1.0"))

# 價格欄位擷取區域（視窗相對座標，由 --track-mouse 校準）
# 格式：x1,y1,x2,y2
# 裝備類型（左側欄）
_equip_region = os.getenv("PRICE_REGION_EQUIP", "1427,412,1713,501").split(",")
PRICE_REGION_EQUIP: tuple[int, int, int, int] = tuple(int(v) for v in _equip_region)  # type: ignore
# 一般道具（卷軸 / 技能書，右側欄）
_default_region = os.getenv("PRICE_REGION_DEFAULT", "1722,415,1975,503").split(",")
PRICE_REGION_DEFAULT: tuple[int, int, int, int] = tuple(int(v) for v in _default_region)  # type: ignore
