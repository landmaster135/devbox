package usecases

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestImageFilterService はImageFilterServiceのテストクラス
type TestImageFilterService struct{}

// createTestImage はテスト用の小さな画像を作成します
func createTestImage(path string, width, height int, format string) error {
	// ディレクトリを作成
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// 小さなテスト画像を作成
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 簡単なパターンを描画
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x+y)%2 == 0 {
				img.Set(x, y, color.RGBA{255, 0, 0, 255}) // 赤
			} else {
				img.Set(x, y, color.RGBA{0, 255, 0, 255}) // 緑
			}
		}
	}

	// ファイルを作成
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// 形式に応じてエンコード
	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	case "png":
		return png.Encode(file, img)
	default:
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	}
}

// setupTestDirectories はテスト用ディレクトリを準備します
func setupTestDirectories(t *testing.T) (string, string, string, string) {
	baseDir := filepath.Join("test_data", "tmp")
	inputDir := filepath.Join(baseDir, "input")
	outputDir := filepath.Join(baseDir, "output")
	archiveDir := filepath.Join(baseDir, "archive")

	// ディレクトリを作成
	require.NoError(t, os.MkdirAll(inputDir, 0755))
	require.NoError(t, os.MkdirAll(outputDir, 0755))
	require.NoError(t, os.MkdirAll(archiveDir, 0755))

	return baseDir, inputDir, outputDir, archiveDir
}

// cleanupTestDirectories はテスト後のクリーンアップを行います
func cleanupTestDirectories(baseDir string) {
	os.RemoveAll(baseDir)
}

// TestNewImageFilterService_Normal は正常な設定でのサービス作成をテストします
func TestNewImageFilterService_Normal(t *testing.T) {
	procConfig := &ProcessingConfig{
		SrcDir:    "test_input",
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

	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, procConfig, service.ProcessingConfig)
	assert.Equal(t, filterConfig, service.FilterConfig)
}

// TestNewImageFilterService_InvalidCoordinates は無効な座標でのサービス作成をテストします
func TestNewImageFilterService_InvalidCoordinates(t *testing.T) {
	procConfig := &ProcessingConfig{
		SrcDir:    "test_input",
		OutDir:    "test_output",
		ArcDir:    "test_archive",
		X1:        100, // X2より大きい
		Y1:        50,
		X2:        50, // X1より小さい
		Y2:        100,
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

	assert.Error(t, err)
	assert.Nil(t, service)
	assert.Contains(t, err.Error(), "invalid coordinates")
}

// TestNewImageFilterService_InvalidWorkers は無効なワーカー数でのサービス作成をテストします
func TestNewImageFilterService_InvalidWorkers(t *testing.T) {
	procConfig := &ProcessingConfig{
		SrcDir:    "test_input",
		OutDir:    "test_output",
		ArcDir:    "test_archive",
		X1:        0,
		Y1:        0,
		X2:        0,
		Y2:        0,
		Suffix:    "filtered",
		Move:      false,
		Recursive: false,
		Workers:   0, // 無効な値
	}

	filterConfig := &FilterConfig{
		Mode:    BlurMode,
		Radius:  2.0,
		RWeight: 0.299,
		GWeight: 0.587,
		BWeight: 0.114,
	}

	service, err := NewImageFilterService(procConfig, filterConfig)

	assert.Error(t, err)
	assert.Nil(t, service)
	assert.Contains(t, err.Error(), "workers must be positive")
}

// TestNewImageFilterService_InvalidFilterMode は無効なフィルターモードでのサービス作成をテストします
func TestNewImageFilterService_InvalidFilterMode(t *testing.T) {
	procConfig := &ProcessingConfig{
		SrcDir:    "test_input",
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
		Mode:    FilterMode("invalid"), // 無効なモード
		Radius:  2.0,
		RWeight: 0.299,
		GWeight: 0.587,
		BWeight: 0.114,
	}

	service, err := NewImageFilterService(procConfig, filterConfig)

	assert.Error(t, err)
	assert.Nil(t, service)
	assert.Contains(t, err.Error(), "unsupported filter mode")
}

// TestNewImageFilterService_InvalidRGBWeights は無効なRGB重みでのサービス作成をテストします
func TestNewImageFilterService_InvalidRGBWeights(t *testing.T) {
	testCases := []struct {
		name     string
		rWeight  float64
		gWeight  float64
		bWeight  float64
		expected string
	}{
		{
			name:     "RWeight too high",
			rWeight:  1.5, // 無効な値
			gWeight:  0.587,
			bWeight:  0.114,
			expected: "r-weight must be 0.0-1.0",
		},
		{
			name:     "GWeight negative",
			rWeight:  0.299,
			gWeight:  -0.1, // 無効な値
			bWeight:  0.114,
			expected: "g-weight must be 0.0-1.0",
		},
		{
			name:     "BWeight too high",
			rWeight:  0.299,
			gWeight:  0.587,
			bWeight:  2.0, // 無効な値
			expected: "b-weight must be 0.0-1.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			procConfig := &ProcessingConfig{
				SrcDir:    "test_input",
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
				Mode:    GrayscaleMode,
				Radius:  2.0,
				RWeight: tc.rWeight,
				GWeight: tc.gWeight,
				BWeight: tc.bWeight,
			}

			service, err := NewImageFilterService(procConfig, filterConfig)

			assert.Error(t, err)
			assert.Nil(t, service)
			assert.Contains(t, err.Error(), tc.expected)
		})
	}
}

