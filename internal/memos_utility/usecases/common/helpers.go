package common

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	memos "github.com/landmaster135/devbox/internal/memos/usecases"
)

type MemoIdentifierSource string

const (
	MemoIdentifierSourceName MemoIdentifierSource = "name"
	MemoIdentifierSourceUID  MemoIdentifierSource = "uid"
	MemoIdentifierSourceID   MemoIdentifierSource = "id"
)

type ResolvedMemoIdentifier struct {
	Identifier string
	Source     MemoIdentifierSource
}

var clipSlugPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

func NormalizeOperation(operation string) string {
	return strings.ToLower(strings.TrimSpace(operation))
}

func NormalizeAttachments(values []string) []string {
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

func BuildDisplayTime(operation, contentFile string) (string, error) {
	baseName := filepath.Base(contentFile)

	var (
		timestamp string
		ok        bool
	)

	switch operation {
	case OperationCreateWebClip:
		timestamp, ok = ParseWebClipDisplayTime(baseName)
	case OperationCreateMovieClip:
		timestamp, ok = ParseMovieClipDisplayTime(baseName)
	case OperationCreateCommonMemos:
		timestamp, ok = ParseCommonMemoDisplayTime(baseName)
	default:
		return "", fmt.Errorf("未対応の operation です: %s", operation)
	}

	if !ok {
		return "", fmt.Errorf("operation=%s に対する content-file の形式が不正です: %s", operation, baseName)
	}

	parsedTime, err := time.ParseInLocation("20060102150405", timestamp, jstLocation)
	if err != nil {
		return "", fmt.Errorf("content-file の日時解析に失敗しました: %w", err)
	}

	return parsedTime.Format(time.RFC3339), nil
}

func ResolveMemoIdentifier(memo *memos.Memo) (string, error) {
	resolved, err := ResolveMemoIdentifierWithSource(memo)
	if err != nil {
		return "", err
	}
	return resolved.Identifier, nil
}

func ResolveMemoIdentifierWithSource(memo *memos.Memo) (ResolvedMemoIdentifier, error) {
	if memo == nil {
		return ResolvedMemoIdentifier{}, fmt.Errorf("メモ情報が空です")
	}

	if memoName := strings.TrimSpace(memo.Name); memoName != "" {
		return ResolvedMemoIdentifier{
			Identifier: memoName,
			Source:     MemoIdentifierSourceName,
		}, nil
	}
	if memoUID := strings.TrimSpace(memo.UID); memoUID != "" {
		return ResolvedMemoIdentifier{
			Identifier: memoUID,
			Source:     MemoIdentifierSourceUID,
		}, nil
	}
	if memo.ID > 0 {
		return ResolvedMemoIdentifier{
			Identifier: strconv.FormatInt(memo.ID, 10),
			Source:     MemoIdentifierSourceID,
		}, nil
	}

	return ResolvedMemoIdentifier{}, fmt.Errorf("name/uid/id のいずれも取得できません")
}

func ResolveContentBaseNameFromAttachment(baseName string) (string, error) {
	if contentBaseName, ok := ParseWebAttachmentContentBaseName(baseName); ok {
		return contentBaseName, nil
	}
	if contentBaseName, ok := ParseMovieAttachmentContentBaseName(baseName); ok {
		return contentBaseName, nil
	}
	return "", fmt.Errorf("attachment-dir 内のファイル名が不正です。web-summary-YYYYMMDD-hhmmss-<slug>_<number>.<extension> または movie-summary-YYYYMMDD-hhmmss-<slug>_<number>.<extension> のみ指定できます: %s", baseName)
}

func ExplainClipContentFileNameIssue(baseName string) string {
	expected := "web-summary-YYYYMMDD-hhmmss-<slug>.md または movie-summary-YYYYMMDD-hhmmss-<slug>.md"
	trimmed := strings.TrimSpace(baseName)
	if trimmed == "" {
		return fmt.Sprintf("直す箇所: ファイル名が空です。期待形式: %s", expected)
	}

	if filepath.Ext(trimmed) != ".md" {
		return fmt.Sprintf("直す箇所: 拡張子を .md にしてください。期待形式: %s", expected)
	}

	stem := strings.TrimSuffix(trimmed, ".md")
	rest := ""
	switch {
	case strings.HasPrefix(stem, "web-summary-"):
		rest = strings.TrimPrefix(stem, "web-summary-")
	case strings.HasPrefix(stem, "movie-summary-"):
		rest = strings.TrimPrefix(stem, "movie-summary-")
	default:
		return fmt.Sprintf("直す箇所: 接頭辞を web-summary- または movie-summary- にしてください。期待形式: %s", expected)
	}

	parts := strings.SplitN(rest, "-", 3)
	if len(parts) != 3 {
		return fmt.Sprintf("直す箇所: 日付 YYYYMMDD・時刻 hhmmss・slug をハイフン区切りで指定してください。期待形式: %s", expected)
	}

	datePart := parts[0]
	timePart := parts[1]
	slugPart := parts[2]

	if len(datePart) != 8 || !isDigits(datePart) {
		return fmt.Sprintf("直す箇所: 日付部 YYYYMMDD を8桁で指定してください。現在: %s", datePart)
	}
	if len(timePart) != 6 || !isDigits(timePart) {
		return fmt.Sprintf("直す箇所: 時刻部 hhmmss を6桁で指定してください。現在: %s", timePart)
	}
	if slugPart == "" {
		return "直す箇所: slug が不足しています。時刻の後ろに -<slug> を付けてください。"
	}
	if !clipSlugPattern.MatchString(slugPart) {
		return fmt.Sprintf("直す箇所: slug は英数字で開始し、英数字・ハイフン・アンダースコアのみ使用できます。現在: %s", slugPart)
	}

	return fmt.Sprintf("直す箇所: 期待形式に合わせてファイル名を修正してください。期待形式: %s", expected)
}

func isDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
