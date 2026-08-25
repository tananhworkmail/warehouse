package api

import (
	"fmt"
	"time"

	"web-api/internal/api/services"
)

const (
	tempHumidityRealCheckInterval = 30 * time.Minute
	tempHumiditySafeCheckInterval = 30 * time.Minute
)

func StartTempHumidityCheckJob() {
	fmt.Println("[TempHumidityCheckJob] Starting monitor (real=30 minute cycle, fake=30 minute cycle, inactive=no notify)...")

	for {
		time.Sleep(nextTempHumidityCheckDelay(time.Now()))
		if shouldRunTempHumidityCheckJob(time.Now()) {
			runTempHumidityCheckJob()
		}
	}
}

func shouldRunTempHumidityCheckJob(now time.Time) bool {
	return services.TempHumidityCheck.IsRealDataWindow(now) || services.TempHumidityCheck.IsFakeDataWindow(now)
}

func nextTempHumidityCheckDelay(now time.Time) time.Duration {
	realStart, realEnd := services.TempHumidityCheck.RealDataWindow(now)
	fakeStart, fakeEnd := services.TempHumidityCheck.FakeDataWindow(now)

	if services.TempHumidityCheck.IsRealDataWindow(now) {
		return nextTempHumiditySlotDelay(now, realEnd, tempHumidityRealCheckInterval)
	}

	if services.TempHumidityCheck.IsFakeDataWindow(now) {
		return nextTempHumiditySlotDelay(now, fakeEnd, tempHumiditySafeCheckInterval)
	}

	if now.Before(realStart) {
		return realStart.Sub(now)
	}
	if now.Before(fakeStart) {
		return fakeStart.Sub(now)
	}

	nextRealStart := realStart.Add(24 * time.Hour)
	return nextRealStart.Sub(now)
}

func nextTempHumiditySlotDelay(now time.Time, windowEnd time.Time, interval time.Duration) time.Duration {
	next := now.Truncate(interval)
	if next.Equal(now) {
		return 0
	}
	if next.Before(now) {
		next = next.Add(interval)
	}
	if next.After(windowEnd) {
		next = windowEnd
	}
	return next.Sub(now)
}

func runTempHumidityCheckJob() {
	alerts, err := services.TempHumidityCheck.CheckLatestAndNotify(nil)
	if err != nil {
		fmt.Println("[TempHumidityCheckJob] check failed:", err)
		return
	}

	if len(alerts) == 0 {
		return
	}

	fmt.Printf("[TempHumidityCheckJob] alert count=%d\n", len(alerts))
}
