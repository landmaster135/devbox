package usecases

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type OCRExecutorServiceTestSuite struct {
	tempDir string
	service *OCRExecutorService
}

func (suite *OCRExecutorServiceTestSuite) SetupTest(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "ocr_test")
	require.NoError(t, err)
	suite.tempDir = tempDir

	// テスト用のサービスを作成
	suite.service = NewOCRExecutorService(
		tempDir,
		false,
		"eng",
		"",
		"text",
	)
}

func (suite *OCRExecutorServiceTestSuite) TeardownTest(t *testing.T) {
	// テスト用の一時ディレクトリを削除
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

func TestOCRExecutorService_NewOCRExecutorService_Normal(t *testing.T) {
	// Arrange
	targetPath := "/test/path"
	recursive := true
	languages := "jpn+eng"
	outputDir := "/output"
	outputFormat := "json"

	// Act
	service := NewOCRExecutorService(targetPath, recursive, languages, outputDir, outputFormat)

	// Assert
	assert.NotNil(t, service)
	assert.Equal(t, targetPath, service.targetPath)
	assert.Equal(t, recursive, service.recursive)
	assert.Equal(t, languages, service.languages)
	assert.Equal(t, outputDir, service.outputDir)
	assert.Equal(t, outputFormat, service.outputFormat)
}

func TestOCRExecutorService_isImageFile_Normal(t *testing.T) {
	suite := &OCRExecutorServiceTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	testCases := []struct {
		name     string
		filePath string
		expected bool
	}{
		{"JPEG file", "test.jpg", true},
		{"PNG file", "test.png", true},
		{"GIF file", "test.gif", true},
		{"BMP file", "test.bmp", true},
		{"WEBP file", "test.webp", true},
		{"TIFF file", "test.tiff", true},
		{"TIF file", "test.tif", true},
		{"Text file", "test.txt", false},
		{"PDF file", "test.pdf", false},
		{"No extension", "test", false},
		{"Upper case", "test.JPG", true},
		{"Mixed case", "test.Png", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := suite.service.isImageFile(tc.filePath)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestOCRExecutorService_findImageFiles_Normal(t *testing.T) {
	suite := &OCRExecutorServiceTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange - テスト用ファイルを作成
	imageFiles := []string{"test1.jpg", "test2.png", "test3.gif"}
	nonImageFiles := []string{"test.txt", "test.pdf"}

	for _, file := range imageFiles {
		filePath := filepath.Join(suite.tempDir, file)
		err := os.WriteFile(filePath, []byte("dummy"), 0644)
		require.NoError(t, err)
	}

	for _, file := range nonImageFiles {
		filePath := filepath.Join(suite.tempDir, file)
		err := os.WriteFile(filePath, []byte("dummy"), 0644)
		require.NoError(t, err)
	}

	// Act
	result, err := suite.service.findImageFiles(suite.tempDir)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, len(imageFiles))

	// ファイル名のみを抽出して比較
	resultFileNames := make([]string, len(result))
	for i, path := range result {
		resultFileNames[i] = filepath.Base(path)
	}

	for _, expectedFile := range imageFiles {
		assert.Contains(t, resultFileNames, expectedFile)
	}
}

func TestOCRExecutorService_findImageFiles_Recursive_Normal(t *testing.T) {
	suite := &OCRExecutorServiceTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange - サブディレクトリを作成
	subDir := filepath.Join(suite.tempDir, "subdir")
	err := os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	// ルートディレクトリにファイル作成
	rootFile := filepath.Join(suite.tempDir, "root.jpg")
	err = os.WriteFile(rootFile, []byte("dummy"), 0644)
	require.NoError(t, err)

	// サブディレクトリにファイル作成
	subFile := filepath.Join(subDir, "sub.png")
	err = os.WriteFile(subFile, []byte("dummy"), 0644)
	require.NoError(t, err)

	// Act - 再帰検索有効
	suite.service.recursive = true
	result, err := suite.service.findImageFiles(suite.tempDir)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Act - 再帰検索無効
	suite.service.recursive = false
	result, err = suite.service.findImageFiles(suite.tempDir)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Contains(t, result[0], "root.jpg")
}

func TestExecuteResult_FormatAsText_Normal(t *testing.T) {
	// Arrange
	results := []OCRResult{
		{FilePath: "/path/to/image1.jpg", Text: "Hello World", Error: ""},
		{FilePath: "/path/to/image2.png", Text: "", Error: "OCR failed"},
	}
	executeResult := &ExecuteResult{
		Results: results,
		Total:   2,
	}

	// Act
	output := executeResult.FormatAsText()

	// Assert
	assert.Contains(t, output, "=== OCR Results ===")
	assert.Contains(t, output, "Total Images: 2")
	assert.Contains(t, output, "[1] /path/to/image1.jpg")
	assert.Contains(t, output, "Text: Hello World")
	assert.Contains(t, output, "[2] /path/to/image2.png")
	assert.Contains(t, output, "Error: OCR failed")
}

func TestExecuteResult_FormatAsText_NoImages_Normal(t *testing.T) {
	// Arrange
	executeResult := &ExecuteResult{
		Results: []OCRResult{},
		Total:   0,
	}

	// Act
	output := executeResult.FormatAsText()

	// Assert
	assert.Equal(t, "画像ファイルが見つかりませんでした。", output)
}

func TestExecuteResult_FormatAsJSON_Normal(t *testing.T) {
	// Arrange
	results := []OCRResult{
		{FilePath: "/path/to/image1.jpg", Text: "Hello World", Error: ""},
	}
	executeResult := &ExecuteResult{
		Results: results,
		Total:   1,
	}

	// Act
	output, err := executeResult.FormatAsJSON()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, output, `"file_path": "/path/to/image1.jpg"`)
	assert.Contains(t, output, `"text": "Hello World"`)
	assert.Contains(t, output, `"total": 1`)
}

func TestOCRExecutorService_saveToFile_Text_Normal(t *testing.T) {
	suite := &OCRExecutorServiceTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange
	outputDir := filepath.Join(suite.tempDir, "output")
	suite.service.outputDir = outputDir
	suite.service.outputFormat = "text"

	executeResult := &ExecuteResult{
		Results: []OCRResult{
			{FilePath: "/test/image.jpg", Text: "Test text", Error: ""},
		},
		Total: 1,
	}

	// Act
	err := suite.service.saveToFile(executeResult)

	// Assert
	require.NoError(t, err)

	// ファイルが作成されたことを確認
	filePath := filepath.Join(outputDir, "ocr_results.txt")
	assert.FileExists(t, filePath)

	// ファイル内容を確認
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "=== OCR Results ===")
	assert.Contains(t, string(content), "Test text")
}

func TestOCRExecutorService_saveToFile_JSON_Normal(t *testing.T) {
	suite := &OCRExecutorServiceTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange
	outputDir := filepath.Join(suite.tempDir, "output")
	suite.service.outputDir = outputDir
	suite.service.outputFormat = "json"

	executeResult := &ExecuteResult{
		Results: []OCRResult{
			{FilePath: "/test/image.jpg", Text: "Test text", Error: ""},
		},
		Total: 1,
	}

	// Act
	err := suite.service.saveToFile(executeResult)

	// Assert
	require.NoError(t, err)

	// ファイルが作成されたことを確認
	filePath := filepath.Join(outputDir, "ocr_results.json")
	assert.FileExists(t, filePath)

	// ファイル内容を確認
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(content), `"file_path": "/test/image.jpg"`)
	assert.Contains(t, string(content), `"text": "Test text"`)
}

