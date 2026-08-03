package docker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cenron/shipdeck/internal/deploy"
)

// composeProject holds the temporary Compose smoke-test inputs used by the Docker adapter.
// These static values will move to the project/config boundary after the adapter path is proven.
type composeProject struct {
	file        string
	projectName string
	env         map[string]string
}

// Client executes Docker operations for deploy runtime actions.
type Client struct{}

// NewClient constructs a Docker runtime adapter.
func NewClient() *Client {
	return &Client{}
}

// StartProject starts the project's Compose workload using the current smoke-test fixture.
func (c *Client) StartProject(ctx context.Context, project deploy.Project) error {

	name := strings.ToLower(project.Name)
	compose := composeProject{
		file:        "examples/compose-test.yml",
		projectName: name,
		env: map[string]string{
			"SHIPDECK_TEST_IMAGE":     "nginx:1.27-alpine",
			"SHIPDECK_TEST_CONTAINER": name + "-web",
			"SHIPDECK_TEST_BIND":      "127.0.0.1",
			"SHIPDECK_TEST_PORT":      "8088",
		},
	}

	out, err := executeDocker(ctx, compose, "up")
	if err != nil {
		return fmt.Errorf("docker compose up: %w: %s", err, string(out))
	}

	fmt.Println(string(out))
	return nil

}

// StopProject stops the project's Compose workload using the current smoke-test fixture.
func (c *Client) StopProject(ctx context.Context, project deploy.Project) error {
	fmt.Println("Stopping project", project.Name)
	name := strings.ToLower(project.Name)
	compose := composeProject{
		file:        "examples/compose-test.yml",
		projectName: name,
		env: map[string]string{
			"SHIPDECK_TEST_IMAGE":     "nginx:1.27-alpine",
			"SHIPDECK_TEST_CONTAINER": name + "-web",
			"SHIPDECK_TEST_BIND":      "127.0.0.1",
			"SHIPDECK_TEST_PORT":      "8088",
		},
	}

	out, err := executeDocker(ctx, compose, "down")
	if err != nil {
		return fmt.Errorf("docker compose down: %w: %s", err, string(out))
	}

	fmt.Println(string(out))
	return nil
}

// DeployRevision deploys a requested project revision.
// It is the remaining adapter method to implement for redeploy and rollback paths.
func (c *Client) DeployRevision(_ context.Context, project deploy.Project, revision string) error {
	return errors.New("not implemented")

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
func executeDocker(ctx context.Context, compose composeProject, action string) ([]byte, error) {

	args := []string{
		"compose",
		"-f",
		compose.file,
		"-p",
		compose.projectName,
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

	cmd.Env = composeEnv(compose.env)
	return cmd.CombinedOutput()
}
