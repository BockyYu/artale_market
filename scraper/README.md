# Artale Market Scraper

自動擷取 MapleStory Worlds-Artale 拍賣場價格，並透過 API 寫入後端資料庫。支援 macOS 與 Windows。

---

## 安裝

### 前置需求

- Python 3.10 以上
- 後端服務正在運行（預設 `http://localhost:8080`）

### 建立虛擬環境並安裝套件

```bash
cd scraper
python -m venv venv

# macOS / Linux
source venv/bin/activate

# Windows
venv\Scripts\activate

pip install -r requirements.txt
```

---

## 首次設定（座標校準）

首次執行前必須校準各個 UI 元素位置，結果存入 `.env` 供後續自動讀取。

```bash
# 1. 校準拍賣按鈕（先回到遊戲主畫面）
python main.py --set-auction-btn

# 2. 手動進入拍賣後，校準搜尋輸入框
python main.py --set-search-box

# 3. 校準「每個價錢」排序欄位標題
python main.py --set-price-sort

# 4. 校準離開拍賣的按鈕
python main.py --set-auction-exit
```

> 換電腦或螢幕解析度改變時需重新校準。

---

## 環境變數（.env）

| 變數 | 預設值 | 說明 |
|---|---|---|
| `API_BASE_URL` | `http://localhost:8080` | 後端 API 位址 |
| `WINDOW_TITLE` | `MapleStory Worlds-Artale` | 遊戲視窗標題（部分符合即可） |
| `SCHEDULE_TIME` | `08:00` | 排程每日執行時間 |
| `SCHEDULE_INTERVAL_MINUTES` | `3` | 間隔模式每次執行間隔（分鐘） |
| `AFTER_SEARCH_DELAY` | `2.5` | 輸入搜尋詞後等待結果載入（秒） |
| `AFTER_SORT_DELAY` | `1.5` | 點擊排序後等待畫面更新（秒） |
| `BETWEEN_ITEMS_DELAY` | `3.0` | 每個商品間隔（秒） |
| `AUCTION_WAIT_SECS` | `300` | 被踢出拍賣後等待重進的秒數（預設 5 分鐘） |
| `AUCTION_BTN_POS` | — | 拍賣按鈕座標（由 `--set-auction-btn` 寫入） |
| `SEARCH_BOX_POS` | — | 搜尋輸入框座標（由 `--set-search-box` 寫入） |
| `PRICE_SORT_POS` | — | 「每個價錢」排序欄位座標（由 `--set-price-sort` 寫入） |
| `AUCTION_EXIT_POS` | — | 離開拍賣按鈕座標（由 `--set-auction-exit` 寫入） |
| `PRICE_REGION_EQUIP` | `1427,412,1713,501` | 裝備價格 OCR 讀取區域（x1,y1,x2,y2） |
| `PRICE_REGION_DEFAULT` | `1722,415,1975,503` | 卷軸／技能書價格 OCR 讀取區域 |
| `NOTIFY_BOT_ID` | `0` | 後台通知機器人 ID |

---

## main.py 指令說明

### 立即執行一次

```bash
python main.py --now
```

進入拍賣場，對所有追蹤商品抓取最低價並寫入 DB。

---

### 漏補模式

```bash
python main.py --now --fill-missing
```

只抓「今天還沒有價格記錄」的商品。適合當天爬蟲中途中斷後補跑，不重複抓已有記錄的道具。

---

### 模擬模式（不寫入資料庫）

```bash
python main.py --now --dry-run
```

只印出抓到的價格，不實際寫入，適合測試。

---

### 排程模式（每日固定時間執行）

```bash
python main.py
python main.py --time 14:30
```

每天在指定時間自動執行一次。按 `Ctrl+C` 停止。

---

### 間隔模式（每 N 分鐘循環執行）

```bash
python main.py --interval
```

依 `SCHEDULE_INTERVAL_MINUTES` 自動循環，啟動後立即執行第一次。按 `Ctrl+C` 停止。

---

### 測試通知

```bash
python main.py --test-notify
```

發送一則測試訊息，確認 TG 通知功能是否正常。

---

### 校準指令

```bash
# 互動式校準各 UI 座標並存入 .env
python main.py --set-auction-btn    # 拍賣按鈕
python main.py --set-search-box     # 搜尋輸入框
python main.py --set-price-sort     # 「每個價錢」排序欄位
python main.py --set-auction-exit   # 離開拍賣按鈕

# 診斷工具
python main.py --track-mouse        # 即時顯示滑鼠相對視窗座標
python main.py --debug-pos 424 76   # 移動滑鼠到指定座標並截圖確認
python main.py --debug-ocr          # 截圖並印出所有 OCR 辨識結果
```

---

## bot.py 指令說明

從後端取得啟用中的價格提醒道具，掃價後寫入 DB，若低於門檻立即發 TG 通知。

```bash
# 啟動（依 SCHEDULE_INTERVAL_MINUTES 循環）
python bot.py

# 校準座標（與 main.py 相同，存入同一個 .env）
python bot.py --set-auction-btn
python bot.py --set-search-box
python bot.py --set-price-sort
python bot.py --set-auction-exit
```

按 `Ctrl+C` 停止。

---

## 自動重進拍賣

若連續 2 筆道具找不到價格，程式會自動判斷是否已被踢出拍賣畫面，若是則等待 `AUCTION_WAIT_SECS` 秒（預設 5 分鐘）後自動重新進入，不需人工干預。

---

## 檔案說明

| 檔案 | 說明 |
|---|---|
| `main.py` | 主程式入口，處理參數與排程 |
| `bot.py` | 價格提醒掃價 Bot |
| `scraper.py` | 遊戲視窗自動化與 OCR 辨識邏輯 |
| `api_client.py` | 後端 API 呼叫（取得商品、寫入價格、查詢今日價格） |
| `notify.py` | TG / LINE 通知發送 |
| `config.py` | 全域設定（從 .env 讀取） |
| `requirements.txt` | Python 套件相依清單 |