// TestCollectImagePaths_NonRecursive は非再帰モードでのファイル収集をテストします
func TestCollectImagePaths_NonRecursive(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// テスト用画像ファイルを作成
	testFiles := []string{"test1.jpg", "test2.png", "test3.webp", "test4.txt"}
	for _, file := range testFiles {
		path := filepath.Join(inputDir, file)
		if filepath.Ext(file) == ".txt" {
			// テキストファイル（対象外）
			require.NoError(t, os.WriteFile(path, []byte("test"), 0644))
		} else {
			// 画像ファイル
			require.NoError(t, createTestImage(path, 10, 10, filepath.Ext(file)[1:]))
		}
	}

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

	assert.NoError(t, err)
	assert.Len(t, paths, 3) // jpg, png, webpの3ファイル

	// パスが正しく収集されているかチェック
	expectedFiles := []string{"test1.jpg", "test2.png", "test3.webp"}
	for _, expectedFile := range expectedFiles {
		expectedPath := filepath.Join(inputDir, expectedFile)
		assert.Contains(t, paths, expectedPath)
	}
}

// TestCollectImagePaths_Recursive は再帰モードでのファイル収集をテストします
func TestCollectImagePaths_Recursive(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// サブディレクトリを作成
	subDir := filepath.Join(inputDir, "subdir")
	require.NoError(t, os.MkdirAll(subDir, 0755))

	// テスト用画像ファイルを作成
	testFiles := map[string]string{
		"test1.jpg":         inputDir,
		"test2.png":         inputDir,
		"subdir/test3.webp": inputDir,
		"subdir/test4.jpeg": inputDir,
	}

	for file, baseDir := range testFiles {
		path := filepath.Join(baseDir, file)
		require.NoError(t, createTestImage(path, 10, 10, filepath.Ext(file)[1:]))
	}

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

	assert.NoError(t, err)
	assert.Len(t, paths, 4) // 4つの画像ファイル

	// サブディレクトリのファイルも含まれているかチェック
	subDirFile := filepath.Join(subDir, "test3.webp")
	assert.Contains(t, paths, subDirFile)
}

// TestMoveOriginal_Normal は正常なファイル移動をテストします
func TestMoveOriginal_Normal(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// テスト用ファイルを作成
	srcFile := filepath.Join(inputDir, "test.jpg")
	require.NoError(t, createTestImage(srcFile, 10, 10, "jpg"))

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

	// ファイル移動を実行
	err = service.MoveOriginal(srcFile)

	assert.NoError(t, err)

	// 元ファイルが存在しないことを確認
	_, err = os.Stat(srcFile)
	assert.True(t, os.IsNotExist(err))

	// アーカイブディレクトリにファイルが移動されていることを確認
	archivedFile := filepath.Join(archiveDir, "test.jpg")
	_, err = os.Stat(archivedFile)
	assert.NoError(t, err)
}

// TestProcessImagesWithWorkers_Normal は正常な並行処理をテストします
func TestProcessImagesWithWorkers_Normal(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// テスト用画像ファイルを作成
	testFiles := []string{"test1.jpg", "test2.png"}
	for _, file := range testFiles {
		path := filepath.Join(inputDir, file)
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
		Suffix:    "filtered",
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

	// パスを収集
	paths, err := service.CollectImagePaths()
	require.NoError(t, err)
	require.Len(t, paths, 2)

	// 並行処理を実行
	result, err := service.ProcessImagesWithWorkers(paths)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 0, result.ErrorCount)

	// 出力ファイルが作成されていることを確認
	expectedOutputs := []string{"test1_filtered.jpg", "test2_filtered.png"}
	for _, expectedOutput := range expectedOutputs {
		outputPath := filepath.Join(outputDir, expectedOutput)
		_, err := os.Stat(outputPath)
		assert.NoError(t, err, "Output file should exist: %s", expectedOutput)
	}
}

