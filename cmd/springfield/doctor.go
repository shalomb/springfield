package main

import (
	"fmt"
	"os"

	"github.com/shalomb/springfield/internal/setup"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check Springfield prerequisites and configuration",
	Long:  `Verifies that all required tools, configuration, and runtime state are present.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot determine working directory: %w", err)
		}

		root, err := setup.GitRoot(cwd)
		if err != nil {
			return err
		}
		fmt.Printf("Repository root: %s\n\n", root)

		checks := setup.RunChecks(root)

		anyFail := false
		for _, c := range checks {
			fmt.Printf("%s %-30s %s\n", c.Prefix(), c.Name, c.Message)
			if c.Status == setup.Fail {
				anyFail = true
			}
		}

		if anyFail {
			fmt.Println("\nSome checks failed. Run `springfield init` to fix setup issues.")
			return fmt.Errorf("doctor found failures")
		}
		fmt.Println("\nAll checks passed.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
