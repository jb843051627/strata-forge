package engine

import (
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func Reviewable(item model.Measurement) error {
	if item.Status != model.MeasurementDone {
		return fmt.Errorf("%w: measurement %d is %s", model.ErrQualityHold, item.ID, item.Status)
	}
	return nil
}

func ReviewOutcome(item model.Review) string {
	if model.ReviewAllowsReport(item.Decision) {
		return "accepted"
	}
	if model.ReviewNeedsFollowup(item.Decision) {
		return "followup"
	}
	return "unknown"
}
