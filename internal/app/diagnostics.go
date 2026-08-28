package app

import (
	"context"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/store"
)

type Diagnostics struct {
	Database string
	Stats    store.Stats
}

func (a *App) Diagnostics(ctx context.Context) (Diagnostics, error) {
	stats, err := a.store.Stats(ctx)
	if err != nil {
		return Diagnostics{}, fmt.Errorf("diagnostics: %w", err)
	}
	return Diagnostics{Database: a.config.DatabasePath, Stats: stats}, nil
}
