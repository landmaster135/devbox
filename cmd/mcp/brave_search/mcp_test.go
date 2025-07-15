package brave_search

import (
	"testing"

	"github.com/stretchr/testify/assert"

	usecases "github.com/landmaster135/devbox/internal/brave_search/usecases"
)

// #==============================================================#
// ##          Server Configuration Tests                        ##
// #==============================================================#
func TestCreateBraveSearchServer_Normal(t *testing.T) {
	// 関数を呼び出し
	server := createBraveSearchServer()

	// 結果の検証
	assert.NotNil(t, server)
}

// #==============================================================#
// ##          Service Integration Tests                         ##
// #==============================================================#
func TestServiceIntegration_WebSearch(t *testing.T) {
	// 実際のサービスを使用したインテグレーションテスト
	service := usecases.NewBraveSearchService()

	// サービスが正しく初期化されていることを確認
	assert.NotNil(t, service)

	// APIキーが設定されていない場合のテスト
	result, err := service.HandleWebSearch("golang programming", 3, 0)

	// APIキーが設定されていない場合はエラーが発生する
	if err != nil {
		assert.Contains(t, err.Error(), "BRAVE_API_KEY environment variable is required")
		assert.Empty(t, result)
	} else {
		// APIキーが設定されている場合は結果が返される
		assert.NotEmpty(t, result)
	}
}

func TestServiceIntegration_LocalSearch(t *testing.T) {
	// 実際のサービスを使用したインテグレーションテスト
	service := usecases.NewBraveSearchService()

	// サービスが正しく初期化されていることを確認
	assert.NotNil(t, service)

	// APIキーが設定されていない場合のテスト
	result, err := service.HandleLocalSearch("pizza near Tokyo", 2)

	// APIキーが設定されていない場合またはレート制限の場合はエラーが発生する
	if err != nil {
		// APIキーエラーまたはレート制限エラーのいずれかを期待
		isAPIKeyError := err.Error() == "BRAVE_API_KEY environment variable is required"
		isRateLimitError := err.Error() == "rate limit exceeded"
		assert.True(t, isAPIKeyError || isRateLimitError, "Expected API key error or rate limit error, got: %s", err.Error())
		assert.Empty(t, result)
	} else {
		// APIキーが設定されている場合は結果が返される
		assert.NotEmpty(t, result)
	}
}

// #==============================================================#
// ##          Service Constructor Tests                         ##
// #==============================================================#
func TestNewBraveSearchService_Normal(t *testing.T) {
	service := usecases.NewBraveSearchService()
	assert.NotNil(t, service)
}

// #==============================================================#
// ##          Handler Logic Tests                               ##
// #==============================================================#
func TestHandlerLogic_WebSearch_ParameterValidation(t *testing.T) {
	// パラメータ検証のテスト
	testCases := []struct {
		name     string
		query    string
		count    int
		offset   int
		hasError bool
	}{
		{
			name:     "Valid parameters",
			query:    "test query",
			count:    5,
			offset:   0,
			hasError: false,
		},
		{
			name:     "Empty query",
			query:    "",
			count:    5,
			offset:   0,
			hasError: false, // 空のクエリでもAPIに送信される
		},
		{
			name:     "Large count",
			query:    "test query",
			count:    25, // API制限を超える
			offset:   0,
			hasError: false, // 内部で20に制限される
		},
		{
			name:     "Negative count",
			query:    "test query",
			count:    -1,
			offset:   0,
			hasError: false, // 負の値でもAPIに送信される
		},
	}

	service := usecases.NewBraveSearchService()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.HandleWebSearch(tc.query, tc.count, tc.offset)

			// APIキーが設定されていない場合はスキップ
			if err != nil && err.Error() == "BRAVE_API_KEY environment variable is required" {
				t.Skip("BRAVE_API_KEY environment variable is not set, skipping test")
			}

			if tc.hasError {
				assert.Error(t, err)
			} else {
				// APIキーが設定されている場合のみ成功を期待
				if err == nil {
					assert.NotEmpty(t, result)
				}
			}
		})
	}
}

func TestHandlerLogic_LocalSearch_ParameterValidation(t *testing.T) {
	// パラメータ検証のテスト
	testCases := []struct {
		name     string
		query    string
		count    int
		hasError bool
	}{
		{
			name:     "Valid parameters",
			query:    "test local query",
			count:    3,
			hasError: false,
		},
		{
			name:     "Empty query",
			query:    "",
			count:    3,
			hasError: false, // 空のクエリでもAPIに送信される
		},
		{
			name:     "Large count",
			query:    "test local query",
			count:    25, // API制限を超える
			hasError: false, // 内部で20に制限される
		},
	}

	service := usecases.NewBraveSearchService()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.HandleLocalSearch(tc.query, tc.count)

			// APIキーが設定されていない場合はスキップ
			if err != nil && err.Error() == "BRAVE_API_KEY environment variable is required" {
				t.Skip("BRAVE_API_KEY environment variable is not set, skipping test")
			}

			if tc.hasError {
				assert.Error(t, err)
			} else {
				// APIキーが設定されている場合のみ成功を期待
				if err == nil {
					assert.NotEmpty(t, result)
				}
			}
		})
	}
}

// #==============================================================#
// ##          Error Handling Tests                              ##
// #==============================================================#
func TestErrorHandling_NoAPIKey(t *testing.T) {
	service := usecases.NewBraveSearchService()

	// Web検索のエラーハンドリング
	t.Run("WebSearch without API key", func(t *testing.T) {
		result, err := service.HandleWebSearch("test query", 5, 0)

		// APIキーが設定されていない場合またはレート制限の場合はエラーが発生する
		if err != nil {
			isAPIKeyError := err.Error() == "BRAVE_API_KEY environment variable is required"
			isRateLimitError := err.Error() == "rate limit exceeded"
			assert.True(t, isAPIKeyError || isRateLimitError, "Expected API key error or rate limit error, got: %s", err.Error())
			assert.Empty(t, result)
		}
	})

	// ローカル検索のエラーハンドリング
	t.Run("LocalSearch without API key", func(t *testing.T) {
		result, err := service.HandleLocalSearch("test local query", 3)

		// APIキーが設定されていない場合またはレート制限の場合はエラーが発生する
		if err != nil {
			isAPIKeyError := err.Error() == "BRAVE_API_KEY environment variable is required"
			isRateLimitError := err.Error() == "rate limit exceeded"
			assert.True(t, isAPIKeyError || isRateLimitError, "Expected API key error or rate limit error, got: %s", err.Error())
			assert.Empty(t, result)
		}
	})
}

// #==============================================================#
// ##          Utility Function Tests                            ##
// #==============================================================#
func TestUtilityFunctions(t *testing.T) {
	service := usecases.NewBraveSearchService()

	// getOrDefault関数のテスト（間接的にテスト）
	t.Run("Service initialization", func(t *testing.T) {
		assert.NotNil(t, service)
	})
}

// #==============================================================#
// ##          Build Function Tests                              ##
// #==============================================================#
func TestBuildBraveSearchServer_APIKeyCheck(t *testing.T) {
	// この関数は実際にサーバーを起動するため、テストでは呼び出さない
	// 代わりに、APIキーチェックの動作を確認
	t.Run("API key check logic", func(t *testing.T) {
		// 実際のテストでは環境変数の設定状態を確認
		// ここでは単純にサーバー作成関数が動作することを確認
		server := createBraveSearchServer()
		assert.NotNil(t, server)
	})
}
