package model

import "time"

type Sample struct {
	ID         int64     `json:"id"`
	Code       string    `json:"code"`
	Site       string    `json:"site"`
	DepthStart float64   `json:"depth_start"`
	DepthEnd   float64   `json:"depth_end"`
	Status     string    `json:"status"`
	ReceivedAt time.Time `json:"received_at"`
	Notes      string    `json:"notes"`
}

type Layer struct {
	ID          int64     `json:"id"`
	SampleID    int64     `json:"sample_id"`
	Sequence    int       `json:"sequence"`
	TopDepth    float64   `json:"top_depth"`
	BottomDepth float64   `json:"bottom_depth"`
	Material    string    `json:"material"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type Measurement struct {
	ID          int64      `json:"id"`
	LayerID     int64      `json:"layer_id"`
	Kind        string     `json:"kind"`
	InputValue  float64    `json:"input_value"`
	Unit        string     `json:"unit"`
	Value       float64    `json:"value"`
	Uncertainty float64    `json:"uncertainty"`
	Status      string     `json:"status"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	Operator    string     `json:"operator"`
}

type Run struct {
	ID          int64      `json:"id"`
	SampleID    int64      `json:"sample_id"`
	State       string     `json:"state"`
	RequestedAt time.Time  `json:"requested_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	CancelNote  string     `json:"cancel_note"`
}

type Review struct {
	ID            int64     `json:"id"`
	MeasurementID int64     `json:"measurement_id"`
	Decision      string    `json:"decision"`
	Comment       string    `json:"comment"`
	Reviewer      string    `json:"reviewer"`
	CreatedAt     time.Time `json:"created_at"`
}

type Report struct {
	ID         int64      `json:"id"`
	SampleID   int64      `json:"sample_id"`
	Version    int        `json:"version"`
	Status     string     `json:"status"`
	Summary    string     `json:"summary"`
	AgeMin     float64    `json:"age_min"`
	AgeMax     float64    `json:"age_max"`
	CreatedAt  time.Time  `json:"created_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty"`
}

type Alert struct {
	ID           int64     `json:"id"`
	SampleID     int64     `json:"sample_id"`
	LayerID      int64     `json:"layer_id"`
	Severity     string    `json:"severity"`
	Code         string    `json:"code"`
	Message      string    `json:"message"`
	Acknowledged bool      `json:"acknowledged"`
	CreatedAt    time.Time `json:"created_at"`
}

type Event struct {
	ID        int64     `json:"id"`
	SampleID  int64     `json:"sample_id"`
	Kind      string    `json:"kind"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type TimelineEntry struct {
	Kind      string    `json:"kind"`
	Reference int64     `json:"reference"`
	State     string    `json:"state"`
	At        time.Time `json:"at"`
}
