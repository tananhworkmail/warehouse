package router_v1

import (
	"web-api/internal/api/controllers"
	"web-api/internal/api/middlewares"

	"github.com/gin-gonic/gin"
)

func WarehouseKTPRouter(router *gin.RouterGroup) {
	router.GET("/racks", controllers.WarehouseKTP.GetRacks)
	router.GET("/racks/search", controllers.WarehouseKTP.SearchRacks)
	router.GET("/search", controllers.WarehouseKTP.SearchInventory)
	router.GET("/inventory-age", controllers.WarehouseKTP.GetInventoryAge)
	router.GET("/racks/:rackCode/order-detail", controllers.WarehouseKTP.GetRackOrderDetail)
	router.GET("/temphumiditybydevices", controllers.WarehouseKTP.GetTempHumidityByDevices)
	router.GET("/temphumiditylatestlogbydevices", controllers.WarehouseKTP.GetLatestTempHumidityLogByDevices)
	router.GET("/orders/by-carton", controllers.WarehouseKTP.GetMoveOrderByCarton)
	scanIn := router.Group("/scan-in")
	scanIn.Use(middlewares.RequireJWT(), middlewares.RequireERPAnyFormPermission("VDH", "21", "SN7", "SN3"))
	{
		scanIn.GET("/access", controllers.WarehouseKTP.GetScanInAccess)
		scanIn.GET("/carton/:barcode", controllers.WarehouseKTP.GetScanInCarton)
		scanIn.GET("/warehouses", controllers.WarehouseKTP.GetScanInWarehouses)
		scanIn.GET("/production-lines", controllers.WarehouseKTP.GetScanInProductionLines)
		scanIn.POST("/carton", controllers.WarehouseKTP.SubmitScanInCarton)
	}
	scanOut := router.Group("/scan-out")
	scanOut.Use(middlewares.RequireJWT(), middlewares.RequireERPAnyFormPermissionForFeature("VDH", "21", "Scan Out KTP", "SN7", "SN3"))
	{
		scanOut.GET("/access", controllers.WarehouseKTP.GetScanOutAccess)
		scanOut.GET("/carton/:barcode", controllers.WarehouseKTP.GetScanOutCarton)
		scanOut.POST("/carton", controllers.WarehouseKTP.SubmitScanOutCarton)
	}
	scanShip := router.Group("/scan-ship")
	scanShip.Use(middlewares.RequireJWT(), middlewares.RequireERPAnyFormPermissionForFeature("VDH", "21", "Scan Ship KTP", "SN7", "SN3"))
	{
		scanShip.GET("/access", controllers.WarehouseKTP.GetScanShipAccess)
		scanShip.GET("/carton/:barcode", controllers.WarehouseKTP.GetScanShipCarton)
		scanShip.POST("/carton", controllers.WarehouseKTP.SubmitScanShipCartons)
	}
	router.POST("/orders/move", controllers.WarehouseKTP.MoveOrderToRack)
	router.POST("/racks/:rackCode/scan", controllers.WarehouseKTP.ScanRack)
	router.POST("/racks/:rackCode/clear", controllers.WarehouseKTP.ClearRack)
}
