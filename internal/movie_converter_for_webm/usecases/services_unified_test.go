package usecases

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewUnifiedMovieConverterService tests the creation of UnifiedMovieConverterService
func TestNewUnifiedMovieConverterService_Normal(t *testing.T) {
	singleConfig := &ConversionConfig{
		InputFile:      "test.mp4",
		OutputFile:     "test.webm",
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewUnifiedMovieConverterService(singleConfig, nil)

	if service == nil {
		t.Fatal("NewUnifiedMovieConverterService returned nil")
	}

	if service.config.Mode != SingleFileMode {
		t.Errorf("Expected SingleFileMode, got %v", service.config.Mode)
	}

	if service.config.SingleConfig == nil {
		t.Error("Expected SingleConfig to be set, got nil")
	}

	if service.config.BatchConfig != nil {
		t.Error("Expected BatchConfig to be nil, got non-nil")
	}
}

// TestNewUnifiedMovieConverterService tests batch mode detection
func TestNewUnifiedMovieConverterService_BatchMode(t *testing.T) {
	batchConfig := &BatchConversionConfig{
		InputDir:       "./input",
		InputExt:       ".mp4",
		OutputDir:      "./output",
		OutputExt:      ".webm",
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewUnifiedMovieConverterService(nil, batchConfig)

	if service == nil {
		t.Fatal("NewUnifiedMovieConverterService returned nil")
	}

	if service.config.Mode != BatchMode {
		t.Errorf("Expected BatchMode, got %v", service.config.Mode)
	}

	if service.config.SingleConfig != nil {
		t.Error("Expected SingleConfig to be nil, got non-nil")
	}

	if service.config.BatchConfig == nil {
		t.Error("Expected BatchConfig to be set, got nil")
	}
}

// TestProcessSingleFile tests single file processing with valid configuration
func TestUnifiedMovieConverterService_ProcessSingleFile_Normal(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.mp4")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	singleConfig := &ConversionConfig{
		InputFile:      testFile,
		OutputFile:     filepath.Join(tempDir, "test.webm"),
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewUnifiedMovieConverterService(singleConfig, nil)
	result := &UnifiedConversionResult{
		Mode: SingleFileMode,
	}

	// Note: This test will fail during actual conversion due to invalid dummy file
	// but we can test the validation and setup logic
	_, err = service.processSingleFile(result)

	// We expect an error due to invalid dummy file (ffmpeg exit status 183)
	if err == nil {
		t.Error("Expected error due to invalid dummy file, got nil")
	}

	if result.Success {
		t.Error("Expected Success to be false due to ffmpeg error")
	}

	if result.SingleResult == nil {
		t.Error("Expected SingleResult to be set")
	}

	if result.SingleResult.InputFile != testFile {
		t.Errorf("Expected InputFile %s, got %s", testFile, result.SingleResult.InputFile)
	}
}

// TestProcessSingleFile tests single file processing with missing configuration
func TestUnifiedMovieConverterService_ProcessSingleFile_MissingConfig(t *testing.T) {
	service := NewUnifiedMovieConverterService(nil, nil)
	service.config.Mode = SingleFileMode // Force single file mode
	result := &UnifiedConversionResult{
		Mode: SingleFileMode,
	}

	_, err := service.processSingleFile(result)

	if err == nil {
		t.Error("Expected error for missing single config, got nil")
	}

	if result.Success {
		t.Error("Expected Success to be false")
	}
}

// TestProcessSingleFile tests single file processing with invalid configuration
func TestUnifiedMovieConverterService_ProcessSingleFile_InvalidConfig(t *testing.T) {
	singleConfig := &ConversionConfig{
		// Missing InputFile
		OutputFile: "test.webm",
	}

	service := NewUnifiedMovieConverterService(singleConfig, nil)
	result := &UnifiedConversionResult{
		Mode: SingleFileMode,
	}

	_, err := service.processSingleFile(result)

	if err == nil {
		t.Error("Expected error for invalid config, got nil")
	}

	if result.Success {
		t.Error("Expected Success to be false")
	}
}

// TestProcessSingleFile tests automatic output file generation
func TestUnifiedMovieConverterService_ProcessSingleFile_AutoOutputFile(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.mp4")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	singleConfig := &ConversionConfig{
		InputFile: testFile,
		// OutputFile is not set, should be auto-generated
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewUnifiedMovieConverterService(singleConfig, nil)
	result := &UnifiedConversionResult{
		Mode: SingleFileMode,
	}

	// Note: This test will fail during actual conversion due to invalid dummy file
	// but we can test the auto-generation logic
	_, err = service.processSingleFile(result)

	// We expect an error due to invalid dummy file (ffmpeg exit status 183)
	if err == nil {
		t.Error("Expected error due to invalid dummy file, got nil")
	}

	if result.Success {
		t.Error("Expected Success to be false due to ffmpeg error")
	}

	expectedOutputFile := filepath.Join(tempDir, "test.webm")
	if service.config.SingleConfig.OutputFile != expectedOutputFile {
		t.Errorf("Expected auto-generated OutputFile %s, got %s", expectedOutputFile, service.config.SingleConfig.OutputFile)
	}
}

// TestProcessBatchFiles tests batch processing with missing configuration
func TestUnifiedMovieConverterService_ProcessBatchFiles_MissingConfig(t *testing.T) {
	service := NewUnifiedMovieConverterService(nil, nil)
	service.config.Mode = BatchMode // Force batch mode
	result := &UnifiedConversionResult{
		Mode: BatchMode,
	}

	_, err := service.processBatchFiles(result)

	if err == nil {
		t.Error("Expected error for missing batch config, got nil")
	}

	if result.Success {
		t.Error("Expected Success to be false")
	}
}

// TestProcessBatchFiles tests batch processing with invalid configuration
func TestUnifiedMovieConverterService_ProcessBatchFiles_InvalidConfig(t *testing.T) {
	batchConfig := &BatchConversionConfig{
		// Missing required fields
		InputDir: "./input",
		// Missing InputExt, OutputDir, OutputExt
	}

	service := NewUnifiedMovieConverterService(nil, batchConfig)
	result := &UnifiedConversionResult{
		Mode: BatchMode,
	}

	_, err := service.processBatchFiles(result)

	if err == nil {
		t.Error("Expected error for invalid batch config, got nil")
	}

	if result.Success {
		t.Error("Expected Success to be false")
	}
}

// TestProcessConversion tests the main ProcessConversion method
func TestUnifiedMovieConverterService_ProcessConversion_SingleFileMode(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.mp4")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	singleConfig := &ConversionConfig{
		InputFile:      testFile,
		OutputFile:     filepath.Join(tempDir, "test.webm"),
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

	if result.Mode != SingleFileMode {
		t.Errorf("Expected SingleFileMode, got %v", result.Mode)
	}

	if result.SingleResult == nil {
		t.Error("Expected SingleResult to be set")
	}
}

// TestProcessConversion tests batch mode
func TestUnifiedMovieConverterService_ProcessConversion_BatchMode(t *testing.T) {
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	err := os.MkdirAll(inputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}

	batchConfig := &BatchConversionConfig{
		InputDir:       inputDir,
		InputExt:       ".mp4",
		OutputDir:      outputDir,
		OutputExt:      ".webm",
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewUnifiedMovieConverterService(nil, batchConfig)
	result, err := service.ProcessConversion()

	if result == nil {
		t.Fatal("ProcessConversion returned nil result")
	}

	if result.Mode != BatchMode {
		t.Errorf("Expected BatchMode, got %v", result.Mode)
	}

	if result.BatchResult == nil {
		t.Error("Expected BatchResult to be set")
	}

	// For empty directory, we should get successful result with 0 files
	if err != nil {
		t.Errorf("Expected no error for empty directory, got %v", err)
	}

	if !result.Success {
		t.Error("Expected Success to be true for empty directory")
	}

	if result.BatchResult.TotalFiles != 0 {
		t.Errorf("Expected TotalFiles 0, got %d", result.BatchResult.TotalFiles)
	}
}

// TestProcessConversion tests unknown mode
func TestUnifiedMovieConverterService_ProcessConversion_UnknownMode(t *testing.T) {
	service := NewUnifiedMovieConverterService(nil, nil)
	service.config.Mode = ProcessingMode(999) // Invalid mode

	result, err := service.ProcessConversion()

	if err == nil {
		t.Error("Expected error for unknown mode, got nil")
	}

	if result == nil {
		t.Fatal("ProcessConversion returned nil result")
	}

	if result.Success {
		t.Error("Expected Success to be false")
	}
}
