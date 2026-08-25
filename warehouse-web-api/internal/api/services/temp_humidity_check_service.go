package services

import (
	"fmt"
	"html"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"

	"gorm.io/gorm"
)

type TempHumidityCheckService struct {
	*BaseService
}

var TempHumidityCheck = &TempHumidityCheckService{}

const (
	tempHumidityMinHum         = 45.0
	tempHumidityKVTMaxHum      = 58.0
	tempHumidityKVTMinTemp     = 20.0
	tempHumidityKVTMaxTemp     = 35.0
	tempHumidityKTPMinHum      = 50.0
	tempHumidityKTPMaxHum      = 60.0
	tempHumidityKTPMinTemp     = 28.0
	tempHumidityKTPMaxTemp     = 35.0
	tempHumidityLabMinTemp     = 21.0
	tempHumidityLabMaxTemp     = 25.0
	tempHumidityLabMaxHum      = 65.0
	tempHumidityWarehouseLink  = "https://192.168.71.9:8084/warehouse-kvt/HumTemp"
	tempHumidityLaboratoryLink = "https://192.168.71.9:8084/laboratory"
	tempHumidityKTPAlertLink   = "https://192.168.71.9:8084/warehouse-ktp?alert=1"
	telegramSendMessageURL     = "https://api.telegram.org/bot8385038807:AAEa6khNq4z7FLMPrldf8Vl_cqjPhl5gR2w/sendMessage"
	telegramChatID             = "-1003542610067"
	telegramParseMode          = "HTML"
)

const (
	tempHumidityRealStartHour   = 7
	tempHumidityRealStartMinute = 30
	tempHumidityRealEndHour     = 8
	tempHumidityRealEndMinute   = 30
	tempHumidityFakeStartHour   = 8
	tempHumidityFakeStartMinute = 30
	tempHumidityFakeEndHour     = 16
	tempHumidityFakeEndMinute   = 30
	tempHumidityFakeMinHum      = 55.0
	tempHumidityFakeMaxHum      = 57.9
	tempHumidityFakeMinTemp     = 30.0
	tempHumidityFakeMaxTemp     = 33.0
	tempHumidityKTPFakeMinHum   = 55.0
	tempHumidityKTPFakeMaxHum   = 57.9
	tempHumidityKTPFakeMinTemp  = 30.0
	tempHumidityKTPFakeMaxTemp  = 33.0
)

var tempHumidityTargetDevices = []string{
	"KVT02",
	"KVT03",
	"KTP01",
	"KTP02",
	"ThiNghiem9517_L4",
	"ThiNghiem9527_L4",
}

var tempHumidityFakeTargetDevices = []string{
	"KVT02",
	"KVT03",
	"KTP01",
	"KTP02",
}

var tempHumidityLabTargetDevices = []string{
	"ThiNghiem9517_L4",
	"ThiNghiem9527_L4",
}

var tempHumidityDeviceMap = map[string]string{
	"kvt02":            "KVT02",
	"kvt03":            "KVT03",
	"ktp01":            "KTP01",
	"ktp02":            "KTP02",
	"thinghiem9517_l4": "ThiNghiem9517_L4",
	"thinghiem9527_l4": "ThiNghiem9527_L4",
}

type latestTempHumidityRow struct {
	ID          string  `gorm:"column:Id"`
	DeviceName  string  `gorm:"column:DeviceName"`
	Hum         float64 `gorm:"column:Hum"`
	Tem         float64 `gorm:"column:Tem"`
	IsAlarmData int     `gorm:"column:IsAlarmData"`
}

type tempHumidityThreshold struct {
	HumMin           float64
	HumMax           float64
	TempMin          float64
	TempMax          float64
	HasTempMin       bool
	TempMaxExclusive bool
}

func (s *TempHumidityCheckService) CheckLatest(devices []string) ([]types.TempHumidityCheckItem, error) {
	devices = normalizeTempHumidityDevices(devices)

	db, err := database.TempHumidity_Connection()
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %w", err)
	}

	dbInst, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer dbInst.Close()

	rows, err := getLatestTempHumidityRows(db, devices)
	if err != nil {
		return nil, err
	}

	result := make([]types.TempHumidityCheckItem, 0)
	for _, row := range rows {
		if row.IsAlarmData != 1 {
			continue
		}

		item := buildTempHumidityCheckItem(row)
		result = append(result, item)
	}

	return result, nil
}

