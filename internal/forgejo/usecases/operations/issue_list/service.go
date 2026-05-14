package issue_list

import (
	"fmt"
	"strings"

	forgejo "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"

	commonUsecases "github.com/landmaster135/devbox/internal/forgejo/usecases/common"
)

const issuesPageSize = 100

// Service は issue list operation の実装です。
type Service struct {
	client   *forgejo.Client
	username string
}

// Options は issue list operation の依存です。
type Options struct {
	Client   *forgejo.Client
	Username string
}

// NewService は issue list operation サービスを作成します。
func NewService(options Options) *Service {
	return &Service{
		client:   options.Client,
		username: strings.TrimSpace(options.Username),
	}
}

// Execute は issue list を実行します。
func (s *Service) Execute() ([]commonUsecases.IssueRecord, error) {
	repos, _, err := s.client.ListUserRepos(s.username, forgejo.ListReposOptions{})
	if err != nil {
		return nil, fmt.Errorf("ユーザーリポジトリ一覧の取得に失敗しました: %w", err)
	}

	records := make([]commonUsecases.IssueRecord, 0)
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

		issues, fetchErr := s.fetchRepoIssues(owner, repoName)
		if fetchErr != nil {
			return nil, fmt.Errorf("issuesの取得に失敗しました (%s/%s): %w", owner, repoName, fetchErr)
		}

		for _, issue := range issues {
			if issue == nil {
				continue
			}
			records = append(records, toIssueRecord(owner, repoName, issue))
		}
	}

	return records, nil
}

func (s *Service) fetchRepoIssues(owner, repo string) ([]*forgejo.Issue, error) {
	issues := make([]*forgejo.Issue, 0)
	page := 1
	for {
		currentPageIssues, response, err := s.client.ListRepoIssues(owner, repo, forgejo.ListIssueOption{
			ListOptions: forgejo.ListOptions{
				Page:     page,
				PageSize: issuesPageSize,
			},
			State: forgejo.StateAll,
			Type:  forgejo.IssueTypeIssue,
		})
		if err != nil {
			return nil, err
		}

		issues = append(issues, currentPageIssues...)

		if response == nil {
			if len(currentPageIssues) < issuesPageSize {
				return issues, nil
			}
			page++
			continue
		}
		if response.NextPage <= page {
			return issues, nil
		}
		page = response.NextPage
	}
}

func toIssueRecord(owner, repo string, issue *forgejo.Issue) commonUsecases.IssueRecord {
	return commonUsecases.IssueRecord{
		RepoFullName: owner + "/" + repo,
		Number:       issue.Index,
		Title:        issue.Title,
		State:        string(issue.State),
		HTMLURL:      issue.HTMLURL,
		Author:       issueAuthor(issue),
		Assignees:    assigneeNames(issue.Assignees),
		Labels:       labelNames(issue.Labels),
		Comments:     issue.Comments,
		IsLocked:     issue.IsLocked,
		CreatedAt:    commonUsecases.FormatDate(issue.Created),
		UpdatedAt:    commonUsecases.FormatDate(issue.Updated),
		ClosedAt:     commonUsecases.FormatDatePtr(issue.Closed),
	}
}

func issueAuthor(issue *forgejo.Issue) string {
	if issue == nil || issue.Poster == nil {
		return ""
	}
	return strings.TrimSpace(issue.Poster.UserName)
}

func assigneeNames(assignees []*forgejo.User) []string {
	if len(assignees) == 0 {
		return []string{}
	}
	names := make([]string, 0, len(assignees))
	for _, assignee := range assignees {
		if assignee == nil {
			continue
		}
		name := strings.TrimSpace(assignee.UserName)
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	return names
}

func labelNames(labels []*forgejo.Label) []string {
	if len(labels) == 0 {
		return []string{}
	}
	names := make([]string, 0, len(labels))
	for _, label := range labels {
		if label == nil {
			continue
		}
		name := strings.TrimSpace(label.Name)
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	return names
}
