package config

import (
	"artale_market/repository"
	"artale_market/service"
	"log"
	"time"
)

// StartAlertScheduler 每 5 分鐘檢查一次所有啟用中的價格提醒
func StartAlertScheduler(alertRepo repository.AlertRepository, priceRepo repository.PriceRepository, alertSvc service.AlertService) {
	go func() {
		runAlertCheck(alertRepo, priceRepo, alertSvc)

		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			runAlertCheck(alertRepo, priceRepo, alertSvc)
		}
	}()
	log.Println("[AlertScheduler] started — interval: 5 minutes")
}

func runAlertCheck(alertRepo repository.AlertRepository, priceRepo repository.PriceRepository, alertSvc service.AlertService) {
	alerts, err := alertRepo.FindAllActive()
	if err != nil {
		log.Printf("[AlertScheduler] failed to fetch active alerts: %v", err)
		return
	}
	if len(alerts) == 0 {
		return
	}

	checked := make(map[uint]bool)
	for _, alert := range alerts {
		if checked[alert.ItemID] {
			continue
		}
		checked[alert.ItemID] = true

		record, err := priceRepo.FindLatestByItem(alert.ItemID)
		if err != nil || record == nil {
			continue
		}
		alertSvc.CheckAndNotify(alert.ItemID, alert.Item.Name, record.Price)
	}
	log.Printf("[AlertScheduler] checked %d item(s)", len(checked))
}
