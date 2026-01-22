package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"
	
	"github.com/landmaster135/devbox/internal/web_scraper/domain"
	"github.com/landmaster135/devbox/internal/web_scraper/usecases/fetcher"
	"github.com/landmaster135/devbox/internal/web_scraper/interfaces/writers"
)

// DOMService はDOM取得ユースケースを提供します。
type DOMService struct {
	fetcher domain.DOMFetcher
	writer  domain.DOMWriter
}

// NewDOMService はDOMServiceを生成します。
func NewDOMService(fetcher domain.DOMFetcher, writer domain.DOMWriter) *DOMService {
	return &DOMService{fetcher: fetcher, writer: writer}
}

// NewDefaultDOMService は標準実装のfetcher/writerを組み込んだDOMServiceを生成します。
func NewDefaultDOMService() *DOMService {
	return NewDOMService(
		fetchers.NewRodDOMFetcher(),
		writers.NewFileWriter(),
	)
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

	result, err := s.fetcher.FetchDOM(ctx, url, wait)
	if err != nil {
		return "", false, fmt.Errorf("DOMツリーの取得に失敗しました: %w", err)
	}
	html := result.HTML

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

// GetMetaProps は対象URLのタイトルと最終URLを取得します。
func (s *DOMService) GetMetaProps(
	ctx context.Context,
	url string,
	wait time.Duration,
) (*domain.MetaProps, error) {
	if s == nil || s.fetcher == nil {
		return nil, fmt.Errorf("DOMFetcherが初期化されていません")
	}
	trimmedURL := strings.TrimSpace(url)
	if trimmedURL == "" {
		return nil, fmt.Errorf("urlを指定してください")
	}

	result, err := s.fetcher.FetchMetadata(ctx, trimmedURL, wait)
	if err != nil {
		return nil, fmt.Errorf("メタ情報の取得に失敗しました: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("メタ情報が取得できませんでした")
	}

	finalURL := strings.TrimSpace(result.URL)
	if finalURL == "" {
		finalURL = trimmedURL
	}

	return &domain.MetaProps{
		URL:   finalURL,
		Title: strings.TrimSpace(result.Title),
	}, nil
}
