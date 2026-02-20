Feature: Agent Execution Loop
  In order to solve complex tasks
  As a user
  I want the agent to think and act until the task is complete

  Scenario: Agent successfully completes a single-action task
    Given an agent named "test-agent" with role "tester"
    And the LLM is configured with responses:
      | thought                    | action | finish |
      | I need to list files       | ls     | false  |
      | I see the files. I am done |        | true   |
    When the agent runs the task "What's in the current directory?"
    Then the agent should have called the LLM 2 times
    And the agent should have executed 1 action in the sandbox
    And the task should be successful

  Scenario: Agent finishes immediately without any action
    Given an agent named "test-agent" with role "tester"
    And the LLM is configured with responses:
      | thought         | action | finish |
      | Nothing to do.  |        | true   |
    When the agent runs the task "Do nothing"
    Then the agent should have called the LLM 1 times
    And the agent should have executed 0 actions in the sandbox
    And the task should be successful

  Scenario: Agent performs multiple sequential actions
    Given an agent named "test-agent" with role "tester"
    And the LLM is configured with responses:
      | thought         | action   | finish |
      | Step one        | echo a   | false  |
      | Step two        | echo b   | false  |
      | All done        |          | true   |
    When the agent runs the task "Do two things"
    Then the agent should have called the LLM 3 times
    And the agent should have executed 2 actions in the sandbox
    And the task should be successful

  Scenario: Agent receives sandbox output in subsequent LLM prompt
    Given an agent named "test-agent" with role "tester"
    And the LLM is configured with responses:
      | thought              | action      | finish |
      | Run the check        | echo hello  | false  |
      | Saw the output, done |             | true   |
    When the agent runs the task "Check output"
    Then the task should be successful
    And the second LLM call should include the sandbox output

  Scenario: Agent system prompt includes name and role
    Given an agent named "Marge" with role "Product Agent"
    And the LLM is configured with responses:
      | thought | action | finish |
      | Done    |        | true   |
    When the agent runs the task "anything"
    Then the system prompt should contain "Marge"
    And the system prompt should contain "Product Agent"
