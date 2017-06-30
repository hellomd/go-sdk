package testutils

import (
	"os"
	"testing"

	"github.com/DATA-DOG/godog"
)

// TestMain -
func TestMain(m *testing.M, FeatureContext func(suite *godog.Suite)) {
	status := godog.RunWithOptions("godogs", FeatureContext, godog.Options{
		Format: "progress",
		Paths:  []string{"features"},
	})
	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
