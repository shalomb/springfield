# PLAN.md - Feature Roadmap

## Feature: Git Branching Model

### Epic 1: Define and Implement Branching Strategy
- [ ] Task 1: Create Feature Brief for Git Branching Model
  - Status: unstarted
  - Assignee: product-discovery
  - Feature Brief: docs/features/git-branching-model.md
- [ ] Task 2: Define ADR for Branching Strategy
  - Status: unstarted
  - Assignee: frink
  - Depends on: Task 1
- [ ] Task 3: Configure Repository Protection Rules (if applicable/simulated)
  - Status: unstarted
  - Assignee: ralph

## Feature: Compatibility with existing agent and sub-agent definitions or skills

e.g. the `.{github,claude,gemini}/{agents,skills}/` directory structures in an existing repository
The springfield agent harness should be able to read and execute agents and
skills defined in these existing directory structures, allowing for seamless
integration with existing repositories and workflows. The repository
definitions take precedence over any agents or skills defined in the
springfield agents or skills directories, allowing for easy overrides and
customizations and allowing for an opt-in of springfield opinionated workflows
and structures.

## Feature: CLI Tooling for Agent Management

- Entry point CLI to launch agents, view logs, manage sessions, etc.
- Corrolaries with `just invoke $agent` and `just utilize $skill` commands for ease of use and discoverability. Allowing for partial matching of agent and skill names for ease of use (e.g. `just invoke review` to invoke an agent named `code-review-agent`).
- There is likely to be a justfile clash/conflict with an existing one in a repository, So we probably want to use an entirely different CLI tool name that is not `just`. What names can we generate that are catchy, easy to remember, and related to the theme of springfield or orchestration?
- We could use a justfile behind the scenes such that the CLI tool is just a thin wrapper around the justfile in the springfield installation, allowing us to leverage the power and flexibility of just while providing a more user-friendly and thematic CLI interface and avoiding the conflict.

## Feature: Tmux Concurrent Multi-Agent Orchestration

### Epic 1: Enable Parallel Agent Execution
- [ ] Task 1: Create Feature Brief for Tmux Orchestration
  - Status: unstarted
  - Assignee: product-discovery
  - Feature Brief: docs/features/tmux-concurrent-multi-agent-orchestration.md
- [ ] Task 2: Prototype Tmux Session Management Script
  - Status: unstarted
  - Assignee: ralph
- [ ] Task 3: Implement Agent Harness integration with Tmux
  - Status: unstarted
  - Assignee: ralph

## Feature: Logging and Observability

### Epic 1: Centralized Logging
- [ ] Task 1: Create Feature Brief for Logging
  - Status: unstarted
  - Assignee: product-discovery
  - Feature Brief: docs/features/logging-and-observability.md
- [ ] Task 2: Define Logging Standard (ADR)
  - Status: unstarted
  - Assignee: frink
- [ ] Task 3: Implement Structured Logger
  - Status: unstarted
  - Assignee: ralph

## Feature: Sandboxing and Agent Execution Context

### Epic 1: Secure Agent Environment
- [ ] Task 1: Create Feature Brief for Sandboxing
  - Status: unstarted
  - Assignee: product-discovery
  - Feature Brief: docs/features/sandboxing-and-agent-execution-context.md
- [ ] Task 2: Research and Select Sandboxing Technology (Docker/chroot/etc)
  - Status: unstarted
  - Assignee: frink
  - ADR: docs/adr/ADR-001-sandboxing.md

## Feature: Cost Control

### Epic 1: Token and Resource Management
- [ ] Task 1: Create Feature Brief for Cost Control
  - Status: unstarted
  - Assignee: product-discovery
  - Feature Brief: docs/features/cost-control-including-agent-selection.md
- [ ] Task 2: Implement Token Counter Middleware
  - Status: unstarted
  - Assignee: ralph

## Feature: Integrations with Existing Project Management Tools

- Integrations with tools like Jira, Github Issues, Apptio TargetProcess, etc to pull in tasks, user stories, or project requirements that agents can then work on. This would allow for a more seamless integration of the agent orchestration system into existing workflows and project management practices, enabling teams to easily leverage the power of agents to automate and assist with their work without having to manually input tasks or requirements into the system.

