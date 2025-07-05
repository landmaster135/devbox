package usecases

import (
	"strings"
	"testing"
)

// #==============================================================#
// ##          Tests                                             ##
// #==============================================================#
// TestDefaultREADMEParser_ParseUsageExamples_Normal は DefaultREADMEParser のテストです
func TestDefaultREADMEParser_ParseUsageExamples_Normal(t *testing.T) {
	t.Run("ParseUsageExamples_WithUsageSection", func(t *testing.T) {
		// Arrange
		parser := &DefaultREADMEParser{}
		content := []byte(`# Test Package

## 使用例

` + "```" + `
./test_package --help
./test_package input.txt
` + "```" + `

## その他のセクション
`)

		// Act
		result, err := parser.ParseUsageExamples(content)

		// Assert
		if err != nil {
			t.Fatalf("ParseUsageExamples() returned error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("Expected 2 usage examples, got %d", len(result))
		}
		expectedExamples := []string{
			"echo \\\"  ./test_package --help\\\"",
			"echo \\\"  ./test_package input.txt\\\"",
		}
		for i, expected := range expectedExamples {
			if result[i] != expected {
				t.Errorf("Expected usage example '%s', got '%s'", expected, result[i])
			}
		}
	})

	t.Run("ParseUsageExamples_WithUsageMethodSection", func(t *testing.T) {
		// Arrange
		parser := &DefaultREADMEParser{}
		content := []byte(`# Test Package

## 使用方法

` + "```" + `
./test_package --version
` + "```" + `
`)

		// Act
		result, err := parser.ParseUsageExamples(content)

		// Assert
		if err != nil {
			t.Fatalf("ParseUsageExamples() returned error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("Expected 1 usage example, got %d", len(result))
		}
		expected := "echo \\\"  ./test_package --version\\\""
		if result[0] != expected {
			t.Errorf("Expected usage example '%s', got '%s'", expected, result[0])
		}
	})

	t.Run("ParseUsageExamples_NoUsageSection", func(t *testing.T) {
		// Arrange
		parser := &DefaultREADMEParser{}
		content := []byte(`# Test Package

## インストール

` + "```" + `
go install ./...
` + "```" + `
`)

		// Act
		result, err := parser.ParseUsageExamples(content)

		// Assert
		if err != nil {
			t.Fatalf("ParseUsageExamples() returned error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected 0 usage examples, got %d", len(result))
		}
	})

	t.Run("ParseUsageExamples_EmptyContent", func(t *testing.T) {
		// Arrange
		parser := &DefaultREADMEParser{}
		content := []byte("")

		// Act
		result, err := parser.ParseUsageExamples(content)

		// Assert
		if err != nil {
			t.Fatalf("ParseUsageExamples() returned error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected 0 usage examples, got %d", len(result))
		}
	})

	t.Run("ParseUsageExamples_WithSubsection", func(t *testing.T) {
		// Arrange
		parser := &DefaultREADMEParser{}
		content := []byte(`# Test Package

## 使用例

### 基本的な使用例

` + "```" + `
./test_package basic
` + "```" + `

### 高度な使用例

` + "```" + `
./test_package advanced
` + "```" + `
`)

		// Act
		result, err := parser.ParseUsageExamples(content)

		// Assert
		if err != nil {
			t.Fatalf("ParseUsageExamples() returned error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("Expected 1 usage example (first code block only), got %d", len(result))
		}
		expected := "echo \\\"  ./test_package basic\\\""
		if result[0] != expected {
			t.Errorf("Expected usage example '%s', got '%s'", expected, result[0])
		}
	})
}

// TestDefaultScriptGenerator_GenerateContent_Normal は DefaultScriptGenerator のテストです
func TestDefaultScriptGenerator_GenerateContent_Normal(t *testing.T) {
	t.Run("GenerateContent_WithUsageExamples", func(t *testing.T) {
		// Arrange
		generator := &DefaultScriptGenerator{}
		packageName := "test-package"
		packagePath := "cmd/cli/test-package"
		usageExamples := []string{
			"echo \\\"  ./test_package --help\\\"",
			"echo \\\"  ./test_package input.txt\\\"",
		}

		// Act
		result := generator.GenerateContent(packageName, packagePath, usageExamples)

		// Assert
		if !strings.Contains(result, "#!/bin/bash") {
			t.Error("Expected script to contain shebang")
		}
		if !strings.Contains(result, "CMD_DIR=\"cmd/cli/test-package\"") {
			t.Error("Expected script to contain correct CMD_DIR")
		}
		if !strings.Contains(result, "OUTPUT_NAME=\"test-package\"") {
			t.Error("Expected script to contain correct OUTPUT_NAME")
		}
		if !strings.Contains(result, "Building ${OUTPUT_NAME}...") {
			t.Error("Expected script to contain build message")
		}
		if !strings.Contains(result, "GOOS=linux GOARCH=amd64") {
			t.Error("Expected script to contain Linux build")
		}
		if !strings.Contains(result, "GOOS=windows GOARCH=amd64") {
			t.Error("Expected script to contain Windows build")
		}
		if !strings.Contains(result, "GOOS=darwin GOARCH=arm64") {
			t.Error("Expected script to contain macOS build")
		}
		for _, example := range usageExamples {
			if !strings.Contains(result, example) {
				t.Errorf("Expected script to contain usage example: %s", example)
			}
		}
	})

	t.Run("GenerateContent_WithoutUsageExamples", func(t *testing.T) {
		// Arrange
		generator := &DefaultScriptGenerator{}
		packageName := "test-package"
		packagePath := "cmd/cli/test-package"
		usageExamples := []string{}

		// Act
		result := generator.GenerateContent(packageName, packagePath, usageExamples)

		// Assert
		if !strings.Contains(result, "#!/bin/bash") {
			t.Error("Expected script to contain shebang")
		}
		if !strings.Contains(result, "echo \\\"  ./test_package [options]\\\"") {
			t.Error("Expected script to contain default usage example")
		}
	})

	t.Run("GenerateContent_WithHyphenatedPackageName", func(t *testing.T) {
		// Arrange
		generator := &DefaultScriptGenerator{}
		packageName := "my-test-package"
		packagePath := "cmd/cli/my-test-package"
		usageExamples := []string{}

		// Act
		result := generator.GenerateContent(packageName, packagePath, usageExamples)

		// Assert
		if !strings.Contains(result, "OUTPUT_NAME=\"my-test-package\"") {
			t.Error("Expected script to contain correct OUTPUT_NAME with hyphens")
		}
		if !strings.Contains(result, "echo \\\"  ./my_test_package [options]\\\"") {
			t.Error("Expected default usage example to use underscores")
		}
	})
}

// TestOSFileSystem_Integration は OSFileSystem の統合テストです
func TestOSFileSystem_Integration(t *testing.T) {
	t.Run("OSFileSystem_Basic", func(t *testing.T) {
		// Arrange
		fs := &OSFileSystem{}

		// Act & Assert - 基本的な操作が実行できることを確認
		// 実際のファイルシステムを使用するため、存在するディレクトリでテスト
		entries, err := fs.ReadDir(".")
		if err != nil {
			t.Fatalf("ReadDir() returned error: %v", err)
		}
		if len(entries) == 0 {
			t.Error("Expected at least one entry in current directory")
		}

		// 現在のディレクトリの情報を取得
		info, err := fs.Stat(".")
		if err != nil {
			t.Fatalf("Stat() returned error: %v", err)
		}
		if !info.IsDir() {
			t.Error("Expected current directory to be a directory")
		}
	})
}

// TestStdinReader_Integration は StdinReader の統合テストです
func TestStdinReader_Integration(t *testing.T) {
	t.Run("StdinReader_Creation", func(t *testing.T) {
		// Arrange & Act
		reader := NewStdinReader()

		// Assert
		if reader == nil {
			t.Fatal("NewStdinReader() returned nil")
		}
		if reader.reader == nil {
			t.Error("Expected reader to be initialized")
		}
	})
}
