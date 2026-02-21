# Adzic Index — BDD Specification Quality Rubric

> *"Behaviour-Driven Development is not about testing. It is about communication.
> A feature file that cannot be read by a non-programmer has failed at its primary job."*
> — Gojko Adzic, *Specification by Example*

**Companion to:** [`farley-index.md`](farley-index.md) (TDD unit test quality)
**Skill:** `adzic-index` *(see `.pi/agent/skills/adzic-index/SKILL.md`)*

---

## Purpose

The Farley Index measures whether your **unit tests** are trustworthy.
The Adzic Index measures whether your **BDD specifications** are *useful* —
as communication tools, as living documentation, and as acceptance criteria.

A `.feature` file can have 100% step-definition coverage and still be worthless
if it is written in a way only a developer can read, or if it tests mechanics
rather than business outcomes.

The core shift: from **"Do our scenarios pass?"** to **"Do our scenarios
communicate the right thing to the right people?"**

---

## The Six Properties (The Adzic Index Dimensions)

Each dimension is scored **0–10** and contributes equally to the aggregate.

---

### 1. Business-Readable
> A scenario must be understandable by a non-technical stakeholder with no
> coaching.

- **Target:** A product manager, domain expert, or end user can read the
  scenario and confirm it matches their intent.
- **Why it matters:** If only developers can read the spec, it is not a shared
  understanding — it is just a test with extra steps.
- **What to penalise:**
  - Steps that expose technical implementation (`Given the database contains a
    row where user_id = 42`).
  - Steps using code-level concepts (`When I POST to /api/v1/users`).
  - Gherkin that reads like test code rather than plain English.

**Good:**
```gherkin
Given a customer with an empty basket
When they attempt to checkout
Then they should see "Your basket is empty"
```

**Bad:**
```gherkin
Given the CartService is initialised with userId="usr_42"
When I call checkout() with an empty List<Item>
Then the HTTP response code should be 400
```

---

### 2. Intention-Revealing
> Each scenario must describe *what* the system should do, not *how* it does it.

- **Target:** Scenarios survive an internal refactor unchanged. The step
  definitions change; the `.feature` file does not.
- **Why it matters:** Scenarios coupled to implementation become brittle
  specifications that punish change — the BDD equivalent of Mock Tautology
  Theatre.
- **What to penalise:**
  - Steps that name internal classes, methods, or modules.
  - Scenarios that describe UI mechanics (`When I click the button with id
    "submit-btn"`) rather than user intent (`When I submit the order`).
  - Steps that mirror the implementation's method signature.

---

### 3. Focused
> Each scenario must demonstrate exactly one behaviour or business rule.

- **Target:** One scenario = one rule. If a scenario title uses "and" to
  describe two outcomes, it should be two scenarios.
- **Why it matters:** When a broad scenario fails, you can't tell which rule
  broke. A focused scenario is a precise diagnostic signal.
- **What to penalise:**
  - More than one `Then` clause covering unrelated outcomes.
  - Scenarios that set up multiple unrelated preconditions.
  - `Scenario Outline` used to cover fundamentally different behaviours rather
    than variations of the same rule.
  - **Step ratio:** If a scenario has more `Given` steps than `Then` steps by
    a ratio > 3:1, it is likely testing too many things at once.

---

### 4. Declarative (Not Imperative)
> Scenarios must declare the *state of the world*, not script a sequence of
> UI actions.

- **Target:** Steps read as domain-level facts and outcomes, not
  click-by-click instructions.
- **Why it matters:** Imperative scenarios are fragile — a UI change breaks
  every scenario that clicks through it. Declarative scenarios survive
  implementation change.
- **What to penalise:**
  - `Given` steps that walk through a multi-step login or navigation flow
    instead of declaring `Given a logged-in admin user`.
  - Sequences of `When` steps that simulate user interaction (`When I click
    "Next"` / `When I fill in "Email"` / `When I click "Submit"`).
  - More than 2 sequential `When` steps (a single user intent should not
    require multiple steps to express).

**Imperative (bad):**
```gherkin
When I navigate to the login page
And I enter "admin" in the username field
And I enter "password123" in the password field
And I click the "Login" button
And I click "Orders" in the navigation menu
```

**Declarative (good):**
```gherkin
Given I am logged in as an admin
When I view the orders dashboard
```

---

### 5. Minimal Context
> A scenario must include only the context *essential* to the behaviour being
> demonstrated. Everything else is noise.

- **Target:** Every `Given` line is necessary to understand *why* the `When`
  produces the `Then`. Remove any `Given` and the scenario should become
  ambiguous or incomplete.
- **Why it matters:** Over-specified context makes scenarios brittle and hard
  to read. It signals that the scenario is testing an over-complex system, or
  that the author didn't understand the minimum precondition.
- **What to penalise:**
  - `Given` steps that set up data or state irrelevant to the specific rule
    being tested.
  - Long `Background` blocks that apply across all scenarios in a feature but
    are only relevant to some.
  - Scenarios that require reading a long chain of `And` clauses before the
    business rule becomes clear.

