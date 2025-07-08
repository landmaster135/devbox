package git

import (
	"fmt"
	"path/filepath"
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

// GetRepositoryName はリポジトリ名を取得する
func (c *Client) GetRepositoryName() (string, error) {
	output, err := c.executor.Execute(c.workingDir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("リポジトリのルートディレクトリを取得できませんでした: %w", err)
	}

	repoPath := strings.TrimSpace(string(output))
	repoName := filepath.Base(repoPath)

	return repoName, nil
}

// GetCurrentBranch は現在のブランチ名を取得する
func (c *Client) GetCurrentBranch() (string, error) {
	output, err := c.executor.Execute(c.workingDir, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("現在のブランチを取得できませんでした: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetLatestCommitHash は最新のコミットハッシュを取得する
func (c *Client) GetLatestCommitHash() (string, error) {
	output, err := c.executor.Execute(c.workingDir, "rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("最新のコミットハッシュを取得できませんでした: %w", err)
	}

	hash := strings.TrimSpace(string(output))
	// 短縮形で返す（最初の8文字）
	if len(hash) > 8 {
		hash = hash[:8]
	}

	return hash, nil
}

// GetDiff は差分を取得する
func (c *Client) GetDiff(stagedOnly bool) (string, error) {
	var output []byte
	var err error

	if stagedOnly {
		// ステージング済みの差分のみ
		output, err = c.executor.Execute(c.workingDir, "diff", "--cached")
	} else {
		// 全ての差分（ステージング済み + 未ステージング）
		output, err = c.executor.Execute(c.workingDir, "diff", "HEAD")
	}

	if err != nil {
		return "", fmt.Errorf("差分を取得できませんでした: %w", err)
	}

	return string(output), nil
}

// GetStatus はgit statusの情報を取得する
func (c *Client) GetStatus() (string, error) {
	output, err := c.executor.Execute(c.workingDir, "status", "--porcelain")
	if err != nil {
		return "", fmt.Errorf("ステータスを取得できませんでした: %w", err)
	}

	return string(output), nil
}

// GetNewFiles は新規ファイルのリストを取得する
func (c *Client) GetNewFiles(stagedOnly bool) ([]string, error) {
	status, err := c.GetStatus()
	if err != nil {
		return nil, err
	}

	var newFiles []string
	lines := strings.Split(status, "\n")

	for _, line := range lines {
		if len(line) >= 3 {
			statusCode := line[:2]
			filename := strings.TrimSpace(line[3:])

			if stagedOnly {
				// ステージング済み新規ファイルのみ（A: added）
				if statusCode == "A " {
					newFiles = append(newFiles, filename)
				}
			} else {
				// 全ての新規ファイル（A: added, ??: untracked）
				if statusCode == "A " || statusCode == "??" {
					newFiles = append(newFiles, filename)
				}
			}
		}
	}

	return newFiles, nil
}

// GetDeletedFiles は削除されたファイルのリストを取得する
func (c *Client) GetDeletedFiles(stagedOnly bool) ([]string, error) {
	status, err := c.GetStatus()
	if err != nil {
		return nil, err
	}

	var deletedFiles []string
	lines := strings.Split(status, "\n")

	for _, line := range lines {
		if len(line) >= 3 {
			statusCode := line[:2]
			filename := strings.TrimSpace(line[3:])

			if stagedOnly {
				// ステージング済み削除ファイルのみ（D: deleted）
				if statusCode == "D " {
					deletedFiles = append(deletedFiles, filename)
				}
			} else {
				// 全ての削除ファイル（D: deleted）
				if statusCode == "D " || statusCode == " D" {
					deletedFiles = append(deletedFiles, filename)
				}
			}
		}
	}

	return deletedFiles, nil
}

// GetModifiedFilesCount は変更されたファイル数を取得する
func (c *Client) GetModifiedFilesCount(stagedOnly bool) (int, error) {
	status, err := c.GetStatus()
	if err != nil {
		return 0, err
	}

	count := 0
	lines := strings.Split(status, "\n")

	for _, line := range lines {
		if len(line) >= 3 {
			statusCode := line[:2]

			if stagedOnly {
				// ステージング済み変更ファイルのみ（M: modified）
				if statusCode == "M " {
					count++
				}
			} else {
				// 全ての変更ファイル（M: modified）
				if statusCode == "M " || statusCode == " M" || statusCode == "MM" {
					count++
				}
			}
		}
	}

	return count, nil
}
