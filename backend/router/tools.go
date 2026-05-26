package router

import "github.com/gin-gonic/gin"

// registerTools 自動／工具專用路由（無需 token）
func registerTools(g *gin.RouterGroup, d *Deps) {
	g.GET("/items/tracked", d.Item.GetTracked)
	g.POST("/items/:id/prices", d.Price.RecordPrice)
	g.GET("/items", d.Item.GetAll)
	g.POST("/items", d.Item.Create)
	g.PUT("/items/:id", d.Item.Update)
	g.PATCH("/items/:id/track", d.Item.SetTracked)
	g.DELETE("/items/:id", d.Item.Delete)
}
