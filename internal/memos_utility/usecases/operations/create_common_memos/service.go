package createcommonmemos

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	memos "github.com/landmaster135/devbox/internal/memos/usecases"
	common "github.com/landmaster135/devbox/internal/memos_utility/usecases/common"
)

// MemosService は create-common-memos operation が利用する Memos サービスの契約。
type MemosService interface {
	CreateMemo(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error)
	PatchFiles(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error)
	AddMemoRelations(ctx context.Context, memo string, relatedMemos []string, replaces bool) (*memos.AddMemoRelationsOutput, error)
}

// ServiceOptions は Service 生成時の入力。
type ServiceOptions struct {
	MemosService     MemosService
	FileSystem       infrastructures.FileSystem
	ProgressReporter func(progress common.CreateClipsProgress)
}

// Service は create-common-memos operation を扱う。
type Service struct {
	memosService     MemosService
	fileSystem       infrastructures.FileSystem
	progressReporter func(progress common.CreateClipsProgress)
}

type commonMemoKey struct {
	timestamp string
	number    int
}

type commonMemoTarget struct {
	key         commonMemoKey
	contentPath string
	baseName    string
	displayTime string
}

// NewService は Service を生成する。
func NewService(opts ServiceOptions) *Service {
	return &Service{
		memosService:     opts.MemosService,
		fileSystem:       opts.FileSystem,
		progressReporter: opts.ProgressReporter,
	}
}

