package services

import (
	"fmt"
	"strings"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"
)

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
