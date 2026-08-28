package regression

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestBug12_CancelledTimelineQueryStops(t *testing.T) {
	fixture, err := NewFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()
	sample, err := fixture.Sample(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = fixture.Store.EventsSince(ctx, sample.ID, fixture.Clock.Now().Add(-time.Hour))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
}
