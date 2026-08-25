package services

import (
	"fmt"
	"strings"
	"time"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"gorm.io/gorm"
)

const (
	scanOutModeRecycle    = "recycle"
	scanOutModeInspection = "inspection"
)

type warehouseKTPScanOutCartonRow struct {
	CartonBar   string     `gorm:"column:carton_bar"`
	LPNBarCode  string     `gorm:"column:lpn_bar_code"`
	DDBH        string     `gorm:"column:ddbh"`
	OrderNo     string     `gorm:"column:order_no"`
	XH          string     `gorm:"column:xh"`
	Qty         float64    `gorm:"column:qty"`
	SB          string     `gorm:"column:sb"`
	OutCount    int        `gorm:"column:out_count"`
	OutDate     *time.Time `gorm:"column:out_date"`
	ExecuteDate *time.Time `gorm:"column:execute_date"`
	Article     string     `gorm:"column:article"`
	XieMing     string     `gorm:"column:xie_ming"`
	Country     string     `gorm:"column:country"`
	KCBH        string     `gorm:"column:kcbh"`
	KVBH        string     `gorm:"column:kvbh"`
	ResolvedBy  string     `gorm:"column:resolved_by"`
}

type warehouseKTPScanOutDetailRow struct {
	OrderNo   string  `gorm:"column:order_no"`
	SKU       string  `gorm:"column:sku"`
	ModelName string  `gorm:"column:model_name"`
	XH        string  `gorm:"column:xh"`
	OrderSize string  `gorm:"column:order_size"`
	ModelSize string  `gorm:"column:model_size"`
	Qty       float64 `gorm:"column:qty"`
	From      int     `gorm:"column:from_no"`
	To        int     `gorm:"column:to_no"`
}

func (s *WarehouseKTPService) GetScanOutCarton(barcode string) (types.WarehouseKTPScanOutCarton, error) {
	barcode = strings.TrimSpace(barcode)
	if barcode == "" {
		return types.WarehouseKTPScanOutCarton{}, ErrWarehouseKTPCartonNotFound
	}

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return types.WarehouseKTPScanOutCarton{}, err
	}
	dbInstance, err := db.DB()
	if err != nil {
		return types.WarehouseKTPScanOutCarton{}, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	row, err := queryWarehouseKTPScanOutCarton(db, barcode)
	if err != nil {
		return types.WarehouseKTPScanOutCarton{}, err
	}
	result := mapWarehouseKTPScanOutCarton(row, barcode)
	result.Details, err = queryWarehouseKTPScanOutDetails(db, result.DDBH)
	return result, err
}

