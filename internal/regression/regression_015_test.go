package regression

import (
	"testing"

	"github.com/jb843051627/strata-forge/internal/service"
)

func TestBug15_SeriesSnapshotDoesNotAliasState(t *testing.T) {
	series := service.NewMeasurementSeries()
	if err := series.Append(15, 1.2); err != nil {
		t.Fatal(err)
	}
	first := series.Snapshot(15)
	first[0] = 99
	second := series.Snapshot(15)
	if second[0] != 1.2 {
		t.Fatalf("series state was polluted: %v", second)
	}
}
