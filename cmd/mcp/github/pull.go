package github

import (
	"context"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	usecases "github.com/landmaster135/devbox/internal/github/usecases"
)

// PullRequestHandler はPull Request関連のMCPハンドラーを管理します
type PullRequestHandler struct {
	pullRequestService *usecases.GitHubPullRequestService
}

// NewPullRequestHandler は新しいPullRequestHandlerを作成します
func NewPullRequestHandler(token string) *PullRequestHandler {
	return &PullRequestHandler{
		pullRequestService: usecases.NewGitHubPullRequestService(token),
	}
}

// handleToCreatePullRequest は新しいプルリクエストを作成して、結果をJSON形式で返します
func (h *PullRequestHandler) handleToCreatePullRequest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}
	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}
	title, err := request.RequireString("title")
	if err != nil {
		return nil, err
	}
	head, err := request.RequireString("head")
	if err != nil {
		return nil, err
	}
	base, err := request.RequireString("base")
	if err != nil {
		return nil, err
	}

	body := request.GetString("body", "")
	draft := request.GetBool("draft", false)

	result, err := h.pullRequestService.HandleToCreatePullRequest(owner, repo, title, head, base, body, draft)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}

// handleToCreatePullRequestReview はプルリクエストにレビューを作成して、結果をJSON形式で返します
func (h *PullRequestHandler) handleToCreatePullRequestReview(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}
	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}
	pullNumber, err := request.RequireInt("pull_number")
	if err != nil {
		return nil, err
	}

	event := request.GetString("event", "")
	body := request.GetString("body", "")

	result, err := h.pullRequestService.HandleToCreatePullRequestReview(owner, repo, pullNumber, event, body)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}

// handleToMergePullRequest はプルリクエストをマージして、結果をJSON形式で返します
func (h *PullRequestHandler) handleToMergePullRequest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}
	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}
	pullNumber, err := request.RequireInt("pull_number")
	if err != nil {
		return nil, err
	}

	commitTitle := request.GetString("commit_title", "")
	commitMessage := request.GetString("commit_message", "")
	mergeMethod := request.GetString("merge_method", "")

	result, err := h.pullRequestService.HandleToMergePullRequest(owner, repo, pullNumber, commitTitle, commitMessage, mergeMethod)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}

// handleToGetPullRequestFiles はプルリクエストで変更されたファイル一覧を取得して、結果をJSON形式で返します
func (h *PullRequestHandler) handleToGetPullRequestFiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}
	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}
	pullNumber, err := request.RequireInt("pull_number")
	if err != nil {
		return nil, err
	}

	result, err := h.pullRequestService.HandleToGetPullRequestFiles(owner, repo, pullNumber)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}

// handleToGetPullRequestStatus はプルリクエストのステータスを取得して、結果をJSON形式で返します
func (h *PullRequestHandler) handleToGetPullRequestStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}
	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}
	pullNumber, err := request.RequireInt("pull_number")
	if err != nil {
		return nil, err
	}

	result, err := h.pullRequestService.HandleToGetPullRequestStatus(owner, repo, pullNumber)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}

// handleToUpdatePullRequestBranch はプルリクエストのブランチを更新して、結果をJSON形式で返します
func (h *PullRequestHandler) handleToUpdatePullRequestBranch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}
	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}
	pullNumber, err := request.RequireInt("pull_number")
	if err != nil {
		return nil, err
	}

	expectedHeadSHA := request.GetString("expected_head_sha", "")

	result, err := h.pullRequestService.HandleToUpdatePullRequestBranch(owner, repo, pullNumber, expectedHeadSHA)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}

// handleToGetPullRequestComments はプルリクエストのコメントを取得して、結果をJSON形式で返します
func (h *PullRequestHandler) handleToGetPullRequestComments(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}
	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}
	pullNumber, err := request.RequireInt("pull_number")
	if err != nil {
		return nil, err
	}

	result, err := h.pullRequestService.HandleToGetPullRequestComments(owner, repo, pullNumber)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}

