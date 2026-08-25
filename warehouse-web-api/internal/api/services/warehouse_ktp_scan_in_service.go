package services

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"gorm.io/gorm"
)

const (
	scanInStatusPassed    = "passed"
	scanInStatusFailed    = "failed"
	scanInStatusNoData    = "no_data"
	scanInModeCarton      = "carton"
	scanInModeCartonInbox = "carton_inbox"
)

var errWarehouseKTPRubberMaterialUnavailable = errors.New("thiếu database LIY_DD cho định mức vật tư cao su")

type warehouseKTPScanInCartonRow struct {
	CartonBar   string     `gorm:"column:carton_bar"`
	LPNBarCode  *string    `gorm:"column:lpn_bar_code"`
	DDBH        string     `gorm:"column:ddbh"`
	OrderNo     *string    `gorm:"column:order_no"`
	CartonNo    *int       `gorm:"column:carton_no"`
	XH          *string    `gorm:"column:xh"`
	Size        *string    `gorm:"column:size"`
	Qty         *float64   `gorm:"column:qty"`
	PackageQty  *float64   `gorm:"column:package_qty"`
	SB          *string    `gorm:"column:sb"`
	CKBH        *string    `gorm:"column:ckbh"`
	KCBH        *string    `gorm:"column:kcbh"`
	KVBH        *string    `gorm:"column:kvbh"`
	INCS        *int       `gorm:"column:incs"`
	INDate      *time.Time `gorm:"column:in_date"`
	LastInDate  *time.Time `gorm:"column:last_in_date"`
	REDate      *time.Time `gorm:"column:re_date"`
	RECS        *int       `gorm:"column:recs"`
	SGW         *float64   `gorm:"column:sgw"`
	RGW         *float64   `gorm:"column:rgw"`
	Article     *string    `gorm:"column:article"`
	XieMing     *string    `gorm:"column:xie_ming"`
	Country     *string    `gorm:"column:country"`
	KHPO        *string    `gorm:"column:khpo"`
	DepNo       *string    `gorm:"column:dep_no"`
	DepName     *string    `gorm:"column:dep_name"`
	INUser      *string    `gorm:"column:in_user"`
	Factory     *string    `gorm:"column:factory"`
	FactoryType *string    `gorm:"column:factory_type"`
	OrderFlag   *string    `gorm:"column:order_flag"`
	Gram        *float64   `gorm:"column:gram"`
	ResolvedBy  string     `gorm:"column:resolved_by"`
}

type warehouseKTPScanInProductionRow struct {
	CompletedQty   float64 `gorm:"column:completed_qty"`
	AlreadyInQty   float64 `gorm:"column:already_in_qty"`
	SessionQty     float64 `gorm:"column:session_qty"`
	MixedOrderRows int     `gorm:"column:mixed_order_rows"`
}

type warehouseKTPScanInMaterialRow struct {
	MaterialOrder        *string  `gorm:"column:material_order"`
	OrderPairs           *float64 `gorm:"column:order_pairs"`
	CoveredPairs         *float64 `gorm:"column:covered_pairs"`
	AlreadyInQty         float64  `gorm:"column:already_in_qty"`
	SessionQty           float64  `gorm:"column:session_qty"`
	RequiredMaterialRows int      `gorm:"column:required_material_rows"`
	IssuedMaterialRows   int      `gorm:"column:issued_material_rows"`
	MissingMaterialRows  int      `gorm:"column:missing_material_rows"`
	LimitingFactory      *string  `gorm:"column:limiting_factory"`
}

type warehouseKTPScanInInboxRow struct {
	InboxBar    *string `gorm:"column:inbox_bar"`
	Size        *string `gorm:"column:size"`
	ExpectedQty int     `gorm:"column:expected_qty"`
}

func (s *WarehouseKTPService) GetScanInCarton(
	barcode string,
	actualWeight *float64,
	sessionCartons []string,
	scanMode string,
	inboxScans []string,
) (types.WarehouseKTPScanInCartonResult, error) {
	normalizedBarcode := strings.TrimSpace(barcode)
	if normalizedBarcode == "" {
		return types.WarehouseKTPScanInCartonResult{}, ErrWarehouseKTPCartonNotFound
	}

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return types.WarehouseKTPScanInCartonResult{}, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return types.WarehouseKTPScanInCartonResult{}, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()
	return getWarehouseKTPScanInCartonWithDB(
		db,
		normalizedBarcode,
		actualWeight,
		sessionCartons,
		scanMode,
		inboxScans,
		true,
	)
}

func getWarehouseKTPScanInCartonWithDB(
	db *gorm.DB,
	barcode string,
	actualWeight *float64,
	sessionCartons []string,
	scanMode string,
	inboxScans []string,
	parallel bool,
) (types.WarehouseKTPScanInCartonResult, error) {
	row, err := queryWarehouseKTPScanInCarton(db, barcode)
	if err != nil {
		return types.WarehouseKTPScanInCartonResult{}, err
	}
	return evaluateWarehouseKTPScanInCartonWithDB(
		db,
		row,
		barcode,
		actualWeight,
		sessionCartons,
		scanMode,
		inboxScans,
		parallel,
	)
}

func evaluateWarehouseKTPScanInCartonWithDB(
	db *gorm.DB,
	row warehouseKTPScanInCartonRow,
	barcode string,
	actualWeight *float64,
	sessionCartons []string,
	scanMode string,
	inboxScans []string,
	parallel bool,
) (types.WarehouseKTPScanInCartonResult, error) {
	carton := mapWarehouseKTPScanInCarton(row, barcode)
	sessionCartons = compactScanInCartons(sessionCartons)
	scanMode = normalizeWarehouseKTPScanInMode(scanMode)

	var (
		productionRow warehouseKTPScanInProductionRow
		productionErr error
		materialRow   warehouseKTPScanInMaterialRow
		materialErr   error
		inboxRows     []warehouseKTPScanInInboxRow
		inboxErr      error
		checks        sync.WaitGroup
	)
	if parallel {
		checkCount := 2
		if scanMode == scanInModeCartonInbox {
			checkCount++
		}
		checks.Add(checkCount)
		go func() {
			defer checks.Done()
			productionRow, productionErr = queryWarehouseKTPScanInProduction(db, carton, sessionCartons)
		}()
		go func() {
			defer checks.Done()
			materialRow, materialErr = queryWarehouseKTPScanInMaterial(db, carton, sessionCartons)
		}()
		if scanMode == scanInModeCartonInbox {
			go func() {
				defer checks.Done()
				inboxRows, inboxErr = queryWarehouseKTPScanInInboxes(db, carton.CartonBar)
			}()
		}
		checks.Wait()
	} else {
		productionRow, productionErr = queryWarehouseKTPScanInProduction(db, carton, sessionCartons)
		materialRow, materialErr = queryWarehouseKTPScanInMaterial(db, carton, sessionCartons)
		if scanMode == scanInModeCartonInbox {
			inboxRows, inboxErr = queryWarehouseKTPScanInInboxes(db, carton.CartonBar)
		}
	}

	eligibility := evaluateWarehouseKTPScanInEligibility(carton, sessionCartons, productionRow)
	weightRule := evaluateWarehouseKTPScanInWeight(carton, row.Gram, actualWeight)
	production := evaluateWarehouseKTPScanInProduction(carton, productionRow, productionErr)
	material := evaluateWarehouseKTPScanInMaterial(carton, productionRow, materialRow, materialErr)
	inbox := evaluateWarehouseKTPScanInInboxes(scanMode, inboxRows, inboxScans, inboxErr)

	return types.WarehouseKTPScanInCartonResult{
		Carton:          carton,
		Eligibility:     eligibility,
		WeightRule:      weightRule,
		ProductionCheck: production,
		MaterialCheck:   material,
		InboxCheck:      inbox,
		CanScanIn:       eligibility.Passed && weightRule.Passed && production.Passed && material.Passed && inbox.Passed,
		ReadOnly:        true,
	}, nil
}

