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
		auth.POST("/refresh", d.Admin.Refresh)
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

		// 道具列表與查詢優先度管理
		// 注意：price 路由已遷移至 Huma v2（router/price_huma.go），此處不再重複註冊。
		auth.GET("/export/excel", d.Item.ExportExcel)
		auth.POST("/export/discord", d.Item.SendExcelToDiscord)
		auth.GET("/items/categories", d.Item.GetCategories)
		auth.GET("/items", d.Item.AdminGetAll)
		auth.POST("/items", d.Item.Create)
		auth.PUT("/items/:id", d.Item.Update)
		auth.PATCH("/items/:id/track", d.Item.SetTracked)
		auth.PATCH("/items/:id/hidden", d.Item.SetHidden)

		// 價格提醒
		auth.GET("/alerts", d.Alert.List)
		auth.POST("/alerts", d.Alert.Create)
		auth.PUT("/alerts/:id", d.Alert.Update)
		auth.DELETE("/alerts/:id", d.Alert.Delete)
		auth.PATCH("/alerts/:id/active", d.Alert.ToggleActive)

		// Discord 測試
		auth.POST("/discord/test", d.Bot.TestDiscord)

		// 通知機器人
		auth.GET("/bots", d.Bot.List)
		auth.POST("/bots", d.Bot.Create)
		auth.PUT("/bots/:id", d.Bot.Update)
		auth.DELETE("/bots/:id", d.Bot.Delete)
		auth.PATCH("/bots/:id/active", d.Bot.ToggleActive)
		auth.POST("/bots/:id/send", d.Bot.SendMessage)
		auth.POST("/bots/:id/test", d.Bot.TestBot)
	}
}
