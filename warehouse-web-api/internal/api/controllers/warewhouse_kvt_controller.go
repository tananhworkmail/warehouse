package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type WarehouseKVTController struct {
	*BaseController
}

var WarehouseKVT = &WarehouseKVTController{}

// Số ô trong mỗi tầng
func (c *WarehouseKVTController) GetRackKVT(ctx *gin.Context) {
	rackCode := ctx.Query("rackCode")
	data, err := services.WarehouseKVT.GetRackKVT(rackCode)
	if err != nil {
		fmt.Println(" GetRack error:", err)
		response.FailWithMessage(ctx, "Error loading rackkvt")
		return
	}

	fmt.Println("GetRack - result count:", len(data))
	response.OkWithData(ctx, data)
}

// Tổng tồn theo kệ
func (c *WarehouseKVTController) GetRackTotalColumnKVT(ctx *gin.Context) {

	data, err := services.WarehouseKVT.GetRackTotalColumnKVT()
	if err != nil {
		fmt.Println(" GetRackTotalColumnKVT error:", err)
		response.FailWithMessage(ctx, "Error loading GetRackTotalColumnKVT")
		return
	}

	fmt.Println("GetRackTotalColumnKVT - result count:", len(data))
	response.OkWithData(ctx, data)
}

// Tồn kho tầng
func (c *WarehouseKVTController) GetTonKhoTangKVT(ctx *gin.Context) {

	data, err := services.WarehouseKVT.GetTonKhoTangKVT()
	if err != nil {
		fmt.Println("GetTonKhoTangKVT error:", err)
		response.FailWithMessage(ctx, "Error loading tồn kho theo kệ")
		return
	}

	fmt.Println("GetTonKhoTangKVT - result count:", len(data))
	response.OkWithData(ctx, data)
}
func (c *WarehouseKVTController) GetTonKhoInRackKVT(ctx *gin.Context) {
	rackPrefix := ctx.Query("rackNo")
	if rackPrefix == "" {
		response.FailWithMessage(ctx, "rackNo query param is required")
		return
	}

	data, err := services.WarehouseKVT.GetTonKhoInRackKVT(rackPrefix)
	if err != nil {
		fmt.Println("GetTonKhoInRackKVT error:", err)
		response.FailWithMessage(ctx, "Error loading ton kho")
		return
	}

	fmt.Println("GetTonKhoInRackKVT - result count:", len(data))
	response.OkWithData(ctx, data)
}
func (c *WarehouseKVTController) GetRackInformationKVT(ctx *gin.Context) {
	var requestParams request.WarehouseRackGetRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		fmt.Println(" Validation failed:", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	fmt.Println("Request rackNo:", requestParams.RackNo)

	result, err := services.WarehouseKVT.GetRackInformationKVT(requestParams.RackNo)
	if err != nil {
		fmt.Println(" GetRackInformationKVT error:", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	fmt.Println("✅ Result length:", len(result))
	response.OkWithData(ctx, result)
}
func (c *WarehouseKVTController) Get3DayRackInformationKVT(ctx *gin.Context) {
	var requestParams request.WarehouseRackGetRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		fmt.Println(" Validation failed:", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	fmt.Println("Request rackNo:", requestParams.RackNo)

	result, err := services.WarehouseKVT.Get3DayRackInformationKVT(requestParams.RackNo)
	if err != nil {
		fmt.Println(" Get3DayRackInformationKVT error:", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	fmt.Println("✅ Result length:", len(result))
	response.OkWithData(ctx, result)
}

// //////////////////

func (c *WarehouseKVTController) GetSearchByCLBHKVT(ctx *gin.Context) {
	clbh := ctx.Query("CLBH")
	if clbh == "" {
		response.FailWithMessage(ctx, "Query parameter 'CLBH' is required")
		return
	}

	data, err := services.WarehouseKVT.SearchByCLBHKVT(clbh)
	if err != nil {
		fmt.Println("SearchByCLBHKVT error:", err)
		response.FailWithMessage(ctx, "Error retrieving CLBH data")
		return
	}

	fmt.Printf("SearchByCLBHKVT - result count: %d\n", len(data))
	response.OkWithData(ctx, data)
}

func (c *WarehouseKVTController) GetSearchN31(ctx *gin.Context) {
	po := strings.TrimSpace(ctx.Query("po"))
	orderNo := strings.TrimSpace(ctx.Query("orderNo"))

	if po == "" && orderNo == "" {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Query parameter 'po' or 'orderNo' is required")
		return
	}

	data, err := services.WarehouseKVT.SearchN31(po, orderNo)
	if err != nil {
		fmt.Println("SearchN31 error:", err)
		response.FailWithMessage(ctx, "Error retrieving N31 data")
		return
	}

	fmt.Printf("SearchN31 - result count: %d\n", len(data))
	response.OkWithData(ctx, data)
}

// func (c *WarehouseKVTController) GetSearchTotalDDBHKVT(ctx *gin.Context) {
// 	ddbh := ctx.Query("DDBH")

// 	data, err := services.Warehouse.SearchDDBH(ddbh)
// 	if err != nil {
// 		fmt.Println("SearchByDDBH error:", err)
// 		response.FailWithMessage(ctx, "Error retrieving DDBH data")
// 		return
// 	}

//		fmt.Printf("SearchByDDBH - result count: %d\n", len(data))
//		response.OkWithData(ctx, data)
//	}
var rackMap = map[string]string{
	"B0101": "KVT03",
	"B0201": "KVT03",
	"B0202": "KVT03",
	"B0203": "KVT03",
	"B0301": "KVT03",
	"B0302": "KVT03",
	"B0303": "KVT03",
	"B0401": "KVT03",
	"B0402": "KVT03",
	"B0403": "KVT03",

	// ================= KVT02 =================
	"A0101": "KVT02",
	"A0102": "KVT02",
	"A0103": "KVT02",
	"A0201": "KVT02",
	"A0202": "KVT02",
	"A0203": "KVT02",
	"A0301": "KVT02",
	"A0302": "KVT02",
	"A0303": "KVT02",
	"A0401": "KVT02",
	"A0402": "KVT02",
	"A0403": "KVT02",
	"A0501": "KVT02",
	"A0502": "KVT02",
	"A0503": "KVT02",
	"A0601": "KVT02",
	"A0602": "KVT02",
	"A0603": "KVT02",
	"A0701": "KVT02",
	"A0702": "KVT02",
	"A0703": "KVT02",
	"A0801": "KVT02",
	"A0802": "KVT02",
	"A0803": "KVT02",
	"A0901": "KVT02",
	"A0902": "KVT02",
	"A0903": "KVT02",

	"A1801": "KVT02",
	"A1802": "KVT02",
	"A1803": "KVT02",
	"A1901": "KVT02",
	"A1902": "KVT02",
	"A1903": "KVT02",
	"A2001": "KVT02",
	"A2002": "KVT02",
	"A2003": "KVT02",
	"A2101": "KVT02",
	"A2102": "KVT02",
	"A2103": "KVT02",
	"A2201": "KVT02",
	"A2202": "KVT02",
	"A2203": "KVT02",
	"A2301": "KVT02",
	"A2302": "KVT02",
	"A2303": "KVT02",
	"A2401": "KVT02",
	"A2402": "KVT02",
	"A2403": "KVT02",
	"A2501": "KVT02",
	"A2502": "KVT02",
	"A2503": "KVT02",
	"A2601": "KVT02",
	"A2602": "KVT02",
	"A2603": "KVT02",

	// ================= KVT02 =================
	"A1001": "KVT02",
	"A1002": "KVT02",
	"A1003": "KVT02",
	"A1101": "KVT02",
	"A1102": "KVT02",
	"A1103": "KVT02",
	"A1201": "KVT02",
	"A1202": "KVT02",
	"A1203": "KVT02",
	"A1301": "KVT02",
	"A1302": "KVT02",
	"A1303": "KVT02",
	"A1401": "KVT02",
	"A1402": "KVT02",
	"A1403": "KVT02",
	"A1501": "KVT02",
	"A1502": "KVT02",
	"A1503": "KVT02",
	"A1601": "KVT02",
	"A1602": "KVT02",
	"A1603": "KVT02",
	"A1701": "KVT02",
	"A1702": "KVT02",
	"A1703": "KVT02",

	"A2701": "KVT02",
	"A2702": "KVT02",
	"A2703": "KVT02",
	"A2801": "KVT02",
	"A2802": "KVT02",
	"A2803": "KVT02",
	"A2901": "KVT02",
	"A2902": "KVT02",
	"A2903": "KVT02",
	"A3001": "KVT02",
	"A3002": "KVT02",
	"A3003": "KVT02",
	"A3101": "KVT02",
	"A3102": "KVT02",
	"A3103": "KVT02",
	"A3201": "KVT02",
	"A3202": "KVT02",
	"A3203": "KVT02",
	"A3301": "KVT02",
	"A3302": "KVT02",
	"A3303": "KVT02",
	"A3401": "KVT02",
	"A3402": "KVT02",
	"A3403": "KVT02",
	"A3501": "KVT02",
	"A3502": "KVT02",
	"A3503": "KVT02",
	"A3601": "KVT02",
	"A3602": "KVT02",
	"A3603": "KVT02",
	"A3701": "KVT02",
	"A3702": "KVT02",
	"A3703": "KVT02",
	"A3801": "KVT02",
	"A3802": "KVT02",
	"A3803": "KVT02",

	"H06": "KVT02",
	"H07": "KVT02",
	"H08": "KVT02",
	"H09": "KVT02",
	"H10": "KVT02",
	"H11": "KVT02",
}
var zoneRackMap = make(map[string][]string)

func init() {
	for rack, zone := range rackMap {
		zoneRackMap[zone] = append(zoneRackMap[zone], rack)
	}
}
func getPrefix(rack string) string {
	if len(rack) >= 3 {
		return rack[:3]
	}
	return rack
}
func (c *WarehouseKVTController) GetRacksByZone(ctx *gin.Context) {

	zone := ctx.Param("zone")

	racks, ok := zoneRackMap[zone]
	if !ok {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Zone không tồn tại")
		return
	}

	// dùng map để loại duplicate
	prefixMap := make(map[string]struct{})

	for _, rack := range racks {
		prefix := getPrefix(rack)
		prefixMap[prefix] = struct{}{}
	}

	result := make([]string, 0, len(prefixMap))
	for k := range prefixMap {
		result = append(result, k)
	}

	response.OkWithData(ctx, gin.H{
		"zone":  zone,
		"racks": result,
	})
}
func (c *WarehouseKVTController) GetTempHumidityController(ctx *gin.Context) {
	var requestParams request.GetTempHumidityRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	mappedRack, ok := rackMap[requestParams.RackNo]
	if !ok {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "RackNo không tồn tại")
		return
	}

	currentDate := time.Now().Format("2006-01-02")

	result, err := services.WarehouseKVT.GetTempHumidity(mappedRack, currentDate)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}
func (c *WarehouseKVTController) GetHumidityAlertKVT01Controller(ctx *gin.Context) {
	queryDate := ctx.Query("date")
	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	if queryDate == "" {
		response.FailWithMessage(ctx, "Query parameter 'date' is required (YYYY-MM-DD)")
		return
	}

	data, err := services.WarehouseKVT.GetHumidityAlertKVT01(
		queryDate,
		startTime,
		endTime,
	)

	if err != nil {
		fmt.Println("GetHumidityAlert error:", err)
		response.FailWithMessage(ctx, "Error retrieving humidity alert data")
		return
	}

	fmt.Printf("KVT01 | date=%s | start=%s | end=%s | count=%d\n",
		queryDate, startTime, endTime, len(data))

	response.OkWithData(ctx, data)
}
func (c *WarehouseKVTController) GetHumidityAlertKVT02Controller(ctx *gin.Context) {
	queryDate := ctx.Query("date")
	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	if queryDate == "" {
		response.FailWithMessage(ctx, "Query parameter 'date' is required (YYYY-MM-DD)")
		return
	}

	data, err := services.WarehouseKVT.GetHumidityAlertKVT02(
		queryDate,
		startTime,
		endTime,
	)

	if err != nil {
		fmt.Println("GetHumidityAlert KVT02 error:", err)
		response.FailWithMessage(ctx, "Error retrieving humidity alert data")
		return
	}

	fmt.Printf("KVT02 | date=%s | start=%s | end=%s | count=%d\n",
		queryDate, startTime, endTime, len(data))

	response.OkWithData(ctx, data)
}
func (c *WarehouseKVTController) GetHumidityAlertKVT03Controller(ctx *gin.Context) {
	queryDate := ctx.Query("date")
	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	if queryDate == "" {
		response.FailWithMessage(ctx, "Query parameter 'date' is required (YYYY-MM-DD)")
		return
	}

	data, err := services.WarehouseKVT.GetHumidityAlertKVT03(
		queryDate,
		startTime,
		endTime,
	)

	if err != nil {
		fmt.Println("GetHumidityAlert KVT03 error:", err)
		response.FailWithMessage(ctx, "Error retrieving humidity alert data")
		return
	}

	fmt.Printf("KVT03 | date=%s | start=%s | end=%s | count=%d\n",
		queryDate, startTime, endTime, len(data))

	response.OkWithData(ctx, data)
}
func (c *WarehouseKVTController) GetTempHumidityByDevicesController(ctx *gin.Context) {

	result, err := services.WarehouseKVT.GetTempHumidityByDevices()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}
func (c *WarehouseKVTController) GetLatestTempHumidityLogByDevicesController(ctx *gin.Context) {
	queryDate := ctx.Query("date")
	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	result, err := services.WarehouseKVT.GetLatestTempHumidityLogByDevices(
		queryDate,
		startTime,
		endTime,
	)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}
func (c *WarehouseKVTController) GetTempHumidityLaboratoryController(ctx *gin.Context) {

	date := ctx.Query("date") // "YYYY-MM-DD", rỗng = hôm nay

	result, err := services.WarehouseKVT.GetTempHumidityLaboratory(date)
	if err != nil {
		fmt.Println("GetTempHumidityLaboratory error:", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	fmt.Printf("GetTempHumidityLaboratory | date=%s | count=%d\n", date, len(result))

	response.OkWithData(ctx, result)
}
func (c *WarehouseKVTController) Get3And180DayKVT(ctx *gin.Context) {
	result, err := services.WarehouseKVT.Get3And180DayKVT()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}
