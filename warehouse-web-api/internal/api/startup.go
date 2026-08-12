package api

import (
	"fmt"

	router "web-api/internal/api/routers"
	"web-api/internal/pkg/config"
	"web-api/internal/pkg/database"
	"web-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Run(configPath string) {
	if configPath == "" {
		configPath = "data/config.yml"
	}
	if err := config.Setup(configPath); err != nil {
		logger.Fatalf("failed to setup config, %s", err)
	}
	if err := database.Setup(); err != nil {
		logger.Fatalf("failed to setup database, %s", err)
	}

	gin.SetMode(config.GetConfig().Server.Mode)

	cfg := config.GetConfig()

	web := router.Setup()

	// ── Background job kiểm tra độ ẩm ──────────────────────────────────────
	// BUG FIX: Gọi InitBaseline TRƯỚC khi start job để tránh Andon kêu
	// mỗi lần restart server (lastProcessedID=0 → luôn "new" → trigger)
	go StartWarehouseKTPRealtimeJob()
	go StartTempHumidityCheckJob()
	go StartHumidityAlertJob()

	fmt.Println("Web API Running on port " + cfg.Server.Port)
	fmt.Println("================================>")
	logger.Fatalf("%v", web.Run(":"+cfg.Server.Port))
}
