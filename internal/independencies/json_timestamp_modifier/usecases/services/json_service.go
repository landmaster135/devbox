package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/landmaster135/devbox/internal/independencies/json_timestamp_modifier/domain/models"
	"github.com/landmaster135/devbox/internal/independencies/json_timestamp_modifier/domain/repositories"
)

// JSONService はJSONファイルを操作するためのサービスです
type JSONService struct {
	FileRepo repositories.FileRepository
}

// NewJSONService は新しいJSONServiceインスタンスを作成します
func NewJSONService(fileRepo repositories.FileRepository) *JSONService {
	return &JSONService{
		FileRepo: fileRepo,
	}
}

// AddKeyValue は指定されたJSONファイルに新しいキーと値を追加します
func (s *JSONService) AddKeyValue(filePath string, key string, value interface{}) error {
	if key == "" {
		return errors.New("キーは空にできません")
	}

	var data map[string]interface{}

	// ファイルが存在するか確認
	if s.FileRepo.FileExists(filePath) {
		// ファイルが存在する場合は読み込む
		jsonData, err := s.FileRepo.ReadJSONFile(filePath)
		if err != nil {
			// JSONとして読み込めない場合は新しいマップを作成
			data = make(map[string]interface{})
		} else {
			// 型アサーション
			var ok bool
			data, ok = jsonData.(map[string]interface{})
			if !ok {
				return fmt.Errorf("JSONデータをマップに変換できません")
			}
		}
	} else {
		// ファイルが存在しない場合は新しいマップを作成
		data = make(map[string]interface{})
	}

	// キーと値を追加
	data[key] = value

	// JSONに変換
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("JSONへの変換に失敗しました: %w", err)
	}

	// ファイルに書き込む
	lines := strings.Split(string(jsonBytes), "\n")
	content := models.NewFileContent(lines)
	if err := s.FileRepo.WriteFile(filePath, content); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました: %w", err)
	}

	return nil
}

// GetValue は指定されたJSONファイルから特定のキーの値を取得します
func (s *JSONService) GetValue(filePath string, key string) (interface{}, error) {
	// ファイルが存在するか確認
	if !s.FileRepo.FileExists(filePath) {
		return nil, errors.New("ファイルが存在しません")
	}

	// JSONファイルを読み込む
	jsonData, err := s.FileRepo.ReadJSONFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("JSONファイルの読み込みに失敗しました: %w", err)
	}

	// 型アサーション
	data, ok := jsonData.(map[string]interface{})
	if !ok {
		return nil, errors.New("JSONデータをマップに変換できません")
	}

	// キーが存在するか確認
	value, exists := data[key]
	if !exists {
		return nil, fmt.Errorf("キー '%s' が存在しません", key)
	}

	return value, nil
}

// GetAllData は指定されたJSONファイルのすべてのデータを取得します
func (s *JSONService) GetAllData(filePath string) (map[string]interface{}, error) {
	// ファイルが存在するか確認
	if !s.FileRepo.FileExists(filePath) {
		return nil, errors.New("ファイルが存在しません")
	}

	// JSONファイルを読み込む
	jsonData, err := s.FileRepo.ReadJSONFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("JSONファイルの読み込みに失敗しました: %w", err)
	}

	// 型アサーション
	data, ok := jsonData.(map[string]interface{})
	if !ok {
		return nil, errors.New("JSONデータをマップに変換できません")
	}

	return data, nil
}

// ConvertISO8601ToUnix はJSONファイル内の指定したキーの値をISO-8601形式からUNIXタイムスタンプに変換します
func (s *JSONService) ConvertISO8601ToUnix(filePath string, key string, isJst bool) error {
	// ファイルが存在するか確認
	if !s.FileRepo.FileExists(filePath) {
		return errors.New("ファイルが存在しません")
	}

	// JSONファイルを読み込む
	jsonData, err := s.FileRepo.ReadJSONFile(filePath)
	if err != nil {
		return fmt.Errorf("JSONファイルの読み込みに失敗しました: %w", err)
	}

	// 型アサーション
	data, ok := jsonData.(map[string]interface{})
	if !ok {
		return errors.New("JSONデータをマップに変換できません")
	}

	// キーが存在するか確認
	value, exists := data[key]
	if !exists {
		return fmt.Errorf("キー '%s' が存在しません", key)
	}

	// 文字列型かどうか確認
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("キー '%s' の値が文字列ではありません", key)
	}

	var unixTimestamp int64

	// 時刻情報が含まれているかどうかを確認
	if strings.Contains(strValue, "T") {
		// ISO-8601形式からUNIXタイムスタンプに変換
		unixStr, err := ISO8601ToUnix(strValue, isJst)
		if err != nil {
			return fmt.Errorf("ISO-8601形式からUNIXタイムスタンプへの変換に失敗しました: %w", err)
		}
		unixTimestamp, err = strconv.ParseInt(unixStr, 10, 64)
		if err != nil {
			return fmt.Errorf("UNIXタイムスタンプの解析に失敗しました: %w", err)
		}
	} else {
		// 日付のみの場合、isJstに基づいて変換
		unixStr, err := DateToUnix(strValue, isJst)
		if err != nil {
			return fmt.Errorf("日付からUNIXタイムスタンプへの変換に失敗しました: %w", err)
		}
		unixTimestamp, err = strconv.ParseInt(unixStr, 10, 64)
		if err != nil {
			return fmt.Errorf("UNIXタイムスタンプの解析に失敗しました: %w", err)
		}
	}

	// 変換した値を設定
	data[key] = unixTimestamp

	// JSONに変換
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("JSONへの変換に失敗しました: %w", err)
	}

	// ファイルに書き込む
	lines := strings.Split(string(jsonBytes), "\n")
	content := models.NewFileContent(lines)
	if err := s.FileRepo.WriteFile(filePath, content); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました: %w", err)
	}

	return nil
}

