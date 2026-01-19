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

// DOMWriter は取得したDOMを永続化するためのインターフェースです。
type DOMWriter interface {
	Write(path string, content string) error
}

// DOMService はDOM取得ユースケースを提供します。
type DOMService struct {
	fetcher DOMFetcher
	writer  DOMWriter
}

// NewDOMService はDOMServiceを生成します。
func NewDOMService(fetcher DOMFetcher, writer DOMWriter) *DOMService {
	return &DOMService{fetcher: fetcher, writer: writer}
}

// GetDOMTree は対象URLのDOMツリーを取得します。
// outputPath が空でない場合はwriterで保存し、保存結果をboolで返します。
func (s *DOMService) GetDOMTree(
	ctx context.Context,
	url string,
	wait time.Duration,
	outputPath string,
) (string, bool, error) {
	if s == nil || s.fetcher == nil {
		return "", false, fmt.Errorf("DOMFetcherが初期化されていません")
	}
	if strings.TrimSpace(url) == "" {
		return "", false, fmt.Errorf("urlを指定してください")
	}

	html, err := s.fetcher.FetchDOM(ctx, url, wait)
	if err != nil {
		return "", false, fmt.Errorf("DOMツリーの取得に失敗しました: %w", err)
	}

	path := strings.TrimSpace(outputPath)
	if path == "" {
		return html, false, nil
	}
	if s.writer == nil {
		return "", false, fmt.Errorf("DOMWriterが初期化されていません")
	}
	if err := s.writer.Write(path, html); err != nil {
		return "", false, fmt.Errorf("DOMの書き込みに失敗しました: %w", err)
	}

	return html, true, nil
}