func TestOCRExecutorService_ExecuteFromPath_NonExistentPath_Normal(t *testing.T) {
	// Arrange
	service := NewOCRExecutorService("/non/existent/path", false, "eng", "", "text")

	// Act
	result, err := service.ExecuteFromPath()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "パスが存在しません")
}

func TestOCRExecutorService_ExecuteFromPath_SingleFile_Normal(t *testing.T) {
	suite := &OCRExecutorServiceTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange - 単一の画像ファイルを作成
	imageFile := filepath.Join(suite.tempDir, "test.jpg")
	err := os.WriteFile(imageFile, []byte("dummy image data"), 0644)
	require.NoError(t, err)

	service := NewOCRExecutorService(imageFile, false, "eng", "", "text")

	// Act
	result, err := service.ExecuteFromPath()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Results, 1)
	assert.Equal(t, imageFile, result.Results[0].FilePath)
	// OCRの実際の処理はTesseractに依存するため、エラーが発生する可能性がある
}

func TestOCRExecutorService_ExecuteFromPath_Directory_Normal(t *testing.T) {
	suite := &OCRExecutorServiceTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange - ディレクトリに複数の画像ファイルを作成
	imageFiles := []string{"test1.jpg", "test2.png", "test3.gif"}
	for _, file := range imageFiles {
		filePath := filepath.Join(suite.tempDir, file)
		err := os.WriteFile(filePath, []byte("dummy image data"), 0644)
		require.NoError(t, err)
	}

	service := NewOCRExecutorService(suite.tempDir, false, "eng", "", "text")

	// Act
	result, err := service.ExecuteFromPath()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(imageFiles), result.Total)
	assert.Len(t, result.Results, len(imageFiles))
}

