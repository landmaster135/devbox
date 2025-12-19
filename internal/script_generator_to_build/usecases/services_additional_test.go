package usecases

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	config "github.com/landmaster135/devbox/internal/script_generator_to_build/config"
)

// TestApp_selectPackage_InvalidInput は selectPackage メソッドの無効入力テストです
func TestApp_selectPackage_InvalidInput(t *testing.T) {
	t.Run("selectPackage_InvalidInput", func(t *testing.T) {
		// Arrange
		cfg := &config.ServiceConfig{}
		mockFS := &MockFileSystem{}
		mockReader := NewMockInputReader([]string{"invalid", "0", "999", "2"}) // 無効な入力の後に有効な入力

		packages := []string{"package-a", "package-b", "package-c"}
		mockEntries := make([]os.DirEntry, len(packages))
		for i, pkg := range packages {
			mockEntries[i] = NewMockDirEntry(pkg, true)
		}

		mockFS.ReadDirFunc = func(dirname string) ([]os.DirEntry, error) {
			return mockEntries, nil
		}

		app := NewServiceWithDependencies(cfg, mockFS, mockReader, nil, nil)
		w := &bytes.Buffer{}

		// Act
		result, err := app.selectPackage(w)

		// Assert
		if err != nil {
			t.Fatalf("selectPackage() returned error: %v", err)
		}
		if result != packages[1] {
			t.Errorf("Expected package '%s', got '%s'", packages[1], result)
		}
		output := w.String()
		if !strings.Contains(output, "無効な選択です") {
			t.Errorf("Expected output to contain invalid selection message, got: %s", output)
		}
	})
}

// TestApp_parseREADMEFile_ReadError は parseREADMEFile メソッドの読み取りエラーテストです
func TestApp_parseREADMEFile_ReadError(t *testing.T) {
	t.Run("parseREADMEFile_ReadError", func(t *testing.T) {
		// Arrange
		cfg := &config.ServiceConfig{}
		mockFS := &MockFileSystem{}

		mockFS.StatFunc = func(name string) (os.FileInfo, error) {
			if strings.Contains(name, "README.md") {
				return &MockFileInfo{name: "README.md", isDir: false}, nil
			}
			return nil, os.ErrNotExist
		}

		mockFS.ReadFileFunc = func(name string) ([]byte, error) {
			if strings.Contains(name, "README.md") {
				return nil, fmt.Errorf("file read error")
			}
			return nil, os.ErrNotExist
		}

		app := NewServiceWithDependencies(cfg, mockFS, nil, nil, nil)
		w := &bytes.Buffer{}

		// Act
		result, err := app.parseREADMEFile("test-package", w)

		// Assert
		if err == nil {
			t.Error("Expected error when file read fails")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty result when error occurs, got: %v", result)
		}
		if !strings.Contains(err.Error(), "READMEファイルの読み取りに失敗しました") {
			t.Errorf("Expected error message about file read failure, got: %s", err.Error())
		}
	})
}

// TestApp_parseREADMEFile_ParseError は parseREADMEFile メソッドの解析エラーテストです
func TestApp_parseREADMEFile_ParseError(t *testing.T) {
	t.Run("parseREADMEFile_ParseError", func(t *testing.T) {
		// Arrange
		cfg := &config.ServiceConfig{}
		mockFS := &MockFileSystem{}
		mockParser := &MockREADMEParser{}

		mockFS.StatFunc = func(name string) (os.FileInfo, error) {
			if strings.Contains(name, "README.md") {
				return &MockFileInfo{name: "README.md", isDir: false}, nil
			}
			return nil, os.ErrNotExist
		}

		mockFS.ReadFileFunc = func(name string) ([]byte, error) {
			if strings.Contains(name, "README.md") {
				return []byte("# Test Package"), nil
			}
			return nil, os.ErrNotExist
		}

		mockParser.ParseUsageExamplesFunc = func(content []byte) ([]string, error) {
			return nil, fmt.Errorf("parse error")
		}

		app := NewServiceWithDependencies(cfg, mockFS, nil, mockParser, nil)
		w := &bytes.Buffer{}

		// Act
		result, err := app.parseREADMEFile("test-package", w)

		// Assert
		if err == nil {
			t.Error("Expected error when parsing fails")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty result when error occurs, got: %v", result)
		}
		if !strings.Contains(err.Error(), "使用例の解析に失敗しました") {
			t.Errorf("Expected error message about parse failure, got: %s", err.Error())
		}
	})
}

