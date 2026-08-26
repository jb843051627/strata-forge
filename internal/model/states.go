package model

const (
	SampleReceived = "received"
	SampleLayered  = "layered"
	SampleRunning  = "running"
	SampleReview   = "review"
	SampleArchived = "archived"
	SampleRejected = "rejected"

	RunQueued    = "queued"
	RunActive    = "active"
	RunCompleted = "completed"
	RunCancelled = "cancelled"
	RunFailed    = "failed"

	MeasurementPending  = "pending"
	MeasurementRunning  = "running"
	MeasurementDone     = "done"
	MeasurementRejected = "rejected"

	ReviewPending = "pending"
	ReviewPass    = "pass"
	ReviewHold    = "hold"
	ReviewFail    = "fail"

	ReportDraft    = "draft"
	ReportFinal    = "final"
	ReportArchived = "archived"
)

func CanSampleTransition(from, to string) bool {
	allowed := map[string]map[string]bool{
		SampleReceived: {SampleLayered: true, SampleRejected: true},
		SampleLayered:  {SampleRunning: true, SampleRejected: true},
		SampleRunning:  {SampleReview: true, SampleRejected: true},
		SampleReview:   {SampleArchived: true, SampleRunning: true, SampleRejected: true},
		SampleArchived: {},
		SampleRejected: {},
	}
	return allowed[from][to]
}

func CanRunTransition(from, to string) bool {
	allowed := map[string]map[string]bool{
		RunQueued:    {RunActive: true, RunCancelled: true, RunFailed: true},
		RunActive:    {RunCompleted: true, RunCancelled: true, RunFailed: true},
		RunCompleted: {},
		RunCancelled: {},
		RunFailed:    {},
	}
	return allowed[from][to]
}

func CanMeasurementTransition(from, to string) bool {
	allowed := map[string]map[string]bool{
		MeasurementPending:  {MeasurementRunning: true, MeasurementRejected: true},
		MeasurementRunning:  {MeasurementDone: true, MeasurementRejected: true},
		MeasurementDone:     {},
		MeasurementRejected: {},
	}
	return allowed[from][to]
}