func (s *WarehouseKTPService) SubmitScanInCarton(
	params request.WarehouseKTPScanInSubmitRequest,
	createdBy string,
) (types.WarehouseKTPScanInCartonResult, error) {
	params.CartonBar = strings.TrimSpace(params.CartonBar)
	params.Warehouse = "A"
	params.ProductionLine = strings.TrimSpace(params.ProductionLine)
	params.Storage = strings.ToUpper(strings.TrimSpace(params.Storage))
	params.Location = strings.ToUpper(strings.TrimSpace(params.Location))
	params.ToSX = "A"
	params.WeightOverrideUserID = strings.TrimSpace(params.WeightOverrideUserID)
	createdBy = strings.TrimSpace(createdBy)
	if createdBy == "" {
		createdBy = strings.TrimSpace(params.InUser)
	}
	if createdBy == "" {
		createdBy = "WEB_KTP"
	}

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return types.WarehouseKTPScanInCartonResult{}, err
	}
	dbInstance, err := db.DB()
	if err != nil {
		return types.WarehouseKTPScanInCartonResult{}, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	actualWeight := params.ActualWeight
	var validation types.WarehouseKTPScanInCartonResult
	err = db.Transaction(func(tx *gorm.DB) error {
		lockedRow, queryErr := queryWarehouseKTPScanInCarton(tx, params.CartonBar)
		if queryErr != nil {
			return queryErr
		}
		lockKey := firstNonEmptyScanInValue(
			trimScanInString(lockedRow.OrderNo),
			strings.TrimSpace(lockedRow.DDBH),
			strings.TrimSpace(lockedRow.CartonBar),
		)
		if lockErr := acquireWarehouseKTPScanInLock(tx, lockKey); lockErr != nil {
			return lockErr
		}
		lockedRow, queryErr = queryWarehouseKTPScanInCarton(tx, params.CartonBar)
		if queryErr != nil {
			return queryErr
		}
		if setupErr := validateWarehouseKTPScanInSetup(tx, params); setupErr != nil {
			return fmt.Errorf("%w: %s", ErrWarehouseKTPScanInRejected, setupErr.Error())
		}

		validation, queryErr = evaluateWarehouseKTPScanInCartonWithDB(
			tx,
			lockedRow,
			params.CartonBar,
			&actualWeight,
			params.SessionCartons,
			params.ScanMode,
			params.InboxScans,
			false,
		)
		if queryErr != nil {
			return queryErr
		}
		if overrideErr := applyWarehouseKTPScanInWeightOverride(&validation, params, createdBy); overrideErr != nil {
			return overrideErr
		}
		if !validation.CanScanIn {
			return fmt.Errorf("%w: %s", ErrWarehouseKTPScanInRejected, firstFailedScanInMessage(validation))
		}

		update := tx.Exec(`
UPDATE YWCP
SET
	SB = '1',
	CKBH = ?,
	DepNo = CASE WHEN INDATE IS NULL THEN ? ELSE DepNo END,
	KCBH = ?,
	KVBH = ?,
	INCS = CASE WHEN INDATE IS NULL THEN 1 ELSE ISNULL(INCS, 0) + 1 END,
	INDATE = CASE WHEN INDATE IS NULL THEN GETDATE() ELSE INDATE END,
	LastInDate = GETDATE(),
	rgw = ?,
	INUser = ?,
	REDATE = CASE WHEN SB IN ('2', '4') THEN GETDATE() ELSE REDATE END,
	RECS = CASE WHEN SB IN ('2', '4') THEN ISNULL(RECS, 0) + 1 ELSE RECS END
WHERE CARTONBAR = ?
	AND (SB IS NULL OR SB IN ('', '2', '4'))`,
			params.Warehouse,
			params.ProductionLine,
			params.Storage,
			params.Location,
			actualWeight,
			createdBy,
			validation.Carton.CartonBar,
		)
		if update.Error != nil {
			return update.Error
		}
		if update.RowsAffected != 1 {
			return ErrWarehouseKTPScanInUpdateFailed
		}

		savedRow, queryErr := queryWarehouseKTPScanInCarton(tx, validation.Carton.CartonBar)
		if queryErr != nil {
			return queryErr
		}
		if verifyErr := verifyWarehouseKTPScanInSaved(
			lockedRow,
			savedRow,
			params,
			createdBy,
			actualWeight,
		); verifyErr != nil {
			return verifyErr
		}
		validation.Carton = mapWarehouseKTPScanInCarton(savedRow, params.CartonBar)
		validation.ReadOnly = false
		if validation.WeightRule.Overridden {
			if logErr := insertWarehouseKTPWeightOverrideLog(validation, params, createdBy); logErr != nil {
				return fmt.Errorf("không thể lưu log xác nhận lệch cân vào LYS_WEB.dbo.KTP_Scan_logs: %w", logErr)
			}
		}
		return nil
	})
	return validation, err
}

func applyWarehouseKTPScanInWeightOverride(
	result *types.WarehouseKTPScanInCartonResult,
	params request.WarehouseKTPScanInSubmitRequest,
	createdBy string,
) error {
	if !params.WeightOverride {
		return nil
	}
	if result.WeightRule.Status != scanInStatusFailed {
		return fmt.Errorf("%w: chỉ được mở khóa khi trọng lượng nằm ngoài giới hạn", ErrWarehouseKTPScanInRejected)
	}
	if params.WeightOverrideUserID == "" || !strings.EqualFold(params.WeightOverrideUserID, strings.TrimSpace(createdBy)) {
		return fmt.Errorf("%w: userid xác nhận mở khóa phải trùng với userid đang đăng nhập", ErrWarehouseKTPScanInRejected)
	}

	result.WeightRule.Status = scanInStatusPassed
	result.WeightRule.Passed = true
	result.WeightRule.Overridden = true
	result.WeightRule.OverrideUserID = strings.TrimSpace(createdBy)
	result.WeightRule.Message = fmt.Sprintf("Sai lệch cân đã được %s xác nhận mở khóa", strings.TrimSpace(createdBy))
	result.CanScanIn = result.Eligibility.Passed &&
		result.WeightRule.Passed &&
		result.ProductionCheck.Passed &&
		result.MaterialCheck.Passed &&
		result.InboxCheck.Passed
	return nil
}

