package code

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"

	test "github.com/landmaster135/devbox/internal/templ_components/core/test"
)

func TestTextRendersCodeWithoutClassAndEscapesText(t *testing.T) {
	t.Parallel()

	dangerous := `fmt.Println("<script>alert('x')</script>")`

	markup := renderCodeText(t, "", dangerous)
	node := test.ParseSingleElement(t, markup)

	if node.Data != "code" {
		t.Fatalf("expected <code> but got <%s>", node.Data)
	}

	if test.HasAttr(node, "class") {
		t.Fatalf("expected no class attribute but found one in %s", markup)
	}

	if strings.Contains(markup, "<script>") {
		t.Fatalf("raw script tag should be escaped, got %s", markup)
	}

	if got := test.GetTextContent(node); got != dangerous {
		t.Fatalf("expected code text %q but got %q", dangerous, got)
	}
}

func TestTextAppliesClassAndEscapesAttribute(t *testing.T) {
	t.Parallel()

	class := `code-block" onclick="alert('x')`

	markup := renderCodeText(t, class, "fmt.Println(42)")
	node := test.ParseSingleElement(t, markup)

	if node.Data != "code" {
		t.Fatalf("expected <code> but got <%s>", node.Data)
	}

	if strings.Contains(markup, `class="`+class+`"`) {
		t.Fatalf("class attribute should be escaped, got %s", markup)
	}

	if attr := test.GetAttr(node, "class"); attr != class {
		t.Fatalf("expected decoded class %q but got %q", class, attr)
	}

	if text := strings.TrimSpace(test.GetTextContent(node)); text != "fmt.Println(42)" {
		t.Fatalf("expected code text fmt.Println(42) but got %q", text)
	}
}

func TestContentRendersBodyComponent(t *testing.T) {
	t.Parallel()

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<span data-lang="go">fmt.Println(7)</span>`)
		return err
	})

	markup := renderCodeContent(t, "code-sample", body)
	node := test.ParseSingleElement(t, markup)

	if attr := test.GetAttr(node, "class"); attr != "code-sample" {
		t.Fatalf("expected class code-sample but got %q", attr)
	}

	child := test.GetFirstChildOfElement(node)
	if child == nil {
		t.Fatalf("expected rendered child element in %s", markup)
	}

	if child.Data != "span" {
		t.Fatalf("expected child span but got %s", child.Data)
	}

	if attr := test.GetAttr(child, "data-lang"); attr != "go" {
		t.Fatalf("expected data-lang go but got %q", attr)
	}

	if text := strings.TrimSpace(test.GetTextContent(child)); text != "fmt.Println(7)" {
		t.Fatalf("expected fmt.Println(7) but got %q", text)
	}
}

func TestContentWithNilBodyRendersEmptyCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		class string
	}{
		{name: "without class", class: ""},
		{name: "with class", class: "inline-code"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			markup := renderCodeContent(t, tt.class, nil)
			node := test.ParseSingleElement(t, markup)

			if node.Data != "code" {
				t.Fatalf("expected <code> but got <%s>", node.Data)
			}

			if tt.class == "" {
				if test.HasAttr(node, "class") {
					t.Fatalf("expected no class attribute but got one in %s", markup)
				}
			} else if attr := test.GetAttr(node, "class"); attr != tt.class {
				t.Fatalf("expected class %q but got %q", tt.class, attr)
			}

			if text := strings.TrimSpace(test.GetTextContent(node)); text != "" {
				t.Fatalf("expected empty code block but got %q", text)
			}
		})
	}
}

func renderCodeText(t *testing.T, class, text string) string {
	t.Helper()

	component := Text(class, text)
	if component == nil {
		t.Fatalf("Text returned nil component")
	}

	var buf strings.Builder
	return test.RenderComponent(t, component, &buf)
}

func renderCodeContent(t *testing.T, class string, body templ.Component) string {
	t.Helper()

	component := Content(class, body)
	if component == nil {
		t.Fatalf("Content returned nil component")
	}

	var buf strings.Builder
	return test.RenderComponent(t, component, &buf)
}
