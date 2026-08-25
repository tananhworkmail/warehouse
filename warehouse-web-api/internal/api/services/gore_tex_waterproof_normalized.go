package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"gorm.io/gorm"
)

type goreTexWaterproofPayload struct {
	Meta struct {
		Line           string `json:"line"`
		StyleName      string `json:"styleName"`
		InspectionDate string `json:"inspectionDate"`
	} `json:"meta"`
	Counts        map[string]interface{}            `json:"counts"`
	SummaryCounts map[string]map[string]interface{} `json:"summaryCounts"`
	Signatures    struct {
		ProductionSupervisor string `json:"productionSupervisor"`
		KCS                  string `json:"kcs"`
		FieldOfficer         string `json:"fieldOfficer"`
		PreparedBy           string `json:"preparedBy"`
	} `json:"signatures"`
	EditHistory []struct {
		EditedAt string          `json:"editedAt"`
		Data     json.RawMessage `json:"data"`
	} `json:"_editHistory"`
}

type goreTexWaterproofHeader struct {
	ReportID                string        `gorm:"column:report_id"`
	Line                    string        `gorm:"column:line"`
	DepName                 string        `gorm:"column:dep_name"`
	StyleName               string        `gorm:"column:style_name"`
	InspectionDate          time.Time     `gorm:"column:inspection_date"`
	TotalInspection0830     sql.NullInt64 `gorm:"column:total_inspection_0830"`
	TotalInspection0930     sql.NullInt64 `gorm:"column:total_inspection_0930"`
	TotalInspection1030     sql.NullInt64 `gorm:"column:total_inspection_1030"`
	TotalInspection1130     sql.NullInt64 `gorm:"column:total_inspection_1130"`
	TotalInspection1330     sql.NullInt64 `gorm:"column:total_inspection_1330"`
	TotalInspection1430     sql.NullInt64 `gorm:"column:total_inspection_1430"`
	TotalInspection1530     sql.NullInt64 `gorm:"column:total_inspection_1530"`
	TotalInspection1630     sql.NullInt64 `gorm:"column:total_inspection_1630"`
	TotalInspectionOvertime sql.NullInt64 `gorm:"column:total_inspection_overtime"`
	PreparedBy              string        `gorm:"column:prepared_by"`
	FieldOfficer            string        `gorm:"column:field_officer"`
	KCSOfficer              string        `gorm:"column:kcs_officer"`
	ProductionSupervisor    string        `gorm:"column:production_supervisor"`
	CreatedBy               string        `gorm:"column:created_by"`
	CreatedAt               time.Time     `gorm:"column:created_at"`
	UpdatedBy               string        `gorm:"column:updated_by"`
	UpdatedAt               time.Time     `gorm:"column:updated_at"`
	VersionNo               int           `gorm:"column:version_no"`
	RevisionID              int64         `gorm:"column:revision_id"`
	ArchivedAt              time.Time     `gorm:"column:archived_at"`
}

type goreTexWaterproofDetail struct {
	ID          int           `gorm:"column:id"`
	ReportID    string        `gorm:"column:report_id"`
	Defect      string        `gorm:"column:defect"`
	Time0830    sql.NullInt64 `gorm:"column:time_0830"`
	Time0930    sql.NullInt64 `gorm:"column:time_0930"`
	Time1030    sql.NullInt64 `gorm:"column:time_1030"`
	Time1130    sql.NullInt64 `gorm:"column:time_1130"`
	Time1330    sql.NullInt64 `gorm:"column:time_1330"`
	Time1430    sql.NullInt64 `gorm:"column:time_1430"`
	Time1530    sql.NullInt64 `gorm:"column:time_1530"`
	Time1630    sql.NullInt64 `gorm:"column:time_1630"`
	Overtime    sql.NullInt64 `gorm:"column:overtime"`
	TotalDefect sql.NullInt64 `gorm:"column:total_defect"`
	DefectRate  float64       `gorm:"column:defect_rate"`
	ShoeSide    uint8         `gorm:"column:shoe_side"`
	CreatedBy   string        `gorm:"column:created_by"`
	CreatedAt   time.Time     `gorm:"column:created_at"`
	UpdatedBy   string        `gorm:"column:updated_by"`
	UpdatedAt   time.Time     `gorm:"column:updated_at"`
}

type goreTexWaterproofInput struct {
	InspectionBySlot map[string]*int
	DetailByRowSlot  map[string]map[string]*int
	Signatures       goreTexWaterproofPayload
}

var goreTexWaterproofSlots = []string{"0830", "0930", "1030", "1130", "1330", "1430", "1530", "1630", "overtime"}

