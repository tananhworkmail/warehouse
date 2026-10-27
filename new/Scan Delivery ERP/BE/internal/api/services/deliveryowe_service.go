package services

import (
	"database/sql"
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"

	"strconv"
	// "mes-web-api/internal/pkg/database"
)

type DiliveryOweService struct {
	*BaseService
}

var DO = &DiliveryOweService{}

func (s *DiliveryOweService) ScanDeliveryOweQuantity(Quantity, Barcode, USERID, LLNO, CLBH, DFL, SCBH string) ([]types.DeliveryOweKCLLScan_RFSS, []types.DeliveryOweKCLLScan_RFSSS, error) {
	var remqty, qty float64
	var CLBH1 string

	db, err := database.TB_ERP_Connection()
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		dbInstance, _ := db.DB()
		dbInstance.Close()
	}()

	if Quantity == "" {
		Quantity = "0"
	}

	var result []types.ScanDeliveryOweLabel
	query := `SELECT ISNULL(KCRKScan_RFSS.RemQty, KCRKScan_RFSS.Qty) AS RemQty, KCRKScan_RFSS.CLBH, CFMDATE
			  FROM KCRKScan_RFSS
			  WHERE KCRKScan_RFSS.barcode = ?`
	err = db.Raw(query, Barcode).Scan(&result).Error
	if err != nil {
		return nil, nil, err
	}

	if len(result) == 0 || result[0].CFMDATE == "" {
		return nil, nil, fmt.Errorf("Vui lòng scan nhập kho !!!")
	}

	remqty = result[0].RemQty
	CLBH1 = result[0].CLBH

	if CLBH1 == CLBH {
		if val, err := strconv.ParseFloat(Quantity, 64); err == nil && val >= 0 && remqty > 0 {
			qty = remqty - val
			fmt.Printf("Computed qty: %.2f\n", qty)

			var y, m string
			query = "SELECT YEAR(GETDATE()) AS FY, MONTH(GETDATE()) AS FM"
			err = db.Raw(query).Row().Scan(&y, &m)
			if err != nil {
				return nil, nil, err
			}

			if len(m) < 2 {
				m = "0" + m
			}

			var a string
			query = "SELECT TOP 1 LLNO_OWE FROM KCLLScan_RFSSS WHERE LLNO_OWE LIKE ? ORDER BY LLNO_OWE DESC"
			err = db.Raw(query, y+m+"%").Row().Scan(&a)
			if err != nil && err != sql.ErrNoRows {
				return nil, nil, err
			}

			if err == sql.ErrNoRows {
				a = y + m + "00001"
			} else {
				num, _ := strconv.Atoi(a[6:])
				a = y + m + fmt.Sprintf("%05d", num+1)
			}

			if val == 0 {
				query = `INSERT INTO KCLLScan_RFSSS (LLNO_OWE, LLNO, CLBH, DFL, SCBH, barcode, SCANQty, SCANUSERID, SCANUSERDATE)
						 VALUES (?, ?, ?,?, ?, ?, ?, ?, GETDATE())`
				err = db.Exec(query, a, LLNO, CLBH, DFL, SCBH, Barcode, remqty, USERID).Error
				if err != nil {
					return nil, nil, err
				}

				query = `UPDATE KCRKScan_RFSS
						 SET RemQty=0, CFMID1=?, CFMDATE1=GETDATE()
						 WHERE barcode=? AND CFMDATE IS NOT NULL`
				err = db.Exec(query, USERID, Barcode).Error
				if err != nil {
					return nil, nil, err
				}
			} else if val > 0 {
				if qty >= 0 {
					query = `INSERT INTO KCLLScan_RFSSS (LLNO_OWE, LLNO, CLBH, DFL, SCBH, barcode, SCANQty, SCANUSERID, SCANUSERDATE)
							 VALUES (?, ?, ?, ?, ?, ?,?,?, GETDATE())`
					err = db.Exec(query, a, LLNO, CLBH, DFL, SCBH, Barcode, val, USERID).Error
					if err != nil {
						return nil, nil, err
					}

					query = `UPDATE KCRKScan_RFSS
							 SET RemQty=?, CFMID1=?, CFMDATE1=GETDATE()
							 WHERE barcode=? AND CFMDATE IS NOT NULL`
					err = db.Exec(query, qty, USERID, Barcode).Error
					if err != nil {
						return nil, nil, err
					}
				}
			}
		}
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

	var resultKCRKScan_RFSS []types.DeliveryOweKCLLScan_RFSS
	err = db.Raw(query, Barcode).Scan(&resultKCRKScan_RFSS).Error
	if err != nil {
		return nil, nil, err
	}

	query = `	select KCLLScan_RFSSS.*, ywpm, dwbh,DDZL.ARTICLE+'/'+DDZL.BUYNO as SKU 
	from KCLLScan_RFSSS
	left join CLZL on CLZL.cldh=KCLLScan_RFSSS.CLBH 
	left join DDZL on DDZL.DDBH=KCLLScan_RFSSS.SCBH
	where LLNO=? and CLBH= ? and DFL= ? and SCBH= ?
	and KCLLScan_RFSSS.SCANQty<>0`

	var resultKCLLScan_RFSSS []types.DeliveryOweKCLLScan_RFSSS
	err = db.Raw(query, LLNO, CLBH, DFL, SCBH).Scan(&resultKCLLScan_RFSSS).Error
	if err != nil {
		return nil, nil, err
	}

	return resultKCRKScan_RFSS, resultKCLLScan_RFSSS, nil
}

