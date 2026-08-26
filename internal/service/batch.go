package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

type BatchResult struct {
	Completed []model.Measurement `json:"completed"`
	Failed    []string            `json:"failed"`
	Stopped   bool                `json:"stopped"`
}

func (s *LabService) BatchRecord(ctx context.Context, inputs []model.MeasurementInput, values []float64) (BatchResult, error) {
	if len(inputs) != len(values) {
		return BatchResult{}, fmt.Errorf("%w: input/value count mismatch", model.ErrInvalidInput)
	}
	result := BatchResult{Completed: make([]model.Measurement, 0, len(inputs)), Failed: make([]string, 0)}
	for i, input := range inputs {
		if err := s.contextError(ctx); err != nil {
			result.Stopped = true
			return result, fmt.Errorf("%w: stopped after %d measurements", model.ErrCancelled, i)
		}
		item, err := s.RecordMeasurement(ctx, input, values[i], 0.05)
		if err != nil {
			result.Failed = append(result.Failed, fmt.Sprintf("layer=%d: %v", input.LayerID, err))
			continue
		}
		result.Completed = append(result.Completed, item)
	}
	return result, nil
}

func (s *LabService) ValidateBatch(inputs []model.MeasurementInput) error {
	for i, input := range inputs {
		if err := model.ValidateMeasurementInput(input); err != nil {
			return fmt.Errorf("batch item %d: %w", i, err)
		}
	}
	return nil
}
