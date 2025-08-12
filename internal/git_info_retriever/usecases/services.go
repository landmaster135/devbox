package usecases

import (
	"context"
	"encoding/json"
	"fmt"
)

// #==============================================================#
// ##          RepoInfo                                          ##
// #==============================================================#
// RepoInfo はリポジトリ情報を保持する構造体
type RepoInfo struct {
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	IsPrivate        bool           `json:"is_private"`
	HttpUrl          string         `json:"http_url"`
	Language         string         `json:"language"`
	Languages        map[string]int `json:"languages"`
	RepoCreatedAt    string         `json:"repo_created_at"`
	RepoUpdatedAt    string         `json:"repo_updated_at"`
	StargazersCount  int32          `json:"stargazers_count"`
	ForksCount       int32          `json:"forks_count"`
	IssuesCount      int32          `json:"issues_count"`
	PullsCount       int32          `json:"pulls_count"`
	Size             int32          `json:"size"`
	SubscribersCount int32          `json:"subscribers_count"`
	IsArchived       bool           `json:"is_archived"`
}

// #==============================================================#
// ##          GitInfoService                                    ##
// #==============================================================#
// GitInfoService はGit情報取得サービスのインターフェース
type GitInfoService interface {
	RetrieveRepositoryInfo(service, token, filePath string) (string, error)
}

// Service はGitInfoServiceの実装
type Service struct {
	githubService GitHubService
	fileWriter    FileWriter
}

// NewServiceWithDependencies は依存関係を注入してServiceインスタンスを作成する
func NewServiceWithDependencies(githubService GitHubService, fileWriter FileWriter) GitInfoService {
	return &Service{
		githubService: githubService,
		fileWriter:    fileWriter,
	}
}

// NewService はデフォルト実装でServiceインスタンスを作成する
func NewService() GitInfoService {
	return NewServiceWithDependencies(
		NewGitHubService(),
		NewFileWriter(),
	)
}

// getRepositoryInfo は指定されたサービスからリポジトリ情報を取得する
func (s *Service) getRepositoryInfo(service, token string) (string, error) {
	if service != "github" {
		return "", fmt.Errorf("サポートされていないサービスです: %s", service)
	}

	ctx := context.Background()
	client := s.githubService.CreateGitHubClient(ctx, token)

	// 認証されたユーザー情報を取得
	user, _, err := client.GetUser(ctx, "")
	if err != nil {
		return "", fmt.Errorf("ユーザー情報の取得に失敗しました: %v", err)
	}

	username := user.GetLogin()
	repos, err := s.githubService.GetRepoInfo(ctx, client, true, username)
	if err != nil {
		return "", fmt.Errorf("リポジトリ情報の取得に失敗しました: %v", err)
	}

	// JSON形式で結果を返却
	result, err := json.MarshalIndent(repos, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON変換に失敗しました: %v", err)
	}

	return string(result), nil
}

// getAndSaveRepositoryInfo はリポジトリ情報を取得し、ファイルに保存する
func (s *Service) getAndSaveRepositoryInfo(service, token, filePath string) (string, error) {
	// リポジトリ情報を取得
	result, err := s.getRepositoryInfo(service, token)
	if err != nil {
		return "", err
	}

	// ファイルに保存
	err = s.fileWriter.WriteToFile(filePath, result)
	if err != nil {
		return "", fmt.Errorf("ファイル保存に失敗しました: %v", err)
	}

	return result, nil
}

func (s *Service) RetrieveRepositoryInfo(service, token, filePath string) (string, error) {
	var result string
	var err error

	// ファイル出力が指定されている場合
	if filePath != "" {
		result, err = s.getAndSaveRepositoryInfo(service, token, filePath)
		if err != nil {
			return "", err
		}
		fmt.Printf("結果をファイルに保存しました: %s\n", filePath)
		return result, nil
	}

	// 標準出力に表示
	result, err = s.getRepositoryInfo(service, token)
	if err != nil {
		return "", err
	}

	return result, nil
}
