package request

type WarehouseKTPScanRequest struct {
	ScanCode string `json:"scanCode"`
	Replace  bool   `json:"replace"`
}

type WarehouseKTPMoveOrderRequest struct {
	CartonBar   string `json:"cartonBar"`
	NewRackCode string `json:"newRackCode"`
}

type WarehouseKTPScanInSubmitRequest struct {
	CartonBar            string   `json:"cartonBar" binding:"required"`
	ActualWeight         float64  `json:"actualWeight" binding:"required,gt=0"`
	Warehouse            string   `json:"warehouse"`
	ProductionLine       string   `json:"productionLine" binding:"required"`
	Storage              string   `json:"storage" binding:"required"`
	Location             string   `json:"location"`
	ToSX                 string   `json:"toSx"`
	InUser               string   `json:"inUser"`
	SessionCartons       []string `json:"sessionCartons"`
	ScanMode             string   `json:"scanMode"`
	InboxScans           []string `json:"inboxScans"`
	WeightOverride       bool     `json:"weightOverride"`
	WeightOverrideUserID string   `json:"weightOverrideUserId"`
}

type WarehouseKTPScanOutSubmitRequest struct {
	CartonBar string `json:"cartonBar" binding:"required"`
	Mode      string `json:"mode" binding:"required,oneof=recycle inspection"`
}

type WarehouseKTPScanShipCartonRequest struct {
	CartonBar      string  `json:"cartonBar" binding:"required"`
	OutgoingWeight float64 `json:"outgoingWeight"`
	WeightOverride bool    `json:"weightOverride"`
}

type WarehouseKTPScanShipSubmitRequest struct {
	ContainerNo string                              `json:"containerNo" binding:"required"`
	PlateID     string                              `json:"plateId" binding:"required"`
	Cartons     []WarehouseKTPScanShipCartonRequest `json:"cartons" binding:"required,len=1,dive"`
}
