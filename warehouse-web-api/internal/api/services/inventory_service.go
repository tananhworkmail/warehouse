package services

import (
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
	"gorm.io/gorm"
	// "strings"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"
)

type DiscrepanciesService struct {
	*BaseService
}

var Discrepancies = &DiscrepanciesService{}

type CheckResult struct {
	Cfmid string `gorm:"column:Cfmid"`
	CLBH  string `gorm:"column:CLBH"`
}

// HÀM POST: Thêm mới Inventory
func (s *DiscrepanciesService) CreateInventory(param *types.Inventory) (string, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return "Connection Database Fail", err
	}
	db = db.Debug()
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Gán giá trị mặc định cho YN nếu phía FE gửi rỗng
	if param.YN == "" {
		param.YN = "1"
	}

	// 1. Kiểm tra Barcode trong kho (KCRKScan_RFSS)
	queryCheck := `SELECT cfmid AS Cfmid, CLBH FROM KCRKScan_RFSS WHERE barcode = ?`
	var result CheckResult
	err = db.Raw(queryCheck, param.Barcode).Scan(&result).Error
	if result.Cfmid == "" {
		return "", errors.New("Chưa Nhập Kho")
	}

	// 2. Kiểm tra trùng lặp (Barcode + Date + GSBH)
	queryCheck2 := `
        SELECT USERID FROM KCRKScan_RFKK 
        WHERE barcode = ? 
        AND CONVERT(varchar(10), USERDATE, 23) = ?
        AND GSBH = ?
    `
	var existingUserID string
	db.Raw(queryCheck2, param.Barcode, param.USERDATE, param.GSBH).Scan(&existingUserID)
	if existingUserID != "" {
		return "", errors.New("Đã được " + existingUserID + " thêm vào cho " + param.GSBH)
	}

	// 3. Thực hiện Insert với đầy đủ các cột
	queryInsert := `
        INSERT INTO KCRKScan_RFKK (barcode, USERDATE, USERID, GSBH, YN)
        VALUES (?, ?, ?, ?, ?)
    `
	err = db.Exec(queryInsert,
		param.Barcode,
		param.USERDATE,
		param.USERID,
		param.GSBH,
		param.YN, // Truyền giá trị "1" từ struct vào đây
	).Error

	if err != nil {
		return "Insert Fail: " + err.Error(), err
	}

	return result.CLBH, nil
}

// kiểm kê theo tháng
func (s *DiscrepanciesService) CreateInventoryByMonth(param *types.Inventory) (string, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return "Connection Database Fail", err
	}
	db = db.Debug()
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Default YN
	if param.YN == "" {
		param.YN = "1"
	}

	// =========================
	// 1. Kiểm tra tồn kho — barcode có trong kho chưa
	// =========================
	queryCheck := `SELECT cfmid AS Cfmid, CLBH FROM KCRKScan_RFSS WHERE barcode = ?`
	var result CheckResult
	err = db.Raw(queryCheck, param.Barcode).Scan(&result).Error
	if err != nil {
		return "Query Check Fail", err
	}
	if result.Cfmid == "" {
		return "", errors.New("Chưa Nhập Kho")
	}

	// =========================
	// 2. Kiểm tra trùng THÁNG
	//    param.USERDATE đã là 'yyyy-MM' (VD: "2025-07")
	//    So sánh FORMAT(USERDATE, 'yyyy-MM') với param.USERDATE
	// =========================
	queryCheck2 := `
		SELECT USERID FROM KCRKScan_RFKK 
		WHERE barcode = ? 
		AND FORMAT(USERDATE, 'yyyy-MM') = ?
		AND GSBH = ?
	`
	var existingUserID string
	db.Raw(queryCheck2, param.Barcode, param.USERDATE, param.GSBH).Scan(&existingUserID)
	if existingUserID != "" {
		return "", errors.New("Đã được " + existingUserID + " thêm vào tháng này cho " + param.GSBH)
	}

	// =========================
	// 3. MERGE bảng tháng
	//    param.USERDATE đã là 'yyyy-MM' — dùng trực tiếp, KHÔNG FORMAT lại
	//    Round1UserID/Round1Date = người + ngày quét xác nhận lần 1 (đơn vị kiểm kê)
	// =========================
	queryMonth := `
		MERGE KCRKScan_RFKK_MONTH AS target
		USING (
			SELECT ? AS YearMonth, ? AS Barcode, ? AS GSBH, ? AS UserID
		) AS source
		ON  target.YearMonth = source.YearMonth
		AND target.Barcode   = source.Barcode
		AND target.GSBH      = source.GSBH

		WHEN MATCHED THEN
			UPDATE SET
				TotalScan    = target.TotalScan + 1,
				Round1Date   = GETDATE(),
				Round1UserID = source.UserID

		WHEN NOT MATCHED THEN
			INSERT (YearMonth, Barcode, GSBH, TotalScan, Round1Date, Round1UserID)
			VALUES (source.YearMonth, source.Barcode, source.GSBH, 1, GETDATE(), source.UserID);
	`
	err = db.Exec(queryMonth,
		param.USERDATE, // 'yyyy-MM' — trực tiếp
		param.Barcode,
		param.GSBH,
		param.USERID,
	).Error
	if err != nil {
		return "Update Month Fail: " + err.Error(), err
	}

	return result.CLBH, nil
}

