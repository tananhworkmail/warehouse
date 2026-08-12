package types

import "encoding/json"

type GoreTexFormListItem struct {
	FormType         string   `json:"formType"`
	Title            string   `json:"title"`
	RecordKey        string   `json:"recordKey"`
	Line             string   `json:"line,omitempty"`
	StyleName        string   `json:"styleName,omitempty"`
	InspectionDate   string   `json:"inspectionDate,omitempty"`
	TestDates        []string `json:"testDates,omitempty"`
	ImprovementDates []string `json:"improvementDates,omitempty"`
	AnalysisID       uint64   `json:"analysisId,omitempty"`
	CreatedAt        string   `json:"createdAt"`
	UpdatedAt        string   `json:"updatedAt"`
}

type GoreTexFormDetail struct {
	FormType       string          `json:"formType"`
	Line           string          `json:"line,omitempty"`
	StyleName      string          `json:"styleName,omitempty"`
	InspectionDate string          `json:"inspectionDate,omitempty"`
	AnalysisID     uint64          `json:"analysisId,omitempty"`
	Data           json.RawMessage `json:"data"`
	CreatedAt      string          `json:"createdAt"`
	UpdatedAt      string          `json:"updatedAt"`
}

type GoreTexSubmitResult struct {
	FormType  string `json:"formType"`
	RecordKey string `json:"recordKey"`
	Created   bool   `json:"created"`
}
