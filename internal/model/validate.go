package model

import (
	"fmt"
	"strings"
)

func ValidateSampleInput(in SampleInput) error {
	if strings.TrimSpace(in.Code) == "" || strings.TrimSpace(in.Site) == "" {
		return fmt.Errorf("%w: code and site are required", ErrInvalidInput)
	}
	if in.DepthStart < 0 || in.DepthEnd <= in.DepthStart {
		return fmt.Errorf("%w: invalid depth interval", ErrInvalidInput)
	}
	if in.DepthEnd-in.DepthStart > 1000 {
		return fmt.Errorf("%w: depth interval too wide", ErrInvalidInput)
	}
	return nil
}

func ValidateLayerInput(in LayerInput) error {
	if in.Sequence <= 0 || strings.TrimSpace(in.Material) == "" {
		return fmt.Errorf("%w: sequence and material are required", ErrInvalidInput)
	}
	if in.TopDepth < 0 || in.BottomDepth <= in.TopDepth {
		return fmt.Errorf("%w: invalid layer depth", ErrInvalidInput)
	}
	return nil
}

func ValidateMeasurementInput(in MeasurementInput) error {
	if in.LayerID <= 0 || strings.TrimSpace(in.Kind) == "" || strings.TrimSpace(in.Operator) == "" {
		return fmt.Errorf("%w: measurement identity is incomplete", ErrInvalidInput)
	}
	if in.InputValue < 0 {
		return fmt.Errorf("%w: input value cannot be negative", ErrInvalidInput)
	}
	return nil
}

func ValidateReviewInput(in ReviewInput) error {
	if strings.TrimSpace(in.Reviewer) == "" {
		return fmt.Errorf("%w: reviewer is required", ErrInvalidInput)
	}
	switch in.Decision {
	case ReviewPass, ReviewHold, ReviewFail:
		return nil
	default:
		return fmt.Errorf("%w: unknown review decision", ErrInvalidInput)
	}
}
