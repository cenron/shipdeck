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

	client := NewClient()
	project := deploy.Project{Name: "Test"}

	if err := client.StartProject(context.Background(), project); err != nil {
		t.Fatalf("StartProject() returned error: %v", err)
	}
	if err := client.StopProject(context.Background(), project); err != nil {
		t.Fatalf("StopProject() returned error: %v", err)
	}

	logBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read docker log: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(logBytes)), "\n")
	want := []string{
		"compose -f examples/compose-test.yml -p test up -d",
		"compose -f examples/compose-test.yml -p test down",
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
	for i, line := range envLines {
		if line != "nginx:1.27-alpine|test-web|127.0.0.1|8088" {
			t.Fatalf("docker env call %d = %q", i, line)
		}
	}
}
