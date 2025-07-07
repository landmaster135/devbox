package usecases

import (
	"os"
	"path/filepath"
	"testing"
)

// #==============================================================#
// ##         Tests for Unified Service                          ##
// #==============================================================#
// TestUnifiedMovieConverterService tests the UnifiedMovieConverterService struct
type TestUnifiedMovieConverterService struct {
	t *testing.T
}

// NewTestUnifiedMovieConverterService creates a new test instance
func NewTestUnifiedMovieConverterService(t *testing.T) *TestUnifiedMovieConverterService {
	return &TestUnifiedMovieConverterService{t: t}
}

// TestNewUnifiedMovieConverterService_SingleFileMode_Normal tests creation with single file mode
func (ts *TestUnifiedMovieConverterService) TestNewUnifiedMovieConverterService_SingleFileMode_Normal() {
	// Arrange
	singleConfig := &ConversionConfig{
		InputFile:   "test.mp4",
		OutputFile:  "test.gif",
		FPS:         30,
		Width:       320,
		Speed:       1.5,
		Loop:        0,
		UseItsScale: true,
	}

	// Act
	service := NewUnifiedMovieConverterService(singleConfig, nil)

	// Assert
	if service == nil {
		ts.t.Error("NewUnifiedMovieConverterService should not return nil")
		return
	}
	if service.config.Mode != SingleFileMode {
		ts.t.Errorf("Expected SingleFileMode, got %d", service.config.Mode)
	}
	if service.config.SingleConfig == nil {
		ts.t.Error("SingleConfig should not be nil")
	}
	if service.config.BatchConfig != nil {
		ts.t.Error("BatchConfig should be nil for single file mode")
	}
}

// TestNewUnifiedMovieConverterService_BatchMode_Normal tests creation with batch mode
func (ts *TestUnifiedMovieConverterService) TestNewUnifiedMovieConverterService_BatchMode_Normal() {
	// Arrange
	batchConfig := &BatchConversionConfig{
		InputDir:    "/test/input",
		InputExt:    ".mp4",
		OutputDir:   "/test/output",
		OutputExt:   ".gif",
		Recursive:   true,
		FPS:         30,
		Width:       320,
		Speed:       1.5,
		Loop:        0,
		UseItsScale: true,
	}

	// Act
	service := NewUnifiedMovieConverterService(nil, batchConfig)

	// Assert
	if service == nil {
		ts.t.Error("NewUnifiedMovieConverterService should not return nil")
		return
	}
	if service.config.Mode != BatchMode {
		ts.t.Errorf("Expected BatchMode, got %d", service.config.Mode)
	}
	if service.config.BatchConfig == nil {
		ts.t.Error("BatchConfig should not be nil")
	}
	if service.config.SingleConfig != nil {
		ts.t.Error("SingleConfig should be nil for batch mode")
	}
}

// TestNewUnifiedMovieConverterService_DefaultToSingleMode_Normal tests default to single mode
func (ts *TestUnifiedMovieConverterService) TestNewUnifiedMovieConverterService_DefaultToSingleMode_Normal() {
	// Arrange
	singleConfig := &ConversionConfig{
		InputFile:  "test.mp4",
		OutputFile: "test.gif",
	}
	emptyBatchConfig := &BatchConversionConfig{} // Empty batch config

	// Act
	service := NewUnifiedMovieConverterService(singleConfig, emptyBatchConfig)

	// Assert
	if service == nil {
		ts.t.Error("NewUnifiedMovieConverterService should not return nil")
		return
	}
	if service.config.Mode != SingleFileMode {
		ts.t.Errorf("Expected SingleFileMode when batch config is empty, got %d", service.config.Mode)
	}
}

// TestProcessingMode tests the ProcessingMode constants
type TestProcessingMode struct {
	t *testing.T
}

// NewTestProcessingMode creates a new test instance
func NewTestProcessingMode(t *testing.T) *TestProcessingMode {
	return &TestProcessingMode{t: t}
}

// TestProcessingMode_Constants_Normal tests ProcessingMode constants
func (ts *TestProcessingMode) TestProcessingMode_Constants_Normal() {
	// Assert
	if SingleFileMode != 0 {
		ts.t.Errorf("Expected SingleFileMode to be 0, got %d", SingleFileMode)
	}
	if BatchMode != 1 {
		ts.t.Errorf("Expected BatchMode to be 1, got %d", BatchMode)
	}
}

// Standard Go test functions for unified service

