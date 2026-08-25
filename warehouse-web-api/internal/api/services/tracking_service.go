package services

import (
	"errors"
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"
	// "mes-web-api/internal/pkg/database"
)

type TrackingService struct {
	*BaseService
}

var TK = &TrackingService{}

func (s *TrackingService) ScanMaterialImport(
    GSBH, SCBH, CGNO, ZSDH, ZSYWJC, CLBH, YWBM,
    MAKE, Barcode,
    SCFMDATE, ECFMDATE string,
) ([]types.ScanMaterialInput, error) {

    if MAKE != "" && Barcode != "" {
        return nil, errors.New("Chỉ được nhập MAKE hoặc Barcode")
    }

    db, err := database.LYS_ERP_Connection()
    if err != nil {
        return nil, err
    }
    dbInstance, _ := db.DB()
    defer dbInstance.Close()

    query := `
    SELECT DISTINCT 
        KCRKScan_RFSS.SCNO,
        KCRKScan_RFS.CGNO,
        CGZL.ZSBH,
        ZSZL.zsywjc,
        KCRKScan_RF.GSBH,
        KCRKScan_RF.CKBH,
        KCRKScan_RFSS.CLBH,
        CLZL.YWPM,
        CLZL.DWBH,
        KCRKScan_RFSS.Qty,
        ISNULL(KCRKScan_RFSS.RemQty, KCRKScan_RFSS.Qty) AS RemQty,
        ROUND(
            KCRKScan_RFSS.Qty - ISNULL(KCRKScan_RFSS.RemQty, KCRKScan_RFSS.Qty),
            4
        ) AS DQty,
        KCRKScan_RFSS.TAGID,
        KCRKScan_RFSS.Pack,
        KCRKScan_RFSS.Memo_RY,
        KCRKScan_RFSS.Memo_Article,
        KCRKScan_RFSS.MAKE,
        KCRKScan_RFSS.barcode,
        KCRKScan_RF.MEMO,
        KCRKScan_RF.DOCNO,
        KCRKScan_RFS.IDCheck,
        KCRKScan_RFSS.CFMID,
        KCRKScan_RFSS.CFMDATE,
        MaterialQCCheck.File_Name_Lab as File_Name
    FROM KCRKScan_RFSS
    LEFT JOIN KCRKScan_RF 
        ON KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO
    LEFT JOIN KCRKScan_RFS 
        ON KCRKScan_RFS.SCNO = KCRKScan_RFSS.SCNO
        AND KCRKScan_RFS.CLBH = KCRKScan_RFSS.CLBH
    LEFT JOIN CGZL 
        ON KCRKScan_RFS.CGNO = CGZL.CGNO
    LEFT JOIN ZSZL 
        ON CGZL.ZSBH = ZSZL.ZSDH
    LEFT JOIN CLZL 
        ON KCRKScan_RFS.CLBH = CLZL.CLDH
    LEFT JOIN KCZLS 
        ON KCRKScan_RF.CKBH = KCZLS.CKBH
        AND KCRKScan_RFSS.CLBH = KCZLS.CLBH
    LEFT JOIN KCLLScan_RFSS 
        ON KCLLScan_RFSS.barcode = KCRKScan_RFSS.barcode
        AND KCLLScan_RFSS.CLBH = KCRKScan_RFSS.CLBH
    LEFT JOIN KCLLScan_RFSSS 
        ON KCLLScan_RFSSS.barcode = KCRKScan_RFSS.barcode
        AND KCLLScan_RFSSS.CLBH = KCRKScan_RFSS.CLBH
    LEFT JOIN MaterialQCCheck
        ON KCRKScan_RFSS.SID = MaterialQCCheck.No_ID
	WHERE 
        KCRKScan_RF.GSBH = ?
        AND CAST(KCRKScan_RFSS.CFMDATE AS DATE) BETWEEN ? AND ?
        AND (
            KCLLScan_RFSS.SCBH LIKE ?
            OR KCLLScan_RFSSS.SCBH LIKE ?
            OR KCLLScan_RFSS.SCBH IS NULL
            OR KCLLScan_RFSSS.SCBH IS NULL
        )
        AND KCRKScan_RFS.CGNO LIKE ?
        AND ZSZL.ZSDH LIKE ?
        AND ZSZL.zsywjc LIKE ?
        AND KCRKScan_RFSS.CLBH LIKE ?
        AND CLZL.YWPM LIKE ?
        AND (
            (? <> '' AND ? = '' AND KCRKScan_RFSS.MAKE = ?)
         OR (? <> '' AND ? = '' AND KCRKScan_RFSS.barcode = ?)
         OR (? = '' AND ? = '')
        )
    ORDER BY 
        KCRKScan_RFSS.SCNO DESC,
        KCRKScan_RFS.CGNO DESC,
        KCRKScan_RFSS.CLBH DESC,
        KCRKScan_RFSS.Pack ASC;
    `

    var result []types.ScanMaterialInput
    err = db.Raw(
        query,
        GSBH,
        SCFMDATE,
        ECFMDATE,
        "%"+SCBH+"%",
        "%"+SCBH+"%",
        "%"+CGNO+"%",
        "%"+ZSDH+"%",
        "%"+ZSYWJC+"%",
        "%"+CLBH+"%",
        "%"+YWBM+"%",
        MAKE, Barcode, MAKE,
        Barcode, MAKE, Barcode,
        MAKE, Barcode,
    ).Scan(&result).Error

    if err != nil {
        return nil, err
    }
    return result, nil
}

