package heading

import (
	"strings"
	"testing"

	test "github.com/landmaster135/devbox/internal/templ_components/core/test"
)

func TestHeadingLevelSelection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		level   int
		text    string
		wantTag string
	}{
		{name: "level one", level: 1, text: "Primary", wantTag: "h1"},
		{name: "non-positive uses h1", level: -2, text: "NonPositive", wantTag: "h1"},
		{name: "level two", level: 2, text: "Secondary", wantTag: "h2"},
		{name: "fallback", level: 5, text: "Fallback", wantTag: "h3"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			markup := renderHeadingHTML(t, tt.level, tt.text)
			node := test.ParseSingleElement(t, markup)

			if node.Data != tt.wantTag {
				t.Fatalf("expected tag %s but got %s", tt.wantTag, node.Data)
			}

			if text := strings.TrimSpace(test.GetTextContent(node)); text != tt.text {
				t.Fatalf("expected text %q but got %q", tt.text, text)
			}
		})
	}
}

func TestHeadingEscapesText(t *testing.T) {
	t.Parallel()

	dangerous := `<script>alert("x")</script>`
	markup := renderHeadingHTML(t, 3, dangerous)

	if strings.Contains(markup, "<script>") {
		t.Fatalf("raw script tag should be escaped, but got %s", markup)
	}

	if !strings.Contains(markup, "&lt;script&gt;") {
		t.Fatalf("escaped script tag not found in %s", markup)
	}

	node := test.ParseSingleElement(t, markup)
	if got := test.GetTextContent(node); got != dangerous {
		t.Fatalf("expected text %q but got %q", dangerous, got)
	}
}

func renderHeadingHTML(t *testing.T, level int, text string) string {
	t.Helper()

	var buf strings.Builder
	component := Heading(level, text)
	if component == nil {
		t.Fatalf("Heading returned nil component")
	}

	return test.RenderComponent(t, component, &buf)
}
