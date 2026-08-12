package services

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestParseGoreTexDate(t *testing.T) {
	value, err := parseGoreTexDate("2026-07-27")
	if err != nil {
		t.Fatalf("parse valid date: %v", err)
	}
	if formatted := formatGoreTexDate(value); formatted != "2026-07-27" {
		t.Fatalf("unexpected formatted date: %s", formatted)
	}
}

func TestFormatGoreTexDateTimeDoesNotApplyTimezoneOffset(t *testing.T) {
	value := time.Date(2026, 7, 27, 15, 30, 45, 0, time.UTC)
	if formatted := formatGoreTexDateTime(value); formatted != "2026-07-27T15:30:45" {
		t.Fatalf("unexpected formatted datetime: %s", formatted)
	}
}

func TestAppendGoreTexEditHistoryPreservesPreviousVersions(t *testing.T) {
	previous := `{
		"value":"old",
		"_editHistory":[{"editedAt":"2026-07-26T08:00:00","data":{"value":"older"}}]
	}`
	next := json.RawMessage(`{"value":"new"}`)
	editedAt := time.Date(2026, 7, 27, 9, 15, 0, 0, time.Local)

	merged, err := appendGoreTexEditHistory(next, previous, editedAt)
	if err != nil {
		t.Fatalf("append history: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(merged), &result); err != nil {
		t.Fatalf("decode merged data: %v", err)
	}
	if result["value"] != "new" {
		t.Fatalf("new data was not preserved")
	}

	history, ok := result["_editHistory"].([]interface{})
	if !ok || len(history) != 2 {
		t.Fatalf("expected 2 history entries, got %#v", result["_editHistory"])
	}
	latest := history[1].(map[string]interface{})
	snapshot := latest["data"].(map[string]interface{})
	if snapshot["value"] != "old" {
		t.Fatalf("unexpected snapshot: %#v", snapshot)
	}
	if _, nested := snapshot["_editHistory"]; nested {
		t.Fatalf("history snapshot must not contain nested history")
	}
}

func TestExtractGoreTexAnalysisDatesReturnsUniqueDates(t *testing.T) {
	data := `{
		"records":[
			{"testDate":"2026-07-20","improvementDate":"2026-07-22"},
			{"testDate":"2026-07-20","improvementDate":"2026-07-25"},
			{"testDate":"2026-07-21","improvementDate":"2026-07-25"}
		]
	}`

	testDates, improvementDates := extractGoreTexAnalysisDates(data)
	if len(testDates) != 2 || testDates[0] != "2026-07-20" || testDates[1] != "2026-07-21" {
		t.Fatalf("unexpected test dates: %#v", testDates)
	}
	if len(improvementDates) != 2 || improvementDates[0] != "2026-07-22" || improvementDates[1] != "2026-07-25" {
		t.Fatalf("unexpected improvement dates: %#v", improvementDates)
	}
}

func TestParseGoreTexDateRejectsInvalidValue(t *testing.T) {
	_, err := parseGoreTexDate("27/07/2026")
	if !errors.Is(err, ErrGoreTexInvalidForm) {
		t.Fatalf("expected invalid form error, got %v", err)
	}
}
