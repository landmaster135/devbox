package git

import (
	"fmt"
	"strings"
)

// Client はGit操作を行うクライアント
type Client struct {
	workingDir string
	executor   GitCommandExecutor
}

// NewClient は新しいGitクライアントを作成する
func NewClient(workingDir string) *Client {
	return &Client{
		workingDir: workingDir,
		executor:   NewStandardGitExecutor(),
	}
}

// NewClientWithExecutor は依存性注入可能なGitクライアントを作成する
func NewClientWithExecutor(workingDir string, executor GitCommandExecutor) *Client {
	return &Client{
		workingDir: workingDir,
		executor:   executor,
	}
}

// GetCommitHistory はコミット履歴を取得する
func (c *Client) GetCommitHistory(keyword, since, until string) (string, error) {
	// 基本的なgit logコマンドの引数を構築
	args := []string{
		"log",
		"--graph",
		"--pretty=format:%h -%d %s (%cr) <%an>",
		"--abbrev-commit",
	}

	// キーワード検索が指定されている場合
	if keyword != "" {
		args = append(args, "--grep="+keyword)
	}

	// 開始日が指定されている場合
	if since != "" {
		args = append(args, "--since="+since)
	}

	// 終了日が指定されている場合
	if until != "" {
		args = append(args, "--until="+until)
	}

	// Gitコマンドを実行
	output, err := c.executor.Execute(c.workingDir, args...)
	if err != nil {
		return "", fmt.Errorf("コミット履歴の取得に失敗しました: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// IsValidGitRepository は指定されたディレクトリが有効なGitリポジトリかどうかを確認する
func (c *Client) IsValidGitRepository() error {
	_, err := c.executor.Execute(c.workingDir, "rev-parse", "--git-dir")
	if err != nil {
		return fmt.Errorf("指定されたディレクトリは有効なGitリポジトリではありません: %s", c.workingDir)
	}
	return nil
}
