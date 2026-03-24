package createclip

import (
	"context"
	"fmt"
	"strings"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	memos "github.com/landmaster135/devbox/internal/memos/usecases"
	common "github.com/landmaster135/devbox/internal/memos_utility/usecases/common"
)

// MemosService は create_clip operation が利用する Memos サービスの契約。
type MemosService interface {
	CreateMemo(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error)
	PatchFiles(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error)
}

// ServiceOptions は Service 生成時の入力。
type ServiceOptions struct {
	MemosService MemosService
	FileSystem   infrastructures.FileSystem
}

// Service は create-web-clip / create-movie-clip の operation を扱う。
type Service struct {
	memosService MemosService
	fileSystem   infrastructures.FileSystem
}

// NewService は Service を生成する。
func NewService(opts ServiceOptions) *Service {
	return &Service{
		memosService: opts.MemosService,
		fileSystem:   opts.FileSystem,
	}
}

// Execute は単一クリップを作成し、必要に応じて添付を追加する。
func (s *Service) Execute(ctx context.Context, input common.CreateClipInput) (*common.CreateClipOutput, error) {
	operation := common.NormalizeOperation(input.Operation)
	contentFile := strings.TrimSpace(input.ContentFile)
	if contentFile == "" {
		return nil, fmt.Errorf("content-file パラメータは必須です")
	}

	displayTime, err := common.BuildDisplayTime(operation, contentFile)
	if err != nil {
		return nil, err
	}

	attachments := common.NormalizeAttachments(input.Attachments)
	if err := s.precheckAttachments(attachments); err != nil {
		return nil, err
	}

	pinned := false
	memo, err := s.memosService.CreateMemo(
		ctx,
		"",
		"",
		contentFile,
		"PRIVATE",
		"NORMAL",
		&pinned,
		displayTime,
	)
	if err != nil {
		return nil, fmt.Errorf("メモの作成に失敗しました: %w", err)
	}
	if memo == nil {
		return nil, fmt.Errorf("メモ作成結果が空です")
	}

	result := &common.CreateClipOutput{
		Operation:   operation,
		DisplayTime: displayTime,
		Memo:        memo,
	}

	if len(attachments) == 0 {
		return result, nil
	}

	memoID, err := common.ResolveMemoIdentifier(memo)
	if err != nil {
		return nil, fmt.Errorf("メモの作成には成功しましたが、添付対象メモの識別子を取得できません: %w", err)
	}

	setOutput, err := s.memosService.PatchFiles(ctx, memoID, attachments, false)
	if err != nil {
		return nil, fmt.Errorf("メモの作成には成功しましたが、添付の追加に失敗しました: %w", err)
	}

	result.Attachments = attachments
	result.SetMemoAttachments = setOutput
	return result, nil
}

func (s *Service) precheckAttachments(attachments []string) error {
	for _, attachment := range attachments {
		if _, err := s.fileSystem.ReadAttachmentFile(attachment); err != nil {
			return fmt.Errorf("--attachments で指定されたファイルが不正です。メモは作成されませんでした (%s): %w", attachment, err)
		}
	}

	return nil
}
