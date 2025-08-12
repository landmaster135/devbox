package usecases

import (
	"os"
	"strings"
	"testing"
)

// 実際のサンプルファイルを使用したテスト
func TestExifViewerService_ExtractSingleFileExif_JPEG(t *testing.T) {
	service := NewExifViewerService()

	// テストファイルの存在確認
	testFile := "test_data/org/sample_01_01.jpg"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatalf("Required test file not found: %s", testFile)
	}

	config := &Config{
		Extensions: []string{"jpg"},
	}

	data, err := service.ExtractSingleFileExif(testFile, config)
	// このJPEGファイルにはEXIFデータがないため、"no exif data"エラーが返される
	if err == nil {
		t.Error("Expected error for JPEG file without EXIF data")
	} else if err.Error() != "no exif data" {
		t.Errorf("Expected 'no exif data' error, got: %v", err)
	}

	// エラーが返されてもファイルパスは設定される
	if data.FilePath != testFile {
		t.Errorf("Expected FilePath %s, got %s", testFile, data.FilePath)
	}

	// Propertiesは初期化される
	if data.Properties == nil {
		t.Error("Properties should not be nil")
	}
}

func TestExifViewerService_ExtractSingleFileExif_PNG(t *testing.T) {
	service := NewExifViewerService()

	// テストファイルの存在確認
	testFile := "test_data/org/sample_01_01.png"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatalf("Required test file not found: %s", testFile)
	}

	config := &Config{
		Extensions: []string{"png"},
	}

	data, err := service.ExtractSingleFileExif(testFile, config)
	if err != nil {
		t.Errorf("ExtractSingleFileExif should not return error for valid PNG file, got: %v", err)
	}

	if data.FilePath != testFile {
		t.Errorf("Expected FilePath %s, got %s", testFile, data.FilePath)
	}

	if data.Properties == nil {
		t.Error("Properties should not be nil")
	}

	// PNG固有のプロパティが含まれていることを確認
	if data.Properties["File Type"] != "PNG" {
		t.Errorf("Expected File Type PNG, got %s", data.Properties["File Type"])
	}

	if data.Properties["MIME Type"] != "image/png" {
		t.Errorf("Expected MIME Type image/png, got %s", data.Properties["MIME Type"])
	}
}

func TestExifViewerService_ExtractSingleFileExif_WebP(t *testing.T) {
	service := NewExifViewerService()

	// テストファイルの存在確認
	testFile := "test_data/org/sample_01_01.webp"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatalf("Required test file not found: %s", testFile)
	}

	config := &Config{
		Extensions: []string{"webp"},
	}

	data, err := service.ExtractSingleFileExif(testFile, config)
	if err != nil {
		t.Errorf("ExtractSingleFileExif should not return error for valid WebP file, got: %v", err)
	}

	if data.FilePath != testFile {
		t.Errorf("Expected FilePath %s, got %s", testFile, data.FilePath)
	}

	if data.Properties == nil {
		t.Error("Properties should not be nil")
	}

	// WebP固有のプロパティが含まれていることを確認
	if data.Properties["File Type"] != "WEBP" {
		t.Errorf("Expected File Type WEBP, got %s", data.Properties["File Type"])
	}
}

func TestExifViewerService_ExtractExifData_MultipleFiles(t *testing.T) {
	service := NewExifViewerService()

	// テストファイルの存在確認
	testFiles := []string{
		"test_data/org/sample_01_01.jpg",
		"test_data/org/sample_01_01.png",
		"test_data/org/sample_01_01.webp",
	}

	for _, testFile := range testFiles {
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Fatalf("Required test file not found: %s", testFile)
		}
	}

	config := &Config{
		Extensions: []string{"jpg", "png", "webp"},
	}

	exifDataList, err := service.ExtractExifData(testFiles, config)
	if err != nil {
		t.Errorf("ExtractExifData should not return error, got: %v", err)
	}

	if len(exifDataList) == 0 {
		t.Error("ExtractExifData should return at least one ExifData")
	}

	// 各ファイルのデータが正しく抽出されていることを確認
	for _, data := range exifDataList {
		if data.Properties == nil {
			t.Error("Properties should not be nil")
		}
		if data.FilePath == "" {
			t.Error("FilePath should not be empty")
		}
	}
}

func TestExifViewerService_ProcessExifViewing_RealFiles(t *testing.T) {
	service := NewExifViewerService()

	config := &Config{
		Directory:  "test_data/org",
		Extensions: []string{"jpg", "png", "webp"},
		Recursive:  false,
	}

	output, err := service.ProcessExifViewing(config)
	if err != nil {
		t.Errorf("ProcessExifViewing should not return error, got: %v", err)
	}

	if output == "" {
		t.Error("ProcessExifViewing should return non-empty output")
	}

	// 出力にファイル名が含まれていることを確認
	if !strings.Contains(output, "sample_01_01") {
		t.Error("Output should contain sample file names")
	}

	// サマリーが含まれていることを確認
	if !strings.Contains(output, "Summary:") {
		t.Error("Output should contain summary")
	}
}

func TestExifViewerService_ProcessExifViewing_ShowProperties(t *testing.T) {
	service := NewExifViewerService()

	config := &Config{
		Directory:      "test_data/org",
		Extensions:     []string{"jpg", "png", "webp"},
		ShowProperties: true,
		Recursive:      false,
	}

	output, err := service.ProcessExifViewing(config)
	if err != nil {
		t.Errorf("ProcessExifViewing should not return error, got: %v", err)
	}

	if output == "" {
		t.Error("ProcessExifViewing should return non-empty output")
	}

	// プロパティ一覧の出力形式を確認
	if !strings.Contains(output, "Property Name") {
		t.Error("Output should contain Property Name header")
	}

	if !strings.Contains(output, "Data Type") {
		t.Error("Output should contain Data Type header")
	}

	if !strings.Contains(output, "Frequency") {
		t.Error("Output should contain Frequency header")
	}
}

func TestExifViewerService_ProcessExifViewing_PropertyFiltering(t *testing.T) {
	service := NewExifViewerService()

	config := &Config{
		Directory:  "test_data/org",
		Extensions: []string{"jpg", "png", "webp"},
		Properties: []string{"File Name", "File Type"},
		Recursive:  false,
	}

	output, err := service.ProcessExifViewing(config)
	if err != nil {
		t.Errorf("ProcessExifViewing should not return error, got: %v", err)
	}

	if output == "" {
		t.Error("ProcessExifViewing should return non-empty output")
	}

	// 指定したプロパティが含まれていることを確認
	if !strings.Contains(output, "File Name") {
		t.Error("Output should contain File Name property")
	}

	if !strings.Contains(output, "File Type") {
		t.Error("Output should contain File Type property")
	}
}
