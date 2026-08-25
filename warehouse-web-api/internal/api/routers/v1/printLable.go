package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func PrintLableRouter(router *gin.RouterGroup) {

	router.POST("/pstprnt", controllers.Pt.PrintZPL)        // Print ZPL Command
	router.GET("/check", controllers.Pt.CheckPrinterExists) // Check Printer Exists

}
