package body

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"golang.org/x/net/html"

	test "github.com/landmaster135/devbox/internal/templ_components/core/test"
)

func TestTagWithoutClassRendersChildrenAndSkipsNil(t *testing.T) {
	t.Parallel()

	first := templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<section data-slot="alpha">Alpha</section>`)
		return err
	})

	second := templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<section data-slot="beta">Beta</section>`)
		return err
	})

	markup := renderBody(t, "", first, nil, second)

	node := test.ParseSingleElement(t, markup, true)
	body := ascendToBody(t, node)

	if test.HasAttr(body, "class") {
		t.Fatalf("expected no class attribute when class input empty, got %s", markup)
	}

	if count := strings.Count(markup, `<section data-slot="`); count != 2 {
		t.Fatalf("expected two section children but found %d in %s", count, markup)
	}

	if strings.Contains(markup, "<nil>") {
		t.Fatalf("nil children should not render markup, got %s", markup)
	}
}

func TestTagWithClassEscapesAttributeAndRendersChildren(t *testing.T) {
	t.Parallel()

	className := `layout main" onclick="alert('x')`

	child := templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<div id="app-root">Ready</div>`)
		return err
	})

	markup := renderBody(t, className, child)

	if strings.Contains(markup, `class="layout main" onclick="alert('x')"`) {
		t.Fatalf("class attribute should be escaped, got %s", markup)
	}

	node := test.ParseSingleElement(t, markup, true)
	body := ascendToBody(t, node)

	if got := test.GetAttr(body, "class"); got != className {
		t.Fatalf("expected decoded class %q but got %q", className, got)
	}

	if text := strings.TrimSpace(test.GetTextContent(body)); text != "Ready" {
		t.Fatalf("expected child content to render, got %q", text)
	}
}

func renderBody(t *testing.T, class string, children ...templ.Component) string {
	t.Helper()

	component := Tag(class, children...)
	if component == nil {
		t.Fatalf("Tag returned nil component")
	}

	var buf strings.Builder
	return test.RenderComponent(t, component, &buf)
}

func ascendToBody(t *testing.T, node *html.Node) *html.Node {
	t.Helper()

	for current := node; current != nil; current = current.Parent {
		if current.Data == "body" {
			return current
		}
	}

	t.Fatalf("body node not found in parsed markup from %s", test.GetTextContent(node))
	return nil
}
