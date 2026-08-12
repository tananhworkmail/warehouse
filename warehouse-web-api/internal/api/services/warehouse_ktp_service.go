package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"

	"gorm.io/gorm"
)

var (
	ErrWarehouseKTPRackNotFound        = errors.New("rack not found")
	ErrWarehouseKTPInvalidScanCode     = errors.New("scan code must be exactly 6 digits")
	ErrWarehouseKTPScanStorageNotFound = errors.New("missing existing storage for KTP rack scan codes")
	ErrWarehouseKTPCartonNotFound      = errors.New("carton not found")
	ErrWarehouseKTPOrderNotFound       = errors.New("order not found")
	ErrWarehouseKTPInvalidRackCode     = errors.New("invalid rack code")
	ErrWarehouseKTPSameRack            = errors.New("new rack is already used by this order")
)

type WarehouseKTPService struct {
	*BaseService
}

var WarehouseKTP = &WarehouseKTPService{}

var warehouseKTPScanPattern = regexp.MustCompile(`^\d{6}$`)
var warehouseKTPRackPattern = regexp.MustCompile(`^[A-Z]+\d+$`)

const warehouseKTPEligibleSBCondition = "SB IS NOT NULL AND LTRIM(RTRIM(SB)) <> '3'"

type warehouseKTPRackRow struct {
	RackCode              string  `gorm:"column:rack_code"`
	DDBH                  string  `gorm:"column:ddbh"`
	TotalQty              float64 `gorm:"column:total_qty"`
	InboundCartonCount    int     `gorm:"column:inbound_carton_count"`
	RecycleCartonCount    int     `gorm:"column:recycle_carton_count"`
	InspectionCartonCount int     `gorm:"column:inspection_carton_count"`
	CodebarCount          int     `gorm:"column:codebar_count"`
	SampleCode1           *string `gorm:"column:sample_code_1"`
	SampleCode2           *string `gorm:"column:sample_code_2"`
	SampleCode3           *string `gorm:"column:sample_code_3"`
}

type warehouseKTPRackSignatureRow struct {
	RackCode              string  `gorm:"column:rack_code"`
	DDBH                  string  `gorm:"column:ddbh"`
	TotalQty              float64 `gorm:"column:total_qty"`
	InboundCartonCount    int     `gorm:"column:inbound_carton_count"`
	RecycleCartonCount    int     `gorm:"column:recycle_carton_count"`
	InspectionCartonCount int     `gorm:"column:inspection_carton_count"`
	CodebarCount          int     `gorm:"column:codebar_count"`
}

type warehouseKTPCartonDetailRow struct {
	DDBH      string  `gorm:"column:ddbh"`
	SB        string  `gorm:"column:sb"`
	CartonBar string  `gorm:"column:carton_bar"`
	CartonNo  string  `gorm:"column:carton_no"`
	Qty       float64 `gorm:"column:qty"`
}

type warehouseKTPOrderRackRow struct {
	RackCode     string  `gorm:"column:rack_code"`
	TotalQty     float64 `gorm:"column:total_qty"`
	CodebarCount int     `gorm:"column:codebar_count"`
}

type warehouseKTPOrderTotalRow struct {
	TotalQty     float64 `gorm:"column:total_qty"`
	CodebarCount int     `gorm:"column:codebar_count"`
}

func (s *WarehouseKTPService) GetRacks() ([]types.WarehouseKTPRack, error) {
	return s.queryRackSummary("")
}

func (s *WarehouseKTPService) SearchRacks(keyword string) ([]types.WarehouseKTPRack, error) {
	return s.queryRackSummary(strings.TrimSpace(keyword))
}

