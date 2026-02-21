Feature: Unified Agent Runner
  As a developer
  I want a single, parameterized agent runner implementation
  So that I can reduce technical debt and enable autonomous loops for all agents

  Scenario: Agent loads file context automatically
    Given an agent "bart" with profile:
      | context_files | ["main.go", "main_test.go"] |
    And the file "main.go" contains "package main"
    When the agent runs
    Then the LLM should receive "package main" in the context

  Scenario: Agent performs autonomous loop with tools
    Given an agent "ralph" with profile:
      | tools_enabled | ["bash", "read"] |
    And a task "fix the bug in main.go"
    When the agent runs
    Then the agent should execute "bash" commands
    And the agent should continue until "[[FINISH]]" is received

  Scenario: Agent writes output to specified target
    Given an agent "lisa" with profile:
      | output_target | "PLAN.md" |
    When the agent finishes with "The updated plan"
    Then the file "PLAN.md" should contain "The updated plan"
