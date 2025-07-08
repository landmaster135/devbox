package usecases

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewMovieConverterService tests the creation of MovieConverterService
func TestNewMovieConverterService_Normal(t *testing.T) {
	config := ConversionConfig{
		InputFile:      "test.mp4",
		OutputFile:     "test.webm",
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewMovieConverterService(config)

	if service == nil {
		t.Fatal("NewMovieConverterService returned nil")
	}

	if service.config.InputFile != config.InputFile {
		t.Errorf("Expected InputFile %s, got %s", config.InputFile, service.config.InputFile)
	}

	if service.config.OutputFile != config.OutputFile {
		t.Errorf("Expected OutputFile %s, got %s", config.OutputFile, service.config.OutputFile)
	}
}

// TestSetMP4ToWEBMDefaults tests default value setting for MP4 to WEBM conversion
func TestMovieConverterService_SetMP4ToWEBMDefaults_Normal(t *testing.T) {
	config := ConversionConfig{
		InputFile:  "test.mp4",
		OutputFile: "test.webm",
	}

	service := NewMovieConverterService(config)
	service.setMP4ToWEBMDefaults()

	if service.config.AudioBitrate != "128k" {
		t.Errorf("Expected AudioBitrate 128k, got %s", service.config.AudioBitrate)
	}

	if service.config.AudioCodec != "opus" {
		t.Errorf("Expected AudioCodec opus, got %s", service.config.AudioCodec)
	}

	if service.config.ConversionMode != "crf" {
		t.Errorf("Expected ConversionMode crf, got %s", service.config.ConversionMode)
	}

	if service.config.CRF != 30 {
		t.Errorf("Expected CRF 30, got %d", service.config.CRF)
	}

	if service.config.VideoQuality != 75 {
		t.Errorf("Expected VideoQuality 75, got %d", service.config.VideoQuality)
	}
}

// TestSetMP4ToWEBMDefaults tests that existing values are not overwritten
func TestMovieConverterService_SetMP4ToWEBMDefaults_ExistingValues(t *testing.T) {
	config := ConversionConfig{
		InputFile:      "test.mp4",
		OutputFile:     "test.webm",
		AudioBitrate:   "256k",
		AudioCodec:     "vorbis",
		ConversionMode: "cbr",
		CRF:            25,
		VideoQuality:   80,
	}

	service := NewMovieConverterService(config)
	service.setMP4ToWEBMDefaults()

	if service.config.AudioBitrate != "256k" {
		t.Errorf("Expected AudioBitrate 256k, got %s", service.config.AudioBitrate)
	}

	if service.config.AudioCodec != "vorbis" {
		t.Errorf("Expected AudioCodec vorbis, got %s", service.config.AudioCodec)
	}

	if service.config.ConversionMode != "cbr" {
		t.Errorf("Expected ConversionMode cbr, got %s", service.config.ConversionMode)
	}

	if service.config.CRF != 25 {
		t.Errorf("Expected CRF 25, got %d", service.config.CRF)
	}

	if service.config.VideoQuality != 80 {
		t.Errorf("Expected VideoQuality 80, got %d", service.config.VideoQuality)
	}
}

// TestValidateSingleConfig tests configuration validation
func TestValidateSingleConfig_Normal(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.mp4")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:  testFile,
		OutputFile: "test.webm",
	}

	err = validateSingleConfig(config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestValidateSingleConfig tests validation with missing input file
func TestValidateSingleConfig_MissingInputFile(t *testing.T) {
	config := ConversionConfig{
		OutputFile: "test.webm",
	}

	err := validateSingleConfig(config)
	if err == nil {
		t.Error("Expected error for missing input file, got nil")
	}
}

// TestValidateSingleConfig tests validation with non-existent input file
func TestValidateSingleConfig_NonExistentFile(t *testing.T) {
	config := ConversionConfig{
		InputFile:  "non_existent_file.mp4",
		OutputFile: "test.webm",
	}

	err := validateSingleConfig(config)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

// TestValidateSingleConfig tests validation with input file without extension
func TestValidateSingleConfig_NoExtension(t *testing.T) {
	config := ConversionConfig{
		InputFile:  "test_file_without_extension",
		OutputFile: "test.webm",
	}

	err := validateSingleConfig(config)
	if err == nil {
		t.Error("Expected error for file without extension, got nil")
	}
}

// TestGenerateOutputFile tests output file name generation
func TestGenerateOutputFile_Normal(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"video.mp4", "video.webm"},
		{"video.mkv", "video.webm"},
		{"video.avi", "video.webm"},
		{"video.mov", "video.webm"},
		{"video.flv", "video.webm"},
		{"video.webm", "video.mp4"},
		{"video.unknown", "video_converted"},
	}

	for _, test := range tests {
		result := generateOutputFile(test.input)
		if result != test.expected {
			t.Errorf("For input %s, expected %s, got %s", test.input, test.expected, result)
		}
	}
}

// TestGetSupportedExtensions tests supported extensions
func TestGetSupportedExtensions_Normal(t *testing.T) {
	extensions := GetSupportedExtensions()

	expectedInput := []string{".mp4", ".mkv", ".avi", ".mov", ".flv", ".webm"}
	expectedOutput := []string{".mp4", ".webm"}

	inputExts := extensions["input"]
	outputExts := extensions["output"]

	if len(inputExts) != len(expectedInput) {
		t.Errorf("Expected %d input extensions, got %d", len(expectedInput), len(inputExts))
	}

	if len(outputExts) != len(expectedOutput) {
		t.Errorf("Expected %d output extensions, got %d", len(expectedOutput), len(outputExts))
	}

	for _, ext := range expectedInput {
		found := false
		for _, inputExt := range inputExts {
			if inputExt == ext {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected input extension %s not found", ext)
		}
	}

	for _, ext := range expectedOutput {
		found := false
		for _, outputExt := range outputExts {
			if outputExt == ext {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected output extension %s not found", ext)
		}
	}
}

// TestNormalizeExtension tests extension normalization
func TestNormalizeExtension_Normal(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"mp4", ".mp4"},
		{".mp4", ".mp4"},
		{"webm", ".webm"},
		{".webm", ".webm"},
		{"", ""},
	}

	for _, test := range tests {
		result := normalizeExtension(test.input)
		if result != test.expected {
			t.Errorf("For input %s, expected %s, got %s", test.input, test.expected, result)
		}
	}
}
