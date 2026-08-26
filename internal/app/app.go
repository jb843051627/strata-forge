package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jb843051627/strata-forge/internal/api"
	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/report"
	"github.com/jb843051627/strata-forge/internal/service"
	"github.com/jb843051627/strata-forge/internal/store"
)

type App struct {
	config  Config
	store   *store.Store
	service *service.LabService
	server  *http.Server
}

func New(config Config) (*App, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	st, err := store.Open(config.DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("open application store: %w", err)
	}
	svc := service.New(st, clock.System{})
	handler := api.NewHandler(svc, report.NewRenderer())
	return &App{config: config, store: st, service: svc, server: &http.Server{Addr: config.Address, Handler: handler}}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.service.StartWorker(ctx)
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5e9)
		defer cancel()
		_ = a.server.Shutdown(shutdownCtx)
	}()
	err := a.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (a *App) Smoke(ctx context.Context) error {
	_, err := a.service.ListSamples(ctx, "")
	return err
}

func (a *App) Close() error {
	a.service.Close()
	return a.store.Close()
}
