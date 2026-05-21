import os

# 後端 API 位址
API_BASE_URL = os.getenv("API_BASE_URL", "http://localhost:8080")

# 遊戲視窗標題（完整或部分符合即可）
WINDOW_TITLE = os.getenv("WINDOW_TITLE", "MapleStory Worlds-Artale")

# 每天自動執行時間，格式 HH:MM（24小時制）
# 也可以在啟動時用 --time 14:30 覆蓋
SCHEDULE_TIME = os.getenv("SCHEDULE_TIME", "08:00")

# ---- 等待時間（秒）----
# 輸入搜尋詞後等待結果載入的時間
AFTER_SEARCH_DELAY = float(os.getenv("AFTER_SEARCH_DELAY", "2.5"))
# 點擊排序按鈕後等待畫面更新的時間
AFTER_SORT_DELAY = float(os.getenv("AFTER_SORT_DELAY", "1.5"))
# 每個商品之間的間隔，避免操作太快
BETWEEN_ITEMS_DELAY = float(os.getenv("BETWEEN_ITEMS_DELAY", "3.0"))
