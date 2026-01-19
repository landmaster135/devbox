package usecases

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockDOMFetcher struct {
	t         *testing.T
	wantURL   string
	wantWait  time.Duration
	result    string
	err       error
	callCount int
}

func (m *mockDOMFetcher) FetchDOM(ctx context.Context, url string, wait time.Duration) (string, error) {
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
		result:   "<html></html>",
	}

	service := NewDOMService(fetcher, nil)
	got, written, err := service.GetDOMTree(context.Background(), "https://example.com", 2*time.Second, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != fetcher.result {
		t.Fatalf("expected %q, got %q", fetcher.result, got)
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
		result:  "<main>ok</main>",
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

	fetcher := &mockDOMFetcher{t: t, result: "<main></main>"}
	writer := &mockDOMWriter{t: t, err: errors.New("fail")}
	service := NewDOMService(fetcher, writer)

	if _, _, err := service.GetDOMTree(context.Background(), "https://example.com", 0, "out.html"); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDOMService_GetDOMTree_WriterMissing(t *testing.T) {
	t.Parallel()

	fetcher := &mockDOMFetcher{t: t, result: "<main></main>"}
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
