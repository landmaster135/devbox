package hidden_input

import (
	"strings"
	"testing"

	test "github.com/landmaster135/devbox/internal/templ_components/core/test"
)

var getAttr = test.GetAttr

func TestHiddenFieldRendersHiddenInput(t *testing.T) {
	t.Parallel()

	markup := renderHiddenField(t, "session_token", "abc123")
	node := test.ParseSingleElement(t, markup, true)

	if node.Data != "input" {
		t.Fatalf("expected <input> but got <%s>", node.Data)
	}

	if got := getAttr(node, "type"); got != "hidden" {
		t.Fatalf("expected type hidden but got %q", got)
	}

	if got := getAttr(node, "name"); got != "session_token" {
		t.Fatalf("expected name session_token but got %q", got)
	}

	if got := getAttr(node, "value"); got != "abc123" {
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

	node := test.ParseSingleElement(t, markup, true)

	if got := getAttr(node, "name"); got != name {
		t.Fatalf("expected decoded name %q but got %q", name, got)
	}

	if got := getAttr(node, "value"); got != value {
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

	return test.RenderComponent(t, component, &buf)
}
