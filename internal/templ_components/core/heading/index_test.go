package heading

import (
	"context"
	"strings"
	"testing"
	"golang.org/x/net/html"
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
			node := parseSingleElement(t, markup)

			if node.Data != tt.wantTag {
				t.Fatalf("expected tag %s but got %s", tt.wantTag, node.Data)
			}

			if text := strings.TrimSpace(textContent(node)); text != tt.text {
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

	node := parseSingleElement(t, markup)
	if got := textContent(node); got != dangerous {
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

	if err := component.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	return buf.String()
}

func parseSingleElement(t *testing.T, markup string) *html.Node {
	t.Helper()

	root, err := html.Parse(strings.NewReader(markup))
	if err != nil {
		t.Fatalf("failed to parse markup: %v", err)
	}

	node := firstElement(root)
	if node == nil {
		t.Fatalf("no element node found: %s", markup)
	}

	return node
}

func firstElement(n *html.Node) *html.Node {
	if n == nil {
		return nil
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "html", "head", "body":
			// Skip the wrapper nodes html.Parse adds for fragments.
		default:
			return n
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if el := firstElement(c); el != nil {
			return el
		}
	}

	return nil
}

func textContent(n *html.Node) string {
	var sb strings.Builder
	accumulateText(n, &sb)
	return sb.String()
}

func accumulateText(n *html.Node, sb *strings.Builder) {
	if n == nil {
		return
	}

	if n.Type == html.TextNode {
		sb.WriteString(n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		accumulateText(c, sb)
	}
}
