package orchestrator

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestTDClient_QueryEpics(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "td-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize td in temp dir
	cmdInit := exec.Command("td", "-w", tempDir, "init")
	if output, err := cmdInit.CombinedOutput(); err != nil {
		t.Fatalf("failed to init td: %v (output: %s)", err, string(output))
	}

	// Initialize td in temp dir by creating an epic
	// Using a long title to satisfy td's requirements
	cmd := exec.Command("td", "-w", tempDir, "create", "--type", "epic", "Implement the new orchestration system")
	outputCreate, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to create epic: %v (output: %s)", err, string(outputCreate))
	}
	// Extract ID from output like "CREATED td-773d3d"
	var id string
	if _, err := fmt.Sscanf(string(outputCreate), "CREATED %s", &id); err != nil {
		t.Fatalf("failed to parse epic ID from output: %v (output: %s)", err, string(outputCreate))
	}
	t.Logf("Created ID: %s", id)

	client := &TDClient{WorkDir: tempDir}
	epics, err := client.QueryEpics("type = epic")
	if err != nil {
		t.Fatalf("QueryEpics failed: %v", err)
	}

	foundEpic := false
	for _, e := range epics {
		if e.ID == id {
			foundEpic = true
			if e.Title != "Implement the new orchestration system" {
				t.Errorf("expected title 'Implement the new orchestration system', got '%s'", e.Title)
			}
			break
		}
	}
	if !foundEpic {
		t.Errorf("created epic %s not found in query results", id)
	}

	// Test GetEpic
	epic, err := client.GetEpic(id)
	if err != nil {
		t.Fatalf("GetEpic failed: %v", err)
	}
	if epic.ID != id {
		t.Errorf("expected ID %s, got %s", id, epic.ID)
	}

	// Test LogDecision and logs retrieval
	err = client.LogDecision(id, "bart_ok")
	if err != nil {
		t.Fatalf("LogDecision failed: %v", err)
	}

	epicWithLogs, err := client.GetEpic(id)
	if err != nil {
		t.Fatalf("GetEpic with logs failed: %v", err)
	}
	found := false
	for _, l := range epicWithLogs.Logs {
		if l.Type == "decision" && l.Message == "bart_ok" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("decision log 'bart_ok' not found in epic logs")
	}

	// Test QueryIDs
	ids, err := client.QueryIDs("type = epic")
	if err != nil {
		t.Fatalf("QueryIDs failed: %v", err)
	}
	foundID := false
	for _, rid := range ids {
		if rid == id {
			foundID = true
			break
		}
	}
	if !foundID {
		t.Errorf("expected ID %s to be in %v", id, ids)
	}

	// Test Update
	err = client.Update(id, "--priority", "P2")
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	updatedEpic, err := client.GetEpic(id)
	if err != nil {
		t.Fatalf("GetEpic after update failed: %v", err)
	}
	if updatedEpic.Priority != "P2" {
		t.Errorf("expected priority P2, got %s", updatedEpic.Priority)
	}
}