func (s *WarehouseKTPService) GetRackOrderDetail(rackCode string) (types.WarehouseKTPRackOrderDetail, error) {
	normalizedRackCode := normalizeWarehouseKTPRackCode(rackCode)
	if !warehouseKTPRackPattern.MatchString(normalizedRackCode) {
		return types.WarehouseKTPRackOrderDetail{}, ErrWarehouseKTPInvalidRackCode
	}

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return types.WarehouseKTPRackOrderDetail{}, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return types.WarehouseKTPRackOrderDetail{}, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	query := fmt.Sprintf(`
WITH status_cartons AS (
	SELECT
		LTRIM(RTRIM(DDBH)) AS DDBH,
		LTRIM(RTRIM(SB)) AS SB,
		LTRIM(RTRIM(CARTONBAR)) AS CARTONBAR,
		MAX(LTRIM(RTRIM(CONVERT(varchar(50), CARTONNO)))) AS carton_no,
		MAX(ISNULL(Qty, 0)) AS qty
	FROM YWCP WITH (NOLOCK)
	WHERE %s
	  AND KCBH = 'A2'
	  AND LTRIM(RTRIM(KVBH)) = ?
	  AND DDBH IS NOT NULL
	  AND LTRIM(RTRIM(DDBH)) <> ''
	  AND CARTONBAR IS NOT NULL
	  AND LTRIM(RTRIM(CARTONBAR)) <> ''
	GROUP BY
		LTRIM(RTRIM(DDBH)),
		LTRIM(RTRIM(SB)),
		LTRIM(RTRIM(CARTONBAR))
)
SELECT
	DDBH,
	SB,
	CARTONBAR AS carton_bar,
	COALESCE(NULLIF(carton_no, ''), CARTONBAR) AS carton_no,
	qty
FROM status_cartons
ORDER BY DDBH, SB, carton_no`, warehouseKTPEligibleSBCondition)

	var rows []warehouseKTPCartonDetailRow
	if err := db.Raw(query, normalizedRackCode).Scan(&rows).Error; err != nil {
		return types.WarehouseKTPRackOrderDetail{}, err
	}

	orders := make([]types.WarehouseKTPRackOrder, 0)
	orderIndexes := make(map[string]int)
	orderCartons := make(map[string]map[string]float64)
	statusIndexes := map[string]int{"1": 0, "2": 1, "4": 2}

	for _, row := range rows {
		ddbh := strings.TrimSpace(row.DDBH)
		if ddbh == "" {
			continue
		}

		orderIndex, exists := orderIndexes[ddbh]
		if !exists {
			orderIndex = len(orders)
			orderIndexes[ddbh] = orderIndex
			orderCartons[ddbh] = make(map[string]float64)
			orders = append(orders, types.WarehouseKTPRackOrder{
				DDBH: ddbh,
				Statuses: []types.WarehouseKTPStatusDetail{
					{SB: "1", Cartons: make([]types.WarehouseKTPCartonDetail, 0)},
					{SB: "2", Cartons: make([]types.WarehouseKTPCartonDetail, 0)},
					{SB: "4", Cartons: make([]types.WarehouseKTPCartonDetail, 0)},
				},
			})
		}

		cartonBar := strings.TrimSpace(row.CartonBar)
		previousQty, counted := orderCartons[ddbh][cartonBar]
		if !counted {
			orderCartons[ddbh][cartonBar] = row.Qty
			orders[orderIndex].CodebarCount++
			orders[orderIndex].TotalQty += row.Qty
		} else if row.Qty > previousQty {
			orderCartons[ddbh][cartonBar] = row.Qty
			orders[orderIndex].TotalQty += row.Qty - previousQty
		}

		statusIndex, ok := statusIndexes[strings.TrimSpace(row.SB)]
		if !ok {
			continue
		}

		status := &orders[orderIndex].Statuses[statusIndex]
		status.Cartons = append(status.Cartons, types.WarehouseKTPCartonDetail{
			CartonNo: strings.TrimSpace(row.CartonNo),
			Qty:      row.Qty,
		})
		status.CartonCount++
		status.TotalQty += row.Qty
	}

	return types.WarehouseKTPRackOrderDetail{
		RackCode: normalizedRackCode,
		Orders:   orders,
	}, nil
}

