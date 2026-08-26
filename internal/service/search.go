package service

import (
	"context"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) SearchSamples(ctx context.Context, term string, limit int) ([]model.Sample, error) {
	if strings.TrimSpace(term) == "" {
		return s.ListSamples(ctx, "")
	}
	items, err := s.store.SearchSamples(ctx, term, limit)
	return items, wrap("search samples", err)
}

func (s *LabService) FindSampleOrEmpty(ctx context.Context, term string) ([]model.Sample, error) {
	items, err := s.SearchSamples(ctx, term, 20)
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []model.Sample{}, nil
	}
	return items, nil
}
