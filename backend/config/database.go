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
			LogLevel:                  logger.Warn,
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
		&model.PriceHistory{},
		&model.AdminUser{},
		&model.Member{},
		&model.SystemSetting{},
		&model.NotifyBot{},
		&model.PriceAlert{},
		&model.Category{},
	); err != nil {
		log.Fatalf("[DB] migration failed: %v", err)
	}

	seedDefaultAdmin(db)
	seedSupMember(db)
	seedSystemSettings(db)
	seedCategories(db)

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

func seedCategories(db *gorm.DB) {
	var count int64
	db.Model(&model.Category{}).Count(&count)
	if count > 0 {
		return
	}

	// 卷軸 (item_type=1) — 防具部位 + 武器種類
	scrollCats := []string{
		"頭盔", "上衣", "下衣", "套服", "鞋子", "手套", "披風", "盾牌",
		"臉部裝飾", "眼部裝飾", "耳環", "戒指", "墜飾", "腰帶", "肩章", "勳章",
		"單手劍", "雙手劍", "單手斧", "雙手斧", "單手棍", "雙手棍",
		"槍", "矛", "短杖", "長杖", "弓", "弩", "短劍", "拳套", "指虎", "火槍",
	}
	// 技能書 (item_type=4) — 職業
	skillBookCats := []string{
		"劍士", "英雄", "聖騎士", "黑騎士",
		"魔法師", "冰雷法師", "火毒法師", "主教",
		"弓手", "神射手", "箭神",
		"盜賊", "夜行者", "暗影俠盜",
		"海盜", "拳擊手", "槍手",
	}

	var rows []model.Category
	for _, name := range scrollCats {
		rows = append(rows, model.Category{Name: name, ItemType: 1})
	}
	for _, name := range skillBookCats {
		rows = append(rows, model.Category{Name: name, ItemType: 4})
	}

	if err := db.Create(&rows).Error; err != nil {
		log.Printf("[DB] failed to seed categories: %v", err)
		return
	}
	log.Printf("[DB] seeded %d categories", len(rows))
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
