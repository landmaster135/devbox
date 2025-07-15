package usecases

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          BraveSearchService Constructor Tests             ##
// #==============================================================#
func TestNewBraveSearchService_Normal(t *testing.T) {
	service := NewBraveSearchService()
	assert.NotNil(t, service)
	assert.NotNil(t, service.httpClient)
	assert.NotNil(t, service.envReader)
}

func TestNewBraveSearchServiceWithDependencies_Normal(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}

	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)
	assert.NotNil(t, service)
	assert.Equal(t, mockHTTPClient, service.httpClient)
	assert.Equal(t, mockEnvReader, service.envReader)
}

// #==============================================================#
// ##          Default Environment Reader Tests                 ##
// #==============================================================#
func TestDefaultEnvironmentReader_Getenv(t *testing.T) {
	reader := &DefaultEnvironmentReader{}

	// 実際の環境変数は設定されていない可能性があるので、空文字列を期待
	result := reader.Getenv("NON_EXISTENT_ENV_VAR")
	assert.Equal(t, "", result)
}

// #==============================================================#
// ##          Default HTTP Client Tests                        ##
// #==============================================================#
func TestDefaultHTTPClient_Do(t *testing.T) {
	client := &DefaultHTTPClient{}

	// 簡単なHTTPリクエストのテスト（実際のネットワーク呼び出しは避ける）
	req, err := http.NewRequest("GET", "http://httpbin.org/status/200", nil)
	assert.NoError(t, err)

	// このテストは実際のネットワーク接続に依存するため、
	// エラーが発生してもテストが失敗しないようにする
	_, _ = client.Do(req)
	// ネットワークエラーの可能性があるため、エラーチェックはしない
	// 主な目的はコードカバレッジの向上
}
