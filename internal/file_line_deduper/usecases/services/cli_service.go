package services

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/file_line_deduper/infrastructures/filesystem"
)

type lineDeduper interface {
	RemoveMatchingLines(filePath string, startPos, endPos int) (int, error)
}

// CLIService はCLI向けの入出力境界をまとめるサービスです。
type CLIService struct {
	fileService lineDeduper
}

// NewCLIService はCLI向けサービスを生成します。
func NewCLIService() *CLIService {
	fileRepo := filesystem.NewRepository()
	fileService := NewFileService(fileRepo)

	return &CLIService{
		fileService: fileService,
	}
}

// NewCLIServiceWithFileService はテスト用に依存を注入してCLIServiceを作成します。
func NewCLIServiceWithFileService(fileService lineDeduper) *CLIService {
	return &CLIService{
		fileService: fileService,
	}
}

// HandleRemoveMatchingLines は重複行削除処理を実行し、CLI表示用メッセージを返します。
func (s *CLIService) HandleRemoveMatchingLines(filePath string, startPos, endPos int) (string, error) {
	count, err := s.fileService.RemoveMatchingLines(filePath, startPos, endPos)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("処理完了: %d行の重複を削除しました\n", count), nil
}
