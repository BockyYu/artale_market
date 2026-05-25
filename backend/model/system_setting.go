package model

import "time"

type SystemSetting struct {
	ID           uint      `json:"id"            gorm:"primaryKey;autoIncrement"`
	Name         string    `json:"name"          gorm:"uniqueIndex;not null"`
	Status       bool      `json:"status"        gorm:"default:false"`
	OperatorName string    `json:"operator_name"`
	UpdatedAt    time.Time `json:"updated_at"`
}