func (s *WarehouseKTPService) GetTempHumidityByDevices() ([]types.GetTempHumidity, error) {
	db, err := database.TempHumidity_Connection()
	if err != nil {
		return nil, err
	}

	var result []types.GetTempHumidity
	slotClause := ""
	if shouldUseTempHumiditySlotRows("") {
		if err := ensureTempHumidityFakeSlotRows([]string{"KTP01", "KTP02"}); err != nil {
			return nil, err
		}
		slotClause = " AND DATEPART(minute, RecordTime) IN (0, 30) AND DATEPART(second, RecordTime) = 0 AND DATEPART(millisecond, RecordTime) = 0"
	}

	err = db.Raw(`
		WITH x AS (
			SELECT
				DeviceName,
				DeviceAddr,
				Hum,
				Tem,
				RecordTime,
				ROW_NUMBER() OVER (
					PARTITION BY DeviceName
					ORDER BY RecordTime DESC
				) AS rn
			FROM tbhistory
			WHERE DeviceName IN ('KTP01','KTP02')
			AND RecordTime >= CAST(GETDATE() AS DATE)
			AND RecordTime < DATEADD(DAY, 1, CAST(GETDATE() AS DATE))
			` + slotClause + `
		)
		SELECT
			DeviceName,
			DeviceAddr,
			Hum,
			Tem,
			RecordTime
		FROM x
		WHERE rn = 1
		ORDER BY DeviceName;
	`).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WarehouseKTPService) GetLatestTempHumidityLogByDevices(queryDate string, startTime string, endTime string) ([]types.GetHumidityAlert, error) {
	db, err := database.TempHumidity_Connection()
	if err != nil {
		return nil, err
	}

	if queryDate == "" {
		queryDate = time.Now().Format("2006-01-02")
	}

	var result []types.GetHumidityAlert
	args := []interface{}{queryDate, queryDate}
	timeClause := ""
	if startTime != "" && endTime != "" {
		timeClause = " AND CONVERT(time, RecordTime) >= CONVERT(time, ?) AND CONVERT(time, RecordTime) <= CONVERT(time, ?)"
		args = append(args, startTime, endTime)
	}
	slotClause := ""
	if shouldUseTempHumiditySlotRows(queryDate) {
		if err := ensureTempHumidityFakeSlotRows([]string{"KTP01", "KTP02"}); err != nil {
			return nil, err
		}
		slotClause = " AND DATEPART(minute, RecordTime) IN (0, 30) AND DATEPART(second, RecordTime) = 0 AND DATEPART(millisecond, RecordTime) = 0"
	}

	query := `
		WITH base AS (
			SELECT
				DeviceName,
				Hum,
				Tem AS Temp,
				CONVERT(varchar, CAST(RecordTime AS date), 23) AS Alert_Date,
				DATEADD(minute, (DATEDIFF(minute, 0, RecordTime) / 30) * 30, 0) AS HalfHourBucket,
				CONVERT(varchar, RecordTime, 120) AS RecordTime,
				RecordTime AS RawRecordTime
			FROM tbhistory
			WHERE DeviceName IN ('KTP01','KTP02')
			  AND RecordTime >= CAST(? AS DATE)
			  AND RecordTime < DATEADD(DAY, 1, CAST(? AS DATE))
	` + timeClause + slotClause + `
		),
		ranked AS (
			SELECT
				*,
				ROW_NUMBER() OVER (
					PARTITION BY DeviceName, HalfHourBucket
					ORDER BY
						ABS(DATEDIFF(second, RawRecordTime, HalfHourBucket)) ASC,
						RawRecordTime DESC
				) AS rn
			FROM base
		)
		SELECT
			DeviceName,
			Hum,
			Temp,
			Alert_Date,
			CONVERT(varchar, HalfHourBucket, 108) AS Alert_Time,
			RecordTime
		FROM ranked
		WHERE rn = 1
		ORDER BY DeviceName, HalfHourBucket DESC
	`

	err = db.Raw(query, args...).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WarehouseKTPService) GetRealtimeSignature() (string, error) {
	rows, err := s.queryRackSignature()
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for _, row := range rows {
		builder.WriteString(fmt.Sprintf(
			"%s:%s:%.3f:%d:%d:%d:%d|",
			strings.TrimSpace(row.RackCode),
			strings.TrimSpace(row.DDBH),
			row.TotalQty,
			row.InboundCartonCount,
			row.RecycleCartonCount,
			row.InspectionCartonCount,
			row.CodebarCount,
		))
	}

	return builder.String(), nil
}

