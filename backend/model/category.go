package model

type Category struct {
	ID       uint     `json:"id"        gorm:"primaryKey;autoIncrement"`
	Name     string   `json:"name"      gorm:"uniqueIndex:idx_cat_name_type;not null"`
	ItemType ItemType `json:"item_type" gorm:"uniqueIndex:idx_cat_name_type;not null"`
}

var ScrollCategories = []string{
	"頭盔", "上衣", "下衣", "套服", "鞋子", "手套", "披風", "盾牌",
	"臉部", "眼部", "耳環", "戒指", "墜飾", "腰帶", "肩章", "勳章",
	"單手劍", "雙手劍", "單手斧", "雙手斧", "單手棍", "雙手棍",
	"槍", "矛", "短杖", "長杖", "弓", "弩", "短劍", "拳套", "指虎", "火槍",
}

var SkillBookCategories = []string{
	"全職業共通",
	"劍士", "英雄", "聖騎士", "黑騎士",
	"法師", "火毒", "冰雷", "主教",
	"弓手", "箭神", "神射手",
	"盜賊", "神偷", "夜使者",
	"槍神", "拳霸",
}

var EquipCategories = []string{
	"全職業",
	"頭盔", "上衣", "下衣", "套服", "鞋子", "手套", "披風", "盾牌",
	"臉部", "眼部", "耳環", "戒指", "墜飾", "腰帶", "肩章", "勳章",
	"武器",
	"單手劍", "雙手劍", "單手斧", "雙手斧", "單手棍", "雙手棍",
	"槍", "矛", "短杖", "長杖", "弓", "弩", "短劍", "拳套", "指虎", "火槍",
}
