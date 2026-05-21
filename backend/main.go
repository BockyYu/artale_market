package main

import (
	"log"

	"artale_market/config"
	"artale_market/handler"
	"artale_market/repository"
	"artale_market/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 等待 PostgreSQL 和 Redis 連線成功後才繼續
	db  := config.NewDB()
	rdb := config.NewRedis()

	// 2. 自動匯入預設卷軸種子資料（已存在則跳過）
	config.SeedScrolls(db)

	// 3. 組裝依賴（Repository → Service → Handler）
	itemRepo  := repository.NewItemRepository(db)
	priceRepo := repository.NewPriceRepository(db)
	queryRepo := repository.NewQueryRepository(rdb)

	itemSvc  := service.NewItemService(itemRepo, priceRepo)
	priceSvc := service.NewPriceService(itemRepo, priceRepo)
	querySvc := service.NewQueryService(queryRepo, itemRepo)

	itemH  := handler.NewItemHandler(itemSvc)
	priceH := handler.NewPriceHandler(priceSvc, querySvc)
	queryH := handler.NewQueryHandler(querySvc)

	// 4. 設定路由並啟動
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "X-User-ID"},
	}))

	api := r.Group("/api")
	{
		api.GET("/items", itemH.GetAll)
		api.GET("/items/tracked", itemH.GetTracked)
		api.POST("/items", itemH.Create)
		api.PUT("/items/:id", itemH.Update)
		api.PATCH("/items/:id/track", itemH.SetTracked)
		api.DELETE("/items/:id", itemH.Delete)

		api.GET("/prices/summary", priceH.GetSummary)
		api.POST("/items/:id/prices", priceH.RecordPrice)
		api.GET("/items/:id/prices", priceH.GetHistory)

		api.GET("/me/frequent-items", queryH.GetFrequent)
	}

	log.Println("[Server] running on :8080")
	r.Run(":8080")
}
