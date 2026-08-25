package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"
	"web-api/internal/realtime"
)

// ─────────────────────────────────────────────
// LOSS SERVICE (code gốc - đã chỉnh kiểu float)
// ─────────────────────────────────────────────

var Loss = &LossService{}

type LossService struct{}

func (s *LossService) GetLossList(req types.LossListRequest) (*types.LossListResponse, error) {
	db, err := database.LYS_ERP_Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}

	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize
	startRow := offset + 1
	endRow := offset + req.PageSize

	whereClause := "1=1"
	args := []interface{}{}

	if req.ZLBH != "" {
		whereClause += " AND d.ZLBH = ?"
		args = append(args, req.ZLBH)
	}
	if req.DDBH != "" {
		whereClause += " AND d.DDBH = ?"
		args = append(args, req.DDBH)
	}

	var total int64
	var rows []types.LossItem

	err = db.Connection(func(tx *gorm.DB) error {
		if err := tx.Exec("SET ARITHABORT ON").Error; err != nil {
			return err
		}

		countSQL := fmt.Sprintf(`
			SELECT COUNT(*)
			FROM DDZL d
			JOIN XXZL x
				ON d.XieXing = x.XieXing
				AND d.SheHao = x.SheHao
			WHERE %s
		`, whereClause)

		if err := tx.Raw(countSQL, args...).Scan(&total).Error; err != nil {
			fmt.Println("Count error:", err)
			return err
		}

		dataSQL := fmt.Sprintf(`
			WITH base AS (
				SELECT
					d.ZLBH    AS zlbh,
					d.DDBH    AS ddbh,
					d.ARTICLE AS article,
					x.XieMing AS xie_ming,
					CAST(ISNULL(d.Pairs, 0) AS FLOAT) AS pairs,
					ISNULL(sz.xxcc, '') AS xxcc,
					ISNULL(sz.ywpm, '') AS ywpm,
					ROW_NUMBER() OVER (ORDER BY d.ZLBH) AS rn
				FROM DDZL d
				JOIN XXZL x
					ON d.XieXing = x.XieXing
					AND d.SheHao = x.SheHao
				OUTER APPLY (
					SELECT
						STUFF((
							SELECT ',' + t.XXCC
							FROM (
								SELECT DISTINCT c2.XXCC
								FROM CGZLSS c2
								WHERE c2.ZLBH = d.ZLBH
									AND c2.CLBH LIKE 'N04%%'
									AND c2.XXCC IS NOT NULL
									AND LTRIM(RTRIM(c2.XXCC)) <> ''
							) t
							ORDER BY t.XXCC
							FOR XML PATH(''), TYPE
						).value('.', 'NVARCHAR(MAX)'), 1, 1, '') AS xxcc,
						MAX(cl.YWPM) AS ywpm
					FROM CGZLSS c
					LEFT JOIN CLZL cl
						ON c.CLBH = cl.CLDH
					WHERE c.ZLBH = d.ZLBH
						AND c.CLBH LIKE 'N04%%'
				) sz
				WHERE %s
			)
			SELECT
				zlbh,
				ddbh,
				article,
				xie_ming,
				pairs,
				xxcc,
				ywpm
			FROM base
			WHERE rn BETWEEN ? AND ?
			ORDER BY rn
		`, whereClause)

		argsData := append([]interface{}{}, args...)
		argsData = append(argsData, startRow, endRow)

		if err := tx.Raw(dataSQL, argsData...).Scan(&rows).Error; err != nil {
			fmt.Println("Query error:", err)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("lỗi lấy danh sách hao hụt: %w", err)
	}

	return &types.LossListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Items:    rows,
	}, nil
}

var LossTemSize = &LossTemSizeService{}

type LossTemSizeService struct{}

