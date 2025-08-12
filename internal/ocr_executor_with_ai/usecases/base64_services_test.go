package usecases

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBase64ExtractorService はBase64ExtractorServiceのテストクラス
type TestBase64ExtractorService struct {
	tempDir string
}

// setupTestEnvironment はテスト環境をセットアップする
func (ts *TestBase64ExtractorService) setupTestEnvironment(t *testing.T) {
	var err error
	ts.tempDir, err = os.MkdirTemp("", "base64_extractor_test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
}

// teardownTestEnvironment はテスト環境をクリーンアップする
func (ts *TestBase64ExtractorService) teardownTestEnvironment(t *testing.T) {
	if err := os.RemoveAll(ts.tempDir); err != nil {
		t.Errorf("一時ディレクトリの削除に失敗しました: %v", err)
	}
}

// createTestImageFile はテスト用の画像ファイルを作成する
func (ts *TestBase64ExtractorService) createTestImageFile(t *testing.T, filename string, content []byte) string {
	filePath := filepath.Join(ts.tempDir, filename)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}
	return filePath
}

// TestNewBase64ExtractorService_Normal はNewBase64ExtractorServiceの正常系テスト
func TestNewBase64ExtractorService_Normal(t *testing.T) {
	// Arrange
	targetPath := "/test/path"
	recursive := true

	// Act
	service := NewBase64ExtractorService(targetPath, recursive)

	// Assert
	if service == nil {
		t.Fatal("サービスがnilです")
	}
	if service.targetPath != targetPath {
		t.Errorf("targetPathが期待値と異なります。期待値: %s, 実際: %s", targetPath, service.targetPath)
	}
	if service.recursive != recursive {
		t.Errorf("recursiveが期待値と異なります。期待値: %t, 実際: %t", recursive, service.recursive)
	}
}

