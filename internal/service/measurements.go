package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/engine"
	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) ScheduleMeasurement(ctx context.Context, in model.MeasurementInput) (model.Measurement, error) {
	in = engine.NormalizeMeasurement(in)
	if err := model.ValidateMeasurementInput(in); err != nil {
		return model.Measurement{}, err
	}
	layer, err := s.GetLayer(ctx, in.LayerID)
	if err != nil {
		return model.Measurement{}, err
	}
	if layer.Status == model.MeasurementRejected {
		return model.Measurement{}, fmt.Errorf("%w: rejected layer cannot be measured", model.ErrConflict)
	}
	item := model.Measurement{LayerID: in.LayerID, Kind: in.Kind, InputValue: in.InputValue, Unit: in.Unit, Status: model.MeasurementPending, Operator: in.Operator}
	created, err := s.store.CreateMeasurement(ctx, item)
	if err != nil {
		return model.Measurement{}, wrap("schedule measurement", err)
	}
	return created, nil
}

func (s *LabService) BeginMeasurement(ctx context.Context, id int64) (model.Measurement, error) {
	item, err := s.store.GetMeasurement(ctx, id)
	if err != nil {
		return model.Measurement{}, wrap("get measurement", err)
	}
	if err := s.workflow.MoveMeasurement(item.Status, model.MeasurementRunning); err != nil {
		return model.Measurement{}, err
	}
	now := s.clock.Now()
	item.Status, item.StartedAt = model.MeasurementRunning, &now
	if err := s.store.UpdateMeasurement(ctx, item); err != nil {
		return model.Measurement{}, wrap("begin measurement", err)
	}
	return item, nil
}

func (s *LabService) FinishMeasurement(ctx context.Context, id int64, value, uncertainty float64) (model.Measurement, error) {
	item, err := s.store.GetMeasurement(ctx, id)
	if err != nil {
		return model.Measurement{}, wrap("get measurement", err)
	}
	if err := s.workflow.MoveMeasurement(item.Status, model.MeasurementDone); err != nil {
		return model.Measurement{}, err
	}
	if err := validateMeasurementValue(value, uncertainty); err != nil {
		return model.Measurement{}, err
	}
	if err := s.quality.Acceptable(model.Measurement{Value: value, InputValue: item.InputValue, Uncertainty: uncertainty}); err != nil {
		item.Status = model.MeasurementRejected
		item.Value, item.Uncertainty = value, uncertainty
		_ = s.store.UpdateMeasurement(ctx, item)
		return model.Measurement{}, err
	}
	now := s.clock.Now()
	item.Value, item.Uncertainty, item.Status, item.FinishedAt = value, uncertainty, model.MeasurementDone, &now
	if err := s.store.UpdateMeasurement(ctx, item); err != nil {
		return model.Measurement{}, wrap("finish measurement", err)
	}
	return item, nil
}

func (s *LabService) GetMeasurement(ctx context.Context, id int64) (model.Measurement, error) {
	item, err := s.store.GetMeasurement(ctx, id)
	return item, wrap("get measurement", err)
}

func (s *LabService) ListMeasurements(ctx context.Context, layerID int64) ([]model.Measurement, error) {
	items, err := s.store.ListMeasurements(ctx, layerID)
	return engine.SnapshotMeasurements(items), wrap("list measurements", err)
}

func (s *LabService) RecordMeasurement(ctx context.Context, in model.MeasurementInput, value, uncertainty float64) (model.Measurement, error) {
	item, err := s.ScheduleMeasurement(ctx, in)
	if err != nil {
		return model.Measurement{}, err
	}
	if _, err := s.BeginMeasurement(ctx, item.ID); err != nil {
		return model.Measurement{}, err
	}
	return s.FinishMeasurement(ctx, item.ID, value, uncertainty)
}
