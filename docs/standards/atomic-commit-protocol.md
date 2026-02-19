# Atomic Commit Protocol (ACP)

## Overview
The Atomic Commit Protocol (ACP) ensures that every commit in the Springfield project is a discrete, functional, and testable unit of work. This practice facilitates easier code reviews, simplifies debugging, and enables reliable rollbacks.

## Core Principles
1. **Atomicity**: A commit should do one thing and one thing only. If a change involves multiple logical components, it should be split into multiple commits.
2. **Testability**: Every commit MUST pass all tests. Never commit code that breaks the build or fails existing tests.
3. **Indivisibility**: A commit should include the implementation, the tests (TDD), and any relevant documentation updates.
4. **Reversibility**: It should be possible to revert a single commit without breaking the system.

## Commit Message Standard

### Structure
Commits follow the Conventional Commits specification with additional Springfield-specific requirements:

```
<type>(<scope>): <Subject line (50 chars max)>

<Body: Detailed explanation of 'why', not 'what'.>
<Reference any issue or task IDs.>

[BREAKING CHANGE: <description>]
```

### Subject Line
- **Length**: Maximum 50 characters.
- **Format**: Capitalized, imperative mood (e.g., "Add feature" instead of "Added feature").
- **Prefix**: Use standard types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`.

### Body
- Wrap at 72 characters.
- Explain the motivation for the change.
- Contrast current behavior with new behavior if applicable.
- Reference `TODO.md` tasks or `PLAN.md` Epics.

## Commit Scope
- `protocol`: Changes to the core Springfield protocol.
- `docs`: Documentation updates.
- `cmd`: Changes to command-line tools.
- `scripts`: Internal scripts and tools.
- `tests`: Test suite improvements.
- `infra`: Infrastructure and CI/CD changes.

## Examples

### Good Commit
```
feat(auth): Implement JWT validation

To improve security, we are moving from session-based auth to JWT.
This commit adds the validation logic and associated unit tests.
Addresses Task 3 in EPIC-002.
```

### Bad Commit
```
fixed some bugs and updated docs
```

## Workflow Integration
- **Ralph (Build)**: Must ensure every commit in a PR follows this protocol.
- **Bart (Quality)**: Will reject PRs that contain non-atomic or poorly described commits.
- **Lovejoy (Release)**: Uses these commits to generate the `CHANGELOG.md`.

---
*Related: [Coding Conventions](coding-conventions.md), [Git Branching Standard](git-branching-standard.md)*
