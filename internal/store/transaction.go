package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *Store) SaveMeasurementReview(ctx context.Context, measurement model.Measurement, review model.Review, sampleID int64, sampleStatus string) error {
	return s.withTx(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `UPDATE measurements SET value = ?, uncertainty = ?, status = ?, started_at = ?, finished_at = ?, operator = ? WHERE id = ?`, measurement.Value, measurement.Uncertainty, measurement.Status, formatOptional(measurement.StartedAt), formatOptional(measurement.FinishedAt), measurement.Operator, measurement.ID); err != nil {
			return fmt.Errorf("save measurement: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO reviews(measurement_id, decision, comment, reviewer, created_at) VALUES(?, ?, ?, ?, ?)`, review.MeasurementID, review.Decision, review.Comment, review.Reviewer, clock.Format(review.CreatedAt)); err != nil {
			return fmt.Errorf("save review: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `UPDATE samples SET status = ? WHERE id = ?`, sampleStatus, sampleID); err != nil {
			return fmt.Errorf("save sample status: %w", err)
		}
		return nil
	})
}

func (s *Store) SaveReportAndEvent(ctx context.Context, report model.Report, event model.Event) (model.Report, model.Event, error) {
	var savedReport model.Report
	var savedEvent model.Event
	err := s.withTx(ctx, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `INSERT INTO reports(sample_id, version, status, summary, age_min, age_max, created_at, archived_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`, report.SampleID, report.Version, report.Status, report.Summary, report.AgeMin, report.AgeMax, clock.Format(report.CreatedAt), nullableTime(report.ArchivedAt))
		if err != nil {
			return fmt.Errorf("save report: %w", err)
		}
		savedReport = report
		savedReport.ID, err = result.LastInsertId()
		if err != nil {
			return err
		}
		result, err = tx.ExecContext(ctx, `INSERT INTO events(sample_id, kind, payload, created_at) VALUES(?, ?, ?, ?)`, event.SampleID, event.Kind, event.Payload, clock.Format(event.CreatedAt))
		if err != nil {
			return fmt.Errorf("save report event: %w", err)
		}
		savedEvent = event
		savedEvent.ID, err = result.LastInsertId()
		return err
	})
	if err != nil {
		return model.Report{}, model.Event{}, err
	}
	return savedReport, savedEvent, nil
}
