package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// RepoInfo はリポジトリ情報を保持する構造体
type RepoInfo struct {
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	IsPrivate        bool           `json:"is_private"`
	HtmlUrl          string         `json:"html_url"`
	Language         string         `json:"language"`
	Languages        map[string]int `json:"languages"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
	StargazersCount  int32          `json:"stargazers_count"`
	ForksCount       int32          `json:"forks_count"`
	IssuesCount      int32          `json:"issues_count"`
	PullsCount       int32          `json:"pulls_count"`
	Size             int32          `json:"size"`
	SubscribersCount int32          `json:"subscribers_count"`
	IsArchived       bool           `json:"is_archived"`
}

// GitInfoService はGit情報取得サービスのインターフェース
type GitInfoService interface {
	GetRepositoryInfo(service, token string) (string, error)
	GetAndSaveRepositoryInfo(service, token, filePath string) (string, error)
}

// GitHubServiceImpl はGitHubServiceの実装
type GitHubServiceImpl struct{}

// FileWriter はファイル書き込み操作のインターフェース
type FileWriter interface {
	WriteToFile(filePath, content string) error
	EnsureDirectory(dirPath string) error
}

// FileWriterImpl はFileWriterの実装
type FileWriterImpl struct{}

// Service はGitInfoServiceの実装
type Service struct {
	githubService GitHubService
	fileWriter    FileWriter
}

// GitHubService はGitHub操作のインターフェース
type GitHubService interface {
	CreateGitHubClient(ctx context.Context, token string) GitHubClient
	GetRepoInfo(ctx context.Context, client GitHubClient, isThreading bool, username string) ([]RepoInfo, error)
}

// GitHubClient はGitHub APIクライアントのインターフェース
type GitHubClient interface {
	ListRepositories(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error)
	ListRepoLanguages(ctx context.Context, owner string, repo string) (map[string]int, *github.Response, error)
	ListPullRequests(ctx context.Context, user string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
	GetUser(ctx context.Context, user string) (*github.User, *github.Response, error)
}

// GitHubClientAdapter はgithub.Clientをラップするアダプター
type GitHubClientAdapter struct {
	*github.Client
}

// NewService は新しいServiceインスタンスを作成する
func NewService() GitInfoService {
	return &Service{
		githubService: &GitHubServiceImpl{},
		fileWriter:    &FileWriterImpl{},
	}
}

// GetRepositoryInfo は指定されたサービスからリポジトリ情報を取得する
func (s *Service) GetRepositoryInfo(service, token string) (string, error) {
	if service != "github" {
		return "", fmt.Errorf("サポートされていないサービスです: %s", service)
	}

	ctx := context.Background()
	client := s.githubService.CreateGitHubClient(ctx, token)

	// 認証されたユーザー情報を取得
	user, _, err := client.GetUser(ctx, "")
	if err != nil {
		return "", fmt.Errorf("ユーザー情報の取得に失敗しました: %v", err)
	}

	username := user.GetLogin()
	repos, err := s.githubService.GetRepoInfo(ctx, client, true, username)
	if err != nil {
		return "", fmt.Errorf("リポジトリ情報の取得に失敗しました: %v", err)
	}

	// JSON形式で結果を返却
	result, err := json.MarshalIndent(repos, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON変換に失敗しました: %v", err)
	}

	return string(result), nil
}

// CreateGitHubClient はGitHubクライアントを作成する
func (gs *GitHubServiceImpl) CreateGitHubClient(ctx context.Context, token string) GitHubClient {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oc := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(oc)
	return &GitHubClientAdapter{Client: client}
}

// GetRepoInfo はリポジトリ情報を取得する
func (gs *GitHubServiceImpl) GetRepoInfo(ctx context.Context, client GitHubClient, isThreading bool, username string) ([]RepoInfo, error) {
	repos, err := gs.fetchRepositories(ctx, client, "owner", username)
	if err != nil {
		return nil, err
	}

	var results []RepoInfo
	if isThreading {
		var wg sync.WaitGroup
		mu := &sync.Mutex{}
		for _, repo := range repos {
			wg.Add(1)
			go func(repo *github.Repository) {
				defer wg.Done()
				repoInfo, err := gs.getRepoInfoInFormat(ctx, client, repo)
				if err != nil {
					// エラーログを出力するが、処理は継続
					fmt.Printf("リポジトリ %s の情報取得でエラー: %v\n", repo.GetName(), err)
					return
				}
				mu.Lock()
				results = append(results, *repoInfo)
				mu.Unlock()
			}(repo)
		}
		wg.Wait()
	} else {
		for _, repo := range repos {
			repoInfo, err := gs.getRepoInfoInFormat(ctx, client, repo)
			if err != nil {
				return nil, err
			}
			results = append(results, *repoInfo)
		}
	}

	return results, nil
}

// fetchRepositories はリポジトリ一覧を取得する
func (gs *GitHubServiceImpl) fetchRepositories(ctx context.Context, client GitHubClient, repoType string, username string) ([]*github.Repository, error) {
	if repoType == "" {
		return nil, fmt.Errorf("リポジトリタイプが空です")
	}

	var allRepos []*github.Repository
	opt := &github.RepositoryListOptions{
		Type:        repoType,
		Sort:        "full_name",
		Direction:   "asc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, resp, err := client.ListRepositories(ctx, username, opt)
		if err != nil {
			if rateLimitErr, ok := err.(*github.RateLimitError); ok {
				return nil, fmt.Errorf("GitHub APIのレート制限に達しました: %v", rateLimitErr)
			}
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}

// getRepoInfoInFormat はリポジトリ情報をRepoInfo形式で取得する
func (gs *GitHubServiceImpl) getRepoInfoInFormat(ctx context.Context, client GitHubClient, repo *github.Repository) (*RepoInfo, error) {
	// プルリクエスト一覧を取得
	pulls, err := gs.getRepoPulls(ctx, client, repo, "all")
	if err != nil {
		return nil, err
	}

	// 言語情報を取得
	languages, _, err := client.ListRepoLanguages(ctx, repo.GetOwner().GetLogin(), repo.GetName())
	if err != nil {
		return nil, err
	}

	// 時刻フォーマット
	createdAt := formatGitHubTime(repo.GetCreatedAt())
	updatedAt := formatGitHubTime(repo.GetUpdatedAt())

	repoInfo := &RepoInfo{
		Name:             repo.GetName(),
		Description:      repo.GetDescription(),
		IsPrivate:        repo.GetPrivate(),
		HtmlUrl:          repo.GetHTMLURL(),
		Language:         repo.GetLanguage(),
		Languages:        languages,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		StargazersCount:  int32(repo.GetStargazersCount()),
		ForksCount:       int32(repo.GetForksCount()),
		IssuesCount:      int32(repo.GetOpenIssuesCount()),
		PullsCount:       int32(len(pulls)),
		Size:             int32(repo.GetSize()),
		SubscribersCount: int32(repo.GetSubscribersCount()),
		IsArchived:       repo.GetArchived(),
	}

	return repoInfo, nil
}

// getRepoPulls はプルリクエスト一覧を取得する
func (gs *GitHubServiceImpl) getRepoPulls(ctx context.Context, client GitHubClient, repo *github.Repository, state string) ([]*github.PullRequest, error) {
	if state == "" {
		return nil, fmt.Errorf("stateが空です")
	}

	opt := &github.PullRequestListOptions{State: state}
	pulls, _, err := client.ListPullRequests(ctx, repo.GetOwner().GetLogin(), repo.GetName(), opt)
	if err != nil {
		return nil, err
	}

	return pulls, nil
}

// formatGitHubTime はGitHubのTimestampを文字列に変換する
func formatGitHubTime(t github.Timestamp) string {
	return t.Time.Format("2006-01-02T15:04:05Z")
}

// GetAndSaveRepositoryInfo はリポジトリ情報を取得し、ファイルに保存する
func (s *Service) GetAndSaveRepositoryInfo(service, token, filePath string) (string, error) {
	// リポジトリ情報を取得
	result, err := s.GetRepositoryInfo(service, token)
	if err != nil {
		return "", err
	}

	// ファイルに保存
	err = s.fileWriter.WriteToFile(filePath, result)
	if err != nil {
		return "", fmt.Errorf("ファイル保存に失敗しました: %v", err)
	}

	return result, nil
}

// WriteToFile はファイルに内容を書き込む
func (fw *FileWriterImpl) WriteToFile(filePath, content string) error {
	// ディレクトリを確保
	dir := filepath.Dir(filePath)
	if err := fw.EnsureDirectory(dir); err != nil {
		return err
	}

	// ファイルに書き込み
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました: %v", err)
	}

	return nil
}

// EnsureDirectory はディレクトリが存在することを確認し、必要に応じて作成する
func (fw *FileWriterImpl) EnsureDirectory(dirPath string) error {
	if dirPath == "." || dirPath == "" {
		return nil
	}

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("ディレクトリの作成に失敗しました: %v", err)
	}

	return nil
}

// GitHubClientAdapterのメソッド実装

func (g *GitHubClientAdapter) ListRepositories(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	return g.Repositories.List(ctx, user, opts)
}

func (g *GitHubClientAdapter) ListRepoLanguages(ctx context.Context, owner string, repo string) (map[string]int, *github.Response, error) {
	return g.Repositories.ListLanguages(ctx, owner, repo)
}

func (g *GitHubClientAdapter) ListPullRequests(ctx context.Context, user string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
	return g.PullRequests.List(ctx, user, repo, opts)
}

func (g *GitHubClientAdapter) GetUser(ctx context.Context, user string) (*github.User, *github.Response, error) {
	return g.Users.Get(ctx, user)
}
