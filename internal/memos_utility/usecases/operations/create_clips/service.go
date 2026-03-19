package createclips

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	"github.com/landmaster135/devbox/internal/memos_utility/usecases/common"
)

// CreateClipService は単一クリップ生成 operation の契約。
type CreateClipService interface {
	Execute(ctx context.Context, input common.CreateClipInput) (*common.CreateClipOutput, error)
}

// ServiceOptions は Service 生成時の入力。
type ServiceOptions struct {
	CreateClipService CreateClipService
	FileSystem        infrastructures.FileSystem
	ProgressReporter  func(progress common.CreateClipsProgress)
}

// Service は create-clips operation を扱う。
type Service struct {
	createClipService CreateClipService
	fileSystem        infrastructures.FileSystem
	progressReporter  func(progress common.CreateClipsProgress)
}

// NewService は Service を生成する。
func NewService(opts ServiceOptions) *Service {
	return &Service{
		createClipService: opts.CreateClipService,
		fileSystem:        opts.FileSystem,
		progressReporter:  opts.ProgressReporter,
	}
}

type clipTarget struct {
	Operation       string
	ContentPath     string
	ContentBaseName string
}

// Execute は content-dir 配下の clip ファイルを走査し、メモを一括作成する。
func (s *Service) Execute(ctx context.Context, input common.CreateClipsInput) (*common.CreateClipsOutput, error) {
	operation := common.NormalizeOperation(input.Operation)
	if operation != common.OperationCreateClips {
		return nil, fmt.Errorf("未対応の operation です: %s", operation)
	}

	contentDir := strings.TrimSpace(input.ContentDir)
	if contentDir == "" {
		return nil, fmt.Errorf("content-dir パラメータは必須です")
	}

	contentFiles, err := listRegularFilesInDir(contentDir, "content-dir")
	if err != nil {
		return nil, err
	}
	if len(contentFiles) == 0 {
		return nil, fmt.Errorf("content-dir に処理対象ファイルがありません: %s", contentDir)
	}

	targets, err := resolveClipTargets(contentFiles)
	if err != nil {
		return nil, err
	}

	attachmentDir := strings.TrimSpace(input.AttachmentDir)
	attachmentsByContent := make(map[string][]string)
	if attachmentDir != "" {
		attachmentFiles, err := listRegularFilesInDir(attachmentDir, "attachment-dir")
		if err != nil {
			return nil, err
		}

		attachmentsByContent, err = resolveAttachmentsByContent(attachmentFiles)
		if err != nil {
			return nil, err
		}
		if err := s.precheckAttachmentFiles(flattenAttachmentFiles(attachmentsByContent)); err != nil {
			return nil, err
		}
	}

	output := &common.CreateClipsOutput{
		Operation:     operation,
		ContentDir:    contentDir,
		AttachmentDir: attachmentDir,
		Clips:         make([]*common.CreateClipOutput, 0, len(targets)),
	}

	for index, target := range targets {
		attachments := attachmentsByContent[target.ContentBaseName]
		if s.progressReporter != nil {
			s.progressReporter(common.CreateClipsProgress{
				Current:         index + 1,
				Total:           len(targets),
				Operation:       target.Operation,
				ContentFile:     target.ContentPath,
				AttachmentCount: len(attachments),
			})
		}

		clipOutput, err := s.createClipService.Execute(ctx, common.CreateClipInput{
			Operation:   target.Operation,
			ContentFile: target.ContentPath,
			Attachments: attachments,
		})
		if err != nil {
			return nil, fmt.Errorf("content-file %s の処理に失敗しました: %w", filepath.Base(target.ContentPath), err)
		}
		output.Clips = append(output.Clips, clipOutput)
	}

	output.Total = len(output.Clips)
	return output, nil
}

func (s *Service) precheckAttachmentFiles(attachments []string) error {
	for _, attachment := range attachments {
		if _, err := s.fileSystem.ReadAttachmentFile(attachment); err != nil {
			return fmt.Errorf("--attachment-dir で指定されたファイルが不正です。メモは作成されませんでした (%s): %w", attachment, err)
		}
	}

	return nil
}

func listRegularFilesInDir(dirPath, flagName string) ([]string, error) {
	cleanDirPath := strings.TrimSpace(dirPath)
	if cleanDirPath == "" {
		return nil, fmt.Errorf("%s パラメータは必須です", flagName)
	}

	info, err := os.Stat(cleanDirPath)
	if err != nil {
		return nil, fmt.Errorf("%s の検証に失敗しました (%s): %w", flagName, cleanDirPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s はディレクトリを指定してください: %s", flagName, cleanDirPath)
	}

	entries, err := os.ReadDir(cleanDirPath)
	if err != nil {
		return nil, fmt.Errorf("%s の読み取りに失敗しました (%s): %w", flagName, cleanDirPath, err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, filepath.Join(cleanDirPath, entry.Name()))
	}
	sort.Strings(files)
	return files, nil
}

func resolveClipTargets(contentFiles []string) ([]clipTarget, error) {
	targets := make([]clipTarget, 0, len(contentFiles))
	for _, contentFile := range contentFiles {
		baseName := filepath.Base(contentFile)
		contentBaseName := strings.TrimSuffix(baseName, filepath.Ext(baseName))

		switch {
		case common.MatchWebClipFile(baseName):
			targets = append(targets, clipTarget{
				Operation:       common.OperationCreateWebClip,
				ContentPath:     contentFile,
				ContentBaseName: contentBaseName,
			})
		case common.MatchMovieClipFile(baseName):
			targets = append(targets, clipTarget{
				Operation:       common.OperationCreateMovieClip,
				ContentPath:     contentFile,
				ContentBaseName: contentBaseName,
			})
		default:
			return nil, fmt.Errorf("content-dir 内のファイル名が不正です。web-summary-YYYYMMDD-hhmmss-<slug>.md または movie-summary-YYYYMMDD-hhmmss-<slug>.md のみ指定できます: %s", baseName)
		}
	}
	return targets, nil
}

func resolveAttachmentsByContent(attachmentFiles []string) (map[string][]string, error) {
	attachmentsByContent := make(map[string][]string)
	for _, attachmentFile := range attachmentFiles {
		contentBaseName, err := common.ResolveContentBaseNameFromAttachment(filepath.Base(attachmentFile))
		if err != nil {
			return nil, err
		}
		attachmentsByContent[contentBaseName] = append(attachmentsByContent[contentBaseName], attachmentFile)
	}

	for contentBaseName := range attachmentsByContent {
		sort.Strings(attachmentsByContent[contentBaseName])
	}
	return attachmentsByContent, nil
}

func flattenAttachmentFiles(attachmentsByContent map[string][]string) []string {
	allAttachments := make([]string, 0)
	for _, attachments := range attachmentsByContent {
		allAttachments = append(allAttachments, attachments...)
	}
	sort.Strings(allAttachments)
	return allAttachments
}