var goreTexWaterproofRows = []struct {
	ID       string
	Defect   string
	ShoeSide uint8
}{
	{ID: "toe-left", Defect: "toe", ShoeSide: 1},
	{ID: "toe-right", Defect: "toe", ShoeSide: 2},
	{ID: "heel-left", Defect: "heel", ShoeSide: 1},
	{ID: "heel-right", Defect: "heel", ShoeSide: 2},
	{ID: "medial-left", Defect: "medial", ShoeSide: 1},
	{ID: "medial-right", Defect: "medial", ShoeSide: 2},
	{ID: "lateral-left", Defect: "lateral", ShoeSide: 1},
	{ID: "lateral-right", Defect: "lateral", ShoeSide: 2},
	{ID: "material-left", Defect: "material", ShoeSide: 1},
	{ID: "material-right", Defect: "material", ShoeSide: 2},
	{ID: "attaching-left", Defect: "attaching", ShoeSide: 1},
	{ID: "attaching-right", Defect: "attaching", ShoeSide: 2},
	{ID: "wrinkled-left", Defect: "wrinkled", ShoeSide: 1},
	{ID: "wrinkled-right", Defect: "wrinkled", ShoeSide: 2},
	{ID: "zigzag-quarter-left", Defect: "zigzag-quarter", ShoeSide: 1},
	{ID: "zigzag-quarter-right", Defect: "zigzag-quarter", ShoeSide: 2},
	{ID: "zigzag-sole-left", Defect: "zigzag-sole", ShoeSide: 1},
	{ID: "zigzag-sole-right", Defect: "zigzag-sole", ShoeSide: 2},
	{ID: "bonding-toe-left", Defect: "bonding-toe", ShoeSide: 1},
	{ID: "bonding-toe-right", Defect: "bonding-toe", ShoeSide: 2},
	{ID: "bonding-heel-left", Defect: "bonding-heel", ShoeSide: 1},
	{ID: "bonding-heel-right", Defect: "bonding-heel", ShoeSide: 2},
}

const goreTexWaterproofHeaderSelect = `
	SELECT
		ReportID AS report_id,
		DepNo AS line,
		COALESCE(DepName, '') AS dep_name,
		StyleName AS style_name,
		InspectionDate AS inspection_date,
		TotalInspection0830 AS total_inspection_0830,
		TotalInspection0930 AS total_inspection_0930,
		TotalInspection1030 AS total_inspection_1030,
		TotalInspection1130 AS total_inspection_1130,
		TotalInspection1330 AS total_inspection_1330,
		TotalInspection1430 AS total_inspection_1430,
		TotalInspection1530 AS total_inspection_1530,
		TotalInspection1630 AS total_inspection_1630,
		TotalInspectionOvertime AS total_inspection_overtime,
		COALESCE(PreparedBy, '') AS prepared_by,
		COALESCE(FieldOfficer, '') AS field_officer,
		COALESCE(KCSOfficer, '') AS kcs_officer,
		COALESCE(ProductionSupervisor, '') AS production_supervisor,
		COALESCE(CreatedBy, '') AS created_by,
		CreatedAt AS created_at,
		COALESCE(UpdatedBy, '') AS updated_by,
		COALESCE(UpdatedAt, CreatedAt) AS updated_at
	FROM dbo.GoreTex_Waterproof_Report
`

const goreTexWaterproofDetailSelect = `
	SELECT
		ID AS id,
		ReportID AS report_id,
		Defect AS defect,
		Time0830 AS time_0830,
		Time0930 AS time_0930,
		Time1030 AS time_1030,
		Time1130 AS time_1130,
		Time1330 AS time_1330,
		Time1430 AS time_1430,
		Time1530 AS time_1530,
		Time1630 AS time_1630,
		Overtime AS overtime,
		TotalDefect AS total_defect,
		COALESCE(DefectRate, 0) AS defect_rate,
		ShoeSide AS shoe_side,
		COALESCE(CreatedBy, '') AS created_by,
		CreatedAt AS created_at,
		COALESCE(UpdatedBy, '') AS updated_by,
		COALESCE(UpdatedAt, CreatedAt) AS updated_at
	FROM dbo.GoreTex_Waterproof_Report_Detail
`

func parseGoreTexWaterproofPayload(data json.RawMessage) (*goreTexWaterproofPayload, map[string]*int, map[string]map[string]*int, error) {
	var payload goreTexWaterproofPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, nil, nil, ErrGoreTexInvalidForm
	}

	inspection := make(map[string]*int, len(goreTexWaterproofSlots))
	detail := make(map[string]map[string]*int, len(goreTexWaterproofRows))
	for _, slot := range goreTexWaterproofSlots {
		value, err := goreTexWaterproofCount(payload.SummaryCounts["inspection"][slot], slot == "overtime")
		if err != nil {
			return nil, nil, nil, err
		}
		inspection[slot] = value
	}
	for _, row := range goreTexWaterproofRows {
		detail[row.ID] = make(map[string]*int, len(goreTexWaterproofSlots))
		for _, slot := range goreTexWaterproofSlots {
			value, err := goreTexWaterproofCount(payload.Counts[row.ID+":"+slot], slot == "overtime")
			if err != nil {
				return nil, nil, nil, err
			}
			detail[row.ID][slot] = value
		}
	}
	for _, slot := range goreTexWaterproofSlots {
		defects := 0
		for _, row := range goreTexWaterproofRows {
			defects += intOrZero(detail[row.ID][slot])
		}
		if defects > intOrZero(inspection[slot]) {
			return nil, nil, nil, ErrGoreTexInvalidForm
		}
	}

	if strings.TrimSpace(payload.Signatures.PreparedBy) == "" ||
		strings.TrimSpace(payload.Signatures.FieldOfficer) == "" ||
		strings.TrimSpace(payload.Signatures.KCS) == "" ||
		strings.TrimSpace(payload.Signatures.ProductionSupervisor) == "" {
		return nil, nil, nil, ErrGoreTexInvalidForm
	}
	return &payload, inspection, detail, nil
}

