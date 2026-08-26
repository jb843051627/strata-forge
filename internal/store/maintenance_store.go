package store

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *Store) MarkStaleRuns(ctx context.Context, before string) (int64, error) {
	result, err := s.execContext(ctx, `UPDATE runs SET state = ?, cancel_note = ? WHERE state IN (?, ?) AND requested_at < ?`, model.RunFailed, "worker timeout", model.RunQueued, model.RunActive, before)
	if err != nil {
		return 0, fmt.Errorf("mark stale runs: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) RemoveUnacknowledgedAlerts(ctx context.Context, sampleID int64) error {
	_, err := s.execContext(ctx, `DELETE FROM alerts WHERE sample_id = ? AND acknowledged = 1`, sampleID)
	if err != nil {
		return fmt.Errorf("remove acknowledged alerts: %w", err)
	}
	return nil
}