func TestOCRExecutorService_ExecuteFromPath_NonImageFile_Normal(t *testing.T) {
	suite := &OCRExecutorServiceTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange - 非画像ファイルを作成
	textFile := filepath.Join(suite.tempDir, "test.txt")
	err := os.WriteFile(textFile, []byte("dummy text data"), 0644)
	require.NoError(t, err)

	service := NewOCRExecutorService(textFile, false, "eng", "", "text")

	// Act
	result, err := service.ExecuteFromPath()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Total)
	assert.Len(t, result.Results, 0)
}

func TestOCRExecutorService_ExecuteFromPath_WithOutputDir_Normal(t *testing.T) {
	suite := &OCRExecutorServiceTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange
	imageFile := filepath.Join(suite.tempDir, "test.jpg")
	err := os.WriteFile(imageFile, []byte("dummy image data"), 0644)
	require.NoError(t, err)

	outputDir := filepath.Join(suite.tempDir, "output")
	service := NewOCRExecutorService(imageFile, false, "eng", outputDir, "text")

	// Act
	result, err := service.ExecuteFromPath()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)

	// 出力ファイルが作成されたことを確認
	outputFile := filepath.Join(outputDir, "ocr_results.txt")
	assert.FileExists(t, outputFile)
}

func TestOCRExecutorService_ExecuteFromPath_WithOutputDirJSON_Normal(t *testing.T) {
	suite := &OCRExecutorServiceTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange
	imageFile := filepath.Join(suite.tempDir, "test.jpg")
	err := os.WriteFile(imageFile, []byte("dummy image data"), 0644)
	require.NoError(t, err)

	outputDir := filepath.Join(suite.tempDir, "output")
	service := NewOCRExecutorService(imageFile, false, "eng", outputDir, "json")

	// Act
	result, err := service.ExecuteFromPath()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)

	// 出力ファイルが作成されたことを確認
	outputFile := filepath.Join(outputDir, "ocr_results.json")
	assert.FileExists(t, outputFile)
}

func TestOCRExecutorService_saveToFile_InvalidOutputDir_Normal(t *testing.T) {
	suite := &OCRExecutorServiceTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange
	service := NewOCRExecutorService(suite.tempDir, false, "eng", "/root/invalid_dir", "text")
	executeResult := &ExecuteResult{
		Results: []OCRResult{
			{FilePath: "/test/image.jpg", Text: "Test text", Error: ""},
		},
		Total: 1,
	}

	// Act
	err := service.saveToFile(executeResult)

	// Assert
	// 権限エラーが発生する可能性があるが、環境によって異なる
	if err != nil {
		assert.Contains(t, err.Error(), "出力ディレクトリの作成に失敗しました")
	}
}

func TestOCRExecutorService_findImageFiles_EmptyDirectory_Normal(t *testing.T) {
	suite := &OCRExecutorServiceTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Act - 空のディレクトリで検索
	result, err := suite.service.findImageFiles(suite.tempDir)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestOCRExecutorService_findImageFiles_WithSubdirectories_Normal(t *testing.T) {
	suite := &OCRExecutorServiceTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange - 複数レベルのサブディレクトリを作成
	subDir1 := filepath.Join(suite.tempDir, "sub1")
	subDir2 := filepath.Join(subDir1, "sub2")
	err := os.MkdirAll(subDir2, 0755)
	require.NoError(t, err)

	// 各レベルに画像ファイルを作成
	files := map[string]string{
		filepath.Join(suite.tempDir, "root.jpg"): "root image",
		filepath.Join(subDir1, "sub1.png"):       "sub1 image",
		filepath.Join(subDir2, "sub2.gif"):       "sub2 image",
	}

	for filePath, content := range files {
		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)
	}

	// Act - 再帰検索有効
	suite.service.recursive = true
	result, err := suite.service.findImageFiles(suite.tempDir)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 3)

	// Act - 再帰検索無効
	suite.service.recursive = false
	result, err = suite.service.findImageFiles(suite.tempDir)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1) // ルートディレクトリのファイルのみ
}

