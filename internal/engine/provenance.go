package engine

import (
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func ValidateProvenance(p model.Provenance) error {
	if !p.Normalize().Complete() {
		return fmt.Errorf("%w: measurement provenance incomplete", model.ErrInvalidInput)
	}
	return nil
}

func ProvenanceLabel(p model.Provenance) string {
	p = p.Normalize()
	return p.Source + " / " + p.Instrument
}
