package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type DeliveryOweController struct {
	*BaseController
}

var Do = &DeliveryOweController{}

func (c *DeliveryOweController) ScanDeliveryOweQuantity(ctx *gin.Context) {
	var requestParams request.DeliveryOweQuantity

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	resultKCRKScan_RFSS, resultKCLLScan_RFSSS, err := services.DO.ScanDeliveryOweQuantity(requestParams.Quantity, requestParams.Barcode, requestParams.USERID, requestParams.LLNO, requestParams.CLBH, requestParams.DFL, requestParams.SCBH)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	responseData := map[string]interface{}{
		"resultKCRKScan_RFSS":  resultKCRKScan_RFSS,
		"resultKCLLScan_RFSSS": resultKCLLScan_RFSSS,
	}

	response.OkWithData(ctx, responseData)
}

func (c *DeliveryOweController) ScanDeliveryOwe(ctx *gin.Context) {
	var requestParams request.DeliveryOweQuantity

	// Validate the request parameters
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	// Call the service function
	resultKCLLScan_RFSSS, resultMaterialPrint, err := services.DO.ScanDeliveryOwe(requestParams.Barcode, requestParams.USERID, requestParams.LLNO, requestParams.CLBH, requestParams.DFL, requestParams.SCBH)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	// Prepare the response data
	responseData := map[string]interface{}{
		"resultKCLLScan_RFSSS": resultKCLLScan_RFSSS,
		"resultMaterialPrint":  resultMaterialPrint,
	}

	// Send the response
	response.OkWithData(ctx, responseData)
}

func (c *DeliveryOweController) DeliveryOweKCLLScan_RFSS(ctx *gin.Context) {
	var requestParams request.DeliveryOweKCLLScan_RFSSS

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	err := services.DO.DeliveryOweKCLLScan_RFSS(
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
