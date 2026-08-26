package engine

import (
	"fmt"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func QualitySummary(items []model.Measurement) string {
	counts := StatusCounts(items)
	kinds := KindCounts(items)
	var b strings.Builder
	fmt.Fprintf(&b, "completed=%d pending=%d rejected=%d", counts[model.MeasurementDone], counts[model.MeasurementPending], counts[model.MeasurementRejected])
	for _, kind := range SortedKinds(kinds) {
		fmt.Fprintf(&b, " %s=%d", kind, kinds[kind])
	}
	return b.String()
}

func QualityReady(items []model.Measurement) bool {
	if len(items) == 0 {
		return false
	}
	for _, item := range items {
		if item.Status != model.MeasurementDone {
			return false
		}
	}
	return true
}
