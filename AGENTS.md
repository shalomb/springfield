# Springfield Agent Command Center

Welcome to the team.

This is your map. You don't need to read the entire documentation tree to do your job—in fact, please don't. Use this file to find exactly what you need, get the context, and get to work.

---

## 🗺️ The System

### 📍 Core
- **[README.md](README.md)**: What this project is.
- **[quick-start.md](docs/tutorials/quick-start.md)**: How to get moving fast.
- **[PLAN.md](PLAN.md)**: The roadmap. What we're building right now.
- **`td`**: The command-line task and state manager.

### 📚 The Docs (Diataxis)
We organize documentation by purpose, not just topic.
- **[Concepts](docs/concepts/)**: The "Why." Architecture and mental models.
- **[How-To](docs/how-to/)**: The "How." Step-by-step guides.
- **[Reference](docs/reference/)**: The "What." Specs, glossaries, and hard data.
- **[Standards](docs/standards/)**: The "Rules." How we code and commit.
- **[ADRs](docs/adr/)**: The "Decisions." Why we did it this way.

### 🤖 The Crew
- **[.pi/agents/](.pi/agents/)**: Who you are.
- **[.pi/skills/](.pi/skills/)**: What you can do.

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
7.  **No Littering.** Do NOT write status reports, analysis logs, or "thinking" files to the repository. Use `/tmp/` for scratchpads. If a finding is critical, summarize it in `PLAN.md` or `td` logs. Only code, tests, and permanent documentation belong in git.

---

## ⚠️ Important Technical Notes

### Git Commits Must Be Signed  
- All commits must be GPG-signed: `git commit -S`
- Signing is enforced by pre-commit hooks
- **Known Issue:** In non-TTY environments (like CI/CD), GPG signing may fail with "gpg: cannot open '/dev/tty'" 
- **Workaround:** Temporarily disable with `git config commit.gpgSign false` if needed, then re-enable

### Interactive Git Operations Require Non-Editor Execution
- **Problem:** `git rebase -i` and similar commands launch `vim` by default, which will hang
- **Solution:** Use environment variable to avoid editor: `GIT_SEQUENCE_EDITOR=true git rebase -i HEAD~N`
- **Alternative:** Use `git reset --soft` and recommit for safer operations without vim

---
*Last Updated: 2026-02-22*
