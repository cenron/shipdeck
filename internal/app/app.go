package app

import (
	"log/slog"

	"github.com/cenron/shipdeck/internal/config"
	"github.com/cenron/shipdeck/internal/state"
)

type App struct {
	cfg   *config.Config
	log   *slog.Logger
	store *state.Store
}

func NewApp(cfg *config.Config, log *slog.Logger, store *state.Store) *App {
	return &App{
		cfg:   cfg,
		log:   log,
		store: store,
	}
}

func (app *App) Run() error {
	return nil
}
