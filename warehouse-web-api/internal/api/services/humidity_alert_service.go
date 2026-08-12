package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"
	"web-api/internal/realtime"
)


type HumidityAlertService struct {
	*BaseService
	andonMap        map[string]string
	lastProcessedID map[string]int64
}

var HumidityAlert = &HumidityAlertService{
	andonMap: map[string]string{
		"KVT01": "10.20.0.135:1035",
		"KVT02": "10.20.0.136:1036",
		"KVT03": "10.20.0.137:1037",
	},
	lastProcessedID: make(map[string]int64),
}

type latestRecord struct {
	Id         int64   `gorm:"column:Id"`
	DeviceName string  `gorm:"column:DeviceName"`
	Hum        float64 `gorm:"column:Hum"`
	Tem        float64 `gorm:"column:Tem"`
	RecordTime string  `gorm:"column:RecordTime"`
}

// ─── InitBaseline ────────────────────────────────────────────────────────────
// Gọi 1 lần khi server start để lấy Id hiện tại làm baseline.
// Không bật Andon — chỉ ghi nhớ Id để lần sau so sánh.
// Tránh Andon kêu mỗi khi restart server.
func (s *HumidityAlertService) InitBaseline() {
	fmt.Println("[HumidityAlert] InitBaseline — lấy Id baseline, KHÔNG bật Andon")

	db, err := database.TempHumidity_Connection()
	if err != nil {
		fmt.Println("[HumidityAlert] InitBaseline DB error:", err)
		return
	}
	dbInst, _ := db.DB()
	if dbInst != nil {
		defer dbInst.Close()
	}

	query := `
		SELECT Id, DeviceName, Hum, Tem,
		       CONVERT(varchar, RecordTime, 120) AS RecordTime
		FROM (
		    SELECT Id, DeviceName, Hum, Tem, RecordTime,
		           ROW_NUMBER() OVER (PARTITION BY DeviceName ORDER BY Id DESC) AS rn
		    FROM tbhistory_alert
		    WHERE DeviceName IN ('KVT01', 'KVT02', 'KVT03')
		) t
		WHERE rn = 1
	`
	var records []latestRecord
	if err := db.Raw(query).Scan(&records).Error; err != nil {
		fmt.Println("[HumidityAlert] InitBaseline query error:", err)
		return
	}

	for _, r := range records {
		s.lastProcessedID[r.DeviceName] = r.Id
		fmt.Printf("[HumidityAlert] Baseline %s Id=%d Hum=%.1f\n",
			r.DeviceName, r.Id, r.Hum)
	}
}


// Logic:
//   1. Lấy bản ghi mới nhất từng thiết bị từ tbhistory_alert
//   2. So Id với lần trước — nếu KHÔNG ĐỔI thì skip (1 tiếng/lần cập nhật)
//   3. Id mới → ghi log vào humidity_alert_log
//   4. Nếu Hum < 45 hoặc Hum > 60 → bật Andon (đèn đỏ + còi)
type AlertMessage struct {
	DeviceName string  `json:"device_name"`
	Hum        float64 `json:"hum,omitempty"`
	Tem        float64 `json:"tem,omitempty"`
	Type       string  `json:"type"`
	IsAlert    bool    `json:"is_alert,omitempty"`
	AlertType  string  `json:"alert_type,omitempty"`
	Time       string  `json:"time"`
}

func (s *HumidityAlertService) CheckAndTriggerAndon() {
	db, err := database.TempHumidity_Connection()
	if err != nil {
		return
	}
	dbInst, _ := db.DB()
	if dbInst != nil {
		defer dbInst.Close()
	}

	// Truy vấn lấy bản ghi mới nhất của 3 thiết bị
	query := `
		SELECT Id, DeviceName, Hum, Tem, RecordTime
		FROM (
		    SELECT Id, DeviceName, Hum, Tem, RecordTime,
		           ROW_NUMBER() OVER (PARTITION BY DeviceName ORDER BY Id DESC) AS rn
		    FROM tbhistory_alert
		    WHERE DeviceName IN ('KVT01', 'KVT02', 'KVT03')
		) t
		WHERE rn = 1
	`

	var records []latestRecord
	if err := db.Raw(query).Scan(&records).Error; err != nil {
		return
	}

	for _, r := range records {
		prevID := s.lastProcessedID[r.DeviceName]
		
		// Chỉ xử lý nếu có Id mới (dữ liệu bảng 1 vừa thay đổi)
		if r.Id <= prevID {
			continue 
		}

		// LOGIC SO SÁNH: Ngưỡng 40 và 60 cho cả Hum và Tem
		isAlert := false
		alertType := "NORMAL"

		if r.Hum > 60 || r.Tem > 60 {
			alertType = "HIGH"
			isAlert = true
		} else if r.Hum < 40 || r.Tem < 40 {
			alertType = "LOW"
			isAlert = true
		}

		// Cập nhật Id đã xử lý để không trùng lặp
		s.lastProcessedID[r.DeviceName] = r.Id

		// 1. Lưu ngay lập tức vào bảng humidity_alert_log (Hình 2)
		s.writeAlertLog(r.DeviceName, r.Hum, r.Tem, alertType, isAlert)

		// 2. Điều khiển đèn Andon
		if isAlert {
			msg := AlertMessage{
			DeviceName: r.DeviceName,
			Hum:        r.Hum,
			Tem:        r.Tem,
			Type:       alertType,
			Time:       time.Now().Format("2006-01-02 15:04:05"),
		}

			jsonData, _ := json.Marshal(msg)

// 🚀 gửi realtime cho FE
			realtime.AlertHub.Broadcast(jsonData)
			addr := s.andonMap[r.DeviceName]
			fmt.Printf("[ALERT] %s vi phạm ngưỡng! Hum:%.1f, Tem:%.1f. Bật đèn đỏ + còi.\n", r.DeviceName, r.Hum, r.Tem)
			go s.callAndon(addr, 3) // Bật Đèn Đỏ
			go s.callAndon(addr, 4) // Bật Còi
		} else {
			// Nếu trong khoảng an toàn (40-60), không làm gì hoặc có thể tắt đèn nếu muốn
			fmt.Printf("[OK] %s bình thường. Hum:%.1f, Tem:%.1f\n", r.DeviceName, r.Hum, r.Tem)
		}
	}
}


