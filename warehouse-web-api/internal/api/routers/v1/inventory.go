// internal/api/routers/v1/inventory.go
package router_v1

import (
	"github.com/gin-gonic/gin"
	"web-api/internal/api/controllers"
)

func InventoryRouter(router *gin.RouterGroup) {
	// Đổi controllers.Dp thành controllers.NumberOfScans
	router.GET("/inventory", controllers.NumberOfScans.SelectInventory)
	router.POST("/inventory", controllers.NumberOfScans.CreateInventory)
	router.POST("/inventory-month", controllers.NumberOfScans.CreateInventoryMonth)
	router.GET("/inventory-month", controllers.NumberOfScans.SelectInventoryByMonth)
	router.GET("/rack-list", controllers.NumberOfScans.GetRackList)
	router.POST("/inventory-month/sync", controllers.NumberOfScans.SyncInventoryMonth)
	router.POST("/inventory-month/confirm-round2", controllers.NumberOfScans.ConfirmRound2)
	router.POST("/inventory-month/confirm-round2-scan", controllers.NumberOfScans.ConfirmRound2ByBarcode)
	router.GET("/inventory-month/last-sync", controllers.NumberOfScans.GetLastSync)
}