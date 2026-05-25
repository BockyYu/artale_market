package handler

import (
	"time"

	"artale_market/dto"
	"artale_market/model"
	"artale_market/service"

	"github.com/gin-gonic/gin"
)

type ItemHandler struct {
	svc service.ItemService
}

func NewItemHandler(svc service.ItemService) *ItemHandler {
	return &ItemHandler{svc: svc}
}

func (h *ItemHandler) GetAll(c *gin.Context) {
	items, err := h.svc.GetAll()
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, items)
}

func (h *ItemHandler) AdminGetAll(c *gin.Context) {
	items, err := h.svc.GetAllWithLatestPrice(c.Query("sort_by"))
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, items)
}

func (h *ItemHandler) Create(c *gin.Context) {
	var item model.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		respBadRequest(c, err)
		return
	}
	if err := h.svc.Create(&item); err != nil {
		respInternal(c, err)
		return
	}
	respCreated(c, item)
}

func (h *ItemHandler) Update(c *gin.Context) {
	var input dto.UpdateItemReq
	if err := c.ShouldBindJSON(&input); err != nil {
		respBadRequest(c, err)
		return
	}
	item, err := h.svc.Update(parseID(c), input.Name, input.Percentage, input.Category, input.Description)
	if err != nil {
		respNotFound(c, err)
		return
	}
	respOK(c, item)
}

func (h *ItemHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(parseID(c)); err != nil {
		respNotFound(c, err)
		return
	}
	respDeleted(c)
}

func (h *ItemHandler) GetByID(c *gin.Context) {
	summary, err := h.svc.GetPriceSummary(parseID(c))
	if err != nil {
		respNotFound(c, err)
		return
	}
	respOK(c, summary)
}

func (h *ItemHandler) GetTracked(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	items, err := h.svc.GetTracked(today)
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, items)
}

func (h *ItemHandler) SetTracked(c *gin.Context) {
	var input dto.SetTrackedReq
	if err := c.ShouldBindJSON(&input); err != nil {
		respBadRequest(c, err)
		return
	}
	item, err := h.svc.SetTracked(parseID(c), input.TrackPriority)
	if err != nil {
		respNotFound(c, err)
		return
	}
	respOK(c, item)
}
