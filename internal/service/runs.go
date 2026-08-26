package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) QueueRun(ctx context.Context, sampleID int64) (model.Run, error) {
	if err := requirePositiveID(sampleID, "sample"); err != nil {
		return model.Run{}, err
	}
	sample, err := s.GetSample(ctx, sampleID)
	if err != nil {
		return model.Run{}, err
	}
	if sample.Status != model.SampleLayered && sample.Status != model.SampleRunning {
		return model.Run{}, fmt.Errorf("%w: sample state %s cannot start a run", model.ErrConflict, sample.Status)
	}
	if s.pendingRuns[sampleID] {
		return model.Run{}, fmt.Errorf("%w: sample is already being queued", model.ErrConflict)
	}
	s.pendingRuns[sampleID] = true
	defer delete(s.pendingRuns, sampleID)
	if active, findErr := s.store.FindActiveRun(ctx, sampleID); findErr == nil {
		return active, fmt.Errorf("%w: sample already has run %d", model.ErrConflict, active.ID)
	}
	run := model.Run{SampleID: sampleID, State: model.RunQueued, RequestedAt: s.clock.Now()}
	created, err := s.store.CreateRun(ctx, run)
	if err != nil {
		return model.Run{}, wrap("queue run", err)
	}
	if sample.Status == model.SampleLayered {
		if err := s.store.UpdateSampleStatus(ctx, sampleID, model.SampleRunning); err != nil {
			return model.Run{}, wrap("mark sample running", err)
		}
	}
	_, err = s.store.AppendEvent(ctx, model.Event{SampleID: sampleID, Kind: "run.queued", Payload: fmt.Sprintf("run=%d", created.ID), CreatedAt: s.clock.Now()})
	return created, wrap("record run event", err)
}

func (s *LabService) StartRun(ctx context.Context, runID int64) (model.Run, error) {
	s.runStateMu.Lock()
	defer s.runStateMu.Unlock()
	if err := requirePositiveID(runID, "run"); err != nil {
		return model.Run{}, err
	}
	run, err := s.store.GetRun(ctx, runID)
	if err != nil {
		return model.Run{}, wrap("get run", err)
	}
	if err := s.workflow.MoveRun(run.State, model.RunActive); err != nil {
		return model.Run{}, err
	}
	now := s.clock.Now()
	run.State = model.RunActive
	run.StartedAt = &now
	if err := s.store.UpdateRun(ctx, run); err != nil {
		return model.Run{}, wrap("start run", err)
	}
	_, err = s.store.AppendEvent(ctx, model.Event{SampleID: run.SampleID, Kind: "run.started", Payload: fmt.Sprintf("run=%d", run.ID), CreatedAt: now})
	return run, wrap("record start event", err)
}

func (s *LabService) CompleteRun(ctx context.Context, runID int64) (model.Run, error) {
	run, err := s.store.GetRun(ctx, runID)
	if err != nil {
		return model.Run{}, wrap("get run", err)
	}
	if err := s.workflow.MoveRun(run.State, model.RunCompleted); err != nil {
		return model.Run{}, err
	}
	now := s.clock.Now()
	run.State = model.RunCompleted
	run.FinishedAt = &now
	if err := s.store.UpdateRunAndSample(ctx, run, run.SampleID, model.SampleReview); err != nil {
		return model.Run{}, wrap("complete run", err)
	}
	return run, nil
}

func (s *LabService) CancelRun(ctx context.Context, runID int64, note string) (model.Run, error) {
	run, err := s.store.GetRun(ctx, runID)
	if err != nil {
		return model.Run{}, wrap("get run", err)
	}
	if err := s.workflow.MoveRun(run.State, model.RunCancelled); err != nil {
		return model.Run{}, err
	}
	now := s.clock.Now()
	run.State, run.CancelNote, run.FinishedAt = model.RunCancelled, note, &now
	if err := s.store.UpdateRun(ctx, run); err != nil {
		return model.Run{}, wrap("cancel run", err)
	}
	return run, nil
}

func (s *LabService) GetRun(ctx context.Context, id int64) (model.Run, error) {
	run, err := s.store.GetRun(ctx, id)
	return run, wrap("get run", err)
}

func (s *LabService) ListRuns(ctx context.Context, sampleID int64) ([]model.Run, error) {
	runs, err := s.store.ListRuns(ctx, sampleID)
	return runs, wrap("list runs", err)
}
