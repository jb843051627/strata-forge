package engine

import (
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func AssignOperator(operator model.Operator, measurement model.Measurement) error {
	if !operator.CanMeasure(measurement.Kind) {
		return fmt.Errorf("%w: operator cannot perform %s", model.ErrInvalidInput, measurement.Kind)
	}
	return nil
}

func OperatorSummary(operator model.Operator) string {
	return operator.ID + ":" + operator.Name
}
