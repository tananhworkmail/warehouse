package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func WarehouseKTPRouter(router *gin.RouterGroup) {
	router.GET("/racks", controllers.WarehouseKTP.GetRacks)
	router.GET("/racks/search", controllers.WarehouseKTP.SearchRacks)
	router.GET("/racks/:rackCode/order-detail", controllers.WarehouseKTP.GetRackOrderDetail)
	router.GET("/temphumiditybydevices", controllers.WarehouseKTP.GetTempHumidityByDevices)
	router.GET("/temphumiditylatestlogbydevices", controllers.WarehouseKTP.GetLatestTempHumidityLogByDevices)
	router.GET("/orders/by-carton", controllers.WarehouseKTP.GetMoveOrderByCarton)
	router.POST("/orders/move", controllers.WarehouseKTP.MoveOrderToRack)
	router.POST("/racks/:rackCode/scan", controllers.WarehouseKTP.ScanRack)
	router.POST("/racks/:rackCode/clear", controllers.WarehouseKTP.ClearRack)
}