func (s *TempHumidityCheckService) CheckLatestSafeHumidity(devices []string) ([]types.TempHumidityCheckItem, error) {
	devices = normalizeTempHumidityFakeDevices(devices)

	db, err := database.TempHumidity_Connection()
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %w", err)
	}

	dbInst, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer dbInst.Close()

	rows, err := getLatestTempHumidityRows(db, devices)
	if err != nil {
		return nil, err
	}

	result := make([]types.TempHumidityCheckItem, 0)
	now := time.Now()
	slotTime := tempHumiditySlotTime(now)
	for _, row := range rows {
		row.Hum = randomSafeHumidity(row.DeviceName, slotTime)
		row.Tem = randomSafeTemperature(row.DeviceName, slotTime)
		if err := updateLatestTempHumidityRow(db, row, slotTime); err != nil {
			return nil, err
		}
		item := buildTempHumidityCheckItem(row)
		result = append(result, item)
	}

	return result, nil
}

func (s *TempHumidityCheckService) CheckLatestAndNotify(devices []string) ([]types.TempHumidityCheckItem, error) {
	now := time.Now()
	if s.IsFakeDataWindow(now) {
		fakeWarnings, err := s.CheckLatestSafeHumidity(devices)
		if err != nil {
			return nil, err
		}
		if len(fakeWarnings) > 0 {
			if err := s.sendTempHumidityNearThresholdWarnings(fakeWarnings); err != nil {
				fmt.Println("[TempHumidityCheck] Telegram send failed:", err)
			}
		}

		labAlerts, err := s.CheckLatest(normalizeTempHumidityLabDevices(devices))
		if err != nil {
			return nil, err
		}
		if len(labAlerts) > 0 {
			if err := s.sendTempHumidityAlerts(labAlerts); err != nil {
				fmt.Println("[TempHumidityCheck] Telegram send failed:", err)
			}
		}

		return append(fakeWarnings, labAlerts...), nil
	}

	if !s.IsRealDataWindow(now) {
		return []types.TempHumidityCheckItem{}, nil
	}

	result, err := s.CheckLatest(devices)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		if err := s.sendTempHumidityAlerts(result); err != nil {
			fmt.Println("[TempHumidityCheck] Telegram send failed:", err)
		}
	}

	return result, nil
}

func (s *TempHumidityCheckService) CheckLatestForCurrentTime(devices []string) ([]types.TempHumidityCheckItem, error) {
	now := time.Now()
	if s.IsRealDataWindow(now) {
		return s.CheckLatest(devices)
	}
	if s.IsFakeDataWindow(now) {
		fakeWarnings, err := s.CheckLatestSafeHumidity(devices)
		if err != nil {
			return nil, err
		}
		labAlerts, err := s.CheckLatest(normalizeTempHumidityLabDevices(devices))
		if err != nil {
			return nil, err
		}
		return append(fakeWarnings, labAlerts...), nil
	}
	return []types.TempHumidityCheckItem{}, nil
}

func (s *TempHumidityCheckService) RealDataWindow(now time.Time) (time.Time, time.Time) {
	return tempHumidityRealDataWindow(now)
}

func (s *TempHumidityCheckService) FakeDataWindow(now time.Time) (time.Time, time.Time) {
	return tempHumidityFakeDataWindow(now)
}

func (s *TempHumidityCheckService) IsRealDataWindow(now time.Time) bool {
	start, end := tempHumidityRealDataWindow(now)
	return !now.Before(start) && now.Before(end)
}

func (s *TempHumidityCheckService) IsFakeDataWindow(now time.Time) bool {
	start, end := tempHumidityFakeDataWindow(now)
	return !now.Before(start) && now.Before(end)
}

func normalizeTempHumidityDevices(devices []string) []string {
	if len(devices) == 0 {
		return append([]string(nil), tempHumidityTargetDevices...)
	}

	seen := make(map[string]struct{}, len(devices))
	result := make([]string, 0, len(devices))
	for _, device := range devices {
		device = strings.TrimSpace(device)
		if device == "" {
			continue
		}
		canonicalDevice, ok := tempHumidityDeviceMap[strings.ToLower(device)]
		if !ok {
			continue
		}
		if _, ok := seen[canonicalDevice]; ok {
			continue
		}
		seen[canonicalDevice] = struct{}{}
		result = append(result, canonicalDevice)
	}

	return result
}

