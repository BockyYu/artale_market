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
| `SCHEDULE_INTERVAL_MINUTES` | `5` | 間隔模式每次執行間隔（分鐘） |
| `AFTER_SEARCH_DELAY` | `2.5` | 輸入搜尋詞後等待結果載入（秒） |
| `AFTER_SORT_DELAY` | `1.5` | 點擊排序後等待畫面更新（秒） |
| `BETWEEN_ITEMS_DELAY` | `3.0` | 每個商品間隔（秒） |
| `PRICE_REGION_EQUIP` | `1427,412,1713,501` | 裝備價格擷取區域（視窗相對座標 x1,y1,x2,y2） |
| `PRICE_REGION_DEFAULT` | `1722,415,1975,503` | 卷軸／技能書價格擷取區域（視窗相對座標 x1,y1,x2,y2） |

---

## 指令說明（main.py）

### 立即執行一次

```bash
python main.py --now
```

從後端取得今日待追蹤商品，開啟遊戲視窗自動搜尋並記錄最低價格。

---

### 模擬模式（不寫入資料庫）

```bash
python main.py --now --dry-run
```

只印出抓到的價格，不實際呼叫 API 寫入，適合測試用。

---

### 排程模式（每日固定時間執行）

```bash
# 使用預設時間（08:00）
python main.py

# 指定每日執行時間
python main.py --time 14:30
```

按 `Ctrl+C` 停止。

---

### 間隔模式（每 N 分鐘執行一次）

```bash
python main.py --interval
python main.py --interval --dry-run
```

依 `SCHEDULE_INTERVAL_MINUTES`（預設 5 分鐘）自動循環執行，啟動後立即執行第一次。按 `Ctrl+C` 停止。

---

### 指定價格擷取區域

```bash
# 覆蓋裝備擷取區域
python main.py --now --equip-region 1427,412,1713,501

# 覆蓋卷軸／技能書擷取區域
python main.py --now --default-region 1722,415,1975,503

# 兩個都指定
python main.py --now --equip-region 1427,412,1713,501 --default-region 1722,415,1975,503
```

不帶這兩個 flag 時，使用 `config.py` 的預設值。換電腦或螢幕解析度改變時，先用 `--track-mouse` 重新校準後再更新 `config.py`。

---

### 校準工具

#### 即時顯示滑鼠座標

```bash
python main.py --track-mouse
```

即時顯示滑鼠相對遊戲視窗的座標，用於校準價格擷取區域的四個角落。按 `Ctrl+C` 停止。

#### 確認點擊位置

```bash
python main.py --debug-pos 1427 412
```

移動滑鼠到指定視窗座標，3 秒後截圖並用紅圈標記位置，存為 `debug_pos_X_Y.png`。用於確認座標是否正確。

#### OCR 診斷

```bash
python main.py --debug-ocr
```

截圖並印出所有 OCR 辨識結果（含座標與信心值），用於排查找不到欄位的問題。截圖存為 `debug_screenshot.png`。

---

## 提醒 Bot（bot.py）

從後端取得啟用中的價格提醒道具，掃描後寫入 DB，由後端自動比對門檻並發送通知。

```bash
# 使用預設設定啟動
python bot.py

# 指定擷取區域
python bot.py --equip-region 1427,412,1713,501 --default-region 1722,415,1975,503
```

依 `SCHEDULE_INTERVAL_MINUTES` 自動循環執行，按 `Ctrl+C` 停止。

---

## 座標校準流程

換電腦、換螢幕解析度或遊戲視窗大小改變時需重新校準：

1. 開遊戲，進入拍賣介面
2. 執行 `python main.py --track-mouse`，移動滑鼠到價格欄位的左上角和右下角，記錄座標
3. 執行 `python main.py --debug-pos X Y` 確認截圖位置正確
4. 更新 `config.py` 中的 `PRICE_REGION_EQUIP` 和 `PRICE_REGION_DEFAULT`

---

## 檔案說明

| 檔案 | 說明 |
|---|---|
| `main.py` | 主程式入口，處理參數與排程 |
| `bot.py` | 提醒道具掃價 Bot |
| `scraper.py` | 遊戲視窗自動化與 OCR 辨識邏輯 |
| `api_client.py` | 後端 API 呼叫（取得商品、寫入價格） |
| `config.py` | 全域設定（API 位址、延遲時間、擷取區域等） |
| `requirements.txt` | Python 套件相依清單 |
