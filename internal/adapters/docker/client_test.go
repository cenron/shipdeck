package docker

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cenron/shipdeck/internal/deploy"
)

func TestClientStartAndStopProjectRunComposeCommands(t *testing.T) {
	tmp := t.TempDir()
	binDir := filepath.Join(tmp, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("make fake bin dir: %v", err)
	}

	logPath := filepath.Join(tmp, "docker.log")
	envPath := filepath.Join(tmp, "docker.env")
	dockerPath := filepath.Join(binDir, "docker")
	script := `#!/bin/sh
printf '%s\n' "$*" >> "$SHIPDECK_DOCKER_LOG"
printf '%s|%s|%s|%s\n' "$SHIPDECK_TEST_IMAGE" "$SHIPDECK_TEST_CONTAINER" "$SHIPDECK_TEST_BIND" "$SHIPDECK_TEST_PORT" >> "$SHIPDECK_DOCKER_ENV_LOG"
printf 'ok\n'
`
	if err := os.WriteFile(dockerPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake docker: %v", err)
	}

	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("SHIPDECK_DOCKER_LOG", logPath)
	t.Setenv("SHIPDECK_DOCKER_ENV_LOG", envPath)

	client := NewClient(SmokeComposeResolver{})
	project := deploy.Project{Name: "Test"}

	if err := client.StartProject(context.Background(), project); err != nil {
		t.Fatalf("StartProject() returned error: %v", err)
	}
	if err := client.StopProject(context.Background(), project); err != nil {
		t.Fatalf("StopProject() returned error: %v", err)
	}
	if err := client.DeployRevision(context.Background(), project, "  nginx:1.28-alpine  "); err != nil {
		t.Fatalf("DeployRevision() returned error: %v", err)
	}

	logBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read docker log: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(logBytes)), "\n")
	want := []string{
		"compose -f examples/compose-test.yml -p test up -d",
		"compose -f examples/compose-test.yml -p test down",
		"compose -f examples/compose-test.yml -p test up -d",
	}
	if len(lines) != len(want) {
		t.Fatalf("expected %d docker calls, got %d: %q", len(want), len(lines), lines)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("docker call %d = %q, want %q", i, lines[i], want[i])
		}
	}

	envBytes, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("read docker env log: %v", err)
	}
	envLines := strings.Split(strings.TrimSpace(string(envBytes)), "\n")
	wantEnv := []string{
		"nginx:1.27-alpine|test-web|127.0.0.1|8088",
		"nginx:1.27-alpine|test-web|127.0.0.1|8088",
		"nginx:1.28-alpine|test-web|127.0.0.1|8088",
	}
	if len(envLines) != len(wantEnv) {
		t.Fatalf("expected %d docker env calls, got %d: %q", len(wantEnv), len(envLines), envLines)
	}
	for i := range wantEnv {
		if envLines[i] != wantEnv[i] {
			t.Fatalf("docker env call %d = %q, want %q", i, envLines[i], wantEnv[i])
		}
	}
}

func TestNewClientPanicsWhenResolverNil(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected NewClient(nil) to panic")
		}
		if msg, ok := r.(string); !ok || msg != "compose spec resolver must not be nil" {
			t.Fatalf("expected panic message %q, got %#v", "compose spec resolver must not be nil", r)
		}
	}()

	_ = NewClient(nil)
}
