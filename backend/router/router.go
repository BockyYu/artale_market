package router

import (
	"artale_market/handler"
	"log"
	"os"
	"strings"

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

func allowedOrigins() []string {
	v := os.Getenv("ALLOW_ORIGINS")
	if v == "" {
		log.Fatal("ALLOW_ORIGINS is not set")
	}
	var origins []string
	for _, p := range strings.Split(v, ",") {
		if s := strings.TrimSpace(p); s != "" {
			origins = append(origins, s)
		}
	}
	if len(origins) == 0 {
		log.Fatal("ALLOW_ORIGINS is empty")
	}
	return origins
}

func Setup(deps *Deps) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins(),
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "X-User-ID", "Authorization"},
	}))

	registerTools(r.Group("/api"), deps)
	registerMemberV1(r.Group("/api/v1"), deps)
	registerAdmin(r.Group("/api/v1/admin"), deps)

	return r
}
