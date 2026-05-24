package router

import "github.com/gin-gonic/gin"

func registerMember(g *gin.RouterGroup, d *Deps) {
	// 前台登入
	g.POST("/member/login", d.Member.Login)

	g.GET("/items", d.Item.GetAll)
	g.GET("/member/items", d.Item.GetAll)
	g.GET("/items/tracked", d.Item.GetTracked)
	g.POST("/items", d.Item.Create)
	g.PUT("/items/:id", d.Item.Update)
	g.PATCH("/items/:id/track", d.Item.SetTracked)
	g.DELETE("/items/:id", d.Item.Delete)

	g.POST("/v1/scrolls/search", d.Price.GetScrollSummary)
	g.POST("/v1/skillbooks/search", d.Price.GetSkillBookSummary)
	g.POST("/items/:id/prices", d.Price.RecordPrice)
	g.GET("/items/:id/prices", d.Item.GetByID)

	// TODO: 常用商品功能，暫不開放
	// g.GET("/me/frequent-items", d.Query.GetFrequent)
}