func (s *WarehouseKTPService) SubmitScanOutCarton(
	params request.WarehouseKTPScanOutSubmitRequest,
) (types.WarehouseKTPScanOutCarton, error) {
	params.CartonBar = strings.TrimSpace(params.CartonBar)
	targetStatus, err := warehouseKTPScanOutTargetStatus(params.Mode)
	if err != nil {
		return types.WarehouseKTPScanOutCarton{}, err
	}

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return types.WarehouseKTPScanOutCarton{}, err
	}
	dbInstance, err := db.DB()
	if err != nil {
		return types.WarehouseKTPScanOutCarton{}, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	var result types.WarehouseKTPScanOutCarton
	err = db.Transaction(func(tx *gorm.DB) error {
		row, queryErr := queryWarehouseKTPScanOutCarton(tx, params.CartonBar)
		if queryErr != nil {
			return queryErr
		}
		if lockErr := acquireWarehouseKTPScanOutLock(tx, row.CartonBar); lockErr != nil {
			return lockErr
		}
		row, queryErr = queryWarehouseKTPScanOutCarton(tx, row.CartonBar)
		if queryErr != nil {
			return queryErr
		}
		carton := mapWarehouseKTPScanOutCarton(row, params.CartonBar)
		if !carton.Eligible {
			return fmt.Errorf("%w: %s (%s)", ErrWarehouseKTPScanOutRejected, carton.Message, carton.CartonBar)
		}

		update := tx.Exec(`
UPDATE YWCP
SET OUTCS = ISNULL(OUTCS, 0) + 1,
	OUTDATE = GETDATE(),
	SB = ?
WHERE CARTONBAR = ?
	AND LTRIM(RTRIM(ISNULL(SB, ''))) = '1'
	AND EXEDATE IS NULL`, targetStatus, carton.CartonBar)
		if update.Error != nil {
			return update.Error
		}
		if update.RowsAffected != 1 {
			return ErrWarehouseKTPScanOutUpdateFailed
		}

		if updateOrder := tx.Exec(`UPDATE DDZL SET YN = '2' WHERE DDBH = ? AND YN = '5'`, carton.DDBH); updateOrder.Error != nil {
			return updateOrder.Error
		}
		if updatePacking := tx.Exec(`UPDATE YWDD SET SB = 'IN' WHERE YSBH = ? AND SB = 'OK'`, carton.DDBH); updatePacking.Error != nil {
			return updatePacking.Error
		}

		savedRow, queryErr := queryWarehouseKTPScanOutCarton(tx, carton.CartonBar)
		if queryErr != nil {
			return queryErr
		}
		if strings.TrimSpace(savedRow.SB) != targetStatus || savedRow.OutCount != row.OutCount+1 || savedRow.OutDate == nil {
			return ErrWarehouseKTPScanOutUpdateFailed
		}
		result = mapWarehouseKTPScanOutCarton(savedRow, params.CartonBar)
		result.Message = "Scan Out thành công"
		result.Details, queryErr = queryWarehouseKTPScanOutDetails(tx, result.DDBH)
		return queryErr
	})
	return result, err
}

func warehouseKTPScanOutTargetStatus(mode string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case scanOutModeRecycle:
		return "2", nil
	case scanOutModeInspection:
		return "4", nil
	default:
		return "", fmt.Errorf("%w: phải chọn Tái chế hoặc Kiểm hàng", ErrWarehouseKTPScanOutRejected)
	}
}

func queryWarehouseKTPScanOutCarton(db *gorm.DB, barcode string) (warehouseKTPScanOutCartonRow, error) {
	row, err := queryWarehouseKTPScanOutCartonByCode(db, "cp.CARTONBAR = ?", barcode, "cartonBar")
	if err != nil || strings.TrimSpace(row.CartonBar) != "" {
		return row, err
	}

	lpnCode := barcode
	if len(lpnCode) > 15 {
		lpnCode = lpnCode[len(lpnCode)-15:]
	}
	row, err = queryWarehouseKTPScanOutCartonByCode(db, "cp.LPNBarCode = ?", lpnCode, "lpnBarcode")
	if err != nil {
		return row, err
	}
	if strings.TrimSpace(row.CartonBar) == "" {
		return row, ErrWarehouseKTPCartonNotFound
	}
	return row, nil
}

