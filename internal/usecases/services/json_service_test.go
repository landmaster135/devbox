package services

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/landmaster135/devbox/internal/domain/models"
)

// MockJSONFileRepository はJSONServiceテスト用のモックです
type MockJSONFileRepository struct {
	FileExistsFunc     func(path string) bool
	ReadJSONFileFunc   func(path string) (interface{}, error)
	WriteFileFunc      func(path string, content *models.FileContent) error
	FindFilesByExtFunc func(dirPath, ext string) ([]string, error)
	HasFilesWithExtFunc func(dirPath, ext string) (bool, error)
	GetDirectoryPathFunc func(path string) string
	CreateDirectoryFunc func(dirPath string) error
	ReadDirFunc        func(dirPath string) ([]*models.DirEntry, error)
}

func (m *MockJSONFileRepository) FileExists(path string) bool {
	return m.FileExistsFunc(path)
}

func (m *MockJSONFileRepository) ReadJSONFile(path string) (interface{}, error) {
	return m.ReadJSONFileFunc(path)
}

func (m *MockJSONFileRepository) WriteFile(path string, content *models.FileContent) error {
	return m.WriteFileFunc(path, content)
}

func (m *MockJSONFileRepository) ReadFile(path string) (*models.FileContent, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJSONFileRepository) FindFilesByExt(dirPath, ext string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (m *MockJSONFileRepository) HasFilesWithExt(dirPath, ext string) (bool, error) {
	return false, errors.New("not implemented")
}

func (m *MockJSONFileRepository) GetDirectoryPath(path string) string {
	return ""
}

func (m *MockJSONFileRepository) CreateDirectory(dirPath string) error {
	return errors.New("not implemented")
}

func (m *MockJSONFileRepository) ReadDir(dirPath string) ([]*models.DirEntry, error) {
	if m.ReadDirFunc != nil {
		return m.ReadDirFunc(dirPath)
	}
	return nil, errors.New("not implemented")
}

func TestJSONService_AddKeyValue_Normal(t *testing.T) {
	// テスト用のデータ
	testData := map[string]interface{}{
		"name": "テスト",
		"age":  30,
	}

	// モックの設定
	mockRepo := &MockJSONFileRepository{
		FileExistsFunc: func(path string) bool {
			return true
		},
		ReadJSONFileFunc: func(path string) (interface{}, error) {
			return testData, nil
		},
		WriteFileFunc: func(path string, content *models.FileContent) error {
			// 書き込まれたデータを検証
			jsonStr := ""
			for _, line := range content.Lines {
				jsonStr += line + "\n"
			}
			jsonStr = jsonStr[:len(jsonStr)-1] // 最後の改行を削除

			var writtenData map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &writtenData)
			if err != nil {
				t.Fatalf("書き込まれたデータのパースに失敗しました: %v", err)
			}

			// 新しいキーと値が追加されているか確認
			if writtenData["description"] != "テスト説明" {
				t.Errorf("期待された値 'テスト説明' ですが、'%v' が取得されました", writtenData["description"])
			}

			return nil
		},
	}

	// テスト対象のサービスを作成
	service := NewJSONService(mockRepo)

	// テスト実行
	err := service.AddKeyValue("test.json", "description", "テスト説明")
	if err != nil {
		t.Fatalf("AddKeyValue()でエラーが発生しました: %v", err)
	}
}

func TestJSONService_AddKeyValue_EmptyKey(t *testing.T) {
	// モックの設定
	mockRepo := &MockJSONFileRepository{}

	// テスト対象のサービスを作成
	service := NewJSONService(mockRepo)

	// テスト実行
	err := service.AddKeyValue("test.json", "", "テスト説明")
	if err == nil {
		t.Error("空のキーでエラーが発生することが期待されましたが、成功しました")
	}
}

func TestJSONService_GetValue_Normal(t *testing.T) {
	// テスト用のデータ
	testData := map[string]interface{}{
		"name": "テスト",
		"age":  30,
	}

	// モックの設定
	mockRepo := &MockJSONFileRepository{
		FileExistsFunc: func(path string) bool {
			return true
		},
		ReadJSONFileFunc: func(path string) (interface{}, error) {
			return testData, nil
		},
	}

	// テスト対象のサービスを作成
	service := NewJSONService(mockRepo)

	// テスト実行
	value, err := service.GetValue("test.json", "name")
	if err != nil {
		t.Fatalf("GetValue()でエラーが発生しました: %v", err)
	}

	// 値が正しく取得されたか確認
	if value != "テスト" {
		t.Errorf("期待された値 'テスト' ですが、'%v' が取得されました", value)
	}
}

func TestJSONService_GetValue_FileNotExists(t *testing.T) {
	// モックの設定
	mockRepo := &MockJSONFileRepository{
		FileExistsFunc: func(path string) bool {
			return false
		},
	}

	// テスト対象のサービスを作成
	service := NewJSONService(mockRepo)

	// テスト実行
	_, err := service.GetValue("test.json", "name")
	if err == nil {
		t.Error("ファイルが存在しない場合にエラーが発生することが期待されましたが、成功しました")
	}
}

