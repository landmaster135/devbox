package usecases

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExifMirrorService_copyExifToWebp_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "exif_mirror_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	service := NewExifMirrorService()
	config := &Config{
		Verbose: true,
	}

	// テスト用のダミーWebPファイルを作成
	webpPath := filepath.Join(tempDir, "test.webp")
	webpData := []byte{
		'R', 'I', 'F', 'F', // RIFF header
		0x20, 0x00, 0x00, 0x00, // file size (32 bytes)
		'W', 'E', 'B', 'P', // WEBP signature
		'V', 'P', '8', ' ', // VP8 chunk
		0x10, 0x00, 0x00, 0x00, // chunk size (16 bytes)
		// VP8 data (16 bytes of dummy data)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	err = os.WriteFile(webpPath, webpData, 0644)
	require.NoError(t, err)

	// テスト用のダミーEXIFデータ
	exifData := []byte{
		0x45, 0x78, 0x69, 0x66, 0x00, 0x00, // Exif header
		0x49, 0x49, 0x2A, 0x00, // TIFF header (little endian)
		0x08, 0x00, 0x00, 0x00, // offset to first IFD
	}

	// WebPファイルにEXIFデータを書き込み
	err = service.writeExifToWebpFile(webpPath, exifData, config)
	assert.NoError(t, err)

	// ファイルが正常に更新されたことを確認
	updatedData, err := os.ReadFile(webpPath)
	require.NoError(t, err)
	assert.Greater(t, len(updatedData), len(webpData), "WebPファイルにEXIFチャンクが追加されている")

	// EXIFチャンクが含まれていることを確認
	assert.Contains(t, string(updatedData), "EXIF", "EXIFチャンクが含まれている")
}

func TestExifMirrorService_ValidateDirectory_Normal(t *testing.T) {
	// 一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "validate_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 正常なディレクトリの検証
	err = ValidateDirectory(tempDir)
	assert.NoError(t, err)
}

func TestExifMirrorService_ValidateDirectory_NotExists(t *testing.T) {
	// 存在しないディレクトリの検証
	err := ValidateDirectory("/non/existent/directory")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ディレクトリが存在しません")
}

