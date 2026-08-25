package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"web-api/internal/api/services"
	"web-api/internal/pkg/auth"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type WarehouseKTPController struct {
	*BaseController
}

var WarehouseKTP = &WarehouseKTPController{}

func (c *WarehouseKTPController) GetScanInAccess(ctx *gin.Context) {
	userID, _ := ctx.Get(auth.UserIDContextKey)
	response.OkWithData(ctx, gin.H{
		"allowed": true,
		"userId":  userID,
	})
}

func (c *WarehouseKTPController) GetScanOutAccess(ctx *gin.Context) {
	userID, _ := ctx.Get(auth.UserIDContextKey)
	response.OkWithData(ctx, gin.H{
		"allowed": true,
		"userId":  userID,
	})
}

func (c *WarehouseKTPController) GetScanShipAccess(ctx *gin.Context) {
	userID, _ := ctx.Get(auth.UserIDContextKey)
	response.OkWithData(ctx, gin.H{
		"allowed": true,
		"userId":  userID,
	})
}

func (c *WarehouseKTPController) GetRacks(ctx *gin.Context) {
	data, err := services.WarehouseKTP.GetRacks()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) SearchRacks(ctx *gin.Context) {
	data, err := services.WarehouseKTP.SearchRacks(ctx.Query("keyword"))
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) SearchInventory(ctx *gin.Context) {
	keyword := strings.TrimSpace(ctx.Query("keyword"))
	if keyword == "" {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "keyword is required")
		return
	}

	data, err := services.WarehouseKTP.SearchInventory(keyword)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) GetInventoryAge(ctx *gin.Context) {
	data, err := services.WarehouseKTP.GetInventoryAge()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) GetRackOrderDetail(ctx *gin.Context) {
	data, err := services.WarehouseKTP.GetRackOrderDetail(ctx.Param("rackCode"))
	if err != nil {
		writeWarehouseKTPError(ctx, err)
		return
	}

	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) GetTempHumidityByDevices(ctx *gin.Context) {
	data, err := services.WarehouseKTP.GetTempHumidityByDevices()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) GetLatestTempHumidityLogByDevices(ctx *gin.Context) {
	data, err := services.WarehouseKTP.GetLatestTempHumidityLogByDevices(
		ctx.Query("date"),
		ctx.Query("startTime"),
		ctx.Query("endTime"),
	)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) GetMoveOrderByCarton(ctx *gin.Context) {
	data, err := services.WarehouseKTP.GetMoveOrderByCarton(ctx.Query("cartonBar"))
	if err != nil {
		writeWarehouseKTPError(ctx, err)
		return
	}

	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) GetScanInCarton(ctx *gin.Context) {
	var actualWeight *float64
	if rawWeight := strings.TrimSpace(ctx.Query("actualWeight")); rawWeight != "" {
		parsedWeight, err := strconv.ParseFloat(rawWeight, 64)
		if err != nil {
			response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Trọng lượng thực tế không hợp lệ")
			return
		}
		actualWeight = &parsedWeight
	}

	sessionCartons := strings.Split(ctx.Query("sessionCartons"), ",")
	data, err := services.WarehouseKTP.GetScanInCarton(
		ctx.Param("barcode"),
		actualWeight,
		sessionCartons,
		ctx.Query("scanMode"),
		nil,
	)
	if err != nil {
		writeWarehouseKTPError(ctx, err)
		return
	}

	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) GetScanInProductionLines(ctx *gin.Context) {
	data, err := services.WarehouseKTP.GetScanInProductionLines(ctx.Query("factory"))
	if err != nil {
		writeWarehouseKTPError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) GetScanInWarehouses(ctx *gin.Context) {
	data, err := services.WarehouseKTP.GetScanInWarehouses()
	if err != nil {
		writeWarehouseKTPError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) SubmitScanInCarton(ctx *gin.Context) {
	var params request.WarehouseKTPScanInSubmitRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Thông tin nhập kho chưa đầy đủ hoặc không hợp lệ")
		return
	}

	data, err := services.WarehouseKTP.SubmitScanInCarton(params, getCreatedBy(ctx))
	if err != nil {
		if errors.Is(err, services.ErrWarehouseKTPScanInRejected) {
			response.FailWithDetailed(ctx, http.StatusConflict, data, err.Error())
			return
		}
		writeWarehouseKTPError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) GetScanOutCarton(ctx *gin.Context) {
	data, err := services.WarehouseKTP.GetScanOutCarton(ctx.Param("barcode"))
	if err != nil {
		writeWarehouseKTPError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) SubmitScanOutCarton(ctx *gin.Context) {
	var params request.WarehouseKTPScanOutSubmitRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Phải có mã carton và chọn Tái chế hoặc Kiểm hàng")
		return
	}

	data, err := services.WarehouseKTP.SubmitScanOutCarton(params)
	if err != nil {
		writeWarehouseKTPError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) GetScanShipCarton(ctx *gin.Context) {
	weightEnabled, _ := strconv.ParseBool(ctx.Query("weightEnabled"))
	var outgoingWeight *float64
	if rawWeight := strings.TrimSpace(ctx.Query("outgoingWeight")); rawWeight != "" {
		parsedWeight, err := strconv.ParseFloat(rawWeight, 64)
		if err != nil {
			response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Trọng lượng xuất không hợp lệ")
			return
		}
		outgoingWeight = &parsedWeight
	}
	data, err := services.WarehouseKTP.GetScanShipCarton(
		ctx.Param("barcode"),
		ctx.Query("containerNo"),
		ctx.Query("plateId"),
		outgoingWeight,
		weightEnabled,
	)
	if err != nil {
		writeWarehouseKTPError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) SubmitScanShipCartons(ctx *gin.Context) {
	var params request.WarehouseKTPScanShipSubmitRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Thông tin Scan Ship chưa đầy đủ hoặc không hợp lệ")
		return
	}
	data, err := services.WarehouseKTP.SubmitScanShipCartons(params, getCreatedBy(ctx))
	if err != nil {
		writeWarehouseKTPError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) MoveOrderToRack(ctx *gin.Context) {
	var params request.WarehouseKTPMoveOrderRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	data, err := services.WarehouseKTP.MoveOrderToRack(params.CartonBar, params.NewRackCode, getCreatedBy(ctx))
	if err != nil {
		writeWarehouseKTPError(ctx, err)
		return
	}

	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) ScanRack(ctx *gin.Context) {
	rackCode := ctx.Param("rackCode")
	var params request.WarehouseKTPScanRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	data, err := services.WarehouseKTP.ScanRack(rackCode, params.ScanCode, params.Replace, getCreatedBy(ctx))
	if err != nil {
		writeWarehouseKTPError(ctx, err)
		return
	}

	response.OkWithData(ctx, data)
}

func (c *WarehouseKTPController) ClearRack(ctx *gin.Context) {
	data, err := services.WarehouseKTP.ClearRack(ctx.Param("rackCode"), getCreatedBy(ctx))
	if err != nil {
		writeWarehouseKTPError(ctx, err)
		return
	}

	response.OkWithData(ctx, data)
}

func writeWarehouseKTPError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrWarehouseKTPInvalidScanCode):
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Mã scan phải đúng 6 chữ số")
	case errors.Is(err, services.ErrWarehouseKTPRackNotFound):
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "Không tìm thấy kệ")
	case errors.Is(err, services.ErrWarehouseKTPScanStorageNotFound):
		response.FailWithDetailed(ctx, http.StatusNotImplemented, nil, "Dữ liệu KTP được nhập từ hệ thống ERP khác")
	case errors.Is(err, services.ErrWarehouseKTPCartonNotFound):
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "Không tìm thấy barcode thùng carton")
	case errors.Is(err, services.ErrWarehouseKTPOrderNotFound):
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "Không tìm thấy đơn hàng")
	case errors.Is(err, services.ErrWarehouseKTPInvalidRackCode):
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Mã kệ không hợp lệ")
	case errors.Is(err, services.ErrWarehouseKTPSameRack):
		response.FailWithDetailed(ctx, http.StatusConflict, nil, "Kệ mới trùng với kệ hiện tại của đơn hàng")
	case errors.Is(err, services.ErrWarehouseKTPScanInUpdateFailed):
		response.FailWithDetailed(ctx, http.StatusConflict, nil, "Không cập nhật được carton vào YWCP")
	case errors.Is(err, services.ErrWarehouseKTPScanOutRejected):
		response.FailWithDetailed(ctx, http.StatusConflict, nil, err.Error())
	case errors.Is(err, services.ErrWarehouseKTPScanOutUpdateFailed):
		response.FailWithDetailed(ctx, http.StatusConflict, nil, "Không cập nhật được carton Scan Out; thay đổi đã được hoàn tác")
	case errors.Is(err, services.ErrWarehouseKTPScanShipRejected):
		response.FailWithDetailed(ctx, http.StatusConflict, nil, err.Error())
	case errors.Is(err, services.ErrWarehouseKTPScanShipUpdateFailed):
		response.FailWithDetailed(ctx, http.StatusConflict, nil, "Không cập nhật được Scan Ship; toàn bộ phiên đã được hoàn tác")
	default:
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
	}
}

func getCreatedBy(ctx *gin.Context) string {
	if value, exists := ctx.Get(auth.UserIDContextKey); exists {
		if userID := strings.TrimSpace(fmt.Sprint(value)); userID != "" {
			return userID
		}
	}
	for _, header := range []string{"X-User", "X-Username"} {
		if value := ctx.GetHeader(header); value != "" {
			return value
		}
	}
	return ""
}
