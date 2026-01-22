package fetchers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/landmaster135/devbox/internal/web_scraper/usecases"
)

const defaultLaunchTimeout = 90 * time.Second

// RodDOMFetcher はgithub.com/go-rod/rodを利用したDOM取得実装です。
type RodDOMFetcher struct {
	launchTimeout time.Duration
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
func (f *RodDOMFetcher) FetchDOM(ctx context.Context, targetURL string, wait time.Duration) (usecases.DOMResult, error) {
	html, resolvedURL, pageTitle, err := f.fetchPage(ctx, targetURL, wait)
	if err != nil {
		return usecases.DOMResult{}, err
	}

	sanitized, err := sanitizeHTMLBody(html, true)
	if err != nil {
		var notFoundErr *MainNotFoundError
		if errors.As(err, &notFoundErr) {
			return usecases.DOMResult{}, errors.New(notFoundErr.Error())
		}
		return usecases.DOMResult{}, fmt.Errorf("HTMLのサニタイズに失敗しました: %w", err)
	}

	return usecases.DOMResult{
		HTML:     sanitized,
		FinalURL: strings.TrimSpace(resolvedURL),
		Title:    strings.TrimSpace(pageTitle),
	}, nil
}

// FetchMetadata は対象URLのメタ情報を取得します（DOMのサニタイズは行いません）。
func (f *RodDOMFetcher) FetchMetadata(ctx context.Context, targetURL string, wait time.Duration) (*usecases.MetaProps, error) {
	_, resolvedURL, pageTitle, err := f.fetchPage(ctx, targetURL, wait)
	if err != nil {
		return nil, err
	}

	return &usecases.MetaProps{
		URL:   strings.TrimSpace(resolvedURL),
		Title: strings.TrimSpace(pageTitle),
	}, nil
}

func (f *RodDOMFetcher) fetchPage(ctx context.Context, targetURL string, wait time.Duration) (string, string, string, error) {
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
