Feature: Git Branching Model
  In order to maintain a consistent and reliable codebase
  As a developer
  I want a standard for naming and managing branches

  Scenario: Create a Feature Branch
    Given I am on the "main" branch
    And the working directory is clean
    When I run "just start-feature 'my-feature'"
    Then a new branch "feat/my-feature" should be created
    And I should be switched to "feat/my-feature"

  Scenario: Create a Fix Branch
    Given I am on the "main" branch
    And the working directory is clean
    When I run "just start-fix 'my-bug'"
    Then a new branch "fix/my-bug" should be created
    And I should be switched to "fix/my-bug"

  Scenario: Enforce Naming Convention
    Given I am on the "main" branch
    When I run "just start-feature 'Invalid Name'"
    Then the command should fail
    And the error message should mention "lowercase-kebab-case"
