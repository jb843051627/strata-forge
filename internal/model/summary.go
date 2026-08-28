package model

type SampleSummary struct {
	Sample        Sample        `json:"sample"`
	Layers        []Layer       `json:"layers"`
	Measurements  []Measurement `json:"measurements"`
	Reviews       []Review      `json:"reviews"`
	Alerts        []Alert       `json:"alerts"`
	Reports       []Report      `json:"reports"`
	Completed     int           `json:"completed"`
	Rejected      int           `json:"rejected"`
	PendingReview int           `json:"pending_review"`
}

type QualityResult struct {
	Pass        bool     `json:"pass"`
	Score       float64  `json:"score"`
	Flags       []string `json:"flags"`
	Explanation string   `json:"explanation"`
}

type AgeEstimate struct {
	Minimum    float64 `json:"minimum"`
	Maximum    float64 `json:"maximum"`
	Midpoint   float64 `json:"midpoint"`
	Confidence float64 `json:"confidence"`
}
