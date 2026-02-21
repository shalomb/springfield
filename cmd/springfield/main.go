package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shalomb/springfield/internal/agent"
	"github.com/shalomb/springfield/internal/config"
	"github.com/shalomb/springfield/internal/llm"
	"github.com/shalomb/springfield/internal/orchestrator"
	"github.com/shalomb/springfield/internal/sandbox"
	"github.com/shalomb/springfield/internal/testutils"
	"github.com/spf13/cobra"
)

var (
	agentName  string
	task       string
	configPath string
)

var rootCmd = &cobra.Command{
	Use:           "springfield",
	Short:         "Springfield is an AI agent orchestration tool",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if agentName == "" || task == "" {
			return cmd.Help()
		}

		roles := map[string]string{
			"marge":   "Product Agent",
			"lisa":    "Planning Agent",
			"ralph":   "Build Agent",
			"bart":    "Quality Agent",
			"lovejoy": "Release Agent",
		}

		role, ok := roles[strings.ToLower(agentName)]
		if !ok {
			role = "Assistant"
		}

		fmt.Printf("Agent: %s (%s)\n", agentName, role)
		fmt.Printf("Task: %s\n", task)

		// Load config
		cfg, err := config.LoadConfig(".")
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		// Get agent-specific config (falls back to defaults if not configured)
		agentCfg := cfg.GetAgentConfig(agentName)

		// Setup dependencies
		var l llm.LLMClient
		if os.Getenv("USE_MOCK_LLM") == "true" {
			l = &testutils.MockLLM{}
		} else {
			primaryModel := agentCfg.PrimaryModel
			if primaryModel == "" {
				primaryModel = agentCfg.Model
			}
			primary := &llm.PiLLM{Model: primaryModel}

			if agentCfg.FallbackModel != "" {
				fallback := &llm.PiLLM{Model: agentCfg.FallbackModel}
				l = &llm.FallbackLLM{Primary: primary, Fallback: fallback}
			} else {
				l = primary
			}
		}
		// Initialize sandbox
		sandboxInst, err := sandbox.NewAxonSandbox(configPath)
		if err != nil {
			return fmt.Errorf("error initializing sandbox: %w", err)
		}

		// Use a 60-second timeout for agent execution
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Create a specialized runner based on the agent type, with budget and sandbox
		runner, err := agent.NewRunnerWithBudget(agentName, task, l, sandboxInst, agentCfg.Budget)
		if err != nil {
			return fmt.Errorf("error creating runner for agent %s: %w", agentName, err)
		}

		fmt.Println("Starting agent loop...")
		if err := runner.Run(ctx); err != nil {
			// Check for quota errors (terminal conditions)
			if llm.IsQuotaExceededError(err) {
				fmt.Fprintf(os.Stderr, "\n🛑 CRITICAL: API QUOTA EXCEEDED\n")
				fmt.Fprintf(os.Stderr, "   %s\n", err.Error())
				fmt.Fprintf(os.Stderr, "\n⚠️  Execution halted to preserve uncommitted changes.\n")
				fmt.Fprintf(os.Stderr, "   Please resolve the quota issue and try again.\n\n")
				return fmt.Errorf("quota exceeded - execution halted")
			}

			// Format other error messages more clearly
			errMsg := fmt.Sprintf("%v", err)
			fmt.Fprintf(os.Stderr, "❌ Error: %s\n", errMsg)
			return fmt.Errorf("error in agent loop: %w", err)
		}

		fmt.Println("✅ Agent completed successfully")
		return nil
	},
}

var orchestrateCmd = &cobra.Command{
	Use:   "orchestrate",
	Short: "Run the orchestration loop",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Orchestration loop starting...")
		tdClient := orchestrator.NewTDClient("")
		worktreeManager := &orchestrator.WorktreeManager{BaseDir: "."}
		agentRunner := &orchestrator.CommandAgentRunner{BinaryPath: os.Args[0]}
		orch := orchestrator.NewOrchestrator(tdClient, agentRunner, worktreeManager)

		return orch.Tick()
	},
}

func init() {
	rootCmd.AddCommand(orchestrateCmd)
	rootCmd.Flags().StringVarP(&agentName, "agent", "a", "", "Name of the agent (marge/lisa/ralph/bart/lovejoy)")
	rootCmd.Flags().StringVarP(&task, "task", "t", "", "Task to execute")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to axon config.toml")
}

func main() {
	if err := runMain(); err != nil {
		// Don't print error here - it's already printed in the RunE function
		os.Exit(1)
	}
}

func runMain() error {
	return rootCmd.Execute()
}
