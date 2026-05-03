package agent

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/shalomb/axon/pkg/types"
	"github.com/shalomb/springfield/internal/governance"
	"github.com/shalomb/springfield/internal/llm"
	"github.com/shalomb/springfield/internal/parser"
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
	Sentinel      string
	EpicID        string
}

// Agent represents an autonomous agent.
type Agent struct {
	Profile              AgentProfile
	Task                 string
	LLM                  llm.LLMClient
	Sandbox              sandbox.Sandbox
	MaxRetries           int
	MaxIterations        int
	BudgetTokens         int   // Max tokens per session (0 = unlimited)
	MaxCostNanoDollars   int64 // Max cost per session in nano-dollars (0 = unlimited)
	DailyBudgetTokens    int   // Max tokens per day (0 = unlimited)
	DailyMaxCostNano     int64 // Max cost per day in nano-dollars (0 = unlimited)
	TotalUsage           int   // Track total tokens used
	TotalCostNanoDollars int64 // Track total cost in nano-dollars
	Tracker              *governance.UsageTracker
	Sentinel             string
}

// New creates a new Agent with default settings.
func New(profile AgentProfile, l llm.LLMClient, s sandbox.Sandbox) *Agent {
	maxIterations := profile.MaxIterations
	if maxIterations == 0 {
		maxIterations = 20
	}

	// Use provided sentinel or generate session sentinel as UUIDv4
	sentinel := profile.Sentinel
	if sentinel == "" {
		// Generate UUIDv4-format sentinel (8-4-4-4-12 hex characters)
		bytes := make([]byte, 16)
		if _, err := rand.Read(bytes); err != nil {
			// Fallback if rand fails (unlikely)
			copy(bytes, []byte("failsafeguardian"))
		}
		// Set version to 4 and variant to RFC 4122
		bytes[6] = (bytes[6] & 0x0f) | 0x40
		bytes[8] = (bytes[8] & 0x3f) | 0x80
		// Format as UUIDv4: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
		sentinel = hex.EncodeToString(bytes[0:4]) + "-" +
			hex.EncodeToString(bytes[4:6]) + "-" +
			hex.EncodeToString(bytes[6:8]) + "-" +
			hex.EncodeToString(bytes[8:10]) + "-" +
			hex.EncodeToString(bytes[10:16])
	}

	// Sync sentinel to profile so Profile.Sentinel matches Agent.Sentinel
	profile.Sentinel = sentinel

	return &Agent{
		Profile:       profile,
		LLM:           l,
		Sandbox:       s,
		MaxRetries:    3,
		MaxIterations: maxIterations,
		Sentinel:      sentinel,
	}
}

