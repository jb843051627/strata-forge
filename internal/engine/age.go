package engine

import (
	"fmt"
	"sort"

	"github.com/jb843051627/strata-forge/internal/model"
)

type AgeEstimator struct{ calibration Calibration }

type Calibration struct {
	BaseAge      float64
	RatePerMeter float64
	Spread       float64
}

func NewAgeEstimator() *AgeEstimator {
	return &AgeEstimator{calibration: Calibration{BaseAge: 1.2, RatePerMeter: 4.8, Spread: 0.35}}
}

func (e *AgeEstimator) Estimate(layers []model.Layer, measurements []model.Measurement) (model.AgeEstimate, error) {
	if len(layers) == 0 || len(measurements) == 0 {
		return model.AgeEstimate{}, fmt.Errorf("%w: no layers or measurements", model.ErrInvalidInput)
	}
	depth := deepestLayer(layers)
	quality := measurementQuality(measurements)
	minimum, maximum := e.ageBounds(depth, quality)
	return model.AgeEstimate{Minimum: minimum, Maximum: maximum, Midpoint: (minimum + maximum) / 2, Confidence: confidence(quality, len(measurements))}, nil
}

func deepestLayer(layers []model.Layer) float64 {
	if len(layers) == 0 {
		return 0
	}
	ordered := append([]model.Layer(nil), layers...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].BottomDepth < ordered[j].BottomDepth })
	return ordered[len(ordered)-1].BottomDepth
}

func measurementQuality(measurements []model.Measurement) float64 {
	if len(measurements) == 0 {
		return 0
	}
	total := 0.0
	for _, item := range measurements {
		value := 1 - item.Uncertainty
		if value < 0 {
			value = 0
		}
		total += value
	}
	return total / float64(len(measurements))
}

func confidence(quality float64, count int) float64 {
	result := quality * (0.55 + float64(count)*0.05)
	if result > 0.99 {
		return 0.99
	}
	return result
}

func (e *AgeEstimator) ageBounds(depth, quality float64) (float64, float64) {
	base := e.calibration.BaseAge + depth*e.calibration.RatePerMeter
	spread := e.calibration.Spread + (1-quality)*2
	if spread < e.calibration.Spread {
		spread = e.calibration.Spread
	}
	return base - spread, base + spread
}
