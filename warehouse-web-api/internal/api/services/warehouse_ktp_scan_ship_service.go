package services

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"gorm.io/gorm"
)

type warehouseKTPScanShipCartonRow struct {
	CartonBar        string     `gorm:"column:carton_bar"`
	LPNBarCode       string     `gorm:"column:lpn_bar_code"`
	DDBH             string     `gorm:"column:ddbh"`
	CartonNo         int        `gorm:"column:carton_no"`
	XH               string     `gorm:"column:xh"`
	Qty              float64    `gorm:"column:qty"`
	SB               string     `gorm:"column:sb"`
	StandardWeight   float64    `gorm:"column:standard_weight"`
	ReceivedWeight   float64    `gorm:"column:received_weight"`
	OutgoingWeight   float64    `gorm:"column:outgoing_weight"`
	ExecuteDate      *time.Time `gorm:"column:execute_date"`
	ContainerNo      string     `gorm:"column:container_no"`
	PlateID          string     `gorm:"column:plate_id"`
	PlanCartonNo     *int       `gorm:"column:plan_carton_no"`
	PlannedContainer string     `gorm:"column:planned_container"`
	PlannedPlateID   string     `gorm:"column:planned_plate_id"`
	ResolvedBy       string     `gorm:"column:resolved_by"`
}

type warehouseKTPScanShipPending struct {
	CartonBar      string
	OutgoingWeight float64
	WeightOverride bool
}

type warehouseKTPScanShipEvaluated struct {
	row     warehouseKTPScanShipCartonRow
	request warehouseKTPScanShipPending
	result  types.WarehouseKTPScanShipCarton
}

func (s *WarehouseKTPService) GetScanShipCarton(
	barcode string,
	containerNo string,
	plateID string,
	outgoingWeight *float64,
	weightEnabled bool,
) (types.WarehouseKTPScanShipCarton, error) {
	barcode = strings.TrimSpace(barcode)
	if barcode == "" {
		return types.WarehouseKTPScanShipCarton{}, ErrWarehouseKTPCartonNotFound
	}

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return types.WarehouseKTPScanShipCarton{}, err
	}
	dbInstance, err := db.DB()
	if err != nil {
		return types.WarehouseKTPScanShipCarton{}, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	row, err := queryWarehouseKTPScanShipCarton(db, barcode)
	if err != nil {
		return types.WarehouseKTPScanShipCarton{}, err
	}
	result := evaluateWarehouseKTPScanShipCarton(row, barcode, containerNo, plateID, outgoingWeight, weightEnabled)
	result.Details, err = queryWarehouseKTPScanShipDetails(db, result.DDBH, result.XH)
	return result, err
}