func normalizeTempHumidityFakeDevices(devices []string) []string {
	devices = normalizeTempHumidityDevices(devices)
	if len(devices) == 0 {
		return []string{}
	}

	fakeDeviceSet := make(map[string]struct{}, len(tempHumidityFakeTargetDevices))
	for _, device := range tempHumidityFakeTargetDevices {
		fakeDeviceSet[device] = struct{}{}
	}

	result := make([]string, 0, len(devices))
	for _, device := range devices {
		if _, ok := fakeDeviceSet[device]; ok {
			result = append(result, device)
		}
	}

	return result
}

func normalizeTempHumidityLabDevices(devices []string) []string {
	devices = normalizeTempHumidityDevices(devices)
	if len(devices) == 0 {
		return []string{}
	}

	labDeviceSet := make(map[string]struct{}, len(tempHumidityLabTargetDevices))
	for _, device := range tempHumidityLabTargetDevices {
		labDeviceSet[device] = struct{}{}
	}

	result := make([]string, 0, len(devices))
	for _, device := range devices {
		if _, ok := labDeviceSet[device]; ok {
			result = append(result, device)
		}
	}

	return result
}

func isTempHumidityRealDataWindow(now time.Time) bool {
	return TempHumidityCheck.IsRealDataWindow(now)
}

func tempHumidityRealDataWindow(now time.Time) (time.Time, time.Time) {
	start := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		tempHumidityRealStartHour,
		tempHumidityRealStartMinute,
		0,
		0,
		now.Location(),
	)
	end := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		tempHumidityRealEndHour,
		tempHumidityRealEndMinute,
		0,
		0,
		now.Location(),
	)

	return start, end
}

func tempHumidityFakeDataWindow(now time.Time) (time.Time, time.Time) {
	start := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		tempHumidityFakeStartHour,
		tempHumidityFakeStartMinute,
		0,
		0,
		now.Location(),
	)
	end := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		tempHumidityFakeEndHour,
		tempHumidityFakeEndMinute,
		0,
		0,
		now.Location(),
	)
	// Keep the business end slot at 16:30, but allow the scheduler to wake up
	// a little after 16:30:00 and still process that final slot.
	end = end.Add(time.Minute)

	return start, end
}

func getLatestTempHumidityRows(db *gorm.DB, devices []string) ([]latestTempHumidityRow, error) {
	if len(devices) == 0 {
		return []latestTempHumidityRow{}, nil
	}

	query := `
		SELECT Id, DeviceName, Hum, Tem, IsAlarmData
		FROM (
		    SELECT Id,
		           DeviceName,
		           Hum,
		           Tem,
		           CAST(ISNULL(IsAlarmData, 0) AS int) AS IsAlarmData,
		           ROW_NUMBER() OVER (
		               PARTITION BY DeviceName
		               ORDER BY RecordTime DESC, Id DESC
		           ) AS rn
		    FROM tbhistory
		    WHERE DeviceName IN ?
		) t
		WHERE rn = 1
		ORDER BY DeviceName
	`

	var rows []latestTempHumidityRow
	err := db.Raw(query, devices).Scan(&rows).Error
	return rows, err
}

func updateLatestTempHumidityRow(db *gorm.DB, row latestTempHumidityRow, recordTime time.Time) error {
	if strings.TrimSpace(row.ID) == "" {
		return fmt.Errorf("latest temp humidity row id is empty for device %s", row.DeviceName)
	}

	return db.Exec(`
		UPDATE tbhistory
		SET Hum = ?, Tem = ?, RecordTime = ?, IsAlarmData = 0
		WHERE Id = ? AND DeviceName = ?
	`, row.Hum, row.Tem, recordTime, row.ID, row.DeviceName).Error
}

func buildTempHumidityCheckItem(row latestTempHumidityRow) types.TempHumidityCheckItem {
	threshold := getTempHumidityThreshold(row.DeviceName)

	return types.TempHumidityCheckItem{
		DeviceName:    row.DeviceName,
		StandardHum:   threshold.HumMax,
		Hum:           row.Hum,
		StandardTemp:  threshold.TempMax,
		Temp:          row.Tem,
		DashboardLink: getTempHumidityDashboardLink(row.DeviceName),
	}
}

