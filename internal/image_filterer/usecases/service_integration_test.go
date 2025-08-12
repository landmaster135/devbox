package usecases

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestImageFilterServiceIntegration は統合テストクラス
type TestImageFilterServiceIntegration struct {}

// TestProcessImages_ErrorHandling はエラーハンドリングをテストします
func TestProcessImages_ErrorHandling(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// 無効な画像ファイル（テキストファイル）を作成
	invalidFile := filepath.Join(inputDir, "invalid.jpg")
	require.NoError(t, os.WriteFile(invalidFile, []byte("not an image"), 0644))

	// 正常な画像ファイルも作成
	validFile := filepath.Join(inputDir, "valid.jpg")
	require.NoError(t, createTestImage(validFile, 20, 20, "jpg"))

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        0,
		Y1:        0,
		X2:        0,
		Y2:        0,
		Suffix:    "filtered",
		Move:      false,
		Recursive: false,
		Workers:   1,
	}

	filterConfig := &FilterConfig{
		Mode:    BlurMode,
		Radius:  1.0,
		RWeight: 0.299,
		GWeight: 0.587,
		BWeight: 0.114,
	}

	service, err := NewImageFilterService(procConfig, filterConfig)
	require.NoError(t, err)

	// パイプラインを実行（エラーが発生するが処理は継続される）
	result, err := service.ProcessImages()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// 1つは成功、1つはエラー
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.ErrorCount)

	// 正常なファイルの出力は作成されている
	validOutput := filepath.Join(outputDir, "valid_filtered.jpg")
	_, err = os.Stat(validOutput)
	assert.NoError(t, err)

	// 無効なファイルの出力は作成されていない
	invalidOutput := filepath.Join(outputDir, "invalid_filtered.jpg")
	_, err = os.Stat(invalidOutput)
	assert.True(t, os.IsNotExist(err))
}

// TestProcessImages_WithMove は移動オプション付きの完全パイプラインをテストします
func TestProcessImages_WithMove(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// テスト用画像ファイルを作成
	testFiles := []string{"move1.jpg", "move2.png"}
	for _, file := range testFiles {
		path := filepath.Join(inputDir, file)
		require.NoError(t, createTestImage(path, 25, 25, filepath.Ext(file)[1:]))
	}

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        0,
		Y1:        0,
		X2:        0,
		Y2:        0,
		Suffix:    "moved",
		Move:      true, // 移動オプション有効
		Recursive: false,
		Workers:   2,
	}

	filterConfig := &FilterConfig{
		Mode:    GrayscaleMode,
		Radius:  0.0,
		RWeight: 0.299,
		GWeight: 0.587,
		BWeight: 0.114,
	}

	service, err := NewImageFilterService(procConfig, filterConfig)
	require.NoError(t, err)

	// 完全なパイプラインを実行
	result, err := service.ProcessImages()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 0, result.ErrorCount)

	// 出力ファイルが作成されていることを確認
	expectedOutputs := []string{"move1_moved.jpg", "move2_moved.png"}
	for _, expectedOutput := range expectedOutputs {
		outputPath := filepath.Join(outputDir, expectedOutput)
		_, err := os.Stat(outputPath)
		assert.NoError(t, err, "Output file should exist: %s", expectedOutput)
	}

	// 元ファイルがアーカイブに移動されていることを確認
	for _, file := range testFiles {
		// 元の場所にはファイルが存在しない
		originalPath := filepath.Join(inputDir, file)
		_, err := os.Stat(originalPath)
		assert.True(t, os.IsNotExist(err), "Original file should be moved: %s", file)

		// アーカイブディレクトリにファイルが存在する
		archivedPath := filepath.Join(archiveDir, file)
		_, err = os.Stat(archivedPath)
		assert.NoError(t, err, "Archived file should exist: %s", file)
	}
}

// TestProcessImages_RecursiveMode は再帰モードでの完全パイプラインをテストします
func TestProcessImages_RecursiveMode(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// サブディレクトリを作成
	subDir1 := filepath.Join(inputDir, "sub1")
	subDir2 := filepath.Join(inputDir, "sub2")
	require.NoError(t, os.MkdirAll(subDir1, 0755))
	require.NoError(t, os.MkdirAll(subDir2, 0755))

	// 各ディレクトリにテスト用画像ファイルを作成
	testFiles := map[string]string{
		"root.jpg":        inputDir,
		"sub1/file1.png":  inputDir,
		"sub2/file2.webp": inputDir,
	}

	for file, baseDir := range testFiles {
		path := filepath.Join(baseDir, file)
		require.NoError(t, createTestImage(path, 20, 20, filepath.Ext(file)[1:]))
	}

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        0,
		Y1:        0,
		X2:        0,
		Y2:        0,
		Suffix:    "recursive",
		Move:      false,
		Recursive: true, // 再帰モード
		Workers:   2,
	}

	filterConfig := &FilterConfig{
		Mode:    BlurMode,
		Radius:  1.5,
		RWeight: 0.299,
		GWeight: 0.587,
		BWeight: 0.114,
	}

	service, err := NewImageFilterService(procConfig, filterConfig)
	require.NoError(t, err)

	// 完全なパイプラインを実行
	result, err := service.ProcessImages()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.ErrorCount)

	// 出力ファイルが作成されていることを確認
	expectedOutputs := []string{"root_recursive.jpg", "file1_recursive.png", "file2_recursive.webp"}
	for _, expectedOutput := range expectedOutputs {
		outputPath := filepath.Join(outputDir, expectedOutput)
		_, err := os.Stat(outputPath)
		assert.NoError(t, err, "Output file should exist: %s", expectedOutput)
	}
}

