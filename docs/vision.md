# Vision

Shipdeck is an SSH-first deployment manager for a single VPS.

## Product intent

Make self-hosted deployment operations approachable from a terminal-first workflow. The operator connects over SSH, lands in a Bubble Tea interface, and manages compose-style projects without exposing shell access.

## User and problem

Primary user: a solo operator managing one VPS.

Problem: deployment workflows are often fragmented across SSH commands, compose files, logs, and ad hoc scripts. Shipdeck provides one operational surface for deployment, update checks, and rollback-aware control.

## MVP boundaries

- Single operator
- Existing SSH keys for access control
- Bubble Tea over Wish
- Compose-style project lifecycle
- SQLite as current metadata source of truth
- Local-only HTTP API

## Non-goals right now

- Multi-tenant team support
- Managed cloud control plane
- Complex role/permission matrix
- Broad plugin ecosystem before core workflows are stable

## Current implementation frontier

Planning is complete. Implementation proceeds in vertical slices, proving end-to-end behavior first, then evolving schema and internals from real usage.