func insertWarehouseKTPWeightOverrideLog(
	result types.WarehouseKTPScanInCartonResult,
	params request.WarehouseKTPScanInSubmitRequest,
	createdBy string,
) error {
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
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULLIF(?, ''), ?)`,
		result.Carton.CartonBar,
		result.Carton.OrderNo,
		result.WeightRule.OverrideUserID,
		strings.TrimSpace(createdBy),
		scanInFloat(result.WeightRule.StandardWeight),
		scanInFloat(result.WeightRule.ActualWeight),
		scanInFloat(result.WeightRule.Gram),
		scanInFloat(result.WeightRule.Tolerance),
		scanInFloat(result.WeightRule.Difference),
		params.Warehouse,
		params.ProductionLine,
		params.Storage,
		params.Location,
		normalizeWarehouseKTPScanInMode(params.ScanMode),
	).Error
}

func acquireWarehouseKTPScanInLock(tx *gorm.DB, key string) error {
	resource := "warehouse-ktp-scan-in:" + strings.ToUpper(strings.TrimSpace(key))
	return tx.Exec(`
DECLARE @lock_result INT;
EXEC @lock_result = sys.sp_getapplock
	@Resource = ?,
	@LockMode = 'Exclusive',
	@LockOwner = 'Transaction',
	@LockTimeout = 15000;
IF @lock_result < 0
	THROW 50001, 'Khong the khoa don hang de scan nhap kho.', 1;`, resource).Error
}

func validateWarehouseKTPScanInSetup(tx *gorm.DB, params request.WarehouseKTPScanInSubmitRequest) error {
	var setup struct {
		WarehouseCount      int `gorm:"column:warehouse_count"`
		ProductionLineCount int `gorm:"column:production_line_count"`
	}
	err := tx.Raw(`
SELECT
	(SELECT COUNT(*) FROM KCCKCP WHERE CKBH = ?) AS warehouse_count,
	(SELECT COUNT(*)
	 FROM BDepartment
	 WHERE ID = ?
	   AND ISNULL(GSBH, '') = 'VDH'
	   AND ISNULL(UserID, '') = '340205'
	   AND ProYN = 1
	   AND ISNULL(YN, '') = '1'
	   AND ISNULL(GXLB, '') = 'A') AS production_line_count`,
		params.Warehouse,
		params.ProductionLine,
	).Scan(&setup).Error
	if err != nil {
		return err
	}
	if setup.WarehouseCount == 0 {
		return fmt.Errorf("Warehouse %s không tồn tại trong KCCKCP", params.Warehouse)
	}
	if setup.ProductionLineCount == 0 {
		return fmt.Errorf("Production line %s không hợp lệ tại VDH", params.ProductionLine)
	}
	return nil
}

func verifyWarehouseKTPScanInSaved(
	previous warehouseKTPScanInCartonRow,
	saved warehouseKTPScanInCartonRow,
	params request.WarehouseKTPScanInSubmitRequest,
	createdBy string,
	actualWeight float64,
) error {
	expectedINCS := 1
	if previous.INDate != nil {
		expectedINCS = scanInInt(previous.INCS) + 1
	}
	expectedDepNo := params.ProductionLine
	if previous.INDate != nil {
		expectedDepNo = trimScanInString(previous.DepNo)
	}
	valid := strings.TrimSpace(saved.CartonBar) == strings.TrimSpace(previous.CartonBar) &&
		trimScanInString(saved.SB) == "1" &&
		strings.EqualFold(trimScanInString(saved.CKBH), params.Warehouse) &&
		strings.EqualFold(trimScanInString(saved.KCBH), params.Storage) &&
		strings.EqualFold(trimScanInString(saved.KVBH), params.Location) &&
		strings.EqualFold(trimScanInString(saved.DepNo), expectedDepNo) &&
		scanInInt(saved.INCS) == expectedINCS &&
		saved.INDate != nil &&
		saved.LastInDate != nil &&
		saved.RGW != nil && math.Abs(*saved.RGW-actualWeight) < 0.000001 &&
		strings.EqualFold(trimScanInString(saved.INUser), createdBy)
	if !valid {
		return ErrWarehouseKTPScanInUpdateFailed
	}
	if previousSB := trimScanInString(previous.SB); previousSB == "2" || previousSB == "4" {
		if saved.REDate == nil || scanInInt(saved.RECS) != scanInInt(previous.RECS)+1 {
			return ErrWarehouseKTPScanInUpdateFailed
		}
	}
	return nil
}

func firstNonEmptyScanInValue(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return "unknown"
}

func firstFailedScanInMessage(result types.WarehouseKTPScanInCartonResult) string {
	for _, check := range []types.WarehouseKTPScanInCheck{
		result.Eligibility.WarehouseKTPScanInCheck,
		result.WeightRule.WarehouseKTPScanInCheck,
		result.ProductionCheck.WarehouseKTPScanInCheck,
		result.MaterialCheck.WarehouseKTPScanInCheck,
		result.InboxCheck.WarehouseKTPScanInCheck,
	} {
		if !check.Passed && strings.TrimSpace(check.Message) != "" {
			return check.Message
		}
	}
	return "Vui lòng quét lại carton và kiểm tra các điều kiện"
}

func (s *WarehouseKTPService) GetScanInProductionLines(factory string) ([]types.WarehouseKTPScanInProductionLine, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	const query = `
SELECT
	LTRIM(RTRIM(line.ID)) AS id,
	LTRIM(RTRIM(line.DepName)) AS name,
	LTRIM(RTRIM(line.GSBH)) AS factory
FROM BDepartment line
WHERE line.ProYN = 1
	AND ISNULL(line.YN, '') = '1'
	AND ISNULL(line.GXLB, '') = 'A'
	AND ISNULL(line.UserID, '') = '340205'
	AND (? = '' OR line.GSBH = ?)