func (s *WarehouseKTPService) GetMoveOrderByCarton(cartonBar string) (types.WarehouseKTPMoveOrderInfo, error) {
	normalizedCartonBar := strings.TrimSpace(cartonBar)
	if normalizedCartonBar == "" {
		return types.WarehouseKTPMoveOrderInfo{}, ErrWarehouseKTPCartonNotFound
	}

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return types.WarehouseKTPMoveOrderInfo{}, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return types.WarehouseKTPMoveOrderInfo{}, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	ddbh, err := s.findDDBHByCarton(db, normalizedCartonBar)
	if err != nil {
		return types.WarehouseKTPMoveOrderInfo{}, err
	}

	return s.getMoveOrderByDDBH(db, ddbh, normalizedCartonBar)
}

func (s *WarehouseKTPService) MoveOrderToRack(cartonBar string, newRackCode string, createdBy string) (types.WarehouseKTPMoveOrderInfo, error) {
	normalizedCartonBar := strings.TrimSpace(cartonBar)
	normalizedRackCode := normalizeWarehouseKTPRackCode(newRackCode)
	if normalizedCartonBar == "" {
		return types.WarehouseKTPMoveOrderInfo{}, ErrWarehouseKTPCartonNotFound
	}
	if !warehouseKTPRackPattern.MatchString(normalizedRackCode) {
		return types.WarehouseKTPMoveOrderInfo{}, ErrWarehouseKTPInvalidRackCode
	}

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return types.WarehouseKTPMoveOrderInfo{}, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return types.WarehouseKTPMoveOrderInfo{}, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	tx := db.Begin()
	if tx.Error != nil {
		return types.WarehouseKTPMoveOrderInfo{}, tx.Error
	}
	defer tx.Rollback()

	ddbh, err := s.findDDBHByCarton(tx, normalizedCartonBar)
	if err != nil {
		return types.WarehouseKTPMoveOrderInfo{}, err
	}

	orderInfo, err := s.getMoveOrderByDDBH(tx, ddbh, normalizedCartonBar)
	if err != nil {
		return types.WarehouseKTPMoveOrderInfo{}, err
	}
	for _, rack := range orderInfo.CurrentRacks {
		if normalizeWarehouseKTPRackCode(rack.RackCode) == normalizedRackCode {
			return types.WarehouseKTPMoveOrderInfo{}, ErrWarehouseKTPSameRack
		}
	}

	updateQuery := fmt.Sprintf(`
UPDATE YWCP
SET KVBH = ?
WHERE KCBH = 'A2'
  AND %s
  AND LTRIM(RTRIM(DDBH)) = ?
  AND CARTONBAR IS NOT NULL
  AND LTRIM(RTRIM(CARTONBAR)) <> ''`, warehouseKTPEligibleSBCondition)

	result := tx.Exec(updateQuery, normalizedRackCode, ddbh)
	if result.Error != nil {
		return types.WarehouseKTPMoveOrderInfo{}, result.Error
	}
	if result.RowsAffected == 0 {
		return types.WarehouseKTPMoveOrderInfo{}, ErrWarehouseKTPOrderNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return types.WarehouseKTPMoveOrderInfo{}, err
	}

	updatedInfo, err := s.getMoveOrderByDDBH(db, ddbh, normalizedCartonBar)
	if err != nil {
		return types.WarehouseKTPMoveOrderInfo{}, err
	}
	updatedInfo.NewRackCode = normalizedRackCode
	updatedInfo.UpdatedRows = result.RowsAffected

	return updatedInfo, nil
}

