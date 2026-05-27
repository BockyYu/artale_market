package handler

import (
	"strconv"
	"time"

	"artale_market/dto"
	"artale_market/model"
	"artale_market/repository"
	"artale_market/service"

	"github.com/gin-gonic/gin"
)

type ItemHandler struct {
	svc       service.ItemService
	queryRepo repository.QueryRepository
}

func NewItemHandler(svc service.ItemService, qr repository.QueryRepository) *ItemHandler {
	return &ItemHandler{svc: svc, queryRepo: qr}
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
	sortBy := c.Query("sort_by")
	search := c.Query("search")
	filterType, _ := strconv.Atoi(c.DefaultQuery("filter_type", "0"))
	filterPriority, _ := strconv.Atoi(c.DefaultQuery("filter_priority", "-1"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}

	result, err := h.svc.GetAllWithLatestPrice(sortBy, search, filterType, filterPriority, page, pageSize)
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, result)
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
	item, err := h.svc.Update(parseID(c), input.Name, input.EnglishName, input.SearchMode, input.ItemType, input.Percentage, input.Category, input.Description)
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
	id := parseID(c)
	summary, err := h.svc.GetPriceSummary(id)
	if err != nil {
		respNotFound(c, err)
		return
	}
	today := time.Now().Format("2006-01-02")
	go func() { _ = h.queryRepo.RecordItemView(id, today) }()
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