---

### 6. Living (Executable and Passing)
> A specification that does not execute is a document. A specification that
> fails silently is a lie.

- **Target:** Every scenario has a corresponding step definition. Every
  scenario is green in CI. Pending/skipped scenarios are explicitly tracked
  as debt.
- **Why it matters:** A `.feature` file with pending steps gives false
  confidence — it looks like a spec but has no enforcement value.
- **What to penalise:**
  - Scenarios with `@pending` or `@skip` tags not backed by a tracked issue.
  - Scenarios whose step definitions are stubbed with `pass` / `pending()`.
  - Feature files that exist only in `docs/` but have no corresponding
    executable counterpart in `tests/`.
  - Scenarios that are green but whose step definitions assert nothing
    (empty step bodies).

---

## Scoring Heuristics (Internal Computation Guide)

### Business-Readable
```
base_score = 10
penalise -2 per step containing a technical identifier (camelCase method name,
           HTTP verb, SQL keyword, class name, URL path)
penalise -1 per step containing an internal ID format (uuid, numeric PK)
penalise -2 if a non-technical reviewer could not parse > 50% of steps unaided
floor at 0
```

### Intention-Revealing
```
base_score = 10
penalise -3 per step that names an internal module, class, or method
penalise -2 per step describing UI mechanics (click, fill in, navigate to)
           rather than user intent
penalise -2 if the .feature file would need to change during an internal
           refactor (vcs_history signal)
floor at 0
```

### Focused
```
base_score = 10
penalise -2 per scenario with > 1 unrelated Then clause
penalise -2 per scenario whose title contains " and " joining two outcomes
penalise -1 if Given:Then step ratio > 3:1 (median across scenarios)
penalise -1 per Scenario Outline covering fundamentally different rules
           (vs. legitimate data variations)
floor at 0
```

### Declarative
```
base_score = 10
penalise -2 per scenario with > 2 sequential When steps
penalise -2 per Given block that walks through a UI flow rather than
           declaring a precondition state
penalise -1 per step starting with "I click", "I navigate", "I fill in",
           "I select from dropdown" (UI mechanics markers)
floor at 0
```

### Minimal Context
```
base_score = 10
penalise -1 per Given step that can be removed without making the scenario
           ambiguous (requires human judgement / LLM analysis)
penalise -2 if Background block contains steps not used by > 50% of scenarios
penalise -1 if median Given count per scenario > 4
floor at 0
```

### Living
```
base_score = 10
penalise -4 per scenario with no step definition (undefined steps)
penalise -3 per scenario tagged @pending or @skip without a linked issue
penalise -3 per step definition with an empty body (no assertions)
penalise -2 if feature file exists only in docs/ with no executable counterpart
bonus +0 (no bonus — living is the baseline expectation, not a positive signal)
floor at 0
```

### Aggregate
```
adzic_index_score = mean(business_readable, intention_revealing, focused,
                         declarative, minimal_context, living)
rounded to 1 decimal place
```

---

## Red Flag Catalogue

| Red Flag                        | Triggered By                                              | Worst-Case Impact              |
| :------------------------------ | :-------------------------------------------------------- | :----------------------------- |
| **Gherkin as Code**             | Steps expose method names, SQL, HTTP verbs, class names   | Business-Readable ↓↓           |
| **UI Script Theatre**           | Sequential click/navigate/fill steps instead of intent    | Declarative + Intention ↓↓     |
| **Omnibus Scenario**            | Multiple unrelated `Then` clauses; title has " and "      | Focused ↓↓                     |
| **Phantom Specification**       | Pending/skipped steps; empty step bodies; docs-only files | Living ↓↓                      |
| **Over-Specified Context**      | `Given` steps irrelevant to the rule under test           | Minimal Context ↓              |
| **Implementation Coupling**     | Scenario breaks on internal refactor (not behaviour change)| Intention-Revealing ↓↓        |
| **ID Leakage**                  | UUIDs, numeric PKs, internal identifiers in steps         | Business-Readable ↓            |

---

## Per-Agent Usage Guide

### Marge (Product Agent) — Primary Author
Marge writes the Feature Brief. The `.feature` file is her **executable
translation** of the brief into Gherkin. She owns the `Business-Readable` and
`Intention-Revealing` dimensions.

> *Marge's question when writing a scenario:*
> *"Can our most non-technical stakeholder read this and say 'yes, that's what I wanted'?"*

**Marge's per-scenario checklist:**
- [ ] Would a non-programmer understand every step without explanation?
- [ ] Does the scenario title describe a business rule, not a test case?
- [ ] Am I describing *what* the system does, not *how* it does it?
- [ ] Have I avoided internal system identifiers (IDs, URLs, class names)?

### Lisa (Planning Agent)
Uses the `Focused` and `Living` dimensions to validate that each scenario maps
to exactly one acceptance criterion in the PLAN, and that all scenarios have
executable counterparts before an Epic is marked `Ready`.