// TestApp_writeScriptFile_MkdirError は writeScriptFile メソッドのディレクトリ作成エラーテストです
func TestApp_writeScriptFile_MkdirError(t *testing.T) {
	t.Run("writeScriptFile_MkdirError", func(t *testing.T) {
		// Arrange
		cfg := &config.ServiceConfig{}
		mockFS := &MockFileSystem{}

		mockFS.MkdirAllFunc = func(path string, perm os.FileMode) error {
			return fmt.Errorf("mkdir error")
		}

		app := NewServiceWithDependencies(cfg, mockFS, nil, nil, nil)
		w := &bytes.Buffer{}
		content := "#!/bin/bash\necho \"test script\""

		// Act
		err := app.writeScriptFile("test-package", content, w)

		// Assert
		if err == nil {
			t.Error("Expected error when mkdir fails")
		}
		if !strings.Contains(err.Error(), "スクリプトディレクトリの作成に失敗しました") {
			t.Errorf("Expected error message about mkdir failure, got: %s", err.Error())
		}
	})
}

// TestApp_writeScriptFile_WriteError は writeScriptFile メソッドのファイル書き込みエラーテストです
func TestApp_writeScriptFile_WriteError(t *testing.T) {
	t.Run("writeScriptFile_WriteError", func(t *testing.T) {
		// Arrange
		cfg := &config.ServiceConfig{}
		mockFS := &MockFileSystem{}

		mockFS.MkdirAllFunc = func(path string, perm os.FileMode) error {
			return nil
		}

		mockFS.WriteFileFunc = func(name string, data []byte, perm os.FileMode) error {
			return fmt.Errorf("write error")
		}

		app := NewServiceWithDependencies(cfg, mockFS, nil, nil, nil)
		w := &bytes.Buffer{}
		content := "#!/bin/bash\necho \"test script\""

		// Act
		err := app.writeScriptFile("test-package", content, w)

		// Assert
		if err == nil {
			t.Error("Expected error when file write fails")
		}
		if !strings.Contains(err.Error(), "ビルドスクリプトの書き込みに失敗しました") {
			t.Errorf("Expected error message about write failure, got: %s", err.Error())
		}
	})
}

// TestApp_generateBuildScript_Integration は generateBuildScript メソッドの統合テストです
func TestApp_generateBuildScript_Integration(t *testing.T) {
	t.Run("generateBuildScript_Integration", func(t *testing.T) {
		// Arrange
		cfg := &config.ServiceConfig{}
		mockFS := &MockFileSystem{}
		mockParser := &MockREADMEParser{}
		mockGenerator := &MockScriptGenerator{}

		// ファイルシステムのモック設定
		mockFS.StatFunc = func(name string) (os.FileInfo, error) {
			if strings.Contains(name, "test-package") {
				return &MockFileInfo{name: "test-package", isDir: true}, nil
			}
			if strings.Contains(name, "README.md") {
				return &MockFileInfo{name: "README.md", isDir: false}, nil
			}
			return nil, os.ErrNotExist
		}

		mockFS.ReadFileFunc = func(name string) ([]byte, error) {
			if strings.Contains(name, "README.md") {
				return []byte("# Test Package\n\n## 使用例\n\n```\n./test_package --help\n```\n"), nil
			}
			return nil, os.ErrNotExist
		}

		mockFS.MkdirAllFunc = func(path string, perm os.FileMode) error {
			return nil
		}

		var writtenContent []byte
		mockFS.WriteFileFunc = func(name string, data []byte, perm os.FileMode) error {
			writtenContent = data
			return nil
		}

		mockParser.ParseUsageExamplesFunc = func(content []byte) ([]string, error) {
			return []string{"echo \"  ./test_package --help\""}, nil
		}

		mockGenerator.GenerateContentFunc = func(packageName, packagePath string, usageExamples []string) string {
			return fmt.Sprintf("#!/bin/bash\necho \"Building %s...\"\n%s", packageName, strings.Join(usageExamples, "\n"))
		}

		app := NewServiceWithDependencies(cfg, mockFS, nil, mockParser, mockGenerator)
		w := &bytes.Buffer{}

		// Act
		err := app.generateBuildScript("test-package", w)

		// Assert
		if err != nil {
			t.Fatalf("generateBuildScript() returned error: %v", err)
		}
		if len(writtenContent) == 0 {
			t.Error("Expected script content to be written")
		}
		scriptContent := string(writtenContent)
		if !strings.Contains(scriptContent, "Building test-package...") {
			t.Errorf("Expected script to contain build message, got: %s", scriptContent)
		}
		if !strings.Contains(scriptContent, "echo \"  ./test_package --help\"") {
			t.Errorf("Expected script to contain usage example, got: %s", scriptContent)
		}
		output := w.String()
		if !strings.Contains(output, "READMEファイルを読み込みました") {
			t.Errorf("Expected output to contain README read message, got: %s", output)
		}
		if !strings.Contains(output, "ビルドスクリプトを生成しました") {
			t.Errorf("Expected output to contain script generation message, got: %s", output)
		}
	})
}

