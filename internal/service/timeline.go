package service

import (
	"context"
	"sort"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) Timeline(ctx context.Context, sampleID int64) ([]model.TimelineEntry, error) {
	events, err := s.store.ListEvents(ctx, sampleID)
	if err != nil {
		return nil, wrap("read sample timeline", err)
	}
	entries := make([]model.TimelineEntry, 0, len(events))
	for _, event := range events {
		entries = append(entries, model.TimelineEntry{Kind: event.Kind, Reference: event.ID, State: event.Payload, At: event.CreatedAt})
	}
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].At.Before(entries[j].At) })
	return entries, nil
}

func (s *LabService) ClearTimeline(ctx context.Context, sampleID int64) error {
	if _, err := s.GetSample(ctx, sampleID); err != nil {
		return err
	}
	return wrap("clear sample timeline", s.store.DeleteEvents(ctx, sampleID))
}
