package types

type ScanGetBarcode struct {
	SCNO         string
	CLBH         string
	YWPM         string
	ZSYWJC       string
	Memo_RY      string
	Memo_Article string
	Qty          float32
	PACK1        string
	Pack         float64
	Memo         string
	DWBH         string
	CGNO         string
	QRCode       string
	FIFO         string
	LotNO        string
	IDCheck      string
	CFMDate      string
	MAKE         string
	File_Name    string
}

type ScanGetLLNO struct {
	LLNO         string
	CLBH         string
	SCBH         string
	DFL          string
	YWPM         string
	DWBH         string
	Qty          float32
	YWSM         string
	SKU          string
	BLSB         string
	SCANUSERID   string
	SCANUSERDATE string
	Allqty       float64
	LLQty        float64
	TempQty      float64
	Checkry      float64
}
type ScanDeliveryLabel struct {
	RemQty  float64
	CLBH    string
	CFMDATE string
}

type RemQtyCLBH struct {
	RemQty float64
	CLBH   string
}
type Scanqty struct {
	SCANQty float64
}

type ScanMaterialInput struct {
	SCNO         string
	CGNO         string
	ZSBH         string
	ZSYWJC       string
	GSBH         string
	CKBH         string
	CLBH         string
	YWPM         string
	DWBH         string
	Qty          float32
	RemQty       float32
	DQty         float32
	TAGID        string
	Pack         float64
	Memo_RY      string
	Memo_Article string
	MAKE         string
	Barcode      string
	MEMO         string
	DOCNO        string
	IDCheck      string
	CFMDATE      string
	CFMID        string
	File_Name    string // Tên file báo cáo IMQC từ MaterialQCCheck
}

type ScanMaterialInput1 struct {
	SCNO         string
	CGNO         string
	ZSBH         string
	ZSYWJC       string
	GSBH         string
	CKBH         string
	CLBH         string
	YWPM         string
	DWBH         string
	Qty          float32
	RemQty       float32
	DQty         float32
	TAGID        string
	Pack         float64
	Memo_RY      string
	Memo_Article string
	MAKE         string
	Barcode      string
	MEMO         string
	DOCNO        string
	IDCheck      string
	CFMDATE1     string
	CFMID        string
}
type ScanMaterialPrint struct {
	SCNO         string
	CLBH         string
	YWPM         string
	ZSYWJC       string
	Memo_RY      string
	Memo_Article string
	Qty          float32
	PACK1        string
	Pack         float64
	DQty         float64
	Memo         string
	DWBH         string
	CGNO         string
	QRCode       string
	FIFO         string
	LotNO        string
	IDCheck      string
}

type InsertKCLLScan_RFSS struct {
	TableName string
	TNO       string
	CLBH      string
	OldID     string
	DelID     string
	DelDate   string
}

type GetKCLLScan_RFSS struct {
	LLNO         string
	SCANQty      float64
	DFL          string
	SCBH         string
	SCANUSERID   string
	SCANUSERDATE string
	CLBH         string
}
type GetKCLLScan_RFSSS struct {
	LLNO         string
	SCANQty      float64
	DFL          string
	SCBH         string
	SCANUSERID   string
	SCANUSERDATE string
	CLBH         string
	LLNO_OWE     string
}
type SelectKCLLScan_RFSS struct {
	LLNO         string
	CLBH         string
	DFL          string
	SCBH         string
	Barcode      string
	SCANQty      float64
	SCANUSERID   string
	SCANUSERDATE string
	LLNO_OWE     string
}

type DeliveryPayExPort struct {
	LLNO_OWE     string
	LLNO         string
	CLBH         string
	DFL          string
	SCBH         string
	Barcode      string
	SCANQty      float64
	SCANUSERID   string
	SCANUSERDATE string
	YWPM         string
	DWBH         string
	SKU          string
}