func (s *DiliveryOweService) DeliveryOweKCLLScan_RFSS(LLNO_OWE, LLNO, SCBH, CLBH, Barcode, SCANUSERID, USERID string, ScanQty_input, SCANQty float64) error {
	// Kết nối đến cơ sở dữ liệu
	db, err := database.TB_ERP_Connection()
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

		// Cập nhật KCRKScan_RFSS
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

	fmt.Println("OK 1")
	return nil
}
func (s *DiliveryOweService) ScanDeliveryOwe(Barcode, USERID, LLNO, CLBH, DFL, SCBH string) ([]types.DeliveryOweKCLLScan_RFSSS, []types.ScanMaterialPrint, error) {
	// 建立資料庫連線
	db, err := database.TB_ERP_Connection()
	if err != nil {
		return nil, nil, err
	}

	// 第一個查詢：掃描交貨欠件
	query1 := `SELECT 
				KCLLScan_RFSSS.*, 
				CLZL.YWPM, 
				DDZL.ARTICLE+'/'+DDZL.BUYNO AS SKU 
			   FROM 
				KCLLScan_RFSSS
			   LEFT JOIN 
				CLZL ON CLZL.cldh=KCLLScan_RFSSS.CLBH 
			   LEFT JOIN 
				DDZL ON DDZL.DDBH=KCLLScan_RFSSS.SCBH
			   WHERE 
				LLNO=? AND CLBH=? AND DFL=? AND SCBH=? AND KCLLScan_RFSSS.SCANQty<>0`

	var resultKCLLScan_RFSSS []types.DeliveryOweKCLLScan_RFSSS
	err = db.Raw(query1, LLNO, CLBH, DFL, SCBH).Scan(&resultKCLLScan_RFSSS).Error
	if err != nil {
		return nil, nil, err
	}

	// 第二個查詢：掃描原料列印
	query2 := `SELECT 
				KCRKScan_RFSS.SCNO, 
				KCRKScan_RFSS.CLBH, 
				CLZL.YWPM, 
				ZSZL.ZSYWJC, 
				KCRKScan_RFSS.Memo_RY, 
				KCRKScan_RFSS.Memo_Article, 
				KCRKScan_RFSS.Qty, 
				CONVERT(varchar, KCRKScan_RFSS.Pack, 101) + ' of ' + CONVERT(varchar, TotalPack, 101) AS PACK1, 
				KCRKScan_RFSS.Pack, 
				KCRKScan_RF.Memo + ' - ' + ISNULL(KCRKScan_RF.DOCNO,'') AS Memo, 
				CONVERT(varchar, ISNULL(KCRKScan_RFSS.RemQty, KCRKScan_RFSS.Qty), 101) + ' ' + CLZL.DWBH AS DWBH, 
				KCRKScan_RFS.CGNO, 
				KCRKScan_RFSS.barcode AS QRCode, 
				CONVERT(varchar, MONTH(GETDATE())) AS FIFO, 
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
					SELECT SCNO, CLBH, COUNT(Pack) AS TotalPack   
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
	err = db.Raw(query2, Barcode).Scan(&resultMaterialPrint).Error
	if err != nil {
		return nil, nil, err
	}

	return resultKCLLScan_RFSSS, resultMaterialPrint, nil
}
