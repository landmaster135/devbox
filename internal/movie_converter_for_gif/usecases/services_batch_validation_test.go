package usecases

import (
	"os"
	"path/filepath"
	"testing"
)

// #==============================================================#
// ##         Tests for Batch Validation                         ##
// #==============================================================#
// TestValidateBatchConfig tests the ValidateBatchConfig function
type TestValidateBatchConfig struct {
	t *testing.T
}

// NewTestValidateBatchConfig creates a new test instance
func NewTestValidateBatchConfig(t *testing.T) *TestValidateBatchConfig {
	return &TestValidateBatchConfig{t: t}
}

// TestValidateBatchConfig_EmptyInputDir tests validation with empty input directory
func (ts *TestValidateBatchConfig) TestValidateBatchConfig_EmptyInputDir() {
	// Arrange
	config := BatchConversionConfig{
		InputDir:  "",
		InputExt:  ".mp4",
		OutputDir: "/test/output",
		OutputExt: ".gif",
	}

	// Act
	err := validateBatchConfig(&config)

	// Assert
	if err == nil {
		ts.t.Error("Expected error for empty input directory, got nil")
	}
	if err.Error() != "入力ディレクトリが指定されていません" {
		ts.t.Errorf("Expected specific error message, got %s", err.Error())
	}
}

// TestValidateBatchConfig_EmptyInputExt tests validation with empty input extension
func (ts *TestValidateBatchConfig) TestValidateBatchConfig_EmptyInputExt() {
	// Arrange
	config := BatchConversionConfig{
		InputDir:  "/test/input",
		InputExt:  "",
		OutputDir: "/test/output",
		OutputExt: ".gif",
	}

	// Act
	err := validateBatchConfig(&config)

	// Assert
	if err == nil {
		ts.t.Error("Expected error for empty input extension, got nil")
	}
	if err.Error() != "入力拡張子が指定されていません" {
		ts.t.Errorf("Expected specific error message, got %s", err.Error())
	}
}

// TestValidateBatchConfig_UnsupportedInputExt tests validation with unsupported input extension
func (ts *TestValidateBatchConfig) TestValidateBatchConfig_UnsupportedInputExt() {
	// Arrange
	tempDir := ts.t.TempDir()
	config := BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".txt",
		OutputDir: "/test/output",
		OutputExt: ".gif",
	}

	// Act
	err := validateBatchConfig(&config)

	// Assert
	if err == nil {
		ts.t.Error("Expected error for unsupported input extension, got nil")
	}
	expectedMsg := "サポートされていない入力拡張子: .txt"
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestValidateBatchConfig_UnsupportedOutputExt tests validation with unsupported output extension
func (ts *TestValidateBatchConfig) TestValidateBatchConfig_UnsupportedOutputExt() {
	// Arrange
	tempDir := ts.t.TempDir()
	config := BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputDir: "/test/output",
		OutputExt: ".txt",
	}

	// Act
	err := validateBatchConfig(&config)

	// Assert
	if err == nil {
		ts.t.Error("Expected error for unsupported output extension, got nil")
	}
	expectedMsg := "サポートされていない出力拡張子: .txt"
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestValidateBatchConfig_ValidConfig tests validation with valid configuration
func (ts *TestValidateBatchConfig) TestValidateBatchConfig_ValidConfig() {
	// Arrange
	tempDir := ts.t.TempDir()
	config := BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputDir: "/test/output",
		OutputExt: ".gif",
	}

	// Act
	err := validateBatchConfig(&config)

	// Assert
	if err != nil {
		ts.t.Errorf("Expected no error for valid config, got %v", err)
	}
}

// TestValidateBatchConfig_ExtensionNormalization tests extension normalization (dot addition)
func (ts *TestValidateBatchConfig) TestValidateBatchConfig_ExtensionNormalization() {
	// Arrange
	tempDir := ts.t.TempDir()
	config := BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  "mp4", // ドットなし
		OutputDir: "/test/output",
		OutputExt: "gif", // ドットなし
	}

	// Act
	err := validateBatchConfig(&config)

	// Assert
	if err != nil {
		ts.t.Errorf("Expected no error for valid config with dot-less extensions, got %v", err)
	}
	// 拡張子が正規化されていることを確認
	if config.InputExt != ".mp4" {
		ts.t.Errorf("Expected InputExt to be normalized to '.mp4', got %s", config.InputExt)
	}
	if config.OutputExt != ".gif" {
		ts.t.Errorf("Expected OutputExt to be normalized to '.gif', got %s", config.OutputExt)
	}
}

// #==============================================================#
// ##         Tests for validateBatchConfig edge cases           ##
// #==============================================================#
// TestValidateBatchConfigEdgeCases tests edge cases for validateBatchConfig
type TestValidateBatchConfigEdgeCases struct {
	t *testing.T
}

// NewTestValidateBatchConfigEdgeCases creates a new test instance
func NewTestValidateBatchConfigEdgeCases(t *testing.T) *TestValidateBatchConfigEdgeCases {
	return &TestValidateBatchConfigEdgeCases{t: t}
}

