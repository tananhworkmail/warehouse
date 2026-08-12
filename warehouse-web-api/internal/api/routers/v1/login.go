package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func LoginRouter(router *gin.RouterGroup) {

	router.POST("/login", controllers.Lg.Login)
	router.GET("/me", controllers.Lg.Login)
}
