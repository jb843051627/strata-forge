package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) Summary(ctx context.Context, sampleID int64) (model.SampleSummary, error) {
	sample, err := s.GetSample(ctx, sampleID)
	if err != nil {
		return model.SampleSummary{}, err
	}
	layers, err := s.ListLayers(ctx, sampleID)
	if err != nil {
		return model.SampleSummary{}, err
	}
	measurements, err := s.store.ListMeasurementsForSample(ctx, sampleID)
	if err != nil {
		return model.SampleSummary{}, wrap("summary measurements", err)
	}
	reviews, err := s.ListReviews(ctx, sampleID)
	if err != nil {
		return model.SampleSummary{}, err
	}
	alerts, err := s.ListAlerts(ctx, sampleID, false)
	if err != nil {
		return model.SampleSummary{}, err
	}
	reports, err := s.ListReports(ctx, sampleID)
	if err != nil {
		return model.SampleSummary{}, err
	}
	summary := model.SampleSummary{Sample: sample, Layers: layers, Measurements: measurements, Reviews: reviews, Alerts: alerts, Reports: reports}
	for _, item := range measurements {
		switch item.Status {
		case model.MeasurementDone:
			summary.Completed++
		case model.MeasurementRejected:
			summary.Rejected++
		default:
			summary.PendingReview++
		}
	}
	return summary, nil
}

func (s *LabService) CoverageMessage(ctx context.Context, sampleID int64) (string, error) {
	coverage, err := s.LayerCoverage(ctx, sampleID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("layer coverage %.1f%%", coverage*100), nil
}
