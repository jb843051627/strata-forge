package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func NormalizeOperator(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" || len(value) > 80 {
		return "", fmt.Errorf("%w: invalid operator", model.ErrInvalidInput)
	}
	return value, nil
}

func (s *LabService) RecordOperatorMeasurement(ctx context.Context, operator string, in model.MeasurementInput, value, uncertainty float64) (model.Measurement, error) {
	name, err := NormalizeOperator(operator)
	if err != nil {
		return model.Measurement{}, err
	}
	in.Operator = name
	return s.RecordMeasurement(ctx, in, value, uncertainty)
}
