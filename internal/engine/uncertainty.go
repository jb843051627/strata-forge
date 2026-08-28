package engine

import (
	"math"

	"github.com/jb843051627/strata-forge/internal/model"
)

func CombineUncertainty(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value * value
	}
	return math.Sqrt(total)
}

func MeasurementConfidence(item model.Measurement) float64 {
	return Clamp(1-item.Uncertainty, 0, 1)
}

func ConfidenceBand(estimate model.AgeEstimate) string {
	spread := estimate.Maximum - estimate.Minimum
	if spread < 1 {
		return "tight"
	}
	if spread < 3 {
		return "moderate"
	}
	return "wide"
}
