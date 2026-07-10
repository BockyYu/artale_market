package router

import "github.com/gin-gonic/gin"

// registerTools 自動／工具專用路由（無需 token）
func registerTools(g *gin.RouterGroup, d *Deps) {
	g.GET("/items/tracked", d.Item.GetTracked)
	g.POST("/items/:id/prices", d.Price.RecordPrice)
	g.GET("/items/:id/prices/latest", d.Price.GetLatest)
	g.POST("/items/prices/latest-batch", d.Price.GetLatestBatch)
	g.GET("/items", d.Item.GetAll)
	g.POST("/items", d.Item.Create)
	g.PUT("/items/:id", d.Item.Update)
	g.PATCH("/items/:id/track", d.Item.SetTracked)
	g.DELETE("/items/:id", d.Item.Delete)
	g.GET("/bot/alert-items", d.Alert.ListBotItems)
	g.GET("/bot/active", d.Bot.ListActiveBots)
	g.POST("/bot/notify", d.Bot.PublicNotify)
}
