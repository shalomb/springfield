Feature: Tmux Agent Orchestration
  In order to manage multiple agents concurrently
  As a developer
  I want a tmux-based orchestration layer

  Scenario: Adopt Existing Tmux Session
    Given I am already inside a tmux session
    When I run "just flow"
    Then no new tmux session should be created
    And new windows should be added to the current session for "marge", "lisa", "ralph", "bart", and "lovejoy"

  Scenario: Create New Session when Out-of-Session
    Given I am not inside a tmux session
    When I run "just flow"
    Then a new tmux session should be created
    And the session should be named after the current directory
    And it should contain the core agent windows

  Scenario: Agent Logs are Visible
    Given the "ralph" agent is running
    When I switch to the "ralph" window
    Then I should see a pane tailing "logs/ralph.log"
