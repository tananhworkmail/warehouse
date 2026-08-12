package controllers

import (
	"errors"
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type WarehouseKTPController struct {
	*BaseController
}

var WarehouseKTP = &WarehouseKTPController{}

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
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Ma scan phai dung 6 chu so")
	case errors.Is(err, services.ErrWarehouseKTPRackNotFound):
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "Khong tim thay ke")
	case errors.Is(err, services.ErrWarehouseKTPScanStorageNotFound):
		response.FailWithDetailed(ctx, http.StatusNotImplemented, nil, "Du lieu KTP duoc nhap tu he thong ERP khac; trang nay chi doc YWCP va khong ghi scan")
	case errors.Is(err, services.ErrWarehouseKTPCartonNotFound):
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "Khong tim thay barcode thung carton")
	case errors.Is(err, services.ErrWarehouseKTPOrderNotFound):
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "Khong tim thay don hang")
	case errors.Is(err, services.ErrWarehouseKTPInvalidRackCode):
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Ma ke khong hop le")
	case errors.Is(err, services.ErrWarehouseKTPSameRack):
		response.FailWithDetailed(ctx, http.StatusConflict, nil, "Ke moi trung voi ke hien tai cua don hang")
	default:
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
	}
}

func getCreatedBy(ctx *gin.Context) string {
	for _, header := range []string{"X-User", "X-Username"} {
		if value := ctx.GetHeader(header); value != "" {
			return value
		}
	}
	return ""
}
