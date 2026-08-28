package service

import (
	"fmt"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func validateMeasurementValue(value float64, uncertainty float64) error {
	if value < 0 || uncertainty < 0 {
		return fmt.Errorf("%w: measurement values must be non-negative", model.ErrInvalidInput)
	}
	return nil
}

func validateSampleLayer(sample model.Sample, layer model.LayerInput) error {
	if layer.TopDepth < sample.DepthStart || layer.BottomDepth > sample.DepthEnd {
		return fmt.Errorf("%w: layer is outside sample interval", model.ErrInvalidInput)
	}
	if layer.BottomDepth-layer.TopDepth < 0.01 {
		return fmt.Errorf("%w: layer is too thin", model.ErrInvalidInput)
	}
	return nil
}

func validateReviewComment(in model.ReviewInput) error {
	if in.Decision == model.ReviewHold && strings.TrimSpace(in.Comment) == "" {
		return fmt.Errorf("%w: hold decisions need a comment", model.ErrInvalidInput)
	}
	return nil
}
