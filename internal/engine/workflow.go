package engine

import (
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

type Workflow struct{}

func NewWorkflow() *Workflow { return &Workflow{} }

func (Workflow) MoveSample(from, to string) error {
	if !model.CanSampleTransition(from, to) {
		return fmt.Errorf("%w: sample %s -> %s", model.ErrInvalidTransition, from, to)
	}
	return nil
}

func (Workflow) MoveRun(from, to string) error {
	if !model.CanRunTransition(from, to) {
		return fmt.Errorf("%w: run %s -> %s", model.ErrInvalidTransition, from, to)
	}
	return nil
}

func (Workflow) MoveMeasurement(from, to string) error {
	if !model.CanMeasurementTransition(from, to) {
		return fmt.Errorf("%w: measurement %s -> %s", model.ErrInvalidTransition, from, to)
	}
	return nil
}

func (Workflow) ReviewStatus(decision string) string {
	switch decision {
	case model.ReviewPass:
		return model.SampleReview
	case model.ReviewFail:
		return model.SampleRejected
	default:
		return model.SampleReview
	}
}
