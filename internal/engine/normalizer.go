package engine

import (
	"math"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func NormalizeKind(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "magnetic", "magnetism", "magnetic susceptibility":
		return "magnetic"
	case "isotope", "isotopic", "isotope ratio":
		return "isotope"
	case "grain", "grain-size", "grain size":
		return "grain"
	default:
		return "other"
	}
}

func NormalizeMeasurement(input model.MeasurementInput) model.MeasurementInput {
	input.Kind = NormalizeKind(input.Kind)
	input.Unit = strings.ToLower(strings.TrimSpace(input.Unit))
	input.Operator = strings.TrimSpace(input.Operator)
	if math.IsNaN(input.InputValue) || math.IsInf(input.InputValue, 0) {
		input.InputValue = 0
	}
	return input
}
