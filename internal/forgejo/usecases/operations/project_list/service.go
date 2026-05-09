package project_list

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	forgejo "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"

	commonUsecases "github.com/landmaster135/devbox/internal/forgejo/usecases/common"
)

// Service は project list operation の実装です。
type Service struct {
	client     *forgejo.Client
	host       string
	username   string
	token      string
	httpClient commonUsecases.HTTPClient
}

// Options は project list operation の依存です。
type Options struct {
	Client     *forgejo.Client
	Host       string
	Username   string
	Token      string
	HTTPClient commonUsecases.HTTPClient
}

// NewService は project list operation サービスを作成します。
func NewService(options Options) *Service {
	return &Service{
		client:     options.Client,
		host:       options.Host,
		username:   options.Username,
		token:      options.Token,
		httpClient: options.HTTPClient,
	}
}

// Execute は project list を実行します。
func (s *Service) Execute() ([]commonUsecases.ProjectRecord, error) {
	repos, _, err := s.client.ListUserRepos(s.username, forgejo.ListReposOptions{})
	if err != nil {
		return nil, fmt.Errorf("ユーザーリポジトリ一覧の取得に失敗しました: %w", err)
	}

	records := make([]commonUsecases.ProjectRecord, 0)
	notFoundCount := 0
	for _, repo := range repos {
		if repo == nil {
			continue
		}
		owner, repoName := commonUsecases.RepositoryName(repo)
		if owner == "" {
			owner = s.username
		}
		if repoName == "" {
			continue
		}

		projects, fetchErr := s.fetchProjects(owner, repoName)
		if fetchErr != nil {
			if commonUsecases.IsNotFoundError(fetchErr) {
				notFoundCount++
				continue
			}
			return nil, fmt.Errorf("projectsの取得に失敗しました (%s/%s): %w", owner, repoName, fetchErr)
		}

		for _, project := range projects {
			records = append(records, commonUsecases.ProjectRecord{
				Name:         commonUsecases.FirstNonEmptyString(project.Name, project.Title),
				Description:  project.Description,
				IsPrivate:    project.IsPrivate,
				IsArchived:   project.Archived,
				RepoFullName: owner + "/" + repoName,
				CreatedAt:    commonUsecases.FormatDate(project.CreatedAt),
				UpdatedAt:    commonUsecases.FormatDate(project.UpdatedAt),
			})
		}
	}

	if len(records) == 0 && len(repos) > 0 && notFoundCount == len(repos) {
		return nil, fmt.Errorf("project list API is not supported on this server")
	}
	return records, nil
}

func (s *Service) fetchProjects(owner, repo string) ([]commonUsecases.ProjectResponse, error) {
	query := url.Values{}
	query.Set("state", "all")
	requestURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/projects", s.host, url.PathEscape(owner), url.PathEscape(repo))
	if encoded := query.Encode(); encoded != "" {
		requestURL = requestURL + "?" + encoded
	}

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	commonUsecases.SetAuthHeader(req, s.token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, &commonUsecases.RequestError{
			Status: resp.StatusCode,
			Body:   string(body),
		}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &commonUsecases.RequestError{
			Status: resp.StatusCode,
			Body:   string(body),
		}
	}

	return commonUsecases.DecodeProjects(body)
}
