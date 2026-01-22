package domain

import (
	"context"
	"time"
)

// DOMResult はDOM取得時の結果を保持します。
type DOMResult struct {
	HTML     string
	FinalURL string
	Title    string
}

// MetaProps はget_meta_propsが返すメタ情報です。
type MetaProps struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

// DOMFetcher は対象ページの取得とメタ情報取得を担います。
type DOMFetcher interface {
	FetchDOM(ctx context.Context, url string, wait time.Duration) (DOMResult, error)
	FetchMetadata(ctx context.Context, url string, wait time.Duration) (*MetaProps, error)
}

// DOMWriter は取得したDOMを永続化するためのインターフェースです。
type DOMWriter interface {
	Write(path string, content string) error
}
