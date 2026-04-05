# Shipdeck

SSH-first Docker deployment management for a single VPS.

Shipdeck gives you a Bubble Tea TUI over SSH for creating, editing, importing, and managing Docker Compose projects without a browser dashboard.

## Demo

```text
$ ssh shipdeck@your-vps

Shipdeck Dashboard
------------------
Projects: 4
Updates available: 1
Last check: 2m ago

[Enter] open project  [r] refresh  [q] quit
```

## What it does

- Connect over SSH and land in a terminal UI
- Create and edit deployment projects in the TUI
- Import existing Compose setups
- Pull, update, and rollback deployments
- Watch image tags like `latest` or pinned versions
- Check for new versions on a schedule
- Keep logs, status, and audit history in one place

## Why Shipdeck

- SSH-native workflow
- Built for self-hosted VPS deployments
- Compose projects, not Kubernetes
- Registry and GHCR image tracking
- Local HTTP API for automation
- Open source and repo-first

## Install

### From source

```bash
git clone git@github.com:cenron/shipdeck.git
cd shipdeck
make build
```

### Run locally

```bash
make run
```

## Development

- `make build`
- `make run`
- `make test`
- `make lint`
- `make dev`

## Docs

- `docs/architecture.md`
- `docs/implementation-plan.md`

## Status

Shipdeck is in active development.

## License

Apache 2.0
