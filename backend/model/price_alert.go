package model

import "time"

type PriceAlert struct {
	ID              uint        `json:"id"               gorm:"primaryKey;autoIncrement"`
	ItemID          uint        `json:"item_id"          gorm:"index"`
	Item            Item        `json:"item"             gorm:"foreignKey:ItemID"`
	BotID           *uint       `json:"bot_id"           gorm:"index"`
	Bot             *NotifyBot  `json:"bot"              gorm:"foreignKey:BotID"`
	ThresholdPrice  float64     `json:"threshold_price"`
	IsActive        bool        `json:"is_active"        gorm:"default:true"`
	Note            string      `json:"note"`
	LastTriggeredAt *time.Time  `json:"last_triggered_at"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}
