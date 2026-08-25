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

type WarehouseKTPInventoryAge struct {
	DDBH           string  `json:"ddbh" gorm:"column:ddbh"`
	Qty            float64 `json:"qty" gorm:"column:qty"`
	RackCode       string  `json:"rackCode" gorm:"column:rack_code"`
	LastInDate     *string `json:"lastInDate" gorm:"column:last_in_date"`
	StorageAgeDays *int    `json:"storageAgeDays" gorm:"column:storage_age_days"`
}

type WarehouseKTPSearchResult struct {
	KHPO           string  `json:"khpo" gorm:"column:khpo"`
	DDBH           string  `json:"ddbh" gorm:"column:ddbh"`
	Article        string  `json:"article" gorm:"column:article"`
	RackCode       string  `json:"rackCode" gorm:"column:rack_code"`
	CartonBar      string  `json:"cartonBar" gorm:"column:carton_bar"`
	CartonNo       string  `json:"cartonNo" gorm:"column:carton_no"`
	Qty            float64 `json:"qty" gorm:"column:qty"`
	HasInbound     bool    `json:"hasInbound" gorm:"column:has_inbound"`
	HasRecycle     bool    `json:"hasRecycle" gorm:"column:has_recycle"`
	HasInspection  bool    `json:"hasInspection" gorm:"column:has_inspection"`
	OtherStatus    string  `json:"otherStatus,omitempty" gorm:"column:other_status"`
	LastInDate     *string `json:"lastInDate" gorm:"column:last_in_date"`
	StorageAgeDays *int    `json:"storageAgeDays" gorm:"column:storage_age_days"`
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
