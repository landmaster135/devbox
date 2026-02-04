package test

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"golang.org/x/net/html"
)

func RenderComponent(t *testing.T, component templ.Component, buf *strings.Builder) string {
	if err := component.Render(context.Background(), buf); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	return buf.String()
}

func ParseSingleElement(t *testing.T, markup string, ignoresBaseNodes bool) *html.Node {
	t.Helper()

	root, err := html.Parse(strings.NewReader(markup))
	if err != nil {
		t.Fatalf("failed to parse markup: %v", err)
	}

	node := GetFirstElement(root, ignoresBaseNodes)
	if node == nil {
		t.Fatalf("no element node found: %s", markup)
	}

	return node
}

func GetFirstElement(n *html.Node, ignoresBaseNodes bool) *html.Node {
	if n == nil {
		return nil
	}

	baseNodes := []string{"html", "head", "body"}

	if n.Type == html.ElementNode {
		switch {
		case ignoresBaseNodes && slices.Contains(baseNodes, n.Data):
			// Skip the wrapper nodes html.Parse adds for fragments.
		default:
			return n
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if el := GetFirstElement(c, ignoresBaseNodes); el != nil {
			return el
		}
	}

	return nil
}

func GetFirstChildOfElement(n *html.Node) *html.Node {
	if n == nil {
		return nil
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			return c
		}
	}

	return nil
}

func GetFirstChildOfElementByData(n *html.Node, data string) *html.Node {
	if n == nil {
		return nil
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == data {
			return c
		}
	}

	return nil
}

func GetNextSiblingOfElement(n *html.Node) *html.Node {
	if n == nil {
		return nil
	}

	for c := n.NextSibling; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			return c
		}
	}

	return nil
}

func GetNextSiblingOfElementByData(n *html.Node, data string) *html.Node {
	if n == nil {
		return nil
	}

	for c := n.NextSibling; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == data {
			return c
		}
	}

	return nil
}

func GetTextContent(n *html.Node) string {
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

func GetAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func HasAttr(n *html.Node, key string) bool {
	for _, a := range n.Attr {
		if a.Key == key {
			return true
		}
	}
	return false
}
