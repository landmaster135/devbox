package repositories

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/landmaster135/devbox/internal/independencies/json_iso8601_converter/domain/models"
)

// テスト用の一時ディレクトリとファイルを作成する
func setupTestFiles(t *testing.T) (string, string, func()) {
	// 一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "file_repository_test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}

	// テスト用のファイルパス
	testFilePath := filepath.Join(tempDir, "test_file.txt")

	// テスト用のファイル内容
	testContent := []string{
		"これはテスト用のファイルです。",
		"2行目のテキスト",
		"3行目のテキスト",
	}

	// ファイルを作成して内容を書き込む
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

	// クリーンアップ関数を返す
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, testFilePath, cleanup
}

// TestNewFileRepository は NewFileRepository 関数をテストします
func TestNewFileRepository(t *testing.T) {
	repo := NewFileRepository()
	if repo == nil {
		t.Error("NewFileRepository() がnilを返しました")
	}
}

// TestFileRepositoryImpl_ReadFile は ReadFile メソッドをテストします
func TestFileRepositoryImpl_ReadFile(t *testing.T) {
	// テスト用のファイルをセットアップ
	_, testFilePath, cleanup := setupTestFiles(t)
	defer cleanup()

	// テスト対象のインスタンスを作成
	repo := NewFileRepository()

	// テスト実行
	content, err := repo.ReadFile(testFilePath)
	if err != nil {
		t.Errorf("ReadFile() エラー = %v", err)
		return
	}

	// 結果の検証
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

	// 存在しないファイルの場合のテスト
	_, err = repo.ReadFile("存在しないファイル.txt")
	if err == nil {
		t.Error("存在しないファイルを読み込もうとしたときにエラーが発生しませんでした")
	}
}

// TestFileRepositoryImpl_WriteFile は WriteFile メソッドをテストします
func TestFileRepositoryImpl_WriteFile(t *testing.T) {
	// テスト用のディレクトリをセットアップ
	tempDir, _, cleanup := setupTestFiles(t)
	defer cleanup()

	// 書き込み先のファイルパス
	outputPath := filepath.Join(tempDir, "output.txt")

	// テスト対象のインスタンスを作成
	repo := NewFileRepository()

	// テスト用のFileContentを作成
	lines := []string{
		"これは書き込みテスト用の1行目です。",
		"これは書き込みテスト用の2行目です。",
		"これは書き込みテスト用の3行目です。",
	}
	content := models.NewFileContent(lines)

	// テスト実行
	err := repo.WriteFile(outputPath, content)
	if err != nil {
		t.Errorf("WriteFile() エラー = %v", err)
		return
	}

	// 書き込まれたファイルを読み込んで検証
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

	// 無効なパスへの書き込みテスト
	invalidPath := filepath.Join(tempDir, "invalid_dir", "invalid.txt")
	err = repo.WriteFile(invalidPath, content)
	if err == nil {
		t.Error("無効なパスへの書き込みでエラーが発生しませんでした")
	}
}

// TestFileRepositoryImpl_FileExists は FileExists メソッドをテストします
func TestFileRepositoryImpl_FileExists(t *testing.T) {
	// テスト用のファイルをセットアップ
	_, testFilePath, cleanup := setupTestFiles(t)
	defer cleanup()

	// テスト対象のインスタンスを作成
	repo := NewFileRepository()

	// 存在するファイルのテスト
	if !repo.FileExists(testFilePath) {
		t.Errorf("FileExists() が存在するファイル %s に対して false を返しました", testFilePath)
	}

	// 存在しないファイルのテスト
	if repo.FileExists("存在しないファイル.txt") {
		t.Error("FileExists() が存在しないファイルに対して true を返しました")
	}
}

// TestFileRepositoryImpl_ReadJSONFile_WithBOM は BOMを含むJSONファイルを読み込むテストです
func TestFileRepositoryImpl_ReadJSONFile_WithBOM(t *testing.T) {
	// テスト対象のインスタンスを作成
	repo := NewFileRepository()

	// BOMを含むJSONファイルのパス
	testFilePath := "./test_data/org/sample_request_with_crlf_02.json"

	// 絶対パスに変換（テスト実行環境によって相対パスが変わる可能性があるため）
	absPath, err := filepath.Abs(testFilePath)
	if err != nil {
		t.Fatalf("絶対パスへの変換に失敗しました: %v", err)
	}

	// ファイルが存在することを確認
	if !repo.FileExists(absPath) {
		t.Fatalf("テスト用のJSONファイルが存在しません: %s", absPath)
	}

	// テスト実行
	jsonData, err := repo.ReadJSONFile(absPath)
	if err != nil {
		t.Errorf("ReadJSONFile() エラー = %v", err)
		return
	}

	// 結果の検証
	// 型アサーションでマップに変換
	jsonMap, ok := jsonData.(map[string]interface{})
	if !ok {
		t.Errorf("JSONデータが期待した型ではありません。got = %T", jsonData)
		return
	}

	// 各フィールドの検証
	expectedName := "テストユーザー"
	if name, ok := jsonMap["name"].(string); !ok || name != expectedName {
		t.Errorf("name フィールドが期待と異なります。got = %v, want = %v", jsonMap["name"], expectedName)
	}

	expectedEmail := "test@example.com"
	if email, ok := jsonMap["email"].(string); !ok || email != expectedEmail {
		t.Errorf("email フィールドが期待と異なります。got = %v, want = %v", jsonMap["email"], expectedEmail)
	}

	expectedAge := float64(30) // JSONのnumberはfloat64として解析される
	if age, ok := jsonMap["age"].(float64); !ok || age != expectedAge {
		t.Errorf("age フィールドが期待と異なります。got = %v, want = %v", jsonMap["age"], expectedAge)
	}

	// interestsフィールド（配列）の検証
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

	// addressフィールド（オブジェクト）の検証
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

// TestFileRepositoryImpl_ReadJSONFile_FileNotFound はファイルが存在しない場合のテストです
func TestFileRepositoryImpl_ReadJSONFile_FileNotFound(t *testing.T) {
	// テスト対象のインスタンスを作成
	repo := NewFileRepository()

	// 存在しないファイルパス
	nonExistentPath := "存在しないファイル.json"

	// テスト実行
	_, err := repo.ReadJSONFile(nonExistentPath)

	// エラーが発生することを確認
	if err == nil {
		t.Error("存在しないファイルを読み込もうとしたときにエラーが発生しませんでした")
	}
}

// TestFileRepositoryImpl_ReadJSONFile_InvalidJSON は無効なJSONファイルを読み込む場合のテストです
func TestFileRepositoryImpl_ReadJSONFile_InvalidJSON(t *testing.T) {
	// テスト用の一時ディレクトリをセットアップ
	tempDir, _, cleanup := setupTestFiles(t)
	defer cleanup()

	// 無効なJSONファイルのパス
	invalidJSONPath := filepath.Join(tempDir, "invalid.json")

	// 無効なJSONファイルを作成
	invalidJSON := "{ これは無効なJSONです }"
	if err := os.WriteFile(invalidJSONPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("無効なJSONファイルの作成に失敗しました: %v", err)
	}

	// テスト対象のインスタンスを作成
	repo := NewFileRepository()

	// テスト実行
	_, err := repo.ReadJSONFile(invalidJSONPath)

	// エラーが発生することを確認
	if err == nil {
		t.Error("無効なJSONファイルを読み込もうとしたときにエラーが発生しませんでした")
	}
}