func (s *WarehouseKTPService) SubmitScanShipCartons(
	params request.WarehouseKTPScanShipSubmitRequest,
	createdBy string,
) (types.WarehouseKTPScanShipSubmitResult, error) {
	params.ContainerNo = strings.TrimSpace(params.ContainerNo)
	params.PlateID = strings.TrimSpace(params.PlateID)
	createdBy = strings.TrimSpace(createdBy)
	if createdBy == "" {
		createdBy = "WEB_KTP"
	}
	if params.ContainerNo == "" || params.PlateID == "" {
		return types.WarehouseKTPScanShipSubmitResult{}, fmt.Errorf("%w: thiếu số container hoặc biển số xe", ErrWarehouseKTPScanShipRejected)
	}

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return types.WarehouseKTPScanShipSubmitResult{}, err
	}
	dbInstance, err := db.DB()
	if err != nil {
		return types.WarehouseKTPScanShipSubmitResult{}, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	result := types.WarehouseKTPScanShipSubmitResult{Cartons: make([]types.WarehouseKTPScanShipCarton, 0, len(params.Cartons))}
	err = db.Transaction(func(tx *gorm.DB) error {
		pending := make([]warehouseKTPScanShipPending, 0, len(params.Cartons))
		seen := make(map[string]struct{}, len(params.Cartons))
		for _, item := range params.Cartons {
			row, resolveErr := queryWarehouseKTPScanShipCarton(tx, strings.TrimSpace(item.CartonBar))
			if resolveErr != nil {
				return resolveErr
			}
			key := strings.ToUpper(strings.TrimSpace(row.CartonBar))
			if _, exists := seen[key]; exists {
				return fmt.Errorf("%w: carton %s bị trùng trong phiên", ErrWarehouseKTPScanShipRejected, row.CartonBar)
			}
			seen[key] = struct{}{}
			pending = append(pending, warehouseKTPScanShipPending{
				CartonBar:      row.CartonBar,
				OutgoingWeight: item.OutgoingWeight,
				WeightOverride: item.WeightOverride,
			})
		}
		sort.Slice(pending, func(i, j int) bool {
			return strings.ToUpper(pending[i].CartonBar) < strings.ToUpper(pending[j].CartonBar)
		})
		for _, item := range pending {
			if lockErr := acquireWarehouseKTPScanShipLock(tx, item.CartonBar); lockErr != nil {
				return lockErr
			}
		}

		evaluated := make([]warehouseKTPScanShipEvaluated, 0, len(pending))
		for _, item := range pending {
			row, readErr := queryWarehouseKTPScanShipCarton(tx, item.CartonBar)
			if readErr != nil {
				return readErr
			}
			weight := &item.OutgoingWeight
			validation := evaluateWarehouseKTPScanShipCarton(
				row,
				item.CartonBar,
				params.ContainerNo,
				params.PlateID,
				weight,
				true,
			)
			if !validation.EligibilityPassed || !validation.PlanPassed {
				return fmt.Errorf("%w: %s %s (%s)", ErrWarehouseKTPScanShipRejected, validation.EligibilityMessage, validation.PlanMessage, validation.CartonBar)
			}
			if !validation.WeightCheck.Passed && !(validation.WeightCheck.RequiresOverride && item.WeightOverride) {
				return fmt.Errorf("%w: %s (%s)", ErrWarehouseKTPScanShipRejected, validation.WeightCheck.Message, validation.CartonBar)
			}
			evaluated = append(evaluated, warehouseKTPScanShipEvaluated{row: row, request: item, result: validation})
		}

		for _, item := range evaluated {
			outgoingWeight := item.request.OutgoingWeight
			update := tx.Exec(`
UPDATE YWCP
SET EXEDATE = GETDATE(),
	SB = '3',
	CON_NO = ?,
	BIEN_SO = ?,
	ogw = ?
WHERE CARTONBAR = ?
	AND LTRIM(RTRIM(ISNULL(SB, ''))) = '1'
	AND EXEDATE IS NULL`,
				params.ContainerNo,
				params.PlateID,
				outgoingWeight,
				item.result.CartonBar,
			)
			if update.Error != nil {
				return update.Error
			}
			if update.RowsAffected != 1 {
				return ErrWarehouseKTPScanShipUpdateFailed
			}

			savedRow, readErr := queryWarehouseKTPScanShipCarton(tx, item.result.CartonBar)
			if readErr != nil {
				return readErr
			}
			if strings.TrimSpace(savedRow.SB) != "3" || savedRow.ExecuteDate == nil ||
				!strings.EqualFold(strings.TrimSpace(savedRow.ContainerNo), params.ContainerNo) ||
				!strings.EqualFold(strings.TrimSpace(savedRow.PlateID), params.PlateID) ||
				math.Abs(savedRow.OutgoingWeight-outgoingWeight) > 0.001 {
				return ErrWarehouseKTPScanShipUpdateFailed
			}
			saved := evaluateWarehouseKTPScanShipCarton(
				savedRow,
				item.result.ScannedCode,
				params.ContainerNo,
				params.PlateID,
				&outgoingWeight,
				true,
			)
			saved.CanScan = true
			saved.EligibilityMessage = "Scan Ship thành công"
			saved.Details, readErr = queryWarehouseKTPScanShipDetails(tx, saved.DDBH, saved.XH)
			if readErr != nil {
				return readErr
			}
			if warehouseKTPScanShipNeedsWeightOverrideLog(item.request, item.result) {
				if logErr := insertWarehouseKTPScanShipWeightOverrideLog(
					item.result,
					params,
					createdBy,
				); logErr != nil {
					return fmt.Errorf("không thể lưu log xác nhận lệch cân Scan Ship vào LYS_WEB.dbo.KTP_Scan_logs: %w", logErr)
				}
			}
			result.Cartons = append(result.Cartons, saved)
			result.UpdatedCount++
		}
		return nil
	})
	return result, err
}

