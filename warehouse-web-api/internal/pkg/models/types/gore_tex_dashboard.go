package types

type GoreTexDashboardRateItem struct {
	Label string  `json:"label"`
	Pass  int     `json:"pass"`
	Total int     `json:"total"`
	Rate  float64 `json:"rate"`
}

type GoreTexDashboardTrendPoint struct {
	Date   string  `json:"date"`
	Label  string  `json:"label"`
	Pass   int     `json:"pass"`
	Total  int     `json:"total"`
	Rate   float64 `json:"rate"`
	Target float64 `json:"target"`
}

type GoreTexDashboardComparison struct {
	Year  int                        `json:"year"`
	Week  int                        `json:"week"`
	Label string                     `json:"label"`
	Items []GoreTexDashboardRateItem `json:"items"`
}

type GoreTexDashboardParetoItem struct {
	Label      string  `json:"label"`
	Count      float64 `json:"count"`
	Cumulative float64 `json:"cumulative"`
}

type GoreTexWeeklyDashboard struct {
	Year                 int                          `json:"year"`
	Week                 int                          `json:"week"`
	FromDate             string                       `json:"fromDate"`
	ToDate               string                       `json:"toDate"`
	SuterByItems         []GoreTexDashboardRateItem   `json:"suterByItems"`
	SuterTrend           []GoreTexDashboardTrendPoint `json:"suterTrend"`
	SuterComparison      []GoreTexDashboardComparison `json:"suterComparison"`
	RRdyTrend            []GoreTexDashboardTrendPoint `json:"rRdyTrend"`
	VisualizationResults []GoreTexDashboardParetoItem `json:"visualizationResults"`
}
