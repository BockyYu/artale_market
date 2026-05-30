package handler

import (
	"artale_market/dto"
	"artale_market/service"

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
		respInternal(c, err)
		return
	}
	respOK(c, alert)
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

	seen := map[uint]bool{}
	type botItem struct {
		ItemID   uint   `json:"item_id"`
		ItemName string `json:"item_name"`
		ItemType int    `json:"item_type"`
	}
	var items []botItem
	for _, a := range alerts {
		if !a.IsActive || seen[a.ItemID] {
			continue
		}
		seen[a.ItemID] = true
		items = append(items, botItem{
			ItemID:   a.ItemID,
			ItemName: a.Item.Name,
			ItemType: int(a.Item.ItemType),
		})
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
