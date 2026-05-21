package handler

import (
	"net/http"
	"strconv"
	"strings"
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
	date := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	var pcts []int
	if pctStr := c.Query("percentage"); pctStr != "" {
		for _, s := range strings.Split(pctStr, ",") {
			if v, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				pcts = append(pcts, v)
			}
		}
	}

	var categories []string
	if catStr := c.Query("category"); catStr != "" {
		for _, s := range strings.Split(catStr, ",") {
			if t := strings.TrimSpace(s); t != "" {
				categories = append(categories, t)
			}
		}
	}

	summaries, err := h.svc.GetSummary(date, pcts, categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summaries)
}

func (h *PriceHandler) RecordPrice(c *gin.Context) {
	var input struct {
		Price float64 `json:"price" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemID := parseID(c)
	record, err := h.svc.Record(itemID, input.Price)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	// 非同步記錄查詢紀錄，不影響主流程
	if userID := c.GetHeader("X-User-ID"); userID != "" {
		go func() { _ = h.querySvc.RecordQuery(userID, itemID) }()
	}

	c.JSON(http.StatusOK, record)
}

func (h *PriceHandler) GetHistory(c *gin.Context) {
	records, err := h.svc.GetHistory(parseID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}
