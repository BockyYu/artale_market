package model

import "time"

// PriceHistory 每次價格異動的完整記錄（含日內多次變動）
type PriceHistory struct {
	ID         uint      `json:"id"          gorm:"primaryKey;autoIncrement"`
	ItemID     uint      `json:"item_id"     gorm:"index:idx_ph_item_time,priority:1"`
	Price      float64   `json:"price"`
	Source     string    `json:"source"      gorm:"default:scraper"` // scraper / admin
	RecordedAt time.Time `json:"recorded_at" gorm:"index:idx_ph_item_time,priority:2;autoCreateTime"`
}
