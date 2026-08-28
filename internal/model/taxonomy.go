package model

import "strings"

var SupportedKinds = []string{"magnetic", "isotope", "grain"}

func IsSupportedKind(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, kind := range SupportedKinds {
		if value == kind {
			return true
		}
	}
	return false
}

func IsTerminalSample(status string) bool {
	return status == SampleArchived || status == SampleRejected
}

func IsTerminalRun(state string) bool {
	return state == RunCompleted || state == RunCancelled || state == RunFailed
}

func IsTerminalMeasurement(status string) bool {
	return status == MeasurementDone || status == MeasurementRejected
}
