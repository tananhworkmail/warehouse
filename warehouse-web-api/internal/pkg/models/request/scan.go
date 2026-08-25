package request

type ScanGetBarcodeRequest struct {
	Barcode string `form:"barcode" json:"barcode"`
}
type UpdateCFMIDCodeRequest struct {
	Barcode string `form:"barcode" json:"barcode"`
	CfmId   string `form:"CFMID" json:"CFMID"`
	GSBH    string `form:"CFMID" json:"GSBH"`
	MAKE    string `form:"MAKE" json:"MAKE"`
}
type ScanGetKCLLRequest struct {
	LLNO []string `form:"LLNO" json:"LLNO"`
}

type ScanDeliveryLabelQuantity struct {
	Barcode  string `form:"barcode" json:"barcode"`
	Quantity string `form:"quantity" json:"quantity"`
	MAKE     string `form:"MAKE" json:"MAKE"`
	USERID   string
}
type ScanDeliveryLabelCLBHSelected struct {
	LLNO string `form:"LLNO" json:"LLNO"`
	CLBH string `form:"CLBH" json:"CLBH"`
	DFL  string `form:"DFL" json:"DFL"`
	SCBH string `form:"SCBH" json:"SCBH"`
	YWPM string `form:"YWPM" json:"YWPM"`
}
type ScanDeliveryLabelCLBH struct {
	Barcode string `form:"barcode" json:"barcode"`
	USERID  string
	Data    []ScanDeliveryLabelCLBHSelected `form:"Data" json:"Data" binding:"required"`
}
type ScanMaterialIpRequest struct {
	GSBH     string `form:"GSBH" json:"GSBH"`
	SCBH     string `form:"SCBH" json:"SCBH"`
	CGNO     string `form:"CGNO" json:"CGNO"`
	ZSDH     string `form:"ZSDH" json:"ZSDH"`
	ZSYWJC   string `form:"zsywjc" json:"zsywjc"`
	CLBH     string `form:"CLBH" json:"CLBH"`
	YWBM     string `form:"YWBM" json:"YWBM"`
	SCFMDATE string
	ECFMDATE string
	MAKE     string `json:"MAKE"`
	Barcode  string `json:"Barcode"`
}
type ScanMaterialPrintRequest struct {
	GSBH     string `form:"GSBH" json:"GSBH"`
	SCBH     string `form:"SCBH" json:"SCBH"`
	CGNO     string `form:"CGNO" json:"CGNO"`
	ZSDH     string `form:"ZSDH" json:"ZSDH"`
	ZSYWJC   string `form:"zsywjc" json:"zsywjc"`
	MAKE     string `json:"MAKE"`
	Barcode  string `json:"Barcode"`
	CLBH     string `form:"CLBH" json:"CLBH"`
	YWBM     string `form:"YWBM" json:"YWBM"`
	SCFMDATE string
	ECFMDATE string
}
type GetKCLLScan_RFSS struct {
	Barcode string `form:"barcode" json:"barcode"`
}
type GetKCLLScan_RFSSS struct {
	Barcode string `form:"barcode" json:"barcode"`
}
type UpdateKCLLScan_RFSSRequest struct {
	LLNO          string
	SCBH          string
	CLBH          string
	Barcode       string
	DFL           string
	SCANQty       float64
	ScanQty_input float64
	SCANUSERID    string
	USERID        string
}
type UpdateKCRKScan_RFSSRequest struct {
	Barcode string
}

type UpdateKCLLScan_RFSSSRequest struct {
	KCLLScan_RFSS string
	LLNO_OWE      string
	LLNO          string
	SCBH          string
	CLBH          string
	Barcode       string
	SCANQty       float64
	ScanQty_input float64
	SCANUSERID    string
	USERID        string
}

type DeliveryPayExPortRequest struct {
	SCBH     string
	CLBH     string
	YWPM     string
	SCFMDATE string
	ECFMDATE string
}

type DeliveryPayKCLLScan_RFSS struct {
	LLNO          string
	LLNO_OWE      string
	SCBH          string
	CLBH          string
	Barcode       string
	DFL           string
	SCANQty       float64
	ScanQty_input float64
	SCANUSERID    string
	USERID        string
}
type DeliveryDblClick struct {
	LLNO_KCLLS     string  `json:"LLNO_KCLLS"`
	LLNO_OWE_LLSSS string  `json:"LLNO_OWE_LLSSS"`
	CLBH_KCLL      string  `json:"CLBH_KCLL"`
	SCBH_KCLL      string  `json:"SCBH_KCLL"`
	DFL_KCLL       string  `json:"DFL_KCLL"`
	USERID         string  `json:"USERID"`
	LLQty_KCLL     float64 `json:"LLQty_KCLL"`
	Qty_KCLL       float64 `json:"Qty_KCLL"`
}
