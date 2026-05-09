package testsupport

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// HandlerResponse はテストサーバーのレスポンス定義です。
type HandlerResponse struct {
	Status  int
	Body    string
	Headers map[string]string
}

// TestServer は in-memory transport のテストサーバーです。
type TestServer struct {
	URL    string
	client *http.Client
}

// Close は互換のための no-op です。
func (s *TestServer) Close() {
}

// Client はテストサーバー向け HTTP クライアントを返します。
func (s *TestServer) Client() *http.Client {
	return s.client
}

type mockTransport struct {
	handler http.Handler
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	recorder := httptest.NewRecorder()
	m.handler.ServeHTTP(recorder, req)
	return recorder.Result(), nil
}

// NewForgejoTestServer は Forgejo API のモックサーバーを返します。
func NewForgejoTestServer(paths map[string]HandlerResponse) (*TestServer, func(string) bool) {
	normalizeRequestPath := func(path, rawQuery string) string {
		path = strings.TrimSuffix(path, "/")
		if rawQuery == "" {
			return path
		}
		values, err := url.ParseQuery(rawQuery)
		if err != nil {
			return path + "?" + rawQuery
		}
		return path + "?" + values.Encode()
	}

	normalizeRequestKey := func(rawKey string) string {
		parts := strings.SplitN(strings.TrimSpace(rawKey), " ", 2)
		if len(parts) != 2 {
			return rawKey
		}
		method := parts[0]
		pathAndQuery := parts[1]
		split := strings.SplitN(pathAndQuery, "?", 2)
		path := split[0]
		rawQuery := ""
		if len(split) == 2 {
			rawQuery = split[1]
		}
		return fmt.Sprintf("%s %s", method, normalizeRequestPath(path, rawQuery))
	}

	normalizedPaths := make(map[string]HandlerResponse, len(paths))
	for key, response := range paths {
		normalizedPaths[normalizeRequestKey(key)] = response
	}

	pathStates := struct {
		mu     sync.Mutex
		counts map[string]int
	}{
		counts: map[string]int{},
	}

	recordPathState := func(key string) {
		pathStates.mu.Lock()
		pathStates.counts[key]++
		pathStates.mu.Unlock()
	}

	buildKey := func(r *http.Request) string {
		return fmt.Sprintf("%s %s", r.Method, normalizeRequestPath(r.URL.Path, r.URL.RawQuery))
	}
	buildPathOnlyKey := func(r *http.Request) string {
		return fmt.Sprintf("%s %s", r.Method, strings.TrimSuffix(r.URL.Path, "/"))
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := buildKey(r)
		recordPathState(key)
		pathOnlyKey := buildPathOnlyKey(r)
		recordPathState(pathOnlyKey)

		if response, ok := normalizedPaths[key]; ok {
			for headerKey, headerValue := range response.Headers {
				w.Header().Set(headerKey, headerValue)
			}
			w.WriteHeader(response.Status)
			_, _ = w.Write([]byte(response.Body))
			return
		}
		if response, ok := normalizedPaths[pathOnlyKey]; ok {
			for headerKey, headerValue := range response.Headers {
				w.Header().Set(headerKey, headerValue)
			}
			w.WriteHeader(response.Status)
			_, _ = w.Write([]byte(response.Body))
			return
		}

		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	})

	called := func(path string) bool {
		pathStates.mu.Lock()
		defer pathStates.mu.Unlock()
		_, ok := pathStates.counts[normalizeRequestKey(path)]
		return ok
	}

	return &TestServer{
		URL: "http://forgejo.local",
		client: &http.Client{
			Transport: &mockTransport{
				handler: handler,
			},
		},
	}, called
}

// NewForgejoTestServerWithRequestDelay は遅延付きモックサーバーを返します。
func NewForgejoTestServerWithRequestDelay(paths map[string]HandlerResponse, activeCount, maxActiveCount *int64, requestDelay time.Duration) (*TestServer, func(string) bool) {
	server, called := NewForgejoTestServer(paths)
	innerHandler := server.client.Transport.(*mockTransport).handler
	server.client.Transport = &mockTransport{
		handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentActive := atomic.AddInt64(activeCount, 1)
			for {
				maxActive := atomic.LoadInt64(maxActiveCount)
				if currentActive <= maxActive {
					break
				}
				if atomic.CompareAndSwapInt64(maxActiveCount, maxActive, currentActive) {
					break
				}
			}
			defer atomic.AddInt64(activeCount, -1)

			if requestDelay > 0 {
				time.Sleep(requestDelay)
			}
			innerHandler.ServeHTTP(w, r)
		}),
	}
	return server, called
}
