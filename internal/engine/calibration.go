package engine

import (
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

type CalibrationEngine struct {
	records map[string]model.CalibrationRecord
}

func NewCalibrationEngine() *CalibrationEngine {
	return &CalibrationEngine{records: make(map[string]model.CalibrationRecord)}
}

func (e *CalibrationEngine) Register(record model.CalibrationRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	e.records[record.Standard] = record
	return nil
}

func (e *CalibrationEngine) Apply(standard string, value float64) (float64, error) {
	record, ok := e.records[standard]
	if !ok {
		return 0, fmt.Errorf("%w: calibration %s", model.ErrNotFound, standard)
	}
	return record.Apply(value), nil
}

func (e *CalibrationEngine) IsWithin(standard string, expected, observed, tolerance float64) (bool, error) {
	if _, ok := e.records[standard]; !ok {
		return false, fmt.Errorf("%w: calibration %s", model.ErrNotFound, standard)
	}
	return model.CalibrationWithinTolerance(expected, observed, tolerance), nil
}

func (e *CalibrationEngine) Standards() []string {
	values := make([]string, 0, len(e.records))
	for key := range e.records {
		values = append(values, key)
	}
	return values
}
