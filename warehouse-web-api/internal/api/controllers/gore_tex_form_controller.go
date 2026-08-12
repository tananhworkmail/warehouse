package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type GoreTexFormController struct {
	*BaseController
}

var GoreTexForms = &GoreTexFormController{}

func (c *GoreTexFormController) List(ctx *gin.Context) {
	data, err := services.GoreTexForms.List()
	if err != nil {
		writeGoreTexFormError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func (c *GoreTexFormController) WeeklyDashboard(ctx *gin.Context) {
	now := time.Now()
	defaultYear, defaultWeek := now.ISOWeek()
	year, week := defaultYear, defaultWeek
	var err error
	if value := ctx.Query("year"); value != "" {
		year, err = strconv.Atoi(value)
		if err != nil {
			writeGoreTexFormError(ctx, services.ErrGoreTexInvalidForm)
			return
		}
	}
	if value := ctx.Query("week"); value != "" {
		week, err = strconv.Atoi(value)
		if err != nil {
			writeGoreTexFormError(ctx, services.ErrGoreTexInvalidForm)
			return
		}
	}

	data, err := services.GoreTexForms.WeeklyDashboard(year, week)
	if err != nil {
		writeGoreTexFormError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func (c *GoreTexFormController) SubmitWaterproof(ctx *gin.Context) {
	var params request.GoreTexWaterproofSubmitRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	data, err := services.GoreTexForms.SaveWaterproof(params)
	if err != nil {
		writeGoreTexFormError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func (c *GoreTexFormController) SubmitCentrifugal(ctx *gin.Context) {
	var params request.GoreTexCentrifugalSubmitRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	data, err := services.GoreTexForms.SaveCentrifugal(params)
	if err != nil {
		writeGoreTexFormError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func (c *GoreTexFormController) SubmitAnalysis(ctx *gin.Context) {
	var params request.GoreTexAnalysisSubmitRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	data, err := services.GoreTexForms.SaveAnalysis(params)
	if err != nil {
		writeGoreTexFormError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func (c *GoreTexFormController) GetWaterproof(ctx *gin.Context) {
	data, err := services.GoreTexForms.GetWaterproof(ctx.Query("line"), ctx.Query("styleName"))
	if err != nil {
		writeGoreTexFormError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func (c *GoreTexFormController) GetCentrifugal(ctx *gin.Context) {
	data, err := services.GoreTexForms.GetCentrifugal(ctx.Param("inspectionDate"))
	if err != nil {
		writeGoreTexFormError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func (c *GoreTexFormController) GetAnalysis(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		writeGoreTexFormError(ctx, services.ErrGoreTexInvalidForm)
		return
	}

	data, err := services.GoreTexForms.GetAnalysis(id)
	if err != nil {
		writeGoreTexFormError(ctx, err)
		return
	}
	response.OkWithData(ctx, data)
}

func writeGoreTexFormError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrGoreTexInvalidForm):
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
	case errors.Is(err, services.ErrGoreTexFormNotFound):
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, err.Error())
	default:
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
	}
}
