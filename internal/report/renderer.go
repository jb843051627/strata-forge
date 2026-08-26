package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

type Renderer struct{}

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (Renderer) JSON(summary model.SampleSummary) ([]byte, error) {
	return json.MarshalIndent(summary, "", "  ")
}

func (Renderer) Text(summary model.SampleSummary) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Sample %s at %s\n", summary.Sample.Code, summary.Sample.Site)
	fmt.Fprintf(&b, "Status: %s\n", summary.Sample.Status)
	fmt.Fprintf(&b, "Layers: %d\nMeasurements: %d\n", len(summary.Layers), len(summary.Measurements))
	for _, item := range summary.Measurements {
		fmt.Fprintf(&b, "- %d %s %.3f %s (%s)\n", item.ID, model.MeasurementLabel(item.Kind), item.Value, item.Unit, item.Status)
	}
	return b.String()
}
