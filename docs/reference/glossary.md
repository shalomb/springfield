# Springfield Protocol Glossary

## Agent Roles
- **Discovery Agent**: The "Product Owner" (e.g., product-discovery). Focuses on "What" and "Why".
- **Planning Agent**: The "Architect/Tech Lead" (e.g., lisa, frink). Focuses on "How" and structure.
- **Implementation Agent**: The "Developer" (e.g., ralph). Focuses on "Doing" and TDD.
- **Review Agent**: The "Critic" (e.g., bart). Focuses on breaking things and security.
- **Verification Agent**: The "QA" (e.g., herb). Focuses on coverage and correctness.
- **Release Agent**: The "Publisher" (e.g., lovejoy). Focuses on versioning and changelogs.

## Core Documents & Ownership Matrix
| Document | Primary Owner (Creator) | Contributors (Updaters) | Purpose |
| :--- | :--- | :--- | :--- |
| **PLAN.md** | **Planning Agent** | All Agents | **Backlog**: Epic roadmap & task status. |
| **TODO.md** | **Planning Agent** | Implementation Agent | **Sprint**: Transient list of tasks for ONE Epic. |
| **Feature.md** | **Discovery Agent** | Architecture, Learning | **Brief**: Problem, requirements, success criteria. |
| **ADR** | **Architecture Agent** | Review, Planning | **Decision**: Architecture Decision Record. |
| **scenarios.feature** | **Discovery/Build Agent** | Implementation, Testing | **Specs**: Executable BDD acceptance criteria. |
| **FEEDBACK.md** | **Review Agent** | Verification Agent | **Review**: Code review issues and gate results. |
| **CHANGELOG.md** | **Release Agent** | Learning Agent | **History**: Release notes and captured learning. |

## Document Definitions
- **PLAN.md**: The master schedule. Tracks Epics and their states.
- **TODO.md**: The worker's checklist. Exists only while an Epic is being implemented. Deleted upon completion.
- **Feature.md**: The "Why" and "What". Found in `docs/features/`.
- **ADR**: The "Why we chose X". Found in `docs/adr/`.
- **scenarios.feature**: The "Definition of Done". Found in `features/`.
- **FEEDBACK.md**: The "Quality Gate". Generated during PR review.
- **CHANGELOG.md**: The "Product History".

## Work Units
- **Feature**: A high-level capability (e.g., "Logging"). Defined in a `Feature.md`.
- **Epic**: A logical grouping of work delivering a **verifiable increment** (e.g., "Implement Structured Logger"). In `PLAN.md`.
- **Task**: An atomic, executable unit of work (e.g., "Write `logger.py`"). In `TODO.md`.

## Processes & Loops
- **Ralph Wiggum Loop**: The core engine. Pick Task -> Clean Context -> Implement -> Verify.
- **Discovery Diamond**: Diverge (Explore Options) -> Converge (Feature Brief).
- **Delivery Diamond**: Diverge (Implement) -> Converge (Verify).
- **Gemba Walk**: Analysis of the current codebase state.
- **Stateless Resampling**: Starting fresh to prevent error accumulation.

## Metrics
- **Flow Score**: `(Value + Urgency) / Effort`. Used to prioritize Epics.
