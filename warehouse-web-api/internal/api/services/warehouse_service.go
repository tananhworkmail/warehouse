package services

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/entities"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"
)

type WarehouseService struct {
	*BaseService
}

var Warehouse = &WarehouseService{}
func getimg(path string) string {
	var base64Encoding string
	f, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	} else {
		img := []byte(f)
		mimeType := http.DetectContentType(img)
		switch mimeType {
		case "image/jpeg":
			base64Encoding += "data:image/jpeg;base64,"
		case "image/png":
			base64Encoding += "data:image/png;base64,"
		}
		base64Encoding += base64.StdEncoding.EncodeToString(img)
	}
	return base64Encoding
}
// Search RKno
func (s *WarehouseService) GetWareHouseList(req request.GetWarehouseListRequest) ([]types.GetWarehouseList, error) {
	var result []types.GetWarehouseList

	// Kết nối database
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Base query
	query := `
SELECT 
	K.ZLBH,
	K.CLBH,
	K.ZSBH,
	K.COLOR,
	K.TYPE,
	SUM(CASE WHEN C.YN = 2 THEN I.QTY ELSE 0 END) OVER (PARTITION BY K.ZLBH, K.CLBH, K.ZSBH) AS STOCKQTY,
	CASE WHEN C.YN = 2 THEN I.SIZE END AS SIZE,
	CASE WHEN C.YN = 2 THEN I.QTY END AS QTY
FROM WEBQR_KCRK K
JOIN WEBQR_CARTON C ON K.RKNO = C.RKNO
JOIN WEBQR_ITEM I ON C.CARTONQR = I.CARTONQR
WHERE C.YN >= 2

	`

	var args []interface{}

	if req.ZLBH != "" {
		query += " AND K.ZLBH = ?"
		args = append(args, req.ZLBH)
	}
	if req.CLBH != "" {
		query += " AND K.CLBH = ?"
		args = append(args, req.CLBH)
	}
	if req.ZSBH != "" {
		query += " AND K.ZSBH = ?"
		args = append(args, req.ZSBH)
	}
	if req.COLOR != "" {
		query += " AND K.COLOR = ?"
		args = append(args, req.COLOR)
	}
	if req.TYPE != "" {
		query += " AND K.TYPE = ?"
		args = append(args, req.TYPE)
	}

	query += `
	ORDER BY 
		K.ZLBH, K.CLBH, K.ZSBH, I.SIZE
	`

	// Thực thi query
	err = db.Raw(query, args...).Scan(&result).Error
	if err != nil {
		fmt.Println("Query error:", err)
		return nil, err
	}

	return result, nil
}

