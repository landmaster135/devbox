package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// CreateIssue は新しいイシューを作成します
func (c *GitHubClient) CreateIssue(owner, repo string, options map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues", apiBaseURL, owner, repo)
	jsonBody, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	data, err := c.doRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ヘルパー関数: オプションマップにパラメータを追加
func addToOptions(options map[string]interface{}, args map[string]interface{}, key string) {
	if val, ok := args[key]; ok {
		options[key] = val
	}
}

// HandleToCreateIssue は新しいイシューを作成して、結果をJSON形式で返します
func (c *GitHubClient) HandleToCreateIssue(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}

	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}

	options := make(map[string]interface{})
	title, err := request.RequireString("title")
	if err != nil {
		return nil, err
	}
	options["title"] = title

	// オプションパラメータを追加
	body := request.GetString("body", "")
	if body != "" {
		options["body"] = body
	}

	// 配列パラメータを追加
	args := request.GetArguments()
	addToOptions(options, args, "labels")
	addToOptions(options, args, "assignees")

	result, err := c.CreateIssue(owner, repo, options)
	if err != nil {
		return nil, err
	}

	return returnJSONResult(result)
}

// ListIssues はリポジトリのイシュー一覧を取得します
func (c *GitHubClient) ListIssues(owner, repo string, options map[string]interface{}) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues", apiBaseURL, owner, repo)

	// クエリパラメータを追加
	queryParams := []string{}
	for k, v := range options {
		queryParams = append(queryParams, fmt.Sprintf("%s=%v", k, v))
	}
	if len(queryParams) > 0 {
		url += "?" + strings.Join(queryParams, "&")
	}

	data, err := c.doRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// HandleToListIssues はリポジトリのイシュー一覧を取得して、結果をJSON形式で返します
func (c *GitHubClient) HandleToListIssues(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}

	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}

	options := make(map[string]interface{})

	// 文字列オプションパラメータを追加
	state := request.GetString("state", "")
	if state != "" {
		options["state"] = state
	}

	sort := request.GetString("sort", "")
	if sort != "" {
		options["sort"] = sort
	}

	direction := request.GetString("direction", "")
	if direction != "" {
		options["direction"] = direction
	}

	// 数値オプションパラメータを追加
	options["per_page"] = request.GetInt("per_page", 30)

	options["page"] = request.GetInt("page", 1)

	result, err := c.ListIssues(owner, repo, options)
	if err != nil {
		return nil, err
	}

	return returnJSONResult(result)
}

// UpdateIssue は既存のイシューを更新します
func (c *GitHubClient) UpdateIssue(owner, repo string, issueNumber int, options map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d", apiBaseURL, owner, repo, issueNumber)

	jsonBody, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	data, err := c.doRequest("PATCH", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// HandleToUpdateIssue は既存のイシューを更新して、結果をJSON形式で返します
func (c *GitHubClient) HandleToUpdateIssue(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}

	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}

	issueNumber := request.GetInt("issue_number", 0)

	options := make(map[string]interface{})

	// 文字列オプションパラメータを追加
	title := request.GetString("title", "")
	if title != "" {
		options["title"] = title
	}

	body := request.GetString("body", "")
	if body != "" {
		options["body"] = body
	}

	state := request.GetString("state", "")
	if state != "" {
		options["state"] = state
	}

	// 配列パラメータを追加
	args := request.GetArguments()
	addToOptions(options, args, "labels")
	addToOptions(options, args, "assignees")

	result, err := c.UpdateIssue(owner, repo, issueNumber, options)
	if err != nil {
		return nil, err
	}

	return returnJSONResult(result)
}

// AddIssueComment はイシューにコメントを追加します
func (c *GitHubClient) AddIssueComment(owner, repo string, issueNumber int, body string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments", apiBaseURL, owner, repo, issueNumber)

	jsonBody, err := json.Marshal(map[string]string{"body": body})
	if err != nil {
		return nil, err
	}

	data, err := c.doRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// HandleToAddIssueComment はイシューにコメントを追加して、結果をJSON形式で返します
func (c *GitHubClient) HandleToAddIssueComment(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	result, err := c.AddIssueComment(owner, repo, issueNumber, body)
	if err != nil {
		return nil, err
	}

	return returnJSONResult(result)
}

// SetGitHubIssueServer は受け取ったMCPサーバにGitHub用のツールを付与して、そのMCPサーバを返します。
func SetGitHubIssueServer(token string, s *server.MCPServer) *server.MCPServer {
	// GitHubクライアントを初期化
	client := NewGitHubClient(token)

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
	s.AddTool(createIssueTool, client.HandleToCreateIssue)

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
	s.AddTool(listIssuesTool, client.HandleToListIssues)

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

	s.AddTool(updateIssueTool, client.HandleToUpdateIssue)

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

	s.AddTool(addIssueCommentTool, client.HandleToAddIssueComment)

	return s
}
