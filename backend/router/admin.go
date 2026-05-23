package router

import (
	"artale_market/middleware"

	"github.com/gin-gonic/gin"
)

func registerAdmin(g *gin.RouterGroup, d *Deps) {
	// 後台登入（公開）
	g.POST("/login", d.Admin.Login)

	// JWT 保護
	auth := g.Group("/")
	auth.Use(middleware.JWTAuth())
	{
		// 管理員帳號管理（需要 admin:manage 權限）
		auth.GET("/admins", d.Admin.List)
		auth.POST("/admins", middleware.CasbinAuth(d.Enforcer, "admin", "manage"), d.Admin.Create)
		auth.PUT("/admins/:id", middleware.CasbinAuth(d.Enforcer, "admin", "manage"), d.Admin.Update)
		auth.DELETE("/admins/:id", middleware.CasbinAuth(d.Enforcer, "admin", "manage"), d.Admin.Delete)

		// 權限管理
		auth.GET("/admins/:id/permissions", d.Permission.Get)
		auth.PUT("/admins/:id/permissions", d.Permission.Update)

		// 會員管理
		auth.GET("/members", d.Member.List)
		auth.PUT("/members/:id/status", d.Member.UpdateStatus)
		auth.DELETE("/members/:id", d.Member.Delete)

		// 後台新增價格（需要 price:write 權限）
		auth.POST("/items/:id/prices", middleware.CasbinAuth(d.Enforcer, "price", "write"), d.Price.RecordPrice)
	}
}
