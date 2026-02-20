# Atomic Commit Protocol (ACP)

**The Golden Rule:** Every commit must leave the repo in a working state.

If I checkout your commit and run `just test`, it must pass. If it doesn't, your commit is invalid.

---

## 1. The Philosophy

We don't do "WIP" commits. We don't do "save point" commits.

A commit is a discrete unit of finished work. It includes:
1.  **The Test:** The proof it works.
2.  **The Code:** The implementation.
3.  **The Docs:** The explanation.

If you commit code without tests, you are breaking the protocol.

## 2. The Format (Conventional Commits)

We follow the [Conventional Commits](https://www.conventionalcommits.org/) spec.

```
<type>(<scope>): <subject>

<body>

<footer>
```

### The Header
*   **Type:** `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`.
*   **Scope:** `core`, `agent`, `docs`, `cli`.
*   **Subject:** 50 chars max. Imperative mood. ("Add feature", not "Added feature").

### The Body
*   **Why:** Explain *why* you made this change. The code explains *what*, the history explains *why*.
*   **Wrap:** 72 characters.

## 3. Examples

### ✅ Good
```
feat(auth): add jwt validation middleware

The legacy session auth was causing race conditions in distributed
deployments. This moves us to stateless JWTs.

Closes #123
```

### ❌ Bad
```
fix bug
```
*(Ralph will reject this commit and laugh at you.)*

### ❌ Bad
```
wip
```
*(Bart will reject this commit and cry.)*

## 4. Workflow Integration

*   **Ralph** creates the commits.
*   **Bart** validates the commits.
*   **Lovejoy** reads the commits to write the Changelog.

If your commits are bad, Lovejoy's changelog will be garbage, and he will be unhappy. Don't make Lovejoy unhappy.
