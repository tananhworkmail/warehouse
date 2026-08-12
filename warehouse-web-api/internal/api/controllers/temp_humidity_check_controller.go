package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"web-api/internal/api/services"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type TempHumidityCheckController struct {
	*BaseController
}

var TempHumidityCheckCtrl = &TempHumidityCheckController{}

func (c *TempHumidityCheckController) Check(ctx *gin.Context) {
	devices := parseTempHumidityDevices(ctx.Query("devices"))

	data, err := services.TempHumidityCheck.CheckLatest(devices)
	if err != nil {
		fmt.Println("TempHumidityCheck error:", err)
		response.FailWithMessage(ctx, "Loi lay nhiet do, do am")
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (c *TempHumidityCheckController) CheckSafeHumidity(ctx *gin.Context) {
	devices := parseTempHumidityDevices(ctx.Query("devices"))

	data, err := services.TempHumidityCheck.CheckLatestSafeHumidity(devices)
	if err != nil {
		fmt.Println("TempHumidityCheck safe humidity error:", err)
		response.FailWithMessage(ctx, "Loi lay nhiet do, do am")
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func parseTempHumidityDevices(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	return strings.Split(raw, ",")
}
