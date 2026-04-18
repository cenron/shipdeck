# Working mode

This is the collaboration contract for implementation sessions in this repo.

## Two-phase workflow

### Phase 1: planning in the vault

Use vault notes for vision, feasibility, architecture direction, and decisions.

### Phase 2: implementation in the workspace

Build in this repository using pair-programming mode.

Default behavior:

- AI acts as a navigator, not an instructor.
- Guidance starts high-level.
- No file-by-file implementation instructions unless explicitly requested.

## Build-session cadence

When the user asks "what's next?", respond in this flow:

1. Global direction
2. Current slice to prove
3. Success signal
4. Likely next move

## Detail controls

- `go concrete` means exact files/functions and concrete implementation detail.
- `stay high level` means directional guidance only.

## Working style defaults

- Build in vertical slices.
- Verify each slice before broadening scope.
- Let schema and internals evolve from proven usage.
- Use runtime feedback to drive the next move.

## Keeping this current

Treat this as a living note. Update it when collaboration friction appears and keep `AGENTS.md` aligned in the same edit.
