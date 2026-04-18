package state

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/cenron/shipdeck/internal/deploy"
)

func TestProjectRepositoryCreateAndGetRoundTrip(t *testing.T) {
	t.Parallel()

	db := mustOpenTestDB(t)
	t.Cleanup(func() { _ = db.Close() })

	s := NewStore(db, testLogger())
	ctx := context.Background()
	if err := s.Migrate(ctx); err != nil {
		t.Fatalf("Migrate() returned error: %v", err)
	}

	now := time.Now().UTC().Truncate(time.Second)
	project := deploy.Project{
		ID:   0,
		Name: "example-project",
		Images: []deploy.ProjectImage{
			{Name: "web", Reference: "ghcr.io/acme/web:latest", Digest: "sha256:web"},
			{Name: "worker", Reference: "ghcr.io/acme/worker:latest", Digest: "sha256:worker"},
		},
		WatchTags: []deploy.WatchTag{
			{ImageName: "web", Tag: "latest"},
			{ImageName: "worker", Tag: "stable"},
		},
		CredentialRefs: []deploy.CredentialRef{
			{Name: "db", SecretRef: "secret/database"},
			{Name: "registry", SecretRef: "secret/registry"},
		},
		Update: deploy.UpdateConfig{
			Enabled:   true,
			AutoApply: false,
			Schedule:  "0 */6 * * *",
		},
		UpdateState: deploy.ProjectUpdateState{
			Available:     true,
			LatestDigest:  "sha256:newest",
			LastCheckedAt: &now,
		},
	}

	projectID, err := s.CreateProject(ctx, project)
	if err != nil {
		t.Fatalf("CreateProject() returned error: %v", err)
	}
	if projectID <= 0 {
		t.Fatalf("expected created project ID > 0, got %d", projectID)
	}
	project.ID = projectID

	got, err := s.GetProject(ctx, projectID)
	if err != nil {
		t.Fatalf("GetProject() returned error: %v", err)
	}

	assertProjectCoreEqual(t, got, project)
	if got.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}
	if got.UpdatedAt.IsZero() {
		t.Fatal("expected UpdatedAt to be set")
	}
}

func TestProjectRepositoryUpdateReplacesCollectionsAndSettings(t *testing.T) {
	t.Parallel()

	db := mustOpenTestDB(t)
	t.Cleanup(func() { _ = db.Close() })

	s := NewStore(db, testLogger())
	ctx := context.Background()
	if err := s.Migrate(ctx); err != nil {
		t.Fatalf("Migrate() returned error: %v", err)
	}

	initial := deploy.Project{
		ID:   0,
		Name: "before",
		Images: []deploy.ProjectImage{
			{Name: "web", Reference: "ghcr.io/acme/web:latest", Digest: "sha256:old"},
			{Name: "worker", Reference: "ghcr.io/acme/worker:latest", Digest: "sha256:old-worker"},
		},
		WatchTags: []deploy.WatchTag{
			{ImageName: "web", Tag: "latest"},
			{ImageName: "worker", Tag: "stable"},
		},
		CredentialRefs: []deploy.CredentialRef{
			{Name: "registry", SecretRef: "secret/registry-v1"},
			{Name: "db", SecretRef: "secret/db-v1"},
		},
		Update: deploy.UpdateConfig{Enabled: true, AutoApply: false, Schedule: "0 * * * *"},
		UpdateState: deploy.ProjectUpdateState{
			Available:    true,
			LatestDigest: "sha256:oldest",
		},
	}

	projectID, err := s.CreateProject(ctx, initial)
	if err != nil {
		t.Fatalf("CreateProject() returned error: %v", err)
	}
	if projectID <= 0 {
		t.Fatalf("expected created project ID > 0, got %d", projectID)
	}
	initial.ID = projectID

	updated := deploy.Project{
		ID:   initial.ID,
		Name: "after",
		Images: []deploy.ProjectImage{
			{Name: "web", Reference: "ghcr.io/acme/web:v2", Digest: "sha256:new"},
		},
		WatchTags: []deploy.WatchTag{
			{ImageName: "web", Tag: "v2"},
		},
		CredentialRefs: []deploy.CredentialRef{
			{Name: "registry", SecretRef: "secret/registry-v2"},
		},
		Update: deploy.UpdateConfig{Enabled: false, AutoApply: false, Schedule: ""},
		UpdateState: deploy.ProjectUpdateState{
			Available:     false,
			LatestDigest:  "",
			LastCheckedAt: nil,
		},
	}

	if err := s.UpdateProject(ctx, updated); err != nil {
		t.Fatalf("UpdateProject() returned error: %v", err)
	}

	got, err := s.GetProject(ctx, updated.ID)
	if err != nil {
		t.Fatalf("GetProject() returned error: %v", err)
	}

	assertProjectCoreEqual(t, got, updated)
	if got.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to stay set")
	}
	if got.UpdatedAt.IsZero() {
		t.Fatal("expected UpdatedAt to stay set")
	}
}

