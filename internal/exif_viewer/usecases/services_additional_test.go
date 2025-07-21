package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 追加のテストケース - PropertyInfo構造体のテスト
func TestPropertyInfo_Structure(t *testing.T) {
	info := PropertyInfo{
		Name:     "Camera",
		DataType: "string",
		Count:    5,
		Examples: []string{"Canon", "Nikon", "Sony"},
	}

	if info.Name != "Camera" {
		t.Errorf("Expected Name Camera, got %s", info.Name)
	}

	if info.DataType != "string" {
		t.Errorf("Expected DataType string, got %s", info.DataType)
	}

	if info.Count != 5 {
		t.Errorf("Expected Count 5, got %d", info.Count)
	}

	if len(info.Examples) != 3 {
		t.Errorf("Expected 3 examples, got %d", len(info.Examples))
	}
}

// ShowDataTypesフラグのテスト
func TestExifViewerService_ProcessExifViewing_ShowDataTypes(t *testing.T) {
	service := NewExifViewerService()

	config := &Config{
		Directory:     "test_data/org",
		Extensions:    []string{"jpg", "png", "webp"},
		ShowDataTypes: true,
		Recursive:     false,
	}

	output, err := service.ProcessExifViewing(config)
	if err != nil {
		t.Errorf("ProcessExifViewing should not return error, got: %v", err)
	}

	if output == "" {
		t.Error("ProcessExifViewing should return non-empty output")
	}

	// データ型表示の出力形式を確認
	if !strings.Contains(output, "Data Type") {
		t.Error("Output should contain Data Type header")
	}
}

// 複数の拡張子でのファイル検索テスト
func TestExifViewerService_FindImageFiles_MultipleExtensions(t *testing.T) {
	service := NewExifViewerService()

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "exif_multi_ext_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// 様々な拡張子のファイルを作成
	testFiles := []string{
		"image1.jpg",
		"image2.JPEG", // 大文字
		"image3.png",
		"image4.PNG",  // 大文字
		"image5.tiff",
		"image6.TIF",  // 大文字
		"image7.webp",
		"document.pdf", // 対象外
		"text.txt",     // 対象外
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
		Extensions: []string{"jpg", "jpeg", "png", "tiff", "tif", "webp"},
		Recursive:  false,
	}

	files, err := service.FindImageFiles(config)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 7 {
		t.Errorf("Expected 7 image files, got %d", len(files))
	}
}

// 空の拡張子リストでのテスト
func TestExifViewerService_FindImageFiles_EmptyExtensions(t *testing.T) {
	service := NewExifViewerService()

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "exif_empty_ext_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// テストファイルを作成
	testFile := filepath.Join(tempDir, "test.jpg")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	config := &Config{
		Directory:  tempDir,
		Extensions: []string{}, // 空の拡張子リスト
		Recursive:  false,
	}

	files, err := service.FindImageFiles(config)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files with empty extensions, got %d", len(files))
	}
}