func parseLegacyGoreTexWaterproofPayload(data json.RawMessage) (*goreTexWaterproofPayload, map[string]*int, map[string]map[string]*int, error) {
	var payload goreTexWaterproofPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, nil, nil, ErrGoreTexInvalidForm
	}
	zero := 0
	inspection := make(map[string]*int, len(goreTexWaterproofSlots))
	detail := make(map[string]map[string]*int, len(goreTexWaterproofRows))
	for _, slot := range goreTexWaterproofSlots {
		value, err := goreTexWaterproofCount(payload.SummaryCounts["inspection"][slot], true)
		if err != nil {
			value = &zero
		}
		if value == nil && slot != "overtime" {
			value = &zero
		}
		inspection[slot] = value
	}
	for _, row := range goreTexWaterproofRows {
		detail[row.ID] = make(map[string]*int, len(goreTexWaterproofSlots))
		for _, slot := range goreTexWaterproofSlots {
			value, err := goreTexWaterproofCount(payload.Counts[row.ID+":"+slot], true)
			if err != nil {
				value = &zero
			}
			if value == nil && slot != "overtime" {
				value = &zero
			}
			detail[row.ID][slot] = value
		}
	}
	return &payload, inspection, detail, nil
}

func goreTexWaterproofCount(raw interface{}, optional bool) (*int, error) {
	value := strings.TrimSpace(fmt.Sprint(raw))
	if raw == nil || value == "" || value == "<nil>" {
		if optional {
			return nil, nil
		}
		return nil, ErrGoreTexInvalidForm
	}
	number, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", "."), 64)
	if err != nil || math.IsNaN(number) || math.IsInf(number, 0) || number < 0 || math.Trunc(number) != number || number > math.MaxInt32 {
		return nil, ErrGoreTexInvalidForm
	}
	result := int(number)
	return &result, nil
}

func nullableIntValue(value *int) interface{} {
	if value == nil {
		return nil
	}
	return *value
}

