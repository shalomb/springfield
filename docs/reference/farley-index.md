# Farley Index Evaluator

> *"A test suite is only as valuable as the confidence it gives you to change the code."*
> — Dave Farley, *Modern Software Engineering*

**Skill:** `farley-index`
**Type:** Diagnostic / Audit Tool
**Track:** Quality (cross-cutting — usable by all agents)

---

## Purpose

The Farley Index is a **quantitative diagnostic** that scores an automated test
suite against Dave Farley's six Properties of Good Tests. It moves beyond code
coverage (a quantity measure) to assess test **health, reliability, and design
integrity**.

The core shift: from **"Do we have tests?"** to **"Can we trust our tests?"**

---

## The Six Properties (The Farley Index Dimensions)

Each dimension is scored **0–10** and contributes equally to the aggregate score.

### 1. Fast
> Tests must provide rapid feedback to support the interactive TDD cycle.

- **Target:** < 10ms per unit test.
- **Why it matters:** Slow tests break the Red–Green–Refactor loop. Developers
  skip or batch slow tests, which eliminates the safety net.
- **What to penalise:**
  - Any test importing or calling real I/O: file system, HTTP, DB, sockets.
  - Test setup that requires external process startup (containers, servers).
  - Execution time distributions with a long tail (> 1s outliers).

### 2. Maintainable
> Tests must survive internal refactoring without breaking.

- **Target:** Tests break only when **observable behaviour** changes, never when
  only internal implementation changes.
- **Why it matters:** Tests coupled to the "how" (not the "what") become
  tripwires, not safety nets. They punish refactoring.
- **What to penalise:**
  - Tests that access private/internal members via reflection.
  - Test names or assertion messages containing internal method names.
  - A high rate of test breakage correlated with non-behavioural commits (VCS
    history signal).
  - Mocks that return exactly the same value they receive — **Mock Tautology
    Theatre** (testing the mock, not the behaviour).

### 3. Repeatable
> The same test must produce the same result on every run, in any environment.

- **Target:** Zero flakiness. Tests must be fully deterministic.
- **Why it matters:** A flaky test is worse than no test — it trains the team
  to ignore failures.
- **What to penalise:**
  - Tests that depend on wall-clock time (`Date.now()`, `datetime.now()`).
  - Tests that depend on execution order (shared mutable state, static globals).
  - Tests that make real network calls or depend on external service availability.
  - Tests that rely on file-system state not owned and reset by the test.

### 4. Atomic
> Each test must verify one thing and be independent of every other test.

- **Target:** One logical assertion per test. Setup is minimal and local.
- **Why it matters:** When a test fails, you need to know *exactly* what broke.
  Fat tests with multiple assertions hide the failure signal.
- **What to penalise:**
  - **Setup-to-Assertion Ratio > 3:1** — if 75%+ of the test is scaffolding,
    the production code is likely not modular enough, or the test is too broad.
  - Shared test fixtures mutated across test cases.
  - `before_all` / `setUpClass` patterns that create global state.
  - Tests with more than one `assert` block covering unrelated concerns.

### 5. Necessary
> Every test must correspond to a real requirement. No test should be redundant.

- **Target:** 100% of production code is demanded by at least one test.
  No tests exist that don't map to a user-observable requirement.
- **Why it matters:** Redundant tests add noise and maintenance cost. Missing
  tests leave real requirements uncovered. This is the TDD proof property.
- **What to penalise:**
  - Production code paths with zero test coverage.
  - Duplicate tests that cover identical behaviour with different data
    (use parametrised tests instead).
  - Tests added after-the-fact that only exercise "happy paths" already covered.

### 6. Understandable
> A test must communicate its intent clearly enough to act as a specification.

- **Target:** A developer (or agent) can understand the test's purpose within
  3–5 lines, without needing to read the production code.
- **Why it matters:** Tests are the living documentation of the system. If they
  require a manual to decode, they fail as specifications.
- **What to penalise:**
  - Test names that describe mechanics (`test_method_returns_value`) rather than
    behaviour (`user_cannot_checkout_with_empty_basket`).
  - Magic literals with no named constant or comment explaining their origin.
  - Deeply nested setup or assertion logic.
  - AAA (Arrange–Act–Assert) structure is not identifiable by visual scan.

---

## Inputs

