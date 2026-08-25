package types

import "time"

type RackKVT struct {
	MAKE string `json:"make"`
}
type RackTotalColumn struct {
	MAKE3       string  `json:"make3"`
	TotalColumn float64 `json:"totalcolumn"`
	Limit       float64 `json:"limit"`
}
type RackTotalTang struct {
	MAKE5     string  `json:"make5"`
	TotalTang float64 `json:"totaltang"`
}
type TonkhoRack struct {
	MAKE       string  `json:"make"`
	TonkhoRack float64 `json:"tonkhorack"`
}
type WarehouseRackInformationKVT struct {
	CLBH         string  `json:"clbh"`
	SCNO         string  `json:"scno"`
	CGNO         string  `json:"cgno"`
	ZSBH         string  `json:"zsbh"`
	ZSYWJC       string  `json:"zsywjc"`
	GSBH         string  `json:"gsbh"`
	YWPM         string  `json:"ywpm"`
	DWBH         string  `json:"dwbh"`
	Qty          float64 `json:"qty"`
	RemQty       float64 `json:"remqty"`
	DQty         float64 `json:"dqty"`
	TAGID        string  `json:"tagid"`
	Pack         string  `json:"pack"`
	Memo_RY      string  `json:"memo_ry"`
	Memo_Article string  `json:"memo_article"`
	KCBH         string  `json:"kcbh"`
	Barcode      string  `json:"barcode"`
	MEMO         string  `json:"memo"`
	DOCNO        string  `json:"docno"`
	IDCheck      string  `json:"idcheck"`
	CFMID        string  `json:"cfmid"`
	CFMDATE      string  `json:"cfmdate"`
	Is_over_180  int     `json:"is_over_180"`
}
type Warehouse3DayRackInformationKVT struct {
	CLBH string `json:"clbh"`
	SCNO string `json:"scno"`
	MAKE string `json:"make"`
}
type SearchByCLBH struct {
	CLBH         string  `json:"clbh"`
	SCNO         string  `json:"scno"`
	CGNO         string  `json:"cgno"`
	ZSBH         string  `json:"zsbh"`
	ZSYWJC       string  `json:"zsywjc"`
	GSBH         string  `json:"gsbh"`
	YWPM         string  `json:"ywpm"`
	DWBH         string  `json:"dwbh"`
	Qty          float64 `json:"qty"`
	RemQty       float64 `json:"remqty"`
	DQty         float64 `json:"dqty"`
	TAGID        string  `json:"tagid"`
	Pack         string  `json:"pack"`
	Memo_RY      string  `json:"memo_ry"`
	Memo_Article string  `json:"memo_article"`
	KCBH         string  `json:"kcbh"`
	Barcode      string  `json:"barcode"`
	MEMO         string  `json:"memo"`
	DOCNO        string  `json:"docno"`
	IDCheck      string  `json:"idcheck"`
	CFMID        string  `json:"cfmid"`
	CFMDATE      string  `json:"cfmdate"`
}
type SearchN31 struct {
	BUYNO    string  `json:"buyno"`
	KHPO     string  `json:"khpo"`
	ZLBH     string  `json:"zlbh"`
	Article  string  `json:"article"`
	XieMing  string  `json:"xieming"`
	Pairs    float64 `json:"pairs"`
	CLBH     string  `json:"clbh"`
	YWPM     string  `json:"ywpm"`
	DWBH     string  `json:"dwbh"`
	CQDH     string  `json:"cqdh"`
	CLSL     float64 `json:"clsl"`
	CGQty    float64 `json:"cgqty"`
	RKQty    float64 `json:"rkqty"`
	UseStock float64 `json:"usestock"`
	CGNO     string  `json:"cgno"`
	DDGB     string  `json:"ddgb"`
	ZSBH     string  `json:"zsbh"`
	ZSYWJC   string  `json:"zsywjc"`
	Invoice  string  `json:"invoice"`
	MAKE     string  `json:"make"`
}
type GetTempHumidity struct {
	DeviceName string    `json:"DeviceName"`
	DeviceAddr string    `json:"DeviceAddr"`
	Hum        float64   `json:"Hum"`
	Tem        float64   `json:"Tem"`
	RecordTime time.Time `json:"RecordTime"`
}
type GetHumidityAlert struct {
	DeviceName string  `json:"DeviceName"`
	Hum        float64 `json:"Hum"`
	Temp       float64 `json:"Temp"`
	Alert_Date string  `json:"Alert_Date"`
	Alert_Time string  `json:"Alert_Time"`
	RecordTime string  `json:"RecordTime"`
}

// ─── MỚI THÊM ────────────────────────────────────────────────
// TbHistory — ánh xạ bảng tbhistory trong DB rkmonitor
// Cột: Id, DeviceName, Hum, Tem, RecordTime, alert_date, alert_time, CreatedAt
type TbHistory struct {
	Id         int64   `gorm:"column:Id"         json:"Id"`
	DeviceName string  `gorm:"column:DeviceName" json:"DeviceName"`
	Hum        float64 `gorm:"column:Hum"        json:"Hum"`
	Tem        float64 `gorm:"column:Tem"        json:"Tem"`
	RecordTime string  `gorm:"column:RecordTime" json:"RecordTime"`
	AlertDate  string  `gorm:"column:alert_date" json:"Alert_Date"`
	AlertTime  string  `gorm:"column:alert_time" json:"Alert_Time"`
	CreatedAt  string  `gorm:"column:CreatedAt"  json:"CreatedAt"`
}

// HumidityLiveStatus — kết quả kiểm tra realtime 1 thiết bị
type HumidityLiveStatus struct {
	DeviceName string  `json:"DeviceName"`
	Hum        float64 `json:"Hum"`
	Tem        float64 `json:"Tem"`
	RecordTime string  `json:"RecordTime"`
	IsAlert    bool    `json:"IsAlert"`   // true khi Hum < 45 hoặc Hum > 60
	AlertType  string  `json:"AlertType"` // "HIGH" | "LOW" | ""
	IsStale    bool    `json:"IsStale"`   // true nếu data cũ hơn 90 phút
}
type Warehouse3And180DayKVT struct {
	MAKE   string `json:"make"`
	Day3   int    `json:"day3"`
	Day180 int    `json:"day180"`
}
