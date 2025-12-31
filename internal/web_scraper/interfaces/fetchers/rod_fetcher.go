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

	var mainHTML string
	err = rod.Try(func() {
		doc, parseErr := goquery.NewDocumentFromReader(strings.NewReader(html))
		if parseErr != nil {
			panic(parseErr)
		}

		selection := doc.Find("main").First()
		if selection.Length() == 0 {
			panic(&mainNotFoundError{})
		}

		var outer string
		outer, parseErr = goquery.OuterHtml(selection)
		if parseErr != nil {
			panic(parseErr)
		}
		mainHTML = outer
	})
	if err != nil {
		var notFoundErr *mainNotFoundError
		if errors.As(err, &notFoundErr) {
			return "", fmt.Errorf("main要素が見つかりません: %w", err)
		}
		return "", fmt.Errorf("main要素の抽出に失敗しました: %w", err)
	}

	sanitized := mainHTML
	if len(denySelectors) == 0 {
		return sanitized, nil
	}

	err = rod.Try(func() {
		doc, parseErr := goquery.NewDocumentFromReader(strings.NewReader(mainHTML))
		if parseErr != nil {
			panic(parseErr)
		}

		selection := doc.Find("main").First()
		if selection.Length() == 0 {
			panic(&mainNotFoundError{})
		}

		for _, selector := range denySelectors {
			trimmed := strings.TrimSpace(selector)
			if trimmed == "" {
				continue
			}
			selection.Find(trimmed).Remove()
		}

		var outer string
		outer, parseErr = goquery.OuterHtml(selection)
		if parseErr != nil {
			panic(parseErr)
		}
		sanitized = outer
	})
	if err != nil {
		return "", fmt.Errorf("denyElementsの適用に失敗しました: %w", err)
	}

	return sanitized, nil
}
