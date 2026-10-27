package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func DeliveryPayRouter(router *gin.RouterGroup) {
	router.POST("/export", controllers.Dp.ScanDeliveryPayExPort)

	router.POST("/getkclls", controllers.Dp.GetScanDeliveryPay)
	router.POST("/getkcllScan_rfss", controllers.Dp.GetScanDeliveryPayKCLLScan_RFSS)
	// DBGridEh2DblClick
	router.POST("/dbclick", controllers.Dp.DBGridEh2DblClick)

	router.POST("/kcllscan_rfsss", controllers.Dp.DeliveryPayKCLLScan_RFSSS)

}
