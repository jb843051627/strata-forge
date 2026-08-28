package regression

import (
	"testing"

	"github.com/jb843051627/strata-forge/internal/engine"
)

func TestBug20_NilProvenanceDoesNotPanic(t *testing.T) {
	if got := engine.FormatProvenance(nil); got != "unknown provenance" {
		t.Fatalf("unexpected label %q", got)
	}
}