func TestJSONService_GetValue_KeyNotExists(t *testing.T) {
	// テスト用のデータ
	testData := map[string]interface{}{
		"name": "テスト",
		"age":  30,
	}

	// モックの設定
	mockRepo := &MockJSONFileRepository{
		FileExistsFunc: func(path string) bool {
			return true
		},
		ReadJSONFileFunc: func(path string) (interface{}, error) {
			return testData, nil
		},
	}

	// テスト対象のサービスを作成
	service := NewJSONService(mockRepo)

	// テスト実行
	_, err := service.GetValue("test.json", "description")
	if err == nil {
		t.Error("存在しないキーの場合にエラーが発生することが期待されましたが、成功しました")
	}
}

func TestJSONService_AddKeyValue_Integer(t *testing.T) {
	// テスト用のデータ
	testData := map[string]interface{}{
		"name": "テスト",
	}

	// モックの設定
	mockRepo := &MockJSONFileRepository{
		FileExistsFunc: func(path string) bool {
			return true
		},
		ReadJSONFileFunc: func(path string) (interface{}, error) {
			return testData, nil
		},
		WriteFileFunc: func(path string, content *models.FileContent) error {
			// 書き込まれたデータを検証
			jsonStr := ""
			for _, line := range content.Lines {
				jsonStr += line + "\n"
			}
			jsonStr = jsonStr[:len(jsonStr)-1] // 最後の改行を削除

			var writtenData map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &writtenData)
			if err != nil {
				t.Fatalf("書き込まれたデータのパースに失敗しました: %v", err)
			}

			// 整数値が追加されているか確認
			count, exists := writtenData["count"].(float64)
			if !exists {
				t.Error("整数値が追加されていません")
			}

			// 整数値が正しいか確認
			if count != 42 {
				t.Errorf("期待された値 42 ですが、'%v' が取得されました", count)
			}

			return nil
		},
	}

	// テスト対象のサービスを作成
	service := NewJSONService(mockRepo)

	// テスト実行
	err := service.AddKeyValue("test.json", "count", int64(42))
	if err != nil {
		t.Fatalf("AddKeyValue()でエラーが発生しました: %v", err)
	}
}

func TestJSONService_AddKeyValue_Float(t *testing.T) {
	// テスト用のデータ
	testData := map[string]interface{}{
		"name": "テスト",
	}

	// モックの設定
	mockRepo := &MockJSONFileRepository{
		FileExistsFunc: func(path string) bool {
			return true
		},
		ReadJSONFileFunc: func(path string) (interface{}, error) {
			return testData, nil
		},
		WriteFileFunc: func(path string, content *models.FileContent) error {
			// 書き込まれたデータを検証
			jsonStr := ""
			for _, line := range content.Lines {
				jsonStr += line + "\n"
			}
			jsonStr = jsonStr[:len(jsonStr)-1] // 最後の改行を削除

			var writtenData map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &writtenData)
			if err != nil {
				t.Fatalf("書き込まれたデータのパースに失敗しました: %v", err)
			}

			// 浮動小数点数が追加されているか確認
			price, exists := writtenData["price"].(float64)
			if !exists {
				t.Error("浮動小数点数が追加されていません")
			}

			// 浮動小数点数が正しいか確認
			if price != 19.99 {
				t.Errorf("期待された値 19.99 ですが、'%v' が取得されました", price)
			}

			return nil
		},
	}

	// テスト対象のサービスを作成
	service := NewJSONService(mockRepo)

	// テスト実行
	err := service.AddKeyValue("test.json", "price", 19.99)
	if err != nil {
		t.Fatalf("AddKeyValue()でエラーが発生しました: %v", err)
	}
}

func TestTimestampService_AddTimestamp_Normal(t *testing.T) {
	// モックの設定
	mockRepo := &MockJSONFileRepository{
		FileExistsFunc: func(path string) bool {
			return true
		},
		ReadJSONFileFunc: func(path string) (interface{}, error) {
			return map[string]interface{}{
				"name": "テスト",
			}, nil
		},
		WriteFileFunc: func(path string, content *models.FileContent) error {
			// 書き込まれたデータを検証
			jsonStr := ""
			for _, line := range content.Lines {
				jsonStr += line + "\n"
			}
			jsonStr = jsonStr[:len(jsonStr)-1] // 最後の改行を削除

			var writtenData map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &writtenData)
			if err != nil {
				t.Fatalf("書き込まれたデータのパースに失敗しました: %v", err)
			}

			// タイムスタンプが追加されているか確認
			timestamp, exists := writtenData["timestamp"].(float64)
			if !exists {
				t.Error("タイムスタンプが追加されていません")
			}

			// タイムスタンプが現在時刻に近いか確認
			now := time.Now().Unix()
			if timestamp < float64(now-60) || timestamp > float64(now+60) {
				t.Errorf("タイムスタンプが現在時刻から離れすぎています: %v", timestamp)
			}

			return nil
		},
	}

	// テスト対象のサービスを作成
	jsonService := NewJSONService(mockRepo)
	timestampService := NewTimestampService(jsonService)

	// テスト実行
	err := timestampService.AddTimestamp("test.json", "timestamp")
	if err != nil {
		t.Fatalf("AddTimestamp()でエラーが発生しました: %v", err)
	}
}

func TestTimestampService_AddTimestamp_EmptyKey(t *testing.T) {
	// モックの設定
	mockRepo := &MockJSONFileRepository{}

	// テスト対象のサービスを作成
	jsonService := NewJSONService(mockRepo)
	timestampService := NewTimestampService(jsonService)

	// テスト実行
	err := timestampService.AddTimestamp("test.json", "")
	if err == nil {
		t.Error("空のキーでエラーが発生することが期待されましたが、成功しました")
	}
}