// TestBase64ExtractorService_ExtractFromPath_SingleFile_Normal は単一ファイルの正常系テスト
func TestBase64ExtractorService_ExtractFromPath_SingleFile_Normal(t *testing.T) {
	// Arrange
	ts := &TestBase64ExtractorService{}
	ts.setupTestEnvironment(t)
	defer ts.teardownTestEnvironment(t)

	testContent := []byte("test image content")
	testFile := ts.createTestImageFile(t, "test.jpg", testContent)
	expectedBase64 := base64.StdEncoding.EncodeToString(testContent)

	service := NewBase64ExtractorService(testFile, true)

	// Act
	result, err := service.ExtractFromPath()

	// Assert
	if err != nil {
		t.Fatalf("ExtractFromPathでエラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}
	if result.Total != 1 {
		t.Errorf("Totalが期待値と異なります。期待値: 1, 実際: %d", result.Total)
	}
	if len(result.Images) != 1 {
		t.Errorf("Imagesの長さが期待値と異なります。期待値: 1, 実際: %d", len(result.Images))
	}
	if result.Images[0].FilePath != testFile {
		t.Errorf("FilePathが期待値と異なります。期待値: %s, 実際: %s", testFile, result.Images[0].FilePath)
	}
	if result.Images[0].Base64 != expectedBase64 {
		t.Errorf("Base64が期待値と異なります。期待値: %s, 実際: %s", expectedBase64, result.Images[0].Base64)
	}
	if result.Images[0].Error != "" {
		t.Errorf("Errorが空でありません: %s", result.Images[0].Error)
	}
}

// TestBase64ExtractorService_ExtractFromPath_Directory_Normal はディレクトリの正常系テスト
func TestBase64ExtractorService_ExtractFromPath_Directory_Normal(t *testing.T) {
	// Arrange
	ts := &TestBase64ExtractorService{}
	ts.setupTestEnvironment(t)
	defer ts.teardownTestEnvironment(t)

	testContent1 := []byte("test image content 1")
	testContent2 := []byte("test image content 2")
	ts.createTestImageFile(t, "test1.jpg", testContent1)
	ts.createTestImageFile(t, "test2.png", testContent2)
	ts.createTestImageFile(t, "test.txt", []byte("not an image")) // 画像ファイルではない

	service := NewBase64ExtractorService(ts.tempDir, true)

	// Act
	result, err := service.ExtractFromPath()

	// Assert
	if err != nil {
		t.Fatalf("ExtractFromPathでエラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}
	if result.Total != 2 {
		t.Errorf("Totalが期待値と異なります。期待値: 2, 実際: %d", result.Total)
	}
	if len(result.Images) != 2 {
		t.Errorf("Imagesの長さが期待値と異なります。期待値: 2, 実際: %d", len(result.Images))
	}
}

// TestBase64ExtractorService_ExtractFromPath_NonRecursive_Normal は非再帰検索の正常系テスト
func TestBase64ExtractorService_ExtractFromPath_NonRecursive_Normal(t *testing.T) {
	// Arrange
	ts := &TestBase64ExtractorService{}
	ts.setupTestEnvironment(t)
	defer ts.teardownTestEnvironment(t)

	// ルートディレクトリにファイル作成
	ts.createTestImageFile(t, "root.jpg", []byte("root image"))

	// サブディレクトリ作成とファイル作成
	subDir := filepath.Join(ts.tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("サブディレクトリの作成に失敗しました: %v", err)
	}
	subFile := filepath.Join(subDir, "sub.jpg")
	if err := os.WriteFile(subFile, []byte("sub image"), 0644); err != nil {
		t.Fatalf("サブディレクトリのファイル作成に失敗しました: %v", err)
	}

	service := NewBase64ExtractorService(ts.tempDir, false) // 非再帰

	// Act
	result, err := service.ExtractFromPath()

	// Assert
	if err != nil {
		t.Fatalf("ExtractFromPathでエラーが発生しました: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("Totalが期待値と異なります。期待値: 1, 実際: %d", result.Total)
	}
}

// TestBase64ExtractorService_ExtractFromPath_NonExistentPath はパスが存在しない場合のテスト
func TestBase64ExtractorService_ExtractFromPath_NonExistentPath(t *testing.T) {
	// Arrange
	service := NewBase64ExtractorService("/non/existent/path", true)

	// Act
	result, err := service.ExtractFromPath()

	// Assert
	if err == nil {
		t.Fatal("エラーが発生するべきです")
	}
	if result != nil {
		t.Error("結果はnilであるべきです")
	}
	if !strings.Contains(err.Error(), "パスが存在しません") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %v", err)
	}
}

// TestBase64ExtractorService_isImageFile_Normal は画像ファイル判定の正常系テスト
func TestBase64ExtractorService_isImageFile_Normal(t *testing.T) {
	// Arrange
	service := NewBase64ExtractorService("", true)
	testCases := []struct {
		filename string
		expected bool
	}{
		{"test.jpg", true},
		{"test.jpeg", true},
		{"test.png", true},
		{"test.gif", true},
		{"test.bmp", true},
		{"test.webp", true},
		{"test.JPG", true}, // 大文字
		{"test.txt", false},
		{"test.pdf", false},
		{"test", false}, // 拡張子なし
	}

	for _, tc := range testCases {
		// Act
		result := service.isImageFile(tc.filename)

		// Assert
		if result != tc.expected {
			t.Errorf("isImageFile(%s)が期待値と異なります。期待値: %t, 実際: %t", tc.filename, tc.expected, result)
		}
	}
}

// TestExtractResult_FormatAsText_Normal はテキスト形式出力の正常系テスト
func TestExtractResult_FormatAsText_Normal(t *testing.T) {
	// Arrange
	result := &ExtractResult{
		Images: []ImageResult{
			{FilePath: "/test/image1.jpg", Base64: "dGVzdDE=", Error: ""},
			{FilePath: "/test/image2.png", Base64: "dGVzdDI=", Error: ""},
		},
		Total: 2,
	}

	// Act
	output := result.FormatAsText()

	// Assert
	if !strings.Contains(output, "=== Base64 Extraction Results ===") {
		t.Error("ヘッダーが含まれていません")
	}
	if !strings.Contains(output, "Total Images: 2") {
		t.Error("総数が含まれていません")
	}
	if !strings.Contains(output, "/test/image1.jpg") {
		t.Error("ファイルパス1が含まれていません")
	}
	if !strings.Contains(output, "/test/image2.png") {
		t.Error("ファイルパス2が含まれていません")
	}
	if !strings.Contains(output, "dGVzdDE=") {
		t.Error("Base64データ1が含まれていません")
	}
	if !strings.Contains(output, "dGVzdDI=") {
		t.Error("Base64データ2が含まれていません")
	}
}

// TestExtractResult_FormatAsText_Empty は空の結果のテキスト形式出力テスト
func TestExtractResult_FormatAsText_Empty(t *testing.T) {
	// Arrange
	result := &ExtractResult{
		Images: []ImageResult{},
		Total:  0,
	}

	// Act
	output := result.FormatAsText()

	// Assert
	expected := "画像ファイルが見つかりませんでした。"
	if output != expected {
		t.Errorf("出力が期待値と異なります。期待値: %s, 実際: %s", expected, output)
	}
}

// TestExtractResult_FormatAsJSON_Normal はJSON形式出力の正常系テスト
func TestExtractResult_FormatAsJSON_Normal(t *testing.T) {
	// Arrange
	result := &ExtractResult{
		Images: []ImageResult{
			{FilePath: "/test/image1.jpg", Base64: "dGVzdDE=", Error: ""},
		},
		Total: 1,
	}

	// Act
	output, err := result.FormatAsJSON()

	// Assert
	if err != nil {
		t.Fatalf("JSON変換でエラーが発生しました: %v", err)
	}
	if !strings.Contains(output, `"file_path": "/test/image1.jpg"`) {
		t.Error("ファイルパスが含まれていません")
	}
	if !strings.Contains(output, `"base64": "dGVzdDE="`) {
		t.Error("Base64データが含まれていません")
	}
	if !strings.Contains(output, `"total": 1`) {
		t.Error("総数が含まれていません")
	}
}

// TestExtractResult_FormatAsJSON_WithError はエラー付きのJSON形式出力テスト
func TestExtractResult_FormatAsJSON_WithError(t *testing.T) {
	// Arrange
	result := &ExtractResult{
		Images: []ImageResult{
			{FilePath: "/test/image1.jpg", Base64: "", Error: "ファイル読み込みエラー"},
		},
		Total: 1,
	}

	// Act
	output, err := result.FormatAsJSON()

	// Assert
	if err != nil {
		t.Fatalf("JSON変換でエラーが発生しました: %v", err)
	}
	if !strings.Contains(output, `"error": "ファイル読み込みエラー"`) {
		t.Error("エラーメッセージが含まれていません")
	}
}

// TestExtractResult_FormatAsJSON_Empty は空の結果のJSON形式出力テスト
func TestExtractResult_FormatAsJSON_Empty(t *testing.T) {
	// Arrange
	result := &ExtractResult{
		Images: []ImageResult{},
		Total:  0,
	}

	// Act
	output, err := result.FormatAsJSON()

	// Assert
	if err != nil {
		t.Fatalf("JSON変換でエラーが発生しました: %v", err)
	}
	if !strings.Contains(output, `"total": 0`) {
		t.Error("総数が含まれていません")
	}
	if !strings.Contains(output, `"images": []`) {
		t.Error("空の画像配列が含まれていません")
	}
}

// TestBase64ExtractorService_ExtractFromPath_FileReadError はファイル読み込みエラーのテスト
func TestBase64ExtractorService_ExtractFromPath_FileReadError(t *testing.T) {
	// Arrange
	ts := &TestBase64ExtractorService{}
	ts.setupTestEnvironment(t)
	defer ts.teardownTestEnvironment(t)

	// 読み込み権限のないファイルを作成
	testFile := filepath.Join(ts.tempDir, "no_read.jpg")
	if err := os.WriteFile(testFile, []byte("test content"), 0000); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	service := NewBase64ExtractorService(testFile, true)

	// Act
	result, err := service.ExtractFromPath()

	// Assert
	if err != nil {
		t.Fatalf("ExtractFromPathでエラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}
	if result.Total != 1 {
		t.Errorf("Totalが期待値と異なります。期待値: 1, 実際: %d", result.Total)
	}
	if len(result.Images) != 1 {
		t.Errorf("Imagesの長さが期待値と異なります。期待値: 1, 実際: %d", len(result.Images))
	}
	if result.Images[0].Error == "" {
		t.Error("エラーが設定されているべきです")
	}
	if result.Images[0].Base64 != "" {
		t.Error("Base64は空であるべきです")
	}
}

// TestBase64ExtractorService_ExtractFromPath_EmptyDirectory は空のディレクトリのテスト
func TestBase64ExtractorService_ExtractFromPath_EmptyDirectory(t *testing.T) {
	// Arrange
	ts := &TestBase64ExtractorService{}
	ts.setupTestEnvironment(t)
	defer ts.teardownTestEnvironment(t)

	service := NewBase64ExtractorService(ts.tempDir, true)

	// Act
	result, err := service.ExtractFromPath()

	// Assert
	if err != nil {
		t.Fatalf("ExtractFromPathでエラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}
	if result.Total != 0 {
		t.Errorf("Totalが期待値と異なります。期待値: 0, 実際: %d", result.Total)
	}
	if len(result.Images) != 0 {
		t.Errorf("Imagesの長さが期待値と異なります。期待値: 0, 実際: %d", len(result.Images))
	}
}

// TestBase64ExtractorService_ExtractFromPath_NonImageFile は画像ファイルでないファイルのテスト
func TestBase64ExtractorService_ExtractFromPath_NonImageFile(t *testing.T) {
	// Arrange
	ts := &TestBase64ExtractorService{}
	ts.setupTestEnvironment(t)
	defer ts.teardownTestEnvironment(t)

	testFile := ts.createTestImageFile(t, "test.txt", []byte("not an image"))
	service := NewBase64ExtractorService(testFile, true)

	// Act
	result, err := service.ExtractFromPath()

	// Assert
	if err != nil {
		t.Fatalf("ExtractFromPathでエラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}
	if result.Total != 0 {
		t.Errorf("Totalが期待値と異なります。期待値: 0, 実際: %d", result.Total)
	}
}

// TestBase64ExtractorService_ExtractFromPath_MixedFiles は画像ファイルと非画像ファイルの混合テスト
func TestBase64ExtractorService_ExtractFromPath_MixedFiles(t *testing.T) {
	// Arrange
	ts := &TestBase64ExtractorService{}
	ts.setupTestEnvironment(t)
	defer ts.teardownTestEnvironment(t)

	ts.createTestImageFile(t, "image1.jpg", []byte("image1"))
	ts.createTestImageFile(t, "image2.png", []byte("image2"))
	ts.createTestImageFile(t, "document.txt", []byte("text"))
	ts.createTestImageFile(t, "data.json", []byte("{}"))
	ts.createTestImageFile(t, "image3.gif", []byte("image3"))

	service := NewBase64ExtractorService(ts.tempDir, true)

	// Act
	result, err := service.ExtractFromPath()

	// Assert
	if err != nil {
		t.Fatalf("ExtractFromPathでエラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}
	if result.Total != 3 {
		t.Errorf("Totalが期待値と異なります。期待値: 3, 実際: %d", result.Total)
	}
	if len(result.Images) != 3 {
		t.Errorf("Imagesの長さが期待値と異なります。期待値: 3, 実際: %d", len(result.Images))
	}
}

// TestBase64ExtractorService_ExtractFromPath_DeepDirectory は深いディレクトリ構造のテスト
func TestBase64ExtractorService_ExtractFromPath_DeepDirectory(t *testing.T) {
	// Arrange
	ts := &TestBase64ExtractorService{}
	ts.setupTestEnvironment(t)
	defer ts.teardownTestEnvironment(t)

	// 深いディレクトリ構造を作成
	deepDir := filepath.Join(ts.tempDir, "level1", "level2", "level3")
	if err := os.MkdirAll(deepDir, 0755); err != nil {
		t.Fatalf("深いディレクトリの作成に失敗しました: %v", err)
	}

	// 各レベルにファイルを作成
	ts.createTestImageFile(t, "root.jpg", []byte("root"))

	level1File := filepath.Join(ts.tempDir, "level1", "level1.png")
	if err := os.WriteFile(level1File, []byte("level1"), 0644); err != nil {
		t.Fatalf("level1ファイルの作成に失敗しました: %v", err)
	}

	level3File := filepath.Join(deepDir, "level3.gif")
	if err := os.WriteFile(level3File, []byte("level3"), 0644); err != nil {
		t.Fatalf("level3ファイルの作成に失敗しました: %v", err)
	}

	service := NewBase64ExtractorService(ts.tempDir, true)

	// Act
	result, err := service.ExtractFromPath()

	// Assert
	if err != nil {
		t.Fatalf("ExtractFromPathでエラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}
	if result.Total != 3 {
		t.Errorf("Totalが期待値と異なります。期待値: 3, 実際: %d", result.Total)
	}
}

// TestBase64ExtractorService_ExtractFromPath_DeepDirectoryNonRecursive は深いディレクトリの非再帰テスト
func TestBase64ExtractorService_ExtractFromPath_DeepDirectoryNonRecursive(t *testing.T) {
	// Arrange
	ts := &TestBase64ExtractorService{}
	ts.setupTestEnvironment(t)
	defer ts.teardownTestEnvironment(t)

	// 深いディレクトリ構造を作成
	deepDir := filepath.Join(ts.tempDir, "level1", "level2")
	if err := os.MkdirAll(deepDir, 0755); err != nil {
		t.Fatalf("深いディレクトリの作成に失敗しました: %v", err)
	}

	// 各レベルにファイルを作成
	ts.createTestImageFile(t, "root.jpg", []byte("root"))

	level1File := filepath.Join(ts.tempDir, "level1", "level1.png")
	if err := os.WriteFile(level1File, []byte("level1"), 0644); err != nil {
		t.Fatalf("level1ファイルの作成に失敗しました: %v", err)
	}

	service := NewBase64ExtractorService(ts.tempDir, false) // 非再帰

	// Act
	result, err := service.ExtractFromPath()

	// Assert
	if err != nil {
		t.Fatalf("ExtractFromPathでエラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}
	if result.Total != 1 {
		t.Errorf("Totalが期待値と異なります。期待値: 1, 実際: %d", result.Total)
	}
}

// TestBase64ExtractorService_isImageFile_CaseInsensitive は大文字小文字を区別しない拡張子テスト
func TestBase64ExtractorService_isImageFile_CaseInsensitive(t *testing.T) {
	// Arrange
	service := NewBase64ExtractorService("", true)
	testCases := []struct {
		filename string
		expected bool
	}{
		{"test.JPG", true},
		{"test.JPEG", true},
		{"test.PNG", true},
		{"test.GIF", true},
		{"test.BMP", true},
		{"test.WEBP", true},
		{"test.Jpg", true},
		{"test.JpEg", true},
		{"TEST.PNG", true},
		{"test.TXT", false},
		{"test.PDF", false},
	}

	for _, tc := range testCases {
		// Act
		result := service.isImageFile(tc.filename)

		// Assert
		if result != tc.expected {
			t.Errorf("isImageFile(%s)が期待値と異なります。期待値: %t, 実際: %t", tc.filename, tc.expected, result)
		}
	}
}

// TestBase64ExtractorService_convertToBase64_Normal はconvertToBase64の正常系テスト
func TestBase64ExtractorService_convertToBase64_Normal(t *testing.T) {
	// Arrange
	ts := &TestBase64ExtractorService{}
	ts.setupTestEnvironment(t)
	defer ts.teardownTestEnvironment(t)

	testContent := []byte("test image content")
	testFile := ts.createTestImageFile(t, "test.jpg", testContent)
	expectedBase64 := base64.StdEncoding.EncodeToString(testContent)

	service := NewBase64ExtractorService("", true)

	// Act
	result := service.convertToBase64(testFile)

	// Assert
	if result.FilePath != testFile {
		t.Errorf("FilePathが期待値と異なります。期待値: %s, 実際: %s", testFile, result.FilePath)
	}
	if result.Base64 != expectedBase64 {
		t.Errorf("Base64が期待値と異なります。期待値: %s, 実際: %s", expectedBase64, result.Base64)
	}
	if result.Error != "" {
		t.Errorf("Errorが空でありません: %s", result.Error)
	}
}

// TestBase64ExtractorService_convertToBase64_FileNotFound はファイルが存在しない場合のテスト
func TestBase64ExtractorService_convertToBase64_FileNotFound(t *testing.T) {
	// Arrange
	service := NewBase64ExtractorService("", true)
	nonExistentFile := "/non/existent/file.jpg"

	// Act
	result := service.convertToBase64(nonExistentFile)

	// Assert
	if result.FilePath != nonExistentFile {
		t.Errorf("FilePathが期待値と異なります。期待値: %s, 実際: %s", nonExistentFile, result.FilePath)
	}
	if result.Base64 != "" {
		t.Error("Base64は空であるべきです")
	}
	if result.Error == "" {
		t.Error("エラーが設定されているべきです")
	}
	if !strings.Contains(result.Error, "ファイル読み込みエラー") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %s", result.Error)
	}
}

// TestExtractResult_FormatAsText_WithError はエラー付きのテキスト形式出力テスト
func TestExtractResult_FormatAsText_WithError(t *testing.T) {
	// Arrange
	result := &ExtractResult{
		Images: []ImageResult{
			{FilePath: "/test/image1.jpg", Base64: "dGVzdDE=", Error: ""},
			{FilePath: "/test/image2.jpg", Base64: "", Error: "ファイル読み込みエラー"},
		},
		Total: 2,
	}

	// Act
	output := result.FormatAsText()

	// Assert
	if !strings.Contains(output, "=== Base64 Extraction Results ===") {
		t.Error("ヘッダーが含まれていません")
	}
	if !strings.Contains(output, "Total Images: 2") {
		t.Error("総数が含まれていません")
	}
	if !strings.Contains(output, "/test/image1.jpg") {
		t.Error("ファイルパス1が含まれていません")
	}
	if !strings.Contains(output, "/test/image2.jpg") {
		t.Error("ファイルパス2が含まれていません")
	}
	if !strings.Contains(output, "dGVzdDE=") {
		t.Error("Base64データが含まれていません")
	}
	if !strings.Contains(output, "Error: ファイル読み込みエラー") {
		t.Error("エラーメッセージが含まれていません")
	}
}

// TestBase64ExtractorService_findImageFiles_WalkError はfilepathWalkエラーのテスト
func TestBase64ExtractorService_findImageFiles_WalkError(t *testing.T) {
	// Arrange
	service := NewBase64ExtractorService("", true)

	// Act
	result, err := service.findImageFiles("/non/existent/directory")

	// Assert
	if err == nil {
		t.Fatal("エラーが発生するべきです")
	}
	if result != nil {
		t.Error("結果はnilであるべきです")
	}
}