func warehouseKTPScanShipNeedsWeightOverrideLog(
	request warehouseKTPScanShipPending,
	result types.WarehouseKTPScanShipCarton,
) bool {
	return request.WeightOverride &&
		result.WeightCheck.RequiresOverride &&
		!result.WeightCheck.Passed
}

func insertWarehouseKTPScanShipWeightOverrideLog(
	result types.WarehouseKTPScanShipCarton,
	params request.WarehouseKTPScanShipSubmitRequest,
	createdBy string,
) error {
	weightCheck := result.WeightCheck
	if weightCheck.ReferenceWeight == nil || weightCheck.OutgoingWeight == nil || weightCheck.Difference == nil {
		return errors.New("thiếu dữ liệu cân để ghi log Scan Ship")
	}

	db, err := database.LYS_WEB_Connection()
	if err != nil {
		return err
	}
	dbInstance, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get LYS_WEB DB instance: %w", err)
	}
	defer dbInstance.Close()

	return db.Exec(`
INSERT INTO dbo.KTP_Scan_logs
(
	CartonBar,
	OrderNo,
	ApprovedBy,
	ScanUserID,
	StandardWeight,
	ActualWeight,
	Gram,
	Tolerance,
	WeightDifference,
	Warehouse,
	ProductionLine,
	Storage,
	Location,
	ScanMode
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'A', LEFT(CONCAT('SCAN_SHIP_', ?), 20), LEFT(?, 20), NULLIF(LEFT(?, 30), ''), 'carton')`,
		result.CartonBar,
		result.DDBH,
		createdBy,
		createdBy,
		*weightCheck.ReferenceWeight,
		*weightCheck.OutgoingWeight,
		weightCheck.Tolerance,
		weightCheck.Tolerance,
		*weightCheck.Difference,
		strings.ToUpper(strings.TrimSpace(weightCheck.Reference)),
		params.ContainerNo,
		params.PlateID,
	).Error
}

func queryWarehouseKTPScanShipCarton(db *gorm.DB, barcode string) (warehouseKTPScanShipCartonRow, error) {
	row, err := queryWarehouseKTPScanShipCartonByCode(db, "cp.CARTONBAR = ?", barcode, "cartonBar")
	if err != nil || strings.TrimSpace(row.CartonBar) != "" {
		return row, err
	}
	lpnCode := strings.TrimSpace(barcode)
	if len(lpnCode) > 15 {
		lpnCode = lpnCode[len(lpnCode)-15:]
	}
	row, err = queryWarehouseKTPScanShipCartonByCode(db, "cp.LPNBarCode = ?", lpnCode, "lpnBarcode")
	if err != nil {
		return row, err
	}
	if strings.TrimSpace(row.CartonBar) == "" {
		return row, ErrWarehouseKTPCartonNotFound
	}
	return row, nil
}

