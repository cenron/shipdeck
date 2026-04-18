package deploy

import (
	"context"
	"fmt"
)

type Runtime interface {
	// StartProject starts the currently configured workload for a project.
	StartProject(ctx context.Context, project Project) error
	// StopProject stops the currently running workload for a project.
	StopProject(ctx context.Context, project Project) error
	// DeployRevision deploys a specific revision for a project.
	DeployRevision(ctx context.Context, project Project, revision string) error
}

type Engine struct {
	runtime Runtime
}

// NewEngine constructs an Engine with the runtime adapter used for deploy actions.
func NewEngine(runtime Runtime) *Engine {
	if runtime == nil {
		panic("runtime must not be nil")
	}

	return &Engine{
		runtime: runtime,
	}
}

// Start delegates project start to the runtime adapter.
func (e *Engine) Start(ctx context.Context, project Project) error {
	return e.runtime.StartProject(ctx, project)
}

// Stop delegates project stop to the runtime adapter.
func (e *Engine) Stop(ctx context.Context, project Project) error {
	return e.runtime.StopProject(ctx, project)
}

// Rollback deploys the project's latest known digest revision.
func (e *Engine) Rollback(ctx context.Context, req RollbackRequest) error {
	return e.runtime.DeployRevision(ctx, req.Project, req.Revision)
}

// Redeploy stops the current workload, deploys the target revision, and rolls back on failure.
func (e *Engine) Redeploy(ctx context.Context, req RedeployRequest) error {

	if req.PreviousRevision == "" {
		return fmt.Errorf("no previous revision specified")
	}

	if req.TargetRevision == "" {
		return fmt.Errorf("no target revision specified")
	}

	err := e.runtime.StopProject(ctx, req.Project)
	if err != nil {
		return fmt.Errorf("stop current: %w", err)
	}

	err = e.runtime.DeployRevision(ctx, req.Project, req.TargetRevision)
	if err != nil {
		rbErr := e.runtime.DeployRevision(ctx, req.Project, req.PreviousRevision)
		if rbErr != nil {
			return fmt.Errorf("deploy failed: %v; rollback failed: %w", err, rbErr)
		}

		return fmt.Errorf("deploy failed, rollback succeeded: %w", err)
	}
	return nil
}
