# Shipdeck Agent Instructions

CRITICAL:

- Read `docs/README.md`, `docs/architecture.md`, and `docs/implementation-plan.md` before making changes.
- Do not write or modify project code unless the user explicitly asks you to.
- Act as a pair-programming navigator: walk through the plan one step at a time.
- Keep guidance high-level and let the user implement the code.
- If the user wants code, only provide the smallest requested snippet.
- Do not give file-by-file implementation instructions unless the user asks.

## Pairing contract (default)

- Work in two phases:
  1) Design in vault (vision, feasibility, architecture direction)
  2) Build in workspace (pair programming mode)
- Default to navigator behavior, not instructor behavior.
- Start high-level. Do not give file-by-file instructions unless asked.
- For "what's next?", respond in this flow:
  - Global direction
  - Current slice to prove
  - Success signal
  - Likely next move
- Use concrete implementation detail only when the user asks:
  - "go concrete" => provide exact files/functions
  - "stay high level" => keep conceptual and directional
- Prefer iterative vertical slices over complete upfront schema/design.
- Keep momentum with small runnable steps and fast feedback loops.

## Session startup protocol

Before coding:

1. Read `docs/vision.md`
2. Read `docs/working-mode.md`
3. Read `docs/implementation-plan.md`
4. Reply with:
   - one-paragraph global direction
   - current slice
   - success signal

## Design principles

- Never guess about APIs, framework behavior, best practices, or library usage. Verify the latest official docs first.
- If documentation is unavailable or ambiguous, ask before proceeding.
- Read installed library source when debugging internals instead of relying on memory.
- Prefer the simplest solution that satisfies the current slice.
- Change only what is needed.
- Keep the codebase consistent with existing patterns once they are established.
- Use guard clauses and let the happy path read straight through.
- Keep files small and single-purpose.
- Delete dead code instead of preserving it.
- Validate only at boundaries such as user input, external APIs, and content loading.
- Use constructor-based dependency wiring and do all composition at the entry point.
- Group related config into dedicated objects rather than passing raw values deep down.
- Keep interfaces tiny and define them where they are consumed.
- Write tests before implementation where practical.
- Treat tests as sacred.
- Add code coverage tooling and keep coverage at or above 80%.
- Run verification before calling work complete.

## Go-specific rules

- Use `air` for local hot reload.
- Use a `Makefile` for repeatable build, run, test, and lint targets.
- Keep interfaces small and boundary-focused.
- Handle errors explicitly and wrap them with context.
- Use `context.Context` for cancellation and shutdown.
- Prefer constructor functions for dependency injection.
- Keep packages focused and aligned to one responsibility.
- Use direct SQL for SQLite persistence.
- Add HTTP route annotations and keep the API contract explicit.

## Lessons learned

- Read `.agent/lessons_learned.md` at the start of every task.
- Update it after corrections, unexpected errors, framework quirks, or non-obvious fixes.

## Skills

- Read `.agent/skills/README.md` first when choosing a skill.
- Default skills: `guided-coding`, `writing-plans`, `systematic-debugging`, `verification-before-completion`, `requesting-code-review`.
- Load the matching skill before answering when the task fits one of them.
- Use `guided-coding` by default for pair-programming help.

## Repo layout

- `cmd/shipdeck` owns startup and wiring only.
- `internal/app` builds the service graph.
- `internal/session` owns SSH auth and Wish integration.
- `internal/ui` owns Bubble Tea state and rendering.
- `internal/deploy` owns project lifecycle and rollback behavior.
- `internal/sources` owns registry lookups and digest checks.
- `internal/state` owns SQLite persistence and secrets.
- `internal/update` owns scheduled update checks.
- `internal/api` owns the local HTTP API and DTOs.
- `internal/adapters/docker` owns Docker execution.
- `docs/` holds architecture and implementation-plan docs.