func (s *WarehouseKTPService) ScanRack(rackCode string, scanCode string, replace bool, createdBy string) (types.WarehouseKTPRack, error) {
	if !warehouseKTPScanPattern.MatchString(strings.TrimSpace(scanCode)) {
		return types.WarehouseKTPRack{}, ErrWarehouseKTPInvalidScanCode
	}

	racks, err := s.SearchRacks(normalizeWarehouseKTPRackCode(rackCode))
	if err != nil {
		return types.WarehouseKTPRack{}, err
	}
	if len(racks) == 0 {
		return types.WarehouseKTPRack{}, ErrWarehouseKTPRackNotFound
	}

	return types.WarehouseKTPRack{}, ErrWarehouseKTPScanStorageNotFound
}

func (s *WarehouseKTPService) ClearRack(rackCode string, createdBy string) (types.WarehouseKTPRack, error) {
	racks, err := s.SearchRacks(normalizeWarehouseKTPRackCode(rackCode))
	if err != nil {
		return types.WarehouseKTPRack{}, err
	}
	if len(racks) == 0 {
		return types.WarehouseKTPRack{}, ErrWarehouseKTPRackNotFound
	}

	return types.WarehouseKTPRack{}, ErrWarehouseKTPScanStorageNotFound
}

func (s *WarehouseKTPService) queryRackSummary(keyword string) ([]types.WarehouseKTPRack, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	query := fmt.Sprintf(`
WITH rack_rows AS (
	SELECT
		LTRIM(RTRIM(KVBH)) AS rack_code,
		LTRIM(RTRIM(CARTONBAR)) AS CARTONBAR,
		LTRIM(RTRIM(DDBH)) AS DDBH,
		LTRIM(RTRIM(SB)) AS SB,
		ISNULL(Qty, 0) AS qty
	FROM YWCP WITH (NOLOCK)
	WHERE %s
	  AND KCBH = 'A2'
	  AND KVBH IS NOT NULL
	  AND LTRIM(RTRIM(KVBH)) <> ''
	  AND LTRIM(RTRIM(KVBH)) <> '0'
	  AND CARTONBAR IS NOT NULL
	  AND LTRIM(RTRIM(CARTONBAR)) <> ''
),
carton_rows AS (
	SELECT
		rack_code,
		CARTONBAR,
		MAX(DDBH) AS DDBH,
		MAX(qty) AS qty
	FROM rack_rows
	GROUP BY rack_code, CARTONBAR
),
ranked AS (
	SELECT
		rack_code,
		CARTONBAR,
		DDBH,
		qty,
		ROW_NUMBER() OVER (PARTITION BY rack_code ORDER BY CARTONBAR) AS rn
	FROM carton_rows
),
rack_summary AS (
	SELECT
		rack_code,
		SUM(qty) AS total_qty,
		COUNT(DISTINCT CARTONBAR) AS codebar_count,
		MAX(CASE WHEN rn = 1 THEN CARTONBAR END) AS sample_code_1,
		MAX(CASE WHEN rn = 2 THEN CARTONBAR END) AS sample_code_2,
		MAX(CASE WHEN rn = 3 THEN CARTONBAR END) AS sample_code_3
	FROM ranked
	GROUP BY rack_code
),
order_cartons AS (
	SELECT
		rack_code,
		DDBH,
		CARTONBAR,
		MAX(qty) AS qty
	FROM rack_rows
	WHERE DDBH IS NOT NULL AND DDBH <> ''
	GROUP BY rack_code, DDBH, CARTONBAR
),
order_totals AS (
	SELECT
		rack_code,
		DDBH,
		SUM(qty) AS order_total_qty
	FROM order_cartons
	GROUP BY rack_code, DDBH
),
ranked_orders AS (
	SELECT
		rack_code,
		DDBH,
		ROW_NUMBER() OVER (
			PARTITION BY rack_code
			ORDER BY order_total_qty DESC, DDBH
		) AS rn
	FROM order_totals
),
selected_orders AS (
	SELECT rack_code, DDBH
	FROM ranked_orders
	WHERE rn = 1
),
selected_order_cartons AS (
	SELECT
		r.rack_code,
		r.CARTONBAR,
		s.DDBH,
		MAX(CASE WHEN r.SB = '1' THEN 1 ELSE 0 END) AS inbound_carton,
		MAX(CASE WHEN r.SB = '2' THEN 1 ELSE 0 END) AS recycle_carton,
		MAX(CASE WHEN r.SB = '4' THEN 1 ELSE 0 END) AS inspection_carton
	FROM rack_rows r
	INNER JOIN selected_orders s
		ON s.rack_code = r.rack_code
		AND s.DDBH = r.DDBH
	GROUP BY r.rack_code, r.CARTONBAR, s.DDBH
),
order_details AS (
	SELECT
		rack_code,
		DDBH,
		SUM(inbound_carton) AS inbound_carton_count,
		SUM(recycle_carton) AS recycle_carton_count,
		SUM(inspection_carton) AS inspection_carton_count
	FROM selected_order_cartons
	GROUP BY rack_code, DDBH
)
SELECT
	rs.rack_code,
	od.DDBH,
	rs.total_qty,
	ISNULL(od.inbound_carton_count, 0) AS inbound_carton_count,
	ISNULL(od.recycle_carton_count, 0) AS recycle_carton_count,
	ISNULL(od.inspection_carton_count, 0) AS inspection_carton_count,
	rs.codebar_count,
	rs.sample_code_1,
	rs.sample_code_2,
	rs.sample_code_3
FROM rack_summary rs
LEFT JOIN order_details od ON od.rack_code = rs.rack_code
`, warehouseKTPEligibleSBCondition)

	args := []interface{}{}
	if keyword != "" {
		query += `WHERE EXISTS (
	SELECT 1
	FROM rack_rows search_rows
	WHERE search_rows.rack_code = rs.rack_code
	  AND (search_rows.rack_code LIKE ? OR search_rows.CARTONBAR LIKE ? OR search_rows.DDBH LIKE ?)
)
`
		likeKeyword := "%" + keyword + "%"
		args = append(args, likeKeyword, likeKeyword, likeKeyword)
	}

	query += `ORDER BY rs.rack_code`

	var rows []warehouseKTPRackRow
	if err := db.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	racks := make([]types.WarehouseKTPRack, 0, len(rows))
	for _, row := range rows {
		status := "empty"
		if row.TotalQty > 0 {
			status = "occupied"
		}

		racks = append(racks, types.WarehouseKTPRack{
			RackCode:              strings.TrimSpace(row.RackCode),
			DDBH:                  strings.TrimSpace(row.DDBH),
			CurrentCode:           nil,
			Status:                status,
			TotalQty:              row.TotalQty,
			InboundCartonCount:    row.InboundCartonCount,
			RecycleCartonCount:    row.RecycleCartonCount,
			InspectionCartonCount: row.InspectionCartonCount,
			CodebarCount:          row.CodebarCount,
			SampleCodes:           compactSampleCodes(row.SampleCode1, row.SampleCode2, row.SampleCode3),
		})
	}

	return racks, nil
}

