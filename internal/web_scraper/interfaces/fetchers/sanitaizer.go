package fetchers

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func sanitizeHTMLBody(body string) (string, error) {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return "", fmt.Errorf("HTMLが空です")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("HTMLの解析に失敗しました: %w", err)
	}

	mainSelection, err := extractMainSelection(doc)
	if err != nil {
		return "", err
	}

	mainHTML, err := goquery.OuterHtml(mainSelection)
	if err != nil {
		return "", fmt.Errorf("main要素の抽出に失敗しました: %w", err)
	}

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(mainHTML))
	if err != nil {
		return "", fmt.Errorf("main要素の再解析に失敗しました: %w", err)
	}

	for _, node := range doc.Selection.Nodes {
		removeHTMLComments(node)
	}

	for _, selector := range getDefaultHTMLDenySelectors() {
		trimmedSelector := strings.TrimSpace(selector)
		if trimmedSelector == "" {
			continue
		}
		doc.Find(trimmedSelector).Remove()
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
		if htmlStr, err := doc.Html(); err == nil {
			builder.WriteString(htmlStr)
		}
	}

	if builder.Len() == 0 {
		return "", fmt.Errorf("サニタイズ後のHTML生成に失敗しました")
	}

	return collapseBlankLines(builder.String()), nil
}

func extractMainSelection(doc *goquery.Document) (*goquery.Selection, error) {
	if doc == nil {
		return nil, fmt.Errorf("HTMLドキュメントが初期化されていません")
	}

	if mainSelection := doc.Find("main").First(); mainSelection.Length() > 0 {
		return mainSelection, nil
	}

	if articleSelection := findLongestArticle(doc); articleSelection != nil {
		return articleSelection, nil
	}

	return nil, &mainNotFoundError{}
}

func findLongestArticle(doc *goquery.Document) *goquery.Selection {
	if doc == nil {
		return nil
	}

	var (
		longest *goquery.Selection
		maxLen  = -1
	)

	doc.Find("article").Each(func(_ int, s *goquery.Selection) {
		length := utf8.RuneCountInString(strings.TrimSpace(s.Text()))
		if length > maxLen {
			longest = s.Clone()
			maxLen = length
		}
	})

	return longest
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
	return append(items, itemsForReddit...)
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
