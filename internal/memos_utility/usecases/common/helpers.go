package common

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	memos "github.com/landmaster135/devbox/internal/memos/usecases"
)

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

func ResolveContentBaseNameFromAttachment(baseName string) (string, error) {
	if contentBaseName, ok := ParseWebAttachmentContentBaseName(baseName); ok {
		return contentBaseName, nil
	}
	if contentBaseName, ok := ParseMovieAttachmentContentBaseName(baseName); ok {
		return contentBaseName, nil
	}
	return "", fmt.Errorf("attachment-dir 内のファイル名が不正です。web-summary-YYYYMMDD-hhmmss-<slug>_<number>.<extension> または movie-summary-YYYYMMDD-hhmmss-<slug>_<number>.<extension> のみ指定できます: %s", baseName)
}
