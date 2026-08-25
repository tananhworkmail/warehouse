package router_v1

import (
	"web-api/internal/api/controllers"
	"web-api/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	v1.GET("/sse", handlers.SSEHandler)

	RegisterCommonRouter(v1.Group(""))

	RegisterReportRouter(v1.Group("/rp"))

	RegisterAdminRouter(v1.Group("/rp/admin"))

	// All web modules use the branch-aware login contract. Keep /scan/login as
	// a compatibility alias for older N8d deployments.
	v1.Group("/auth").POST("/login", controllers.User_Login.Login)

	// N8d material tracking and inventory APIs.
	scan := v1.Group("/scan")
	N8dLoginRouter(scan)
	TrackingRouter(scan.Group("/tracking"))
	InventoryRouter(scan.Group("/number-scan"))

	TempHumidityCheckRouter(v1.Group(""))

	GoreTexFormRouter(v1.Group("/gore-tex"))
	WarehouseRouter(v1.Group("/warehouse"))
	WarehouseKTPRouter(v1.Group("/warehouse-ktp"))
	LossRouter(v1.Group("/loss"))
}
