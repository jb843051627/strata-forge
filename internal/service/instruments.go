package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/engine"
	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) RegisterInstrument(item model.Instrument) error {
	return wrap("register instrument", s.instruments.Register(item))
}

func (s *LabService) MeasurementInstrument(ctx context.Context, measurementID int64) (string, error) {
	measurement, err := s.GetMeasurement(ctx, measurementID)
	if err != nil {
		return "", err
	}
	instrument, err := s.instruments.Find(measurement.Kind)
	if err != nil {
		return "", wrap("find instrument", err)
	}
	if instrument == nil {
		return "", fmt.Errorf("%w: no instrument for %s", model.ErrNotFound, measurement.Kind)
	}
	return engine.InstrumentLabel(instrument), nil
}

func (s *LabService) Instruments() []model.Instrument {
	return s.instruments.List()
}
