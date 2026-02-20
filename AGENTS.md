# Springfield Agent Command Center

Welcome to the team.

This is your map. You don't need to read the entire documentation tree to do your job—in fact, please don't. Use this file to find exactly what you need, get the context, and get to work.

---

## 🗺️ The System

### 📍 Core
- **[README.md](README.md)**: What this project is.
- **[QUICK_START.md](QUICK_START.md)**: How to get moving fast.
- **[PLAN.md](PLAN.md)**: The roadmap. What we're building right now.
- **[TODO.md](TODO.md)**: The task list. What you're building right now.

### 📚 The Docs (Diataxis)
We organize documentation by purpose, not just topic.
- **[Concepts](docs/concepts/)**: The "Why." Architecture and mental models.
- **[How-To](docs/how-to/)**: The "How." Step-by-step guides.
- **[Reference](docs/reference/)**: The "What." Specs, glossaries, and hard data.
- **[Standards](docs/standards/)**: The "Rules." How we code and commit.
- **[ADRs](docs/adr/)**: The "Decisions." Why we did it this way.

### 🤖 The Crew
- **[.github/agents/](.github/agents/)**: Who you are.
- **[.pi/agent/skills/](.pi/agent/skills/)**: What you can do.

---

## 🚦 Who are you?

| Role | Read this first |
| :--- | :--- |
| **@Marge (Product)** | `Feature.md` and `docs/concepts/model.md`. Understand the user. |
| **@Lisa (Planning)** | `PLAN.md` and `docs/standards/`. Keep us organized. |
| **@Ralph (Build)** | `TODO.md`. Read the `coding-conventions.md` before you write a line of code. |
| **@Bart (Quality)** | `FEEDBACK.md`. Break things. Check `repository-protection.md`. |
| **@Lovejoy (Release)** | `CHANGELOG.md`. Get it shipped. |

---

## 🧠 How to think (LLM Guidance)

1.  **Start here.** Always. `AGENTS.md` and `PLAN.md` give you the lay of the land.
2.  **No memory.** You have zero memory between cycles. If it's not written down, it didn't happen.
3.  **Save tokens.** Don't read the whole `docs/` folder. Only load what matches your task.
4.  **Trust the TODO.** If `TODO.md` says one thing and the docs say another, follow the TODO.
5.  **Follow the rules.** The [Atomic Commit Protocol](docs/standards/atomic-commit-protocol.md) isn't a suggestion.
6.  **Test first.** Ralph, if you write code without a failing test, you're doing it wrong.

---
*Last Updated: 2026-02-19*
