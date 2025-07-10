package git

import (
	"fmt"
	"regexp"
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

// ExtractCommitHashes はコミット履歴からコミットハッシュを抽出する
func (c *Client) ExtractCommitHashes(history string) []string {
	// git logの出力形式 "%h -%d %s (%cr) <%an>" からハッシュを抽出
	// * の後にある7文字のコミットハッシュを抽出する正規表現
	// マージコミットやブランチ表示を考慮して、* の後に任意の文字（スペース、パイプなど）が続く場合も対応
	hashPattern := regexp.MustCompile(`\*\s*[|\s]*([a-f0-9]{7})\s+-`)

	matches := hashPattern.FindAllStringSubmatch(history, -1)
	var hashes []string

	for _, match := range matches {
		if len(match) > 1 {
			hashes = append(hashes, match[1])
		}
	}

	return hashes
}

// GetCommitDetails はコミットハッシュリストから詳細情報を取得する
func (c *Client) GetCommitDetails(commitHashes []string) (string, error) {
	if len(commitHashes) == 0 {
		return "", nil
	}

	// git show --stat コマンドの引数を構築
	args := append([]string{"show", "--stat"}, commitHashes...)

	// Gitコマンドを実行
	output, err := c.executor.Execute(c.workingDir, args...)
	if err != nil {
		return "", fmt.Errorf("コミット詳細の取得に失敗しました: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
