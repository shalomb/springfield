Feature: Agent Sandboxing
  In order to safely execute autonomous agent code
  As a system administrator
  I want agent execution to be isolated from the host system

  @requires_axon
  Scenario: Agent cannot read host files outside workspace (Host File Isolation)
    Given a sandbox environment is configured
    And a file "/tmp/springfield_host_only.txt" exists on the host
    When I attempt to run "cat /tmp/springfield_host_only.txt" within the sandbox
    Then the operation should fail

  @requires_axon
  Scenario: Agent cannot write to host system directories
    Given a sandbox environment is configured
    When I attempt to run "touch /usr/bin/malicious_tool" within the sandbox
    Then the operation should fail

  @requires_axon
  Scenario: Agent can read and write within the workspace
    Given a sandbox environment is configured
    And a file "sandbox_test.txt" exists in the workspace
    When I run the command "cat sandbox_test.txt" within the sandbox
    Then the output should contain "Initial content"
    And I run the command "echo 'updated in sandbox' > sandbox_test.txt" within the sandbox
    And the file "sandbox_test.txt" in the host workspace should contain "updated in sandbox"

  @requires_axon
  Scenario: Agent resource usage is limited
    Given a sandbox environment is configured
    When I run a memory-intensive process in the sandbox that exceeds 512MB
    Then the process should be killed or restricted
