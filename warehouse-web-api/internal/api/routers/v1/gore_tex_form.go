package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func GoreTexFormRouter(router *gin.RouterGroup) {
	router.GET("/dashboard/weekly", controllers.GoreTexForms.WeeklyDashboard)

	forms := router.Group("/forms")
	{
		forms.GET("", controllers.GoreTexForms.List)

		forms.POST("/waterproof", controllers.GoreTexForms.SubmitWaterproof)
		forms.GET("/waterproof", controllers.GoreTexForms.GetWaterproof)

		forms.POST("/centrifugal", controllers.GoreTexForms.SubmitCentrifugal)
		forms.GET("/centrifugal/:inspectionDate", controllers.GoreTexForms.GetCentrifugal)

		forms.POST("/analysis", controllers.GoreTexForms.SubmitAnalysis)
		forms.GET("/analysis/:id", controllers.GoreTexForms.GetAnalysis)
	}
}
