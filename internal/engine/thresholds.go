package engine

import "github.com/jb843051627/strata-forge/internal/model"

type Thresholds struct {
	NegativeTolerance float64
	UncertaintyWarn   float64
	UncertaintyFail   float64
	OutlierDistance   float64
}

func DefaultThresholds() Thresholds {
	return Thresholds{NegativeTolerance: 0, UncertaintyWarn: 0.25, UncertaintyFail: 0.5, OutlierDistance: 3}
}

func (t Thresholds) Flags(item model.Measurement) model.FlagSet {
	flags := model.NewFlagSet()
	if item.Value < t.NegativeTolerance {
		flags.Add(model.FlagNegativeValue)
	}
	if item.Uncertainty > t.UncertaintyWarn {
		flags.Add(model.FlagHighUncertainty)
	}
	return flags
}

func (t Thresholds) IsCritical(item model.Measurement) bool {
	return item.Uncertainty > t.UncertaintyFail || item.Value < t.NegativeTolerance
}
