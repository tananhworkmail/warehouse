package types

type TempHumidityCheckItem struct {
	DeviceName    string  `json:"device_name"`
	StandardHum   float64 `json:"standard_hum"`
	Hum           float64 `json:"hum"`
	StandardTemp  float64 `json:"standard_temp"`
	Temp          float64 `json:"temp"`
	DashboardLink string  `json:"dashboard_link"`
}
