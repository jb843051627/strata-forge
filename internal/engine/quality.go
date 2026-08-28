package engine

import (
	"fmt"
	"math"

	"github.com/jb843051627/strata-forge/internal/model"
)

type QualityEngine struct {
	rules RuleSet
}

func NewQualityEngine() *QualityEngine {
	return &QualityEngine{rules: DefaultRules()}
}

func (e *QualityEngine) Evaluate(measurement model.Measurement) model.QualityResult {
	flags := make([]string, 0, 3)
	score := 100.0
	if measurement.InputValue < 0 {
		flags = append(flags, "negative_input")
		score -= 50
	}
	if measurement.Uncertainty < 0 {
		flags = append(flags, "negative_uncertainty")
		score -= 30
	}
	if measurement.Uncertainty > e.rules.MaxUncertainty {
		flags = append(flags, "uncertainty_high")
		score -= 35
	}
	if math.IsNaN(measurement.Value) || math.IsInf(measurement.Value, 0) {
		flags = append(flags, "non_finite_value")
		score -= 80
	}
	if score < 0 {
		score = 0
	}
	return model.QualityResult{
		Pass:        len(flags) == 0 && score >= e.rules.MinimumScore,
		Score:       score,
		Flags:       flags,
		Explanation: qualityExplanation(flags, score),
	}
}

func qualityExplanation(flags []string, score float64) string {
	if len(flags) == 0 {
		return fmt.Sprintf("quality score %.1f is within the laboratory limits", score)
	}
	return fmt.Sprintf("quality score %.1f raised flags: %v", score, flags)
}

func (e *QualityEngine) Acceptable(measurement model.Measurement) error {
	result := e.Evaluate(measurement)
	if !result.Pass {
		return fmt.Errorf("%w: %s", model.ErrQualityHold, result.Explanation)
	}
	return nil
}