ORDER BY line.DepName, line.ID`

	var lines []types.WarehouseKTPScanInProductionLine
	normalizedFactory := strings.TrimSpace(factory)
	if normalizedFactory == "" {
		normalizedFactory = "VDH"
	}
	if err := db.Raw(query, normalizedFactory, normalizedFactory).Scan(&lines).Error; err != nil {
		return nil, err
	}
	return lines, nil
}

func (s *WarehouseKTPService) GetScanInWarehouses() ([]types.WarehouseKTPScanInWarehouse, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	const query = `
SELECT DISTINCT LTRIM(RTRIM(CKBH)) AS code
FROM KCCKCP
WHERE ISNULL(CKBH, '') <> ''
ORDER BY code`

	var warehouses []types.WarehouseKTPScanInWarehouse
	if err := db.Raw(query).Scan(&warehouses).Error; err != nil {
		return nil, err
	}
	return warehouses, nil
}

func queryWarehouseKTPScanInCarton(db *gorm.DB, barcode string) (warehouseKTPScanInCartonRow, error) {
	row, err := queryWarehouseKTPScanInCartonByCode(db, "cp.CARTONBAR = ?", barcode, "cartonBar")
	if err != nil {
		return row, err
	}
	if strings.TrimSpace(row.CartonBar) != "" {
		return row, nil
	}

	lpnCode := barcode
	if len(lpnCode) > 15 {
		lpnCode = lpnCode[len(lpnCode)-15:]
	}
	row, err = queryWarehouseKTPScanInCartonByCode(db, "cp.LPNBarCode = ?", lpnCode, "lpnBarcode")
	if err != nil {
		return row, err
	}
	if strings.TrimSpace(row.CartonBar) == "" {
		return row, ErrWarehouseKTPCartonNotFound
	}
	return row, nil
}

func queryWarehouseKTPScanInCartonByCode(
	db *gorm.DB,
	lookupCondition string,
	lookupValue string,
	resolvedBy string,
) (warehouseKTPScanInCartonRow, error) {
	query := `
SELECT TOP 1
	LTRIM(RTRIM(cp.CARTONBAR)) AS carton_bar,
	LTRIM(RTRIM(ISNULL(cp.LPNBarCode, ''))) AS lpn_bar_code,
	LTRIM(RTRIM(cp.DDBH)) AS ddbh,
	LTRIM(RTRIM(ISNULL(ywdd.YSBH, ''))) AS order_no,
	cp.CARTONNO AS carton_no,
	LTRIM(RTRIM(ISNULL(cp.XH, ''))) AS xh,
	LTRIM(RTRIM(ISNULL(pos.DDCC, ''))) AS size,
	CONVERT(float, ISNULL(cp.Qty, 0)) AS qty,
	CONVERT(float, ISNULL(pos.Qty, cp.Qty)) AS package_qty,
	LTRIM(RTRIM(ISNULL(cp.SB, ''))) AS sb,
	LTRIM(RTRIM(ISNULL(cp.CKBH, ''))) AS ckbh,
	LTRIM(RTRIM(ISNULL(cp.KCBH, ''))) AS kcbh,
	LTRIM(RTRIM(ISNULL(cp.KVBH, ''))) AS kvbh,
	ISNULL(cp.INCS, 0) AS incs,
	cp.INDATE AS in_date,
	cp.LastInDate AS last_in_date,
	cp.REDATE AS re_date,
	ISNULL(cp.RECS, 0) AS recs,
	CONVERT(float, ISNULL(cp.sgw, 0)) AS sgw,
	CONVERT(float, cp.rgw) AS rgw,
	LTRIM(RTRIM(ISNULL(ddzl.Article, ''))) AS article,
	LTRIM(RTRIM(ISNULL(shoe.XieMing, ''))) AS xie_ming,
	LTRIM(RTRIM(COALESCE(NULLIF(country.ywsm, ''), ddzl.DDGB, ''))) AS country,
	LTRIM(RTRIM(ISNULL(ddzl.KHPO, ''))) AS khpo,
	LTRIM(RTRIM(ISNULL(cp.DepNo, ''))) AS dep_no,
	LTRIM(RTRIM(ISNULL(department.DepName, ''))) AS dep_name,
	LTRIM(RTRIM(ISNULL(cp.INUSER, ''))) AS in_user,
	LTRIM(RTRIM(ISNULL(ddzl.GSBH, ''))) AS factory,
	LTRIM(RTRIM(ISNULL(factory.SFL, ''))) AS factory_type,
	LTRIM(RTRIM(ISNULL(ddzl.Flag, ''))) AS order_flag,
	CONVERT(float,
		CASE
			WHEN LTRIM(RTRIM(ISNULL(ddzl.GSBH, ''))) IN ('VB1', 'VB7')
				THEN ROUND(ISNULL(cp.sgw, 0) * 0.03, 3)
			ELSE ISNULL(weight_data.NWeigh, 0) * 1.04
		END
	) AS gram,
	? AS resolved_by
FROM YWCP cp
LEFT JOIN YWDD ywdd ON ywdd.DDBH = cp.DDBH
LEFT JOIN DDZL ddzl ON ddzl.DDBH = ywdd.YSBH
LEFT JOIN XXZL shoe ON shoe.XieXing = ddzl.XieXing AND shoe.SheHao = ddzl.SheHao
LEFT JOIN BDepartment department ON department.ID = cp.DepNo
LEFT JOIN LBZLS country ON country.lb = '06' AND country.lbdh = ddzl.DDGB
LEFT JOIN YWBZPOS pos ON pos.DDBH = cp.DDBH AND pos.XH = cp.XH
OUTER APPLY (
	SELECT TOP 1 Bgszl.SFL
	FROM Bgszl
	WHERE Bgszl.GSDH = ddzl.GSBH
) factory
OUTER APPLY (
	SELECT TOP 1 YWXXNew.NWeigh
	FROM YWXXNew
	WHERE YWXXNew.XieXing = ddzl.XieXing
	  AND YWXXNew.SheHao = ddzl.SheHao
	  AND YWXXNew.XXCC = pos.DDCC
) weight_data
WHERE ` + lookupCondition

	var row warehouseKTPScanInCartonRow
	if err := db.Raw(query, resolvedBy, lookupValue).Scan(&row).Error; err != nil {
		return row, err
	}
	return row, nil
}

func queryWarehouseKTPScanInProduction(
	db *gorm.DB,
	carton types.WarehouseKTPScanInCarton,
	sessionCartons []string,
) (warehouseKTPScanInProductionRow, error) {
	excludedCartons := compactScanInCartons(append(append([]string{}, sessionCartons...), carton.CartonBar))
	exclusionClause, exclusionArgs := scanInCartonExclusionClause(excludedCartons)
	sessionClause, sessionArgs := scanInSessionClause(sessionCartons, carton.OrderNo, carton.Size, false)
	mixedOrderClause, mixedOrderArgs := scanInSessionClause(sessionCartons, carton.OrderNo, carton.Size, true)
	query := `
