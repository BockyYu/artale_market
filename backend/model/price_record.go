package model

import "time"

// PriceRecord 單筆價格記錄
type PriceRecord struct {
	ID           uint      `json:"id"            gorm:"primaryKey;autoIncrement"`
	ItemID       uint      `json:"item_id"       gorm:"index:idx_pr_item_date,priority:1"` // 對應 items.id
	Price        float64   `json:"price"`                                                   // 記錄價格
	RecordedDate time.Time `json:"recorded_date" gorm:"type:date;index:idx_pr_item_date,priority:2"` // 記錄日期（複合索引）
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PagedSummary 分頁查詢結果
type PagedSummary struct {
	Data     []PriceSummary `json:"data"`      // 當頁資料
	Total    int64          `json:"total"`     // 總筆數
	Page     int            `json:"page"`      // 當前頁碼
	PageSize int            `json:"page_size"` // 每頁筆數
}

// PriceSummary 商品價格摘要（跨多日聚合，用於列表頁）
type PriceSummary struct {
	ItemID            uint       `json:"item_id"              gorm:"column:item_id"`
	ItemName          string     `json:"item_name"            gorm:"column:item_name"`
	ItemPercentage    int        `json:"item_percentage"      gorm:"column:item_percentage"`  // 卷軸成功率
	ItemType          ItemType   `json:"item_type"            gorm:"column:item_type"`
	Category          string     `json:"category"             gorm:"column:category"`
	Description       string     `json:"description"          gorm:"column:description"`
	TodayPrice        *float64   `json:"today_price"          gorm:"column:today_price"`          // 今日價格，無記錄為 null
	TodayCreatedAt    *time.Time `json:"today_created_at"     gorm:"column:today_created_at"`     // 今日首次記錄時間
	TodayUpdatedAt    *time.Time `json:"today_updated_at"     gorm:"column:today_updated_at"`     // 今日最後更新時間
	YesterdayPrice    *float64   `json:"yesterday_price"      gorm:"column:yesterday_price"`      // 昨日價格，無記錄為 null
	ThreeDaysAgoPrice *float64   `json:"three_days_ago_price" gorm:"column:three_days_ago_price"` // 三天前價格，無記錄為 null
	ChangePercent     *float64   `json:"change_percent"       gorm:"-"`                           // 漲跌幅（%），由 service 層計算，不對應欄位
}
