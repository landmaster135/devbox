package github

import (
	"context"
	"testing"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

func TestNewIssueHandler(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "トークンあり",
			token: "testtoken",
		},
		{
			name:  "トークンなし",
			token: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewIssueHandler(tc.token)

			if handler == nil {
				t.Fatal("ハンドラーがnilです")
			}

			if handler.issueService == nil {
				t.Fatal("イシューサービスがnilです")
			}
		})
	}
}

func TestSetGitHubIssueServer(t *testing.T) {
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
			result := SetGitHubIssueServer(tc.token, s)

			// 関数が正常に実行され、サーバーが返されることを確認
			if result == nil {
				t.Fatal("SetGitHubIssueServerがnilを返しました")
			}

			// 注意: 実際のツールの追加や設定の検証は、MCPサーバーの内部実装に依存するため、
			// ここでは基本的な動作確認のみを行います
		})
	}
}

func TestIssueHandlerMethods(t *testing.T) {
	// 基本的なハンドラーメソッドの存在確認
	handler := NewIssueHandler("test_token")

	// リクエストの作成（最小限のパラメータ）
	request := mcp.CallToolRequest{}
	request.Params.Name = "create_issue"
	request.Params.Arguments = map[string]interface{}{
		"owner": "test_user",
		"repo":  "test_repo",
		"title": "テストイシュー",
	}

	ctx := context.Background()

	// HandleToCreateIssueメソッドが存在し、呼び出し可能であることを確認
	// 注意: 実際のGitHub APIを呼び出すため、エラーが発生することが予想されます
	// ここではメソッドが存在し、パニックしないことを確認します
	_, err := handler.HandleToCreateIssue(ctx, request)
	// エラーが発生することは予想されるため、パニックしなければOK
	if err == nil {
		t.Log("HandleToCreateIssue: 予期しない成功（実際のAPIが呼ばれた可能性があります）")
	} else {
		t.Logf("HandleToCreateIssue: 期待されるエラー: %v", err)
	}

	// 他のメソッドも同様にテスト
	request.Params.Name = "list_issues"
	_, err = handler.HandleToListIssues(ctx, request)
	if err == nil {
		t.Log("HandleToListIssues: 予期しない成功")
	} else {
		t.Logf("HandleToListIssues: 期待されるエラー: %v", err)
	}

	request.Params.Name = "update_issue"
	request.Params.Arguments = map[string]interface{}{
		"issue_number": float64(1),
	}
	_, err = handler.HandleToUpdateIssue(ctx, request)
	if err == nil {
		t.Log("HandleToUpdateIssue: 予期しない成功")
	} else {
		t.Logf("HandleToUpdateIssue: 期待されるエラー: %v", err)
	}

	request.Params.Name = "add_issue_comment"
	request.Params.Arguments = map[string]interface{}{
		"body": "テストコメント",
	}
	_, err = handler.HandleToAddIssueComment(ctx, request)
	if err == nil {
		t.Log("HandleToAddIssueComment: 予期しない成功")
	} else {
		t.Logf("HandleToAddIssueComment: 期待されるエラー: %v", err)
	}
}
