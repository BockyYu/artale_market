package handler

import (
	"errors"

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
		respBadRequest(c, errors.New("missing X-User-ID header"))
		return
	}
	items, err := h.svc.GetFrequent(userID)
	if err != nil {
		respInternal(c, err)
		return
	}
	respOK(c, items)
}
