package app

import (
	"io"
	"log/slog"
	"testing"

	"github.com/cenron/shipdeck/internal/config"
	"github.com/cenron/shipdeck/internal/state"
)

func TestNewWiresDependencies(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{DBPath: "./data/test.sqlite"}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	store := &state.Store{}

	a := NewApp(cfg, log, store)
	if a == nil {
		t.Fatal("expected app instance")
	}
	if a.cfg != cfg {
		t.Fatalf("expected cfg %p, got %p", cfg, a.cfg)
	}
	if a.log != log {
		t.Fatalf("expected log %p, got %p", log, a.log)
	}
	if a.store != store {
		t.Fatalf("expected store %p, got %p", store, a.store)
	}
}
