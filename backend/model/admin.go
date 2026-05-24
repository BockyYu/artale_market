package model

import "time"

// AdminUser 後台管理員帳號
type AdminUser struct {
	ID        uint      `json:"id"         gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username"   gorm:"uniqueIndex;not null"` // 登入帳號，唯一
	Password  string    `json:"-"          gorm:"not null"`             // 密碼雜湊，不輸出 JSON
	Role      string    `json:"role"       gorm:"default:'admin'"`      // 角色：admin / super_admin
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
