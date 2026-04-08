package main

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestMainSuccessDoesNotExit(t *testing.T) {
	oldRun := runFn
	oldExit := exitFn
	t.Cleanup(func() {
		runFn = oldRun
		exitFn = oldExit
	})

	runFn = func() error { return nil }
	exited := false
	exitFn = func(_ int) { exited = true }

	main()
	if exited {
		t.Fatal("main should not exit on success")
	}
}

func TestMainFailureExits(t *testing.T) {
	oldRun := runFn
	oldExit := exitFn
	t.Cleanup(func() {
		runFn = oldRun
		exitFn = oldExit
	})

	runFn = func() error { return errors.New("boom") }
	code := 0
	exitFn = func(c int) { code = c }

	main()
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestRunShutsDownOnSignal(t *testing.T) {
	tmp := t.TempDir()
	authPath := filepath.Join(tmp, "authorized_keys")

	const keyLine = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICgVP9MWumPiw3gETqb5Z+CiQ0LtRddpmNQDmUWbxbJc test@example.com\n"
	if err := os.WriteFile(authPath, []byte(keyLine), 0o600); err != nil {
		t.Fatalf("write authorized keys: %v", err)
	}

	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("PORT", "0")
	t.Setenv("DB_PATH", filepath.Join(tmp, "shipdeck.sqlite"))
	t.Setenv("AUTHORIZED_KEYS_PATH", authPath)

	go func() {
		time.Sleep(150 * time.Millisecond)
		_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}()

	if err := run(); err != nil {
		t.Fatalf("run() returned error: %v", err)
	}
}
