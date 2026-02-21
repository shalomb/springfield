package agent

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/shalomb/axon/pkg/types"
	"github.com/shalomb/springfield/internal/llm"
	"github.com/shalomb/springfield/internal/sandbox"
	"github.com/shalomb/springfield/pkg/logger"
)

const FinishMarker = "[[FINISH]]"

// Runner defines the interface for agent runners.
type Runner interface {
	Run(ctx context.Context) error
}

// AgentProfile defines the personality and behavior constraints for an agent.
type AgentProfile struct {
	Name          string
	Role          string
	SystemPrompt  string
	ContextFiles  []string
	OutputTarget  string
	ToolsEnabled  []string
	FinishMarker  string
	MaxIterations int
}

// Agent represents an autonomous agent.
type Agent struct {
	Profile       AgentProfile
	Task          string
	LLM           llm.LLMClient
	Sandbox       sandbox.Sandbox
	MaxRetries    int
	MaxIterations int
	Budget        int // Max tokens per session (0 = unlimited)
	TotalUsage    int // Track total tokens used
}

// New creates a new Agent with default settings.
func New(profile AgentProfile, l llm.LLMClient, s sandbox.Sandbox) *Agent {
	maxIterations := profile.MaxIterations
	if maxIterations == 0 {
		maxIterations = 20
	}
	return &Agent{
		Profile:       profile,
		LLM:           l,
		Sandbox:       s,
		MaxRetries:    3,
		MaxIterations: maxIterations,
	}
}

func (a *Agent) log(message, level string, tokenUsage interface{}, cost float64) {
	if err := logger.Log(message, level, a.Profile.Name, "", "", tokenUsage, cost, nil); err != nil {
		fmt.Fprintf(os.Stderr, "CRITICAL: Logger failed: %v\nMessage was: %s\n", err, message)
	}
}

// Run executes the agent's task.
// It implements the Runner interface.
func (a *Agent) Run(ctx context.Context) error {
	task := a.Task
	a.log(fmt.Sprintf("Starting task: %s", task), "INFO", nil, 0)

	systemPrompt := a.Profile.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = fmt.Sprintf("You are %s, a %s.", a.Profile.Name, a.Profile.Role)
	}

	messages := []llm.Message{
		{Role: "system", Content: systemPrompt},
	}

	// Load context files if specified
	if len(a.Profile.ContextFiles) > 0 {
		fileContext := a.loadFilesContext()
		if fileContext != "" {
			messages = append(messages, llm.Message{Role: "user", Content: fileContext})
		}
	}

	messages = append(messages, llm.Message{Role: "user", Content: task})

	for iteration := 0; iteration < a.MaxIterations; iteration++ {
		var resp llm.Response
		var err error
		for i := 0; i <= a.MaxRetries; i++ {
			resp, err = a.LLM.Chat(ctx, messages)
			if err == nil {
				break
			}
			a.log(fmt.Sprintf("LLM error (attempt %d/%d): %v", i+1, a.MaxRetries+1, err), "WARNING", nil, 0)
			if i == a.MaxRetries {
				a.log("Max retries reached for LLM call.", "ERROR", nil, 0)
				return err
			}
		}

		a.TotalUsage += resp.TokenUsage.TotalTokens
		if a.Budget > 0 && a.TotalUsage > a.Budget {
			a.log(fmt.Sprintf("Budget exceeded: %d > %d", a.TotalUsage, a.Budget), "ERROR", nil, 0)
			return fmt.Errorf("session budget exceeded: %d tokens used", a.TotalUsage)
		}

		cost := a.calculateCost(resp.TokenUsage)
		a.log(fmt.Sprintf("LLM response: %s", resp.Content), "DEBUG", resp.TokenUsage, cost)
		messages = append(messages, llm.Message{Role: "assistant", Content: resp.Content})

		if a.isFinished(resp.Content) {
			a.log("Task complete.", "INFO", nil, 0)

			// Persist output if target is specified
			if a.Profile.OutputTarget != "" {
				if err := a.persistOutput(resp.Content); err != nil {
					a.log(fmt.Sprintf("Error persisting output to %s: %v", a.Profile.OutputTarget, err), "ERROR", nil, 0)
					return err
				}
			}
			return nil
		}

		// Very basic action extraction
		if strings.Contains(resp.Content, "ACTION:") {
			action := extractAction(resp.Content)
			if action == "" {
				a.log("Extracted action is empty.", "WARNING", nil, 0)
				continue
			}

			if isUnsafeAction(action) {
				a.log(fmt.Sprintf("Blocked unsafe action: %s", action), "ERROR", nil, 0)
				messages = append(messages, llm.Message{Role: "user", Content: "Action blocked for security reasons."})
				continue
			}

			a.log(fmt.Sprintf("Executing action: %s", action), "INFO", nil, 0)
			var result *types.Result
			for i := 0; i <= a.MaxRetries; i++ {
				result, err = a.Sandbox.Execute(ctx, action)
				if err == nil {
					break
				}
				a.log(fmt.Sprintf("Sandbox error (attempt %d/%d): %v", i+1, a.MaxRetries+1, err), "WARNING", nil, 0)
				if i == a.MaxRetries {
					a.log("Max retries reached for Sandbox execution.", "ERROR", nil, 0)
					return err
				}
			}

			resultStr := fmt.Sprintf("STDOUT: %s\nSTDERR: %s\nEXIT CODE: %d", result.Stdout, result.Stderr, result.ExitCode)
			if ctxInfo := formatContext(result.Context); ctxInfo != "" {
				resultStr += "\nSANDBOX CONTEXT: " + ctxInfo
			}
			a.log(fmt.Sprintf("Action result: %s", resultStr), "DEBUG", nil, 0)
			messages = append(messages, llm.Message{Role: "user", Content: resultStr})
		} else {
			// If no action and no finish, we might be stuck or just talking
			// For now, let's just continue to the next loop
			a.log("No action or finish detected.", "WARNING", nil, 0)
		}
	}

	return fmt.Errorf("max iterations reached")
}

