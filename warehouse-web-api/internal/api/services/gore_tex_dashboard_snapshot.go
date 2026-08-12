package services

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"web-api/internal/pkg/models/types"

	"gorm.io/gorm"
)

func goreTexDashboardSnapshotsAvailable(db *gorm.DB) bool {
	var result struct {
		Ready int `gorm:"column:ready"`
	}
	query := db.Raw(`
		SELECT CASE WHEN
			OBJECT_ID(N'dbo.gore_tex_dashboard_pass_snapshots', N'U') IS NOT NULL AND
			OBJECT_ID(N'dbo.gore_tex_dashboard_defect_snapshots', N'U') IS NOT NULL
		THEN 1 ELSE 0 END AS ready
	`).Scan(&result)
	return query.Error == nil && result.Ready == 1
}

func syncGoreTexPassDailySnapshot(tx *gorm.DB, inspectionDate time.Time, dataJSON string, now time.Time) error {
	if !goreTexDashboardSnapshotsAvailable(tx) {
		return nil
	}
	var data goreTexDashboardCentrifugalData
	if err := json.Unmarshal([]byte(dataJSON), &data); err != nil {
		return err
	}
	byStyle := make(map[string]*goreTexRateAccumulator)
	for _, entry := range data.Entries {
		style := normalizeGoreTexStyle(entry.Style)
		result := normalizeGoreTexResult(entry.Result)
		if style != "" && (result == "PASS" || result == "FAIL") {
			addGoreTexRate(byStyle, style, result)
		}
	}
	if err := tx.Exec(`
		DELETE FROM dbo.gore_tex_dashboard_pass_snapshots
		WHERE period_type = 'D' AND period_start = ?
	`, inspectionDate).Error; err != nil {
		return err
	}
	year, week := inspectionDate.ISOWeek()
	for style, value := range byStyle {
		if err := tx.Exec(`
			INSERT INTO dbo.gore_tex_dashboard_pass_snapshots
				(period_type, period_start, iso_year, iso_week, style_name, pass_count, fail_count, total_count, updated_at)
			VALUES ('D', ?, ?, ?, ?, ?, ?, ?, ?)
		`, inspectionDate, year, week, style, value.Pass, value.Total-value.Pass, value.Total, now).Error; err != nil {
			return err
		}
	}
	return rebuildGoreTexWeeklySnapshots(tx, inspectionDate, now)
}

func syncGoreTexDefectDailySnapshot(tx *gorm.DB, inspectionDate time.Time, recordKey, dataJSON string, now time.Time) error {
	if !goreTexDashboardSnapshotsAvailable(tx) {
		return nil
	}
	var data goreTexDashboardWaterproofData
	if err := json.Unmarshal([]byte(dataJSON), &data); err != nil {
		return err
	}
	if err := tx.Exec(`
		DELETE FROM dbo.gore_tex_dashboard_defect_snapshots
		WHERE period_type = 'D' AND period_start = ? AND source_record_key = ?
	`, inspectionDate, recordKey).Error; err != nil {
		return err
	}
	year, week := inspectionDate.ISOWeek()
	for index, group := range goreTexDefectGroups {
		count := goreTexDefectGroupCount(data, group.Rows)
		if count <= 0 {
			continue
		}
		if err := tx.Exec(`
			INSERT INTO dbo.gore_tex_dashboard_defect_snapshots
				(period_type, period_start, iso_year, iso_week, source_record_key, defect_code, defect_name, defect_count, updated_at)
			VALUES ('D', ?, ?, ?, ?, ?, ?, ?, ?)
		`, inspectionDate, year, week, recordKey, goreTexDefectCode(index), group.Label, count, now).Error; err != nil {
			return err
		}
	}
	return rebuildGoreTexWeeklySnapshots(tx, inspectionDate, now)
}