func queryWarehouseKTPScanOutCartonByCode(
	db *gorm.DB,
	lookupCondition string,
	lookupValue string,
	resolvedBy string,
) (warehouseKTPScanOutCartonRow, error) {
	query := `
SELECT TOP 1
	LTRIM(RTRIM(cp.CARTONBAR)) AS carton_bar,
	LTRIM(RTRIM(ISNULL(cp.LPNBarCode, ''))) AS lpn_bar_code,
	LTRIM(RTRIM(ISNULL(cp.DDBH, ''))) AS ddbh,
	LTRIM(RTRIM(ISNULL(ywdd.YSBH, ''))) AS order_no,
	LTRIM(RTRIM(ISNULL(cp.XH, ''))) AS xh,
	CONVERT(float, ISNULL(cp.Qty, 0)) AS qty,
	LTRIM(RTRIM(ISNULL(cp.SB, ''))) AS sb,
	ISNULL(cp.OUTCS, 0) AS out_count,
	cp.OUTDATE AS out_date,
	cp.EXEDATE AS execute_date,
	LTRIM(RTRIM(ISNULL(ddzl.Article, ''))) AS article,
	LTRIM(RTRIM(ISNULL(shoe.XieMing, ''))) AS xie_ming,
	LTRIM(RTRIM(COALESCE(NULLIF(country.ywsm, ''), ddzl.DDGB, ''))) AS country,
	LTRIM(RTRIM(ISNULL(cp.KCBH, ''))) AS kcbh,
	LTRIM(RTRIM(ISNULL(cp.KVBH, ''))) AS kvbh,
	? AS resolved_by
FROM YWCP cp
LEFT JOIN YWDD ywdd ON ywdd.DDBH = cp.DDBH
LEFT JOIN DDZL ddzl ON ddzl.DDBH = ywdd.YSBH
LEFT JOIN XXZL shoe ON shoe.XieXing = ddzl.XieXing AND shoe.SheHao = ddzl.SheHao
LEFT JOIN LBZLS country ON country.lb = '06' AND country.lbdh = ddzl.DDGB
WHERE ` + lookupCondition

	var row warehouseKTPScanOutCartonRow
	if err := db.Raw(query, resolvedBy, lookupValue).Scan(&row).Error; err != nil {
		return row, err
	}
	return row, nil
}

func queryWarehouseKTPScanOutDetails(db *gorm.DB, ddbh string) ([]types.WarehouseKTPScanOutDetail, error) {
	return queryWarehouseKTPDetails(db, ddbh, "", false)
}

func queryWarehouseKTPScanShipDetails(db *gorm.DB, ddbh string, xh string) ([]types.WarehouseKTPScanOutDetail, error) {
	return queryWarehouseKTPDetails(db, ddbh, strings.TrimSpace(xh), true)
}

func queryWarehouseKTPDetails(
	db *gorm.DB,
	ddbh string,
	xh string,
	exactCarton bool,
) ([]types.WarehouseKTPScanOutDetail, error) {
	const query = `
SELECT
	LTRIM(RTRIM(ISNULL(pos.DDBH, ''))) AS order_no,
	LTRIM(RTRIM(ISNULL(ddzl.Article, ''))) AS sku,
	LTRIM(RTRIM(ISNULL(shoe.XieMing, ''))) AS model_name,
	LTRIM(RTRIM(ISNULL(pos.XH, ''))) AS xh,
	LTRIM(RTRIM(ISNULL(pos.DDCC, ''))) AS order_size,
	LTRIM(RTRIM(ISNULL(size_map.XXCC, ''))) AS model_size,
	CONVERT(float, ISNULL(pos.Qty, 0)) AS qty,
	ISNULL(pos.CTQ, 0) AS from_no,
	ISNULL(pos.CTZ, 0) AS to_no
FROM YWBZPOS pos
LEFT JOIN YWDD ywdd ON ywdd.DDBH = pos.DDBH
LEFT JOIN DDZL ddzl ON ddzl.DDBH = ywdd.YSBH
LEFT JOIN XXZL shoe ON shoe.XieXing = ddzl.XieXing AND shoe.SheHao = ddzl.SheHao
OUTER APPLY (
	SELECT TOP 1 sizes.XXCC
	FROM YWDDS sizes
	WHERE sizes.DDBH = pos.DDBH
		AND sizes.DDCC = pos.DDCC
	ORDER BY sizes.XXCC
) size_map
WHERE pos.DDBH = ?
	AND (? = 0 OR LTRIM(RTRIM(ISNULL(pos.XH, ''))) = ?)
ORDER BY pos.XH, pos.DDCC`

	var rows []warehouseKTPScanOutDetailRow
	exactCartonFlag := 0
	if exactCarton {
		exactCartonFlag = 1
	}
	if err := db.Raw(query, ddbh, exactCartonFlag, xh).Scan(&rows).Error; err != nil {
		return nil, err
	}
	if exactCarton && len(rows) > 1 {
		rows = rows[:1]
	}
	details := make([]types.WarehouseKTPScanOutDetail, 0, len(rows))
	for _, row := range rows {
		details = append(details, types.WarehouseKTPScanOutDetail{
			OrderNo:   strings.TrimSpace(row.OrderNo),
			SKU:       strings.TrimSpace(row.SKU),
			ModelName: strings.TrimSpace(row.ModelName),
			XH:        strings.TrimSpace(row.XH),
			OrderSize: strings.TrimSpace(row.OrderSize),
			ModelSize: strings.TrimSpace(row.ModelSize),
			Qty:       row.Qty,
			From:      row.From,
			To:        row.To,
		})
	}
	return details, nil
}