// ConvertUnixToISO8601 はJSONファイル内の指定したキーの値をUNIXタイムスタンプからISO-8601形式に変換します
func (s *JSONService) ConvertUnixToISO8601(filePath string, key string) error {
	// ファイルが存在するか確認
	if !s.FileRepo.FileExists(filePath) {
		return errors.New("ファイルが存在しません")
	}

	// JSONファイルを読み込む
	jsonData, err := s.FileRepo.ReadJSONFile(filePath)
	if err != nil {
		return fmt.Errorf("JSONファイルの読み込みに失敗しました: %w", err)
	}

	// 型アサーション
	data, ok := jsonData.(map[string]interface{})
	if !ok {
		return errors.New("JSONデータをマップに変換できません")
	}

	// キーが存在するか確認
	value, exists := data[key]
	if !exists {
		return fmt.Errorf("キー '%s' が存在しません", key)
	}

	// 数値型に変換
	var numValue int64
	switch v := value.(type) {
	case float64:
		numValue = int64(v)
	case int64:
		numValue = v
	case int:
		numValue = int64(v)
	case string:
		var err error
		numValue, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("キー '%s' の値をUNIXタイムスタンプとして解析できません: %w", key, err)
		}
	default:
		return fmt.Errorf("キー '%s' の値がUNIXタイムスタンプとして解析できません", key)
	}

	// UNIXタイムスタンプからISO-8601形式に変換
	iso8601Str, err := UnixToISO8601(strconv.FormatInt(numValue, 10))
	if err != nil {
		return fmt.Errorf("UNIXタイムスタンプからISO-8601形式への変換に失敗しました: %w", err)
	}

	// 変換した値を設定
	data[key] = iso8601Str

	// JSONに変換
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("JSONへの変換に失敗しました: %w", err)
	}

	// ファイルに書き込む
	lines := strings.Split(string(jsonBytes), "\n")
	content := models.NewFileContent(lines)
	if err := s.FileRepo.WriteFile(filePath, content); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました: %w", err)
	}

	return nil
}

// AddKeyValueToAllFiles は指定されたディレクトリ内の全てのJSONファイルに新しいキーと値を追加します
func (s *JSONService) AddKeyValueToAllFiles(dirPath string, key string, value interface{}, recursive bool) (int, error) {
	if key == "" {
		return 0, errors.New("キーは空にできません")
	}

	// ディレクトリ内のJSONファイルを検索
	jsonFiles, err := s.FileRepo.FindFilesByExt(dirPath, ".json")
	if err != nil {
		return 0, fmt.Errorf("JSONファイルの検索に失敗しました: %w", err)
	}

	// 再帰的に検索する場合
	if recursive {
		// ディレクトリ内のサブディレクトリを取得
		entries, err := s.FileRepo.ReadDir(dirPath)
		if err != nil {
			return 0, fmt.Errorf("ディレクトリの読み込みに失敗しました: %w", err)
		}

		// 各サブディレクトリを処理
		for _, entry := range entries {
			if entry.IsDir {
				// サブディレクトリ内のJSONファイルを処理
				subCount, err := s.AddKeyValueToAllFiles(entry.Path, key, value, recursive)
				if err != nil {
					return 0, err
				}
				// 処理したファイル数を加算
				jsonFiles = append(jsonFiles, make([]string, subCount)...)
			}
		}
	}

	// 各JSONファイルを処理
	processedCount := 0
	for _, filePath := range jsonFiles {
		err := s.AddKeyValue(filePath, key, value)
		if err != nil {
			// エラーが発生しても処理を続行
			fmt.Printf("ファイル %s の処理中にエラーが発生しました: %v\n", filePath, err)
			continue
		}
		processedCount++
	}

	return processedCount, nil
}

