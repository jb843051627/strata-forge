package model

import "fmt"

type DepthInterval struct {
	Top    float64 `json:"top"`
	Bottom float64 `json:"bottom"`
}

func (d DepthInterval) Width() float64 {
	return d.Bottom - d.Top
}

func (d DepthInterval) Contains(value float64) bool {
	return value >= d.Top && value <= d.Bottom
}

func (d DepthInterval) Overlaps(other DepthInterval) bool {
	return d.Top < other.Bottom && other.Top < d.Bottom
}

func (d DepthInterval) Validate() error {
	if d.Top < 0 || d.Bottom <= d.Top {
		return fmt.Errorf("%w: invalid interval", ErrInvalidInput)
	}
	return nil
}

func MergeIntervals(values []DepthInterval) []DepthInterval {
	if len(values) < 2 {
		return append([]DepthInterval(nil), values...)
	}
	result := make([]DepthInterval, 0, len(values))
	for _, item := range values {
		if len(result) == 0 || !result[len(result)-1].Overlaps(item) {
			result = append(result, item)
			continue
		}
		if item.Bottom > result[len(result)-1].Bottom {
			result[len(result)-1].Bottom = item.Bottom
		}
	}
	return result
}
