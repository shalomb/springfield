package integration

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var opts = godog.Options{
	Format: "pretty",
	Paths:  []string{"features"},
	Tags:   "~@wip",
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestFeatures(t *testing.T) {
	opts.TestingT = t
	opts.Output = colors.Colored(os.Stdout)
	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			InitializeSandboxingScenario(ctx)
			InitializeCliScenario(ctx)
			InitializeAgentLoopScenario(ctx)
			InitializeFeedbackLoopScenario(ctx)
			InitializeGovernanceScenario(ctx)
		},
		Options: &opts,
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
