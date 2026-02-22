package orchestrator

import (
	"fmt"
	"os"
	"testing"
)

type mockAgentRunner struct {
	runs []string
}

func (m *mockAgentRunner) Run(agent string, epicID string, worktreeDir string) error {
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
	output, err := client.runTD("create", "--type", "epic", "--labels", "ready", "Implement the new orchestration system")
	if err != nil {
		t.Fatal(err)
	}
	var id string
	if _, err := fmt.Sscanf(string(output), "CREATED %s", &id); err != nil {
		t.Fatalf("failed to parse epic ID: %v", err)
	}

	agentRunner := &mockAgentRunner{}
	orch := NewOrchestrator(client, agentRunner, nil)

	// 1. Ready -> InProgress
	err = orch.Tick()
	if err != nil {
		t.Fatal(err)
	}
	// Verify our specific epic transitioned
	ids, _ := client.QueryIDs("status = in_progress")
	foundID := false
	var foundIDStr string
	for _, rid := range ids {
		if rid == id {
			foundID = true
			foundIDStr = rid
			break
		}
	}
	if !foundID {
		t.Fatalf("expected status in_progress for epic %s, but not found in %v", id, ids)
	}
	if len(agentRunner.runs) < 1 {
		t.Errorf("expected at least one run, got %d", len(agentRunner.runs))
	}
	foundRun := false
	for _, run := range agentRunner.runs {
		if run == "ralph:"+foundIDStr {
			foundRun = true
			break
		}
	}
	if !foundRun {
		t.Errorf("expected ralph to be run for epic %s, runs: %v", id, agentRunner.runs)
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
	foundBart := false
	for _, run := range agentRunner.runs {
		if run == "bart:"+id {
			foundBart = true
			break
		}
	}
	if !foundBart {
		t.Errorf("expected bart to be run for epic %s, got %v", id, agentRunner.runs)
	}
	agentRunner.runs = nil // reset

	// 3. Implemented -> Blocked (Failure)
	err = client.LogDecision(id, "bart_fail_implementation")
	if err != nil {
		t.Fatal(err)
	}
	err = orch.Tick()
	if err != nil {
		t.Fatal(err)
	}
	epic, _ = client.GetEpic(id)
	if epic.Status != "blocked" {
		t.Errorf("expected blocked status after failure, got %s", epic.Status)
	}
	foundLisa := false
	for _, run := range agentRunner.runs {
		if run == "lisa:"+id {
			foundLisa = true
			break
		}
	}
	if !foundLisa {
		t.Errorf("expected lisa to be run for epic %s, got %v", id, agentRunner.runs)
	}
	agentRunner.runs = nil // reset

	// 4. Blocked -> Ready (Lisa fixes it)
	// Lisa will set it back to ready
	err = client.Update(id, "--labels", "ready")
	if err != nil {
		t.Fatal(err)
	}
	err = orch.Tick()
	if err != nil {
		t.Fatal(err)
	}
	epic, _ = client.GetEpic(id)
	if epic.Status != "in_progress" {
		t.Errorf("expected in_progress status after Lisa fix, got %s", epic.Status)
	}
	foundRalphFix := false
	for _, run := range agentRunner.runs {
		if run == "ralph:"+id {
			foundRalphFix = true
			break
		}
	}
	if !foundRalphFix {
		t.Errorf("expected ralph to be run for epic %s after Lisa fix, got %v", id, agentRunner.runs)
	}
	agentRunner.runs = nil // reset

	// 5. InProgress -> Implemented (Retry)
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
	foundBartRetry := false
	for _, run := range agentRunner.runs {
		if run == "bart:"+id {
			foundBartRetry = true
			break
		}
	}
	if !foundBartRetry {
		t.Errorf("expected bart to be run for epic %s after retry, got %v", id, agentRunner.runs)
	}
	agentRunner.runs = nil // reset

	// 6. Implemented -> Blocked
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
	foundLisaFinal := false
	for _, run := range agentRunner.runs {
		if run == "lisa:"+id {
			foundLisaFinal = true
			break
		}
	}
	if !foundLisaFinal {
		t.Errorf("expected lisa to be run for epic %s, got %v", id, agentRunner.runs)
	}
}

func TestOrchestrator_StrictHandoff(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	tempDir, err := os.MkdirTemp("", "orchestrator-handoff-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	client := NewTDClient(tempDir)
	_, err = client.runTD("init")
	if err != nil {
		t.Fatal(err)
	}

	// Initialize git repo for worktrees
	runCmd(t, tempDir, "git", "init")
	runCmd(t, tempDir, "git", "config", "user.email", "test@example.com")
	runCmd(t, tempDir, "git", "config", "user.name", "Test User")
	runCmd(t, tempDir, "git", "commit", "--allow-empty", "-m", "initial commit")

	// Create a ready epic
	_, err = client.runTD("create", "--type", "epic", "--labels", "ready", "Implement the new handoff system")
	if err != nil {
		t.Fatal(err)
	}
	ids, _ := client.QueryIDs("status = open")
	id := ids[0]

	// Setup WorktreeManager
	wm := &WorktreeManager{BaseDir: tempDir}
	// NO handoff file created!

	agentRunner := &mockAgentRunner{}
	orch := NewOrchestrator(client, agentRunner, wm)

	// Tick should fail because handoff file is missing
	err = orch.Tick()
	if err == nil {
		t.Error("expected Tick to fail due to missing handoff file, but it succeeded")
	}

	// Verify it's still 'open' (or at least not in_progress if it failed early)
	epic, _ := client.GetEpic(id)
	if epic.Status == "in_progress" {
		t.Error("expected status not to be in_progress after failed handoff deposit")
	}
}
