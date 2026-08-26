package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) AddLayer(ctx context.Context, sampleID int64, in model.LayerInput) (model.Layer, error) {
	if err := requirePositiveID(sampleID, "sample"); err != nil {
		return model.Layer{}, err
	}
	if err := model.ValidateLayerInput(in); err != nil {
		return model.Layer{}, err
	}
	sample, err := s.GetSample(context.Background(), sampleID)
	if err != nil {
		return model.Layer{}, err
	}
	if err := validateSampleLayer(sample, in); err != nil {
		return model.Layer{}, err
	}
	item := model.Layer{SampleID: sampleID, Sequence: in.Sequence, TopDepth: in.TopDepth, BottomDepth: in.BottomDepth, Material: in.Material, Status: model.MeasurementPending, CreatedAt: s.clock.Now()}
	created, err := s.store.CreateLayer(ctx, item)
	if err != nil {
		return model.Layer{}, wrap("add layer", err)
	}
	if sample.Status == model.SampleReceived {
		if err := s.workflow.MoveSample(sample.Status, model.SampleLayered); err != nil {
			return model.Layer{}, err
		}
		if err := s.store.UpdateSampleStatus(ctx, sampleID, model.SampleLayered); err != nil {
			return model.Layer{}, wrap("mark sample layered", err)
		}
	}
	_, err = s.store.AppendEvent(ctx, model.Event{SampleID: sampleID, Kind: "layer.added", Payload: fmt.Sprintf("layer=%d", created.ID), CreatedAt: s.clock.Now()})
	return created, wrap("record layer event", err)
}

func (s *LabService) GetLayer(ctx context.Context, id int64) (model.Layer, error) {
	if err := requirePositiveID(id, "layer"); err != nil {
		return model.Layer{}, err
	}
	item, err := s.store.GetLayer(ctx, id)
	return item, wrap("get layer", err)
}

func (s *LabService) ListLayers(ctx context.Context, sampleID int64) ([]model.Layer, error) {
	if err := requirePositiveID(sampleID, "sample"); err != nil {
		return nil, err
	}
	items, err := s.store.ListLayers(ctx, sampleID)
	return items, wrap("list layers", err)
}

func (s *LabService) LayerCoverage(ctx context.Context, sampleID int64) (float64, error) {
	top, bottom, err := s.store.LayerDepths(ctx, sampleID)
	if err != nil {
		return 0, wrap("read layer coverage", err)
	}
	sample, err := s.GetSample(ctx, sampleID)
	if err != nil {
		return 0, err
	}
	width := sample.DepthEnd - sample.DepthStart
	if width <= 0 {
		return 0, fmt.Errorf("%w: invalid sample width", model.ErrInvalidInput)
	}
	return (bottom - top) / width, nil
}
