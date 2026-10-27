package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type DeliveryPayController struct {
	*BaseController
}

var Dp = &DeliveryPayController{}

func (c *DeliveryPayController) GetScanDeliveryPay(ctx *gin.Context) {

	var requestParams request.ScanGetKCLLSRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.DP.GetScanDeliveryPay(requestParams.LLNO)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *DeliveryPayController) GetScanDeliveryPayKCLLScan_RFSS(ctx *gin.Context) {

	var requestParams request.ScanGetKCLLScan_RFSSRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.DP.GetScanDeliveryPayKCLLScan_RFSS(requestParams.LLNO, requestParams.CLBH, requestParams.DFL, requestParams.SCBH)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *DeliveryPayController) ScanDeliveryPayExPort(ctx *gin.Context) {

	var requestParams request.DeliveryPayExPortRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.DP.ScanDeliveryPayExPort(requestParams.SCBH, requestParams.CLBH, requestParams.YWPM, requestParams.SCFMDATE, requestParams.ECFMDATE)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *DeliveryPayController) DBGridEh2DblClick(ctx *gin.Context) {
	var requestParams request.DeliveryDblClick

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	err := services.DP.DBGridEh2DblClick(
		requestParams.LLNO_KCLLS,
		requestParams.LLNO_OWE_LLSSS,
		requestParams.CLBH_KCLL,
		requestParams.SCBH_KCLL,
		requestParams.DFL_KCLL,
		requestParams.USERID,
		requestParams.LLQty_KCLL,
		requestParams.Qty_KCLL,
	)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *DeliveryPayController) DeliveryPayKCLLScan_RFSSS(ctx *gin.Context) {
	var requestParams request.DeliveryPayKCLLScan_RFSS

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	err := services.DP.DeliveryPayKCLLScan_RFSSS(

		requestParams.LLNO,
		requestParams.LLNO_OWE,
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
