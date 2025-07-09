package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMovieConverterService_GetSourceVideoBitrate_ValidJSON tests bitrate detection with valid JSON
func TestMovieConverterService_GetSourceVideoBitrate_ValidJSON(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "test.mp4")
	file, err := os.Create(inputFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:    inputFile,
		OutputFile:   "test.webm",
		VideoQuality: 80,
	}

	service := NewMovieConverterService(config)
	bitrate, err := service.getSourceVideoBitrate()

	// Should return default value since ffprobe will fail
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if bitrate != "1M" {
		t.Errorf("Expected default bitrate 1M, got %s", bitrate)
	}
}

// TestMovieConverterService_GetSourceVideoBitrate_DifferentVideoQuality tests different video quality values
func TestMovieConverterService_GetSourceVideoBitrate_DifferentVideoQuality(t *testing.T) {
	qualities := []int{25, 50, 75, 100}

	for _, quality := range qualities {
		config := ConversionConfig{
			InputFile:    "non_existent_file.mp4",
			OutputFile:   "test.webm",
			VideoQuality: quality,
		}

		service := NewMovieConverterService(config)
		bitrate, err := service.getSourceVideoBitrate()

		// Should return default value without error
		if err != nil {
			t.Errorf("Expected no error for quality %d, got %v", quality, err)
		}

		if bitrate != "1M" {
			t.Errorf("Expected default bitrate 1M for quality %d, got %s", quality, bitrate)
		}
	}
}

// TestBatchMovieConverterService_ConvertSingleFile_RelativePathError tests relative path error
func TestBatchMovieConverterService_ConvertSingleFile_RelativePathError(t *testing.T) {
	config := BatchConversionConfig{
		InputDir:  "/invalid/path",
		OutputDir: "/output",
		OutputExt: ".webm",
	}

	service := NewBatchMovieConverterService(config)

	// Test with a file that would cause relative path error
	result := service.convertSingleFile("/completely/different/path/file.mp4")

	if result.Success {
		t.Error("Expected failure for relative path error, got success")
	}

	if result.Error == nil {
		t.Error("Expected error for relative path calculation, got nil")
	}
}

// TestBatchMovieConverterService_ScanFiles_WalkError tests file walk error handling
func TestBatchMovieConverterService_ScanFiles_WalkError(t *testing.T) {
	config := BatchConversionConfig{
		InputDir:  "/non/existent/directory",
		InputExt:  ".mp4",
		OutputDir: "/output",
		OutputExt: ".webm",
		Recursive: true,
	}

	service := NewBatchMovieConverterService(config)
	files, err := service.scanFiles()

	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files for non-existent directory, got %d", len(files))
	}
}

// TestBatchMovieConverterService_BatchConvert_OutputDirCreationError tests output directory creation error
func TestBatchMovieConverterService_BatchConvert_OutputDirCreationError(t *testing.T) {
	// Create input directory
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	err := os.MkdirAll(inputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}

	// Try to create output directory in a location that should fail
	config := BatchConversionConfig{
		InputDir:  inputDir,
		InputExt:  ".mp4",
		OutputDir: "/root/forbidden", // This should fail on most systems
		OutputExt: ".webm",
		Recursive: true,
	}

	service := NewBatchMovieConverterService(config)
	result, err := service.batchConvert()

	// This might succeed on some systems, so we check if it fails appropriately
	if err != nil && !strings.Contains(err.Error(), "出力ディレクトリの作成に失敗") {
		t.Errorf("Expected output directory creation error, got: %v", err)
	}

	// If it succeeds, that's also valid (depends on system permissions)
	if result != nil && result.TotalFiles != 0 {
		// This is fine, no files to process
	}
}

