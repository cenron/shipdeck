package session

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/cenron/shipdeck/internal/config"
)

func TestNewServer(t *testing.T) {
	ctx := context.Background()
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{Host: "127.0.0.1", Port: "2222", AuthorizedKeysPath: "./data/authorized_keys"}

	s := NewServer(ctx, log, cfg)
	if s == nil {
		t.Fatal("expected server")
	}
	if s.ctx != ctx {
		t.Fatal("expected context to be stored")
	}
	if s.log != log {
		t.Fatal("expected logger to be stored")
	}
	if s.config.Host != cfg.Host || s.config.Port != cfg.Port || s.config.AuthorizedKeysPath != cfg.AuthorizedKeysPath {
		t.Fatalf("unexpected config stored: %+v", s.config)
	}
}

func TestRunErrorsWhenAuthorizedKeysMissing(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	s := NewServer(ctx, log, &config.Config{
		Host:               "127.0.0.1",
		Port:               "0",
		AuthorizedKeysPath: "./does-not-exist/authorized_keys",
	})

	err := s.Run()
	if err == nil {
		t.Fatal("expected run error")
	}
	if !strings.Contains(err.Error(), "create SSH server") {
		t.Fatalf("expected create SSH server error, got %v", err)
	}
}

func TestRunStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	s := NewServer(ctx, log, &config.Config{
		Host:               "127.0.0.1",
		Port:               "0",
		AuthorizedKeysPath: "../../data/authorized_keys",
	})

	err := s.Run()
	if err != nil {
		t.Fatalf("expected clean shutdown on canceled context, got %v", err)
	}
}
