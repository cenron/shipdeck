package app

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cenron/shipdeck/internal/config"
)

func TestWireWithCancelledContext(t *testing.T) {
	tmp := t.TempDir()
	authPath := filepath.Join(tmp, "authorized_keys")

	const keyLine = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICgVP9MWumPiw3gETqb5Z+CiQ0LtRddpmNQDmUWbxbJc test@example.com\n"
	if err := os.WriteFile(authPath, []byte(keyLine), 0o600); err != nil {
		t.Fatalf("write authorized keys: %v", err)
	}

	cfg := config.Config{
		Host:               "127.0.0.1",
		Port:               "0",
		DBPath:             filepath.Join(tmp, "shipdeck.sqlite"),
		AuthorizedKeysPath: authPath,
	}

	binDir := filepath.Join(tmp, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("make fake bin dir: %v", err)
	}
	fakeDocker := filepath.Join(binDir, "docker")
	if err := os.WriteFile(fakeDocker, []byte("#!/bin/sh\nprintf 'ok\\n'\n"), 0o755); err != nil {
		t.Fatalf("write fake docker: %v", err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(2500 * time.Millisecond)
		cancel()
	}()

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	if err := Wire(ctx, cfg, log); err != nil {
		t.Fatalf("Wire() returned error: %v", err)
	}
}
