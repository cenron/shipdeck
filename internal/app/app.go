package app

import (
	"context"
	"log/slog"

	"github.com/cenron/shipdeck/internal/config"
	"github.com/cenron/shipdeck/internal/deploy"
	"github.com/cenron/shipdeck/internal/state"
)

type App struct {
	cfg     *config.Config
	log     *slog.Logger
	store   *state.Store
	service *deploy.Service
}

func NewApp(cfg *config.Config, log *slog.Logger, store *state.Store, service *deploy.Service) *App {
	return &App{
		cfg:     cfg,
		log:     log,
		store:   store,
		service: service,
	}
}

func (app *App) Run(ctx context.Context) error {

	return nil
}
