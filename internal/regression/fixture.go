package regression

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/model"
	"github.com/jb843051627/strata-forge/internal/service"
	"github.com/jb843051627/strata-forge/internal/store"
)

type Fixture struct {
	Dir   string
	Store *store.Store
	Lab   *service.LabService
	Clock *clock.Fake
}

func NewFixture() (*Fixture, error) {
	dir, err := os.MkdirTemp("", "strata-forge-regression-")
	if err != nil {
		return nil, err
	}
	st, err := store.Open(filepath.Join(dir, "fixture.db"))
	if err != nil {
		_ = os.RemoveAll(dir)
		return nil, err
	}
	fake := clock.NewFake(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
	return &Fixture{Dir: dir, Store: st, Lab: service.New(st, fake), Clock: fake}, nil
}

func (f *Fixture) Close() {
	f.Lab.Close()
	_ = f.Store.Close()
	_ = os.RemoveAll(f.Dir)
}

func (f *Fixture) Sample(ctx context.Context) (model.Sample, error) {
	return f.Lab.ReceiveSample(ctx, model.SampleInput{Code: "SF-CORE-01", Site: "North Basin", DepthStart: 0, DepthEnd: 12, Notes: "sealed tube"})
}

func (f *Fixture) Layer(ctx context.Context, sampleID int64, seq int) (model.Layer, error) {
	top := float64(seq - 1)
	return f.Lab.AddLayer(ctx, sampleID, model.LayerInput{Sequence: seq, TopDepth: top, BottomDepth: top + 1, Material: "laminated silt"})
}

func (f *Fixture) Measurement(ctx context.Context, layerID int64, kind string, value float64) (model.Measurement, error) {
	return f.Lab.RecordMeasurement(ctx, model.MeasurementInput{LayerID: layerID, Kind: kind, InputValue: value, Unit: "si", Operator: "lab-tech"}, value, 0.05)
}