func TestProjectRepositoryGetReturnsErrorWhenMissing(t *testing.T) {
	t.Parallel()

	db := mustOpenTestDB(t)
	t.Cleanup(func() { _ = db.Close() })

	s := NewStore(db, testLogger())
	ctx := context.Background()
	if err := s.Migrate(ctx); err != nil {
		t.Fatalf("Migrate() returned error: %v", err)
	}

	if _, err := s.GetProject(ctx, int64(9999)); err == nil {
		t.Fatal("expected GetProject() to return an error for missing project")
	}
}

func TestProjectRepositoryUpdateReturnsErrorWhenMissing(t *testing.T) {
	t.Parallel()

	db := mustOpenTestDB(t)
	t.Cleanup(func() { _ = db.Close() })

	s := NewStore(db, testLogger())
	ctx := context.Background()
	if err := s.Migrate(ctx); err != nil {
		t.Fatalf("Migrate() returned error: %v", err)
	}

	err := s.UpdateProject(ctx, deploy.Project{ID: int64(9999), Name: "missing"})
	if err == nil {
		t.Fatal("expected UpdateProject() to return an error for missing project")
	}
}

func assertProjectCoreEqual(t *testing.T, got, want deploy.Project) {
	t.Helper()

	if got.ID != want.ID {
		t.Fatalf("ID mismatch: got %d want %d", got.ID, want.ID)
	}
	if got.Name != want.Name {
		t.Fatalf("Name mismatch: got %q want %q", got.Name, want.Name)
	}
	if !reflect.DeepEqual(got.Images, want.Images) {
		t.Fatalf("Images mismatch: got %#v want %#v", got.Images, want.Images)
	}
	if !reflect.DeepEqual(got.WatchTags, want.WatchTags) {
		t.Fatalf("WatchTags mismatch: got %#v want %#v", got.WatchTags, want.WatchTags)
	}
	if !reflect.DeepEqual(got.CredentialRefs, want.CredentialRefs) {
		t.Fatalf("CredentialRefs mismatch: got %#v want %#v", got.CredentialRefs, want.CredentialRefs)
	}
	if !reflect.DeepEqual(got.Update, want.Update) {
		t.Fatalf("Update mismatch: got %#v want %#v", got.Update, want.Update)
	}

	if got.UpdateState.Available != want.UpdateState.Available {
		t.Fatalf("UpdateState.Available mismatch: got %t want %t", got.UpdateState.Available, want.UpdateState.Available)
	}
	if got.UpdateState.LatestDigest != want.UpdateState.LatestDigest {
		t.Fatalf("UpdateState.LatestDigest mismatch: got %q want %q", got.UpdateState.LatestDigest, want.UpdateState.LatestDigest)
	}
	assertTimePtrEqual(t, got.UpdateState.LastCheckedAt, want.UpdateState.LastCheckedAt)
}

func assertTimePtrEqual(t *testing.T, got, want *time.Time) {
	t.Helper()

	if got == nil && want == nil {
		return
	}
	if got == nil || want == nil {
		t.Fatalf("time pointer mismatch: got %v want %v", got, want)
	}
	if !got.Equal(*want) {
		t.Fatalf("time mismatch: got %s want %s", got.UTC().Format(time.RFC3339Nano), want.UTC().Format(time.RFC3339Nano))
	}
}
