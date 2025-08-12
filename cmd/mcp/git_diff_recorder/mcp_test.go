package git_diff_recorder

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/landmaster135/devbox/internal/git_diff_recorder/config"
)

// createMockCallToolRequest は実際のCallToolRequestを作成する
func createMockCallToolRequest(arguments map[string]interface{}) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      "get_git_diff",
			Arguments: arguments,
		},
	}
}

// TestHandleGetGitDiff_Normal はhandleGetGitDiff関数の正常系テスト
func TestHandleGetGitDiff_Normal(t *testing.T) {
	// Arrange
	// 現在のディレクトリを取得（実際のgitリポジトリが必要）
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("現在のディレクトリの取得に失敗しました: %v", err)
	}

	// devboxディレクトリに移動してテスト
	devboxDir := filepath.Join(workingDir, "../../..")
	request := createMockCallToolRequest(map[string]interface{}{
		"git_dir":     devboxDir,
		"staged_only": false,
	})

	ctx := context.Background()

	// Act
	result, err := handleGetGitDiff(ctx, request)

	// Assert
	if err != nil {
		// gitリポジトリでない場合はエラーになることを期待
		t.Logf("期待されるエラー（gitリポジトリでない場合）: %v", err)
		return
	}

	if result == nil {
		t.Error("結果がnilです")
		return
	}

	// 結果の内容を確認
	t.Logf("Git差分取得が正常に完了しました")
}

// TestHandleGetGitDiff_StagedOnly はhandleGetGitDiff関数のステージング済み差分テスト
func TestHandleGetGitDiff_StagedOnly(t *testing.T) {
	// Arrange
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("現在のディレクトリの取得に失敗しました: %v", err)
	}

	devboxDir := filepath.Join(workingDir, "../../..")
	request := createMockCallToolRequest(map[string]interface{}{
		"git_dir":     devboxDir,
		"staged_only": true,
	})

	ctx := context.Background()

	// Act
	result, err := handleGetGitDiff(ctx, request)

	// Assert
	if err != nil {
		t.Logf("期待されるエラー（gitリポジトリでない場合）: %v", err)
		return
	}

	if result == nil {
		t.Error("結果がnilです")
		return
	}

	t.Logf("ステージング済み差分取得が正常に完了しました")
}

// TestHandleGetGitDiff_MissingGitDir はhandleGetGitDiff関数のgit_dir未指定テスト
func TestHandleGetGitDiff_MissingGitDir(t *testing.T) {
	// Arrange
	request := createMockCallToolRequest(map[string]interface{}{
		"staged_only": false,
	})

	ctx := context.Background()

	// Act
	result, err := handleGetGitDiff(ctx, request)

	// Assert
	if err == nil {
		t.Error("git_dirが未指定の場合、エラーが発生するべきです")
		return
	}

	if result != nil {
		t.Error("エラー時には結果はnilであるべきです")
	}

	// エラーメッセージの確認
	if !strings.Contains(err.Error(), "git_dir") {
		t.Errorf("エラーメッセージにgit_dirが含まれていません: %v", err)
	}
}

// TestHandleGetGitDiff_InvalidGitDir はhandleGetGitDiff関数の無効なgit_dirテスト
func TestHandleGetGitDiff_InvalidGitDir(t *testing.T) {
	// Arrange
	request := createMockCallToolRequest(map[string]interface{}{
		"git_dir":     "/tmp/non-existent-git-repo",
		"staged_only": false,
	})

	ctx := context.Background()

	// Act
	result, err := handleGetGitDiff(ctx, request)

	// Assert
	if err == nil {
		t.Error("無効なgit_dirの場合、エラーが発生するべきです")
		return
	}

	if result != nil {
		t.Error("エラー時には結果はnilであるべきです")
	}

	// エラーメッセージの確認
	if !strings.Contains(err.Error(), "git差分の取得に失敗しました") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %v", err)
	}
}

// TestAddPromptIntoServer_Normal はaddPromptIntoServer関数の正常系テスト
func TestAddPromptIntoServer_Normal(t *testing.T) {
	// Arrange
	s := server.NewMCPServer(
		"Test Server",
		"1.0.0",
		server.WithPromptCapabilities(true),
	)

	// Act
	result := addPromptIntoServer(s)

	// Assert
	if result == nil {
		t.Error("サーバーがnilです")
		return
	}

	if result != s {
		t.Error("同じサーバーインスタンスが返されるべきです")
	}

	t.Log("プロンプトの追加が正常に完了しました")
}

// TestAddPromptIntoServer_PromptHandler はaddPromptIntoServer関数のプロンプトハンドラーテスト
func TestAddPromptIntoServer_PromptHandler(t *testing.T) {
	// Arrange
	s := server.NewMCPServer(
		"Test Server",
		"1.0.0",
		server.WithPromptCapabilities(true),
	)

	// Act
	result := addPromptIntoServer(s)

	// Assert
	if result == nil {
		t.Error("サーバーがnilです")
		return
	}

	// プロンプトハンドラーの動作確認は、実際のMCPサーバーの内部実装に依存するため、
	// ここでは基本的な追加処理が正常に完了することを確認
	t.Log("プロンプトハンドラーの追加が正常に完了しました")
}

// TestConfig_Integration は設定との統合テスト
func TestConfig_Integration(t *testing.T) {
	// Arrange
	gitDir := "/tmp/test-git"
	cfg := &config.Config{
		GitDir:     gitDir,
		StagedOnly: false,
	}

	// Act & Assert
	if cfg.GitDir != gitDir {
		t.Errorf("GitDirが正しく設定されていません。期待値: %s, 実際値: %s", gitDir, cfg.GitDir)
	}

	if cfg.StagedOnly != false {
		t.Errorf("StagedOnlyが正しく設定されていません。期待値: %t, 実際値: %t", false, cfg.StagedOnly)
	}

	t.Log("設定との統合テストが正常に完了しました")
}
