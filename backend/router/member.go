package router

import (
	"artale_market/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func registerMemberV1(g *gin.RouterGroup, d *Deps) {
	g.GET("/system", d.System.GetStatus)

	// 公開（不需 token）
	g.POST("/member/login", d.Member.Login)
	g.POST("/member/register", d.Member.Register)
	g.POST("/member/logout", d.Member.Logout)

	// 注意：/scrolls/search、/skillbooks/search、/equips/search 已遷移至 Huma v2（router/price_huma.go）
	auth := g.Group("/member")
	if os.Getenv("APP_MODE") == "prod" {
		auth.Use(middleware.MemberJWTAuth())
	}
	{
		auth.GET("/me", d.Member.Me)
		auth.GET("/items", d.Item.GetAll)
		auth.GET("/items/:id/prices", d.Item.GetByID)
	}
}
