package types

// HumidityAlertLog — ánh xạ bảng humidity_alert_log
// Bảng này lưu lịch sử mỗi lần Andon bị kích hoạt do độ ẩm vượt ngưỡng
type HumidityAlertLog struct {
	ID              int64   `gorm:"column:id"               json:"id"`
	DeviceName      string  `gorm:"column:device_name"      json:"device_name"`
	Hum             float64 `gorm:"column:hum"              json:"hum"`
	Tem             float64 `gorm:"column:tem"              json:"tem"`
	AlertType       string  `gorm:"column:alert_type"       json:"alert_type"`       // "HIGH" | "LOW"
	AndonTriggered  bool    `gorm:"column:andon_triggered"  json:"andon_triggered"`
	TriggeredAt     string  `gorm:"column:triggered_at"     json:"triggered_at"`
	AndonOffBy      string  `gorm:"column:andon_off_by"     json:"andon_off_by"`
	AndonOffAt      string  `gorm:"column:andon_off_at"     json:"andon_off_at"`
}