func TestUnifiedMovieConverterServiceCreation(t *testing.T) {
	testService := NewTestUnifiedMovieConverterService(t)
	testService.TestNewUnifiedMovieConverterService_SingleFileMode_Normal()
	testService.TestNewUnifiedMovieConverterService_BatchMode_Normal()
	testService.TestNewUnifiedMovieConverterService_DefaultToSingleMode_Normal()
}

func TestProcessingModeConstants(t *testing.T) {
	testService := NewTestProcessingMode(t)
	testService.TestProcessingMode_Constants_Normal()
}

// #==============================================================#
// ##         Tests for Unified Service Edge Cases               ##
// #==============================================================#
// TestUnifiedMovieConverterServiceProcessConversion tests the ProcessConversion method
type TestUnifiedMovieConverterServiceProcessConversion struct {
	t *testing.T
}

// NewTestUnifiedMovieConverterServiceProcessConversion creates a new test instance
func NewTestUnifiedMovieConverterServiceProcessConversion(t *testing.T) *TestUnifiedMovieConverterServiceProcessConversion {
	return &TestUnifiedMovieConverterServiceProcessConversion{t: t}
}

// TestUnifiedMovieConverterService_ProcessConversion_SingleFileMode_InvalidConfig tests single file mode with invalid config
func (ts *TestUnifiedMovieConverterServiceProcessConversion) TestUnifiedMovieConverterService_ProcessConversion_SingleFileMode_InvalidConfig() {
	// Arrange
	singleConfig := &ConversionConfig{
		InputFile:  "", // Empty input file should cause error
		OutputFile: "test.gif",
	}
	service := NewUnifiedMovieConverterService(singleConfig, nil)

	// Act
	result, err := service.ProcessConversion()

	// Assert
	if err == nil {
		ts.t.Error("Expected error for invalid single file config, got nil")
	}
	if result == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	if result.Success {
		ts.t.Error("Expected Success to be false for invalid config")
	}
	if result.Mode != SingleFileMode {
		ts.t.Errorf("Expected SingleFileMode, got %d", result.Mode)
	}
}

// TestUnifiedMovieConverterService_ProcessConversion_BatchMode_InvalidConfig tests batch mode with invalid config
func (ts *TestUnifiedMovieConverterServiceProcessConversion) TestUnifiedMovieConverterService_ProcessConversion_BatchMode_InvalidConfig() {
	// Arrange
	batchConfig := &BatchConversionConfig{
		InputDir:  "", // Empty input directory should cause error
		InputExt:  ".mp4",
		OutputDir: "/test/output",
		OutputExt: ".gif",
	}
	service := NewUnifiedMovieConverterService(nil, batchConfig)

	// Act
	result, err := service.ProcessConversion()

	// Assert
	if err == nil {
		ts.t.Error("Expected error for invalid batch config, got nil")
	}
	if result == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	if result.Success {
		ts.t.Error("Expected Success to be false for invalid config")
	}
	if result.Mode != BatchMode {
		ts.t.Errorf("Expected BatchMode, got %d", result.Mode)
	}
}

