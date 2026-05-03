package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents the severity of a check result.
type Status int

const (
	OK   Status = iota
	Warn        // present but something is off
	Fail        // missing or broken
)

// Check is a single diagnostic result.
type Check struct {
	Name    string
	Status  Status
	Message string
}

func (c Check) Prefix() string {
	switch c.Status {
	case OK:
		return "[ok]  "
	case Warn:
		return "[warn]"
	case Fail:
		return "[fail]"
	default:
		return "[?]   "
	}
}

// RunChecks executes all diagnostic checks from the given repo root and
// returns the results. root must be an absolute path to the git repo root.
func RunChecks(root string) []Check {
	var checks []Check

	checks = append(checks, checkTD(root))
	checks = append(checks, checkTDDB(root))
	checks = append(checks, checkConfig(root))
	checks = append(checks, checkGitignore(root))
	checks = append(checks, checkAgentCLIs()...)

	return checks
}

func checkTD(root string) Check {
	path, err := exec.LookPath("td")
	if err != nil {
		return Check{
			Name:    "td binary",
			Status:  Fail,
			Message: "td not found on PATH — install from https://github.com/shalomb/td",
		}
	}

	out, err := exec.Command("td", "--version").Output()
	version := strings.TrimSpace(string(out))
	if err != nil || version == "" {
		return Check{Name: "td binary", Status: Warn, Message: fmt.Sprintf("found at %s but --version failed", path)}
	}

	return Check{Name: "td binary", Status: OK, Message: fmt.Sprintf("%s (%s)", path, version)}
}

func checkTDDB(root string) Check {
	db := filepath.Join(root, ".todos", "issues.db")
	if _, err := os.Stat(db); os.IsNotExist(err) {
		return Check{
			Name:    "td database",
			Status:  Fail,
			Message: fmt.Sprintf("not initialised — run: td init -w %s", root),
		}
	}
	return Check{Name: "td database", Status: OK, Message: db}
}

func checkConfig(root string) Check {
	for _, name := range []string{".springfield.toml", "config.toml"} {
		p := filepath.Join(root, name)
		if _, err := os.Stat(p); err == nil {
			return Check{Name: "springfield config", Status: OK, Message: p}
		}
	}
	return Check{
		Name:    "springfield config",
		Status:  Warn,
		Message: "no .springfield.toml found — run: springfield init",
	}
}

func checkGitignore(root string) Check {
	gi := filepath.Join(root, ".gitignore")
	data, err := os.ReadFile(gi)
	if err != nil {
		return Check{Name: ".gitignore", Status: Warn, Message: ".gitignore not found"}
	}

	missing := []string{}
	for _, entry := range []string{".todos/", "worktrees/"} {
		if !strings.Contains(string(data), entry) {
			missing = append(missing, entry)
		}
	}
	if len(missing) > 0 {
		return Check{
			Name:    ".gitignore",
			Status:  Warn,
			Message: fmt.Sprintf("missing entries: %s", strings.Join(missing, ", ")),
		}
	}
	return Check{Name: ".gitignore", Status: OK, Message: gi}
}

// knownCLIs is the ordered list of agent CLIs Springfield probes for.
// pi is the preferred default; others are alternatives.
var knownCLIs = []struct {
	bin      string
	label    string
	provider string // default provider when this CLI is chosen
}{
	{"pi", "pi (mariozechner/pi-coding-agent)", "anthropic"},
	{"claude", "claude (Anthropic Claude Code)", "anthropic"},
	{"gemini", "gemini (Google Gemini CLI)", "google-gemini-cli"},
	{"gh copilot", "gh copilot (GitHub Copilot)", "github-copilot"},
	{"kiro", "kiro (Amazon Kiro)", "openai"},
}

// AgentCLI describes a discovered agent CLI.
type AgentCLI struct {
	Bin      string
	Label    string
	Provider string
	Path     string
}

// DiscoverCLIs returns all agent CLIs found on PATH, in preference order.
func DiscoverCLIs() []AgentCLI {
	var found []AgentCLI
	for _, c := range knownCLIs {
		// "gh copilot" is a subcommand — check for `gh` then verify copilot extension
		if c.bin == "gh copilot" {
			if p, err := exec.LookPath("gh"); err == nil {
				out, err := exec.Command("gh", "copilot", "--version").Output()
				if err == nil && len(out) > 0 {
					found = append(found, AgentCLI{Bin: "gh copilot", Label: c.label, Provider: c.provider, Path: p + " (extension)"})
				}
			}
			continue
		}
		if p, err := exec.LookPath(c.bin); err == nil {
			found = append(found, AgentCLI{Bin: c.bin, Label: c.label, Provider: c.provider, Path: p})
		}
	}
	return found
}

func checkAgentCLIs() []Check {
	found := DiscoverCLIs()
	if len(found) == 0 {
		return []Check{{
			Name:    "agent CLI",
			Status:  Fail,
			Message: "no agent CLI found — install pi, claude, gemini, or gh copilot",
		}}
	}

	checks := make([]Check, len(found))
	for i, c := range found {
		checks[i] = Check{
			Name:    fmt.Sprintf("agent CLI: %s", c.Bin),
			Status:  OK,
			Message: c.Path,
		}
	}
	return checks
}