func (a *Agent) calculateCost(usage llm.TokenUsage) float64 {
	// Simple cost calculation. In the future this should be based on the model from config.
	// Prices per 1M tokens.
	const promptPrice = 0.075 / 1000000.0    // $0.075 / 1M
	const completionPrice = 0.30 / 1000000.0 // $0.30 / 1M
	return float64(usage.PromptTokens)*promptPrice + float64(usage.CompletionTokens)*completionPrice
}

func (a *Agent) persistOutput(content string) error {
	// Strip finish marker
	marker := a.Profile.FinishMarker
	if marker == "" {
		marker = FinishMarker
	}
	content = strings.Replace(content, marker, "", -1)
	content = strings.TrimSpace(content)

	a.log(fmt.Sprintf("Persisting output to %s", a.Profile.OutputTarget), "INFO", nil, 0)
	return os.WriteFile(a.Profile.OutputTarget, []byte(content), 0644)
}

func (a *Agent) loadFilesContext() string {
	var parts []string
	for _, file := range a.Profile.ContextFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			a.log(fmt.Sprintf("Warning: Could not read context file %s: %v", file, err), "WARNING", nil, 0)
			continue
		}
		parts = append(parts, fmt.Sprintf("FILE: %s\nCONTENT:\n%s\n---", file, string(content)))
	}
	if len(parts) == 0 {
		return ""
	}
	return "CURRENT CONTEXT FILES:\n\n" + strings.Join(parts, "\n\n")
}

func (a *Agent) isFinished(resp string) bool {
	marker := a.Profile.FinishMarker
	if marker == "" {
		marker = FinishMarker
	}
	return strings.HasSuffix(strings.TrimSpace(resp), marker)
}

func formatContext(c types.ContextMetadata) string {
	if c.ProjectType == "" && c.BuildTool == "" && c.GitStatus == "" {
		return ""
	}
	return fmt.Sprintf("Project type: %s, Build tool: %s, Test framework: %s, Git status: %s",
		c.ProjectType, c.BuildTool, c.TestFramework, c.GitStatus)
}

var actionRegex = regexp.MustCompile(`(?m)^ACTION:\s*(.+)$`)

func extractAction(resp string) string {
	match := actionRegex.FindStringSubmatch(resp)
	if len(match) < 2 {
		return ""
	}
	return strings.TrimSpace(match[1])
}

func isUnsafeAction(action string) bool {
	// Block obviously malicious sequential commands that try to escape the intended action.
	// We allow pipes '|' and logical AND '&&' as they are common in legitimate agent tasks.

	// Still block:
	// ;  - command separator (can be used to run arbitrary secondary commands)
	// `  - backticks (command substitution)
	// $() - command substitution
	// || - logical OR (can be used for malicious branching)

	blockedPatterns := []string{
		";", "`", "$(", "||",
	}

	for _, p := range blockedPatterns {
		if strings.Contains(action, p) {
			return true
		}
	}

	return false
}
