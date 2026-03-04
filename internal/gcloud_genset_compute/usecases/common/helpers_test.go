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

func TestShellQuoteSSHKeyPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{name: "default path", value: "$HOME/.ssh/google_compute_engine", expected: "\"$HOME/.ssh/google_compute_engine\""},
		{name: "custom path", value: "/tmp/key", expected: "'/tmp/key'"},
		{name: "contains quote", value: "/tmp/key'o", expected: "'/tmp/key'\"'\"'o'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShellQuoteSSHKeyPath(tt.value); got != tt.expected {
				t.Fatalf("ShellQuoteSSHKeyPath mismatch: expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestBuildSSHAgentSetupCommand(t *testing.T) {
	t.Parallel()

	command := BuildSSHAgentSetupCommand("$HOME/.ssh/google_compute_engine")
	expected := "if [ -z \"${SSH_AUTH_SOCK:-}\" ]; then eval \"$(ssh-agent -s)\" >/dev/null; fi && ssh-add \"$HOME/.ssh/google_compute_engine\""
	if command != expected {
		t.Fatalf("BuildSSHAgentSetupCommand mismatch:\nexpected: %s\nactual:   %s", expected, command)
	}
}
