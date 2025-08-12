package usecases

import (
	"os"
	"strings"
	"testing"
)

// NewConfig関数のテスト
func TestNewConfig_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "new_config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	config, err := NewConfig(tempDir, "jpg,png", "Camera,DateTime", 10, true, false, true, false)
	if err != nil {
		t.Errorf("NewConfig should not return error for valid parameters, got: %v", err)
	}

	if config.Directory != tempDir {
		t.Errorf("Expected Directory %s, got %s", tempDir, config.Directory)
	}

	if len(config.Extensions) != 2 {
		t.Errorf("Expected 2 extensions, got %d", len(config.Extensions))
	}

	if config.Extensions[0] != "jpg" || config.Extensions[1] != "png" {
		t.Errorf("Expected extensions [jpg, png], got %v", config.Extensions)
	}

	if len(config.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(config.Properties))
	}

	if config.Properties[0] != "Camera" || config.Properties[1] != "DateTime" {
		t.Errorf("Expected properties [Camera, DateTime], got %v", config.Properties)
	}

	if config.MaxProps != 10 {
		t.Errorf("Expected MaxProps 10, got %d", config.MaxProps)
	}

	if !config.Verbose {
		t.Error("Expected Verbose to be true")
	}

	if config.Recursive {
		t.Error("Expected Recursive to be false")
	}

	if !config.ShowProperties {
		t.Error("Expected ShowProperties to be true")
	}

	if config.ShowDataTypes {
		t.Error("Expected ShowDataTypes to be false")
	}
}

func TestNewConfig_EmptyProperties(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "new_config_empty_props_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	config, err := NewConfig(tempDir, "jpg,png", "", 0, false, true, false, true)
	if err != nil {
		t.Errorf("NewConfig should not return error for empty properties, got: %v", err)
	}

	if config.Properties != nil {
		t.Errorf("Expected Properties to be nil for empty string, got %v", config.Properties)
	}
}

func TestNewConfig_InvalidDirectory(t *testing.T) {
	_, err := NewConfig("/non/existent/directory", "jpg,png", "", 0, false, false, false, false)
	if err == nil {
		t.Error("NewConfig should return error for non-existent directory")
	}

	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Expected error message to contain 'does not exist', got: %v", err)
	}
}

func TestNewConfig_EmptyExtensions(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "new_config_empty_ext_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	_, err = NewConfig(tempDir, "", "", 0, false, false, false, false)
	if err == nil {
		t.Error("NewConfig should return error for empty extensions")
	}

	if !strings.Contains(err.Error(), "extensions string cannot be empty") {
		t.Errorf("Expected error message to contain 'extensions string cannot be empty', got: %v", err)
	}
}

func TestNewConfig_NegativeMaxProps(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "new_config_negative_max_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	_, err = NewConfig(tempDir, "jpg,png", "", -1, false, false, false, false)
	if err == nil {
		t.Error("NewConfig should return error for negative maxProps")
	}

	if !strings.Contains(err.Error(), "maxProps cannot be negative") {
		t.Errorf("Expected error message to contain 'maxProps cannot be negative', got: %v", err)
	}
}

// validateConfigParams関数のテスト
func TestValidateConfigParams_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "validate_config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	err = validateConfigParams(tempDir, "jpg,png,tiff", 10)
	if err != nil {
		t.Errorf("validateConfigParams should not return error for valid parameters, got: %v", err)
	}
}

func TestValidateConfigParams_NonExistentDirectory(t *testing.T) {
	err := validateConfigParams("/non/existent/directory", "jpg,png", 0)
	if err == nil {
		t.Error("validateConfigParams should return error for non-existent directory")
	}

	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Expected error message to contain 'does not exist', got: %v", err)
	}
}

func TestValidateConfigParams_EmptyExtensions(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "validate_empty_ext_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	testCases := []string{
		"",
		"   ",
		"\t\n",
	}

	for _, extensionsStr := range testCases {
		err = validateConfigParams(tempDir, extensionsStr, 0)
		if err == nil {
			t.Errorf("validateConfigParams should return error for empty extensions: '%s'", extensionsStr)
		}

		if !strings.Contains(err.Error(), "extensions string cannot be empty") {
			t.Errorf("Expected error message to contain 'extensions string cannot be empty', got: %v", err)
		}
	}
}

func TestValidateConfigParams_NegativeMaxProps(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "validate_negative_max_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	testCases := []int{-1, -10, -100}

	for _, maxProps := range testCases {
		err = validateConfigParams(tempDir, "jpg,png", maxProps)
		if err == nil {
			t.Errorf("validateConfigParams should return error for negative maxProps: %d", maxProps)
		}

		if !strings.Contains(err.Error(), "maxProps cannot be negative") {
			t.Errorf("Expected error message to contain 'maxProps cannot be negative', got: %v", err)
		}
	}
}

func TestValidateConfigParams_EmptyExtensionInList(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "validate_empty_ext_in_list_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	testCases := []string{
		"jpg,,png",
		"jpg, ,png",
		",jpg,png",
		"jpg,png,",
		"jpg,\t,png",
	}

	for _, extensionsStr := range testCases {
		err = validateConfigParams(tempDir, extensionsStr, 0)
		if err == nil {
			t.Errorf("validateConfigParams should return error for empty extension in list: '%s'", extensionsStr)
		}

		if !strings.Contains(err.Error(), "extension at position") || !strings.Contains(err.Error(), "is empty") {
			t.Errorf("Expected error message to contain 'extension at position X is empty', got: %v", err)
		}
	}
}

func TestValidateConfigParams_ValidExtensionsWithSpaces(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "validate_ext_spaces_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// 空白を含む拡張子リストは正常に処理されるべき
	err = validateConfigParams(tempDir, " jpg , png , tiff ", 0)
	if err != nil {
		t.Errorf("validateConfigParams should not return error for extensions with spaces, got: %v", err)
	}
}

func TestValidateConfigParams_ZeroMaxProps(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "validate_zero_max_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// MaxProps = 0 は有効（制限なしを意味する）
	err = validateConfigParams(tempDir, "jpg,png", 0)
	if err != nil {
		t.Errorf("validateConfigParams should not return error for maxProps = 0, got: %v", err)
	}
}
