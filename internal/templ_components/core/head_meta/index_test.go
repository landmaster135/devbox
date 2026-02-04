package head_meta

import (
	"strings"
	"testing"

	test "github.com/landmaster135/devbox/internal/templ_components/core/test"
)

func TestBaseRendersMetaTags(t *testing.T) {
	t.Parallel()

	markup := renderBase(t)

	charsetMeta := test.ParseSingleElement(t, markup, true)
	if charsetMeta.Data != "meta" {
		t.Fatalf("expected first element to be <meta> but got <%s>", charsetMeta.Data)
	}

	if got := test.GetAttr(charsetMeta, "charset"); got != "utf-8" {
		t.Fatalf("expected charset utf-8 but got %q", got)
	}

	viewportMeta := test.GetNextSiblingOfElement(charsetMeta)
	if viewportMeta == nil {
		t.Fatalf("expected second <meta> element but none found")
	}

	if viewportMeta.Data != "meta" {
		t.Fatalf("expected viewport element to be <meta> but got <%s>", viewportMeta.Data)
	}

	if got := test.GetAttr(viewportMeta, "name"); got != "viewport" {
		t.Fatalf("expected viewport meta name but got %q", got)
	}

	const wantContent = "width=device-width, initial-scale=1"
	if got := test.GetAttr(viewportMeta, "content"); got != wantContent {
		t.Fatalf("expected viewport content %q but got %q", wantContent, got)
	}

	if tail := test.GetNextSiblingOfElement(viewportMeta); tail != nil {
		t.Fatalf("expected exactly two meta tags but found additional <%s>", tail.Data)
	}
}

func renderBase(t *testing.T) string {
	t.Helper()

	var buf strings.Builder
	component := Base()
	if component == nil {
		t.Fatalf("Base returned nil component")
	}

	return test.RenderComponent(t, component, &buf)
}
