package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/shalomb/springfield/internal/agent"
	"github.com/shalomb/springfield/internal/config"
	"github.com/shalomb/springfield/internal/llm"
	"github.com/shalomb/springfield/internal/sandbox"
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

		role := "Assistant"
		switch strings.ToLower(agentName) {
		case "marge":
			role = "Product Agent"
		case "lisa":
			role = "Planning Agent"
		case "ralph":
			role = "Build Agent"
		case "bart":
			role = "Quality Agent"
		case "lovejoy":
			role = "Release Agent"
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
			l = &mockLLM{}
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
		return nil
	},
}

type mockLLM struct{}

func (m *mockLLM) Chat(ctx context.Context, messages []llm.Message) (llm.Response, error) {
	if os.Getenv("MOCK_LLM_ERROR") == "true" {
		return llm.Response{}, fmt.Errorf("mock llm error")
	}
	// Very simple mock response to allow the loop to finish in tests
	return llm.Response{
		Content: "THOUGHT: I am a mock agent. [[FINISH]]",
		TokenUsage: llm.TokenUsage{
			PromptTokens:     10,
			CompletionTokens: 10,
			TotalTokens:      20,
		},
	}, nil
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
