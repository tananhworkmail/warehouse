package types

import "time"

type WarehouseKTPScanOutCarton struct {
	ScannedCode string                      `json:"scannedCode"`
	ResolvedBy  string                      `json:"resolvedBy"`
	CartonBar   string                      `json:"cartonBar"`
	LPNBarCode  string                      `json:"lpnBarCode"`
	DDBH        string                      `json:"ddbh"`
	OrderNo     string                      `json:"orderNo"`
	XH          string                      `json:"xh"`
	Qty         float64                     `json:"qty"`
	SB          string                      `json:"sb"`
	OutCount    int                         `json:"outCount"`
	OutDate     *time.Time                  `json:"outDate"`
	ExecuteDate *time.Time                  `json:"executeDate"`
	Article     string                      `json:"article"`
	XieMing     string                      `json:"xieMing"`
	Country     string                      `json:"country"`
	KCBH        string                      `json:"kcbh"`
	KVBH        string                      `json:"kvbh"`
	Eligible    bool                        `json:"eligible"`
	Message     string                      `json:"message"`
	Details     []WarehouseKTPScanOutDetail `json:"details"`
}

type WarehouseKTPScanOutDetail struct {
	OrderNo   string  `json:"orderNo"`
	SKU       string  `json:"sku"`
	ModelName string  `json:"modelName"`
	XH        string  `json:"xh"`
	OrderSize string  `json:"orderSize"`
	ModelSize string  `json:"modelSize"`
	Qty       float64 `json:"qty"`
	From      int     `json:"from"`
	To        int     `json:"to"`
}
