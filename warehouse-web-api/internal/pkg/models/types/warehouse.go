package types

type GetWarehouseList struct {
	ZLBH     string `json:"zlbh"`
	CLBH     string `json:"clbh"`
	ZSBH     string `json:"zsbh"`
	COLOR    string `json:"color"`
	TYPE     string `json:"type"`
	STOCKQTY int    `json:"stockqty"`
	SIZE     string `json:"size"`
	QTY      int    `json:"qty"`
}
type WarehouseRackInformation struct {
	DDBH    string `json:"ddbh"`
	YSBH    string `json:"ysbh"`
	XieMing string `json:"XieMing"`
	XXCC    string `json:"XXCC"`
	QtyIn   string `json:"QtyIn"`
	QtyOut  string `json:"QtyOut"`
	TonKho  string `json:"TonKho"`
	Codebar string `json:"Codebar"`
	Ma_Ke   string `json:"Ma_Ke"`
	Name    string `json:"Name"`
	DonVi   string `json:"DonVi"`
}
type Warehouse3DayRackInformation struct {
	RY    string `json:"ry"`
	Ma_Ke string `json:"Ma_Ke"`
}
type RackTonKho struct {
	MA_KE      string  `gorm:"column:Ma_Ke" json:"MA_KE"`
	TonKhoTang float64 `gorm:"column:TonKhoTang" json:"TonKhoTang"`
}
type RackTonKhoKe struct {
	MA_KE    string  `json:"MA_KE"`
	TonKhoKe float64 `json:"TonKhoKe"`
	Limit    float64 `json:"Limit"`
}
type RackTonKhoTang struct {
	MA_KE      string  `json:"MA_KE"`
	TonKhoTang float64 `json:"TonKhoKe"`
}
type SearchByDDBH struct {
	DDBH    string `json:"ddbh"`
	XieMing string `json:"XieMing"`
	XXCC    string `json:"xxcc"`
	QtyIn   string `json:"qtyin"`
	QtyOut  string `json:"qtyout"`
	TonKho  string `json:"tonkho"`
	Ma_Ke   string `json:"make"`
}
type SearchTotalDDBH struct {
	TotalQtyIn    string `json:"totalqtyin"`
	TotalQtyOut   string `json:"totalqtyout"`
	TotalOrderQty string `json:"totalorderqty"`
	// TongTonKho  string `json:"tongtonkho"`
	// Remaining   string `json:"remaining "`
}
type TraceOrderExistenceStatus struct {
	DDBH string // 訂單編號
	KHPO string // 客戶PO
}
type TraceProductionInformation struct {
	SCBH     string // 訂單編號
	KHPO     string // 客戶PO
	GXLB     string // 工段類別
	Qty      string // 報工數量
	USERDATE string // 報工日期
}

type TraceOrderHead struct {
	ZLBH     string
	XieMing  string
	ShipDate string
	KHPO     string
	ARTICLE  string
	BZCC     string
	Kfjc     string
	XieXing  string
	SheHao   string
	XTMH     string
	Ywsm     string
	YSSM     string
	DAOMH    string
	IMGName  string
}

type TraceOrderSize struct {
	XXCC string  // 訂單編號
	QTY  float32 // 客戶PO
}

type TraceOrderDetail struct {
	BWBH   string
	BWMC   string
	CLBH   string
	CLMC   string
	DWBH   string
	ZSYWJC string
	CLSL   float32
	SumQty float32
	CQDH   string
}

type TraceCGNOAll struct {
	CGNO     string
	GSBH     string
	CGDate   string
	DH       string
	CZ       string
	ZSYWJC   string
	CLBH     string
	Qty      float32
	YQDate   string
	Ywpm     string
	DWBH     string
	Order_ID int32
}

type TraceCGNOHead struct {
	CGNO     string
	GSBH     string
	CGDate   string
	DH       string
	CZ       string
	ZSYWJC   string
	Order_ID int32
}

type TraceCGNODetail struct {
	CLBH   string
	Qty    float32
	YQDate string
	Ywpm   string
	DWBH   string
}

type TraceOrderTmpInfo struct {
	Cghzbh  string
	Dgdh    string
	Clbh    string
	Ywpm    string
	Dwbh    string
	Zl_qty  float32
	Zd_date string
	Zsbh    string
	Kfqm    string
	Dh      string
	Cz      string
	Dg_date string
}

type TraceOrderTmpDetail struct {
	Dgdh    string
	Clbh    string
	Ywpm    string
	Dwbh    string
	Zl_qty  float32
	Zd_date string
	Zsbh    string
}

type TraceOrderTmpHead struct {
	Cghzbh  string
	Kfqm    string
	Dh      string
	Cz      string
	Dg_date string
}

type TraceOrderTmpAll struct {
	Head   TraceOrderTmpHead
	Detail []TraceOrderTmpDetail
}

// 物料採購資料
type TraceCGNOInfo struct {
	Title  TraceCGNOHead
	Rows   []TraceCGNODetail
	Orders TraceOrderTmpAll
}

// BOM 資料
type TraceOrderInfo struct {
	Head TraceOrderHead   // 訂單編號
	Size []TraceOrderSize // 客戶PO
	Body []TraceOrderDetail
}
type TraceInvoice struct {
	Head TraceInvoiceHead
	Body []TraceInvoice1
}
type TraceInvoice1 struct {
	Body     TraceInvoiceTitle
	Total    TracePackSum
	PackList []TracePackList
}
type TraceProductHead struct {
	XieMing string
	XieXing string
	ARTICLE string
	YSSM    string
	Qty     string
}

type TraceProductBody struct {
	XXCC   string
	Qty    int32
	INDATE string
	TQty   int32
}

type TraceProductBody2 struct {
	INDATE    string
	CARTONBAR string
	Qty       string
	Sgw       string
	Rgw       string
	TQty      string
}

type TraceProduc struct {
	Items  []TraceProductBody
	Items2 []TraceProductBody2
}

type TraceCustomerHead struct {
	KHDDBH1     string
	BUYER       string
	ShipTo      string
	Customer_PO string
	XieMing     string
	XieXing     string
	ARTICLE     string
	YSSM        string
	Qty         string
	Cp          string
}

type TraceCustomerTitle struct {
	//XXCC Size,Qty Size_Qty,'' Case_Qty,Unit Unit_Case
	Size      string
	Size_Qty  string
	Case_Qty  string
	Unit_Case string
}

type TraceCustomer struct {
	Head TraceCustomerHead
	Body []TraceCustomerTitle
}

type TracePackList struct {
	Number string
	SIZ    string
	QTY    string
	CTS    float32
	PAIRS  float32
	NW     float32
	TNW    float32
	GW     float32
	TGW    float32
	L      float32
	W      float32
	H      float32
}

type TracePackSum struct {
	CTS   float32
	Pairs float32
	NW    float32
	GW    float32
	CBM   float32
}

type TraceInvoiceHead struct {
	INV_NO     string
	Date       string
	Goods      string
	Shipped    string
	FROM_WHERE string
	TO_WHERE   string
}

type TraceInvoiceTitle struct {
	RYNO    string
	Pairs   string
	KHDDBH1 string
	XieMing string
	YSSM    string
	ARTICLE string
}

type TraceServiceAllDAta struct {
	Bom      TraceOrderInfo
	Material []TraceCGNOInfo
	Orders   []TraceOrderTmpAll
	Invoice  TraceInvoice
	Customer TraceCustomer
	Product  TraceProduc
}
