package engine

import (
	"fmt"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func StateLabel(value string) string {
	labels := map[string]string{
		model.SampleReceived: "待分层",
		model.SampleLayered:  "待测量",
		model.SampleRunning:  "测量中",
		model.SampleReview:   "质量复核",
		model.SampleArchived: "已归档",
		model.SampleRejected: "已退回",
	}
	if label, ok := labels[value]; ok {
		return label
	}
	return strings.ToUpper(value)
}

func MeasurementDescription(item model.Measurement) string {
	return fmt.Sprintf("%s %.3f %s (uncertainty %.3f)", model.MeasurementLabel(item.Kind), item.Value, item.Unit, item.Uncertainty)
}
