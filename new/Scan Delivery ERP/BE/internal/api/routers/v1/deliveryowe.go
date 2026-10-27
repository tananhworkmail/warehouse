package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func DeliveryOweRouter(router *gin.RouterGroup) {
	router.POST("/quantity", controllers.Do.ScanDeliveryOweQuantity)
	router.POST("/kcllscan_rfss", controllers.Do.DeliveryOweKCLLScan_RFSS)
	router.POST("/owe", controllers.Do.ScanDeliveryOwe)
}
