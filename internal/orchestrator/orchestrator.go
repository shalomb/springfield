package orchestrator

import (
	"fmt"
	"log"
)

// AgentRunner provides an interface for running agents.
type AgentRunner interface {
	Run(agent string, epicID string) error
}

// Orchestrator manages the execution of Epics.
type Orchestrator struct {
	TD    *TDClient
	Agent AgentRunner
}

// NewOrchestrator creates a new Orchestrator.
func NewOrchestrator(td *TDClient, agent AgentRunner) *Orchestrator {
	return &Orchestrator{TD: td, Agent: agent}
}

// CommandAgentRunner runs agents by executing the springfield binary.
type CommandAgentRunner struct {
	BinaryPath string
}

func (r *CommandAgentRunner) Run(agent string, epicID string) error {
	// For now, we just log. In a real implementation, we would execute
	// the springfield binary with the appropriate flags.
	log.Printf("INVOKING AGENT: %s for Epic %s", agent, epicID)
	return nil
}

// Tick performs one iteration of the orchestration loop.
func (o *Orchestrator) Tick() error {
	// 1. Find Epics that might need processing
	ids, err := o.TD.QueryIDs("type = epic AND status != closed")
	if err != nil {
		return fmt.Errorf("failed to query epics: %w", err)
	}

	for _, id := range ids {
		log.Printf("Processing Epic %s", id)
		if err := o.processEpic(id); err != nil {
			log.Printf("Error processing Epic %s: %v", id, err)
		}
	}

	return nil
}

func (o *Orchestrator) processEpic(id string) error {
	epic, err := o.TD.GetEpic(id)
	if err != nil {
		return err
	}

	state := o.determineState(epic)
	
	log.Printf("Epic %s is in state %s", id, state)

	switch state {
	case StatusReady:
		log.Printf("Transitioning Epic %s to in_progress", id)
		// Update td status to in_progress and remove 'ready' label
		if err := o.TD.Update(id, "--status", "in_progress", "--labels", ""); err != nil {
			return err
		}
		if o.Agent != nil {
			return o.Agent.Run("ralph", id)
		}
		return nil
	case StatusInProgress:
		// Check for completion signals from Ralph
		if o.hasDecision(epic, "ralph_done") {
			log.Printf("Ralph complete for Epic %s. Transitioning to implemented.", id)
			// td status: in_review, label: implemented
			if err := o.TD.Update(id, "--status", "in_review", "--labels", "implemented"); err != nil {
				return err
			}
			if o.Agent != nil {
				return o.Agent.Run("bart", id)
			}
			return nil
		}
	}

	return nil
}

func (o *Orchestrator) hasDecision(epic *Issue, decision string) bool {
	for _, l := range epic.Logs {
		if l.Type == "decision" && l.Message == decision {
			return true
		}
	}
	return false
}

func (o *Orchestrator) determineState(epic *Issue) EpicStatus {
	// Simple mapping for now. ADR-008 states Lisa updates td Epic status.
	// Since td only supports a few, we use labels to supplement.
	
	for _, label := range epic.Labels {
		switch label {
		case "ready":
			return StatusReady
		case "implemented":
			return StatusImplemented
		case "verified":
			return StatusVerified
		}
	}

	switch epic.Status {
	case "open":
		return StatusPlanned
	case "in_progress":
		return StatusInProgress
	case "blocked":
		return StatusBlocked
	case "closed":
		return StatusDone
	}

	return StatusPlanned
}
