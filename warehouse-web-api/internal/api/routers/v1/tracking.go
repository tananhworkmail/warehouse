package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func TrackingRouter(router *gin.RouterGroup) {
	// ScanMaterial
	// scan nhập -	BitBtn2Click
	router.POST("/meterialimport", controllers.Tk.ScanMaterialImport)
	router.POST("/updatekcbhbybarcode", controllers.Tk.UpdateKCBHByBarcode)
	// scan xuất -BitBtn2Click
	router.POST("/meterialexport", controllers.Tk.ScanMaterialExport)
	// imgSave00Click
	router.POST("/getkcllscan_rfss", controllers.Tk.GetKCLLScan_RFSS)
	router.POST("/getkcllscan_rfsss", controllers.Tk.GetKCLLScan_RFSSS)

	router.POST("/merterialprint", controllers.Tk.ScanMaterialPrint)
	//lưu dữ liệu bị xóa -Save1Click
	router.POST("/updatekcllscan_rfss", controllers.Tk.UpdateKCLLScan_RFSS)
	router.POST("/updatekcrkscan_rfss", controllers.Tk.UpdateKCRKScan_RFSS)
	// Save1Click - 	// ScanQRDeliveryOWE1 => TScanQRDeliveryOWE.SaveClick
	router.POST("/updatekcllscan_rfsss", controllers.Tk.UpdateKCLLScan_RFSSS)
}
