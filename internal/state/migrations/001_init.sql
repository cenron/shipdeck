-- projects: aggregate root
CREATE TABLE IF NOT EXISTS projects (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  update_enabled INTEGER NOT NULL DEFAULT 0 CHECK (update_enabled IN (0, 1)),
  update_auto_apply INTEGER NOT NULL DEFAULT 0 CHECK (update_auto_apply IN (0, 1)),
  update_schedule TEXT NOT NULL DEFAULT '',
  update_available INTEGER NOT NULL DEFAULT 0 CHECK (update_available IN (0, 1)),
  latest_digest TEXT NOT NULL DEFAULT '',
  last_checked_at TEXT,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- one project -> many images
CREATE TABLE IF NOT EXISTS project_images (
  project_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  reference TEXT NOT NULL,
  digest TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (project_id, name),
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- one project -> many watched tags
CREATE TABLE IF NOT EXISTS project_watch_tags (
  project_id INTEGER NOT NULL,
  image_name TEXT NOT NULL,
  tag TEXT NOT NULL,
  PRIMARY KEY (project_id, image_name, tag),
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- one project -> many credential references (refs only, no secret values)
CREATE TABLE IF NOT EXISTS project_credentials (
  project_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  secret_ref TEXT NOT NULL,
  PRIMARY KEY (project_id, name),
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_project_images_project_id ON project_images(project_id);
CREATE INDEX IF NOT EXISTS idx_project_watch_tags_project_id ON project_watch_tags(project_id);
CREATE INDEX IF NOT EXISTS idx_project_credentials_project_id ON project_credentials(project_id);
