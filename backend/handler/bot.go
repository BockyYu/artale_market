package handler

import (
	"artale_market/dto"
	"artale_market/service"

	"github.com/gin-gonic/gin"
)

type BotHandler struct {
	svc service.BotService
}

func NewBotHandler(svc service.BotService) *BotHandler {
	return &BotHandler{svc: svc}
}

func (h *BotHandler) List(c *gin.Context) {
	bots, err := h.svc.List()
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, bots)
}

func (h *BotHandler) Create(c *gin.Context) {
	var req dto.CreateBotReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respBadRequest(c, err)
		return
	}
	bot, err := h.svc.Create(req.Name, req.Platform, req.Token, req.ChatID)
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, bot)
}

func (h *BotHandler) Update(c *gin.Context) {
	var req dto.UpdateBotReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respBadRequest(c, err)
		return
	}
	if err := h.svc.Update(parseID(c), req.Name, req.Platform, req.Token, req.ChatID); err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, gin.H{"message": "updated"})
}

func (h *BotHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(parseID(c)); err != nil {
		respInternal(c, err)
		return
	}
	respDeleted(c)
}

func (h *BotHandler) ToggleActive(c *gin.Context) {
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
