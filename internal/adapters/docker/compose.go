package docker

import (
	"strings"

	"github.com/cenron/shipdeck/internal/deploy"
)

type ComposeSpec struct {
	File        string
	ProjectName string
	Env         map[string]string
}

type ComposeSpecResolver interface {
	Resolve(project deploy.Project, revision string) (ComposeSpec, error)
}

type SmokeComposeResolver struct{}

func (sc SmokeComposeResolver) Resolve(project deploy.Project, revision string) (ComposeSpec, error) {
	name := strings.ToLower(project.Name)
	image := "nginx:1.27-alpine"
	if strings.TrimSpace(revision) != "" {
		image = strings.TrimSpace(revision)
	}

	return ComposeSpec{
		File:        "examples/compose-test.yml",
		ProjectName: name,
		Env: map[string]string{
			"SHIPDECK_TEST_IMAGE":     image,
			"SHIPDECK_TEST_CONTAINER": name + "-web",
			"SHIPDECK_TEST_BIND":      "127.0.0.1",
			"SHIPDECK_TEST_PORT":      "8088",
		},
	}, nil

}
