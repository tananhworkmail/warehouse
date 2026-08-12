package services

import (
	"fmt"
	"time"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"
)

type WarehouseKVTService struct {
	*BaseService
}

var WarehouseKVT = &WarehouseKVTService{}

func (s *WarehouseKVTService) GetRackKVT(rackCode string) ([]types.RackKVT, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	query := `
	SELECT * FROM MAKE WHERE DONVI ='KVT' AND MAKE LIKE ? AND LEN(MAKE) > 3
	`

	likeRackcode := "%" + rackCode + "%"
	var result []types.RackKVT
	if err := db.Raw(query, likeRackcode).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (s *WarehouseKVTService) GetRackTotalColumnKVT() ([]types.RackTotalColumn, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	query := `
SELECT 
    MK.MAKE AS MAKE3,
    ISNULL(ROUND(SUM(K.RemQty), 1), 0) AS TotalColumn,
    MK.[Limit] AS [Limit]
FROM MAKE MK
LEFT JOIN KCRKScan_RFSS K
    ON LEFT(K.MAKE, 3) = MK.MAKE
    AND K.RemQty > 0
    AND EXISTS (
        SELECT 1
        FROM KCRKScan_RF RF
        WHERE RF.SCNO = K.SCNO
          AND RF.GSBH = 'VDH'
    )
WHERE LEN(LTRIM(RTRIM(MK.MAKE))) = 3
GROUP BY 
    MK.MAKE,
    MK.[Limit]
ORDER BY 
    MK.MAKE;
	`

	var result []types.RackTotalColumn
	if err := db.Raw(query).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WarehouseKVTService) GetTonKhoTangKVT() ([]types.RackTotalTang, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	query := `
	SELECT 
    LEFT(RFSS.MAKE, 5) AS MAKE5,
    ROUND(SUM(RFSS.RemQty), 1) AS TotalTang
FROM KCRKScan_RFSS RFSS
INNER JOIN KCRKScan_RF RF
    ON RFSS.SCNO = RF.SCNO
WHERE RFSS.RemQty > 0
  AND RF.GSBH = 'VDH'
GROUP BY LEFT(RFSS.MAKE, 5)
`

	var result []types.RackTotalTang
	if err := db.Raw(query).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (s *WarehouseKVTService) GetTonKhoInRackKVT(rackPrefix string) ([]types.TonkhoRack, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	query := `
SELECT 
    RFSS.MAKE,
    ROUND(SUM(RFSS.RemQty), 1, 1) AS TonkhoRack
FROM KCRKScan_RFSS RFSS
INNER JOIN KCRKScan_RF RF
    ON RFSS.SCNO = RF.SCNO
WHERE RFSS.RemQty > 0
  AND RFSS.MAKE LIKE ?
  AND RF.GSBH = 'VDH'
GROUP BY RFSS.MAKE
	`

	likeValue := "%" + rackPrefix + "%"
	var result []types.TonkhoRack
	if err := db.Raw(query, likeValue).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (s *WarehouseKVTService) GetRackInformationKVT(rackNo string) ([]types.WarehouseRackInformationKVT, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	query := `
SELECT DISTINCT KCRKScan_RFSS.SCNO,KCRKScan_RFS.CGNO,CGZL.ZSBH,ZSZL.zsywjc,KCRKScan_RF.GSBH,KCRKScan_RF.CKBH,KCRKScan_RFSS.CLBH,CLZL.YWPM,CLZL.DWBH, 
       KCRKScan_RFSS.Qty,isnull(KCRKScan_RFSS.RemQty,KCRKScan_RFSS.Qty) RemQty,Round(KCRKScan_RFSS.Qty-isnull(KCRKScan_RFSS.RemQty,KCRKScan_RFSS.Qty),4) AS DQty,KCRKScan_RFSS.TAGID, 
       KCRKScan_RFSS.Pack,KCRKScan_RFSS.Memo_RY,KCRKScan_RFSS.Memo_Article, KCRKScan_RFSS.MAKE AS KCBH, KCRKScan_RFSS.barcode,KCRKScan_RF.MEMO,KCRKScan_RF.DOCNO,KCRKScan_RFS.IDCheck 
   ,KCRKScan_RFSS.CFMID,KCRKScan_RFSS.CFMDATE,
   CASE 
  WHEN KCRKScan_RFSS.CFMDATE IS NOT NULL 
       AND DATEDIFF(DAY, KCRKScan_RFSS.CFMDATE, GETDATE()) > 180 
  THEN 1
  ELSE 0
END AS is_over_180
FROM KCRKScan_RFSS 
left join KCRKScan_RF ON KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO 
left join KCRKScan_RFS ON KCRKScan_RFS.SCNO = KCRKScan_RFSS.SCNO AND KCRKScan_RFS.CLBH = KCRKScan_RFSS.CLBH 
left join CGZL ON KCRKScan_RFS.CGNO = CGZL.CGNO 
left join ZSZL ON CGZL.ZSBH = ZSZL.ZSDH 
left join CLZL ON KCRKScan_RFS.CLBH = CLZL.CLDH 
left join KCZLS ON KCRKScan_RF.CKBH=KCZLS.CKBH AND KCRKScan_RFSS.CLBH=KCZLS.CLBH 
left join KCLLScan_RFSS ON KCLLScan_RFSS.barcode=KCRKScan_RFSS.barcode AND KCLLScan_RFSS.CLBH=KCRKScan_RFSS.CLBH 
left join KCLLScan_RFSSS ON KCLLScan_RFSSS.barcode=KCRKScan_RFSS.barcode AND KCLLScan_RFSSS.CLBH=KCRKScan_RFSS.CLBH 
WHERE KCRKScan_RF.GSBH='VDH' AND  KCRKScan_RFSS.MAKE is not null AND KCRKScan_RFSS.RemQty >0 AND KCRKScan_RFSS.MAKE = ?
ORDER BY KCRKScan_RFSS.SCNO DESC, KCRKScan_RFS.CGNO DESC, KCRKScan_RFSS.CLBH DESC, KCRKScan_RFSS.Pack ASC
	`

	var result []types.WarehouseRackInformationKVT
	if err := db.Raw(query, rackNo).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (s *WarehouseKVTService) Get3DayRackInformationKVT(rackNo string) ([]types.Warehouse3DayRackInformationKVT, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	query := `
		select k.SCNO, k.CLBH, k.MAKE
		from ProductionPlan p
		join KCRKScan_RFSS k
			on ',' + k.Memo_RY + ',' like '%,' + p.RY + ',%'
		where p.PlanType = '3-day'
		and convert(date, p.PlanDate) between convert(date, getdate()+1)
											and convert(date, getdate()+3)
		and k.RemQty>0
		and k.MAKE=?
		order by p.PlanDate
	`

	var result []types.Warehouse3DayRackInformationKVT
	if err := db.Raw(query, rackNo).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// //////
func (s *WarehouseKVTService) SearchByCLBHKVT(CLBH string) ([]types.SearchByCLBH, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()
	query := `SELECT DISTINCT KCRKScan_RFSS.SCNO,KCRKScan_RFS.CGNO,CGZL.ZSBH,ZSZL.zsywjc,KCRKScan_RF.GSBH,KCRKScan_RF.CKBH,KCRKScan_RFSS.CLBH,CLZL.YWPM,CLZL.DWBH, 
       KCRKScan_RFSS.Qty,isnull(KCRKScan_RFSS.RemQty,KCRKScan_RFSS.Qty) RemQty,Round(KCRKScan_RFSS.Qty-isnull(KCRKScan_RFSS.RemQty,KCRKScan_RFSS.Qty),4) AS DQty,KCRKScan_RFSS.TAGID, 
       KCRKScan_RFSS.Pack,KCRKScan_RFSS.Memo_RY,KCRKScan_RFSS.Memo_Article, KCRKScan_RFSS.MAKE AS KCBH, KCRKScan_RFSS.barcode,KCRKScan_RF.MEMO,KCRKScan_RF.DOCNO,KCRKScan_RFS.IDCheck 
   ,KCRKScan_RFSS.CFMID,KCRKScan_RFSS.CFMDATE 
FROM KCRKScan_RFSS 
left join KCRKScan_RF ON KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO 
left join KCRKScan_RFS ON KCRKScan_RFS.SCNO = KCRKScan_RFSS.SCNO AND KCRKScan_RFS.CLBH = KCRKScan_RFSS.CLBH 
left join CGZL ON KCRKScan_RFS.CGNO = CGZL.CGNO 
left join ZSZL ON CGZL.ZSBH = ZSZL.ZSDH 
left join CLZL ON KCRKScan_RFS.CLBH = CLZL.CLDH 
left join KCZLS ON KCRKScan_RF.CKBH=KCZLS.CKBH AND KCRKScan_RFSS.CLBH=KCZLS.CLBH 
left join KCLLScan_RFSS ON KCLLScan_RFSS.barcode=KCRKScan_RFSS.barcode AND KCLLScan_RFSS.CLBH=KCRKScan_RFSS.CLBH 
left join KCLLScan_RFSSS ON KCLLScan_RFSSS.barcode=KCRKScan_RFSS.barcode AND KCLLScan_RFSSS.CLBH=KCRKScan_RFSS.CLBH 
WHERE KCRKScan_RF.GSBH='VDH' AND  KCRKScan_RFSS.MAKE is not null AND KCRKScan_RFSS.RemQty >0 AND KCRKScan_RFSS.CLBH LIKE ?
ORDER BY KCRKScan_RFSS.SCNO DESC, KCRKScan_RFS.CGNO DESC, KCRKScan_RFSS.CLBH DESC, KCRKScan_RFSS.Pack ASC`
	likeCLBH := "%" + CLBH + "%"
	var result []types.SearchByCLBH
	if err := db.Raw(query, likeCLBH).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// func (s *WarehouseKVTService) SearchDDBHKVT(DDBH string) ([]types.SearchTotalDDBH, error) {
// 	db, err := database.LYS_ERP_Connection()
// 	if err != nil {
// 		fmt.Println("Database connection error:", err)
// 		return nil, err
// 	}

// 	dbInstance, err := db.DB()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get DB instance: %w", err)
// 	}
// 	defer dbInstance.Close()
// 	query := `SELECT

//   (SELECT SUM(QtyIn)
//    FROM SMZL_MAKE
//    WHERE CODEBAR IN (
//      SELECT CODEBAR FROM smddss WHERE DDBH LIKE ?
//    )
//   ) AS TotalQtyIn,

//   (SELECT SUM(QtyOut)
//    FROM SMZL_MAKE
//    WHERE CODEBAR IN (
//      SELECT CODEBAR FROM smddss WHERE DDBH LIKE ?
//    )
//   ) AS TotalQtyOut,

//   (SELECT SUM(Qty)
//    FROM smdd
//    WHERE DDBH LIKE ?
//  AND GXLB = 'O'
//   ) AS TotalOrderQty;

// `
// 	likeDDBH := "%" + DDBH + "%"
// 	var result []types.SearchTotalDDBH
// 	if err := db.Raw(query, likeDDBH, likeDDBH, likeDDBH).Scan(&result).Error; err != nil {
// 		return nil, err
// 	}

//		return result, nil
//	}
func (s *WarehouseKVTService) GetTempHumidity(
	buildingNo string,
	queryDate string, // "2025-12-19"
) (types.GetTempHumidity, error) {

	db, err := database.TempHumidity_Connection()
	if err != nil {
		return types.GetTempHumidity{}, err
	}

	var result types.GetTempHumidity
	slotClause := ""
	if shouldUseTempHumiditySlotRows(queryDate) {
		if err := ensureTempHumidityFakeSlotRows([]string{buildingNo}); err != nil {
			return types.GetTempHumidity{}, err
		}
		slotClause = " AND DATEPART(minute, RecordTime) IN (0, 30) AND DATEPART(second, RecordTime) = 0 AND DATEPART(millisecond, RecordTime) = 0"
	}

	err = db.Raw(`
		SELECT TOP 1
			DeviceName,
			DeviceAddr,
			Hum,
			Tem,
			RecordTime
		FROM tbhistory
		WHERE DeviceName LIKE ?
		  AND RecordTime >= ?
		  AND RecordTime < DATEADD(day, 1, ?)
	`+slotClause+`
		ORDER BY RecordTime DESC
	`,
		"%"+buildingNo+"%",
		queryDate, // 2025-12-19 00:00:00
		queryDate, // +1 day
	).Scan(&result).Error

	if err != nil {
		return types.GetTempHumidity{}, err
	}

	return result, nil
}
func (s *WarehouseKVTService) GetHumidityAlertKVT01(
	queryDate string,
	startTime string,
	endTime string,
) ([]types.GetHumidityAlert, error) {

	db, err := database.TempHumidity_Connection()
	if err != nil {
		return nil, err
	}

	var result []types.GetHumidityAlert

	query := `
	SELECT 
		DeviceName,
		Hum,
		Tem as Temp,
		CONVERT(varchar, alert_date, 23) as Alert_Date,
		CONVERT(varchar, alert_time, 108) as Alert_Time,
		CONVERT(varchar, RecordTime, 120) as RecordTime
	FROM TBHISTORY_ALERT 
	WHERE ALERT_DATE = ? 
	AND DeviceName = 'KVT01'
	`

	args := []interface{}{queryDate}

	// 👇 thêm filter time nếu có
	if startTime != "" && endTime != "" {
		query += " AND alert_time BETWEEN ? AND ?"
		args = append(args, startTime, endTime)
	}

	query += " ORDER BY RecordTime DESC"

	err = db.Raw(query, args...).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (s *WarehouseKVTService) GetHumidityAlertKVT02(
	queryDate string,
	startTime string,
	endTime string,
) ([]types.GetHumidityAlert, error) {

	db, err := database.TempHumidity_Connection()
	if err != nil {
		return nil, err
	}
	var result []types.GetHumidityAlert

	query := `
	SELECT 
		DeviceName,
		Hum,
		Tem as Temp,
		CONVERT(varchar, alert_date, 23) as Alert_Date, -- yyyy-MM-dd
		CONVERT(varchar, alert_time, 108) as Alert_Time, -- HH:mm:ss
		CONVERT(varchar, RecordTime, 120) as RecordTime -- yyyy-MM-dd HH:mm:ss
	FROM TBHISTORY_ALERT 
	WHERE ALERT_DATE = ? AND DeviceName ='KVT02'
`

	args := []interface{}{queryDate}

	// 👇 thêm filter time nếu có
	if startTime != "" && endTime != "" {
		query += " AND alert_time BETWEEN ? AND ?"
		args = append(args, startTime, endTime)
	}

	query += " ORDER BY RecordTime DESC"

	err = db.Raw(query, args...).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (s *WarehouseKVTService) GetHumidityAlertKVT03(
	queryDate string,
	startTime string,
	endTime string,
) ([]types.GetHumidityAlert, error) {

	db, err := database.TempHumidity_Connection()
	if err != nil {
		return nil, err
	}

	var result []types.GetHumidityAlert

	query := `
	SELECT 
		DeviceName,
		Hum,
		Tem as Temp,
		CONVERT(varchar, alert_date, 23) as Alert_Date, -- yyyy-MM-dd
		CONVERT(varchar, alert_time, 108) as Alert_Time, -- HH:mm:ss
		CONVERT(varchar, RecordTime, 120) as RecordTime -- yyyy-MM-dd HH:mm:ss
	FROM TBHISTORY_ALERT 
	WHERE ALERT_DATE = ? AND DeviceName ='KVT03'
	`

	args := []interface{}{queryDate}

	// 👇 thêm filter time nếu có
	if startTime != "" && endTime != "" {
		query += " AND alert_time BETWEEN ? AND ?"
		args = append(args, startTime, endTime)
	}

	query += " ORDER BY RecordTime DESC"

	err = db.Raw(query, args...).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WarehouseKVTService) GetTempHumidityByDevices() ([]types.GetTempHumidity, error) {
	db, err := database.TempHumidity_Connection()
	if err != nil {
		return nil, err
	}

	var result []types.GetTempHumidity
	slotClause := ""
	if shouldUseTempHumiditySlotRows("") {
		if err := ensureTempHumidityFakeSlotRows([]string{"KVT02", "KVT03"}); err != nil {
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
			WHERE DeviceName IN ('KVT01','KVT02','KVT03')
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

func (s *WarehouseKVTService) GetLatestTempHumidityLogByDevices(
	queryDate string,
	startTime string,
	endTime string,
) ([]types.GetHumidityAlert, error) {

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
		if err := ensureTempHumidityFakeSlotRows([]string{"KVT02", "KVT03"}); err != nil {
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
			WHERE DeviceName IN ('KVT02','KVT03')
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
	`

	query += " ORDER BY DeviceName, HalfHourBucket DESC"

	err = db.Raw(query, args...).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WarehouseKVTService) GetTempHumidityLaboratory(date string) ([]types.GetTempHumidity, error) {

	db, err := database.TempHumidity_Connection()
	if err != nil {
		return nil, err
	}

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	var result []types.GetTempHumidity

	err = db.Raw(`
        SELECT
            DeviceName,
            DeviceAddr,
            Hum,
            Tem,
            RecordTime
        FROM tbhistory
        WHERE DeviceName IN (
            'ThiNghiem9517_L4',
            'ThiNghiem9527_L4'
        )
        AND RecordTime >= CAST(? AS DATE)
        AND RecordTime < DATEADD(DAY, 1, CAST(? AS DATE))
        ORDER BY RecordTime DESC
    `, date, date).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
func (s *WarehouseKVTService) Get3And180DayKVT() ([]types.Warehouse3And180DayKVT, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	query := `
SELECT
    MAKE,
    MAX([3day]) AS Day3,
    MAX([180day]) AS Day180
FROM (
    -- 3-day
    SELECT
        LEFT(k.MAKE, 5) AS MAKE,
        1 AS [3day],
        0 AS [180day]
    FROM ProductionPlan p
    JOIN KCRKScan_RFSS k
        ON ',' + k.Memo_RY + ',' LIKE '%,' + p.RY + ',%'
    WHERE p.PlanType = '3-day'
      AND CONVERT(date, p.PlanDate)
            BETWEEN CONVERT(date, DATEADD(DAY, 1, GETDATE()))
                AND CONVERT(date, DATEADD(DAY, 3, GETDATE()))
      AND k.RemQty > 0

    UNION ALL

    -- 180-day
    SELECT
        LEFT(a.MAKE, 5) AS MAKE,
        0 AS [3day],
        1 AS [180day]
    FROM KCRKScan_RFSS a
    INNER JOIN KCRKScan_RF b
        ON a.SCNO = b.SCNO
    WHERE DATEDIFF(DAY, a.CFMDate, GETDATE()) > 180
      AND a.RemQty > 0
      AND a.MAKE IS NOT NULL
      AND b.GSBH = 'VDH'
) AS T
GROUP BY MAKE
ORDER BY MAKE;
`

	var result []types.Warehouse3And180DayKVT
	if err := db.Raw(query).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("query 3 and 180 day KVT failed: %w", err)
	}

	return result, nil
}
