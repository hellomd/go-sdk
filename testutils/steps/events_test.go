package steps

import (
	"testing"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/hellomd/go-sdk/testutils/fakes"
)

func TestTheFollowingEventsShouldHaveBeenPublished(t *testing.T) {
	t.Run("pass with two matching events", func(t *testing.T) {
		// Arrange
		table := buildTable([][]string{
			{"Key", "Body"},
			{"some.key", `"foo"`},
			{"some.key", `"bar"`},
		})

		pub := new(fakes.Publisher)
		pub.Publish("some.key", "foo")
		pub.Publish("some.key", "bar")

		// Act
		err := TheFollowingEventsShouldHaveBeenPublished(table, pub)

		// Assert
		if err != nil {
			t.Error("unexpected error:", err)
		}
	})

	t.Run("pass with one complex json", func(t *testing.T) {
		// Arrange
		table := buildTable([][]string{
			{"Key", "Body"},
			{"some.key", `{"id": "one", "size": 3}`},
		})

		pub := new(fakes.Publisher)
		pub.Publish("some.key", map[string]interface{}{"id": "one", "size": 3})

		// Act
		err := TheFollowingEventsShouldHaveBeenPublished(table, pub)

		// Assert
		if err != nil {
			t.Error("unexpected error:", err)
		}
	})

	t.Run("pass with no events", func(t *testing.T) {
		// Arrange
		table := buildTable([][]string{{"Key", "Body"}})

		pub := new(fakes.Publisher)

		// Act
		err := TheFollowingEventsShouldHaveBeenPublished(table, pub)

		// Assert
		if err != nil {
			t.Error("unexpected error:", err)
		}
	})

	t.Run("fail with fewer events than expected", func(t *testing.T) {
		// Arrange
		table := buildTable([][]string{
			{"Key", "Body"},
			{"some.key", `"foo"`},
			{"some.key", `"bar"`},
		})

		pub := new(fakes.Publisher)
		pub.Publish("some.key", "foo")

		// Act
		err := TheFollowingEventsShouldHaveBeenPublished(table, pub)

		// Assert
		if err == nil {
			t.Error("expected an error, but got none")
		}
	})

	t.Run("fail with more events than expected", func(t *testing.T) {
		// Arrange
		table := buildTable([][]string{
			{"Key", "Body"},
			{"some.key", `"foo"`},
		})

		pub := new(fakes.Publisher)
		pub.Publish("some.key", "foo")
		pub.Publish("some.key", "bar")

		// Act
		err := TheFollowingEventsShouldHaveBeenPublished(table, pub)

		// Assert
		if err == nil {
			t.Error("expected an error, but got none")
		}
	})

	t.Run("fail with one complex json", func(t *testing.T) {
		// Arrange
		table := buildTable([][]string{
			{"Key", "Body"},
			{"some.key", `{"id": "one", "size": 3}`},
		})

		pub := new(fakes.Publisher)
		pub.Publish("some.key", map[string]interface{}{"id": "oh no!", "size": 999})

		// Act
		err := TheFollowingEventsShouldHaveBeenPublished(table, pub)

		// Assert
		if err == nil {
			t.Error("expected an error, but got none")
		}
	})
}

func buildTable(src [][]string) *gherkin.DataTable {
	tableRows := []*gherkin.TableRow{}
	for _, row := range src {
		tableCells := []*gherkin.TableCell{}
		for _, cell := range row {
			tableCells = append(tableCells, &gherkin.TableCell{Value: cell})
		}

		tableRows = append(tableRows, &gherkin.TableRow{Cells: tableCells})
	}

	return &gherkin.DataTable{Rows: tableRows}
}
