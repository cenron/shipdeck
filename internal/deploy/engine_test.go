package deploy

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestNewEngineSetsRuntime(t *testing.T) {
	t.Parallel()

	rt := &fakeRuntime{}
	e := NewEngine(rt)

	if e == nil {
		t.Fatal("expected NewEngine() to return a non-nil engine")
	}
	if e.runtime != rt {
		t.Fatal("expected NewEngine() to keep the provided runtime")
	}
}

func TestNewEnginePanicsWhenRuntimeNil(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected NewEngine(nil) to panic")
		}
		if msg, ok := r.(string); !ok || msg != "runtime must not be nil" {
			t.Fatalf("expected panic message %q, got %#v", "runtime must not be nil", r)
		}
	}()

	_ = NewEngine(nil)
}

func TestEngineStartDelegatesToRuntime(t *testing.T) {
	t.Parallel()

	project := Project{ID: 42, Name: "demo"}
	rt := &fakeRuntime{}
	e := NewEngine(rt)

	if err := e.Start(context.Background(), project); err != nil {
		t.Fatalf("Start() returned error: %v", err)
	}

	if len(rt.startCalls) != 1 {
		t.Fatalf("expected one StartProject call, got %d", len(rt.startCalls))
	}
	if rt.startCalls[0].ID != project.ID || rt.startCalls[0].Name != project.Name {
		t.Fatalf("StartProject called with wrong project: got %#v want %#v", rt.startCalls[0], project)
	}
}

func TestEngineStopDelegatesToRuntime(t *testing.T) {
	t.Parallel()

	project := Project{ID: 42, Name: "demo"}
	rt := &fakeRuntime{}
	e := NewEngine(rt)

	if err := e.Stop(context.Background(), project); err != nil {
		t.Fatalf("Stop() returned error: %v", err)
	}

	if len(rt.stopCalls) != 1 {
		t.Fatalf("expected one StopProject call, got %d", len(rt.stopCalls))
	}
	if rt.stopCalls[0].ID != project.ID || rt.stopCalls[0].Name != project.Name {
		t.Fatalf("StopProject called with wrong project: got %#v want %#v", rt.stopCalls[0], project)
	}
}

func TestEngineRollbackDeploysRequestedRevision(t *testing.T) {
	t.Parallel()

	project := Project{
		ID:   42,
		Name: "demo",
	}
	rt := &fakeRuntime{}
	e := NewEngine(rt)
	req := RollbackRequest{
		Project:  project,
		Revision: "sha256:previous",
	}

	if err := e.Rollback(context.Background(), req); err != nil {
		t.Fatalf("Rollback() returned error: %v", err)
	}

	if len(rt.deployCalls) != 1 {
		t.Fatalf("expected one DeployRevision call, got %d", len(rt.deployCalls))
	}
	if rt.deployCalls[0].revision != req.Revision {
		t.Fatalf("DeployRevision called with wrong revision: got %q want %q", rt.deployCalls[0].revision, req.Revision)
	}
}

func TestEngineRedeployStopsThenDeploysTargetRevision(t *testing.T) {
	t.Parallel()

	project := Project{ID: 42, Name: "demo"}
	rt := &fakeRuntime{}
	e := NewEngine(rt)

	req := RedeployRequest{
		Project:          project,
		TargetRevision:   "new_revision",
		PreviousRevision: "previous_revision",
	}

	if err := e.Redeploy(context.Background(), req); err != nil {
		t.Fatalf("Redeploy() returned error: %v", err)
	}

	if len(rt.stopCalls) != 1 {
		t.Fatalf("expected one StopProject call, got %d", len(rt.stopCalls))
	}
	if len(rt.deployCalls) != 1 {
		t.Fatalf("expected one DeployRevision call, got %d", len(rt.deployCalls))
	}
	if rt.deployCalls[0].revision != req.TargetRevision {
		t.Fatalf("expected target revision %q, got %q", req.TargetRevision, rt.deployCalls[0].revision)
	}

	wantOrder := []string{"stop", "deploy"}
	if !reflect.DeepEqual(rt.actions, wantOrder) {
		t.Fatalf("unexpected call order: got %v want %v", rt.actions, wantOrder)
	}
}

func TestEngineRedeployReturnsErrorWhenTargetRevisionMissing(t *testing.T) {
	t.Parallel()

	project := Project{ID: 42, Name: "demo"}
	rt := &fakeRuntime{}
	e := NewEngine(rt)

	req := RedeployRequest{
		Project:          project,
		TargetRevision:   "",
		PreviousRevision: "previous_revision",
	}

	err := e.Redeploy(context.Background(), req)
	if err == nil {
		t.Fatal("expected Redeploy() to fail when target revision is missing")
	}
	if !strings.Contains(err.Error(), "no target revision specified") {
		t.Fatalf("expected target revision validation error, got %v", err)
	}
	if len(rt.actions) != 0 {
		t.Fatalf("expected no runtime actions on validation error, got %v", rt.actions)
	}
}

