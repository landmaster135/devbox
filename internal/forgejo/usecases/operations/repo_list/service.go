package repo_list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	forgejo "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"

	commonUsecases "github.com/landmaster135/devbox/internal/forgejo/usecases/common"
)

const defaultReposWorkers = 4

type repoJob struct {
	index    int
	owner    string
	repoName string
	repo     *forgejo.Repository
}

type repoResult struct {
	index  int
	record commonUsecases.RepoRecord
}

// Service は repo list operation の実装です。
type Service struct {
	client       *forgejo.Client
	host         string
	username     string
	token        string
	httpClient   commonUsecases.HTTPClient
	reposWorkers int
}

// Options は repo list operation の依存です。
type Options struct {
	Client       *forgejo.Client
	Host         string
	Username     string
	Token        string
	HTTPClient   commonUsecases.HTTPClient
	ReposWorkers int
}

// NewService は repo list operation サービスを作成します。
func NewService(options Options) *Service {
	return &Service{
		client:       options.Client,
		host:         options.Host,
		username:     strings.TrimSpace(options.Username),
		token:        options.Token,
		httpClient:   options.HTTPClient,
		reposWorkers: resolveReposWorkers(options.ReposWorkers),
	}
}

// Execute は repo list を実行します。
func (s *Service) Execute() ([]commonUsecases.RepoRecord, error) {
	const reposPageSize = 100

	repos := make([]*forgejo.Repository, 0)
	page := 1
	for {
		currentPageRepos, response, err := s.client.ListMyRepos(forgejo.ListReposOptions{
			ListOptions: forgejo.ListOptions{
				Page:     page,
				PageSize: reposPageSize,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("リポジトリ一覧の取得に失敗しました: %w", err)
		}

		repos = append(repos, currentPageRepos...)

		if response == nil {
			if len(currentPageRepos) < reposPageSize {
				break
			}
			page++
			continue
		}
		if response.NextPage <= page {
			break
		}
		page = response.NextPage
	}

	repoJobs := make([]repoJob, 0, len(repos))
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
		repoJobs = append(repoJobs, repoJob{
			index:    len(repoJobs),
			owner:    owner,
			repoName: repoName,
			repo:     repo,
		})
	}

	records := make([]commonUsecases.RepoRecord, len(repoJobs))
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
			record, fetchErr := s.fetchRepoRecord(job.owner, job.repoName, job.repo)
			if fetchErr != nil {
				reportErr(fetchErr)
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
	for index := 0; index < workerCount; index++ {
		go worker()
	}

	wg.Wait()
	close(results)

	for result := range results {
		records[result.index] = result.record
	}

	select {
	case resultErr := <-errCh:
		return nil, resultErr
	default:
	}

	return records, nil
}

func resolveReposWorkers(raw int) int {
	if raw <= 0 {
		return defaultReposWorkers
	}
	return raw
}

func (s *Service) fetchTopics(owner, repo string) ([]string, error) {
	topics, _, err := s.client.ListRepoTopics(owner, repo, forgejo.ListRepoTopicsOptions{})
	return topics, err
}

func (s *Service) fetchLanguages(owner, repo string) (map[string]float64, error) {
	requestURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/languages", s.host, url.PathEscape(owner), url.PathEscape(repo))
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
	if resp.StatusCode != http.StatusOK {
		return nil, &commonUsecases.RequestError{
			Status: resp.StatusCode,
			Body:   string(body),
		}
	}

	var languages map[string]float64
	if err := json.Unmarshal(body, &languages); err != nil {
		return nil, err
	}
	return languages, nil
}

func (s *Service) fetchRepoRecord(owner, repoName string, repo *forgejo.Repository) (commonUsecases.RepoRecord, error) {
	topics, err := s.fetchTopics(owner, repoName)
	if err != nil {
		return commonUsecases.RepoRecord{}, fmt.Errorf("topicsの取得に失敗しました (%s/%s): %w", owner, repoName, err)
	}
	languages, err := s.fetchLanguages(owner, repoName)
	if err != nil {
		return commonUsecases.RepoRecord{}, fmt.Errorf("languagesの取得に失敗しました (%s/%s): %w", owner, repoName, err)
	}
	closedPullsCount, err := s.fetchClosedPullsCount(owner, repoName)
	if err != nil {
		return commonUsecases.RepoRecord{}, fmt.Errorf("pullsの取得に失敗しました (%s/%s): %w", owner, repoName, err)
	}
	closedIssuesCount, err := s.fetchClosedIssuesCount(owner, repoName)
	if err != nil {
		return commonUsecases.RepoRecord{}, fmt.Errorf("issuesの取得に失敗しました (%s/%s): %w", owner, repoName, err)
	}

	return commonUsecases.RepoRecord{
		Name:              repoName,
		Description:       repo.Description,
		IsPrivate:         repo.Private,
		HTTPURL:           repo.HTMLURL,
		OpenIssuesCount:   repo.OpenIssues,
		ClosedIssuesCount: closedIssuesCount,
		OpenPullsCount:    repo.OpenPulls,
		ClosedPullsCount:  closedPullsCount,
		ForksCount:        repo.Forks,
		StargazersCount:   repo.Stars,
		SubscribersCount:  repo.Watchers,
		Language:          commonUsecases.PrimaryLanguage(languages),
		Languages:         languages,
		Size:              repo.Size,
		RepoCreatedAt:     commonUsecases.FormatDate(repo.Created),
		RepoUpdatedAt:     commonUsecases.FormatDate(repo.Updated),
		IsArchived:        repo.Archived,
		Tags:              strings.Join(topics, ","),
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
		totalClosedPulls, convErr := strconv.Atoi(totalCountFromHeader)
		if convErr == nil {
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

func (s *Service) fetchClosedIssuesCount(owner, repo string) (int, error) {
	return s.fetchIssuesCount(owner, repo, forgejo.StateClosed)
}

func (s *Service) fetchIssuesCount(owner, repo string, state forgejo.StateType) (int, error) {
	firstIssues, response, err := s.client.ListRepoIssues(owner, repo, forgejo.ListIssueOption{
		ListOptions: forgejo.ListOptions{
			Page:     1,
			PageSize: 1,
		},
		State: state,
		Type:  forgejo.IssueTypeIssue,
	})
	if err != nil {
		return 0, err
	}

	if response == nil {
		return len(firstIssues), nil
	}

	totalCountFromHeader := strings.TrimSpace(response.Header.Get("X-Total-Count"))
	if totalCountFromHeader != "" {
		totalClosedIssues, convErr := strconv.Atoi(totalCountFromHeader)
		if convErr == nil {
			return totalClosedIssues, nil
		}
	}
	if len(firstIssues) == 0 {
		return 0, nil
	}

	if response.LastPage > 1 {
		return response.LastPage, nil
	}

	return len(firstIssues), nil
}
