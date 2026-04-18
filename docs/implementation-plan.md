# Shipdeck Implementation Plan

**Goal:** Build an SSH-first, open source Docker deployment manager for a single VPS, with a Bubble Tea TUI, a local HTTP API, and a modular Go backend.

**Architecture:** Shipdeck is split into small Go packages with clear boundaries: SSH/session wiring, TUI presentation, deployment domain logic, registry/source adapters, SQLite state, local API, and Docker execution adapters. The TUI and API are thin clients over the same service layer, and future-changing systems like Docker access, registry lookups, update checks, and secret storage are isolated behind interfaces so they can evolve without rewriting the core.

**Tech Stack:** Go, Bubble Tea, Wish, SQLite, Docker Engine API/CLI hybrid, OpenAPI for `/api/v1`.

---

## Execution mode

Implementation in this repo follows pair-programming navigator mode.

- Use `docs/working-mode.md` as the collaboration contract.
- Work in vertical slices with fast verification loops.
- Use implementation details only when requested.
- Prefer direction and outcome signals over rigid micro-task instructions.

## Active slice tracking

Keep this section updated as work moves so new sessions can continue smoothly.

- Current direction: complete Task 4 by wiring deploy service orchestration to a concrete Docker adapter.
- Current slice: Task 4 service layer plus Docker adapter implementation.
- Success signal: deploy service invokes engine paths through a Docker adapter boundary, with tests proving start, stop, redeploy, and rollback behavior from service entrypoints.

---

### Task 1: Foundation and wiring

**Files:**
- `cmd/shipdeck/main.go`
- `internal/app/app.go`
- `internal/app/wire.go`
- `internal/config/config.go`
- `internal/state/store.go`

Start by writing a small test that proves the composition root can build the app from config without reaching into SSH, Docker, or SQLite internals directly. Then wire the top-level constructors and config shape only as far as needed to assemble the app. Keep this task focused on bootstrapping and dependency boundaries.

Status: the initial constructor wiring and proof test are in place. The current implementation loads config, creates the logger, opens SQLite, runs migrations, and starts the app through `internal/app.Wire(...)`.

### Task 2: SSH entrypoint and TUI shell

**Files:**
- `internal/session/server.go`
- `internal/session/auth.go`
- `internal/ui/model.go`
- `internal/ui/view.go`
- `internal/ui/update.go`

Add tests that prove valid SSH keys start a Bubble Tea session and invalid keys are rejected. Implement the Wish server, key lookup, and a minimal TUI shell that can render a placeholder dashboard. Keep the session layer separate from presentation.

Status: complete. SSH entrypoint, authorized key enforcement, and Bubble Tea session startup are implemented and covered by tests for authorized and unauthorized keys. A minimal dashboard placeholder is in place in the TUI shell.

### Task 3: Project model and SQLite persistence

**Files:**
- `internal/deploy/project.go`
- `internal/deploy/service.go`
- `internal/state/sqlite.go`
- `internal/state/migrations/*.sql`
- `internal/state/repository.go`

Add repository tests for creating, loading, and updating a project with images, watched tags, credentials, and per-project update settings. Then implement the SQLite schema and repository methods needed for project CRUD and secret references. SQLite is the source of truth for project metadata and update state.

Decision (2026-04-11): keep `state.Store` as the persistence facade name to stay backend-agnostic, while using SQLite as the current implementation detail for migrations and repository behavior.

Status: complete. Project aggregate modeling and SQLite persistence are implemented, including migrations and repository methods for create/load/update with images, watched tags, credential references, and per-project update settings/state. Repository tests cover happy paths and transactional rollback behavior.

### Task 4: Docker adapter and deployment engine

**Files:**
- `internal/adapters/docker/client.go`
- `internal/adapters/docker/cli.go`
- `internal/adapters/docker/api.go`
- `internal/deploy/engine.go`
- `internal/deploy/rollback.go`
- `internal/deploy/strategy.go`

Write tests for start, stop, redeploy, and rollback behavior against a fake Docker adapter. Implement the hybrid Docker backend and the deployment engine with a global deployment strategy. Compose-style projects should be manageable through one domain service.

Decision (2026-04-18): deployment revision selection is request-driven in the engine (`RedeployRequest`, `RollbackRequest`). `Project.UpdateState` remains persisted metadata and should not be mixed as an implicit runtime control source for rollback/redeploy orchestration.

Status: in progress. Engine domain contracts and behavior are implemented and covered by tests against a fake runtime adapter, including validation guards and rollback success/failure branches. Remaining work for Task 4 is wiring the deploy service entrypoints and implementing the concrete Docker adapter (`internal/adapters/docker`) plus strategy/rollback extraction files.

### Task 5: Registry/source monitoring and update checks

**Files:**
- `internal/sources/registry.go`
- `internal/sources/digests.go`
- `internal/update/checker.go`
- `internal/update/scheduler.go`
- `internal/update/notifications.go`

Add tests proving a project with watched tags detects a newer digest, records it, and surfaces it as an available update without auto-applying it unless enabled. Implement per-project scheduled checks, OR semantics for watched tags, manual refresh, and optional auto-update behavior.

### Task 6: HTTP API and DTO boundary

**Files:**
- `internal/api/server.go`
- `internal/api/routes.go`
- `internal/api/dto/*.go`
- `internal/api/openapi.yaml` or generated spec output

Add API contract tests for versioned routes, read endpoints, and internal-only write action handling. Expose `/api/v1`, generate OpenAPI, map internal models to DTOs, and keep write operations behind the service layer. The API should have its own contract and not leak TUI-only state.

### Task 7: TUI project workflows

**Files:**
- `internal/ui/model.go`
- `internal/ui/update.go`
- `internal/ui/view.go`
- `internal/ui/projects/*.go`

Add UI tests for dashboard-first navigation, project creation/edit/import, log viewing modes, and manual refresh. Implement the dashboard, project list, project detail, images, logs, and settings screens with keyboard and mouse support. Make the TUI usable for creating, editing, importing, and inspecting projects.

### Task 8: Open source readiness

**Files:**
- `README.md`
- `.github/ISSUE_TEMPLATE/*.md`
- `.github/PULL_REQUEST_TEMPLATE.md`
- `LICENSE`

Add the repo-facing OSS basics after the core MVP is stable. Keep the README focused on the product value and make sure the repository is ready for public iteration without overbuilding launch infrastructure.
