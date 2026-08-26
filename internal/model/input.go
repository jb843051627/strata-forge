package model

type SampleInput struct {
	Code       string  `json:"code"`
	Site       string  `json:"site"`
	DepthStart float64 `json:"depth_start"`
	DepthEnd   float64 `json:"depth_end"`
	Notes      string  `json:"notes"`
}

type LayerInput struct {
	Sequence    int     `json:"sequence"`
	TopDepth    float64 `json:"top_depth"`
	BottomDepth float64 `json:"bottom_depth"`
	Material    string  `json:"material"`
}

type MeasurementInput struct {
	LayerID    int64   `json:"layer_id"`
	Kind       string  `json:"kind"`
	InputValue float64 `json:"input_value"`
	Unit       string  `json:"unit"`
	Operator   string  `json:"operator"`
}

type ReviewInput struct {
	Decision string `json:"decision"`
	Comment  string `json:"comment"`
	Reviewer string `json:"reviewer"`
}

type RunInput struct {
	SampleID int64 `json:"sample_id"`
}

type AlertInput struct {
	SampleID int64  `json:"sample_id"`
	LayerID  int64  `json:"layer_id"`
	Severity string `json:"severity"`
	Code     string `json:"code"`
	Message  string `json:"message"`
}
