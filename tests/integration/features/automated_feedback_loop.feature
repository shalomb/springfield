Feature: Automated Feedback Loop
  As a developer (user)
  I want the system to automatically attempt to fix issues found during verification
  So that I don't have to manually intervene for minor errors or simple bugs

  @wip
  Scenario: Ralph fixes a bug found by Bart
    Given an Epic is in state "implemented"
    When Bart logs an "implementation failure" in td
    Then the Springfield binary should transition the Epic to "in_progress"
    And the system should trigger Ralph again

  @wip
  Scenario: Lisa defers a minor issue
    Given an Epic is in state "implemented"
    When Bart logs a "minor issue" signal in td
    Then the Springfield binary should transition the Epic to "verified"
    And Lisa should record the issue in PLAN.md under "Technical Debt"
    And the system should proceed to Release

  @wip
  Scenario: System halts on repeated failure (Circuit Breaker)
    Given an Epic is in state "implemented"
    And Ralph has already attempted to fix the same issue 2 times
    When Bart logs an "implementation failure" in td again
    Then the Springfield binary should transition the Epic to "blocked"
    And the system should exit with an error
