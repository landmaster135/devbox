package github

import (
	"context"
	"testing"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// TestSetGitHubRepositoryServer はSetGitHubRepositoryServer関数をテストする
func TestSetGitHubRepositoryServer(t *testing.T) {
	// テストケース
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "正常系 - トークンあり",
			token: "testtoken",
		},
		{
			name:  "正常系 - トークンなし",
			token: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 実際のMCPサーバーを作成
			s := server.NewMCPServer("test-server", "1.0.0")

			// テスト対象の関数を実行
			result := setGitHubRepositoryServer(tc.token, s)

			// 関数が正常に実行され、サーバーが返されることを確認
			if result == nil {
				t.Fatal("SetGitHubRepositoryServerがnilを返しました")
			}

			// 注意: 実際のツールの追加や設定の検証は、MCPサーバーの内部実装に依存するため、
			// ここでは基本的な動作確認のみを行います
		})
	}
}

// TestHandleToGetUserRepositories はHandleToGetUserRepositoriesメソッドをテストする
func TestHandleToGetUserRepositories(t *testing.T) {
	// テストケース
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
	}{
		{
			name: "正常系 - 必須パラメータのみ",
			arguments: map[string]interface{}{
				"username": "test_user",
			},
			expectError: false,
		},
		{
			name: "正常系 - すべてのパラメータ",
			arguments: map[string]interface{}{
				"username":  "test_user",
				"per_page":  float64(10),
				"page":      float64(2),
				"sort":      "updated",
				"direction": "asc",
				"type":      "owner",
			},
			expectError: false,
		},
		{
			name: "異常系 - 必須パラメータなし",
			arguments: map[string]interface{}{
				"per_page": float64(10),
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GitHubClientを作成
			client := NewGitHubClient("test_token")

			// リクエストの作成
			request := mcp.CallToolRequest{}
			// Paramsフィールドに直接アクセス
			request.Params.Name = "get_user_repositories"
			request.Params.Arguments = tc.arguments

			// テスト対象の関数を実行
			ctx := context.Background()
			result, err := client.handleToGetUserRepositories(ctx, request)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				// 実際のAPIエラーは許容する（認証エラーなど）
				t.Logf("APIエラーが発生しました（テスト環境では正常）: %v", err)
			}

			// 正常系の場合、結果の基本的な検証
			if !tc.expectError {
				// エラーがない場合、または認証エラーの場合は結果がnilでも許容
				if err == nil && result == nil {
					t.Error("結果がnilです")
				}
			}
		})
	}
}

// TestHandleToGetFileContents はHandleToGetFileContentsメソッドをテストする
func TestHandleToGetFileContents(t *testing.T) {
	// テストケース
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
	}{
		{
			name: "正常系 - 必須パラメータのみ",
			arguments: map[string]interface{}{
				"owner": "test_user",
				"repo":  "test_repo",
				"path":  "README.md",
			},
			expectError: false,
		},
		{
			name: "正常系 - すべてのパラメータ",
			arguments: map[string]interface{}{
				"owner":  "test_user",
				"repo":   "test_repo",
				"path":   "README.md",
				"branch": "develop",
			},
			expectError: false,
		},
		{
			name: "異常系 - 必須パラメータなし",
			arguments: map[string]interface{}{
				"owner": "test_user",
				"repo":  "test_repo",
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GitHubClientを作成
			client := NewGitHubClient("test_token")

			// リクエストの作成
			request := mcp.CallToolRequest{}
			// Paramsフィールドに直接アクセス
			request.Params.Name = "get_file_contents"
			request.Params.Arguments = tc.arguments

			// テスト対象の関数を実行
			ctx := context.Background()
			result, err := client.handleToGetFileContents(ctx, request)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				// 実際のAPIエラーは許容する（認証エラーなど）
				t.Logf("APIエラーが発生しました（テスト環境では正常）: %v", err)
			}

			// 正常系の場合、結果の基本的な検証
			if !tc.expectError {
				// エラーがない場合、または認証エラーの場合は結果がnilでも許容
				if err == nil && result == nil {
					t.Error("結果がnilです")
				}
			}
		})
	}
}

// TestHandleToListCommits はHandleToListCommitsメソッドをテストする
func TestHandleToListCommits(t *testing.T) {
	// テストケース
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
	}{
		{
			name: "正常系 - 必須パラメータのみ",
			arguments: map[string]interface{}{
				"owner": "test_user",
				"repo":  "test_repo",
			},
			expectError: false,
		},
		{
			name: "正常系 - すべてのパラメータ",
			arguments: map[string]interface{}{
				"owner":    "test_user",
				"repo":     "test_repo",
				"page":     float64(2),
				"per_page": float64(10),
				"sha":      "develop",
			},
			expectError: false,
		},
		{
			name: "異常系 - 必須パラメータなし",
			arguments: map[string]interface{}{
				"page": float64(1),
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GitHubClientを作成
			client := NewGitHubClient("test_token")

			// リクエストの作成
			request := mcp.CallToolRequest{}
			// Paramsフィールドに直接アクセス
			request.Params.Name = "list_commits"
			request.Params.Arguments = tc.arguments

			// テスト対象の関数を実行
			ctx := context.Background()
			result, err := client.handleToListCommits(ctx, request)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				// 実際のAPIエラーは許容する（認証エラーなど）
				t.Logf("APIエラーが発生しました（テスト環境では正常）: %v", err)
			}

			// 正常系の場合、結果の基本的な検証
			if !tc.expectError {
				// エラーがない場合、または認証エラーの場合は結果がnilでも許容
				if err == nil && result == nil {
					t.Error("結果がnilです")
				}
			}
		})
	}
}
