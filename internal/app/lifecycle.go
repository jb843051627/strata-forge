package app

import (
	"context"
	"time"
)

func RunFor(ctx context.Context, app *App, duration time.Duration) error {
	child, cancel := context.WithTimeout(ctx, duration)
	defer cancel()
	return app.Run(child)
}

func (a *App) Config() Config {
	return a.config
}
