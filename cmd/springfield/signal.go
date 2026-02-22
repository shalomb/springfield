package main

import (
	"fmt"
	"os"

	"github.com/shalomb/springfield/internal/orchestrator"
	"github.com/spf13/cobra"
)

var (
	sentinelFlag string
	statusFlag   string
	reasonFlag   string
	epicIDFlag   string
)

var signalCmd = &cobra.Command{
	Use:   "signal",
	Short: "Signal completion or status change for the current session",
	Long:  `Agents use this command to terminate their session and signal a state transition.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Validation
		if sentinelFlag == "" {
			return fmt.Errorf("missing required flag: --sentinel")
		}
		if statusFlag == "" {
			return fmt.Errorf("missing required flag: --status")
		}

		// Validate status enum
		validStatuses := map[string]bool{
			"complete":    true, // Lisa
			"success":     true, // Ralph, Bart
			"failed":      true, // Bart
			"blocked":     true, // Any
			"released":    true, // Lovejoy
			"implemented": true, // Ralph (alias for success)
			"verified":    true, // Bart (alias for success)
		}

		if !validStatuses[statusFlag] {
			return fmt.Errorf("invalid status '%s'. Allowed: complete, success, failed, blocked, released, implemented, verified", statusFlag)
		}

		// 2. Sentinel Verification (Inside Sandbox)
		// The orchestrator sets this env var when spawning the agent
		expectedSentinel := os.Getenv("SPRINGFIELD_SENTINEL")
		if expectedSentinel == "" {
			// If not set, we might be running locally/manually. Warn but allow?
			// For rigorous control plane, we should fail or warn loudly.
			// Let's print to stderr so it doesn't pollute stdout capture if we care.
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Warning: SPRINGFIELD_SENTINEL env var not set. Skipping internal validation.\n")
		} else if sentinelFlag != expectedSentinel {
			return fmt.Errorf("unauthorized: sentinel token mismatch")
		}

		// 3. Resolve Epic ID
		epicID := epicIDFlag
		if epicID == "" {
			epicID = os.Getenv("SPRINGFIELD_EPIC_ID")
		}

		if epicID != "" {
			cmd.Printf("Signaling status: %s for Epic %s\n", statusFlag, epicID)

			// 4. Update td (if available)
			// We create a client for the current directory
			td := orchestrator.NewTDClient("")

			reason := reasonFlag
			if reason == "" {
				reason = fmt.Sprintf("Agent signaled %s", statusFlag)
			}

			// Log decision
			if logErr := td.LogDecision(epicID, fmt.Sprintf("%s: %s", statusFlag, reason)); logErr != nil {
				cmd.Printf("Warning: Failed to log decision to td: %v\n", logErr)
			}

			// Update Status based on Signal
			var nextStatus string
			var labels []string

			agentRole := os.Getenv("SPRINGFIELD_AGENT")

			switch statusFlag {
			case "complete":
				nextStatus = "ready"
			case "success", "implemented", "verified":
				// Context-dependent mapping
				if agentRole == "ralph" {
					nextStatus = "in_review" // Mapped to 'implemented' via label in Orchestrator logic usually
					labels = []string{"implemented"}
				} else if agentRole == "bart" {
					// Bart success -> Verified
					labels = []string{"verified"}
				}
			case "released":
				nextStatus = "done"
			case "failed":
				nextStatus = "in_progress"
			case "blocked":
				nextStatus = "blocked"
			}

			if nextStatus != "" || len(labels) > 0 {
				// We rely on Orchestrator Daemon to do the heavy lifting of state transitions.
				cmd.Printf("Signal maps to status '%s' (transition handled by orchestrator)\n", nextStatus)
			}
		}

		if reasonFlag != "" {
			cmd.Printf("Reason: %s\n", reasonFlag)
		}

		// The specific string "SIGNAL_SENT" is useful for the orchestrator to parse if needed,
		// but exit code 0 is the primary signal.
		cmd.Println("SIGNAL_SENT")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(signalCmd)
	signalCmd.Flags().StringVar(&sentinelFlag, "sentinel", "", "Session sentinel token")
	signalCmd.Flags().StringVar(&statusFlag, "status", "", "Outcome status (complete, success, failed, blocked, released)")
	signalCmd.Flags().StringVar(&reasonFlag, "reason", "", "Description of the result")
	signalCmd.Flags().StringVar(&epicIDFlag, "epic", "", "Epic ID associated with this session")
}
