package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var twLoc = time.FixedZone("Asia/Taipei", 8*60*60)

func twToday() string {
	return time.Now().In(twLoc).Format("2006-01-02")
}

func parseID(c *gin.Context) uint {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(id)
}