func mapWarehouseKTPScanOutCarton(row warehouseKTPScanOutCartonRow, scannedCode string) types.WarehouseKTPScanOutCarton {
	result := types.WarehouseKTPScanOutCarton{
		ScannedCode: strings.TrimSpace(scannedCode),
		ResolvedBy:  strings.TrimSpace(row.ResolvedBy),
		CartonBar:   strings.TrimSpace(row.CartonBar),
		LPNBarCode:  strings.TrimSpace(row.LPNBarCode),
		DDBH:        strings.TrimSpace(row.DDBH),
		OrderNo:     strings.TrimSpace(row.OrderNo),
		XH:          strings.TrimSpace(row.XH),
		Qty:         row.Qty,
		SB:          strings.TrimSpace(row.SB),
		OutCount:    row.OutCount,
		OutDate:     row.OutDate,
		ExecuteDate: row.ExecuteDate,
		Article:     strings.TrimSpace(row.Article),
		XieMing:     strings.TrimSpace(row.XieMing),
		Country:     strings.TrimSpace(row.Country),
		KCBH:        strings.TrimSpace(row.KCBH),
		KVBH:        strings.TrimSpace(row.KVBH),
		Details:     make([]types.WarehouseKTPScanOutDetail, 0),
	}
	result.Eligible, result.Message = evaluateWarehouseKTPScanOutEligibility(result)
	return result
}

func evaluateWarehouseKTPScanOutEligibility(carton types.WarehouseKTPScanOutCarton) (bool, string) {
	if carton.CartonBar == "" {
		return false, "Không tìm thấy carton"
	}
	if carton.DDBH == "" {
		return false, "Carton chưa có thông tin đơn hàng"
	}
	if carton.ExecuteDate != nil {
		return false, "Carton đã có ngày xuất hàng EXEDATE"
	}
	if carton.SB != "1" {
		switch carton.SB {
		case "2":
			return false, "Carton đã đưa đi tái chế (SB = 2)"
		case "4":
			return false, "Carton đã đưa đi kiểm hàng (SB = 4)"
		case "3":
			return false, "Carton đã xuất khỏi kho thành phẩm (SB = 3)"
		case "":
			return false, "Carton chưa Scan In vào kho thành phẩm"
		default:
			return false, fmt.Sprintf("Trạng thái carton SB = %s không cho phép Scan Out", carton.SB)
		}
	}
	return true, "Carton hợp lệ để Scan Out"
}

func acquireWarehouseKTPScanOutLock(tx *gorm.DB, cartonBar string) error {
	resource := "warehouse-ktp-scan-out:" + strings.ToUpper(strings.TrimSpace(cartonBar))
	return tx.Exec(`
DECLARE @lock_result INT;
EXEC @lock_result = sys.sp_getapplock
	@Resource = ?,
	@LockMode = 'Exclusive',
	@LockOwner = 'Transaction',
	@LockTimeout = 15000;
IF @lock_result < 0
	THROW 50001, 'Khong the khoa carton de Scan Out.', 1;`, resource).Error
}
