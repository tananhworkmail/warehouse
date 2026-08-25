package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func LossRouter(router *gin.RouterGroup) { // ✅ bỏ db *gorm.DB
	router.GET("/list", controllers.Loss.GetLossList)
	
	// Loss Tem Size routes
	temSizeGroup := router.Group("/loss-tem-size")
	temSizeGroup.POST("/save", controllers.LossTemSize.Save)
	temSizeGroup.GET("/list", controllers.LossTemSize.GetList)
	temSizeGroup.GET("/summary", controllers.LossTemSize.GetSummary)
}
