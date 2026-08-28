package regression

import (
	"testing"

	"github.com/jb843051627/strata-forge/internal/engine"
	"github.com/jb843051627/strata-forge/internal/model"
)

func TestBug26_MergeMeasurementsDoesNotMutatePrimary(t *testing.T) {
	primary := []model.Measurement{{ID: 1, Value: 1}}
	merged := engine.MergeMeasurements(primary, []model.Measurement{{ID: 1, Value: 2}})
	merged[0].Value = 99
	if primary[0].Value == 99 {
		t.Fatal("primary measurement was aliased")
	}
}
