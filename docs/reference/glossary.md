# Springfield Protocol Glossary

## Agent Roles
- **Product Agent**: The "Product Owner" (Marge). Focuses on "What" and "Why".
- **Planning Agent**: The "Architect/Tech Lead" (Lisa). Focuses on "How" and structure.
- **Build Agent**: The "Developer" (Ralph). Focuses on "Doing" and TDD.
- **Quality Agent**: The "Critic" (Bart). Focuses on breaking things and verification.
- **Release Agent**: The "Publisher" (Lovejoy). Focuses on versioning and changelogs.

## Core Documents & Ownership Matrix
| Document | Primary Owner (Creator) | Contributors (Updaters) | Purpose |
| :--- | :--- | :--- | :--- |
| **PLAN.md** | **Planning Agent (Lisa)** | All Agents | **Backlog**: Epic roadmap & task status. |
| **TODO.md** | **Planning Agent (Lisa)** | Build Agent | **Sprint**: Transient list of tasks for ONE Epic. |
| **Feature.md** | **Product Agent (Marge)** | Architecture, Learning | **Brief**: Problem, requirements, success criteria. |
| **ADR** | **Planning Agent (Lisa)** | Quality, Planning | **Decision**: Architecture Decision Record. |
| **scenarios.feature** | **Product/Build Agent** | Build, Testing | **Specs**: Executable BDD acceptance criteria. |
| **FEEDBACK.md** | **Quality Agent (Bart)** | Build Agent | **Review**: Code review issues and gate results. |
| **CHANGELOG.md** | **Release Agent (Lovejoy)** | Learning Agent | **History**: Release notes and captured learning. |

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
- **User Story**: Mapped to **BDD Scenarios** in `features/*.feature`. We treat stories as *executable specifications*, not just tickets.
- **Task**: An atomic, executable unit of work (e.g., "Write `logger.py`"). In `TODO.md`.

## Processes & Loops
- **Ralph Wiggum Loop**: The core engine. Pick Task -> Clean Context -> Implement -> Verify.
- **Discovery Diamond**: Diverge (Explore Options) -> Converge (Feature Brief).
- **Delivery Diamond**: Diverge (Implement) -> Converge (Verify).
- **Gemba Walk**: Analysis of the current codebase state.
- **Stateless Resampling**: Starting fresh to prevent error accumulation.

## Metrics
- **Flow Score**: `(Value + Urgency) / Effort`. Used to prioritize Epics.