func (a *Agent) log(message, level string, tokenUsage interface{}, costNanoDollars int64) {
	costDollars := float64(costNanoDollars) / 1000000000.0
	if err := logger.Log(message, level, a.Profile.Name, "", "", tokenUsage, costDollars, nil); err != nil {
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

	// Template injection for system prompt
	tmpl, err := template.New("systemPrompt").Parse(systemPrompt)
	if err != nil {
		a.log(fmt.Sprintf("Warning: Template parse error: %v. Using raw prompt.", err), "WARNING", nil, 0)
	} else {
		var buf bytes.Buffer
		data := struct {
			Name     string
			Role     string
			Sentinel string
			EpicID   string
		}{
			Name:     a.Profile.Name,
			Role:     a.Profile.Role,
			Sentinel: a.Sentinel,
			EpicID:   a.Profile.EpicID,
		}
		if err := tmpl.Execute(&buf, data); err != nil {
			a.log(fmt.Sprintf("Warning: Template execution error: %v. Using raw prompt.", err), "WARNING", nil, 0)
		} else {
			systemPrompt = buf.String()
		}
	}

	// Always ensure the agent knows its sentinel if it's set
	if a.Sentinel != "" && !strings.Contains(systemPrompt, a.Sentinel) {
		systemPrompt += fmt.Sprintf("\n\nYour session sentinel is: %s\nTo exit, use: ACTION: springfield signal --sentinel %s --status <success|failed|blocked> --reason \"...\"", a.Sentinel, a.Sentinel)
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
		cost := a.calculateCost(resp.TokenUsage)
		a.TotalCostNanoDollars += cost

		// Record usage if tracker is provided
		if a.Tracker != nil {
			if err := a.Tracker.RecordUsage(resp.TokenUsage.TotalTokens, cost); err != nil {
				a.log(fmt.Sprintf("Warning: Failed to record usage: %v", err), "WARNING", nil, 0)
			}
		}

		// Enforce session budgets
		if a.BudgetTokens > 0 && a.TotalUsage > a.BudgetTokens {
			a.log(fmt.Sprintf("Token budget exceeded: %d > %d", a.TotalUsage, a.BudgetTokens), "ERROR", nil, 0)
			return fmt.Errorf("session token budget exceeded: %d tokens used", a.TotalUsage)
		}

		if a.MaxCostNanoDollars > 0 && a.TotalCostNanoDollars > a.MaxCostNanoDollars {
			a.log(fmt.Sprintf("Cost budget exceeded: %.6f > %.6f",
				float64(a.TotalCostNanoDollars)/1000000000.0,
				float64(a.MaxCostNanoDollars)/1000000000.0), "ERROR", nil, 0)
			return fmt.Errorf("session cost budget exceeded: $%.6f used", float64(a.TotalCostNanoDollars)/1000000000.0)
		}

		// Enforce daily budgets if tracker is provided
		if a.Tracker != nil {
			daily, err := a.Tracker.GetDailyUsage()
			if err == nil {
				if a.DailyBudgetTokens > 0 && daily.TotalTokens > a.DailyBudgetTokens {
					a.log(fmt.Sprintf("Daily token budget exceeded: %d > %d", daily.TotalTokens, a.DailyBudgetTokens), "ERROR", nil, 0)
					return fmt.Errorf("daily token budget exceeded: %d tokens used today", daily.TotalTokens)
				}
				if a.DailyMaxCostNano > 0 && daily.TotalCostNano > a.DailyMaxCostNano {
					a.log(fmt.Sprintf("Daily cost budget exceeded: %.6f > %.6f",
						float64(daily.TotalCostNano)/1000000000.0,
						float64(a.DailyMaxCostNano)/1000000000.0), "ERROR", nil, 0)
					return fmt.Errorf("daily cost budget exceeded: $%.6f used today", float64(daily.TotalCostNano)/1000000000.0)
				}
			}
		}

		a.log(fmt.Sprintf("LLM response: %s", resp.Content), "DEBUG", resp.TokenUsage, cost)

		// Extract thought if present
		thought := extractThought(resp.Content)
		if thought != "" {
			a.log(fmt.Sprintf("Thought: %s", thought), "INFO", nil, 0)
		}

		// Check for promise (v0.6.0 standard)
		promise, pErr := parser.ExtractPromise(resp.Content)
		if pErr == nil {
			if promise == parser.PromiseComplete {
				a.log("Detected promise: COMPLETE. Task complete.", "INFO", nil, 0)
				return a.finish(resp.Content, thought)
			}
			if promise == parser.PromiseFailed {
				a.log("Agent promised failure. Stopping.", "ERROR", nil, 0)
				return fmt.Errorf("agent promised failure")
			}
		} else {
			a.log(fmt.Sprintf("Warning: Promise extraction error: %v", pErr), "WARNING", nil, 0)
		}

		// Fallback to legacy finish marker detection
		if a.isFinished(resp.Content) {
			a.log("Detected legacy finish marker. Task complete.", "INFO", nil, 0)
			return a.finish(resp.Content, thought)
		}

		// If we reach here and no promise was found, we might want to warn
		if pErr == nil && promise == parser.PromiseUnknown {
			a.log("No promise found in response. Continuing loop.", "WARNING", nil, 0)
		}

		messages = append(messages, llm.Message{Role: "assistant", Content: resp.Content})

		// Improved action extraction
		action := extractAction(resp.Content)
		if action != "" {
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

			// Intercept signal command *after* sandbox execution to decide if we exit
			if isSignalAction(action) && result.ExitCode == 0 {
				sentinel, ok := extractSentinel(action)
				if ok && sentinel == a.Sentinel {
					a.log("Intercepted successful signal, exiting loop.", "INFO", nil, 0)
					return a.finish(resp.Content, thought)
				}
			}

			resultStr := fmt.Sprintf("STDOUT: %s\nSTDERR: %s\nEXIT CODE: %d", result.Stdout, result.Stderr, result.ExitCode)
			if ctxInfo := formatContext(result.Context); ctxInfo != "" {
				resultStr += "\nSANDBOX CONTEXT: " + ctxInfo
			}
			a.log(fmt.Sprintf("Action result: %s", resultStr), "DEBUG", nil, 0)
			messages = append(messages, llm.Message{Role: "user", Content: resultStr})
		} else {
			a.log("No action or finish detected.", "WARNING", nil, 0)
		}
	}

	return fmt.Errorf("max iterations reached")
}

func (a *Agent) finish(content, thought string) error {
	// Persist output if target is specified
	if a.Profile.OutputTarget != "" {
		// Clean up tags from final response if present
		cleanContent := content
		if thought != "" {
			cleanContent = strings.Replace(cleanContent, fmt.Sprintf("<thought>%s</thought>", thought), "", 1)
		}
		// We also want to strip thought tags entirely if they are in different formatting
		cleanContent = thoughtTagRegex.ReplaceAllString(cleanContent, "")
		// Strip promise tags too
		cleanContent = promiseTagRegex.ReplaceAllString(cleanContent, "")

		if err := a.persistOutput(cleanContent); err != nil {
			a.log(fmt.Sprintf("Error persisting output to %s: %v", a.Profile.OutputTarget, err), "ERROR", nil, 0)
			return err
		}
	}
	return nil
}

func (a *Agent) calculateCost(usage llm.TokenUsage) int64 {
	if usage.CostNanoDollars > 0 {
		return usage.CostNanoDollars
	}

	// Simple fallback cost calculation using nano-dollars for precision ($1.00 = 1,000,000,000 nano-dollars)
	// Prices per token in nano-dollars:
	const promptPricePerTokenNano = 75      // $0.075 / 1M = 75 nano-dollars per token
	const completionPricePerTokenNano = 300 // $0.30 / 1M = 300 nano-dollars per token

	return int64(usage.PromptTokens)*promptPricePerTokenNano + int64(usage.CompletionTokens)*completionPricePerTokenNano
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
	trimmed := strings.TrimSpace(resp)
	finished := strings.HasSuffix(trimmed, marker)
	return finished
}

func formatContext(c types.ContextMetadata) string {
	if c.ProjectType == "" && c.BuildTool == "" && c.GitStatus == "" {
		return ""
	}
	return fmt.Sprintf("Project type: %s, Build tool: %s, Test framework: %s, Git status: %s",
		c.ProjectType, c.BuildTool, c.TestFramework, c.GitStatus)
}

var actionRegex = regexp.MustCompile(`(?m)^ACTION:\s*(.+)$`)
var actionTagRegex = regexp.MustCompile(`(?s)<action>\s*(.*?)\s*</action>`)
var thoughtTagRegex = regexp.MustCompile(`(?s)<thought>\s*(.*?)\s*</thought>`)
var promiseTagRegex = regexp.MustCompile(`(?si)<promise>(.*?)</promise>`)

func extractAction(resp string) string {
	s := parser.NewMarkdownSanitizer()
	sanitized := s.StripCodeBlocks(resp)

	// Try tag-based extraction first (more robust)
	match := actionTagRegex.FindStringSubmatch(sanitized)
	if len(match) >= 2 {
		return strings.TrimSpace(match[1])
	}

	// Fallback to legacy ACTION: prefix
	match = actionRegex.FindStringSubmatch(sanitized)
	if len(match) >= 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func extractThought(resp string) string {
	s := parser.NewMarkdownSanitizer()
	sanitized := s.StripCodeBlocks(resp)

	match := thoughtTagRegex.FindStringSubmatch(sanitized)
	if len(match) >= 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
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

func isSignalAction(action string) bool {
	// Must be a signal command at the start AND must not contain newlines
	// (signals are single-line commands; multiline indicates a malformed action)
	if !strings.HasPrefix(action, "springfield signal") {
		return false
	}
	return !strings.Contains(action, "\n")
}

func extractSentinel(action string) (string, bool) {
	// Extract sentinel only if it matches UUIDv4 format (8-4-4-4-12 hex chars)
	// This prevents injection of arbitrary values like flags, paths, or command substitutions
	re := regexp.MustCompile(`--sentinel\s+([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})`)
	matches := re.FindStringSubmatch(action)
	if len(matches) < 2 {
		return "", false
	}
	return matches[1], true
}
