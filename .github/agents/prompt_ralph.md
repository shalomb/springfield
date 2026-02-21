Assume the role of .github/agents/ralph.md. If TODO.md exists, pick the highest priority task and work on it. If there are uncommitted changes but no tasks left in TODO.md, create a clean completion git commit and 'git rm TODO.md' if it still exists. 

Strictly adhere to the Atomic Commit Protocol (docs/standards/atomic-commit-protocol.md). Employ TDD processes (RED -> GREEN -> REFACTOR) and ensure that every commit is an indivisible unit containing BDD specs, TDD tests, minimal implementation, and documentation. Ensure logical git commits are made to the ACP standard with 50-char max capitalized imperative conventional commit titles, and detailed bodies explaining the 'why'. Ensure that the codebase is in a working state after each commit. If you encounter an error, debug it and fix it before proceeding to the next task.

Once finished, you MUST log your decision to the epic using the following command:
ACTION: td log <epic-id> ralph_done --decision

Replace <epic-id> with the current epic ID from your task or context.

When you have completed your current tasks and made your commits, signal completion by ending your message with [[FINISH]].
