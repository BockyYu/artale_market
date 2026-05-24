package handler

import (
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
		body.Date = time.Now().Format("2006-01-02")
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
		body.Date = time.Now().Format("2006-01-02")
	}
	if body.Page < 1 {
		body.Page = 1
	}

	result, err := h.svc.GetSkillBookSummary(body.Date, body.Category, body.SortBy, body.Page, body.PageSize)
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
		input.Date = time.Now().Format("2006-01-02")
	} else if _, err := time.Parse("2006-01-02", input.Date); err != nil {
		respBadRequest(c, err)
		return
	}

	itemID := parseID(c)
	record, err := h.svc.Record(itemID, input.Price, input.Date)
	if err != nil {
		respNotFound(c, err)
		return
	}

	if userID := c.GetHeader("X-User-ID"); userID != "" {
		go func() { _ = h.querySvc.RecordQuery(userID, itemID) }()
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
