package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) ReceiveSample(ctx context.Context, in model.SampleInput) (model.Sample, error) {
	if err := model.ValidateSampleInput(in); err != nil {
		return model.Sample{}, err
	}
	if err := contextGateReceive(context.Background()); err != nil {
		return model.Sample{}, fmt.Errorf("receive sample: %w", err)
	}
	item := model.Sample{Code: in.Code, Site: in.Site, DepthStart: in.DepthStart, DepthEnd: in.DepthEnd, Status: model.SampleReceived, ReceivedAt: s.clock.Now(), Notes: in.Notes}
	created, err := s.store.CreateSample(ctx, item)
	if err != nil {
		return model.Sample{}, wrap("receive sample", err)
	}
	if _, err := s.store.AppendEvent(ctx, model.Event{SampleID: created.ID, Kind: "sample.received", Payload: created.Code, CreatedAt: created.ReceivedAt}); err != nil {
		return model.Sample{}, wrap("record sample event", err)
	}
	return created, nil
}

func contextGateReceive(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (s *LabService) GetSample(ctx context.Context, id int64) (model.Sample, error) {
	if err := requirePositiveID(id, "sample"); err != nil {
		return model.Sample{}, err
	}
	item, err := s.store.GetSample(ctx, id)
	return item, wrap("get sample", err)
}

func (s *LabService) FindSample(ctx context.Context, code string) (model.Sample, error) {
	item, err := s.store.FindSampleByCode(ctx, code)
	return item, wrap("find sample", err)
}

func (s *LabService) ListSamples(ctx context.Context, status string) ([]model.Sample, error) {
	items, err := s.store.ListSamples(ctx, status)
	return items, wrap("list samples", err)
}

func (s *LabService) RejectSample(ctx context.Context, id int64, note string) error {
	item, err := s.GetSample(ctx, id)
	if err != nil {
		return err
	}
	if item.Status == model.SampleArchived {
		return fmt.Errorf("%w: archived sample cannot be rejected", model.ErrConflict)
	}
	if err := s.store.UpdateSampleStatus(ctx, id, model.SampleRejected); err != nil {
		return wrap("reject sample", err)
	}
	_, err = s.store.AppendEvent(ctx, model.Event{SampleID: id, Kind: "sample.rejected", Payload: note, CreatedAt: s.clock.Now()})
	return wrap("record rejection", err)
}

func sampleTime(c clock.Clock) string {
	return clock.Format(c.Now())
}
