package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

type LayerBatch struct {
	SampleID int64
	Layers   []model.LayerInput
}

func (s *LabService) IngestLayers(ctx context.Context, batch LayerBatch) ([]model.Layer, error) {
	if batch.SampleID <= 0 || len(batch.Layers) == 0 {
		return nil, fmt.Errorf("%w: layer batch is empty", model.ErrInvalidInput)
	}
	result := make([]model.Layer, 0, len(batch.Layers))
	for _, input := range batch.Layers {
		if err := contextGateIngest(ctx); err != nil {
			return result, fmt.Errorf("%w: ingest interrupted", model.ErrCancelled)
		}
		layer, err := s.AddLayer(ctx, batch.SampleID, input)
		if err != nil {
			return result, err
		}
		result = append(result, layer)
	}
	return result, nil
}

func contextGateIngest(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func ParseLayerMaterial(value string) string {
	parts := strings.Fields(value)
	if len(parts) == 0 {
		return "unknown"
	}
	return strings.ToLower(strings.Join(parts, " "))
}
