// internal/pkg/models/types/inventory.go
package types

type Inventory struct {
	Barcode  string `json:"Barcode"`
	USERDATE string `json:"USERDATE"`
	USERID   string `json:"USERID"`
	GSBH     string `json:"GSBH"`
	YN       string `json:"YN"`
}

type ListInventory struct {
	GSBH         string
	Rack         string
	Location     string
	CLBH         string
	Pack         string
	QtyScan      string
	Memo_Article string
	Memo_RY      string
	Qty          string
	RemQty       string
	DQty         string
	Status       string
	Barcode      string
	CFMDate      string
	BuyNos       string

	// --- Kiểm kê theo tháng: theo dõi người thực hiện từng mốc ---
	SyncUserID   string `json:"SyncUserID"`   // người bấm Đồng bộ
	SyncDate     string `json:"SyncDate"`     // ngày bấm Đồng bộ
	Round1UserID string `json:"Round1UserID"` // người quét xác nhận lần 1 (đơn vị kiểm kê)
	Round1Date   string `json:"Round1Date"`   // ngày xác nhận lần 1
	Round2UserID string `json:"Round2UserID"` // người xác nhận lần 2 (kế toán)
	Round2Date   string `json:"Round2Date"`   // ngày xác nhận lần 2
}

type ListInventoryMonth struct {
	GSBH string `json:"GSBH"`

	Barcode string `json:"Barcode"`
	Rack    string `json:"Rack"`
	CLBH    string `json:"CLBH"`
	Pack    string `json:"Pack"`

	Qty     float64 `json:"Qty"` // Qty gốc lấy từ KCRKScan_RFSS
	RemQty  float64 `json:"RemQty"`
	CFMDate string  `json:"CFMDate"`
	QtyScan float64 `json:"QtyScan"`

	Status string `json:"Status"`

	Memo_Article string `json:"Memo_Article"`
	Memo_RY      string `json:"Memo_RY"`
	BuyNos       string `json:"BuyNos"`

	SyncUserID   string `json:"SyncUserID"`
	SyncDate     string `json:"SyncDate"`
	Round1UserID string `json:"Round1UserID"`
	Round1Date   string `json:"Round1Date"`
	Round2UserID string `json:"Round2UserID"`
	Round2Date   string `json:"Round2Date"`
}

// ConfirmRound2Request: dùng khi kế toán xác nhận kiểm kê lần 2 cho 1 barcode cụ thể.
type ConfirmRound2Request struct {
	Barcode string `json:"Barcode"`
	Date    string `json:"Date"` // YearMonth, định dạng yyyy-MM
	GSBH    string `json:"GSBH"`
	USERID  string `json:"USERID"`
}