package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) RecoverStaleRuns(ctx context.Context) (int64, error) {
	cutoff := sampleTime(s.clock)
	count, err := s.store.MarkStaleRuns(ctx, cutoff)
	return count, wrap("recover stale runs", err)
}

func (s *LabService) RetryRun(ctx context.Context, runID int64) (model.Run, error) {
	run, err := s.GetRun(ctx, runID)
	if err != nil {
		return model.Run{}, err
	}
	if run.State != model.RunFailed && run.State != model.RunCancelled {
		return model.Run{}, fmt.Errorf("%w: only failed runs can be retried", model.ErrConflict)
	}
	return s.QueueRun(ctx, run.SampleID)
}
