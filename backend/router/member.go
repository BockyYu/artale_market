package router

import (
	"artale_market/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

// registerScraper 爬蟲專用路由（無需 token）
func registerScraper(g *gin.RouterGroup, d *Deps) {
	g.GET("/items/tracked", d.Item.GetTracked)
	g.POST("/items/:id/prices", d.Price.RecordPrice)
	g.GET("/items", d.Item.GetAll)
	g.POST("/items", d.Item.Create)
	g.PUT("/items/:id", d.Item.Update)
	g.PATCH("/items/:id/track", d.Item.SetTracked)
	g.DELETE("/items/:id", d.Item.Delete)
}

func registerMemberV1(g *gin.RouterGroup, d *Deps) {
	g.GET("/system", d.System.GetStatus)

	// 公開（不需 token）
	g.POST("/member/login", d.Member.Login)
	g.POST("/member/register", d.Member.Register)
	g.POST("/member/logout", d.Member.Logout)

	auth := g.Group("/member")
	if os.Getenv("APP_MODE") == "prod" {
		auth.Use(middleware.MemberJWTAuth())
	}
	{
		auth.GET("/me", d.Member.Me)
		auth.GET("/items", d.Item.GetAll)
		auth.GET("/items/:id/prices", d.Item.GetByID)
		auth.POST("/scrolls/search", d.Price.GetScrollSummary)
		auth.POST("/skillbooks/search", d.Price.GetSkillBookSummary)
	}
}
