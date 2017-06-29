package config

import (
	"os"
	"testing"
)

const TestKey Key = "CONFIG_TEST_KEY"

func TestConfig(t *testing.T) {
	result := Get(TestKey)
	if result != "" {
		t.Errorf("Expected \"\", got: %s", result)
	}

	expected := "expected"
	os.Setenv(string(TestKey), expected)
	defer os.Setenv(string(TestKey), "")

	result = Get(TestKey)
	if result != expected {
		t.Errorf("Expected %s, got: %s", expected, result)
	}
}
