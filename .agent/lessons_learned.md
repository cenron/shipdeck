# Lessons Learned

Keep this file updated with mistakes, discoveries, and gotchas that matter for Shipdeck.

## Format

- Date
- What happened
- Takeaway

## Entries

- 2026-04-05
  - What happened: Task 1 started with a startup/composition test, and the app wiring ended up carrying DB open/migrate/bootstrap logic in `internal/app/wire.go`.
  - Takeaway: keep `main` thin, but accept a small bootstrap helper when the startup path is still the simplest place to compose config, logger, and persistence.

- 2026-04-11
  - What happened: While planning Task 3, naming was clarified to keep `state.Store` instead of renaming to a SQLite-specific type.
  - Takeaway: keep top-level persistence naming backend-agnostic and treat SQLite as an implementation detail so future storage swaps do not force broad type renames.

- 2026-04-18
  - What happened: Task 4 engine work initially mixed two rollback inputs (`Project.UpdateState.LatestDigest` and `RedeployRequest.PreviousRevision`), which created inconsistent expectations between implementation and tests.
  - Takeaway: use a single request-driven revision model for deploy flows (`RedeployRequest` and `RollbackRequest`) and keep persisted project update state as metadata, not orchestration control input.

- 2026-04-18
  - What happened: Some tests instantiated `Engine` directly and bypassed `NewEngine`, weakening constructor invariant coverage.
  - Takeaway: in tests and production wiring, instantiate through constructors so dependency invariants (for example, non-nil runtime) are enforced consistently.
