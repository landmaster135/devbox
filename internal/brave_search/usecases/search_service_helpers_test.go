package usecases

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/stretchr/testify/mock"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#
// MockHTTPClient はHTTPクライアントのモック実装
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// MockEnvironmentReader は環境変数読み取りのモック実装
type MockEnvironmentReader struct {
	mock.Mock
}

func (m *MockEnvironmentReader) Getenv(key string) string {
	args := m.Called(key)
	return args.String(0)
}

// MockRateLimiter はレート制限のモック実装
type MockRateLimiter struct {
	mock.Mock
}

func (m *MockRateLimiter) CheckLimit() error {
	args := m.Called()
	return args.Error(0)
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
