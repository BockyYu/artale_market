package main

import (
	"log"

	"artale_market/config"
	"artale_market/handler"
	"artale_market/repository"
	"artale_market/router"
	"artale_market/service"
)

func main() {
	db  := config.NewDB()
	rdb := config.NewRedis()

	config.SeedScrolls(db)

	enforcer := config.NewEnforcer(db)

	itemRepo   := repository.NewItemRepository(db)
	priceRepo  := repository.NewPriceRepository(db)
	queryRepo  := repository.NewQueryRepository(rdb)
	adminRepo  := repository.NewAdminRepository(db)
	memberRepo := repository.NewMemberRepository(db)

	itemSvc   := service.NewItemService(itemRepo, priceRepo)
	priceSvc  := service.NewPriceService(itemRepo, priceRepo)
	querySvc  := service.NewQueryService(queryRepo, itemRepo)
	adminSvc  := service.NewAdminService(adminRepo)
	memberSvc := service.NewMemberService(memberRepo)

	deps := &router.Deps{
		Item:       handler.NewItemHandler(itemSvc),
		Price:      handler.NewPriceHandler(priceSvc, querySvc),
		Query:      handler.NewQueryHandler(querySvc),
		Admin:      handler.NewAdminHandler(adminSvc),
		Member:     handler.NewMemberHandler(memberSvc),
		Permission: handler.NewPermissionHandler(enforcer, adminSvc),
		Enforcer:   enforcer,
	}

	r := router.Setup(deps)

	log.Println("[Server] running on :8080")
	r.Run(":8080")
}
