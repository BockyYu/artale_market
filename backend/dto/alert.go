package dto

type CreateAlertReq struct {
	ItemID         uint    `json:"item_id"         binding:"required"`
	ThresholdPrice float64 `json:"threshold_price" binding:"required,gt=0"`
	BotID          *uint   `json:"bot_id"`
	Note           string  `json:"note"`
}
