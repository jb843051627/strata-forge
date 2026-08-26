package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *Store) EventsSince(ctx context.Context, sampleID int64, since time.Time) ([]model.Event, error) {
	rows, err := s.queryContext(ctx, `SELECT `+eventColumns+` FROM events WHERE sample_id = ? AND created_at >= ? ORDER BY created_at, id`, sampleID, clock.Format(since))
	if err != nil {
		return nil, fmt.Errorf("events since: %w", err)
	}
	defer rows.Close()
	items := make([]model.Event, 0)
	for rows.Next() {
		item, scanErr := scanEvent(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) EventKinds(ctx context.Context, sampleID int64) ([]string, error) {
	rows, err := s.queryContext(ctx, `SELECT DISTINCT kind FROM events WHERE sample_id = ? ORDER BY kind`, sampleID)
	if err != nil {
		return nil, fmt.Errorf("event kinds: %w", err)
	}
	defer rows.Close()
	result := make([]string, 0)
	for rows.Next() {
		var kind string
		if err := rows.Scan(&kind); err != nil {
			return nil, err
		}
		result = append(result, kind)
	}
	return result, rows.Err()
}
