package regression

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/strata-forge/internal/model"
)

func TestBug30_UpdatingMissingRunReturnsError(t *testing.T) {
	fixture, err := NewFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()
	err = fixture.Store.UpdateRun(context.Background(), model.Run{ID: 99999, State: model.RunActive})
	if !errors.Is(err, model.ErrNotFound) {
		t.Fatalf("expected not-found update error, got %v", err)
	}
}
