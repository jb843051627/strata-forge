package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/model"
)

const sampleColumns = `id, code, site, depth_start, depth_end, status, received_at, notes`

func (s *Store) CreateSample(ctx context.Context, item model.Sample) (model.Sample, error) {
	result, err := s.execContext(context.Background(), `INSERT INTO samples(code, site, depth_start, depth_end, status, received_at, notes) VALUES(?, ?, ?, ?, ?, ?, ?)`, item.Code, item.Site, item.DepthStart, item.DepthEnd, item.Status, clock.Format(item.ReceivedAt), item.Notes)
	if err != nil {
		return model.Sample{}, fmt.Errorf("create sample: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return model.Sample{}, fmt.Errorf("read sample id: %w", err)
	}
	item.ID = id
	return item, nil
}

func (s *Store) GetSample(ctx context.Context, id int64) (model.Sample, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+sampleColumns+` FROM samples WHERE id = ?`, id)
	item, err := scanSample(row)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Sample{}, model.ErrNotFound
	}
	if err != nil {
		return model.Sample{}, fmt.Errorf("get sample: %w", err)
	}
	return item, nil
}

func (s *Store) FindSampleByCode(ctx context.Context, code string) (model.Sample, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+sampleColumns+` FROM samples WHERE code = ?`, code)
	item, err := scanSample(row)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Sample{}, model.ErrNotFound
	}
	if err != nil {
		return model.Sample{}, fmt.Errorf("find sample: %w", err)
	}
	return item, nil
}

func (s *Store) ListSamples(ctx context.Context, status string) ([]model.Sample, error) {
	query := `SELECT ` + sampleColumns + ` FROM samples ORDER BY received_at DESC, id DESC`
	args := []any{}
	if status != "" {
		query = `SELECT ` + sampleColumns + ` FROM samples WHERE status = ? ORDER BY received_at DESC, id DESC`
		args = append(args, status)
	}
	rows, err := s.queryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list samples: %w", err)
	}
	defer rows.Close()
	items := make([]model.Sample, 0)
	for rows.Next() {
		item, scanErr := scanSample(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan sample: %w", scanErr)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate samples: %w", err)
	}
	return items, nil
}

func (s *Store) UpdateSampleStatus(ctx context.Context, id int64, status string) error {
	result, err := s.execContext(ctx, `UPDATE samples SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("update sample status: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read sample update count: %w", err)
	}
	if count == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (s *Store) CountSamples(ctx context.Context, status string) (int, error) {
	var count int
	query := `SELECT COUNT(1) FROM samples`
	args := []any{}
	if status != "" {
		query += ` WHERE status = ?`
		args = append(args, status)
	}
	if err := s.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("count samples: %w", err)
	}
	return count, nil
}
