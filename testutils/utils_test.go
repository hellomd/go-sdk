package testutils

import "testing"
import "time"

func TestFillStruct(t *testing.T) {
	t.Run("string fields", func(t *testing.T) {
		// Arrange
		table := map[string]string{"Name": "John"}
		value := struct{ Name string }{}

		// Act
		err := FillStruct(&value, table)

		// Assert
		if err != nil {
			t.Error(err)
			return
		}

		if value.Name != "John" {
			t.Error("expected Name to be John, but was", value.Name)
			return
		}
	})

	t.Run("int fields", func(t *testing.T) {
		// Arrange
		table := map[string]string{"Age": "25"}
		value := struct{ Age int }{}

		// Act
		err := FillStruct(&value, table)

		// Assert
		if err != nil {
			t.Error(err)
			return
		}

		if value.Age != 25 {
			t.Error("expected Age to be 25, but was", value.Age)
			return
		}
	})

	t.Run("Time fields", func(t *testing.T) {
		// Arrange
		table := map[string]string{"When": "25 Jul 17 12:27 UTC"}
		value := struct{ When time.Time }{}

		// Act
		err := FillStruct(&value, table)

		// Assert
		if err != nil {
			t.Error(err)
			return
		}

		expected := time.Date(2017, time.Month(7), 25, 12, 27, 0, 0, time.UTC)
		if value.When != expected {
			t.Error("expected When to be", expected, ", but was", value.When)
			return
		}
	})
}
