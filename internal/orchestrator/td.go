package orchestrator

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// TDClient is a client for the td(1) binary.
type TDClient struct {
	WorkDir string
}

// NewTDClient creates a new TDClient.
func NewTDClient(workDir string) *TDClient {
	if workDir == "" {
		wd, _ := os.Getwd()
		workDir = wd
	}
	return &TDClient{WorkDir: workDir}
}

// Issue represents a td issue.
type Issue struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Status      string   `json:"status"`
	Type        string   `json:"type"`
	Priority    string   `json:"priority"`
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
	if err := c.robustUnmarshal(output, &issue); err != nil {
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
		// Try single object fallback for Query too, though Query usually returns array
		var issue Issue
		if err := json.Unmarshal(output, &issue); err == nil {
			issues = []Issue{issue}
		} else {
			return nil, fmt.Errorf("failed to unmarshal td output: %w", err)
		}
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
	_, err := c.runTD("log", id, decision, "--decision")
	return err
}

// QueryIDs executes a td query and returns matching issue IDs.
func (c *TDClient) QueryIDs(expression string) ([]string, error) {
	output, err := c.runTD("query", expression, "--output", "ids")
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" || strings.Contains(trimmed, "No issues matching query") {
		return []string{}, nil
	}
	return strings.Split(trimmed, "\n"), nil
}

// Update updates one or more fields of an issue.
func (c *TDClient) Update(id string, flags ...string) error {
	args := append([]string{"update", id}, flags...)
	_, err := c.runTD(args...)
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

func (c *TDClient) robustUnmarshal(data []byte, target interface{}) error {
	// Try to unmarshal as single object first
	if err := json.Unmarshal(data, target); err == nil {
		return nil
	}

	// If that fails, try to unmarshal as array and take first element
	var list []json.RawMessage
	if err := json.Unmarshal(data, &list); err == nil {
		if len(list) == 0 {
			return fmt.Errorf("empty array")
		}
		if len(list) > 1 {
			return fmt.Errorf("multiple objects in array, expected one")
		}
		return json.Unmarshal(list[0], target)
	}

	// Return original error if both failed
	return json.Unmarshal(data, target)
}
