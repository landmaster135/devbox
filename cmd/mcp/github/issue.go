package github

import (
	"context"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
	"github.com/landmaster135/devbox/internal/github/usecases"
)

// IssueHandler はIssue関連のMCPハンドラーを管理します
type IssueHandler struct {
	issueService *usecases.GitHubIssueService
}

// NewIssueHandler は新しいIssueHandlerを作成します
func NewIssueHandler(token string) *IssueHandler {
	return &IssueHandler{
		issueService: usecases.NewGitHubIssueService(token),
	}
}

// HandleToCreateIssue は新しいイシューを作成して、結果をJSON形式で返します
func (h *IssueHandler) HandleToCreateIssue(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	body := request.GetString("body", "")

	// 配列パラメータを取得
	args := request.GetArguments()
	var labels []interface{}
	var assignees []interface{}

	if labelsVal, ok := args["labels"]; ok {
		if labelsArray, ok := labelsVal.([]interface{}); ok {
			labels = labelsArray
		}
	}

	if assigneesVal, ok := args["assignees"]; ok {
		if assigneesArray, ok := assigneesVal.([]interface{}); ok {
			assignees = assigneesArray
		}
	}

	result, err := h.issueService.HandleToCreateIssue(owner, repo, title, body, labels, assignees)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}

// HandleToListIssues はリポジトリのイシュー一覧を取得して、結果をJSON形式で返します
func (h *IssueHandler) HandleToListIssues(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	perPage := request.GetInt("per_page", 30)
	page := request.GetInt("page", 1)

	result, err := h.issueService.HandleToListIssues(owner, repo, state, sort, direction, perPage, page)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}

// HandleToUpdateIssue は既存のイシューを更新して、結果をJSON形式で返します
func (h *IssueHandler) HandleToUpdateIssue(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}

	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}

	issueNumber := request.GetInt("issue_number", 0)
	title := request.GetString("title", "")
	body := request.GetString("body", "")
	state := request.GetString("state", "")

	// 配列パラメータを取得
	args := request.GetArguments()
	var labels []interface{}
	var assignees []interface{}

	if labelsVal, ok := args["labels"]; ok {
		if labelsArray, ok := labelsVal.([]interface{}); ok {
			labels = labelsArray
		}
	}

	if assigneesVal, ok := args["assignees"]; ok {
		if assigneesArray, ok := assigneesVal.([]interface{}); ok {
			assignees = assigneesArray
		}
	}

	result, err := h.issueService.HandleToUpdateIssue(owner, repo, issueNumber, title, body, state, labels, assignees)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}

// HandleToAddIssueComment はイシューにコメントを追加して、結果をJSON形式で返します
func (h *IssueHandler) HandleToAddIssueComment(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}

	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}

	issueNumber := request.GetInt("issue_number", 0)

	body, err := request.RequireString("body")
	if err != nil {
		return nil, err
	}

	result, err := h.issueService.HandleToAddIssueComment(owner, repo, issueNumber, body)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}

// SetGitHubIssueServer は受け取ったMCPサーバにGitHub用のツールを付与して、そのMCPサーバを返します。
func SetGitHubIssueServer(token string, s *server.MCPServer) *server.MCPServer {
	// IssueHandlerを初期化
	handler := NewIssueHandler(token)

	// ツール1: イシューの作成
	createIssueTool := mcp.NewTool("create_issue",
		mcp.WithDescription("Create a new issue in a GitHub repository"),
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
			mcp.Description("Issue title"),
		),
		mcp.WithString("body",
			mcp.Description("Issue body"),
		),
		mcp.WithArray("labels",
			mcp.Description("Issue labels"),
		),
		mcp.WithArray("assignees",
			mcp.Description("Users to assign to this issue"),
		),
	)
	s.AddTool(createIssueTool, handler.HandleToCreateIssue)

	// ツール2: イシュー一覧の取得
	listIssuesTool := mcp.NewTool("list_issues",
		mcp.WithDescription("List issues in a GitHub repository"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithString("state",
			mcp.Description("Issue state: open, closed, or all (default: open)"),
			mcp.Enum("open", "closed", "all"),
		),
		mcp.WithString("sort",
			mcp.Description("Sort field: created, updated, or comments (default: created)"),
			mcp.Enum("created", "updated", "comments"),
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
	)
	s.AddTool(listIssuesTool, handler.HandleToListIssues)

	// ツール3: イシューの更新
	updateIssueTool := mcp.NewTool("update_issue",
		mcp.WithDescription("Update an existing issue in a GitHub repository"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithNumber("issue_number",
			mcp.Required(),
			mcp.Description("Issue number"),
		),
		mcp.WithString("title",
			mcp.Description("New issue title"),
		),
		mcp.WithString("body",
			mcp.Description("New issue body"),
		),
		mcp.WithString("state",
			mcp.Description("State of the issue: open or closed"),
			mcp.Enum("open", "closed"),
		),
		mcp.WithArray("labels",
			mcp.Description("New labels for the issue"),
		),
		mcp.WithArray("assignees",
			mcp.Description("New assignees for the issue"),
		),
	)

	s.AddTool(updateIssueTool, handler.HandleToUpdateIssue)

	// ツール4: イシューコメントの追加
	addIssueCommentTool := mcp.NewTool("add_issue_comment",
		mcp.WithDescription("Add a comment to an existing issue"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithNumber("issue_number",
			mcp.Required(),
			mcp.Description("Issue number"),
		),
		mcp.WithString("body",
			mcp.Required(),
			mcp.Description("Comment body"),
		),
	)

	s.AddTool(addIssueCommentTool, handler.HandleToAddIssueComment)

	return s
}
