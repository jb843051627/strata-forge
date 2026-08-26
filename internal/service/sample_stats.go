package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

type SampleStats struct {
	DepthCoverage float64
	OpenAlerts    int
	Kinds         []string
	Timeline      []model.TimelineEntry
}

func (s *LabService) SampleStats(ctx context.Context, sampleID int64) (SampleStats, error) {
	coverage, err := s.LayerCoverage(ctx, sampleID)
	if err != nil {
		return SampleStats{}, err
	}
	alerts, err := s.store.CountOpenAlerts(ctx, sampleID)
	if err != nil {
		return SampleStats{}, wrap("sample alert stats", err)
	}
	kinds, err := s.store.EventKinds(ctx, sampleID)
	if err != nil {
		return SampleStats{}, wrap("sample event stats", err)
	}
	timeline, err := s.Timeline(ctx, sampleID)
	if err != nil {
		return SampleStats{}, err
	}
	return SampleStats{DepthCoverage: coverage, OpenAlerts: alerts, Kinds: kinds, Timeline: timeline}, nil
}

func FormatStats(stats SampleStats) string {
	return fmt.Sprintf("coverage=%.3f open_alerts=%d events=%d", stats.DepthCoverage, stats.OpenAlerts, len(stats.Timeline))
}
