package usecases

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// #==============================================================#
// ##          GitHubService                                     ##
// #==============================================================#
// GitHubService はGitHub操作のインターフェース
type GitHubService interface {
	CreateGitHubClient(ctx context.Context, token string) GitHubClient
	GetRepoInfo(ctx context.Context, client GitHubClient, isThreading bool, username string) ([]RepoInfo, error)
}

// GitHubServiceImpl はGitHubServiceの実装
type GitHubServiceImpl struct{}

// NewGitHubService は新しいGitHubServiceインスタンスを作成する
func NewGitHubService() GitHubService {
	return &GitHubServiceImpl{}
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
		HttpUrl:          repo.GetHTMLURL(),
		Language:         repo.GetLanguage(),
		Languages:        languages,
		RepoCreatedAt:    createdAt,
		RepoUpdatedAt:    updatedAt,
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

// GetRepoInfo はリポジトリ情報を取得する
func (gs *GitHubServiceImpl) GetRepoInfo(ctx context.Context, client GitHubClient, isThreading bool, username string) ([]RepoInfo, error) {
	repoType := "all"
	repos, err := gs.fetchRepositories(ctx, client, repoType, username)
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

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

// #==============================================================#
// ##         GitHubClient                                       ##
// #==============================================================#
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

// CreateGitHubClient はGitHubクライアントを作成する
func (gs *GitHubServiceImpl) CreateGitHubClient(ctx context.Context, token string) GitHubClient {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oc := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(oc)
	return &GitHubClientAdapter{Client: client}
}

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