// ─── GetLatestHumidityStatus — API realtime cho frontend ─────────────────────
func (s *HumidityAlertService) GetLatestHumidityStatus() ([]types.HumidityLiveStatus, error) {
	db, err := database.TempHumidity_Connection()
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %w", err)
	}
	dbInst, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer dbInst.Close()

	query := `
		SELECT Id, DeviceName, Hum, Tem,
		       CONVERT(varchar, RecordTime, 120) AS RecordTime
		FROM (
		    SELECT Id, DeviceName, Hum, Tem, RecordTime,
		           ROW_NUMBER() OVER (PARTITION BY DeviceName ORDER BY Id DESC) AS rn
		    FROM tbhistory_alert
		    WHERE DeviceName IN ('KVT01', 'KVT02', 'KVT03')
		) t
		WHERE rn = 1
	`

	var rows []struct {
		Id         int64   `gorm:"column:Id"`
		DeviceName string  `gorm:"column:DeviceName"`
		Hum        float64 `gorm:"column:Hum"`
		Tem        float64 `gorm:"column:Tem"`
		RecordTime string  `gorm:"column:RecordTime"`
	}

	if err := db.Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}

	now := time.Now() // Đã sử dụng biến này ở dưới
	result := make([]types.HumidityLiveStatus, 0, len(rows))

	for _, r := range rows {
		status := types.HumidityLiveStatus{
			DeviceName: r.DeviceName,
			Hum:        r.Hum,
			Tem:        r.Tem,
			RecordTime: r.RecordTime,
		}

		// --- SỬ DỤNG BIẾN 'now' ĐỂ TÍNH ISSTALE ---
		// Nếu dữ liệu trong DB đã quá 90 phút chưa cập nhật mới -> IsStale = true
		recTime, parseErr := time.ParseInLocation("2006-01-02 15:04:05", r.RecordTime, time.Local)
		if parseErr == nil {
			status.IsStale = now.Sub(recTime).Minutes() > 90
		} else {
			status.IsStale = true
		}

		// Cập nhật ngưỡng hiển thị trên giao diện (40 - 60)
		if r.Hum > 60 || r.Tem > 60 {
			status.IsAlert = true
			status.AlertType = "HIGH"
		} else if r.Hum < 40 || r.Tem < 40 {
			status.IsAlert = true
			status.AlertType = "LOW"
		} else {
			status.IsAlert = false
			status.AlertType = "NORMAL"
		}
		result = append(result, status)
	}
	return result, nil
}
// ─── GetAlertLog ──────────────────────────────────────────────────────────────
func (s *HumidityAlertService) GetAlertLog(deviceName string, date string) ([]types.HumidityAlertLog, error) {
	db, err := database.TempHumidity_Connection()
	if err != nil {
		return nil, err
	}
	dbInst, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer dbInst.Close()

	query := `
		SELECT id, device_name, hum, tem, alert_type, andon_triggered,
		       ISNULL(andon_off_by, '')                          AS andon_off_by,
		       ISNULL(CONVERT(varchar, andon_off_at, 120), '')   AS andon_off_at,
		       CONVERT(varchar, triggered_at, 120)               AS triggered_at
		FROM humidity_alert_log
		WHERE 1=1
	`
	args := []interface{}{}
	if deviceName != "" {
		query += " AND device_name = ?"
		args = append(args, deviceName)
	}
	if date != "" {
		query += " AND CONVERT(date, triggered_at) = ?"
		args = append(args, date)
	}
	query += " ORDER BY triggered_at DESC"

	var result []types.HumidityAlertLog
	if err := db.Raw(query, args...).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// ─── TurnOffAndon ────────────────────────────────────────────────────────────
func (s *HumidityAlertService) TurnOffAndon(deviceName string, offBy string) error {
	addr, ok := s.andonMap[deviceName]
	if !ok {
		return fmt.Errorf("không tìm thấy Andon cho thiết bị %s", deviceName)
	}

	url := fmt.Sprintf("http://%s/andon/led/0/off", addr)
	client := &http.Client{Timeout: 3 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		// ESP chết vẫn cho qua
		fmt.Printf("[TurnOffAndon] ESP %s không phản hồi: %v — vẫn ghi DB\n", deviceName, err)
	} else {
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("[TurnOffAndon] ESP %s trả về status %d\n", deviceName, resp.StatusCode)
		}
	}

	// ✅ update DB
	s.markAndonOff(deviceName, offBy)

	// 🚀 broadcast SSE cho toàn bộ client
	msg := AlertMessage{
		DeviceName: deviceName,
		Type:       "ANDON_OFF",
		IsAlert:    false,
		AlertType:  "NORMAL",
		Time:       time.Now().Format("2006-01-02 15:04:05"),
	}

	jsonData, err := json.Marshal(msg)
	if err == nil {
		realtime.AlertHub.Broadcast(jsonData)
	} else {
		fmt.Println("[TurnOffAndon] marshal lỗi:", err)
	}

	fmt.Printf("[HumidityAlert] Tắt Andon %s bởi %s\n", deviceName, offBy)

	return nil
}
// ─── TurnOnAndon ─────────────────────────────────────────────────────────────
func (s *HumidityAlertService) TurnOnAndon(deviceName string, ledID int) error {
	addr, ok := s.andonMap[deviceName]
	if !ok {
		return fmt.Errorf("không tìm thấy Andon cho thiết bị %s", deviceName)
	}

	var url string
	if ledID == 0 {
		url = fmt.Sprintf("http://%s/andon/led/0/off", addr)
	} else {
		url = fmt.Sprintf("http://%s/andon/led/%d/toggle", addr, ledID)
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("ESP %s không phản hồi: %w", deviceName, err)
	}
	defer resp.Body.Close()
	return nil
}