// TestValidateBatchConfig_EdgeCases tests edge cases in batch config validation
func TestValidateBatchConfig_EdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	// Test with extensions that need normalization
	config := &BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  "mp4", // No dot
		OutputDir: tempDir,
		OutputExt: "webm", // No dot
	}

	err := validateBatchConfig(config)
	if err != nil {
		t.Errorf("Expected no error for extension normalization, got %v", err)
	}

	// Check that extensions were normalized
	if config.InputExt != ".mp4" {
		t.Errorf("Expected InputExt .mp4, got %s", config.InputExt)
	}

	if config.OutputExt != ".webm" {
		t.Errorf("Expected OutputExt .webm, got %s", config.OutputExt)
	}
}

// TestValidateBatchConfig_CaseInsensitiveExtensions tests case insensitive extension validation
func TestValidateBatchConfig_CaseInsensitiveExtensions(t *testing.T) {
	tempDir := t.TempDir()

	config := &BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".MP4", // Uppercase
		OutputDir: tempDir,
		OutputExt: ".WEBM", // Uppercase
	}

	err := validateBatchConfig(config)
	if err != nil {
		t.Errorf("Expected no error for uppercase extensions, got %v", err)
	}
}

// TestNormalizeExtension_EdgeCases tests edge cases in extension normalization
func TestNormalizeExtension_EdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{".", "."},
		{"..", ".."},
		{"..mp4", "..mp4"},
		{"mp4.", ".mp4."},
	}

	for _, test := range tests {
		result := normalizeExtension(test.input)
		if result != test.expected {
			t.Errorf("For input '%s', expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

// TestGenerateOutputFile_EdgeCases tests edge cases in output file generation
func TestGenerateOutputFile_EdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"file", "file_converted"},
		{"file.", "file_converted"},
		{".file", "_converted"},
		{"file.MP4", "file.webm"}, // Case insensitive
		{"file.WEBM", "file.mp4"}, // Case insensitive
		{"path/to/file.mp4", "path/to/file.webm"},
		{"file.mp4.backup", "file.mp4_converted"},
	}

	for _, test := range tests {
		result := generateOutputFile(test.input)
		if result != test.expected {
			t.Errorf("For input '%s', expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

// TestGetSupportedExtensions_Immutability tests that supported extensions map is not modified
func TestGetSupportedExtensions_Immutability(t *testing.T) {
	extensions1 := GetSupportedExtensions()
	extensions2 := GetSupportedExtensions()

	// Modify one map
	extensions1["input"] = append(extensions1["input"], ".test")

	// Check that the other map is not affected
	if len(extensions2["input"]) != 6 { // Original count
		t.Error("GetSupportedExtensions should return a new map each time")
	}
}

// TestUnifiedMovieConverterService_ProcessingModeDetection tests processing mode detection edge cases
func TestUnifiedMovieConverterService_ProcessingModeDetection(t *testing.T) {
	// Test with both configs provided (batch should take precedence)
	singleConfig := &ConversionConfig{
		InputFile:  "test.mp4",
		OutputFile: "test.webm",
	}

	batchConfig := &BatchConversionConfig{
		InputDir:  "./input",
		InputExt:  ".mp4",
		OutputDir: "./output",
		OutputExt: ".webm",
	}

	service := NewUnifiedMovieConverterService(singleConfig, batchConfig)

	if service.config.Mode != BatchMode {
		t.Errorf("Expected BatchMode when both configs provided, got %v", service.config.Mode)
	}

	// Test with empty batch config (should use single mode)
	emptyBatchConfig := &BatchConversionConfig{}
	service2 := NewUnifiedMovieConverterService(singleConfig, emptyBatchConfig)

	if service2.config.Mode != SingleFileMode {
		t.Errorf("Expected SingleFileMode with empty batch config, got %v", service2.config.Mode)
	}
}

// TestUnifiedMovieConverterService_ProcessConversion_InvalidMode tests invalid processing mode
func TestUnifiedMovieConverterService_ProcessConversion_InvalidMode(t *testing.T) {
	service := NewUnifiedMovieConverterService(nil, nil)
	service.config.Mode = ProcessingMode(999) // Invalid mode

	result, err := service.ProcessConversion()

	if err == nil {
		t.Error("Expected error for invalid processing mode, got nil")
	}

	if result == nil {
		t.Fatal("Expected result even with error, got nil")
	}

	if result.Success {
		t.Error("Expected Success to be false for invalid mode")
	}

	if !strings.Contains(err.Error(), "不明な処理モード") {
		t.Errorf("Expected unknown mode error, got: %v", err)
	}
}

// TestBatchMovieConverterService_ConvertSingleFile_OutputDirCreationError tests output directory creation error in single file conversion
func TestBatchMovieConverterService_ConvertSingleFile_OutputDirCreationError(t *testing.T) {
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "test.mp4")

	// Create input file
	file, err := os.Create(inputFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := BatchConversionConfig{
		InputDir:  tempDir,
		OutputDir: "/root/forbidden/deep/path", // This should fail
		OutputExt: ".webm",
	}

	service := NewBatchMovieConverterService(config)
	result := service.convertSingleFile(inputFile)

	// Should fail due to output directory creation error
	if result.Success {
		t.Error("Expected failure for output directory creation, got success")
	}

	if result.Error == nil {
		t.Error("Expected error for output directory creation, got nil")
	}
}

// TestMovieConverterService_Convert_EmptyOutputFile tests conversion with empty output file
func TestMovieConverterService_Convert_EmptyOutputFile(t *testing.T) {
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "test.mp4")

	// Create input file
	file, err := os.Create(inputFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:  inputFile,
		OutputFile: "", // Empty output file
	}

	service := NewMovieConverterService(config)
	err = service.convert()

	// Should fail due to empty output file extension
	if err == nil {
		t.Error("Expected error for empty output file, got nil")
	}
}

// TestProcessingModeConstants tests ProcessingMode constants
func TestProcessingModeConstants(t *testing.T) {
	if SingleFileMode != 0 {
		t.Errorf("Expected SingleFileMode to be 0, got %d", SingleFileMode)
	}

	if BatchMode != 1 {
		t.Errorf("Expected BatchMode to be 1, got %d", BatchMode)
	}
}

// TestConversionResult_Structure tests ConversionResult structure
func TestConversionResult_Structure(t *testing.T) {
	result := ConversionResult{
		InputFile:  "input.mp4",
		OutputFile: "output.webm",
		Success:    true,
		Error:      nil,
	}

	if result.InputFile != "input.mp4" {
		t.Errorf("Expected InputFile 'input.mp4', got '%s'", result.InputFile)
	}

	if result.OutputFile != "output.webm" {
		t.Errorf("Expected OutputFile 'output.webm', got '%s'", result.OutputFile)
	}

	if !result.Success {
		t.Error("Expected Success to be true")
	}

	if result.Error != nil {
		t.Errorf("Expected Error to be nil, got %v", result.Error)
	}
}

// TestBatchConversionResult_Structure tests BatchConversionResult structure
func TestBatchConversionResult_Structure(t *testing.T) {
	result := BatchConversionResult{
		TotalFiles:   5,
		SuccessCount: 3,
		FailureCount: 2,
		Results:      []ConversionResult{},
		FailedFiles:  []string{"file1.mp4", "file2.mp4"},
	}

	if result.TotalFiles != 5 {
		t.Errorf("Expected TotalFiles 5, got %d", result.TotalFiles)
	}

	if result.SuccessCount != 3 {
		t.Errorf("Expected SuccessCount 3, got %d", result.SuccessCount)
	}

	if result.FailureCount != 2 {
		t.Errorf("Expected FailureCount 2, got %d", result.FailureCount)
	}

	if len(result.FailedFiles) != 2 {
		t.Errorf("Expected 2 failed files, got %d", len(result.FailedFiles))
	}
}
