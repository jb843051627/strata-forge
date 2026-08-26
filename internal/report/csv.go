package report

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/jb843051627/strata-forge/internal/model"
)

func WriteMeasurements(w io.Writer, items []model.Measurement) error {
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"id", "layer_id", "kind", "value", "unit", "uncertainty", "status"}); err != nil {
		return err
	}
	for _, item := range items {
		if err := writer.Write([]string{fmt.Sprint(item.ID), fmt.Sprint(item.LayerID), item.Kind, fmt.Sprintf("%.6f", item.Value), item.Unit, fmt.Sprintf("%.6f", item.Uncertainty), item.Status}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
