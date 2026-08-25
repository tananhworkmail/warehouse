package types

import "time"

type WarehouseKTPScanShipWeightCheck struct {
	Enabled          bool     `json:"enabled"`
	Passed           bool     `json:"passed"`
	RequiresOverride bool     `json:"requiresOverride"`
	Reference        string   `json:"reference"`
	ReferenceWeight  *float64 `json:"referenceWeight"`
	OutgoingWeight   *float64 `json:"outgoingWeight"`
	Tolerance        float64  `json:"tolerance"`
	Difference       *float64 `json:"difference"`
	Message          string   `json:"message"`
}

type WarehouseKTPScanShipCarton struct {
	ScannedCode        string                          `json:"scannedCode"`
	ResolvedBy         string                          `json:"resolvedBy"`
	CartonBar          string                          `json:"cartonBar"`
	LPNBarCode         string                          `json:"lpnBarCode"`
	DDBH               string                          `json:"ddbh"`
	CartonNo           int                             `json:"cartonNo"`
	XH                 string                          `json:"xh"`
	Qty                float64                         `json:"qty"`
	SB                 string                          `json:"sb"`
	StandardWeight     float64                         `json:"standardWeight"`
	ReceivedWeight     float64                         `json:"receivedWeight"`
	OutgoingWeight     float64                         `json:"outgoingWeight"`
	ExecuteDate        *time.Time                      `json:"executeDate"`
	ContainerNo        string                          `json:"containerNo"`
	PlateID            string                          `json:"plateId"`
	HasShippingPlan    bool                            `json:"hasShippingPlan"`
	PlannedContainer   string                          `json:"plannedContainer"`
	PlannedPlateID     string                          `json:"plannedPlateId"`
	EligibilityPassed  bool                            `json:"eligibilityPassed"`
	EligibilityMessage string                          `json:"eligibilityMessage"`
	PlanPassed         bool                            `json:"planPassed"`
	PlanMessage        string                          `json:"planMessage"`
	WeightCheck        WarehouseKTPScanShipWeightCheck `json:"weightCheck"`
	CanScan            bool                            `json:"canScan"`
	Details            []WarehouseKTPScanOutDetail     `json:"details"`
}

type WarehouseKTPScanShipSubmitResult struct {
	UpdatedCount int                          `json:"updatedCount"`
	Cartons      []WarehouseKTPScanShipCarton `json:"cartons"`
}
