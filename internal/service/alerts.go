package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/engine"
	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) EvaluateAlerts(ctx context.Context, measurementID int64) ([]model.Alert, error) {
	measurement, err := s.GetMeasurement(ctx, measurementID)
	if err != nil {
		return nil, err
	}
	layer, err := s.GetLayer(ctx, measurement.LayerID)
	if err != nil {
		return nil, err
	}
	sample, err := s.GetSample(ctx, layer.SampleID)
	if err != nil {
		return nil, err
	}
	rule := engine.EvaluateMeasurementAlert(measurement)
	if rule == nil {
		return []model.Alert{}, nil
	}
	alert := engine.BuildAlert(sample.ID, layer.ID, rule, sampleTime(s.clock))
	created, err := s.store.CreateAlert(ctx, alert)
	if err != nil {
		return nil, wrap("create measurement alert", err)
	}
	return []model.Alert{created}, nil
}

func (s *LabService) ListAlerts(ctx context.Context, sampleID int64, openOnly bool) ([]model.Alert, error) {
	alerts, err := s.store.ListAlerts(ctx, sampleID, openOnly)
	return engine.SnapshotAlerts(alerts), wrap("list alerts", err)
}

func (s *LabService) AcknowledgeAlert(ctx context.Context, id int64) error {
	if err := requirePositiveID(id, "alert"); err != nil {
		return err
	}
	if err := s.store.AcknowledgeAlert(ctx, id); err != nil {
		return wrap("acknowledge alert", err)
	}
	return nil
}

func (s *LabService) RaiseAlert(ctx context.Context, in model.AlertInput) (model.Alert, error) {
	if in.SampleID <= 0 || in.Severity == "" || in.Code == "" {
		return model.Alert{}, fmt.Errorf("%w: incomplete alert", model.ErrInvalidInput)
	}
	item, err := s.store.CreateAlert(ctx, model.Alert{SampleID: in.SampleID, LayerID: in.LayerID, Severity: in.Severity, Code: in.Code, Message: in.Message, CreatedAt: s.clock.Now()})
	return item, wrap("raise alert", err)
}