func (s *WarehouseKTPService) queryRackSignature() ([]warehouseKTPRackSignatureRow, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	query := fmt.Sprintf(`
WITH rack_rows AS (
	SELECT
		LTRIM(RTRIM(KVBH)) AS rack_code,
		LTRIM(RTRIM(CARTONBAR)) AS CARTONBAR,
		LTRIM(RTRIM(DDBH)) AS DDBH,
		LTRIM(RTRIM(SB)) AS SB,
		ISNULL(Qty, 0) AS qty
	FROM YWCP WITH (NOLOCK)
	WHERE %s
	  AND KCBH = 'A2'
	  AND KVBH IS NOT NULL
	  AND LTRIM(RTRIM(KVBH)) <> ''
	  AND LTRIM(RTRIM(KVBH)) <> '0'
	  AND CARTONBAR IS NOT NULL
	  AND LTRIM(RTRIM(CARTONBAR)) <> ''
),
carton_rows AS (
	SELECT
		rack_code,
		CARTONBAR,
		MAX(qty) AS qty
	FROM rack_rows
	GROUP BY rack_code, CARTONBAR
),
rack_summary AS (
	SELECT
		rack_code,
		SUM(qty) AS total_qty,
		COUNT(DISTINCT CARTONBAR) AS codebar_count
	FROM carton_rows
	GROUP BY rack_code
),
order_cartons AS (
	SELECT
		rack_code,
		DDBH,
		CARTONBAR,
		MAX(qty) AS qty
	FROM rack_rows
	WHERE DDBH IS NOT NULL AND DDBH <> ''
	GROUP BY rack_code, DDBH, CARTONBAR
),
order_totals AS (
	SELECT
		rack_code,
		DDBH,
		SUM(qty) AS order_total_qty
	FROM order_cartons
	GROUP BY rack_code, DDBH
),
ranked_orders AS (
	SELECT
		rack_code,
		DDBH,
		ROW_NUMBER() OVER (
			PARTITION BY rack_code
			ORDER BY order_total_qty DESC, DDBH
		) AS rn
	FROM order_totals
),
selected_orders AS (
	SELECT rack_code, DDBH
	FROM ranked_orders
	WHERE rn = 1
),
selected_order_cartons AS (
	SELECT
		r.rack_code,
		r.CARTONBAR,
		s.DDBH,
		MAX(CASE WHEN r.SB = '1' THEN 1 ELSE 0 END) AS inbound_carton,
		MAX(CASE WHEN r.SB = '2' THEN 1 ELSE 0 END) AS recycle_carton,
		MAX(CASE WHEN r.SB = '4' THEN 1 ELSE 0 END) AS inspection_carton
	FROM rack_rows r
	INNER JOIN selected_orders s
		ON s.rack_code = r.rack_code
		AND s.DDBH = r.DDBH
	GROUP BY r.rack_code, r.CARTONBAR, s.DDBH
),
order_details AS (
	SELECT
		rack_code,
		DDBH,
		SUM(inbound_carton) AS inbound_carton_count,
		SUM(recycle_carton) AS recycle_carton_count,
		SUM(inspection_carton) AS inspection_carton_count
	FROM selected_order_cartons
	GROUP BY rack_code, DDBH
)
SELECT
	rs.rack_code,
	od.DDBH,
	rs.total_qty,
	ISNULL(od.inbound_carton_count, 0) AS inbound_carton_count,
	ISNULL(od.recycle_carton_count, 0) AS recycle_carton_count,
	ISNULL(od.inspection_carton_count, 0) AS inspection_carton_count,
	rs.codebar_count
FROM rack_summary rs
LEFT JOIN order_details od ON od.rack_code = rs.rack_code
ORDER BY rs.rack_code`, warehouseKTPEligibleSBCondition)

	var rows []warehouseKTPRackSignatureRow
	if err := db.Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *WarehouseKTPService) findDDBHByCarton(db *gorm.DB, cartonBar string) (string, error) {
	query := fmt.Sprintf(`
SELECT TOP 1 LTRIM(RTRIM(DDBH)) AS ddbh
FROM YWCP WITH (NOLOCK)
WHERE KCBH = 'A2'
  AND %s
  AND CARTONBAR IS NOT NULL
  AND LTRIM(RTRIM(CARTONBAR)) = ?
  AND DDBH IS NOT NULL
  AND LTRIM(RTRIM(DDBH)) <> ''`, warehouseKTPEligibleSBCondition)

	var ddbh string
	if err := db.Raw(query, cartonBar).Scan(&ddbh).Error; err != nil {
		return "", err
	}
	ddbh = strings.TrimSpace(ddbh)
	if ddbh == "" {
		return "", ErrWarehouseKTPCartonNotFound
	}

	return ddbh, nil
}

