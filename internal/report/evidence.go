package report

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

type Evidence struct {
	SampleCode string
	Checks     []string
	Warnings   []string
}

func BuildEvidence(summary model.SampleSummary) Evidence {
	evidence := Evidence{SampleCode: summary.Sample.Code, Checks: make([]string, 0), Warnings: make([]string, 0)}
	if len(summary.Layers) > 0 {
		evidence.Checks = append(evidence.Checks, "layer sequence recorded")
	} else {
		evidence.Warnings = append(evidence.Warnings, "no layers")
	}
	for _, item := range summary.Measurements {
		if item.Status == model.MeasurementDone {
			evidence.Checks = append(evidence.Checks, fmt.Sprintf("measurement %d complete", item.ID))
		} else {
			evidence.Warnings = append(evidence.Warnings, fmt.Sprintf("measurement %d is %s", item.ID, item.Status))
		}
	}
	sort.Strings(evidence.Checks)
	sort.Strings(evidence.Warnings)
	return evidence
}

func (e Evidence) Text() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Sample: %s\nChecks:\n", e.SampleCode)
	for _, check := range e.Checks {
		fmt.Fprintf(&b, "- %s\n", check)
	}
	if len(e.Warnings) > 0 {
		b.WriteString("Warnings:\n")
		for _, warning := range e.Warnings {
			fmt.Fprintf(&b, "- %s\n", warning)
		}
	}
	return b.String()
}
