// internal/api/controllers/inventory_controller.go
package controllers

import (
	"net/http"
	"strings"
	"web-api/internal/api/services"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/response"
	"web-api/internal/pkg/models/types"

	"github.com/gin-gonic/gin"
)

type NumberOfScansController struct {
	*BaseController
}

var NumberOfScans = &NumberOfScansController{}

func (c *NumberOfScansController) CreateInventory(ctx *gin.Context) {
	var req types.Inventory
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}
	var result string
	result, err := services.Discrepancies.CreateInventory(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"type":    "warning",
			"message": err.Error(),
		})

		return
	}

	response.OkWithData(ctx, result)
}

func (c *NumberOfScansController) CreateInventoryMonth(ctx *gin.Context) {
	var req types.Inventory

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	// Bắt buộc: barcode, USERDATE dạng 'yyyy-MM', GSBH
	if req.Barcode == "" || req.USERDATE == "" || req.GSBH == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"type":    "warning",
			"message": "Thiếu dữ liệu đầu vào",
		})
		return
	}

	result, err := services.Discrepancies.CreateInventoryByMonth(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"type":    "warning",
			"message": err.Error(),
		})
		return
	}

	response.OkWithData(ctx, result)
}

func (c *NumberOfScansController) SelectInventory(ctx *gin.Context) {
	// 1. Lấy tham số Date từ URL (?Date=...)
	Date := ctx.Query("Date")

	// 2. Lấy thêm tham số GSBH từ URL (?GSBH=...)
	// Code Vue của bạn gửi storage_GSBH, nên ở đây cần lấy đúng key đó
	GSBH := ctx.Query("GSBH")

	if Date == "" {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Date parameter is required")
		return
	}

	// Nếu GSBH trống, bạn có thể báo lỗi hoặc gán mặc định tùy nghiệp vụ
	if GSBH == "" {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "GSBH parameter is required")
		return
	}

	// 3. Truyền ĐỦ 2 tham số (Date, GSBH) vào service
	result, err := services.Discrepancies.SelectInventory(Date, GSBH)

	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *NumberOfScansController) SelectInventoryByMonth(ctx *gin.Context) {
	// Frontend gửi query:
	// ?Date=2026-04&GSBH=...&Rack=A09
	month := strings.TrimSpace(ctx.Query("Date"))
	if month == "" {
		// cho phép dùng Month nếu frontend đổi tên param
		month = strings.TrimSpace(ctx.Query("Month"))
	}

	gsbh := strings.TrimSpace(ctx.Query("GSBH"))
	rack := strings.TrimSpace(ctx.Query("Rack"))

	if month == "" {
		response.FailWithDetailed(
			ctx,
			http.StatusBadRequest,
			nil,
			"Date parameter is required (format: YYYY-MM)",
		)
		return
	}

	if gsbh == "" {
		response.FailWithDetailed(
			ctx,
			http.StatusBadRequest,
			nil,
			"GSBH parameter is required",
		)
		return
	}

	if rack == "" {
		response.FailWithDetailed(
			ctx,
			http.StatusBadRequest,
			nil,
			"Rack parameter is required",
		)
		return
	}

	result, err := services.Discrepancies.SelectInventoryByMonth(month, gsbh, rack)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *NumberOfScansController) SyncInventoryMonth(ctx *gin.Context) {
	// Rack không bắt buộc: để trống sẽ đồng bộ toàn bộ các kệ trong GSBH.
	// USERID: người đăng nhập bấm nút Đồng bộ — bắt buộc để lưu lại người thực hiện.
	month := strings.TrimSpace(ctx.Query("Date"))
	gsbh := strings.TrimSpace(ctx.Query("GSBH"))
	rack := strings.TrimSpace(ctx.Query("Rack"))
	userID := strings.TrimSpace(ctx.Query("USERID"))

	if month == "" || gsbh == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"type":    "warning",
			"message": "Date và GSBH là bắt buộc",
		})
		return
	}

	if userID == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"type":    "warning",
			"message": "Thiếu USERID của người thực hiện đồng bộ",
		})
		return
	}

	var hasSN8 string
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Không thể kiểm tra quyền SN8"})
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Không thể kiểm tra quyền SN8"})
		return
	}
	defer sqlDB.Close()

	if err := db.Raw(`
		SELECT TOP 1 '1'
		FROM BLimits
		WHERE (USERID = ? OR SUSERID = ?)
		  AND GSBH = ?
		  AND UPPER(LTRIM(RTRIM(FORMID))) = 'SN8'
	`, userID, userID, gsbh).Scan(&hasSN8).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Không thể kiểm tra quyền SN8"})
		return
	}

	switch strings.ToUpper(strings.TrimSpace(hasSN8)) {
	case "1", "Y", "YES", "TRUE":
	default:
		ctx.JSON(http.StatusForbidden, gin.H{"message": "Tài khoản không có quyền SN8"})
		return
	}

	rowsAffected, err := services.Discrepancies.SyncInventoryMonth(month, gsbh, rack, userID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"type":    "warning",
			"message": err.Error(),
		})
		return
	}

	response.OkWithData(ctx, gin.H{
		"message": "Đồng bộ thành công",
		"synced":  rowsAffected,
	})
}

