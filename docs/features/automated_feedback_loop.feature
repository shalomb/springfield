Feature: Automated Feedback Loop
  As a developer (user)
  I want the system to automatically attempt to fix issues found during verification
  So that I don't have to manually intervene for minor errors or simple bugs

  Scenario: Ralph fixes a bug found by Bart
    Given a feature branch "feat/epic-test" exists
    And Ralph has implemented a change
    But Bart finds a failure in "test-unit" and updates FEEDBACK.md
    When Lisa analyzes FEEDBACK.md
    Then she should identify the failure as "Fixable"
    And she should add a "Fix Bug" task to TODO.md
    And the system should trigger Ralph again

  Scenario: Bart defers a minor issue for later
    Given a feature branch "feat/epic-test" exists
    And Ralph has implemented a change
    But Bart finds a minor code style issue and updates FEEDBACK.md
    When Lisa analyzes FEEDBACK.md
    Then she should identify the issue as "Minor"
    And she should add a note to PLAN.md under "Technical Debt"
    And she should clear FEEDBACK.md
    And the system should proceed to Release

  Scenario: System halts on repeated failure (Circuit Breaker)
    Given a feature branch "feat/epic-test" exists
    And Ralph has attempted to fix the same issue 2 times
    But Bart still finds the same failure
    When Lisa analyzes FEEDBACK.md
    Then she should identify the failure as "Persistent"
    And she should NOT trigger Ralph again
    And the system should exit with an error
