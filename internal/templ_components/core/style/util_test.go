package style

import (
	"testing"

	"github.com/a-h/templ"
)

func TestMustHelpersHandleErrors(t *testing.T) {
	tests := []struct {
		name string
		call func() (templ.SafeCSS, error)
	}{
		{name: "box shadow", call: func() (templ.SafeCSS, error) { return MustSafeBoxShadow("box-shadow", nil, []string{"#fff"}) }},
		{name: "border", call: func() (templ.SafeCSS, error) { return MustSafeBorder("border", "", "solid", []string{"#000"}) }},
		{name: "color", call: func() (templ.SafeCSS, error) { return MustSafeColorProperty("color", "not-a-color") }},
		{name: "font", call: func() (templ.SafeCSS, error) { return MustSafeFontFamily("font-family", nil) }},
		{name: "linear gradient", call: func() (templ.SafeCSS, error) { return MustSafeLinearGradient("background", "90deg", "#fff") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.call(); err == nil {
				t.Fatalf("expected error for %s", tt.name)
			}
		})
	}
}

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
