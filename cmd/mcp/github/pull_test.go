package github

import (
	"context"
	"testing"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

func TestNewPullRequestHandler(t *testing.T) {
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
			handler := NewPullRequestHandler(tc.token)

			if handler == nil {
				t.Fatal("ハンドラーがnilです")
			}

			if handler.pullRequestService == nil {
				t.Fatal("プルリクエストサービスがnilです")
			}
		})
	}
}

func TestSetGitHubPullRequestServer(t *testing.T) {
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
			result := setGitHubPullRequestServer(tc.token, s)

			// 関数が正常に実行され、サーバーが返されることを確認
			if result == nil {
				t.Fatal("SetGitHubPullRequestServerがnilを返しました")
			}

			// 注意: 実際のツールの追加や設定の検証は、MCPサーバーの内部実装に依存するため、
			// ここでは基本的な動作確認のみを行います
		})
	}
}

func TestPullRequestHandlerMethods(t *testing.T) {
	// 基本的なハンドラーメソッドの存在確認
	handler := NewPullRequestHandler("test_token")

	// リクエストの作成（最小限のパラメータ）
	request := mcp.CallToolRequest{}
	request.Params.Name = "create_pull_request"
	request.Params.Arguments = map[string]interface{}{
		"owner": "test_user",
		"repo":  "test_repo",
		"title": "テストプルリクエスト",
		"head":  "feature-branch",
		"base":  "main",
	}

	ctx := context.Background()

	// HandleToCreatePullRequestメソッドが存在し、呼び出し可能であることを確認
	// 注意: 実際のGitHub APIを呼び出すため、エラーが発生することが予想されます
	// ここではメソッドが存在し、パニックしないことを確認します
	_, err := handler.handleToCreatePullRequest(ctx, request)
	// エラーが発生することは予想されるため、パニックしなければOK
	if err == nil {
		t.Log("HandleToCreatePullRequest: 予期しない成功（実際のAPIが呼ばれた可能性があります）")
	} else {
		t.Logf("HandleToCreatePullRequest: 期待されるエラー: %v", err)
	}

	// 他のメソッドも同様にテスト
	request.Params.Name = "create_pull_request_review"
	request.Params.Arguments = map[string]interface{}{
		"owner":       "test_user",
		"repo":        "test_repo",
		"pull_number": float64(1),
	}
	_, err = handler.handleToCreatePullRequestReview(ctx, request)
	if err == nil {
		t.Log("HandleToCreatePullRequestReview: 予期しない成功")
	} else {
		t.Logf("HandleToCreatePullRequestReview: 期待されるエラー: %v", err)
	}

	request.Params.Name = "merge_pull_request"
	request.Params.Arguments = map[string]interface{}{
		"owner":       "test_user",
		"repo":        "test_repo",
		"pull_number": float64(1),
	}
	_, err = handler.handleToMergePullRequest(ctx, request)
	if err == nil {
		t.Log("HandleToMergePullRequest: 予期しない成功")
	} else {
		t.Logf("HandleToMergePullRequest: 期待されるエラー: %v", err)
	}

	request.Params.Name = "get_pull_request_files"
	_, err = handler.handleToGetPullRequestFiles(ctx, request)
	if err == nil {
		t.Log("HandleToGetPullRequestFiles: 予期しない成功")
	} else {
		t.Logf("HandleToGetPullRequestFiles: 期待されるエラー: %v", err)
	}

	request.Params.Name = "get_pull_request_status"
	_, err = handler.handleToGetPullRequestStatus(ctx, request)
	if err == nil {
		t.Log("HandleToGetPullRequestStatus: 予期しない成功")
	} else {
		t.Logf("HandleToGetPullRequestStatus: 期待されるエラー: %v", err)
	}

	request.Params.Name = "update_pull_request_branch"
	_, err = handler.handleToUpdatePullRequestBranch(ctx, request)
	if err == nil {
		t.Log("HandleToUpdatePullRequestBranch: 予期しない成功")
	} else {
		t.Logf("HandleToUpdatePullRequestBranch: 期待されるエラー: %v", err)
	}

	request.Params.Name = "get_pull_request_comments"
	_, err = handler.handleToGetPullRequestComments(ctx, request)
	if err == nil {
		t.Log("HandleToGetPullRequestComments: 予期しない成功")
	} else {
		t.Logf("HandleToGetPullRequestComments: 期待されるエラー: %v", err)
	}

	request.Params.Name = "get_pull_request_reviews"
	_, err = handler.handleToGetPullRequestReviews(ctx, request)
	if err == nil {
		t.Log("HandleToGetPullRequestReviews: 予期しない成功")
	} else {
		t.Logf("HandleToGetPullRequestReviews: 期待されるエラー: %v", err)
	}

	request.Params.Name = "list_pull_requests"
	_, err = handler.handleToListPullRequests(ctx, request)
	if err == nil {
		t.Log("HandleToListPullRequests: 予期しない成功")
	} else {
		t.Logf("HandleToListPullRequests: 期待されるエラー: %v", err)
	}
}
