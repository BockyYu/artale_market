# Artale Market Scraper

自動擷取 MapleStory Worlds-Artale 拍賣場價格，並透過 API 寫入後端資料庫。

---

## 安裝

### 前置需求

- Python 3.10 以上
- 後端服務正在運行（預設 `http://localhost:8080`）
- 遊戲視窗標題需包含 `MapleStory Worlds-Artale`

### 建立虛擬環境並安裝套件

```bash
cd scraper
python -m venv .venv

# Windows
.venv\Scripts\activate

# macOS / Linux
source .venv/bin/activate

pip install -r requirements.txt
```

---

## 環境設定

可透過環境變數覆寫預設值，或直接修改 `config.py`：

| 環境變數 | 預設值 | 說明 |
|---|---|---|
| `API_BASE_URL` | `http://localhost:8080` | 後端 API 位址 |
| `WINDOW_TITLE` | `MapleStory Worlds-Artale` | 遊戲視窗標題（部分符合即可） |
| `SCHEDULE_TIME` | `08:00` | 排程每日執行時間 |
| `AFTER_SEARCH_DELAY` | `2.5` | 輸入搜尋詞後等待結果載入（秒） |
| `AFTER_SORT_DELAY` | `1.5` | 點擊排序後等待畫面更新（秒） |
| `BETWEEN_ITEMS_DELAY` | `3.0` | 每個商品間隔（秒） |
| `AFTER_ROW_CLICK_DELAY` | `1.0` | 點擊列表列後等待詳情載入（秒） |

---

## 指令說明

### 立即執行一次

```bash
python main.py --now
```

從後端取得今日待追蹤商品，開啟遊戲視窗自動搜尋並記錄最低價格。

---

### 排程模式（自動每日執行）

```bash
python main.py
```

預設每天 `08:00` 自動執行，按 `Ctrl+C` 停止。

```bash
python main.py --time 14:30
```

指定每日執行時間。

---

### 模擬模式（不寫入資料庫）

```bash
python main.py --now --dry-run
```

只印出抓到的價格，不實際呼叫 API 寫入，適合測試用。

---

### 自動重新抓取 priority=1 商品

```bash
python main.py --autoupdate
```

針對 **查詢優先度＝優先（priority=1）** 的所有商品，自動開啟遊戲視窗以 OCR 重新抓取最新價格並寫入資料庫。

與 `--now` 的差異：`--now` 只抓今天還沒有價格的追蹤商品；`--autoupdate` 不論今天是否已有記錄，都會重新抓取 priority=1 的全部商品。

搭配 `--dry-run` 可先確認抓到的價格再決定是否寫入：

```bash
python main.py --autoupdate --dry-run
```

---

### OCR 診斷模式

```bash
python main.py --debug-ocr
```

截圖並印出所有 OCR 辨識結果（含座標與信心值），用於排查找不到欄位的問題。截圖會存為 `debug_screenshot.png`。

---

## 檔案說明

| 檔案 | 說明 |
|---|---|
| `main.py` | 主程式入口，處理參數與排程 |
| `scraper.py` | 遊戲視窗自動化與 OCR 辨識邏輯 |
| `api_client.py` | 後端 API 呼叫（取得商品、寫入價格） |
| `config.py` | 全域設定（API 位址、延遲時間等） |
| `requirements.txt` | Python 套件相依清單 |
