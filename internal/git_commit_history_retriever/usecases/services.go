package usecases

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/git_commit_history_retriever/config"
	"github.com/landmaster135/devbox/internal/git_commit_history_retriever/git"
)

// GitCommitHistoryService はGitコミット履歴を取得するサービス
type GitCommitHistoryService struct {
	gitClient *git.Client
	config    *config.Config
}

// NewGitCommitHistoryService は新しいGitCommitHistoryServiceを作成する
func NewGitCommitHistoryService(gitDir string, cfg *config.Config) *GitCommitHistoryService {
	return &GitCommitHistoryService{
		gitClient: git.NewClient(gitDir),
		config:    cfg,
	}
}

// GetCommitHistory はコミット履歴を取得して見出し付きで返す
func (s *GitCommitHistoryService) GetCommitHistory() (string, error) {
	// Gitリポジトリの有効性を確認
	if err := s.gitClient.IsValidGitRepository(); err != nil {
		return "", err
	}

	// コミット履歴を取得
	history, err := s.gitClient.GetCommitHistory(s.config.Keyword, s.config.Since, s.config.Until)
	if err != nil {
		return "", err
	}

	// 結果が空の場合
	if history == "" {
		return fmt.Sprintf("%s\n指定された条件に一致するコミットが見つかりませんでした。", config.HeaderCommitHistory), nil
	}

	// 見出し付きで返す
	return fmt.Sprintf("%s\n%s", config.HeaderCommitHistory, history), nil
}

// GetCommitHistoryWithDetails はコミット履歴と詳細を取得して統合出力する
func (s *GitCommitHistoryService) GetCommitHistoryWithDetails() (string, error) {
	// Gitリポジトリの有効性を確認
	if err := s.gitClient.IsValidGitRepository(); err != nil {
		return "", err
	}

	// コミット履歴を取得
	history, err := s.gitClient.GetCommitHistory(s.config.Keyword, s.config.Since, s.config.Until)
	if err != nil {
		return "", err
	}

	// 結果が空の場合
	if history == "" {
		return fmt.Sprintf("%s\n指定された条件に一致するコミットが見つかりませんでした。", config.HeaderCommitHistory), nil
	}

	// コミットハッシュを抽出
	hashes := s.gitClient.ExtractCommitHashes(history)

	// コミット詳細を取得
	details, err := s.gitClient.GetCommitDetails(hashes)
	if err != nil {
		return "", err
	}

	// 統合出力を作成
	result := fmt.Sprintf("%s\n%s", config.HeaderCommitHistory, history)

	if details != "" {
		result += fmt.Sprintf("\n\n%s\n%s", config.HeaderCommitDetails, details)
	}

	return result, nil
}