// TestProcessImages_LargeCoordinates は大きな座標指定でのテストを行います
func TestProcessImages_LargeCoordinates(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// 大きなテスト用画像ファイルを作成
	inputFile := filepath.Join(inputDir, "large.jpg")
	require.NoError(t, createTestImage(inputFile, 100, 100, "jpg"))

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        20,
		Y1:        20,
		X2:        80,
		Y2:        80,
		Suffix:    "cropped",
		Move:      false,
		Recursive: false,
		Workers:   1,
	}

	filterConfig := &FilterConfig{
		Mode:    GrayscaleMode,
		Radius:  0.0,
		RWeight: 0.299,
		GWeight: 0.587,
		BWeight: 0.114,
	}

	service, err := NewImageFilterService(procConfig, filterConfig)
	require.NoError(t, err)

	// パイプラインを実行
	result, err := service.ProcessImages()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 0, result.ErrorCount)

	// 出力ファイルが作成されていることを確認
	outputFile := filepath.Join(outputDir, "large_cropped.jpg")
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)
}

// TestProcessImages_MultipleWorkers は複数ワーカーでの並行処理をテストします
func TestProcessImages_MultipleWorkers(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// 多数のテスト用画像ファイルを作成
	testFiles := []string{"worker1.jpg", "worker2.png", "worker3.webp", "worker4.jpeg", "worker5.jpg"}
	for _, file := range testFiles {
		path := filepath.Join(inputDir, file)
		require.NoError(t, createTestImage(path, 15, 15, filepath.Ext(file)[1:]))
	}

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        0,
		Y1:        0,
		X2:        0,
		Y2:        0,
		Suffix:    "parallel",
		Move:      false,
		Recursive: false,
		Workers:   3, // 3つのワーカー
	}

	filterConfig := &FilterConfig{
		Mode:    BlurMode,
		Radius:  1.0,
		RWeight: 0.299,
		GWeight: 0.587,
		BWeight: 0.114,
	}

	service, err := NewImageFilterService(procConfig, filterConfig)
	require.NoError(t, err)

	// パイプラインを実行
	result, err := service.ProcessImages()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 5, result.SuccessCount)
	assert.Equal(t, 0, result.ErrorCount)

	// すべての出力ファイルが作成されていることを確認
	for _, file := range testFiles {
		baseName := file[:len(file)-len(filepath.Ext(file))]
		expectedOutput := baseName + "_parallel" + filepath.Ext(file)
		outputPath := filepath.Join(outputDir, expectedOutput)
		_, err := os.Stat(outputPath)
		assert.NoError(t, err, "Output file should exist: %s", expectedOutput)
	}
}

// TestProcessImages_MixedFormats は異なる画像形式の混在処理をテストします
func TestProcessImages_MixedFormats(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// 異なる形式のテスト用画像ファイルを作成
	testFiles := map[string]string{
		"test.jpg":  "jpg",
		"test.jpeg": "jpeg",
		"test.png":  "png",
		"test.webp": "webp",
	}

	for file, format := range testFiles {
		path := filepath.Join(inputDir, file)
		require.NoError(t, createTestImage(path, 30, 30, format))
	}

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        5,
		Y1:        5,
		X2:        25,
		Y2:        25,
		Suffix:    "mixed",
		Move:      false,
		Recursive: false,
		Workers:   2,
	}

	filterConfig := &FilterConfig{
		Mode:    GrayscaleMode,
		Radius:  0.0,
		RWeight: 0.299,
		GWeight: 0.587,
		BWeight: 0.114,
	}

	service, err := NewImageFilterService(procConfig, filterConfig)
	require.NoError(t, err)

	// パイプラインを実行
	result, err := service.ProcessImages()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 4, result.SuccessCount)
	assert.Equal(t, 0, result.ErrorCount)

	// すべての出力ファイルが作成されていることを確認
	for file := range testFiles {
		baseName := file[:len(file)-len(filepath.Ext(file))]
		expectedOutput := baseName + "_mixed" + filepath.Ext(file)
		outputPath := filepath.Join(outputDir, expectedOutput)
		_, err := os.Stat(outputPath)
		assert.NoError(t, err, "Output file should exist: %s", expectedOutput)
	}
}