// ─── GetAndonStatus ───────────────────────────────────────────────────────────
func (s *HumidityAlertService) GetAndonStatus(deviceName string) (string, error) {
	addr, ok := s.andonMap[deviceName]
	if !ok {
		return "", fmt.Errorf("không tìm thấy Andon cho %s", deviceName)
	}
	url := fmt.Sprintf("http://%s/andon/check", addr)
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "GREEN: OFF\nYELLOW: OFF\nRED: OFF\nRING: OFF", nil
	}
	defer resp.Body.Close()
	var buf [512]byte
	n, _ := resp.Body.Read(buf[:])
	return string(buf[:n]), nil
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func (s *HumidityAlertService) callAndon(addr string, ledID int) {
	var url string
	if ledID == 0 {
		url = fmt.Sprintf("http://%s/andon/led/0/off", addr)
	} else {
		url = fmt.Sprintf("http://%s/andon/led/%d/toggle", addr, ledID)
	}
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("[callAndon] Lỗi → %s: %v\n", url, err)
		return
	}
	resp.Body.Close()
	fmt.Printf("[callAndon] OK → %s\n", url)
}

func (s *HumidityAlertService) writeAlertLog(
	deviceName string, hum, tem float64,
	alertType string, andonTriggered bool,
) {
	db, err := database.TempHumidity_Connection()
	if err != nil {
		fmt.Println("[writeAlertLog] DB error:", err)
		return
	}
	dbInst, _ := db.DB()
	if dbInst != nil {
		defer dbInst.Close()
	}

	triggered := 0
	if andonTriggered {
		triggered = 1
	}

	res := db.Exec(`
		INSERT INTO humidity_alert_log
		    (device_name, hum, tem, alert_type, andon_triggered, triggered_at)
		VALUES (?, ?, ?, ?, ?, GETDATE())
	`, deviceName, hum, tem, alertType, triggered)

	if res.Error != nil {
		fmt.Printf("[writeAlertLog] Insert error %s: %v\n", deviceName, res.Error)
	} else {
		fmt.Printf("[writeAlertLog] OK → %s Hum=%.1f Type=%s Andon=%d\n",
			deviceName, hum, alertType, triggered)
	}
}

func (s *HumidityAlertService) markAndonOff(deviceName string, offBy string) {
	db, err := database.TempHumidity_Connection()
	if err != nil {
		return
	}
	dbInst, _ := db.DB()
	if dbInst != nil {
		defer dbInst.Close()
	}
	db.Exec(`
		UPDATE humidity_alert_log
		SET andon_off_by = ?,
		    andon_off_at = GETDATE()
		WHERE device_name = ?
		  AND andon_off_at IS NULL
		  AND andon_triggered = 1
	`, offBy, deviceName)
}