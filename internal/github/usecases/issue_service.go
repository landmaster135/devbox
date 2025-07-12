package usecases

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
// ##          Helper Functions                                 ##
// #==============================================================#

// AddToOptions はオプションマップにパラメータを追加するヘルパー関数です
func AddToOptions(options map[string]interface{}, args map[string]interface{}, key string) {
	if val, ok := args[key]; ok {
		options[key] = val
	}
}

// #==============================================================#
// ##          Issue Operations                                 ##
// #==============================================================#

// createIssue は新しいイシューを作成します
func (s *GitHubIssueService) createIssue(owner, repo string, options map[string]interface{}) (map[string]interface{}, error) {
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
func (s *GitHubIssueService) listIssues(owner, repo string, options map[string]interface{}) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues", apiBaseURL, owner, repo)

	// クエリパラメータを追加
	queryParams := []string{}
	for k, v := range options {
		queryParams = append(queryParams, fmt.Sprintf("%s=%v", k, v))
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

// updateIssue は既存のイシューを更新します
func (s *GitHubIssueService) updateIssue(owner, repo string, issueNumber int, options map[string]interface{}) (map[string]interface{}, error) {
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
	options := make(map[string]interface{})
	options["title"] = title

	if body != "" {
		options["body"] = body
	}

	if labels != nil {
		options["labels"] = labels
	}

	if assignees != nil {
		options["assignees"] = assignees
	}

	result, err := s.createIssue(owner, repo, options)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// HandleToListIssues はリポジトリのイシュー一覧を取得して、結果をJSON形式で返します
func (s *GitHubIssueService) HandleToListIssues(owner, repo, state, sort, direction string, perPage, page int) (string, error) {
	options := make(map[string]interface{})

	if state != "" {
		options["state"] = state
	}

	if sort != "" {
		options["sort"] = sort
	}

	if direction != "" {
		options["direction"] = direction
	}

	options["per_page"] = perPage
	options["page"] = page

	result, err := s.listIssues(owner, repo, options)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// HandleToUpdateIssue は既存のイシューを更新して、結果をJSON形式で返します
func (s *GitHubIssueService) HandleToUpdateIssue(owner, repo string, issueNumber int, title, body, state string, labels, assignees []interface{}) (string, error) {
	options := make(map[string]interface{})

	if title != "" {
		options["title"] = title
	}

	if body != "" {
		options["body"] = body
	}

	if state != "" {
		options["state"] = state
	}

	if labels != nil {
		options["labels"] = labels
	}

	if assignees != nil {
		options["assignees"] = assignees
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
