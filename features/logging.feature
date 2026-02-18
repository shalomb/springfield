Feature: Structured Logging
  In order to debug and audit agent activities
  As an operator
  I want agents to emit structured JSON logs

  Scenario: Agent emits an INFO log
    Given the "ralph-1" agent is working on "EPIC-003"
    When Ralph performs a successful implementation step
    Then a new entry should appear in "logs/ralph-1.log"
    And the entry should be valid JSON
    And the "level" should be "INFO"
    And the "agent" should be "ralph-1"
    And the "epic" should be "EPIC-003"

  Scenario: Orchestrator logs session startup
    When I run "just flow"
    Then a new entry should appear in "logs/orchestrator.log"
    And the "message" should mention "session startup"
    And the "agent" should be "orchestrator"

  Scenario: Log tailing via CLI
    When I run "just logs"
    Then I should see a combined stream of JSON logs from all agents
