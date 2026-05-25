package config

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"time"

	"artale_market/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")

	var (
		db  *gorm.DB
		err error
	)

	sqlLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	for i := 1; i <= 15; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: sqlLogger})
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
		&model.SystemSetting{},
	); err != nil {
		log.Fatalf("[DB] migration failed: %v", err)
	}

	seedDefaultAdmin(db)
	seedSupMember(db)
	seedSystemSettings(db)

	log.Println("[DB] connected and migrated successfully")
	return db
}

func md5Hash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func seedSupMember(db *gorm.DB) {
	var count int64
	db.Model(&model.Member{}).Where("username = ?", "sup_member").Count(&count)
	if count > 0 {
		return
	}
	member := model.Member{
		Nickname: "sup_member",
		Username: "sup_member",
		Password: md5Hash("sup_member"),
		Email:    "sup_member@artale.dev",
		Status:   1,
	}
	if err := db.Create(&member).Error; err != nil {
		log.Printf("[DB] failed to seed sup_member: %v", err)
		return
	}
	log.Println("[DB] sup_member created — username: sup_member  password: sup_member")
}

func seedSystemSettings(db *gorm.DB) {
	var count int64
	db.Model(&model.SystemSetting{}).Where("name = ?", "maintenance").Count(&count)
	if count > 0 {
		return
	}
	db.Create(&model.SystemSetting{Name: "maintenance", Status: false, OperatorName: "system"})
	log.Println("[DB] system setting 'maintenance' initialized")
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