SELECT
	ISNULL((
		SELECT SUM(CONVERT(float, ISNULL(smddss.Qty, 0)) * CONVERT(float, ISNULL(smddss.okCTS, 0)))
		FROM SMDDSS smddss
		WHERE smddss.DDBH LIKE ?
		  AND smddss.GXLB = 'A'
		  AND smddss.XXCC = ?
	), 0) AS completed_qty,
	ISNULL((
		SELECT SUM(CONVERT(float, ISNULL(pos.Qty, cp.Qty)))
		FROM YWCP cp
		LEFT JOIN YWDD ywdd ON ywdd.DDBH = cp.DDBH
		LEFT JOIN YWBZPOS pos ON pos.DDBH = cp.DDBH AND pos.XH = cp.XH
		WHERE ywdd.YSBH = ?
		  AND pos.DDCC = ?
		  AND cp.LastInDate IS NOT NULL
		  AND cp.SB = '1'
		  ` + exclusionClause + `
	), 0) AS already_in_qty,
	` + sessionClause + ` AS session_qty,
	` + mixedOrderClause + ` AS mixed_order_rows`

	args := []interface{}{carton.OrderNo + "%", carton.Size, carton.OrderNo, carton.Size}
	args = append(args, exclusionArgs...)
	args = append(args, sessionArgs...)
	args = append(args, mixedOrderArgs...)

	var row warehouseKTPScanInProductionRow
	err := db.Raw(query, args...).Scan(&row).Error
	return row, err
}

func scanInCartonExclusionClause(cartons []string) (string, []interface{}) {
	if len(cartons) == 0 {
		return "", nil
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(cartons)), ",")
	args := make([]interface{}, 0, len(cartons))
	for _, carton := range cartons {
		args = append(args, carton)
	}
	return "AND cp.CARTONBAR NOT IN (" + placeholders + ")", args
}

func scanInSessionClause(
	sessionCartons []string,
	orderNo string,
	size string,
	mixedOrder bool,
) (string, []interface{}) {
	if len(sessionCartons) == 0 {
		return "CONVERT(float, 0)", nil
	}

	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(sessionCartons)), ",")
	aggregate := "SUM(CONVERT(float, ISNULL(pos.Qty, cp.Qty)))"
	orderComparison := "="
	sizeCondition := "AND pos.DDCC = ?"
	if mixedOrder {
		aggregate = "COUNT(*)"
		orderComparison = "<>"
		sizeCondition = ""
	}
	clause := `ISNULL((
		SELECT ` + aggregate + `
		FROM YWCP cp
		LEFT JOIN YWDD ywdd ON ywdd.DDBH = cp.DDBH
		LEFT JOIN YWBZPOS pos ON pos.DDBH = cp.DDBH AND pos.XH = cp.XH
		WHERE cp.CARTONBAR IN (` + placeholders + `)
		  AND ywdd.YSBH ` + orderComparison + ` ?
		  ` + sizeCondition + `
	), 0)`

	args := make([]interface{}, 0, len(sessionCartons)+2)
	for _, carton := range sessionCartons {
		args = append(args, carton)
	}
	args = append(args, orderNo)
	if !mixedOrder {
		args = append(args, size)
	}
	return clause, args
}

func queryWarehouseKTPScanInMaterial(
	db *gorm.DB,
	carton types.WarehouseKTPScanInCarton,
	sessionCartons []string,
) (warehouseKTPScanInMaterialRow, error) {
	orderNo := carton.OrderNo
	var capability struct {
		NeedsRubber int `gorm:"column:needs_rubber"`
		HasLIYDD    int `gorm:"column:has_liy_dd"`
	}
	const capabilityQuery = `
DECLARE @material_order varchar(50) = COALESCE(
	(SELECT TOP 1 LTRIM(RTRIM(DDZLTR.DDBH)) FROM DDZLTR WHERE DDZLTR.DDBH1 = ?),
	?
);
SELECT
	CASE WHEN EXISTS (
		SELECT 1
		FROM DDZL
		INNER JOIN Bgszl source_factory ON source_factory.GSDH = DDZL.GSBH
		INNER JOIN Bgszl scope_factory
			ON scope_factory.DFL = source_factory.DFL
			AND scope_factory.SFL = 'RB'
			AND ISNULL(scope_factory.tybh, '') <> 'Y'
		WHERE DDZL.DDBH = @material_order
	) THEN 1 ELSE 0 END AS needs_rubber,
	CASE WHEN DB_ID('LIY_DD') IS NOT NULL
		AND OBJECT_ID('LIY_DD.dbo.ZLZLS3') IS NOT NULL
		THEN 1 ELSE 0 END AS has_liy_dd`
	if err := db.Raw(capabilityQuery, orderNo, orderNo).Scan(&capability).Error; err != nil {
		return warehouseKTPScanInMaterialRow{}, err
	}
	if capability.NeedsRubber == 1 && capability.HasLIYDD == 0 {
		return warehouseKTPScanInMaterialRow{}, errWarehouseKTPRubberMaterialUnavailable
	}

	requiredRubberCTE := `required_rubber AS (
	SELECT
		CAST('' AS varchar(20)) AS factory,
		CAST('' AS varchar(20)) AS CLBH,
		CONVERT(float, 0) AS required_qty
	WHERE 1 = 0
)`
	if capability.HasLIYDD == 1 {
		requiredRubberCTE = `required_rubber AS (
	SELECT
		factory_scope.factory,
		RIGHT(ZLZLS3.cldhz, 10) AS CLBH,
		SUM(CONVERT(float, ZLZLS3.TCLYL)) AS required_qty
	FROM order_info
	INNER JOIN factory_scope ON factory_scope.SFL = 'RB'
	INNER JOIN LIY_DD.dbo.ZLZLS3 ZLZLS3
		ON ZLZLS3.ZLBH1 = order_info.material_order
		AND ZLZLS3.CQDH = factory_scope.factory
	WHERE ZLZLS3.TCLYL > 0
	GROUP BY factory_scope.factory, RIGHT(ZLZLS3.cldhz, 10)
)`
	}

	excludedCartons := compactScanInCartons(append(append([]string{}, sessionCartons...), carton.CartonBar))
	exclusionClause, exclusionArgs := scanInCartonExclusionClause(excludedCartons)
	sessionClause, sessionArgs := scanInMaterialSessionClause(sessionCartons, orderNo)
	query := `
