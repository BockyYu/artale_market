package model

import "time"

// ItemType 商品類型
type ItemType int

const (
	ItemTypeScroll    ItemType = 1 // 卷軸
	ItemTypeMaterial  ItemType = 2 // 素材
	ItemTypeConsume   ItemType = 3 // 消耗品
	ItemTypeSkillBook ItemType = 4 // 技能書
	ItemTypePoints    ItemType = 5 // 商城
	ItemTypeEquip    ItemType = 6 // 裝備
	ItemTypeEvent    ItemType = 7 // 活動道具
)

// TrackPriority 追蹤優先等級
type TrackPriority int

const (
	TrackPriorityOff       TrackPriority = 0 // 不追蹤
	TrackPriorityPrimary   TrackPriority = 1 // 優先查詢
	TrackPrioritySecondary TrackPriority = 2 // 次要查詢
	TrackPriorityNotSeen   TrackPriority = 3 // 尚未出現
)

// ItemAdminRow 後台道具列表（含最新價格與今日統計）
type ItemAdminRow struct {
	Item
	LatestPrice   *float64   `json:"latest_price"    gorm:"column:latest_price"`
	LatestPriceAt *time.Time `json:"latest_price_at" gorm:"column:latest_price_at"`
	TodayChanges  int        `json:"today_changes"   gorm:"column:today_changes"`
	TodayViews    int        `json:"today_views"     gorm:"-"`
}

// AdminItemPage 後台道具分頁回傳
type AdminItemPage struct {
	Data     []ItemAdminRow `json:"data"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// Item 商品
type Item struct {
	ID            uint          `json:"id"             gorm:"primaryKey;autoIncrement"`
	Name          string        `json:"name"           binding:"required"`
	EnglishName   string        `json:"english_name"`                        // 英文名稱
	SearchMode    int           `json:"search_mode"    gorm:"default:1"`     // 查詢方式 1=中文 2=英文
	Percentage    int           `json:"percentage"`                          // 卷軸成功率（技能書不使用）
	Category      string        `json:"category"`                            // 分類：防具部位 / 武器種類 / 職業
	Description   string        `json:"description"`                         // 備註
	ItemType      ItemType      `json:"item_type"      gorm:"default:1;index"` // 商品類型
	TrackPriority TrackPriority `json:"track_priority" gorm:"default:0"`     // 追蹤優先等級
	IsHidden      bool          `json:"-"              gorm:"default:false"`  // 隱藏（軟刪除）
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
