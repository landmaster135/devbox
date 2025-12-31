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
	wantDeny  []string
	result    string
	err       error
	callCount int
}

func (m *mockDOMFetcher) FetchDOM(ctx context.Context, url string, wait time.Duration, denySelectors []string) (string, error) {
	m.callCount++
	if m.wantURL != "" && url != m.wantURL {
		m.t.Fatalf("expected url %s, got %s", m.wantURL, url)
	}
	if m.wantWait != wait {
		m.t.Fatalf("expected wait %s, got %s", m.wantWait, wait)
	}
	if m.wantDeny != nil {
		if len(m.wantDeny) != len(denySelectors) {
			m.t.Fatalf("expected denySelectors %v, got %v", m.wantDeny, denySelectors)
		}
		for i, sel := range m.wantDeny {
			if sel != denySelectors[i] {
				m.t.Fatalf("expected denySelectors[%d]=%s, got %s", i, sel, denySelectors[i])
			}
		}
	}
	if ctx == nil {
		m.t.Fatalf("context is nil")
	}
	return m.result, m.err
}

func TestDOMService_GetDOMTree_Success(t *testing.T) {
	t.Parallel()

	fetcher := &mockDOMFetcher{
		t:        t,
		wantURL:  "https://example.com",
		wantWait: 2 * time.Second,
		wantDeny: []string{".ad", "script"},
		result:   "<html></html>",
	}

	service := NewDOMService(fetcher)
	got, err := service.GetDOMTree(context.Background(), "https://example.com", 2*time.Second, []string{".ad", "script"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != fetcher.result {
		t.Fatalf("expected %q, got %q", fetcher.result, got)
	}
	if fetcher.callCount != 1 {
		t.Fatalf("expected fetcher to be called once, got %d", fetcher.callCount)
	}
}

func TestDOMService_GetDOMTree_FetcherError(t *testing.T) {
	t.Parallel()

	fetcher := &mockDOMFetcher{
		t:       t,
		wantURL: "https://example.org",
		err:     errors.New("boom"),
	}

	service := NewDOMService(fetcher)
	if _, err := service.GetDOMTree(context.Background(), "https://example.org", 0, nil); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDOMService_GetDOMTree_EmptyURL(t *testing.T) {
	t.Parallel()

	service := NewDOMService(&mockDOMFetcher{t: t})
	if _, err := service.GetDOMTree(context.Background(), " ", 0, nil); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDOMService_GetDOMTree_FetcherMissing(t *testing.T) {
	t.Parallel()

	var service *DOMService
	if _, err := service.GetDOMTree(context.Background(), "https://example.com", 0, nil); err == nil {
		t.Fatalf("expected error, got nil")
	}

	service = &DOMService{}
	if _, err := service.GetDOMTree(context.Background(), "https://example.com", 0, nil); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
