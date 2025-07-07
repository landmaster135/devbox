package usecases

import (
	"os"
	"path/filepath"
	"testing"
)

// #==============================================================#
// ##         Tests for batch convertSingleFile method           ##
// #==============================================================#
// TestBatchMovieConverterServiceConvertSingleFile tests the convertSingleFile method
type TestBatchMovieConverterServiceConvertSingleFile struct {
	t *testing.T
}

// NewTestBatchMovieConverterServiceConvertSingleFile creates a new test instance
func NewTestBatchMovieConverterServiceConvertSingleFile(t *testing.T) *TestBatchMovieConverterServiceConvertSingleFile {
	return &TestBatchMovieConverterServiceConvertSingleFile{t: t}
}

// TestBatchMovieConverterService_convertSingleFile_RelativePathError tests convertSingleFile with relative path error
func (ts *TestBatchMovieConverterServiceConvertSingleFile) TestBatchMovieConverterService_convertSingleFile_RelativePathError() {
	// Arrange
	tempDir := ts.t.TempDir()
	config := BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputDir: "/test/output",
		OutputExt: ".gif",
	}
	service := NewBatchMovieConverterService(config)

	// Use a file path that cannot be made relative to InputDir
	invalidFile := "/completely/different/path/test.mp4"

	// Act
	result := service.convertSingleFile(invalidFile)

	// Assert
	if result.Success {
		ts.t.Error("Expected convertSingleFile to fail for invalid relative path")
	}
	if result.Error == nil {
		ts.t.Error("Expected error for invalid relative path, got nil")
	}
	if result.InputFile != invalidFile {
		ts.t.Errorf("Expected InputFile %s, got %s", invalidFile, result.InputFile)
	}
}

// #==============================================================#
// ##         Tests for batch processing with output dir creation ##
// #==============================================================#
// TestBatchMovieConverterServiceOutputDirCreation tests output directory creation
type TestBatchMovieConverterServiceOutputDirCreation struct {
	t *testing.T
}

// NewTestBatchMovieConverterServiceOutputDirCreation creates a new test instance
func NewTestBatchMovieConverterServiceOutputDirCreation(t *testing.T) *TestBatchMovieConverterServiceOutputDirCreation {
	return &TestBatchMovieConverterServiceOutputDirCreation{t: t}
}

// TestBatchMovieConverterService_convertSingleFile_OutputDirCreation tests convertSingleFile with output directory creation
func (ts *TestBatchMovieConverterServiceOutputDirCreation) TestBatchMovieConverterService_convertSingleFile_OutputDirCreation() {
	// Arrange
	tempDir := ts.t.TempDir()
	outputDir := filepath.Join(tempDir, "output")

	// Create input file
	inputFile := filepath.Join(tempDir, "test.mp4")
	file, err := os.Create(inputFile)
	if err != nil {
		ts.t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputDir: outputDir,
		OutputExt: ".gif",
	}
	service := NewBatchMovieConverterService(config)

	// Act
	result := service.convertSingleFile(inputFile)

	// Assert
	// This will fail because ffmpeg is not available, but we can test the setup
	if result.InputFile != inputFile {
		ts.t.Errorf("Expected InputFile %s, got %s", inputFile, result.InputFile)
	}
	expectedOutput := filepath.Join(outputDir, "test.gif")
	if result.OutputFile != expectedOutput {
		ts.t.Errorf("Expected OutputFile %s, got %s", expectedOutput, result.OutputFile)
	}
	// Check that output directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		ts.t.Error("Expected output directory to be created")
	}
}

// #==============================================================#
// ##         Tests for additional generateOutputFile cases      ##
// #==============================================================#
// TestGenerateOutputFileAdditional tests additional cases for generateOutputFile
type TestGenerateOutputFileAdditional struct {
	t *testing.T
}

// NewTestGenerateOutputFileAdditional creates a new test instance
func NewTestGenerateOutputFileAdditional(t *testing.T) *TestGenerateOutputFileAdditional {
	return &TestGenerateOutputFileAdditional{t: t}
}

