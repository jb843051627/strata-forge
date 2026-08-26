package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (s *LabService) InspectSample(ctx context.Context, sampleID int64, fn func(model.Sample) error) error {
	if fn == nil {
		return fmt.Errorf("%w: inspection callback is nil", model.ErrInvalidInput)
	}
	sample, err := s.GetSample(ctx, sampleID)
	if err != nil {
		return err
	}
	return fn(sample)
}
