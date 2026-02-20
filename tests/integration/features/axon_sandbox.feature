Feature: Agent Sandboxing via Axon
  In order to run tasks safely and consistently
  As a developer
  I want agents to execute in isolated environments via Axon

  @requires_axon
  Scenario: Execute a valid command in the sandbox
    Given a sandbox environment is configured
    When I run the command "echo 'Hello Sandbox'" within the sandbox
    Then the command should succeed
    And the output should contain "Hello Sandbox"
    And the execution metadata should be populated (duration > 0)

  @requires_axon @wip
  Scenario: Non-zero exit code is captured correctly
    # WIP: Axon currently exits 0 regardless of inner command exit code.
    # Requires Axon pkg/executor work (ADR-004 Phase 3d).
    Given a sandbox environment is configured
    When I run the command "exit 42" within the sandbox
    Then the exit code should be 42

  @requires_axon @wip
  Scenario: Stderr is captured separately from stdout
    # WIP: Axon merges TTY warning into stderr channel, masking user stderr.
    # Requires Axon to suppress internal warnings in --format json mode.
    Given a sandbox environment is configured
    When I run the command "echo 'err' >&2" within the sandbox
    Then the stderr should contain "err"

  @requires_axon @wip
  Scenario: Prevent host file access outside workspace
    # WIP: Axon returns exit 0 even when ls fails inside the container.
    # Requires exit code passthrough fix (ADR-004 Phase 3d).
    Given a sandbox environment is configured
    And the host has a directory "/etc"
    When I attempt to run "ls /etc/shadow" within the sandbox
    Then the operation should fail

  @requires_axon @wip
  Scenario: Preserve workspace state within sandbox
    # WIP: Volume mount write-back not yet verified end-to-end from Springfield.
    Given a sandbox environment is configured with a workspace volume
    And a file "test.txt" exists in the workspace
    When I run the command "echo 'Updated' > test.txt" within the sandbox
    Then the file "test.txt" in the host workspace should contain "Updated"
    And the changes should persist after the sandbox execution

  @requires_axon @wip
  Scenario: Prevent unauthorized access to root (guardrail)
    # WIP: Axon guardrails block based on instruction text patterns, not
    # runtime filesystem enforcement. "ls /root" isn't matched by current
    # threat patterns so it passes. Requires guardrail tuning in Axon.
    Given a sandbox environment is configured
    When I attempt to run "ls /root" within the sandbox
    Then the operation should fail
    And the output should contain "BLOCKED"

  @requires_axon @wip
  Scenario: Context intelligence detects project type
    Given a sandbox environment is configured
    When I run the command "pwd" within the sandbox
    Then the execution context should identify project_type as "go"
    And the execution context should identify build_tool as "go"

  @requires_axon @wip
  Scenario: Tool discovery finds common devops tools
    Given a sandbox environment is configured
    When I run the command "which git" within the sandbox
    Then the command should succeed
    And the tools list should contain "git"
    And the tools list should contain "go"
