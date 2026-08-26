package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/model"
)

func parseStoredTime(raw string) (time.Time, error) {
	value, err := clock.Parse(raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse timestamp %q: %w", raw, err)
	}
	return value, nil
}

func scanSample(row interface{ Scan(...any) error }) (model.Sample, error) {
	var item model.Sample
	var received string
	if err := row.Scan(&item.ID, &item.Code, &item.Site, &item.DepthStart, &item.DepthEnd, &item.Status, &received, &item.Notes); err != nil {
		return item, err
	}
	var err error
	item.ReceivedAt, err = parseStoredTime(received)
	return item, err
}

func scanLayer(row interface{ Scan(...any) error }) (model.Layer, error) {
	var item model.Layer
	var created string
	if err := row.Scan(&item.ID, &item.SampleID, &item.Sequence, &item.TopDepth, &item.BottomDepth, &item.Material, &item.Status, &created); err != nil {
		return item, err
	}
	var err error
	item.CreatedAt, err = parseStoredTime(created)
	return item, err
}

func scanOptional(raw sql.NullString) (*time.Time, error) {
	if !raw.Valid || raw.String == "" {
		return nil, nil
	}
	t, err := parseStoredTime(raw.String)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func scanMeasurement(row interface{ Scan(...any) error }) (model.Measurement, error) {
	var item model.Measurement
	var started, finished sql.NullString
	if err := row.Scan(&item.ID, &item.LayerID, &item.Kind, &item.InputValue, &item.Unit, &item.Value, &item.Uncertainty, &item.Status, &started, &finished, &item.Operator); err != nil {
		return item, err
	}
	var err error
	if item.StartedAt, err = scanOptional(started); err != nil {
		return item, err
	}
	if item.FinishedAt, err = scanOptional(finished); err != nil {
		return item, err
	}
	return item, nil
}

func scanRun(row interface{ Scan(...any) error }) (model.Run, error) {
	var item model.Run
	var requested, started, finished sql.NullString
	if err := row.Scan(&item.ID, &item.SampleID, &item.State, &requested, &started, &finished, &item.CancelNote); err != nil {
		return item, err
	}
	var err error
	if item.RequestedAt, err = parseStoredTime(requested.String); err != nil {
		return item, err
	}
	if item.StartedAt, err = scanOptional(started); err != nil {
		return item, err
	}
	if item.FinishedAt, err = scanOptional(finished); err != nil {
		return item, err
	}
	return item, nil
}

func scanReview(row interface{ Scan(...any) error }) (model.Review, error) {
	var item model.Review
	var created string
	if err := row.Scan(&item.ID, &item.MeasurementID, &item.Decision, &item.Comment, &item.Reviewer, &created); err != nil {
		return item, err
	}
	var err error
	item.CreatedAt, err = parseStoredTime(created)
	return item, err
}

func scanReport(row interface{ Scan(...any) error }) (model.Report, error) {
	var item model.Report
	var created, archived sql.NullString
	if err := row.Scan(&item.ID, &item.SampleID, &item.Version, &item.Status, &item.Summary, &item.AgeMin, &item.AgeMax, &created, &archived); err != nil {
		return item, err
	}
	var err error
	if item.CreatedAt, err = parseStoredTime(created.String); err != nil {
		return item, err
	}
	if item.ArchivedAt, err = scanOptional(archived); err != nil {
		return item, err
	}
	return item, nil
}

func scanAlert(row interface{ Scan(...any) error }) (model.Alert, error) {
	var item model.Alert
	var acknowledged int
	var created string
	if err := row.Scan(&item.ID, &item.SampleID, &item.LayerID, &item.Severity, &item.Code, &item.Message, &acknowledged, &created); err != nil {
		return item, err
	}
	item.Acknowledged = acknowledged != 0
	var err error
	item.CreatedAt, err = parseStoredTime(created)
	return item, err
}

func scanEvent(row interface{ Scan(...any) error }) (model.Event, error) {
	var item model.Event
	var created string
	if err := row.Scan(&item.ID, &item.SampleID, &item.Kind, &item.Payload, &created); err != nil {
		return item, err
	}
	var err error
	item.CreatedAt, err = parseStoredTime(created)
	return item, err
}
