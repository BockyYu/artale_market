package handler

import (
	"artale_market/repository"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	repo repository.SystemRepository
}

func NewSystemHandler(repo repository.SystemRepository) *SystemHandler {
	return &SystemHandler{repo}
}

func (h *SystemHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, h.buildStatus())
}

func (h *SystemHandler) buildStatus() gin.H {
	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = "test"
	}

	maintenance := false
	setting, err := h.repo.FindByName("maintenance")
	if err == nil {
		maintenance = setting.Status
	}

	return gin.H{
		"mode":        mode,
		"maintenance": maintenance,
		"message":     os.Getenv("APP_MESSAGE"),
	}
}
