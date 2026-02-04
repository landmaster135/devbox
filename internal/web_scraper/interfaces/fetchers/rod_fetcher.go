package fetchers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

const defaultLaunchTimeout = 90 * time.Second

// RodDOMFetcher はgithub.com/go-rod/rodを利用したDOM取得実装です。
type RodDOMFetcher struct {
	launchTimeout time.Duration
}

// NewRodDOMFetcher はRodDOMFetcherのインスタンスを生成します。
func NewRodDOMFetcher(timeout time.Duration) *RodDOMFetcher {
	if timeout == 0 {
		timeout = defaultLaunchTimeout
	}
	return &RodDOMFetcher{
		launchTimeout: timeout,
	}
}

func (f *RodDOMFetcher) FetchPage(ctx context.Context, targetURL string, wait time.Duration) (string, string, string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		html        string
		pageTitle   string
		resolvedURL string
	)

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

		pageTitle = page.MustEval(`() => document.title || ""`).Str()
		resolvedURL = page.MustEval(`() => window.location.href`).Str()
		html = page.MustEval(`() => document.documentElement.outerHTML`).Str()
	})
	if err != nil {
		return "", "", "", fmt.Errorf("DOMツリーの取得に失敗しました: %w", err)
	}

	return html, resolvedURL, pageTitle, nil
}
