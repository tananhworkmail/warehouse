package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

// HumidityAlertRouter — đăng ký các route liên quan đến cảnh báo độ ẩm + Andon
// Gọi hàm này trong WarehouseRouter hoặc Register()
func HumidityAlertRouter(router *gin.RouterGroup) {
	humidity := router.Group("/humidity")
	{
		// Trạng thái realtime độ ẩm 3 kho
		humidity.GET("/status", controllers.HumidityAlertCtrl.GetLatestHumidityStatus)

		// Lịch sử cảnh báo
		humidity.GET("/log", controllers.HumidityAlertCtrl.GetAlertLog)

		// Điều khiển Andon
		humidity.POST("/andon/off", controllers.HumidityAlertCtrl.TurnOffAndon)
		// humidity.POST("/andon/on", controllers.HumidityAlertCtrl.TurnOnAndon)
		humidity.GET("/andon/status", controllers.HumidityAlertCtrl.GetAndonStatus)
	}
}