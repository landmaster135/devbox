package common

import "testing"

func TestIsBlank(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{name: "empty", value: "", expected: true},
		{name: "spaces", value: "   ", expected: true},
		{name: "tab and newline", value: "\t\n", expected: true},
		{name: "non blank", value: "value", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBlank(tt.value); got != tt.expected {
				t.Fatalf("IsBlank mismatch: expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestShellQuote(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{name: "simple", value: "abc", expected: "'abc'"},
		{name: "contains quote", value: "a'b", expected: "'a'\"'\"'b'"},
		{name: "empty", value: "", expected: "''"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShellQuote(tt.value); got != tt.expected {
				t.Fatalf("ShellQuote mismatch: expected %q, got %q", tt.expected, got)
			}
		})
	}
}
