package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) ReviewMeasurement(ctx context.Context, id int64, in model.ReviewInput) (model.Review, error) {
	if err := model.ValidateReviewInput(in); err != nil {
		return model.Review{}, err
	}
	if err := validateReviewComment(in); err != nil {
		return model.Review{}, err
	}
	measurement, err := s.GetMeasurement(ctx, id)
	if err != nil {
		return model.Review{}, err
	}
	if measurement.Status != model.MeasurementDone {
		return model.Review{}, fmt.Errorf("%w: measurement is not complete", model.ErrConflict)
	}
	layer, err := s.GetLayer(ctx, measurement.LayerID)
	if err != nil {
		return model.Review{}, err
	}
	now := s.clock.Now()
	review := model.Review{MeasurementID: id, Decision: in.Decision, Comment: in.Comment, Reviewer: in.Reviewer, CreatedAt: now}
	status := model.SampleReview
	if err := s.store.SaveMeasurementReview(ctx, measurement, review, layer.SampleID, status); err != nil {
		return model.Review{}, wrap("review measurement", err)
	}
	if in.Decision == model.ReviewFail {
		_, _ = s.store.CreateAlert(ctx, model.Alert{SampleID: layer.SampleID, LayerID: layer.ID, Severity: "critical", Code: "review-failed", Message: in.Comment, CreatedAt: now})
	}
	return review, nil
}

func (s *LabService) LatestReview(ctx context.Context, measurementID int64) (model.Review, error) {
	item, err := s.store.LatestReview(ctx, measurementID)
	return item, wrap("latest review", err)
}

func (s *LabService) ListReviews(ctx context.Context, sampleID int64) ([]model.Review, error) {
	items, err := s.store.ListReviewsForSample(ctx, sampleID)
	return items, wrap("list reviews", err)
}