| Parameter            | Type    | Required | Description                                                  |
| :------------------- | :------ | :------- | :----------------------------------------------------------- |
| `source_code_path`   | String  | Yes      | Path to the production source code directory.                |
| `test_suite_path`    | String  | Yes      | Path to the automated test suite directory.                  |
| `execution_logs_path`| String  | No       | Path to CI/CD or local test run logs (for speed + flakiness).|
| `vcs_history`        | Boolean | No       | If `true`, analyse git log to detect implementation coupling.|

---

## Outputs

```json
{
  "farley_index_score": 7.4,
  "metrics": {
    "fast":           8.0,
    "maintainable":   6.5,
    "repeatable":     9.0,
    "atomic":         7.0,
    "necessary":      6.5,
    "understandable": 7.5
  },
  "red_flags": [
    "Mock Tautology Theatre detected in 4 tests",
    "Setup-to-Assertion ratio > 3:1 in 12 tests",
    "2 tests depend on wall-clock time"
  ],
  "recommendations": [
    "[Maintainable] Replace mock-return-self patterns with behavioural fakes.",
    "[Atomic] Extract shared setup into factory helpers; reduce test surface.",
    "[Repeatable] Inject a clock interface; remove direct DateTime.Now() calls."
  ]
}
```

---

## Scoring Heuristics (Internal Computation Guide)

### Fast
```
base_score = 10
penalise -2 per unique real-I/O import pattern found (HttpClient, SQLConnection,
           File.ReadAll, os.open, requests.get, etc.)
penalise -1 if median test execution time > 100ms (from logs)
penalise -2 if any test execution time > 1000ms (from logs)
floor at 0
```

### Maintainable
```
base_score = 10
penalise -3 if any test accesses private/internal members via reflection
penalise -1 per test whose name contains an internal method name (not a
           public API boundary)
penalise -2 if vcs_history=true and > 20% of refactor commits broke tests
penalise -2 per detected Mock Tautology (mock returns the exact input it receives,
           or mock assertion mirrors the mock setup with no transformation)
floor at 0
```

### Repeatable
```
base_score = 10
penalise -2 per test using real time (Date.now, datetime.now, time.time, new Date())
penalise -2 per test with shared mutable state (global vars, static fields
           mutated in test body)
penalise -2 per test making real network/DB calls not controlled by the test
penalise -1 per test with filesystem dependency outside a temp dir
floor at 0
```

### Atomic
```
base_score = 10
compute setup_lines : assertion_lines ratio per test
penalise -2 if median ratio > 3:1
penalise -1 if median ratio > 2:1
penalise -2 per test with multiple unrelated assert groups
penalise -2 if before_all / setUpClass mutates state shared between tests
floor at 0
```

### Necessary
```
base_score = coverage_percentage / 10  (e.g., 85% coverage = 8.5 base)
penalise -1 per cluster of > 3 tests covering identical code paths with
           non-parametrised duplication
penalise -1 if > 10% of tests are trivial getters/setters with no logic
bonus +0.5 if parametrised test patterns are used consistently
cap at 10
```

### Understandable
```
base_score = 10
penalise -1 per test whose name matches pattern: test_<method>_returns / test_<method>_called
penalise -1 per test with > 2 magic literals (unexplained numeric/string constants)
penalise -2 if AAA structure (Arrange / Act / Assert) is not discernible by
           blank lines or comments
penalise -1 if median test length > 20 lines
floor at 0
```

### Aggregate
```
farley_index_score = mean(fast, maintainable, repeatable, atomic, necessary, understandable)
rounded to 1 decimal place
```

---

## Red Flag Catalogue

| Red Flag                        | Triggered By                                          | Worst-Case Impact       |
| :------------------------------ | :---------------------------------------------------- | :---------------------- |
| **Mock Tautology Theatre**      | Mock returns exactly what it received; no real logic tested | Maintainable ↓↓       |
| **Implementation Coupling**     | Tests break on internal refactor; private member access | Maintainable ↓↓        |
| **Slow Feedback Loop**          | Tests > 1s; real I/O imports                          | Fast ↓↓                 |
| **Clock Dependency**            | `Date.now()` / `datetime.now()` in test body          | Repeatable ↓            |
| **Shared Mutable State**        | `static` / global mutated across tests                | Repeatable + Atomic ↓↓  |
| **Bloated Test Setup**          | Setup-to-Assertion ratio > 3:1                        | Atomic ↓                |
| **Trivial Coverage Gaming**     | Tests on getters/setters with no logic                | Necessary ↓             |
| **Opaque Test Names**           | Names describe mechanics, not behaviour               | Understandable ↓        |
| **Magic Literals**              | Unexplained numeric/string constants in assertions    | Understandable ↓        |

