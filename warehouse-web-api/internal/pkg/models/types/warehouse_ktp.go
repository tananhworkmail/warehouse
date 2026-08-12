package types

type WarehouseKTPRack struct {
	RackCode              string   `json:"rackCode"`
	DDBH                  string   `json:"ddbh,omitempty"`
	CurrentCode           *string  `json:"currentCode"`
	Status                string   `json:"status"`
	TotalQty              float64  `json:"totalQty"`
	InboundCartonCount    int      `json:"inboundCartonCount"`
	RecycleCartonCount    int      `json:"recycleCartonCount"`
	InspectionCartonCount int      `json:"inspectionCartonCount"`
	CodebarCount          int      `json:"codebarCount"`
	SampleCodes           []string `json:"sampleCodes,omitempty"`
}

type WarehouseKTPCartonDetail struct {
	CartonNo string  `json:"cartonNo"`
	Qty      float64 `json:"qty"`
}

type WarehouseKTPStatusDetail struct {
	SB          string                     `json:"sb"`
	CartonCount int                        `json:"cartonCount"`
	TotalQty    float64                    `json:"totalQty"`
	Cartons     []WarehouseKTPCartonDetail `json:"cartons"`
}

type WarehouseKTPRackOrder struct {
	DDBH         string                     `json:"ddbh"`
	TotalQty     float64                    `json:"totalQty"`
	CodebarCount int                        `json:"codebarCount"`
	Statuses     []WarehouseKTPStatusDetail `json:"statuses"`
}

type WarehouseKTPRackOrderDetail struct {
	RackCode string                  `json:"rackCode"`
	Orders   []WarehouseKTPRackOrder `json:"orders"`
}

type WarehouseKTPOrderRack struct {
	RackCode     string  `json:"rackCode"`
	TotalQty     float64 `json:"totalQty"`
	CodebarCount int     `json:"codebarCount"`
}

type WarehouseKTPMoveOrderInfo struct {
	ScannedCartonBar string                  `json:"scannedCartonBar"`
	DDBH             string                  `json:"ddbh"`
	TotalQty         float64                 `json:"totalQty"`
	CodebarCount     int                     `json:"codebarCount"`
	CurrentRacks     []WarehouseKTPOrderRack `json:"currentRacks"`
	NewRackCode      string                  `json:"newRackCode,omitempty"`
	UpdatedRows      int64                   `json:"updatedRows,omitempty"`
}
