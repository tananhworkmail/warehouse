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
)

const goreTexDashboardTarget = 99.0

type goreTexDashboardCentrifugalData struct {
	Entries []struct {
		Style  string `json:"style"`
		Result string `json:"result"`
	} `json:"entries"`
}

type goreTexDashboardWaterproofData struct {
	Counts map[string]interface{} `json:"counts"`
	Totals map[string]interface{} `json:"totals"`
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

	var centrifugalRows []goreTexCentrifugalRow
	query := db.Raw(`
		SELECT inspection_date, data_json
		FROM dbo.gore_tex_centrifugal_forms
		WHERE inspection_date >= ? AND inspection_date < ?
		ORDER BY inspection_date
	`, previousStart, weekEndExclusive).Scan(&centrifugalRows)
	if query.Error != nil {
		return nil, query.Error
	}

	var waterproofRows []goreTexWaterproofRow
	query = db.Raw(`
		SELECT [line], style_name, inspection_date, data_json
		FROM dbo.gore_tex_waterproof_forms
		WHERE inspection_date >= ? AND inspection_date < ?
		ORDER BY inspection_date
	`, previousStart, weekEndExclusive).Scan(&waterproofRows)
	if query.Error != nil {
		return nil, query.Error
	}

	if goreTexDashboardSnapshotsAvailable(db) {
		if err := backfillGoreTexDashboardSnapshots(db, centrifugalRows, waterproofRows); err != nil {
			return nil, err
		}
		return readGoreTexWeeklyDashboardSnapshots(db, year, week, weekStart, previousStart)
	}

	selectedByStyle := make(map[string]*goreTexRateAccumulator)
	previousByStyle := make(map[string]*goreTexRateAccumulator)
	selectedByDate := make(map[string]*goreTexRateAccumulator)
	rRdyByDate := make(map[string]*goreTexRateAccumulator)

	for _, row := range centrifugalRows {
		var data goreTexDashboardCentrifugalData
		if json.Unmarshal([]byte(row.DataJSON), &data) != nil {
			continue
		}
		isSelected := !row.InspectionDate.Before(weekStart) && row.InspectionDate.Before(weekEndExclusive)
		for _, entry := range data.Entries {
			style := strings.TrimSpace(entry.Style)
			result := strings.ToUpper(strings.TrimSpace(entry.Result))
			if style == "" || (result != "PASS" && result != "FAIL") {
				continue
			}
			styleTarget := previousByStyle
			if isSelected {
				styleTarget = selectedByStyle
				dateKey := row.InspectionDate.Format("2006-01-02")
				addGoreTexRate(selectedByDate, dateKey, result)
				if isGoreTexRRdyStyle(style) {
					addGoreTexRate(rRdyByDate, dateKey, result)
				}
			}
			addGoreTexRate(styleTarget, style, result)
		}
	}

	selectedWaterproofRows := make([]goreTexWaterproofRow, 0, len(waterproofRows))
	for _, row := range waterproofRows {
		if !row.InspectionDate.Before(weekStart) && row.InspectionDate.Before(weekEndExclusive) {
			selectedWaterproofRows = append(selectedWaterproofRows, row)
		}
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
		RRdyTrend:            goreTexTrendPoints(rRdyByDate),
		VisualizationResults: goreTexParetoItems(selectedWaterproofRows),
	}, nil
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
	keys := make([]string, 0, len(source))
	for key := range source {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	points := make([]types.GoreTexDashboardTrendPoint, 0, len(keys))
	for _, key := range keys {
		date, _ := time.Parse("2006-01-02", key)
		value := source[key]
		points = append(points, types.GoreTexDashboardTrendPoint{
			Date: key, Label: date.Format("02/01"), Pass: value.Pass, Total: value.Total,
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
	for _, row := range rows {
		var data goreTexDashboardWaterproofData
		if json.Unmarshal([]byte(row.DataJSON), &data) != nil {
			continue
		}
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

	items := make([]types.GoreTexDashboardParetoItem, 0, len(counts))
	for label, count := range counts {
		if count > 0 {
			items = append(items, types.GoreTexDashboardParetoItem{Label: label, Count: count})
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Label < items[j].Label
		}
		return items[i].Count > items[j].Count
	})
	total := 0.0
	for _, item := range items {
		total += item.Count
	}
	running := 0.0
	for index := range items {
		running += items[index].Count
		items[index].Cumulative = math.Round(running*1000/total) / 10
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