func TestExifMirrorService_ValidateExtension_Normal(t *testing.T) {
	testCases := []struct {
		name      string
		extension string
		expectErr bool
	}{
		{"jpg extension", "jpg", false},
		{"jpeg extension", "jpeg", false},
		{"webp extension", "webp", false},
		{"png extension", "png", false},
		{"dot jpg extension", ".jpg", false},
		{"unsupported extension", "xyz", true},
		{"empty extension", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateExtension(tc.extension)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExifMirrorService_isTargetFile_Normal(t *testing.T) {
	service := NewExifMirrorService()

	testCases := []struct {
		name            string
		filePath        string
		targetExtension string
		expected        bool
	}{
		{"jpg file with jpg target", "test.jpg", "jpg", true},
		{"jpeg file with jpeg target", "test.jpeg", "jpeg", true},
		{"webp file with webp target", "test.webp", "webp", true},
		{"jpg file with webp target", "test.jpg", "webp", false},
		{"uppercase extension", "test.JPG", "jpg", true},
		{"dot extension", "test.jpg", ".jpg", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.isTargetFile(tc.filePath, tc.targetExtension)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestExifMirrorService_removeExtension_Normal(t *testing.T) {
	testCases := []struct {
		name     string
		fileName string
		expected string
	}{
		{"jpg file", "test.jpg", "test"},
		{"jpeg file", "test.jpeg", "test"},
		{"webp file", "test.webp", "test"},
		{"no extension", "test", "test"},
		{"multiple dots", "test.backup.jpg", "test.backup"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := removeExtension(tc.fileName)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestExifMirrorService_buildWebpFile_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "build_webp_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	service := NewExifMirrorService()
	config := &Config{
		Verbose: true,
	}

	webpPath := filepath.Join(tempDir, "test.webp")

	// テスト用のチャンクデータ
	chunks := []webpChunk{
		{
			FourCC: [4]byte{'V', 'P', '8', ' '},
			Data:   []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
		},
		{
			FourCC: [4]byte{'E', 'X', 'I', 'F'},
			Data:   []byte{0x45, 0x78, 0x69, 0x66, 0x00, 0x00},
		},
	}

	// WebPファイルを構築
	err = service.buildWebpFile(webpPath, chunks, config)
	assert.NoError(t, err)

	// ファイルが作成されたことを確認
	_, err = os.Stat(webpPath)
	assert.NoError(t, err)

	// ファイル内容を確認
	data, err := os.ReadFile(webpPath)
	require.NoError(t, err)

	// RIFFヘッダーとWEBPシグネチャを確認
	assert.Equal(t, "RIFF", string(data[0:4]))
	assert.Equal(t, "WEBP", string(data[8:12]))

	// チャンクが含まれていることを確認
	assert.Contains(t, string(data), "VP8 ")
	assert.Contains(t, string(data), "EXIF")
}

func TestExifMirrorService_MirrorExifData_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "mirror_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	service := NewExifMirrorService()

	// テスト用のソースディレクトリとターゲットディレクトリを作成
	sourceDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")
	err = os.MkdirAll(sourceDir, 0755)
	require.NoError(t, err)
	err = os.MkdirAll(targetDir, 0755)
	require.NoError(t, err)

	// テスト用のダミーJPGファイルを作成（ソース）
	sourceFile := filepath.Join(sourceDir, "test.jpg")
	jpgData := []byte{
		0xFF, 0xD8, 0xFF, 0xE0, // JPEG header
		0x00, 0x10, 'J', 'F', 'I', 'F', 0x00, 0x01, 0x01, 0x01, 0x00, 0x48, 0x00, 0x48, 0x00, 0x00,
		0xFF, 0xD9, // JPEG end
	}
	err = os.WriteFile(sourceFile, jpgData, 0644)
	require.NoError(t, err)

	// テスト用のダミーWebPファイルを作成（ターゲット）
	targetFile := filepath.Join(targetDir, "test.webp")
	webpData := []byte{
		'R', 'I', 'F', 'F', // RIFF header
		0x20, 0x00, 0x00, 0x00, // file size
		'W', 'E', 'B', 'P', // WEBP signature
		'V', 'P', '8', ' ', // VP8 chunk
		0x10, 0x00, 0x00, 0x00, // chunk size
		// VP8 data
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	err = os.WriteFile(targetFile, webpData, 0644)
	require.NoError(t, err)

	config := &Config{
		SourceFolderPath: sourceDir,
		TargetFolderPath: targetDir,
		SourceExtension:  "jpg",
		TargetExtension:  "webp",
		Recursive:        false,
		DryRun:           true, // ドライランでテスト
		Verbose:          false,
		WorkerCount:      1,
	}

	// MirrorExifDataを実行
	processedCount, errorCount, skipCount, err := service.MirrorExifData(config)
	assert.NoError(t, err)
	assert.Equal(t, 1, processedCount)
	assert.Equal(t, 0, errorCount)
	assert.Equal(t, 0, skipCount)
}

func TestExifMirrorService_findImageFiles_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "find_files_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	service := NewExifMirrorService()

	// テスト用のファイルを作成
	files := []string{"test1.jpg", "test2.webp", "test3.png", "test4.txt"}
	for _, file := range files {
		err = os.WriteFile(filepath.Join(tempDir, file), []byte("dummy"), 0644)
		require.NoError(t, err)
	}

	// JPGファイルを検索
	jpgFiles, err := service.findImageFiles(tempDir, "jpg", false)
	assert.NoError(t, err)
	assert.Len(t, jpgFiles, 1)
	assert.Contains(t, jpgFiles[0], "test1.jpg")

	// WebPファイルを検索
	webpFiles, err := service.findImageFiles(tempDir, "webp", false)
	assert.NoError(t, err)
	assert.Len(t, webpFiles, 1)
	assert.Contains(t, webpFiles[0], "test2.webp")

	// 全ての画像ファイルを検索
	allFiles, err := service.findImageFiles(tempDir, "", false)
	assert.NoError(t, err)
	assert.Len(t, allFiles, 3) // jpg, webp, png
}

func TestExifMirrorService_findCorrespondingSourceFile_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "corresponding_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	service := NewExifMirrorService()

	// テスト用のソースファイルを作成
	sourceDir := filepath.Join(tempDir, "source")
	err = os.MkdirAll(sourceDir, 0755)
	require.NoError(t, err)
	sourceFile := filepath.Join(sourceDir, "test.jpg")
	err = os.WriteFile(sourceFile, []byte("dummy"), 0644)
	require.NoError(t, err)

	// ターゲットファイルパス
	targetDir := filepath.Join(tempDir, "target")
	targetFile := filepath.Join(targetDir, "test.webp")

	config := &Config{
		SourceFolderPath: sourceDir,
		TargetFolderPath: targetDir,
		SourceExtension:  "jpg",
		TargetExtension:  "webp",
	}

	// 対応するソースファイルを検索
	foundSource := service.findCorrespondingSourceFile(targetFile, config)
	assert.Equal(t, sourceFile, foundSource)

	// 存在しないファイルの場合
	nonExistentTarget := filepath.Join(targetDir, "nonexistent.webp")
	foundSource = service.findCorrespondingSourceFile(nonExistentTarget, config)
	assert.Empty(t, foundSource)
}

func TestExifMirrorService_CopyFileExifSimple_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "copy_simple_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	service := NewExifMirrorService()

	// ソースファイルを作成
	sourceFile := filepath.Join(tempDir, "source.jpg")
	err = os.WriteFile(sourceFile, []byte("source data"), 0644)
	require.NoError(t, err)

	// ターゲットファイルを作成
	targetFile := filepath.Join(tempDir, "target.webp")
	err = os.WriteFile(targetFile, []byte("target data"), 0644)
	require.NoError(t, err)

	// ファイル時刻をコピー
	err = service.CopyFileExifSimple(sourceFile, targetFile)
	assert.NoError(t, err)

	// ファイル時刻が同じになったことを確認
	sourceInfo, err := os.Stat(sourceFile)
	require.NoError(t, err)
	targetInfo, err := os.Stat(targetFile)
	require.NoError(t, err)

	assert.Equal(t, sourceInfo.ModTime().Unix(), targetInfo.ModTime().Unix())
}

func TestExifMirrorService_BackupFile_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "backup_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	service := NewExifMirrorService()

	// テスト用のファイルを作成
	originalFile := filepath.Join(tempDir, "original.txt")
	originalData := []byte("original content")
	err = os.WriteFile(originalFile, originalData, 0644)
	require.NoError(t, err)

	// バックアップを作成
	backupPath, err := service.BackupFile(originalFile)
	assert.NoError(t, err)
	assert.Equal(t, originalFile+".backup", backupPath)

	// バックアップファイルが作成されたことを確認
	_, err = os.Stat(backupPath)
	assert.NoError(t, err)

	// バックアップファイルの内容が同じことを確認
	backupData, err := os.ReadFile(backupPath)
	require.NoError(t, err)
	assert.Equal(t, originalData, backupData)
}

func TestExifMirrorService_writeExifToWebpFile_InvalidFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "invalid_webp_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	service := NewExifMirrorService()
	config := &Config{Verbose: false}

	// 無効なWebPファイルを作成
	invalidFile := filepath.Join(tempDir, "invalid.webp")
	invalidData := []byte("not a webp file")
	err = os.WriteFile(invalidFile, invalidData, 0644)
	require.NoError(t, err)

	// EXIFデータの書き込みを試行（エラーになるはず）
	exifData := []byte{0x45, 0x78, 0x69, 0x66}
	err = service.writeExifToWebpFile(invalidFile, exifData, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "有効なWebPファイルではありません")
}

func TestExifMirrorService_ValidateDirectory_EmptyPath(t *testing.T) {
	err := ValidateDirectory("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ディレクトリパスが空です")
}

func TestExifMirrorService_ValidateDirectory_NotDirectory(t *testing.T) {
	// テスト用の一時ファイルを作成
	tempFile, err := os.CreateTemp("", "not_dir_test")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	err = ValidateDirectory(tempFile.Name())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "指定されたパスはディレクトリではありません")
}
