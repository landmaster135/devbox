package usecases

import (
	"os"
	"path/filepath"
	"testing"
)

// #==============================================================#
// ##         Tests                                              ##
// #==============================================================#
// TestMovieConverterService tests the MovieConverterService struct
type TestMovieConverterService struct {
	t *testing.T
}

// NewTestMovieConverterService creates a new test instance
func NewTestMovieConverterService(t *testing.T) *TestMovieConverterService {
	return &TestMovieConverterService{t: t}
}

// TestNewMovieConverterService_Normal tests normal creation of MovieConverterService
func (ts *TestMovieConverterService) TestNewMovieConverterService_Normal() {
	// Arrange
	config := ConversionConfig{
		InputFile:   "test.mp4",
		OutputFile:  "test.gif",
		FPS:         30,
		Width:       320,
		Speed:       1.5,
		Loop:        0,
		UseItsScale: true,
	}

	// Act
	service := NewMovieConverterService(config)

	// Assert
	if service == nil {
		ts.t.Error("NewMovieConverterService should not return nil")
	}
	if service.config.InputFile != config.InputFile {
		ts.t.Errorf("Expected InputFile %s, got %s", config.InputFile, service.config.InputFile)
	}
	if service.config.OutputFile != config.OutputFile {
		ts.t.Errorf("Expected OutputFile %s, got %s", config.OutputFile, service.config.OutputFile)
	}
	if service.config.FPS != config.FPS {
		ts.t.Errorf("Expected FPS %d, got %d", config.FPS, service.config.FPS)
	}
}

// TestGenerateOutputFile tests the GenerateOutputFile function
type TestGenerateOutputFile struct {
	t *testing.T
}

// NewTestGenerateOutputFile creates a new test instance
func NewTestGenerateOutputFile(t *testing.T) *TestGenerateOutputFile {
	return &TestGenerateOutputFile{t: t}
}

