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
