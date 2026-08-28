package clock

import "time"

type Clock interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

func Format(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

func Parse(value string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, value)
}