// TestProcessImages_FullPipeline は完全なパイプラインをテストします
func TestProcessImages_FullPipeline(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// テスト用画像ファイルを作成
	testFiles := []string{"image1.jpg", "image2.png", "image3.webp"}
	for _, file := range testFiles {
		path := filepath.Join(inputDir, file)
		require.NoError(t, createTestImage(path, 30, 30, filepath.Ext(file)[1:]))
	}

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        5,
		Y1:        5,
		X2:        25,
		Y2:        25,
		Suffix:    "processed",
		Move:      true,
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

	// 完全なパイプラインを実行
	result, err := service.ProcessImages()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.ErrorCount)

	// 出力ファイルが作成されていることを確認
	expectedOutputs := []string{"image1_processed.jpg", "image2_processed.png", "image3_processed.webp"}
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

// TestCollectImagePaths_EmptyDirectory は空のディレクトリでのファイル収集をテストします
func TestCollectImagePaths_EmptyDirectory(t *testing.T) {
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

	paths, err := service.CollectImagePaths()

	assert.NoError(t, err)
	assert.Len(t, paths, 0) // 空のディレクトリなので0ファイル
}

// TestCollectImagePaths_NonExistentDirectory は存在しないディレクトリでのファイル収集をテストします
func TestCollectImagePaths_NonExistentDirectory(t *testing.T) {
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

	paths, err := service.CollectImagePaths()

	assert.Error(t, err)
	assert.Nil(t, paths)
	assert.Contains(t, err.Error(), "directory read failed")
}

// TestApplyFilterAndSave_BlurMode はぼかしフィルターの適用をテストします
func TestApplyFilterAndSave_BlurMode(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// テスト用画像ファイルを作成
	inputFile := filepath.Join(inputDir, "test_blur.jpg")
	require.NoError(t, createTestImage(inputFile, 50, 50, "jpg"))

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        10,
		Y1:        10,
		X2:        40,
		Y2:        40,
		Suffix:    "blurred",
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

	// フィルターを適用
	err = service.ApplyFilterAndSave(inputFile)

	assert.NoError(t, err)

	// 出力ファイルが作成されていることを確認
	outputFile := filepath.Join(outputDir, "test_blur_blurred.jpg")
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)
}

// TestApplyFilterAndSave_GrayscaleMode はグレースケールフィルターの適用をテストします
func TestApplyFilterAndSave_GrayscaleMode(t *testing.T) {
	baseDir, inputDir, outputDir, archiveDir := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// テスト用画像ファイルを作成
	inputFile := filepath.Join(inputDir, "test_gray.png")
	require.NoError(t, createTestImage(inputFile, 40, 40, "png"))

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    outputDir,
		ArcDir:    archiveDir,
		X1:        0,
		Y1:        0,
		X2:        0,
		Y2:        0, // 画像全体を対象
		Suffix:    "gray",
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

	// フィルターを適用
	err = service.ApplyFilterAndSave(inputFile)

	assert.NoError(t, err)

	// 出力ファイルが作成されていることを確認
	outputFile := filepath.Join(outputDir, "test_gray_gray.png")
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)
}

// TestMoveOriginal_DirectoryCreation はディレクトリ自動作成をテストします
func TestMoveOriginal_DirectoryCreation(t *testing.T) {
	baseDir, inputDir, _, _ := setupTestDirectories(t)
	defer cleanupTestDirectories(baseDir)

	// 存在しないアーカイブディレクトリを指定
	nonExistentArchiveDir := filepath.Join(baseDir, "new_archive", "subdir")

	// テスト用ファイルを作成
	srcFile := filepath.Join(inputDir, "test.jpg")
	require.NoError(t, createTestImage(srcFile, 10, 10, "jpg"))

	procConfig := &ProcessingConfig{
		SrcDir:    inputDir,
		OutDir:    "test_output",
		ArcDir:    nonExistentArchiveDir,
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

	// ファイル移動を実行（ディレクトリが自動作成される）
	err = service.MoveOriginal(srcFile)

	assert.NoError(t, err)

	// アーカイブディレクトリが作成されていることを確認
	_, err = os.Stat(nonExistentArchiveDir)
	assert.NoError(t, err)

	// ファイルが移動されていることを確認
	archivedFile := filepath.Join(nonExistentArchiveDir, "test.jpg")
	_, err = os.Stat(archivedFile)
	assert.NoError(t, err)
}
