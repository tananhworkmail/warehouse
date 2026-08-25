package types

import "time"

type WarehouseKTPScanInCarton struct {
	ScannedCode string     `json:"scannedCode"`
	ResolvedBy  string     `json:"resolvedBy"`
	CartonBar   string     `json:"cartonBar"`
	LPNBarCode  string     `json:"lpnBarCode"`
	DDBH        string     `json:"ddbh"`
	OrderNo     string     `json:"orderNo"`
	CartonNo    *int       `json:"cartonNo"`
	XH          string     `json:"xh"`
	Size        string     `json:"size"`
	Qty         float64    `json:"qty"`
	PackageQty  float64    `json:"packageQty"`
	SB          string     `json:"sb"`
	CKBH        string     `json:"ckbh"`
	KCBH        string     `json:"kcbh"`
	KVBH        string     `json:"kvbh"`
	INCS        int        `json:"incs"`
	INDate      *time.Time `json:"inDate"`
	LastInDate  *time.Time `json:"lastInDate"`
	SGW         float64    `json:"sgw"`
	RGW         *float64   `json:"rgw"`
	Article     string     `json:"article"`
	XieMing     string     `json:"xieMing"`
	Country     string     `json:"country"`
	KHPO        string     `json:"khpo"`
	DepNo       string     `json:"depNo"`
	DepName     string     `json:"depName"`
	INUser      string     `json:"inUser"`
	Factory     string     `json:"factory"`
	FactoryType string     `json:"factoryType"`
	OrderFlag   string     `json:"orderFlag"`
}

type WarehouseKTPScanInCheck struct {
	Status  string `json:"status"`
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

type WarehouseKTPScanInEligibility struct {
	WarehouseKTPScanInCheck
	SessionDuplicate bool `json:"sessionDuplicate"`
	MixedOrder       bool `json:"mixedOrder"`
}

type WarehouseKTPScanInWeightRule struct {
	WarehouseKTPScanInCheck
	StandardWeight *float64 `json:"standardWeight"`
	ActualWeight   *float64 `json:"actualWeight"`
	Gram           *float64 `json:"gram"`
	Tolerance      *float64 `json:"tolerance"`
	Difference     *float64 `json:"difference"`
	Comparison     string   `json:"comparison"`
	Overridden     bool     `json:"overridden"`
	OverrideUserID string   `json:"overrideUserId,omitempty"`
}

type WarehouseKTPScanInProductionCheck struct {
	WarehouseKTPScanInCheck
	Enforced         bool    `json:"enforced"`
	Size             string  `json:"size"`
	CompletedQty     float64 `json:"completedQty"`
	AlreadyInQty     float64 `json:"alreadyInQty"`
	SessionQty       float64 `json:"sessionQty"`
	CurrentCartonQty float64 `json:"currentCartonQty"`
	AfterScanQty     float64 `json:"afterScanQty"`
}

type WarehouseKTPScanInMaterialCheck struct {
	WarehouseKTPScanInCheck
	MaterialOrder        string  `json:"materialOrder"`
	OrderPairs           float64 `json:"orderPairs"`
	CoveredPairs         float64 `json:"coveredPairs"`
	AlreadyInQty         float64 `json:"alreadyInQty"`
	SessionQty           float64 `json:"sessionQty"`
	CurrentCartonQty     float64 `json:"currentCartonQty"`
	AfterScanQty         float64 `json:"afterScanQty"`
	RequiredMaterialRows int     `json:"requiredMaterialRows"`
	IssuedMaterialRows   int     `json:"issuedMaterialRows"`
	MissingMaterialRows  int     `json:"missingMaterialRows"`
	LimitingFactory      string  `json:"limitingFactory"`
}

type WarehouseKTPScanInInboxItem struct {
	InboxBar    string `json:"inboxBar"`
	Size        string `json:"size"`
	ExpectedQty int    `json:"expectedQty"`
	ScannedQty  int    `json:"scannedQty"`
}

type WarehouseKTPScanInInboxCheck struct {
	WarehouseKTPScanInCheck
	Enforced     bool                          `json:"enforced"`
	ExpectedQty  int                           `json:"expectedQty"`
	ScannedQty   int                           `json:"scannedQty"`
	RemainingQty int                           `json:"remainingQty"`
	Items        []WarehouseKTPScanInInboxItem `json:"items"`
}

type WarehouseKTPScanInCartonResult struct {
	Carton          WarehouseKTPScanInCarton          `json:"carton"`
	Eligibility     WarehouseKTPScanInEligibility     `json:"eligibility"`
	WeightRule      WarehouseKTPScanInWeightRule      `json:"weightRule"`
	ProductionCheck WarehouseKTPScanInProductionCheck `json:"productionCheck"`
	MaterialCheck   WarehouseKTPScanInMaterialCheck   `json:"materialCheck"`
	InboxCheck      WarehouseKTPScanInInboxCheck      `json:"inboxCheck"`
	CanScanIn       bool                              `json:"canScanIn"`
	ReadOnly        bool                              `json:"readOnly"`
}

type WarehouseKTPScanInProductionLine struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Factory string `json:"factory"`
}

type WarehouseKTPScanInWarehouse struct {
	Code string `json:"code"`
}
