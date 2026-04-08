package state

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestMigrateCreatesAndSeedsSchemaVersion(t *testing.T) {
	db := mustOpenTestDB(t)
	t.Cleanup(func() { _ = db.Close() })

	s := NewStore(db)
	if err := s.Migrate(context.Background()); err != nil {
		t.Fatalf("Migrate() returned error: %v", err)
	}

	var version int
	if err := db.Get(&version, `SELECT version FROM schema_version WHERE id = 1`); err != nil {
		t.Fatalf("query schema_version failed: %v", err)
	}
	if version != 1 {
		t.Fatalf("expected schema version 1, got %d", version)
	}
}

func TestMigrateIdempotent(t *testing.T) {
	db := mustOpenTestDB(t)
	t.Cleanup(func() { _ = db.Close() })

	s := NewStore(db)
	if err := s.Migrate(context.Background()); err != nil {
		t.Fatalf("first Migrate() failed: %v", err)
	}
	if err := s.Migrate(context.Background()); err != nil {
		t.Fatalf("second Migrate() failed: %v", err)
	}
}

func TestLoadMigrations(t *testing.T) {
	s := &Store{}
	migrations, err := s.loadMigrations()
	if err != nil {
		t.Fatalf("loadMigrations() returned error: %v", err)
	}
	if len(migrations) == 0 {
		t.Fatal("expected at least one migration")
	}
	if migrations[0].version != 1 {
		t.Fatalf("expected first migration version 1, got %d", migrations[0].version)
	}
}

func TestValidateMigrations(t *testing.T) {
	s := &Store{}

	if err := s.validateMigrations([]migration{{version: 1, name: "001_init.sql"}, {version: 2, name: "002_next.sql"}}); err != nil {
		t.Fatalf("expected valid migrations, got error: %v", err)
	}

	if err := s.validateMigrations([]migration{{version: 2, name: "002_init.sql"}}); err == nil {
		t.Fatal("expected error for first migration not version 1")
	}

	if err := s.validateMigrations([]migration{{version: 1, name: "001.sql"}, {version: 1, name: "001_dup.sql"}}); err == nil {
		t.Fatal("expected error for duplicate versions")
	}

	if err := s.validateMigrations([]migration{{version: 1, name: "001.sql"}, {version: 3, name: "003.sql"}}); err == nil {
		t.Fatal("expected error for missing version")
	}
}

func TestParseVersion(t *testing.T) {
	s := &Store{}

	v, err := s.parseVersion("001_init.sql")
	if err != nil {
		t.Fatalf("expected parse success, got error: %v", err)
	}
	if v != 1 {
		t.Fatalf("expected version 1, got %d", v)
	}

	for _, tc := range []string{"no_underscore.sql", "abc_init.sql", "_init.sql"} {
		if _, err := s.parseVersion(tc); err == nil {
			t.Fatalf("expected parseVersion error for %q", tc)
		}
	}
}

func mustOpenTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	db, err := sqlx.Connect("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	return db
}
