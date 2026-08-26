package report

import (
	"fmt"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func MarkdownSummary(summary model.SampleSummary) string {
	if summary.Sample.ID == 0 {
		return "# 未选择样品\n"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", summary.Sample.Code)
	fmt.Fprintf(&b, "- Site: %s\n- Status: %s\n- Depth: %s\n\n", summary.Sample.Site, summary.Sample.Status, model.FormatDepth(summary.Sample.DepthStart, summary.Sample.DepthEnd))
	b.WriteString("## Layers\n\n| Sequence | Depth | Material |\n|---:|---|---|\n")
	for _, layer := range summary.Layers {
		fmt.Fprintf(&b, "| %d | %s | %s |\n", layer.Sequence, model.FormatDepth(layer.TopDepth, layer.BottomDepth), layer.Material)
	}
	b.WriteString("\n## Measurements\n\n| ID | Kind | Value | Status |\n|---:|---|---:|---|\n")
	for _, item := range summary.Measurements {
		fmt.Fprintf(&b, "| %d | %s | %.3f %s | %s |\n", item.ID, model.MeasurementLabel(item.Kind), item.Value, item.Unit, item.Status)
	}
	return b.String()
}
