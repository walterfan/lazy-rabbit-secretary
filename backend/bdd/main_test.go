package bdd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/walterfan/lazy-rabbit-reminder/bdd/steps"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can be changed to "pretty", "cucumber", "junit", etc.
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestFeatures(t *testing.T) {
	o := opts
	o.TestingT = t

	status := godog.TestSuite{
		Name:                "lazy-rabbit-reminder BDD tests",
		ScenarioInitializer: InitializeScenario,
		Options:             &o,
	}.Run()

	if status == 2 {
		t.SkipNow()
	}

	if status != 0 {
		t.Fatalf("zero status code expected, %d received", status)
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	// Initialize user registration context
	userRegCtx := steps.NewUserRegistrationContext()
	userRegCtx.InitializeScenario(sc)

	// Global before scenario hook
	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		fmt.Printf("Starting scenario: %s\n", sc.Name)
		return ctx, nil
	})

	// Global after scenario hook
	sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if err != nil {
			fmt.Printf("Scenario failed: %s - %v\n", sc.Name, err)
		} else {
			fmt.Printf("Scenario passed: %s\n", sc.Name)
		}
		return ctx, nil
	})
}

func TestMain(m *testing.M) {
	flag.Parse()
	opts.Paths = []string{"../features"}

	status := godog.TestSuite{
		Name:                "lazy-rabbit-reminder BDD tests",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
