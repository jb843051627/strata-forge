package regression

import (
	"context"
	"sync"
	"testing"
)

func TestBug13_ConcurrentQueueCreatesOneRun(t *testing.T) {
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
	var wg sync.WaitGroup
	var mu sync.Mutex
	successes := 0
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, queueErr := fixture.Lab.QueueRun(context.Background(), sample.ID); queueErr == nil {
				mu.Lock()
				successes++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if successes != 1 {
		t.Fatalf("expected one successful queue, got %d", successes)
	}
	runs, err := fixture.Lab.ListRuns(context.Background(), sample.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected one run row, got %d", len(runs))
	}
}

func TestBug13_ConcurrentQueueKeepsOneStoredRun(t *testing.T) {
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
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = fixture.Lab.QueueRun(context.Background(), sample.ID)
		}()
	}
	wg.Wait()
	runs, err := fixture.Lab.ListRuns(context.Background(), sample.ID)
	if err != nil || len(runs) != 1 {
		t.Fatalf("expected one stored run, runs=%d err=%v", len(runs), err)
	}
}

func TestBug13_ConcurrentQueueLeavesOneActiveCandidate(t *testing.T) {
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
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = fixture.Lab.QueueRun(context.Background(), sample.ID)
		}()
	}
	wg.Wait()
	run, err := fixture.Store.FindActiveRun(context.Background(), sample.ID)
	if err != nil || run.ID == 0 {
		t.Fatalf("expected one active candidate, run=%#v err=%v", run, err)
	}
}