func rebuildGoreTexWeeklySnapshots(tx *gorm.DB, date time.Time, now time.Time) error {
	year, week := date.ISOWeek()
	weekStart, err := goreTexISOWeekStart(year, week)
	if err != nil {
		return err
	}
	weekEnd := weekStart.AddDate(0, 0, 7)
	lockResource := fmt.Sprintf("gore-tex-dashboard-week-%04d-%02d", year, week)
	if err := tx.Exec(`
		DECLARE @lock_result INT;
		EXEC @lock_result = sys.sp_getapplock
			@Resource = ?,
			@LockMode = 'Exclusive',
			@LockOwner = 'Transaction',
			@LockTimeout = 15000;
		IF @lock_result < 0
			THROW 50001, 'Khong the khoa snapshot dashboard theo tuan.', 1;
	`, lockResource).Error; err != nil {
		return err
	}
	if err := tx.Exec(`
		MERGE dbo.gore_tex_dashboard_pass_snapshots WITH (HOLDLOCK) AS target
		USING (
			SELECT
				CAST(? AS DATE) AS period_start,
				? AS iso_year,
				? AS iso_week,
				style_name,
				SUM(pass_count) AS pass_count,
				SUM(fail_count) AS fail_count,
				SUM(total_count) AS total_count
			FROM dbo.gore_tex_dashboard_pass_snapshots
			WHERE period_type = 'D' AND period_start >= ? AND period_start < ?
			GROUP BY style_name
		) AS source
		ON target.period_type = 'W'
			AND target.period_start = source.period_start
			AND target.style_name = source.style_name
		WHEN MATCHED THEN UPDATE SET
			target.iso_year = source.iso_year,
			target.iso_week = source.iso_week,
			target.pass_count = source.pass_count,
			target.fail_count = source.fail_count,
			target.total_count = source.total_count,
			target.updated_at = ?
		WHEN NOT MATCHED THEN INSERT
			(period_type, period_start, iso_year, iso_week, style_name, pass_count, fail_count, total_count, updated_at)
		VALUES
			('W', source.period_start, source.iso_year, source.iso_week, source.style_name, source.pass_count, source.fail_count, source.total_count, ?);
	`, weekStart, year, week, weekStart, weekEnd, now, now).Error; err != nil {
		return fmt.Errorf("cap nhat snapshot PASS theo tuan: %w", err)
	}
	if err := tx.Exec(`
		DELETE target
		FROM dbo.gore_tex_dashboard_pass_snapshots AS target
		WHERE target.period_type = 'W'
			AND target.period_start = ?
			AND NOT EXISTS (
				SELECT 1
				FROM dbo.gore_tex_dashboard_pass_snapshots AS daily
				WHERE daily.period_type = 'D'
					AND daily.period_start >= ? AND daily.period_start < ?
					AND daily.style_name = target.style_name
			);
	`, weekStart, weekStart, weekEnd).Error; err != nil {
		return fmt.Errorf("xoa snapshot PASS tuan khong con du lieu: %w", err)
	}
	if err := tx.Exec(`
		MERGE dbo.gore_tex_dashboard_defect_snapshots WITH (HOLDLOCK) AS target
		USING (
			SELECT
				CAST(? AS DATE) AS period_start,
				? AS iso_year,
				? AS iso_week,
				defect_code,
				MAX(defect_name) AS defect_name,
				SUM(defect_count) AS defect_count
			FROM dbo.gore_tex_dashboard_defect_snapshots
			WHERE period_type = 'D' AND period_start >= ? AND period_start < ?
			GROUP BY defect_code
		) AS source
		ON target.period_type = 'W'
			AND target.period_start = source.period_start
			AND target.source_record_key = ''
			AND target.defect_code = source.defect_code
		WHEN MATCHED THEN UPDATE SET
			target.iso_year = source.iso_year,
			target.iso_week = source.iso_week,
			target.defect_name = source.defect_name,
			target.defect_count = source.defect_count,
			target.updated_at = ?
		WHEN NOT MATCHED THEN INSERT
			(period_type, period_start, iso_year, iso_week, source_record_key, defect_code, defect_name, defect_count, updated_at)
		VALUES
			('W', source.period_start, source.iso_year, source.iso_week, '', source.defect_code, source.defect_name, source.defect_count, ?);
	`, weekStart, year, week, weekStart, weekEnd, now, now).Error; err != nil {
		return fmt.Errorf("cap nhat snapshot loi theo tuan: %w", err)
	}
	if err := tx.Exec(`
		DELETE target
		FROM dbo.gore_tex_dashboard_defect_snapshots AS target
		WHERE target.period_type = 'W'
			AND target.period_start = ?
			AND NOT EXISTS (
				SELECT 1
				FROM dbo.gore_tex_dashboard_defect_snapshots AS daily
				WHERE daily.period_type = 'D'
					AND daily.period_start >= ? AND daily.period_start < ?
					AND daily.defect_code = target.defect_code
			);
	`, weekStart, weekStart, weekEnd).Error; err != nil {
		return fmt.Errorf("xoa snapshot loi tuan khong con du lieu: %w", err)
	}
	return nil
}