---

## Per-Agent Usage Guide

### Bart (Adversarial Reviewer — Refactor Judge)
Invoke the full evaluation. Focus your output on `red_flags` and
`recommendations`. Use the Farley Index as the evidence base for blocking
a merge if `farley_index_score < 5.0` or if any Maintainability red flag
is present.

> *Bart's question:* "Ralph's tests are Green. But are they Clean?"

### Ralph (TDD Executor)
Run the evaluator on your own test suite **before raising a PR**. Target:
- `fast` ≥ 8.0 (no real I/O in unit tests)
- `necessary` ≥ 8.0 (every line of code is demanded by a test)
- `understandable` ≥ 7.0 (tests read as specifications)

Use the output to self-correct before Bart reviews.

### Lisa (Planning Agent)
Use `maintainable` and `atomic` scores to inform architecture decisions.
A low `maintainable` score is a signal that the current module boundaries
are wrong — tests are coupling to internals because the public API is
insufficient.

### Lovejoy (Release Agent)
Use `repeatable` as a **release gate signal**. If `repeatable < 7.0`,
flaky tests may produce false positives or negatives in the release
pipeline. Flag as a go/no-go risk.

### Marge (Product Agent)
Translate the aggregate score into user-impact language for stakeholders:

| Score Range | Plain-Language Signal                                       |
| :---------- | :---------------------------------------------------------- |
| 8.0 – 10.0  | High confidence. Tests are a reliable safety net.           |
| 6.0 – 7.9   | Moderate confidence. Some brittleness; refactor planned.    |
| 4.0 – 5.9   | Low confidence. Tests may be misleading. Ship with caution. |
| 0.0 – 3.9   | Test Theatre. Coverage metrics are false signals. Stop ship.|

---

## Examples

### Example 1 — Brittle Legacy Suite

**Context:** 90% code coverage, but tests are slow and mock private methods.

```json
{
  "farley_index_score": 3.2,
  "metrics": {
    "fast": 3.0, "maintainable": 2.0, "repeatable": 5.0,
    "atomic": 4.0, "necessary": 4.0, "understandable": 3.0
  },
  "red_flags": [
    "Implementation Coupling — 14 tests access private members via reflection",
    "Slow Feedback Loop — median test time 2.3s; 6 tests hit real DB",
    "Mock Tautology Theatre — 9 mocks return their own input unchanged"
  ],
  "recommendations": [
    "[Maintainable] Refactor tests to assert on public API outcomes only.",
    "[Fast] Replace real DB calls with in-memory fakes or test doubles.",
    "[Maintainable] Delete mock-return-self patterns; test the transformation."
  ]
}
```

### Example 2 — High-Quality TDD Suite

**Context:** Tests run in milliseconds; stable through 3 refactor sprints.

```json
{
  "farley_index_score": 9.4,
  "metrics": {
    "fast": 9.8, "maintainable": 9.5, "repeatable": 9.5,
    "atomic": 9.2, "necessary": 9.0, "understandable": 9.4
  },
  "red_flags": [],
  "recommendations": [
    "Suite follows TDD best practices. Maintain current granularity.",
    "[Necessary] Consider adding parametrised edge-case coverage for boundary inputs."
  ]
}
```

---

## Relationship to the Springfield Delivery Cycle

```
Marge (Feature Brief)
  → Lisa (td status: planned)
    → Ralph (Red)
        ↳ Farley per-test checklist     ← applied to EVERY test at point of writing
    → Ralph (Green → Refactor)
        ↳ Farley Index self-audit       ← full suite score before raising PR
          → Bart (Refactor Judge)       ← Farley Index as adversarial evidence base
            → Lovejoy (Release)         ← Repeatable score as go/no-go gate
```

The Farley Index operates at **two points**, not one:

1. **Shift-left (Ralph):** The per-test checklist is applied inline, at the
   moment of writing. Ralph does not commit a test that fails the checklist.
   This prevents low-quality tests from ever entering the suite.

2. **Gate (Bart):** The full suite audit is run by Bart as the Refactor Judge.
   By the time Bart sees the code, Ralph's inline discipline should mean the
   score is already high. Bart's role is adversarial verification, not first
   discovery.

---

*Reference: Dave Farley — [Modern Software Engineering](https://www.davefarley.net), Continuous Delivery (with Jez Humble)*
