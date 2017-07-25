package testutils

import (
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

var supportedLayouts = []string{
	time.RFC822,
	time.RFC3339,
	time.RFC3339Nano,
}

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

// ParseTable receives a godog gherkin table and returns a map
// containing the rows of the table
func ParseTable(table *gherkin.DataTable) []map[string]string {
	if len(table.Rows) == 0 {
		return []map[string]string{}
	}

	headRow := table.Rows[0]

	valueRows := table.Rows[1:]
	values := make([]map[string]string, len(valueRows))
	for i := 0; i < len(valueRows); i++ {
		rowMap := map[string]string{}
		for i, cell := range valueRows[i].Cells {
			rowMap[headRow.Cells[i].Value] = cell.Value
		}
		values[i] = rowMap
	}

	return values
}
