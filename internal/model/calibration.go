package model

import (
	"fmt"
	"strings"
)

type CalibrationRecord struct {
	Standard string  `json:"standard"`
	Observed float64 `json:"observed"`
	Offset   float64 `json:"offset"`
	Scale    float64 `json:"scale"`
}

func (c CalibrationRecord) Apply(value float64) float64 {
	return (value - c.Offset) * c.Scale
}

func (c CalibrationRecord) Validate() error {
	if strings.TrimSpace(c.Standard) == "" || c.Scale <= 0 {
		return fmt.Errorf("%w: invalid calibration record", ErrInvalidInput)
	}
	return nil
}

func CalibrationDelta(expected, observed float64) float64 {
	return observed - expected
}

func CalibrationWithinTolerance(expected, observed, tolerance float64) bool {
	return absolute(expected-observed) <= tolerance
}

func absolute(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}
