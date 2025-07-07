package usecases

import (
	"os"
	"path/filepath"
	"testing"
)

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
		return
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

// #==============================================================#
// ##         Tests for batchConvert method with files           ##
// #==============================================================#
// TestBatchMovieConverterServiceBatchConvertWithFiles tests batchConvert with actual files
type TestBatchMovieConverterServiceBatchConvertWithFiles struct {
	t *testing.T
}

// NewTestBatchMovieConverterServiceBatchConvertWithFiles creates a new test instance
func NewTestBatchMovieConverterServiceBatchConvertWithFiles(t *testing.T) *TestBatchMovieConverterServiceBatchConvertWithFiles {
	return &TestBatchMovieConverterServiceBatchConvertWithFiles{t: t}
}

// TestBatchMovieConverterService_batchConvert_WithFiles tests batchConvert with actual files
func (ts *TestBatchMovieConverterServiceBatchConvertWithFiles) TestBatchMovieConverterService_batchConvert_WithFiles() {
	// Arrange
	tempDir := ts.t.TempDir()
	outputDir := filepath.Join(tempDir, "output")

	// Create test files
	testFiles := []string{"video1.mp4", "video2.mp4", "video3.mp4"}
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
		FPS:       30,
		Width:     320,
		Speed:     1.5,
	}
	service := NewBatchMovieConverterService(config)

	// Act
	result, err := service.batchConvert()

	// Assert
	if err != nil {
		ts.t.Errorf("Expected no error for batchConvert setup, got %v", err)
	}
	if result == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	if result.TotalFiles != 3 {
		ts.t.Errorf("Expected TotalFiles to be 3, got %d", result.TotalFiles)
	}
	// All conversions will fail due to ffmpeg not being available, but we test the setup
	if result.FailureCount != 3 {
		ts.t.Errorf("Expected FailureCount to be 3 (due to ffmpeg), got %d", result.FailureCount)
	}
	if result.SuccessCount != 0 {
		ts.t.Errorf("Expected SuccessCount to be 0 (due to ffmpeg), got %d", result.SuccessCount)
	}
	if len(result.Results) != 3 {
		ts.t.Errorf("Expected 3 results, got %d", len(result.Results))
	}
	if len(result.FailedFiles) != 3 {
		ts.t.Errorf("Expected 3 failed files, got %d", len(result.FailedFiles))
	}
}

// TestBatchMovieConverterService_batchConvert_OutputDirCreationFailure tests output directory creation failure
func (ts *TestBatchMovieConverterServiceBatchConvertWithFiles) TestBatchMovieConverterService_batchConvert_OutputDirCreationFailure() {
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

// TestBatchMovieConverterService_batchConvert_MixedFileTypes tests batchConvert with mixed file types
func (ts *TestBatchMovieConverterServiceBatchConvertWithFiles) TestBatchMovieConverterService_batchConvert_MixedFileTypes() {
	// Arrange
	tempDir := ts.t.TempDir()
	outputDir := filepath.Join(tempDir, "output")

	// Create mixed files (only .mp4 should be processed)
	testFiles := []string{"video1.mp4", "video2.mp4", "image.jpg", "document.txt", "animation.gif"}
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
		ts.t.Errorf("Expected no error for batchConvert setup, got %v", err)
	}
	if result == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	// Only 2 MP4 files should be processed
	if result.TotalFiles != 2 {
		ts.t.Errorf("Expected TotalFiles to be 2, got %d", result.TotalFiles)
	}
}

