package orchestrator

import (
	"fmt"
	"log"
)

// MockTDClient is a mock implementation of TDClient for testing.
type MockTDClient struct {
	Issues []Issue
	logs   map[string][]Log
}

// QueryIDs returns all issue IDs matching a query.
func (m *MockTDClient) QueryIDs(query string) ([]string, error) {
	var ids []string
	for _, issue := range m.Issues {
		// Simple query matching for testing
		if query == "type = epic AND status != closed" {
			if issue.Type == "epic" && issue.Status != "closed" {
				ids = append(ids, issue.ID)
			}
		}
	}
	return ids, nil
}

// GetEpic retrieves an epic by ID.
func (m *MockTDClient) GetEpic(id string) (*Issue, error) {
	for i, issue := range m.Issues {
		if issue.ID == id && issue.Type == "epic" {
			// Attach logs to the issue
			issue.Logs = m.logs[id]
			m.Issues[i] = issue
			return &m.Issues[i], nil
		}
	}
	return nil, fmt.Errorf("epic not found: %s", id)
}

// GetIssue retrieves an issue by ID.
func (m *MockTDClient) GetIssue(id string) (*Issue, error) {
	for i, issue := range m.Issues {
		if issue.ID == id {
			// Attach logs to the issue
			issue.Logs = m.logs[id]
			m.Issues[i] = issue
			return &m.Issues[i], nil
		}
	}
	return nil, fmt.Errorf("issue not found: %s", id)
}

// Update updates an issue.
func (m *MockTDClient) Update(id string, args ...string) error {
	for i, issue := range m.Issues {
		if issue.ID == id {
			// Simple mock: parse args for status and labels
			for j := 0; j < len(args); j += 2 {
				if j+1 < len(args) {
					switch args[j] {
					case "--status":
						issue.Status = args[j+1]
					case "--labels":
						if args[j+1] == "" {
							issue.Labels = []string{}
						} else {
							issue.Labels = []string{args[j+1]}
						}
					}
				}
			}
			m.Issues[i] = issue
			return nil
		}
	}
	return fmt.Errorf("issue not found: %s", id)
}

// LogDecision logs a decision.
func (m *MockTDClient) LogDecision(id string, decision string) error {
	if m.logs == nil {
		m.logs = make(map[string][]Log)
	}
	m.logs[id] = append(m.logs[id], Log{Type: "decision", Message: decision})
	return nil
}

// AddLog adds a log entry for testing.
func (m *MockTDClient) AddLog(id string, message string) {
	if m.logs == nil {
		m.logs = make(map[string][]Log)
	}
	m.logs[id] = append(m.logs[id], Log{Type: "log", Message: message})
}

// MockAgentRunner is a mock implementation of AgentRunner for testing.
type MockAgentRunner struct {
	calls map[string]int
}

// Run executes an agent.
func (m *MockAgentRunner) Run(agent string, epicID string, worktreeDir string) error {
	log.Printf("MockAgentRunner: running %s for epic %s", agent, epicID)
	if m.calls == nil {
		m.calls = make(map[string]int)
	}
	m.calls[agent]++
	return nil
}

// WasCalled checks if an agent was called.
func (m *MockAgentRunner) WasCalled(agent string) bool {
	if m.calls == nil {
		return false
	}
	return m.calls[agent] > 0
}
