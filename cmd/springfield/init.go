package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/shalomb/springfield/internal/setup"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise Springfield in the current git repository",
	Long: `Sets up the td database, .gitignore entries, and a starter config.
Safe to re-run — all steps are idempotent.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot determine working directory: %w", err)
		}

		root, err := setup.GitRoot(cwd)
		if err != nil {
			warnf(cmd, "error: %v", err)
			warnf(cmd, "Springfield must be initialised inside a git repository.")
			return err
		}
		fmt.Printf("Repository root: %s\n\n", root)

		if err := runInit(cmd, root); err != nil {
			warnf(cmd, "\nerror: %v", err)
			warnf(cmd, "Initialisation incomplete. Fix the error above and re-run `springfield init`.")
			return err
		}

		fmt.Println("\nDone. Run `springfield doctor` to verify your setup.")
		fmt.Println("Next: create an epic with `td epic add 'Your feature'`, then `springfield orchestrate`.")
		return nil
	},
}

// warnf writes a formatted message to cmd's stderr, ignoring write errors
// (there is nothing useful to do if stderr itself is broken).
func warnf(cmd *cobra.Command, format string, a ...any) {
	_, _ = fmt.Fprintf(cmd.ErrOrStderr(), format+"\n", a...)
}

func runInit(cmd *cobra.Command, root string) error {
	// 1. td database
	fmt.Println("--- td database")
	if err := initTD(root); err != nil {
		warnf(cmd, "  failed: %v", err)
		return fmt.Errorf("td setup failed")
	}

	// 2. .gitignore
	fmt.Println("\n--- .gitignore")
	if err := ensureGitignore(root); err != nil {
		warnf(cmd, "  failed: %v", err)
		return fmt.Errorf(".gitignore setup failed")
	}

	// 3. config
	fmt.Println("\n--- config")
	if err := initConfig(cmd, root); err != nil {
		warnf(cmd, "  failed: %v", err)
		return fmt.Errorf("config setup failed")
	}

	// 4. .pi/agents
	fmt.Println("\n--- .pi/agents")
	if err := scaffoldPiAgents(cmd, root); err != nil {
		warnf(cmd, "  failed: %v", err)
		return fmt.Errorf(".pi/agents setup failed")
	}

	// 5. .todos
	fmt.Println("\n--- .todos")
	if err := scaffoldTodos(cmd, root); err != nil {
		warnf(cmd, "  failed: %v", err)
		return fmt.Errorf(".todos setup failed")
	}

	// 6. Justfile
	fmt.Println("\n--- Justfile")
	if err := scaffoldJustfile(cmd, root); err != nil {
		warnf(cmd, "  failed: %v", err)
		return fmt.Errorf("Justfile setup failed")
	}

	// 7. .codemap/config.json
	fmt.Println("\n--- .codemap/config.json")
	if err := scaffoldCodemapConfig(cmd, root); err != nil {
		warnf(cmd, "  failed: %v", err)
		return fmt.Errorf(".codemap/config.json setup failed")
	}

	return nil
}

func initTD(root string) error {
	dbPath := filepath.Join(root, ".todos", "issues.db")
	if _, err := os.Stat(dbPath); err == nil {
		fmt.Println("  already initialised")
		return nil
	}

	if _, err := exec.LookPath("td"); err != nil {
		return fmt.Errorf("td not found on PATH — install from https://github.com/shalomb/td")
	}

	out, err := exec.Command("td", "-w", root, "init").CombinedOutput()
	if err != nil {
		return fmt.Errorf("td init failed: %w\n%s", err, string(out))
	}
	fmt.Println("  initialised")
	return nil
}

var gitignoreEntries = []string{".todos/", "worktrees/"}

func ensureGitignore(root string) error {
	giPath := filepath.Join(root, ".gitignore")

	existing := ""
	data, err := os.ReadFile(giPath)
	if err == nil {
		existing = string(data)
	}

	var toAdd []string
	for _, entry := range gitignoreEntries {
		if !strings.Contains(existing, entry) {
			toAdd = append(toAdd, entry)
		}
	}

	if len(toAdd) == 0 {
		fmt.Println("  already up to date")
		return nil
	}

	f, err := os.OpenFile(giPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot write .gitignore: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if err := writeln(w, existing, toAdd); err != nil {
		return fmt.Errorf("cannot write .gitignore: %w", err)
	}

	fmt.Printf("  added: %s\n", strings.Join(toAdd, ", "))
	return nil
}

// writeln appends the gitignore entries to w, ensuring a leading newline when
// existing content doesn't end with one.
func writeln(w *bufio.Writer, existing string, entries []string) error {
	if len(existing) > 0 && !strings.HasSuffix(existing, "\n") {
		if _, err := io.WriteString(w, "\n"); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(w, "# Springfield runtime artifacts\n"); err != nil {
		return err
	}
	for _, entry := range entries {
		if _, err := io.WriteString(w, entry+"\n"); err != nil {
			return err
		}
	}
	return w.Flush()
}

func initConfig(cmd *cobra.Command, root string) error {
	for _, name := range []string{".springfield.toml", "config.toml"} {
		if _, err := os.Stat(filepath.Join(root, name)); err == nil {
			fmt.Printf("  found existing %s — skipping\n", name)
			return nil
		}
	}

	clis := setup.DiscoverCLIs()
	chosenCLI := ""
	chosenProvider := "anthropic"

	if len(clis) > 0 {
		chosenCLI = clis[0].Bin
		chosenProvider = clis[0].Provider
		fmt.Printf("  detected agent CLIs:")
		for _, c := range clis {
			fmt.Printf(" %s", c.Bin)
		}
		fmt.Println()
		fmt.Printf("  using: %s (provider: %s)\n", chosenCLI, chosenProvider)
	} else {
		warnf(cmd, "  warning: no agent CLI found — config will use defaults, edit before use")
	}

	_ = chosenCLI // informational; pi is invoked by the binary directly

	configPath := filepath.Join(root, ".springfield.toml")
	if err := os.WriteFile(configPath, []byte(buildDefaultConfig(chosenProvider)), 0644); err != nil {
		return fmt.Errorf("cannot write %s: %w", configPath, err)
	}
	fmt.Printf("  wrote %s\n", configPath)
	return nil
}

func buildDefaultConfig(provider string) string {
	return fmt.Sprintf(`# Springfield configuration
