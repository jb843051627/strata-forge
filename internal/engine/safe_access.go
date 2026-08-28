package engine

import (
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func FormatProvenance(provenance *model.Provenance) string {
	return provenance.Source + " / " + provenance.Instrument
}

func TakeLatest(values []float64, count int) ([]float64, error) {
	if count < 0 || count > len(values) {
		return nil, fmt.Errorf("%w: invalid latest count", model.ErrInvalidInput)
	}
	start := len(values) - count
	return append([]float64(nil), values[start:]...), nil
}
