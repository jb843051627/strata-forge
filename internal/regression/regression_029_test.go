package regression

import (
	"context"
	"testing"

	"github.com/jb843051627/strata-forge/internal/model"
)

func TestBug29_FailedReviewRejectsSample(t *testing.T) {
	fixture, err := NewFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()
	sample, measurement, err := SeedReviewedSample(context.Background(), fixture)
	if err != nil {
		t.Fatal(err)
	}
	_, err = fixture.Lab.ReviewMeasurement(context.Background(), measurement.ID, model.ReviewInput{Decision: model.ReviewFail, Reviewer: "reviewer-2", Comment: "outlier"})
	if err != nil {
		t.Fatal(err)
	}
	current, err := fixture.Lab.GetSample(context.Background(), sample.ID)
	if err != nil {
		t.Fatal(err)
	}
	if current.Status != model.SampleRejected {
		t.Fatalf("expected rejected sample, got %s", current.Status)
	}
}
