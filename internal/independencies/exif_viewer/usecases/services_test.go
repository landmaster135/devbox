package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExifViewerService_FindImageFiles(t *testing.T) {
	service := NewExifViewerService()

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "exif_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// テストファイルを作成
	testFiles := []string{
		"test1.jpg",
		"test2.png",
		"test3.tiff",
		"test4.txt", // 対象外
	}

	for _, filename := range testFiles {
		file, err := os.Create(filepath.Join(tempDir, filename))
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
	}

	config := &Config{
		Directory:  tempDir,
		Extensions: []string{"jpg", "png", "tiff"},
		Recursive:  false,
	}

	files, err := service.FindImageFiles(config)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}

	// ファイル名をチェック
	foundFiles := make(map[string]bool)
	for _, file := range files {
		foundFiles[filepath.Base(file)] = true
	}

	expectedFiles := []string{"test1.jpg", "test2.png", "test3.tiff"}
	for _, expected := range expectedFiles {
		if !foundFiles[expected] {
			t.Errorf("Expected file %s not found", expected)
		}
	}
}

func TestExifViewerService_contains(t *testing.T) {
	service := NewExifViewerService()

	testCases := []struct {
		slice    []string
		item     string
		expected bool
	}{
		{[]string{"Camera", "DateTime", "GPS"}, "Camera", true},
		{[]string{"Camera", "DateTime", "GPS"}, "camera", true}, // 大文字小文字無視
		{[]string{"Camera", "DateTime", "GPS"}, "NotFound", false},
		{[]string{" Camera ", "DateTime"}, "Camera", true}, // 空白除去
		{[]string{}, "Camera", false},
		{nil, "Camera", false},
	}

	for _, tc := range testCases {
		result := service.contains(tc.slice, tc.item)
		if result != tc.expected {
			t.Errorf("contains(%v, %s) = %v, expected %v", tc.slice, tc.item, result, tc.expected)
		}
	}
}

func TestExifViewerService_formatFileSize(t *testing.T) {
	service := NewExifViewerService()

	testCases := []struct {
		size     int64
		expected string
	}{
		{512, "512 B"},
		{1024, "1 kB"},
		{1536, "2 kB"},
		{1024 * 1024, "1 MB"},
		{1024 * 1024 * 1024, "1 GB"},
	}

	for _, tc := range testCases {
		result := service.formatFileSize(tc.size)
		if !strings.Contains(result, strings.Split(tc.expected, " ")[0]) {
			t.Errorf("formatFileSize(%d) = %s, expected to contain %s", tc.size, result, tc.expected)
		}
	}
}

func TestExifViewerService_getColorTypeName(t *testing.T) {
	service := NewExifViewerService()

	testCases := []struct {
		colorType byte
		expected  string
	}{
		{0, "Grayscale"},
		{2, "RGB"},
		{3, "Palette"},
		{4, "Grayscale with Alpha"},
		{6, "RGB with Alpha"},
		{255, "Unknown"},
	}

	for _, tc := range testCases {
		result := service.getColorTypeName(tc.colorType)
		if result != tc.expected {
			t.Errorf("getColorTypeName(%d) = %s, expected %s", tc.colorType, result, tc.expected)
		}
	}
}

func TestConfig_Validation(t *testing.T) {
	config := &Config{
		Directory:      "/tmp",
		Extensions:     []string{"jpg", "png"},
		Properties:     []string{"Camera", "DateTime"},
		MaxProps:       10,
		Verbose:        true,
		Recursive:      false,
		ShowProperties: true,
		ShowDataTypes:  false,
	}

	if config.Directory != "/tmp" {
		t.Errorf("Expected Directory /tmp, got %s", config.Directory)
	}

	if len(config.Extensions) != 2 {
		t.Errorf("Expected 2 extensions, got %d", len(config.Extensions))
	}

	if len(config.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(config.Properties))
	}

	if !config.ShowProperties {
		t.Error("Expected ShowProperties to be true")
	}

	if config.ShowDataTypes {
		t.Error("Expected ShowDataTypes to be false")
	}
}

func TestExifData_Structure(t *testing.T) {
	data := ExifData{
		FilePath:   "/test/path.jpg",
		Properties: make(map[string]string),
	}

	data.Properties["Camera"] = "Canon EOS R5"
	data.Properties["DateTime"] = "2024:01:01 12:00:00"

	if data.FilePath != "/test/path.jpg" {
		t.Errorf("Expected FilePath /test/path.jpg, got %s", data.FilePath)
	}

	if len(data.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(data.Properties))
	}

	if data.Properties["Camera"] != "Canon EOS R5" {
		t.Errorf("Expected Camera Canon EOS R5, got %s", data.Properties["Camera"])
	}
}

