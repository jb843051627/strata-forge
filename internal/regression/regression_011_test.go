package regression

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/strata-forge/internal/model"
)

func TestBug11_CancelledReviewDoesNotPersist(t *testing.T) {
	fixture, err := NewFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()
	sample, measurement, err := SeedReviewedSample(context.Background(), fixture)
	if err != nil {
		t.Fatal(err)
	}
	_ = sample
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = fixture.Lab.ReviewMeasurement(ctx, measurement.ID, model.ReviewInput{Decision: model.ReviewPass, Reviewer: "cancelled"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
	reviews, err := fixture.Store.ListReviewsForSample(context.Background(), measurement.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 1 {
		t.Fatalf("cancelled review changed review count to %d", len(reviews))
	}
}
