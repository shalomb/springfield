Feature: Agent Promise Semantic Contract
  As an orchestrator
  I want agents to explicitly promise completion or failure
  To ensure reliable state transitions

  Scenario: Agent completes with promise
    Given an agent "ralph" with profile configured for promises
    When the agent returns "<promise>COMPLETE</promise>"
    Then the agent loop should terminate successfully

  Scenario: Agent fails with promise
    Given an agent "ralph" with profile configured for promises
    When the agent returns "<promise>FAILED</promise>"
    Then the agent loop should terminate with error "agent promised failure"

  Scenario: Agent backward compatible with [[FINISH]]
    Given an agent "ralph" with profile configured for legacy finish
    When the agent returns "[[FINISH]]"
    Then the agent loop should terminate successfully

  Scenario: Agent ignores promise in code block
    Given an agent "ralph" with profile configured for promises
    When the agent returns:
      """
      I am still working.
      ```
      <promise>COMPLETE</promise>
      ```
      <action>ls</action>
      """
    Then the agent loop should not terminate
    And the action "ls" should be extracted