func intOrZero(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

func nullIntPointer(value sql.NullInt64) *int {
	if !value.Valid {
		return nil
	}
	result := int(value.Int64)
	return &result
}

func headerInspectionSlots(header goreTexWaterproofHeader) map[string]*int {
	return map[string]*int{
		"0830": nullIntPointer(header.TotalInspection0830), "0930": nullIntPointer(header.TotalInspection0930),
		"1030": nullIntPointer(header.TotalInspection1030), "1130": nullIntPointer(header.TotalInspection1130),
		"1330": nullIntPointer(header.TotalInspection1330), "1430": nullIntPointer(header.TotalInspection1430),
		"1530": nullIntPointer(header.TotalInspection1530), "1630": nullIntPointer(header.TotalInspection1630),
		"overtime": nullIntPointer(header.TotalInspectionOvertime),
	}
}

func detailSlots(detail goreTexWaterproofDetail) map[string]*int {
	return map[string]*int{
		"0830": nullIntPointer(detail.Time0830), "0930": nullIntPointer(detail.Time0930),
		"1030": nullIntPointer(detail.Time1030), "1130": nullIntPointer(detail.Time1130),
		"1330": nullIntPointer(detail.Time1330), "1430": nullIntPointer(detail.Time1430),
		"1530": nullIntPointer(detail.Time1530), "1630": nullIntPointer(detail.Time1630),
		"overtime": nullIntPointer(detail.Overtime),
	}
}

func goreTexWaterproofRowID(defect string, side uint8) string {
	suffix := "left"
	if side == 2 {
		suffix = "right"
	}
	return strings.TrimSpace(defect) + "-" + suffix
}

func goreTexWaterproofNumber(value float64) string {
	if math.Abs(value-math.Round(value)) < 0.0000001 {
		return strconv.FormatInt(int64(math.Round(value)), 10)
	}
	return strconv.FormatFloat(math.Round(value*100)/100, 'f', -1, 64)
}

func goreTexWaterproofRate(numerator, denominator int) float64 {
	if denominator <= 0 {
		return 0
	}
	return float64(numerator) / float64(denominator) * 100
}

func buildGoreTexWaterproofData(header goreTexWaterproofHeader, details []goreTexWaterproofDetail) map[string]interface{} {
	inspectionSlots := headerInspectionSlots(header)
	counts := make(map[string]interface{}, len(goreTexWaterproofRows)*len(goreTexWaterproofSlots))
	totals := make(map[string]interface{}, len(goreTexWaterproofRows))
	rates := make(map[string]interface{}, len(goreTexWaterproofRows))
	defectsBySlot := make(map[string]int, len(goreTexWaterproofSlots))
	detailByKey := make(map[string]goreTexWaterproofDetail, len(details))
	for _, detail := range details {
		detailByKey[goreTexWaterproofRowID(detail.Defect, detail.ShoeSide)] = detail
	}

	totalInspection := 0
	for _, slot := range goreTexWaterproofSlots {
		totalInspection += intOrZero(inspectionSlots[slot])
	}
	for _, definition := range goreTexWaterproofRows {
		detail, found := detailByKey[definition.ID]
		slots := make(map[string]*int)
		if found {
			slots = detailSlots(detail)
		}
		total := 0
		for _, slot := range goreTexWaterproofSlots {
			value := slots[slot]
			if slot == "overtime" && value == nil {
				counts[definition.ID+":"+slot] = ""
			} else {
				counts[definition.ID+":"+slot] = strconv.Itoa(intOrZero(value))
			}
			total += intOrZero(value)
			defectsBySlot[slot] += intOrZero(value)
		}
		totals[definition.ID] = strconv.Itoa(total)
		rates[definition.ID] = goreTexWaterproofNumber(goreTexWaterproofRate(total, totalInspection))
	}

	inspectionValues := make(map[string]interface{}, len(goreTexWaterproofSlots))
	defectValues := make(map[string]interface{}, len(goreTexWaterproofSlots))
	rateValues := make(map[string]interface{}, len(goreTexWaterproofSlots))
	passedValues := make(map[string]interface{}, len(goreTexWaterproofSlots))
	passRateValues := make(map[string]interface{}, len(goreTexWaterproofSlots))
	totalDefects := 0
	totalPassed := 0
	for _, slot := range goreTexWaterproofSlots {
		inspection := intOrZero(inspectionSlots[slot])
		defects := defectsBySlot[slot]
		passed := inspection - defects
		if slot == "overtime" && inspectionSlots[slot] == nil {
			inspectionValues[slot] = ""
		} else {
			inspectionValues[slot] = strconv.Itoa(inspection)
		}
		defectValues[slot] = strconv.Itoa(defects)
		rateValues[slot] = goreTexWaterproofNumber(goreTexWaterproofRate(defects, inspection))
		passedValues[slot] = strconv.Itoa(passed)
		passRateValues[slot] = goreTexWaterproofNumber(goreTexWaterproofRate(passed, inspection))
		totalDefects += defects
		totalPassed += passed
	}

	return map[string]interface{}{
		"meta": map[string]interface{}{
			"line": header.Line, "styleName": normalizeGoreTexStyle(header.StyleName), "inspectionDate": formatGoreTexDate(header.InspectionDate),
		},
		"counts": counts,
		"totals": totals,
		"rates":  rates,
		"summaryCounts": map[string]interface{}{
			"inspection": inspectionValues, "defects": defectValues, "rates": rateValues,
			"passed": passedValues, "passRates": passRateValues,
		},
		"summaryTotals": map[string]interface{}{
			"inspection": strconv.Itoa(totalInspection), "defects": strconv.Itoa(totalDefects),
			"rate":   goreTexWaterproofNumber(goreTexWaterproofRate(totalDefects, totalInspection)),
			"passed": strconv.Itoa(totalPassed), "passRate": goreTexWaterproofNumber(goreTexWaterproofRate(totalPassed, totalInspection)),
		},
		"signatures": map[string]interface{}{
			"productionSupervisor": header.ProductionSupervisor, "kcs": header.KCSOfficer,
			"fieldOfficer": header.FieldOfficer, "preparedBy": header.PreparedBy,
		},
	}
}

func marshalGoreTexWaterproofData(header goreTexWaterproofHeader, details []goreTexWaterproofDetail, history []interface{}) ([]byte, error) {
	data := buildGoreTexWaterproofData(header, details)
	if len(history) > 0 {
		data["_editHistory"] = history
	}
	return json.Marshal(data)
}

func loadGoreTexWaterproofHeaderByBusinessKey(db *gorm.DB, line, styleName string) (goreTexWaterproofHeader, bool, error) {
	var header goreTexWaterproofHeader
	query := db.Raw(goreTexWaterproofHeaderSelect+`
		WHERE DepNo = ? AND StyleName = ?
		ORDER BY COALESCE(UpdatedAt, CreatedAt) DESC
	`, line, styleName).Scan(&header)
	return header, query.RowsAffected > 0, query.Error
}

func loadGoreTexWaterproofDetails(db *gorm.DB, reportID string) ([]goreTexWaterproofDetail, error) {
	var details []goreTexWaterproofDetail
	query := db.Raw(goreTexWaterproofDetailSelect+` WHERE ReportID = ? ORDER BY ID`, reportID).Scan(&details)
	return details, query.Error
}

func archiveGoreTexWaterproofRevision(tx *gorm.DB, reportID, userID string, now time.Time) error {
	var revision struct {
		RevisionID int64 `gorm:"column:revision_id"`
	}
	result := tx.Raw(`
		DECLARE @VersionNo int = ISNULL((
			SELECT MAX(VersionNo) FROM dbo.GoreTex_Waterproof_Report_Revision WHERE ReportID = ?
		), 0) + 1;
		INSERT INTO dbo.GoreTex_Waterproof_Report_Revision (
			ReportID, VersionNo, DepNo, DepName, StyleName, InspectionDate,
			TotalInspection0830, TotalInspection0930, TotalInspection1030, TotalInspection1130,
			TotalInspection1330, TotalInspection1430, TotalInspection1530, TotalInspection1630,
			TotalInspectionOvertime, PreparedBy, FieldOfficer, KCSOfficer, ProductionSupervisor,
			CreatedBy, CreatedAt, UpdatedBy, UpdatedAt, ArchivedBy, ArchivedAt
		)
		OUTPUT INSERTED.RevisionID AS revision_id
		SELECT ReportID, @VersionNo, DepNo, DepName, StyleName, InspectionDate,
			TotalInspection0830, TotalInspection0930, TotalInspection1030, TotalInspection1130,
			TotalInspection1330, TotalInspection1430, TotalInspection1530, TotalInspection1630,
			TotalInspectionOvertime, PreparedBy, FieldOfficer, KCSOfficer, ProductionSupervisor,
			CreatedBy, CreatedAt, UpdatedBy, UpdatedAt, ?, ?
		FROM dbo.GoreTex_Waterproof_Report WHERE ReportID = ?;
	`, reportID, userID, now, reportID).Scan(&revision)
	if result.Error != nil {
		return result.Error
	}
	if revision.RevisionID == 0 {
		return ErrGoreTexFormNotFound
	}
	return tx.Exec(`
		INSERT INTO dbo.GoreTex_Waterproof_Report_Detail_Revision (
			RevisionID, DetailID, ReportID, Defect, Time0830, Time0930, Time1030, Time1130,
			Time1330, Time1430, Time1530, Time1630, Overtime, TotalDefect, DefectRate, ShoeSide,
			CreatedBy, CreatedAt, UpdatedBy, UpdatedAt
		)
		SELECT ?, ID, ReportID, Defect, Time0830, Time0930, Time1030, Time1130,
			Time1330, Time1430, Time1530, Time1630, Overtime, TotalDefect, DefectRate, ShoeSide,
			CreatedBy, CreatedAt, UpdatedBy, UpdatedAt
		FROM dbo.GoreTex_Waterproof_Report_Detail WHERE ReportID = ?
	`, revision.RevisionID, reportID).Error
}

func insertGoreTexWaterproofDetails(tx *gorm.DB, reportID string, rows map[string]map[string]*int, totalInspection int, userID string, createdAt, updatedAt time.Time) error {
	for _, definition := range goreTexWaterproofRows {
		slots := rows[definition.ID]
		total := 0
		for _, slot := range goreTexWaterproofSlots {
			total += intOrZero(slots[slot])
		}
		rate := goreTexWaterproofRate(total, totalInspection)
		if err := tx.Exec(`
			INSERT INTO dbo.GoreTex_Waterproof_Report_Detail (
				ReportID, Defect, Time0830, Time0930, Time1030, Time1130, Time1330,
				Time1430, Time1530, Time1630, Overtime, TotalDefect, DefectRate, ShoeSide,
				CreatedBy, CreatedAt, UpdatedBy, UpdatedAt
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, reportID, definition.Defect,
			nullableIntValue(slots["0830"]), nullableIntValue(slots["0930"]), nullableIntValue(slots["1030"]), nullableIntValue(slots["1130"]),
			nullableIntValue(slots["1330"]), nullableIntValue(slots["1430"]), nullableIntValue(slots["1530"]), nullableIntValue(slots["1630"]),
			nullableIntValue(slots["overtime"]), total, rate, definition.ShoeSide, userID, createdAt, userID, updatedAt).Error; err != nil {
			return err
		}
	}
	return nil
}

func saveNormalizedGoreTexWaterproof(params request.GoreTexWaterproofSubmitRequest, inspectionDate time.Time) (*types.GoreTexSubmitResult, error) {
	line := strings.TrimSpace(params.Line)
	styleName := normalizeGoreTexStyle(params.StyleName)
	payload, inspection, detail, err := parseGoreTexWaterproofPayload(params.Data)
	if err != nil || line == "" || styleName == "" {
		return nil, ErrGoreTexInvalidForm
	}

	db, closeDB, err := openGoreTexDB()
	if err != nil {
		return nil, err
	}
	defer closeDB()

	now := goreTexNow()
	userID := "web"
	created := false
	reportID := ""
	err = db.Transaction(func(tx *gorm.DB) error {
		existing, found, findErr := loadGoreTexWaterproofHeaderByBusinessKey(tx, line, styleName)
		if findErr != nil {
			return findErr
		}
		if !found && params.IsEdit {
			return ErrGoreTexFormNotFound
		}

		totalInspection := 0
		for _, slot := range goreTexWaterproofSlots {
			totalInspection += intOrZero(inspection[slot])
		}
		if !found {
			created = true
			var inserted struct {
				ReportID string `gorm:"column:report_id"`
			}
			result := tx.Raw(`
				DECLARE @ReportID varchar(50) = CONVERT(varchar(36), NEWID());
				INSERT INTO dbo.GoreTex_Waterproof_Report (
					ReportID, DepNo, DepName, StyleName, InspectionDate,
					TotalInspection0830, TotalInspection0930, TotalInspection1030, TotalInspection1130,
					TotalInspection1330, TotalInspection1430, TotalInspection1530, TotalInspection1630,
					TotalInspectionOvertime, PreparedBy, FieldOfficer, KCSOfficer, ProductionSupervisor,
					CreatedBy, CreatedAt, UpdatedBy, UpdatedAt
				)
				OUTPUT INSERTED.ReportID AS report_id
				VALUES (@ReportID, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, line, styleName, inspectionDate,
				nullableIntValue(inspection["0830"]), nullableIntValue(inspection["0930"]), nullableIntValue(inspection["1030"]), nullableIntValue(inspection["1130"]),
				nullableIntValue(inspection["1330"]), nullableIntValue(inspection["1430"]), nullableIntValue(inspection["1530"]), nullableIntValue(inspection["1630"]),
				nullableIntValue(inspection["overtime"]), strings.TrimSpace(payload.Signatures.PreparedBy), strings.TrimSpace(payload.Signatures.FieldOfficer),
				strings.TrimSpace(payload.Signatures.KCS), strings.TrimSpace(payload.Signatures.ProductionSupervisor), userID, now, userID, now).Scan(&inserted)
			if result.Error != nil {
				return result.Error
			}
			reportID = inserted.ReportID
		} else {
			reportID = existing.ReportID
			if err := archiveGoreTexWaterproofRevision(tx, reportID, userID, now); err != nil {
				return err
			}
			if err := tx.Exec(`
				UPDATE dbo.GoreTex_Waterproof_Report SET
					DepNo = ?, StyleName = ?, InspectionDate = ?,
					TotalInspection0830 = ?, TotalInspection0930 = ?, TotalInspection1030 = ?, TotalInspection1130 = ?,
					TotalInspection1330 = ?, TotalInspection1430 = ?, TotalInspection1530 = ?, TotalInspection1630 = ?,
					TotalInspectionOvertime = ?, PreparedBy = ?, FieldOfficer = ?, KCSOfficer = ?,
					ProductionSupervisor = ?, UpdatedBy = ?, UpdatedAt = ?
				WHERE ReportID = ?
			`, line, styleName, inspectionDate,
				nullableIntValue(inspection["0830"]), nullableIntValue(inspection["0930"]), nullableIntValue(inspection["1030"]), nullableIntValue(inspection["1130"]),
				nullableIntValue(inspection["1330"]), nullableIntValue(inspection["1430"]), nullableIntValue(inspection["1530"]), nullableIntValue(inspection["1630"]),
				nullableIntValue(inspection["overtime"]), strings.TrimSpace(payload.Signatures.PreparedBy), strings.TrimSpace(payload.Signatures.FieldOfficer),
				strings.TrimSpace(payload.Signatures.KCS), strings.TrimSpace(payload.Signatures.ProductionSupervisor), userID, now, reportID).Error; err != nil {
				return err
			}
			if err := tx.Exec(`DELETE FROM dbo.GoreTex_Waterproof_Report_Detail WHERE ReportID = ?`, reportID).Error; err != nil {
				return err
			}
		}

		if err := insertGoreTexWaterproofDetails(tx, reportID, detail, totalInspection, userID, now, now); err != nil {
			return err
		}
		header, _, err := loadGoreTexWaterproofHeaderByBusinessKey(tx, line, styleName)
		if err != nil {
			return err
		}
		details, err := loadGoreTexWaterproofDetails(tx, reportID)
		if err != nil {
			return err
		}
		dataJSON, err := marshalGoreTexWaterproofData(header, details, nil)
		if err != nil {
			return err
		}
		return syncGoreTexDefectDailySnapshot(tx, inspectionDate, reportID, string(dataJSON), now)
	})
	if err != nil {
		return nil, err
	}
	return &types.GoreTexSubmitResult{FormType: "waterproof", RecordKey: line + " / " + styleName, Created: created}, nil
}

