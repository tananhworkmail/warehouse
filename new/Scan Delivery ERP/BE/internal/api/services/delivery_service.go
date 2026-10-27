package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"strconv"
	"strings"
	// "mes-web-api/internal/pkg/database"
)

type DeliveryService struct {
	*BaseService
}

var DL = &DeliveryService{}

func (s *DeliveryService) GetScanDelivery(LLNOs []string) ([]types.ScanGetLLNO, error) {
	db, err := database.TB_ERP_Connection()
	if err != nil {
		return nil, err
	}
	defer func() {
		dbInstance, _ := db.DB()
		dbInstance.Close()
	}()

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

	// fmt.Println("Executing query:", query)

	var result []types.ScanGetLLNO
	err = db.Raw(query).Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	return result, nil
}
func (s *DeliveryService) GetScanDeliveryQuantity(
	Quantity string,
	Barcode string,
	KCBH string,
	USERID string,
) ([]types.ScanMaterialPrint, error) {

	db, err := database.TB_ERP_Connection()
	if err != nil {
		return nil, err
	}
	defer func() {
		dbInstance, _ := db.DB()
		dbInstance.Close()
	}()

	if Quantity == "" {
		Quantity = "0"
	}

	type ScanDeliveryLabel struct {
		RemQty  float64
		CLBH    string
		KCBH    string
	}

	var scanData ScanDeliveryLabel

	query := `
		SELECT
			ISNULL(RemQty, Qty) AS RemQty,
			CLBH,
			KCBH
		FROM KCRKScan_RFSS
		WHERE barcode = ?
			AND CFMDATE IS NOT NULL
	`

	err = db.Raw(query, Barcode).Scan(&scanData).Error
	if err != nil {
		return nil, err
	}

	if scanData.CLBH == "" {
		return nil, errors.New("Barcode không tồn tại hoặc chưa xác nhận")
	}

	if scanData.KCBH != KCBH {
		return nil, errors.New("Sai kệ: barcode không thuộc kệ đã scan")
	}

	val, err := strconv.ParseFloat(Quantity, 64)
	if err != nil || val < 0 {
		return nil, errors.New("Số lượng không hợp lệ")
	}

	if val > scanData.RemQty {
		return nil, errors.New("Số lượng nhập lớn hơn số lượng còn lại")
	}

	if scanData.RemQty > 0 {

		tx := db.Begin()
		if tx.Error != nil {
			return nil, tx.Error
		}

		now := time.Now()
		prefix := now.Format("200601")

		var lastLLNO sql.NullString
		var LLNO string

		err = tx.Raw(`
			SELECT MAX(LLNO_OWE)
			FROM KCLLScan_RFSSS
			WHERE LLNO_OWE LIKE ?
		`, prefix+"%").Row().Scan(&lastLLNO)

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if lastLLNO.Valid && lastLLNO.String != "" {
			num, _ := strconv.Atoi(lastLLNO.String[6:])
			LLNO = prefix + fmt.Sprintf("%05d", num+1)
		} else {
			LLNO = prefix + "00001"
		}

		scanQty := val
		remainQty := scanData.RemQty - val

		if val == 0 {
			scanQty = scanData.RemQty
			remainQty = 0
		}

		insertQuery := `
			INSERT INTO KCLLScan_RFSSS
			(
				LLNO_OWE,
				LLNO,
				CLBH,
				DFL,
				SCBH,
				barcode,
				SCANQty,
				SCANUSERID,
				SCANUSERDATE,
				KCBH
			)
			VALUES
			(
				?,
				'NO',
				?,
				CONVERT(VARCHAR(15), GETDATE(), 114),
				CONVERT(VARCHAR(15), GETDATE(), 111),
				?,
				?,
				?,
				GETDATE(),
				?
			)
		`

		err = tx.Exec(
			insertQuery,
			LLNO,
			scanData.CLBH,
			Barcode,
			scanQty,
			USERID,
			KCBH,
		).Error

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		updateQuery := `
			UPDATE KCRKScan_RFSS
			SET
				RemQty = ?,
				CFMID1 = ?,
				CFMDATE1 = GETDATE()
			WHERE barcode = ?
				AND CFMDATE IS NOT NULL
		`

		err = tx.Exec(
			updateQuery,
			remainQty,
			USERID,
			Barcode,
		).Error

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if err = tx.Commit().Error; err != nil {
			return nil, err
		}
	}

	var resultMaterialPrint []types.ScanMaterialPrint

	query = `
		SELECT
			KCRKScan_RFSS.SCNO,
			KCRKScan_RFSS.CLBH,
			CLZL.YWPM,
			ZSZL.ZSYWJC,
			KCRKScan_RFSS.Memo_RY,
			KCRKScan_RFSS.Memo_Article,
			KCRKScan_RFSS.Qty,
			CONVERT(varchar, KCRKScan_RFSS.Pack, 101) + ' of ' +
			CONVERT(varchar, TotalPack, 101) AS PACK1,
			KCRKScan_RFSS.Pack,
			ROUND(
				KCRKScan_RFSS.Qty -
				ISNULL(KCRKScan_RFSS.RemQty, KCRKScan_RFSS.Qty),
				4
			) AS DQty,
			KCRKScan_RF.Memo + ' - ' +
			ISNULL(KCRKScan_RF.DOCNO,'') AS Memo,
			CONVERT(
				varchar,
				ISNULL(KCRKScan_RFSS.RemQty, KCRKScan_RFSS.Qty),
				101
			) + ' ' + CLZL.DWBH AS DWBH,
			KCRKScan_RFS.CGNO,
			KCRKScan_RFSS.barcode AS QRCode,
			CONVERT(varchar, MONTH(GETDATE())) AS FIFO,
			KCRKScan_RFSS.LotNO,
			KCRKScan_RFS.IDCheck
		FROM KCRKScan_RFSS
		LEFT JOIN KCRKScan_RF
			ON KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO
		LEFT JOIN KCRKScan_RFS
			ON KCRKScan_RFSS.SCNO = KCRKScan_RFS.SCNO
			AND KCRKScan_RFSS.CLBH = KCRKScan_RFS.CLBH
		LEFT JOIN (
			SELECT
				SCNO,
				CLBH,
				COUNT(Pack) AS TotalPack
			FROM KCRKScan_RFSS
			GROUP BY SCNO, CLBH
		) KCRKScan
			ON KCRKScan.SCNO = KCRKScan_RFS.SCNO
			AND KCRKScan.CLBH = KCRKScan_RFS.CLBH
		LEFT JOIN CGZL
			ON KCRKScan_RFS.CGNO = CGZL.CGNO
		LEFT JOIN ZSZL
			ON CGZL.ZSBH = ZSZL.ZSDH
		LEFT JOIN CLZL
			ON KCRKScan_RFSS.CLBH = CLZL.CLDH
		WHERE KCRKScan_RFSS.barcode = ?
	`

	err = db.Raw(query, Barcode).Scan(&resultMaterialPrint).Error
	if err != nil {
		return nil, err
	}

	return resultMaterialPrint, nil
}

