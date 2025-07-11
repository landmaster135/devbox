package github

import (
	"context"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	usecases "github.com/landmaster135/devbox/internal/github/usecases"
)

// HandleToListCommits はリポジトリのコミット一覧を取得して、結果をJSON形式で返します
func (c *GitHubClient) HandleToListCommits(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}
	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}
	page := request.GetInt("page", 1)
	perPage := request.GetInt("per_page", 30)
	sha := request.GetString("sha", "")

	// usecasesレイヤーのサービスを使用
	service := usecases.NewGitHubRepositoryService(c.token)
	result, err := service.HandleToListCommits(owner, repo, page, perPage, sha)
	if err != nil {
		return nil, err
	}

	return returnJSONResult(result)
}

// HandleToSearchRepositories はGitHubリポジトリを検索して、結果をJSON形式で返します
func (c *GitHubClient) HandleToSearchRepositories(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := request.RequireString("query")
	if err != nil {
		return nil, err
	}
	page := request.GetInt("page", 1)
	perPage := request.GetInt("per_page", 30)

	// usecasesレイヤーのサービスを使用
	service := usecases.NewGitHubRepositoryService(c.token)
	result, err := service.HandleToSearchRepositories(query, page, perPage)
	if err != nil {
		return nil, err
	}

	return returnJSONResult(result)
}

// HandleToGetUserRepositories はユーザーのリポジトリ一覧を取得して、結果をJSON形式で返します
func (c *GitHubClient) HandleToGetUserRepositories(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	username, err := request.RequireString("username")
	if err != nil {
		return nil, err
	}

	// パラメータを取得
	sort := request.GetString("sort", "")
	direction := request.GetString("direction", "")
	type_ := request.GetString("type", "")
	perPage := request.GetInt("per_page", 0)
	page := request.GetInt("page", 0)

	// usecasesレイヤーのサービスを使用
	service := usecases.NewGitHubRepositoryService(c.token)
	result, err := service.HandleToGetUserRepositories(username, sort, direction, type_, perPage, page)
	if err != nil {
		return nil, err
	}

	return returnJSONResult(result)
}

// HandleToGetFileContents はリポジトリからファイルの内容を取得して、結果をJSON形式で返します
func (c *GitHubClient) HandleToGetFileContents(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	owner, err := request.RequireString("owner")
	if err != nil {
		return nil, err
	}
	repo, err := request.RequireString("repo")
	if err != nil {
		return nil, err
	}
	path, err := request.RequireString("path")
	if err != nil {
		return nil, err
	}
	branch := request.GetString("branch", "")

	// usecasesレイヤーのサービスを使用
	service := usecases.NewGitHubRepositoryService(c.token)
	result, err := service.HandleToGetFileContents(owner, repo, path, branch)
	if err != nil {
		return nil, err
	}

	return returnJSONResult(result)
}

// SetGitHubRepositoryServer は受け取ったMCPサーバにGitHubリポジトリ用のツールを付与して、そのMCPサーバを返します。
func SetGitHubRepositoryServer(token string, s *server.MCPServer) *server.MCPServer {
	// GitHubクライアントを初期化
	client := NewGitHubClient(token)

	// ツール1: リポジトリ検索
	searchRepositoriesTool := mcp.NewTool("search_repositories",
		mcp.WithDescription("Search for GitHub repositories"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number (default: 1)"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Results per page (default: 30, max: 100)"),
		),
	)
	s.AddTool(searchRepositoriesTool, client.HandleToSearchRepositories)

	// ツール2: ユーザーリポジトリの検索
	getUserRepositoriesTool := mcp.NewTool("get_user_repositories",
		mcp.WithDescription("Get repositories for a specific GitHub user"),
		mcp.WithString("username",
			mcp.Required(),
			mcp.Description("GitHub username"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Results per page (default: 30, max: 100)"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number (default: 1)"),
		),
		mcp.WithString("sort",
			mcp.Description("Sort field: created, updated, pushed, full_name (default: full_name)"),
			mcp.Enum("created", "updated", "pushed", "full_name"),
		),
		mcp.WithString("direction",
			mcp.Description("Sort direction: asc or desc (default: desc)"),
			mcp.Enum("asc", "desc"),
		),
		mcp.WithString("type",
			mcp.Description("Type of repositories to include: all, owner, member, public, private (default: all)"),
			mcp.Enum("all", "owner", "member", "public", "private"),
		),
	)
	s.AddTool(getUserRepositoriesTool, client.HandleToGetUserRepositories)

	// ツール3: ファイル内容の取得
	getFileContentsTool := mcp.NewTool("get_file_contents",
		mcp.WithDescription("Get the contents of a file from a GitHub repository"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("File path within the repository"),
		),
		mcp.WithString("branch",
			mcp.Description("Branch name (default: repository's default branch)"),
		),
	)
	s.AddTool(getFileContentsTool, client.HandleToGetFileContents)

	// ツール4: コミット一覧の取得
	listCommitsTool := mcp.NewTool("list_commits",
		mcp.WithDescription("Get list of commits of a branch in a GitHub repository"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number (default: 1)"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Results per page (default: 30, max: 100)"),
		),
		mcp.WithString("sha",
			mcp.Description("SHA or branch name to start listing commits from"),
		),
	)
	s.AddTool(listCommitsTool, client.HandleToListCommits)

	return s
}