func (s *TrackingService) UpdateKCBHByBarcode(barcode, newKCBH string) error {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return fmt.Errorf("kết nối database thất bại: %w", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	if barcode == "" || newKCBH == "" {
		return errors.New("barcode, MAKE mới không được rỗng")
	}

	query := `
	UPDATE S
	SET S.MAKE = ?
	FROM KCRKScan_RFSS S
	INNER JOIN KCRKScan_RF RF ON S.SCNO = RF.SCNO
	WHERE S.barcode = ?
	`

	result := db.Exec(query, newKCBH, barcode)

	if result.Error != nil {
		return fmt.Errorf("lỗi cập nhật MAKE: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("không tìm thấy barcode hoặc dữ liệu không cho phép cập nhật")
	}

	return nil
}

func (s *TrackingService) ScanMaterialExport(GSBH, SCBH, CGNO, ZSDH, ZSYWJC, CLBH, YWBM, SCFMDATE, ECFMDATE string) ([]types.ScanMaterialInput1, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()

	query := `
	SELECT DISTINCT 
	KCRKScan_RFSS.SCNO,
	KCRKScan_RFS.CGNO,
	CGZL.ZSBH,
	ZSZL.zsywjc,
	KCRKScan_RF.GSBH,
	KCRKScan_RF.CKBH,
	KCRKScan_RFSS.CLBH,
	CLZL.YWPM,
	CLZL.DWBH,
	KCRKScan_RFSS.Qty,
	ISNULL(KCRKScan_RFSS.RemQty, KCRKScan_RFSS.Qty) AS RemQty,
	ROUND(KCRKScan_RFSS.Qty - ISNULL(KCRKScan_RFSS.RemQty, KCRKScan_RFSS.Qty), 4) AS DQty,
	KCRKScan_RFSS.TAGID,
	CONVERT(varchar, KCRKScan_RFSS.Pack, 101) + ' of ' + CONVERT(varchar, TotalPack, 101) AS PACK1,
	KCRKScan_RFSS.Pack,
	KCRKScan_RFSS.Memo_RY,
	KCRKScan_RFSS.Memo_Article,
	KCZLS.MAKE,
	KCRKScan_RFSS.barcode,
	KCRKScan_RF.MEMO,
	KCRKScan_RF.DOCNO,
	KCRKScan_RFS.IDCheck,
	KCRKScan_RFSS.CFMID,
	KCRKScan_RFSS.CFMDATE CFMDATE1,
	KCRKScan_RFSS.LotNO
FROM 
	KCRKScan_RFSS
LEFT JOIN 
	KCRKScan_RF ON KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO
LEFT JOIN 
	KCRKScan_RFS ON KCRKScan_RFS.SCNO = KCRKScan_RFSS.SCNO AND KCRKScan_RFS.CLBH = KCRKScan_RFSS.CLBH
LEFT JOIN 
	CGZL ON KCRKScan_RFS.CGNO = CGZL.CGNO
LEFT JOIN 
	ZSZL ON CGZL.ZSBH = ZSZL.ZSDH
	LEFT JOIN 
(
	SELECT SCNO, CLBH, COUNT(Pack) AS TotalPack
	FROM KCRKScan_RFSS
	GROUP BY SCNO, CLBH
) KCRKScan ON KCRKScan.SCNO = KCRKScan_RFS.SCNO AND KCRKScan.CLBH = KCRKScan_RFS.CLBH
LEFT JOIN 
	CLZL ON KCRKScan_RFS.CLBH = CLZL.CLDH
LEFT JOIN 
	KCZLS ON KCRKScan_RF.CKBH = KCZLS.CKBH AND KCRKScan_RFSS.CLBH = KCZLS.CLBH
LEFT JOIN 
	KCLLScan_RFSS ON KCLLScan_RFSS.barcode = KCRKScan_RFSS.barcode AND KCLLScan_RFSS.CLBH = KCRKScan_RFSS.CLBH
LEFT JOIN 
	KCLLScan_RFSSS ON KCLLScan_RFSSS.barcode = KCRKScan_RFSS.barcode AND KCLLScan_RFSSS.CLBH = KCRKScan_RFSS.CLBH
WHERE 
	KCRKScan_RF.GSBH = ?
	AND CAST(KCRKScan_RFSS.CFMDATE1 AS DATE) BETWEEN ? AND ?
	AND (KCLLScan_RFSS.SCBH LIKE ? OR KCLLScan_RFSSS.SCBH LIKE ? OR KCLLScan_RFSS.SCBH IS NULL OR KCLLScan_RFSSS.SCBH IS NULL)
	AND KCRKScan_RFS.CGNO LIKE ?
	AND ZSZL.ZSDH LIKE ?
	AND ZSZL.zsywjc LIKE ?
	AND KCRKScan_RFS.CLBH LIKE ?
	AND CLZL.YWPM LIKE ?
ORDER BY 
	KCRKScan_RFSS.SCNO DESC, 
	KCRKScan_RFS.CGNO DESC, 
	KCRKScan_RFSS.CLBH DESC, 
	KCRKScan_RFSS.Pack ASC;`

	var result []types.ScanMaterialInput1
	err = db.Raw(query, GSBH, SCFMDATE, ECFMDATE, "%"+SCBH+"%", "%"+SCBH+"%", "%"+CGNO+"%", "%"+ZSDH+"%", "%"+ZSYWJC+"%", "%"+CLBH+"%", "%"+YWBM+"%").Scan(&result).Error
	dbInstance.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *TrackingService) ScanMaterialPrint(
	GSBH, SCBH, CGNO, ZSDH, ZSYWJC, CLBH, YWBM,
	MAKE, Barcode, // Thêm 2 tham số này
	SCFMDATE, ECFMDATE string,
) ([]types.ScanMaterialPrint, error) {

	// Kiểm tra ràng buộc giống hàm Import
	if MAKE != "" && Barcode != "" {
		return nil, errors.New("Chỉ được nhập MAKE hoặc Barcode")
	}

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	query := `
    SELECT DISTINCT 
        KCRKScan_RFSS.SCNO,
        KCRKScan_RFSS.CLBH,
        CLZL.YWPM,
        ZSZL.ZSYWJC,
        KCRKScan_RFSS.Memo_RY,
        KCRKScan_RFSS.Memo_Article,
        KCRKScan_RFSS.Qty,
        CONVERT(varchar, KCRKScan_RFSS.Pack, 101) + ' of ' + CONVERT(varchar, TotalPack, 101) AS PACK1,
        KCRKScan_RFSS.Pack,
        KCRKScan_RF.Memo + ' - ' + KCRKScan_RF.DOCNO AS Memo,
        CONVERT(varchar, ISNULL(KCRKScan_RFSS.RemQty, KCRKScan_RFSS.Qty), 101) + ' ' + CLZL.DWBH AS DWBH,
        KCRKScan_RFS.CGNO,
        KCRKScan_RFSS.barcode AS QRCode,
        CONVERT(varchar, MONTH(GETDATE())) AS FIFO,
        KCRKScan_RFSS.LotNO,
        KCRKScan_RFS.IDCheck
    FROM KCRKScan_RFSS
    LEFT JOIN KCRKScan_RF ON KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO
    LEFT JOIN KCRKScan_RFS ON KCRKScan_RFS.SCNO = KCRKScan_RFSS.SCNO AND KCRKScan_RFS.CLBH = KCRKScan_RFSS.CLBH
    LEFT JOIN (
        SELECT SCNO, CLBH, COUNT(Pack) AS TotalPack
        FROM KCRKScan_RFSS
        GROUP BY SCNO, CLBH
    ) KCRKScan ON KCRKScan.SCNO = KCRKScan_RFS.SCNO AND KCRKScan.CLBH = KCRKScan_RFS.CLBH
    LEFT JOIN CGZL ON KCRKScan_RFS.CGNO = CGZL.CGNO
    LEFT JOIN ZSZL ON CGZL.ZSBH = ZSZL.ZSDH
    LEFT JOIN CLZL ON KCRKScan_RFSS.CLBH = CLZL.CLDH
    LEFT JOIN KCLLScan_RFSS ON KCLLScan_RFSS.barcode = KCRKScan_RFSS.barcode AND KCLLScan_RFSS.CLBH = KCRKScan_RFSS.CLBH
    LEFT JOIN KCLLScan_RFSSS ON KCLLScan_RFSSS.barcode = KCRKScan_RFSS.barcode AND KCLLScan_RFSSS.CLBH = KCRKScan_RFSS.CLBH
    WHERE 
        KCRKScan_RF.GSBH = ?
        AND CAST(KCRKScan_RFSS.CFMDATE AS DATE) BETWEEN ? AND ?
        AND (
            KCLLScan_RFSS.SCBH LIKE ? OR KCLLScan_RFSSS.SCBH LIKE ? 
            OR KCLLScan_RFSS.SCBH IS NULL OR KCLLScan_RFSSS.SCBH IS NULL
        )
        AND KCRKScan_RFS.CGNO LIKE ?
        AND ZSZL.ZSDH LIKE ?
        AND ZSZL.zsywjc LIKE ?
        AND KCRKScan_RFSS.CLBH LIKE ?
        AND CLZL.YWPM LIKE ?
        AND (
            (? <> '' AND ? = '' AND KCRKScan_RFSS.MAKE = ?)
            OR (? <> '' AND ? = '' AND KCRKScan_RFSS.barcode = ?)
            OR (? = '' AND ? = '')
        )
    ORDER BY 
        KCRKScan_RFSS.SCNO DESC, 
        KCRKScan_RFS.CGNO DESC, 
        KCRKScan_RFSS.CLBH DESC, 
        KCRKScan_RFSS.Pack ASC`

	var result []types.ScanMaterialPrint

	err = db.Raw(
		query,
		GSBH, SCFMDATE, ECFMDATE,
		"%"+SCBH+"%", "%"+SCBH+"%",
		"%"+CGNO+"%", "%"+ZSDH+"%", "%"+ZSYWJC+"%", "%"+CLBH+"%", "%"+YWBM+"%",
		MAKE, Barcode, MAKE, // Cho tham số lọc MAKE
		Barcode, MAKE, Barcode, // Cho tham số lọc Barcode
		MAKE, Barcode, // Cho tham số rỗng
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *TrackingService) GetKCLLScan_RFSS(Barcode string) ([]types.GetKCLLScan_RFSS, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()
	queryGet := (`
    SELECT LLNO,SCANQty,DFL ,SCBH ,SCANUSERDATE ,SCANUSERID,CLBH FROM KCLLScan_RFSS
    where KCLLScan_RFSS.barcode = ?
    and KCLLScan_RFSS.SCANQty<>0
	`)
	var result []types.GetKCLLScan_RFSS
	err = db.Raw(queryGet, Barcode).Scan(&result).Error
	dbInstance.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *TrackingService) GetKCLLScan_RFSSS(Barcode string) ([]types.GetKCLLScan_RFSSS, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()
	queryGet := (`
    SELECT LLNO,SCANQty,DFL ,SCBH ,SCANUSERDATE ,SCANUSERID,CLBH,LLNO_OWE FROM KCLLScan_RFSSS
    where KCLLScan_RFSSS.barcode = ?
    and KCLLScan_RFSSS.SCANQty<>0
	`)
	var result []types.GetKCLLScan_RFSSS
	err = db.Raw(queryGet, Barcode).Scan(&result).Error
	dbInstance.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *TrackingService) UpdateKCLLScan_RFSS(LLNO, SCBH, CLBH, Barcode, DFL, SCANUSERID, USERID string, ScanQty_input, SCANQty float64) error {

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return err
	}

	if ScanQty_input == 0 {
		// Thực hiện chèn vào BDelRec nếu SCANQty = 0
		queryInsert := `
				INSERT INTO BDelRec (TableName, TNO, CLBH, OldID, DelID, DelDate)
				VALUES (?, ?, ?, ?, ?, GETDATE())
			`
		tno := LLNO + "*" + SCBH
		clbh := CLBH + "*" + Barcode + "*" + fmt.Sprintf("%.2f", SCANQty)
		oldID := SCANUSERID
		delID := USERID
		err = db.Exec(queryInsert, "KCLLScan_RFSS", tno, clbh, oldID, delID, Barcode).Error
		if err != nil {
			return err
		}
		// Thực hiện thao tác xóa bản ghi
		err = db.Exec(`DELETE FROM KCLLScan_RFSS WHERE LLNO = ? AND CLBH = ? AND SCBH = ? AND DFL = ? AND barcode = ?`, LLNO, CLBH, SCBH, DFL, Barcode).Error
		if err != nil {
			return err
		}
	} else {
		// Cập nhật SCANUSERID và SCANUSERDATE
		queryUpdate := `
				UPDATE KCLLScan_RFSS 
				SET SCANQty = ?, SCANUSERID = ?, SCANUSERDATE = GETDATE()
				WHERE LLNO = ? AND CLBH = ? AND SCBH = ? AND DFL = ? AND barcode = ?
			`
		err = db.Exec(queryUpdate, ScanQty_input, USERID, LLNO, CLBH, SCBH, DFL, Barcode).Error
		if err != nil {
			return err
		}
	}

	fmt.Println("OK")
	return nil
}
func (s *TrackingService) UpdateKCRKScan_RFSS(Barcode string) error {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return err
	}
	// Cập nhật KCRKScan_RFSS
	queryUpdateKCRK := `
		UPDATE KCRKScan_RFSS 
		SET RemQty = ROUND(Qty - ISNULL((select SUM(CONVERT(float, SCANQty)) from KCLLScan_RFSSS where barcode=?),0) - ISNULL((select SUM(CONVERT(float, SCANQty)) from KCLLScan_RFSS where barcode=?),0), 2)
		WHERE barcode = ? AND CFMDATE IS NOT NULL
	`
	err = db.Exec(queryUpdateKCRK, Barcode, Barcode, Barcode).Error
	if err != nil {
		return err
	}

	fmt.Println("OK")
	return nil
}

func (s *TrackingService) UpdateKCLLScan_RFSSS(LLNO_OWE, LLNO, SCBH, CLBH, Barcode, SCANUSERID, USERID string, ScanQty_input, SCANQty float64) error {
	// Kết nối đến cơ sở dữ liệu

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return err
	}
	// qty := 0.0

	if ScanQty_input == 0 {
		// Thực hiện chèn vào BDelRec nếu SCANQty = 0
		queryInsert := `
				INSERT INTO BDelRec (TableName, TNO, CLBH, OldID, DelID, DelDate)
				VALUES (?, ?, ?, ?, ?, GETDATE())
			`
		tno := LLNO_OWE + "*" + LLNO + "*" + SCBH
		clbh := CLBH + "*" + Barcode + "*" + fmt.Sprintf("%.2f", SCANQty)
		oldID := SCANUSERID
		delID := USERID

		err = db.Exec(queryInsert, "KCLLScan_RFSSS", tno, clbh, oldID, delID, Barcode).Error
		if err != nil {
			return err
		}
		// Thực hiện thao tác xóa bản ghi

		err = db.Exec(` Delete From KCLLScan_RFSSS where LLNO_OWE = ?`, LLNO_OWE).Error
		if err != nil {
			return err
		}

		queryUpdateKCRK := `
		UPDATE KCRKScan_RFSS 
		SET RemQty = RemQty + ?
		WHERE barcode = ? AND CFMDATE IS NOT NULL
	`
		err = db.Exec(queryUpdateKCRK, SCANQty, Barcode).Error
		if err != nil {
			return err
		}
	}
	// else {
	// 	// Cập nhật SCANUSERID và SCANUSERDATE
	// 	queryUpdate := `
	//         Update
	//         KCLLScan_RFSSS
	//         Set
	//         SCANQty =?,
	//         SCANUSERID =?,
	//         SCANUSERDATE = GETDATE()
	//         where
	//         LLNO_OWE = ?
	// 		`
	// 	err = db.Exec(queryUpdate, ScanQty_input, USERID, LLNO_OWE).Error
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// Cập nhật KCRKScan_RFSS

	fmt.Println("OK 1")
	return nil
}
