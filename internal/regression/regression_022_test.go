package regression

import (
	"testing"

	"github.com/jb843051627/strata-forge/internal/engine"
	"github.com/jb843051627/strata-forge/internal/model"
)

func TestBug22_EmptyChronologyIsSafe(t *testing.T) {
	if got := engine.BuildChronology(nil, model.AgeEstimate{}); got != nil {
		t.Fatalf("expected empty chronology, got %#v", got)
	}
}
