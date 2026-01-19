package services

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func sanitizeHTMLBody(body string) (string, bool) {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return body, false
	}
	if !strings.Contains(trimmed, "<") || !strings.Contains(trimmed, ">") {
		return body, false
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return body, false
	}

	mainSelection := doc.Find("main").First()
	if mainSelection.Length() == 0 {
		return body, false
	}

	mainHTML, err := goquery.OuterHtml(mainSelection)
	if err != nil {
		return body, false
	}

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(mainHTML))
	if err != nil {
		return body, false
	}

	for _, node := range doc.Selection.Nodes {
		removeHTMLComments(node)
	}

	for _, selector := range getDefaultHTMLDenySelectors() {
		if sel := strings.TrimSpace(selector); sel != "" {
			doc.Find(sel).Remove()
		}
	}

	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		s.RemoveAttr("class")
		s.RemoveAttr("style")
		s.RemoveAttr("data-testid")

		switch strings.ToLower(goquery.NodeName(s)) {
		case "span":
			s.RemoveAttr("data-allow-missmatch")
			s.RemoveAttr("data-allow-mismatch")
		case "img":
			s.RemoveAttr("onerror")
			s.RemoveAttr("data-nuxt-img")
			s.RemoveAttr("sizes")
			s.RemoveAttr("srcset")
		}
	})

	removeEmptyDivs(doc)

	var builder strings.Builder
	if bodyNode := doc.Find("body").First(); bodyNode.Length() > 0 {
		if outer, err := goquery.OuterHtml(bodyNode); err == nil {
			builder.WriteString(outer)
		}
	}

	if builder.Len() == 0 {
		if htmlNode := doc.Find("html").First(); htmlNode.Length() > 0 {
			if outer, err := goquery.OuterHtml(htmlNode); err == nil {
				builder.WriteString(outer)
			}
		}
	}

	if builder.Len() == 0 {
		if html, err := doc.Html(); err == nil {
			builder.WriteString(html)
		}
	}

	if builder.Len() == 0 {
		return body, false
	}

	return collapseBlankLines(builder.String()), true
}

func getDefaultHTMLDenySelectors() []string {
	itemsForReddit := []string{
		"pdp-back-button",
		"faceplate-loader",
		"faceplate-tracker",
		"faceplate-perfmark",
		"faceplate-number",
		"faceplate-dropdown-menu",
		"faceplate-partial",
		"shreddit-comments-page-ad",
		"shreddit-async-loader",
		"shreddit-comment-tree-ad",
		"button",
	}
	items := []string{
		"svg",
		"script",
		"header",
		"footer",
	}
	items = append(items, itemsForReddit...)

	return items
}

func removeHTMLComments(node *html.Node) {
	if node == nil {
		return
	}

	for child := node.FirstChild; child != nil; {
		next := child.NextSibling
		if child.Type == html.CommentNode {
			node.RemoveChild(child)
		} else {
			removeHTMLComments(child)
		}
		child = next
	}
}

func removeEmptyDivs(doc *goquery.Document) {
	if doc == nil {
		return
	}
	for {
		removed := false
		doc.Find("div").Each(func(_ int, s *goquery.Selection) {
			if s.Children().Length() == 0 && strings.TrimSpace(s.Text()) == "" {
				s.Remove()
				removed = true
			}
		})
		if !removed {
			break
		}
	}
}

func collapseBlankLines(src string) string {
	lines := strings.Split(src, "\n")
	if len(lines) == 1 {
		return src
	}

	var builder strings.Builder
	consecutiveBlank := 0

	for _, line := range lines {
		isBlank := strings.TrimSpace(line) == ""
		if isBlank {
			if consecutiveBlank > 0 {
				continue
			}
			consecutiveBlank = 1
		} else {
			consecutiveBlank = 0
		}

		if builder.Len() > 0 {
			builder.WriteByte('\n')
		}
		builder.WriteString(line)
	}

	return builder.String()
}
