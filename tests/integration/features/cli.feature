Feature: CLI
  As a developer
  I want a command-line interface for Springfield
  So that I can run agents and tasks easily

  Scenario: Run an agent with a task
    Given a task "say hello"
    When I run "springfield --agent ralph --task 'say hello'"
    Then the agent "ralph" should receive the task "say hello"
    And the CLI output should contain "Agent: ralph"
    And the CLI output should contain "Task: say hello"
