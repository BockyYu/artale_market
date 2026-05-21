package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseID(c *gin.Context) uint {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(id)
}
