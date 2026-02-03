package hidden_input

import (
	"context"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestHiddenFieldRendersHiddenInput(t *testing.T) {
	t.Parallel()

	markup := renderHiddenField(t, "session_token", "abc123")
	node := parseSingleInput(t, markup)

	if node.Data != "input" {
		t.Fatalf("expected <input> but got <%s>", node.Data)
	}

	if got := attr(node, "type"); got != "hidden" {
		t.Fatalf("expected type hidden but got %q", got)
	}

	if got := attr(node, "name"); got != "session_token" {
		t.Fatalf("expected name session_token but got %q", got)
	}

	if got := attr(node, "value"); got != "abc123" {
		t.Fatalf("expected value abc123 but got %q", got)
	}
}

func TestHiddenFieldEscapesAttributeValues(t *testing.T) {
	t.Parallel()

	name := `"onfocus=alert('x')`
	value := `</script><script>alert(1)</script>`
	markup := renderHiddenField(t, name, value)

	if strings.Contains(markup, value) {
		t.Fatalf("raw value should be escaped in markup: %s", markup)
	}

	node := parseSingleInput(t, markup)

	if got := attr(node, "name"); got != name {
		t.Fatalf("expected decoded name %q but got %q", name, got)
	}

	if got := attr(node, "value"); got != value {
		t.Fatalf("expected decoded value %q but got %q", value, got)
	}
}

func renderHiddenField(t *testing.T, name, value string) string {
	t.Helper()

	var buf strings.Builder
	component := HiddenField(name, value)
	if component == nil {
		t.Fatalf("HiddenField returned nil component")
	}

	if err := component.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	return buf.String()
}

func parseSingleInput(t *testing.T, markup string) *html.Node {
	t.Helper()

	root, err := html.Parse(strings.NewReader(markup))
	if err != nil {
		t.Fatalf("failed to parse markup: %v", err)
	}

	node := firstElement(root)
	if node == nil {
		t.Fatalf("no element in markup: %s", markup)
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

func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
