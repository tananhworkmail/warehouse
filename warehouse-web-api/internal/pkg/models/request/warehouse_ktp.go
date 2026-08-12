package request

type WarehouseKTPScanRequest struct {
	ScanCode string `json:"scanCode"`
	Replace  bool   `json:"replace"`
}

type WarehouseKTPMoveOrderRequest struct {
	CartonBar   string `json:"cartonBar"`
	NewRackCode string `json:"newRackCode"`
}
