package api

import (
	"fmt"
	"time"
	"web-api/internal/api/services"
)

func StartHumidityAlertJob() {
	go func() {
		fmt.Println("[HumidityJob] Starting Monitor (30s cycle)...")

		services.HumidityAlert.InitBaseline()

		services.HumidityAlert.CheckAndTriggerAndon()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			services.HumidityAlert.CheckAndTriggerAndon()
		}
	}()
}