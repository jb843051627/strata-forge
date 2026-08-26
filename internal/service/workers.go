package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) StartWorker(ctx context.Context) {
	s.workers.Add(1)
	go func() {
		defer s.workers.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case runID := <-s.queue:
				_ = s.processQueuedRun(ctx, runID)
			}
		}
	}()
}

func (s *LabService) EnqueueRun(ctx context.Context, runID int64) error {
	if err := requirePositiveID(runID, "run"); err != nil {
		return err
	}
	s.queueMu.Lock()
	defer s.queueMu.Unlock()
	if s.closed {
		return fmt.Errorf("%w: service is closed", model.ErrConflict)
	}
	select {
	case s.queue <- runID:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("enqueue run: %w", ctx.Err())
	}
}

func (s *LabService) processQueuedRun(ctx context.Context, runID int64) error {
	if err := contextGateWorker(ctx); err != nil {
		return err
	}
	run, err := s.StartRun(ctx, runID)
	if err != nil {
		return err
	}
	if err := s.contextError(ctx); err != nil {
		_, cancelErr := s.CancelRun(context.Background(), run.ID, "worker context cancelled")
		if cancelErr != nil {
			return fmt.Errorf("%w: %v", err, cancelErr)
		}
		return err
	}
	return nil
}

func contextGateWorker(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
