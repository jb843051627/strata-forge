package engine

import (
	"math"
	"sort"

	"github.com/jb843051627/strata-forge/internal/model"
)

type Outlier struct {
	Index  int
	Value  float64
	Reason string
}

func FindOutliers(values []float64, threshold float64) []Outlier {
	if len(values) < 3 {
		return nil
	}
	center := median(values)
	result := make([]Outlier, 0)
	for i, value := range values {
		if math.Abs(value-center) > threshold {
			result = append(result, Outlier{Index: i, Value: value, Reason: "distance from median"})
		}
	}
	return result
}

func median(values []float64) float64 {
	ordered := append([]float64(nil), values...)
	sort.Float64s(ordered)
	middle := len(ordered) / 2
	if len(ordered)%2 == 0 {
		return (ordered[middle-1] + ordered[middle]) / 2
	}
	return ordered[middle]
}

func MeasurementOutliers(items []model.Measurement, threshold float64) []Outlier {
	values := make([]float64, len(items))
	for i, item := range items {
		values[i] = item.Value
	}
	return FindOutliers(values, threshold)
}
