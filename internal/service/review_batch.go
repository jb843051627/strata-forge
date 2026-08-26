package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) BatchReview(ctx context.Context, measurementIDs []int64, input model.ReviewInput) (int, error) {
	completed := 0
	for _, id := range measurementIDs {
		if err := s.contextError(ctx); err != nil {
			return completed, fmt.Errorf("%w: review batch stopped", model.ErrCancelled)
		}
		if _, err := s.ReviewMeasurement(ctx, id, input); err != nil {
			return completed, err
		}
		completed++
	}
	return completed, nil
}