func loadGoreTexWaterproofHistory(db *gorm.DB, reportID string) ([]interface{}, error) {
	var revisions []goreTexWaterproofHeader
	query := db.Raw(`
		SELECT RevisionID AS revision_id, ReportID AS report_id, VersionNo AS version_no,
			DepNo AS line, COALESCE(DepName, '') AS dep_name, StyleName AS style_name, InspectionDate AS inspection_date,
			TotalInspection0830 AS total_inspection_0830, TotalInspection0930 AS total_inspection_0930,
			TotalInspection1030 AS total_inspection_1030, TotalInspection1130 AS total_inspection_1130,
			TotalInspection1330 AS total_inspection_1330, TotalInspection1430 AS total_inspection_1430,
			TotalInspection1530 AS total_inspection_1530, TotalInspection1630 AS total_inspection_1630,
			TotalInspectionOvertime AS total_inspection_overtime,
			COALESCE(PreparedBy, '') AS prepared_by, COALESCE(FieldOfficer, '') AS field_officer,
			COALESCE(KCSOfficer, '') AS kcs_officer, COALESCE(ProductionSupervisor, '') AS production_supervisor,
			COALESCE(CreatedBy, '') AS created_by, CreatedAt AS created_at,
			COALESCE(UpdatedBy, '') AS updated_by, UpdatedAt AS updated_at, ArchivedAt AS archived_at
		FROM dbo.GoreTex_Waterproof_Report_Revision
		WHERE ReportID = ? ORDER BY VersionNo
	`, reportID).Scan(&revisions)
	if query.Error != nil {
		return nil, query.Error
	}
	history := make([]interface{}, 0, len(revisions))
	for _, revision := range revisions {
		var details []goreTexWaterproofDetail
		if err := db.Raw(`
			SELECT DetailID AS id, ReportID AS report_id, Defect AS defect,
				Time0830 AS time_0830, Time0930 AS time_0930, Time1030 AS time_1030, Time1130 AS time_1130,
				Time1330 AS time_1330, Time1430 AS time_1430, Time1530 AS time_1530, Time1630 AS time_1630,
				Overtime AS overtime, TotalDefect AS total_defect, COALESCE(DefectRate, 0) AS defect_rate,
				ShoeSide AS shoe_side, COALESCE(CreatedBy, '') AS created_by, CreatedAt AS created_at,
				COALESCE(UpdatedBy, '') AS updated_by, UpdatedAt AS updated_at
			FROM dbo.GoreTex_Waterproof_Report_Detail_Revision WHERE RevisionID = ? ORDER BY DetailID
		`, revision.RevisionID).Scan(&details).Error; err != nil {
			return nil, err
		}
		history = append(history, map[string]interface{}{
			"editedAt": formatGoreTexDateTime(revision.ArchivedAt),
			"data":     buildGoreTexWaterproofData(revision, details),
		})
	}
	return history, nil
}

