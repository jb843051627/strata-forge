package engine

import "math"

func WeightedAge(values, weights []float64) float64 {
	if len(values) == 0 || len(values) != len(weights) {
		return 0
	}
	total, weight := 0.0, 0.0
	for i, value := range values {
		if math.IsNaN(value) || math.IsInf(value, 0) || weights[i] <= 0 {
			continue
		}
		total += value * weights[i]
		weight += weights[i]
	}
	if weight == 0 {
		return 0
	}
	return total / weight
}

func Clamp(value, low, high float64) float64 {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}

func RelativeSpread(low, high float64) float64 {
	center := (low + high) / 2
	if center == 0 {
		return 0
	}
	return (high - low) / math.Abs(center)
}