# See: docs/how-to/configure-agent-models.md

[agent]
model = "%s/claude-sonnet"
fallback_model = "%s/claude-haiku"
max_iterations = 20
budget = 100000

# Per-agent overrides — uncomment and adjust as needed
# [agents.lisa]
# model = "%s/claude-opus"   # Lisa does deep planning; use a stronger model
# max_iterations = 10
# budget = 200000

# [agents.ralph]
# max_iterations = 30
# budget = 150000

# [agents.bart]
# budget = 80000
`, provider, provider, provider)
}
func scaffoldPiAgents(cmd *cobra.Command, root string) error {
	dir := filepath.Join(root, ".pi", "agents")
	if _, err := os.Stat(dir); err == nil {
		fmt.Println("  already initialised")
		return nil
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, ".gitkeep"), []byte(""), 0644); err != nil {
		return err
	}
	fmt.Println("  initialised")
	return nil
}

func scaffoldTodos(cmd *cobra.Command, root string) error {
	dir := filepath.Join(root, ".todos")
	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	fmt.Println("  initialised")
	return nil
}

func scaffoldJustfile(cmd *cobra.Command, root string) error {
	jf := filepath.Join(root, "Justfile")
	if _, err := os.Stat(jf); err == nil {
		fmt.Println("  already initialised")
		return nil
	}
	content := []byte("default:\n\t@just --list\n")
	if err := os.WriteFile(jf, content, 0644); err != nil {
		return err
	}
	fmt.Println("  initialised")
	return nil
}

func scaffoldCodemapConfig(cmd *cobra.Command, root string) error {
	dir := filepath.Join(root, ".codemap")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	cf := filepath.Join(dir, "config.json")
	if _, err := os.Stat(cf); err == nil {
		fmt.Println("  already initialised")
		return nil
	}
	content := []byte("{\n  \"ignore\": [\n    \"node_modules\",\n    \"vendor\"\n  ]\n}\n")
	if err := os.WriteFile(cf, content, 0644); err != nil {
		return err
	}
	fmt.Println("  initialised")
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
