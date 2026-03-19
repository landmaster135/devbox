package usecases

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	memos "github.com/landmaster135/devbox/internal/memos/usecases"
)

const (
	operationCreateWebClip   = "create-web-clip"
	operationCreateMovieClip = "create-movie-clip"
)

var webClipFilePattern = regexp.MustCompile(`^web-summary-(\d{8})-(\d{6})-([A-Za-z0-9][A-Za-z0-9_-]*)\.md$`)
var movieClipFilePattern = regexp.MustCompile(`^movie-summary-(\d{8})-(\d{6})-([A-Za-z0-9][A-Za-z0-9_-]*)\.md$`)

var jstLocation = time.FixedZone("JST", 9*60*60)

// ServiceOptions は Service 生成時の入力。
type ServiceOptions struct {
	BaseURL      string
	APIToken     string
	Timeout      time.Duration
	MemosService MemosService
	FileSystem   infrastructures.FileSystem
}

// CreateClipInput は create-web-clip / create-movie-clip の入力。
type CreateClipInput struct {
	Operation   string
	ContentFile string
	Attachments []string
}

// CreateClipOutput は create-web-clip / create-movie-clip の出力。
type CreateClipOutput struct {
	Operation          string                          `json:"operation"`
	DisplayTime        string                          `json:"displayTime"`
	Memo               *memos.Memo                     `json:"memo,omitempty"`
	Attachments        []string                        `json:"attachments,omitempty"`
	SetMemoAttachments *memos.SetMemoAttachmentsOutput `json:"setMemoAttachments,omitempty"`
}

// Service は memos-utility のユースケースを提供する。
type Service struct {
	memosService MemosService
	fileSystem   infrastructures.FileSystem
}

// NewService は Service を生成する。
func NewService(opts ServiceOptions) *Service {
	service := opts.MemosService
	if service == nil {
		service = memos.NewService(memos.ServiceOptions{
			BaseURL:  opts.BaseURL,
			APIToken: opts.APIToken,
			Timeout:  opts.Timeout,
		})
	}

	fileSystem := opts.FileSystem
	if fileSystem == nil {
		fileSystem = infrastructures.NewOSFileSystem()
	}

	return &Service{
		memosService: service,
		fileSystem:   fileSystem,
	}
}

// CreateClip はメモを作成し、必要に応じて添付を追加する。
func (s *Service) CreateClip(ctx context.Context, input CreateClipInput) (*CreateClipOutput, error) {
	operation := normalizeOperation(input.Operation)
	contentFile := strings.TrimSpace(input.ContentFile)
	if contentFile == "" {
		return nil, fmt.Errorf("content-file パラメータは必須です")
	}

	displayTime, err := buildDisplayTime(operation, contentFile)
	if err != nil {
		return nil, err
	}

	attachments := normalizeAttachments(input.Attachments)
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

	result := &CreateClipOutput{
		Operation:   operation,
		DisplayTime: displayTime,
		Memo:        memo,
	}

	if len(attachments) == 0 {
		return result, nil
	}

	memoID, err := resolveMemoIdentifier(memo)
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

func buildDisplayTime(operation, contentFile string) (string, error) {
	baseName := filepath.Base(contentFile)

	var pattern *regexp.Regexp
	switch operation {
	case operationCreateWebClip:
		pattern = webClipFilePattern
	case operationCreateMovieClip:
		pattern = movieClipFilePattern
	default:
		return "", fmt.Errorf("未対応の operation です: %s", operation)
	}

	matches := pattern.FindStringSubmatch(baseName)
	if len(matches) != 4 {
		return "", fmt.Errorf("operation=%s に対する content-file の形式が不正です: %s", operation, baseName)
	}

	timestamp := matches[1] + matches[2]
	parsedTime, err := time.ParseInLocation("20060102150405", timestamp, jstLocation)
	if err != nil {
		return "", fmt.Errorf("content-file の日時解析に失敗しました: %w", err)
	}

	return parsedTime.Format(time.RFC3339), nil
}

func normalizeOperation(operation string) string {
	return strings.ToLower(strings.TrimSpace(operation))
}

func normalizeAttachments(values []string) []string {
	attachments := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		attachments = append(attachments, trimmed)
	}
	if len(attachments) == 0 {
		return nil
	}
	return attachments
}

func resolveMemoIdentifier(memo *memos.Memo) (string, error) {
	if memo == nil {
		return "", fmt.Errorf("メモ情報が空です")
	}

	if memoName := strings.TrimSpace(memo.Name); memoName != "" {
		return memoName, nil
	}
	if memoUID := strings.TrimSpace(memo.UID); memoUID != "" {
		return memoUID, nil
	}
	if memo.ID > 0 {
		return strconv.FormatInt(memo.ID, 10), nil
	}

	return "", fmt.Errorf("name/uid/id のいずれも取得できません")
}

var _ MemosUtilityService = (*Service)(nil)
