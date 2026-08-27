package regression

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/jb843051627/strata-forge/internal/model"
)

func TestBug14_ConcurrentStartTransitionsOnce(t *testing.T) {
	fixture, err := NewFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()
	sample, err := fixture.Sample(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = fixture.Layer(context.Background(), sample.ID, 1); err != nil {
		t.Fatal(err)
	}
	run, err := fixture.Lab.QueueRun(context.Background(), sample.ID)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	successes := 0
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, startErr := fixture.Lab.StartRun(context.Background(), run.ID); startErr == nil {
				mu.Lock()
				successes++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if successes != 1 {
		t.Fatalf("expected one successful transition, got %d", successes)
	}
	current, err := fixture.Lab.GetRun(context.Background(), run.ID)
	if err != nil {
		t.Fatal(err)
	}
	if current.State != model.RunActive {
		t.Fatalf("expected active run, got %s", current.State)
	}
}

func TestBug14_CancelledStartReturnsContextError(t *testing.T) {
	fixture, err := NewFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()
	sample, err := fixture.Sample(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = fixture.Layer(context.Background(), sample.ID, 1); err != nil {
		t.Fatal(err)
	}
	run, err := fixture.Lab.QueueRun(context.Background(), sample.ID)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = fixture.Lab.StartRun(ctx, run.ID)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
}

func TestBug14_CancelledStartKeepsRunQueued(t *testing.T) {
	fixture, err := NewFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()
	sample, err := fixture.Sample(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = fixture.Layer(context.Background(), sample.ID, 1); err != nil {
		t.Fatal(err)
	}
	run, err := fixture.Lab.QueueRun(context.Background(), sample.ID)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _ = fixture.Lab.StartRun(ctx, run.ID)
	stored, err := fixture.Lab.GetRun(context.Background(), run.ID)
	if err != nil || stored.State != model.RunQueued {
		t.Fatalf("expected queued run, got %#v err=%v", stored, err)
	}
}
