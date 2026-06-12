package model

import "time"

// NotifyBot 通知機器人設定
type NotifyBot struct {
	ID        uint      `json:"id"        gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name"`
	Platform  string    `json:"platform"` // tg | line
	Token     string    `json:"token"`    // TG bot token / LINE Notify token / Discord webhook URL
	ChatID    string    `json:"chat_id"`  // TG 專用，其他平台留空
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
