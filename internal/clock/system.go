package clock

import "time"

type System struct{}

func (System) Now() time.Time {
	return time.Now().UTC()
}

func (System) Since(start time.Time) time.Duration {
	return time.Since(start)
}
