package service

import (
	"context"
	"time"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) WaitForRun(ctx context.Context, runID int64, timeout time.Duration) (model.Run, error) {
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	ticker := time.NewTicker(15 * time.Millisecond)
	defer ticker.Stop()
	for {
		run, err := s.GetRun(ctx, runID)
		if err != nil {
			return model.Run{}, err
		}
		if run.State == model.RunCompleted || run.State == model.RunCancelled || run.State == model.RunFailed {
			return run, nil
		}
		select {
		case <-ctx.Done():
			return model.Run{}, ctx.Err()
		case <-deadline.C:
			return model.Run{}, context.DeadlineExceeded
		case <-ticker.C:
		}
	}
}
