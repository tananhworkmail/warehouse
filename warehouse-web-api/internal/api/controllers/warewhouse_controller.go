package controllers

import (
	"fmt"
	"net/http"

	"web-api/internal/api/services"
	"web-api/internal/pkg/models/entities"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type WarehouseController struct {
	*BaseController
}

var Warehouse = &WarehouseController{}

func (c *WarehouseController) GetWareHouseList(ctx *gin.Context) {
	var req request.GetWarehouseListRequest

	// Kiểm tra dữ liệu đầu vào
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	// Gọi service để tìm đơn hàng
	warehouse, err := services.Warehouse.GetWareHouseList(req)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	// Trả về dữ liệu đơn hàng
	response.OkWithData(ctx, warehouse)
}
func (c *WarehouseController) GetRackInformation(ctx *gin.Context) {
	var requestParams request.WarehouseRackGetRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		fmt.Println("❌ Validation failed:", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	fmt.Println("✅ Request rackNo:", requestParams.RackNo)

	result, err := services.Warehouse.GetRackInformation(requestParams.RackNo)
	if err != nil {
		fmt.Println("❌ GetRackInformation error:", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	fmt.Println("✅ Result length:", len(result))
	response.OkWithData(ctx, result)
}
func (c *WarehouseController) Get3DayRackInformation(ctx *gin.Context) {
	var requestParams request.WarehouseRackGetRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		fmt.Println("❌ Validation failed:", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	fmt.Println("✅ Request rackNo:", requestParams.RackNo)

	result, err := services.Warehouse.Get3DayRackInformation(requestParams.RackNo)
	if err != nil {
		fmt.Println("❌ Get3DayRackInformation error:", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	fmt.Println("✅ Result length:", len(result))
	response.OkWithData(ctx, result)
}
func (c *WarehouseController) GetTonKhoByRack(ctx *gin.Context) {
	rackPrefix := ctx.Query("rackNo")
	if rackPrefix == "" {
		response.FailWithMessage(ctx, "rackNo query param is required")
		return
	}

	data, err := services.Warehouse.GetTonKhoByRackPrefix(rackPrefix)
	if err != nil {
		fmt.Println("GetTonKhoByRack error:", err)
		response.FailWithMessage(ctx, "Error loading ton kho")
		return
	}

	fmt.Println("GetTonKhoByRack - result count:", len(data))
	response.OkWithData(ctx, data)
}
func (c *WarehouseController) GetTonKhoKeByRack(ctx *gin.Context) {

	data, err := services.Warehouse.GetTonKhoKeByRack()
	if err != nil {
		fmt.Println("GetTonKhoKeByRack error:", err)
		response.FailWithMessage(ctx, "Error loading tồn kho theo kệ")
		return
	}

	fmt.Println("GetTonKhoKeByRack - result count:", len(data))
	response.OkWithData(ctx, data)
}

func (c *WarehouseController) GetTonKhoTangByRack(ctx *gin.Context) {

	data, err := services.Warehouse.GetTonKhoTangByRack()
	if err != nil {
		fmt.Println(" GetTonKhoTangByRack error:", err)
		response.FailWithMessage(ctx, "Error loading tồn kho theo kệ")
		return
	}

	fmt.Println("GetTonKhoTangByRack - result count:", len(data))
	response.OkWithData(ctx, data)
}
func (c *WarehouseController) GetSearchByDDBH(ctx *gin.Context) {
	ddbh := ctx.Query("DDBH")
	if ddbh == "" {
		response.FailWithMessage(ctx, "Query parameter 'DDBH' is required")
		return
	}

	data, err := services.Warehouse.SearchByDDBH(ddbh)
	if err != nil {
		fmt.Println("SearchByDDBH error:", err)
		response.FailWithMessage(ctx, "Error retrieving DDBH data")
		return
	}

	fmt.Printf("SearchByDDBH - result count: %d\n", len(data))
	response.OkWithData(ctx, data)
}
func (c *WarehouseController) GetSearchTotalDDBH(ctx *gin.Context) {
	ddbh := ctx.Query("DDBH")

	data, err := services.Warehouse.SearchDDBH(ddbh)
	if err != nil {
		fmt.Println("SearchByDDBH error:", err)
		response.FailWithMessage(ctx, "Error retrieving DDBH data")
		return
	}

	fmt.Printf("SearchByDDBH - result count: %d\n", len(data))
	response.OkWithData(ctx, data)
}
func (c *WarehouseController) UploadSchedule(ctx *gin.Context) {
	var requestParams request.UploadScheduleRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.Warehouse.UploadSchedule(requestParams.BuildingNo, file); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}
func (c *WarehouseController) GetSchedule(ctx *gin.Context) {
	var requestParams request.GetScheduleGetRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var result entities.ScheduleFile
	if err := services.Warehouse.GetSchedule(entities.ScheduleFile{BuildingNo: requestParams.BuildingNo}, &result, []string{}); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	fileData := []byte(result.FileStream)


	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileData)
}
func (c *WarehouseController) GetOrderExistenceStatus(ctx *gin.Context) {
	var requestParams request.TraceOrderExistenceStatusRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.Warehouse.GetOrderExistenceStatus(requestParams.OrderNo)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}
func (c *WarehouseController) GetBomInformation(ctx *gin.Context) {
	var requestParams request.TraceProductionInformationRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.Warehouse.GetTraceServiceAll(requestParams.OrderNo)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}