// TestValidateBatchConfig_EmptyOutputDir tests validation with empty output directory
func (ts *TestValidateBatchConfigEdgeCases) TestValidateBatchConfig_EmptyOutputDir() {
	// Arrange
	tempDir := ts.t.TempDir()
	config := BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputDir: "",
		OutputExt: ".gif",
	}

	// Act
	err := validateBatchConfig(&config)

	// Assert
	if err == nil {
		ts.t.Error("Expected error for empty output directory, got nil")
	}
	expectedMsg := "出力ディレクトリが指定されていません"
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestValidateBatchConfig_EmptyOutputExt tests validation with empty output extension
func (ts *TestValidateBatchConfigEdgeCases) TestValidateBatchConfig_EmptyOutputExt() {
	// Arrange
	tempDir := ts.t.TempDir()
	config := BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputDir: "/test/output",
		OutputExt: "",
	}

	// Act
	err := validateBatchConfig(&config)

	// Assert
	if err == nil {
		ts.t.Error("Expected error for empty output extension, got nil")
	}
	expectedMsg := "出力拡張子が指定されていません"
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// #==============================================================#
// ##         Tests for batch batchConvert method                ##
// #==============================================================#
// TestBatchMovieConverterServiceBatchConvert tests the batchConvert method
type TestBatchMovieConverterServiceBatchConvert struct {
	t *testing.T
}

// NewTestBatchMovieConverterServiceBatchConvert creates a new test instance
func NewTestBatchMovieConverterServiceBatchConvert(t *testing.T) *TestBatchMovieConverterServiceBatchConvert {
	return &TestBatchMovieConverterServiceBatchConvert{t: t}
}

// TestBatchMovieConverterService_batchConvert_InputDirNotExists tests batchConvert with non-existent input directory
func (ts *TestBatchMovieConverterServiceBatchConvert) TestBatchMovieConverterService_batchConvert_InputDirNotExists() {
	// Arrange
	config := BatchConversionConfig{
		InputDir:  "/nonexistent/directory",
		InputExt:  ".mp4",
		OutputDir: "/test/output",
		OutputExt: ".gif",
	}
	service := NewBatchMovieConverterService(config)

	// Act
	result, err := service.batchConvert()

	// Assert
	if err == nil {
		ts.t.Error("Expected error for non-existent input directory, got nil")
	}
	if result != nil {
		ts.t.Error("Expected nil result for error case, got non-nil")
	}
	expectedMsg := "入力ディレクトリが見つかりません: /nonexistent/directory"
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestBatchMovieConverterService_batchConvert_NoFiles tests batchConvert with no matching files
func (ts *TestBatchMovieConverterServiceBatchConvert) TestBatchMovieConverterService_batchConvert_NoFiles() {
	// Arrange
	tempDir := ts.t.TempDir()
	outputDir := ts.t.TempDir()

	// Create non-matching files
	testFiles := []string{"document.txt", "image.jpg"}
	for _, fileName := range testFiles {
		file, err := os.Create(filepath.Join(tempDir, fileName))
		if err != nil {
			ts.t.Fatalf("Failed to create test file %s: %v", fileName, err)
		}
		file.Close()
	}

	config := BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputDir: outputDir,
		OutputExt: ".gif",
	}
	service := NewBatchMovieConverterService(config)

	// Act
	result, err := service.batchConvert()

	// Assert
	if err != nil {
		ts.t.Errorf("Expected no error for no files case, got %v", err)
	}
	if result == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	if result.TotalFiles != 0 {
		ts.t.Errorf("Expected TotalFiles to be 0, got %d", result.TotalFiles)
	}
	if result.SuccessCount != 0 {
		ts.t.Errorf("Expected SuccessCount to be 0, got %d", result.SuccessCount)
	}
	if result.FailureCount != 0 {
		ts.t.Errorf("Expected FailureCount to be 0, got %d", result.FailureCount)
	}
}

// TestBatchMovieConverterService_batchConvert_OutputDirCreationFailure tests output directory creation failure
func (ts *TestBatchMovieConverterServiceBatchConvert) TestBatchMovieConverterService_batchConvert_OutputDirCreationFailure() {
	// Arrange
	tempDir := ts.t.TempDir()

	// Create a file where we want to create the output directory (this will cause mkdir to fail)
	invalidOutputDir := filepath.Join(tempDir, "invalid_output")
	file, err := os.Create(invalidOutputDir)
	if err != nil {
		ts.t.Fatalf("Failed to create blocking file: %v", err)
	}
	file.Close()

	// Create test file
	testFile := filepath.Join(tempDir, "test.mp4")
	file, err = os.Create(testFile)
	if err != nil {
		ts.t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputDir: invalidOutputDir, // This will fail because it's a file, not a directory
		OutputExt: ".gif",
	}
	service := NewBatchMovieConverterService(config)

	// Act
	result, err := service.batchConvert()

	// Assert
	if err == nil {
		ts.t.Error("Expected error for output directory creation failure, got nil")
	}
	if result != nil {
		ts.t.Error("Expected nil result for error case, got non-nil")
	}
}

// Standard Go test functions for batch validation tests

func TestBatchConfigValidation(t *testing.T) {
	testService := NewTestValidateBatchConfig(t)
	testService.TestValidateBatchConfig_EmptyInputDir()
	testService.TestValidateBatchConfig_EmptyInputExt()
	testService.TestValidateBatchConfig_UnsupportedInputExt()
	testService.TestValidateBatchConfig_UnsupportedOutputExt()
	testService.TestValidateBatchConfig_ValidConfig()
	testService.TestValidateBatchConfig_ExtensionNormalization()
}

func TestValidateBatchConfigEdgeCasesMethod(t *testing.T) {
	testService := NewTestValidateBatchConfigEdgeCases(t)
	testService.TestValidateBatchConfig_EmptyOutputDir()
	testService.TestValidateBatchConfig_EmptyOutputExt()
}

func TestBatchMovieConverterServiceBatchConvertMethod(t *testing.T) {
	testService := NewTestBatchMovieConverterServiceBatchConvert(t)
	testService.TestBatchMovieConverterService_batchConvert_InputDirNotExists()
	testService.TestBatchMovieConverterService_batchConvert_NoFiles()
	testService.TestBatchMovieConverterService_batchConvert_OutputDirCreationFailure()
}