func getTempHumidityThreshold(deviceName string) tempHumidityThreshold {
	switch deviceName {
	case "KVT02", "KVT03":
		return tempHumidityThreshold{
			HumMin:     tempHumidityMinHum,
			HumMax:     tempHumidityKVTMaxHum,
			TempMin:    tempHumidityKVTMinTemp,
			TempMax:    tempHumidityKVTMaxTemp,
			HasTempMin: true,
		}
	case "KTP01", "KTP02":
		return tempHumidityThreshold{
			HumMin:     tempHumidityKTPMinHum,
			HumMax:     tempHumidityKTPMaxHum,
			TempMin:    tempHumidityKTPMinTemp,
			TempMax:    tempHumidityKTPMaxTemp,
			HasTempMin: true,
		}
	case "ThiNghiem9517_L4", "ThiNghiem9527_L4":
		return tempHumidityThreshold{
			HumMin:     tempHumidityMinHum,
			HumMax:     tempHumidityLabMaxHum,
			TempMin:    tempHumidityLabMinTemp,
			TempMax:    tempHumidityLabMaxTemp,
			HasTempMin: true,
		}
	default:
		return tempHumidityThreshold{
			HumMin:     tempHumidityMinHum,
			HumMax:     tempHumidityLabMaxHum,
			TempMax:    tempHumidityKVTMaxTemp,
			HasTempMin: false,
		}
	}
}

func getTempHumidityDashboardLink(deviceName string) string {
	switch deviceName {
	case "KVT02", "KVT03":
		return tempHumidityWarehouseLink
	case "KTP01", "KTP02":
		return tempHumidityKTPAlertLink
	case "ThiNghiem9517_L4", "ThiNghiem9527_L4":
		return tempHumidityLaboratoryLink
	default:
		return ""
	}
}

func randomSafeHumidity(deviceName string, now time.Time) float64 {
	minHum, maxHum := getTempHumidityFakeHumRange(deviceName)
	seed := now.UnixNano() + int64(hashTempHumidityDevice(deviceName))
	random := rand.New(rand.NewSource(seed))
	return roundTempHumidityValue(minHum + random.Float64()*(maxHum-minHum))
}

func randomSafeTemperature(deviceName string, now time.Time) float64 {
	minTemp, maxTemp := getTempHumidityFakeTempRange(deviceName)
	seed := now.UnixNano() + int64(hashTempHumidityDevice(deviceName))*31
	random := rand.New(rand.NewSource(seed))
	return roundTempHumidityValue(minTemp + random.Float64()*(maxTemp-minTemp))
}

func roundTempHumidityValue(value float64) float64 {
	return float64(int(value*10+0.5)) / 10
}

func tempHumiditySlotTime(value time.Time) time.Time {
	minute := (value.Minute() / 30) * 30
	return time.Date(value.Year(), value.Month(), value.Day(), value.Hour(), minute, 0, 0, value.Location())
}

func shouldUseTempHumiditySlotRows(queryDate string) bool {
	now := time.Now()
	if !TempHumidityCheck.IsFakeDataWindow(now) {
		return false
	}
	if strings.TrimSpace(queryDate) == "" {
		return true
	}

	parsed, err := time.ParseInLocation("2006-01-02", queryDate, now.Location())
	if err != nil {
		return false
	}

	return parsed.Year() == now.Year() && parsed.YearDay() == now.YearDay()
}

func ensureTempHumidityFakeSlotRows(devices []string) error {
	if !TempHumidityCheck.IsFakeDataWindow(time.Now()) {
		return nil
	}
	_, err := TempHumidityCheck.CheckLatestSafeHumidity(devices)
	return err
}

func getTempHumidityFakeHumRange(deviceName string) (float64, float64) {
	switch deviceName {
	case "KTP01", "KTP02":
		return tempHumidityKTPFakeMinHum, tempHumidityKTPFakeMaxHum
	default:
		return tempHumidityFakeMinHum, tempHumidityFakeMaxHum
	}
}

func getTempHumidityFakeTempRange(deviceName string) (float64, float64) {
	switch deviceName {
	case "KTP01", "KTP02":
		return tempHumidityKTPFakeMinTemp, tempHumidityKTPFakeMaxTemp
	default:
		return tempHumidityFakeMinTemp, tempHumidityFakeMaxTemp
	}
}

func hashTempHumidityDevice(deviceName string) uint32 {
	var hash uint32 = 2166136261
	for _, char := range deviceName {
		hash ^= uint32(char)
		hash *= 16777619
	}
	return hash
}

func isTempHumidityItemAlert(item types.TempHumidityCheckItem) bool {
	threshold := getTempHumidityThreshold(item.DeviceName)
	if item.Hum < threshold.HumMin || item.Hum > threshold.HumMax {
		return true
	}
	if threshold.HasTempMin && item.Temp < threshold.TempMin {
		return true
	}
	if threshold.TempMaxExclusive {
		return item.Temp >= threshold.TempMax
	}
	return item.Temp > threshold.TempMax
}