type RawInventory struct {
	GSBH         string
	Rack         string
	CLBH         string
	Pack         string
	Qty          float64
	RemQty       *float64
	Memo_Article string
	Memo_RY      string
}

func formatRound1(val float64) string {
	rounded := math.Round(val*10) / 10
	return strconv.FormatFloat(rounded, 'f', 1, 64)
}

// truncateStr cắt bớt chuỗi theo số ký tự tối đa để tránh lỗi
// "String or binary data would be truncated" khi insert vào SQL Server.
func truncateStr(s string, max int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= max {
		return string(r)
	}
	return string(r[:max])
}

func (s *DiscrepanciesService) SelectInventory(Date string, GSBH string) ([]types.ListInventory, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	if GSBH == "" {
		return []types.ListInventory{}, nil
	}

	// =========================
	// 1. CLBH đã scan trong ngày
	// =========================
	var scannedCLBHs []string
	err = db.Raw(`
		SELECT DISTINCT s2.CLBH 
		FROM KCRKScan_RFKK k2
		INNER JOIN KCRKScan_RFSS s2 ON k2.barcode = s2.barcode
		WHERE CONVERT(varchar(10), k2.USERDATE, 23) = ? AND k2.GSBH = ?
	`, Date, GSBH).Scan(&scannedCLBHs).Error

	if err != nil || len(scannedCLBHs) == 0 {
		return []types.ListInventory{}, err
	}

	// =========================
	// 2. Lấy tồn kho (KHÔNG xử lý BuyNos ở SQL)
	// =========================
	type RFSSRaw struct {
		Barcode      string
		MAKE         string
		CLBH         string
		Pack         string
		Memo_Article string
		Memo_RY      string
		RemQty       float64
		CFMDate      time.Time
		BuyNos       string
	}

	var allSS []RFSSRaw

	err = db.Raw(`
		SELECT 
			a.barcode, 
			a.MAKE, 
			a.CLBH, 
			a.Pack, 
			a.Memo_Article, 
			a.Memo_RY, 
			a.RemQty,
			a.CFMDate
		FROM KCRKScan_RFSS a
		LEFT JOIN KCRKScan_RF b ON a.SCNO = b.SCNO 
		WHERE a.CLBH IN (?) 
			AND a.CFMDate IS NOT NULL 
			AND ISNULL(a.MAKE,'') <> ''
			AND a.RemQty > 0
			AND b.GSBH = ?
	`, scannedCLBHs, GSBH).Scan(&allSS).Error

	if err != nil {
		return nil, err
	}

	// =========================
	// 3. Lấy tất cả ddbh từ Memo_RY
	// =========================
	ddbhSet := make(map[string]bool)

	for _, ss := range allSS {
		if ss.Memo_RY == "" {
			continue
		}
		for _, p := range strings.Split(ss.Memo_RY, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				ddbhSet[p] = true
			}
		}
	}

	var ddbhList []string
	for k := range ddbhSet {
		ddbhList = append(ddbhList, k)
	}

	// =========================
	// 4. Query ddzl 1 lần
	// =========================
	type DdRow struct {
		Ddbh  string
		BuyNo string
	}

	ddMap := make(map[string]map[string]bool)

	if len(ddbhList) > 0 {
		var ddRows []DdRow

		err = db.Raw(`
			SELECT ddbh, BuyNo 
			FROM ddzl 
			WHERE ddbh IN (?)
		`, ddbhList).Scan(&ddRows).Error

		if err != nil {
			return nil, err
		}

		for _, d := range ddRows {
			if _, ok := ddMap[d.Ddbh]; !ok {
				ddMap[d.Ddbh] = make(map[string]bool)
			}
			ddMap[d.Ddbh][d.BuyNo] = true
		}
	}

	// =========================
	// 5. Gắn BuyNos vào allSS
	// =========================
	for i, ss := range allSS {
		if ss.Memo_RY == "" {
			continue
		}

		buySet := make(map[string]bool)

		for _, p := range strings.Split(ss.Memo_RY, ",") {
			p = strings.TrimSpace(p)
			if list, ok := ddMap[p]; ok {
				for b := range list {
					buySet[b] = true
				}
			}
		}

		var buys []string
		for b := range buySet {
			buys = append(buys, b)
		}

		sort.Strings(buys)
		allSS[i].BuyNos = strings.Join(buys, ",")
	}

	// =========================
	// 6. Barcode scan hôm nay
	// =========================
	var todayBarcodes []string
	err = db.Raw(`
		SELECT barcode 
		FROM KCRKScan_RFKK 
		WHERE CONVERT(varchar(10), USERDATE, 23) = ? AND GSBH = ?
	`, Date, GSBH).Scan(&todayBarcodes).Error

	if err != nil {
		return nil, err
	}

	isScannedToday := make(map[string]bool)
	for _, b := range todayBarcodes {
		isScannedToday[b] = true
	}

	// =========================
	// 7. Build result
	// =========================
	type FinalRow struct {
		RFSSRaw
		Status string
	}

	var midRows []FinalRow

	for _, ss := range allSS {
		if ss.RemQty <= 0 {
			continue
		}

		status := "0"
		if isScannedToday[ss.Barcode] {
			status = "1"
		}

		midRows = append(midRows, FinalRow{
			RFSSRaw: ss,
			Status:  status,
		})
	}

	// sort
	sort.Slice(midRows, func(i, j int) bool {
		if midRows[i].CLBH != midRows[j].CLBH {
			return midRows[i].CLBH < midRows[j].CLBH
		}
		if midRows[i].Status != midRows[j].Status {
			return midRows[i].Status > midRows[j].Status
		}
		return midRows[i].Barcode < midRows[j].Barcode
	})

	// =========================
	// 8. Response
	// =========================
	var result []types.ListInventory

	for _, r := range midRows {
		qtyScan := 0.0
		if r.Status == "1" {
			qtyScan = r.RemQty
		}

		cfmDate := ""
		if !r.CFMDate.IsZero() {
			cfmDate = r.CFMDate.Format("2006-01-02")
		}

		result = append(result, types.ListInventory{
			GSBH: GSBH,

			Barcode: r.Barcode,
			Rack:    r.MAKE,
			CLBH:    r.CLBH,
			Pack:    r.Pack,

			RemQty:  formatRound1(r.RemQty),
			CFMDate: cfmDate,
			QtyScan: formatRound1(qtyScan),

			Status:       r.Status,
			Memo_Article: r.Memo_Article,
			Memo_RY:      r.Memo_RY,
			BuyNos:       r.BuyNos,
		})
	}

	return result, nil
}

