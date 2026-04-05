package app

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/cenron/shipdeck/internal/config"
	"github.com/cenron/shipdeck/internal/state"
	"github.com/jmoiron/sqlx"
)

func Wire(ctx context.Context, cfg config.Config, log *slog.Logger) error {
	// Create our store
	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0o755); err != nil {
		return err
	}

	db, err := sqlx.Connect("sqlite", cfg.DBPath)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	store := state.NewStore(db)
	if err := store.Migrate(ctx); err != nil {
		return err
	}

	a := NewApp(&cfg, log, store)
	err = a.Run()
	if err != nil {
		return err
	}

	return nil
}