// TestGenerateOutputFile_MP4ToGIF_Normal tests MP4 to GIF filename generation
func (ts *TestGenerateOutputFile) TestGenerateOutputFile_MP4ToGIF_Normal() {
	// Arrange
	inputFile := "video.mp4"
	expected := "video.gif"

	// Act
	result := GenerateOutputFile(inputFile)

	// Assert
	if result != expected {
		ts.t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestGenerateOutputFile_GIFToMP4_Normal tests GIF to MP4 filename generation
func (ts *TestGenerateOutputFile) TestGenerateOutputFile_GIFToMP4_Normal() {
	// Arrange
	inputFile := "animation.gif"
	expected := "animation.mp4"

	// Act
	result := GenerateOutputFile(inputFile)

	// Assert
	if result != expected {
		ts.t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestGenerateOutputFile_MKVToGIF_Normal tests MKV to GIF filename generation
func (ts *TestGenerateOutputFile) TestGenerateOutputFile_MKVToGIF_Normal() {
	// Arrange
	inputFile := "video.mkv"
	expected := "video.gif"

	// Act
	result := GenerateOutputFile(inputFile)

	// Assert
	if result != expected {
		ts.t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestGenerateOutputFile_UnsupportedExtension_Normal tests unsupported extension handling
func (ts *TestGenerateOutputFile) TestGenerateOutputFile_UnsupportedExtension_Normal() {
	// Arrange
	inputFile := "document.txt"
	expected := "document_converted"

	// Act
	result := GenerateOutputFile(inputFile)

	// Assert
	if result != expected {
		ts.t.Errorf("Expected %s, got %s", expected, result)
	}
}

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
	err := ValidateConfig(config)

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
	err := ValidateConfig(config)

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
	err := ValidateConfig(config)

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
	err = ValidateConfig(config)

	// Assert
	if err != nil {
		ts.t.Errorf("Expected no error for valid file, got %v", err)
	}
}

// TestGetSupportedExtensions tests the GetSupportedExtensions function
type TestGetSupportedExtensions struct {
	t *testing.T
}

// NewTestGetSupportedExtensions creates a new test instance
func NewTestGetSupportedExtensions(t *testing.T) *TestGetSupportedExtensions {
	return &TestGetSupportedExtensions{t: t}
}

// TestGetSupportedExtensions_Normal tests getting supported extensions
func (ts *TestGetSupportedExtensions) TestGetSupportedExtensions_Normal() {
	// Act
	extensions := GetSupportedExtensions()

	// Assert
	if extensions == nil {
		ts.t.Error("GetSupportedExtensions should not return nil")
	}

	inputExts, exists := extensions["input"]
	if !exists {
		ts.t.Error("Expected 'input' key in extensions map")
	}
	expectedInputExts := []string{".mp4", ".mkv", ".gif"}
	if len(inputExts) != len(expectedInputExts) {
		ts.t.Errorf("Expected %d input extensions, got %d", len(expectedInputExts), len(inputExts))
	}

	outputExts, exists := extensions["output"]
	if !exists {
		ts.t.Error("Expected 'output' key in extensions map")
	}
	expectedOutputExts := []string{".mp4", ".gif"}
	if len(outputExts) != len(expectedOutputExts) {
		ts.t.Errorf("Expected %d output extensions, got %d", len(expectedOutputExts), len(outputExts))
	}
}

// TestMovieConverterService_setMP4ToGIFDefaults tests default value setting for MP4 to GIF
type TestMovieConverterServiceDefaults struct {
	t *testing.T
}

// NewTestMovieConverterServiceDefaults creates a new test instance
func NewTestMovieConverterServiceDefaults(t *testing.T) *TestMovieConverterServiceDefaults {
	return &TestMovieConverterServiceDefaults{t: t}
}

// TestMovieConverterService_setMP4ToGIFDefaults_Normal tests setting MP4 to GIF defaults
func (ts *TestMovieConverterServiceDefaults) TestMovieConverterService_setMP4ToGIFDefaults_Normal() {
	// Arrange
	config := ConversionConfig{
		InputFile:   "test.mp4",
		OutputFile:  "test.gif",
		FPS:         0, // Should be set to default
		Speed:       0, // Should be set to default
		UseItsScale: true,
	}
	service := NewMovieConverterService(config)

	// Act
	service.setMP4ToGIFDefaults()

	// Assert
	if service.config.FPS != 60 {
		ts.t.Errorf("Expected FPS to be set to 60, got %d", service.config.FPS)
	}
	if service.config.Speed != 2.0 {
		ts.t.Errorf("Expected Speed to be set to 2.0, got %f", service.config.Speed)
	}
}

// TestMovieConverterService_setGIFToMP4Defaults_Normal tests setting GIF to MP4 defaults
func (ts *TestMovieConverterServiceDefaults) TestMovieConverterService_setGIFToMP4Defaults_Normal() {
	// Arrange
	config := ConversionConfig{
		InputFile:  "test.gif",
		OutputFile: "test.mp4",
		FPS:        0, // Should be set to default
	}
	service := NewMovieConverterService(config)

	// Act
	service.setGIFToMP4Defaults()

	// Assert
	if service.config.FPS != 15 {
		ts.t.Errorf("Expected FPS to be set to 15, got %d", service.config.FPS)
	}
}

// Standard Go test functions that call the test methods

func TestMovieConverterServiceCreation(t *testing.T) {
	testService := NewTestMovieConverterService(t)
	testService.TestNewMovieConverterService_Normal()
}

func TestOutputFileGeneration(t *testing.T) {
	testService := NewTestGenerateOutputFile(t)
	testService.TestGenerateOutputFile_MP4ToGIF_Normal()
	testService.TestGenerateOutputFile_GIFToMP4_Normal()
	testService.TestGenerateOutputFile_MKVToGIF_Normal()
	testService.TestGenerateOutputFile_UnsupportedExtension_Normal()
}

func TestConfigValidation(t *testing.T) {
	testService := NewTestValidateConfig(t)
	testService.TestValidateConfig_EmptyInputFile()
	testService.TestValidateConfig_NoExtension()
	testService.TestValidateConfig_FileNotExists()
	testService.TestValidateConfig_ValidFile()
}

func TestSupportedExtensions(t *testing.T) {
	testService := NewTestGetSupportedExtensions(t)
	testService.TestGetSupportedExtensions_Normal()
}

func TestServiceDefaults(t *testing.T) {
	testService := NewTestMovieConverterServiceDefaults(t)
	testService.TestMovieConverterService_setMP4ToGIFDefaults_Normal()
	testService.TestMovieConverterService_setGIFToMP4Defaults_Normal()
}
