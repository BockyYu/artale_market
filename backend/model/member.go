package model

import "time"

// Member 一般會員帳號
type Member struct {
	ID        uint      `json:"id"         gorm:"primaryKey;autoIncrement"`
	Nickname  string    `json:"nickname"   gorm:"not null"`
	Username  string    `json:"username"   gorm:"uniqueIndex;not null"`
	Password  string    `json:"-"          gorm:"not null"`
	Email     string    `json:"email"      gorm:"uniqueIndex;not null"`
	Status    int       `json:"status"     gorm:"default:1"` // 1=正常 0=封禁
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