WITH target_order AS (
	SELECT COALESCE(
		(SELECT TOP 1 LTRIM(RTRIM(DDZLTR.DDBH)) FROM DDZLTR WHERE DDZLTR.DDBH1 = ?),
		?
	) AS material_order
),
order_info AS (
	SELECT TOP 1
		target_order.material_order,
		CONVERT(float, ISNULL(DDZL.Pairs, 0)) AS order_pairs,
		LTRIM(RTRIM(ISNULL(DDZL.GSBH, ''))) AS factory
	FROM target_order
	LEFT JOIN DDZL ON DDZL.DDBH = target_order.material_order
),
factory_scope AS (
	SELECT DISTINCT Bgszl.GSDH AS factory, Bgszl.SFL
	FROM order_info
	INNER JOIN Bgszl ON Bgszl.DFL IN (
		SELECT source_factory.DFL
		FROM Bgszl source_factory
		WHERE source_factory.GSDH = order_info.factory
	)
	WHERE ISNULL(Bgszl.tybh, '') <> 'Y'
),
required_main AS (
	SELECT
		order_info.factory,
		ZLZLS2.CLBH,
		SUM(CONVERT(float, ZLZLS2.CLSL)) AS required_qty
	FROM order_info
	INNER JOIN DDZL_ZLZLS2 ZLZLS2 ON ZLZLS2.ZLBH = order_info.material_order
	WHERE ZLZLS2.CLSL > 0
	GROUP BY order_info.factory, ZLZLS2.CLBH
),
required_gc AS (
	SELECT
		factory_scope.factory,
		ZLZLS2.CLBH,
		SUM(CONVERT(float, ZLZLS2.CLSL)) AS required_qty
	FROM order_info
	INNER JOIN factory_scope ON factory_scope.SFL = 'GC'
	INNER JOIN ZLZLS2_YM ZLZLS2 ON ZLZLS2.ZLBH = order_info.material_order
	WHERE ZLZLS2.CLSL > 0
	GROUP BY factory_scope.factory, ZLZLS2.CLBH
),
` + requiredRubberCTE + `,
requirements AS (
	SELECT factory, CLBH, SUM(required_qty) AS required_qty
	FROM (
		SELECT * FROM required_main
		UNION ALL
		SELECT * FROM required_gc
		UNION ALL
		SELECT * FROM required_rubber
	) all_requirements
	GROUP BY factory, CLBH
),
issued AS (
	SELECT
		KCLL.GSBH AS factory,
		KCLLS.CLBH,
		SUM(CONVERT(float, ISNULL(KCLLS.Qty, 0))) AS issued_qty
	FROM order_info
	INNER JOIN KCLLS ON KCLLS.SCBH = order_info.material_order
	INNER JOIN KCLL ON KCLL.LLNO = KCLLS.LLNO AND KCLL.CFMID <> 'NO'
	GROUP BY KCLL.GSBH, KCLLS.CLBH
),
coverage AS (
	SELECT
		requirements.factory,
		requirements.CLBH,
		requirements.required_qty,
		ISNULL(issued.issued_qty, 0) AS issued_qty,
		CASE
			WHEN requirements.required_qty > 0
				THEN ISNULL(issued.issued_qty, 0) * order_info.order_pairs / requirements.required_qty
			ELSE 0
		END AS covered_pairs
	FROM requirements
	CROSS JOIN order_info
	LEFT JOIN issued ON issued.factory = requirements.factory AND issued.CLBH = requirements.CLBH
),
limiting AS (
	SELECT TOP 1 factory, covered_pairs
	FROM coverage
	ORDER BY covered_pairs, factory, CLBH
)
SELECT
	order_info.material_order,
	order_info.order_pairs,
	CONVERT(float, ISNULL(limiting.covered_pairs, 0)) AS covered_pairs,
	ISNULL((
		SELECT SUM(CONVERT(float, ISNULL(cp.Qty, 0)))
		FROM YWCP cp
		INNER JOIN YWDD ywdd ON ywdd.DDBH = cp.DDBH
		WHERE ywdd.YSBH = ?
		  AND ISNULL(cp.INCS, 0) >= 1
		  ` + exclusionClause + `
	), 0) AS already_in_qty,
	` + sessionClause + ` AS session_qty,
	(SELECT COUNT(*) FROM coverage) AS required_material_rows,
	(SELECT COUNT(*) FROM coverage WHERE issued_qty > 0) AS issued_material_rows,
	(SELECT COUNT(*) FROM coverage WHERE issued_qty <= 0) AS missing_material_rows,
	LTRIM(RTRIM(ISNULL(limiting.factory, ''))) AS limiting_factory
FROM order_info
LEFT JOIN limiting ON 1 = 1`

	args := []interface{}{orderNo, orderNo, orderNo}
	args = append(args, exclusionArgs...)
	args = append(args, sessionArgs...)
	var row warehouseKTPScanInMaterialRow
	err := db.Raw(query, args...).Scan(&row).Error
	return row, err
}

func scanInMaterialSessionClause(sessionCartons []string, orderNo string) (string, []interface{}) {
	if len(sessionCartons) == 0 {
		return "CONVERT(float, 0)", nil
	}

	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(sessionCartons)), ",")
	clause := `ISNULL((
		SELECT SUM(CONVERT(float, ISNULL(cp.Qty, 0)))
		FROM YWCP cp
		INNER JOIN YWDD ywdd ON ywdd.DDBH = cp.DDBH
		WHERE cp.CARTONBAR IN (` + placeholders + `)
		  AND ywdd.YSBH = ?
	), 0)`
	args := make([]interface{}, 0, len(sessionCartons)+1)
	for _, carton := range sessionCartons {
		args = append(args, carton)
	}
	args = append(args, orderNo)
	return clause, args
}

func queryWarehouseKTPScanInInboxes(db *gorm.DB, cartonBar string) ([]warehouseKTPScanInInboxRow, error) {
	const query = `
SELECT
	LTRIM(RTRIM(COALESCE(
		NULLIF(LTRIM(RTRIM(pos.InBoxCode)), ''),
		mapping.InBoxBar,
		''
	))) AS inbox_bar,
	LTRIM(RTRIM(ISNULL(pos.DDCC, ''))) AS size,
	ISNULL(pos.Qty, 0) AS expected_qty
