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
