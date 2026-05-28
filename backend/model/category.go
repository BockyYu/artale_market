package model

type Category struct {
	ID       uint   `json:"id"        gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name"      gorm:"uniqueIndex:idx_cat_name_type;not null"`
	ItemType int    `json:"item_type" gorm:"uniqueIndex:idx_cat_name_type;not null"`
}
