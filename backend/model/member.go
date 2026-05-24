package model

import "time"

// Member 一般會員帳號
type Member struct {
	ID        uint      `json:"id"         gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username"   gorm:"uniqueIndex;not null"` // 登入帳號，唯一
	Password  string    `json:"-"          gorm:"not null"`             // 密碼雜湊，不輸出 JSON
	Email     string    `json:"email"`                                  // 電子郵件（選填）
	Status    int       `json:"status"     gorm:"default:1"`            // 狀態：1=正常 0=封禁
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
