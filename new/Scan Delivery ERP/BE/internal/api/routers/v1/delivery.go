package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func DeliveryRouter(router *gin.RouterGroup) {

	// ScanQRDelivery1
	//LLbarCodeEditKeyPress (phiếu xuất)
	router.POST("/KCLL", controllers.Dl.GetScanDelivery)
	//LLbarCodeEdit1KeyPress (scan nhãn) + ScanQRDeliveryOWE1 LLBarCodeEditKeyPress
	router.POST("/scanlabel", controllers.Dl.GetScanDeliveryQuantity)

	router.POST("/scanlabelclbh", controllers.Dl.GetScanDeliveryCLBH)

}
