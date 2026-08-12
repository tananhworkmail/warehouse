package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func WarehouseRouter(router *gin.RouterGroup) {
	router.POST("/getwarehouselist", controllers.Warehouse.GetWareHouseList)

	router.GET("/rack", controllers.Warehouse.GetRackInformation)
	router.GET("/3dayrack", controllers.Warehouse.Get3DayRackInformation)
	router.GET("/tonkho", controllers.Warehouse.GetTonKhoByRack)
	router.GET("/tonkhoke", controllers.Warehouse.GetTonKhoKeByRack)
	router.GET("/tonkhotang", controllers.Warehouse.GetTonKhoTangByRack)
	router.GET("/search", controllers.Warehouse.GetSearchByDDBH)
	router.GET("/searchtotal", controllers.Warehouse.GetSearchTotalDDBH)
	router.POST("/schedule", controllers.Warehouse.UploadSchedule)
	router.GET("/schedule", controllers.Warehouse.GetSchedule)
	router.GET("/check_order", controllers.Warehouse.GetOrderExistenceStatus)
	router.GET("bom_information", controllers.Warehouse.GetBomInformation)

	// KVT
	router.GET("/rackkvt", controllers.WarehouseKVT.GetRackKVT)
	router.GET("/totalcolumnkvt", controllers.WarehouseKVT.GetRackTotalColumnKVT)
	router.GET("/totaltangkvt", controllers.WarehouseKVT.GetTonKhoTangKVT)
	router.GET("/tonkhorackkvt", controllers.WarehouseKVT.GetTonKhoInRackKVT)
	router.GET("/rackinforkvt", controllers.WarehouseKVT.GetRackInformationKVT)
	router.GET("/3dayrackinforkvt", controllers.WarehouseKVT.Get3DayRackInformationKVT)
	router.GET("/3and180daykvt", controllers.WarehouseKVT.Get3And180DayKVT)
	router.GET("/searchkvt", controllers.WarehouseKVT.GetSearchByCLBHKVT)
	router.GET("/searchn31", controllers.WarehouseKVT.GetSearchN31)
	router.GET("/zonekvt/:zone", controllers.WarehouseKVT.GetRacksByZone)
	router.GET("/temphumidity", controllers.WarehouseKVT.GetTempHumidityController)
	router.GET("/temphumiditybydevices", controllers.WarehouseKVT.GetTempHumidityByDevicesController)
	router.GET("/temphumiditylatestlogbydevices", controllers.WarehouseKVT.GetLatestTempHumidityLogByDevicesController)
	router.GET("/humidityalert01", controllers.WarehouseKVT.GetHumidityAlertKVT01Controller)
	router.GET("/humidityalert02", controllers.WarehouseKVT.GetHumidityAlertKVT02Controller)
	router.GET("/humidityalert03", controllers.WarehouseKVT.GetHumidityAlertKVT03Controller)
	router.GET("/laboratory/environment", controllers.WarehouseKVT.GetTempHumidityLaboratoryController)
	// ✅ Andon + realtime humidity — đăng ký từ humidity_alert_router.go
	// Các route được tạo ra:
	//   GET  /api/v1/warehouse/humidity/status
	//   GET  /api/v1/warehouse/humidity/log
	//   POST /api/v1/warehouse/humidity/andon/off
	//   POST /api/v1/warehouse/humidity/andon/on
	//   GET  /api/v1/warehouse/humidity/andon/status
	HumidityAlertRouter(router)
}
