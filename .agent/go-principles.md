# Go Principles

## Tooling

- Use `air` for hot reloading during development.
- Keep a `Makefile` with `build`, `run`, `test`, and `lint`.

## Idioms

- Keep interfaces small.
- Define interfaces where they are consumed.
- Handle errors explicitly.
- Wrap errors with context.
- Use constructor functions for dependency injection.
- Bundle config into dedicated structs.
- Keep packages focused and small.
- Use `context.Context` for cancellation and shutdown.
- Prefer channels for communication and mutexes for shared state.
- Use `errgroup` for coordinated concurrency when needed.
- Keep files small and single-purpose.
- Delete dead code instead of preserving it.

## Persistence and API

- Use direct SQL for SQLite.
- Do not introduce an ORM.
- Keep HTTP route wiring in the composition root.
- Annotate API endpoints explicitly when they are added.

## Testing

- Write tests before implementation where practical.
- Treat tests as sacred.
- Keep coverage at or above 80%.
