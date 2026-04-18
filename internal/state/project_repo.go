package state

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cenron/shipdeck/internal/deploy"
	"github.com/jmoiron/sqlx"
)

func (s *Store) CreateProject(ctx context.Context, p deploy.Project) (int64, error) {
	// Write the aggregate root and all child collections in one transaction so
	// callers never observe a partially-created project.
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin create project transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	result, err := tx.ExecContext(ctx, `
		INSERT INTO projects (
			name,
			update_enabled,
			update_auto_apply,
			update_schedule,
			update_available,
			latest_digest,
			last_checked_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		p.Name,
		boolToInt(p.Update.Enabled),
		boolToInt(p.Update.AutoApply),
		p.Update.Schedule,
		boolToInt(p.UpdateState.Available),
		p.UpdateState.LatestDigest,
		nullableTimeString(p.UpdateState.LastCheckedAt),
	)
	if err != nil {
		return 0, fmt.Errorf("insert project: %w", err)
	}

	projectID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("read inserted project id: %w", err)
	}

	if err := insertProjectImages(ctx, tx, projectID, p.Images); err != nil {
		return 0, err
	}
	if err := insertProjectWatchTags(ctx, tx, projectID, p.WatchTags); err != nil {
		return 0, err
	}
	if err := insertProjectCredentials(ctx, tx, projectID, p.CredentialRefs); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit create project transaction: %w", err)
	}

	return projectID, nil
}

func (s *Store) GetProject(ctx context.Context, id int64) (deploy.Project, error) {
	type projectRow struct {
		ID              int64          `db:"id"`
		Name            string         `db:"name"`
		UpdateEnabled   int            `db:"update_enabled"`
		UpdateAutoApply int            `db:"update_auto_apply"`
		UpdateSchedule  string         `db:"update_schedule"`
		UpdateAvailable int            `db:"update_available"`
		LatestDigest    string         `db:"latest_digest"`
		LastCheckedAt   sql.NullString `db:"last_checked_at"`
		CreatedAt       string         `db:"created_at"`
		UpdatedAt       string         `db:"updated_at"`
	}

	var row projectRow
	err := s.db.GetContext(ctx, &row, `
		SELECT id, name, update_enabled, update_auto_apply, update_schedule, update_available, latest_digest, last_checked_at, created_at, updated_at
		FROM projects
		WHERE id = ?
	`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return deploy.Project{}, fmt.Errorf("project %d not found: %w", id, err)
		}
		return deploy.Project{}, fmt.Errorf("select project %d: %w", id, err)
	}

	createdAt, err := parseDBTime(row.CreatedAt)
	if err != nil {
		return deploy.Project{}, fmt.Errorf("parse created_at for project %d: %w", id, err)
	}
	updatedAt, err := parseDBTime(row.UpdatedAt)
	if err != nil {
		return deploy.Project{}, fmt.Errorf("parse updated_at for project %d: %w", id, err)
	}

	project := deploy.Project{
		ID:   row.ID,
		Name: row.Name,
		Update: deploy.UpdateConfig{
			Enabled:   row.UpdateEnabled != 0,
			AutoApply: row.UpdateAutoApply != 0,
			Schedule:  row.UpdateSchedule,
		},
		UpdateState: deploy.ProjectUpdateState{
			Available:    row.UpdateAvailable != 0,
			LatestDigest: row.LatestDigest,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	if row.LastCheckedAt.Valid {
		t, err := parseDBTime(row.LastCheckedAt.String)
		if err != nil {
			return deploy.Project{}, fmt.Errorf("parse last_checked_at for project %d: %w", id, err)
		}
		project.UpdateState.LastCheckedAt = &t
	}

	if err := s.db.SelectContext(ctx, &project.Images, `
		SELECT name, reference, digest
		FROM project_images
		WHERE project_id = ?
		ORDER BY name
	`, id); err != nil {
		return deploy.Project{}, fmt.Errorf("select project images for project %d: %w", id, err)
	}

	if err := s.db.SelectContext(ctx, &project.WatchTags, `
		-- Use aliases that match sqlx's default name mapping for struct fields.
		SELECT image_name AS imagename, tag
		FROM project_watch_tags
		WHERE project_id = ?
		ORDER BY image_name, tag
	`, id); err != nil {
		return deploy.Project{}, fmt.Errorf("select project watch tags for project %d: %w", id, err)
	}

	if err := s.db.SelectContext(ctx, &project.CredentialRefs, `
		-- Use aliases that match sqlx's default name mapping for struct fields.
		SELECT name, secret_ref AS secretref
		FROM project_credentials
		WHERE project_id = ?
		ORDER BY name
	`, id); err != nil {
		return deploy.Project{}, fmt.Errorf("select project credentials for project %d: %w", id, err)
	}

	return project, nil
}

func (s *Store) UpdateProject(ctx context.Context, p deploy.Project) error {
	if p.ID <= 0 {
		return fmt.Errorf("invalid project id %d", p.ID)
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin update project transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	result, err := tx.ExecContext(ctx, `
		UPDATE projects
		SET
			name = ?,
			update_enabled = ?,
			update_auto_apply = ?,
			update_schedule = ?,
			update_available = ?,
			latest_digest = ?,
			last_checked_at = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		p.Name,
		boolToInt(p.Update.Enabled),
		boolToInt(p.Update.AutoApply),
		p.Update.Schedule,
		boolToInt(p.UpdateState.Available),
		p.UpdateState.LatestDigest,
		nullableTimeString(p.UpdateState.LastCheckedAt),
		p.ID,
	)
	if err != nil {
		return fmt.Errorf("update project %d: %w", p.ID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read update result for project %d: %w", p.ID, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("project %d not found", p.ID)
	}

	// Replace collection rows atomically to keep update behavior simple and
	// deterministic for this MVP slice.
	if _, err := tx.ExecContext(ctx, `DELETE FROM project_images WHERE project_id = ?`, p.ID); err != nil {
		return fmt.Errorf("delete project images for project %d: %w", p.ID, err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM project_watch_tags WHERE project_id = ?`, p.ID); err != nil {
		return fmt.Errorf("delete project watch tags for project %d: %w", p.ID, err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM project_credentials WHERE project_id = ?`, p.ID); err != nil {
		return fmt.Errorf("delete project credentials for project %d: %w", p.ID, err)
	}

	if err := insertProjectImages(ctx, tx, p.ID, p.Images); err != nil {
		return err
	}
	if err := insertProjectWatchTags(ctx, tx, p.ID, p.WatchTags); err != nil {
		return err
	}
	if err := insertProjectCredentials(ctx, tx, p.ID, p.CredentialRefs); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit update project transaction: %w", err)
	}

	return nil
}

func insertProjectImages(ctx context.Context, tx *sqlx.Tx, projectID int64, images []deploy.ProjectImage) error {
	stmt, err := tx.PreparexContext(ctx, `
		INSERT INTO project_images (project_id, name, reference, digest)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("prepare project_images insert: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	for _, image := range images {
		if _, err := stmt.ExecContext(ctx, projectID, image.Name, image.Reference, image.Digest); err != nil {
			return fmt.Errorf("insert project image %q for project %d: %w", image.Name, projectID, err)
		}
	}

	return nil
}

func insertProjectWatchTags(ctx context.Context, tx *sqlx.Tx, projectID int64, tags []deploy.WatchTag) error {
	stmt, err := tx.PreparexContext(ctx, `
		INSERT INTO project_watch_tags (project_id, image_name, tag)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("prepare project_watch_tags insert: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	for _, tag := range tags {
		if _, err := stmt.ExecContext(ctx, projectID, tag.ImageName, tag.Tag); err != nil {
			return fmt.Errorf("insert watch tag %q for image %q and project %d: %w", tag.Tag, tag.ImageName, projectID, err)
		}
	}

	return nil
}

func insertProjectCredentials(ctx context.Context, tx *sqlx.Tx, projectID int64, refs []deploy.CredentialRef) error {
	stmt, err := tx.PreparexContext(ctx, `
		INSERT INTO project_credentials (project_id, name, secret_ref)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("prepare project_credentials insert: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	for _, ref := range refs {
		if _, err := stmt.ExecContext(ctx, projectID, ref.Name, ref.SecretRef); err != nil {
			return fmt.Errorf("insert credential ref %q for project %d: %w", ref.Name, projectID, err)
		}
	}

	return nil
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func nullableTimeString(v *time.Time) any {
	if v == nil {
		return nil
	}
	return v.UTC().Format(time.RFC3339Nano)
}

func parseDBTime(value string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.UTC); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("unsupported time format %q", value)
}
