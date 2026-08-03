package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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

	prj := deploy.Project{
		ID:             0,
		Name:           "Test",
		Images:         nil,
		WatchTags:      nil,
		CredentialRefs: nil,
		Update:         deploy.UpdateConfig{},
		UpdateState:    deploy.ProjectUpdateState{},
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
	}

	err := app.service.StartProject(ctx, prj)
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Down")
	err = app.service.StopProject(ctx, prj)
	return nil
}
