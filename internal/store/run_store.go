package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/model"
)

const runColumns = `id, sample_id, state, requested_at, started_at, finished_at, cancel_note`

func (s *Store) CreateRun(ctx context.Context, item model.Run) (model.Run, error) {
	result, err := s.execContext(ctx, `INSERT INTO runs(sample_id, state, requested_at, started_at, finished_at, cancel_note) VALUES(?, ?, ?, ?, ?, ?)`, item.SampleID, item.State, clock.Format(item.RequestedAt), nil, nil, item.CancelNote)
	if err != nil {
		return model.Run{}, fmt.Errorf("create run: %w", err)
	}
	item.ID, err = result.LastInsertId()
	if err != nil {
		return model.Run{}, fmt.Errorf("read run id: %w", err)
	}
	return item, nil
}

func (s *Store) GetRun(ctx context.Context, id int64) (model.Run, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+runColumns+` FROM runs WHERE id = ?`, id)
	item, err := scanRun(row)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Run{}, model.ErrNotFound
	}
	if err != nil {
		return model.Run{}, fmt.Errorf("get run: %w", err)
	}
	return item, nil
}

func (s *Store) ListRuns(ctx context.Context, sampleID int64) ([]model.Run, error) {
	rows, err := s.queryContext(ctx, `SELECT `+runColumns+` FROM runs WHERE sample_id = ? ORDER BY requested_at, id`, sampleID)
	if err != nil {
		return nil, fmt.Errorf("list runs: %w", err)
	}
	defer rows.Close()
	items := make([]model.Run, 0)
	for rows.Next() {
		item, scanErr := scanRun(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan run: %w", scanErr)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) UpdateRun(ctx context.Context, item model.Run) error {
	var started, finished any
	if item.StartedAt != nil {
		started = clock.Format(*item.StartedAt)
	}
	if item.FinishedAt != nil {
		finished = clock.Format(*item.FinishedAt)
	}
	result, err := s.execContext(ctx, `UPDATE runs SET state = ?, started_at = ?, finished_at = ?, cancel_note = ? WHERE id = ?`, item.State, started, finished, item.CancelNote, item.ID)
	if err != nil {
		return fmt.Errorf("update run: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read run update count: %w", err)
	}
	if count == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (s *Store) FindActiveRun(ctx context.Context, sampleID int64) (model.Run, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+runColumns+` FROM runs WHERE sample_id = ? AND state IN (?, ?) ORDER BY id DESC LIMIT 1`, sampleID, model.RunQueued, model.RunActive)
	item, err := scanRun(row)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Run{}, model.ErrNotFound
	}
	if err != nil {
		return model.Run{}, fmt.Errorf("find active run: %w", err)
	}
	return item, nil
}

func (s *Store) UpdateRunAndSample(ctx context.Context, run model.Run, sampleID int64, sampleStatus string) error {
	if _, err := s.execContext(ctx, `UPDATE runs SET state = ?, cancel_note = ? WHERE id = ?`, run.State, run.CancelNote, run.ID); err != nil {
		return err
	}
	result, err := s.execContext(ctx, `UPDATE samples SET status = ? WHERE id = ?`, sampleStatus, sampleID)
	if err != nil {
		return err
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
