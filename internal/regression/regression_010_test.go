package regression

import (
	"context"
	"errors"
	"testing"
)

func TestBug10_CancelledReportBuildDoesNotPersist(t *testing.T) {
	fixture, err := NewFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()
	sample, _, err := SeedReviewedSample(context.Background(), fixture)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = fixture.Lab.BuildReport(ctx, sample.ID)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
	reports, err := fixture.Store.ListReports(context.Background(), sample.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 0 {
		t.Fatalf("cancelled report build persisted %d reports", len(reports))
	}
}

func TestBug10_CancelledReportReturnsContextError(t *testing.T) {
	fixture, err := NewFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()
	sample, _, err := SeedReviewedSample(context.Background(), fixture)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = fixture.Lab.BuildReport(ctx, sample.ID)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
}

func TestBug10_LiveReportBuildCreatesOneReport(t *testing.T) {
	fixture, err := NewFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()
	sample, _, err := SeedReviewedSample(context.Background(), fixture)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = fixture.Lab.BuildReport(context.Background(), sample.ID); err != nil {
		t.Fatal(err)
	}
	reports, err := fixture.Store.ListReports(context.Background(), sample.ID)
	if err != nil || len(reports) != 1 {
		t.Fatalf("expected one report, reports=%d err=%v", len(reports), err)
	}
}