func goreTexDefectGroupCount(data goreTexDashboardWaterproofData, rowIDs []string) float64 {
	total := 0.0
	for _, rowID := range rowIDs {
		if value, ok := goreTexDashboardNumber(data.Totals[rowID]); ok {
			total += value
			continue
		}
		for key, rawValue := range data.Counts {
			if len(key) > len(rowID) && key[:len(rowID)+1] == rowID+":" {
				if value, ok := goreTexDashboardNumber(rawValue); ok {
					total += value
				}
			}
		}
	}
	return total
}

func goreTexDefectCode(index int) string {
	codes := []string{"TOE", "HEEL", "MEDIAL", "LATERAL", "MATERIAL", "ATTACHING", "WRINKLED", "ZIGZAG", "BONDING"}
	return codes[index]
}

func normalizeGoreTexStyle(value string) string  { return strings.TrimSpace(value) }
func normalizeGoreTexResult(value string) string { return strings.ToUpper(strings.TrimSpace(value)) }

type goreTexPassSnapshotRow struct {
	PeriodStart time.Time `gorm:"column:period_start"`
	StyleName   string    `gorm:"column:style_name"`
	PassCount   int       `gorm:"column:pass_count"`
	TotalCount  int       `gorm:"column:total_count"`
}

type goreTexDefectSnapshotRow struct {
	DefectName  string  `gorm:"column:defect_name"`
	DefectCount float64 `gorm:"column:defect_count"`
}

func backfillGoreTexDashboardSnapshots(db *gorm.DB, centrifugalRows []goreTexCentrifugalRow, waterproofRows []goreTexWaterproofRow) error {
	now := goreTexNow()
	return db.Transaction(func(tx *gorm.DB) error {
		for _, row := range centrifugalRows {
			if err := syncGoreTexPassDailySnapshot(tx, row.InspectionDate, row.DataJSON, now); err != nil {
				return err
			}
		}
		for _, row := range waterproofRows {
			if err := syncGoreTexDefectDailySnapshot(tx, row.InspectionDate, row.Line+" / "+row.StyleName, row.DataJSON, now); err != nil {
				return err
			}
		}
		return nil
	})
}

func readGoreTexWeeklyDashboardSnapshots(db *gorm.DB, year, week int, weekStart, previousStart time.Time) (*types.GoreTexWeeklyDashboard, error) {
	selectedItems, err := readGoreTexWeeklyPassItems(db, weekStart)
	if err != nil {
		return nil, err
	}
	previousItems, err := readGoreTexWeeklyPassItems(db, previousStart)
	if err != nil {
		return nil, err
	}
	selectedTrend, selectedRRdy, err := readGoreTexDailyTrends(db, weekStart)
	if err != nil {
		return nil, err
	}
	pareto, err := readGoreTexWeeklyPareto(db, weekStart)
	if err != nil {
		return nil, err
	}
	previousYear, previousWeek := previousStart.ISOWeek()
	return &types.GoreTexWeeklyDashboard{
		Year: year, Week: week,
		FromDate: weekStart.Format("2006-01-02"), ToDate: weekStart.AddDate(0, 0, 6).Format("2006-01-02"),
		SuterByItems: selectedItems,
		SuterTrend:   selectedTrend,
		SuterComparison: []types.GoreTexDashboardComparison{
			{Year: previousYear, Week: previousWeek, Label: "Tuần " + strconv.Itoa(previousWeek), Items: previousItems},
			{Year: year, Week: week, Label: "Tuần " + strconv.Itoa(week), Items: selectedItems},
		},
		RRdyTrend: selectedRRdy, VisualizationResults: pareto,
	}, nil
}

