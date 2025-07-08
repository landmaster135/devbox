package usecases

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewBatchMovieConverterService tests the creation of BatchMovieConverterService
func TestNewBatchMovieConverterService_Normal(t *testing.T) {
	config := BatchConversionConfig{
		InputDir:       "./input",
		InputExt:       ".mp4",
		OutputDir:      "./output",
		OutputExt:      ".webm",
		Recursive:      true,
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewBatchMovieConverterService(config)

	if service == nil {
		t.Fatal("NewBatchMovieConverterService returned nil")
	}

	if service.config.InputDir != config.InputDir {
		t.Errorf("Expected InputDir %s, got %s", config.InputDir, service.config.InputDir)
	}

	if service.config.OutputDir != config.OutputDir {
		t.Errorf("Expected OutputDir %s, got %s", config.OutputDir, service.config.OutputDir)
	}

	if service.config.Recursive != config.Recursive {
		t.Errorf("Expected Recursive %v, got %v", config.Recursive, service.config.Recursive)
	}
}

// TestValidateBatchConfig tests batch configuration validation
func TestValidateBatchConfig_Normal(t *testing.T) {
	// Create temporary directories
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	err := os.MkdirAll(inputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}

	config := &BatchConversionConfig{
		InputDir:  inputDir,
		InputExt:  "mp4",
		OutputDir: filepath.Join(tempDir, "output"),
		OutputExt: "webm",
	}

	err = validateBatchConfig(config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that extensions are normalized
	if config.InputExt != ".mp4" {
		t.Errorf("Expected InputExt .mp4, got %s", config.InputExt)
	}

	if config.OutputExt != ".webm" {
		t.Errorf("Expected OutputExt .webm, got %s", config.OutputExt)
	}
}

// TestValidateBatchConfig tests validation with missing input directory
func TestValidateBatchConfig_MissingInputDir(t *testing.T) {
	config := &BatchConversionConfig{
		InputExt:  ".mp4",
		OutputDir: "./output",
		OutputExt: ".webm",
	}

	err := validateBatchConfig(config)
	if err == nil {
		t.Error("Expected error for missing input directory, got nil")
	}
}

// TestValidateBatchConfig tests validation with missing input extension
func TestValidateBatchConfig_MissingInputExt(t *testing.T) {
	tempDir := t.TempDir()
	config := &BatchConversionConfig{
		InputDir:  tempDir,
		OutputDir: "./output",
		OutputExt: ".webm",
	}

	err := validateBatchConfig(config)
	if err == nil {
		t.Error("Expected error for missing input extension, got nil")
	}
}

// TestValidateBatchConfig tests validation with missing output directory
func TestValidateBatchConfig_MissingOutputDir(t *testing.T) {
	tempDir := t.TempDir()
	config := &BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputExt: ".webm",
	}

	err := validateBatchConfig(config)
	if err == nil {
		t.Error("Expected error for missing output directory, got nil")
	}
}

// TestValidateBatchConfig tests validation with missing output extension
func TestValidateBatchConfig_MissingOutputExt(t *testing.T) {
	tempDir := t.TempDir()
	config := &BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputDir: "./output",
	}

	err := validateBatchConfig(config)
	if err == nil {
		t.Error("Expected error for missing output extension, got nil")
	}
}

// TestValidateBatchConfig tests validation with non-existent input directory
func TestValidateBatchConfig_NonExistentInputDir(t *testing.T) {
	config := &BatchConversionConfig{
		InputDir:  "./non_existent_directory",
		InputExt:  ".mp4",
		OutputDir: "./output",
		OutputExt: ".webm",
	}

	err := validateBatchConfig(config)
	if err == nil {
		t.Error("Expected error for non-existent input directory, got nil")
	}
}

// TestValidateBatchConfig tests validation with unsupported input extension
func TestValidateBatchConfig_UnsupportedInputExt(t *testing.T) {
	tempDir := t.TempDir()
	config := &BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".unsupported",
		OutputDir: "./output",
		OutputExt: ".webm",
	}

	err := validateBatchConfig(config)
	if err == nil {
		t.Error("Expected error for unsupported input extension, got nil")
	}
}

