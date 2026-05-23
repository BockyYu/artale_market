package config

import (
	"log"
	"os"
	"time"

	"artale_market/model"

	"golang.org/x/crypto/bcrypt"
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

	if err := db.AutoMigrate(
		&model.Item{},
		&model.PriceRecord{},
		&model.AdminUser{},
		&model.Member{},
	); err != nil {
		log.Fatalf("[DB] migration failed: %v", err)
	}

	seedDefaultAdmin(db)

	log.Println("[DB] connected and migrated successfully")
	return db
}

func seedDefaultAdmin(db *gorm.DB) {
	var count int64
	db.Model(&model.AdminUser{}).Count(&count)
	if count > 0 {
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("Admin1234"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("[DB] failed to hash default admin password: %v", err)
	}
	admin := model.AdminUser{
		Username: "admin",
		Password: string(hash),
		Role:     "superadmin",
	}
	if err := db.Create(&admin).Error; err != nil {
		log.Printf("[DB] failed to seed default admin: %v", err)
		return
	}
	log.Println("[DB] default admin created — username: admin  password: Admin1234")
}