func TestEngineRedeployReturnsErrorWhenPreviousRevisionMissing(t *testing.T) {
	t.Parallel()

	project := Project{ID: 42, Name: "demo"}
	rt := &fakeRuntime{}
	e := NewEngine(rt)

	req := RedeployRequest{
		Project:          project,
		TargetRevision:   "new_revision",
		PreviousRevision: "",
	}

	err := e.Redeploy(context.Background(), req)
	if err == nil {
		t.Fatal("expected Redeploy() to fail when previous revision is missing")
	}
	if !strings.Contains(err.Error(), "no previous revision specified") {
		t.Fatalf("expected previous revision validation error, got %v", err)
	}
	if len(rt.actions) != 0 {
		t.Fatalf("expected no runtime actions on validation error, got %v", rt.actions)
	}
}

func TestEngineRedeployReturnsStopErrorWithoutDeploy(t *testing.T) {
	t.Parallel()

	project := Project{ID: 42, Name: "demo"}
	rt := &fakeRuntime{stopErr: errors.New("stop failed")}
	e := NewEngine(rt)

	req := RedeployRequest{
		Project:          project,
		TargetRevision:   "new_revision",
		PreviousRevision: "previous_revision",
	}

	err := e.Redeploy(context.Background(), req)
	if err == nil {
		t.Fatal("expected Redeploy() to fail when stop fails")
	}
	if !strings.Contains(err.Error(), "stop current") {
		t.Fatalf("expected wrapped stop error, got %v", err)
	}
	if len(rt.deployCalls) != 0 {
		t.Fatalf("expected no deploy calls after stop error, got %d", len(rt.deployCalls))
	}

	wantOrder := []string{"stop"}
	if !reflect.DeepEqual(rt.actions, wantOrder) {
		t.Fatalf("unexpected call order: got %v want %v", rt.actions, wantOrder)
	}
}

func TestEngineRedeployRollsBackWhenTargetDeployFailsAndRollbackSucceeds(t *testing.T) {
	t.Parallel()

	project := Project{ID: 42, Name: "demo"}
	rt := &fakeRuntime{deployErrs: []error{errors.New("deploy failed"), nil}}
	e := NewEngine(rt)

	req := RedeployRequest{
		Project:          project,
		TargetRevision:   "new_revision",
		PreviousRevision: "previous_revision",
	}

	err := e.Redeploy(context.Background(), req)
	if err == nil {
		t.Fatal("expected Redeploy() to return error when target deploy fails")
	}
	if !strings.Contains(err.Error(), "rollback succeeded") {
		t.Fatalf("expected rollback success context, got %v", err)
	}

	if len(rt.deployCalls) != 2 {
		t.Fatalf("expected target and rollback DeployRevision calls, got %d", len(rt.deployCalls))
	}
	if rt.deployCalls[0].revision != req.TargetRevision {
		t.Fatalf("target deploy used wrong revision: got %q want %q", rt.deployCalls[0].revision, req.TargetRevision)
	}
	if rt.deployCalls[1].revision != req.PreviousRevision {
		t.Fatalf("rollback used wrong revision: got %q want %q", rt.deployCalls[1].revision, req.PreviousRevision)
	}

	wantOrder := []string{"stop", "deploy", "deploy"}
	if !reflect.DeepEqual(rt.actions, wantOrder) {
		t.Fatalf("unexpected call order: got %v want %v", rt.actions, wantOrder)
	}
}

func TestEngineRedeployReturnsCombinedErrorWhenRollbackFails(t *testing.T) {
	t.Parallel()

	project := Project{ID: 42, Name: "demo"}
	rt := &fakeRuntime{deployErrs: []error{errors.New("target deploy failed"), errors.New("rollback failed")}}
	e := NewEngine(rt)

	req := RedeployRequest{
		Project:          project,
		TargetRevision:   "new_revision",
		PreviousRevision: "previous_revision",
	}

	err := e.Redeploy(context.Background(), req)
	if err == nil {
		t.Fatal("expected Redeploy() to return error when rollback fails")
	}
	if !strings.Contains(err.Error(), "deploy failed") || !strings.Contains(err.Error(), "rollback failed") {
		t.Fatalf("expected combined deploy and rollback error context, got %v", err)
	}

	if len(rt.deployCalls) != 2 {
		t.Fatalf("expected target and rollback DeployRevision calls, got %d", len(rt.deployCalls))
	}

	wantOrder := []string{"stop", "deploy", "deploy"}
	if !reflect.DeepEqual(rt.actions, wantOrder) {
		t.Fatalf("unexpected call order: got %v want %v", rt.actions, wantOrder)
	}
}

type deployCall struct {
	project  Project
	revision string
}

type fakeRuntime struct {
	startCalls  []Project
	stopCalls   []Project
	deployCalls []deployCall
	actions     []string

	startErr   error
	stopErr    error
	deployErr  error
	deployErrs []error
}

func (f *fakeRuntime) StartProject(_ context.Context, project Project) error {
	f.startCalls = append(f.startCalls, project)
	f.actions = append(f.actions, "start")
	return f.startErr
}

func (f *fakeRuntime) StopProject(_ context.Context, project Project) error {
	f.stopCalls = append(f.stopCalls, project)
	f.actions = append(f.actions, "stop")
	return f.stopErr
}

func (f *fakeRuntime) DeployRevision(_ context.Context, project Project, revision string) error {
	f.deployCalls = append(f.deployCalls, deployCall{project: project, revision: revision})
	f.actions = append(f.actions, "deploy")
	if len(f.deployErrs) > 0 {
		err := f.deployErrs[0]
		f.deployErrs = f.deployErrs[1:]
		return err
	}
	return f.deployErr
}
