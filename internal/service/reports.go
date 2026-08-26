package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/engine"
	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) BuildReport(ctx context.Context, sampleID int64) (model.Report, error) {
	sample, err := s.GetSample(context.Background(), sampleID)
	if err != nil {
		return model.Report{}, err
	}
	layers, err := s.ListLayers(context.Background(), sampleID)
	if err != nil {
		return model.Report{}, err
	}
	measurements, err := s.store.ListMeasurementsForSample(context.Background(), sampleID)
	if err != nil {
		return model.Report{}, wrap("read report measurements", err)
	}
	if err := s.ensureReviewable(context.Background(), sampleID, measurements); err != nil {
		return model.Report{}, err
	}
	estimate, err := s.age.Estimate(engine.SnapshotLayers(layers), engine.SnapshotMeasurements(measurements))
	if err != nil {
		return model.Report{}, wrap("estimate sample age", err)
	}
	version, err := s.store.NextReportVersion(context.Background(), sampleID)
	if err != nil {
		return model.Report{}, wrap("allocate report version", err)
	}
	summary := fmt.Sprintf("%s: %d layers, %d measurements, estimated %.2f-%.2f ka", sample.Code, len(layers), len(measurements), estimate.Minimum, estimate.Maximum)
	report := model.Report{SampleID: sampleID, Version: version, Status: model.ReportFinal, Summary: summary, AgeMin: estimate.Minimum, AgeMax: estimate.Maximum, CreatedAt: s.clock.Now()}
	eventPayload, _ := json.Marshal(map[string]any{"report_version": version, "age_min": estimate.Minimum, "age_max": estimate.Maximum})
	saved, _, err := s.store.SaveReportAndEvent(context.Background(), report, model.Event{SampleID: sampleID, Kind: "report.created", Payload: string(eventPayload), CreatedAt: s.clock.Now()})
	if err != nil {
		return model.Report{}, wrap("save report", err)
	}
	if err := s.store.UpdateSampleStatus(context.Background(), sampleID, model.SampleReview); err != nil {
		return model.Report{}, wrap("mark sample review", err)
	}
	return saved, nil
}

func (s *LabService) ensureReviewable(ctx context.Context, sampleID int64, measurements []model.Measurement) error {
	if len(measurements) == 0 {
		return fmt.Errorf("%w: sample has no measurements", model.ErrQualityHold)
	}
	for _, measurement := range measurements {
		if measurement.Status != model.MeasurementDone {
			return fmt.Errorf("%w: measurement %d is not complete", model.ErrQualityHold, measurement.ID)
		}
		review, err := s.store.LatestReview(ctx, measurement.ID)
		if err != nil {
			return fmt.Errorf("%w: measurement %d has no review", model.ErrQualityHold, measurement.ID)
		}
		if review.Decision != model.ReviewPass {
			return fmt.Errorf("%w: measurement %d is not accepted", model.ErrQualityHold, measurement.ID)
		}
	}
	return nil
}

func (s *LabService) GetReport(ctx context.Context, id int64) (model.Report, error) {
	report, err := s.store.GetReport(ctx, id)
	return report, wrap("get report", err)
}

func (s *LabService) ListReports(ctx context.Context, sampleID int64) ([]model.Report, error) {
	reports, err := s.store.ListReports(ctx, sampleID)
	return reports, wrap("list reports", err)
}

func (s *LabService) ArchiveReport(ctx context.Context, id int64) error {
	if _, err := s.GetReport(ctx, id); err != nil {
		return err
	}
	if err := s.store.ArchiveReport(ctx, id, sampleTime(s.clock)); err != nil {
		return wrap("archive report", err)
	}
	return nil
}
