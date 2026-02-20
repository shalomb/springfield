package main

import (
	"context"
	"fmt"
	"os"
	"strings"

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
	Use:   "springfield",
	Short: "Springfield is an AI agent orchestration tool",
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

		// Setup dependencies
		var l llm.LLMClient
		if os.Getenv("USE_MOCK_LLM") == "true" {
			l = &testutils.MockLLM{}
		} else {
			primaryModel := cfg.Agent.PrimaryModel
			if primaryModel == "" {
				primaryModel = cfg.Agent.Model
			}
			primary := &llm.PiLLM{Model: primaryModel}

			if cfg.Agent.FallbackModel != "" {
				fallback := &llm.PiLLM{Model: cfg.Agent.FallbackModel}
				l = &llm.FallbackLLM{Primary: primary, Fallback: fallback}
			} else {
				l = primary
			}
		}
		s, err := sandbox.NewAxonSandbox(configPath)
		if err != nil {
			return fmt.Errorf("error initializing sandbox: %w", err)
		}

		a := agent.New(agentName, role, l, s)
		a.MaxIterations = cfg.Agent.MaxIterations
		a.Budget = cfg.Agent.Budget
		ctx := context.Background()

		fmt.Println("Starting agent loop...")
		if err := a.Run(ctx, task); err != nil {
			return fmt.Errorf("error in agent loop: %w", err)
		}

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
		fmt.Println(err)
		os.Exit(1)
	}
}

func runMain() error {
	return rootCmd.Execute()
}
