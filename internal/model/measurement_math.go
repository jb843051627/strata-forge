package model

import "math"

func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}

func StandardDeviation(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	mean := Mean(values)
	variance := 0.0
	for _, value := range values {
		delta := value - mean
		variance += delta * delta
	}
	return math.Sqrt(variance / float64(len(values)-1))
}

func MeasurementValues(items []Measurement) []float64 {
	values := make([]float64, len(items))
	for i, item := range items {
		values[i] = item.Value
	}
	return values
}

func MeasurementUncertainties(items []Measurement) []float64 {
	values := make([]float64, len(items))
	for i, item := range items {
		values[i] = item.Uncertainty
	}
	return values
}
