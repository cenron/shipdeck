# Architecture

Shipdeck is an SSH-first deployment manager for a single VPS.

## Core pieces

- `cmd/shipdeck` - composition root and program entrypoint
- `internal/session` - SSH auth and Wish server wiring
- `internal/ui` - Bubble Tea screens and interaction state
- `internal/deploy` - project lifecycle, update, and rollback logic
- `internal/sources` - registry tag and digest lookup
- `internal/state` - SQLite persistence and secret storage
- `internal/update` - scheduled update checks and notifications
- `internal/api` - local HTTP API and DTO mapping
- `internal/adapters/docker` - Docker CLI/API hybrid backend

## Folder structure

```text
shipdeck/
├── cmd/
│   └── shipdeck/
│       └── main.go
├── internal/
│   ├── app/
│   ├── session/
│   ├── ui/
│   ├── deploy/
│   ├── sources/
│   ├── state/
│   ├── update/
│   ├── api/
│   └── adapters/
│       └── docker/
├── docs/
├── deploy/
├── build/
└── bin/
```

### Folder uses

- `cmd/shipdeck` - owns startup, config loading, and dependency wiring only.
- `internal/app` - creates the service graph and coordinates startup.
- `internal/session` - owns SSH session lifecycle, auth, and Wish integration.
- `internal/ui` - owns Bubble Tea state, rendering, and user interaction.
- `internal/deploy` - owns project state transitions, deployment actions, and rollback behavior.
- `internal/sources` - owns registry lookups, tag watching, and digest comparison.
- `internal/state` - owns SQLite schema, repositories, and secret storage.
- `internal/update` - owns scheduled checks, update availability, and optional auto-update behavior.
- `internal/api` - owns the local HTTP API surface and DTO contracts.
- `internal/adapters/docker` - owns the Docker CLI/API hybrid execution layer.
- `docs` - repository documentation for architecture and implementation plan.
- `deploy` - deployment-related assets such as service definitions or scripts.
- `build` - build helpers and release artifacts.
- `bin` - compiled local binaries.

## MVP boundaries

- Single operator
- Bare metal install
- Existing SSH keys
- Compose-style projects
- Local-only HTTP API
- SQLite-backed metadata

## Extensibility

- Keep external systems behind narrow interfaces.
- Keep TUI-only state out of API DTOs.
- Use small packages with one responsibility each.
