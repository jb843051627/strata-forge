package store

import (
	"context"
	"fmt"
)

func (s *Store) migrate(ctx context.Context) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS samples (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            code TEXT NOT NULL UNIQUE,
            site TEXT NOT NULL,
            depth_start REAL NOT NULL,
            depth_end REAL NOT NULL,
            status TEXT NOT NULL,
            received_at TEXT NOT NULL,
            notes TEXT NOT NULL DEFAULT ''
        )`,
		`CREATE TABLE IF NOT EXISTS layers (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            sample_id INTEGER NOT NULL REFERENCES samples(id) ON DELETE CASCADE,
            sequence_no INTEGER NOT NULL,
            top_depth REAL NOT NULL,
            bottom_depth REAL NOT NULL,
            material TEXT NOT NULL,
            status TEXT NOT NULL,
            created_at TEXT NOT NULL,
            UNIQUE(sample_id, sequence_no)
        )`,
		`CREATE TABLE IF NOT EXISTS runs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            sample_id INTEGER NOT NULL REFERENCES samples(id) ON DELETE CASCADE,
            state TEXT NOT NULL,
            requested_at TEXT NOT NULL,
            started_at TEXT,
            finished_at TEXT,
            cancel_note TEXT NOT NULL DEFAULT ''
        )`,
		`CREATE TABLE IF NOT EXISTS measurements (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            layer_id INTEGER NOT NULL REFERENCES layers(id) ON DELETE CASCADE,
            kind TEXT NOT NULL,
            input_value REAL NOT NULL,
            unit TEXT NOT NULL,
            value REAL NOT NULL DEFAULT 0,
            uncertainty REAL NOT NULL DEFAULT 0,
            status TEXT NOT NULL,
            started_at TEXT,
            finished_at TEXT,
            operator TEXT NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS reviews (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            measurement_id INTEGER NOT NULL REFERENCES measurements(id) ON DELETE CASCADE,
            decision TEXT NOT NULL,
            comment TEXT NOT NULL DEFAULT '',
            reviewer TEXT NOT NULL,
            created_at TEXT NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS reports (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            sample_id INTEGER NOT NULL REFERENCES samples(id) ON DELETE CASCADE,
            version INTEGER NOT NULL,
            status TEXT NOT NULL,
            summary TEXT NOT NULL,
            age_min REAL NOT NULL,
            age_max REAL NOT NULL,
            created_at TEXT NOT NULL,
            archived_at TEXT,
            UNIQUE(sample_id, version)
        )`,
		`CREATE TABLE IF NOT EXISTS alerts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            sample_id INTEGER NOT NULL REFERENCES samples(id) ON DELETE CASCADE,
            layer_id INTEGER NOT NULL DEFAULT 0,
            severity TEXT NOT NULL,
            code TEXT NOT NULL,
            message TEXT NOT NULL,
            acknowledged INTEGER NOT NULL DEFAULT 0,
            created_at TEXT NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS events (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            sample_id INTEGER NOT NULL REFERENCES samples(id) ON DELETE CASCADE,
            kind TEXT NOT NULL,
            payload TEXT NOT NULL,
            created_at TEXT NOT NULL
        )`,
		`CREATE INDEX IF NOT EXISTS idx_layers_sample ON layers(sample_id, sequence_no)`,
		`CREATE INDEX IF NOT EXISTS idx_measurements_layer ON measurements(layer_id, status)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_sample ON alerts(sample_id, acknowledged)`,
		`CREATE INDEX IF NOT EXISTS idx_events_sample ON events(sample_id, created_at)`,
	}
	for _, statement := range statements {
		if _, err := s.execContext(ctx, statement); err != nil {
			return fmt.Errorf("migrate schema: %w", err)
		}
	}
	return nil
}