// TestApp_showHelp_Content は showHelp メソッドの内容テストです
func TestApp_showHelp_Content(t *testing.T) {
	t.Run("showHelp_Content", func(t *testing.T) {
		// Arrange
		cfg := &config.ServiceConfig{}
		app := NewService(cfg)
		w := &bytes.Buffer{}

		// Act
		app.showHelp(w)

		// Assert
		output := w.String()
		expectedContents := []string{
			"使用方法: script-generator-to-build [パッケージ名]",
			"このツールは、指定されたGoパッケージのビルドスクリプトを生成します。",
			"パッケージ名が指定されない場合は、利用可能なパッケージの一覧から選択できます。",
			"例:",
			"script-generator-to-build code-analyzer",
			"script-generator-to-build",
		}
		for _, expected := range expectedContents {
			if !strings.Contains(output, expected) {
				t.Errorf("Expected help output to contain '%s', got: %s", expected, output)
			}
		}
	})
}

// TestApp_getAvailablePackages_OnlyDirectories は getAvailablePackages メソッドのディレクトリフィルタリングテストです
func TestApp_getAvailablePackages_OnlyDirectories(t *testing.T) {
	t.Run("getAvailablePackages_OnlyDirectories", func(t *testing.T) {
		// Arrange
		cfg := &config.ServiceConfig{}
		mockFS := &MockFileSystem{}

		// ディレクトリとファイルが混在するエントリを作成
		mockEntries := []os.DirEntry{
			NewMockDirEntry("package-a", true),  // ディレクトリ
			NewMockDirEntry("file.txt", false),  // ファイル
			NewMockDirEntry("package-b", true),  // ディレクトリ
			NewMockDirEntry("README.md", false), // ファイル
			NewMockDirEntry("package-c", true),  // ディレクトリ
		}

		mockFS.ReadDirFunc = func(dirname string) ([]os.DirEntry, error) {
			return mockEntries, nil
		}

		app := NewServiceWithDependencies(cfg, mockFS, nil, nil, nil)

		// Act
		result, err := app.getAvailablePackages()

		// Assert
		if err != nil {
			t.Fatalf("getAvailablePackages() returned error: %v", err)
		}
		expectedPackages := []string{"package-a", "package-b", "package-c"}
		if len(result) != len(expectedPackages) {
			t.Errorf("Expected %d packages (directories only), got %d", len(expectedPackages), len(result))
		}
		for _, expected := range expectedPackages {
			found := false
			for _, r := range result {
				if r == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected package '%s' not found in result", expected)
			}
		}
		// ファイルが除外されていることを確認
		for _, r := range result {
			if r == "file.txt" || r == "README.md" {
				t.Errorf("Expected files to be excluded, but found '%s' in result", r)
			}
		}
	})
}
