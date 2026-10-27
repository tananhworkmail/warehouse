package services

import (
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"

	"strings"
	// "mes-web-api/internal/pkg/database"
)

type DiliveryPayService struct {
	*BaseService
}

var DP = &DiliveryPayService{}

func (s *DiliveryPayService) ScanDeliveryPayExPort(SCBH, CLBH, YWPM, SCFMDATE, ECFMDATE string) ([]types.DeliveryPayExPort, error) {
	db, err := database.TB_ERP_Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()

	query := `	SELECT KCLLScan_RFSSS.*, ywpm, dwbh, DDZL.ARTICLE + '/' + DDZL.BUYNO AS SKU
	FROM KCLLScan_RFSSS
	LEFT JOIN CLZL ON CLZL.cldh = KCLLScan_RFSSS.CLBH
	LEFT JOIN DDZL ON DDZL.DDBH = KCLLScan_RFSSS.SCBH
	WHERE KCLLScan_RFSSS.SCANQty <> 0 
		AND CAST(KCLLScan_RFSSS.SCANUSERDATE AS DATE) BETWEEN ? AND ?
		AND KCLLScan_RFSSS.SCBH LIKE  ?
		AND KCLLScan_RFSSS.CLBH LIKE ?
		AND EXISTS (SELECT 1 FROM CLZL WHERE CLZL.YWPM LIKE ? AND CLZL.cldh = KCLLScan_RFSSS.CLBH)
	ORDER BY KCLLScan_RFSSS.LLNO_OWE DESC, KCLLScan_RFSSS.CLBH DESC`
	var result []types.DeliveryPayExPort
	err = db.Raw(query, SCFMDATE, ECFMDATE, SCBH+"%", CLBH+"%", "%"+YWPM+"%").Scan(&result).Error
	dbInstance.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *DiliveryPayService) GetScanDeliveryPay(LLNOs []string) ([]types.ScanPayGetLLNO, error) {
	db, err := database.TB_ERP_Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()

	// Join LLNOs into a single string with single quotes around each value
	quotedLLNOs := "'" + strings.Join(LLNOs, "','") + "'"

	query := fmt.Sprintf(`
        SELECT 
        KCLLS.LLNO,
        KCLLS.CLBH,
        KCLLS.SCBH,
        KCLLS.DFL,
        CLZL.ywpm,
        CLZL.DWBH,
        KCLLS.Qty,
        SCBLYY.YWSM,
        DDZL.ARTICLE + '/' + DDZL.BUYNO AS SKU,
        KCLLS.BLSB,
        KCLLS.SCANUSERID,
        KCLLS.SCANUSERDATE,
        KCLLSAll.Allqty,
        SUM(ISNULL(KCLLScan_RFSS.SCANQty, 0)) AS LLQty,
        KCLLSAll.Allqty - ISNULL(KCLLScan_RFSSAll.LLAllqty, 0) AS TempQty,
        ISNULL(KCLLS.SCANQty, 0) AS checkry
    FROM 
        KCLLS 
    LEFT JOIN 
        (SELECT LLNO, CLBH, DFL, SCBH, SCANQty FROM KCLLScan_RFSS 
        UNION ALL 
        SELECT LLNO, CLBH, DFL, SCBH, SCANQty FROM KCLLScan_RFSSS) KCLLScan_RFSS 
    ON 
        KCLLScan_RFSS.LLNO = KCLLS.LLNO 
        AND KCLLScan_RFSS.CLBH = KCLLS.CLBH 
        AND KCLLScan_RFSS.DFL = KCLLS.DFL 
        AND KCLLScan_RFSS.SCBH = KCLLS.SCBH 
    LEFT JOIN 
        (SELECT CLBH, SUM(KCLLS.Qty) AS Allqty  
        FROM KCLLS 
        WHERE KCLLS.LLNO IN (%s) 
        GROUP BY CLBH) KCLLSAll ON KCLLSAll.CLBH = KCLLS.CLBH 
    LEFT JOIN 
        (SELECT CLBH, SUM(ISNULL(KCLLScan_RFSS.SCANQty, 0)) AS LLAllqty  
        FROM KCLLScan_RFSS 
        WHERE KCLLScan_RFSS.LLNO IN (%s) 
        GROUP BY CLBH) KCLLScan_RFSSAll ON KCLLScan_RFSSAll.CLBH = KCLLS.CLBH 
    LEFT JOIN 
        CLZL ON CLZL.cldh = KCLLS.CLBH 
    LEFT JOIN 
        DDZL ON DDZL.DDBH = KCLLS.SCBH 
    LEFT JOIN 
        SCBLYY ON SCBLYY.YYBH = KCLLS.YYBH 
    WHERE 
        KCLLS.LLNO IN (%s) 
    GROUP BY 
        KCLLS.LLNO,
        KCLLS.CLBH,
        KCLLS.SCBH,
        KCLLS.DFL,
        CLZL.ywpm,
        CLZL.DWBH,
        DDZL.ARTICLE,
        KCLLS.Qty,
        DDZL.BUYNO,
        KCLLS.SCANQty,
        SCBLYY.YWSM,
        KCLLS.BLSB,
        KCLLS.SCANUSERID,
        KCLLS.SCANUSERDATE,
        KCLLSAll.Allqty,
        KCLLScan_RFSSAll.LLAllqty 
    ORDER BY 
        KCLLS.CLBH`, quotedLLNOs, quotedLLNOs, quotedLLNOs)
	var result []types.ScanPayGetLLNO
	err = db.Raw(query).Scan(&result).Error
	dbInstance.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *DiliveryPayService) GetScanDeliveryPayKCLLScan_RFSS(LLNO, CLBH, DFL, SCBH string) ([]types.ScanPayGetKCLLScan_RFSS, error) {
	db, err := database.TB_ERP_Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()

	query := fmt.Sprintf(`
       select * from KCLLScan_RFSS
where LLNO='%s' and CLBH='%s' and DFL='%s' and SCBH='%s'
and KCLLScan_RFSS.SCANQty<>0`, LLNO, CLBH, DFL, SCBH)
	var result []types.ScanPayGetKCLLScan_RFSS
	err = db.Raw(query).Scan(&result).Error
	dbInstance.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *DiliveryPayService) DBGridEh2DblClick(LLNO_KCLLS, LLNO_OWE_LLSSS, CLBH_KCLL, SCBH_KCLL, DFL_KCLL, USERID string, LLQty_KCLL, Qty_KCLL float64) error {
	var qty, SCANQty1, SCANQty_LLSSS float64
	var CLBH_LLSSS, Barcode_LLSSS string
	// Kết nối đến cơ sở dữ liệu
	db, err := database.TB_ERP_Connection()
	if err != nil {
		return err
	}

	var result []types.ScanPayGetKCLLScan_RFSSS
	query := `SELECT CLBH, Barcode, SCANQty
			  FROM KCLLScan_RFSSS
			  WHERE KCLLScan_RFSSS.LLNO_OWE = ?`
	err = db.Raw(query, LLNO_OWE_LLSSS).Scan(&result).Error
	if err != nil {
		return err
	}

	CLBH_LLSSS = result[0].CLBH
	Barcode_LLSSS = result[0].Barcode
	SCANQty_LLSSS = result[0].SCANQty

	if LLNO_OWE_LLSSS != "" && CLBH_KCLL != "" {
		if CLBH_LLSSS == CLBH_KCLL && LLQty_KCLL < Qty_KCLL {
			SCANQty1 = Qty_KCLL - LLQty_KCLL
			qty = SCANQty_LLSSS - SCANQty1

			if SCANQty1 > 0 && SCANQty_LLSSS > 0 {
				updateKCLLS := `Update KCLLS set SCANUSERID = ?, SCANUSERDATE = GETDATE() where KCLLS.LLNO=? AND KCLLS.CLBH=? AND KCLLS.DFL=? AND KCLLS.SCBH=?`
				err = db.Exec(updateKCLLS, USERID, LLNO_KCLLS, CLBH_LLSSS, DFL_KCLL, SCBH_KCLL).Error
				if err != nil {
					return err
				}
				if qty > 0 {
					InsertLLSS := `INSERT INTO KCLLScan_RFSS (LLNO, CLBH, DFL, SCBH, barcode, SCANQty, SCANUSERID, SCANUSERDATE, LLNO_OWE) Values(?,?,?,?,?,?,?,GETDATE(), ?)`
					err = db.Exec(InsertLLSS, LLNO_KCLLS, CLBH_LLSSS, DFL_KCLL, SCBH_KCLL, Barcode_LLSSS, SCANQty1, USERID, LLNO_OWE_LLSSS).Error
					if err != nil {
						return err
					}

					updateLLSSS := `Update KCLLScan_RFSSS set SCANQty=? Where KCLLScan_RFSSS.LLNO_OWE=?`
					err = db.Exec(updateLLSSS, qty, LLNO_OWE_LLSSS).Error
					if err != nil {
						return err
					}
				} else {
					InsertLLSS := `INSERT INTO KCLLScan_RFSS (LLNO, CLBH, DFL, SCBH, barcode, SCANQty, SCANUSERID, SCANUSERDATE, LLNO_OWE) Values(?,?,?,?,?,?, ?,GETDATE(), ?)`
					err = db.Exec(InsertLLSS, LLNO_KCLLS, CLBH_LLSSS, DFL_KCLL, SCBH_KCLL, Barcode_LLSSS, SCANQty_LLSSS, USERID, LLNO_OWE_LLSSS).Error
					if err != nil {
						return err
					}

					updateLLSSS := `Update KCLLScan_RFSSS set SCANQty=0 Where KCLLScan_RFSSS.LLNO_OWE=?`
					err = db.Exec(updateLLSSS, LLNO_OWE_LLSSS).Error
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (s *DiliveryPayService) DeliveryPayKCLLScan_RFSSS(LLNO, LLNO_OWE, SCBH, CLBH, Barcode, DFL, SCANUSERID, USERID string, ScanQty_input, SCANQty float64) error {
	var qty float64
	// Kết nối đến cơ sở dữ liệu
	db, err := database.TB_ERP_Connection()
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
		updateKCLLSSS := `Update KCLLScan_RFSSS set SCANQty = SCANQty + ? where LLNO_OWE=?`
		err = db.Exec(updateKCLLSSS, SCANQty, LLNO_OWE).Error
		if err != nil {
			return err
		}
		// Thực hiện thao tác xóa bản ghi
		err = db.Exec(`DELETE FROM KCLLScan_RFSS WHERE LLNO = ? AND CLBH = ? AND SCBH = ? AND DFL = ? AND barcode = ?`, LLNO, CLBH, SCBH, DFL, Barcode).Error
		if err != nil {
			return err
		}
	} else {
		qty = SCANQty - ScanQty_input
		// Cập nhật KCRKScan_RFSS
		queryUpdateKCRK := `
		UPDATE KCLLScan_RFSSS 
		SET SCANQty = SCANQty + ?
		WHERE LLNO_OWE = ? 
	`
		err = db.Exec(queryUpdateKCRK, qty, LLNO_OWE).Error
		if err != nil {
			return err
		}
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
