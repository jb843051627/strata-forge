package regression

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/strata-forge/internal/model"
)

func TestBug06_AcknowledgingMissingAlertKeepsErrorIdentity(t *testing.T) {
	fixture, err := NewFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()
	err = fixture.Lab.AcknowledgeAlert(context.Background(), 1002)
	if !errors.Is(err, model.ErrNotFound) {
		t.Fatalf("expected not-found identity, got %v", err)
	}
}
