package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/landmaster135/devbox/internal/web_scraper/domain"
)

type mockDOMFetcher struct {
	t         *testing.T
	wantURL   string
	wantWait  time.Duration
	result    domain.DOMResult
	err       error
	callCount int

	metaWantURL  string
	metaWantWait time.Duration
	metaResult   *domain.MetaProps
	metaErr      error
	metaCalls    int
}

func (m *mockDOMFetcher) FetchDOM(ctx context.Context, url string, wait time.Duration) (domain.DOMResult, error) {
	m.callCount++
	if m.wantURL != "" && url != m.wantURL {
		m.t.Fatalf("expected url %s, got %s", m.wantURL, url)
	}
	if m.wantWait != wait {
		m.t.Fatalf("expected wait %s, got %s", m.wantWait, wait)
	}
	if ctx == nil {
		m.t.Fatalf("context is nil")
	}
	return m.result, m.err
}

func (m *mockDOMFetcher) FetchMetadata(ctx context.Context, url string, wait time.Duration) (*domain.MetaProps, error) {
	m.metaCalls++
	if m.metaWantURL != "" && url != m.metaWantURL {
		m.t.Fatalf("expected meta url %s, got %s", m.metaWantURL, url)
	}
	if m.metaWantWait != wait {
		m.t.Fatalf("expected meta wait %s, got %s", m.metaWantWait, wait)
	}
	if ctx == nil {
		m.t.Fatalf("context is nil")
	}
	return m.metaResult, m.metaErr
}

type mockDOMWriter struct {
	t         *testing.T
	wantPath  string
	wantValue string
	err       error
	callCount int
}

func (m *mockDOMWriter) Write(path string, content string) error {
	m.callCount++
	if m.wantPath != "" && path != m.wantPath {
		m.t.Fatalf("expected path %s, got %s", m.wantPath, path)
	}
	if m.wantValue != "" && content != m.wantValue {
		m.t.Fatalf("expected content %s, got %s", m.wantValue, content)
	}
	return m.err
}

func TestDOMService_GetDOMTree_Success(t *testing.T) {
	t.Parallel()

	fetcher := &mockDOMFetcher{
		t:        t,
		wantURL:  "https://example.com",
		wantWait: 2 * time.Second,
		result:   domain.DOMResult{HTML: "<html></html>"},
	}

	service := NewDOMService(fetcher, nil)
	got, written, err := service.GetDOMTree(context.Background(), "https://example.com", 2*time.Second, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != fetcher.result.HTML {
		t.Fatalf("expected %q, got %q", fetcher.result.HTML, got)
	}
	if fetcher.callCount != 1 {
		t.Fatalf("expected fetcher to be called once, got %d", fetcher.callCount)
	}
	if written {
		t.Fatalf("expected written to be false")
	}
}

func TestDOMService_GetDOMTree_WriteSuccess(t *testing.T) {
	t.Parallel()

	fetcher := &mockDOMFetcher{
		t:       t,
		wantURL: "https://example.io",
		result:  domain.DOMResult{HTML: "<main>ok</main>"},
	}
	writer := &mockDOMWriter{
		t:         t,
		wantPath:  "out.html",
		wantValue: "<main>ok</main>",
	}

	service := NewDOMService(fetcher, writer)
	_, written, err := service.GetDOMTree(context.Background(), "https://example.io", 0, "out.html")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !written {
		t.Fatalf("expected written to be true")
	}
	if writer.callCount != 1 {
		t.Fatalf("expected writer to be called once, got %d", writer.callCount)
	}
}

func TestDOMService_GetDOMTree_WriteError(t *testing.T) {
	t.Parallel()

	fetcher := &mockDOMFetcher{t: t, result: domain.DOMResult{HTML: "<main></main>"}}
	writer := &mockDOMWriter{t: t, err: errors.New("fail")}
	service := NewDOMService(fetcher, writer)

	if _, _, err := service.GetDOMTree(context.Background(), "https://example.com", 0, "out.html"); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDOMService_GetDOMTree_WriterMissing(t *testing.T) {
	t.Parallel()

	fetcher := &mockDOMFetcher{t: t, result: domain.DOMResult{HTML: "<main></main>"}}
	service := NewDOMService(fetcher, nil)
	if _, _, err := service.GetDOMTree(context.Background(), "https://example.com", 0, "out.html"); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDOMService_GetDOMTree_FetcherError(t *testing.T) {
	t.Parallel()

	fetcher := &mockDOMFetcher{
		t:       t,
		wantURL: "https://example.org",
		err:     errors.New("boom"),
	}

	service := NewDOMService(fetcher, nil)
	if _, _, err := service.GetDOMTree(context.Background(), "https://example.org", 0, ""); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDOMService_GetDOMTree_EmptyURL(t *testing.T) {
	t.Parallel()

	service := NewDOMService(&mockDOMFetcher{t: t}, nil)
	if _, _, err := service.GetDOMTree(context.Background(), " ", 0, ""); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDOMService_GetDOMTree_FetcherMissing(t *testing.T) {
	t.Parallel()

	var service *DOMService
	if _, _, err := service.GetDOMTree(context.Background(), "https://example.com", 0, ""); err == nil {
		t.Fatalf("expected error, got nil")
	}

	service = &DOMService{}
	if _, _, err := service.GetDOMTree(context.Background(), "https://example.com", 0, ""); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDOMService_GetMetaProps_Success(t *testing.T) {
	t.Parallel()

	fetcher := &mockDOMFetcher{
		t:            t,
		metaWantURL:  "https://example.com",
		metaWantWait: time.Second,
		metaResult: &domain.MetaProps{
			URL:   "https://example.net/page",
			Title: " Example Title ",
		},
	}
	service := NewDOMService(fetcher, nil)
	props, err := service.GetMetaProps(context.Background(), "https://example.com", time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if props.URL != "https://example.net/page" {
		t.Fatalf("expected final url to be https://example.net/page, got %s", props.URL)
	}
	if props.Title != "Example Title" {
		t.Fatalf("expected title to be trimmed, got %q", props.Title)
	}
}

func TestDOMService_GetMetaProps_FallbackURL(t *testing.T) {
	t.Parallel()

	fetcher := &mockDOMFetcher{
		t: t,
		metaResult: &domain.MetaProps{
			URL:   "",
			Title: "Sample",
		},
	}
	service := NewDOMService(fetcher, nil)
	props, err := service.GetMetaProps(context.Background(), " https://fallback.example ", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if props.URL != "https://fallback.example" {
		t.Fatalf("expected fallback url, got %s", props.URL)
	}
}

func TestDOMService_GetMetaProps_FetcherError(t *testing.T) {
	t.Parallel()

	fetcher := &mockDOMFetcher{
		t:       t,
		metaErr: errors.New("failed"),
	}
	service := NewDOMService(fetcher, nil)
	if _, err := service.GetMetaProps(context.Background(), "https://example.com", 0); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDOMService_GetMetaProps_FetcherMissing(t *testing.T) {
	t.Parallel()

	var service *DOMService
	if _, err := service.GetMetaProps(context.Background(), "https://example.com", 0); err == nil {
		t.Fatalf("expected error, got nil")
	}

	service = &DOMService{}
	if _, err := service.GetMetaProps(context.Background(), "https://example.com", 0); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