// SelectInventoryByMonth: đọc trực tiếp từ KCRKScan_RFKK_MONTH (đã có đủ dữ liệu chi tiết
// sau khi bấm nút Đồng bộ), KHÔNG join trực tiếp KCRKScan_RFSS nữa.
// Trả kèm Qty gốc và thông tin người+ngày thực hiện Sync / Round1 (đơn vị kiểm kê) / Round2 (kế toán).
func (s *DiscrepanciesService) SelectInventoryByMonth(Month string, GSBH string, Rack string) ([]types.ListInventory, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	if strings.TrimSpace(GSBH) == "" {
		return []types.ListInventory{}, nil
	}

	rackPattern := strings.TrimSpace(Rack) + "%"

	type MonthRow struct {
		Barcode      string
		Rack         string
		CLBH         string
		Pack         string
		Memo_Article string
		Memo_RY      string
		Qty          float64
		RemQty       float64
		CFMDate      time.Time
		BuyNos       string
		TotalScan    int

		SyncUserID   string
		SyncDate     time.Time
		Round1UserID string
		Round1Date   time.Time
		Round2UserID string
		Round2Date   time.Time
	}

	var rows []MonthRow
	err = db.Raw(`
		SELECT
			Barcode,
			ISNULL(Rack, '')         AS Rack,
			ISNULL(CLBH, '')         AS CLBH,
			ISNULL(Pack, '')         AS Pack,
			ISNULL(Memo_Article, '') AS Memo_Article,
			ISNULL(Memo_RY, '')      AS Memo_RY,
			ISNULL(Qty, 0)           AS Qty,
			ISNULL(RemQty, 0)        AS RemQty,
			CFMDate,
			ISNULL(BuyNos, '')       AS BuyNos,
			ISNULL(TotalScan, 0)     AS TotalScan,
			ISNULL(SyncUserID, '')   AS SyncUserID,
			SyncDate,
			ISNULL(Round1UserID, '') AS Round1UserID,
			Round1Date,
			ISNULL(Round2UserID, '') AS Round2UserID,
			Round2Date
		FROM KCRKScan_RFKK_MONTH
		WHERE YearMonth = ?
		  AND GSBH = ?
		  AND Rack LIKE ?
		  AND ISNULL(RemQty, 0) > 0
	`, Month, GSBH, rackPattern).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].CLBH != rows[j].CLBH {
			return rows[i].CLBH < rows[j].CLBH
		}
		iScanned := rows[i].TotalScan > 0
		jScanned := rows[j].TotalScan > 0
		if iScanned != jScanned {
			return iScanned
		}
		return rows[i].Barcode < rows[j].Barcode
	})

	formatDate := func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Format("2006-01-02")
	}

	result := make([]types.ListInventory, 0, len(rows))
	for _, r := range rows {
		status := "0"
		qtyScan := 0.0
		if r.TotalScan > 0 {
			status = "1"
			qtyScan = r.RemQty
		}

		result = append(result, types.ListInventory{
			GSBH:         GSBH,
			Barcode:      r.Barcode,
			Rack:         r.Rack,
			CLBH:         r.CLBH,
			Pack:         r.Pack,
			Qty:          formatRound1(r.Qty),
			RemQty:       formatRound1(r.RemQty),
			CFMDate:      formatDate(r.CFMDate),
			QtyScan:      formatRound1(qtyScan),
			Status:       status,
			Memo_Article: r.Memo_Article,
			Memo_RY:      r.Memo_RY,
			BuyNos:       r.BuyNos,

			SyncUserID:   r.SyncUserID,
			SyncDate:     formatDate(r.SyncDate),
			Round1UserID: r.Round1UserID,
			Round1Date:   formatDate(r.Round1Date),
			Round2UserID: r.Round2UserID,
			Round2Date:   formatDate(r.Round2Date),
		})
	}

	return result, nil
}

