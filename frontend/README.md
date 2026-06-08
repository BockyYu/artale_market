# Artale Market — Frontend

React + Vite 前端，提供市場行情查詢介面與後台管理系統。

## 技術棧

| 套件 | 版本 |
|------|------|
| React | 18 |
| React Router DOM | 7 |
| Vite | 5 |

## 啟動開發

```bash
cd frontend
npm install      # 首次需要
npm run dev      # http://localhost:5173
```

其他指令：

```bash
npm run build    # 打包至 dist/
npm run preview  # 預覽打包結果
```

## 目錄結構

```
src/
├── main.jsx              # 進入點，設定 Router
├── App.jsx               # 主頁面（市場行情）
├── App.css               # 主頁樣式
├── member-api.js         # 一般用戶 API 封裝
├── MemberAuth.jsx        # 會員登入元件
├── Portfolio.jsx         # 自選清單
├── PotionTable.jsx       # 藥水參考表
├── potionData.js         # 藥水靜態資料
└── admin/
    ├── Layout.jsx        # 後台共用 Layout（側欄導覽）
    ├── Login.jsx         # 後台登入頁
    ├── api.js            # 後台 API 封裝
    ├── admin.css         # 後台樣式
    ├── Admins.jsx        # 管理員帳號管理
    ├── Members.jsx       # 會員列表
    ├── Items.jsx         # 道具列表
    ├── NotifyBots.jsx    # 通知機器人管理
    └── PriceAlerts.jsx   # 價格提醒管理
```

## 頁面路由

| 路徑 | 說明 |
|------|------|
| `/` | 市場行情（卷軸 / 職業技能書），支援搜尋與分類篩選 |
| `/admin/login` | 後台登入 |
| `/admin/admins` | 管理員帳號 CRUD |
| `/admin/members` | 會員列表 |
| `/admin/items` | 道具列表 |
| `/admin/bots` | 通知機器人（Telegram / LINE / Discord）|
| `/admin/alerts` | 價格提醒，低於門檻時透過機器人發送通知 |

## API 對接

- 一般用戶：`/api/v1/member/...`，封裝於 `member-api.js`
- 後台管理：`/api/v1/admin/...`，封裝於 `admin/api.js`，需 JWT 驗證

後台 token 儲存於 `localStorage`，過期前 4 小時自動 refresh（由 `Layout.jsx` 定時觸發）。

## 注意事項

- 僅支援桌機瀏覽器（螢幕寬度 < 768px 時顯示提示並擋住畫面）
- 後台需先以管理員帳號登入，token 過期後會自動導回登入頁
