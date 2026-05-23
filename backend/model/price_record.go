package model

import "time"

type PriceRecord struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ItemID       uint      `json:"item_id"`
	Price        float64   `json:"price"`
	RecordedDate time.Time `json:"recorded_date" gorm:"type:date"`
	CreatedAt    time.Time `json:"created_at"`
}

type PriceSummary struct {
	ItemID            uint     `json:"item_id"`
	ItemName          string   `json:"item_name"`
	ItemPercentage    int      `json:"item_percentage"`
	ItemType          ItemType `json:"item_type"`
	Category          string   `json:"category"`
	Description       string   `json:"description"`
	TodayPrice        *float64 `json:"today_price"`
	YesterdayPrice    *float64 `json:"yesterday_price"`
	ThreeDaysAgoPrice *float64 `json:"three_days_ago_price"`
	ChangePercent     *float64 `json:"change_percent"`
}
