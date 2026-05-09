package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	forgejo "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

const (
	defaultHTTPTimeout  = 30 * time.Second
	timeFormatDate      = time.RFC3339
	pullsPageSize       = 100
	defaultReposWorkers = 4
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

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
	httpClient   httpClient
	reposWorkers int
}

// RepoRecord は repo list の出力レコードです。
type RepoRecord struct {
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	IsPrivate        bool               `json:"is_private"`
	HTTPURL          string             `json:"http_url"`
	IssuesCount      int                `json:"issues_count"`
	OpenPullsCount   int                `json:"open_pulls_count"`
	ClosedPullsCount int                `json:"closed_pulls_count"`
	ForksCount       int                `json:"forks_count"`
	StargazersCount  int                `json:"stargazers_count"`
	SubscribersCount int                `json:"subscribers_count"`
	Language         string             `json:"language"`
	Languages        map[string]float64 `json:"languages"`
	Size             int                `json:"size"`
	RepoCreatedAt    string             `json:"repo_created_at"`
	RepoUpdatedAt    string             `json:"repo_updated_at"`
	IsArchived       bool               `json:"is_archived"`
	Tags             string             `json:"tags"`
}

// ProjectRecord は project list の出力レコードです。
type ProjectRecord struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsPrivate    bool   `json:"is_private"`
	IsArchived   bool   `json:"is_archived"`
	RepoFullName string `json:"repo_full_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type projectResponse struct {
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"is_private"`
	Archived    bool      `json:"is_archived"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

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
		reposWorkers: resolveReposWorkers(options.ReposWorkers),
	}, nil
}

// ListRepos はリポジトリ一覧を取得します。
func (s *Service) ListRepos() ([]RepoRecord, error) {
	repos, _, err := s.client.ListMyRepos(forgejo.ListReposOptions{})
	if err != nil {
		return nil, fmt.Errorf("リポジトリ一覧の取得に失敗しました: %w", err)
	}

	type repoJob struct {
		index    int
		owner    string
		repoName string
		repo     *forgejo.Repository
	}

	type repoResult struct {
		index  int
		record RepoRecord
	}

	repoJobs := make([]repoJob, 0, len(repos))
	for _, repo := range repos {
		if repo == nil {
			continue
		}
		owner, repoName := repositoryName(repo)
		if owner == "" {
			owner = s.username
		}
		if repoName == "" {
			continue
		}
		repoJobs = append(repoJobs, repoJob{
			index:    len(repoJobs),
			owner:    owner,
			repoName: repoName,
			repo:     repo,
		})
	}

	records := make([]RepoRecord, len(repoJobs))
	if len(repoJobs) == 0 {
		return records, nil
	}

	jobs := make(chan repoJob, len(repoJobs))
	results := make(chan repoResult, len(repoJobs))
	errCh := make(chan error, 1)
	done := make(chan struct{})
	var doneOnce sync.Once
	var wg sync.WaitGroup

	reportErr := func(fetchErr error) {
		if fetchErr == nil {
			return
		}
		doneOnce.Do(func() {
			close(done)
		})
		select {
		case errCh <- fetchErr:
		default:
		}
	}

	worker := func() {
		defer wg.Done()
		for job := range jobs {
			select {
			case <-done:
				return
			default:
			}
			record, err := s.fetchRepoRecord(job.owner, job.repoName, job.repo)
			if err != nil {
				reportErr(err)
				return
			}
			select {
			case <-done:
				return
			case results <- repoResult{
				index:  job.index,
				record: record,
			}:
			}
		}
	}

	for _, job := range repoJobs {
		jobs <- job
	}
	close(jobs)

	workerCount := s.reposWorkers
	if workerCount > len(repoJobs) {
		workerCount = len(repoJobs)
	}
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker()
	}

	wg.Wait()
	close(results)

	for result := range results {
		records[result.index] = result.record
	}

	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	return records, nil
}

// ListProjects はプロジェクト一覧を取得します。
func (s *Service) ListProjects() ([]ProjectRecord, error) {
	repos, _, err := s.client.ListUserRepos(s.username, forgejo.ListReposOptions{})
	if err != nil {
		return nil, fmt.Errorf("ユーザーリポジトリ一覧の取得に失敗しました: %w", err)
	}

	records := make([]ProjectRecord, 0)
	notFoundCount := 0
	for _, repo := range repos {
		if repo == nil {
			continue
		}
		owner, repoName := repositoryName(repo)
		if owner == "" {
			owner = s.username
		}
		if repoName == "" {
			continue
		}
		projects, err := s.fetchProjects(owner, repoName)
		if err != nil {
			if isNotFoundError(err) {
				notFoundCount++
				continue
			}
			return nil, fmt.Errorf("projectsの取得に失敗しました (%s/%s): %w", owner, repoName, err)
		}
		for _, project := range projects {
			records = append(records, ProjectRecord{
				Name:         firstNonEmptyString(project.Name, project.Title),
				Description:  project.Description,
				IsPrivate:    project.IsPrivate,
				IsArchived:   project.Archived,
				RepoFullName: owner + "/" + repoName,
				CreatedAt:    formatDate(project.CreatedAt),
				UpdatedAt:    formatDate(project.UpdatedAt),
			})
		}
	}

	if len(records) == 0 && len(repos) > 0 && notFoundCount == len(repos) {
		return nil, fmt.Errorf("project list API is not supported on this server")
	}
	return records, nil
}

func (s *Service) fetchTopics(owner, repo string) ([]string, error) {
	topics, _, err := s.client.ListRepoTopics(owner, repo, forgejo.ListRepoTopicsOptions{})
	return topics, err
}

func resolveReposWorkers(raw int) int {
	if raw <= 0 {
		return defaultReposWorkers
	}
	return raw
}

func (s *Service) fetchLanguages(owner, repo string) (map[string]float64, error) {
	requestURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/languages", s.host, url.PathEscape(owner), url.PathEscape(repo))
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	s.setAuthHeader(req)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &requestError{
			status: resp.StatusCode,
			body:   string(body),
		}
	}

	var languages map[string]float64
	if err := json.Unmarshal(body, &languages); err != nil {
		return nil, err
	}
	return languages, nil
}

func (s *Service) fetchRepoRecord(owner, repoName string, repo *forgejo.Repository) (RepoRecord, error) {
	topics, err := s.fetchTopics(owner, repoName)
	if err != nil {
		return RepoRecord{}, fmt.Errorf("topicsの取得に失敗しました (%s/%s): %w", owner, repoName, err)
	}
	languages, err := s.fetchLanguages(owner, repoName)
	if err != nil {
		return RepoRecord{}, fmt.Errorf("languagesの取得に失敗しました (%s/%s): %w", owner, repoName, err)
	}
	closedPullsCount, err := s.fetchClosedPullsCount(owner, repoName)
	if err != nil {
		return RepoRecord{}, fmt.Errorf("pullsの取得に失敗しました (%s/%s): %w", owner, repoName, err)
	}

	return RepoRecord{
		Name:             repoName,
		Description:      repo.Description,
		IsPrivate:        repo.Private,
		HTTPURL:          repo.HTMLURL,
		IssuesCount:      repo.OpenIssues,
		OpenPullsCount:   repo.OpenPulls,
		ClosedPullsCount: closedPullsCount,
		ForksCount:       repo.Forks,
		StargazersCount:  repo.Stars,
		SubscribersCount: repo.Watchers,
		Language:         primaryLanguage(languages),
		Languages:        languages,
		Size:             repo.Size,
		RepoCreatedAt:    formatDate(repo.Created),
		RepoUpdatedAt:    formatDate(repo.Updated),
		IsArchived:       repo.Archived,
		Tags:             strings.Join(topics, ","),
	}, nil
}

func (s *Service) fetchClosedPullsCount(owner, repo string) (int, error) {
	firstPulls, response, err := s.client.ListRepoPullRequests(owner, repo, forgejo.ListPullRequestsOptions{
		ListOptions: forgejo.ListOptions{
			Page:     1,
			PageSize: 1,
		},
		State: forgejo.StateClosed,
	})
	if err != nil {
		return 0, err
	}

	if response == nil {
		return len(firstPulls), nil
	}

	totalCountFromHeader := strings.TrimSpace(response.Header.Get("X-Total-Count"))
	if totalCountFromHeader != "" {
		totalClosedPulls, err := strconv.Atoi(totalCountFromHeader)
		if err == nil {
			return totalClosedPulls, nil
		}
	}
	if len(firstPulls) == 0 {
		return 0, nil
	}

	if response.LastPage > 1 {
		return response.LastPage, nil
	}

	return len(firstPulls), nil
}

func (s *Service) fetchProjects(owner, repo string) ([]projectResponse, error) {
	query := url.Values{}
	query.Set("state", "all")
	requestURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/projects", s.host, url.PathEscape(owner), url.PathEscape(repo))
	if q := query.Encode(); q != "" {
		requestURL = requestURL + "?" + q
	}

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	s.setAuthHeader(req)
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
		return nil, &requestError{
			status: resp.StatusCode,
			body:   string(body),
		}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &requestError{
			status: resp.StatusCode,
			body:   string(body),
		}
	}

	projects, err := decodeProjects(body)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func decodeProjects(data []byte) ([]projectResponse, error) {
	var projects []projectResponse
	if err := json.Unmarshal(data, &projects); err == nil {
		return projects, nil
	}

	var wrapper struct {
		Data     []projectResponse `json:"data"`
		Projects []projectResponse `json:"projects"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}
	if len(wrapper.Data) > 0 {
		return wrapper.Data, nil
	}
	return wrapper.Projects, nil
}

