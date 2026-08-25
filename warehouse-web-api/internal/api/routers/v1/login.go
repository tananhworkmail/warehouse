package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func LoginRouter(router *gin.RouterGroup) {

	router.POST("/login", controllers.Lg.Login)
	router.GET("/me", controllers.Lg.Login)
}

// N8dLoginRouter is kept as a compatibility alias for existing N8d clients.
func N8dLoginRouter(router *gin.RouterGroup) {
	router.POST("/Login", controllers.User_Login.Login)
	router.POST("/login", controllers.User_Login.Login)
}
