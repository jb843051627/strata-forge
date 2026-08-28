package regression

import (
	"errors"
	"testing"

	"github.com/jb843051627/strata-forge/internal/engine"
	"github.com/jb843051627/strata-forge/internal/model"
)

func TestBug23_EmptyLayerAgeEstimateReturnsError(t *testing.T) {
	_, err := engine.NewAgeEstimator().Estimate(nil, []model.Measurement{{Value: 1, Uncertainty: 0.1}})
	if err == nil {
		t.Fatal("expected empty-layer error")
	}
}

func TestBug23_EmptyLayerEstimateKeepsInvalidInputIdentity(t *testing.T) {
	_, err := engine.NewAgeEstimator().Estimate(nil, []model.Measurement{{Value: 2, Uncertainty: 0.2}})
	if !errors.Is(err, model.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestBug23_SeveralMeasurementsWithoutLayerStillFailSafely(t *testing.T) {
	_, err := engine.NewAgeEstimator().Estimate(nil, []model.Measurement{{Value: 1}, {Value: 2}})
	if err == nil {
		t.Fatal("expected empty-layer error")
	}
}
