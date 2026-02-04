package fetchers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	sanitizer "github.com/landmaster135/devbox/internal/html_sanitizer/usecases/sanitizer"

	domain "github.com/landmaster135/devbox/internal/web_scraper/domain"
	rodFetcher "github.com/landmaster135/devbox/internal/web_scraper/interfaces/fetchers"
)

const defaultLaunchTimeout = 90 * time.Second

// DOMFetcher はgithub.com/go-rod/rodを利用したDOM取得実装です。
type DOMFetcher struct {
	fetchInfra    *rodFetcher.RodDOMFetcher
}

// Option はRodDOMFetcherのオプションです。
type Option func(*DOMFetcher)

// NewRodDOMFetcher はRodDOMFetcherのインスタンスを生成します。
func NewRodDOMFetcher(opts ...Option) *DOMFetcher {
	f := rodFetcher.NewRodDOMFetcher(defaultLaunchTimeout)
	fetcher := &DOMFetcher{
		fetchInfra:    f,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(fetcher)
		}
	}
	return fetcher
}

// FetchDOM は対象URLのDOMツリーを取得します。
func (f *DOMFetcher) FetchDOM(ctx context.Context, targetURL string, wait time.Duration) (domain.DOMResult, error) {
	html, resolvedURL, pageTitle, err := f.fetchInfra.FetchPage(ctx, targetURL, wait)
	if err != nil {
		return domain.DOMResult{}, err
	}

	sanitized, err := sanitizer.SanitizeHTMLBody(html, true)
	if err != nil {
		var notFoundErr *sanitizer.MainNotFoundError
		if errors.As(err, &notFoundErr) {
			return domain.DOMResult{}, errors.New(notFoundErr.Error())
		}
		return domain.DOMResult{}, fmt.Errorf("HTMLのサニタイズに失敗しました: %w", err)
	}

	return domain.DOMResult{
		HTML:     sanitized,
		FinalURL: strings.TrimSpace(resolvedURL),
		Title:    strings.TrimSpace(pageTitle),
	}, nil
}

// FetchMetadata は対象URLのメタ情報を取得します（DOMのサニタイズは行いません）。
func (f *DOMFetcher) FetchMetadata(ctx context.Context, targetURL string, wait time.Duration) (*domain.MetaProps, error) {
	_, resolvedURL, pageTitle, err := f.fetchInfra.FetchPage(ctx, targetURL, wait)
	if err != nil {
		return nil, err
	}

	return &domain.MetaProps{
		URL:   strings.TrimSpace(resolvedURL),
		Title: strings.TrimSpace(pageTitle),
	}, nil
}
