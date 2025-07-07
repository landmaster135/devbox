package usecases

import (
	"os"
	"path/filepath"
	"testing"
)

// #==============================================================#
// ##         Tests for Single File Conversion                   ##
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
		return
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
	result := generateOutputFile(inputFile)

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
	result := generateOutputFile(inputFile)

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
	result := generateOutputFile(inputFile)

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
	result := generateOutputFile(inputFile)

	// Assert
	if result != expected {
		ts.t.Errorf("Expected %s, got %s", expected, result)
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

// #==============================================================#
// ##         Tests for normalizeExtension function              ##
// #==============================================================#
// TestNormalizeExtension tests the normalizeExtension function
type TestNormalizeExtension struct {
	t *testing.T
}

// NewTestNormalizeExtension creates a new test instance
func NewTestNormalizeExtension(t *testing.T) *TestNormalizeExtension {
	return &TestNormalizeExtension{t: t}
}

// TestNormalizeExtension_EmptyString_Normal tests empty string handling
func (ts *TestNormalizeExtension) TestNormalizeExtension_EmptyString_Normal() {
	// Arrange
	ext := ""
	expected := ""

	// Act
	result := normalizeExtension(ext)

	// Assert
	if result != expected {
		ts.t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestNormalizeExtension_WithDot_Normal tests extension with dot
func (ts *TestNormalizeExtension) TestNormalizeExtension_WithDot_Normal() {
	// Arrange
	ext := ".mp4"
	expected := ".mp4"

	// Act
	result := normalizeExtension(ext)

	// Assert
	if result != expected {
		ts.t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestNormalizeExtension_WithoutDot_Normal tests extension without dot
func (ts *TestNormalizeExtension) TestNormalizeExtension_WithoutDot_Normal() {
	// Arrange
	ext := "mp4"
	expected := ".mp4"

	// Act
	result := normalizeExtension(ext)

	// Assert
	if result != expected {
		ts.t.Errorf("Expected %s, got %s", expected, result)
	}
}

// #==============================================================#
// ##         Tests for additional edge cases                    ##
// #==============================================================#
// TestMovieConverterServiceEdgeCases tests edge cases for MovieConverterService
type TestMovieConverterServiceEdgeCases struct {
	t *testing.T
}

// NewTestMovieConverterServiceEdgeCases creates a new test instance
func NewTestMovieConverterServiceEdgeCases(t *testing.T) *TestMovieConverterServiceEdgeCases {
	return &TestMovieConverterServiceEdgeCases{t: t}
}

// TestMovieConverterService_setMP4ToGIFDefaults_NonZeroValues tests setMP4ToGIFDefaults with non-zero values
func (ts *TestMovieConverterServiceEdgeCases) TestMovieConverterService_setMP4ToGIFDefaults_NonZeroValues() {
	// Arrange
	config := ConversionConfig{
		InputFile:   "test.mp4",
		OutputFile:  "test.gif",
		FPS:         30,  // Non-zero value should not be changed
		Speed:       1.5, // Non-zero value should not be changed
		UseItsScale: true,
	}
	service := NewMovieConverterService(config)

	// Act
	service.setMP4ToGIFDefaults()

	// Assert
	if service.config.FPS != 30 {
		ts.t.Errorf("Expected FPS to remain 30, got %d", service.config.FPS)
	}
	if service.config.Speed != 1.5 {
		ts.t.Errorf("Expected Speed to remain 1.5, got %f", service.config.Speed)
	}
}

// TestMovieConverterService_setGIFToMP4Defaults_NonZeroValues tests setGIFToMP4Defaults with non-zero values
func (ts *TestMovieConverterServiceEdgeCases) TestMovieConverterService_setGIFToMP4Defaults_NonZeroValues() {
	// Arrange
	config := ConversionConfig{
		InputFile:  "test.gif",
		OutputFile: "test.mp4",
		FPS:        30, // Non-zero value should not be changed
	}
	service := NewMovieConverterService(config)

	// Act
	service.setGIFToMP4Defaults()

	// Assert
	if service.config.FPS != 30 {
		ts.t.Errorf("Expected FPS to remain 30, got %d", service.config.FPS)
	}
}

// #==============================================================#
// ##         Tests for convertGIFToMP4 method                   ##
// #==============================================================#
// TestMovieConverterServiceConvertGIFToMP4 tests the convertGIFToMP4 method
type TestMovieConverterServiceConvertGIFToMP4 struct {
	t *testing.T
}

// NewTestMovieConverterServiceConvertGIFToMP4 creates a new test instance
func NewTestMovieConverterServiceConvertGIFToMP4(t *testing.T) *TestMovieConverterServiceConvertGIFToMP4 {
	return &TestMovieConverterServiceConvertGIFToMP4{t: t}
}

// TestMovieConverterService_convertGIFToMP4_WithSpaceInFilename tests convertGIFToMP4 with space in filename
func (ts *TestMovieConverterServiceConvertGIFToMP4) TestMovieConverterService_convertGIFToMP4_WithSpaceInFilename() {
	// Arrange
	tempDir := ts.t.TempDir()
	testFile := filepath.Join(tempDir, "test file.gif")
	file, err := os.Create(testFile)
	if err != nil {
		ts.t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:  testFile,
		OutputFile: filepath.Join(tempDir, "output.mp4"),
		FPS:        0, // Will be set to default
	}
	service := NewMovieConverterService(config)

	// Act
	err = service.convertGIFToMP4()

	// Assert
	// This will fail because ffmpeg is not available, but we test the setup and warning
	if err == nil {
		ts.t.Error("Expected error due to ffmpeg not being available")
	}
	// The method should have been called and defaults should be set
	if service.config.FPS != 15 {
		ts.t.Errorf("Expected FPS to be set to default 15, got %d", service.config.FPS)
	}
}

// TestMovieConverterService_convertGIFToMP4_DefaultFPS tests convertGIFToMP4 with default FPS
func (ts *TestMovieConverterServiceConvertGIFToMP4) TestMovieConverterService_convertGIFToMP4_DefaultFPS() {
	// Arrange
	tempDir := ts.t.TempDir()
	testFile := filepath.Join(tempDir, "test.gif")
	file, err := os.Create(testFile)
	if err != nil {
		ts.t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:  testFile,
		OutputFile: filepath.Join(tempDir, "output.mp4"),
		FPS:        0, // Will be set to default
	}
	service := NewMovieConverterService(config)

	// Act
	err = service.convertGIFToMP4()

	// Assert
	// This will fail because ffmpeg is not available, but we test the setup
	if err == nil {
		ts.t.Error("Expected error due to ffmpeg not being available")
	}
	// The method should have been called and defaults should be set
	if service.config.FPS != 15 {
		ts.t.Errorf("Expected FPS to be set to default 15, got %d", service.config.FPS)
	}
}

// TestMovieConverterService_convertGIFToMP4_CustomFPS tests convertGIFToMP4 with custom FPS
func (ts *TestMovieConverterServiceConvertGIFToMP4) TestMovieConverterService_convertGIFToMP4_CustomFPS() {
	// Arrange
	tempDir := ts.t.TempDir()
	testFile := filepath.Join(tempDir, "test.gif")
	file, err := os.Create(testFile)
	if err != nil {
		ts.t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:  testFile,
		OutputFile: filepath.Join(tempDir, "output.mp4"),
		FPS:        30, // Custom FPS
	}
	service := NewMovieConverterService(config)

	// Act
	err = service.convertGIFToMP4()

	// Assert
	// This will fail because ffmpeg is not available, but we test the setup
	if err == nil {
		ts.t.Error("Expected error due to ffmpeg not being available")
	}
	// The custom FPS should remain unchanged
	if service.config.FPS != 30 {
		ts.t.Errorf("Expected FPS to remain 30, got %d", service.config.FPS)
	}
}

// #==============================================================#
// ##         Tests for convertMP4ToGIF method edge cases        ##
// #==============================================================#
// TestMovieConverterServiceConvertMP4ToGIFEdgeCases tests edge cases for convertMP4ToGIF
type TestMovieConverterServiceConvertMP4ToGIFEdgeCases struct {
	t *testing.T
}

// NewTestMovieConverterServiceConvertMP4ToGIFEdgeCases creates a new test instance
func NewTestMovieConverterServiceConvertMP4ToGIFEdgeCases(t *testing.T) *TestMovieConverterServiceConvertMP4ToGIFEdgeCases {
	return &TestMovieConverterServiceConvertMP4ToGIFEdgeCases{t: t}
}

// TestMovieConverterService_convertMP4ToGIF_WithWidth tests convertMP4ToGIF with width specified
func (ts *TestMovieConverterServiceConvertMP4ToGIFEdgeCases) TestMovieConverterService_convertMP4ToGIF_WithWidth() {
	// Arrange
	tempDir := ts.t.TempDir()
	testFile := filepath.Join(tempDir, "test.mp4")
	file, err := os.Create(testFile)
	if err != nil {
		ts.t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:   testFile,
		OutputFile:  filepath.Join(tempDir, "output.gif"),
		FPS:         0,   // Will be set to default
		Width:       320, // Specified width
		Speed:       0,   // Will be set to default
		UseItsScale: true,
	}
	service := NewMovieConverterService(config)

	// Act
	err = service.convertMP4ToGIF()

	// Assert
	// This will fail because ffmpeg is not available, but we test the setup
	if err == nil {
		ts.t.Error("Expected error due to ffmpeg not being available")
	}
	// Defaults should be set
	if service.config.FPS != 60 {
		ts.t.Errorf("Expected FPS to be set to default 60, got %d", service.config.FPS)
	}
	if service.config.Speed != 2.0 {
		ts.t.Errorf("Expected Speed to be set to default 2.0, got %f", service.config.Speed)
	}
}

// TestMovieConverterService_convertMP4ToGIF_WithoutItsScale tests convertMP4ToGIF without itsscale
func (ts *TestMovieConverterServiceConvertMP4ToGIFEdgeCases) TestMovieConverterService_convertMP4ToGIF_WithoutItsScale() {
	// Arrange
	tempDir := ts.t.TempDir()
	testFile := filepath.Join(tempDir, "test.mp4")
	file, err := os.Create(testFile)
	if err != nil {
		ts.t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:   testFile,
		OutputFile:  filepath.Join(tempDir, "output.gif"),
		FPS:         30,
		Width:       0, // Default quality
		Speed:       1.5,
		UseItsScale: false, // Use setpts instead
	}
	service := NewMovieConverterService(config)

	// Act
	err = service.convertMP4ToGIF()

	// Assert
	// This will fail because ffmpeg is not available, but we test the setup
	if err == nil {
		ts.t.Error("Expected error due to ffmpeg not being available")
	}
	// Values should remain as set
	if service.config.FPS != 30 {
		ts.t.Errorf("Expected FPS to remain 30, got %d", service.config.FPS)
	}
	if service.config.Speed != 1.5 {
		ts.t.Errorf("Expected Speed to remain 1.5, got %f", service.config.Speed)
	}
}

// TestMovieConverterService_convertMP4ToGIF_WithSpaceInFilename tests convertMP4ToGIF with space in filename
func (ts *TestMovieConverterServiceConvertMP4ToGIFEdgeCases) TestMovieConverterService_convertMP4ToGIF_WithSpaceInFilename() {
	// Arrange
	tempDir := ts.t.TempDir()
	testFile := filepath.Join(tempDir, "test file.mp4")
	file, err := os.Create(testFile)
	if err != nil {
		ts.t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:   testFile,
		OutputFile:  filepath.Join(tempDir, "output.gif"),
		FPS:         30,
		Width:       0,
		Speed:       2.0,
		UseItsScale: true,
	}
	service := NewMovieConverterService(config)

	// Act
	err = service.convertMP4ToGIF()

	// Assert
	// This will fail because ffmpeg is not available, but we test the setup and warning
	if err == nil {
		ts.t.Error("Expected error due to ffmpeg not being available")
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

func TestSupportedExtensions(t *testing.T) {
	testService := NewTestGetSupportedExtensions(t)
	testService.TestGetSupportedExtensions_Normal()
}

func TestServiceDefaults(t *testing.T) {
	testService := NewTestMovieConverterServiceDefaults(t)
	testService.TestMovieConverterService_setMP4ToGIFDefaults_Normal()
	testService.TestMovieConverterService_setGIFToMP4Defaults_Normal()
}

func TestNormalizeExtensionFunction(t *testing.T) {
	testService := NewTestNormalizeExtension(t)
	testService.TestNormalizeExtension_EmptyString_Normal()
	testService.TestNormalizeExtension_WithDot_Normal()
	testService.TestNormalizeExtension_WithoutDot_Normal()
}

func TestMovieConverterServiceEdgeCasesMethod(t *testing.T) {
	testService := NewTestMovieConverterServiceEdgeCases(t)
	testService.TestMovieConverterService_setMP4ToGIFDefaults_NonZeroValues()
	testService.TestMovieConverterService_setGIFToMP4Defaults_NonZeroValues()
}

func TestMovieConverterServiceConvertGIFToMP4Method(t *testing.T) {
	testService := NewTestMovieConverterServiceConvertGIFToMP4(t)
	testService.TestMovieConverterService_convertGIFToMP4_WithSpaceInFilename()
	testService.TestMovieConverterService_convertGIFToMP4_DefaultFPS()
	testService.TestMovieConverterService_convertGIFToMP4_CustomFPS()
}

func TestMovieConverterServiceConvertMP4ToGIFEdgeCasesMethod(t *testing.T) {
	testService := NewTestMovieConverterServiceConvertMP4ToGIFEdgeCases(t)
	testService.TestMovieConverterService_convertMP4ToGIF_WithWidth()
	testService.TestMovieConverterService_convertMP4ToGIF_WithoutItsScale()
	testService.TestMovieConverterService_convertMP4ToGIF_WithSpaceInFilename()
}
