package services

import (
	"testing"
	"time"
)

func TestGoreTexISOWeekStart(t *testing.T) {
	start, err := goreTexISOWeekStart(2021, 27)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if actual := start.Format("2006-01-02"); actual != "2021-07-05" {
		t.Fatalf("expected 2021-07-05, got %s", actual)
	}
	if start.Weekday() != time.Monday {
		t.Fatalf("expected Monday, got %s", start.Weekday())
	}
}

func TestGoreTexISOWeekStartRejectsMissingWeek53(t *testing.T) {
	if _, err := goreTexISOWeekStart(2021, 53); err == nil {
		t.Fatal("expected invalid week error")
	}
}

func TestGoreTexDashboardRateAndRRdyStyle(t *testing.T) {
	value := &goreTexRateAccumulator{Pass: 7, Total: 8}
	if actual := goreTexRate(value); actual != 87.5 {
		t.Fatalf("expected 87.5, got %v", actual)
	}
	for _, style := range []string{"TERREX R.RDY", "R RDY 2.0", "rrdy_kids"} {
		if !isGoreTexRRdyStyle(style) {
			t.Fatalf("expected %q to match R.RDY", style)
		}
	}
}
