package store

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

type MeasurementStats struct {
	Done        int
	Pending     int
	Rejected    int
	Mean        float64
	Uncertainty float64
}

func (s *Store) MeasurementStats(ctx context.Context, sampleID int64) (MeasurementStats, error) {
	var result MeasurementStats
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(CASE WHEN status = ? THEN 1 END), COUNT(CASE WHEN status = ? THEN 1 END), COUNT(CASE WHEN status = ? THEN 1 END), COALESCE(AVG(value), 0), COALESCE(AVG(uncertainty), 0) FROM measurements WHERE layer_id IN (SELECT id FROM layers WHERE sample_id = ?)`, model.MeasurementDone, model.MeasurementPending, model.MeasurementRejected, sampleID).Scan(&result.Done, &result.Pending, &result.Rejected, &result.Mean, &result.Uncertainty); err != nil {
		return MeasurementStats{}, fmt.Errorf("measurement stats: %w", err)
	}
	return result, nil
}

func (s *Store) MeasurementKinds(ctx context.Context, sampleID int64) ([]string, error) {
	rows, err := s.queryContext(ctx, `SELECT DISTINCT kind FROM measurements WHERE layer_id IN (SELECT id FROM layers WHERE sample_id = ?) ORDER BY kind`, sampleID)
	if err != nil {
		return nil, fmt.Errorf("measurement kinds: %w", err)
	}
	defer rows.Close()
	kinds := make([]string, 0)
	for rows.Next() {
		var kind string
		if err := rows.Scan(&kind); err != nil {
			return nil, err
		}
		kinds = append(kinds, kind)
	}
	return kinds, rows.Err()
}