// ─────────────────────────────────────────────
func (s *LossTemSizeService) Save(req types.LossTemSizeSaveRequest) error {
	req.DDBH = strings.TrimSpace(req.DDBH)
	req.MSNV = strings.TrimSpace(req.MSNV)
	req.MSNVOut = strings.TrimSpace(req.MSNVOut)
	req.YWPM = strings.TrimSpace(req.YWPM)
	req.Article = strings.TrimSpace(req.Article)
	req.XieMing = strings.TrimSpace(req.XieMing)
	req.Reason = strings.TrimSpace(req.Reason)
	req.Date = strings.TrimSpace(req.Date)
	req.Note = strings.TrimSpace(req.Note)
	req.ProductCode = strings.TrimSpace(req.ProductCode)

	if req.DDBH == "" {
		return fmt.Errorf("DDBH là bắt buộc")
	}
	if req.MSNV == "" {
		return fmt.Errorf("MSNV là bắt buộc")
	}
	if req.Mode != "IN" && req.Mode != "OUT" {
		return fmt.Errorf("mode phải là 'IN' hoặc 'OUT'")
	}
	if req.Mode == "OUT" {
		if req.Reason == "" {
			return fmt.Errorf("nguyên nhân là bắt buộc khi mode = OUT")
		}
		if req.MSNVOut == "" {
			return fmt.Errorf("MSNV người lấy là bắt buộc khi mode = OUT")
		}
	}

	type sizeRow struct {
		size string
		qty  float64
	}

	var sizes []sizeRow
	for sz, qty := range req.Sizes {
		sz = strings.TrimSpace(sz)
		if sz == "" {
			continue
		}
		if qty < 0 {
			continue
		}
		sizes = append(sizes, sizeRow{size: sz, qty: qty})
	}

	if len(sizes) == 0 {
		return fmt.Errorf("không có size hợp lệ để lưu")
	}

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return fmt.Errorf("lỗi kết nối DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("lỗi lấy sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.Connection(func(tx *gorm.DB) error {
		if err := tx.Exec("SET ARITHABORT ON").Error; err != nil {
			return err
		}

		if req.Mode == "IN" {
			mergeInfoSQL := `
				MERGE LossTemSize_Info AS target
				USING (
					SELECT
						? AS DDBH,
						? AS MSNV,
						? AS YWPM,
						? AS ARTICLE,
						? AS XIEMING,
						? AS PAIRS
				) AS source
				ON target.DDBH = source.DDBH
				WHEN NOT MATCHED THEN
					INSERT (DDBH, MSNV, YWPM, ARTICLE, XIEMING, PAIRS)
					VALUES (source.DDBH, source.MSNV, source.YWPM, source.ARTICLE, source.XIEMING, source.PAIRS);
			`

			if err := tx.Exec(
				mergeInfoSQL,
				req.DDBH,
				req.MSNV,
				req.YWPM,
				req.Article,
				req.XieMing,
				req.Pairs,
			).Error; err != nil {
				return fmt.Errorf("lỗi upsert LossTemSize_Info: %w", err)
			}

			for _, r := range sizes {
				var totalOut float64
				if err := tx.Raw(`
					SELECT COALESCE(SUM(CAST(QTYOUT AS FLOAT)), 0)
					FROM LossTemSize_Out
					WHERE DDBH = ? AND XXCC = ?
				`, req.DDBH, r.size).Scan(&totalOut).Error; err != nil {
					return fmt.Errorf("lỗi kiểm tra OUT size=%s: %w", r.size, err)
				}

				if r.qty == 0 && totalOut > 0 {
					return fmt.Errorf("❌ Size %s: số lượng nhập (0) nhỏ hơn số đã lấy ra (%.2f) — không cho phép nhập lại", r.size, totalOut)
				}

				if r.qty > 0 && r.qty < totalOut {
					return fmt.Errorf("❌ Size %s: số lượng nhập (%.2f) nhỏ hơn số đã lấy ra (%.2f) — không cho phép nhập lại", r.size, r.qty, totalOut)
				}

				mergeInSQL := `
					MERGE LossTemSize_In AS target
					USING (
						SELECT
							? AS DDBH,
							? AS XXCC
					) AS source
					ON target.DDBH = source.DDBH
					AND target.XXCC = source.XXCC
					WHEN MATCHED THEN
						UPDATE SET
							USERID = ?,
							QTYIN  = ?,
							PRODUCT_CODE = ?,
							[DATE] = GETDATE()
					WHEN NOT MATCHED THEN
						INSERT (DDBH, XXCC, USERID, QTYIN, PRODUCT_CODE, [DATE])
						VALUES (?, ?, ?, ?, ?, GETDATE());
				`

				if err := tx.Exec(
					mergeInSQL,
					req.DDBH,
					r.size,
					req.MSNV,
					r.qty,
					req.ProductCode,
					req.DDBH,
					r.size,
					req.MSNV,
					r.qty,
					req.ProductCode,
				).Error; err != nil {
					return fmt.Errorf("lỗi upsert LossTemSize_In size=%s: %w", r.size, err)
				}
			}
		} else {
			var sb strings.Builder
			sb.WriteString(`
				INSERT INTO LossTemSize_Out (DDBH, XXCC, USERID, USERIDOUT, QTYOUT, NN, [DATE])
				VALUES
			`)

			args := []interface{}{}
			for i, r := range sizes {
				if i > 0 {
					sb.WriteString(",\n")
				}
				sb.WriteString("(?, ?, ?, ?, ?, ?, GETDATE())")
				args = append(args, req.DDBH, r.size, req.MSNV, req.MSNVOut, r.qty, req.Reason)
			}

			if err := tx.Exec(sb.String(), args...).Error; err != nil {
				return fmt.Errorf("lỗi bulk INSERT LossTemSize_Out: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	broadcastMsg, _ := json.Marshal(map[string]any{
		"event":    "loss_tem_size_updated",
		"ddbh":     req.DDBH,
		"msnv":     req.MSNV,
		"msnv_out": req.MSNVOut,
		"mode":     req.Mode,
		"reason":   req.Reason,
		"date":     now,
	})
	realtime.AlertHub.Broadcast(broadcastMsg)

	return nil
}
// ─────────────────────────────────────────────
// GetList – phân trang, JOIN 3 bảng
// ─────────────────────────────────────────────
func (s *LossTemSizeService) GetList(req types.LossTemSizeListRequest) (*types.LossTemSizeListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	offset := (req.Page - 1) * req.PageSize
	startRow := offset + 1
	endRow := offset + req.PageSize

	infoWhere := "1=1"
	args := []interface{}{}

	if req.DDBH != "" {
		infoWhere += " AND i.DDBH = ?"
		args = append(args, req.DDBH)
	}
	if req.MSNV != "" {
		infoWhere += " AND i.MSNV = ?"
		args = append(args, req.MSNV)
	}

	dateWhere := "1=1"
	if req.FromDate != "" {
		dateWhere += " AND t.[DATE] >= ?"
		args = append(args, req.FromDate)
	}
	if req.ToDate != "" {
		dateWhere += " AND t.[DATE] <= ?"
		args = append(args, req.ToDate)
	}

	type rawRow struct {
		ID          int64   `gorm:"column:id"`
		DDBH        string  `gorm:"column:ddbh"`
		MSNV        string  `gorm:"column:msnv"`
		YWPM        string  `gorm:"column:ywpm"`
		Article     string  `gorm:"column:article"`
		XieMing     string  `gorm:"column:xie_ming"`
		Pairs       float64 `gorm:"column:pairs"`
		XXCC        string  `gorm:"column:xxcc"`
		Mode        string  `gorm:"column:mode"`
		Reason      string  `gorm:"column:reason"`
		QtyIn       float64 `gorm:"column:qty_in"`
		QtyOut      float64 `gorm:"column:qty_out"`
		UserID      string  `gorm:"column:userid"`
		Date        string  `gorm:"column:date"`
		ProductCode string  `gorm:"column:product_code"`
	}

	var total int64
	var rawRows []rawRow

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối DB: %w", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	err = db.Connection(func(tx *gorm.DB) error {
		if err := tx.Exec("SET ARITHABORT ON").Error; err != nil {
			return err
		}

		unionSQL := fmt.Sprintf(`
			SELECT COUNT(*) FROM (
				SELECT t.ID
				FROM LossTemSize_In t
				JOIN LossTemSize_Info i ON i.DDBH = t.DDBH
				WHERE %s AND %s
				UNION ALL
				SELECT t.ID
				FROM LossTemSize_Out t
				JOIN LossTemSize_Info i ON i.DDBH = t.DDBH
				WHERE %s AND %s
			) u
		`, infoWhere, dateWhere, infoWhere, dateWhere)

		countArgs := append(args, args...)
		if err := tx.Raw(unionSQL, countArgs...).Scan(&total).Error; err != nil {
			fmt.Println("Count error:", err)
			return err
		}

		dataSQL := fmt.Sprintf(`
			WITH combined AS (
				SELECT
					t.ID                               AS id,
					i.DDBH                             AS ddbh,
					i.MSNV                             AS msnv,
					ISNULL(i.YWPM, '')                 AS ywpm,
					ISNULL(i.ARTICLE, '')              AS article,
					ISNULL(i.XIEMING, '')              AS xie_ming,
					CAST(ISNULL(i.PAIRS, 0) AS FLOAT)  AS pairs,
					ISNULL(t.XXCC, '')                 AS xxcc,
					'IN'                               AS mode,
					''                                 AS reason,
					CAST(ISNULL(t.QTYIN, 0) AS FLOAT)  AS qty_in,
					CAST(0 AS FLOAT)                   AS qty_out,
					ISNULL(t.USERID, '')               AS userid,
					CONVERT(VARCHAR(10), t.[DATE], 120) AS date,
					ISNULL(t.PRODUCT_CODE, '')         AS product_code,
					ROW_NUMBER() OVER (ORDER BY t.ID DESC) AS rn
				FROM LossTemSize_In t
				JOIN LossTemSize_Info i ON i.DDBH = t.DDBH
				WHERE %s AND %s

				UNION ALL

				SELECT
					t.ID                               AS id,
					i.DDBH                             AS ddbh,
					i.MSNV                             AS msnv,
					ISNULL(i.YWPM, '')                 AS ywpm,
					ISNULL(i.ARTICLE, '')              AS article,
					ISNULL(i.XIEMING, '')              AS xie_ming,
					CAST(ISNULL(i.PAIRS, 0) AS FLOAT)  AS pairs,
					ISNULL(t.XXCC, '')                 AS xxcc,
					'OUT'                              AS mode,
					ISNULL(t.NN, '')                   AS reason,
					CAST(0 AS FLOAT)                   AS qty_in,
					CAST(ISNULL(t.QTYOUT, 0) AS FLOAT) AS qty_out,
					ISNULL(t.USERIDOUT, '')            AS userid,
					CONVERT(VARCHAR(10), t.[DATE], 120) AS date,
					ISNULL((
						SELECT TOP 1 li.PRODUCT_CODE
						FROM LossTemSize_In li
						WHERE li.DDBH = t.DDBH
					), '') AS product_code,
					ROW_NUMBER() OVER (ORDER BY t.ID DESC) AS rn
				FROM LossTemSize_Out t
				JOIN LossTemSize_Info i ON i.DDBH = t.DDBH
				WHERE %s AND %s
			)
			SELECT id, ddbh, msnv, ywpm, article, xie_ming, pairs,
			       xxcc, mode, reason, qty_in, qty_out, userid, date, product_code
			FROM combined
			WHERE rn BETWEEN ? AND ?
			ORDER BY rn
		`, infoWhere, dateWhere, infoWhere, dateWhere)

		dataArgs := append(append(append(args, args...), args...), args...)
		dataArgs = append(dataArgs, startRow, endRow)

		if err := tx.Raw(dataSQL, dataArgs...).Scan(&rawRows).Error; err != nil {
			fmt.Println("Query error:", err)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("lỗi lấy danh sách LossTemSize: %w", err)
	}

	items := make([]types.LossTemSizeItem, 0, len(rawRows))
	for _, r := range rawRows {
		items = append(items, types.LossTemSizeItem{
			ID:          r.ID,
			DDBH:        r.DDBH,
			MSNV:        r.MSNV,
			YWPM:        r.YWPM,
			Article:     r.Article,
			XieMing:     r.XieMing,
			Pairs:       r.Pairs,
			XXCC:        r.XXCC,
			Mode:        r.Mode,
			Reason:      r.Reason,
			QtyIn:       r.QtyIn,
			QtyOut:      r.QtyOut,
			UserID:      r.UserID,
			Date:        r.Date,
			ProductCode: r.ProductCode,
		})
	}

	return &types.LossTemSizeListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Items:    items,
	}, nil
}

// ─────────────────────────────────────────────
// GetSummary – tổng hợp theo DDBH
// ─────────────────────────────────────────────
func (s *LossTemSizeService) GetSummary(req types.LossTemSizeListRequest) ([]types.LossTemSizeSummaryItem, error) {
	infoWhere := "1=1"
	args := []interface{}{}

	if req.DDBH != "" {
		infoWhere += " AND i.DDBH = ?"
		args = append(args, req.DDBH)
	}
	if req.MSNV != "" {
		infoWhere += " AND i.MSNV = ?"
		args = append(args, req.MSNV)
	}

	dateWhere := "1=1"
	if req.FromDate != "" {
		dateWhere += " AND t.[DATE] >= ?"
		args = append(args, req.FromDate)
	}
	if req.ToDate != "" {
		dateWhere += " AND t.[DATE] <= ?"
		args = append(args, req.ToDate)
	}

	type rawRow struct {
		DDBH        string  `gorm:"column:ddbh"`
		MSNV        string  `gorm:"column:msnv"`
		YWPM        string  `gorm:"column:ywpm"`
		Article     string  `gorm:"column:article"`
		XieMing     string  `gorm:"column:xie_ming"`
		Pairs       float64 `gorm:"column:pairs"`
		XXCC        string  `gorm:"column:xxcc"`
		Mode        string  `gorm:"column:mode"`
		Reason      string  `gorm:"column:reason"`
		QtyIn       float64 `gorm:"column:qty_in"`
		QtyOut      float64 `gorm:"column:qty_out"`
		UserID      string  `gorm:"column:userid"`
		Date        string  `gorm:"column:date"`
		ProductCode string  `gorm:"column:product_code"`
	}

	var rows []rawRow

	db, err := database.LYS_ERP_Connection()
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.Connection(func(tx *gorm.DB) error {
		if err := tx.Exec("SET ARITHABORT ON").Error; err != nil {
			return err
		}

		query := fmt.Sprintf(`
			SELECT
				i.DDBH                               AS ddbh,
				i.MSNV                               AS msnv,
				ISNULL(i.YWPM, '')                   AS ywpm,
				ISNULL(i.ARTICLE, '')                AS article,
				ISNULL(i.XIEMING, '')                AS xie_ming,
				CAST(ISNULL(i.PAIRS, 0) AS FLOAT)    AS pairs,
				ISNULL(t.XXCC, '')                   AS xxcc,
				'IN'                                 AS mode,
				''                                   AS reason,
				CAST(ISNULL(t.QTYIN, 0) AS FLOAT)    AS qty_in,
				CAST(0 AS FLOAT)                     AS qty_out,
				ISNULL(t.USERID, '')                 AS userid,
				CONVERT(VARCHAR(10), t.[DATE], 120)  AS date,
				ISNULL(t.PRODUCT_CODE, '')           AS product_code
			FROM LossTemSize_In t
			JOIN LossTemSize_Info i ON i.DDBH = t.DDBH
			WHERE %s AND %s

			UNION ALL

			SELECT
				i.DDBH                               AS ddbh,
				i.MSNV                               AS msnv,
				ISNULL(i.YWPM, '')                   AS ywpm,
				ISNULL(i.ARTICLE, '')                AS article,
				ISNULL(i.XIEMING, '')                AS xie_ming,
				CAST(ISNULL(i.PAIRS, 0) AS FLOAT)    AS pairs,
				ISNULL(t.XXCC, '')                   AS xxcc,
				'OUT'                                AS mode,
				ISNULL(t.NN, '')                     AS reason,
				CAST(0 AS FLOAT)                     AS qty_in,
				CAST(ISNULL(t.QTYOUT, 0) AS FLOAT)   AS qty_out,
				ISNULL(t.USERIDOUT, '')              AS userid,
				CONVERT(VARCHAR(10), t.[DATE], 120)  AS date,
				ISNULL((
					SELECT TOP 1 li.PRODUCT_CODE
					FROM LossTemSize_In li
					WHERE li.DDBH = t.DDBH
				), '') AS product_code
			FROM LossTemSize_Out t
			JOIN LossTemSize_Info i ON i.DDBH = t.DDBH
			WHERE %s AND %s

			ORDER BY ddbh, xxcc
		`, infoWhere, dateWhere, infoWhere, dateWhere)

		queryArgs := append(args, args...)

		if err := tx.Raw(query, queryArgs...).Scan(&rows).Error; err != nil {
			fmt.Println("Summary query error:", err)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("lỗi GetSummary: %w", err)
	}

	orderKeys := []string{}
	summaryMap := map[string]*types.LossTemSizeSummaryItem{}

	for _, r := range rows {
		sm, ok := summaryMap[r.DDBH]
		if !ok {
			orderKeys = append(orderKeys, r.DDBH)

			sm = &types.LossTemSizeSummaryItem{
				DDBH:        r.DDBH,
				MSNV:        r.MSNV,
				Article:     r.Article,
				XieMing:     r.XieMing,
				YWPM:        r.YWPM,
				Pairs:       r.Pairs,
				ProductCode: r.ProductCode,
				SizesIn:     map[string]float64{},
				SizesOut:    map[string]float64{},
				SizesNet:    map[string]float64{},
				OutByReason: map[string]map[string]float64{},
			}
			summaryMap[r.DDBH] = sm
		} else if sm.ProductCode == "" && r.ProductCode != "" {
			sm.ProductCode = r.ProductCode
		}

		if r.Date > sm.DATE {
			sm.DATE = r.Date
		}

		if r.Mode == "IN" {
			sm.SizesIn[r.XXCC] += r.QtyIn
			sm.SizesNet[r.XXCC] += r.QtyIn
			sm.TotalIn += r.QtyIn
		}

		if r.Mode == "OUT" {
			sm.SizesOut[r.XXCC] += r.QtyOut
			sm.SizesNet[r.XXCC] -= r.QtyOut
			sm.TotalOut += r.QtyOut

			if r.Reason != "" {
				if sm.OutByReason[r.Reason] == nil {
					sm.OutByReason[r.Reason] = map[string]float64{}
				}
				sm.OutByReason[r.Reason][r.XXCC] += r.QtyOut
			}
		}
	}

	result := make([]types.LossTemSizeSummaryItem, 0, len(orderKeys))
	for _, key := range orderKeys {
		sm := summaryMap[key]
		sm.TotalNet = sm.TotalIn - sm.TotalOut
		result = append(result, *sm)
	}

	return result, nil
}