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
	Percentage float64 `json:"percentage"`
	Cumulative float64 `json:"cumulative"`
}

type GoreTexDashboardPeriod struct {
	Period               string                       `json:"period"`
	Value                string                       `json:"value"`
	FromDate             string                       `json:"fromDate"`
	ToDate               string                       `json:"toDate"`
	Styles               []string                     `json:"styles"`
	SelectedStyle        string                       `json:"selectedStyle"`
	SuterTrend           []GoreTexDashboardTrendPoint `json:"suterTrend"`
	VisualizationResults []GoreTexDashboardParetoItem `json:"visualizationResults"`
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