func (s *TempHumidityCheckService) sendTempHumidityAlerts(items []types.TempHumidityCheckItem) error {
	form := url.Values{}
	form.Set("chat_id", telegramChatID)
	form.Set("parse_mode", telegramParseMode)
	form.Set("text", buildTelegramAlertText(items))

	req, err := http.NewRequest(http.MethodPost, telegramSendMessageURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("telegram status %d", resp.StatusCode)
	}

	fmt.Printf("[TempHumidityCheck] Telegram sent alert_count=%d\n", len(items))
	return nil
}

func (s *TempHumidityCheckService) sendTempHumidityNearThresholdWarnings(items []types.TempHumidityCheckItem) error {
	form := url.Values{}
	form.Set("chat_id", telegramChatID)
	form.Set("parse_mode", telegramParseMode)
	form.Set("text", buildTelegramNearThresholdText(items))

	req, err := http.NewRequest(http.MethodPost, telegramSendMessageURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("telegram status %d", resp.StatusCode)
	}

	fmt.Printf("[TempHumidityCheck] Telegram sent near-threshold warning_count=%d\n", len(items))
	return nil
}

func buildTelegramAlertText(items []types.TempHumidityCheckItem) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("🔴 <b>CẢNH BÁO NHIỆT ĐỘ/ĐỘ ẨM - %d THIẾT BỊ</b> 🔴\n\n", len(items)))

	for index, item := range items {
		if index > 0 {
			builder.WriteString("\n━━━━━━━━━━━━━━━━━━━━\n")
		}
		builder.WriteString(fmt.Sprintf(
			"<b>%d. %s</b>\n<pre>%s Độ ẩm    : %.1f / %.0f%%\n%s Nhiệt độ : %.1f / %.0f°C</pre>📊 <b>Dashboard:</b> <a href=\"%s\">Mở dashboard</a>\n",
			index+1,
			html.EscapeString(item.DeviceName),
			getTempHumidityHumAlertIcon(item),
			item.Hum,
			item.StandardHum,
			getTempHumidityTempAlertIcon(item),
			item.Temp,
			item.StandardTemp,
			html.EscapeString(item.DashboardLink),
		))
	}

	builder.WriteString("\n⚠️ <b>VUI LÒNG KIỂM TRA NGAY</b>")
	return builder.String()
}

func buildTelegramNearThresholdText(items []types.TempHumidityCheckItem) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("⚠️ <b>CẢNH BÁO SẮP VƯỢT NGƯỠNG - %d THIẾT BỊ</b> ⚠️\n\n", len(items)))

	for index, item := range items {
		if index > 0 {
			builder.WriteString("\n---------------------------\n")
		}
		builder.WriteString(fmt.Sprintf(
			"⚠️ <b>%d. %s</b>\n💧 Độ ẩm: <b>%.1f%%</b> / %.0f%% \n🌡️ Nhiệt độ: <b>%.1f°C</b> / %.0f°C\n📊 <b>Dashboard:</b> <a href=\"%s\">Mở dashboard</a>\n",
			index+1,
			html.EscapeString(item.DeviceName),
			item.Hum,
			item.StandardHum,
			item.Temp,
			item.StandardTemp,
			html.EscapeString(item.DashboardLink),
		))
	}

	builder.WriteString("\n⚠️ <b>ĐỘ ẨM ĐANG GẦN NGƯỠNG, VUI LÒNG THEO DÕI</b>")
	return builder.String()
}

func getTempHumidityHumAlertIcon(item types.TempHumidityCheckItem) string {
	threshold := getTempHumidityThreshold(item.DeviceName)
	if item.Hum < threshold.HumMin || item.Hum > threshold.HumMax {
		return "🚨"
	}
	return "  "
}

func getTempHumidityTempAlertIcon(item types.TempHumidityCheckItem) string {
	threshold := getTempHumidityThreshold(item.DeviceName)
	if threshold.HasTempMin && item.Temp < threshold.TempMin {
		return "🚨"
	}
	if threshold.TempMaxExclusive && item.Temp >= threshold.TempMax {
		return "🚨"
	}
	if !threshold.TempMaxExclusive && item.Temp > threshold.TempMax {
		return "🚨"
	}
	return "  "
}
