package engine

import "time"

type Window struct {
	Start time.Time
	End   time.Time
}

func (w Window) Contains(value time.Time) bool {
	return !value.Before(w.Start) && value.Before(w.End)
}

func (w Window) Duration() time.Duration {
	return w.End.Sub(w.Start)
}

func (w Window) Valid() bool {
	return w.End.After(w.Start)
}

func NewWindow(start time.Time, duration time.Duration) Window {
	return Window{Start: start, End: start.Add(duration)}
}
