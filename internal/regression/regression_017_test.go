package regression

import (
	"sync"
	"testing"
	"time"

	"github.com/jb843051627/strata-forge/internal/engine"
)

func TestBug17_AuditAppendIsSynchronized(t *testing.T) {
	audit := engine.NewRunAudit()
	var wg sync.WaitGroup
	for i := 0; i < 24; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				audit.Add(int64(id*20+j), "state", time.Now())
				_ = audit.Snapshot()
			}
		}(i)
	}
	wg.Wait()
	if len(audit.Snapshot()) != 480 {
		t.Fatalf("audit entries lost")
	}
}
