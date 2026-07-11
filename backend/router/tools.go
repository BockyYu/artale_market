package router

import "github.com/gin-gonic/gin"

// registerTools 自動／工具專用路由（無需 token）
// 注意：price 路由已遷移至 Huma v2（router/price_huma.go），此處不再重複註冊。
func registerTools(g *gin.RouterGroup, d *Deps) {
	g.GET("/items/tracked", d.Item.GetTracked)
	g.GET("/items", d.Item.GetAll)
	g.POST("/items", d.Item.Create)
	g.PUT("/items/:id", d.Item.Update)
	g.PATCH("/items/:id/track", d.Item.SetTracked)
	g.DELETE("/items/:id", d.Item.Delete)
	g.GET("/bot/alert-items", d.Alert.ListBotItems)
	g.GET("/bot/active", d.Bot.ListActiveBots)
	g.POST("/bot/notify", d.Bot.PublicNotify)
}
