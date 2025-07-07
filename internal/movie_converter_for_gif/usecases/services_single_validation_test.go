package usecases

import (
	"os"
	"path/filepath"
	"testing"
)

// #==============================================================#
// ##         Tests for Single File Validation                   ##
// #==============================================================#
// TestValidateConfig tests the ValidateConfig function
type TestValidateConfig struct {
	t *testing.T
}

// NewTestValidateConfig creates a new test instance
func NewTestValidateConfig(t *testing.T) *TestValidateConfig {
	return &TestValidateConfig{t: t}
}

// TestValidateConfig_EmptyInputFile tests validation with empty input file
func (ts *TestValidateConfig) TestValidateConfig_EmptyInputFile() {
	// Arrange
	config := ConversionConfig{
		InputFile:  "",
		OutputFile: "output.gif",
	}

	// Act
	err := validateSingleConfig(config)

	// Assert
	if err == nil {
		ts.t.Error("Expected error for empty input file, got nil")
	}
	if err.Error() != "入力ファイルが指定されていません" {
		ts.t.Errorf("Expected specific error message, got %s", err.Error())
	}
}

// TestValidateConfig_NoExtension tests validation with no file extension
func (ts *TestValidateConfig) TestValidateConfig_NoExtension() {
	// Arrange
	config := ConversionConfig{
		InputFile:  "videofile",
		OutputFile: "output.gif",
	}

	// Act
	err := validateSingleConfig(config)

	// Assert
	if err == nil {
		ts.t.Error("Expected error for file without extension, got nil")
	}
	expectedMsg := "入力ファイル名に拡張子が含まれていません: videofile"
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestValidateConfig_FileNotExists tests validation with non-existent file
func (ts *TestValidateConfig) TestValidateConfig_FileNotExists() {
	// Arrange
	config := ConversionConfig{
		InputFile:  "nonexistent.mp4",
		OutputFile: "output.gif",
	}

	// Act
	err := validateSingleConfig(config)

	// Assert
	if err == nil {
		ts.t.Error("Expected error for non-existent file, got nil")
	}
	expectedMsg := "入力ファイルが見つかりません: nonexistent.mp4"
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestValidateConfig_ValidFile tests validation with valid file
func (ts *TestValidateConfig) TestValidateConfig_ValidFile() {
	// Arrange
	// Create a temporary test file
	tempDir := ts.t.TempDir()
	testFile := filepath.Join(tempDir, "test.mp4")
	file, err := os.Create(testFile)
	if err != nil {
		ts.t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:  testFile,
		OutputFile: "output.gif",
	}

	// Act
	err = validateSingleConfig(config)

	// Assert
	if err != nil {
		ts.t.Errorf("Expected no error for valid file, got %v", err)
	}
}

// #==============================================================#
// ##         Tests for convert method                           ##
// #==============================================================#
// TestMovieConverterServiceConvert tests the convert method
type TestMovieConverterServiceConvert struct {
	t *testing.T
}

// NewTestMovieConverterServiceConvert creates a new test instance
func NewTestMovieConverterServiceConvert(t *testing.T) *TestMovieConverterServiceConvert {
	return &TestMovieConverterServiceConvert{t: t}
}

// TestMovieConverterService_convert_FileNotExists tests convert with non-existent file
func (ts *TestMovieConverterServiceConvert) TestMovieConverterService_convert_FileNotExists() {
	// Arrange
	config := ConversionConfig{
		InputFile:  "nonexistent.mp4",
		OutputFile: "output.gif",
	}
	service := NewMovieConverterService(config)

	// Act
	err := service.convert()

	// Assert
	if err == nil {
		ts.t.Error("Expected error for non-existent file, got nil")
	}
	expectedMsg := "入力ファイルが見つかりません: nonexistent.mp4"
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestMovieConverterService_convert_NoExtension tests convert with no extension
func (ts *TestMovieConverterServiceConvert) TestMovieConverterService_convert_NoExtension() {
	// Arrange
	// Create a temporary test file without extension
	tempDir := ts.t.TempDir()
	testFile := filepath.Join(tempDir, "testfile")
	file, err := os.Create(testFile)
	if err != nil {
		ts.t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:  testFile,
		OutputFile: "output.gif",
	}
	service := NewMovieConverterService(config)

	// Act
	err = service.convert()

	// Assert
	if err == nil {
		ts.t.Error("Expected error for file without extension, got nil")
	}
	expectedMsg := "入力ファイル名に拡張子が含まれていません: " + testFile
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestMovieConverterService_convert_UnsupportedConversion tests unsupported conversion
func (ts *TestMovieConverterServiceConvert) TestMovieConverterService_convert_UnsupportedConversion() {
	// Arrange
	// Create a temporary test file with unsupported extension
	tempDir := ts.t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	file, err := os.Create(testFile)
	if err != nil {
		ts.t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:  testFile,
		OutputFile: "output.gif",
	}
	service := NewMovieConverterService(config)

	// Act
	err = service.convert()

	// Assert
	if err == nil {
		ts.t.Error("Expected error for unsupported conversion, got nil")
	}
	expectedMsg := "サポートされていない変換: .txt -> .gif"
	if err.Error() != expectedMsg {
		ts.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// Standard Go test functions for validation tests

func TestConfigValidation(t *testing.T) {
	testService := NewTestValidateConfig(t)
	testService.TestValidateConfig_EmptyInputFile()
	testService.TestValidateConfig_NoExtension()
	testService.TestValidateConfig_FileNotExists()
	testService.TestValidateConfig_ValidFile()
}

func TestMovieConverterServiceConvertMethod(t *testing.T) {
	testService := NewTestMovieConverterServiceConvert(t)
	testService.TestMovieConverterService_convert_FileNotExists()
	testService.TestMovieConverterService_convert_NoExtension()
	testService.TestMovieConverterService_convert_UnsupportedConversion()
}