// TestValidateBatchConfig tests validation with unsupported output extension
func TestValidateBatchConfig_UnsupportedOutputExt(t *testing.T) {
	tempDir := t.TempDir()
	config := &BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputDir: "./output",
		OutputExt: ".unsupported",
	}

	err := validateBatchConfig(config)
	if err == nil {
		t.Error("Expected error for unsupported output extension, got nil")
	}
}

// TestScanFiles tests file scanning functionality
func TestBatchMovieConverterService_ScanFiles_Normal(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	subDir := filepath.Join(inputDir, "subdir")

	err := os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	// Create test files
	testFiles := []string{
		filepath.Join(inputDir, "video1.mp4"),
		filepath.Join(inputDir, "video2.mp4"),
		filepath.Join(inputDir, "document.txt"), // Should be ignored
		filepath.Join(subDir, "video3.mp4"),
	}

	for _, file := range testFiles {
		f, err := os.Create(file)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
		f.Close()
	}

	config := BatchConversionConfig{
		InputDir:  inputDir,
		InputExt:  ".mp4",
		OutputDir: filepath.Join(tempDir, "output"),
		OutputExt: ".webm",
		Recursive: true,
	}

	service := NewBatchMovieConverterService(config)
	files, err := service.scanFiles()

	if err != nil {
		t.Fatalf("scanFiles returned error: %v", err)
	}

	expectedCount := 3 // video1.mp4, video2.mp4, video3.mp4
	if len(files) != expectedCount {
		t.Errorf("Expected %d files, got %d", expectedCount, len(files))
	}

	// Check that all found files have .mp4 extension
	for _, file := range files {
		if filepath.Ext(file) != ".mp4" {
			t.Errorf("Found file with wrong extension: %s", file)
		}
	}
}

// TestScanFiles tests file scanning without recursion
func TestBatchMovieConverterService_ScanFiles_NonRecursive(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	subDir := filepath.Join(inputDir, "subdir")

	err := os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	// Create test files
	testFiles := []string{
		filepath.Join(inputDir, "video1.mp4"),
		filepath.Join(inputDir, "video2.mp4"),
		filepath.Join(subDir, "video3.mp4"), // Should be ignored due to non-recursive
	}

	for _, file := range testFiles {
		f, err := os.Create(file)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
		f.Close()
	}

	config := BatchConversionConfig{
		InputDir:  inputDir,
		InputExt:  ".mp4",
		OutputDir: filepath.Join(tempDir, "output"),
		OutputExt: ".webm",
		Recursive: false, // Non-recursive
	}

	service := NewBatchMovieConverterService(config)
	files, err := service.scanFiles()

	if err != nil {
		t.Fatalf("scanFiles returned error: %v", err)
	}

	expectedCount := 2 // Only video1.mp4, video2.mp4 (not video3.mp4 in subdir)
	if len(files) != expectedCount {
		t.Errorf("Expected %d files, got %d", expectedCount, len(files))
	}
}

// TestScanFiles tests file scanning with empty directory
func TestBatchMovieConverterService_ScanFiles_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")

	err := os.MkdirAll(inputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}

	config := BatchConversionConfig{
		InputDir:  inputDir,
		InputExt:  ".mp4",
		OutputDir: filepath.Join(tempDir, "output"),
		OutputExt: ".webm",
		Recursive: true,
	}

	service := NewBatchMovieConverterService(config)
	files, err := service.scanFiles()

	if err != nil {
		t.Fatalf("scanFiles returned error: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files, got %d", len(files))
	}
}

// TestBatchConvert tests batch conversion with empty directory
func TestBatchMovieConverterService_BatchConvert_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	err := os.MkdirAll(inputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}

	config := BatchConversionConfig{
		InputDir:       inputDir,
		InputExt:       ".mp4",
		OutputDir:      outputDir,
		OutputExt:      ".webm",
		Recursive:      true,
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewBatchMovieConverterService(config)
	result, err := service.batchConvert()

	if err != nil {
		t.Fatalf("batchConvert returned error: %v", err)
	}

	if result.TotalFiles != 0 {
		t.Errorf("Expected TotalFiles 0, got %d", result.TotalFiles)
	}

	if result.SuccessCount != 0 {
		t.Errorf("Expected SuccessCount 0, got %d", result.SuccessCount)
	}

	if result.FailureCount != 0 {
		t.Errorf("Expected FailureCount 0, got %d", result.FailureCount)
	}
}