// SyncInventoryMonth: kéo dữ liệu tồn kho từ KCRKScan_RFSS (kèm Qty gốc) và MERGE (insert dòng mới /
// update dòng đã có) vào KCRKScan_RFKK_MONTH theo (YearMonth, Barcode, GSBH).
// Đồng thời lưu lại người bấm Đồng bộ (SyncUserID) và ngày bấm (SyncDate).
func (s *DiscrepanciesService) SyncInventoryMonth(Month string, GSBH string, Rack string, UserID string) (int64, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return 0, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	if strings.TrimSpace(Month) == "" || strings.TrimSpace(GSBH) == "" {
		return 0, errors.New("Month và GSBH là bắt buộc")
	}
	if strings.TrimSpace(UserID) == "" {
		return 0, errors.New("UserID là bắt buộc để lưu người thực hiện đồng bộ")
	}

	rackPattern := strings.TrimSpace(Rack) + "%"
	const insertBatch = 500

	type RFSSRaw struct {
		Barcode      string
		MAKE         string
		CLBH         string
		Pack         string
		Memo_Article string
		Memo_RY      string
		Qty          float64
		RemQty       float64
		CFMDate      time.Time
		BuyNos       string
	}

	var synced int64

	err = db.Transaction(func(tx *gorm.DB) error {
		insertToTemp := func(tableName, colName string, values []string) error {
			if len(values) == 0 {
				return nil
			}
			for i := 0; i < len(values); i += insertBatch {
				end := i + insertBatch
				if end > len(values) {
					end = len(values)
				}
				batch := values[i:end]

				placeholders := make([]string, len(batch))
				args := make([]interface{}, len(batch))
				for j, v := range batch {
					placeholders[j] = "(?)"
					args[j] = v
				}

				sql := "INSERT INTO " + tableName + " (" + colName + ") VALUES " + strings.Join(placeholders, ",")
				if err := tx.Exec(sql, args...).Error; err != nil {
					return err
				}
			}
			return nil
		}

		// =========================
		// 1) Lấy toàn bộ tồn kho đang còn hàng, đúng rack, để đồng bộ (kèm Qty gốc)
		// =========================
		var allSS []RFSSRaw
		if err := tx.Raw(`
			SELECT
				r.barcode      AS Barcode,
				r.MAKE         AS MAKE,
				r.CLBH         AS CLBH,
				r.Pack         AS Pack,
				r.Memo_Article AS Memo_Article,
				r.Memo_RY      AS Memo_RY,
				r.Qty          AS Qty,
				r.RemQty       AS RemQty,
				r.CFMDate      AS CFMDate
			FROM KCRKScan_RFSS r
			LEFT JOIN KCRKScan_RF f ON r.SCNO = f.SCNO
			WHERE r.MAKE LIKE ?
			  AND r.RemQty > 0
			  AND f.GSBH = ?
		`, rackPattern, GSBH).Scan(&allSS).Error; err != nil {
			return err
		}

		if len(allSS) == 0 {
			synced = 0
			return nil
		}

		// =========================
		// 2) Gắn BuyNos qua ddzl (giống logic cũ)
		// =========================
		ddbhSet := make(map[string]struct{})
		for _, ss := range allSS {
			if strings.TrimSpace(ss.Memo_RY) == "" {
				continue
			}
			for _, p := range strings.Split(ss.Memo_RY, ",") {
				p = strings.TrimSpace(p)
				if p != "" {
					ddbhSet[p] = struct{}{}
				}
			}
		}

		ddMap := make(map[string]map[string]bool)
		if len(ddbhSet) > 0 {
			if err := tx.Exec(`IF OBJECT_ID('tempdb..#TempDDBH') IS NOT NULL DROP TABLE #TempDDBH;`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`
				CREATE TABLE #TempDDBH (
					Ddbh NVARCHAR(100) COLLATE DATABASE_DEFAULT NOT NULL
				);
			`).Error; err != nil {
				return err
			}

			ddbhList := make([]string, 0, len(ddbhSet))
			for k := range ddbhSet {
				ddbhList = append(ddbhList, k)
			}
			if err := insertToTemp("#TempDDBH", "Ddbh", ddbhList); err != nil {
				return err
			}

			var ddRows []struct {
				Ddbh  string
				BuyNo string
			}
			if err := tx.Raw(`
				SELECT d.ddbh AS Ddbh, d.BuyNo AS BuyNo
				FROM ddzl d
				INNER JOIN #TempDDBH t ON d.ddbh COLLATE DATABASE_DEFAULT = t.Ddbh
			`).Scan(&ddRows).Error; err != nil {
				return err
			}

			for _, d := range ddRows {
				if _, ok := ddMap[d.Ddbh]; !ok {
					ddMap[d.Ddbh] = make(map[string]bool)
				}
				ddMap[d.Ddbh][d.BuyNo] = true
			}
		}

		for i, ss := range allSS {
			if strings.TrimSpace(ss.Memo_RY) == "" {
				continue
			}
			buySet := make(map[string]bool)
			for _, p := range strings.Split(ss.Memo_RY, ",") {
				p = strings.TrimSpace(p)
				if list, ok := ddMap[p]; ok {
					for b := range list {
						buySet[b] = true
					}
				}
			}
			var buys []string
			for b := range buySet {
				buys = append(buys, b)
			}
			sort.Strings(buys)
			allSS[i].BuyNos = strings.Join(buys, ",")
		}

		// =========================
		// 3) Đổ toàn bộ allSS vào temp table rồi MERGE 1 lần vào KCRKScan_RFKK_MONTH
		// =========================
		if err := tx.Exec(`IF OBJECT_ID('tempdb..#TempSyncMonth') IS NOT NULL DROP TABLE #TempSyncMonth;`).Error; err != nil {
			return err
		}
		if err := tx.Exec(`
			CREATE TABLE #TempSyncMonth (
				Barcode      NVARCHAR(100) COLLATE DATABASE_DEFAULT NOT NULL,
				Rack         NVARCHAR(100) COLLATE DATABASE_DEFAULT NULL,
				CLBH         NVARCHAR(150) COLLATE DATABASE_DEFAULT NULL,
				Pack         NVARCHAR(100) COLLATE DATABASE_DEFAULT NULL,
				Memo_Article NVARCHAR(500) COLLATE DATABASE_DEFAULT NULL,
				Memo_RY      NVARCHAR(500) COLLATE DATABASE_DEFAULT NULL,
				Qty          FLOAT NULL,
				RemQty       FLOAT NULL,
				CFMDate      DATETIME NULL,
				BuyNos       NVARCHAR(1000) COLLATE DATABASE_DEFAULT NULL
			);
		`).Error; err != nil {
			return err
		}

		const syncInsertBatch = 190 // 10 cột/dòng => 190*10 = 1900 tham số, dưới giới hạn 2100 của SQL Server
		for i := 0; i < len(allSS); i += syncInsertBatch {
			end := i + syncInsertBatch
			if end > len(allSS) {
				end = len(allSS)
			}
			batch := allSS[i:end]

			placeholders := make([]string, len(batch))
			args := make([]interface{}, 0, len(batch)*10)
			for j, ss := range batch {
				placeholders[j] = "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
				var cfmDate interface{}
				if !ss.CFMDate.IsZero() {
					cfmDate = ss.CFMDate
				} else {
					cfmDate = nil
				}
				args = append(args,
					truncateStr(ss.Barcode, 100),
					truncateStr(ss.MAKE, 100),
					truncateStr(ss.CLBH, 150),
					truncateStr(ss.Pack, 100),
					truncateStr(ss.Memo_Article, 500),
					truncateStr(ss.Memo_RY, 500),
					ss.Qty,
					ss.RemQty,
					cfmDate,
					truncateStr(ss.BuyNos, 1000),
				)
			}

			sqlInsert := `INSERT INTO #TempSyncMonth
				(Barcode, Rack, CLBH, Pack, Memo_Article, Memo_RY, Qty, RemQty, CFMDate, BuyNos)
				VALUES ` + strings.Join(placeholders, ",")
			if err := tx.Exec(sqlInsert, args...).Error; err != nil {
				return err
			}
		}

		queryMerge := `
			MERGE KCRKScan_RFKK_MONTH AS target
			USING #TempSyncMonth AS source
				ON target.YearMonth = ?
				AND target.GSBH     = ?
				AND target.Barcode  = source.Barcode

			WHEN MATCHED THEN
				UPDATE SET
					Rack         = source.Rack,
					CLBH         = source.CLBH,
					Pack         = source.Pack,
					Memo_Article = source.Memo_Article,
					Memo_RY      = source.Memo_RY,
					Qty          = source.Qty,
					RemQty       = source.RemQty,
					CFMDate      = source.CFMDate,
					BuyNos       = source.BuyNos,
					SyncUserID   = ?,
					SyncDate     = GETDATE()

			WHEN NOT MATCHED BY TARGET THEN
				INSERT (YearMonth, Barcode, GSBH, TotalScan,
					Rack, CLBH, Pack, Memo_Article, Memo_RY, Qty, RemQty, CFMDate, BuyNos,
					SyncUserID, SyncDate, CreatedAt)
				VALUES (?, source.Barcode, ?, 0,
					source.Rack, source.CLBH, source.Pack, source.Memo_Article, source.Memo_RY,
					source.Qty, source.RemQty, source.CFMDate, source.BuyNos,
					?, GETDATE(), GETDATE());
		`
		mergeTx := tx.Exec(queryMerge, Month, GSBH, UserID, Month, GSBH, UserID)
		if mergeTx.Error != nil {
			return mergeTx.Error
		}
		synced = mergeTx.RowsAffected

		// =========================
		// 4) Dọn các dòng đã xuất kho hết (RemQty = 0, không còn trong #TempSyncMonth)
		//    - Chưa có Round2UserID (kế toán chưa xác nhận vòng 2) -> XÓA, kể cả đã có Round1UserID
		//    - Đã có Round2UserID -> GIỮ LẠI, không đụng tới
		//    Chỉ xóa trong đúng phạm vi YearMonth + GSBH đang sync.
		// =========================
		queryCleanup := `
			DELETE target
			FROM KCRKScan_RFKK_MONTH AS target
			WHERE target.YearMonth = ?
			  AND target.GSBH     = ?
			  AND (target.Round2UserID IS NULL OR LTRIM(RTRIM(target.Round2UserID)) = '')
			  AND NOT EXISTS (
					SELECT 1 FROM #TempSyncMonth s
					WHERE s.Barcode = target.Barcode
			  );
		`
		if err := tx.Exec(queryCleanup, Month, GSBH).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return synced, nil
}

// ConfirmRound2ByBarcode: kế toán quét barcode để xác nhận kiểm kê lần 2 — cùng cơ chế quét
// như đơn vị kiểm kê (CreateInventoryByMonth) nhưng ghi vào Round2UserID/Round2Date thay vì Round1.
// Chỉ xác nhận được cho barcode đã có dòng trong KCRKScan_RFKK_MONTH (đã Sync hoặc đã có Round1).
func (s *DiscrepanciesService) ConfirmRound2ByBarcode(param *types.Inventory) (string, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return "Connection Database Fail", err
	}
	db = db.Debug()
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// 1. Kiểm tra barcode có trong kho (KCRKScan_RFSS) không
	queryCheck := `SELECT cfmid AS Cfmid, CLBH FROM KCRKScan_RFSS WHERE barcode = ?`
	var checkResult CheckResult
	err = db.Raw(queryCheck, param.Barcode).Scan(&checkResult).Error
	if err != nil {
		return "Query Check Fail", err
	}
	if checkResult.Cfmid == "" {
		return "", errors.New("Chưa Nhập Kho")
	}

	// 2. Đã được kế toán xác nhận lần 2 trước đó chưa (Barcode + Month + GSBH)
	queryCheck2 := `
		SELECT ISNULL(Round2UserID, '') FROM KCRKScan_RFKK_MONTH
		WHERE Barcode = ? AND YearMonth = ? AND GSBH = ?
	`
	var existingRound2UserID string
	db.Raw(queryCheck2, param.Barcode, param.USERDATE, param.GSBH).Scan(&existingRound2UserID)
	if existingRound2UserID != "" {
		return "", errors.New("Đã được " + existingRound2UserID + " xác nhận lần 2 cho " + param.GSBH)
	}

	// 3. Ghi nhận Round2UserID/Round2Date — dòng phải đã tồn tại (đã Sync hoặc đã kiểm kê lần 1)
	result := db.Exec(`
		UPDATE KCRKScan_RFKK_MONTH
		SET Round2UserID = ?,
			Round2Date    = GETDATE()
		WHERE Barcode = ? AND YearMonth = ? AND GSBH = ?
	`, param.USERID, param.Barcode, param.USERDATE, param.GSBH)
	if result.Error != nil {
		return "Update Fail: " + result.Error.Error(), result.Error
	}
	if result.RowsAffected == 0 {
		return "", errors.New("Vật tư chưa được đồng bộ / kiểm kê lần 1, không thể xác nhận lần 2")
	}

	return checkResult.CLBH, nil
}

// ConfirmRound2: kế toán xác nhận kiểm kê lần 2 cho 1 barcode cụ thể trong tháng/GSBH (xác nhận thủ công qua nút bấm).
func (s *DiscrepanciesService) ConfirmRound2(Month string, GSBH string, Barcode string, UserID string) (int64, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return 0, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	Month = strings.TrimSpace(Month)
	GSBH = strings.TrimSpace(GSBH)
	Barcode = strings.TrimSpace(Barcode)
	UserID = strings.TrimSpace(UserID)

	if Month == "" || GSBH == "" || Barcode == "" || UserID == "" {
		return 0, errors.New("Date, GSBH, Barcode và UserID là bắt buộc")
	}

	result := db.Exec(`
		UPDATE KCRKScan_RFKK_MONTH
		SET Round2UserID = ?,
			Round2Date    = GETDATE()
		WHERE YearMonth = ?
		  AND GSBH      = ?
		  AND Barcode   = ?
	`, UserID, Month, GSBH, Barcode)

	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, errors.New("Không tìm thấy dòng kiểm kê phù hợp để xác nhận lần 2")
	}

	return result.RowsAffected, nil
}

