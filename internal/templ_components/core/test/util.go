package test

import (
	"context"
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

func ParseSingleElement(t *testing.T, markup string) *html.Node {
	t.Helper()

	root, err := html.Parse(strings.NewReader(markup))
	if err != nil {
		t.Fatalf("failed to parse markup: %v", err)
	}

	node := GetFirstElement(root)
	if node == nil {
		t.Fatalf("no element node found: %s", markup)
	}

	return node
}

func GetFirstElement(n *html.Node) *html.Node {
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
		if el := GetFirstElement(c); el != nil {
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
