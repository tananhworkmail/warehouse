package services

import (
	"fmt"
	"strings"
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

// N31 purchasing report.
// SearchN31 returns the 20 columns defined by the N31 purchasing-system report.
func (s *WarehouseKVTService) SearchN31(po string, orderNo string) ([]types.SearchN31, error) {
	po = strings.TrimSpace(po)
	orderNo = strings.TrimSpace(orderNo)
	if po == "" && orderNo == "" {
		return nil, fmt.Errorf("N31 PO or order number is required")
	}

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	query := `
WITH FilteredOrders AS (
	SELECT
		DDZL.DDBH,
		DDZL.BUYNO,
		DDZL.KHPO,
		DDZL.Pairs,
		DDZL.DDGB
	FROM DDZL WITH (NOLOCK)
	WHERE DDZL.GSBH = 'VDH'
	  AND (? = '' OR DDZL.KHPO LIKE ?)
	  AND (? = '' OR DDZL.DDBH LIKE ?)
),
UseStock AS (
	SELECT
		CGKCUSES.ZLBH,
		CGKCUSES.CLBH,
		SUM(CGKCUSES.Qty) AS Qty
	FROM CGKCUSES WITH (NOLOCK)
	INNER JOIN FilteredOrders FO ON FO.DDBH = CGKCUSES.ZLBH
	WHERE CGKCUSES.GSBH = 'VDH'
	GROUP BY CGKCUSES.ZLBH, CGKCUSES.CLBH
),
PurchaseInfo AS (
	SELECT
		CGZLSS.ZLBH,
		CGZLSS.CLBH,
		SUM(CGZLSS.Qty) AS CGQty,
		MAX(CGZL.CGNO) AS CGNO
	FROM CGZLSS WITH (NOLOCK)
	INNER JOIN FilteredOrders FO ON FO.DDBH = CGZLSS.ZLBH
	LEFT JOIN CGZLS WITH (NOLOCK)
		ON CGZLS.CGNO = CGZLSS.CGNO
		AND CGZLS.CLBH = CGZLSS.CLBH
	LEFT JOIN CGZL WITH (NOLOCK) ON CGZL.CGNO = CGZLSS.CGNO
	WHERE CGZL.CGLX IN ('1', '2', '5')
	GROUP BY CGZLSS.ZLBH, CGZLSS.CLBH
),
ReceiptInfo AS (
	SELECT
		KCRKS.CGBH AS ZLBH,
		KCRKS.CLBH,
		SUM(KCRKS.Qty) AS RKQty
	FROM KCRKS WITH (NOLOCK)
	INNER JOIN FilteredOrders FO ON FO.DDBH = KCRKS.CGBH
	LEFT JOIN KCRK WITH (NOLOCK) ON KCRK.RKNO = KCRKS.RKNO
	WHERE ISNULL(KCRK.SFL, '') <> 'THU HOI'
	GROUP BY KCRKS.CGBH, KCRKS.CLBH
),
InvoiceInfo AS (
	SELECT
		HG_INVS.CGBH AS ZLBH,
		HG_INVS.CLBH,
		MAX(HG_INV.DOCNO) AS Invoice
	FROM HG_INVS WITH (NOLOCK)
	INNER JOIN FilteredOrders FO ON FO.DDBH = HG_INVS.CGBH
	LEFT JOIN HG_INV WITH (NOLOCK) ON HG_INV.RKNO = HG_INVS.RKNO
	GROUP BY HG_INVS.CGBH, HG_INVS.CLBH
),
Materials AS (
	SELECT
		ISNULL(FO.BUYNO, '') AS BUYNO,
		ISNULL(FO.KHPO, '') AS KHPO,
		ZLZLS2.ZLBH,
		ISNULL(XXZL.Article, '') AS Article,
		ISNULL(XXZL.XieMing, '') AS XieMing,
		ISNULL(FO.Pairs, 0) AS Pairs,
		ZLZLS2.CLBH,
		ISNULL(CLZL.YWPM, '') AS YWPM,
		ISNULL(CLZL.DWBH, '') AS DWBH,
		ISNULL(CLZL.CQDH, '') AS CQDH,
		CASE
			WHEN DDZLTW.DDZT <> 'C' AND ISNULL(SCZL.IsCGZLS, '') <> 'N'
				THEN ISNULL(SUM(ZLZLS2.CLSL), 0)
			ELSE 0
		END AS CLSL,
		ISNULL(PurchaseInfo.CGQty, 0) AS CGQty,
		ISNULL(ReceiptInfo.RKQty, 0) AS RKQty,
		ISNULL(UseStock.Qty, 0) AS UseStock,
		ISNULL(PurchaseInfo.CGNO, '') AS CGNO,
		ISNULL(FO.DDGB, '') AS DDGB,
		ISNULL(MAX(ZSZL.ZSDH), '') AS ZSBH,
		ISNULL(MAX(ZSZL.ZSYWJC), '') AS ZSYWJC,
		ISNULL(InvoiceInfo.Invoice, '') AS Invoice
	FROM ZLZLS2 WITH (NOLOCK)
	INNER JOIN FilteredOrders FO ON FO.DDBH = ZLZLS2.ZLBH
	LEFT JOIN DDZLTW WITH (NOLOCK) ON DDZLTW.DDBH = ZLZLS2.ZLBH
	LEFT JOIN SCZL WITH (NOLOCK) ON SCZL.SCBH = ZLZLS2.ZLBH
	LEFT JOIN CLZL WITH (NOLOCK) ON CLZL.CLDH = ZLZLS2.CLBH
	LEFT JOIN XXZL WITH (NOLOCK)
		ON XXZL.XieXing = DDZLTW.XieXing
		AND XXZL.SheHao = DDZLTW.SheHao
	LEFT JOIN XXBWFL WITH (NOLOCK)
		ON XXBWFL.XieXing = XXZL.XieXing
		AND XXBWFL.BWBH = ZLZLS2.BWBH
	LEFT JOIN XXBWFLS WITH (NOLOCK) ON XXBWFLS.FLBH = XXBWFL.FLBH
	LEFT JOIN UseStock
		ON UseStock.ZLBH = ZLZLS2.ZLBH
		AND UseStock.CLBH = ZLZLS2.CLBH
	LEFT JOIN PurchaseInfo
		ON PurchaseInfo.ZLBH = ZLZLS2.ZLBH
		AND PurchaseInfo.CLBH = ZLZLS2.CLBH
	LEFT JOIN ReceiptInfo
		ON ReceiptInfo.ZLBH = ZLZLS2.ZLBH
		AND ReceiptInfo.CLBH = ZLZLS2.CLBH
	LEFT JOIN InvoiceInfo
		ON InvoiceInfo.ZLBH = ZLZLS2.ZLBH
		AND InvoiceInfo.CLBH = ZLZLS2.CLBH
	LEFT JOIN ZSZL WITH (NOLOCK) ON ZSZL.ZSDH = ZLZLS2.CSBH
	WHERE ZLZLS2.CLBH NOT LIKE 'W%'
	  AND NOT EXISTS (
			SELECT 1
			FROM KCSAFE WITH (NOLOCK)
			WHERE KCSAFE.CLBH = ZLZLS2.CLBH
	  )
	  AND ZLZLS2.ZMLB = 'N'
	  AND CLZL.CQDH = 'VN'
	  AND ZLZLS2.CLSL <> 0
	  AND (XXBWFLS.DFL <> 'G' OR XXBWFLS.DFL IS NULL)
	GROUP BY
		FO.BUYNO,
		FO.KHPO,
		FO.Pairs,
		ZLZLS2.ZLBH,
		XXZL.Article,
		XXZL.XieMing,
		ZLZLS2.CLBH,
		CLZL.YWPM,
		CLZL.DWBH,
		CLZL.CQDH,
		DDZLTW.DDZT,
		SCZL.IsCGZLS,
		PurchaseInfo.CGQty,
		PurchaseInfo.CGNO,
		ReceiptInfo.RKQty,
		UseStock.Qty,
		FO.DDGB,
		InvoiceInfo.Invoice
)
SELECT
	Materials.BUYNO,
	Materials.KHPO,
	Materials.ZLBH,
	Materials.Article,
	Materials.XieMing,
	Materials.Pairs,
	Materials.CLBH,
	Materials.YWPM,
	Materials.DWBH,
	Materials.CQDH,
	Materials.CLSL,
	Materials.CGQty,
	Materials.RKQty,
	Materials.UseStock,
	Materials.CGNO,
	Materials.DDGB,
	Materials.ZSBH,
	Materials.ZSYWJC,
	Materials.Invoice,
	ISNULL(Rack.MAKE, '') AS MAKE
FROM Materials
OUTER APPLY (
	SELECT TOP 1 RFSS.MAKE
	FROM KCRKScan_RFSS RFSS WITH (NOLOCK)
	WHERE RFSS.CLBH = Materials.CLBH
	  AND RFSS.RemQty > 0
	  AND ISNULL(RFSS.MAKE, '') <> ''
	  AND (
			RFSS.Memo_RY = Materials.ZLBH
			OR ',' + REPLACE(ISNULL(RFSS.Memo_RY, ''), ' ', '') + ','
				LIKE '%,' + REPLACE(Materials.ZLBH, ' ', '') + ',%'
	  )
	ORDER BY RFSS.CFMDate DESC, RFSS.SCNO DESC
) Rack
ORDER BY Materials.ZLBH, Materials.CQDH, Materials.CLBH
`

	poPattern := "%" + po + "%"
	orderPattern := "%" + orderNo + "%"
	var result []types.SearchN31
	if err := db.Raw(query, po, poPattern, orderNo, orderPattern).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("search N31 failed: %w", err)
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
