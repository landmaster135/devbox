package services

import (
	"fmt"
	"time"

	domainRepo "github.com/landmaster135/devbox/internal/independencies/json_iso8601_converter/domain/repositories"
)

// ISO8601Service はJSONファイル内のISO8601形式の日時文字列をUNIXタイムスタンプに変換するためのサービスです
type ISO8601Service struct {
	jsonRepo domainRepo.JSONRepository
}

// NewISO8601Service は新しいISO8601Serviceインスタンスを作成します
func NewISO8601Service(jsonRepo domainRepo.JSONRepository) *ISO8601Service {
	return &ISO8601Service{
		jsonRepo: jsonRepo,
	}
}

// ConvertISO8601ToTimestamp はディレクトリ内のJSONファイルのキーの値をISO8601形式からUNIXタイムスタンプに変換します
func (s *ISO8601Service) ConvertISO8601ToTimestamp(dirPath, key string, recursive, dryRun bool) (int, error) {
	// 処理したファイル数
	processedCount := 0

	// ディレクトリを処理する関数
	processDir := func(currentDir string) error {
		// ディレクトリ内のJSONファイルを検索
		jsonFiles, err := s.findJSONFiles(currentDir, recursive)
		if err != nil {
			return fmt.Errorf("JSONファイルの検索に失敗しました: %w", err)
		}

		// 各JSONファイルを処理
		for _, path := range jsonFiles {
			// JSONファイルを処理
			converted, err := s.convertFile(path, key, dryRun)
			if err != nil {
				return fmt.Errorf("ファイル '%s' の処理中にエラーが発生しました: %w", path, err)
			}

			if converted {
				processedCount++
			}
		}

		return nil
	}

	// ディレクトリを処理
	if err := processDir(dirPath); err != nil {
		return processedCount, err
	}

	return processedCount, nil
}

// convertFile は単一のJSONファイルを処理します
func (s *ISO8601Service) convertFile(filePath, key string, dryRun bool) (bool, error) {
	return s.jsonRepo.ConvertFile(filePath, key, dryRun)
}

// processJSONData はJSONデータを再帰的に処理します
func (s *ISO8601Service) processJSONData(data interface{}, targetKey string) (interface{}, bool) {
	return s.jsonRepo.ProcessJSONData(data, targetKey)
}

// findJSONFiles はディレクトリ内のJSONファイルを検索します
func (s *ISO8601Service) findJSONFiles(dirPath string, recursive bool) ([]string, error) {
	return s.jsonRepo.FindJSONFiles(dirPath, recursive)
}

// parseISO8601 はISO8601形式の日時文字列をUNIXタイムスタンプに変換します
func (s *ISO8601Service) parseISO8601(dateStr string) (int64, error) {
	// 複数のフォーマットを試す
	formats := []string{
		time.RFC3339,      // 2006-01-02T15:04:05Z07:00
		time.RFC3339Nano,  // 2006-01-02T15:04:05.999999999Z07:00
		"2006-01-02 15:04:05Z07:00",      // スペース区切りの標準形式
		"2006-01-02 15:04:05.999999999Z07:00", // スペース区切りのナノ秒形式
		"2006-01-02 15:04:05-07:00",      // スペース区切り、タイムゾーンあり
		"2006-01-02 15:04:05+07:00",      // スペース区切り、タイムゾーンあり（+）
		"2006-01-02 15:04:05.999999999-07:00", // スペース区切り、ナノ秒、タイムゾーンあり
		"2006-01-02 15:04:05.999999999+07:00", // スペース区切り、ナノ秒、タイムゾーンあり（+）
		"2006-01-02 15:04:05-0700",       // スペース区切り、タイムゾーンあり（省略形式）
		"2006-01-02 15:04:05+0700",       // スペース区切り、タイムゾーンあり（省略形式、+）
		"2006-01-02 15:04:05.999999999-0700", // スペース区切り、ナノ秒、タイムゾーンあり（省略形式）
		"2006-01-02 15:04:05.999999999+0700", // スペース区切り、ナノ秒、タイムゾーンあり（省略形式、+）
	}

	var parseErr error
	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			// UNIXタイムスタンプに変換
			return t.Unix(), nil
		}
		parseErr = err
	}

	// 特殊なケース: PostgreSQLの日時形式（例: 2025-04-10 09:42:29.059877+00）
	if len(dateStr) > 2 && (dateStr[len(dateStr)-3:] == "+00" || dateStr[len(dateStr)-3:] == "-00") {
		// +00 または -00 を +0000 または -0000 に変換
		modifiedDateStr := dateStr[:len(dateStr)-3] + "+0000"
		t, err := time.Parse("2006-01-02 15:04:05.999999999-0700", modifiedDateStr)
		if err == nil {
			return t.Unix(), nil
		}
		// 秒までの場合
		t, err = time.Parse("2006-01-02 15:04:05-0700", modifiedDateStr)
		if err == nil {
			return t.Unix(), nil
		}
	}

	return 0, fmt.Errorf("日時文字列のパースに失敗しました: %w", parseErr)
}
