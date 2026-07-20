package handler

import (
	"strconv"
	"time"

	"artale_market/dto"
	"artale_market/service"

	"github.com/gin-gonic/gin"
)

type PriceHandler struct {
	svc      service.PriceService
	querySvc service.QueryService
}

func NewPriceHandler(svc service.PriceService, querySvc service.QueryService) *PriceHandler {
	return &PriceHandler{svc: svc, querySvc: querySvc}
}

func (h *PriceHandler) GetScrollSummary(c *gin.Context) {
	var body dto.ScrollSearchReq
	if err := c.ShouldBindJSON(&body); err != nil {
		respBadRequest(c, err)
		return
	}
	if body.Date == "" {
		body.Date = twToday()
	}
	if body.Page < 1 {
		body.Page = 1
	}
	if len(body.Category) == 1 && body.Category[0] == "scroll_all" {
		body.Category = nil
	}

	result, err := h.svc.GetScrollSummary(body.Date, body.Percentage, body.Category, body.SortBy, body.Page, body.PageSize)
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, result)
}

func (h *PriceHandler) GetSkillBookSummary(c *gin.Context) {
	var body dto.SkillBookSearchReq
	if err := c.ShouldBindJSON(&body); err != nil {
		respBadRequest(c, err)
		return
	}
	if body.Date == "" {
		body.Date = twToday()
	}
	if body.Page < 1 {
		body.Page = 1
	}
	if body.SortBy == "" {
		body.SortBy = "price_desc"
	}

	result, err := h.svc.GetSkillBookSummary(body.Date, body.Category, body.SortBy, body.Page, body.PageSize)
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, result)
}

func (h *PriceHandler) GetEquipSummary(c *gin.Context) {
	var body dto.EquipSearchReq
	if err := c.ShouldBindJSON(&body); err != nil {
		respBadRequest(c, err)
		return
	}
	if body.Date == "" {
		body.Date = twToday()
	}
	if body.Page < 1 {
		body.Page = 1
	}
	if body.SortBy == "" {
		body.SortBy = "price_desc"
	}

	result, err := h.svc.GetEquipSummary(body.Date, body.Category, body.SortBy, body.Page, body.PageSize)
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, result)
}

func (h *PriceHandler) RecordPrice(c *gin.Context) {
	var input dto.RecordPriceReq
	if err := c.ShouldBindJSON(&input); err != nil {
		respBadRequest(c, err)
		return
	}

	if input.Date == "" {
		input.Date = twToday()
	} else if _, err := time.Parse("2006-01-02", input.Date); err != nil {
		respBadRequest(c, err)
		return
	}

	source := "scraper"
	if c.GetHeader("Authorization") != "" {
		source = "admin"
	}

	itemID := parseID(c)
	record, err := h.svc.Record(itemID, input.Price, input.Date, source)
	if err != nil {
		respNotFound(c, err)
		return
	}

	if userID := c.GetHeader("X-User-ID"); userID != "" {
		go func() { _ = h.querySvc.RecordQuery(userID, itemID) }()
	}

	respOK(c, record)
}

func (h *PriceHandler) GetLatestBatch(c *gin.Context) {
	var req struct {
		ItemIDs []uint `json:"item_ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respBadRequest(c, err)
		return
	}
	records, err := h.svc.GetLatestBatch(req.ItemIDs)
	if err != nil {
		respInternal(c, err)
		return
	}
	result := make(map[uint]float64, len(records))
	for _, r := range records {
		result[r.ItemID] = r.Price
	}
	respOK(c, gin.H{"data": result})
}

// GetMemberHistory GET /member/items/:id/price-history?days=7
// 回傳最近 N 天（預設 7，最多 30）的每日最低價記錄。
func (h *PriceHandler) GetMemberHistory(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	if days < 1 {
		days = 7
	}
	if days > 30 {
		days = 30
	}
	records, err := h.svc.GetHistory(parseID(c))
	if err != nil {
		respNotFound(c, err)
		return
	}
	if len(records) > days {
		records = records[:days]
	}
	respOK(c, gin.H{"data": records})
}

func (h *PriceHandler) GetLatest(c *gin.Context) {
	record, err := h.svc.GetLatest(parseID(c))
	if err != nil {
		respNotFound(c, err)
		return
	}
	respOK(c, record)
}

func (h *PriceHandler) GetHistory(c *gin.Context) {
	records, err := h.svc.GetHistory(parseID(c))
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, records)
}

func (h *PriceHandler) AdminGetHistory(c *gin.Context) {
	records, err := h.svc.GetAllHistory(parseID(c))
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, records)
}

func (h *PriceHandler) AdminGetPriceHistories(c *gin.Context) {
	records, err := h.svc.GetPriceHistories(parseID(c))
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, records)
}

func (h *PriceHandler) DeletePriceHistory(c *gin.Context) {
	if err := h.svc.DeletePriceHistory(parseID(c)); err != nil {
		respInternal(c, err)
		return
	}
	respDeleted(c)
}

func (h *PriceHandler) TogglePriceHistoryHidden(c *gin.Context) {
	var req struct {
		IsHidden bool `json:"is_hidden"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respBadRequest(c, err)
		return
	}
	if err := h.svc.TogglePriceHistoryHidden(parseID(c), req.IsHidden); err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, gin.H{"message": "updated"})
}
