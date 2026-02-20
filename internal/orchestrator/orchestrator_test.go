package orchestrator

import (
	"os"
	"testing"
)

type mockAgentRunner struct {
	runs []string
}

func (m *mockAgentRunner) Run(agent string, epicID string) error {
	m.runs = append(m.runs, agent+":"+epicID)
	return nil
}

func TestOrchestrator_Tick(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "orchestrator-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	client := NewTDClient(tempDir)
	// Initialize td
	_, err = client.runTD("init")
	if err != nil {
		t.Fatal(err)
	}

	// Create a ready epic
	// Since we can't easily use labels with the current TDClient.Update (it replaces),
	// and we don't have a CreateEpic method yet, let's just use runTD.
	_, err = client.runTD("create", "--type", "epic", "--labels", "ready", "Implement the new orchestration system")
	if err != nil {
		t.Fatal(err)
	}

	agentRunner := &mockAgentRunner{}
	orch := NewOrchestrator(client, agentRunner, nil)
	err = orch.Tick()
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}

	// Verify that the epic transitioned to in_progress
	ids, err := client.QueryIDs("status = in_progress")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 {
		t.Errorf("expected 1 in_progress epic, got %d", len(ids))
	}
	id := ids[0]

	if len(agentRunner.runs) != 1 || agentRunner.runs[0] != "ralph:"+id {
		t.Errorf("expected ralph run, got %v", agentRunner.runs)
	}

	// Now log Ralph's completion
	err = client.LogDecision(id, "ralph_done")
	if err != nil {
		t.Fatal(err)
	}

	// Run Tick again
	err = orch.Tick()
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}

	// Verify that the epic transitioned to implemented (status in_review)
	ids, err = client.QueryIDs("status = in_review AND labels ~ implemented")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 {
		t.Errorf("expected 1 implemented epic, got %d", len(ids))
	}

	if len(agentRunner.runs) != 2 || agentRunner.runs[1] != "bart:"+id {
		t.Errorf("expected bart run, got %v", agentRunner.runs)
	}
}