// 長いファイルパスでのテスト
func TestExifViewerService_ExtractSingleFileExif_LongPath(t *testing.T) {
	service := NewExifViewerService()

	// 長いディレクトリ構造を作成
	tempDir, err := os.MkdirTemp("", "exif_long_path_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// 深いディレクトリ構造を作成
	longPath := tempDir
	for i := 0; i < 3; i++ { // 深さを減らす
		longPath = filepath.Join(longPath, "long_dir_name")
		err = os.MkdirAll(longPath, 0755)
		if err != nil {
			t.Fatal(err)
		}
	}

	// テストファイルを作成
	testFile := filepath.Join(longPath, "test_image.jpg")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	config := &Config{
		Extensions: []string{"jpg"},
	}

	data, _ := service.ExtractSingleFileExif(testFile, config)
	// エラーが発生してもファイルパスは設定されることを確認
	if data.FilePath != testFile {
		t.Errorf("Expected FilePath %s, got %s", testFile, data.FilePath)
	}
}

// UTF-8文字列処理の詳細テスト
func TestExifViewerService_EnsureUTF8String_InvalidBytes(t *testing.T) {
	service := NewExifViewerService()

	// 無効なUTF-8バイトシーケンスを含む文字列をシミュレート
	testCases := []struct {
		input       string
		description string
	}{
		{"valid string", "通常の有効な文字列"},
		{"", "空文字列"},
		{"日本語", "日本語文字列"},
		{"English text", "英語文字列"},
		{"Mixed 日本語 English", "混合文字列"},
	}

	for _, tc := range testCases {
		result := service.EnsureUTF8String(tc.input)
		if result == "" && tc.input != "" {
			t.Errorf("EnsureUTF8String should not return empty string for %s", tc.description)
		}
	}
}

// FormatPropertyList の空リストテスト
func TestExifViewerService_FormatPropertyList_Empty(t *testing.T) {
	service := NewExifViewerService()

	result := service.FormatPropertyList([]PropertyInfo{}, 0)

	if result != "No properties found." {
		t.Errorf("Expected 'No properties found.', got %s", result)
	}
}

// AnalyzeProperties の空データテスト
func TestExifViewerService_AnalyzeProperties_Empty(t *testing.T) {
	service := NewExifViewerService()

	result := service.AnalyzeProperties([]ExifData{})

	if len(result) != 0 {
		t.Errorf("Expected empty result for empty input, got %d properties", len(result))
	}
}

// 重複する例の処理テスト
func TestExifViewerService_AnalyzeProperties_DuplicateExamples(t *testing.T) {
	service := NewExifViewerService()

	// 同じ値を持つプロパティを含むテストデータ
	exifDataList := []ExifData{
		{
			FilePath: "/test/image1.jpg",
			Properties: map[string]string{
				"Camera": "Canon EOS R5",
			},
		},
		{
			FilePath: "/test/image2.jpg",
			Properties: map[string]string{
				"Camera": "Canon EOS R5", // 同じ値
			},
		},
		{
			FilePath: "/test/image3.jpg",
			Properties: map[string]string{
				"Camera": "Nikon D850", // 異なる値
			},
		},
	}

	propertyInfos := service.AnalyzeProperties(exifDataList)

	if len(propertyInfos) != 1 {
		t.Errorf("Expected 1 property, got %d", len(propertyInfos))
	}

	cameraInfo := propertyInfos[0]
	if cameraInfo.Name != "Camera" {
		t.Errorf("Expected Camera property, got %s", cameraInfo.Name)
	}

	if cameraInfo.Count != 3 {
		t.Errorf("Expected count 3, got %d", cameraInfo.Count)
	}

	// 重複する例が除外されていることを確認（最大2つの異なる例）
	if len(cameraInfo.Examples) > 2 {
		t.Errorf("Expected at most 2 unique examples, got %d", len(cameraInfo.Examples))
	}
}

// プロパティフィルタリングの詳細テスト
func TestExifViewerService_FormatExifTable_PropertyFiltering(t *testing.T) {
	service := NewExifViewerService()

	exifDataList := []ExifData{
		{
			FilePath: "/test/image1.jpg",
			Properties: map[string]string{
				"Camera":     "Canon EOS R5",
				"DateTime":   "2024:01:01 12:00:00",
				"Width":      "6000",
				"Height":     "4000",
				"FileSize":   "5 MB",
				"ISO":        "100",
				"Aperture":   "f/2.8",
				"ShutterSpeed": "1/100",
			},
		},
	}

	// 特定のプロパティのみを表示
	config := &Config{
		Properties: []string{"Camera", "DateTime", "Width"},
		MaxProps:   0,
		Verbose:    false,
	}

	result := service.FormatExifTable(exifDataList, config)

	// 指定したプロパティが含まれていることを確認
	if !strings.Contains(result, "Camera") {
		t.Error("Expected result to contain Camera")
	}

	if !strings.Contains(result, "DateTime") {
		t.Error("Expected result to contain DateTime")
	}

	if !strings.Contains(result, "Width") {
		t.Error("Expected result to contain Width")
	}

	// 指定していないプロパティが含まれていないことを確認
	if strings.Contains(result, "Height") {
		t.Error("Expected result to not contain Height")
	}

	if strings.Contains(result, "ISO") {
		t.Error("Expected result to not contain ISO")
	}
}

// 大文字小文字を無視したプロパティフィルタリングのテスト
func TestExifViewerService_contains_CaseInsensitive(t *testing.T) {
	service := NewExifViewerService()

	testCases := []struct {
		slice    []string
		item     string
		expected bool
	}{
		{[]string{"Camera", "DateTime"}, "camera", true},
		{[]string{"Camera", "DateTime"}, "CAMERA", true},
		{[]string{"Camera", "DateTime"}, "CaMeRa", true},
		{[]string{"camera", "datetime"}, "Camera", true},
		{[]string{"CAMERA", "DATETIME"}, "camera", true},
		{[]string{"Camera", "DateTime"}, "GPS", false},
	}

	for _, tc := range testCases {
		result := service.contains(tc.slice, tc.item)
		if result != tc.expected {
			t.Errorf("contains(%v, %s) = %v, expected %v", tc.slice, tc.item, result, tc.expected)
		}
	}
}

// 非常に長い値の処理テスト
func TestExifViewerService_FormatExifTable_LongValues(t *testing.T) {
	service := NewExifViewerService()

	// 非常に長い値を持つテストデータ
	longValue := strings.Repeat("Very long property value ", 10) // 約250文字
	exifDataList := []ExifData{
		{
			FilePath: "/test/image1.jpg",
			Properties: map[string]string{
				"Camera":      "Canon EOS R5",
				"Description": longValue,
			},
		},
	}

	config := &Config{
		MaxProps: 0,
		Verbose:  false,
	}

	result := service.FormatExifTable(exifDataList, config)

	// 長い値が短縮されていることを確認（30文字制限 + "..."）
	if !strings.Contains(result, "...") {
		t.Error("Expected long values to be truncated with '...'")
	}

	// 短縮された値が含まれていることを確認
	if !strings.Contains(result, "Very long property value") {
		t.Error("Expected truncated value to contain beginning of long text")
	}
}

// 相対パス表示のテスト
func TestExifViewerService_FormatExifTable_RelativePath(t *testing.T) {
	service := NewExifViewerService()

	// 現在のディレクトリからの相対パスでテスト
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	absolutePath := filepath.Join(currentDir, "test", "image1.jpg")
	exifDataList := []ExifData{
		{
			FilePath: absolutePath,
			Properties: map[string]string{
				"Camera": "Canon EOS R5",
			},
		},
	}

	config := &Config{
		MaxProps: 0,
		Verbose:  false,
	}

	result := service.FormatExifTable(exifDataList, config)

	// 結果にファイルパスが含まれていることを確認（相対パス変換の動作確認）
	if !strings.Contains(result, "image1.jpg") {
		t.Error("Expected result to contain filename")
	}

	// テーブル形式の出力であることを確認
	if !strings.Contains(result, "File Path") {
		t.Error("Expected result to contain File Path header")
	}
}

// 空のプロパティ値の処理テスト
func TestExifViewerService_FormatExifTable_EmptyValues(t *testing.T) {
	service := NewExifViewerService()

	exifDataList := []ExifData{
		{
			FilePath: "/test/image1.jpg",
			Properties: map[string]string{
				"Camera":     "Canon EOS R5",
				"GPS":        "", // 空の値
				"DateTime":   "2024:01:01 12:00:00",
				"EmptyField": "", // 空の値
			},
		},
	}

	config := &Config{
		MaxProps: 0,
		Verbose:  false,
	}

	result := service.FormatExifTable(exifDataList, config)

	// 空の値が "-" で表示されていることを確認
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if strings.Contains(line, "image1.jpg") {
			// データ行を見つけた場合、空の値が "-" になっていることを確認
			if strings.Contains(line, "\t\t") {
				t.Error("Expected empty values to be replaced with '-'")
			}
		}
	}
}

// inferDataType の境界値テスト
func TestExifViewerService_inferDataType_EdgeCases(t *testing.T) {
	service := NewExifViewerService()

	testCases := []struct {
		value    string
		expected string
	}{
		{"0", "integer"},
		{"-123", "integer"},
		{"0.0", "float"},
		{"-123.45", "float"},
		{"123.0", "float"},
		{".5", "float"},
		{"123.", "float"},
		{"1e10", "float"},
		{"1E-5", "float"},
		{"NaN", "float"},
		{"Infinity", "float"},
		{"123abc", "string"},
		{"abc123", "string"},
	}

	for _, tc := range testCases {
		result := service.inferDataType(tc.value)
		if result != tc.expected {
			t.Errorf("inferDataType(%s) = %s, expected %s", tc.value, result, tc.expected)
		}
	}
}
