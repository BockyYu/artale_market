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
)

// TrackPriority 追蹤優先等級
type TrackPriority int

const (
	TrackPriorityOff       TrackPriority = 0 // 不追蹤
	TrackPriorityPrimary   TrackPriority = 1 // 優先查詢
	TrackPrioritySecondary TrackPriority = 2 // 次要查詢
	TrackPriorityNotSeen   TrackPriority = 3 // 尚未出現
)

// ItemAdminRow 後台道具列表（含最新價格）
type ItemAdminRow struct {
	Item
	LatestPrice *float64 `json:"latest_price" gorm:"column:latest_price"`
}

// Item 商品
type Item struct {
	ID            uint          `json:"id"             gorm:"primaryKey;autoIncrement"`
	Name          string        `json:"name"           binding:"required"`
	Percentage    int           `json:"percentage"`                          // 卷軸成功率（技能書不使用）
	Category      string        `json:"category"       binding:"required"`   // 分類：防具部位 / 武器種類 / 職業
	Description   string        `json:"description"`                         // 備註
	ItemType      ItemType      `json:"item_type"      gorm:"default:1;index"` // 商品類型
	TrackPriority TrackPriority `json:"track_priority" gorm:"default:0"`     // 追蹤優先等級
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
