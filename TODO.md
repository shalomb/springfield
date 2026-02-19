# TODO - EPIC-005: Agent Governance & Selection

This TODO list covers the implementation of the Agent Governance Layer, focusing on configuration, budget enforcement, and model selection.

## Status Legend
- ⬜ To Do
- 🏃 In Progress
- ✅ Done
- ❌ Blocked

## Milestone 1: Unified Configuration Layer
Goal: Implement the core configuration loading logic and schema validation.

- [ ] **Task 1: Define Configuration Schema & ADR**
  - **Description:** Formalize the `.springfield.yaml` schema and document the governance strategy in `docs/adr/ADR-005-agent-governance.md`.
  - **Success Criteria:** `docs/adr/ADR-005-agent-governance.md` exists. Schema supports model selection and budget limits.
  - **ACP Alignment:** 1 commit: `docs(governance): define ADR-005 and configuration schema`

- [ ] **Task 2: Implement Config Loader**
  - **Description:** Implement logic to load config from project root `.springfield.yaml`.
  - **Success Criteria:** `internal/config` package can parse the yaml.
  - **ACP Alignment:** 1 commit: `feat(config): implement YAML configuration loader with TDD`

## Milestone 2: Budget Enforcement & Token Tracking
Goal: Prevent runaway costs by tracking usage and enforcing limits.

- [ ] **Task 3: Token Counting Middleware**
  - **Description:** Intercept LLM calls to count tokens.
  - **Success Criteria:** Logs contain `input_tokens` and `output_tokens`.
  - **ACP Alignment:** 1 commit: `feat(llm): add token counting middleware and audit logging`

- [ ] **Task 4: Budget Enforcer**
  - **Description:** Check cumulative usage against limits.
  - **Success Criteria:** BDD scenario for budget limit passes.
  - **ACP Alignment:** 1 commit: `feat(governance): implement budget enforcer with BDD verification`

## Milestone 3: Model Selection
Goal: Allow project-specific model routing.

- [ ] **Task 5: Model Routing**
  - **Description:** Route tasks to specific models based on config.
  - **Success Criteria:** Agents use the model specified in `.springfield.yaml`.
  - **ACP Alignment:** 1 commit: `feat(llm): implement dynamic model selection based on config`

## Milestone 4: Verification & Marge Gates
- [ ] **Task 6: Budget Reporting CLI**
  - **Description:** Implement `just budget` command.
  - **Success Criteria:** Displays current spend.
  - **ACP Alignment:** 1 commit: `feat(cli): add budget reporting command`

## Moral Compass & Compliance
- **Audit Logging**: Ensured via Task 3.
- **RBAC**: Config overrides allow for role-based model/tool access (future).
- **Safety**: Budget enforcer prevents financial "infinite loops".
