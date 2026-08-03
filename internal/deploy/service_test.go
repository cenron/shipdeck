package deploy

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestServiceStartProjectDelegatesToEngine(t *testing.T) {
	t.Parallel()

	project := Project{ID: 42, Name: "demo"}
	rt := &serviceFakeRuntime{}
	svc := NewService(NewEngine(rt))

	if err := svc.StartProject(context.Background(), project); err != nil {
		t.Fatalf("StartProject() returned error: %v", err)
	}

	if len(rt.startCalls) != 1 {
		t.Fatalf("expected one StartProject call, got %d", len(rt.startCalls))
	}
	if !reflect.DeepEqual(rt.startCalls[0], project) {
		t.Fatalf("StartProject called with wrong project: got %#v want %#v", rt.startCalls[0], project)
	}
}

func TestServiceStopProjectDelegatesToEngine(t *testing.T) {
	t.Parallel()

	project := Project{ID: 42, Name: "demo"}
	rt := &serviceFakeRuntime{}
	svc := NewService(NewEngine(rt))

	if err := svc.StopProject(context.Background(), project); err != nil {
		t.Fatalf("StopProject() returned error: %v", err)
	}

	if len(rt.stopCalls) != 1 {
		t.Fatalf("expected one StopProject call, got %d", len(rt.stopCalls))
	}
	if !reflect.DeepEqual(rt.stopCalls[0], project) {
		t.Fatalf("StopProject called with wrong project: got %#v want %#v", rt.stopCalls[0], project)
	}
}

func TestServiceRedeployProjectDelegatesThroughEngine(t *testing.T) {
	t.Parallel()

	project := Project{ID: 42, Name: "demo"}
	rt := &serviceFakeRuntime{}
	svc := NewService(NewEngine(rt))
	req := RedeployRequest{
		Project:          project,
		TargetRevision:   "nginx:1.28-alpine",
		PreviousRevision: "nginx:1.27-alpine",
	}

	if err := svc.RedeployProject(context.Background(), req); err != nil {
		t.Fatalf("RedeployProject() returned error: %v", err)
	}

	wantActions := []string{"stop", "deploy"}
	if !reflect.DeepEqual(rt.actions, wantActions) {
		t.Fatalf("unexpected runtime actions: got %v want %v", rt.actions, wantActions)
	}
	if len(rt.deployCalls) != 1 {
		t.Fatalf("expected one DeployRevision call, got %d", len(rt.deployCalls))
	}
	if !reflect.DeepEqual(rt.deployCalls[0].project, project) {
		t.Fatalf("DeployRevision called with wrong project: got %#v want %#v", rt.deployCalls[0].project, project)
	}
	if rt.deployCalls[0].revision != req.TargetRevision {
		t.Fatalf("DeployRevision called with wrong revision: got %q want %q", rt.deployCalls[0].revision, req.TargetRevision)
	}
}

func TestServiceRollbackProjectDelegatesThroughEngine(t *testing.T) {
	t.Parallel()

	project := Project{ID: 42, Name: "demo"}
	rt := &serviceFakeRuntime{}
	svc := NewService(NewEngine(rt))
	req := RollbackRequest{
		Project:  project,
		Revision: "nginx:1.27-alpine",
	}

	if err := svc.RollbackProject(context.Background(), req); err != nil {
		t.Fatalf("RollbackProject() returned error: %v", err)
	}

	if len(rt.deployCalls) != 1 {
		t.Fatalf("expected one DeployRevision call, got %d", len(rt.deployCalls))
	}
	if !reflect.DeepEqual(rt.deployCalls[0].project, project) {
		t.Fatalf("DeployRevision called with wrong project: got %#v want %#v", rt.deployCalls[0].project, project)
	}
	if rt.deployCalls[0].revision != req.Revision {
		t.Fatalf("DeployRevision called with wrong revision: got %q want %q", rt.deployCalls[0].revision, req.Revision)
	}
}

func TestServiceReturnsRuntimeErrors(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("start failed")
	rt := &serviceFakeRuntime{startErr: wantErr}
	svc := NewService(NewEngine(rt))

	err := svc.StartProject(context.Background(), Project{ID: 42, Name: "demo"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("StartProject() error = %v, want %v", err, wantErr)
	}
}

type serviceDeployCall struct {
	project  Project
	revision string
}

type serviceFakeRuntime struct {
	startCalls  []Project
	stopCalls   []Project
	deployCalls []serviceDeployCall
	actions     []string

	startErr  error
	stopErr   error
	deployErr error
}

func (f *serviceFakeRuntime) StartProject(_ context.Context, project Project) error {
	f.startCalls = append(f.startCalls, project)
	f.actions = append(f.actions, "start")
	return f.startErr
}

func (f *serviceFakeRuntime) StopProject(_ context.Context, project Project) error {
	f.stopCalls = append(f.stopCalls, project)
	f.actions = append(f.actions, "stop")
	return f.stopErr
}

func (f *serviceFakeRuntime) DeployRevision(_ context.Context, project Project, revision string) error {
	f.deployCalls = append(f.deployCalls, serviceDeployCall{project: project, revision: revision})
	f.actions = append(f.actions, "deploy")
	return f.deployErr
}
