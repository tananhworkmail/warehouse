package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type TrackingController struct {
	*BaseController
}

var Tk = &TrackingController{}

func (c *TrackingController) ScanMaterialImport(ctx *gin.Context) {

	var requestParams request.ScanMaterialIpRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.TK.ScanMaterialImport(
		requestParams.GSBH,
		requestParams.SCBH,
		requestParams.CGNO,
		requestParams.ZSDH,
		requestParams.ZSYWJC,
		requestParams.CLBH,
		requestParams.YWBM,
		requestParams.MAKE,     // ✅ đúng
		requestParams.Barcode,  // ✅ đúng
		requestParams.SCFMDATE, // ✅ đúng
		requestParams.ECFMDATE, // ✅ đúng
	)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

type UpdateKCBHRequest struct {
	Barcode string `json:"barcode" binding:"required"`
	NewKCBH string `json:"newKCBH" binding:"required"`
}

func (c *TrackingController) UpdateKCBHByBarcode(ctx *gin.Context) {

	var requestParams UpdateKCBHRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	err := services.TK.UpdateKCBHByBarcode(
		requestParams.Barcode,
		requestParams.NewKCBH,
	)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithMessage(ctx, "Cập nhật MAKE thành công")
}

func (c *TrackingController) ScanMaterialExport(ctx *gin.Context) {

	var requestParams request.ScanMaterialIpRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.TK.ScanMaterialExport(requestParams.GSBH, requestParams.SCBH, requestParams.CGNO, requestParams.ZSDH, requestParams.ZSYWJC, requestParams.CLBH, requestParams.YWBM, requestParams.SCFMDATE, requestParams.ECFMDATE)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *TrackingController) ScanMaterialPrint(ctx *gin.Context) {
	var requestParams request.ScanMaterialPrintRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	// SỬA TẠI ĐÂY: Thêm requestParams.MAKE và requestParams.Barcode vào giữa YWBM và SCFMDATE
	result, err := services.TK.ScanMaterialPrint(
		requestParams.GSBH,
		requestParams.SCBH,
		requestParams.CGNO,
		requestParams.ZSDH,
		requestParams.ZSYWJC,
		requestParams.CLBH,
		requestParams.YWBM,
		requestParams.MAKE,    // Tham số mới 1
		requestParams.Barcode, // Tham số mới 2
		requestParams.SCFMDATE,
		requestParams.ECFMDATE,
	)

	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *TrackingController) GetKCLLScan_RFSS(ctx *gin.Context) {
	var requestParams request.GetKCLLScan_RFSS

	// 验证请求参数
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	// 调用服务方法并传递请求参数
	result, err := services.TK.GetKCLLScan_RFSS(requestParams.Barcode)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	// 返回成功响应和结果
	response.OkWithData(ctx, result)
}

func (c *TrackingController) GetKCLLScan_RFSSS(ctx *gin.Context) {
	var requestParams request.GetKCLLScan_RFSSS

	// 验证请求参数
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	// 调用服务方法并传递请求参数
	result, err := services.TK.GetKCLLScan_RFSSS(requestParams.Barcode)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	// 返回成功响应和结果
	response.OkWithData(ctx, result)
}

func (c *TrackingController) UpdateKCLLScan_RFSS(ctx *gin.Context) {
	var requestParams request.UpdateKCLLScan_RFSSRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	err := services.TK.UpdateKCLLScan_RFSS(
		requestParams.LLNO,
		requestParams.SCBH,
		requestParams.CLBH,
		requestParams.Barcode,
		requestParams.DFL,
		requestParams.SCANUSERID,
		requestParams.USERID,
		requestParams.ScanQty_input,
		requestParams.SCANQty,
	)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *TrackingController) UpdateKCRKScan_RFSS(ctx *gin.Context) {
	var requestParams request.UpdateKCRKScan_RFSSRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	err := services.TK.UpdateKCRKScan_RFSS(
		requestParams.Barcode,
	)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *TrackingController) UpdateKCLLScan_RFSSS(ctx *gin.Context) {
	var requestParams request.UpdateKCLLScan_RFSSSRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	err := services.TK.UpdateKCLLScan_RFSSS(
		requestParams.LLNO_OWE,
		requestParams.LLNO,
		requestParams.SCBH,
		requestParams.CLBH,
		requestParams.Barcode,

		requestParams.SCANUSERID,
		requestParams.USERID,
		requestParams.ScanQty_input,
		requestParams.SCANQty,
	)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.Ok(ctx)
}