// ConfirmRound2ByBarcode: kế toán quét barcode để xác nhận kiểm kê lần 2 (giống luồng quét của đơn vị).
func (c *NumberOfScansController) ConfirmRound2ByBarcode(ctx *gin.Context) {
	var req types.Inventory
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	if req.Barcode == "" || req.USERDATE == "" || req.GSBH == "" || req.USERID == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"type":    "warning",
			"message": "Thiếu Barcode, USERDATE, GSBH hoặc USERID",
		})
		return
	}

	result, err := services.Discrepancies.ConfirmRound2ByBarcode(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"type":    "warning",
			"message": err.Error(),
		})
		return
	}

	response.OkWithData(ctx, result)
}

// ConfirmRound2: kế toán xác nhận kiểm kê lần 2 cho 1 barcode (xác nhận thủ công qua nút bấm trong bảng).
func (c *NumberOfScansController) ConfirmRound2(ctx *gin.Context) {
	var req types.ConfirmRound2Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	if req.Barcode == "" || req.Date == "" || req.GSBH == "" || req.USERID == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"type":    "warning",
			"message": "Thiếu Barcode, Date, GSBH hoặc USERID",
		})
		return
	}

	rowsAffected, err := services.Discrepancies.ConfirmRound2(req.Date, req.GSBH, req.Barcode, req.USERID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"type":    "warning",
			"message": err.Error(),
		})
		return
	}

	response.OkWithData(ctx, gin.H{
		"message": "Đã xác nhận kiểm kê lần 2",
		"updated": rowsAffected,
	})
}

// GetLastSync: trả về SyncDate/SyncUserID gần nhất của Month + GSBH, để FE hiện cạnh nút Đồng bộ.
func (c *NumberOfScansController) GetLastSync(ctx *gin.Context) {
	month := strings.TrimSpace(ctx.Query("Date"))
	if month == "" {
		month = strings.TrimSpace(ctx.Query("Month"))
	}
	gsbh := strings.TrimSpace(ctx.Query("GSBH"))

	if month == "" {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Date parameter is required (format: YYYY-MM)")
		return
	}
	if gsbh == "" {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "GSBH parameter is required")
		return
	}

	result, err := services.Discrepancies.GetLastSync(month, gsbh)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *NumberOfScansController) GetRackList(ctx *gin.Context) {
	gsbh := strings.TrimSpace(ctx.Query("GSBH"))
	if gsbh == "" {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "GSBH is required")
		return
	}

	result, err := services.Discrepancies.GetRackList(gsbh)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}