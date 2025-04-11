package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/landmaster135/devbox/internal/domain/models"
	"github.com/landmaster135/devbox/internal/domain/repositories"
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