// TestUnifiedMovieConverterService_ProcessConversion_NilSingleConfig tests single file mode with nil config
func (ts *TestUnifiedMovieConverterServiceProcessConversion) TestUnifiedMovieConverterService_ProcessConversion_NilSingleConfig() {
	// Arrange
	service := NewUnifiedMovieConverterService(nil, nil)

	// Act
	result, err := service.ProcessConversion()

	// Assert
	if err == nil {
		ts.t.Error("Expected error for nil single config, got nil")
	}
	if result == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	if result.Success {
		ts.t.Error("Expected Success to be false for nil config")
	}
	expectedMsg := "単一ファイル処理の設定が指定されていません"
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestUnifiedMovieConverterService_ProcessConversion_NilBatchConfig tests batch mode with nil config
func (ts *TestUnifiedMovieConverterServiceProcessConversion) TestUnifiedMovieConverterService_ProcessConversion_NilBatchConfig() {
	// Arrange
	batchConfig := &BatchConversionConfig{
		InputDir:  "/test/input",
		InputExt:  ".mp4",
		OutputDir: "/test/output",
		OutputExt: ".gif",
	}
	service := NewUnifiedMovieConverterService(nil, batchConfig)
	// Force batch config to nil to test error handling
	service.config.BatchConfig = nil

	// Act
	result, err := service.ProcessConversion()

	// Assert
	if err == nil {
		ts.t.Error("Expected error for nil batch config, got nil")
	}
	if result == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	if result.Success {
		ts.t.Error("Expected Success to be false for nil config")
	}
	expectedMsg := "バッチ処理の設定が指定されていません"
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// #==============================================================#
// ##         Tests for unknown processing mode                  ##
// #==============================================================#
// TestUnifiedMovieConverterServiceUnknownMode tests unknown processing mode
type TestUnifiedMovieConverterServiceUnknownMode struct {
	t *testing.T
}

// NewTestUnifiedMovieConverterServiceUnknownMode creates a new test instance
func NewTestUnifiedMovieConverterServiceUnknownMode(t *testing.T) *TestUnifiedMovieConverterServiceUnknownMode {
	return &TestUnifiedMovieConverterServiceUnknownMode{t: t}
}

// TestUnifiedMovieConverterService_ProcessConversion_UnknownMode tests ProcessConversion with unknown mode
func (ts *TestUnifiedMovieConverterServiceUnknownMode) TestUnifiedMovieConverterService_ProcessConversion_UnknownMode() {
	// Arrange
	service := NewUnifiedMovieConverterService(nil, nil)
	// Force an unknown mode
	service.config.Mode = ProcessingMode(999)

	// Act
	result, err := service.ProcessConversion()

	// Assert
	if err == nil {
		ts.t.Error("Expected error for unknown processing mode, got nil")
	}
	if result == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	if result.Success {
		ts.t.Error("Expected Success to be false for unknown mode")
	}
	expectedMsg := "不明な処理モード: 999"
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// #==============================================================#
// ##         Tests for single file processing with valid file   ##
// #==============================================================#
// TestUnifiedMovieConverterServiceValidFile tests single file processing with valid file
type TestUnifiedMovieConverterServiceValidFile struct {
	t *testing.T
}

// NewTestUnifiedMovieConverterServiceValidFile creates a new test instance
func NewTestUnifiedMovieConverterServiceValidFile(t *testing.T) *TestUnifiedMovieConverterServiceValidFile {
	return &TestUnifiedMovieConverterServiceValidFile{t: t}
}

// TestUnifiedMovieConverterService_ProcessConversion_ValidFile_AutoGenerateOutput tests ProcessConversion with valid file and auto-generated output
func (ts *TestUnifiedMovieConverterServiceValidFile) TestUnifiedMovieConverterService_ProcessConversion_ValidFile_AutoGenerateOutput() {
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
	}
	service := NewUnifiedMovieConverterService(singleConfig, nil)

	// Act
	result, _ := service.ProcessConversion()

	// Assert
	// This will fail because ffmpeg is not available, but we can test the setup
	if result == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	if result.Mode != SingleFileMode {
		ts.t.Errorf("Expected SingleFileMode, got %d", result.Mode)
	}
	if result.SingleResult == nil {
		ts.t.Error("Expected SingleResult, got nil")
	}
	// Check that output file was auto-generated
	expectedOutput := filepath.Join(tempDir, "test.gif")
	if result.SingleResult.OutputFile != expectedOutput {
		ts.t.Errorf("Expected auto-generated output file %s, got %s", expectedOutput, result.SingleResult.OutputFile)
	}
}

// #==============================================================#
// ##         Tests for generateOutputFile function              ##
// #==============================================================#
// TestGenerateOutputFileFunction tests the generateOutputFile function
type TestGenerateOutputFileFunction struct {
	t *testing.T
}

// NewTestGenerateOutputFileFunction creates a new test instance
func NewTestGenerateOutputFileFunction(t *testing.T) *TestGenerateOutputFileFunction {
	return &TestGenerateOutputFileFunction{t: t}
}

// TestGenerateOutputFile_WithPath_Normal tests generateOutputFile with path
func (ts *TestGenerateOutputFileFunction) TestGenerateOutputFile_WithPath_Normal() {
	// Arrange
	inputFile := "/path/to/video.mp4"
	expected := "/path/to/video.gif"

	// Act
	result := generateOutputFile(inputFile)

	// Assert
	if result != expected {
		ts.t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestGenerateOutputFile_CaseInsensitive_Normal tests case insensitive extension handling
func (ts *TestGenerateOutputFileFunction) TestGenerateOutputFile_CaseInsensitive_Normal() {
	// Arrange
	inputFile := "VIDEO.MP4"
	expected := "VIDEO.gif"

	// Act
	result := generateOutputFile(inputFile)

	// Assert
	if result != expected {
		ts.t.Errorf("Expected %s, got %s", expected, result)
	}
}

// #==============================================================#
// ##         Tests for edge cases in unified service            ##
// #==============================================================#
// TestUnifiedMovieConverterServiceEdgeCases tests edge cases for unified service
type TestUnifiedMovieConverterServiceEdgeCases struct {
	t *testing.T
}

// NewTestUnifiedMovieConverterServiceEdgeCases creates a new test instance
func NewTestUnifiedMovieConverterServiceEdgeCases(t *testing.T) *TestUnifiedMovieConverterServiceEdgeCases {
	return &TestUnifiedMovieConverterServiceEdgeCases{t: t}
}

// TestUnifiedMovieConverterService_ProcessConversion_BatchModeWithPartialConfig tests batch mode with partial config
func (ts *TestUnifiedMovieConverterServiceEdgeCases) TestUnifiedMovieConverterService_ProcessConversion_BatchModeWithPartialConfig() {
	// Arrange
	tempDir := ts.t.TempDir()
	batchConfig := &BatchConversionConfig{
		InputDir: tempDir, // Only InputDir is set, should still trigger batch mode
	}
	service := NewUnifiedMovieConverterService(nil, batchConfig)

	// Act
	result, err := service.ProcessConversion()

	// Assert
	if err == nil {
		ts.t.Error("Expected error for incomplete batch config, got nil")
	}
	if result == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	if result.Mode != BatchMode {
		ts.t.Errorf("Expected BatchMode, got %d", result.Mode)
	}
	if result.Success {
		ts.t.Error("Expected Success to be false for incomplete config")
	}
}

// TestUnifiedMovieConverterService_ProcessConversion_SingleModeWithMKVFile tests single mode with MKV file
func (ts *TestUnifiedMovieConverterServiceEdgeCases) TestUnifiedMovieConverterService_ProcessConversion_SingleModeWithMKVFile() {
	// Arrange
	tempDir := ts.t.TempDir()
	testFile := filepath.Join(tempDir, "test.mkv")
	file, err := os.Create(testFile)
	if err != nil {
		ts.t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	singleConfig := &ConversionConfig{
		InputFile: testFile,
		// OutputFile is empty, should be auto-generated
	}
	service := NewUnifiedMovieConverterService(singleConfig, nil)

	// Act
	result, err := service.ProcessConversion()

	// Assert
	if result == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	if result.Mode != SingleFileMode {
		ts.t.Errorf("Expected SingleFileMode, got %d", result.Mode)
	}
	// The conversion will fail due to ffmpeg not being available, but we test the setup
	if err == nil {
		ts.t.Error("Expected error due to ffmpeg not being available")
	}
	if result.SingleResult == nil {
		ts.t.Error("Expected SingleResult, got nil")
	}
	// Check that output file was auto-generated for MKV -> GIF
	expectedOutput := filepath.Join(tempDir, "test.gif")
	if result.SingleResult.OutputFile != expectedOutput {
		ts.t.Errorf("Expected auto-generated output file %s, got %s", expectedOutput, result.SingleResult.OutputFile)
	}
}

// Standard Go test functions for unified edge tests

func TestGenerateOutputFileFunctionUnified(t *testing.T) {
	testService := NewTestGenerateOutputFileFunction(t)
	testService.TestGenerateOutputFile_WithPath_Normal()
	testService.TestGenerateOutputFile_CaseInsensitive_Normal()
}
func TestUnifiedMovieConverterServiceProcessConversionMethod(t *testing.T) {
	testService := NewTestUnifiedMovieConverterServiceProcessConversion(t)
	testService.TestUnifiedMovieConverterService_ProcessConversion_SingleFileMode_InvalidConfig()
	testService.TestUnifiedMovieConverterService_ProcessConversion_BatchMode_InvalidConfig()
	testService.TestUnifiedMovieConverterService_ProcessConversion_NilSingleConfig()
	testService.TestUnifiedMovieConverterService_ProcessConversion_NilBatchConfig()
}

func TestUnifiedMovieConverterServiceUnknownModeMethod(t *testing.T) {
	testService := NewTestUnifiedMovieConverterServiceUnknownMode(t)
	testService.TestUnifiedMovieConverterService_ProcessConversion_UnknownMode()
}

func TestUnifiedMovieConverterServiceValidFileMethod(t *testing.T) {
	testService := NewTestUnifiedMovieConverterServiceValidFile(t)
	testService.TestUnifiedMovieConverterService_ProcessConversion_ValidFile_AutoGenerateOutput()
}

func TestUnifiedMovieConverterServiceEdgeCasesMethod(t *testing.T) {
	testService := NewTestUnifiedMovieConverterServiceEdgeCases(t)
	testService.TestUnifiedMovieConverterService_ProcessConversion_BatchModeWithPartialConfig()
	testService.TestUnifiedMovieConverterService_ProcessConversion_SingleModeWithMKVFile()
}
