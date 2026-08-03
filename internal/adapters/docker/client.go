package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cenron/shipdeck/internal/deploy"
)

// Client executes Docker operations for deploy runtime actions.
type Client struct {
	resolver ComposeSpecResolver
}

// NewClient constructs a Docker runtime adapter.
func NewClient(r ComposeSpecResolver) *Client {
	if r == nil {
		panic("compose spec resolver must not be nil")
	}

	return &Client{
		resolver: r,
	}
}

// StartProject starts the project's resolved Compose workload.
func (c *Client) StartProject(ctx context.Context, project deploy.Project) error {

	spec, err := c.resolver.Resolve(project, "")
	if err != nil {
		return err
	}

	out, err := executeDocker(ctx, spec, "up")
	if err != nil {
		return fmt.Errorf("docker compose up: %w: %s", err, string(out))
	}

	return nil
}

// StopProject stops the project's resolved Compose workload.
func (c *Client) StopProject(ctx context.Context, project deploy.Project) error {

	spec, err := c.resolver.Resolve(project, "")
	if err != nil {
		return err
	}

	out, err := executeDocker(ctx, spec, "down")
	if err != nil {
		return fmt.Errorf("docker compose down: %w: %s", err, string(out))
	}

	return nil
}

// DeployRevision deploys a requested project revision.
func (c *Client) DeployRevision(ctx context.Context, project deploy.Project, revision string) error {
	revision = strings.TrimSpace(revision)
	if revision == "" {
		return fmt.Errorf("revision must not be empty")
	}

	spec, err := c.resolver.Resolve(project, revision)
	if err != nil {
		return err
	}

	out, err := executeDocker(ctx, spec, "up")
	if err != nil {
		return fmt.Errorf("docker compose up: %w: %s", err, string(out))
	}

	return nil
}

// composeEnv merges Compose-specific environment values into the current process environment.
func composeEnv(values map[string]string) []string {
	env := os.Environ()
	for k, v := range values {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env
}

// executeDocker runs a Docker Compose action and returns the combined stdout/stderr output.
func executeDocker(ctx context.Context, spec ComposeSpec, action string) ([]byte, error) {

	args := []string{
		"compose",
		"-f",
		spec.File,
		"-p",
		spec.ProjectName,
		action,
	}

	if action == "up" {
		args = append(args, "-d")
	}

	cmd := exec.CommandContext(
		ctx,
		"docker",
		args...,
	)

	cmd.Env = composeEnv(spec.Env)
	return cmd.CombinedOutput()
}
