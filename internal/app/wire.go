package app

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/cenron/shipdeck/internal/adapters/docker"
	"github.com/cenron/shipdeck/internal/config"
	"github.com/cenron/shipdeck/internal/deploy"
	"github.com/cenron/shipdeck/internal/session"
	"github.com/cenron/shipdeck/internal/state"
	"github.com/jmoiron/sqlx"
)

func Wire(ctx context.Context, cfg config.Config, log *slog.Logger) error {
	// Create our store
	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0755); err != nil {
		return err
	}

	db, err := sqlx.Connect("sqlite", cfg.DBPath)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	store := state.NewStore(db, log)
	if err := store.Migrate(ctx); err != nil {
		return err
	}

	// Create our temporary config resolver
	resolver := docker.SmokeComposeResolver{}

	r := docker.NewClient(resolver)
	e := deploy.NewEngine(r)
	service := deploy.NewService(e)

	a := NewApp(&cfg, log, store, service)
	err = a.Run(ctx)
	if err != nil {
		return err
	}

	server := session.NewServer(ctx, log, &cfg)
	err = server.Run()
	if err != nil {
		return err
	}

	return nil
}
