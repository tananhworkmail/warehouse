package request

import "encoding/json"

type GoreTexWaterproofSubmitRequest struct {
	Line           string          `json:"line" binding:"required"`
	StyleName      string          `json:"styleName" binding:"required"`
	InspectionDate string          `json:"inspectionDate" binding:"required"`
	Data           json.RawMessage `json:"data" binding:"required"`
	IsEdit         bool            `json:"isEdit"`
}

type GoreTexCentrifugalSubmitRequest struct {
	InspectionDate string          `json:"inspectionDate" binding:"required"`
	Data           json.RawMessage `json:"data" binding:"required"`
	IsEdit         bool            `json:"isEdit"`
}

type GoreTexAnalysisSubmitRequest struct {
	AnalysisID uint64          `json:"analysisId"`
	Data       json.RawMessage `json:"data" binding:"required"`
	IsEdit     bool            `json:"isEdit"`
}
