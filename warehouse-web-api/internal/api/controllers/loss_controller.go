package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/types"

	"github.com/gin-gonic/gin"
)

// ─────────────────────────────────────────────
// LOSS CONTROLLER (code gốc - không sửa)
// ─────────────────────────────────────────────

type LossController struct {
	service *services.LossService
}

var Loss = &LossController{}

func NewLossController(service *services.LossService) *LossController {
	return &LossController{service: service}
}

func (ctrl *LossController) GetLossList(c *gin.Context) {
	var req types.LossListRequest

	// Bind query params: /loss?page=1&page_size=20&zlbh=JHS2607-025
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	result, err := ctrl.service.GetLossList(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Lỗi lấy danh sách hao hụt: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    result,
	})
}

// ─────────────────────────────────────────────
// LOSS TEM SIZE CONTROLLER (thêm mới)
// ─────────────────────────────────────────────

type LossTemSizeController struct {
	service *services.LossTemSizeService
}

var LossTemSize = &LossTemSizeController{
	service: services.LossTemSize,
}

func NewLossTemSizeController(service *services.LossTemSizeService) *LossTemSizeController {
	return &LossTemSizeController{service: service}
}

// POST /loss-tem-size/save
func (ctrl *LossTemSizeController) Save(c *gin.Context) {
	var req types.LossTemSizeSaveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	if err := ctrl.service.Save(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(), // ← trả thẳng message từ service, không prefix
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Lưu thành công",
	})
}

// GET /loss-tem-size/list
func (ctrl *LossTemSizeController) GetList(c *gin.Context) {
	var req types.LossTemSizeListRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	result, err := ctrl.service.GetList(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Lỗi lấy danh sách: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    result,
	})
}

// GET /loss-tem-size/summary
func (ctrl *LossTemSizeController) GetSummary(c *gin.Context) {
	var req types.LossTemSizeListRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	result, err := ctrl.service.GetSummary(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Lỗi lấy summary: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    result,
	})
}