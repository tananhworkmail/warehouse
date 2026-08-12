package router_v1

import (
	"web-api/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	v1.GET("/sse", handlers.SSEHandler)

	RegisterCommonRouter(v1.Group(""))

	RegisterReportRouter(v1.Group("/rp"))

	RegisterAdminRouter(v1.Group("/rp/admin"))

	LoginRouter(v1.Group("/auth"))

	TempHumidityCheckRouter(v1.Group(""))

	GoreTexFormRouter(v1.Group("/gore-tex"))
	WarehouseRouter(v1.Group("/warehouse"))
	WarehouseKTPRouter(v1.Group("/warehouse-ktp"))
	LossRouter(v1.Group("/loss"))
}
