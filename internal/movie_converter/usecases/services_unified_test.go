package usecases

import (
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
