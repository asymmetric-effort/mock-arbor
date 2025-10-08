package util

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	const key = "TEST_GETENV_KEY"
	// Unset -> default
	_ = os.Unsetenv(key)
	if got := GetEnv(key, "def"); got != "def" {
		t.Fatalf("GetEnv unset got %q; want def", got)
	}
	// Set -> value
	if err := os.Setenv(key, "value"); err != nil {
		t.Fatalf("setenv: %v", err)
	}
	if got := GetEnv(key, "def"); got != "value" {
		t.Fatalf("GetEnv set got %q; want value", got)
	}
}
