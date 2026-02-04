package head

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"golang.org/x/net/html"

	test "github.com/landmaster135/devbox/internal/templ_components/core/test"
)

func TestTagRendersHeadAndChildren(t *testing.T) {
	t.Parallel()

	meta := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<meta charset="utf-8">`)
		return err
	})

	title := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<title>Playground</title>`)
		return err
	})

	markup := renderHeadTag(t, meta, title)

	htmlNode := test.ParseSingleElement(t, markup, false)
	headNode := test.GetFirstChildOfElementByData(htmlNode, "head")
	if headNode == nil {
		t.Fatalf("expected <head> as first element but got %v", headNode)
	}

	metaNode := test.GetFirstChildOfElement(headNode)
	if metaNode == nil || metaNode.Data != "meta" {
		t.Fatalf("expected first child element to be <meta>, got %v", metaNode)
	}

	if got := test.GetAttr(metaNode, "charset"); got != "utf-8" {
		t.Fatalf("expected charset utf-8 but got %q", got)
	}

	titleNode := test.GetNextSiblingOfElement(metaNode)
	if titleNode == nil || titleNode.Data != "title" {
		t.Fatalf("expected <title> sibling, got %v", titleNode)
	}

	if text := strings.TrimSpace(test.GetTextContent(titleNode)); text != "Playground" {
		t.Fatalf("expected title text Playground but got %q", text)
	}
}

func TestTagSkipsNilChildrenAndPreservesOrder(t *testing.T) {
	t.Parallel()

	link := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<link rel="stylesheet" href="/app.css">`)
		return err
	})

	script := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<script src="/app.js"></script>`)
		return err
	})

	markup := renderHeadTag(t, nil, link, nil, script, nil)

	htmlNode := test.ParseSingleElement(t, markup, false)
	headNode := test.GetFirstChildOfElementByData(htmlNode, "head")
	if headNode == nil {
		t.Fatalf("expected <head> as first element but got %v", headNode)
	}

	var childTags []string
	for child := headNode.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			childTags = append(childTags, child.Data)
		} else if child.Type == html.TextNode && strings.Contains(child.Data, "nil") {
			t.Fatalf("nil child leaked into markup: %q", child.Data)
		}
	}

	if len(childTags) != 2 {
		t.Fatalf("expected 2 child elements but got %d (%v)", len(childTags), childTags)
	}

	if childTags[0] != "link" || childTags[1] != "script" {
		t.Fatalf("expected child order [link script] but got %v", childTags)
	}

	linkNode := test.GetFirstChildOfElement(headNode)
	if linkNode == nil || linkNode.Data != "link" {
		t.Fatalf("expected <link> as first child, got %v", linkNode)
	}

	if got := test.GetAttr(linkNode, "href"); got != "/app.css" {
		t.Fatalf("expected href /app.css but got %q", got)
	}

	scriptNode := test.GetNextSiblingOfElementByData(linkNode, "script")

	if scriptNode == nil {
		t.Fatalf("expected to find <script> sibling, got nil")
	}

	if got := test.GetAttr(scriptNode, "src"); got != "/app.js" {
		t.Fatalf("expected script src /app.js but got %q", got)
	}

	if strings.Contains(markup, "<nil>") {
		t.Fatalf("markup should not include placeholders for nil components: %s", markup)
	}
}

func renderHeadTag(t *testing.T, children ...templ.Component) string {
	t.Helper()

	component := Tag(children...)
	if component == nil {
		t.Fatalf("Tag returned nil component")
	}

	var buf strings.Builder
	return test.RenderComponent(t, component, &buf)
}