func (s *WarehouseKTPService) getMoveOrderByDDBH(db *gorm.DB, ddbh string, scannedCartonBar string) (types.WarehouseKTPMoveOrderInfo, error) {
	totalQuery := fmt.Sprintf(`
WITH carton_rows AS (
	SELECT
		LTRIM(RTRIM(CARTONBAR)) AS CARTONBAR,
		MAX(ISNULL(Qty, 0)) AS qty
	FROM YWCP WITH (NOLOCK)
	WHERE KCBH = 'A2'
	  AND %s
	  AND LTRIM(RTRIM(DDBH)) = ?
	  AND CARTONBAR IS NOT NULL
	  AND LTRIM(RTRIM(CARTONBAR)) <> ''
	GROUP BY LTRIM(RTRIM(CARTONBAR))
)
SELECT
	SUM(qty) AS total_qty,
	COUNT(*) AS codebar_count
FROM carton_rows`, warehouseKTPEligibleSBCondition)

	var total warehouseKTPOrderTotalRow
	if err := db.Raw(totalQuery, ddbh).Scan(&total).Error; err != nil {
		return types.WarehouseKTPMoveOrderInfo{}, err
	}
	if total.CodebarCount == 0 {
		return types.WarehouseKTPMoveOrderInfo{}, ErrWarehouseKTPOrderNotFound
	}

	rackQuery := fmt.Sprintf(`
WITH carton_rows AS (
	SELECT
		LTRIM(RTRIM(KVBH)) AS rack_code,
		LTRIM(RTRIM(CARTONBAR)) AS CARTONBAR,
		MAX(ISNULL(Qty, 0)) AS qty
	FROM YWCP WITH (NOLOCK)
	WHERE KCBH = 'A2'
	  AND %s
	  AND LTRIM(RTRIM(DDBH)) = ?
	  AND CARTONBAR IS NOT NULL
	  AND LTRIM(RTRIM(CARTONBAR)) <> ''
	  AND KVBH IS NOT NULL
	  AND LTRIM(RTRIM(KVBH)) <> ''
	  AND LTRIM(RTRIM(KVBH)) <> '0'
	GROUP BY LTRIM(RTRIM(KVBH)), LTRIM(RTRIM(CARTONBAR))
)
SELECT
	rack_code,
	SUM(qty) AS total_qty,
	COUNT(*) AS codebar_count
FROM carton_rows
GROUP BY rack_code
ORDER BY rack_code`, warehouseKTPEligibleSBCondition)

	var rackRows []warehouseKTPOrderRackRow
	if err := db.Raw(rackQuery, ddbh).Scan(&rackRows).Error; err != nil {
		return types.WarehouseKTPMoveOrderInfo{}, err
	}

	currentRacks := make([]types.WarehouseKTPOrderRack, 0, len(rackRows))
	for _, row := range rackRows {
		currentRacks = append(currentRacks, types.WarehouseKTPOrderRack{
			RackCode:     strings.TrimSpace(row.RackCode),
			TotalQty:     row.TotalQty,
			CodebarCount: row.CodebarCount,
		})
	}

	return types.WarehouseKTPMoveOrderInfo{
		ScannedCartonBar: strings.TrimSpace(scannedCartonBar),
		DDBH:             strings.TrimSpace(ddbh),
		TotalQty:         total.TotalQty,
		CodebarCount:     total.CodebarCount,
		CurrentRacks:     currentRacks,
	}, nil
}

func compactSampleCodes(values ...*string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}
		trimmed := strings.TrimSpace(*value)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func normalizeWarehouseKTPRackCode(rackCode string) string {
	return strings.ToUpper(strings.TrimSpace(rackCode))
}
