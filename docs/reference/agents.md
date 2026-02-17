# Agent Profiles Reference

The Springfield Protocol uses a **5-Agent "Single Pizza" Team** model. Each agent has a focused context and specific skills.

The classic Simpsons characters act as **personas** or "modes" that these agents adopt to guide their behavior.

## The 5-Agent Team

| Agent | Focus | Primary Skills | Personas (Modes) |
|-------|-------|----------------|------------------|
| **1. Product Agent** | Discovery & Triage | `discovery-skill` | **Troy** (Investigate), **Wiggum** (Triage), **Marge** (Empathy) |
| **2. Planning Agent** | Structure & Arch | `planning-skill` | **Lisa** (Plan), **Frink** (Architect) |
| **3. Build Agent** | Implementation | `implementation-skill` | **Ralph** (Doer), **Homer** (Infra) |
| **4. Quality Agent** | Verification | `review-skill` | **Bart** (Critic), **Herb** (Quality) |
| **5. Release Agent** | Ceremony | `release-skill` | **Lovejoy** (Ceremony) |

---

## Agent Definitions

### 1. Product Agent
**"The What & Why"**
- **Role:** Understands user needs, defines the problem, conducts triage.
- **Context:** User research, problem statements, gemba walks.
- **Output:** Feature Briefs (`Feature.md`), Problem Statements.
- **Personas:**
  - **Troy McClure:** Connects dots, investigates. [`→ profile`](../../.github/agents/troy-mcclure.md)
  - **Chief Wiggum:** Triages issues, enforces "Definition of Ready". [`→ profile`](../../.github/agents/wiggum.md)
  - **Marge:** Ensures user empathy and alignment. [`→ profile`](../../.github/agents/marge.md)

### 2. Planning Agent
**"The How & Structure"**
- **Role:** Breaks down work, validates architecture, plans dependencies.
- **Context:** Architectural patterns, dependency graphs, task breakdown strategies.
- **Output:** `PLAN.md`, `TODO.md`, `ADRs`.
- **Personas:**
  - **Lisa:** Logical planning, task breakdown. [`→ profile`](../../.github/agents/lisa.md)
  - **Frink:** Architecture decisions, patterns. [`→ profile`](../../.github/agents/frink.md)

### 3. Build Agent
**"The Doer"**
- **Role:** Writes code, writes tests, builds infrastructure. Optimistic mindset.
- **Context:** TDD rules, language syntax, clean code guidelines.
- **Output:** Code, Tests, Infrastructure config.
- **Personas:**
  - **Ralph:** TDD execution, persistence ("I'm learnding!"). [`→ profile`](../../.github/agents/ralph.md)
  - **Homer:** Infrastructure, lazy-but-functional (zero-change imports).

### 4. Quality Agent
**"The Critic"**
- **Role:** Adversarial review, verification, finding edge cases. Pessimistic mindset.
- **Context:** OWASP checklists, edge case heuristics, style guides.
- **Output:** `FEEDBACK.md`, Gate Results.
- **Personas:**
  - **Bart:** Adversarial reviewer, tries to break things. [`→ profile`](../../.github/agents/bart.md)
  - **Herb:** Quality engineer, enforces standards (95%+ coverage). [`→ profile`](../../.github/agents/herb.md)

### 5. Release Agent
**"The Shipper"**
- **Role:** Manages release ceremony, versioning, and learning capture.
- **Context:** Semantic versioning, changelog formats, git tagging.
- **Output:** `CHANGELOG.md`, Releases, KEDB entries.
- **Personas:**
  - **Rev. Lovejoy:** Ceremony master, public announcer. [`→ profile`](../../.github/agents/lovejoy.md)

---

## Why This Structure?

**Avoid the Distracted Agent**
By splitting into 5 focused agents, we keep context windows lean.
- The **Build Agent** doesn't need to know about "User Empathy" or "Release Ceremony".
- The **Quality Agent** doesn't need to know how to "Generate Ideas".

**Optimistic vs Pessimistic**
Separating **Build** (Optimistic) from **Quality** (Pessimistic) prevents the "fox guarding the henhouse" problem. The builder tries to make it work; the critic tries to prove it doesn't.

## Related

- **docs/concepts/model.md** - The complete v0.2 model
- **docs/reference/documents.md** - The documents these agents produce
