package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/model"
)

const reviewColumns = `id, measurement_id, decision, comment, reviewer, created_at`

func (s *Store) CreateReview(ctx context.Context, item model.Review) (model.Review, error) {
	result, err := s.execContext(ctx, `INSERT INTO reviews(measurement_id, decision, comment, reviewer, created_at) VALUES(?, ?, ?, ?, ?)`, item.MeasurementID, item.Decision, item.Comment, item.Reviewer, clock.Format(item.CreatedAt))
	if err != nil {
		return model.Review{}, fmt.Errorf("create review: %w", err)
	}
	item.ID, err = result.LastInsertId()
	if err != nil {
		return model.Review{}, fmt.Errorf("read review id: %w", err)
	}
	return item, nil
}

func (s *Store) GetReview(ctx context.Context, id int64) (model.Review, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+reviewColumns+` FROM reviews WHERE id = ?`, id)
	item, err := scanReview(row)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Review{}, model.ErrNotFound
	}
	if err != nil {
		return model.Review{}, fmt.Errorf("get review: %w", err)
	}
	return item, nil
}

func (s *Store) ListReviewsForSample(ctx context.Context, sampleID int64) ([]model.Review, error) {
	rows, err := s.queryContext(ctx, `SELECT `+reviewColumns+` FROM reviews WHERE measurement_id IN (SELECT m.id FROM measurements m JOIN layers l ON l.id = m.layer_id WHERE l.sample_id = ?) ORDER BY id`, sampleID)
	if err != nil {
		return nil, fmt.Errorf("list reviews: %w", err)
	}
	defer rows.Close()
	items := make([]model.Review, 0)
	for rows.Next() {
		item, scanErr := scanReview(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan review: %w", scanErr)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) LatestReview(ctx context.Context, measurementID int64) (model.Review, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+reviewColumns+` FROM reviews WHERE measurement_id = ? ORDER BY id DESC LIMIT 1`, measurementID)
	item, err := scanReview(row)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Review{}, fmt.Errorf("latest review: %v", model.ErrNotFound)

	}
	if err != nil {
		return model.Review{}, fmt.Errorf("latest review: %w", err)
	}
	return item, nil
}