> *Lisa's gate:* No Epic moves to `In Progress` if any scenario in its feature
> file is `@pending` without a linked task in TODO.md.

### Ralph (Build Agent)
Ralph implements step definitions. He uses the `Declarative` and
`Minimal Context` dimensions to validate that he is not writing step definitions
that reach into implementation details, and that scenarios don't require
excessive fixture setup.

> *Ralph's signal:* If writing a step definition requires more than 5 lines of
> setup code, the scenario may be over-specified or the system may not be
> modular enough. Escalate to Lisa.

**Ralph's pre-PR self-audit thresholds:**
- `living` ≥ 9.0 — all steps defined and asserting something
- `declarative` ≥ 7.0 — no UI script scenarios
- `focused` ≥ 7.0 — no omnibus scenarios

### Bart (Quality Agent) — Refactor Judge
Bart runs the full Adzic Index audit as part of his review. A low
`intention_revealing` score is a signal that scenarios will break during
refactoring — the BDD equivalent of implementation-coupled unit tests.

> *Bart's block threshold:* `adzic_index_score < 5.0` or any `Phantom
> Specification` red flag = recommend blocking merge.

### Lovejoy (Release Agent)
Uses `living` as a release gate. A scenario that is `@pending` at release time
means an acceptance criterion is unverified. This is a go/no-go signal.

| Score Range | Plain-Language Signal                                               |
| :---------- | :------------------------------------------------------------------ |
| 8.0 – 10.0  | Specifications are clear, executable, and stakeholder-readable.     |
| 6.0 – 7.9   | Mostly good. Some imperative or over-specified scenarios. Refactor. |
| 4.0 – 5.9   | Specifications are brittle or unclear. Ship with explicit risk note.|
| 0.0 – 3.9   | Specification Theatre. Feature files are not living documentation.  |

---

## Examples

### Example 1 — Specification Theatre

**Context:** Feature file with pending steps, UI-scripted scenarios, and
technical identifiers in steps.

```json
{
  "adzic_index_score": 2.8,
  "metrics": {
    "business_readable": 3.0,
    "intention_revealing": 2.0,
    "focused": 5.0,
    "declarative": 2.0,
    "minimal_context": 4.0,
    "living": 1.0
  },
  "red_flags": [
    "Gherkin as Code — 6 steps contain HTTP verbs or class names",
    "UI Script Theatre — 4 scenarios use sequential click/navigate steps",
    "Phantom Specification — 3 scenarios are @pending with no linked issue"
  ],
  "recommendations": [
    "[Business-Readable] Replace 'POST to /api/v1/orders' with 'a customer places an order'.",
    "[Declarative] Collapse login click-sequence into 'Given I am logged in as <role>'.",
    "[Living] Link pending scenarios to TODO.md tasks or delete them."
  ]
}
```

### Example 2 — High-Quality Living Specification

**Context:** Scenarios written collaboratively with a product owner; all steps
defined and green.

```json
{
  "adzic_index_score": 9.1,
  "metrics": {
    "business_readable": 9.5,
    "intention_revealing": 9.0,
    "focused": 9.5,
    "declarative": 9.0,
    "minimal_context": 8.5,
    "living": 9.0
  },
  "red_flags": [],
  "recommendations": [
    "Consider adding Scenario Outlines for edge-case data variations to reduce duplication."
  ]
}
```

---

## Relationship to the Springfield Delivery Cycle

```
Marge (Feature Brief)
  ↳ Adzic per-scenario checklist     ← applied when WRITING the .feature file
    → Lisa (td status: planned)
        ↳ Adzic Living check         ← no Epic is Ready if scenarios are @pending
          → Ralph (step definitions)
              ↳ Adzic self-audit     ← declarative ≥ 7.0, living ≥ 9.0 before PR
                → Bart (audit)       ← full Adzic Index as adversarial evidence base
                  ↳ ADR status change (Proposed -> Accepted)
                    → Lovejoy        ← living score as release go/no-go gate
```

### Farley + Adzic: The Complete Quality Picture

| Layer | Tool | Question |
| :---- | :--- | :------- |
| BDD Specification | **Adzic Index** | Is this scenario a clear, executable, stakeholder-readable specification? |
| TDD Unit Test | **Farley Index** | Is this test fast, honest, and implementation-independent? |

The two indices are complementary. A system with a high Farley score but a low
Adzic score has trustworthy unit tests but unmaintainable acceptance criteria.
A system with a high Adzic score but a low Farley score has clear specifications
but brittle implementation-level tests.

**Target state:** Both indices ≥ 7.0. This is "Clean Code that Works,
clearly specified."

---

*Reference: Gojko Adzic — [Specification by Example](https://gojko.net/books/specification-by-example/),
[Fifty Quick Ideas to Improve Your Tests](https://fiftyquickideas.com/fifty-quick-ideas-to-improve-your-tests/)*