func getNormalizedGoreTexWaterproof(db *gorm.DB, line, styleName string, withHistory bool) (*types.GoreTexFormDetail, error) {
	header, found, err := loadGoreTexWaterproofHeaderByBusinessKey(db, line, styleName)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrGoreTexFormNotFound
	}
	details, err := loadGoreTexWaterproofDetails(db, header.ReportID)
	if err != nil {
		return nil, err
	}
	var history []interface{}
	if withHistory {
		history, err = loadGoreTexWaterproofHistory(db, header.ReportID)
		if err != nil {
			return nil, err
		}
	}
	data, err := marshalGoreTexWaterproofData(header, details, history)
	if err != nil {
		return nil, err
	}
	return &types.GoreTexFormDetail{
		FormType: "waterproof", Line: header.Line, StyleName: normalizeGoreTexStyle(header.StyleName),
		InspectionDate: formatGoreTexDate(header.InspectionDate), Data: data,
		CreatedAt: formatGoreTexDateTime(header.CreatedAt), UpdatedAt: formatGoreTexDateTime(header.UpdatedAt),
	}, nil
}

func loadNormalizedGoreTexWaterproofDashboardRows(db *gorm.DB, from, to time.Time) ([]goreTexWaterproofRow, error) {
	var headers []goreTexWaterproofHeader
	query := db.Raw(goreTexWaterproofHeaderSelect+`
		WHERE InspectionDate >= ? AND InspectionDate < ? ORDER BY InspectionDate
	`, from, to).Scan(&headers)
	if query.Error != nil {
		return nil, query.Error
	}
	rows := make([]goreTexWaterproofRow, 0, len(headers))
	for _, header := range headers {
		details, err := loadGoreTexWaterproofDetails(db, header.ReportID)
		if err != nil {
			return nil, err
		}
		data, err := marshalGoreTexWaterproofData(header, details, nil)
		if err != nil {
			return nil, err
		}
		rows = append(rows, goreTexWaterproofRow{
			RecordKey: header.ReportID, Line: header.Line, StyleName: normalizeGoreTexStyle(header.StyleName), InspectionDate: header.InspectionDate,
			DataJSON: string(data), CreatedAt: header.CreatedAt, UpdatedAt: header.UpdatedAt,
		})
	}
	return rows, nil
}

