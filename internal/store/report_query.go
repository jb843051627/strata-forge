package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *Store) LatestReport(ctx context.Context, sampleID int64) (model.Report, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+reportColumns+` FROM reports WHERE sample_id = ? ORDER BY version DESC LIMIT 1`, sampleID)
	item, err := scanReport(row)
	if err == sql.ErrNoRows {
		return model.Report{}, model.ErrNotFound
	}
	if err != nil {
		return model.Report{}, fmt.Errorf("latest report: %w", err)
	}
	return item, nil
}

func (s *Store) ReportsByStatus(ctx context.Context, status string) ([]model.Report, error) {
	rows, err := s.queryContext(ctx, `SELECT `+reportColumns+` FROM reports WHERE status = ? ORDER BY created_at DESC`, status)
	if err != nil {
		return nil, fmt.Errorf("reports by status: %w", err)
	}
	defer rows.Close()
	items := make([]model.Report, 0)
	for rows.Next() {
		item, scanErr := scanReport(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
