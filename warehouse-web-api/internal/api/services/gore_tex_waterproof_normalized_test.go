package services

import (
	"database/sql"
	"encoding/json"
	"testing"
)

func waterproofTestPayload() map[string]interface{} {
	counts := make(map[string]interface{})
	for _, row := range goreTexWaterproofRows {
		for _, slot := range goreTexWaterproofSlots {
			if slot == "overtime" {
				counts[row.ID+":"+slot] = ""
			} else {
				counts[row.ID+":"+slot] = "0"
			}
		}
	}
	inspection := make(map[string]interface{})
	for _, slot := range goreTexWaterproofSlots {
		if slot == "overtime" {
			inspection[slot] = ""
		} else {
			inspection[slot] = "0"
		}
	}
	inspection["0830"] = "10"
	return map[string]interface{}{
		"meta":          map[string]interface{}{"line": "Ủng", "styleName": "Test", "inspectionDate": "2026-08-20"},
		"counts":        counts,
		"summaryCounts": map[string]interface{}{"inspection": inspection},
		"signatures": map[string]interface{}{
			"productionSupervisor": "A", "kcs": "B", "fieldOfficer": "C", "preparedBy": "D",
		},
	}
}

func TestGoreTexDashboardComparisonRange(t *testing.T) {
	tests := []struct {
		period string
		value  string
		from   string
		to     string
	}{
		{period: "day", value: "2026-W34", from: "2026-08-17", to: "2026-08-24"},
		{period: "week", value: "2026-08", from: "2026-08-01", to: "2026-09-01"},
		{period: "month", value: "2026", from: "2026-01-01", to: "2027-01-01"},
	}
	for _, test := range tests {
		start, end, err := goreTexDashboardComparisonRange(test.period, test.value)
		if err != nil {
			t.Fatalf("%s %s: %v", test.period, test.value, err)
		}
		if start.Format("2006-01-02") != test.from || end.Format("2006-01-02") != test.to {
			t.Fatalf("%s %s: got %s to %s", test.period, test.value, start, end)
		}
	}
}

func TestGoreTexDashboardComparisonRangeRejectsInvalidValue(t *testing.T) {
	for _, test := range [][2]string{{"day", "2026-W99"}, {"week", "2026-13"}, {"month", "26"}, {"year", "2026"}} {
		if _, _, err := goreTexDashboardComparisonRange(test[0], test[1]); err == nil {
			t.Fatalf("expected %s %s to be rejected", test[0], test[1])
		}
	}
}

func TestGoreTexDashboardSelectedRange(t *testing.T) {
	tests := []struct {
		period string
		value  string
		from   string
		to     string
	}{
		{period: "day", value: "2026-08-21", from: "2026-08-21", to: "2026-08-22"},
		{period: "week", value: "2026-W34", from: "2026-08-17", to: "2026-08-24"},
		{period: "month", value: "2026-08", from: "2026-08-01", to: "2026-09-01"},
	}
	for _, test := range tests {
		start, end, err := goreTexDashboardSelectedRange(test.period, test.value)
		if err != nil {
			t.Fatalf("%s %s: %v", test.period, test.value, err)
		}
		if start.Format("2006-01-02") != test.from || end.Format("2006-01-02") != test.to {
			t.Fatalf("%s %s: got %s to %s", test.period, test.value, start, end)
		}
	}
}

func TestGoreTexParetoItemsKeepsAllIssueNames(t *testing.T) {
	items := goreTexParetoItems([]goreTexWaterproofRow{{DataJSON: `{"totals":{"toe-left":2}}`}})
	if len(items) != len(goreTexDefectGroups) {
		t.Fatalf("expected %d issues, got %d", len(goreTexDefectGroups), len(items))
	}
	if items[0].Label != goreTexDefectGroups[0].Label || items[0].Count != 2 || items[0].Percentage != 100 || items[0].Cumulative != 100 {
		t.Fatalf("unexpected leading issue: %+v", items[0])
	}
}