func readGoreTexWeeklyPassItems(db *gorm.DB, weekStart time.Time) ([]types.GoreTexDashboardRateItem, error) {
	var rows []goreTexPassSnapshotRow
	query := db.Raw(`
		SELECT style_name, SUM(pass_count) AS pass_count, SUM(total_count) AS total_count
		FROM dbo.gore_tex_dashboard_pass_snapshots
		WHERE period_type = 'D' AND period_start >= ? AND period_start < ?
		GROUP BY style_name
	`, weekStart, weekStart.AddDate(0, 0, 7)).Scan(&rows)
	if query.Error != nil {
		return nil, query.Error
	}
	items := make([]types.GoreTexDashboardRateItem, 0, len(rows))
	for _, row := range rows {
		value := &goreTexRateAccumulator{Pass: row.PassCount, Total: row.TotalCount}
		items = append(items, types.GoreTexDashboardRateItem{Label: row.StyleName, Pass: row.PassCount, Total: row.TotalCount, Rate: goreTexRate(value)})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Rate == items[j].Rate {
			return items[i].Label < items[j].Label
		}
		return items[i].Rate > items[j].Rate
	})
	return items, nil
}

func readGoreTexDailyTrends(db *gorm.DB, weekStart time.Time) ([]types.GoreTexDashboardTrendPoint, []types.GoreTexDashboardTrendPoint, error) {
	var rows []goreTexPassSnapshotRow
	query := db.Raw(`
		SELECT period_start, style_name, pass_count, total_count
		FROM dbo.gore_tex_dashboard_pass_snapshots
		WHERE period_type = 'D' AND period_start >= ? AND period_start < ?
		ORDER BY period_start
	`, weekStart, weekStart.AddDate(0, 0, 7)).Scan(&rows)
	if query.Error != nil {
		return nil, nil, query.Error
	}
	all := make(map[string]*goreTexRateAccumulator)
	rRdy := make(map[string]*goreTexRateAccumulator)
	for _, row := range rows {
		key := row.PeriodStart.Format("2006-01-02")
		addGoreTexSnapshotCounts(all, key, row.PassCount, row.TotalCount)
		if isGoreTexRRdyStyle(row.StyleName) {
			addGoreTexSnapshotCounts(rRdy, key, row.PassCount, row.TotalCount)
		}
	}
	return goreTexTrendPoints(all), goreTexTrendPoints(rRdy), nil
}

func addGoreTexSnapshotCounts(target map[string]*goreTexRateAccumulator, key string, pass, total int) {
	item := target[key]
	if item == nil {
		item = &goreTexRateAccumulator{}
		target[key] = item
	}
	item.Pass += pass
	item.Total += total
}

func readGoreTexWeeklyPareto(db *gorm.DB, weekStart time.Time) ([]types.GoreTexDashboardParetoItem, error) {
	var rows []goreTexDefectSnapshotRow
	query := db.Raw(`
		SELECT MAX(defect_name) AS defect_name, SUM(defect_count) AS defect_count
		FROM dbo.gore_tex_dashboard_defect_snapshots
		WHERE period_type = 'D' AND period_start >= ? AND period_start < ?
		GROUP BY defect_code
		ORDER BY SUM(defect_count) DESC, MAX(defect_name)
	`, weekStart, weekStart.AddDate(0, 0, 7)).Scan(&rows)
	if query.Error != nil {
		return nil, query.Error
	}
	items := make([]types.GoreTexDashboardParetoItem, 0, len(rows))
	total := 0.0
	for _, row := range rows {
		total += row.DefectCount
	}
	running := 0.0
	for _, row := range rows {
		running += row.DefectCount
		items = append(items, types.GoreTexDashboardParetoItem{Label: row.DefectName, Count: row.DefectCount, Cumulative: math.Round(running*1000/total) / 10})
	}
	return items, nil
}
