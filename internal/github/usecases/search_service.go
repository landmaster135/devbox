package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
)

// #==============================================================#
// ##          GitHubSearchService                               ##
// #==============================================================#

func validateParams(token string) error {
	if token == "" {
		return errors.New("token must not be empty")
	}
	return nil
}

// GitHubSearchService はGitHub検索機能のビジネスロジックを提供します
type GitHubSearchService struct {
	clientService *GitHubClientService
}

// NewGitHubSearchService は新しいGitHubSearchServiceを作成します
func NewGitHubSearchService(token string) (*GitHubSearchService, error) {
	err := validateParams(token)
	if err != nil {
		return nil, err
	}
	return &GitHubSearchService{
		clientService: NewGitHubClientService(token),
	}, nil
}

// NewGitHubSearchServiceWithDependencies はテスト用に依存性を注入できるGitHubSearchServiceを作成します
func NewGitHubSearchServiceWithDependencies(clientService *GitHubClientService) *GitHubSearchService {
	return &GitHubSearchService{
		clientService: clientService,
	}
}

// SearchCode はGitHub全体でコードを検索します
func (s *GitHubSearchService) SearchCode(query string, options map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/search/code?q=%s", apiBaseURL, query)

	// クエリパラメータを追加
	for k, v := range options {
		url += fmt.Sprintf("&%s=%v", k, v)
	}

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

// HandleToSearchCode はGitHub全体でコードを検索して、結果をJSON形式で返します
func (s *GitHubSearchService) HandleToSearchCode(query string, page, perPage int) (string, error) {
	options := make(map[string]interface{})

	// 数値オプションパラメータを追加
	if page != 1 {
		options["page"] = page
	}

	if perPage != 30 {
		options["per_page"] = perPage
	}

	result, err := s.SearchCode(query, options)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}
