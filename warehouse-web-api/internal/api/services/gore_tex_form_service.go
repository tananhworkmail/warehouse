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

func (s *GoreTexFormService) SaveWaterproof(params request.GoreTexWaterproofSubmitRequest) (*types.GoreTexSubmitResult, error) {
	line := strings.TrimSpace(params.Line)
	styleName := strings.TrimSpace(params.StyleName)
	if line == "" || styleName == "" || !json.Valid(params.Data) || !validateGoreTexWaterproofNumbers(params.Data) {
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
			InspectionDate time.Time `gorm:"column:inspection_date"`
			DataJSON       string    `gorm:"column:data_json"`
		}
		found := tx.Raw(`
			SELECT inspection_date, data_json
			FROM dbo.gore_tex_waterproof_forms
			WHERE [line] = ? AND style_name = ?
		`, line, styleName).Scan(&existing)
		if found.Error != nil {
			return found.Error
		}

		if found.RowsAffected == 0 {
			if params.IsEdit {
				return ErrGoreTexFormNotFound
			}
			created = true
			if insertErr := tx.Exec(`
				INSERT INTO dbo.gore_tex_waterproof_forms
					([line], style_name, inspection_date, data_json, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?)
			`, line, styleName, inspectionDate, string(params.Data), now, now).Error; insertErr != nil {
				return insertErr
			}
			return syncGoreTexDefectDailySnapshot(tx, inspectionDate, line+" / "+styleName, string(params.Data), now)
		}

		if syncErr := syncGoreTexDefectDailySnapshot(tx, existing.InspectionDate, line+" / "+styleName, existing.DataJSON, now); syncErr != nil {
			return syncErr
		}

		dataJSON, mergeErr := appendGoreTexEditHistory(params.Data, existing.DataJSON, now)
		if mergeErr != nil {
			return mergeErr
		}
		if updateErr := tx.Exec(`
			UPDATE dbo.gore_tex_waterproof_forms
			SET inspection_date = ?, data_json = ?, updated_at = ?
			WHERE [line] = ? AND style_name = ?
		`, inspectionDate, dataJSON, now, line, styleName).Error; updateErr != nil {
			return updateErr
		}
		return syncGoreTexDefectDailySnapshot(tx, inspectionDate, line+" / "+styleName, dataJSON, now)
	})
	if err != nil {
		return nil, err
	}

	return &types.GoreTexSubmitResult{
		FormType:  "waterproof",
		RecordKey: line + " / " + styleName,
		Created:   created,
	}, nil
}

func (s *GoreTexFormService) SaveCentrifugal(params request.GoreTexCentrifugalSubmitRequest) (*types.GoreTexSubmitResult, error) {
	if !json.Valid(params.Data) || !validateGoreTexCentrifugalResults(params.Data) {
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
			return syncGoreTexPassDailySnapshot(tx, inspectionDate, string(params.Data), now)
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
		return syncGoreTexPassDailySnapshot(tx, inspectionDate, dataJSON, now)
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
	if !json.Valid(params.Data) || (params.IsEdit && params.AnalysisID == 0) {
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
		SELECT [line], style_name, inspection_date, created_at, updated_at
		FROM dbo.gore_tex_waterproof_forms
		ORDER BY updated_at DESC
	`).Scan(&waterproof)
	if query.Error != nil {
		return nil, query.Error
	}
	for _, form := range waterproof {
		result = append(result, types.GoreTexFormListItem{
			FormType:       "waterproof",
			Title:          "Kiểm tra chất lượng giày chống thấm",
			RecordKey:      form.Line + " / " + form.StyleName,
			Line:           form.Line,
			StyleName:      form.StyleName,
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
	styleName = strings.TrimSpace(styleName)
	if line == "" || styleName == "" {
		return nil, ErrGoreTexInvalidForm
	}

	var form goreTexWaterproofRow
	db, closeDB, err := openGoreTexDB()
	if err != nil {
		return nil, err
	}
	defer closeDB()

	query := db.Raw(`
		SELECT [line], style_name, inspection_date, data_json, created_at, updated_at
		FROM dbo.gore_tex_waterproof_forms
		WHERE [line] = ? AND style_name = ?
	`, line, styleName).Scan(&form)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return nil, ErrGoreTexFormNotFound
	}

	return &types.GoreTexFormDetail{
		FormType:       "waterproof",
		Line:           form.Line,
		StyleName:      form.StyleName,
		InspectionDate: formatGoreTexDate(form.InspectionDate),
		Data:           json.RawMessage(form.DataJSON),
		CreatedAt:      formatGoreTexDateTime(form.CreatedAt),
		UpdatedAt:      formatGoreTexDateTime(form.UpdatedAt),
	}, nil
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
