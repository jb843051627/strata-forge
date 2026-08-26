package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) ArchiveSample(ctx context.Context, sampleID int64) (model.Sample, error) {
	sample, err := s.GetSample(ctx, sampleID)
	if err != nil {
		return model.Sample{}, err
	}
	if sample.Status != model.SampleReview {
		return model.Sample{}, fmt.Errorf("%w: sample must be in review", model.ErrConflict)
	}
	if err := s.store.UpdateSampleStatus(ctx, sampleID, model.SampleArchived); err != nil {
		return model.Sample{}, wrap("archive sample", err)
	}
	sample.Status = model.SampleArchived
	return sample, nil
}

func (s *LabService) ReopenSample(ctx context.Context, sampleID int64) (model.Sample, error) {
	sample, err := s.GetSample(ctx, sampleID)
	if err != nil {
		return model.Sample{}, err
	}
	if sample.Status != model.SampleReview && sample.Status != model.SampleRejected {
		return model.Sample{}, fmt.Errorf("%w: sample cannot be reopened from %s", model.ErrConflict, sample.Status)
	}
	if err := s.store.UpdateSampleStatus(ctx, sampleID, model.SampleRunning); err != nil {
		return model.Sample{}, wrap("reopen sample", err)
	}
	sample.Status = model.SampleRunning
	return sample, nil
}
