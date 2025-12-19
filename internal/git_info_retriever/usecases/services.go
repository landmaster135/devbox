package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	security "github.com/landmaster135/devbox/internal/git_info_retriever/security"
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
	ArchiveRepositories(service, token, outputFilePath, archiveDir, srcFile string) (string, error)
}

// PathValidator はサービス内で利用するバリデーターのインターフェース
type PathValidator interface {
	ValidateArchiveDirectory(path string) (string, error)
	ValidatePathComponent(component string) (string, error)
	ValidateURLString(rawURL string) (string, error)
}

// Service はGitInfoServiceの実装
type Service struct {
	githubService GitHubService
	fileWriter    FileWriter
	pathValidator PathValidator
}

// NewServiceWithDependencies は依存関係を注入してServiceインスタンスを作成する
func NewServiceWithDependencies(githubService GitHubService, fileWriter FileWriter) GitInfoService {
	return &Service{
		githubService: githubService,
		fileWriter:    fileWriter,
		pathValidator: security.NewDefaultPathValidator(),
	}
}

// NewService はデフォルト実装でServiceインスタンスを作成する
func NewService() GitInfoService {
	return NewServiceWithDependencies(
		NewGitHubService(),
		NewFileWriter(),
	)
}

func (s *Service) getPathValidator() PathValidator {
	if s.pathValidator == nil {
		s.pathValidator = security.NewDefaultPathValidator()
	}
	return s.pathValidator
}

// getRepositoryInfo は指定されたサービスからリポジトリ情報を取得する
func (s *Service) getRepositoryInfo(service, token string) (string, error) {
	if service != "github" {
		return "", fmt.Errorf("サポートされていないサービスです: %s", service)
	}

	ctx := context.Background()
	client := s.githubService.CreateGitHubClient(ctx, token)

	// プライベートリポジトリも含めて取得するため、空文字列を渡す
	username := ""
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

// loadRepoInfoFromFile はファイルからリポジトリ情報を読み込む
func (s *Service) loadRepoInfoFromFile(filePath string) ([]RepoInfo, error) {
	// ファイルの存在確認
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("指定されたファイルが存在しません: %s", filePath)
	}

	// ファイル読み込み
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗しました: %v", err)
	}

	// JSON解析
	var repos []RepoInfo
	err = json.Unmarshal(data, &repos)
	if err != nil {
		return nil, fmt.Errorf("JSONの解析に失敗しました: %v", err)
	}

	return repos, nil
}

// extractRepoName はHTTP URLからリポジトリ名を抽出する
func extractRepoName(httpUrl string) string {
	// URLの最後の部分を取得
	parts := strings.Split(httpUrl, "/")
	if len(parts) == 0 {
		return ""
	}

	repoName := parts[len(parts)-1]
	// .gitサフィックスを除去
	repoName = strings.TrimSuffix(repoName, ".git")

	return repoName
}

// generateBashFunctions はBash関数を生成する
func (s *Service) generateBashFunctions(repos []RepoInfo, archiveDir string) (string, error) {
	validator := s.getPathValidator()

	sanitizedArchiveDir, err := validator.ValidateArchiveDirectory(archiveDir)
	if err != nil {
		return "", fmt.Errorf("アーカイブディレクトリの検証に失敗しました: %w", err)
	}

	type sanitizedRepo struct {
		name string
		url  string
	}

	sanitizedRepos := make([]sanitizedRepo, 0, len(repos))
	for _, repo := range repos {
		validatedURL, err := validator.ValidateURLString(repo.HttpUrl)
		if err != nil {
			return "", fmt.Errorf("リポジトリURLの検証に失敗しました (%s): %w", repo.HttpUrl, err)
		}

		repoName := extractRepoName(validatedURL)
		if repoName == "" {
			continue
		}

		sanitizedName, err := validator.ValidatePathComponent(repoName)
		if err != nil {
			return "", fmt.Errorf("リポジトリ名の検証に失敗しました (%s): %w", repoName, err)
		}

		sanitizedRepos = append(sanitizedRepos, sanitizedRepo{
			name: sanitizedName,
			url:  validatedURL,
		})
	}

	var result strings.Builder

	// archive_repos関数
	result.WriteString("function archive_repos() {\n")
	result.WriteString(fmt.Sprintf("\tmkdir -p %s\n", singleQuote(sanitizedArchiveDir)))
	for _, repo := range sanitizedRepos {
		archiveZip := fmt.Sprintf("%s/%s.zip", sanitizedArchiveDir, repo.name)
		result.WriteString(fmt.Sprintf("\tgit clone %s\n", singleQuote(repo.url)))
		result.WriteString(fmt.Sprintf("\tzip -rq %s %s\n", singleQuote(archiveZip), singleQuote("./"+repo.name)))
	}
	result.WriteString("}\n\n")

	// display_zipinfo関数
	result.WriteString("function display_zipinfo() {\n")
	for _, repo := range sanitizedRepos {
		archiveZip := fmt.Sprintf("%s/%s.zip", sanitizedArchiveDir, repo.name)
		result.WriteString(fmt.Sprintf("\tzipinfo %s\n", singleQuote(archiveZip)))
	}
	result.WriteString("}\n\n")

	// unzip_repos関数
	result.WriteString("function unzip_repos() {\n")
	for _, repo := range sanitizedRepos {
		archiveZip := fmt.Sprintf("%s/%s.zip", sanitizedArchiveDir, repo.name)
		result.WriteString(fmt.Sprintf("\tunzip %s -d %s\n", singleQuote(archiveZip), singleQuote("./unarchived")))
	}
	result.WriteString("}\n")

	return result.String(), nil
}

// ArchiveRepositories はリポジトリアーカイブ用のBash関数を生成する
func (s *Service) ArchiveRepositories(service, token, outputFilePath, archiveDir, srcFile string) (string, error) {
	var repos []RepoInfo
	var err error

	// リポジトリ情報の取得
	if srcFile != "" {
		// ファイルから読み込み
		repos, err = s.loadRepoInfoFromFile(srcFile)
		if err != nil {
			return "", err
		}
	} else {
		// GitHubから取得
		if service != "github" {
			return "", fmt.Errorf("サポートされていないサービスです: %s", service)
		}

		ctx := context.Background()
		client := s.githubService.CreateGitHubClient(ctx, token)

		// プライベートリポジトリも含めて取得するため、空文字列を渡す
		username := ""
		repos, err = s.githubService.GetRepoInfo(ctx, client, true, username)
		if err != nil {
			return "", fmt.Errorf("リポジトリ情報の取得に失敗しました: %v", err)
		}
	}

	// アーカイブディレクトリのパス正規化
	archiveDir = filepath.Clean(archiveDir)

	// Bash関数を生成
	bashFunctions, err := s.generateBashFunctions(repos, archiveDir)
	if err != nil {
		return "", err
	}

	// ファイル出力が指定されている場合
	if outputFilePath != "" {
		err = s.fileWriter.WriteToFile(outputFilePath, bashFunctions)
		if err != nil {
			return "", fmt.Errorf("ファイル保存に失敗しました: %v", err)
		}
		fmt.Printf("Bash関数をファイルに保存しました: %s\n", outputFilePath)
	}

	return bashFunctions, nil
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

func singleQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\\''")
	return "'" + escaped + "'"
}
