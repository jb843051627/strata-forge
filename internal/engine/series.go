package engine

import (
	"math"

	"github.com/jb843051627/strata-forge/internal/model"
)

func Smooth(values []float64, radius int) []float64 {
	if radius <= 0 || len(values) == 0 {
		return append([]float64(nil), values...)
	}
	result := make([]float64, len(values))
	for i := range values {
		start, end := i-radius, i+radius+1
		if start < 0 {
			start = 0
		}
		if end > len(values) {
			end = len(values)
		}
		total := 0.0
		for _, value := range values[start:end] {
			total += value
		}
		result[i] = total / float64(end-start)
	}
	return result
}

func Derivative(values []float64, step float64) []float64 {
	if len(values) < 2 || step <= 0 {
		return nil
	}
	result := make([]float64, len(values)-1)
	for i := range result {
		result[i] = (values[i+1] - values[i]) / step
	}
	return result
}

func SeriesVariance(items []model.Measurement) float64 {
	if len(items) < 2 {
		return 0
	}
	mean := 0.0
	for _, item := range items {
		mean += item.Value
	}
	mean /= float64(len(items))
	variance := 0.0
	for _, item := range items {
		delta := item.Value - mean
		variance += delta * delta
	}
	return variance / float64(len(items)-1)
}

func IsFiniteSeries(values []float64) bool {
	for _, value := range values {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return false
		}
	}
	return true
}
