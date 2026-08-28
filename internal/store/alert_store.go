package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/model"
)

const alertColumns = `id, sample_id, layer_id, severity, code, message, acknowledged, created_at`

func (s *Store) CreateAlert(ctx context.Context, item model.Alert) (model.Alert, error) {
	result, err := s.execContext(ctx, `INSERT INTO alerts(sample_id, layer_id, severity, code, message, acknowledged, created_at) VALUES(?, ?, ?, ?, ?, ?, ?)`, item.SampleID, item.LayerID, item.Severity, item.Code, item.Message, boolInt(item.Acknowledged), clock.Format(item.CreatedAt))
	if err != nil {
		return model.Alert{}, fmt.Errorf("create alert: %w", err)
	}
	item.ID, err = result.LastInsertId()
	if err != nil {
		return model.Alert{}, fmt.Errorf("read alert id: %w", err)
	}
	return item, nil
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func (s *Store) GetAlert(ctx context.Context, id int64) (model.Alert, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+alertColumns+` FROM alerts WHERE id = ?`, id)
	item, err := scanAlert(row)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Alert{}, model.ErrNotFound
	}
	if err != nil {
		return model.Alert{}, fmt.Errorf("get alert: %w", err)
	}
	return item, nil
}

func (s *Store) ListAlerts(ctx context.Context, sampleID int64, openOnly bool) ([]model.Alert, error) {
	query := `SELECT ` + alertColumns + ` FROM alerts WHERE sample_id = ? ORDER BY created_at DESC, id DESC`
	if openOnly {
		query = `SELECT ` + alertColumns + ` FROM alerts WHERE sample_id = ? AND acknowledged = 0 ORDER BY created_at DESC, id DESC`
	}
	rows, err := s.queryContext(ctx, query, sampleID)
	if err != nil {
		return nil, fmt.Errorf("list alerts: %w", err)
	}
	defer rows.Close()
	items := make([]model.Alert, 0)
	for rows.Next() {
		item, scanErr := scanAlert(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan alert: %w", scanErr)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) AcknowledgeAlert(ctx context.Context, id int64) error {
	result, err := s.execContext(ctx, `UPDATE alerts SET acknowledged = 1 WHERE id = ? AND acknowledged = 0`, id)
	if err != nil {
		return fmt.Errorf("acknowledge alert: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("alert %d not found: %v", id, model.ErrNotFound)
	}
	return nil
}
