package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/engine"
	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) QualityDecision(ctx context.Context, measurementID int64) (engine.Decision, error) {
	measurement, err := s.GetMeasurement(ctx, measurementID)
	if err != nil {
		return engine.Decision{}, err
	}
	result := s.quality.Evaluate(measurement)
	decision := engine.DecideQuality(result)
	if decision.Status == model.ReviewFail {
		return decision, fmt.Errorf("%w: %s", model.ErrQualityHold, decision.Reason)
	}
	return decision, nil
}

func (s *LabService) CanArchive(ctx context.Context, sampleID int64) (bool, error) {
	sample, err := s.GetSample(ctx, sampleID)
	if err != nil {
		return false, err
	}
	return sample.Status == model.SampleReview, nil
}
