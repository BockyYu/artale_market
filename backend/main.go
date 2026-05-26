package main

import (
	"log"

	"artale_market/config"
	"artale_market/handler"
	"artale_market/repository"
	"artale_market/router"
	"artale_market/service"

	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	if err := godotenv.Load(fmt.Sprintf(".env.%s", env)); err != nil {
		log.Printf("[Config] no .env.%s file found, using system environment", env)
	}
	log.Printf("[Config] running in %s mode", env)
	if env == "dev" {
		logConnInfo()
	}
	db  := config.NewDB()
	rdb := config.NewRedis()


	enforcer := config.NewEnforcer(db)

	itemRepo        := repository.NewItemRepository(db)
	priceRepo       := repository.NewPriceRepository(db)
	priceHistoryRepo := repository.NewPriceHistoryRepository(db)
	queryRepo       := repository.NewQueryRepository(rdb)
	adminRepo       := repository.NewAdminRepository(db)
	memberRepo      := repository.NewMemberRepository(db)
	systemRepo      := repository.NewSystemRepository(db)

	itemSvc   := service.NewItemService(itemRepo, priceRepo, queryRepo)
	priceSvc  := service.NewPriceService(itemRepo, priceRepo, priceHistoryRepo)
	querySvc  := service.NewQueryService(queryRepo, itemRepo)
	adminSvc  := service.NewAdminService(adminRepo)
	memberSvc := service.NewMemberService(memberRepo)

	deps := &router.Deps{
		Item:       handler.NewItemHandler(itemSvc, queryRepo),
		Price:      handler.NewPriceHandler(priceSvc, querySvc),
		Query:      handler.NewQueryHandler(querySvc),
		Admin:      handler.NewAdminHandler(adminSvc),
		Member:     handler.NewMemberHandler(memberSvc),
		Permission: handler.NewPermissionHandler(enforcer, adminSvc),
		System:     handler.NewSystemHandler(systemRepo),
		Enforcer:   enforcer,
	}

	r := router.Setup(deps)

	log.Println("[Server] running on :8080")
	r.Run(":8080")
}

func logConnInfo() {
	dsn := os.Getenv("DATABASE_URL")
	parts := strings.Fields(dsn)
	safe := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.HasPrefix(p, "password=") {
			safe = append(safe, "password=****")
		} else {
			safe = append(safe, p)
		}
	}
	log.Printf("[Config] DB  → %s", strings.Join(safe, " "))
	log.Printf("[Config] Redis → %s", os.Getenv("REDIS_URL"))
}
