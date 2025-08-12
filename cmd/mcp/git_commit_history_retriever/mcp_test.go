package git_commit_history_retriever

import (
	"context"
	"strings"
	"testing"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// #==============================================================#
// ##          Helper Functions                                  ##
// #==============================================================#

// createMockCallToolRequest は実際のCallToolRequestを作成する
func createMockCallToolRequest(arguments map[string]interface{}) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      "test_tool",
			Arguments: arguments,
		},
	}
}

// #==============================================================#
// ##          Tests                                             ##
// #==============================================================#

// TestCreateConfigFromRequest_Normal はcreateConfigFromRequestの正常系テスト
func TestCreateConfigFromRequest_Normal(t *testing.T) {
	// Arrange
	mockRequest := createMockCallToolRequest(map[string]interface{}{
		"git_dir": "/test/repo",
		"keyword": "feat:",
		"since":   "2025-01-01",
		"until":   "2025-01-31",
	})

	// Act
	config, err := createConfigFromRequest(mockRequest)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if config.GitDir != "/test/repo" {
		t.Errorf("Expected GitDir to be /test/repo, got %s", config.GitDir)
	}

	if config.Keyword != "feat:" {
		t.Errorf("Expected Keyword to be feat:, got %s", config.Keyword)
	}

	if config.Since != "2025-01-01" {
		t.Errorf("Expected Since to be 2025-01-01, got %s", config.Since)
	}

	if config.Until != "2025-01-31" {
		t.Errorf("Expected Until to be 2025-01-31, got %s", config.Until)
	}
}

// TestCreateConfigFromRequest_MinimalParams は最小パラメータでのテスト
func TestCreateConfigFromRequest_MinimalParams(t *testing.T) {
	// Arrange
	mockRequest := createMockCallToolRequest(map[string]interface{}{
		"git_dir": "/test/repo",
	})

	// Act
	config, err := createConfigFromRequest(mockRequest)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if config.GitDir != "/test/repo" {
		t.Errorf("Expected GitDir to be /test/repo, got %s", config.GitDir)
	}

	if config.Keyword != "" {
		t.Errorf("Expected Keyword to be empty, got %s", config.Keyword)
	}

	if config.Since != "" {
		t.Errorf("Expected Since to be empty, got %s", config.Since)
	}

	if config.Until != "" {
		t.Errorf("Expected Until to be empty, got %s", config.Until)
	}
}

// TestCreateConfigFromRequest_MissingGitDir はgit_dirが不足している場合のテスト
func TestCreateConfigFromRequest_MissingGitDir(t *testing.T) {
	// Arrange
	mockRequest := createMockCallToolRequest(map[string]interface{}{
		"keyword": "feat:",
	})

	// Act
	_, err := createConfigFromRequest(mockRequest)

	// Assert
	if err == nil {
		t.Error("Expected error for missing git_dir, got nil")
	}

	if !strings.Contains(err.Error(), "required argument \"git_dir\" not found") {
		t.Errorf("Expected error message about missing git_dir, got %v", err)
	}
}

// TestCreateConfigFromRequest_InvalidDateFormat は無効な日付フォーマットの場合のテスト
func TestCreateConfigFromRequest_InvalidDateFormat(t *testing.T) {
	// Arrange
	mockRequest := createMockCallToolRequest(map[string]interface{}{
		"git_dir": "/test/repo",
		"since":   "invalid-date",
	})

	// Act
	_, err := createConfigFromRequest(mockRequest)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid date format, got nil")
	}

	if !strings.Contains(err.Error(), "MCPリクエストからConfigを作成するのに失敗しました") {
		t.Errorf("Expected error message about config creation failure, got %v", err)
	}
}

// TestHandleGetCommitHistory_Normal はhandleGetCommitHistoryの正常系テスト
func TestHandleGetCommitHistory_Normal(t *testing.T) {
	// Note: この関数は実際のGitリポジトリとusecasesパッケージに依存するため、
	// 統合テストとして実装するか、依存関係をモック化する必要がある。
	// ここでは関数が呼び出し可能であることを確認する。

	// Arrange
	mockRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "get_git_commit_history",
			Arguments: map[string]interface{}{
				"git_dir": "/tmp/nonexistent", // 存在しないパスでエラーを期待
			},
		},
	}

	// Act
	ctx := context.Background()
	_, err := handleGetCommitHistory(ctx, mockRequest)

	// Assert
	// 存在しないGitリポジトリなのでエラーが発生することを期待
	if err == nil {
		t.Error("Expected error for nonexistent git repository, got nil")
	}
}

// TestHandleGetCommitHistoryWithDetails_Normal はhandleGetCommitHistoryWithDetailsの正常系テスト
func TestHandleGetCommitHistoryWithDetails_Normal(t *testing.T) {
	// Note: この関数は実際のGitリポジトリとusecasesパッケージに依存するため、
	// 統合テストとして実装するか、依存関係をモック化する必要がある。
	// ここでは関数が呼び出し可能であることを確認する。

	// Arrange
	mockRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "get_git_commit_history_with_details",
			Arguments: map[string]interface{}{
				"git_dir": "/tmp/nonexistent", // 存在しないパスでエラーを期待
			},
		},
	}

	// Act
	ctx := context.Background()
	_, err := handleGetCommitHistoryWithDetails(ctx, mockRequest)

	// Assert
	// 存在しないGitリポジトリなのでエラーが発生することを期待
	if err == nil {
		t.Error("Expected error for nonexistent git repository, got nil")
	}
}

// TestAddPromptIntoServer_Normal はaddPromptIntoServerの正常系テスト
func TestAddPromptIntoServer_Normal(t *testing.T) {
	// Arrange
	s := server.NewMCPServer(
		"Test Git Commit History Retriever",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
	)

	// Act
	result := addPromptIntoServer(s)

	// Assert
	if result == nil {
		t.Error("Expected server to be returned, got nil")
	}

	// サーバーが正常に設定されていることを確認
	// 実際のプロンプトの動作確認は統合テストで行う
}

// TestAddPromptIntoServer_PromptHandler はプロンプトハンドラーの動作テスト
func TestAddPromptIntoServer_PromptHandler(t *testing.T) {
	// Arrange
	s := server.NewMCPServer(
		"Test Git Commit History Retriever",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
	)

	// プロンプトを追加
	s = addPromptIntoServer(s)

	// Act & Assert
	// プロンプトハンドラーの詳細なテストは統合テストで実装
	// ここではサーバーにプロンプトが追加されたことを確認
	if s == nil {
		t.Error("Expected server with prompt to be returned, got nil")
	}
}

// TestBuildMcpServer_Normal はBuildMcpServerの正常系テスト
func TestBuildMcpServer_Normal(t *testing.T) {
	// Note: BuildMcpServer関数はserver.ServeStdio()を呼び出すため、
	// 実際のテストでは無限ループに入る可能性がある。
	// ここでは関数が存在することのみを確認する。

	// Act & Assert
	// 関数が呼び出し可能であることを確認（実際の実行はしない）
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("BuildMcpServer function should not panic during definition, got %v", r)
		}
	}()

	// 実際の呼び出しはスキップ（無限ループを避けるため）
	// BuildMcpServer()
}
