package report

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func TimelineText(entries []model.TimelineEntry) string {
	ordered := append([]model.TimelineEntry(nil), entries...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].At.Before(ordered[j].At) })
	var b strings.Builder
	for _, entry := range ordered {
		fmt.Fprintf(&b, "%s | %s | %d | %s\n", entry.At.Format("2006-01-02T15:04:05Z07:00"), entry.Kind, entry.Reference, entry.State)
	}
	return b.String()
}

func AlertText(alerts []model.Alert) string {
	var b strings.Builder
	for _, alert := range alerts {
		state := "open"
		if alert.Acknowledged {
			state = "acknowledged"
		}
		fmt.Fprintf(&b, "%s %s %s %s\n", state, alert.Severity, alert.Code, alert.Message)
	}
	return b.String()
}
