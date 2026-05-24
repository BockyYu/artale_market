package config

import (
	"log"

	"artale_market/model"

	"gorm.io/gorm"
)

type seedScroll struct {
	Name          string
	Percentage    int
	Category      string
	Description   string
	ItemType      model.ItemType
	TrackPriority model.TrackPriority
}

// SeedScrolls 在專案啟動時自動執行：
//   - 新資料：直接 insert
//   - 已存在：依 seed 資料同步 track_priority
func SeedScrolls(db *gorm.DB) {
	all := make([]seedScroll, 0, len(scrollSeeds)+len(miscSeeds)+len(skillBookSeeds))
	all = append(all, scrollSeeds...)
	all = append(all, miscSeeds...)
	all = append(all, skillBookSeeds...)

	inserted, patched := 0, 0
	for _, s := range all {
		var existing model.Item
		err := db.Where("name = ? AND percentage = ?", s.Name, s.Percentage).First(&existing).Error
		if err != nil {
			item := model.Item{
				Name:          s.Name,
				Percentage:    s.Percentage,
				Category:      s.Category,
				Description:   s.Description,
				ItemType:      s.ItemType,
				TrackPriority: s.TrackPriority,
			}
			if err := db.Create(&item).Error; err != nil {
				log.Printf("[Seed] failed to insert %s: %v", s.Name, err)
			} else {
				inserted++
			}
			continue
		}

		// 已存在，同步 track_priority 與 item_type
		updates := map[string]any{}
		if existing.TrackPriority != s.TrackPriority {
			updates["track_priority"] = s.TrackPriority
		}
		if existing.ItemType != s.ItemType {
			updates["item_type"] = s.ItemType
		}
		if len(updates) > 0 {
			if err := db.Model(&existing).Updates(updates).Error; err != nil {
				log.Printf("[Seed] failed to patch %s: %v", s.Name, err)
			} else {
				patched++
			}
		}
	}

	if inserted > 0 {
		log.Printf("[Seed] inserted %d item(s)", inserted)
	}
	if patched > 0 {
		log.Printf("[Seed] track_priority patched for %d item(s)", patched)
	}
	if inserted == 0 && patched == 0 {
		log.Println("[Seed] nothing to update")
	}
}
