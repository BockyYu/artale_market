package router

import "github.com/gin-gonic/gin"

func registerMember(g *gin.RouterGroup, d *Deps) {
	// 前台登入
	g.POST("/member/login", d.Member.Login)

	g.GET("/items", d.Item.GetAll)
	g.GET("/items/tracked", d.Item.GetTracked)
	g.POST("/items", d.Item.Create)
	g.PUT("/items/:id", d.Item.Update)
	g.PATCH("/items/:id/track", d.Item.SetTracked)
	g.DELETE("/items/:id", d.Item.Delete)

	g.GET("/prices/summary", d.Price.GetSummary)
	g.POST("/items/:id/prices", d.Price.RecordPrice)
	g.GET("/items/:id/prices", d.Price.GetHistory)

	g.GET("/me/frequent-items", d.Query.GetFrequent)
}
