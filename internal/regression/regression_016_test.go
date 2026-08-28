package regression

import (
	"testing"

	"github.com/jb843051627/strata-forge/internal/engine"
)

func TestBug16_StateCacheSnapshotDoesNotAliasState(t *testing.T) {
	cache := engine.NewStateCache()
	cache.Set(16, "review")
	snapshot := cache.Snapshot()
	snapshot[16] = "archived"
	value, ok := cache.Get(16)
	if !ok || value != "review" {
		t.Fatalf("cache state was polluted: %q", value)
	}
}
