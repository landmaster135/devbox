package usecases

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestImageFilterServiceEdgeCase はエッジケースのテストクラス
type TestImageFilterServiceEdgeCase struct {}

// TestProcessImages_CollectPathsError はパス収集エラーのテストを行います
func TestProcessImages_CollectPathsError(t *testing.T) {
	procConfig := &ProcessingConfig{
		SrcDir:    "/non/existent/directory",
		OutDir:    "test_output",
		ArcDir:    "test_archive",
		X1:        0,
		Y1:        0,
		X2:        0,
		Y2:        0,
		Suffix:    "filtered",
		Move:      false,
		Recursive: false,
		Workers:   2,
	}

	filterConfig := &FilterConfig{
		Mode:    BlurMode,
		Radius:  2.0,
		RWeight: 0.299,
		GWeight: 0.587,
		BWeight: 0.114,
	}

	service, err := NewImageFilterService(procConfig, filterConfig)
	require.NoError(t, err)

	// パス収集でエラーが発生する
	result, err := service.ProcessImages()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "directory read failed")
}

// TestCollectImagePaths_RecursiveError は再帰モードでのエラーテストを行います
func TestCollectImagePaths_RecursiveError(t *testing.T) {
	procConfig := &ProcessingConfig{
		SrcDir:    "/non/existent/directory",
		OutDir:    "test_output",
		ArcDir:    "test_archive",
		X1:        0,
		Y1:        0,
		X2:        0,
		Y2:        0,
		Suffix:    "filtered",
		Move:      false,
		Recursive: true, // 再帰モード
		Workers:   2,
	}

	filterConfig := &FilterConfig{
		Mode:    BlurMode,
		Radius:  2.0,
		RWeight: 0.299,
		GWeight: 0.587,
		BWeight: 0.114,
	}

	service, err := NewImageFilterService(procConfig, filterConfig)
	require.NoError(t, err)

	paths, err := service.CollectImagePaths()

	assert.Error(t, err)
	assert.Empty(t, paths)
	assert.Contains(t, err.Error(), "directory walk failed")
}

// TestApplyFilterAndSave_InvalidCoordinates は無効な座標でのテストを行います
func TestApplyFilterAndSave_InvalidCoordinates(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// テスト用画像ファイルを作成
	inputFile := filepath.Join(inputDir, "test_invalid_coords.jpg")
	require.NoError(t, createTestImage(inputFile, 50, 50, "jpg"))

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        60, // 画像サイズを超える座標
		Y1:        60,
		X2:        70,
		Y2:        70,
		Suffix:    "invalid",
		Move:      false,
		Recursive: false,
		Workers:   1,
	}

	filterConfig := &FilterConfig{
		Mode:    BlurMode,
		Radius:  2.0,
		RWeight: 0.299,
		GWeight: 0.587,
		BWeight: 0.114,
	}

	service, err := NewImageFilterService(procConfig, filterConfig)
	require.NoError(t, err)

	// 無効な座標でエラーが発生する
	err = service.ApplyFilterAndSave(inputFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid rectangle")
}

// TestApplyFilterAndSave_NonExistentFile は存在しないファイルでのテストを行います
func TestApplyFilterAndSave_NonExistentFile(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        0,
		Y1:        0,
		X2:        0,
		Y2:        0,
		Suffix:    "nonexistent",
		Move:      false,
		Recursive: false,
		Workers:   1,
	}

	filterConfig := &FilterConfig{
		Mode:    BlurMode,
		Radius:  2.0,
		RWeight: 0.299,
		GWeight: 0.587,
		BWeight: 0.114,
	}

	service, err := NewImageFilterService(procConfig, filterConfig)
	require.NoError(t, err)

	// 存在しないファイルでエラーが発生する
	nonExistentFile := filepath.Join(inputDir, "nonexistent.jpg")
	err = service.ApplyFilterAndSave(nonExistentFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "open")
}

// TestApplyFilterAndSave_UnsupportedFormat は未対応形式でのテストを行います
func TestApplyFilterAndSave_UnsupportedFormat(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// 未対応の拡張子のファイルを作成
	unsupportedFile := filepath.Join(inputDir, "test.bmp")
	require.NoError(t, createTestImage(unsupportedFile, 20, 20, "jpg")) // 内容はJPEGだが拡張子はBMP

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        0,
		Y1:        0,
		X2:        0,
		Y2:        0,
		Suffix:    "unsupported",
		Move:      false,
		Recursive: false,
		Workers:   1,
	}

	filterConfig := &FilterConfig{
		Mode:    BlurMode,
		Radius:  2.0,
		RWeight: 0.299,
		GWeight: 0.587,
		BWeight: 0.114,
	}

	service, err := NewImageFilterService(procConfig, filterConfig)
	require.NoError(t, err)

	// 未対応形式でエラーが発生する
	err = service.ApplyFilterAndSave(unsupportedFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported extension")
}

