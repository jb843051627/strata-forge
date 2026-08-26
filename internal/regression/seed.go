package regression

import (
	"context"

	"github.com/jb843051627/strata-forge/internal/model"
)

func SeedReviewedSample(ctx context.Context, f *Fixture) (model.Sample, model.Measurement, error) {
	sample, err := f.Sample(ctx)
	if err != nil {
		return model.Sample{}, model.Measurement{}, err
	}
	layer, err := f.Layer(ctx, sample.ID, 1)
	if err != nil {
		return model.Sample{}, model.Measurement{}, err
	}
	measurement, err := f.Measurement(ctx, layer.ID, "magnetic", 1.4)
	if err != nil {
		return model.Sample{}, model.Measurement{}, err
	}
	if _, err := f.Lab.ReviewMeasurement(ctx, measurement.ID, model.ReviewInput{Decision: model.ReviewPass, Reviewer: "reviewer-1", Comment: "within expected range"}); err != nil {
		return model.Sample{}, model.Measurement{}, err
	}
	return sample, measurement, nil
}