FROM YWCP cp
INNER JOIN YWBZPOS pos ON pos.DDBH = cp.DDBH AND pos.XH = cp.XH
LEFT JOIN YWDD ywdd ON ywdd.DDBH = cp.DDBH
LEFT JOIN DDZL ddzl ON ddzl.DDBH = ywdd.YSBH
OUTER APPLY (
	SELECT TOP 1 LTRIM(RTRIM(history.InBoxBar)) AS InBoxBar
	FROM YWCPLH history
	WHERE history.XieXing = ddzl.XieXing
		AND history.SheHao = ddzl.SheHao
		AND history.XXCC = pos.DDCC
		AND history.GBBH IN (NULLIF(LTRIM(RTRIM(ddzl.DDGB)), ''), 'NORMAL')
		AND ISNULL(history.InBoxBar, '') <> ''
	ORDER BY
		CASE WHEN history.GBBH = NULLIF(LTRIM(RTRIM(ddzl.DDGB)), '') THEN 0 ELSE 1 END,
		history.USERDATE DESC
) mapping
WHERE cp.CARTONBAR = ?
ORDER BY pos.DDCC`

	var rows []warehouseKTPScanInInboxRow
	err := db.Raw(query, cartonBar).Scan(&rows).Error
	return rows, err
}

func normalizeWarehouseKTPScanInMode(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), scanInModeCartonInbox) {
		return scanInModeCartonInbox
	}
	return scanInModeCarton
}

func evaluateWarehouseKTPScanInInboxes(
	scanMode string,
	rows []warehouseKTPScanInInboxRow,
	scans []string,
	queryErr error,
) types.WarehouseKTPScanInInboxCheck {
	result := types.WarehouseKTPScanInInboxCheck{
		WarehouseKTPScanInCheck: types.WarehouseKTPScanInCheck{Status: scanInStatusPassed, Passed: true},
		Enforced:                normalizeWarehouseKTPScanInMode(scanMode) == scanInModeCartonInbox,
		Items:                   make([]types.WarehouseKTPScanInInboxItem, 0),
	}
	if !result.Enforced {
		return result
	}
	result.Status = scanInStatusNoData
	result.Passed = false
	if queryErr != nil {
		result.Message = "Không đọc được dữ liệu hộp con của carton"
		return result
	}

	itemIndex := make(map[string]int)
	for _, row := range rows {
		if row.ExpectedQty <= 0 {
			continue
		}
		inboxBar := trimScanInString(row.InboxBar)
		if inboxBar == "" {
			result.Message = "Carton có size chưa được khai báo mã hộp con trong ERP"
			return result
		}
		key := strings.ToUpper(inboxBar)
		if index, exists := itemIndex[key]; exists {
			result.Items[index].ExpectedQty += row.ExpectedQty
			result.ExpectedQty += row.ExpectedQty
			continue
		}
		itemIndex[key] = len(result.Items)
		result.Items = append(result.Items, types.WarehouseKTPScanInInboxItem{
			InboxBar:    inboxBar,
			Size:        trimScanInString(row.Size),
			ExpectedQty: row.ExpectedQty,
		})
		result.ExpectedQty += row.ExpectedQty
	}
	if result.ExpectedQty == 0 {
		result.Message = "Carton chưa có dữ liệu mã hộp con trong ERP"
		return result
	}

	for _, scannedValue := range scans {
		inboxBar := strings.TrimSpace(scannedValue)
		if inboxBar == "" {
			continue
		}
		index, exists := itemIndex[strings.ToUpper(inboxBar)]
		if !exists {
			result.Status = scanInStatusFailed
			result.Message = fmt.Sprintf("Mã hộp con %s không thuộc carton này", inboxBar)
			result.RemainingQty = result.ExpectedQty - result.ScannedQty
			return result
		}
		item := &result.Items[index]
		if item.ScannedQty >= item.ExpectedQty {
			result.Status = scanInStatusFailed
			result.Message = fmt.Sprintf("Mã hộp con %s đã quét đủ %d hộp", item.InboxBar, item.ExpectedQty)
			result.RemainingQty = result.ExpectedQty - result.ScannedQty
			return result
		}
		item.ScannedQty++
		result.ScannedQty++
	}

	result.RemainingQty = result.ExpectedQty - result.ScannedQty
	if result.RemainingQty > 0 {
		result.Message = fmt.Sprintf(
			"Đã quét %d/%d hộp con, còn %d hộp",
			result.ScannedQty,
			result.ExpectedQty,
			result.RemainingQty,
		)
		return result
	}
	result.Status = scanInStatusPassed
	result.Passed = true
	result.Message = "Đã quét đủ hộp con trong carton"
	return result
}

func mapWarehouseKTPScanInCarton(row warehouseKTPScanInCartonRow, scannedCode string) types.WarehouseKTPScanInCarton {
	return types.WarehouseKTPScanInCarton{
		ScannedCode: scannedCode,
		ResolvedBy:  strings.TrimSpace(row.ResolvedBy),
		CartonBar:   strings.TrimSpace(row.CartonBar),
		LPNBarCode:  trimScanInString(row.LPNBarCode),
		DDBH:        strings.TrimSpace(row.DDBH),
		OrderNo:     trimScanInString(row.OrderNo),
		CartonNo:    row.CartonNo,
		XH:          trimScanInString(row.XH),
		Size:        trimScanInString(row.Size),
		Qty:         scanInFloat(row.Qty),
		PackageQty:  scanInFloat(row.PackageQty),
		SB:          trimScanInString(row.SB),
		CKBH:        trimScanInString(row.CKBH),
		KCBH:        trimScanInString(row.KCBH),
		KVBH:        trimScanInString(row.KVBH),
		INCS:        scanInInt(row.INCS),
		INDate:      row.INDate,
		LastInDate:  row.LastInDate,
		SGW:         scanInFloat(row.SGW),
		RGW:         row.RGW,
		Article:     trimScanInString(row.Article),
		XieMing:     trimScanInString(row.XieMing),
		Country:     trimScanInString(row.Country),
		KHPO:        trimScanInString(row.KHPO),
		DepNo:       trimScanInString(row.DepNo),
		DepName:     trimScanInString(row.DepName),
		INUser:      trimScanInString(row.INUser),
		Factory:     trimScanInString(row.Factory),
		FactoryType: trimScanInString(row.FactoryType),
		OrderFlag:   trimScanInString(row.OrderFlag),
	}
}

func evaluateWarehouseKTPScanInEligibility(
	carton types.WarehouseKTPScanInCarton,
	sessionCartons []string,
	production warehouseKTPScanInProductionRow,
) types.WarehouseKTPScanInEligibility {
	duplicate := false
	for _, sessionCarton := range sessionCartons {
		if strings.EqualFold(sessionCarton, carton.CartonBar) {
			duplicate = true
			break
		}
	}
	mixedOrder := production.MixedOrderRows > 0
	check := types.WarehouseKTPScanInEligibility{
		WarehouseKTPScanInCheck: types.WarehouseKTPScanInCheck{Status: scanInStatusPassed, Passed: true},
		SessionDuplicate:        duplicate,
		MixedOrder:              mixedOrder,
	}
	if duplicate {
		check.Status = scanInStatusFailed
		check.Passed = false
		check.Message = "Carton đã có trong phiên scan hiện tại"
		return check
	}
	if carton.OrderNo == "" || carton.DDBH == "" {
		check.Status = scanInStatusNoData
		check.Passed = false
		check.Message = "Carton chưa có đủ thông tin đơn hàng"
		return check
	}
	switch carton.SB {
	case "":
		check.Message = "Carton hợp lệ cho lần nhập KTP đầu tiên"
	case "2", "4":
		check.Message = "Carton thuộc luồng tái nhập KTP"
	case "1":
		check.Status = scanInStatusFailed
		check.Passed = false
		check.Message = "Carton đang ở trong kho thành phẩm (SB = 1)"
	case "3":
		check.Status = scanInStatusFailed
		check.Passed = false
		check.Message = "Carton đã xuất khỏi kho thành phẩm (SB = 3)"
	default:
		check.Status = scanInStatusFailed
		check.Passed = false
		check.Message = fmt.Sprintf("Trạng thái carton SB = %s không được phép Scan In", carton.SB)
	}
	return check
}

func evaluateWarehouseKTPScanInWeight(
	carton types.WarehouseKTPScanInCarton,
	gram *float64,
	actualWeight *float64,
) types.WarehouseKTPScanInWeightRule {
	standardWeight := carton.SGW
	rule := types.WarehouseKTPScanInWeightRule{
		WarehouseKTPScanInCheck: types.WarehouseKTPScanInCheck{Status: scanInStatusNoData},
		StandardWeight:          &standardWeight,
		ActualWeight:            actualWeight,
		Gram:                    gram,
	}
	if standardWeight <= 0 {
		rule.Status = scanInStatusFailed
		rule.Message = "Không có trọng lượng tiêu chuẩn thùng"
		return rule
	}
	if gram == nil || *gram <= 0 {
		rule.Status = scanInStatusFailed
		rule.Message = "Không có trọng lượng tiêu chuẩn giày"
		return rule
	}
	tolerance := *gram / 3
	rule.Tolerance = &tolerance
	rule.Comparison = "absolute"
	if actualWeight == nil || *actualWeight <= 0 {
		rule.Message = "Chưa đọc được G.W. từ cân"
		return rule
	}

	difference := standardWeight - *actualWeight
	rule.Difference = &difference
	failed := math.Abs(difference) > tolerance
	if failed {
		rule.Status = scanInStatusFailed
		rule.Message = "Độ lệch trọng lượng vượt giới hạn cho phép"
		return rule
	}
	rule.Status = scanInStatusPassed
	rule.Passed = true
	// rule.Message = "Trọng lượng nằm trong giới hạn ERP"
	return rule
}

func evaluateWarehouseKTPScanInProduction(
	carton types.WarehouseKTPScanInCarton,
	row warehouseKTPScanInProductionRow,
	queryErr error,
) types.WarehouseKTPScanInProductionCheck {
	currentQty := carton.PackageQty
	result := types.WarehouseKTPScanInProductionCheck{
		WarehouseKTPScanInCheck: types.WarehouseKTPScanInCheck{Status: scanInStatusNoData},
		Enforced:                carton.FactoryType == "PD",
		Size:                    carton.Size,
		CompletedQty:            row.CompletedQty,
		AlreadyInQty:            row.AlreadyInQty,
		SessionQty:              row.SessionQty,
		CurrentCartonQty:        currentQty,
		AfterScanQty:            row.AlreadyInQty + row.SessionQty + currentQty,
	}
	if queryErr != nil {
		result.Message = "Không đọc được dữ liệu sản lượng Gò"
		return result
	}
	if !result.Enforced {
		result.Status = scanInStatusPassed
		result.Passed = true
		result.Message = "ERP chỉ khóa sản lượng Gò tại nhà máy được phân loại PD"
		return result
	}
	if carton.OrderFlag == "1" {
		result.Status = scanInStatusPassed
		result.Passed = true
		result.Message = "Đơn hàng có cờ cho phép của ERP"
		return result
	}
	if row.CompletedQty <= 0 {
		result.Status = scanInStatusFailed
		result.Message = "Công đoạn Gò chưa có sản lượng hoàn thành cho size này"
		return result
	}
	if result.AfterScanQty > row.CompletedQty {
		result.Status = scanInStatusFailed
		result.Message = "Số lượng KTP sau carton này vượt sản lượng Gò hoàn thành"
		return result
	}
	result.Status = scanInStatusPassed
	result.Passed = true
	result.Message = "Sản lượng KTP không vượt sản lượng Gò hoàn thành"
	return result
}

func evaluateWarehouseKTPScanInMaterial(
	carton types.WarehouseKTPScanInCarton,
	_ warehouseKTPScanInProductionRow,
	row warehouseKTPScanInMaterialRow,
	queryErr error,
) types.WarehouseKTPScanInMaterialCheck {
	currentQty := carton.Qty
	result := types.WarehouseKTPScanInMaterialCheck{
		WarehouseKTPScanInCheck: types.WarehouseKTPScanInCheck{Status: scanInStatusNoData},
		MaterialOrder:           trimScanInString(row.MaterialOrder),
		OrderPairs:              scanInFloat(row.OrderPairs),
		CoveredPairs:            scanInFloat(row.CoveredPairs),
		AlreadyInQty:            row.AlreadyInQty,
		SessionQty:              row.SessionQty,
		CurrentCartonQty:        currentQty,
		AfterScanQty:            row.AlreadyInQty + row.SessionQty + currentQty,
		RequiredMaterialRows:    row.RequiredMaterialRows,
		IssuedMaterialRows:      row.IssuedMaterialRows,
		MissingMaterialRows:     row.MissingMaterialRows,
		LimitingFactory:         trimScanInString(row.LimitingFactory),
	}
	if queryErr != nil {
		if errors.Is(queryErr, errWarehouseKTPRubberMaterialUnavailable) {
			result.Message = "Máy thử nghiệm thiếu dữ liệu định mức vật tư cao su (LIY_DD)"
		} else {
			result.Message = "Không đọc được dữ liệu cấp phát vật tư"
		}
		return result
	}
	if result.MaterialOrder == "" || result.OrderPairs <= 0 || result.RequiredMaterialRows == 0 {
		result.Message = "Đơn hàng chưa có đủ định mức vật tư để đối chiếu"
		return result
	}
	if result.MissingMaterialRows > 0 {
		result.Status = scanInStatusFailed
		result.Message = "Còn vật tư chưa mở phiếu cấp phát"
		return result
	}
	if carton.SB != "" {
		result.Status = scanInStatusPassed
		result.Passed = true
		result.Message = "Carton tái nhập: ERP không khóa lại điều kiện vật tư"
		return result
	}
	if result.CoveredPairs >= result.OrderPairs {
		result.Status = scanInStatusPassed
		result.Passed = true
		result.Message = "Vật tư đã cấp đủ cho toàn bộ đơn hàng"
		return result
	}
	if result.AfterScanQty < result.CoveredPairs {
		result.Status = scanInStatusPassed
		result.Passed = true
		result.Message = "Số lượng sau carton vẫn nhỏ hơn mức vật tư đã cấp"
		return result
	}
	if carton.OrderFlag == "1" {
		result.Status = scanInStatusPassed
		result.Passed = true
		result.Message = "Đơn hàng có cờ cho phép của ERP"
		return result
	}
	result.Status = scanInStatusFailed
	result.Message = "Vật tư cấp phát chưa đủ cho carton hiện tại"
	return result
}

func compactScanInCartons(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{})
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		key := strings.ToUpper(trimmed)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, trimmed)
		if len(result) == 100 {
			break
		}
	}
	return result
}

func trimScanInString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func scanInFloat(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func scanInInt(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}
