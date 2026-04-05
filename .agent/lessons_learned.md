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
