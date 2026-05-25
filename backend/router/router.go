package router

import (
	"artale_market/handler"

	"github.com/casbin/casbin/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Deps 聚合所有 handler，由 main.go 組裝後傳入。
type Deps struct {
	Item       *handler.ItemHandler
	Price      *handler.PriceHandler
	Query      *handler.QueryHandler
	Admin      *handler.AdminHandler
	Member     *handler.MemberHandler
	Permission *handler.PermissionHandler
	System     *handler.SystemHandler
	Enforcer   *casbin.Enforcer
}

func Setup(deps *Deps) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "X-User-ID", "Authorization"},
	}))

	registerScraper(r.Group("/api"), deps)
	registerMemberV1(r.Group("/api/v1"), deps)
	registerAdmin(r.Group("/api/admin"), deps)

	return r
}