func (s *Service) setAuthHeader(req *http.Request) {
	if req == nil {
		return
	}
	if strings.TrimSpace(s.token) != "" {
		req.Header.Set("Authorization", "token "+s.token)
	}
	req.Header.Set("Accept", "application/json")
}

func repositoryName(repo *forgejo.Repository) (owner string, name string) {
	name = strings.TrimSpace(repo.Name)
	if repo.Owner != nil {
		owner = strings.TrimSpace(repo.Owner.UserName)
	}
	if owner == "" {
		parts := strings.SplitN(repo.FullName, "/", 2)
		if len(parts) == 2 {
			owner = strings.TrimSpace(parts[0])
		}
	}
	return owner, name
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

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(timeFormatDate)
}

func primaryLanguage(langs map[string]float64) string {
	if len(langs) == 0 {
		return ""
	}
	type item struct {
		name  string
		value float64
	}
	items := make([]item, 0, len(langs))
	for name, value := range langs {
		items = append(items, item{name: name, value: value})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].value == items[j].value {
			return items[i].name < items[j].name
		}
		return items[i].value > items[j].value
	})
	return items[0].name
}

func firstNonEmptyString(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if requestErr, ok := err.(*requestError); ok {
		return requestErr.status == http.StatusNotFound
	}
	return false
}

type requestError struct {
	status int
	body   string
}

func (e *requestError) Error() string {
	return fmt.Sprintf("request failed with status %d: %s", e.status, e.body)
}
