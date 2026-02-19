# Coding Conventions

## Overview
To maintain a high standard of quality and consistency across the Springfield Protocol, all code implementations must follow these language-specific conventions.

## Core Values
- **Readability Over Cleverness**: Code is read much more often than it is written.
- **TDD Always**: Never implement logic without a corresponding test.
- **Atomic Commits**: Every change must adhere to the [Atomic Commit Protocol](atomic-commit-protocol.md).
- **Self-Documenting Code**: Use descriptive variable and function names.

## Go Standards (`cmd/springfield`)
The Springfield command-line tool is built with Go.

1.  **Standard Formatting**: Always run `gofmt` or `goimports` on save.
2.  **Explicit Errors**: Always handle errors explicitly. Do not ignore them with `_`.
3.  **No Package-Level State**: Prefer dependency injection over global variables.
4.  **Interface Segregation**: Keep interfaces small and focused.
5.  **Documentation**: Use `godoc` style comments for exported functions and types.
6.  **Linting**: Must pass `golangci-lint`.

## Python Standards (`scripts/`, `tests/`)
Python is used for automation scripts, testing, and BDD scenarios.

1.  **PEP 8**: Adhere strictly to PEP 8 style guidelines.
2.  **Type Hinting**: Use type hints for all function signatures.
3.  **Modern Tools**: Use `ruff` for linting and formatting.
4.  **Context Managers**: Use `with` statements for resource management (files, network, etc.).
5.  **Docstrings**: Use Google-style docstrings for all modules and major functions.
6.  **Testing**: Use `pytest` and `pytest-bdd` for all tests.

## General Practices
- **KISS (Keep It Simple, Stupid)**: Avoid over-engineering. Focus on the minimal implementation required to pass the tests.
- **DRY (Don't Repeat Yourself)**: Refactor common logic into reusable functions or modules, but avoid premature abstraction.
- **Comment "Why", Not "What"**: Assume the reader can understand the code itself; use comments to explain the reasoning behind complex or non-obvious decisions.

## Enforcement
- **Ralph**: Responsible for applying these conventions during implementation.
- **Bart**: Will reject any PR that does not adhere to these standards.
- **CI**: Automated linting checks will block merging of non-compliant code.

---
*Related: [Atomic Commit Protocol](atomic-commit-protocol.md), [Git Branching Standard](git-branching.md)*
