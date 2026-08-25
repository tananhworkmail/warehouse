package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

// PrintLableController 控制器結構
type PrintLableController struct {
	*BaseController
}

// Pt 是 PrintLableController 的全局實例
var Pt = &PrintLableController{}

// PrintZPL 處理打印標籤的請求
func (c *PrintLableController) PrintZPL(ctx *gin.Context) {
	var requestParams request.PrintLableRequest

	// 驗證請求參數
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	// 調用服務層打印 ZPL
	err := services.PrintZPL(
		requestParams.IP,
		requestParams.ZPLCommand,
	)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	// 返回成功響應
	response.Ok(ctx)
}

// CheckPrinterExists 處理檢查打印機存在的請求
func (c *PrintLableController) CheckPrinterExists(ctx *gin.Context) {
	var requestParams request.CheckPrinterRequest

	// 驗證請求參數
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	// 調用服務層檢查打印機是否存在
	exists, err := services.CheckPrinterExists(requestParams.IP)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	// 返回成功響應，包含打印機存在的狀態
	response.OkWithData(ctx, gin.H{"exists": exists})
}
