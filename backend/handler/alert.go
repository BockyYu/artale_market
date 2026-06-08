package handler

import (
	"artale_market/dto"
	"artale_market/service"
	"errors"

	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	svc service.AlertService
}

func NewAlertHandler(svc service.AlertService) *AlertHandler {
	return &AlertHandler{svc: svc}
}

func (h *AlertHandler) List(c *gin.Context) {
	alerts, err := h.svc.List()
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, gin.H{"data": alerts})
}

func (h *AlertHandler) Create(c *gin.Context) {
	var req dto.CreateAlertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respBadRequest(c, err)
		return
	}
	alert, err := h.svc.Create(req.ItemID, req.BotID, req.ThresholdPrice, req.Note)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateAlert) {
			respConflict(c, err)
			return
		}
		respInternal(c, err)
		return
	}
	respOK(c, alert)
}

func (h *AlertHandler) Update(c *gin.Context) {
	var req dto.UpdateAlertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respBadRequest(c, err)
		return
	}
	if err := h.svc.Update(parseID(c), req.BotID, req.ThresholdPrice, req.Note); err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, gin.H{"message": "updated"})
}

func (h *AlertHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(parseID(c)); err != nil {
		respInternal(c, err)
		return
	}
	respDeleted(c)
}

// ListBotItems 回傳所有啟用提醒的道具清單（bot 專用，不需驗證）
func (h *AlertHandler) ListBotItems(c *gin.Context) {
	alerts, err := h.svc.List()
	if err != nil {
		respInternal(c, err)
		return
	}

	// 以 item_id 為 key，記錄最低門檻（一個道具可能有多筆 alert）
	type botItem struct {
		ItemID         uint    `json:"item_id"`
		ItemName       string  `json:"item_name"`
		EnglishName    string  `json:"english_name"`
		SearchMode     int     `json:"search_mode"`
		ItemType       int     `json:"item_type"`
		ThresholdPrice float64 `json:"threshold_price"`
		BotID          *uint   `json:"bot_id"`
	}
	itemMap := map[uint]*botItem{}
	for _, a := range alerts {
		if !a.IsActive {
			continue
		}
		if existing, ok := itemMap[a.ItemID]; ok {
			if a.ThresholdPrice < existing.ThresholdPrice {
				existing.ThresholdPrice = a.ThresholdPrice
				existing.BotID = a.BotID
			}
		} else {
			itemMap[a.ItemID] = &botItem{
				ItemID:         a.ItemID,
				ItemName:       a.Item.Name,
				EnglishName:    a.Item.EnglishName,
				SearchMode:     a.Item.SearchMode,
				ItemType:       int(a.Item.ItemType),
				ThresholdPrice: a.ThresholdPrice,
				BotID:          a.BotID,
			}
		}
	}
	var items []botItem
	for _, v := range itemMap {
		items = append(items, *v)
	}
	if items == nil {
		items = []botItem{}
	}
	respOK(c, gin.H{"data": items})
}

func (h *AlertHandler) ToggleActive(c *gin.Context) {
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respBadRequest(c, err)
		return
	}
	if err := h.svc.ToggleActive(parseID(c), req.IsActive); err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, gin.H{"message": "updated"})
}
