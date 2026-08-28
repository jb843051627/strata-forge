package model

import "time"

type MeasurementWindow struct {
	Earliest time.Time `json:"earliest"`
	Latest   time.Time `json:"latest"`
}

func (w MeasurementWindow) Includes(value time.Time) bool {
	return !value.Before(w.Earliest) && !value.After(w.Latest)
}

func (w MeasurementWindow) Empty() bool {
	return w.Earliest.IsZero() || w.Latest.IsZero() || w.Latest.Before(w.Earliest)
}

func (w MeasurementWindow) Duration() time.Duration {
	if w.Empty() {
		return 0
	}
	return w.Latest.Sub(w.Earliest)
}
