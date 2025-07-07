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

// #==============================================================#
// ##         Tests for Batch                                    ##
// #==============================================================#
// TestBatchMovieConverterService tests the BatchMovieConverterService struct
type TestBatchMovieConverterService struct {
	t *testing.T
}

// NewTestBatchMovieConverterService creates a new test instance
func NewTestBatchMovieConverterService(t *testing.T) *TestBatchMovieConverterService {
	return &TestBatchMovieConverterService{t: t}
}

// TestNewBatchMovieConverterService_Normal tests normal creation of BatchMovieConverterService
func (ts *TestBatchMovieConverterService) TestNewBatchMovieConverterService_Normal() {
	// Arrange
	config := BatchConversionConfig{
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
	service := NewBatchMovieConverterService(config)

	// Assert
	if service == nil {
		ts.t.Error("NewBatchMovieConverterService should not return nil")
	}
	if service.config.InputDir != config.InputDir {
		ts.t.Errorf("Expected InputDir %s, got %s", config.InputDir, service.config.InputDir)
	}
	if service.config.InputExt != config.InputExt {
		ts.t.Errorf("Expected InputExt %s, got %s", config.InputExt, service.config.InputExt)
	}
	if service.config.Recursive != config.Recursive {
		ts.t.Errorf("Expected Recursive %t, got %t", config.Recursive, service.config.Recursive)
	}
}

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

// Standard Go test functions for batch processing

func TestBatchMovieConverterServiceCreation(t *testing.T) {
	testService := NewTestBatchMovieConverterService(t)
	testService.TestNewBatchMovieConverterService_Normal()
}

func TestBatchConfigValidation(t *testing.T) {
	testService := NewTestValidateBatchConfig(t)
	testService.TestValidateBatchConfig_EmptyInputDir()
	testService.TestValidateBatchConfig_EmptyInputExt()
	testService.TestValidateBatchConfig_UnsupportedInputExt()
	testService.TestValidateBatchConfig_UnsupportedOutputExt()
	testService.TestValidateBatchConfig_ValidConfig()
	testService.TestValidateBatchConfig_ExtensionNormalization()
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

// #==============================================================#
// ##         Tests for batch scanFiles method                   ##
// #==============================================================#
// TestBatchMovieConverterServiceScanFiles tests the scanFiles method
type TestBatchMovieConverterServiceScanFiles struct {
	t *testing.T
}

// NewTestBatchMovieConverterServiceScanFiles creates a new test instance
func NewTestBatchMovieConverterServiceScanFiles(t *testing.T) *TestBatchMovieConverterServiceScanFiles {
	return &TestBatchMovieConverterServiceScanFiles{t: t}
}

// TestBatchMovieConverterService_scanFiles_Normal tests normal file scanning
func (ts *TestBatchMovieConverterServiceScanFiles) TestBatchMovieConverterService_scanFiles_Normal() {
	// Arrange
	tempDir := ts.t.TempDir()

	// Create test files
	testFiles := []string{"video1.mp4", "video2.mp4", "image.jpg"}
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
		OutputDir: "/test/output",
		OutputExt: ".gif",
		Recursive: false,
	}
	service := NewBatchMovieConverterService(config)

	// Act
	files, err := service.scanFiles()

	// Assert
	if err != nil {
		ts.t.Errorf("Expected no error, got %v", err)
	}
	if len(files) != 2 {
		ts.t.Errorf("Expected 2 MP4 files, got %d", len(files))
	}
}

// TestBatchMovieConverterService_scanFiles_Recursive tests recursive file scanning
func (ts *TestBatchMovieConverterServiceScanFiles) TestBatchMovieConverterService_scanFiles_Recursive() {
	// Arrange
	tempDir := ts.t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	err := os.MkdirAll(subDir, 0755)
	if err != nil {
		ts.t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create test files in root and subdirectory
	rootFile := filepath.Join(tempDir, "root.mp4")
	subFile := filepath.Join(subDir, "sub.mp4")

	for _, fileName := range []string{rootFile, subFile} {
		file, err := os.Create(fileName)
		if err != nil {
			ts.t.Fatalf("Failed to create test file %s: %v", fileName, err)
		}
		file.Close()
	}

	config := BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputDir: "/test/output",
		OutputExt: ".gif",
		Recursive: true,
	}
	service := NewBatchMovieConverterService(config)

	// Act
	files, err := service.scanFiles()

	// Assert
	if err != nil {
		ts.t.Errorf("Expected no error, got %v", err)
	}
	if len(files) != 2 {
		ts.t.Errorf("Expected 2 MP4 files (recursive), got %d", len(files))
	}
}

// TestBatchMovieConverterService_scanFiles_NonRecursive tests non-recursive file scanning
func (ts *TestBatchMovieConverterServiceScanFiles) TestBatchMovieConverterService_scanFiles_NonRecursive() {
	// Arrange
	tempDir := ts.t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	err := os.MkdirAll(subDir, 0755)
	if err != nil {
		ts.t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create test files in root and subdirectory
	rootFile := filepath.Join(tempDir, "root.mp4")
	subFile := filepath.Join(subDir, "sub.mp4")

	for _, fileName := range []string{rootFile, subFile} {
		file, err := os.Create(fileName)
		if err != nil {
			ts.t.Fatalf("Failed to create test file %s: %v", fileName, err)
		}
		file.Close()
	}

	config := BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputDir: "/test/output",
		OutputExt: ".gif",
		Recursive: false,
	}
	service := NewBatchMovieConverterService(config)

	// Act
	files, err := service.scanFiles()

	// Assert
	if err != nil {
		ts.t.Errorf("Expected no error, got %v", err)
	}
	if len(files) != 1 {
		ts.t.Errorf("Expected 1 MP4 file (non-recursive), got %d", len(files))
	}
}

// Standard Go test functions for new tests

func TestNormalizeExtensionFunction(t *testing.T) {
	testService := NewTestNormalizeExtension(t)
	testService.TestNormalizeExtension_EmptyString_Normal()
	testService.TestNormalizeExtension_WithDot_Normal()
	testService.TestNormalizeExtension_WithoutDot_Normal()
}

func TestMovieConverterServiceConvertMethod(t *testing.T) {
	testService := NewTestMovieConverterServiceConvert(t)
	testService.TestMovieConverterService_convert_FileNotExists()
	testService.TestMovieConverterService_convert_NoExtension()
	testService.TestMovieConverterService_convert_UnsupportedConversion()
}

func TestBatchMovieConverterServiceScanFilesMethod(t *testing.T) {
	testService := NewTestBatchMovieConverterServiceScanFiles(t)
	testService.TestBatchMovieConverterService_scanFiles_Normal()
	testService.TestBatchMovieConverterService_scanFiles_Recursive()
	testService.TestBatchMovieConverterService_scanFiles_NonRecursive()
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
		FPS:         30, // Non-zero value should not be changed
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

// Standard Go test functions for new tests

func TestBatchMovieConverterServiceBatchConvertMethod(t *testing.T) {
	testService := NewTestBatchMovieConverterServiceBatchConvert(t)
	testService.TestBatchMovieConverterService_batchConvert_InputDirNotExists()
	testService.TestBatchMovieConverterService_batchConvert_NoFiles()
}

func TestMovieConverterServiceEdgeCasesMethod(t *testing.T) {
	testService := NewTestMovieConverterServiceEdgeCases(t)
	testService.TestMovieConverterService_setMP4ToGIFDefaults_NonZeroValues()
	testService.TestMovieConverterService_setGIFToMP4Defaults_NonZeroValues()
}

func TestValidateBatchConfigEdgeCasesMethod(t *testing.T) {
	testService := NewTestValidateBatchConfigEdgeCases(t)
	testService.TestValidateBatchConfig_EmptyOutputDir()
	testService.TestValidateBatchConfig_EmptyOutputExt()
}
