# Springfield Agent Command Center (AGENTS.md)

Welcome, Agent. This document is your primary entry point for navigating the Springfield Protocol. Use it to quickly find the context you need without reading the entire documentation tree.

---

## 🗺️ System Map

### 📍 Core Protocol
- **[README.md](README.md)**: Main project documentation and overview.
- **[QUICK_START.md](QUICK_START.md)**: Fast-track for common workflows.
- **[PLAN.md](PLAN.md)**: The current roadmap and epic status.
- **[TODO.md](TODO.md)**: Your current active tasks.

### 📚 Documentation (Diataxis)
- **[Concepts (The "Why")](docs/concepts/)**: Understanding architecture and models.
- **[How-To (The "How")](docs/how-to/)**: Procedures and walkthroughs.
- **[Reference (The "What")](docs/reference/)**: Technical specs, agents, and glossary.
- **[Standards (The "Rules")](docs/standards/)**: Coding, commits, and protection rules.
- **[ADRs (The "Decisions")](docs/adr/)**: Architectural Decision Records.

### 🤖 Personas & Skills
- **[.github/agents/](.github/agents/)**: All agent persona profiles.
- **[.github/skills/](.github/skills/)**: Reusable automation capabilities.

---

## 🚦 Navigation by Agent

| If you are... | Focus on... |
| :--- | :--- |
| **@Marge (Product)** | `Feature.md`, `docs/concepts/model.md` |
| **@Lisa (Planning)** | `PLAN.md`, `TODO.md`, `docs/standards/` |
| **@Ralph (Build)** | `TODO.md`, `docs/standards/coding-conventions.md`, `docs/standards/atomic-commit-protocol.md` |
| **@Bart (Quality)** | `FEEDBACK.md`, `docs/features/README.md`, `docs/standards/repository-protection.md` |
| **@Lovejoy (Release)** | `CHANGELOG.md`, `docs/standards/git-branching.md` |

---

## 🧠 LLM Guidance (How to use this documentation)

When interacting with this repository, follow these best practices for optimal context management:

1.  **Start here**: Always read `AGENTS.md` and `PLAN.md` first to understand the current state.
2.  **Stateless Mindset**: Assume you have zero memory between cycles. The documentation *is* your memory.
3.  **Atomic Context**: Only read the files relevant to your current task. Avoid reading the entire `docs/` directory to save tokens.
4.  **Truth in Documents**: If `TODO.md` contradicts your instructions, follow `TODO.md`.
5.  **Strict Standards**: Adhere to the [Atomic Commit Protocol](docs/standards/atomic-commit-protocol.md) and [Coding Conventions](docs/standards/coding-conventions.md) at all times.
6.  **TDD First**: For Ralph, never implement code without first writing a failing test and validating it with the appropriate script.

---
*Last Updated: 2026-02-19*
