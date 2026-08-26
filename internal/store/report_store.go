package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/model"
)

const reportColumns = `id, sample_id, version, status, summary, age_min, age_max, created_at, archived_at`

func (s *Store) CreateReport(ctx context.Context, item model.Report) (model.Report, error) {
	result, err := s.execContext(ctx, `INSERT INTO reports(sample_id, version, status, summary, age_min, age_max, created_at, archived_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`, item.SampleID, item.Version, item.Status, item.Summary, item.AgeMin, item.AgeMax, clock.Format(item.CreatedAt), nullableTime(item.ArchivedAt))
	if err != nil {
		return model.Report{}, fmt.Errorf("create report: %w", err)
	}
	item.ID, err = result.LastInsertId()
	if err != nil {
		return model.Report{}, fmt.Errorf("read report id: %w", err)
	}
	return item, nil
}

func (s *Store) GetReport(ctx context.Context, id int64) (model.Report, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+reportColumns+` FROM reports WHERE id = ?`, id)
	item, err := scanReport(row)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Report{}, fmt.Errorf("get report: %v", model.ErrNotFound)
	}
	if err != nil {
		return model.Report{}, fmt.Errorf("get report: %w", err)
	}
	return item, nil
}

func (s *Store) ListReports(ctx context.Context, sampleID int64) ([]model.Report, error) {
	rows, err := s.queryContext(ctx, `SELECT `+reportColumns+` FROM reports WHERE sample_id = ? ORDER BY version DESC`, sampleID)
	if err != nil {
		return nil, fmt.Errorf("list reports: %w", err)
	}
	defer rows.Close()
	items := make([]model.Report, 0)
	for rows.Next() {
		item, scanErr := scanReport(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan report: %w", scanErr)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) NextReportVersion(ctx context.Context, sampleID int64) (int, error) {
	var version int
	if err := s.db.QueryRowContext(ctx, `SELECT COALESCE(MAX(version), 0) + 1 FROM reports WHERE sample_id = ?`, sampleID).Scan(&version); err != nil {
		return 0, fmt.Errorf("next report version: %w", err)
	}
	return version, nil
}

func (s *Store) ArchiveReport(ctx context.Context, id int64, at string) error {
	result, err := s.execContext(ctx, `UPDATE reports SET status = ?, archived_at = ? WHERE id = ? AND status <> ?`, model.ReportArchived, at, id, model.ReportArchived)
	if err != nil {
		return fmt.Errorf("archive report: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return model.ErrAlreadyArchived
	}
	return nil
}
