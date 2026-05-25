package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetAppConfig(c *gin.Context) {
	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = "public"
	}

	message := os.Getenv("APP_MESSAGE")

	c.JSON(http.StatusOK, gin.H{
		"mode":    mode,
		"message": message,
	})
}
