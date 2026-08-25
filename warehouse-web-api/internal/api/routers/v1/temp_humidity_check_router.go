package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func TempHumidityCheckRouter(router *gin.RouterGroup) {
	tempHumidity := router.Group("/temp-humidity")
	{
		tempHumidity.GET("/check", controllers.TempHumidityCheckCtrl.Check)
		tempHumidity.GET("/check-safe-humidity", controllers.TempHumidityCheckCtrl.CheckSafeHumidity)
	}
}