func queryWarehouseKTPScanShipCartonByCode(
	db *gorm.DB,
	lookupCondition string,
	lookupValue string,
	resolvedBy string,
) (warehouseKTPScanShipCartonRow, error) {
	query := `
SELECT TOP 1
	LTRIM(RTRIM(cp.CARTONBAR)) AS carton_bar,
	LTRIM(RTRIM(ISNULL(cp.LPNBarCode, ''))) AS lpn_bar_code,
	LTRIM(RTRIM(ISNULL(cp.DDBH, ''))) AS ddbh,
	ISNULL(cp.CARTONNO, 0) AS carton_no,
	LTRIM(RTRIM(ISNULL(cp.XH, ''))) AS xh,
	CONVERT(float, ISNULL(cp.Qty, 0)) AS qty,
	LTRIM(RTRIM(ISNULL(cp.SB, ''))) AS sb,
	CONVERT(float, ISNULL(cp.sgw, 0)) AS standard_weight,
	CONVERT(float, ISNULL(cp.rgw, 0)) AS received_weight,
	CONVERT(float, ISNULL(cp.ogw, 0)) AS outgoing_weight,
	cp.EXEDATE AS execute_date,
	LTRIM(RTRIM(ISNULL(cp.CON_NO, ''))) AS container_no,
	LTRIM(RTRIM(ISNULL(cp.BIEN_SO, ''))) AS plate_id,
	ship_plan.cartonno AS plan_carton_no,
	LTRIM(RTRIM(ISNULL(ship_plan.Con_no, ''))) AS planned_container,
	LTRIM(RTRIM(ISNULL(ship_plan.PlateID, ''))) AS planned_plate_id,
	? AS resolved_by
FROM YWCP cp
OUTER APPLY (
	SELECT TOP 1 planned.cartonno, planned.Con_no, planned.PlateID
	FROM ShippingPlanCarton planned
	WHERE planned.ry = cp.DDBH
		AND planned.cartonno = cp.CARTONNO
) ship_plan
WHERE ` + lookupCondition

	var row warehouseKTPScanShipCartonRow
	if err := db.Raw(query, resolvedBy, lookupValue).Scan(&row).Error; err != nil {
		return row, err
	}
	return row, nil
}

func evaluateWarehouseKTPScanShipCarton(
	row warehouseKTPScanShipCartonRow,
	scannedCode string,
	containerNo string,
	plateID string,
	outgoingWeight *float64,
	weightEnabled bool,
) types.WarehouseKTPScanShipCarton {
	result := types.WarehouseKTPScanShipCarton{
		ScannedCode:      strings.TrimSpace(scannedCode),
		ResolvedBy:       strings.TrimSpace(row.ResolvedBy),
		CartonBar:        strings.TrimSpace(row.CartonBar),
		LPNBarCode:       strings.TrimSpace(row.LPNBarCode),
		DDBH:             strings.TrimSpace(row.DDBH),
		CartonNo:         row.CartonNo,
		XH:               strings.TrimSpace(row.XH),
		Qty:              row.Qty,
		SB:               strings.TrimSpace(row.SB),
		StandardWeight:   row.StandardWeight,
		ReceivedWeight:   row.ReceivedWeight,
		OutgoingWeight:   row.OutgoingWeight,
		ExecuteDate:      row.ExecuteDate,
		ContainerNo:      strings.TrimSpace(row.ContainerNo),
		PlateID:          strings.TrimSpace(row.PlateID),
		HasShippingPlan:  row.PlanCartonNo != nil,
		PlannedContainer: strings.TrimSpace(row.PlannedContainer),
		PlannedPlateID:   strings.TrimSpace(row.PlannedPlateID),
		Details:          make([]types.WarehouseKTPScanOutDetail, 0),
	}
	result.EligibilityPassed, result.EligibilityMessage = evaluateWarehouseKTPScanShipEligibility(result)
	result.PlanPassed, result.PlanMessage = evaluateWarehouseKTPScanShipPlan(result, containerNo, plateID)
	result.WeightCheck = evaluateWarehouseKTPScanShipWeight(row, outgoingWeight, weightEnabled)
	result.CanScan = result.EligibilityPassed && result.PlanPassed && result.WeightCheck.Passed
	return result
}

