package state

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

type migration struct {
	version int
	name    string
	sql     string
}

type Store struct {
	db  *sqlx.DB
	log *slog.Logger
}

// NewStore constructs a persistence store backed by the provided database handle.
func NewStore(db *sqlx.DB, log *slog.Logger) *Store {
	return &Store{
		db:  db,
		log: log,
	}
}

// Migrate applies embedded SQL migrations in version order within a single transaction.
func (s *Store) Migrate(ctx context.Context) error {
	s.log.Info("running migrations...")
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration transaction: %w", err)
	}
	defer func(tx *sqlx.Tx) { _ = tx.Rollback() }(tx)

	err = s.createSchemaVersionTable(ctx, tx)
	if err != nil {
		return fmt.Errorf("ensure schema_version table: %w", err)
	}

	currentVersion, err := s.getSchemaVersion(ctx, tx)
	if err != nil {
		return fmt.Errorf("get schema version: %w", err)
	}

	migrations, err := s.loadMigrations()
	if err != nil {
		return fmt.Errorf("load embedded migrations: %w", err)
	}

	for _, m := range migrations {
		if m.version <= currentVersion {
			s.log.Info("skipping migration", "name", m.name)
			continue
		}

		s.log.Info("applying migration", "name", m.name)
		if _, err := tx.ExecContext(ctx, m.sql); err != nil {
			return fmt.Errorf("apply migration v%d (%s): %w", m.version, m.name, err)
		}

		if err := s.setSchemaVersion(ctx, tx, m.version); err != nil {
			return fmt.Errorf("set schema version to %d: %w", m.version, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration transaction: %w", err)
	}
	return nil
}

// executeTx executes a query within a transaction, and rolls back on error
func (s *Store) executeTx(ctx context.Context, query string, args ...any) (int64, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("starting transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("inserting snapshot: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting last insert id: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("committing transaction: %w", err)
	}

	return id, nil
}

// createSchemaVersionTable ensures the schema_version table exists and is seeded with version 0.
func (s *Store) createSchemaVersionTable(ctx context.Context, tx *sqlx.Tx) error {
	_, err := tx.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS schema_version (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			version INTEGER NOT NULL
			)
	`)

	if err != nil {
		return fmt.Errorf("create schema_version table: %w", err)
	}

	// Check to see if we have entry into the schema.
	var version int
	err = tx.GetContext(ctx, &version, `SELECT version FROM schema_version WHERE id = 1`)
	if err == nil {
		return nil
	}

	// If we didn't get rows back, then we insert one. If it's not a no rows error then we fail.
	if errors.Is(err, sql.ErrNoRows) {
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_version (id, version) VALUES (1, 0)`); err != nil {
			return fmt.Errorf("seed schema_version: %w", err)
		}
	} else {
		return fmt.Errorf("read schema_version row: %w", err)
	}

	return nil
}

// loadMigrations reads embedded migration files and returns them sorted by version.
// Files must be named NNN_description.sql where NNN is a zero-padded integer.
func (s *Store) loadMigrations() ([]migration, error) {

	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("read migrations directory: %w", err)
	}

	migrations := make([]migration, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		v, err := s.parseVersion(name)
		if err != nil {
			return nil, err
		}

		// Read the contents of the file
		b, err := migrationFS.ReadFile(fmt.Sprintf("migrations/%s", name))
		if err != nil {
			return nil, fmt.Errorf("read migration file %q: %w", name, err)
		}

		migrations = append(migrations, migration{
			version: v,
			name:    name,
			sql:     string(b),
		})
	}

	sort.Slice(migrations, func(i, j int) bool { return migrations[i].version < migrations[j].version })
	if err := s.validateMigrations(migrations); err != nil {
		return nil, err
	}

	return migrations, nil
}

func (s *Store) validateMigrations(migrations []migration) error {
	for i := range migrations {
		if i == 0 {
			if migrations[i].version != 1 {
				return fmt.Errorf("first migration must be version 1, got %d (%s)", migrations[i].version, migrations[i].name)
			}
			continue
		}

		prev := migrations[i-1]
		curr := migrations[i]

		if curr.version == prev.version {
			return fmt.Errorf("duplicate migration version %d (%s and %s)", curr.version, prev.name, curr.name)
		}
		if curr.version != prev.version+1 {
			return fmt.Errorf("missing migration version between %d (%s) and %d (%s)", prev.version, prev.name, curr.version, curr.name)
		}
	}

	return nil
}

// parseVersion extracts the numeric version prefix from a migration filename.
func (s *Store) parseVersion(name string) (int, error) {
	base := strings.TrimSuffix(name, ".sql")
	prefix, _, ok := strings.Cut(base, "_")
	if !ok || prefix == "" {
		return 0, fmt.Errorf("invalid migration name %q: expected NNN_description.sql", name)
	}

	v, err := strconv.Atoi(prefix)
	if err != nil {
		return 0, fmt.Errorf("invalid migration version %q in %q: %w", prefix, name, err)
	}
	if v < 0 {
		return 0, fmt.Errorf("migration version must be >= 0 in %q", name)
	}

	return v, nil
}

// getSchemaVersion returns the currently applied schema version from schema_version.
func (s *Store) getSchemaVersion(ctx context.Context, tx *sqlx.Tx) (int, error) {
	var v int
	err := tx.GetContext(ctx, &v, `SELECT version FROM schema_version WHERE id = 1`)
	if err != nil {
		return 0, fmt.Errorf("query schema_version row: %w", err)
	}

	return v, err
}

// setSchemaVersion updates the currently applied schema version in schema_version.
func (s *Store) setSchemaVersion(ctx context.Context, tx *sqlx.Tx, v int) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE schema_version
		SET version = ?
		WHERE id = 1
	`, v)
	if err != nil {
		return fmt.Errorf("update schema_version to %d: %w", v, err)
	}

	return nil
}
