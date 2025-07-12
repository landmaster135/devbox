package usecases

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// GitHubRepositoryService はGitHubリポジトリ関連の操作を提供します
type GitHubRepositoryService struct {
	clientService *GitHubClientService
}

// NewGitHubRepositoryService は新しいGitHubRepositoryServiceを作成します
func NewGitHubRepositoryService(token string) *GitHubRepositoryService {
	return &GitHubRepositoryService{
		clientService: NewGitHubClientService(token),
	}
}

// ListCommits はリポジトリのコミット一覧を取得します
func (s *GitHubRepositoryService) ListCommits(owner, repo string, page, perPage int, sha string) ([]map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 30
	}

	url := fmt.Sprintf("%s/repos/%s/%s/commits?page=%d&per_page=%d", apiBaseURL, owner, repo, page, perPage)
	if sha != "" {
		url += fmt.Sprintf("&sha=%s", sha)
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

// HandleToListCommits はリポジトリのコミット一覧を取得して、結果をJSON形式で返します
func (s *GitHubRepositoryService) HandleToListCommits(owner, repo string, page, perPage int, sha string) (string, error) {
	result, err := s.ListCommits(owner, repo, page, perPage, sha)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// SearchRepositories はGitHubリポジトリを検索します
func (s *GitHubRepositoryService) SearchRepositories(query string, page, perPage int) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 30
	}

	url := fmt.Sprintf("%s/search/repositories?q=%s&page=%d&per_page=%d", apiBaseURL, query, page, perPage)
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

// HandleToSearchRepositories はGitHubリポジトリを検索して、結果をJSON形式で返します
func (s *GitHubRepositoryService) HandleToSearchRepositories(query string, page, perPage int) (string, error) {
	result, err := s.SearchRepositories(query, page, perPage)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// GetUserRepositories はユーザーのリポジトリ一覧を取得します
func (s *GitHubRepositoryService) GetUserRepositories(username string, options map[string]interface{}) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/users/%s/repos", apiBaseURL, username)

	// クエリパラメータを追加（バグ修正: "?=" を "?" に変更）
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

// HandleToGetUserRepositories はユーザーのリポジトリ一覧を取得して、結果をJSON形式で返します
func (s *GitHubRepositoryService) HandleToGetUserRepositories(username, sort, direction, type_ string, perPage, page int) (string, error) {
	options := make(map[string]interface{})

	// 数値オプションパラメータを追加
	if perPage > 0 {
		options["per_page"] = perPage
	}
	if page > 0 {
		options["page"] = page
	}

	// 文字列オプションパラメータを追加
	if sort != "" {
		options["sort"] = sort
	}
	if direction != "" {
		options["direction"] = direction
	}
	if type_ != "" {
		options["type"] = type_
	}

	result, err := s.GetUserRepositories(username, options)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}

// GetFileContents はリポジトリからファイルの内容を取得します
func (s *GitHubRepositoryService) GetFileContents(owner, repo, path, branch string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", apiBaseURL, owner, repo, path)
	if branch != "" {
		url += fmt.Sprintf("?ref=%s", branch)
	}

	data, err := s.clientService.DoRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	// ファイルの内容をデコードする
	if content, ok := result["content"].(string); ok {
		decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(content, "\n", ""))
		if err != nil {
			return nil, err
		}
		result["decoded_content"] = string(decoded)
	}

	return result, nil
}

// HandleToGetFileContents はリポジトリからファイルの内容を取得して、結果をJSON形式で返します
func (s *GitHubRepositoryService) HandleToGetFileContents(owner, repo, path, branch string) (string, error) {
	result, err := s.GetFileContents(owner, repo, path, branch)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}
