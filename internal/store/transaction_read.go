package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *Store) LoadSampleGraph(ctx context.Context, sampleID int64) (model.Sample, []model.Layer, []model.Measurement, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var sample model.Sample
	row := s.db.QueryRowContext(ctx, `SELECT `+sampleColumns+` FROM samples WHERE id = ?`, sampleID)
	var err error
	sample, err = scanSample(row)
	if err == sql.ErrNoRows {
		return model.Sample{}, nil, nil, model.ErrNotFound
	}
	if err != nil {
		return model.Sample{}, nil, nil, fmt.Errorf("load graph sample: %w", err)
	}
	layerRows, err := s.db.QueryContext(ctx, `SELECT `+layerColumns+` FROM layers WHERE sample_id = ? ORDER BY sequence_no`, sampleID)
	if err != nil {
		return model.Sample{}, nil, nil, err
	}
	defer layerRows.Close()
	layers := make([]model.Layer, 0)
	measurements := make([]model.Measurement, 0)
	for layerRows.Next() {
		layer, scanErr := scanLayer(layerRows)
		if scanErr != nil {
			return model.Sample{}, nil, nil, scanErr
		}
		layers = append(layers, layer)
		measurementRows, queryErr := s.db.QueryContext(ctx, `SELECT `+measurementColumns+` FROM measurements WHERE layer_id = ? ORDER BY id`, layer.ID)
		if queryErr != nil {
			return model.Sample{}, nil, nil, queryErr
		}
		for measurementRows.Next() {
			item, measurementErr := scanMeasurement(measurementRows)
			if measurementErr != nil {
				measurementRows.Close()
				return model.Sample{}, nil, nil, measurementErr
			}
			measurements = append(measurements, item)
		}
		measurementRows.Close()
	}
	return sample, layers, measurements, layerRows.Err()
}
