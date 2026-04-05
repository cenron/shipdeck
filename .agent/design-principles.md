# Design Principles

## General

- Never guess about APIs, framework behavior, best practices, or library usage.
- If documentation is unavailable or ambiguous, ask before proceeding.
- Read installed library source when debugging internals instead of relying on memory.
- Prefer the simplest solution that satisfies the current slice.
- Only change what is needed.
- Keep the codebase consistent with existing patterns.
- Use guard clauses and keep the happy path straight through.
- Keep files small and single-purpose.
- Delete dead code.
- Validate only at boundaries.
- Use constructor-based dependency wiring.
- Group related config into dedicated objects.
- Keep interfaces tiny and define them where they are consumed.
- Write tests before implementation where practical.
- Treat tests as sacred.
- Add code coverage tooling and keep coverage at or above 80%.
- Run verification before calling work complete.

## Architecture

- Follow SOLID principles.
- Prefer composition over inheritance.
- Separate what's yours from what's theirs with narrow interfaces only at real boundaries.
- Use one implementation, one location for shared abstractions.
- Keep external adapters out of core business logic.

## Docs

- Read `docs/README.md` before planning work.
- Read `docs/architecture.md` before changing structure.
- Read `docs/implementation-plan.md` before implementing.

## Lessons learned

- Read `.agent/lessons_learned.md` at the start of every task.
- Update it after corrections, unexpected errors, framework quirks, or non-obvious fixes.
