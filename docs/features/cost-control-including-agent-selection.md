# Feature.md - Cost Control and Agent Selection

## Problem
Currently, there is no standardized way to track or control the cost of running agents. Using expensive LLM models without budget controls or selection criteria can lead to unexpected expenses.

## Requirements
- Track token usage and estimate cost for each agent execution.
- Implement cost limits or budgets per session, task, or user.
- Provide a mechanism to select the most cost-effective agent for a given task (e.g., smaller models for simple tasks).

## Acceptance Criteria
- [ ] Centralized token counting and cost estimation logic.
- [ ] Configurable budget limits.
- [ ] Agent selection strategy based on task complexity and cost.
- [ ] Integration with `just` commands to report current usage/costs.

## Constraints & Unknowns
- **Constraint:** Must accurately track usage from various LLM providers (OpenAI, Anthropic, etc.).
- **Unknown:** Handling of rate limits and API errors in cost calculation.

## Options Considered
- [ ] Manual tracking: Error-prone, hard to maintain.
- [ ] Third-party service (e.g., Helicone): Powerful, but external dependency.
- [ ] Built-in tracking: Requires custom implementation, but integrated.

## Scope
✅ Token counting implementation
✅ Cost estimation logic
✅ Budget configuration file/ENV variables
❌ Full billing integration (future)

## Success Criteria
- Ability to see current session cost.
- Warning or stop when budget exceeded.
- Reduced overall token usage by selecting appropriate agents.