type GoreTexWaterproofMigrationResult struct {
	MigratedReports   int
	MigratedRevisions int
	SkippedReports    int
}

func parseLegacyGoreTexTime(value string, fallback time.Time) time.Time {
	value = strings.TrimSpace(value)
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05"} {
		if parsed, err := time.ParseInLocation(layout, value, goreTexNow().Location()); err == nil {
			return parsed
		}
	}
	return fallback
}

func insertLegacyGoreTexWaterproofRevision(
	tx *gorm.DB,
	reportID string,
	versionNo int,
	data json.RawMessage,
	archivedAt time.Time,
	createdAt time.Time,
	fallbackLine string,
	fallbackStyle string,
	fallbackDate time.Time,
) error {
	payload, inspection, details, err := parseLegacyGoreTexWaterproofPayload(data)
	if err != nil {
		return fmt.Errorf("legacy waterproof revision %d is invalid: %w", versionNo, err)
	}
	line := strings.TrimSpace(payload.Meta.Line)
	if line == "" {
		line = fallbackLine
	}
	styleName := strings.TrimSpace(payload.Meta.StyleName)
	if styleName == "" {
		styleName = fallbackStyle
	}
	inspectionDate := fallbackDate
	if parsed, parseErr := parseGoreTexDate(payload.Meta.InspectionDate); parseErr == nil {
		inspectionDate = parsed
	}

	var inserted struct {
		RevisionID int64 `gorm:"column:revision_id"`
	}
	result := tx.Raw(`
		INSERT INTO dbo.GoreTex_Waterproof_Report_Revision (
			ReportID, VersionNo, DepNo, DepName, StyleName, InspectionDate,
			TotalInspection0830, TotalInspection0930, TotalInspection1030, TotalInspection1130,
			TotalInspection1330, TotalInspection1430, TotalInspection1530, TotalInspection1630,
			TotalInspectionOvertime, PreparedBy, FieldOfficer, KCSOfficer, ProductionSupervisor,
			CreatedBy, CreatedAt, UpdatedBy, UpdatedAt, ArchivedBy, ArchivedAt
		)
		OUTPUT INSERTED.RevisionID AS revision_id
		VALUES (?, ?, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, reportID, versionNo, line, styleName, inspectionDate,
		nullableIntValue(inspection["0830"]), nullableIntValue(inspection["0930"]), nullableIntValue(inspection["1030"]), nullableIntValue(inspection["1130"]),
		nullableIntValue(inspection["1330"]), nullableIntValue(inspection["1430"]), nullableIntValue(inspection["1530"]), nullableIntValue(inspection["1630"]),
		nullableIntValue(inspection["overtime"]), strings.TrimSpace(payload.Signatures.PreparedBy), strings.TrimSpace(payload.Signatures.FieldOfficer),
		strings.TrimSpace(payload.Signatures.KCS), strings.TrimSpace(payload.Signatures.ProductionSupervisor),
		"legacy-migration", createdAt, "legacy-migration", archivedAt, "legacy-migration", archivedAt).Scan(&inserted)
	if result.Error != nil {
		return result.Error
	}

	totalInspection := 0
	for _, slot := range goreTexWaterproofSlots {
		totalInspection += intOrZero(inspection[slot])
	}
	for _, definition := range goreTexWaterproofRows {
		slots := details[definition.ID]
		total := 0
		for _, slot := range goreTexWaterproofSlots {
			total += intOrZero(slots[slot])
		}
		if err := tx.Exec(`
			INSERT INTO dbo.GoreTex_Waterproof_Report_Detail_Revision (
				RevisionID, DetailID, ReportID, Defect, Time0830, Time0930, Time1030, Time1130,
				Time1330, Time1430, Time1530, Time1630, Overtime, TotalDefect, DefectRate, ShoeSide,
				CreatedBy, CreatedAt, UpdatedBy, UpdatedAt
			) VALUES (?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, inserted.RevisionID, reportID, definition.Defect,
			nullableIntValue(slots["0830"]), nullableIntValue(slots["0930"]), nullableIntValue(slots["1030"]), nullableIntValue(slots["1130"]),
			nullableIntValue(slots["1330"]), nullableIntValue(slots["1430"]), nullableIntValue(slots["1530"]), nullableIntValue(slots["1630"]),
			nullableIntValue(slots["overtime"]), total, goreTexWaterproofRate(total, totalInspection), definition.ShoeSide,
			"legacy-migration", createdAt, "legacy-migration", archivedAt).Error; err != nil {
			return err
		}
	}
	return nil
}