// handleToGetPullRequestReviews はプルリクエストのレビューを取得して、結果をJSON形式で返します
func (h *PullRequestHandler) handleToGetPullRequestReviews(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}
	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}
	pullNumber, err := request.RequireInt("pull_number")
	if err != nil {
		return nil, err
	}

	result, err := h.pullRequestService.HandleToGetPullRequestReviews(owner, repo, pullNumber)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}

// handleToListPullRequests はリポジトリのプルリクエスト一覧を取得して、結果をJSON形式で返します
func (h *PullRequestHandler) handleToListPullRequests(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}
	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}

	state := request.GetString("state", "")
	sort := request.GetString("sort", "")
	direction := request.GetString("direction", "")
	head := request.GetString("head", "")
	base := request.GetString("base", "")
	perPage := request.GetInt("per_page", 30)
	page := request.GetInt("page", 1)

	result, err := h.pullRequestService.HandleToListPullRequests(owner, repo, state, sort, direction, head, base, perPage, page)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}

// handleToCreatePullRequestWithCurrentBranch は指定されたリポジトリパスの現在のブランチを使用してプルリクエストを作成して、結果をJSON形式で返します
func (h *PullRequestHandler) handleToCreatePullRequestWithCurrentBranch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}
	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}
	title, err := request.RequireString("title")
	if err != nil {
		return nil, err
	}
	base, err := request.RequireString("base")
	if err != nil {
		return nil, err
	}
	repoPath, err := request.RequireString("repo_path")
	if err != nil {
		return nil, err
	}

	body := request.GetString("body", "")
	draft := request.GetBool("draft", false)

	result, err := h.pullRequestService.HandleToCreatePullRequestWithCurrentBranch(owner, repo, title, base, body, draft, repoPath)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}

