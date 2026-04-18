package state

import (
	"context"
	"strings"
	"testing"

	"github.com/cenron/shipdeck/internal/deploy"
)

func TestCreateProjectRollsBackWhenChildInsertFails(t *testing.T) {
	t.Parallel()

	db := mustOpenTestDB(t)
	t.Cleanup(func() { _ = db.Close() })

	s := NewStore(db, testLogger())
	ctx := context.Background()
	if err := s.Migrate(ctx); err != nil {
		t.Fatalf("Migrate() returned error: %v", err)
	}

	_, err := s.CreateProject(ctx, deploy.Project{
		Name: "rollback-case",
		Images: []deploy.ProjectImage{
			{Name: "web", Reference: "repo/web:latest"},
			{Name: "web", Reference: "repo/web:v2"},
		},
	})
	if err == nil {
		t.Fatal("expected CreateProject() to fail for duplicate image names")
	}

	var count int
	if err := db.GetContext(ctx, &count, `SELECT COUNT(*) FROM projects`); err != nil {
		t.Fatalf("count projects: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected transaction rollback, found %d projects", count)
	}
}

func TestGetProjectReturnsErrorOnInvalidLastCheckedAt(t *testing.T) {
	t.Parallel()

	db := mustOpenTestDB(t)
	t.Cleanup(func() { _ = db.Close() })

	s := NewStore(db, testLogger())
	ctx := context.Background()
	if err := s.Migrate(ctx); err != nil {
		t.Fatalf("Migrate() returned error: %v", err)
	}

	projectID, err := s.CreateProject(ctx, deploy.Project{Name: "bad-time"})
	if err != nil {
		t.Fatalf("CreateProject() returned error: %v", err)
	}

	if _, err := db.ExecContext(ctx, `UPDATE projects SET last_checked_at = 'not-a-time' WHERE id = ?`, projectID); err != nil {
		t.Fatalf("set invalid last_checked_at: %v", err)
	}

	_, err = s.GetProject(ctx, projectID)
	if err == nil {
		t.Fatal("expected GetProject() to fail on invalid last_checked_at")
	}
	if !strings.Contains(err.Error(), "parse last_checked_at") {
		t.Fatalf("expected last_checked_at parse error, got %v", err)
	}
}

func TestUpdateProjectRejectsInvalidID(t *testing.T) {
	t.Parallel()

	db := mustOpenTestDB(t)
	t.Cleanup(func() { _ = db.Close() })

	s := NewStore(db, testLogger())
	err := s.UpdateProject(context.Background(), deploy.Project{ID: 0, Name: "invalid"})
	if err == nil {
		t.Fatal("expected UpdateProject() to reject ID 0")
	}
}

func TestUpdateProjectRollsBackWhenChildInsertFails(t *testing.T) {
	t.Parallel()

	db := mustOpenTestDB(t)
	t.Cleanup(func() { _ = db.Close() })

	s := NewStore(db, testLogger())
	ctx := context.Background()
	if err := s.Migrate(ctx); err != nil {
		t.Fatalf("Migrate() returned error: %v", err)
	}

	projectID, err := s.CreateProject(ctx, deploy.Project{
		Name: "before",
		CredentialRefs: []deploy.CredentialRef{
			{Name: "registry", SecretRef: "secret/v1"},
		},
	})
	if err != nil {
		t.Fatalf("CreateProject() returned error: %v", err)
	}

	err = s.UpdateProject(ctx, deploy.Project{
		ID:   projectID,
		Name: "after",
		CredentialRefs: []deploy.CredentialRef{
			{Name: "registry", SecretRef: "secret/v2"},
			{Name: "registry", SecretRef: "secret/v3"},
		},
	})
	if err == nil {
		t.Fatal("expected UpdateProject() to fail for duplicate credential names")
	}

	got, err := s.GetProject(ctx, projectID)
	if err != nil {
		t.Fatalf("GetProject() returned error after failed update: %v", err)
	}
	if got.Name != "before" {
		t.Fatalf("expected rollback to preserve name, got %q", got.Name)
	}
	if len(got.CredentialRefs) != 1 || got.CredentialRefs[0].SecretRef != "secret/v1" {
		t.Fatalf("expected rollback to preserve credentials, got %#v", got.CredentialRefs)
	}
}
