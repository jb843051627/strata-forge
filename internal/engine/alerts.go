package engine

import (
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

type AlertRule struct {
	Code     string
	Severity string
	Message  string
}

func EvaluateMeasurementAlert(measurement model.Measurement) *AlertRule {
	if measurement.Uncertainty > 0.5 {
		return &AlertRule{Code: "uncertainty-critical", Severity: "critical", Message: "measurement uncertainty is outside review limits"}
	}
	if measurement.Uncertainty > 0.25 {
		return &AlertRule{Code: "uncertainty-high", Severity: "warning", Message: "measurement requires a second reading"}
	}
	if measurement.Value < 0 {
		return &AlertRule{Code: "negative-value", Severity: "critical", Message: "measurement value cannot describe a valid layer"}
	}
	return nil
}

func BuildAlert(sampleID, layerID int64, rule *AlertRule, now string) model.Alert {
	if rule == nil {
		return model.Alert{SampleID: sampleID, LayerID: layerID, Severity: "warning", Code: "missing-rule", Message: "alert rule unavailable"}
	}
	return model.Alert{SampleID: sampleID, LayerID: layerID, Severity: rule.Severity, Code: rule.Code, Message: fmt.Sprintf("%s (%s)", rule.Message, now)}
}
