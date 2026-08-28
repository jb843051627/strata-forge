package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/model"
)

const measurementColumns = `id, layer_id, kind, input_value, unit, value, uncertainty, status, started_at, finished_at, operator`

func (s *Store) CreateMeasurement(ctx context.Context, item model.Measurement) (model.Measurement, error) {
	result, err := s.execContext(ctx, `INSERT INTO measurements(layer_id, kind, input_value, unit, value, uncertainty, status, started_at, finished_at, operator) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, item.LayerID, item.Kind, item.InputValue, item.Unit, item.Value, item.Uncertainty, item.Status, nullableTime(item.StartedAt), nullableTime(item.FinishedAt), item.Operator)
	if err != nil {
		return model.Measurement{}, fmt.Errorf("create measurement: %w", err)
	}
	item.ID, err = result.LastInsertId()
	if err != nil {
		return model.Measurement{}, fmt.Errorf("read measurement id: %w", err)
	}
	return item, nil
}

func nullableTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return clock.Format(*value)
}

func (s *Store) GetMeasurement(ctx context.Context, id int64) (model.Measurement, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+measurementColumns+` FROM measurements WHERE id = ?`, id)
	item, err := scanMeasurement(row)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Measurement{}, model.ErrNotFound
	}
	if err != nil {
		return model.Measurement{}, fmt.Errorf("get measurement: %w", err)
	}
	return item, nil
}

func (s *Store) ListMeasurements(ctx context.Context, layerID int64) ([]model.Measurement, error) {
	rows, err := s.queryContext(ctx, `SELECT `+measurementColumns+` FROM measurements WHERE layer_id = ? ORDER BY id`, layerID)
	if err != nil {
		return nil, fmt.Errorf("list measurements: %w", err)
	}
	defer rows.Close()
	items := make([]model.Measurement, 0)
	for rows.Next() {
		item, scanErr := scanMeasurement(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan measurement: %w", scanErr)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate measurements: %w", err)
	}
	return items, nil
}

func (s *Store) ListMeasurementsForSample(ctx context.Context, sampleID int64) ([]model.Measurement, error) {
	rows, err := s.queryContext(ctx, `SELECT `+measurementColumns+` FROM measurements WHERE layer_id IN (SELECT id FROM layers WHERE sample_id = ?) ORDER BY id`, sampleID)
	if err != nil {
		return nil, fmt.Errorf("list sample measurements: %w", err)
	}
	defer rows.Close()
	items := make([]model.Measurement, 0)
	for rows.Next() {
		item, scanErr := scanMeasurement(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan sample measurement: %w", scanErr)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) UpdateMeasurement(ctx context.Context, item model.Measurement) error {
	result, err := s.execContext(ctx, `UPDATE measurements SET value = ?, uncertainty = ?, status = ?, started_at = ?, finished_at = ?, operator = ? WHERE id = ?`, item.Value, item.Uncertainty, item.Status, formatOptional(item.StartedAt), formatOptional(item.FinishedAt), item.Operator, item.ID)
	if err != nil {
		return fmt.Errorf("update measurement: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return model.ErrNotFound
	}
	return nil
}

func formatOptional(value *time.Time) any {
	if value == nil {
		return nil
	}
	return clock.Format(*value)
}

func (s *Store) CountMeasurements(ctx context.Context, sampleID int64, status string) (int, error) {
	var count int
	query := `SELECT COUNT(1) FROM measurements WHERE layer_id IN (SELECT id FROM layers WHERE sample_id = ?)`
	args := []any{sampleID}
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	if err := s.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("count measurements: %w", err)
	}
	return count, nil
}