// Execute は content-dir 配下の共通メモファイルを走査し、メモを一括作成する。
func (s *Service) Execute(ctx context.Context, input common.CreateCommonMemosInput) (*common.CreateCommonMemosOutput, error) {
	operation := common.NormalizeOperation(input.Operation)
	if operation != common.OperationCreateCommonMemos {
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

	targets, err := resolveCommonMemoTargets(contentFiles)
	if err != nil {
		return nil, err
	}

	attachmentDir := strings.TrimSpace(input.AttachmentDir)
	attachmentsByMemo := make(map[commonMemoKey][]string)
	if attachmentDir != "" {
		attachmentFiles, err := listRegularFilesInDir(attachmentDir, "attachment-dir")
		if err != nil {
			return nil, err
		}

		attachmentsByMemo, err = resolveAttachmentsByMemo(attachmentFiles)
		if err != nil {
			return nil, err
		}
		if err := s.precheckAttachmentFiles(flattenAttachmentFiles(attachmentsByMemo)); err != nil {
			return nil, err
		}
	}

	output := &common.CreateCommonMemosOutput{
		Operation:     operation,
		ContentDir:    contentDir,
		AttachmentDir: attachmentDir,
		Memos:         make([]*common.CreateCommonMemoOutput, 0, len(targets)),
	}

	memoIDByKey := make(map[commonMemoKey]string, len(targets))
	resultByKey := make(map[commonMemoKey]*common.CreateCommonMemoOutput, len(targets))

	for _, target := range targets {
		attachments := attachmentsByMemo[target.key]
		if s.progressReporter != nil {
			s.progressReporter(common.CreateClipsProgress{
				Current:         len(output.Memos) + 1,
				Total:           len(targets),
				Operation:       operation,
				ContentFile:     target.contentPath,
				AttachmentCount: len(attachments),
			})
		}

		pinned := false
		memo, err := s.memosService.CreateMemo(
			ctx,
			"",
			"",
			target.contentPath,
			"PRIVATE",
			"NORMAL",
			&pinned,
			target.displayTime,
		)
		if err != nil {
			return nil, fmt.Errorf("content-file %s のメモ作成に失敗しました: %w", target.baseName, err)
		}
		if memo == nil {
			return nil, fmt.Errorf("content-file %s のメモ作成結果が空です", target.baseName)
		}

		memoID, err := common.ResolveMemoIdentifier(memo)
		if err != nil {
			return nil, fmt.Errorf("content-file %s のメモ識別子を取得できません: %w", target.baseName, err)
		}
		memoIDByKey[target.key] = memoID

		memoOutput := &common.CreateCommonMemoOutput{
			ContentFile: target.contentPath,
			DisplayTime: target.displayTime,
			Memo:        memo,
		}

		if len(attachments) > 0 {
			setOutput, err := s.memosService.PatchFiles(ctx, memoID, attachments, false)
			if err != nil {
				return nil, fmt.Errorf("content-file %s のメモ作成には成功しましたが、添付の追加に失敗しました: %w", target.baseName, err)
			}
			memoOutput.Attachments = attachments
			memoOutput.SetMemoAttachments = setOutput
		}

		output.Memos = append(output.Memos, memoOutput)
		resultByKey[target.key] = memoOutput
	}

	for _, target := range targets {
		prevKey := commonMemoKey{timestamp: target.key.timestamp, number: target.key.number - 1}
		prevMemoID, ok := memoIDByKey[prevKey]
		if !ok {
			continue
		}
		currentMemoID := memoIDByKey[target.key]
		if _, err := s.memosService.AddMemoRelations(ctx, currentMemoID, []string{prevMemoID}, false); err != nil {
			return nil, fmt.Errorf("content-file %s の relation 追加に失敗しました: %w", target.baseName, err)
		}
		if result := resultByKey[target.key]; result != nil {
			result.RelatedToPreviousBy = prevMemoID
		}
	}

	output.Total = len(output.Memos)
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

func resolveCommonMemoTargets(contentFiles []string) ([]commonMemoTarget, error) {
	targets := make([]commonMemoTarget, 0, len(contentFiles))
	seen := make(map[commonMemoKey]struct{}, len(contentFiles))
	for _, contentFile := range contentFiles {
		baseName := filepath.Base(contentFile)
		timestamp, number, ok := common.ParseCommonMemoFileKey(baseName)
		if !ok {
			return nil, fmt.Errorf("content-dir 内のファイル名が不正です。YYYYMMDDhhmmss_<number>.md のみ指定できます: %s", baseName)
		}

		displayTime, err := common.BuildDisplayTime(common.OperationCreateCommonMemos, contentFile)
		if err != nil {
			return nil, err
		}

		key := commonMemoKey{timestamp: timestamp, number: number}
		if _, exists := seen[key]; exists {
			return nil, fmt.Errorf("content-dir 内で同一キーのファイルが重複しています: %s", baseName)
		}
		seen[key] = struct{}{}

		targets = append(targets, commonMemoTarget{
			key:         key,
			contentPath: contentFile,
			baseName:    baseName,
			displayTime: displayTime,
		})
	}

	sort.Slice(targets, func(i, j int) bool {
		if targets[i].key.timestamp != targets[j].key.timestamp {
			return targets[i].key.timestamp < targets[j].key.timestamp
		}
		if targets[i].key.number != targets[j].key.number {
			return targets[i].key.number < targets[j].key.number
		}
		return targets[i].baseName < targets[j].baseName
	})

	return targets, nil
}

func resolveAttachmentsByMemo(attachmentFiles []string) (map[commonMemoKey][]string, error) {
	attachmentsByMemo := make(map[commonMemoKey][]string)
	for _, attachmentFile := range attachmentFiles {
		baseName := filepath.Base(attachmentFile)
		timestamp, number, ok := common.ParseCommonMemoAttachmentKey(baseName)
		if !ok {
			return nil, fmt.Errorf("attachment-dir 内のファイル名が不正です。YYYYMMDDhhmmss_<number>_<index>.<extension> のみ指定できます: %s", baseName)
		}
		key := commonMemoKey{timestamp: timestamp, number: number}
		attachmentsByMemo[key] = append(attachmentsByMemo[key], attachmentFile)
	}

	for key := range attachmentsByMemo {
		sort.Strings(attachmentsByMemo[key])
	}
	return attachmentsByMemo, nil
}

func flattenAttachmentFiles(attachmentsByMemo map[commonMemoKey][]string) []string {
	allAttachments := make([]string, 0)
	for _, attachments := range attachmentsByMemo {
		allAttachments = append(allAttachments, attachments...)
	}
	sort.Strings(allAttachments)
	return allAttachments
}
