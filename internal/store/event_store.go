package store

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/model"
)

const eventColumns = `id, sample_id, kind, payload, created_at`

func (s *Store) AppendEvent(ctx context.Context, item model.Event) (model.Event, error) {
	result, err := s.execContext(ctx, `INSERT INTO events(sample_id, kind, payload, created_at) VALUES(?, ?, ?, ?)`, item.SampleID, item.Kind, item.Payload, clock.Format(item.CreatedAt))
	if err != nil {
		return model.Event{}, fmt.Errorf("append event: %w", err)
	}
	item.ID, err = result.LastInsertId()
	if err != nil {
		return model.Event{}, fmt.Errorf("read event id: %w", err)
	}
	return item, nil
}

func (s *Store) ListEvents(ctx context.Context, sampleID int64) ([]model.Event, error) {
	rows, err := s.queryContext(ctx, `SELECT `+eventColumns+` FROM events WHERE sample_id = ? ORDER BY created_at, id`, sampleID)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	defer rows.Close()
	items := make([]model.Event, 0)
	for rows.Next() {
		item, scanErr := scanEvent(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan event: %w", scanErr)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) DeleteEvents(ctx context.Context, sampleID int64) error {
	if _, err := s.execContext(ctx, `DELETE FROM events WHERE sample_id = ?`, sampleID); err != nil {
		return fmt.Errorf("delete events: %w", err)
	}
	return nil
}
