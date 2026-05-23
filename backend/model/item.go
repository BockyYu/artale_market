package model

import "time"

type ItemType int

const (
	ItemTypeScroll   ItemType = 1
	ItemTypeMaterial ItemType = 2
	ItemTypeConsume  ItemType = 3
	ItemTypeSkillBook ItemType = 4
)

type TrackPriority int

const (
	TrackPriorityOff       TrackPriority = 0 // 不追蹤
	TrackPriorityPrimary   TrackPriority = 1 // 優先查詢
	TrackPrioritySecondary TrackPriority = 2 // 次要查詢
	TrackPriorityNotSeen   TrackPriority = 3 // 尚未出現
)

type Item struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string    `json:"name" binding:"required"`
	Percentage    int       `json:"percentage"`
	Category      string    `json:"category" binding:"required"`
	Description   string    `json:"description"`
	ItemType      ItemType  `json:"item_type" gorm:"default:1"`
	TrackPriority TrackPriority `json:"track_priority" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