// GetLastSync: lấy thời điểm + người thực hiện Đồng bộ gần nhất cho 1 Month + GSBH.
// Dùng để hiển thị "lần đồng bộ cuối" cạnh nút Đồng bộ trên FE.
func (s *DiscrepanciesService) GetLastSync(Month string, GSBH string) (map[string]interface{}, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	Month = strings.TrimSpace(Month)
	GSBH = strings.TrimSpace(GSBH)
	if Month == "" || GSBH == "" {
		return nil, errors.New("Month và GSBH là bắt buộc")
	}

	type LastSyncRow struct {
		SyncUserID string
		SyncDate   time.Time
	}

	var row LastSyncRow
	err = db.Raw(`
		SELECT TOP 1
			ISNULL(SyncUserID, '') AS SyncUserID,
			SyncDate
		FROM KCRKScan_RFKK_MONTH
		WHERE YearMonth = ?
		  AND GSBH = ?
		  AND SyncDate IS NOT NULL
		ORDER BY SyncDate DESC
	`, Month, GSBH).Scan(&row).Error
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"SyncUserID": row.SyncUserID,
		"SyncDate":   "",
	}
	if !row.SyncDate.IsZero() {
		result["SyncDate"] = row.SyncDate.Format("2006-01-02 15:04:05")
	}
	return result, nil
}

func (s *DiscrepanciesService) GetRackList(GSBH string) ([]map[string]interface{}, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	type RackRow struct {
		Rack  string
		Total int
	}

	var rows []RackRow
	err = db.Raw(`
        SELECT
            LEFT(LTRIM(RTRIM(MAKE)), 3) AS Rack,
            COUNT(*) AS Total
        FROM KCRKScan_RFSS
        WHERE RemQty <> 0
          AND ISNULL(MAKE, '') <> ''
        GROUP BY LEFT(LTRIM(RTRIM(MAKE)), 3)
        ORDER BY LEFT(LTRIM(RTRIM(MAKE)), 3)
    `).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		result = append(result, map[string]interface{}{
			"Rack":       r.Rack,
			"TotalCount": r.Total,
		})
	}
	return result, nil
}