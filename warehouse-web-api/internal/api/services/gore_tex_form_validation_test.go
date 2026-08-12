package services

import (
	"encoding/json"
	"testing"
)

func TestIsGoreTexNonNegativeNumber(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  bool
	}{
		{name: "zero string", value: "0", want: true},
		{name: "positive decimal with comma", value: "1,5", want: true},
		{name: "positive JSON number", value: float64(2), want: true},
		{name: "negative string", value: "-1", want: false},
		{name: "not a number", value: "abc", want: false},
		{name: "empty", value: "", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isGoreTexNonNegativeNumber(test.value); got != test.want {
				t.Fatalf("isGoreTexNonNegativeNumber(%v) = %v, want %v", test.value, got, test.want)
			}
		})
	}
}

func TestValidateGoreTexWaterproofNumbers(t *testing.T) {
	valid := json.RawMessage(`{
		"counts":{"row:slot":"0"},
		"totals":{"row":"1.5"},
		"rates":{"row":0},
		"summaryCounts":{"inspection":{"slot":"0"},"defects":{"slot":"0"},"rates":{"slot":"0"}},
		"summaryTotals":{"inspection":"0","defects":"0","rate":"0"},
		"summaryEdges":{"inspection":{"side":"0","end":"0"},"defects":{"side":"0","end":"0"},"rates":{"side":"0","total":"0"}}
	}`)
	if !validateGoreTexWaterproofNumbers(valid) {
		t.Fatal("expected valid waterproof numeric data")
	}

	invalid := json.RawMessage(`{
		"counts":{"row:slot":"-1"},
		"totals":{"row":"0"},
		"rates":{"row":"0"},
		"summaryCounts":{"inspection":{"slot":"0"}},
		"summaryTotals":{"inspection":"0"},
		"summaryEdges":{"inspection":{"side":"0"}}
	}`)
	if validateGoreTexWaterproofNumbers(invalid) {
		t.Fatal("expected negative waterproof value to be rejected")
	}
}

func TestValidateGoreTexCentrifugalResults(t *testing.T) {
	valid := json.RawMessage(`{
		"entries":[{"result":"PASS","issueValues":{"toe":"0","heel":"1","medial":"0","lateral":"0"},"otherIssue":"0"}]
	}`)
	if !validateGoreTexCentrifugalResults(valid) {
		t.Fatal("expected valid centrifugal data")
	}

	invalid := json.RawMessage(`{
		"entries":[{"result":"FAIL","issueValues":{"toe":"-1","heel":"0","medial":"0","lateral":"0"},"otherIssue":"0"}]
	}`)
	if validateGoreTexCentrifugalResults(invalid) {
		t.Fatal("expected negative centrifugal value to be rejected")
	}
}
