package service

import (
	"bytes"
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
	"github.com/jb843051627/strata-forge/internal/report"
)

func (s *LabService) ExportSample(ctx context.Context, sampleID int64, format string) ([]byte, string, error) {
	summary, err := s.Summary(ctx, sampleID)
	if err != nil {
		return nil, "", err
	}
	switch format {
	case "json", "":
		data, encodeErr := report.NewRenderer().JSON(summary)
		return data, "application/json", encodeErr
	case "text":
		return []byte(report.NewRenderer().Text(summary)), "text/plain; charset=utf-8", nil
	case "csv":
		measurements, listErr := s.store.ListMeasurementsForSample(ctx, sampleID)
		if listErr != nil {
			return nil, "", wrap("read export measurements", listErr)
		}
		var buffer bytes.Buffer
		if writeErr := report.WriteMeasurements(&buffer, measurements); writeErr != nil {
			return nil, "", wrap("write csv export", writeErr)
		}
		return buffer.Bytes(), "text/csv; charset=utf-8", nil
	default:
		return nil, "", fmt.Errorf("%w: unsupported export format %s", model.ErrInvalidInput, format)
	}
}
