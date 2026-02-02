package style

import "testing"

func TestSanitizeFontFamily(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "keyword", input: "Sans-Serif", expected: "sans-serif"},
		{name: "simple", input: "Inter", expected: "Inter"},
		{name: "multi word", input: "JetBrains   Mono", expected: "\"JetBrains Mono\""},
		{name: "system", input: "system-ui", expected: "system-ui"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sanitizeFontFamily(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestSanitizeFontFamilyErrors(t *testing.T) {
	tests := []string{"", "Comic Sans, Arial", "Bad@Font", "\"Inter\""}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			if _, err := sanitizeFontFamily(input); err == nil {
				t.Fatalf("expected error for %q", input)
			}
		})
	}
}

func TestSafeFontFamily(t *testing.T) {
	css, err := safeFontFamily("font-family", []string{"Inter", "Segoe UI", "system-ui", "monospace"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "font-family:Inter, \"Segoe UI\", system-ui, monospace;"
	if string(css) != expected {
		t.Fatalf("expected %q, got %q", expected, css)
	}
}