func evaluateWarehouseKTPScanShipEligibility(carton types.WarehouseKTPScanShipCarton) (bool, string) {
	if carton.CartonBar == "" {
		return false, "Không tìm thấy carton"
	}
	if carton.DDBH == "" {
		return false, "Carton chưa có RY"
	}
	if carton.ExecuteDate != nil {
		return false, "Carton đã có EXEDATE"
	}
	if carton.SB != "1" {
		return false, fmt.Sprintf("Carton không ở trong kho thành phẩm (SB = %s)", carton.SB)
	}
	return true, "Carton hợp lệ để Scan Ship"
}

func evaluateWarehouseKTPScanShipPlan(
	carton types.WarehouseKTPScanShipCarton,
	containerNo string,
	plateID string,
) (bool, string) {
	if strings.TrimSpace(containerNo) == "" || strings.TrimSpace(plateID) == "" {
		return false, "Thiếu số container hoặc biển số xe"
	}
	if !carton.HasShippingPlan {
		return true, "Carton chưa có kế hoạch xe/container"
	}
	if !strings.EqualFold(strings.TrimSpace(carton.PlannedContainer), strings.TrimSpace(containerNo)) ||
		!strings.EqualFold(strings.TrimSpace(carton.PlannedPlateID), strings.TrimSpace(plateID)) {
		return false, "Sai biển số xe hoặc số container so với Shipping Plan"
	}
	return true, "Đúng Shipping Plan"
}

func evaluateWarehouseKTPScanShipWeight(
	row warehouseKTPScanShipCartonRow,
	outgoingWeight *float64,
	weightEnabled bool,
) types.WarehouseKTPScanShipWeightCheck {
	check := types.WarehouseKTPScanShipWeightCheck{Enabled: weightEnabled, Passed: !weightEnabled}
	if !weightEnabled {
		check.Message = "Không sử dụng cân"
		return check
	}
	check.OutgoingWeight = outgoingWeight
	if outgoingWeight == nil || *outgoingWeight <= 0 {
		check.Message = "Chưa đọc được trọng lượng xuất từ cân"
		return check
	}
	reference := row.ReceivedWeight
	check.Reference = "rgw"
	check.Tolerance = 0.05
	if reference <= 0 {
		reference = row.StandardWeight
		check.Reference = "sgw"
		check.Tolerance = 0.5
	}
	if reference <= 0 {
		check.Message = "Carton chưa có trọng lượng tham chiếu rgw/sgw"
		return check
	}
	check.ReferenceWeight = &reference
	difference := reference - *outgoingWeight
	check.Difference = &difference
	if math.Abs(difference) > check.Tolerance {
		check.RequiresOverride = true
		check.Message = fmt.Sprintf("Trọng lượng lệch quá %.2f kg so với %s", check.Tolerance, check.Reference)
		return check
	}
	check.Passed = true
	check.Message = "Trọng lượng đạt"
	return check
}

func queryWarehouseKTPServerDate(db *gorm.DB) (time.Time, error) {
	var row struct {
		ServerDate time.Time `gorm:"column:server_date"`
	}
	if err := db.Raw(`SELECT CONVERT(date, GETDATE()) AS server_date`).Scan(&row).Error; err != nil {
		return time.Time{}, err
	}
	return row.ServerDate, nil
}

func acquireWarehouseKTPScanShipLock(tx *gorm.DB, cartonBar string) error {
	resource := "warehouse-ktp-scan-ship:" + strings.ToUpper(strings.TrimSpace(cartonBar))
	return tx.Exec(`
DECLARE @lock_result INT;
EXEC @lock_result = sys.sp_getapplock
	@Resource = ?,
	@LockMode = 'Exclusive',
	@LockOwner = 'Transaction',
	@LockTimeout = 15000;
IF @lock_result < 0
	THROW 50001, 'Khong the khoa carton de Scan Ship.', 1;`, resource).Error
}
