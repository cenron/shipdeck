package app

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/cenron/shipdeck/internal/config"
	"github.com/cenron/shipdeck/internal/deploy"
	"github.com/cenron/shipdeck/internal/state"
)

func TestNewWiresDependencies(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{DBPath: "./data/test.sqlite"}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	store := &state.Store{}
	service := deploy.NewService(deploy.NewEngine(&fakeRuntime{}))

	a := NewApp(cfg, log, store, service)
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
	if a.service != service {
		t.Fatalf("expected service %p, got %p", service, a.service)
	}
}

func TestAppRun(t *testing.T) {
	runtime := &fakeRuntime{}
	a := NewApp(&config.Config{}, slog.New(slog.NewTextHandler(io.Discard, nil)), &state.Store{}, deploy.NewService(deploy.NewEngine(runtime)))
	if err := a.Run(context.Background()); err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	wantActions := []string{"start", "stop", "stop", "deploy"}
	if len(runtime.actions) != len(wantActions) {
		t.Fatalf("expected %d runtime actions, got %d: %v", len(wantActions), len(runtime.actions), runtime.actions)
	}
	for i := range wantActions {
		if runtime.actions[i] != wantActions[i] {
			t.Fatalf("runtime action %d = %q, want %q", i, runtime.actions[i], wantActions[i])
		}
	}
	if len(runtime.deployCalls) != 1 {
		t.Fatalf("expected one deploy call, got %d", len(runtime.deployCalls))
	}
	if runtime.deployCalls[0] != "nginx:1.28-alpine" {
		t.Fatalf("deploy revision = %q, want %q", runtime.deployCalls[0], "nginx:1.28-alpine")
	}
}

type fakeRuntime struct {
	actions     []string
	deployCalls []string
}

func (f *fakeRuntime) StartProject(context.Context, deploy.Project) error {
	f.actions = append(f.actions, "start")
	return nil
}

func (f *fakeRuntime) StopProject(context.Context, deploy.Project) error {
	f.actions = append(f.actions, "stop")
	return nil
}

func (f *fakeRuntime) DeployRevision(_ context.Context, _ deploy.Project, revision string) error {
	f.actions = append(f.actions, "deploy")
	f.deployCalls = append(f.deployCalls, revision)
	return nil
}
