package dto

import "artale_market/model"

// UpdateItemReq 更新商品資訊請求
type UpdateItemReq struct {
	Name        string           `json:"name"         binding:"required"` // 商品名稱
	EnglishName string           `json:"english_name"`                    // 英文名稱（選填）
	SearchMode  int              `json:"search_mode"  binding:"min=1,max=2"` // 查詢方式 1=中文 2=英文
	ItemType    model.ItemType   `json:"item_type"    binding:"required"` // 商品類型
	Percentage  int              `json:"percentage"   binding:"min=0,max=100"` // 卷軸成功率
	Category    string           `json:"category"     binding:"required"` // 分類
	Description string           `json:"description"` // 備註（選填）
}

// SetTrackedReq 設定追蹤優先等級請求
type SetTrackedReq struct {
	TrackPriority model.TrackPriority `json:"track_priority" binding:"min=0,max=2"` // 0=不追蹤 1=優先 2=次要
}
