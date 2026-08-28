package store

import (
	"context"
	"fmt"
)

type Stats struct {
	Samples       int
	Layers        int
	Measurements  int
	Reviews       int
	OpenAlerts    int
	ArchivedFiles int
}

func (s *Store) Stats(ctx context.Context) (Stats, error) {
	var result Stats
	checks := []struct {
		query string
		value *int
	}{
		{`SELECT COUNT(1) FROM samples`, &result.Samples},
		{`SELECT COUNT(1) FROM layers`, &result.Layers},
		{`SELECT COUNT(1) FROM measurements`, &result.Measurements},
		{`SELECT COUNT(1) FROM reviews`, &result.Reviews},
		{`SELECT COUNT(1) FROM alerts WHERE acknowledged = 0`, &result.OpenAlerts},
		{`SELECT COUNT(1) FROM reports WHERE status = 'archived'`, &result.ArchivedFiles},
	}
	for _, check := range checks {
		if err := s.db.QueryRowContext(ctx, check.query).Scan(check.value); err != nil {
			return Stats{}, fmt.Errorf("stats query: %w", err)
		}
	}
	return result, nil
}

func (s *Store) CountOpenAlerts(ctx context.Context, sampleID int64) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM alerts WHERE sample_id = ? AND acknowledged = 0`, sampleID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count open alerts: %w", err)
	}
	return count, nil
}
