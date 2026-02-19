# Atomic Commit Protocol (ACP)

## Core Principles
1. **Atomicity**: Each commit should represent a single logical change. It must include the implementation, tests (TDD/BDD), and necessary documentation.
2. **Stability**: The codebase must be in a working state (all tests pass) after every commit.
3. **Traceability**: Commits must follow the Conventional Commits specification and link to tasks/issues where possible.

## Commit Structure
- **Title**: Max 50 characters, capitalized, imperative mood (e.g., `Feat: Add budget enforcer middleware`).
- **Body**: Explain the "why" and "how". Include any breaking changes or technical decisions.
- **Footer**: Reference task IDs or ADRs.

## Workflow (TDD-driven)
1. **RED**: Write a failing test (BDD scenario or unit test).
2. **GREEN**: Implement the minimal code to make the test pass.
3. **REFACTOR**: Clean up the code while keeping tests green.
4. **COMMIT**: Once the feature is complete and verified, commit the atomic unit.

## Standards
- Use `feat:`, `fix:`, `docs:`, `test:`, `refactor:`, `chore:`.
- No "WIP" commits.
- Ensure `just test` passes before committing.
