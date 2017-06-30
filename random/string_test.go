package random

import "testing"

func TestString(t *testing.T) {
	expectedSize := 10

	strings := []string{
		String(expectedSize),
		String(expectedSize),
		String(expectedSize),
	}

	for _, value := range strings {
		if size := len(value); size != expectedSize {
			t.Errorf(
				"Random string length mismatch. Expected %d, got %d",
				expectedSize,
				size,
			)
		}
	}
}
