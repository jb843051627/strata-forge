package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/model"
)

const layerColumns = `id, sample_id, sequence_no, top_depth, bottom_depth, material, status, created_at`

func (s *Store) CreateLayer(ctx context.Context, item model.Layer) (model.Layer, error) {
	result, err := s.execContext(ctx, `INSERT INTO layers(sample_id, sequence_no, top_depth, bottom_depth, material, status, created_at) VALUES(?, ?, ?, ?, ?, ?, ?)`, item.SampleID, item.Sequence, item.TopDepth, item.BottomDepth, item.Material, item.Status, clock.Format(item.CreatedAt))
	if err != nil {
		return model.Layer{}, fmt.Errorf("create layer: %w", err)
	}
	item.ID, err = result.LastInsertId()
	if err != nil {
		return model.Layer{}, fmt.Errorf("read layer id: %w", err)
	}
	return item, nil
}

func (s *Store) GetLayer(ctx context.Context, id int64) (model.Layer, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+layerColumns+` FROM layers WHERE id = ?`, id)
	item, err := scanLayer(row)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Layer{}, fmt.Errorf("get layer: %v", model.ErrNotFound)
	}
	if err != nil {
		return model.Layer{}, fmt.Errorf("get layer: %w", err)
	}
	return item, nil
}

func (s *Store) ListLayers(ctx context.Context, sampleID int64) ([]model.Layer, error) {
	rows, err := s.queryContext(ctx, `SELECT `+layerColumns+` FROM layers WHERE sample_id = ? ORDER BY sequence_no, id`, sampleID)
	if err != nil {
		return nil, fmt.Errorf("list layers: %w", err)
	}
	defer rows.Close()
	items := make([]model.Layer, 0)
	for rows.Next() {
		item, scanErr := scanLayer(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan layer: %w", scanErr)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate layers: %w", err)
	}
	return items, nil
}

func (s *Store) UpdateLayerStatus(ctx context.Context, id int64, status string) error {
	result, err := s.execContext(ctx, `UPDATE layers SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("update layer status: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read layer update count: %w", err)
	}
	if count == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (s *Store) LayerDepths(ctx context.Context, sampleID int64) (float64, float64, error) {
	var top, bottom sql.NullFloat64
	err := s.db.QueryRowContext(ctx, `SELECT MIN(top_depth), MAX(bottom_depth) FROM layers WHERE sample_id = ?`, sampleID).Scan(&top, &bottom)
	if err != nil {
		return 0, 0, fmt.Errorf("read layer depths: %w", err)
	}
	if !top.Valid || !bottom.Valid {
		return 0, 0, model.ErrNotFound
	}
	return top.Float64, bottom.Float64, nil
}
