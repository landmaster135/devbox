package usecases

import (
	"encoding/json"
	"fmt"
	"strings"
)

// #==============================================================#
// ##          GitHubPullRequestService                         ##
// #==============================================================#

// GitHubPullRequestService はGitHub Pull Request関連機能のビジネスロジックを提供します
type GitHubPullRequestService struct {
	clientService *GitHubClientService
}

// NewGitHubPullRequestService は新しいGitHubPullRequestServiceを作成します
func NewGitHubPullRequestService(token string) *GitHubPullRequestService {
	return &GitHubPullRequestService{
		clientService: NewGitHubClientService(token),
	}
}

// NewGitHubPullRequestServiceWithDependencies はテスト用に依存性を注入できるGitHubPullRequestServiceを作成します
func NewGitHubPullRequestServiceWithDependencies(clientService *GitHubClientService) *GitHubPullRequestService {
	return &GitHubPullRequestService{
		clientService: clientService,
	}
}

// #==============================================================#
// ##          Pull Request Operations                          ##
// #==============================================================#

// createPullRequest は新しいプルリクエストを作成します
func (s *GitHubPullRequestService) createPullRequest(owner, repo string, options map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls", apiBaseURL, owner, repo)

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

// createPullRequestReview はプルリクエストにレビューを作成します
func (s *GitHubPullRequestService) createPullRequestReview(owner, repo string, pullNumber int, options map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/reviews", apiBaseURL, owner, repo, pullNumber)

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

// mergePullRequest はプルリクエストをマージします
func (s *GitHubPullRequestService) mergePullRequest(owner, repo string, pullNumber int, options map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/merge", apiBaseURL, owner, repo, pullNumber)

	jsonBody, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	data, err := s.clientService.DoRequest("PUT", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// getPullRequestFiles はプルリクエストで変更されたファイル一覧を取得します
func (s *GitHubPullRequestService) getPullRequestFiles(owner, repo string, pullNumber int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/files", apiBaseURL, owner, repo, pullNumber)

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

// getPullRequestStatus はプルリクエストのステータスを取得します
func (s *GitHubPullRequestService) getPullRequestStatus(owner, repo string, pullNumber int) (map[string]interface{}, error) {
	// プルリクエストの詳細を取得
	prURL := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", apiBaseURL, owner, repo, pullNumber)
	prData, err := s.clientService.DoRequest("GET", prURL, nil)
	if err != nil {
		return nil, err
	}

	var pr map[string]interface{}
	if err := json.Unmarshal(prData, &pr); err != nil {
		return nil, err
	}

	// ステータスチェックを取得
	headSHA, ok := pr["head"].(map[string]interface{})["sha"].(string)
	if !ok {
		return nil, fmt.Errorf("could not get head SHA from pull request")
	}

	statusURL := fmt.Sprintf("%s/repos/%s/%s/commits/%s/status", apiBaseURL, owner, repo, headSHA)
	statusData, err := s.clientService.DoRequest("GET", statusURL, nil)
	if err != nil {
		return nil, err
	}

	var status map[string]interface{}
	if err := json.Unmarshal(statusData, &status); err != nil {
		return nil, err
	}

	// 結果を組み合わせる
	result := map[string]interface{}{
		"pull_request": pr,
		"status":       status,
	}

	return result, nil
}

// updatePullRequestBranch はプルリクエストのブランチを更新します
func (s *GitHubPullRequestService) updatePullRequestBranch(owner, repo string, pullNumber int, expectedHeadSHA string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/update-branch", apiBaseURL, owner, repo, pullNumber)

	options := map[string]interface{}{}
	if expectedHeadSHA != "" {
		options["expected_head_sha"] = expectedHeadSHA
	}

	jsonBody, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	data, err := s.clientService.DoRequest("PUT", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// getPullRequestComments はプルリクエストのコメントを取得します
func (s *GitHubPullRequestService) getPullRequestComments(owner, repo string, pullNumber int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/comments", apiBaseURL, owner, repo, pullNumber)

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

// getPullRequestReviews はプルリクエストのレビューを取得します
func (s *GitHubPullRequestService) getPullRequestReviews(owner, repo string, pullNumber int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/reviews", apiBaseURL, owner, repo, pullNumber)

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

// listPullRequests はリポジトリのプルリクエスト一覧を取得します
func (s *GitHubPullRequestService) listPullRequests(owner, repo string, options map[string]interface{}) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls", apiBaseURL, owner, repo)

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

// #==============================================================#
// ##          Handler Methods                                  ##
// #==============================================================#

// HandleToCreatePullRequest は新しいプルリクエストを作成して、結果をJSON形式で返します
func (s *GitHubPullRequestService) HandleToCreatePullRequest(owner, repo, title, head, base, body string, draft bool) (string, error) {
	options := make(map[string]interface{})
	options["title"] = title
	options["head"] = head
	options["base"] = base

	if body != "" {
		options["body"] = body
	}

	options["draft"] = draft

	result, err := s.createPullRequest(owner, repo, options)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// HandleToCreatePullRequestReview はプルリクエストにレビューを作成して、結果をJSON形式で返します
func (s *GitHubPullRequestService) HandleToCreatePullRequestReview(owner, repo string, pullNumber int, event, body string) (string, error) {
	options := make(map[string]interface{})

	if event != "" {
		options["event"] = event
	}
	if body != "" {
		options["body"] = body
	}

	result, err := s.createPullRequestReview(owner, repo, pullNumber, options)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// HandleToMergePullRequest はプルリクエストをマージして、結果をJSON形式で返します
func (s *GitHubPullRequestService) HandleToMergePullRequest(owner, repo string, pullNumber int, commitTitle, commitMessage, mergeMethod string) (string, error) {
	options := make(map[string]interface{})

	if commitTitle != "" {
		options["commit_title"] = commitTitle
	}
	if commitMessage != "" {
		options["commit_message"] = commitMessage
	}
	if mergeMethod != "" {
		options["merge_method"] = mergeMethod
	}

	result, err := s.mergePullRequest(owner, repo, pullNumber, options)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// HandleToGetPullRequestFiles はプルリクエストで変更されたファイル一覧を取得して、結果をJSON形式で返します
func (s *GitHubPullRequestService) HandleToGetPullRequestFiles(owner, repo string, pullNumber int) (string, error) {
	result, err := s.getPullRequestFiles(owner, repo, pullNumber)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// HandleToGetPullRequestStatus はプルリクエストのステータスを取得して、結果をJSON形式で返します
func (s *GitHubPullRequestService) HandleToGetPullRequestStatus(owner, repo string, pullNumber int) (string, error) {
	result, err := s.getPullRequestStatus(owner, repo, pullNumber)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// HandleToUpdatePullRequestBranch はプルリクエストのブランチを更新して、結果をJSON形式で返します
func (s *GitHubPullRequestService) HandleToUpdatePullRequestBranch(owner, repo string, pullNumber int, expectedHeadSHA string) (string, error) {
	result, err := s.updatePullRequestBranch(owner, repo, pullNumber, expectedHeadSHA)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// HandleToGetPullRequestComments はプルリクエストのコメントを取得して、結果をJSON形式で返します
func (s *GitHubPullRequestService) HandleToGetPullRequestComments(owner, repo string, pullNumber int) (string, error) {
	result, err := s.getPullRequestComments(owner, repo, pullNumber)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// HandleToGetPullRequestReviews はプルリクエストのレビューを取得して、結果をJSON形式で返します
func (s *GitHubPullRequestService) HandleToGetPullRequestReviews(owner, repo string, pullNumber int) (string, error) {
	result, err := s.getPullRequestReviews(owner, repo, pullNumber)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// HandleToListPullRequests はリポジトリのプルリクエスト一覧を取得して、結果をJSON形式で返します
func (s *GitHubPullRequestService) HandleToListPullRequests(owner, repo, state, sort, direction, head, base string, perPage, page int) (string, error) {
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
	if head != "" {
		options["head"] = head
	}
	if base != "" {
		options["base"] = base
	}
	if perPage > 0 {
		options["per_page"] = perPage
	}
	if page > 0 {
		options["page"] = page
	}

	result, err := s.listPullRequests(owner, repo, options)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}
