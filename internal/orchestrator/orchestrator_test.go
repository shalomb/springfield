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
	_, err = client.runTD("create", "--type", "epic", "--labels", "ready", "Implement the new orchestration system")
	if err != nil {
		t.Fatal(err)
	}

	agentRunner := &mockAgentRunner{}
	orch := NewOrchestrator(client, agentRunner, nil)

	// 1. Ready -> InProgress
	err = orch.Tick()
	if err != nil {
		t.Fatal(err)
	}
	ids, _ := client.QueryIDs("status = in_progress")
	if len(ids) != 1 {
		t.Fatalf("expected 1 in_progress epic, got %d", len(ids))
	}
	id := ids[0]
	if len(agentRunner.runs) != 1 || agentRunner.runs[0] != "ralph:"+id {
		t.Errorf("expected ralph to be run for epic %s, got %v", id, agentRunner.runs)
	}
	agentRunner.runs = nil // reset

	// 2. InProgress -> Implemented
	err = client.LogDecision(id, "ralph_done")
	if err != nil {
		t.Fatal(err)
	}
	err = orch.Tick()
	if err != nil {
		t.Fatal(err)
	}
	epic, _ := client.GetEpic(id)
	if epic.Status != "in_review" {
		t.Errorf("expected in_review status, got %s", epic.Status)
	}
	if len(agentRunner.runs) != 1 || agentRunner.runs[0] != "bart:"+id {
		t.Errorf("expected bart to be run for epic %s, got %v", id, agentRunner.runs)
	}
	agentRunner.runs = nil // reset

	// 3. Implemented -> InProgress (Failure)
	err = client.LogDecision(id, "bart_fail_implementation")
	if err != nil {
		t.Fatal(err)
	}
	err = orch.Tick()
	if err != nil {
		t.Fatal(err)
	}
	epic, _ = client.GetEpic(id)
	if epic.Status != "in_progress" {
		t.Errorf("expected in_progress status after failure, got %s", epic.Status)
	}
	if len(agentRunner.runs) != 1 || agentRunner.runs[0] != "ralph:"+id {
		t.Errorf("expected ralph to be run for epic %s, got %v", id, agentRunner.runs)
	}
	agentRunner.runs = nil // reset

	// 4. InProgress -> Implemented (Retry)
	err = client.LogDecision(id, "ralph_done")
	if err != nil {
		t.Fatal(err)
	}
	err = orch.Tick()
	if err != nil {
		t.Fatal(err)
	}
	epic, _ = client.GetEpic(id)
	if epic.Status != "in_review" {
		t.Errorf("expected in_review status after retry, got %s", epic.Status)
	}
	if len(agentRunner.runs) != 1 || agentRunner.runs[0] != "bart:"+id {
		t.Errorf("expected bart to be run for epic %s, got %v", id, agentRunner.runs)
	}
	agentRunner.runs = nil // reset

	// 5. Implemented -> Blocked
	err = client.LogDecision(id, "bart_fail_viability")
	if err != nil {
		t.Fatal(err)
	}
	err = orch.Tick()
	if err != nil {
		t.Fatal(err)
	}
	epic, _ = client.GetEpic(id)
	if epic.Status != "blocked" {
		t.Errorf("expected blocked status, got %s", epic.Status)
	}
	if len(agentRunner.runs) != 1 || agentRunner.runs[0] != "lisa:"+id {
		t.Errorf("expected lisa to be run for epic %s, got %v", id, agentRunner.runs)
	}
}

func TestOrchestrator_Planned(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "orchestrator-test-planned")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	client := NewTDClient(tempDir)
	_, err = client.runTD("init")
	if err != nil {
		t.Fatal(err)
	}

	// Create a planned epic (open, no ready label)
	_, err = client.runTD("create", "--type", "epic", "Implement new feature")
	if err != nil {
		t.Fatal(err)
	}

	agentRunner := &mockAgentRunner{}
	orch := NewOrchestrator(client, agentRunner, nil)

	// Planned -> Lisa
	err = orch.Tick()
	if err != nil {
		t.Fatal(err)
	}
	ids, _ := client.QueryIDs("status = open")
	if len(ids) != 1 {
		t.Fatalf("expected 1 open epic, got %d", len(ids))
	}
	id := ids[0]
	if len(agentRunner.runs) != 1 || agentRunner.runs[0] != "lisa:"+id {
		t.Errorf("expected lisa to be run for epic %s, got %v", id, agentRunner.runs)
	}
}

