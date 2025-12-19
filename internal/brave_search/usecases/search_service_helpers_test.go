package usecases

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#
// MockHTTPClient はHTTPクライアントのモック実装
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// MockEnvironmentReader は環境変数読み取りのモック実装
type MockEnvironmentReader struct {
	GetenvFunc func(key string) string
}

func (m *MockEnvironmentReader) Getenv(key string) string {
	return m.GetenvFunc(key)
}

// MockRateLimiter はレート制限のモック実装
type MockRateLimiter struct {
	CheckLimitFunc func() error
}

func (m *MockRateLimiter) CheckLimit() error {
	return m.CheckLimitFunc()
}

// #==============================================================#
// ##          Helper Functions                                  ##
// #==============================================================#
// createMockResponse はモックHTTPレスポンスを作成するヘルパー関数
func createMockResponse(statusCode int, body interface{}) *http.Response {
	var bodyBytes []byte
	var err error

	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			panic(err)
		}
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(bodyBytes)),
		Header:     make(http.Header),
	}
}
