package orchestrator

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// TDClient is a client for the td(1) binary.
type TDClient struct {
	WorkDir string
}

// Issue represents a td issue.
type Issue struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Status      string   `json:"status"`
	Type        string   `json:"type"`
	Labels      []string `json:"labels"`
	Description string   `json:"description"`
	Logs        []Log    `json:"logs"`
}

// Log represents a td log entry.
type Log struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
}

// GetEpic retrieves a specific epic by ID.
func (c *TDClient) GetEpic(id string) (*Issue, error) {
	output, err := c.runTD("show", id, "--json")
	if err != nil {
		return nil, err
	}

	var issue Issue
	if err := json.Unmarshal(output, &issue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal td output: %w", err)
	}

	if issue.Type != "epic" {
		return nil, fmt.Errorf("issue %s is not an epic", id)
	}

	return &issue, nil
}

// QueryEpics queries td for epics matching the given expression.
func (c *TDClient) QueryEpics(expression string) ([]Issue, error) {
	output, err := c.runTD("query", expression, "--output", "json")
	if err != nil {
		return nil, err
	}

	var issues []Issue
	if err := json.Unmarshal(output, &issues); err != nil {
		return nil, fmt.Errorf("failed to unmarshal td output: %w", err)
	}

	// Filter only epics
	var epics []Issue
	for _, issue := range issues {
		if issue.Type == "epic" {
			epics = append(epics, issue)
		}
	}

	return epics, nil
}

// LogDecision logs a decision to an issue.
func (c *TDClient) LogDecision(id string, decision string) error {
	_, err := c.runTD("log", id, "--decision", decision)
	return err
}

func (c *TDClient) runTD(args ...string) ([]byte, error) {
	fullArgs := args
	if c.WorkDir != "" {
		fullArgs = append([]string{"-w", c.WorkDir}, args...)
	}

	cmd := exec.Command("td", fullArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("td %v failed: %w (output: %s)", args, err, string(output))
	}
	return output, nil
}
