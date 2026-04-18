package deploy

import "time"

// WatchTag defines a tag on a specific image that is monitored for updates.
type WatchTag struct {
	ImageName string // which project image this applies to
	Tag       string // mutable tag to track (for example, "latest")
}

// ProjectImage describes one container image used by a project.
type ProjectImage struct {
	Name      string // image alias or service label within the project
	Reference string // pull reference, such as registry/repository:tag
	Digest    string // known immutable digest, if available
}

// CredentialRef points to credential data stored outside the project record.
type CredentialRef struct {
	Name      string // logical credential identifier used by configuration
	SecretRef string // reference key/path to externally stored secret value
}

// UpdateConfig defines per-project update-check and auto-apply behavior.
type UpdateConfig struct {
	Enabled   bool   // whether update checks are active for this project
	AutoApply bool   // whether detected updates may be applied automatically
	Schedule  string // check cadence definition (for example, cron expression)
}

// ProjectUpdateState captures the latest known update-check result for a project.
type ProjectUpdateState struct {
	Available     bool       // whether a newer digest is currently known
	LatestDigest  string     // newest discovered digest from watched sources
	LastCheckedAt *time.Time // when sources were last checked; nil means never checked
}

// Project is the deploy domain aggregate for metadata, images, credentials, and update settings.
type Project struct {
	ID             int64              // stable project identifier
	Name           string             // human-readable project name
	Images         []ProjectImage     // container images associated with this project
	WatchTags      []WatchTag         // tags monitored for update availability
	CredentialRefs []CredentialRef    // credential references used by this project
	Update         UpdateConfig       // update-check policy for this project
	UpdateState    ProjectUpdateState // latest known update-check result
	CreatedAt      time.Time          // when the project record was created
	UpdatedAt      time.Time          // when the project record was last modified
}