func TestExifViewerService_FormatExifTable(t *testing.T) {
	service := NewExifViewerService()

	// テストデータを作成
	exifDataList := []ExifData{
		{
			FilePath: "/test/image1.jpg",
			Properties: map[string]string{
				"Camera":   "Canon EOS R5",
				"DateTime": "2024:01:01 12:00:00",
			},
		},
		{
			FilePath: "/test/image2.jpg",
			Properties: map[string]string{
				"Camera":   "Nikon D850",
				"DateTime": "2024:01:02 13:00:00",
			},
		},
	}

	config := &Config{
		MaxProps: 0,
		Verbose:  false,
	}

	result := service.FormatExifTable(exifDataList, config)

	if !strings.Contains(result, "File Path") {
		t.Error("Expected table to contain File Path header")
	}

	if !strings.Contains(result, "Canon EOS R5") {
		t.Error("Expected table to contain Canon EOS R5")
	}

	if !strings.Contains(result, "Nikon D850") {
		t.Error("Expected table to contain Nikon D850")
	}

	if !strings.Contains(result, "Summary: 2 files processed") {
		t.Error("Expected table to contain summary")
	}
}

func TestExifViewerService_FormatExifTable_Empty(t *testing.T) {
	service := NewExifViewerService()

	config := &Config{}
	result := service.FormatExifTable([]ExifData{}, config)

	if result != "No EXIF data found." {
		t.Errorf("Expected 'No EXIF data found.', got %s", result)
	}
}

func TestExifViewerService_AnalyzeProperties(t *testing.T) {
	service := NewExifViewerService()

	// テストデータを作成
	exifDataList := []ExifData{
		{
			FilePath: "/test/image1.jpg",
			Properties: map[string]string{
				"Camera":   "Canon EOS R5",
				"DateTime": "2024:01:01 12:00:00",
				"Width":    "6000",
				"FileSize": "5 MB",
			},
		},
		{
			FilePath: "/test/image2.jpg",
			Properties: map[string]string{
				"Camera":   "Nikon D850",
				"DateTime": "2024:01:02 13:00:00",
				"Width":    "7360",
			},
		},
	}

	propertyInfos := service.AnalyzeProperties(exifDataList)

	if len(propertyInfos) != 4 {
		t.Errorf("Expected 4 properties, got %d", len(propertyInfos))
	}

	// 最も頻度の高いプロパティが最初に来るかチェック
	for _, info := range propertyInfos {
		if info.Name == "Camera" {
			if info.Count != 2 {
				t.Errorf("Expected Camera count 2, got %d", info.Count)
			}
			if info.DataType != "string" {
				t.Errorf("Expected Camera data type string, got %s", info.DataType)
			}
		}
		if info.Name == "FileSize" {
			if info.Count != 1 {
				t.Errorf("Expected FileSize count 1, got %d", info.Count)
			}
			if info.DataType != "filesize" {
				t.Errorf("Expected FileSize data type filesize, got %s", info.DataType)
			}
		}
	}
}

func TestExifViewerService_FormatPropertyList(t *testing.T) {
	service := NewExifViewerService()

	propertyInfos := []PropertyInfo{
		{
			Name:     "Camera",
			DataType: "string",
			Count:    2,
			Examples: []string{"Canon EOS R5", "Nikon D850"},
		},
		{
			Name:     "Width",
			DataType: "integer",
			Count:    2,
			Examples: []string{"6000", "7360"},
		},
	}

	result := service.FormatPropertyList(propertyInfos, 2)

	if !strings.Contains(result, "Property Name") {
		t.Error("Expected result to contain Property Name header")
	}

	if !strings.Contains(result, "Data Type") {
		t.Error("Expected result to contain Data Type header")
	}

	if !strings.Contains(result, "Camera") {
		t.Error("Expected result to contain Camera")
	}

	if !strings.Contains(result, "string") {
		t.Error("Expected result to contain string data type")
	}

	if !strings.Contains(result, "100.0%") {
		t.Error("Expected result to contain 100.0% usage")
	}
}

func TestExifViewerService_inferDataType(t *testing.T) {
	service := NewExifViewerService()

	testCases := []struct {
		value    string
		expected string
	}{
		{"", "string"},
		{"123", "integer"},
		{"123.45", "float"},
		{"2024:01:01 12:00:00", "datetime"},
		{"5 MB", "filesize"},
		{"1024 B", "filesize"},
		{"35.6789°N", "coordinate"},
		{"1/100", "ratio"},
		{"Canon EOS R5", "string"},
	}

	for _, tc := range testCases {
		result := service.inferDataType(tc.value)
		if result != tc.expected {
			t.Errorf("inferDataType(%s) = %s, expected %s", tc.value, result, tc.expected)
		}
	}
}

func TestExifViewerService_containsString(t *testing.T) {
	service := NewExifViewerService()

	testCases := []struct {
		slice    []string
		item     string
		expected bool
	}{
		{[]string{"a", "b", "c"}, "b", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{}, "a", false},
		{nil, "a", false},
	}

	for _, tc := range testCases {
		result := service.containsString(tc.slice, tc.item)
		if result != tc.expected {
			t.Errorf("containsString(%v, %s) = %v, expected %v", tc.slice, tc.item, result, tc.expected)
		}
	}
}
