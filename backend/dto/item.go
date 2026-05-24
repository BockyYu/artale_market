package dto

import "artale_market/model"

// UpdateItemReq 更新商品資訊請求
type UpdateItemReq struct {
	Name        string `json:"name"        binding:"required"`              // 商品名稱
	Percentage  int    `json:"percentage"  binding:"required,min=1,max=100"` // 卷軸成功率（1-100）
	Category    string `json:"category"    binding:"required"`              // 分類
	Description string `json:"description"`                                 // 備註（選填）
}

// SetTrackedReq 設定追蹤優先等級請求
type SetTrackedReq struct {
	TrackPriority model.TrackPriority `json:"track_priority" binding:"min=0,max=2"` // 0=不追蹤 1=優先 2=次要
}
