package router

import (
	"artale_market/handler"
	"artale_market/middleware"
	"log"
	"os"
	"strings"

	"github.com/casbin/casbin/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Deps 聚合所有 handler，由 main.go 組裝後傳入。
type Deps struct {
	Item       *handler.ItemHandler
	Price      *handler.PriceHandler
	PriceHuma  *handler.PriceHumaHandler
	Query      *handler.QueryHandler
	Admin      *handler.AdminHandler
	Member     *handler.MemberHandler
	Permission *handler.PermissionHandler
	System     *handler.SystemHandler
	Alert      *handler.AlertHandler
	Bot        *handler.BotHandler
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
	origins := allowedOrigins()
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			for _, allowed := range origins {
				if strings.Contains(allowed, "*") {
					// 萬用字元：https://*.example.com 匹配 https://foo.example.com
					parts := strings.SplitN(allowed, "*", 2)
					if strings.HasPrefix(origin, parts[0]) && strings.HasSuffix(origin, parts[1]) {
						return true
					}
				} else if allowed == origin {
					return true
				}
			}
			return false
		},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "X-User-ID", "Authorization"},
	}))

	// ── Huma v2：Price handler 路由（OpenAPI 文件在 /api/docs）────────
	humaConfig := huma.DefaultConfig("Artale Market API", "1.0.0")
	humaConfig.OpenAPI.Servers = []*huma.Server{{URL: "/api", Description: "Artale Market"}}
	humaConfig.OpenAPI.Components = &huma.Components{
		Schemas: huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer),
		SecuritySchemes: map[string]*huma.SecurityScheme{
			"adminBearerAuth":  {Type: "http", Scheme: "bearer", BearerFormat: "JWT"},
			"memberBearerAuth": {Type: "http", Scheme: "bearer", BearerFormat: "JWT"},
		},
	}
	// 停用 Huma 內建 docs，改用下方自訂 gin endpoint 確保 spec URL 正確
	humaConfig.DocsPath = ""
	humaConfig.SchemasPath = ""

	// 公開 API：spec 掛在 /api 下（GET /api/openapi.yaml）
	publicHumaGroup := r.Group("/api")
	publicApi := humagin.NewWithGroup(r, publicHumaGroup, humaConfig)

	// 後續 API 共用同一份 OpenAPI spec（不重複開 spec 端點）
	// 三個 group 都掛在 /api 基底，middleware 各自加；
	// operation path 在 price_huma.go 裡帶完整前綴（/v1/member/... / /v1/admin/...），
	// 確保 OpenAPI spec 路徑與實際 gin 路由一致。
	humaConfig.OpenAPIPath = ""

	memberHumaGroup := r.Group("/api")
	if os.Getenv("APP_MODE") == "prod" {
		memberHumaGroup.Use(middleware.MemberJWTAuth())
	}
	memberApi := humagin.NewWithGroup(r, memberHumaGroup, humaConfig)

	adminHumaGroup := r.Group("/api")
	adminHumaGroup.Use(middleware.JWTAuth())
	adminApi := humagin.NewWithGroup(r, adminHumaGroup, humaConfig)

	// 自訂 docs 頁面：明確指定 spec URL 為 /api/openapi.yaml
	r.GET("/api/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Artale Market API Reference</title>
    <link rel="stylesheet" href="https://unpkg.com/@stoplight/elements@9.0.15/styles.min.css"
          crossorigin integrity="sha384-iVQBHadsD+eV0M5+ubRCEVXrXEBj+BqcuwjUwPoVJc0Pb1fmrhYSAhL+BFProHdV">
    <script src="https://unpkg.com/@stoplight/elements@9.0.15/web-components.min.js"
            crossorigin integrity="sha384-xjOcq9PZ/k+pGtPS/xcsCRXGjKKfTlIa4H1IYEnC+97jNa6sAMWTNrV6hY08W3GL"></script>
  </head>
  <body style="height:100vh;">
    <elements-api
      apiDescriptionUrl="/api/openapi.yaml"
      router="hash"
      layout="sidebar"
      tryItCredentialsPolicy="same-origin"
    ></elements-api>
  </body>
</html>`)
	})

	registerPriceHuma(publicApi, memberApi, adminApi, deps)

	// ── Gin 路由（其餘非 price 的 handler）───────────────────────────
	registerTools(r.Group("/api"), deps)
	registerMemberV1(r.Group("/api/v1"), deps)
	registerAdmin(r.Group("/api/v1/admin"), deps)

	return r
}
