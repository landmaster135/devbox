package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMovieConverterService_Convert_MP4ToWEBM tests MP4 to WEBM conversion
func TestMovieConverterService_Convert_MP4ToWEBM_Normal(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "test.mp4")
	outputFile := filepath.Join(tempDir, "test.webm")

	// Create a dummy input file
	file, err := os.Create(inputFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:      inputFile,
		OutputFile:     outputFile,
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewMovieConverterService(config)

	// Note: This will fail due to missing ffmpeg, but we test the logic flow
	err = service.convert()

	// We expect an error due to missing ffmpeg, but the validation should pass
	if err == nil {
		t.Error("Expected error due to missing ffmpeg, got nil")
	}

	// Check that the error is related to ffmpeg execution, not validation
	if !strings.Contains(err.Error(), "ffmpeg") && !strings.Contains(err.Error(), "executable file not found") {
		t.Errorf("Expected ffmpeg-related error, got: %v", err)
	}
}

// TestMovieConverterService_Convert_WEBMToMP4 tests WEBM to MP4 conversion
func TestMovieConverterService_Convert_WEBMToMP4_Normal(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "test.webm")
	outputFile := filepath.Join(tempDir, "test.mp4")

	// Create a dummy input file
	file, err := os.Create(inputFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:  inputFile,
		OutputFile: outputFile,
	}

	service := NewMovieConverterService(config)

	// Note: This will fail due to missing ffmpeg, but we test the logic flow
	err = service.convert()

	// We expect an error due to missing ffmpeg, but the validation should pass
	if err == nil {
		t.Error("Expected error due to missing ffmpeg, got nil")
	}

	// Check that the error is related to ffmpeg execution, not validation
	if !strings.Contains(err.Error(), "ffmpeg") && !strings.Contains(err.Error(), "executable file not found") {
		t.Errorf("Expected ffmpeg-related error, got: %v", err)
	}
}

// TestMovieConverterService_Convert_UnsupportedConversion tests unsupported conversion
func TestMovieConverterService_Convert_UnsupportedConversion(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "test.txt")
	outputFile := filepath.Join(tempDir, "test.gif")

	// Create a dummy input file
	file, err := os.Create(inputFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:  inputFile,
		OutputFile: outputFile,
	}

	service := NewMovieConverterService(config)
	err = service.convert()

	if err == nil {
		t.Error("Expected error for unsupported conversion, got nil")
	}

	if !strings.Contains(err.Error(), "サポートされていない変換") {
		t.Errorf("Expected unsupported conversion error, got: %v", err)
	}
}

// TestMovieConverterService_Convert_MissingInputFile tests missing input file
func TestMovieConverterService_Convert_MissingInputFile(t *testing.T) {
	config := ConversionConfig{
		InputFile:  "non_existent_file.mp4",
		OutputFile: "test.webm",
	}

	service := NewMovieConverterService(config)
	err := service.convert()

	if err == nil {
		t.Error("Expected error for missing input file, got nil")
	}

	if !strings.Contains(err.Error(), "入力ファイルが見つかりません") {
		t.Errorf("Expected missing file error, got: %v", err)
	}
}

// TestMovieConverterService_Convert_NoExtension tests file without extension
func TestMovieConverterService_Convert_NoExtension(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "test_no_ext")
	outputFile := filepath.Join(tempDir, "test.webm")

	// Create a dummy input file
	file, err := os.Create(inputFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:  inputFile,
		OutputFile: outputFile,
	}

	service := NewMovieConverterService(config)
	err = service.convert()

	if err == nil {
		t.Error("Expected error for file without extension, got nil")
	}

	if !strings.Contains(err.Error(), "拡張子が含まれていません") {
		t.Errorf("Expected extension error, got: %v", err)
	}
}

// TestMovieConverterService_SetMP4ToWEBMDefaults_CRFMode tests CRF mode defaults
func TestMovieConverterService_SetMP4ToWEBMDefaults_CRFMode(t *testing.T) {
	config := ConversionConfig{
		InputFile:      "test.mp4",
		OutputFile:     "test.webm",
		ConversionMode: "crf",
		CRF:            0, // Should be set to default
	}

	service := NewMovieConverterService(config)
	service.setMP4ToWEBMDefaults()

	if service.config.CRF != 30 {
		t.Errorf("Expected CRF 30 for CRF mode, got %d", service.config.CRF)
	}
}

// TestMovieConverterService_SetMP4ToWEBMDefaults_CBRMode tests CBR mode defaults
func TestMovieConverterService_SetMP4ToWEBMDefaults_CBRMode(t *testing.T) {
	config := ConversionConfig{
		InputFile:      "test.mp4",
		OutputFile:     "test.webm",
		ConversionMode: "cbr",
		CRF:            0, // Should remain 0 for CBR mode
	}

	service := NewMovieConverterService(config)
	service.setMP4ToWEBMDefaults()

	if service.config.CRF != 0 {
		t.Errorf("Expected CRF 0 for CBR mode, got %d", service.config.CRF)
	}
}

// TestMovieConverterService_SetWEBMToMP4Defaults tests WEBM to MP4 defaults
func TestMovieConverterService_SetWEBMToMP4Defaults_Normal(t *testing.T) {
	config := ConversionConfig{
		InputFile:  "test.webm",
		OutputFile: "test.mp4",
	}

	service := NewMovieConverterService(config)
	service.setWEBMToMP4Defaults()

	// WEBM to MP4 conversion doesn't set special defaults
	// This test ensures the method doesn't crash
}

