package usecases

import (
	"encoding/json"
	"fmt"
	"strings"
)

// #==============================================================#
// ##          Option Structs                                   ##
// #==============================================================#

// CreateIssueOptions はイシュー作成時のオプションを定義します
type CreateIssueOptions struct {
	Title     string   `json:"title"`
	Body      string   `json:"body,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
}

// ListIssuesOptions はイシュー一覧取得時のオプションを定義します
type ListIssuesOptions struct {
	State     string `json:"state,omitempty"`
	Sort      string `json:"sort,omitempty"`
	Direction string `json:"direction,omitempty"`
	PerPage   int    `json:"per_page,omitempty"`
	Page      int    `json:"page,omitempty"`
}

// UpdateIssueOptions はイシュー更新時のオプションを定義します
type UpdateIssueOptions struct {
	Title     string   `json:"title,omitempty"`
	Body      string   `json:"body,omitempty"`
	State     string   `json:"state,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
}

// #==============================================================#
// ##          GitHubIssueService                               ##
// #==============================================================#

// GitHubIssueService はGitHub Issue関連機能のビジネスロジックを提供します
type GitHubIssueService struct {
	clientService *GitHubClientService
}

// NewGitHubIssueService は新しいGitHubIssueServiceを作成します
func NewGitHubIssueService(token string) *GitHubIssueService {
	return &GitHubIssueService{
		clientService: NewGitHubClientService(token),
	}
}

// NewGitHubIssueServiceWithDependencies はテスト用に依存性を注入できるGitHubIssueServiceを作成します
func NewGitHubIssueServiceWithDependencies(clientService *GitHubClientService) *GitHubIssueService {
	return &GitHubIssueService{
		clientService: clientService,
	}
}

// #==============================================================#
// ##          Issue Operations                                 ##
// #==============================================================#

// createIssue は新しいイシューを作成します
func (s *GitHubIssueService) createIssue(owner, repo string, options CreateIssueOptions) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues", apiBaseURL, owner, repo)
	jsonBody, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	data, err := s.clientService.DoRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// listIssues はリポジトリのイシュー一覧を取得します
func (s *GitHubIssueService) listIssues(owner, repo string, options ListIssuesOptions) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues", apiBaseURL, owner, repo)

	// クエリパラメータを追加
	queryParams := []string{}
	if options.State != "" {
		queryParams = append(queryParams, fmt.Sprintf("state=%s", options.State))
	}
	if options.Sort != "" {
		queryParams = append(queryParams, fmt.Sprintf("sort=%s", options.Sort))
	}
	if options.Direction != "" {
		queryParams = append(queryParams, fmt.Sprintf("direction=%s", options.Direction))
	}
	if options.PerPage > 0 {
		queryParams = append(queryParams, fmt.Sprintf("per_page=%d", options.PerPage))
	}
	if options.Page > 0 {
		queryParams = append(queryParams, fmt.Sprintf("page=%d", options.Page))
	}

	if len(queryParams) > 0 {
		url += "?" + strings.Join(queryParams, "&")
	}

	data, err := s.clientService.DoRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// getIssue は特定のイシューを取得します
func (s *GitHubIssueService) getIssue(owner, repo string, issueNumber int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d", apiBaseURL, owner, repo, issueNumber)

	data, err := s.clientService.DoRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// updateIssue は既存のイシューを更新します
func (s *GitHubIssueService) updateIssue(owner, repo string, issueNumber int, options UpdateIssueOptions) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d", apiBaseURL, owner, repo, issueNumber)

	jsonBody, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	data, err := s.clientService.DoRequest("PATCH", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// addIssueComment はイシューにコメントを追加します
func (s *GitHubIssueService) addIssueComment(owner, repo string, issueNumber int, body string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments", apiBaseURL, owner, repo, issueNumber)

	jsonBody, err := json.Marshal(map[string]string{"body": body})
	if err != nil {
		return nil, err
	}

	data, err := s.clientService.DoRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// #==============================================================#
// ##          Handler Methods                                  ##
// #==============================================================#

// HandleToCreateIssue は新しいイシューを作成して、結果をJSON形式で返します
func (s *GitHubIssueService) HandleToCreateIssue(owner, repo, title, body string, labels, assignees []interface{}) (string, error) {
	// []interface{}を[]stringに変換
	var labelStrings []string
	for _, label := range labels {
		if str, ok := label.(string); ok {
			labelStrings = append(labelStrings, str)
		}
	}

	var assigneeStrings []string
	for _, assignee := range assignees {
		if str, ok := assignee.(string); ok {
			assigneeStrings = append(assigneeStrings, str)
		}
	}

	options := CreateIssueOptions{
		Title:     title,
		Body:      body,
		Labels:    labelStrings,
		Assignees: assigneeStrings,
	}

	result, err := s.createIssue(owner, repo, options)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// HandleToListIssues はリポジトリのイシュー一覧を取得して、結果をJSON形式で返します
// issueNumberが指定された場合は特定のイシューを取得します
func (s *GitHubIssueService) HandleToListIssues(owner, repo, state, sort, direction string, perPage, page, issueNumber int) (string, error) {
	// issueNumberが指定された場合は特定のイシューを取得
	if issueNumber > 0 {
		result, err := s.getIssue(owner, repo, issueNumber)
		if err != nil {
			return "", err
		}
		return s.clientService.ReturnJSONResult(result)
	}

	// issueNumberが指定されていない場合はイシュー一覧を取得
	options := ListIssuesOptions{
		State:     state,
		Sort:      sort,
		Direction: direction,
		PerPage:   perPage,
		Page:      page,
	}

	result, err := s.listIssues(owner, repo, options)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// HandleToUpdateIssue は既存のイシューを更新して、結果をJSON形式で返します
func (s *GitHubIssueService) HandleToUpdateIssue(owner, repo string, issueNumber int, title, body, state string, labels, assignees []interface{}) (string, error) {
	// []interface{}を[]stringに変換
	var labelStrings []string
	if labels != nil {
		for _, label := range labels {
			if str, ok := label.(string); ok {
				labelStrings = append(labelStrings, str)
			}
		}
	}

	var assigneeStrings []string
	if assignees != nil {
		for _, assignee := range assignees {
			if str, ok := assignee.(string); ok {
				assigneeStrings = append(assigneeStrings, str)
			}
		}
	}

	options := UpdateIssueOptions{
		Title:     title,
		Body:      body,
		State:     state,
		Labels:    labelStrings,
		Assignees: assigneeStrings,
	}

	result, err := s.updateIssue(owner, repo, issueNumber, options)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// HandleToAddIssueComment はイシューにコメントを追加して、結果をJSON形式で返します
func (s *GitHubIssueService) HandleToAddIssueComment(owner, repo string, issueNumber int, body string) (string, error) {
	result, err := s.addIssueComment(owner, repo, issueNumber, body)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}
