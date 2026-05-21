package handler

import (
	"net/http"

	"artale_market/service"

	"github.com/gin-gonic/gin"
)

type QueryHandler struct {
	svc service.QueryService
}

func NewQueryHandler(svc service.QueryService) *QueryHandler {
	return &QueryHandler{svc: svc}
}

func (h *QueryHandler) GetFrequent(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing X-User-ID header"})
		return
	}
	items, err := h.svc.GetFrequent(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
