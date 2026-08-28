package service

import (
	"fmt"
	"sync"

	"github.com/jb843051627/strata-forge/internal/engine"
	"github.com/jb843051627/strata-forge/internal/model"
)

type MeasurementSeries struct {
	mu     sync.RWMutex
	values map[int64][]float64
}

func NewMeasurementSeries() *MeasurementSeries {
	return &MeasurementSeries{values: make(map[int64][]float64)}
}

func (s *MeasurementSeries) Append(layerID int64, value float64) error {
	if layerID <= 0 {
		return fmt.Errorf("layer id must be positive")
	}
	s.mu.Lock()
	s.values[layerID] = append(s.values[layerID], value)
	s.mu.Unlock()
	return nil
}

func (s *MeasurementSeries) Snapshot(layerID int64) []float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]float64, len(s.values[layerID]))
	copy(result, s.values[layerID])
	return result
}

func (s *MeasurementSeries) Smooth(layerID int64, radius int) []float64 {
	return engine.Smooth(s.Snapshot(layerID), radius)
}

func (s *MeasurementSeries) Variance(layerID int64) float64 {
	values := s.Snapshot(layerID)
	return engine.SeriesVariance(valuesToMeasurements(values))
}

func valuesToMeasurements(values []float64) []model.Measurement {
	items := make([]model.Measurement, len(values))
	for i, value := range values {
		items[i].ID = int64(i + 1)
		items[i].Value = value
	}
	return items
}
