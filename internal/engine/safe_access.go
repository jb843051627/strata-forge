package engine

import (
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func FormatProvenance(provenance *model.Provenance) string {
	if provenance == nil {
		return "unknown provenance"
	}
	return ProvenanceLabel(*provenance)
}

func TakeLatest(values []float64, count int) ([]float64, error) {
	if count < 0 {
		return nil, fmt.Errorf("%w: negative latest count", model.ErrInvalidInput)
	}
	start := len(values) - count
	return values[start:], nil
}
