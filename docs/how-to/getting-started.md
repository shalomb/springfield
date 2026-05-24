# Getting Started

So you've joined the team. Welcome to Springfield.

This guide assumes you have the repository cloned and your environment set up. If not, ask your lead where the coffee machine is.

## 0. Prerequisites: `td` (Task Daemon)

Springfield requires `td` (Task Daemon) to run. It is our decentralized task management system that agents use to read instructions and report progress.

- **What it is:** A lightweight CLI for managing `TODO.md` state.
- **Why it's required:** The orchestration loops (like `just do`) poll `td` to know what to build next and when to stop.
- **How to install it:** Run `go install github.com/shalomb/td@latest` (requires Go) or install from the source repository at `https://github.com/shalomb/td`.

## 1. The Interface: `just`

We use `just` to run everything. It's our command center.

```bash
just help       # Show all commands
just plan       # Plan tasks with Lisa
just do         # Run the full agent loop
```

## 2. Your First Feature

Ready to build something? Here is the lifecycle of a task.

### Step 1: Start a Branch
Don't work on main. We aren't savages.
```bash
just start-feature my-cool-feature
```

### Step 2: Define the Work (Marge & Lisa)
Before writing code, we need a plan.
1.  Create or update **Feature.md** to define the "Why" and "What."
2.  Run the planner to generate the task list.
```bash
just lisa "Analyze Feature.md and create a TODO.md for the implementation"
```

### Step 3: The Loop (Ralph & Bart)
Now we build. You can run the autonomous loop or drive it manually.

**The "Do It All" Command:**
```bash
just do
```
*This runs Lisa (Plan) → Ralph (Build) → Bart (Quality Review & Verify) in a cycle.*

**Or Manual Mode:**
```bash
just ralph    # Build the next task in TODO.md
just test     # Run the test ladder
```

### Step 4: Ship It (Lovejoy)
When `TODO.md` is empty and `FEEDBACK.md` is clean:
```bash
just lovejoy
```

## 3. The Rules of Engagement

1.  **Trust the Tests:** If `just test` fails, you are not done.
2.  **Atomic Commits:** Ralph will yell at you if you make giant commits. Keep them small.
3.  **Read the Docs:** If you get stuck, check `docs/reference/loops.md` to find the right tool for the job.

## 4. Troubleshooting

*   **"Agent is hallucinating":** Kill the process. Check `TODO.md` for ambiguous instructions. Clarify them. Retry.
*   **"Tests failing":** Run `just test-unit` to narrow it down.
*   **"I don't know what to do":** Run `just plan` and let Lisa tell you.
