package engine

import (
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

type Decision struct {
	Status string
	Reason string
}

func DecideQuality(result model.QualityResult) Decision {
	if result.Pass {
		return Decision{Status: model.ReviewPass, Reason: result.Explanation}
	}
	if result.Score >= 50 {
		return Decision{Status: model.ReviewHold, Reason: result.Explanation}
	}
	return Decision{Status: model.ReviewFail, Reason: result.Explanation}
}

func RequireDecision(result model.QualityResult, expected string) error {
	decision := DecideQuality(result)
	if decision.Status != expected {
		return fmt.Errorf("%w: expected %s, got %s", model.ErrQualityHold, expected, decision.Status)
	}
	return nil
}
