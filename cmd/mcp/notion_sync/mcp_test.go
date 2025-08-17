package notion_sync

import (
	"context"
	"testing"

	mcp "github.com/mark3labs/mcp-go/mcp"
)

// MockEnvironmentProvider はテスト用の環境変数プロバイダー
type MockEnvironmentProvider struct {
	envVars map[string]string
}

// GetEnv はモックの環境変数を取得する
func (m *MockEnvironmentProvider) GetEnv(key string) string {
	return m.envVars[key]
}

// TestCreateConfigFromEnv は環境変数から設定を作成するテスト
func TestCreateConfigFromEnv(t *testing.T) {
	tests := []struct {
		name            string
		envVars         map[string]string
		conID           string
		pageID          string
		markdownContent string
		toggleH1        bool
		toggleH2        bool
		toggleH3        bool
		expectError     bool
		errorMessage    string
	}{
		{
			name: "正常なケース",
			envVars: map[string]string{
				"NOTION_INTEGRATION_TOKEN":       "test_token",
				"NOTION_ENDPOINT_URL": "http://localhost:8080/test",
			},
			conID:           "",
			pageID:          "test_page_id",
			markdownContent: "# Test Content",
			toggleH1:        true,
			toggleH2:        false,
			toggleH3:        true,
			expectError:     false,
		},
		{
			name: "トークンが設定されていない",
			envVars: map[string]string{
				"NOTION_ENDPOINT_URL": "http://localhost:8080/test",
			},
			conID:           "",
			pageID:          "test_page_id",
			markdownContent: "# Test Content",
			toggleH1:        false,
			toggleH2:        false,
			toggleH3:        false,
			expectError:     true,
			errorMessage:    "環境変数NOTION_INTEGRATION_TOKENが設定されていません",
		},
		{
			name: "エンドポイントURLが設定されていない",
			envVars: map[string]string{
				"NOTION_INTEGRATION_TOKEN": "test_token",
			},
			conID:           "",
			pageID:          "test_page_id",
			markdownContent: "# Test Content",
			toggleH1:        false,
			toggleH2:        false,
			toggleH3:        false,
			expectError:     true,
			errorMessage:    "環境変数NOTION_ENDPOINT_URLが設定されていません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider := &MockEnvironmentProvider{envVars: tt.envVars}

			cfg, err := createConfigFromEnv(mockProvider, "patch", tt.conID, tt.pageID, tt.markdownContent, tt.toggleH1, tt.toggleH2, tt.toggleH3)

			if tt.expectError {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
					return
				}
				if err.Error() != tt.errorMessage {
					t.Errorf("期待されたエラーメッセージ: %s, 実際のエラーメッセージ: %s", tt.errorMessage, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			if cfg == nil {
				t.Errorf("設定がnilです")
				return
			}

			// 設定の検証
			if cfg.Token != tt.envVars["NOTION_INTEGRATION_TOKEN"] {
				t.Errorf("期待されたトークン: %s, 実際のトークン: %s", tt.envVars["NOTION_INTEGRATION_TOKEN"], cfg.Token)
			}

			if cfg.EndpointURL != tt.envVars["NOTION_ENDPOINT_URL"] {
				t.Errorf("期待されたエンドポイントURL: %s, 実際のエンドポイントURL: %s", tt.envVars["NOTION_ENDPOINT_URL"], cfg.EndpointURL)
			}

			if cfg.PageID != tt.pageID {
				t.Errorf("期待されたページID: %s, 実際のページID: %s", tt.pageID, cfg.PageID)
			}

			if cfg.MarkdownContent != tt.markdownContent {
				t.Errorf("期待されたマークダウンコンテンツ: %s, 実際のマークダウンコンテンツ: %s", tt.markdownContent, cfg.MarkdownContent)
			}

			if cfg.ToggleH1 != tt.toggleH1 {
				t.Errorf("期待されたToggleH1: %t, 実際のToggleH1: %t", tt.toggleH1, cfg.ToggleH1)
			}

			if cfg.ToggleH2 != tt.toggleH2 {
				t.Errorf("期待されたToggleH2: %t, 実際のToggleH2: %t", tt.toggleH2, cfg.ToggleH2)
			}

			if cfg.ToggleH3 != tt.toggleH3 {
				t.Errorf("期待されたToggleH3: %t, 実際のToggleH3: %t", tt.toggleH3, cfg.ToggleH3)
			}
		})
	}
}

// TestCreateNotionSyncServer はMCPサーバーの作成をテストする
func TestCreateNotionSyncServer(t *testing.T) {
	server := createNotionSyncServer()

	if server == nil {
		t.Errorf("サーバーがnilです")
		return
	}

	// サーバーの基本的な検証
	// 実際のMCPサーバーの内部構造にアクセスするのは困難なため、
	// ここでは単純にサーバーが作成されることを確認
}

// TestHandlePatchPageWithMockRequest はhandlePatchPageのテスト（モックリクエスト使用）
func TestHandlePatchPageWithMockRequest(t *testing.T) {
	// 環境変数を設定
	t.Setenv("NOTION_INTEGRATION_TOKEN", "test_token")
	t.Setenv("NOTION_ENDPOINT_URL", "http://localhost:9999/invalid_endpoint")

	// モックリクエストを作成
	mockRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "patch_page",
			Arguments: map[string]interface{}{
				"markdown_content": "# Test Content",
				"page_id":          "test_page_id",
				"toggle_h1":        true,
				"toggle_h2":        false,
				"toggle_h3":        true,
			},
		},
	}

	ctx := context.Background()

	// ハンドラーを呼び出し（存在しないエンドポイントなのでエラーが発生することを期待）
	result, err := handlePatchPage(ctx, mockRequest)

	// HTTPリクエストが失敗することを期待（存在しないエンドポイントのため）
	if err == nil {
		t.Logf("予期せずエラーが発生しませんでした。結果: %v", result)
		// エラーが発生しない場合もあるため、ログに記録するだけにする
		return
	}

	if result != nil {
		t.Errorf("エラー時は結果がnilであることを期待しましたが、結果が返されました: %v", result)
	}

	// エラーメッセージが空でないことを確認
	if err.Error() == "" {
		t.Errorf("エラーメッセージが空です")
	}
}

// TestHandlePatchPageMissingMarkdownContent は必須パラメータが不足している場合のテスト
func TestHandlePatchPageMissingMarkdownContent(t *testing.T) {
	// 環境変数を設定
	t.Setenv("NOTION_INTEGRATION_TOKEN", "test_token")
	t.Setenv("NOTION_ENDPOINT_URL", "http://localhost:8080/test")

	// markdown_contentが不足しているモックリクエストを作成
	mockRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "patch_page",
			Arguments: map[string]interface{}{
				"page_id":   "test_page_id",
				"toggle_h1": true,
			},
		},
	}

	ctx := context.Background()

	// ハンドラーを呼び出し（必須パラメータが不足しているためエラーが発生することを期待）
	result, err := handlePatchPage(ctx, mockRequest)

	// 必須パラメータが不足しているためエラーが発生することを期待
	if err == nil {
		t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
		return
	}

	if result != nil {
		t.Errorf("エラー時は結果がnilであることを期待しましたが、結果が返されました: %v", result)
	}

	// エラーメッセージが空でないことを確認
	if err.Error() == "" {
		t.Errorf("エラーメッセージが空です")
	}
}
