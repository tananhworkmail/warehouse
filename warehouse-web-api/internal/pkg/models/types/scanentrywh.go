package types

type RSTList struct {
	PDFilePage string
}
type MaterialQCcheck struct {
	File_Name_Lab string
}

type KCRKScan_RFSStype struct {
	ARTICLE string
	IMGName string
}

type Ordertype struct {
	SCNO     string
	CGNO     string
	KHPO     string
	ZLBH     string
	Qty      float32
	ARTICLE  string
	Pairs    float32
	ShipDate string
	Assemble string
	ADate    string
	DepName  string
	DepMemo  string
	Xuong    string
	YWSM     string
}
