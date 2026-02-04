package main_component

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"golang.org/x/net/html"

	test "github.com/landmaster135/devbox/internal/templ_components/core/test"
)

func TestTagRendersMainWithoutClass(t *testing.T) {
	t.Parallel()

	child := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<section data-test="child">content</section>`)
		return err
	})

	markup := renderMainTag(t, "", child)
	node := test.ParseSingleElement(t, markup, true)

	if node.Data != "main" {
		t.Fatalf("expected <main> but got <%s>", node.Data)
	}

	if test.HasAttr(node, "class") {
		t.Fatalf("expected no class attribute when class is empty: %s", markup)
	}

	childNode := test.GetFirstChildOfElement(node)
	if childNode == nil {
		t.Fatalf("expected child element but none rendered: %s", markup)
	}

	if childNode.Data != "section" {
		t.Fatalf("expected child section but got <%s>", childNode.Data)
	}

	if attr := test.GetAttr(childNode, "data-test"); attr != "child" {
		t.Fatalf("expected data-test attribute to be child but got %q", attr)
	}

	if text := strings.TrimSpace(test.GetTextContent(childNode)); text != "content" {
		t.Fatalf("expected child content text but got %q", text)
	}
}

func TestTagEscapesClassAttribute(t *testing.T) {
	t.Parallel()

	class := `layout-main" onclick="alert('x')`

	markup := renderMainTag(t, class)
	node := test.ParseSingleElement(t, markup, true)

	if node.Data != "main" {
		t.Fatalf("expected <main> but got <%s>", node.Data)
	}

	if strings.Contains(markup, `class="`+class+`"`) {
		t.Fatalf("class attribute should be escaped, got %s", markup)
	}

	if attr := test.GetAttr(node, "class"); attr != class {
		t.Fatalf("expected decoded class %q but got %q", class, attr)
	}
}

func TestTagSkipsNilChildren(t *testing.T) {
	t.Parallel()

	first := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<section data-order="first"></section>`)
		return err
	})

	last := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<span data-order="last">done</span>`)
		return err
	})

	markup := renderMainTag(t, "wrapper", first, nil, last)
	node := test.ParseSingleElement(t, markup, true)

	if attr := test.GetAttr(node, "class"); attr != "wrapper" {
		t.Fatalf("expected wrapper class but got %q", attr)
	}

	orders := collectChildOrders(node)
	if len(orders) != 2 {
		t.Fatalf("expected 2 rendered children but got %d in %s", len(orders), markup)
	}

	if orders[0] != "first" || orders[1] != "last" {
		t.Fatalf("expected child order [first last] but got %v", orders)
	}
}

func renderMainTag(t *testing.T, class string, children ...templ.Component) string {
	t.Helper()

	component := Tag(class, children...)
	if component == nil {
		t.Fatalf("Tag returned nil component")
	}

	var buf strings.Builder
	return test.RenderComponent(t, component, &buf)
}

func collectChildOrders(node *html.Node) []string {
	var orders []string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.ElementNode {
			continue
		}

		attr := test.GetAttr(child, "data-order")
		if attr != "" {
			orders = append(orders, attr)
		}
	}

	return orders
}