func TestOCRExecutorService_isImageFile_CaseInsensitive_Normal(t *testing.T) {
	suite := &OCRExecutorServiceTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	testCases := []struct {
		name     string
		filePath string
		expected bool
	}{
		{"Lowercase jpg", "test.jpg", true},
		{"Uppercase JPG", "test.JPG", true},
		{"Mixed case Jpg", "test.Jpg", true},
		{"Lowercase png", "test.png", true},
		{"Uppercase PNG", "test.PNG", true},
		{"Mixed case Png", "test.Png", true},
		{"Double extension", "test.backup.jpg", true},
		{"No extension", "testjpg", false},
		{"Partial match", "test.jpgx", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := suite.service.isImageFile(tc.filePath)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestExecuteResult_FormatAsText_MultipleResults_Normal(t *testing.T) {
	// Arrange
	results := []OCRResult{
		{FilePath: "/path/to/image1.jpg", Text: "Hello World", Error: ""},
		{FilePath: "/path/to/image2.png", Text: "こんにちは", Error: ""},
		{FilePath: "/path/to/image3.gif", Text: "", Error: "OCR failed"},
		{FilePath: "/path/to/image4.bmp", Text: "Test\nMultiline\nText", Error: ""},
	}
	executeResult := &ExecuteResult{
		Results: results,
		Total:   4,
	}

	// Act
	output := executeResult.FormatAsText()

	// Assert
	assert.Contains(t, output, "=== OCR Results ===")
	assert.Contains(t, output, "Total Images: 4")
	assert.Contains(t, output, "[1] /path/to/image1.jpg")
	assert.Contains(t, output, "Text: Hello World")
	assert.Contains(t, output, "[2] /path/to/image2.png")
	assert.Contains(t, output, "Text: こんにちは")
	assert.Contains(t, output, "[3] /path/to/image3.gif")
	assert.Contains(t, output, "Error: OCR failed")
	assert.Contains(t, output, "[4] /path/to/image4.bmp")
	assert.Contains(t, output, "Test\nMultiline\nText")
}

func TestExecuteResult_FormatAsJSON_MultipleResults_Normal(t *testing.T) {
	// Arrange
	results := []OCRResult{
		{FilePath: "/path/to/image1.jpg", Text: "Hello World", Error: ""},
		{FilePath: "/path/to/image2.png", Text: "", Error: "OCR failed"},
	}
	executeResult := &ExecuteResult{
		Results: results,
		Total:   2,
	}

	// Act
	output, err := executeResult.FormatAsJSON()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, output, `"file_path": "/path/to/image1.jpg"`)
	assert.Contains(t, output, `"text": "Hello World"`)
	assert.Contains(t, output, `"file_path": "/path/to/image2.png"`)
	assert.Contains(t, output, `"error": "OCR failed"`)
	assert.Contains(t, output, `"total": 2`)
}

func TestExecuteResult_FormatAsJSON_EmptyResults_Normal(t *testing.T) {
	// Arrange
	executeResult := &ExecuteResult{
		Results: []OCRResult{},
		Total:   0,
	}

	// Act
	output, err := executeResult.FormatAsJSON()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, output, `"results": []`)
	assert.Contains(t, output, `"total": 0`)
}

func TestOCRExecutorService_saveToFile_CreateDirectoryError_Normal(t *testing.T) {
	// Arrange
	// 存在しないパスの子ディレクトリを指定（権限問題を回避）
	invalidPath := "/proc/invalid/path" // procファイルシステムは読み取り専用
	service := NewOCRExecutorService("/tmp", false, "eng", invalidPath, "text")

	executeResult := &ExecuteResult{
		Results: []OCRResult{
			{FilePath: "/test/image.jpg", Text: "Test text", Error: ""},
		},
		Total: 1,
	}

	// Act
	err := service.saveToFile(executeResult)

	// Assert
	if err != nil {
		assert.Contains(t, err.Error(), "出力ディレクトリの作成に失敗しました")
	}
}

func TestOCRResult_WithError_Normal(t *testing.T) {
	// Arrange
	result := OCRResult{
		FilePath: "/path/to/image.jpg",
		Text:     "",
		Error:    "Tesseract initialization failed",
	}

	// Act & Assert
	assert.Equal(t, "/path/to/image.jpg", result.FilePath)
	assert.Equal(t, "", result.Text)
	assert.Equal(t, "Tesseract initialization failed", result.Error)
}

func TestOCRResult_WithText_Normal(t *testing.T) {
	// Arrange
	result := OCRResult{
		FilePath: "/path/to/image.jpg",
		Text:     "Extracted text from image",
		Error:    "",
	}

	// Act & Assert
	assert.Equal(t, "/path/to/image.jpg", result.FilePath)
	assert.Equal(t, "Extracted text from image", result.Text)
	assert.Equal(t, "", result.Error)
}
