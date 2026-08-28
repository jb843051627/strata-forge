package model

import "time"

func OptionalTime(value time.Time) *time.Time {
	copy := value
	return &copy
}

func TimeOrZero(value *time.Time) time.Time {
	if value == nil {
		return time.Time{}
	}
	return *value
}
