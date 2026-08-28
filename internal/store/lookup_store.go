package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *Store) SearchSamples(ctx context.Context, term string, limit int) ([]model.Sample, error) {
	term = strings.TrimSpace(term)
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	pattern := "%" + term + "%"
	rows, err := s.queryContext(ctx, `SELECT `+sampleColumns+` FROM samples WHERE code LIKE ? OR site LIKE ? ORDER BY received_at DESC, id DESC LIMIT ?`, pattern, pattern, limit)
	if err != nil {
		return nil, fmt.Errorf("search samples: %w", err)
	}
	defer rows.Close()
	items := make([]model.Sample, 0)
	for rows.Next() {
		item, scanErr := scanSample(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan search result: %w", scanErr)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) SampleIDs(ctx context.Context, status string) ([]int64, error) {
	query := `SELECT id FROM samples ORDER BY id`
	args := []any{}
	if status != "" {
		query = `SELECT id FROM samples WHERE status = ? ORDER BY id`
		args = append(args, status)
	}
	rows, err := s.queryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list sample ids: %w", err)
	}
	defer rows.Close()
	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