func (s *DeliveryService) GetScanDeliveryCLBH(Barcode, USERID string, Data []request.ScanDeliveryLabelCLBHSelected) ([]types.ScanMaterialPrint, error) {
	var remqty, qty float64
	var CLBH1 string

	db, err := database.TB_ERP_Connection()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return nil, err
	}
	defer func() {
		dbInstance, _ := db.DB()
		dbInstance.Close()
	}()

	var result []types.ScanDeliveryLabel
	query := `SELECT ISNULL(KCRKScan_RFSS.RemQty, KCRKScan_RFSS.Qty) AS RemQty, KCRKScan_RFSS.CLBH, CFMDATE
              FROM KCRKScan_RFSS
              WHERE KCRKScan_RFSS.barcode = ?`
	err = db.Raw(query, Barcode).Scan(&result).Error
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}

	if len(result) == 0 || result[0].CFMDATE == "" || result[0].RemQty == 0 {
		log.Println("No data found for the given barcode.")
		return []types.ScanMaterialPrint{}, nil
	}

	remqty = result[0].RemQty
	CLBH1 = result[0].CLBH
	for i := 0; i < len(Data); i++ {
		if CLBH1 == Data[i].CLBH {
			query = `SELECT (MAX(KCLLS.Qty) - SUM(ISNULL(KCLLScan_RFSS.SCANQty, 0))) AS SCANQty 
                 FROM KCLLS 
                 LEFT JOIN KCLLScan_RFSS ON KCLLScan_RFSS.LLNO = KCLLS.LLNO AND KCLLScan_RFSS.CLBH = KCLLS.CLBH 
                 AND KCLLScan_RFSS.DFL = KCLLS.DFL AND KCLLScan_RFSS.SCBH = KCLLS.SCBH 
                 WHERE KCLLS.LLNO = ? AND KCLLS.CLBH = ? AND KCLLS.DFL = ? AND KCLLS.SCBH = ?`
			var resultScanQty []types.Scanqty
			err = db.Raw(query, Data[i].LLNO, CLBH1, Data[i].DFL, Data[i].SCBH).Scan(&resultScanQty).Error
			if err != nil {
				log.Printf("Query error (resultScanQty): %v", err)
				return nil, err
			}
			SCANQty := resultScanQty[0].SCANQty
			qty = remqty - SCANQty

			log.Printf("SCANQty: %f, remqty: %f, qty: %f", SCANQty, remqty, qty)

			if SCANQty > 0 && remqty > 0 {
				log.Println("SCANQty and remqty are greater than 0")
				query = `UPDATE KCLLS SET SCANUSERID = ?, SCANUSERDATE = GETDATE() 
                     WHERE KCLLS.LLNO = ? AND KCLLS.CLBH = ? AND KCLLS.DFL = ? AND KCLLS.SCBH = ?`
				err = db.Exec(query, USERID, Data[i].LLNO, CLBH1, Data[i].DFL, Data[i].SCBH).Error
				if err != nil {
					log.Printf("Update KCLLS error: %v", err)
					return nil, err
				}
				if qty > 0 {
					log.Println("Qty is greater than 0")
					query = `INSERT INTO KCLLScan_RFSS (LLNO, CLBH, DFL, SCBH, barcode, SCANQty, SCANUSERID, SCANUSERDATE) 
                         VALUES (?, ?, ?, ?, ?, ?, ?, GETDATE())`
					err = db.Exec(query, Data[i].LLNO, CLBH1, Data[i].DFL, Data[i].SCBH, Barcode, SCANQty, USERID).Error
					if err != nil {
						log.Printf("Insert KCLLScan_RFSS error: %v", err)
						return nil, err
					}
					remqty = qty
				} else {
					log.Println("Qty is not greater than 0, inserting remaining qty")
					query = `INSERT INTO KCLLScan_RFSS (LLNO, CLBH, DFL, SCBH, barcode, SCANQty, SCANUSERID, SCANUSERDATE) 
                         VALUES (?, ?, ?, ?, ?, ?, ?, GETDATE())`
					err = db.Exec(query, Data[i].LLNO, CLBH1, Data[i].DFL, Data[i].SCBH, Barcode, remqty, USERID).Error
					if err != nil {
						log.Printf("Insert KCLLScan_RFSS error: %v", err)
						return nil, err
					}
					remqty = 0
					query = `UPDATE KCRKScan_RFSS SET RemQty = 0, CFMID1 = ?, CFMDATE1 = GETDATE() 
                         WHERE barcode = ? AND CFMDATE IS NOT NULL`
					err = db.Exec(query, USERID, Barcode).Error
					if err != nil {
						log.Printf("Update KCRKScan_RFSS error: %v", err)
						return nil, err
					}
				}
			}
		} else {
			log.Printf("CLBH1 does not match CLBH: %s != %s", CLBH1, Data[i].CLBH)
		}
	}

	query = `UPDATE KCRKScan_RFSS SET RemQty = ?, CFMID1 = ?, CFMDATE1 = GETDATE() 
		WHERE barcode = ? AND CFMDATE IS NOT NULL`
	err = db.Exec(query, remqty, USERID, Barcode).Error
	if err != nil {
		log.Printf("Update KCRKScan_RFSS error: %v", err)
		return nil, err
	}

	query = `SELECT 
				KCRKScan_RFSS.SCNO,
				KCRKScan_RFSS.CLBH,
				CLZL.YWPM,
				ZSZL.ZSYWJC,
				KCRKScan_RFSS.Memo_RY,
				KCRKScan_RFSS.Memo_Article,
				KCRKScan_RFSS.Qty,
				CONVERT(varchar, KCRKScan_RFSS.Pack, 101) + ' of ' + CONVERT(varchar, TotalPack, 101) as PACK1,
				KCRKScan_RFSS.Pack,
				ROUND(KCRKScan_RFSS.Qty - ISNULL(KCRKScan_RFSS.RemQty, KCRKScan_RFSS.Qty), 4) AS DQty,
				KCRKScan_RF.Memo + ' - ' + isnull(KCRKScan_RF.DOCNO,'') as Memo,
				CONVERT(varchar, ISNULL(KCRKScan_RFSS.RemQty, KCRKScan_RFSS.Qty), 101) + ' ' + CLZL.DWBH as DWBH,
				KCRKScan_RFS.CGNO,
				KCRKScan_RFSS.barcode as QRCode,
				CONVERT(varchar, MONTH(GETDATE())) as FIFO,
				KCRKScan_RFSS.LotNO,
				KCRKScan_RFS.IDCheck 
			  FROM 
				KCRKScan_RFSS    
			  LEFT JOIN 
				KCRKScan_RF ON KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO 
			  LEFT JOIN 
				KCRKScan_RFS ON KCRKScan_RFSS.SCNO = KCRKScan_RFS.SCNO AND KCRKScan_RFSS.CLBH = KCRKScan_RFS.CLBH 
			  LEFT JOIN 
				(
					SELECT SCNO, CLBH, COUNT(Pack) as TotalPack   
					FROM KCRKScan_RFSS  
					GROUP BY SCNO, CLBH
				) KCRKScan ON KCRKScan.SCNO = KCRKScan_RFS.SCNO AND KCRKScan.CLBH = KCRKScan_RFS.CLBH 
			  LEFT JOIN 
				CGZL ON KCRKScan_RFS.CGNO = CGZL.CGNO   
			  LEFT JOIN 
				ZSZL ON CGZL.ZSBH = ZSZL.ZSDH  
			  LEFT JOIN 
				CLZL ON KCRKScan_RFSS.CLBH = CLZL.CLDH  
			  WHERE 
				KCRKScan_RFSS.barcode = ?`

	var resultMaterialPrint []types.ScanMaterialPrint
	err = db.Raw(query, Barcode).Scan(&resultMaterialPrint).Error
	if err != nil {
		log.Printf("Query error (resultMaterialPrint): %v", err)
		return nil, err
	}

	return resultMaterialPrint, nil
}