// TestGenerateOutputFile_EmptyString_Normal tests generateOutputFile with empty string
func (ts *TestGenerateOutputFileAdditional) TestGenerateOutputFile_EmptyString_Normal() {
	// Arrange
	inputFile := ""
	expected := "_converted"

	// Act
	result := generateOutputFile(inputFile)

	// Assert
	if result != expected {
		ts.t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestGenerateOutputFile_NoExtension_Normal tests generateOutputFile with no extension
func (ts *TestGenerateOutputFileAdditional) TestGenerateOutputFile_NoExtension_Normal() {
	// Arrange
	inputFile := "filename"
	expected := "filename_converted"

	// Act
	result := generateOutputFile(inputFile)

	// Assert
	if result != expected {
		ts.t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestGenerateOutputFile_ComplexPath_Normal tests generateOutputFile with complex path
func (ts *TestGenerateOutputFileAdditional) TestGenerateOutputFile_ComplexPath_Normal() {
	// Arrange
	inputFile := "/path/to/my video file.MP4"
	expected := "/path/to/my video file.gif"

	// Act
	result := generateOutputFile(inputFile)

	// Assert
	if result != expected {
		ts.t.Errorf("Expected %s, got %s", expected, result)
	}
}

// #==============================================================#
// ##         Tests for processBatchFiles method                 ##
// #==============================================================#
// TestUnifiedMovieConverterServiceProcessBatchFiles tests the processBatchFiles method
type TestUnifiedMovieConverterServiceProcessBatchFiles struct {
	t *testing.T
}

// NewTestUnifiedMovieConverterServiceProcessBatchFiles creates a new test instance
func NewTestUnifiedMovieConverterServiceProcessBatchFiles(t *testing.T) *TestUnifiedMovieConverterServiceProcessBatchFiles {
	return &TestUnifiedMovieConverterServiceProcessBatchFiles{t: t}
}

// TestUnifiedMovieConverterService_processBatchFiles_ValidConfig tests processBatchFiles with valid config
func (ts *TestUnifiedMovieConverterServiceProcessBatchFiles) TestUnifiedMovieConverterService_processBatchFiles_ValidConfig() {
	// Arrange
	tempDir := ts.t.TempDir()
	outputDir := filepath.Join(tempDir, "output")

	// Create test files
	testFiles := []string{"video1.mp4", "video2.mp4"}
	for _, fileName := range testFiles {
		file, err := os.Create(filepath.Join(tempDir, fileName))
		if err != nil {
			ts.t.Fatalf("Failed to create test file %s: %v", fileName, err)
		}
		file.Close()
	}

	batchConfig := &BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputDir: outputDir,
		OutputExt: ".gif",
		FPS:       30,
		Width:     320,
		Speed:     1.5,
	}
	service := NewUnifiedMovieConverterService(nil, batchConfig)

	// Act
	result := &UnifiedConversionResult{
		Mode: BatchMode,
	}
	finalResult, err := service.processBatchFiles(result)

	// Assert
	if finalResult == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	if finalResult.Mode != BatchMode {
		ts.t.Errorf("Expected BatchMode, got %d", finalResult.Mode)
	}
	// The conversion will fail due to ffmpeg not being available, but we test the setup
	if err != nil {
		ts.t.Fatalf("Failed to do processBatchFiles: %v", err)
	}
	if finalResult.BatchResult == nil {
		ts.t.Error("Expected BatchResult, got nil")
	}
	if finalResult.BatchResult.TotalFiles != 2 {
		ts.t.Errorf("Expected TotalFiles to be 2, got %d", finalResult.BatchResult.TotalFiles)
	}
}

// TestUnifiedMovieConverterService_processBatchFiles_InvalidInputDir tests processBatchFiles with invalid input directory
func (ts *TestUnifiedMovieConverterServiceProcessBatchFiles) TestUnifiedMovieConverterService_processBatchFiles_InvalidInputDir() {
	// Arrange
	batchConfig := &BatchConversionConfig{
		InputDir:  "", // Empty input directory
		InputExt:  ".mp4",
		OutputDir: "/test/output",
		OutputExt: ".gif",
	}
	service := NewUnifiedMovieConverterService(nil, batchConfig)

	// Act
	result := &UnifiedConversionResult{
		Mode: BatchMode,
	}
	finalResult, err := service.processBatchFiles(result)

	// Assert
	if err == nil {
		ts.t.Error("Expected error for invalid input directory, got nil")
	}
	if finalResult == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	if finalResult.Success {
		ts.t.Error("Expected Success to be false for invalid config")
	}
	expectedMsg := "バッチ設定エラー: 入力ディレクトリが指定されていません"
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestUnifiedMovieConverterService_processBatchFiles_UnsupportedExtension tests processBatchFiles with unsupported extension
func (ts *TestUnifiedMovieConverterServiceProcessBatchFiles) TestUnifiedMovieConverterService_processBatchFiles_UnsupportedExtension() {
	// Arrange
	tempDir := ts.t.TempDir()
	batchConfig := &BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".txt", // Unsupported extension
		OutputDir: "/test/output",
		OutputExt: ".gif",
	}
	service := NewUnifiedMovieConverterService(nil, batchConfig)

	// Act
	result := &UnifiedConversionResult{
		Mode: BatchMode,
	}
	finalResult, err := service.processBatchFiles(result)

	// Assert
	if err == nil {
		ts.t.Error("Expected error for unsupported extension, got nil")
	}
	if finalResult == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	if finalResult.Success {
		ts.t.Error("Expected Success to be false for unsupported extension")
	}
	expectedMsg := "バッチ設定エラー: サポートされていない入力拡張子: .txt"
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// #==============================================================#
// ##         Tests for processSingleFile success cases          ##
// #==============================================================#
// TestUnifiedMovieConverterServiceProcessSingleFileSuccess tests processSingleFile success cases
type TestUnifiedMovieConverterServiceProcessSingleFileSuccess struct {
	t *testing.T
}

// NewTestUnifiedMovieConverterServiceProcessSingleFileSuccess creates a new test instance
func NewTestUnifiedMovieConverterServiceProcessSingleFileSuccess(t *testing.T) *TestUnifiedMovieConverterServiceProcessSingleFileSuccess {
	return &TestUnifiedMovieConverterServiceProcessSingleFileSuccess{t: t}
}

// TestUnifiedMovieConverterService_processSingleFile_ValidFileWithOutputGeneration tests processSingleFile with output generation
func (ts *TestUnifiedMovieConverterServiceProcessSingleFileSuccess) TestUnifiedMovieConverterService_processSingleFile_ValidFileWithOutputGeneration() {
	// Arrange
	tempDir := ts.t.TempDir()
	testFile := filepath.Join(tempDir, "test.mp4")
	file, err := os.Create(testFile)
	if err != nil {
		ts.t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	singleConfig := &ConversionConfig{
		InputFile:  testFile,
		OutputFile: "", // Empty output file should trigger auto-generation
		FPS:        30,
		Width:      320,
		Speed:      1.5,
	}
	service := NewUnifiedMovieConverterService(singleConfig, nil)

	// Act
	result := &UnifiedConversionResult{
		Mode: SingleFileMode,
	}
	finalResult, err := service.processSingleFile(result)

	// Assert
	if finalResult == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	if finalResult.Mode != SingleFileMode {
		ts.t.Errorf("Expected SingleFileMode, got %d", finalResult.Mode)
	}
	// The conversion will fail due to ffmpeg not being available, but we test the setup
	if err == nil {
		ts.t.Error("Expected error due to ffmpeg not being available")
	}
	if finalResult.SingleResult == nil {
		ts.t.Error("Expected SingleResult, got nil")
	}
	// Check that output file was auto-generated
	expectedOutput := filepath.Join(tempDir, "test.gif")
	if finalResult.SingleResult.OutputFile != expectedOutput {
		ts.t.Errorf("Expected auto-generated output file %s, got %s", expectedOutput, finalResult.SingleResult.OutputFile)
	}
}

// TestUnifiedMovieConverterService_processSingleFile_GIFToMP4Conversion tests GIF to MP4 conversion
func (ts *TestUnifiedMovieConverterServiceProcessSingleFileSuccess) TestUnifiedMovieConverterService_processSingleFile_GIFToMP4Conversion() {
	// Arrange
	tempDir := ts.t.TempDir()
	testFile := filepath.Join(tempDir, "test.gif")
	file, err := os.Create(testFile)
	if err != nil {
		ts.t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	singleConfig := &ConversionConfig{
		InputFile:  testFile,
		OutputFile: filepath.Join(tempDir, "output.mp4"),
		FPS:        15,
	}
	service := NewUnifiedMovieConverterService(singleConfig, nil)

	// Act
	result := &UnifiedConversionResult{
		Mode: SingleFileMode,
	}
	finalResult, err := service.processSingleFile(result)

	// Assert
	if finalResult == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	// The conversion will fail due to ffmpeg not being available, but we test the setup
	if err == nil {
		ts.t.Error("Expected error due to ffmpeg not being available")
	}
	if finalResult.SingleResult == nil {
		ts.t.Error("Expected SingleResult, got nil")
	}
	if finalResult.SingleResult.InputFile != testFile {
		ts.t.Errorf("Expected InputFile %s, got %s", testFile, finalResult.SingleResult.InputFile)
	}
}

// Standard Go test functions for unified edge tests

func TestBatchMovieConverterServiceConvertSingleFileMethod(t *testing.T) {
	testService := NewTestBatchMovieConverterServiceConvertSingleFile(t)
	testService.TestBatchMovieConverterService_convertSingleFile_RelativePathError()
}

func TestBatchMovieConverterServiceOutputDirCreationMethod(t *testing.T) {
	testService := NewTestBatchMovieConverterServiceOutputDirCreation(t)
	testService.TestBatchMovieConverterService_convertSingleFile_OutputDirCreation()
}

func TestGenerateOutputFileAdditionalMethod(t *testing.T) {
	testService := NewTestGenerateOutputFileAdditional(t)
	testService.TestGenerateOutputFile_EmptyString_Normal()
	testService.TestGenerateOutputFile_NoExtension_Normal()
	testService.TestGenerateOutputFile_ComplexPath_Normal()
}

func TestUnifiedMovieConverterServiceProcessBatchFilesMethod(t *testing.T) {
	testService := NewTestUnifiedMovieConverterServiceProcessBatchFiles(t)
	testService.TestUnifiedMovieConverterService_processBatchFiles_ValidConfig()
	testService.TestUnifiedMovieConverterService_processBatchFiles_InvalidInputDir()
	testService.TestUnifiedMovieConverterService_processBatchFiles_UnsupportedExtension()
}

func TestUnifiedMovieConverterServiceProcessSingleFileSuccessMethod(t *testing.T) {
	testService := NewTestUnifiedMovieConverterServiceProcessSingleFileSuccess(t)
	testService.TestUnifiedMovieConverterService_processSingleFile_ValidFileWithOutputGeneration()
	testService.TestUnifiedMovieConverterService_processSingleFile_GIFToMP4Conversion()
}