// ConvertISO8601ToUnixInAllFiles は指定されたディレクトリ内の全てのJSONファイルに指定したキーの値をISO-8601形式からUNIXタイムスタンプに変換します
func (s *JSONService) ConvertISO8601ToUnixInAllFiles(dirPath string, key string, isJst bool, recursive bool) (int, error) {
	if key == "" {
		return 0, errors.New("キーは空にできません")
	}

	// ディレクトリ内のJSONファイルを検索
	jsonFiles, err := s.FileRepo.FindFilesByExt(dirPath, ".json")
	if err != nil {
		return 0, fmt.Errorf("JSONファイルの検索に失敗しました: %w", err)
	}

	// 再帰的に検索する場合
	if recursive {
		// ディレクトリ内のサブディレクトリを取得
		entries, err := s.FileRepo.ReadDir(dirPath)
		if err != nil {
			return 0, fmt.Errorf("ディレクトリの読み込みに失敗しました: %w", err)
		}

		// 各サブディレクトリを処理
		for _, entry := range entries {
			if entry.IsDir {
				// サブディレクトリ内のJSONファイルを処理
				subCount, err := s.ConvertISO8601ToUnixInAllFiles(entry.Path, key, isJst, recursive)
				if err != nil {
					return 0, err
				}
				// 処理したファイル数を加算
				jsonFiles = append(jsonFiles, make([]string, subCount)...)
			}
		}
	}

	// 各JSONファイルを処理
	processedCount := 0
	for _, filePath := range jsonFiles {
		err := s.ConvertISO8601ToUnix(filePath, key, isJst)
		if err != nil {
			// キーが存在しない場合はスキップ
			if strings.Contains(err.Error(), "キー") && strings.Contains(err.Error(), "が存在しません") {
				continue
			}
			// その他のエラーが発生しても処理を続行
			fmt.Printf("ファイル %s の処理中にエラーが発生しました: %v\n", filePath, err)
			continue
		}
		processedCount++
	}

	return processedCount, nil
}

// ConvertUnixToISO8601InAllFiles は指定されたディレクトリ内の全てのJSONファイルに指定したキーの値をUNIXタイムスタンプからISO-8601形式に変換します
func (s *JSONService) ConvertUnixToISO8601InAllFiles(dirPath string, key string, recursive bool) (int, error) {
	if key == "" {
		return 0, errors.New("キーは空にできません")
	}

	// ディレクトリ内のJSONファイルを検索
	jsonFiles, err := s.FileRepo.FindFilesByExt(dirPath, ".json")
	if err != nil {
		return 0, fmt.Errorf("JSONファイルの検索に失敗しました: %w", err)
	}

	// 再帰的に検索する場合
	if recursive {
		// ディレクトリ内のサブディレクトリを取得
		entries, err := s.FileRepo.ReadDir(dirPath)
		if err != nil {
			return 0, fmt.Errorf("ディレクトリの読み込みに失敗しました: %w", err)
		}

		// 各サブディレクトリを処理
		for _, entry := range entries {
			if entry.IsDir {
				// サブディレクトリ内のJSONファイルを処理
				subCount, err := s.ConvertUnixToISO8601InAllFiles(entry.Path, key, recursive)
				if err != nil {
					return 0, err
				}
				// 処理したファイル数を加算
				jsonFiles = append(jsonFiles, make([]string, subCount)...)
			}
		}
	}

	// 各JSONファイルを処理
	processedCount := 0
	for _, filePath := range jsonFiles {
		err := s.ConvertUnixToISO8601(filePath, key)
		if err != nil {
			// キーが存在しない場合はスキップ
			if strings.Contains(err.Error(), "キー") && strings.Contains(err.Error(), "が存在しません") {
				continue
			}
			// その他のエラーが発生しても処理を続行
			fmt.Printf("ファイル %s の処理中にエラーが発生しました: %v\n", filePath, err)
			continue
		}
		processedCount++
	}

	return processedCount, nil
}

// TimestampService はJSONファイルにタイムスタンプを追加するためのサービスです
type TimestampService struct {
	jsonService *JSONService
}

// NewTimestampService は新しいTimestampServiceインスタンスを作成します
func NewTimestampService(jsonService *JSONService) *TimestampService {
	return &TimestampService{
		jsonService: jsonService,
	}
}

// AddTimestamp は指定されたJSONファイルに現在の日時のタイムスタンプを追加します
func (s *TimestampService) AddTimestamp(filePath string, key string) error {
	if key == "" {
		return fmt.Errorf("キーは空にできません")
	}

	// 現在の日時を取得（UNIXタイムスタンプ）
	now := time.Now()
	timestamp := now.Unix()

	// JSONファイルにキーと値を追加
	return s.jsonService.AddKeyValue(filePath, key, timestamp)
}

// AddTimestampToAllFiles は指定されたディレクトリ内の全てのJSONファイルに現在の日時のタイムスタンプを追加します
func (s *TimestampService) AddTimestampToAllFiles(dirPath string, key string, recursive bool) (int, error) {
	if key == "" {
		return 0, errors.New("キーは空にできません")
	}

	// 現在の日時を取得（UNIXタイムスタンプ）
	now := time.Now()
	timestamp := now.Unix()

	// ディレクトリ内の全てのJSONファイルにキーと値を追加
	return s.jsonService.AddKeyValueToAllFiles(dirPath, key, timestamp, recursive)
}