func TestParseGoreTexWaterproofPayloadAllowsBlankOvertimeAndRejectsNegative(t *testing.T) {
	payload := waterproofTestPayload()
	encoded, _ := json.Marshal(payload)
	if _, inspection, _, err := parseGoreTexWaterproofPayload(encoded); err != nil {
		t.Fatalf("blank overtime should be accepted: %v", err)
	} else if inspection["overtime"] != nil {
		t.Fatal("blank overtime should remain NULL")
	}

	payload["counts"].(map[string]interface{})["toe-left:0830"] = "-1"
	encoded, _ = json.Marshal(payload)
	if _, _, _, err := parseGoreTexWaterproofPayload(encoded); err == nil {
		t.Fatal("negative input must be rejected")
	}

	payload = waterproofTestPayload()
	payload["counts"].(map[string]interface{})["toe-left:0830"] = "11"
	encoded, _ = json.Marshal(payload)
	if _, _, _, err := parseGoreTexWaterproofPayload(encoded); err == nil {
		t.Fatal("defects greater than inspected quantity must be rejected")
	}
}

func TestBuildGoreTexWaterproofDataCalculatesYellowCells(t *testing.T) {
	header := goreTexWaterproofHeader{
		Line: "Ủng", StyleName: "Test",
		TotalInspection0830: sql.NullInt64{Int64: 10, Valid: true},
	}
	detail := goreTexWaterproofDetail{
		Defect: "toe", ShoeSide: 1,
		Time0830: sql.NullInt64{Int64: 2, Valid: true},
	}
	data := buildGoreTexWaterproofData(header, []goreTexWaterproofDetail{detail})

	if got := data["totals"].(map[string]interface{})["toe-left"]; got != "2" {
		t.Fatalf("toe total = %v, want 2", got)
	}
	if got := data["rates"].(map[string]interface{})["toe-left"]; got != "20" {
		t.Fatalf("toe defect rate = %v, want 20", got)
	}
	summary := data["summaryCounts"].(map[string]interface{})
	if got := summary["passed"].(map[string]interface{})["0830"]; got != "8" {
		t.Fatalf("pass quantity = %v, want 8", got)
	}
	if got := summary["passRates"].(map[string]interface{})["0830"]; got != "80" {
		t.Fatalf("pass rate = %v, want 80", got)
	}
}

func TestGoreTexWaterproofPassCountsUsesCalculatedSummary(t *testing.T) {
	data := goreTexDashboardWaterproofData{SummaryTotals: map[string]interface{}{
		"inspection": "540",
		"defects":    "12",
	}}
	pass, total, ok := goreTexWaterproofPassCounts(data)
	if !ok || pass != 528 || total != 540 {
		t.Fatalf("pass/total = %d/%d (ok=%v), want 528/540", pass, total, ok)
	}

	data.SummaryTotals["inspection"] = "0"
	if pass, total, ok := goreTexWaterproofPassCounts(data); !ok || pass != 0 || total != 0 {
		t.Fatal("a newly submitted zero-total report must still become a dashboard item")
	}
}

func TestNormalizeGoreTexJSONStylesUppercasesNames(t *testing.T) {
	raw := json.RawMessage(`{"entries":[{"style":" Clifton 11 gtx ","result":"PASS"}],"keep":"value"}`)
	normalized, err := normalizeGoreTexJSONStyles(raw, "entries", "style")
	if err != nil {
		t.Fatal(err)
	}
	var payload struct {
		Entries []struct {
			Style string `json:"style"`
		} `json:"entries"`
		Keep string `json:"keep"`
	}
	if err := json.Unmarshal(normalized, &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Entries[0].Style != "CLIFTON 11 GTX" || payload.Keep != "value" {
		t.Fatalf("unexpected normalized payload: %s", normalized)
	}
}
