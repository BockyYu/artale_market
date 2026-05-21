package config

import (
	"log"
	"os"
	"time"

	"artale_market/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")

	var (
		db  *gorm.DB
		err error
	)

	for i := 1; i <= 15; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, e := db.DB()
			if e == nil {
				e = sqlDB.Ping()
			}
			if e == nil {
				break
			}
			err = e
		}
		log.Printf("[DB] attempt %d/15 failed: %v — retrying in 2s", i, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("[DB] could not connect after 15 attempts: %v", err)
	}

	if err := db.AutoMigrate(&model.Item{}, &model.PriceRecord{}); err != nil {
		log.Fatalf("[DB] migration failed: %v", err)
	}

	log.Println("[DB] connected and migrated successfully")
	return db
}