// TestMoveOriginal_SourceNotExists は存在しないソースファイルでのテストを行います
func TestMoveOriginal_SourceNotExists(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        0,
		Y1:        0,
		X2:        0,
		Y2:        0,
		Suffix:    "filtered",
		Move:      true,
		Recursive: false,
		Workers:   2,
	}

	filterConfig := &FilterConfig{
		Mode:    BlurMode,
		Radius:  2.0,
		RWeight: 0.299,
		GWeight: 0.587,
		BWeight: 0.114,
	}

	service, err := NewImageFilterService(procConfig, filterConfig)
	require.NoError(t, err)

	// 存在しないファイルの移動でエラーが発生する
	nonExistentFile := filepath.Join(inputDir, "nonexistent.jpg")
	err = service.MoveOriginal(nonExistentFile)

	assert.Error(t, err)
}

// TestProcessImagesWithWorkers_MoveError は移動エラーでのテストを行います
func TestProcessImagesWithWorkers_MoveError(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// テスト用画像ファイルを作成
	testFile := filepath.Join(inputDir, "move_error.jpg")
	require.NoError(t, createTestImage(testFile, 20, 20, "jpg"))

	// アーカイブディレクトリを読み取り専用にしてエラーを発生させる
	require.NoError(t, os.Chmod(archiveDir, 0444))
	defer os.Chmod(archiveDir, 0755) // テスト後に権限を戻す

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        0,
		Y1:        0,
		X2:        0,
		Y2:        0,
		Suffix:    "move_error",
		Move:      true, // 移動オプション有効
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

	paths := []string{testFile}

	// 移動エラーが発生するが処理は継続される
	result, err := service.ProcessImagesWithWorkers(paths)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// フィルター適用は成功するが、移動でエラーが発生する
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.ErrorCount)
}

// TestApplyFilterAndSave_OutputDirectoryCreationError は出力ディレクトリ作成エラーのテストを行います
func TestApplyFilterAndSave_OutputDirectoryCreationError(t *testing.T) {
	baseDir, inputDir, _, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// テスト用画像ファイルを作成
	inputFile := filepath.Join(inputDir, "test_output_error.jpg")
	require.NoError(t, createTestImage(inputFile, 20, 20, "jpg"))

	// 無効な出力ディレクトリ（ファイルと同じ名前）を指定
	invalidOutputDir := filepath.Join(baseDir, "invalid_output")
	require.NoError(t, os.WriteFile(invalidOutputDir, []byte("not a directory"), 0644))

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    invalidOutputDir,
		ArcDir:    archiveDir,
		X1:        0,
		Y1:        0,
		X2:        0,
		Y2:        0,
		Suffix:    "output_error",
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

	// 出力ディレクトリ作成でエラーが発生する
	err = service.ApplyFilterAndSave(inputFile)

	assert.Error(t, err)
}

// TestApplyFilterAndSave_WebPFormat はWebP形式でのテストを行います
func TestApplyFilterAndSave_WebPFormat(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// WebP形式のテスト用画像ファイルを作成
	inputFile := filepath.Join(inputDir, "test_webp.webp")
	require.NoError(t, createTestImage(inputFile, 30, 30, "webp"))

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        5,
		Y1:        5,
		X2:        25,
		Y2:        25,
		Suffix:    "webp_test",
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

	// WebP形式の処理
	err = service.ApplyFilterAndSave(inputFile)

	assert.NoError(t, err)

	// 出力ファイルが作成されていることを確認
	outputFile := filepath.Join(outputDir, "test_webp_webp_test.webp")
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)
}

// TestProcessImagesWithWorkers_EmptyPaths は空のパスリストでのテストを行います
func TestProcessImagesWithWorkers_EmptyPaths(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        0,
		Y1:        0,
		X2:        0,
		Y2:        0,
		Suffix:    "empty",
		Move:      false,
		Recursive: false,
		Workers:   2,
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

	// 空のパスリストで処理
	emptyPaths := []string{}
	result, err := service.ProcessImagesWithWorkers(emptyPaths)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 0, result.ErrorCount)
}
