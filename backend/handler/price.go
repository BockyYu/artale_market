package handler

import (
	"time"

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

func (h *PriceHandler) GetSummary(c *gin.Context) {
	var body struct {
		Date       string   `json:"date"`
		Percentage []int    `json:"percentage"`
		Category   []string `json:"category"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respBadRequest(c, err)
		return
	}
	if body.Date == "" {
		body.Date = time.Now().Format("2006-01-02")
	}

	summaries, err := h.svc.GetSummary(body.Date, body.Percentage, body.Category)
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, summaries)
}

func (h *PriceHandler) RecordPrice(c *gin.Context) {
	var input struct {
		Price float64 `json:"price" binding:"required,gt=0"`
		Date  string  `json:"date"`
	}
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
