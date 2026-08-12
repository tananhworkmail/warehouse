package controllers

import (
	"fmt"
	"net/http"

	"web-api/internal/api/services"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

// HumidityAlertController — controller riêng cho cảnh báo độ ẩm + điều khiển Andon
type HumidityAlertController struct {
	*BaseController
}

var HumidityAlertCtrl = &HumidityAlertController{}

// GetLatestHumidityStatus — lấy trạng thái độ ẩm realtime 3 kho
// GET /api/v1/warehouse/humidity/status
func (c *HumidityAlertController) GetLatestHumidityStatus(ctx *gin.Context) {
	data, err := services.HumidityAlert.GetLatestHumidityStatus()
	if err != nil {
		fmt.Println("GetLatestHumidityStatus error:", err)
		response.FailWithMessage(ctx, "Lỗi lấy trạng thái độ ẩm")
		return
	}
	response.OkWithData(ctx, data)
}

// GetAlertLog — lấy lịch sử cảnh báo
// GET /api/v1/warehouse/humidity/log?deviceName=KVT02&date=2026-04-14
func (c *HumidityAlertController) GetAlertLog(ctx *gin.Context) {
	deviceName := ctx.Query("deviceName")
	date := ctx.Query("date")

	data, err := services.HumidityAlert.GetAlertLog(deviceName, date)
	if err != nil {
		fmt.Println("GetAlertLog error:", err)
		response.FailWithMessage(ctx, "Lỗi lấy lịch sử cảnh báo")
		return
	}
	response.OkWithData(ctx, data)
}

// TurnOffAndon — tắt Andon từ UI
// POST /api/v1/warehouse/humidity/andon/off
// Body: { "deviceName": "KVT02", "offBy": "admin" }
func (c *HumidityAlertController) TurnOffAndon(ctx *gin.Context) {
	var body struct {
		DeviceName string `json:"deviceName" binding:"required"`
		OffBy      string `json:"offBy"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	if body.OffBy == "" {
		body.OffBy = "manual"
	}

	if err := services.HumidityAlert.TurnOffAndon(body.DeviceName, body.OffBy); err != nil {
		fmt.Println("TurnOffAndon error:", err)
		response.FailWithMessage(ctx, err.Error())
		return
	}
	response.OkWithData(ctx, gin.H{"message": "Đã tắt Andon " + body.DeviceName})
}

// TurnOnAndon — bật/toggle Andon thủ công từ UI
// POST /api/v1/warehouse/humidity/andon/on
// Body: { "deviceName": "KVT02", "ledId": 3 }
// ledId: 1=Xanh, 2=Vàng, 3=Đỏ, 4=Còi, 0=Tắt tất cả
// func (c *HumidityAlertController) TurnOnAndon(ctx *gin.Context) {
// 	var body struct {
// 		DeviceName string `json:"deviceName" binding:"required"`
// 		LedID      int    `json:"ledId"`
// 	}
// 	if err := ctx.ShouldBindJSON(&body); err != nil {
// 		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
// 		return
// 	}

// 	if err := services.HumidityAlert.TurnOnAndon(body.DeviceName, body.LedID); err != nil {
// 		fmt.Println("TurnOnAndon error:", err)
// 		response.FailWithMessage(ctx, err.Error())
// 		return
// 	}
// 	response.OkWithData(ctx, gin.H{"message": "Đã gửi lệnh Andon " + body.DeviceName})
// }

// GetAndonStatus — lấy trạng thái đèn hiện tại của ESP
// GET /api/v1/warehouse/humidity/andon/status?deviceName=KVT02
func (c *HumidityAlertController) GetAndonStatus(ctx *gin.Context) {
	deviceName := ctx.Query("deviceName")
	if deviceName == "" {
		response.FailWithMessage(ctx, "deviceName is required")
		return
	}

	status, err := services.HumidityAlert.GetAndonStatus(deviceName)
	if err != nil {
		response.FailWithMessage(ctx, err.Error())
		return
	}
	response.OkWithData(ctx, gin.H{"raw": status})
}