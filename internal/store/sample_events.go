package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *Store) AppendStateEvent(ctx context.Context, sampleID int64, from, to, by, note string, at string) (model.Event, error) {
	payload, err := model.EncodeStateEvent(model.StateEvent{From: from, To: to, By: by, Note: note})
	if err != nil {
		return model.Event{}, err
	}
	if at == "" {
		at = clock.Format(time.Now().UTC())
	}
	when, err := clock.Parse(at)
	if err != nil {
		return model.Event{}, fmt.Errorf("parse state event time: %w", err)
	}
	return s.AppendEvent(ctx, model.Event{SampleID: sampleID, Kind: "sample.state", Payload: payload, CreatedAt: when})
}

func (s *Store) EventsForKind(ctx context.Context, sampleID int64, kind string) ([]model.Event, error) {
	rows, err := s.queryContext(ctx, `SELECT `+eventColumns+` FROM events WHERE sample_id = ? AND kind = ? ORDER BY created_at, id`, sampleID, kind)
	if err != nil {
		return nil, fmt.Errorf("events for kind: %w", err)
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
