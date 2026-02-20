Feature: Agent Governance
  As a budget owner
  I want to control agent costs and model selection
  So that I don't exceed my financial limits

  Scenario: Session budget is enforced
    Given a project configuration with a budget of 30 tokens
    And an agent is configured with a model that uses 20 tokens per call
    When the agent attempts to perform a task that requires 2 LLM calls
    Then the first call should succeed
    And the second call should fail with a "budget exceeded" error

  Scenario: Model fallback on failure
    Given a project configuration with a primary model "primary-llm" and fallback "fallback-llm"
    And the primary model is failing
    When the agent performs a task
    Then the fallback model should be called
    And the task should complete successfully
