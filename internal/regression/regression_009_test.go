package regression

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/strata-forge/internal/model"
)

func TestBug09_CancelledBatchDoesNotRecordMeasurement(t *testing.T) {
	fixture, err := NewFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()
	sample, err := fixture.Sample(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	layer, err := fixture.Layer(context.Background(), sample.ID, 1)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = fixture.Lab.BatchRecord(ctx, []model.MeasurementInput{{LayerID: layer.ID, Kind: "magnetic", InputValue: 1, Unit: "si", Operator: "op"}}, []float64{1})
	if !errors.Is(err, model.ErrCancelled) {
		t.Fatalf("expected cancelled batch, got %v", err)
	}
	count, err := fixture.Store.CountMeasurements(context.Background(), sample.ID, "")
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("cancelled batch recorded %d measurements", count)
	}
}

func TestBug09_CancelledBatchReturnsCancelledError(t *testing.T) {
	fixture, err := NewFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()
	sample, err := fixture.Sample(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	layer, err := fixture.Layer(context.Background(), sample.ID, 1)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result, err := fixture.Lab.BatchRecord(ctx, []model.MeasurementInput{{LayerID: layer.ID, Kind: "magnetic", InputValue: 2, Unit: "si", Operator: "op"}}, []float64{2})
	if !errors.Is(err, model.ErrCancelled) || !result.Stopped {
		t.Fatalf("expected stopped cancellation result=%#v err=%v", result, err)
	}
}

func TestBug09_LiveBatchRecordsOneMeasurement(t *testing.T) {
	fixture, err := NewFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()
	sample, err := fixture.Sample(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	layer, err := fixture.Layer(context.Background(), sample.ID, 1)
	if err != nil {
		t.Fatal(err)
	}
	result, err := fixture.Lab.BatchRecord(context.Background(), []model.MeasurementInput{{LayerID: layer.ID, Kind: "magnetic", InputValue: 1, Unit: "si", Operator: "op"}}, []float64{1})
	if err != nil || len(result.Completed) != 1 {
		t.Fatalf("expected one completed measurement, result=%#v err=%v", result, err)
	}
}