func (s *WarehouseService) GetRackInformation(rackNo string) ([]types.WarehouseRackInformation, error) {
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
    c.DDBH, 
	d.YSBH,
    d.XieMing, 
    c.XXCC, 
    -- Tinh tong nhap
    SUM(CASE WHEN a.Status = 'IN' THEN a.Qty ELSE 0 END) AS QtyIn, 
    -- Tinh tong xuat
    SUM(CASE WHEN a.Status = 'OUT' THEN a.Qty ELSE 0 END) AS QtyOut, 
    -- Tinh ton kho thuc te
    SUM(CASE WHEN a.Status = 'IN' THEN a.Qty ELSE -a.Qty END) AS TonKho, 
    a.Codebar, 
    a.Ma_Ke, 
    MAX(e.Name) AS Name, -- Dung MAX vi chung ta dang GROUP BY
    'PRS' AS DonVi
FROM SMZL_MAKE a WITH (NOLOCK)
LEFT JOIN smddss c WITH (NOLOCK) ON a.CODEBAR = c.CODEBAR 
LEFT JOIN smdd d WITH (NOLOCK) ON c.DDBH = d.DDBH AND c.GXLB = d.GXLB
LEFT JOIN MAKE e WITH (NOLOCK) ON e.MAKE = a.MA_KE
WHERE c.okCTS > 0 
  AND c.GXLB = 'O'
  AND a.Ma_Ke LIKE ?
  and status in ('IN', 'OUT')
GROUP BY 
    c.DDBH, 
    d.XieMing, 
    c.XXCC, 
    a.Codebar, 
    a.Ma_Ke,
	d.YSBH

-- HAVING SUM(CASE WHEN a.Status = 'IN' THEN a.Qty ELSE -a.Qty END) > 0
ORDER BY c.DDBH, c.XXCC, a.Ma_Ke
	`

	likeRackNo := "%" + rackNo + "%"

	var result []types.WarehouseRackInformation
	if err := db.Raw(query, likeRackNo).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (s *WarehouseService) Get3DayRackInformation(rackNo string) ([]types.Warehouse3DayRackInformation, error) {
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
SELECT DISTINCT
    d.YSBH AS RY,
    a.Ma_Ke
FROM SMZL_MAKE a WITH (NOLOCK)

LEFT JOIN smddss c WITH (NOLOCK) 
    ON a.CODEBAR = c.CODEBAR 

LEFT JOIN smdd d WITH (NOLOCK) 
    ON c.DDBH = d.DDBH 
    AND c.GXLB = d.GXLB

JOIN ProductionPlan p
    ON p.RY = d.YSBH

WHERE 
    c.okCTS > 0
    AND c.GXLB = 'O'
    AND a.Status IN ('IN','OUT')
    AND a.Ma_Ke LIKE ?

    AND p.PlanType = '3-day'
    AND CONVERT(date, p.PlanDate) 
        BETWEEN CONVERT(date, GETDATE()+1)
        AND CONVERT(date, GETDATE()+3)
	`

	likeRackNo := "%" + rackNo + "%"

	var result []types.Warehouse3DayRackInformation
	if err := db.Raw(query, likeRackNo).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (s *WarehouseService) GetTonKhoByRackPrefix(rackPrefix string) ([]types.RackTonKho, error) {
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
    SUM(CASE WHEN a.Status = 'IN' THEN a.Qty ELSE -a.Qty END) AS TonKhoTang,
    a.Ma_Ke
FROM SMZL_MAKE a WITH (NOLOCK)
LEFT JOIN smddss c WITH (NOLOCK) ON a.CODEBAR = c.CODEBAR 
WHERE c.okCTS > 0 
  AND c.GXLB = 'O'
  AND a.Ma_Ke LIKE ?
  and status in ('IN', 'OUT')
GROUP BY a.Ma_Ke
	`

	likeValue := "%" + rackPrefix + "%"
	var result []types.RackTonKho
	if err := db.Raw(query, likeValue).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (s *WarehouseService) GetTonKhoKeByRack() ([]types.RackTonKhoKe, error) {
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
    LEFT(MAKETON.MA_KE, 1) AS MA_KE, -- Lay ky tu dau cua ma ke (vi du: A, B, C)
    MAX(e.[Limit]) AS Limit,        -- Su dung MAX vi cot nay nam ngoai GROUP BY
    SUM(MAKETON.TonKho) AS TonKhoKe
FROM (
    SELECT 
        a.MA_KE,
        -- Tinh ton kho thuc te tu bang smzl_make moi
        SUM(CASE WHEN a.Status = 'IN' THEN a.Qty 
                 WHEN a.Status = 'OUT' THEN -a.Qty 
                 ELSE 0 END) AS TonKho
    FROM SMZL_MAKE a WITH (NOLOCK)
    LEFT JOIN smddss c WITH (NOLOCK) ON a.CODEBAR = c.CODEBAR 
    WHERE ISNULL(c.okCTS, 0) > 0 
      AND c.GXLB = 'O'
      AND a.userdate IS NOT NULL
    GROUP BY a.MA_KE, a.CODEBAR
) AS MAKETON
LEFT JOIN MAKE e ON e.MAKE = LEFT(MAKETON.MA_KE, 1)
GROUP BY 
    LEFT(MAKETON.MA_KE, 1)
-- Loai bo cac khu vuc da het hang (Ton kho = 0)
HAVING SUM(MAKETON.TonKho) > 0
ORDER BY LEFT(MAKETON.MA_KE, 1);`

	var result []types.RackTonKhoKe
	if err := db.Raw(query).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (s *WarehouseService) GetTonKhoTangByRack() ([]types.RackTonKhoTang, error) {
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
    LEFT(MAKETON.MA_KE, 3) AS MA_KE, -- Lay ky tu dau cua ma ke (vi du: A, B, C)
       
    SUM(MAKETON.TonKho) AS TonKhoTang
FROM (
    SELECT 
        a.MA_KE,
        -- Tinh ton kho thuc te tu bang smzl_make moi
        SUM(CASE WHEN a.Status = 'IN' THEN a.Qty 
                 WHEN a.Status = 'OUT' THEN -a.Qty 
                 ELSE 0 END) AS TonKho
    FROM SMZL_MAKE a WITH (NOLOCK)
    LEFT JOIN smddss c WITH (NOLOCK) ON a.CODEBAR = c.CODEBAR 
    WHERE ISNULL(c.okCTS, 0) > 0 
      AND c.GXLB = 'O'
      AND a.userdate IS NOT NULL
    GROUP BY a.MA_KE, a.CODEBAR
) AS MAKETON
LEFT JOIN MAKE e ON e.MAKE = LEFT(MAKETON.MA_KE, 3)
GROUP BY 
    LEFT(MAKETON.MA_KE, 3)
-- Loai bo cac khu vuc da het hang (Ton kho = 0)
HAVING SUM(MAKETON.TonKho) > 0
ORDER BY LEFT(MAKETON.MA_KE, 3);`

	var result []types.RackTonKhoTang
	if err := db.Raw(query).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (s *WarehouseService) SearchByDDBH(DDBH string) ([]types.SearchByDDBH, error) {
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
    c.DDBH, 
    d.XieMing, 
    c.XXCC, 
    -- Tinh tong so luong Nhap
    SUM(CASE WHEN a.Status = 'IN' THEN a.Qty ELSE 0 END) AS QtyIn, 
    -- Tinh tong so luong Xuat
    SUM(CASE WHEN a.Status = 'OUT' THEN a.Qty ELSE 0 END) AS QtyOut, 
    -- Tinh Ton Kho (Nhap - Xuat)
    SUM(CASE WHEN a.Status = 'IN' THEN a.Qty ELSE -a.Qty END) AS TonKho, 
    a.Codebar, 
    a.Ma_Ke
FROM SMZL_MAKE a WITH (NOLOCK)
LEFT JOIN smddss c WITH (NOLOCK) ON a.CODEBAR = c.CODEBAR 
LEFT JOIN smdd d WITH (NOLOCK) ON c.DDBH = d.DDBH AND c.GXLB = d.GXLB
WHERE c.okCTS > 0 
  AND c.GXLB = 'O'
  AND c.DDBH LIKE ?
    and status in ('IN', 'OUT')
GROUP BY 
    c.DDBH, 
    d.XieMing, 
    c.XXCC, 
    a.Codebar, 
    a.Ma_Ke
ORDER BY c.DDBH, c.XXCC, a.Ma_Ke`
	likeDDBH := "%" + DDBH + "%"
	var result []types.SearchByDDBH
	if err := db.Raw(query, likeDDBH).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (s *WarehouseService) SearchDDBH(DDBH string) ([]types.SearchTotalDDBH, error) {
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
	query := `SELECT 
    -- Tong so luong don hang (tinh rieng)
    (SELECT SUM(Qty) 
     FROM smdd WITH (NOLOCK)
     WHERE DDBH LIKE ?
       AND GXLB = 'O') AS TotalOrderQty,

    -- Tong nhap
    SUM(CASE WHEN m.Status = 'IN' THEN m.Qty ELSE 0 END) AS TotalQtyIn,

    -- Tong xuat
    SUM(CASE WHEN m.Status = 'OUT' THEN m.Qty ELSE 0 END) AS TotalQtyOut,

    -- Ton kho
    SUM(CASE WHEN m.Status = 'IN' THEN m.Qty ELSE -m.Qty END) AS TotalBalance

FROM smddss s WITH (NOLOCK)

LEFT JOIN SMZL_MAKE m WITH (NOLOCK) 
    ON s.CODEBAR = m.CODEBAR

WHERE s.DDBH LIKE ?

`
	likeDDBH := "%" + DDBH + "%"
	var result []types.SearchTotalDDBH
	if err := db.Raw(query, likeDDBH, likeDDBH).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
func (s *WarehouseService)UploadSchedule(buildingNo string, file *multipart.FileHeader) error {
	_, err := s.DeleteByWhere(entities.ScheduleFile{}, entities.ScheduleFile{BuildingNo: buildingNo})
    // Kết nối database
    db, err := database.LYS_ERP_Connection()
    if err != nil {
        fmt.Println("Database connection error:", err)
        return err
    }

    dbInstance, err := db.DB()
    if err != nil {
        return fmt.Errorf("failed to get DB instance: %w", err)
    }
    defer dbInstance.Close()

    // Mở file
    fileReader, err := file.Open()
    if err != nil {
        return fmt.Errorf("failed to open uploaded file: %w", err)
    }
    defer fileReader.Close()

    // Đọc nội dung file
    data, err := io.ReadAll(fileReader)
    if err != nil {
        return fmt.Errorf("failed to read file content: %w", err)
    }

    // Tạo đối tượng ScheduleFile
    object := entities.ScheduleFile{
		BuildingNo: buildingNo,
        FileStream: data,
    }

    // Lưu bản ghi vào DB
    if err := s.Create(&object); err != nil {
        return fmt.Errorf("failed to create schedule file: %w", err)
    }

    return nil
}
func (s *WarehouseService) GetSchedule(where interface{}, out interface{}, associations []string, orders ...string) (err error) {
	 db, err := database.LYS_ERP_Connection()
    if err != nil {
        fmt.Println("Database connection error:", err)
        return err
    }

    dbInstance, err := db.DB()
    if err != nil {
        return fmt.Errorf("failed to get DB instance: %w", err)
    }
    defer dbInstance.Close()
	for _, a := range associations {
		db = db.Preload(a)
	}

	db = db.Where(where)
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	return db.Find(out).Error
}
func (s *WarehouseService) GetOrderExistenceStatus(orderNo string) (types.TraceOrderExistenceStatus, error) {
	var result types.TraceOrderExistenceStatus

	// Kết nối database
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return result, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return result, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	// Query an toàn với binding để tránh SQL Injection
	query := `SELECT DDBH, KHPO FROM DDZL WHERE DDBH = ? OR KHPO = ?`
	if err := db.Raw(query, orderNo, orderNo).Scan(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (s *WarehouseService) GetTraceServiceAll(OrderNo string) (types.TraceServiceAllDAta, error) {
	var tmps types.TraceServiceAllDAta
	tmps.Bom, _ = GetBomInformation(OrderNo)
	tmps.Material, tmps.Orders, _ = GetMaterialPo(OrderNo)
	tmps.Invoice, _ = GetInvoice(OrderNo)
	tmps.Customer, _ = GetCustomer(OrderNo)
	tmps.Product, _ = GetProducts(OrderNo)
	return tmps, nil
}
func GetBomInformation(orderNo string) (types.TraceOrderInfo, error) {
	var data types.TraceOrderInfo

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return data, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return data, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	query := `
	select distinct SCZL.ZLBH,XieMing,ShipDate,KHPO,XXZL.ARTICLE,BZCC,kfjc Kfjc,XXZL.XieXing,XXZL.SheHao,XTMH,ywsm Ywsm,YSSM,DAOMH,IMGName
	from SCZL
	left join ZLZL on ZLZL.ZLBH=SCZL.ZLBH
	left join DDZL on ZLZL.ZLBH=DDZL.DDBH
	left join KFZL on KFZL.KFDH=DDZL.KHBH 
	left join XXZL on DDZL.XieXing=XXZL.XieXing and DDZL.SheHao=XXZL.SheHao
	left join LBZLS on LBZLS.LB='06' and LBZLS.LBDH=DDZL.DDGB
	where SCZL.ZLBH=?`

	var result types.TraceOrderHead
	if err := db.Raw(query, orderNo).Scan(&result).Error; err != nil {
		return data, err
	}

	if result.ZLBH != "" {
		data.Head = result

		// Lấy size
		query = "select XXCC,QTY from sczls where SCBH=?"
		var size []types.TraceOrderSize
		if err := db.Raw(query, orderNo).Scan(&size).Error; err != nil {
			return data, err
		}

		var sumQty float32
		for _, item := range size {
			sumQty += item.QTY
		}
		data.Size = size

		// Lấy chi tiết BOM
		query = `
		select XXZLS.BWBH,BWZL.ywsm BWMC,XXZLS.CLBH,CLZL.YWPM as CLMC,CLZL.DWBH,ZSZL.ZSYWJC,
		XXZLS.CLSL,XXZLS.CLSL * ? SumQty ,CLZL.CQDH
		FROM XXZLS
		LEFT join (select XieXing,SheHao,BWBH,Max(CLBH) as CLBH,Max(CLSL) as CLSL,Max(Log_DateTime) as Log_DateTime from XXZLS_Log
			where XXZLS_Log.XieXing=? and XXZLS_Log.SheHao=? and Log_DateTime>=GetDate()-1
			Group by XieXing,SheHao,BWBH) XXZLS_Log 
			on XXZLS_Log.XieXing=XXZLS.XieXing and XXZLS_Log.SheHao=XXZLS.SheHao and XXZLS_Log.BWBH=XXZLS.BWBH
		LEFT JOIN BWZL ON XXZLS.BWBH = BWZL.bwdh
		LEFT JOIN CLZL ON XXZLS.CLBH = CLZL.cldh
		LEFT JOIN LBZLS ON SUBSTRING(XXZLS.CLBH,1,1) = LBZLS.lbdh AND LBZLS.LB='05'
		LEFT JOIN ZSZL on ZSZL.ZSDH=XXZLS.CSBH
		where XXZLS.XieXing=? and XXZLS.SheHao=?
		ORDER BY XXZLS.XH`

		var detail []types.TraceOrderDetail
		if err := db.Raw(query, sumQty, result.XieXing, result.SheHao, result.XieXing, result.SheHao).Scan(&detail).Error; err != nil {
			return data, err
		}
		data.Body = detail
	}

	data.Head.IMGName = getimg(data.Head.IMGName)
	return data, nil
}

func GetMaterialPo(Ry string) ([]types.TraceCGNOInfo, []types.TraceOrderTmpAll, error) {
	var cgnoinfos []types.TraceCGNOInfo

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, nil, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	query := `
	select a.CGNO,a.GSBH,convert(varchar,a.CGDate,23) CGDate,a.DH,a.CZ ,a.ZSYWJC,CLBH,sum(Qty) Qty,
	convert(varchar,YQDate,23) YQDate, ywpm Ywpm,dwbh DWBH,
	(select count(*) from cghzzlss where clbh=CGZLS.CLBH) Order_ID
	from CGZLS
	join CGZL on CGZLS.CGNO=CGZL.CGNO
	join (
		select CGZL.CGNO,CGZL.GSBH,CGDate,ZSZL.DH,ZSZL.CZ ,ZSZL.ZSYWJC
		from CGZL
		left join ZSZL on CGZL.ZSBH=ZSZL.ZSDH
		left join ZSZL_DEV on ZSZL_DEV.zsdh=ZSZL.zsdh and ZSZL_DEV.GSBH='VA12'
		left join ZSZL ZSZLTW on ZSZLTW.zsdh=ZSZL_DEV.Zsdh_TW
		where ISNULL(flowflag,'')<>'X' and CGZL.CGNO in(
			select distinct CGNO from CGZLSS where ZLBH=?
		)
	) a on CGZLS.CGNO=a.CGNO
	join clzl on clzl.cldh=CLBH
	group by  a.CGNO,a.GSBH,a.CGDate,a.DH,a.CZ ,a.ZSYWJC,CLBH,YQDate, ywpm ,dwbh
	order by a.ZSYWJC,YQDate,a.CGNO
	`

	var result []types.TraceCGNOAll
	if err := db.Raw(query, Ry).Scan(&result).Error; err != nil {
		return nil, nil, err
	}

	// Gom nhóm theo ZSYWJC
	var cgnoinfo types.TraceCGNOInfo
	for _, item := range result {
		if cgnoinfo.Title.ZSYWJC != "" && cgnoinfo.Title.ZSYWJC != item.ZSYWJC {
			cgnoinfos = append(cgnoinfos, cgnoinfo)
			cgnoinfo = types.TraceCGNOInfo{}
		}
		cgnoinfo.Title.CGNO = item.CGNO
		cgnoinfo.Title.GSBH = item.GSBH
		cgnoinfo.Title.CGDate = item.CGDate
		cgnoinfo.Title.DH = item.DH
		cgnoinfo.Title.CZ = item.CZ
		cgnoinfo.Title.ZSYWJC = item.ZSYWJC
		cgnoinfo.Title.Order_ID += item.Order_ID
		cgnoinfo.Rows = append(cgnoinfo.Rows, types.TraceCGNODetail{
			CLBH: item.CLBH, Qty: item.Qty, YQDate: item.YQDate, Ywpm: item.Ywpm, DWBH: item.DWBH,
		})
	}
	if len(cgnoinfo.Rows) > 0 {
		cgnoinfos = append(cgnoinfos, cgnoinfo)
	}

	// Lấy Order details
	var Orders []types.TraceOrderTmpAll
	for _, item := range cgnoinfos {
		if item.Title.Order_ID > 0 {
			for _, row := range item.Rows {
				query = `
				SELECT cghzbh Cghzbh,dgdh Dgdh,clbh Clbh,ywpm Ywpm,dwbh Dwbh,zl_qty Zl_qty,zd_date Zd_date,zsbh Zsbh,dg_date Dg_date
				FROM cghzzlss a
				inner join clzl b on a.clbh=b.cldh
				inner join kfzl c on a.zsbh=c.kfdh
				WHERE clbh=? AND zl_qty=?
				ORDER BY cghzbh,dg_date,clbh
				`
				var Oresult []types.TraceOrderTmpInfo
				if err := db.Raw(query, row.CLBH, row.Qty).Scan(&Oresult).Error; err != nil {
					return cgnoinfos, nil, err
				}
				if len(Oresult) > 0 {
					var TmpDetails []types.TraceOrderTmpDetail
					var TmpHead types.TraceOrderTmpHead
					TmpHead.Cghzbh = Oresult[0].Cghzbh
					TmpHead.Kfqm = item.Title.ZSYWJC
					TmpHead.Dh = item.Title.DH
					TmpHead.Cz = item.Title.CZ
					TmpHead.Dg_date = Oresult[0].Dg_date
					for _, d := range Oresult {
						TmpDetails = append(TmpDetails, types.TraceOrderTmpDetail{
							Dgdh: d.Dgdh, Clbh: d.Clbh, Ywpm: d.Ywpm, Dwbh: d.Dwbh, Zl_qty: d.Zl_qty, Zd_date: d.Zd_date, Zsbh: d.Zsbh,
						})
					}
					Orders = append(Orders, types.TraceOrderTmpAll{Head: TmpHead, Detail: TmpDetails})
				}
			}
		}
	}

	return cgnoinfos, Orders, nil
}

func GetInvoice(Ry string) (types.TraceInvoice, error) {
	var data types.TraceInvoice

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return data, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return data, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	// Lấy thông tin đầu invoice
	query := `
	select top 1 INV_NO ,CONVERT(char(10),Inv_date,23) [Date],[TOTAL_PAIRS] [Goods],RISK Shipped,FROM_WHERE,TO_WHERE
	FROM [INVOICE_M]
	where INV_NO=(select top 1 INV_NO FROM [INVOICE_D] where RYNO=?)
	`
	var resultH types.TraceInvoiceHead
	if err := db.Raw(query, Ry).Scan(&resultH).Error; err != nil {
		return data, err
	}
	data.Head = resultH

	// Lấy thông tin chi tiết invoice
	query = `
	SELECT [RYNO],[Pairs],KHDDBH1,XieMing ,YSSM,ARTICLE
	FROM [PACKING_D]
	left join(
		select KHDDBH1,XieMing ,YSSM,DDZL.ARTICLE,DDZL.DDBH
		from YWDD
		left join DDZL on YWDD.YSBH=DDZL.DDBH
		left join XXZL on DDZL.XieXing=XXZL.XieXing and DDZL.SheHao=XXZL.SheHao
	) a on a.DDBH=[PACKING_D].RYNO
	where INV_NO=?
	`
	var resultB []types.TraceInvoiceTitle
	if err := db.Raw(query, data.Head.INV_NO).Scan(&resultB).Error; err != nil {
		return data, err
	}

	var tmps []types.TraceInvoice1
	for _, item := range resultB {
		var tcc types.TraceInvoice1
		tcc.Body = item

		// Lấy danh sách packing
		packQuery := `
		SELECT cast(CTQ as char(2))+'-'+cast(CTZ as char(20)) [Number], convert(decimal,SIZ) SIZ,QTY,CTS,PAIRS,
		NW,TNW,GW,TGW,L,W,H
		FROM [PACKING] where INV_NO=? and [RYNO]=?
		`
		if err := db.Raw(packQuery, data.Head.INV_NO, item.RYNO).Scan(&tcc.PackList).Error; err != nil {
			return data, err
		}

		// Lấy tổng packing
		totalQuery := `
		SELECT TOP 1 CTS,Pairs,NW,GW,CBM
		FROM [PACKING_D] where INV_NO=? and [RYNO]=?
		`
		if err := db.Raw(totalQuery, data.Head.INV_NO, item.RYNO).Scan(&tcc.Total).Error; err != nil {
			return data, err
		}

		tmps = append(tmps, tcc)
	}
	data.Body = tmps

	return data, nil
}


func GetCustomer(Ry string) (types.TraceCustomer, error) {
	var data types.TraceCustomer

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return data, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return data, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	// Lấy thông tin khách hàng
	var resultH types.TraceCustomerHead
	query := `
	select KHDDBH1,BUYER,ShipTo,KHDDBH2 Customer_PO,XieMing,XXZL.XieXing,XXZL.ARTICLE,YSSM,Qty,kfqm Cp
	from YWDD
	left join DDZL on YWDD.YSBH=DDZL.DDBH
    left join XXZL on DDZL.XieXing=XXZL.XieXing and DDZL.SheHao=XXZL.SheHao
	left join kfzl on DDZL.KHBH=kfdh
    where YWDD.DDBH=?
	`
	if err := db.Raw(query, Ry).Scan(&resultH).Error; err != nil {
		return data, err
	}
	data.Head = resultH

	// Lấy thông tin size
	var resultB []types.TraceCustomerTitle
	query = "select XXCC Size,Qty Size_Qty,'' Case_Qty,Unit Unit_Case from YWDDS where DDBH=?"
	if err := db.Raw(query, Ry).Scan(&resultB).Error; err != nil {
		return data, err
	}
	data.Body = resultB

	return data, nil
}

func GetProducts(Ry string) (types.TraceProduc, error) {
	var result types.TraceProduc

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return result, err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return result, fmt.Errorf("failed to get DB instance: %w", err)
	}
	defer dbInstance.Close()

	// Lấy danh sách sản phẩm
	query := `
	SELECT [XXCC],[Qty],CONVERT(char(10), [USERDate],120) INDATE,
	(SELECT isnull(sum([Qty]),0) FROM [SMDDS] where [DDBH]=?) TQty
	FROM [SMDDS] where [DDBH]=?
	`
	if err := db.Raw(query, Ry, Ry).Scan(&result.Items).Error; err != nil {
		return result, err
	}

	// Lấy sản phẩm đóng gói
	query = `
	select CONVERT(char(10), [INDATE],120) INDATE,CARTONBAR,Qty,sgw Sgw,rgw Rgw,
	isnull((select sum(Qty) FROM YWCP where [DDBH]=?),0) TQty
	FROM YWCP where [DDBH]=?
	`
	if err := db.Raw(query, Ry, Ry).Scan(&result.Items2).Error; err != nil {
		return result, err
	}

	return result, nil
}
