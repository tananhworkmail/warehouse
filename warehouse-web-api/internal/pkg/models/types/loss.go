package types

// ─────────────────────────────────────────────
// LOSS LIST (code gốc - không sửa)
// ─────────────────────────────────────────────

type LossListRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	FromDate string `form:"from_date"`
	ToDate   string `form:"to_date"`
	ZLBH     string `form:"zlbh"`
	DDBH     string `form:"ddbh"`
	Ywpm     string `json:"ywpm"`
}

type LossItem struct {
	ZLBH    string `json:"zlbh"`
	DDBH    string `json:"ddbh"`
	Article string `json:"article"`
	XieMing string `json:"xie_ming"`
	Pairs   float64 `json:"pairs"`
	XXCC    string `json:"xxcc"`
	Ywpm    string `json:"ywpm"`
}

type LossListResponse struct {
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Items    []LossItem `json:"items"`
}

// ── Save request ─────────────────────────────
// Mode = "IN"  → INSERT LossTemSize_Info (upsert) + LossTemSize_In
// Mode = "OUT" → INSERT LossTemSize_Info (upsert) + LossTemSize_Out
type LossTemSizeSaveRequest struct {
	DDBH    string             `json:"ddbh"     binding:"required"`
	
	MSNV    string             `json:"msnv"     binding:"required"`
	MSNVOut string             `json:"msnv_out"`
	YWPM    string             `json:"ywpm"`
	Article string             `json:"article"`
	XieMing string             `json:"xie_ming"`
	Pairs   float64            `json:"pairs"`
	Mode    string             `json:"mode"     binding:"required"`
	Reason  string             `json:"reason"`
	Sizes   map[string]float64 `json:"sizes"    binding:"required"` // {"36":10.5,"37":5}
	Date    string             `json:"date"`
	Note    string             `json:"note"`
	ProductCode string `json:"product_code"`
}

// ── List / Summary filter ─────────────────────
type LossTemSizeListRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	DDBH     string `form:"ddbh"`
	MSNV     string `form:"msnv"`
	FromDate string `form:"from_date"`
	ToDate   string `form:"to_date"`
}

// ── 1 dòng trả về từ GetList ─────────────────
type LossTemSizeItem struct {
	ID      int64   `json:"id"`
	DDBH    string  `json:"ddbh"`
	MSNV    string  `json:"msnv"`
	YWPM    string  `json:"ywpm"`
	Article string  `json:"article"`
	XieMing string  `json:"xie_ming"`
	Pairs   float64 `json:"pairs"`
	XXCC    string  `json:"xxcc"`
	Mode    string  `json:"mode"`
	Reason  string  `json:"reason"`
	QtyIn   float64 `json:"qty_in"`  // đổi int -> float64
	QtyOut  float64 `json:"qty_out"` // đổi int -> float64
	UserID  string  `json:"userid"`
	Date    string  `json:"date"`
	ProductCode string `json:"product_code"`

}

// ── Phân trang ────────────────────────────────
type LossTemSizeListResponse struct {
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Items    []LossTemSizeItem `json:"items"`
}

// ── Summary theo DDBH ─────────────────────────
type LossTemSizeSummaryItem struct {
	DDBH        string                        `json:"ddbh"`
	MSNV        string                        `json:"msnv"`
	Article     string                        `json:"article"`
	XieMing     string                        `json:"xie_ming"`
	YWPM        string                        `json:"ywpm"`
	Pairs       float64                       `json:"pairs"`
	SizesIn     map[string]float64            `json:"sizes_in"`      
	SizesOut    map[string]float64            `json:"sizes_out"`     
	SizesNet    map[string]float64            `json:"sizes_net"`     
	TotalIn     float64                       `json:"total_in"`
	TotalOut    float64                       `json:"total_out"`
	TotalNet    float64                       `json:"total_net"`
	OutByReason map[string]map[string]float64 `json:"out_by_reason"`
	DATE        string                        `json:"date"`
	ProductCode string `json:"product_code"`

}