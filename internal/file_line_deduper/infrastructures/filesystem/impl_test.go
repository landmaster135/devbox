package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/landmaster135/devbox/internal/file_line_deduper/domain/models"
)

// テスト用の一時ディレクトリとファイルを作成する
func setupTestFiles(t *testing.T) (string, string, func()) {
	tempDir, err := os.MkdirTemp("", "file_repository_test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}

	testFilePath := filepath.Join(tempDir, "test_file.txt")
	testContent := []string{
		"これはテスト用のファイルです。",
		"2行目のテキスト",
		"3行目のテキスト",
	}

	file, err := os.Create(testFilePath)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}
	defer file.Close()

	for _, line := range testContent {
		if _, err := file.WriteString(line + "\n"); err != nil {
			t.Fatalf("テストファイルへの書き込みに失敗しました: %v", err)
		}
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, testFilePath, cleanup
}

func TestNewRepository(t *testing.T) {
	repo := NewRepository()
	if repo == nil {
		t.Error("NewRepository() がnilを返しました")
	}
}

func TestOSRepository_ReadFile(t *testing.T) {
	_, testFilePath, cleanup := setupTestFiles(t)
	defer cleanup()

	repo := NewRepository()

	content, err := repo.ReadFile(testFilePath)
	if err != nil {
		t.Errorf("ReadFile() エラー = %v", err)
		return
	}

	expectedLines := []string{
		"これはテスト用のファイルです。",
		"2行目のテキスト",
		"3行目のテキスト",
	}

	if len(content.Lines) != len(expectedLines) {
		t.Errorf("ReadFile() の行数が期待と異なります。got = %v, want = %v", len(content.Lines), len(expectedLines))
	}

	for i, line := range content.Lines {
		if line != expectedLines[i] {
			t.Errorf("ReadFile() の行 %d が期待と異なります。got = %v, want = %v", i, line, expectedLines[i])
		}
	}

	_, err = repo.ReadFile("存在しないファイル.txt")
	if err == nil {
		t.Error("存在しないファイルを読み込もうとしたときにエラーが発生しませんでした")
	}
}

func TestOSRepository_WriteFile(t *testing.T) {
	tempDir, _, cleanup := setupTestFiles(t)
	defer cleanup()

	outputPath := filepath.Join(tempDir, "output.txt")
	repo := NewRepository()

	lines := []string{
		"これは書き込みテスト用の1行目です。",
		"これは書き込みテスト用の2行目です。",
		"これは書き込みテスト用の3行目です。",
	}
	content := models.NewFileContent(lines)

	err := repo.WriteFile(outputPath, content)
	if err != nil {
		t.Errorf("WriteFile() エラー = %v", err)
		return
	}

	readContent, err := repo.ReadFile(outputPath)
	if err != nil {
		t.Errorf("検証のためのファイル読み込みに失敗しました: %v", err)
		return
	}

	if len(readContent.Lines) != len(lines) {
		t.Errorf("WriteFile() 後の行数が期待と異なります。got = %v, want = %v", len(readContent.Lines), len(lines))
	}

	for i, line := range readContent.Lines {
		if line != lines[i] {
			t.Errorf("WriteFile() 後の行 %d が期待と異なります。got = %v, want = %v", i, line, lines[i])
		}
	}

	invalidPath := filepath.Join(tempDir, "invalid_dir", "invalid.txt")
	err = repo.WriteFile(invalidPath, content)
	if err == nil {
		t.Error("無効なパスへの書き込みでエラーが発生しませんでした")
	}
}

func TestOSRepository_FileExists(t *testing.T) {
	_, testFilePath, cleanup := setupTestFiles(t)
	defer cleanup()

	repo := NewRepository()

	if !repo.FileExists(testFilePath) {
		t.Errorf("FileExists() が存在するファイル %s に対して false を返しました", testFilePath)
	}

	if repo.FileExists("存在しないファイル.txt") {
		t.Error("FileExists() が存在しないファイルに対して true を返しました")
	}
}

func TestOSRepository_ReadJSONFile_WithBOM(t *testing.T) {
	repo := NewRepository()
	testFilePath := "./test_data/org/sample_request_with_crlf_02.json"

	absPath, err := filepath.Abs(testFilePath)
	if err != nil {
		t.Fatalf("絶対パスへの変換に失敗しました: %v", err)
	}

	if !repo.FileExists(absPath) {
		t.Fatalf("テスト用のJSONファイルが存在しません: %s", absPath)
	}

	jsonData, err := repo.ReadJSONFile(absPath)
	if err != nil {
		t.Errorf("ReadJSONFile() エラー = %v", err)
		return
	}

	jsonMap, ok := jsonData.(map[string]interface{})
	if !ok {
		t.Errorf("JSONデータが期待した型ではありません。got = %T", jsonData)
		return
	}

	expectedName := "テストユーザー"
	if name, ok := jsonMap["name"].(string); !ok || name != expectedName {
		t.Errorf("name フィールドが期待と異なります。got = %v, want = %v", jsonMap["name"], expectedName)
	}

	expectedEmail := "test@example.com"
	if email, ok := jsonMap["email"].(string); !ok || email != expectedEmail {
		t.Errorf("email フィールドが期待と異なります。got = %v, want = %v", jsonMap["email"], expectedEmail)
	}

	expectedAge := float64(30)
	if age, ok := jsonMap["age"].(float64); !ok || age != expectedAge {
		t.Errorf("age フィールドが期待と異なります。got = %v, want = %v", jsonMap["age"], expectedAge)
	}

	interests, ok := jsonMap["interests"].([]interface{})
	if !ok {
		t.Errorf("interests フィールドが配列ではありません。got = %T", jsonMap["interests"])
		return
	}

	expectedInterests := []string{"プログラミング", "読書", "旅行"}
	if len(interests) != len(expectedInterests) {
		t.Errorf("interests の長さが期待と異なります。got = %v, want = %v", len(interests), len(expectedInterests))
	}

	for i, interest := range interests {
		if interestStr, ok := interest.(string); !ok || interestStr != expectedInterests[i] {
			t.Errorf("interests[%d] が期待と異なります。got = %v, want = %v", i, interest, expectedInterests[i])
		}
	}

	address, ok := jsonMap["address"].(map[string]interface{})
	if !ok {
		t.Errorf("address フィールドがオブジェクトではありません。got = %T", jsonMap["address"])
		return
	}

	expectedCountry := "日本"
	if country, ok := address["country"].(string); !ok || country != expectedCountry {
		t.Errorf("address.country が期待と異なります。got = %v, want = %v", address["country"], expectedCountry)
	}

	expectedCity := "東京"
	if city, ok := address["city"].(string); !ok || city != expectedCity {
		t.Errorf("address.city が期待と異なります。got = %v, want = %v", address["city"], expectedCity)
	}

	expectedPostalCode := "100-0001"
	if postalCode, ok := address["postalCode"].(string); !ok || postalCode != expectedPostalCode {
		t.Errorf("address.postalCode が期待と異なります。got = %v, want = %v", address["postalCode"], expectedPostalCode)
	}
}

func TestOSRepository_ReadJSONFile_FileNotFound(t *testing.T) {
	repo := NewRepository()

	_, err := repo.ReadJSONFile("存在しないファイル.json")
	if err == nil {
		t.Error("存在しないファイルを読み込もうとしたときにエラーが発生しませんでした")
	}
}

func TestOSRepository_ReadJSONFile_InvalidJSON(t *testing.T) {
	tempDir, _, cleanup := setupTestFiles(t)
	defer cleanup()

	invalidJSONPath := filepath.Join(tempDir, "invalid.json")
	invalidJSON := "{ これは無効なJSONです }"
	if err := os.WriteFile(invalidJSONPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("無効なJSONファイルの作成に失敗しました: %v", err)
	}

	repo := NewRepository()
	_, err := repo.ReadJSONFile(invalidJSONPath)
	if err == nil {
		t.Error("無効なJSONファイルを読み込もうとしたときにエラーが発生しませんでした")
	}
}