// MigrateLegacyGoreTexWaterproofForms moves JSON-backed waterproof forms and
// their revision snapshots into the normalized report/detail tables. It is
// idempotent by the report business key (line + style name).
func MigrateLegacyGoreTexWaterproofForms() (*GoreTexWaterproofMigrationResult, error) {
	db, closeDB, err := openGoreTexDB()
	if err != nil {
		return nil, err
	}
	defer closeDB()

	var legacyRows []goreTexWaterproofRow
	if err := db.Raw(`
		SELECT line, style_name, inspection_date, data_json, created_at, updated_at
		FROM dbo.gore_tex_waterproof_forms
		ORDER BY created_at, line, style_name
	`).Scan(&legacyRows).Error; err != nil {
		return nil, err
	}
	result := &GoreTexWaterproofMigrationResult{}
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, legacy := range legacyRows {
			_, found, findErr := loadGoreTexWaterproofHeaderByBusinessKey(tx, legacy.Line, legacy.StyleName)
			if findErr != nil {
				return findErr
			}
			if found {
				result.SkippedReports++
				continue
			}

			payload, inspection, details, parseErr := parseLegacyGoreTexWaterproofPayload(json.RawMessage(legacy.DataJSON))
			if parseErr != nil {
				return fmt.Errorf("legacy waterproof form %s / %s is invalid: %w", legacy.Line, legacy.StyleName, parseErr)
			}
			var inserted struct {
				ReportID string `gorm:"column:report_id"`
			}
			insertResult := tx.Raw(`
				DECLARE @ReportID varchar(50) = CONVERT(varchar(36), NEWID());
				INSERT INTO dbo.GoreTex_Waterproof_Report (
					ReportID, DepNo, DepName, StyleName, InspectionDate,
					TotalInspection0830, TotalInspection0930, TotalInspection1030, TotalInspection1130,
					TotalInspection1330, TotalInspection1430, TotalInspection1530, TotalInspection1630,
					TotalInspectionOvertime, PreparedBy, FieldOfficer, KCSOfficer, ProductionSupervisor,
					CreatedBy, CreatedAt, UpdatedBy, UpdatedAt
				)
				OUTPUT INSERTED.ReportID AS report_id
				VALUES (@ReportID, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, legacy.Line, legacy.StyleName, legacy.InspectionDate,
				nullableIntValue(inspection["0830"]), nullableIntValue(inspection["0930"]), nullableIntValue(inspection["1030"]), nullableIntValue(inspection["1130"]),
				nullableIntValue(inspection["1330"]), nullableIntValue(inspection["1430"]), nullableIntValue(inspection["1530"]), nullableIntValue(inspection["1630"]),
				nullableIntValue(inspection["overtime"]), strings.TrimSpace(payload.Signatures.PreparedBy), strings.TrimSpace(payload.Signatures.FieldOfficer),
				strings.TrimSpace(payload.Signatures.KCS), strings.TrimSpace(payload.Signatures.ProductionSupervisor),
				"legacy-migration", legacy.CreatedAt, "legacy-migration", legacy.UpdatedAt).Scan(&inserted)
			if insertResult.Error != nil {
				return insertResult.Error
			}
			totalInspection := 0
			for _, slot := range goreTexWaterproofSlots {
				totalInspection += intOrZero(inspection[slot])
			}
			if err := insertGoreTexWaterproofDetails(tx, inserted.ReportID, details, totalInspection, "legacy-migration", legacy.CreatedAt, legacy.UpdatedAt); err != nil {
				return err
			}

			for versionIndex, revision := range payload.EditHistory {
				archivedAt := parseLegacyGoreTexTime(revision.EditedAt, legacy.UpdatedAt)
				if err := insertLegacyGoreTexWaterproofRevision(tx, inserted.ReportID, versionIndex+1, revision.Data, archivedAt,
					legacy.CreatedAt, legacy.Line, legacy.StyleName, legacy.InspectionDate); err != nil {
					return err
				}
				result.MigratedRevisions++
			}

			if err := tx.Exec(`
				DELETE FROM dbo.gore_tex_dashboard_defect_snapshots
				WHERE period_type = 'D' AND source_record_key = ?
			`, legacy.Line+" / "+legacy.StyleName).Error; err != nil {
				return err
			}
			header, _, err := loadGoreTexWaterproofHeaderByBusinessKey(tx, legacy.Line, legacy.StyleName)
			if err != nil {
				return err
			}
			normalizedDetails, err := loadGoreTexWaterproofDetails(tx, inserted.ReportID)
			if err != nil {
				return err
			}
			dataJSON, err := marshalGoreTexWaterproofData(header, normalizedDetails, nil)
			if err != nil {
				return err
			}
			if err := syncGoreTexDefectDailySnapshot(tx, legacy.InspectionDate, inserted.ReportID, string(dataJSON), legacy.UpdatedAt); err != nil {
				return err
			}
			result.MigratedReports++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
