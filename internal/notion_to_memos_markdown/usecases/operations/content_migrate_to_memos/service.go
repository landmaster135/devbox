package migratetomemos

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
	memos "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/memos"
	common "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/common"
)

type ClientFactory func(baseURL, apiToken string) memos.Client
type ProgressReporter interface {
	Report(message string)
}

type Service struct {
	fileSystem    filesystem.Repository
	clientFactory ClientFactory
	reporter      ProgressReporter
}

func NewService(fileSystem filesystem.Repository, clientFactory ClientFactory, reporter ProgressReporter) *Service {
	factory := clientFactory
	if factory == nil {
		factory = memos.NewClient
	}

	return &Service{
		fileSystem:    fileSystem,
		clientFactory: factory,
		reporter:      reporter,
	}
}

func (s *Service) Execute(pageType, baseURL, apiToken, srcBodyDir, srcResourceDir string) (string, error) {
	trimmedPageType := strings.TrimSpace(pageType)
	if trimmedPageType != common.SupportedPageType {
		return "", fmt.Errorf("未対応のpage-typeです: %s", trimmedPageType)
	}

	trimmedBaseURL := strings.TrimSpace(baseURL)
	if trimmedBaseURL == "" {
		return "", fmt.Errorf("base-url パラメータは必須です")
	}
	trimmedAPIToken := strings.TrimSpace(apiToken)
	if trimmedAPIToken == "" {
		return "", fmt.Errorf("api-token パラメータは必須です")
	}

	trimmedSrcBodyDir := strings.TrimSpace(srcBodyDir)
	if trimmedSrcBodyDir == "" {
		return "", fmt.Errorf("src-body-dir パラメータは必須です")
	}
	trimmedSrcResourceDir := strings.TrimSpace(srcResourceDir)
	if trimmedSrcResourceDir == "" {
		return "", fmt.Errorf("src-resource-dir パラメータは必須です")
	}

	bodyPaths, err := s.fileSystem.ListFilesRecursive(trimmedSrcBodyDir)
	if err != nil {
		return "", fmt.Errorf("src-body-dir の読み取りに失敗しました: %w", err)
	}
	bodyFiles := filterMarkdownFiles(bodyPaths)

	resourceFiles, err := s.fileSystem.ListFilesRecursive(trimmedSrcResourceDir)
	if err != nil {
		return "", fmt.Errorf("src-resource-dir の読み取りに失敗しました: %w", err)
	}

	client := s.clientFactory(trimmedBaseURL, trimmedAPIToken)
	if client == nil {
		return "", fmt.Errorf("memos client の初期化に失敗しました")
	}

	ctx := context.Background()
	createdMemos := 0
	attachedFiles := 0
	skippedNoResources := 0

	s.reportProgress("[migrate-to-memos] 開始: 対象body件数=%d", len(bodyFiles))
	for _, bodyFile := range bodyFiles {
		conID := common.ExtractConIDFromPath(bodyFile)
		if conID == "" {
			return "", fmt.Errorf("con_id の抽出に失敗しました: %s", bodyFile)
		}
		s.reportProgress("[migrate-to-memos] con_id=%s メモ作成開始", conID)

		bodyData, err := s.fileSystem.ReadFile(bodyFile)
		if err != nil {
			return "", fmt.Errorf("body ファイルの読み込みに失敗しました (%s): %w", bodyFile, err)
		}

		memoName, err := client.CreateMemo(ctx, string(bodyData))
		if err != nil {
			return "", fmt.Errorf("メモ作成に失敗しました (con_id=%s): %w", conID, err)
		}
		createdMemos++
		s.reportProgress("[migrate-to-memos] con_id=%s メモ作成完了: memo=%s", conID, memoName)

		matchedResources := collectResourceFilesByConID(resourceFiles, conID)
		if len(matchedResources) == 0 {
			skippedNoResources++
			s.reportProgress("[migrate-to-memos] con_id=%s 添付対象なしのためスキップ", conID)
			continue
		}

		if err := client.PatchFiles(ctx, memoName, matchedResources); err != nil {
			return "", fmt.Errorf("添付ファイルの登録に失敗しました (con_id=%s): %w", conID, err)
		}
		attachedFiles += len(matchedResources)
		s.reportProgress("[migrate-to-memos] con_id=%s 添付完了: files=%d", conID, len(matchedResources))
	}

	s.reportProgress(
		"[migrate-to-memos] 完了: 対象body件数=%d メモ作成成功=%d 添付ファイル総数=%d 添付スキップ(リソースなし)=%d",
		len(bodyFiles),
		createdMemos,
		attachedFiles,
		skippedNoResources,
	)
	return fmt.Sprintf(
		"処理完了\n対象body件数=%d\nメモ作成成功=%d\n添付ファイル総数=%d\n添付スキップ(リソースなし)=%d",
		len(bodyFiles),
		createdMemos,
		attachedFiles,
		skippedNoResources,
	), nil
}

func (s *Service) reportProgress(format string, args ...any) {
	if s.reporter == nil {
		return
	}
	s.reporter.Report(fmt.Sprintf(format, args...))
}

func filterMarkdownFiles(paths []string) []string {
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		if strings.EqualFold(filepath.Ext(path), ".md") {
			result = append(result, path)
		}
	}
	return result
}

func collectResourceFilesByConID(resourceFiles []string, conID string) []string {
	result := make([]string, 0)
	for _, resourceFile := range resourceFiles {
		if strings.HasPrefix(filepath.Base(resourceFile), conID) {
			result = append(result, resourceFile)
		}
	}
	return result
}
