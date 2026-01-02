package fetchers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"golang.org/x/net/html"
)

const defaultLaunchTimeout = 90 * time.Second

// RodDOMFetcher はgithub.com/go-rod/rodを利用したDOM取得実装です。
type RodDOMFetcher struct {
	launchTimeout time.Duration
}

type mainNotFoundError struct{}

func (e *mainNotFoundError) Error() string {
	return "main要素が見つかりません"
}

// Option はRodDOMFetcherのオプションです。
type Option func(*RodDOMFetcher)

// WithLaunchTimeout はブラウザ起動のタイムアウトを設定します。
func WithLaunchTimeout(d time.Duration) Option {
	return func(f *RodDOMFetcher) {
		if d > 0 {
			f.launchTimeout = d
		}
	}
}

// NewRodDOMFetcher はRodDOMFetcherのインスタンスを生成します。
func NewRodDOMFetcher(opts ...Option) *RodDOMFetcher {
	fetcher := &RodDOMFetcher{launchTimeout: defaultLaunchTimeout}
	for _, opt := range opts {
		if opt != nil {
			opt(fetcher)
		}
	}
	return fetcher
}

// FetchDOM は対象URLのDOMツリーを取得します。
func (f *RodDOMFetcher) FetchDOM(ctx context.Context, targetURL string, wait time.Duration, denySelectors []string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var html string
	err := rod.Try(func() {
		launchCtx := ctx
		var cancel context.CancelFunc
		if f.launchTimeout > 0 {
			launchCtx, cancel = context.WithTimeout(ctx, f.launchTimeout)
			defer cancel()
		}

		l := launcher.New().Context(launchCtx)
		controlURL := l.MustLaunch()
		defer l.Cleanup()

		browser := rod.New().ControlURL(controlURL).Context(ctx).MustConnect()
		defer browser.MustClose()

		page := browser.MustPage(targetURL)
		page.MustWaitLoad()

		if wait > 0 {
			timer := time.NewTimer(wait)
			defer timer.Stop()

			select {
			case <-ctx.Done():
				panic(ctx.Err())
			case <-timer.C:
			}
		}

		html = page.MustEval(`() => document.documentElement.outerHTML`).Str()
	})
	if err != nil {
		return "", fmt.Errorf("DOMツリーの取得に失敗しました: %w", err)
	}

	sanitized, err := sanitizeHTMLBody(html, denySelectors)
	if err != nil {
		var notFoundErr *mainNotFoundError
		if errors.As(err, &notFoundErr) {
			return "", fmt.Errorf("main要素が見つかりません: %w", err)
		}
		return "", fmt.Errorf("HTMLのサニタイズに失敗しました: %w", err)
	}

	return sanitized, nil
}

func sanitizeHTMLBody(body string, denySelectors []string) (string, error) {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return "", fmt.Errorf("HTMLが空です")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("HTMLの解析に失敗しました: %w", err)
	}

	mainSelection := doc.Find("main").First()
	if mainSelection.Length() == 0 {
		return "", &mainNotFoundError{}
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

	for _, selector := range denySelectors {
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
