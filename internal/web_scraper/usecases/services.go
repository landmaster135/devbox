package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// DOMFetcher は指定したURLのDOMツリーを取得するためのインターフェースです。
type DOMFetcher interface {
	FetchDOM(ctx context.Context, url string, wait time.Duration) (string, error)
}

// DOMService はDOM取得ユースケースを提供します。
type DOMService struct {
	fetcher DOMFetcher
}

// NewDOMService はDOMServiceを生成します。
func NewDOMService(fetcher DOMFetcher) *DOMService {
	return &DOMService{fetcher: fetcher}
}

// GetDOMTree は対象URLのDOMツリーを取得します。
func (s *DOMService) GetDOMTree(ctx context.Context, url string, wait time.Duration) (string, error) {
	if s == nil || s.fetcher == nil {
		return "", fmt.Errorf("DOMFetcherが初期化されていません")
	}
	if strings.TrimSpace(url) == "" {
		return "", fmt.Errorf("urlを指定してください")
	}

	html, err := s.fetcher.FetchDOM(ctx, url, wait)
	if err != nil {
		return "", fmt.Errorf("DOMツリーの取得に失敗しました: %w", err)
	}

	return html, nil
}
