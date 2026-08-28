package model

func StatusGroup(status string) string {
	switch status {
	case SampleReceived, SampleLayered:
		return "preparation"
	case SampleRunning:
		return "measurement"
	case SampleReview:
		return "quality"
	case SampleArchived:
		return "archive"
	case SampleRejected:
		return "exception"
	default:
		return "unknown"
	}
}

func StatusTerminal(status string) bool {
	return IsTerminalSample(status) || IsTerminalRun(status) || IsTerminalMeasurement(status)
}

func ValidSeverity(value string) bool {
	return value == "critical" || value == "warning" || value == "info"
}
