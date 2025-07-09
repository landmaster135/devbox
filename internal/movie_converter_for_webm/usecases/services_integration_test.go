package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFullWorkflow_SingleFileConversion tests complete single file conversion workflow
func TestFullWorkflow_SingleFileConversion_Normal(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "input.mp4")
	outputFile := filepath.Join(tempDir, "output.webm")

	// Create a dummy input file
	file, err := os.Create(inputFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	// Test complete workflow
	singleConfig := &ConversionConfig{
		InputFile:      inputFile,
		OutputFile:     outputFile,
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewUnifiedMovieConverterService(singleConfig, nil)
	result, err := service.ProcessConversion()

	// Verify result structure
	if result == nil {
		t.Fatal("ProcessConversion returned nil result")
	}

	if result.Mode != SingleFileMode {
		t.Errorf("Expected SingleFileMode, got %v", result.Mode)
	}

	if result.SingleResult == nil {
		t.Error("Expected SingleResult to be set")
	}

	if result.BatchResult != nil {
		t.Error("Expected BatchResult to be nil for single file mode")
	}

	// We expect an error due to missing ffmpeg, but structure should be correct
	if err == nil {
		t.Error("Expected error due to missing ffmpeg, got nil")
	}

	if result.Success {
		t.Error("Expected Success to be false due to ffmpeg error")
	}
}

// TestFullWorkflow_BatchConversion tests complete batch conversion workflow
func TestFullWorkflow_BatchConversion_Normal(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	err := os.MkdirAll(inputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}

	// Create test files
	testFiles := []string{"video1.mp4", "video2.mp4", "video3.mp4"}
	for _, fileName := range testFiles {
		file, err := os.Create(filepath.Join(inputDir, fileName))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", fileName, err)
		}
		file.Close()
	}

	// Test complete batch workflow
	batchConfig := &BatchConversionConfig{
		InputDir:       inputDir,
		InputExt:       ".mp4",
		OutputDir:      outputDir,
		OutputExt:      ".webm",
		Recursive:      false,
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewUnifiedMovieConverterService(nil, batchConfig)
	result, err := service.ProcessConversion()

	if err != nil {
		t.Fatalf("Failed to do ProcessConversion: %v", err)
	}

	// Verify result structure
	if result == nil {
		t.Fatal("ProcessConversion returned nil result")
	}

	if result.Mode != BatchMode {
		t.Errorf("Expected BatchMode, got %v", result.Mode)
	}

	if result.BatchResult == nil {
		t.Error("Expected BatchResult to be set")
	}

	if result.SingleResult != nil {
		t.Error("Expected SingleResult to be nil for batch mode")
	}

	// Check batch result details
	if result.BatchResult.TotalFiles != 3 {
		t.Errorf("Expected TotalFiles 3, got %d", result.BatchResult.TotalFiles)
	}

	if result.BatchResult.FailureCount != 3 {
		t.Errorf("Expected FailureCount 3 (due to missing ffmpeg), got %d", result.BatchResult.FailureCount)
	}

	if result.BatchResult.SuccessCount != 0 {
		t.Errorf("Expected SuccessCount 0 (due to missing ffmpeg), got %d", result.BatchResult.SuccessCount)
	}

	if len(result.BatchResult.Results) != 3 {
		t.Errorf("Expected 3 conversion results, got %d", len(result.BatchResult.Results))
	}

	if len(result.BatchResult.FailedFiles) != 3 {
		t.Errorf("Expected 3 failed files, got %d", len(result.BatchResult.FailedFiles))
	}
}

// TestFullWorkflow_AutoOutputGeneration tests automatic output file generation
func TestFullWorkflow_AutoOutputGeneration_Normal(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "test.mp4")

	// Create a dummy input file
	file, err := os.Create(inputFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	// Test with no output file specified
	singleConfig := &ConversionConfig{
		InputFile: inputFile,
		// OutputFile not specified - should be auto-generated
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewUnifiedMovieConverterService(singleConfig, nil)
	result, err := service.ProcessConversion()

	// We expect an error due to invalid dummy file (ffmpeg exit status 183)
	if err == nil {
		t.Error("Expected error due to invalid dummy file, got nil")
	}

	if result == nil {
		t.Fatal("ProcessConversion returned nil result")
	}

	if result.Success {
		t.Error("Expected Success to be false due to ffmpeg error")
	}

	// Verify auto-generated output file
	expectedOutputFile := filepath.Join(tempDir, "test.webm")
	if service.config.SingleConfig.OutputFile != expectedOutputFile {
		t.Errorf("Expected auto-generated OutputFile %s, got %s", expectedOutputFile, service.config.SingleConfig.OutputFile)
	}

	if result.SingleResult.OutputFile != expectedOutputFile {
		t.Errorf("Expected result OutputFile %s, got %s", expectedOutputFile, result.SingleResult.OutputFile)
	}
}

// TestFullWorkflow_RecursiveBatchConversion tests recursive batch conversion
func TestFullWorkflow_RecursiveBatchConversion_Normal(t *testing.T) {
	// Create temporary directory structure with subdirectories
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	subDir := filepath.Join(inputDir, "subdir")
	outputDir := filepath.Join(tempDir, "output")

	err := os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create test files in main directory and subdirectory
	mainFiles := []string{"video1.mp4", "video2.mp4"}
	subFiles := []string{"video3.mp4", "video4.mp4"}

	for _, fileName := range mainFiles {
		file, err := os.Create(filepath.Join(inputDir, fileName))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", fileName, err)
		}
		file.Close()
	}

	for _, fileName := range subFiles {
		file, err := os.Create(filepath.Join(subDir, fileName))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", fileName, err)
		}
		file.Close()
	}

	// Test recursive batch conversion
	batchConfig := &BatchConversionConfig{
		InputDir:       inputDir,
		InputExt:       ".mp4",
		OutputDir:      outputDir,
		OutputExt:      ".webm",
		Recursive:      true, // Enable recursive processing
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewUnifiedMovieConverterService(nil, batchConfig)
	result, err := service.ProcessConversion()

	if err != nil {
		t.Fatalf("Failed to do ProcessConversion: %v", err)
	}

	// Should find all 4 files (2 in main dir + 2 in subdir)
	if result.BatchResult.TotalFiles != 4 {
		t.Errorf("Expected TotalFiles 4 (recursive), got %d", result.BatchResult.TotalFiles)
	}

	// Test non-recursive for comparison
	batchConfig.Recursive = false
	service2 := NewUnifiedMovieConverterService(nil, batchConfig)
	result2, err2 := service2.ProcessConversion()

	if err2 != nil {
		t.Errorf("Failed to do ProcessConversion: %v", err2)
	}

	// Should find only 2 files (main dir only)
	if result2.BatchResult.TotalFiles != 2 {
		t.Errorf("Expected TotalFiles 2 (non-recursive), got %d", result2.BatchResult.TotalFiles)
	}
}

// TestFullWorkflow_ConfigurationInheritance tests configuration inheritance in batch processing
func TestFullWorkflow_ConfigurationInheritance_Normal(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	err := os.MkdirAll(inputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(inputDir, "test.mp4")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	// Test with custom configuration
	batchConfig := &BatchConversionConfig{
		InputDir:       inputDir,
		InputExt:       ".mp4",
		OutputDir:      outputDir,
		OutputExt:      ".webm",
		Recursive:      false,
		VideoBitrate:   "2M",
		AudioBitrate:   "192k",
		AudioCodec:     "vorbis",
		ConversionMode: "cbr",
		CRF:            25,
		VideoQuality:   90,
	}

	service := NewUnifiedMovieConverterService(nil, batchConfig)
	result, err := service.ProcessConversion()

	if err != nil {
		t.Fatalf("Failed to do ProcessConversion: %v", err)
	}

	// Verify that configuration was inherited in conversion results
	if len(result.BatchResult.Results) > 0 {
		convResult := result.BatchResult.Results[0]
		// The individual conversion should have inherited the batch config
		// (We can't verify the exact config values since they're internal to the conversion process)
		if convResult.InputFile == "" {
			t.Error("Expected InputFile to be set in conversion result")
		}
		if convResult.OutputFile == "" {
			t.Error("Expected OutputFile to be set in conversion result")
		}
	}
}

// TestFullWorkflow_ErrorPropagation tests error propagation through the system
func TestFullWorkflow_ErrorPropagation_Normal(t *testing.T) {
	// Test with invalid input file
	singleConfig := &ConversionConfig{
		InputFile:  "non_existent_file.mp4",
		OutputFile: "output.webm",
	}

	service := NewUnifiedMovieConverterService(singleConfig, nil)
	result, err := service.ProcessConversion()

	// Verify error propagation
	if err == nil {
		t.Error("Expected error for non-existent input file, got nil")
	}

	if result.Success {
		t.Error("Expected Success to be false")
	}

	if result.Error == nil {
		t.Error("Expected Error to be set in result")
	}

	if result.SingleResult == nil {
		t.Error("Expected SingleResult to be set even with error")
		return
	}

	if result.SingleResult.Success {
		t.Error("Expected SingleResult.Success to be false")
	}

	if result.SingleResult.Error == nil {
		t.Error("Expected SingleResult.Error to be set")
	}
}

// TestFullWorkflow_MixedFileTypes tests handling of mixed file types in batch processing
func TestFullWorkflow_MixedFileTypes_Normal(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	err := os.MkdirAll(inputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}

	// Create mixed file types
	testFiles := []string{
		"video1.mp4",   // Should be processed
		"video2.mp4",   // Should be processed
		"document.txt", // Should be ignored
		"image.jpg",    // Should be ignored
		"video3.avi",   // Should be ignored (different extension)
	}

	for _, fileName := range testFiles {
		file, err := os.Create(filepath.Join(inputDir, fileName))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", fileName, err)
		}
		file.Close()
	}

	// Test batch conversion with .mp4 filter
	batchConfig := &BatchConversionConfig{
		InputDir:  inputDir,
		InputExt:  ".mp4",
		OutputDir: outputDir,
		OutputExt: ".webm",
		Recursive: false,
	}

	service := NewUnifiedMovieConverterService(nil, batchConfig)
	result, err := service.ProcessConversion()

	if err != nil {
		t.Fatalf("Failed to do ProcessConversion: %v", err)
	}

	// Should only process .mp4 files
	if result.BatchResult.TotalFiles != 2 {
		t.Errorf("Expected TotalFiles 2 (.mp4 files only), got %d", result.BatchResult.TotalFiles)
	}

	// Verify that only .mp4 files were processed
	for _, convResult := range result.BatchResult.Results {
		if !strings.HasSuffix(convResult.InputFile, ".mp4") {
			t.Errorf("Expected only .mp4 files to be processed, found: %s", convResult.InputFile)
		}
	}
}

// TestFullWorkflow_OutputDirectoryCreation tests automatic output directory creation
func TestFullWorkflow_OutputDirectoryCreation_Normal(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output", "nested", "path") // Deep nested path

	err := os.MkdirAll(inputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(inputDir, "test.mp4")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	// Test batch conversion with non-existent output directory
	batchConfig := &BatchConversionConfig{
		InputDir:  inputDir,
		InputExt:  ".mp4",
		OutputDir: outputDir, // This directory doesn't exist yet
		OutputExt: ".webm",
		Recursive: false,
	}

	service := NewUnifiedMovieConverterService(nil, batchConfig)
	result, err := service.ProcessConversion()

	if err != nil {
		t.Fatalf("Failed to do ProcessConversion: %v", err)
	}

	// Verify that output directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Error("Expected output directory to be created, but it doesn't exist")
	}

	// The conversion itself will fail due to missing ffmpeg, but directory creation should succeed
	if result.BatchResult.TotalFiles != 1 {
		t.Errorf("Expected TotalFiles 1, got %d", result.BatchResult.TotalFiles)
	}
}

// TestFullWorkflow_EmptyDirectory tests handling of empty directories
func TestFullWorkflow_EmptyDirectory_Normal(t *testing.T) {
	// Create empty input directory
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	err := os.MkdirAll(inputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}

	// Test batch conversion with empty directory
	batchConfig := &BatchConversionConfig{
		InputDir:  inputDir,
		InputExt:  ".mp4",
		OutputDir: outputDir,
		OutputExt: ".webm",
		Recursive: true,
	}

	service := NewUnifiedMovieConverterService(nil, batchConfig)
	result, err := service.ProcessConversion()

	// Should succeed with no files to process
	if err != nil {
		t.Errorf("Expected no error for empty directory, got %v", err)
	}

	if !result.Success {
		t.Error("Expected Success to be true for empty directory")
	}

	if result.BatchResult.TotalFiles != 0 {
		t.Errorf("Expected TotalFiles 0, got %d", result.BatchResult.TotalFiles)
	}

	if result.BatchResult.SuccessCount != 0 {
		t.Errorf("Expected SuccessCount 0, got %d", result.BatchResult.SuccessCount)
	}

	if result.BatchResult.FailureCount != 0 {
		t.Errorf("Expected FailureCount 0, got %d", result.BatchResult.FailureCount)
	}
}
