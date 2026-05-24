Feature: Springfield Setup Commands
  As a developer adopting Springfield
  I want init and doctor commands that work correctly
  So that I can get the tool running without reading all the docs

  Scenario: init in a fresh git repository
    Given a fresh git repository
    When I run setup command "springfield init"
    Then the setup command should succeed
    And the td database should exist at ".todos/issues.db"
    And setup file ".gitignore" should contain ".todos/"
    And setup file ".gitignore" should contain "worktrees/"
    And setup file ".springfield.toml" should exist
    And setup file ".pi/agents" should exist
    And setup file ".todos" should exist
    And setup file "Justfile" should exist
    And setup file ".codemap/config.json" should exist
    And the setup output should contain "Done."

  Scenario: init is idempotent
    Given a fresh git repository
    And "springfield init" has already been run
    When I run setup command "springfield init"
    Then the setup command should succeed
    And the setup output should contain "already initialised"
    And the setup output should contain "already up to date"

  Scenario: init outside a git repository
    Given a directory that is not a git repository
    When I run setup command "springfield init"
    Then the setup command should fail
    And the setup output should contain "not a git repository"

  Scenario: doctor passes after init
    Given a fresh git repository
    And "springfield init" has already been run
    When I run setup command "springfield doctor"
    Then the setup command should succeed
    And the setup output should contain "[ok]"
    And the setup output should contain "td binary"
    And the setup output should contain "td database"
    And the setup output should contain "springfield config"

  Scenario: doctor reports missing td database
    Given a fresh git repository
    When I run setup command "springfield doctor"
    Then the setup output should contain "[fail]"
    And the setup output should contain "td database"

  Scenario: doctor reports missing springfield config
    Given a fresh git repository
    And the td database has been initialised
    When I run setup command "springfield doctor"
    Then the setup output should contain "[warn]"
    And the setup output should contain "springfield config"