// TestMovieConverterService_GetSourceVideoBitrate_MissingFile tests bitrate detection with missing file
func TestMovieConverterService_GetSourceVideoBitrate_MissingFile(t *testing.T) {
	config := ConversionConfig{
		InputFile:    "non_existent_file.mp4",
		OutputFile:   "test.webm",
		VideoQuality: 75,
	}

	service := NewMovieConverterService(config)
	bitrate, err := service.getSourceVideoBitrate()

	// Should return default value without error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if bitrate != "1M" {
		t.Errorf("Expected default bitrate 1M, got %s", bitrate)
	}
}

// TestMovieConverterService_ConvertMP4ToWEBM_SpaceInFilename tests space in filename warning
func TestMovieConverterService_ConvertMP4ToWEBM_SpaceInFilename(t *testing.T) {
	// Create temporary test files with space in name
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "test file.mp4")
	outputFile := filepath.Join(tempDir, "test file.webm")

	// Create a dummy input file
	file, err := os.Create(inputFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:      inputFile,
		OutputFile:     outputFile,
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewMovieConverterService(config)

	// This will fail due to missing ffmpeg, but should process the space warning
	err = service.convertMP4ToWEBM()

	// We expect an error due to missing ffmpeg
	if err == nil {
		t.Error("Expected error due to missing ffmpeg, got nil")
	}
}

// TestMovieConverterService_ConvertMP4ToWEBM_VorbisCodec tests Vorbis audio codec
func TestMovieConverterService_ConvertMP4ToWEBM_VorbisCodec(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "test.mp4")
	outputFile := filepath.Join(tempDir, "test.webm")

	// Create a dummy input file
	file, err := os.Create(inputFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:      inputFile,
		OutputFile:     outputFile,
		AudioBitrate:   "128k",
		AudioCodec:     "vorbis", // Test Vorbis codec
		ConversionMode: "crf",
		CRF:            30,
		VideoQuality:   75,
	}

	service := NewMovieConverterService(config)

	// This will fail due to missing ffmpeg, but should process the codec setting
	err = service.convertMP4ToWEBM()

	// We expect an error due to missing ffmpeg
	if err == nil {
		t.Error("Expected error due to missing ffmpeg, got nil")
	}
}

// TestMovieConverterService_ConvertMP4ToWEBM_CBRMode tests CBR conversion mode
func TestMovieConverterService_ConvertMP4ToWEBM_CBRMode(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "test.mp4")
	outputFile := filepath.Join(tempDir, "test.webm")

	// Create a dummy input file
	file, err := os.Create(inputFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:      inputFile,
		OutputFile:     outputFile,
		VideoBitrate:   "2M",
		AudioBitrate:   "128k",
		AudioCodec:     "opus",
		ConversionMode: "cbr", // Test CBR mode
		VideoQuality:   75,
	}

	service := NewMovieConverterService(config)

	// This will fail due to missing ffmpeg, but should process CBR mode
	err = service.convertMP4ToWEBM()

	// We expect an error due to missing ffmpeg
	if err == nil {
		t.Error("Expected error due to missing ffmpeg, got nil")
	}
}

// TestMovieConverterService_ConvertWEBMToMP4_SpaceInFilename tests WEBM to MP4 with space in filename
func TestMovieConverterService_ConvertWEBMToMP4_SpaceInFilename(t *testing.T) {
	// Create temporary test files with space in name
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "test file.webm")
	outputFile := filepath.Join(tempDir, "test file.mp4")

	// Create a dummy input file
	file, err := os.Create(inputFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := ConversionConfig{
		InputFile:  inputFile,
		OutputFile: outputFile,
	}

	service := NewMovieConverterService(config)

	// This will fail due to missing ffmpeg, but should process the space warning
	err = service.convertWEBMToMP4()

	// We expect an error due to missing ffmpeg
	if err == nil {
		t.Error("Expected error due to missing ffmpeg, got nil")
	}
}

// TestMovieConverterService_Convert_AllSupportedInputFormats tests all supported input formats
func TestMovieConverterService_Convert_AllSupportedInputFormats(t *testing.T) {
	supportedExts := []string{".mp4", ".mkv", ".avi", ".mov", ".flv"}
	tempDir := t.TempDir()

	for _, ext := range supportedExts {
		inputFile := filepath.Join(tempDir, "test"+ext)
		outputFile := filepath.Join(tempDir, "test.webm")

		// Create a dummy input file
		file, err := os.Create(inputFile)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		file.Close()

		config := ConversionConfig{
			InputFile:  inputFile,
			OutputFile: outputFile,
		}

		service := NewMovieConverterService(config)
		err = service.convert()

		// We expect an error due to missing ffmpeg, but validation should pass
		if err == nil {
			t.Errorf("Expected error due to missing ffmpeg for %s, got nil", ext)
		}

		// Check that the error is related to ffmpeg execution, not validation
		if !strings.Contains(err.Error(), "ffmpeg") && !strings.Contains(err.Error(), "executable file not found") {
			t.Errorf("Expected ffmpeg-related error for %s, got: %v", ext, err)
		}
	}
}
