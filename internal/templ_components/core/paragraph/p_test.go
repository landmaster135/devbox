package paragraph

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"

	test "github.com/landmaster135/devbox/internal/templ_components/core/test"
)

func TestTextRendersParagraphAndEscapesText(t *testing.T) {
	t.Parallel()

	dangerous := `Next <script>alert("x")</script> Step`

	markup := renderParagraphText(t, "", dangerous)
	node := test.ParseSingleElement(t, markup)

	if node.Data != "p" {
		t.Fatalf("expected <p> but got <%s>", node.Data)
	}

	if class := test.GetAttr(node, "class"); class != "" {
		t.Fatalf("expected empty class but got %q", class)
	}

	if strings.Contains(markup, "<script>") {
		t.Fatalf("raw script tag should be escaped, got %s", markup)
	}

	if got := test.GetTextContent(node); got != dangerous {
		t.Fatalf("expected text %q but got %q", dangerous, got)
	}
}

func TestTextAppliesClassAndEscapesAttribute(t *testing.T) {
	t.Parallel()

	class := `body-copy" onclick="alert('x')`

	markup := renderParagraphText(t, class, "Safe")
	node := test.ParseSingleElement(t, markup)

	if node.Data != "p" {
		t.Fatalf("expected <p> but got <%s>", node.Data)
	}

	if strings.Contains(markup, `class="`+class+`"`) {
		t.Fatalf("class attribute should be escaped, got %s", markup)
	}

	if attr := test.GetAttr(node, "class"); attr != class {
		t.Fatalf("expected decoded class %q but got %q", class, attr)
	}

	if text := test.GetTextContent(node); text != "Safe" {
		t.Fatalf("expected text Safe but got %q", text)
	}
}

func TestContentRendersBodyComponent(t *testing.T) {
	t.Parallel()

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<span class="inner">embedded</span>`)
		return err
	})

	markup := renderParagraphContent(t, "with-component", body)
	node := test.ParseSingleElement(t, markup)

	if attr := test.GetAttr(node, "class"); attr != "with-component" {
		t.Fatalf("expected class with-component but got %q", attr)
	}

	child := test.GetFirstChildOfElement(node)
	if child == nil {
		t.Fatalf("expected child element but found none in %s", markup)
	}

	if child.Data != "span" {
		t.Fatalf("expected child span but got %s", child.Data)
	}

	if attr := test.GetAttr(child, "class"); attr != "inner" {
		t.Fatalf("expected inner class but got %q", attr)
	}

	if text := test.GetTextContent(child); text != "embedded" {
		t.Fatalf("expected embedded text but got %q", text)
	}
}

func TestContentWithNilBodyRendersEmptyParagraph(t *testing.T) {
	t.Parallel()

	markup := renderParagraphContent(t, "", nil)
	node := test.ParseSingleElement(t, markup)

	if node.Data != "p" {
		t.Fatalf("expected <p> but got <%s>", node.Data)
	}

	if got := strings.TrimSpace(test.GetTextContent(node)); got != "" {
		t.Fatalf("expected empty paragraph but got %q", got)
	}
}

func TestStatusRendersDataAttribute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		class string
	}{
		{name: "without class", class: ""},
		{name: "with class", class: "manual"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			markup := renderParagraphStatus(t, tt.class)
			node := test.ParseSingleElement(t, markup)

			if node.Data != "p" {
				t.Fatalf("expected <p> but got <%s>", node.Data)
			}

			if attr := test.GetAttr(node, "class"); attr != tt.class {
				t.Fatalf("expected class %q but got %q", tt.class, attr)
			}

			if !test.HasAttr(node, "data-manual-run-status") {
				t.Fatalf("expected data-manual-run-status attribute in %s", markup)
			}

			if text := strings.TrimSpace(test.GetTextContent(node)); text != "" {
				t.Fatalf("status paragraph should be empty but got %q", text)
			}
		})
	}
}

func renderParagraphText(t *testing.T, class, text string) string {
	t.Helper()

	component := Text(class, text)
	if component == nil {
		t.Fatalf("Text returned nil component")
	}

	var buf strings.Builder
	return test.RenderComponent(t, component, &buf)
}

func renderParagraphContent(t *testing.T, class string, body templ.Component) string {
	t.Helper()

	component := Content(class, body)
	if component == nil {
		t.Fatalf("Content returned nil component")
	}

	var buf strings.Builder
	return test.RenderComponent(t, component, &buf)
}

func renderParagraphStatus(t *testing.T, class string) string {
	t.Helper()

	component := Status(class)
	if component == nil {
		t.Fatalf("Status returned nil component")
	}

	var buf strings.Builder
	return test.RenderComponent(t, component, &buf)
}
