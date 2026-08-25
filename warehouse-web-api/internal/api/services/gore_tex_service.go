package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"gorm.io/gorm"
)

var (
	ErrGoreTexInvalidForm  = errors.New("du lieu bieu mau khong hop le")
	ErrGoreTexFormNotFound = errors.New("khong tim thay bieu mau")
)

type GoreTexFormService struct {
	*BaseService
}

var GoreTexForms = &GoreTexFormService{}

type goreTexWaterproofRow struct {
	RecordKey      string    `gorm:"column:record_key"`
	Line           string    `gorm:"column:line"`
	StyleName      string    `gorm:"column:style_name"`
	InspectionDate time.Time `gorm:"column:inspection_date"`
	DataJSON       string    `gorm:"column:data_json"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

type goreTexCentrifugalRow struct {
	InspectionDate time.Time `gorm:"column:inspection_date"`
	DataJSON       string    `gorm:"column:data_json"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

type goreTexAnalysisRow struct {
	ID        uint64    `gorm:"column:id"`
	DataJSON  string    `gorm:"column:data_json"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type goreTexCentrifugalData struct {
	Entries []struct {
		Result      string                 `json:"result"`
		IssueValues map[string]interface{} `json:"issueValues"`
		OtherIssue  interface{}            `json:"otherIssue"`
	} `json:"entries"`
}

type goreTexWaterproofNumericData struct {
	Counts        map[string]interface{}            `json:"counts"`
	Totals        map[string]interface{}            `json:"totals"`
	Rates         map[string]interface{}            `json:"rates"`
	SummaryCounts map[string]map[string]interface{} `json:"summaryCounts"`
	SummaryTotals map[string]interface{}            `json:"summaryTotals"`
	SummaryEdges  map[string]map[string]interface{} `json:"summaryEdges"`
}

func isGoreTexNonNegativeNumber(value interface{}) bool {
	var number float64
	switch typed := value.(type) {
	case float64:
		number = typed
	case string:
		normalized := strings.ReplaceAll(strings.TrimSpace(typed), ",", ".")
		if normalized == "" {
			return false
		}
		parsed, err := strconv.ParseFloat(normalized, 64)
		if err != nil {
			return false
		}
		number = parsed
	default:
		return false
	}
	return !math.IsNaN(number) && !math.IsInf(number, 0) && number >= 0
}

func validateGoreTexNumericMap(values map[string]interface{}) bool {
	if len(values) == 0 {
		return false
	}
	for _, value := range values {
		if !isGoreTexNonNegativeNumber(value) {
			return false
		}
	}
	return true
}

func validateGoreTexNestedNumericMap(groups map[string]map[string]interface{}) bool {
	if len(groups) == 0 {
		return false
	}
	for _, values := range groups {
		if !validateGoreTexNumericMap(values) {
			return false
		}
	}
	return true
}

func validateGoreTexWaterproofNumbers(data json.RawMessage) bool {
	var form goreTexWaterproofNumericData
	if err := json.Unmarshal(data, &form); err != nil {
		return false
	}
	return validateGoreTexNumericMap(form.Counts) &&
		validateGoreTexNumericMap(form.Totals) &&
		validateGoreTexNumericMap(form.Rates) &&
		validateGoreTexNestedNumericMap(form.SummaryCounts) &&
		validateGoreTexNumericMap(form.SummaryTotals) &&
		validateGoreTexNestedNumericMap(form.SummaryEdges)
}

func validateGoreTexCentrifugalResults(data json.RawMessage) bool {
	var form goreTexCentrifugalData
	if err := json.Unmarshal(data, &form); err != nil || len(form.Entries) == 0 {
		return false
	}

	for _, entry := range form.Entries {
		if (entry.Result != "PASS" && entry.Result != "FAIL") ||
			!validateGoreTexNumericMap(entry.IssueValues) ||
			!isGoreTexNonNegativeNumber(entry.OtherIssue) {
			return false
		}
	}
	return true
}

func normalizeGoreTexJSONStyles(data json.RawMessage, collectionKey, styleKey string) (json.RawMessage, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, ErrGoreTexInvalidForm
	}
	if rows, ok := payload[collectionKey].([]interface{}); ok {
		for _, rawRow := range rows {
			row, ok := rawRow.(map[string]interface{})
			if !ok {
				continue
			}
			if style, ok := row[styleKey].(string); ok {
				row[styleKey] = normalizeGoreTexStyle(style)
			}
		}
	}
	normalized, err := json.Marshal(payload)
	if err != nil {
		return nil, ErrGoreTexInvalidForm
	}
	return normalized, nil
}

func (s *GoreTexFormService) SaveWaterproof(params request.GoreTexWaterproofSubmitRequest) (*types.GoreTexSubmitResult, error) {
	line := strings.TrimSpace(params.Line)
	styleName := normalizeGoreTexStyle(params.StyleName)
	if line == "" || styleName == "" || !json.Valid(params.Data) {
		return nil, ErrGoreTexInvalidForm
	}
	params.Line = line
	params.StyleName = styleName

	inspectionDate, err := parseGoreTexDate(params.InspectionDate)
	if err != nil {
		return nil, err
	}
	return saveNormalizedGoreTexWaterproof(params, inspectionDate)
}

func (s *GoreTexFormService) SaveCentrifugal(params request.GoreTexCentrifugalSubmitRequest) (*types.GoreTexSubmitResult, error) {
	normalizedData, err := normalizeGoreTexJSONStyles(params.Data, "entries", "style")
	if err != nil {
		return nil, err
	}
	params.Data = normalizedData
	if !validateGoreTexCentrifugalResults(params.Data) {
		return nil, ErrGoreTexInvalidForm
	}

	inspectionDate, err := parseGoreTexDate(params.InspectionDate)
	if err != nil {
		return nil, err
	}

	db, closeDB, err := openGoreTexDB()
	if err != nil {
		return nil, err
	}
	defer closeDB()

	now := goreTexNow()
	created := false
	err = db.Transaction(func(tx *gorm.DB) error {
		var existing struct {
			DataJSON string `gorm:"column:data_json"`
		}
		found := tx.Raw(`
			SELECT data_json
			FROM dbo.gore_tex_centrifugal_forms
			WHERE inspection_date = ?
		`, inspectionDate).Scan(&existing)
		if found.Error != nil {
			return found.Error
		}

		if found.RowsAffected == 0 {
			if params.IsEdit {
				return ErrGoreTexFormNotFound
			}
			created = true
			if insertErr := tx.Exec(`
				INSERT INTO dbo.gore_tex_centrifugal_forms
					(inspection_date, data_json, created_at, updated_at)
				VALUES (?, ?, ?, ?)
			`, inspectionDate, string(params.Data), now, now).Error; insertErr != nil {
				return insertErr
			}
			return nil
		}

		dataJSON, mergeErr := appendGoreTexEditHistory(params.Data, existing.DataJSON, now)
		if mergeErr != nil {
			return mergeErr
		}
		if updateErr := tx.Exec(`
			UPDATE dbo.gore_tex_centrifugal_forms
			SET data_json = ?, updated_at = ?
			WHERE inspection_date = ?
		`, dataJSON, now, inspectionDate).Error; updateErr != nil {
			return updateErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &types.GoreTexSubmitResult{
		FormType:  "centrifugal",
		RecordKey: params.InspectionDate,
		Created:   created,
	}, nil
}

func (s *GoreTexFormService) SaveAnalysis(params request.GoreTexAnalysisSubmitRequest) (*types.GoreTexSubmitResult, error) {
	normalizedData, err := normalizeGoreTexJSONStyles(params.Data, "records", "styleName")
	if err != nil {
		return nil, err
	}
	params.Data = normalizedData
	if params.IsEdit && params.AnalysisID == 0 {
		return nil, ErrGoreTexInvalidForm
	}

	now := goreTexNow()
	db, closeDB, err := openGoreTexDB()
	if err != nil {
		return nil, err
	}
	defer closeDB()

	recordID := params.AnalysisID
	created := recordID == 0
	if created {
		var inserted struct {
			ID uint64 `gorm:"column:id"`
		}
		result := db.Raw(`
			INSERT INTO dbo.gore_tex_analysis_forms
				(data_json, created_at, updated_at)
			OUTPUT INSERTED.id
			VALUES (?, ?, ?)
		`, string(params.Data), now, now).Scan(&inserted)
		if result.Error != nil {
			return nil, result.Error
		}
		if inserted.ID == 0 {
			return nil, errors.New("khong lay duoc ID bieu mau phan tich")
		}
		recordID = inserted.ID
	} else {
		err = db.Transaction(func(tx *gorm.DB) error {
			var existing struct {
				DataJSON string `gorm:"column:data_json"`
			}
			found := tx.Raw(`
				SELECT data_json
				FROM dbo.gore_tex_analysis_forms
				WHERE id = ?
			`, recordID).Scan(&existing)
			if found.Error != nil {
				return found.Error
			}
			if found.RowsAffected == 0 {
				return ErrGoreTexFormNotFound
			}

			dataJSON, mergeErr := appendGoreTexEditHistory(params.Data, existing.DataJSON, now)
			if mergeErr != nil {
				return mergeErr
			}
			return tx.Exec(`
				UPDATE dbo.gore_tex_analysis_forms
				SET data_json = ?, updated_at = ?
				WHERE id = ?
			`, dataJSON, now, recordID).Error
		})
		if err != nil {
			return nil, err
		}
	}

	return &types.GoreTexSubmitResult{
		FormType:  "analysis",
		RecordKey: strconv.FormatUint(recordID, 10),
		Created:   created,
	}, nil
}

func (s *GoreTexFormService) List() ([]types.GoreTexFormListItem, error) {
	db, closeDB, err := openGoreTexDB()
	if err != nil {
		return nil, err
	}
	defer closeDB()

	result := make([]types.GoreTexFormListItem, 0)

	var waterproof []goreTexWaterproofRow
	query := db.Raw(`
		SELECT
			DepNo AS [line], StyleName AS style_name, InspectionDate AS inspection_date,
			CreatedAt AS created_at, COALESCE(UpdatedAt, CreatedAt) AS updated_at
		FROM dbo.GoreTex_Waterproof_Report
		ORDER BY COALESCE(UpdatedAt, CreatedAt) DESC
	`).Scan(&waterproof)
	if query.Error != nil {
		return nil, query.Error
	}
	for _, form := range waterproof {
		styleName := normalizeGoreTexStyle(form.StyleName)
		result = append(result, types.GoreTexFormListItem{
			FormType:       "waterproof",
			Title:          "Kiểm tra chất lượng giày chống thấm",
			RecordKey:      form.Line + " / " + styleName,
			Line:           form.Line,
			StyleName:      styleName,
			InspectionDate: formatGoreTexDate(form.InspectionDate),
			CreatedAt:      formatGoreTexDateTime(form.CreatedAt),
			UpdatedAt:      formatGoreTexDateTime(form.UpdatedAt),
		})
	}

	var centrifugal []goreTexCentrifugalRow
	query = db.Raw(`
		SELECT inspection_date, created_at, updated_at
		FROM dbo.gore_tex_centrifugal_forms
		ORDER BY updated_at DESC
	`).Scan(&centrifugal)
	if query.Error != nil {
		return nil, query.Error
	}
	for _, form := range centrifugal {
		date := formatGoreTexDate(form.InspectionDate)
		result = append(result, types.GoreTexFormListItem{
			FormType:       "centrifugal",
			Title:          "Giày thành phẩm thử nghiệm li tâm",
			RecordKey:      date,
			InspectionDate: date,
			CreatedAt:      formatGoreTexDateTime(form.CreatedAt),
			UpdatedAt:      formatGoreTexDateTime(form.UpdatedAt),
		})
	}

	var analysis []goreTexAnalysisRow
	query = db.Raw(`
		SELECT id, data_json, created_at, updated_at
		FROM dbo.gore_tex_analysis_forms
		ORDER BY updated_at DESC
	`).Scan(&analysis)
	if query.Error != nil {
		return nil, query.Error
	}
	for _, form := range analysis {
		testDates, improvementDates := extractGoreTexAnalysisDates(form.DataJSON)
		result = append(result, types.GoreTexFormListItem{
			FormType:         "analysis",
			Title:            "Phân tích nguyên nhân và cải thiện",
			RecordKey:        strconv.FormatUint(form.ID, 10),
			AnalysisID:       form.ID,
			TestDates:        testDates,
			ImprovementDates: improvementDates,
			CreatedAt:        formatGoreTexDateTime(form.CreatedAt),
			UpdatedAt:        formatGoreTexDateTime(form.UpdatedAt),
		})
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].UpdatedAt > result[j].UpdatedAt
	})
	return result, nil
}

func (s *GoreTexFormService) GetWaterproof(line, styleName string) (*types.GoreTexFormDetail, error) {
	line = strings.TrimSpace(line)
	styleName = normalizeGoreTexStyle(styleName)
	if line == "" || styleName == "" {
		return nil, ErrGoreTexInvalidForm
	}

	db, closeDB, err := openGoreTexDB()
	if err != nil {
		return nil, err
	}
	defer closeDB()

	return getNormalizedGoreTexWaterproof(db, line, styleName, true)
}

func (s *GoreTexFormService) GetCentrifugal(dateValue string) (*types.GoreTexFormDetail, error) {
	inspectionDate, err := parseGoreTexDate(dateValue)
	if err != nil {
		return nil, err
	}

	var form goreTexCentrifugalRow
	db, closeDB, err := openGoreTexDB()
	if err != nil {
		return nil, err
	}
	defer closeDB()

	query := db.Raw(`
		SELECT inspection_date, data_json, created_at, updated_at
		FROM dbo.gore_tex_centrifugal_forms
		WHERE inspection_date = ?
	`, inspectionDate).Scan(&form)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return nil, ErrGoreTexFormNotFound
	}

	return &types.GoreTexFormDetail{
		FormType:       "centrifugal",
		InspectionDate: formatGoreTexDate(form.InspectionDate),
		Data:           json.RawMessage(form.DataJSON),
		CreatedAt:      formatGoreTexDateTime(form.CreatedAt),
		UpdatedAt:      formatGoreTexDateTime(form.UpdatedAt),
	}, nil
}

func (s *GoreTexFormService) GetAnalysis(id uint64) (*types.GoreTexFormDetail, error) {
	if id == 0 {
		return nil, ErrGoreTexInvalidForm
	}

	var form goreTexAnalysisRow
	db, closeDB, err := openGoreTexDB()
	if err != nil {
		return nil, err
	}
	defer closeDB()

	query := db.Raw(`
		SELECT id, data_json, created_at, updated_at
		FROM dbo.gore_tex_analysis_forms
		WHERE id = ?
	`, id).Scan(&form)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return nil, ErrGoreTexFormNotFound
	}

	return &types.GoreTexFormDetail{
		FormType:   "analysis",
		AnalysisID: form.ID,
		Data:       json.RawMessage(form.DataJSON),
		CreatedAt:  formatGoreTexDateTime(form.CreatedAt),
		UpdatedAt:  formatGoreTexDateTime(form.UpdatedAt),
	}, nil
}

func parseGoreTexDate(value string) (time.Time, error) {
	result, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: ngay phai co dinh dang YYYY-MM-DD", ErrGoreTexInvalidForm)
	}
	return result, nil
}

func formatGoreTexDate(value time.Time) string {
	return value.Format("2006-01-02")
}

func formatGoreTexDateTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format("2006-01-02T15:04:05")
}

func goreTexNow() time.Time {
	return time.Now().In(time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60))
}

func appendGoreTexEditHistory(newData json.RawMessage, previousData string, editedAt time.Time) (string, error) {
	var next map[string]interface{}
	if err := json.Unmarshal(newData, &next); err != nil {
		return "", ErrGoreTexInvalidForm
	}

	var previous map[string]interface{}
	if err := json.Unmarshal([]byte(previousData), &previous); err != nil {
		return "", fmt.Errorf("du lieu bieu mau cu khong hop le: %w", err)
	}

	history := make([]interface{}, 0)
	if existingHistory, ok := previous["_editHistory"].([]interface{}); ok {
		history = append(history, existingHistory...)
	}
	delete(previous, "_editHistory")
	history = append(history, map[string]interface{}{
		"editedAt": formatGoreTexDateTime(editedAt),
		"data":     previous,
	})
	next["_editHistory"] = history

	merged, err := json.Marshal(next)
	if err != nil {
		return "", err
	}
	return string(merged), nil
}

func extractGoreTexAnalysisDates(dataJSON string) ([]string, []string) {
	var data struct {
		Records []struct {
			TestDate        string `json:"testDate"`
			ImprovementDate string `json:"improvementDate"`
		} `json:"records"`
	}
	if err := json.Unmarshal([]byte(dataJSON), &data); err != nil {
		return nil, nil
	}

	testDates := make([]string, 0)
	improvementDates := make([]string, 0)
	seenTestDates := make(map[string]struct{})
	seenImprovementDates := make(map[string]struct{})
	for _, record := range data.Records {
		testDate := strings.TrimSpace(record.TestDate)
		if testDate != "" {
			if _, exists := seenTestDates[testDate]; !exists {
				seenTestDates[testDate] = struct{}{}
				testDates = append(testDates, testDate)
			}
		}

		improvementDate := strings.TrimSpace(record.ImprovementDate)
		if improvementDate != "" {
			if _, exists := seenImprovementDates[improvementDate]; !exists {
				seenImprovementDates[improvementDate] = struct{}{}
				improvementDates = append(improvementDates, improvementDate)
			}
		}
	}
	return testDates, improvementDates
}

func openGoreTexDB() (*gorm.DB, func(), error) {
	db, err := database.LYS_WEB_Connection()
	if err != nil {
		return nil, nil, fmt.Errorf("LYS_WEB connection error: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	closeDB := func() {
		_ = sqlDB.Close()
	}
	return db, closeDB, nil
}

// Dashboard aggregation.
const goreTexDashboardTarget = 99.0

type goreTexDashboardCentrifugalData struct {
	Entries []struct {
		Style  string `json:"style"`
		Result string `json:"result"`
	} `json:"entries"`
}

type goreTexDashboardWaterproofData struct {
	Counts        map[string]interface{} `json:"counts"`
	Totals        map[string]interface{} `json:"totals"`
	SummaryTotals map[string]interface{} `json:"summaryTotals"`
}

type goreTexRateAccumulator struct {
	Pass  int
	Total int
}

var goreTexDefectGroups = []struct {
	Label string
	Rows  []string
}{
	{Label: "Mũi thấm nước", Rows: []string{"toe-left", "toe-right"}},
	{Label: "Gót thấm nước", Rows: []string{"heel-left", "heel-right"}},
	{Label: "Hong trong thấm nước", Rows: []string{"medial-left", "medial-right"}},
	{Label: "Hong ngoài thấm nước", Rows: []string{"lateral-left", "lateral-right"}},
	{Label: "Vật tư không đạt", Rows: []string{"material-left", "material-right"}},
	{Label: "Dán đế lệch", Rows: []string{"attaching-left", "attaching-right"}},
	{Label: "Ép đế nhăn", Rows: []string{"wrinkled-left", "wrinkled-right"}},
	{Label: "Zíc zắc hở", Rows: []string{"zigzag-quarter-left", "zigzag-quarter-right", "zigzag-sole-left", "zigzag-sole-right"}},
	{Label: "Hở keo đế", Rows: []string{"bonding-toe-left", "bonding-toe-right", "bonding-heel-left", "bonding-heel-right"}},
}

func (s *GoreTexFormService) DashboardPeriod(period, value, selectedStyle, mode string) (*types.GoreTexDashboardPeriod, error) {
	period = strings.ToLower(strings.TrimSpace(period))
	value = strings.TrimSpace(value)
	selectedStyle = normalizeGoreTexStyle(selectedStyle)
	mode = strings.ToLower(strings.TrimSpace(mode))
	var start, end time.Time
	var err error
	if mode == "compare" {
		start, end, err = goreTexDashboardComparisonRange(period, value)
	} else {
		start, end, err = goreTexDashboardSelectedRange(period, value)
	}
	if err != nil {
		return nil, err
	}
	selectedStart, selectedEnd := start, end

	db, closeDB, err := openGoreTexDB()
	if err != nil {
		return nil, err
	}
	defer closeDB()

	waterproofRows, err := loadNormalizedGoreTexWaterproofDashboardRows(db, start, end)
	if err != nil {
		return nil, err
	}

	byPeriod := make(map[string]*goreTexRateAccumulator)
	periodLabels := make(map[string]string)
	selectedStyleRows := make([]goreTexWaterproofRow, 0)
	styleSet := make(map[string]struct{})
	for _, row := range waterproofRows {
		if !row.InspectionDate.Before(selectedStart) && row.InspectionDate.Before(selectedEnd) {
			style := normalizeGoreTexStyle(row.StyleName)
			if style != "" {
				styleSet[style] = struct{}{}
			}
			if selectedStyle == "" || style == selectedStyle {
				selectedStyleRows = append(selectedStyleRows, row)
			}
		}
		var data goreTexDashboardWaterproofData
		if json.Unmarshal([]byte(row.DataJSON), &data) != nil {
			continue
		}
		pass, total, ok := goreTexWaterproofPassCounts(data)
		if !ok {
			continue
		}
		bucket := goreTexDashboardBucketStart(row.InspectionDate, period)
		key := bucket.Format("2006-01-02")
		periodLabels[key] = goreTexDashboardBucketLabel(bucket, period)
		addGoreTexSnapshotCounts(byPeriod, key, pass, total)
	}
	styles := make([]string, 0, len(styleSet))
	for style := range styleSet {
		styles = append(styles, style)
	}
	sort.Strings(styles)

	return &types.GoreTexDashboardPeriod{
		Period: period, Value: value,
		FromDate: start.Format("2006-01-02"), ToDate: end.AddDate(0, 0, -1).Format("2006-01-02"),
		Styles: styles, SelectedStyle: selectedStyle,
		SuterTrend: goreTexTrendPointsWithLabels(byPeriod, periodLabels), VisualizationResults: goreTexParetoItems(selectedStyleRows),
	}, nil
}

func goreTexDashboardSelectedRange(period, value string) (time.Time, time.Time, error) {
	switch period {
	case "day":
		start, err := time.ParseInLocation("2006-01-02", value, time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, ErrGoreTexInvalidForm
		}
		return start, start.AddDate(0, 0, 1), nil
	case "week":
		var year, week int
		if _, err := fmt.Sscanf(value, "%d-W%d", &year, &week); err != nil || fmt.Sprintf("%04d-W%02d", year, week) != value {
			return time.Time{}, time.Time{}, ErrGoreTexInvalidForm
		}
		start, err := goreTexISOWeekStart(year, week)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return start, start.AddDate(0, 0, 7), nil
	case "month":
		start, err := time.ParseInLocation("2006-01", value, time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, ErrGoreTexInvalidForm
		}
		return start, start.AddDate(0, 1, 0), nil
	default:
		return time.Time{}, time.Time{}, ErrGoreTexInvalidForm
	}
}

func goreTexDashboardComparisonRange(period, value string) (time.Time, time.Time, error) {
	switch period {
	case "day":
		var year, week int
		if _, err := fmt.Sscanf(value, "%d-W%d", &year, &week); err != nil || fmt.Sprintf("%04d-W%02d", year, week) != value {
			return time.Time{}, time.Time{}, ErrGoreTexInvalidForm
		}
		start, err := goreTexISOWeekStart(year, week)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return start, start.AddDate(0, 0, 7), nil
	case "week":
		start, err := time.ParseInLocation("2006-01", value, time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, ErrGoreTexInvalidForm
		}
		return start, start.AddDate(0, 1, 0), nil
	case "month":
		year, err := strconv.Atoi(value)
		if err != nil || year < 2000 || year > 2200 || strconv.Itoa(year) != value {
			return time.Time{}, time.Time{}, ErrGoreTexInvalidForm
		}
		yearStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
		return yearStart, yearStart.AddDate(1, 0, 0), nil
	default:
		return time.Time{}, time.Time{}, ErrGoreTexInvalidForm
	}
}

func goreTexDashboardBucketStart(date time.Time, period string) time.Time {
	switch period {
	case "week":
		daysSinceMonday := (int(date.Weekday()) + 6) % 7
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()).AddDate(0, 0, -daysSinceMonday)
	case "month":
		return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	default:
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	}
}

func goreTexDashboardBucketLabel(date time.Time, period string) string {
	switch period {
	case "week":
		year, week := date.ISOWeek()
		return fmt.Sprintf("W%02d/%d", week, year)
	case "month":
		return date.Format("01/2006")
	default:
		return date.Format("02/01")
	}
}

func (s *GoreTexFormService) WeeklyDashboard(year, week int) (*types.GoreTexWeeklyDashboard, error) {
	weekStart, err := goreTexISOWeekStart(year, week)
	if err != nil {
		return nil, err
	}
	previousStart := weekStart.AddDate(0, 0, -7)
	weekEndExclusive := weekStart.AddDate(0, 0, 7)
	previousYear, previousWeek := previousStart.ISOWeek()

	db, closeDB, err := openGoreTexDB()
	if err != nil {
		return nil, err
	}
	defer closeDB()

	waterproofRows, err := loadNormalizedGoreTexWaterproofDashboardRows(db, previousStart, weekEndExclusive)
	if err != nil {
		return nil, err
	}

	selectedByStyle := make(map[string]*goreTexRateAccumulator)
	previousByStyle := make(map[string]*goreTexRateAccumulator)
	selectedByDate := make(map[string]*goreTexRateAccumulator)

	selectedWaterproofRows := make([]goreTexWaterproofRow, 0, len(waterproofRows))
	for _, row := range waterproofRows {
		isSelected := !row.InspectionDate.Before(weekStart) && row.InspectionDate.Before(weekEndExclusive)
		if isSelected {
			selectedWaterproofRows = append(selectedWaterproofRows, row)
		}
		var data goreTexDashboardWaterproofData
		if json.Unmarshal([]byte(row.DataJSON), &data) != nil {
			continue
		}
		pass, total, ok := goreTexWaterproofPassCounts(data)
		style := normalizeGoreTexStyle(row.StyleName)
		if !ok || style == "" {
			continue
		}
		styleTarget := previousByStyle
		if isSelected {
			styleTarget = selectedByStyle
			addGoreTexSnapshotCounts(selectedByDate, row.InspectionDate.Format("2006-01-02"), pass, total)
		}
		addGoreTexSnapshotCounts(styleTarget, style, pass, total)
	}

	return &types.GoreTexWeeklyDashboard{
		Year:         year,
		Week:         week,
		FromDate:     weekStart.Format("2006-01-02"),
		ToDate:       weekEndExclusive.AddDate(0, 0, -1).Format("2006-01-02"),
		SuterByItems: goreTexRateItems(selectedByStyle),
		SuterTrend:   goreTexTrendPoints(selectedByDate),
		SuterComparison: []types.GoreTexDashboardComparison{
			{Year: previousYear, Week: previousWeek, Label: fmt.Sprintf("Tuần %d", previousWeek), Items: goreTexRateItems(previousByStyle)},
			{Year: year, Week: week, Label: fmt.Sprintf("Tuần %d", week), Items: goreTexRateItems(selectedByStyle)},
		},
		RRdyTrend:            []types.GoreTexDashboardTrendPoint{},
		VisualizationResults: goreTexParetoItems(selectedWaterproofRows),
	}, nil
}

func goreTexWaterproofPassCounts(data goreTexDashboardWaterproofData) (int, int, bool) {
	inspectionValue, inspectionOK := goreTexDashboardNumber(data.SummaryTotals["inspection"])
	defectValue, defectOK := goreTexDashboardNumber(data.SummaryTotals["defects"])
	if !inspectionOK || inspectionValue < 0 {
		return 0, 0, false
	}
	if !defectOK || defectValue < 0 {
		defectValue = 0
	}
	total := int(math.Round(inspectionValue))
	defects := int(math.Round(defectValue))
	pass := total - defects
	if pass < 0 {
		pass = 0
	}
	return pass, total, true
}

func goreTexISOWeekStart(year, week int) (time.Time, error) {
	if year < 2000 || year > 2200 || week < 1 || week > 53 {
		return time.Time{}, ErrGoreTexInvalidForm
	}
	jan4 := time.Date(year, time.January, 4, 0, 0, 0, 0, time.Local)
	daysSinceMonday := (int(jan4.Weekday()) + 6) % 7
	start := jan4.AddDate(0, 0, -daysSinceMonday+(week-1)*7)
	actualYear, actualWeek := start.ISOWeek()
	if actualYear != year || actualWeek != week {
		return time.Time{}, ErrGoreTexInvalidForm
	}
	return start, nil
}

func addGoreTexRate(target map[string]*goreTexRateAccumulator, key, result string) {
	item := target[key]
	if item == nil {
		item = &goreTexRateAccumulator{}
		target[key] = item
	}
	item.Total++
	if result == "PASS" {
		item.Pass++
	}
}

func goreTexRateItems(source map[string]*goreTexRateAccumulator) []types.GoreTexDashboardRateItem {
	items := make([]types.GoreTexDashboardRateItem, 0, len(source))
	for label, value := range source {
		items = append(items, types.GoreTexDashboardRateItem{Label: label, Pass: value.Pass, Total: value.Total, Rate: goreTexRate(value)})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Rate == items[j].Rate {
			return items[i].Label < items[j].Label
		}
		return items[i].Rate > items[j].Rate
	})
	return items
}

func goreTexTrendPoints(source map[string]*goreTexRateAccumulator) []types.GoreTexDashboardTrendPoint {
	labels := make(map[string]string, len(source))
	for key := range source {
		date, _ := time.Parse("2006-01-02", key)
		labels[key] = date.Format("02/01")
	}
	return goreTexTrendPointsWithLabels(source, labels)
}

func goreTexTrendPointsWithLabels(source map[string]*goreTexRateAccumulator, labels map[string]string) []types.GoreTexDashboardTrendPoint {
	keys := make([]string, 0, len(source))
	for key := range source {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	points := make([]types.GoreTexDashboardTrendPoint, 0, len(keys))
	for _, key := range keys {
		value := source[key]
		points = append(points, types.GoreTexDashboardTrendPoint{
			Date: key, Label: labels[key], Pass: value.Pass, Total: value.Total,
			Rate: goreTexRate(value), Target: goreTexDashboardTarget,
		})
	}
	return points
}

func goreTexRate(value *goreTexRateAccumulator) float64 {
	if value == nil || value.Total == 0 {
		return 0
	}
	return math.Round(float64(value.Pass)*1000/float64(value.Total)) / 10
}

func isGoreTexRRdyStyle(style string) bool {
	normalized := strings.NewReplacer(".", "", "-", "", "_", "", " ", "").Replace(strings.ToUpper(style))
	return strings.Contains(normalized, "RRDY")
}

func goreTexParetoItems(rows []goreTexWaterproofRow) []types.GoreTexDashboardParetoItem {
	counts := make(map[string]float64)
	for _, group := range goreTexDefectGroups {
		counts[group.Label] = 0
	}
	hasData := false
	for _, row := range rows {
		var data goreTexDashboardWaterproofData
		if json.Unmarshal([]byte(row.DataJSON), &data) != nil {
			continue
		}
		hasData = true
		for _, group := range goreTexDefectGroups {
			for _, rowID := range group.Rows {
				if value, ok := goreTexDashboardNumber(data.Totals[rowID]); ok {
					counts[group.Label] += value
					continue
				}
				for key, rawValue := range data.Counts {
					if strings.HasPrefix(key, rowID+":") {
						if value, ok := goreTexDashboardNumber(rawValue); ok {
							counts[group.Label] += value
						}
					}
				}
			}
		}
	}
	if !hasData {
		return []types.GoreTexDashboardParetoItem{}
	}

	items := make([]types.GoreTexDashboardParetoItem, 0, len(goreTexDefectGroups))
	for _, group := range goreTexDefectGroups {
		items = append(items, types.GoreTexDashboardParetoItem{Label: group.Label, Count: counts[group.Label]})
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})
	total := 0.0
	for _, item := range items {
		total += item.Count
	}
	running := 0.0
	for index := range items {
		running += items[index].Count
		if total > 0 {
			items[index].Percentage = math.Round(items[index].Count*1000/total) / 10
			items[index].Cumulative = math.Round(running*1000/total) / 10
		}
	}
	return items
}

func goreTexDashboardNumber(value interface{}) (float64, bool) {
	switch typed := value.(type) {
	case float64:
		return typed, !math.IsNaN(typed) && !math.IsInf(typed, 0)
	case string:
		normalized := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(typed, ",", ""), "%", ""))
		if normalized == "" {
			return 0, false
		}
		parsed, err := strconv.ParseFloat(normalized, 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}

// Dashboard snapshots.
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

func normalizeGoreTexStyle(value string) string  { return strings.ToUpper(strings.TrimSpace(value)) }
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
			recordKey := row.RecordKey
			if recordKey == "" {
				recordKey = row.Line + " / " + row.StyleName
			}
			if err := syncGoreTexDefectDailySnapshot(tx, row.InspectionDate, recordKey, row.DataJSON, now); err != nil {
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