// setGitHubPullRequestServer は受け取ったMCPサーバにGitHubプルリクエスト用のツールを付与して、そのMCPサーバを返します。
func setGitHubPullRequestServer(token string, s *server.MCPServer) *server.MCPServer {
	// PullRequestHandlerを初期化
	handler := NewPullRequestHandler(token)

	// ツール1: プルリクエストの作成
	createPullRequestTool := mcp.NewTool("create_pull_request",
		mcp.WithDescription("Create a new pull request in a GitHub repository"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("Pull request title"),
		),
		mcp.WithString("head",
			mcp.Required(),
			mcp.Description("The name of the branch where your changes are implemented"),
		),
		mcp.WithString("base",
			mcp.Required(),
			mcp.Description("The name of the branch you want the changes pulled into"),
		),
		mcp.WithString("body",
			mcp.Description("Pull request body"),
		),
		mcp.WithBoolean("draft",
			mcp.Description("Whether to create a draft pull request"),
		),
	)
	s.AddTool(createPullRequestTool, handler.handleToCreatePullRequest)

	// ツール2: プルリクエストレビューの作成
	createPullRequestReviewTool := mcp.NewTool("create_pull_request_review",
		mcp.WithDescription("Create a review on a pull request"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithNumber("pull_number",
			mcp.Required(),
			mcp.Description("Pull request number"),
		),
		mcp.WithString("event",
			mcp.Description("Review event: APPROVE, REQUEST_CHANGES, COMMENT"),
			mcp.Enum("APPROVE", "REQUEST_CHANGES", "COMMENT"),
		),
		mcp.WithString("body",
			mcp.Description("Review body"),
		),
	)
	s.AddTool(createPullRequestReviewTool, handler.handleToCreatePullRequestReview)

	// ツール3: プルリクエストのマージ
	mergePullRequestTool := mcp.NewTool("merge_pull_request",
		mcp.WithDescription("Merge a pull request"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithNumber("pull_number",
			mcp.Required(),
			mcp.Description("Pull request number"),
		),
		mcp.WithString("commit_title",
			mcp.Description("Title for the automatic commit message"),
		),
		mcp.WithString("commit_message",
			mcp.Description("Extra detail to append to automatic commit message"),
		),
		mcp.WithString("merge_method",
			mcp.Description("Merge method to use: merge, squash, rebase"),
			mcp.Enum("merge", "squash", "rebase"),
		),
	)
	s.AddTool(mergePullRequestTool, handler.handleToMergePullRequest)

	// ツール4: プルリクエストのファイル一覧取得
	getPullRequestFilesTool := mcp.NewTool("get_pull_request_files",
		mcp.WithDescription("Get the list of files changed in a pull request"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithNumber("pull_number",
			mcp.Required(),
			mcp.Description("Pull request number"),
		),
	)
	s.AddTool(getPullRequestFilesTool, handler.handleToGetPullRequestFiles)

	// ツール5: プルリクエストのステータス取得
	getPullRequestStatusTool := mcp.NewTool("get_pull_request_status",
		mcp.WithDescription("Get the combined status of all status checks for a pull request"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithNumber("pull_number",
			mcp.Required(),
			mcp.Description("Pull request number"),
		),
	)
	s.AddTool(getPullRequestStatusTool, handler.handleToGetPullRequestStatus)

	// ツール6: プルリクエストブランチの更新
	updatePullRequestBranchTool := mcp.NewTool("update_pull_request_branch",
		mcp.WithDescription("Update a pull request branch with the latest changes from the base branch"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithNumber("pull_number",
			mcp.Required(),
			mcp.Description("Pull request number"),
		),
		mcp.WithString("expected_head_sha",
			mcp.Description("The expected SHA of the pull request head"),
		),
	)
	s.AddTool(updatePullRequestBranchTool, handler.handleToUpdatePullRequestBranch)

	// ツール7: プルリクエストコメントの取得
	getPullRequestCommentsTool := mcp.NewTool("get_pull_request_comments",
		mcp.WithDescription("Get the review comments on a pull request"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithNumber("pull_number",
			mcp.Required(),
			mcp.Description("Pull request number"),
		),
	)
	s.AddTool(getPullRequestCommentsTool, handler.handleToGetPullRequestComments)

	// ツール8: プルリクエストレビューの取得
	getPullRequestReviewsTool := mcp.NewTool("get_pull_request_reviews",
		mcp.WithDescription("Get the reviews on a pull request"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithNumber("pull_number",
			mcp.Required(),
			mcp.Description("Pull request number"),
		),
	)
	s.AddTool(getPullRequestReviewsTool, handler.handleToGetPullRequestReviews)

	// ツール9: プルリクエスト一覧の取得
	listPullRequestsTool := mcp.NewTool("list_pull_requests",
		mcp.WithDescription("List and filter repository pull requests"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithString("state",
			mcp.Description("Pull request state: open, closed, or all (default: open)"),
			mcp.Enum("open", "closed", "all"),
		),
		mcp.WithString("sort",
			mcp.Description("Sort field: created, updated, popularity, long-running (default: created)"),
			mcp.Enum("created", "updated", "popularity", "long-running"),
		),
		mcp.WithString("direction",
			mcp.Description("Sort direction: asc or desc (default: desc)"),
			mcp.Enum("asc", "desc"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Results per page (default: 30, max: 100)"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number (default: 1)"),
		),
		mcp.WithString("head",
			mcp.Description("Filter by head user or head organization and branch name in the format of 'user:ref-name' or 'organization:ref-name'"),
		),
		mcp.WithString("base",
			mcp.Description("Filter by base branch name"),
		),
	)
	s.AddTool(listPullRequestsTool, handler.handleToListPullRequests)

	// ツール10: 現在のブランチを使用したプルリクエストの作成
	createPullRequestWithCurrentBranchTool := mcp.NewTool("create_pull_request_with_current_branch",
		mcp.WithDescription("Create a new pull request using the current Git branch from specified repository path"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("Pull request title"),
		),
		mcp.WithString("base",
			mcp.Required(),
			mcp.Description("The name of the branch you want the changes pulled into"),
		),
		mcp.WithString("repo_path",
			mcp.Required(),
			mcp.Description("Absolute path to the local Git repository"),
		),
		mcp.WithString("body",
			mcp.Description("Pull request body"),
		),
		mcp.WithBoolean("draft",
			mcp.Description("Whether to create a draft pull request"),
		),
	)
	s.AddTool(createPullRequestWithCurrentBranchTool, handler.handleToCreatePullRequestWithCurrentBranch)

	return s
}
