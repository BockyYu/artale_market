package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func respOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func respCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, data)
}

func respDeleted(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func respBadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func respUnauthorized(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
}

func respNotFound(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
}

func respConflict(c *gin.Context, err error) {
	c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
}

func respInternal(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
