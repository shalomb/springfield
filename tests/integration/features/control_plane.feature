Feature: Autonomous Control Plane Signaling
  As a System Administrator
  I want the Springfield Orchestrator to manage agent lifecycle via strict signaling
  So that agents are stateless, robust, and cannot corrupt the environment

  Scenario: Successful Agent Handoff via Sentinel
    Given the Orchestrator has started "Ralph" with sentinel "auth-123"
    And the Epic "td-epic-001" is currently "in_progress"
    When the agent executes "springfield signal --sentinel auth-123 --status success"
    Then the agent process should terminate immediately
    And the Epic "td-epic-001" should transition to "implemented"
    And the Orchestrator should schedule "Bart" for the next tick

  Scenario: Unauthorized Signal Rejection
    Given the Orchestrator has started "Ralph" with sentinel "auth-123"
    When the agent executes "springfield signal --sentinel fake-999 --status success"
    Then the command should fail with "Unauthorized: Sentinel mismatch"
    And the agent process should NOT terminate
    And the Epic state should remain "in_progress"

  Scenario: Bart Rejection Triggers Correction Loop
    Given the Epic "td-epic-001" is "implemented"
    And "Bart" is running with sentinel "bart-auth"
    When the agent executes "springfield signal --sentinel bart-auth --status failed --reason 'Tests failed'"
    Then the Epic "td-epic-001" should transition to "in_progress"
    And the Orchestrator should schedule "Ralph" for the next tick
    And the previous worktree should be preserved for Ralph

  Scenario: Lisa Planning Triggers Environment Provisioning
    Given the Epic "td-epic-002" is "planned"
    And "Lisa" is running with sentinel "lisa-auth"
    When the agent executes "springfield signal --sentinel lisa-auth --status complete"
    Then the Epic "td-epic-002" should transition to "ready"
    And the Orchestrator should create a git worktree "worktrees/epic-td-epic-002"
    And the Orchestrator should schedule "Ralph" in "worktrees/epic-td-epic-002"
