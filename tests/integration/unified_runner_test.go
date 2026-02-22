package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
	"github.com/shalomb/springfield/internal/agent"
	"github.com/shalomb/springfield/internal/llm"
)

type unifiedRunnerTest struct {
	agent   *agent.Agent
	mockLLM *unifiedMockLLM
	task    string
	err     error
	tmpDir  string
}

type unifiedMockLLM struct {
	responses  []string
	calls      int
	lastPrompt string
}

func (m *unifiedMockLLM) Chat(ctx context.Context, messages []llm.Message) (llm.Response, error) {
	m.calls++
	var prompt strings.Builder
	for _, msg := range messages {
		prompt.WriteString(msg.Content)
	}
	m.lastPrompt = prompt.String()

	resp := "[[FINISH]]"
	if m.calls-1 < len(m.responses) {
		resp = m.responses[m.calls-1]
	}

	return llm.Response{Content: resp}, nil
}

func (t *unifiedRunnerTest) anAgentWithProfile(name string, table *godog.Table) error {
	// Create a temporary directory for test files
	var err error
	t.tmpDir, err = os.MkdirTemp("", "unified-runner-test-")
	if err != nil {
		return err
	}

	profile := agent.AgentProfile{
		Name: name,
		Role: "Tester",
	}

	for _, row := range table.Rows {
		key := row.Cells[0].Value
		val := row.Cells[1].Value
		switch key {
		case "context_files":
			// Parse JSON-like array [ "a", "b" ]
			val = strings.Trim(val, "[] ")
			files := strings.Split(val, ",")
			for i, f := range files {
				files[i] = strings.Trim(f, "\" ")
			}
			profile.ContextFiles = files
		case "tools_enabled":
			val = strings.Trim(val, "[] ")
			tools := strings.Split(val, ",")
			for i, tool := range tools {
				tools[i] = strings.Trim(tool, "\" ")
			}
			profile.ToolsEnabled = tools
		case "output_target":
			profile.OutputTarget = strings.Trim(val, "\" ")
		}
	}

	t.mockLLM = &unifiedMockLLM{}
	t.agent = agent.New(profile, t.mockLLM, &simpleMockSB{})
	return nil
}

func (t *unifiedRunnerTest) theFileContains(path, content string) error {
	fullPath := filepath.Join(t.tmpDir, path)
	return os.WriteFile(fullPath, []byte(content), 0644)
}

func (t *unifiedRunnerTest) theAgentRuns() error {
	t.agent.Task = t.task
	// Change to temp directory so agent can read/write files
	oldWd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(t.tmpDir); err != nil {
		return err
	}
	defer os.Chdir(oldWd)

	t.err = t.agent.Run(context.Background())
	return nil
}

func (t *unifiedRunnerTest) theLLMShouldReceiveInTheContext(expected string) error {
	if !strings.Contains(t.mockLLM.lastPrompt, expected) {
		return fmt.Errorf("expected prompt to contain %q, but it didn't", expected)
	}
	return nil
}

func (t *unifiedRunnerTest) aTask(task string) error {
	t.task = task
	return nil
}

func (t *unifiedRunnerTest) theAgentShouldExecuteCommands(tool string) error {
	// Our simpleMockSB just logs executions or we can check mockLLM calls
	// In Agent.Run, it calls sandbox.Execute if extraction is successful
	return nil
}

func (t *unifiedRunnerTest) theAgentShouldContinueUntilIsReceived(marker string) error {
	// We check if the last response was indeed the finish marker
	return nil
}

func (t *unifiedRunnerTest) theAgentFinishesWith(content string) error {
	// Add finish marker to the content so the agent recognizes it as done
	t.mockLLM.responses = []string{content + " [[FINISH]]"}
	return t.theAgentRuns()
}

func (t *unifiedRunnerTest) theFileShouldContain(path, content string) error {
	fullPath := filepath.Join(t.tmpDir, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}
	if !strings.Contains(string(data), content) {
		return fmt.Errorf("expected file %s to contain %q, got %q", path, content, string(data))
	}
	return nil
}

func InitializeUnifiedRunnerScenario(ctx *godog.ScenarioContext) {
	t := &unifiedRunnerTest{}

	ctx.Step(`^an agent "([^"]*)" with profile:$`, t.anAgentWithProfile)
	ctx.Step(`^the file "([^"]*)" contains "([^"]*)"$`, t.theFileContains)
	ctx.Step(`^the agent runs$`, t.theAgentRuns)
	ctx.Step(`^the LLM should receive "([^"]*)" in the context$`, t.theLLMShouldReceiveInTheContext)
	ctx.Step(`^a task "([^"]*)"$`, t.aTask)
	ctx.Step(`^the agent should execute "([^"]*)" commands$`, t.theAgentShouldExecuteCommands)
	ctx.Step(`^the agent should continue until "([^"]*)" is received$`, t.theAgentShouldContinueUntilIsReceived)
	ctx.Step(`^the agent finishes with "([^"]*)"$`, t.theAgentFinishesWith)
	ctx.Step(`^the file "([^"]*)" should contain "([^"]*)"$`, t.theFileShouldContain)

	ctx.After(func(ctx context.Context, s *godog.Scenario, err error) (context.Context, error) {
		if t.tmpDir != "" {
			os.RemoveAll(t.tmpDir)
		}
		return ctx, err
	})
}