// TestBatchMovieConverterService_batchConvert_RecursiveProcessing tests recursive processing
func (ts *TestBatchMovieConverterServiceBatchConvertWithFiles) TestBatchMovieConverterService_batchConvert_RecursiveProcessing() {
	// Arrange
	tempDir := ts.t.TempDir()
	subDir1 := filepath.Join(tempDir, "subdir1")
	subDir2 := filepath.Join(tempDir, "subdir2")
	outputDir := filepath.Join(tempDir, "output")

	// Create subdirectories
	err := os.MkdirAll(subDir1, 0755)
	if err != nil {
		ts.t.Fatalf("Failed to create subdirectory: %v", err)
	}
	err = os.MkdirAll(subDir2, 0755)
	if err != nil {
		ts.t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create test files in different directories
	testFiles := []string{
		filepath.Join(tempDir, "root.mp4"),
		filepath.Join(subDir1, "sub1.mp4"),
		filepath.Join(subDir2, "sub2.mp4"),
	}
	for _, fileName := range testFiles {
		file, err := os.Create(fileName)
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
		Recursive: true, // Enable recursive processing
	}
	service := NewBatchMovieConverterService(config)

	// Act
	result, err := service.batchConvert()

	// Assert
	if err != nil {
		ts.t.Errorf("Expected no error for recursive batchConvert, got %v", err)
	}
	if result == nil {
		ts.t.Error("Expected result object, got nil")
		return
	}
	// All 3 files should be processed recursively
	if result.TotalFiles != 3 {
		ts.t.Errorf("Expected TotalFiles to be 3 (recursive), got %d", result.TotalFiles)
	}
}

// #==============================================================#
// ##         Tests for scanFiles error handling                 ##
// #==============================================================#
// TestBatchMovieConverterServiceScanFilesErrors tests scanFiles error handling
type TestBatchMovieConverterServiceScanFilesErrors struct {
	t *testing.T
}

// NewTestBatchMovieConverterServiceScanFilesErrors creates a new test instance
func NewTestBatchMovieConverterServiceScanFilesErrors(t *testing.T) *TestBatchMovieConverterServiceScanFilesErrors {
	return &TestBatchMovieConverterServiceScanFilesErrors{t: t}
}

// TestBatchMovieConverterService_scanFiles_EmptyDirectory tests scanFiles with empty directory
func (ts *TestBatchMovieConverterServiceScanFilesErrors) TestBatchMovieConverterService_scanFiles_EmptyDirectory() {
	// Arrange
	tempDir := ts.t.TempDir()

	config := BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4",
		OutputDir: "/test/output",
		OutputExt: ".gif",
	}
	service := NewBatchMovieConverterService(config)

	// Act
	files, err := service.scanFiles()

	// Assert
	if err != nil {
		ts.t.Errorf("Expected no error for empty directory, got %v", err)
	}
	if len(files) != 0 {
		ts.t.Errorf("Expected 0 files in empty directory, got %d", len(files))
	}
}

// TestBatchMovieConverterService_scanFiles_CaseInsensitive tests case insensitive extension matching
func (ts *TestBatchMovieConverterServiceScanFilesErrors) TestBatchMovieConverterService_scanFiles_CaseInsensitive() {
	// Arrange
	tempDir := ts.t.TempDir()

	// Create files with different case extensions
	testFiles := []string{"video1.mp4", "video2.MP4", "video3.Mp4", "video4.mP4"}
	for _, fileName := range testFiles {
		file, err := os.Create(filepath.Join(tempDir, fileName))
		if err != nil {
			ts.t.Fatalf("Failed to create test file %s: %v", fileName, err)
		}
		file.Close()
	}

	config := BatchConversionConfig{
		InputDir:  tempDir,
		InputExt:  ".mp4", // lowercase
		OutputDir: "/test/output",
		OutputExt: ".gif",
	}
	service := NewBatchMovieConverterService(config)

	// Act
	files, err := service.scanFiles()

	// Assert
	if err != nil {
		ts.t.Errorf("Expected no error, got %v", err)
	}
	// All 4 files should be found regardless of case
	if len(files) != 4 {
		ts.t.Errorf("Expected 4 files (case insensitive), got %d", len(files))
	}
}

// Standard Go test functions for new batch tests

func TestBatchMovieConverterServiceBatchConvertWithFilesMethod(t *testing.T) {
	testService := NewTestBatchMovieConverterServiceBatchConvertWithFiles(t)
	testService.TestBatchMovieConverterService_batchConvert_WithFiles()
	testService.TestBatchMovieConverterService_batchConvert_OutputDirCreationFailure()
	testService.TestBatchMovieConverterService_batchConvert_MixedFileTypes()
	testService.TestBatchMovieConverterService_batchConvert_RecursiveProcessing()
}

func TestBatchMovieConverterServiceScanFilesErrorsMethod(t *testing.T) {
	testService := NewTestBatchMovieConverterServiceScanFilesErrors(t)
	testService.TestBatchMovieConverterService_scanFiles_EmptyDirectory()
	testService.TestBatchMovieConverterService_scanFiles_CaseInsensitive()
}
