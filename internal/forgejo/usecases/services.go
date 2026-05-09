package usecases

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	forgejo "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"

	commonUsecases "github.com/landmaster135/devbox/internal/forgejo/usecases/common"
	projectList "github.com/landmaster135/devbox/internal/forgejo/usecases/operations/project_list"
	repoList "github.com/landmaster135/devbox/internal/forgejo/usecases/operations/repo_list"
)

const (
	defaultHTTPTimeout = 30 * time.Second
	pullsPageSize      = 100
)

// ServiceOptions はForgejoサービス作成時の依存を保持します。
type ServiceOptions struct {
	Host       string
	Username   string
	Token      string
	HTTPClient *http.Client
	// ReposWorkers はrepo list 取得時の同時実行ワーカー数です。
	ReposWorkers int
}

// Service はCLI向けのForgejo処理を担当します。
type Service struct {
	client       *forgejo.Client
	host         string
	username     string
	token        string
	httpClient   commonUsecases.HTTPClient
	reposWorkers int
}

var (
	newRepoListOperation = func(dependencies repoListDependencies) repoListOperation {
		return repoList.NewService(repoList.Options{
			Client:       dependencies.Client,
			Host:         dependencies.Host,
			Username:     dependencies.Username,
			Token:        dependencies.Token,
			HTTPClient:   dependencies.HTTPClient,
			ReposWorkers: dependencies.ReposWorkers,
		})
	}
	newProjectListOperation = func(dependencies projectListDependencies) projectListOperation {
		return projectList.NewService(projectList.Options{
			Client:     dependencies.Client,
			Host:       dependencies.Host,
			Username:   dependencies.Username,
			Token:      dependencies.Token,
			HTTPClient: dependencies.HTTPClient,
		})
	}
)

// NewService はForgejo用Serviceを作成します。
func NewService(options ServiceOptions) (*Service, error) {
	host := normalizeHost(options.Host)
	if host == "" {
		return nil, fmt.Errorf("host is required")
	}
	client := &http.Client{Timeout: defaultHTTPTimeout}
	if options.HTTPClient != nil {
		client = options.HTTPClient
	}

	forgejoClient, err := forgejo.NewClient(
		host,
		forgejo.SetHTTPClient(client),
		forgejo.SetToken(options.Token),
		forgejo.SetForgejoVersion(""),
	)
	if err != nil {
		return nil, err
	}

	return &Service{
		client:       forgejoClient,
		host:         host,
		username:     strings.TrimSpace(options.Username),
		token:        options.Token,
		httpClient:   client,
		reposWorkers: options.ReposWorkers,
	}, nil
}

// ListRepos はリポジトリ一覧を取得します。
func (s *Service) ListRepos() ([]RepoRecord, error) {
	operation := newRepoListOperation(repoListDependencies{
		Client:       s.client,
		Host:         s.host,
		Username:     s.username,
		Token:        s.token,
		HTTPClient:   s.httpClient,
		ReposWorkers: s.reposWorkers,
	})
	return operation.Execute()
}

// ListProjects はプロジェクト一覧を取得します。
func (s *Service) ListProjects() ([]ProjectRecord, error) {
	operation := newProjectListOperation(projectListDependencies{
		Client:     s.client,
		Host:       s.host,
		Username:   s.username,
		Token:      s.token,
		HTTPClient: s.httpClient,
	})
	return operation.Execute()
}

func normalizeHost(host string) string {
	host = strings.TrimSpace(host)
	if host == "" {
		return ""
	}
	if !strings.Contains(host, "://") {
		host = "https://" + host
	}
	return strings.TrimRight(host, "/")
}
