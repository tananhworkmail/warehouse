package controllers

import (

	// "mes-web-api/internal/pkg/models/request"

	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type DeliveryController struct {
	*BaseController
}

var Dl = &DeliveryController{}

func (c *DeliveryController) GetScanDelivery(ctx *gin.Context) {
	var requestParams request.ScanGetKCLLRequest

	// Validate request parameters
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	// Pass LLNO slice directly to the service method
	result, err := services.DL.GetScanDelivery(requestParams.LLNO)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	// Send the result back to the client
	response.OkWithData(ctx, result)
}

func (c *DeliveryController) GetScanDeliveryQuantity(ctx *gin.Context) {

	var requestParams request.ScanDeliveryLabelQuantity

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.DL.GetScanDeliveryQuantity(requestParams.Quantity, requestParams.Barcode,requestParams.KCBH, requestParams.USERID)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *DeliveryController) GetScanDeliveryCLBH(ctx *gin.Context) {

	var requestParams request.ScanDeliveryLabelCLBH

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.DL.GetScanDeliveryCLBH(requestParams.Barcode, requestParams.USERID, requestParams.Data)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}
