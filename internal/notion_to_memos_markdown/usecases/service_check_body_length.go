package usecases

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type bodyLengthResult struct {
	path      string
	runeCount int
}

func (s *Service) CheckBodyLength(srcBodyDir string, threshold int) (string, error) {
	trimmedSrcBodyDir := strings.TrimSpace(srcBodyDir)
	if trimmedSrcBodyDir == "" {
		return "", fmt.Errorf("src-body-dir パラメータは必須です")
	}
	if threshold < 0 {
		return "", fmt.Errorf("threshold パラメータは0以上で必須です")
	}

	filePaths, err := s.fileSystem.ListFilesRecursive(trimmedSrcBodyDir)
	if err != nil {
		return "", fmt.Errorf("src-body-dir の読み取りに失敗しました: %w", err)
	}

	exceededFiles := make([]bodyLengthResult, 0)
	for _, filePath := range filePaths {
		data, err := s.fileSystem.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("ファイルの読み込みに失敗しました (%s): %w", filePath, err)
		}

		runeCount := utf8.RuneCount(data)
		if runeCount > threshold {
			exceededFiles = append(exceededFiles, bodyLengthResult{
				path:      filePath,
				runeCount: runeCount,
			})
		}
	}

	var output strings.Builder
	output.WriteString("処理完了\n")
	output.WriteString(fmt.Sprintf("対象ファイル総数=%d\n", len(filePaths)))
	output.WriteString(fmt.Sprintf("閾値超過ファイル総数=%d\n", len(exceededFiles)))
	output.WriteString("閾値超過ファイル一覧:\n")

	if len(exceededFiles) == 0 {
		output.WriteString("(なし)")
		return output.String(), nil
	}

	for i, file := range exceededFiles {
		if i > 0 {
			output.WriteString("\n")
		}
		output.WriteString(fmt.Sprintf("%s: %d", file.path, file.runeCount))
	}

	return output.String(), nil
}
