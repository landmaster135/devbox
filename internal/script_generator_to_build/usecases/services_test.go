package usecases

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/landmaster135/devbox/internal/script_generator_to_build/config"
)

// #==============================================================#
// ##          Mocks                                             ##
// #==============================================================#
// MockFileSystem はテスト用のファイルシステムモックです
type MockFileSystem struct {
	ReadDirFunc   func(dirname string) ([]os.DirEntry, error)
	StatFunc      func(name string) (os.FileInfo, error)
	WriteFileFunc func(name string, data []byte, perm os.FileMode) error
	ReadFileFunc  func(name string) ([]byte, error)
	MkdirAllFunc  func(path string, perm os.FileMode) error
}

// ReadDir はディレクトリの内容を読み取ります
func (m *MockFileSystem) ReadDir(dirname string) ([]os.DirEntry, error) {
	if m.ReadDirFunc != nil {
		return m.ReadDirFunc(dirname)
	}
	return nil, nil
}

// Stat はファイル情報を取得します
func (m *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	if m.StatFunc != nil {
		return m.StatFunc(name)
	}
	return nil, nil
}

// WriteFile はファイルに書き込みます
func (m *MockFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(name, data, perm)
	}
	return nil
}

// ReadFile はファイルを読み取ります
func (m *MockFileSystem) ReadFile(name string) ([]byte, error) {
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(name)
	}
	return nil, nil
}

// MkdirAll はディレクトリを作成します
func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	if m.MkdirAllFunc != nil {
		return m.MkdirAllFunc(path, perm)
	}
	return nil
}

// MockInputReader はテスト用の入力読み取りモックです
type MockInputReader struct {
	ReadStringFunc func(delim byte) (string, error)
	inputs         []string
	index          int
}

// NewMockInputReader は新しい MockInputReader を作成します
func NewMockInputReader(inputs []string) *MockInputReader {
	return &MockInputReader{
		inputs: inputs,
		index:  0,
	}
}

// ReadString は文字列を読み取ります
func (m *MockInputReader) ReadString(delim byte) (string, error) {
	if m.ReadStringFunc != nil {
		return m.ReadStringFunc(delim)
	}

	if m.index >= len(m.inputs) {
		return "", nil
	}

	result := m.inputs[m.index] + string(delim)
	m.index++
	return result, nil
}

// MockREADMEParser はテスト用のREADME解析モックです
type MockREADMEParser struct {
	ParseUsageExamplesFunc func(content []byte) ([]string, error)
}

// ParseUsageExamples はREADMEファイルから使用例を抽出します
func (m *MockREADMEParser) ParseUsageExamples(content []byte) ([]string, error) {
	if m.ParseUsageExamplesFunc != nil {
		return m.ParseUsageExamplesFunc(content)
	}
	return []string{}, nil
}

// MockScriptGenerator はテスト用のスクリプト生成モックです
type MockScriptGenerator struct {
	GenerateContentFunc func(packageName, packagePath string, usageExamples []string) string
}

// GenerateContent はビルドスクリプトの内容を生成します
func (m *MockScriptGenerator) GenerateContent(packageName, packagePath string, usageExamples []string) string {
	if m.GenerateContentFunc != nil {
		return m.GenerateContentFunc(packageName, packagePath, usageExamples)
	}
	return "# Mock script content"
}

// MockDirEntry はテスト用のディレクトリエントリモックです
type MockDirEntry struct {
	name  string
	isDir bool
}

// NewMockDirEntry は新しい MockDirEntry を作成します
func NewMockDirEntry(name string, isDir bool) *MockDirEntry {
	return &MockDirEntry{
		name:  name,
		isDir: isDir,
	}
}

// Name はエントリ名を返します
func (m *MockDirEntry) Name() string {
	return m.name
}

// IsDir はディレクトリかどうかを返します
func (m *MockDirEntry) IsDir() bool {
	return m.isDir
}

// Type はファイルタイプを返します
func (m *MockDirEntry) Type() os.FileMode {
	if m.isDir {
		return os.ModeDir
	}
	return 0
}

// Info はファイル情報を返します
func (m *MockDirEntry) Info() (os.FileInfo, error) {
	return &MockFileInfo{
		name:  m.name,
		isDir: m.isDir,
	}, nil
}

// MockFileInfo はテスト用のファイル情報モックです
type MockFileInfo struct {
	name  string
	isDir bool
}

// Name はファイル名を返します
func (m *MockFileInfo) Name() string {
	return m.name
}

// Size はファイルサイズを返します
func (m *MockFileInfo) Size() int64 {
	return 0
}

// Mode はファイルモードを返します
func (m *MockFileInfo) Mode() os.FileMode {
	if m.isDir {
		return os.ModeDir
	}
	return 0
}

// ModTime は変更時刻を返します
func (m *MockFileInfo) ModTime() time.Time {
	return time.Time{}
}

// IsDir はディレクトリかどうかを返します
func (m *MockFileInfo) IsDir() bool {
	return m.isDir
}

// Sys はシステム固有の情報を返します
func (m *MockFileInfo) Sys() interface{} {
	return nil
}

// #==============================================================#
// ##          Tests                                             ##
// #==============================================================#
// TestNewApp_Normal は NewApp 関数のテストです
func TestNewApp_Normal(t *testing.T) {
	t.Run("NewApp_Normal", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{
			PackageName: "test-package",
			ShowHelp:    false,
		}

		// Act
		app := NewApp(cfg)

		// Assert
		if app == nil {
			t.Fatal("NewApp() returned nil")
		}
		if app.Config != cfg {
			t.Errorf("Expected config %v, got %v", cfg, app.Config)
		}
		// デフォルト値が設定されているか確認
		if app.Config.BaseDir == "" {
			t.Error("Expected BaseDir to be set by default")
		}
		if app.FileSystem == nil {
			t.Error("Expected FileSystem to be set")
		}
		if app.InputReader == nil {
			t.Error("Expected InputReader to be set")
		}
		if app.READMEParser == nil {
			t.Error("Expected READMEParser to be set")
		}
		if app.ScriptGenerator == nil {
			t.Error("Expected ScriptGenerator to be set")
		}
	})
}

// TestNewAppWithDependencies_Normal は NewAppWithDependencies 関数のテストです
func TestNewAppWithDependencies_Normal(t *testing.T) {
	t.Run("NewAppWithDependencies_Normal", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{
			PackageName: "test-package",
			ShowHelp:    false,
		}
		mockFS := &MockFileSystem{}
		mockReader := NewMockInputReader([]string{"1"})
		mockParser := &MockREADMEParser{}
		mockGenerator := &MockScriptGenerator{}

		// Act
		app := NewAppWithDependencies(cfg, mockFS, mockReader, mockParser, mockGenerator)

		// Assert
		if app == nil {
			t.Fatal("NewAppWithDependencies() returned nil")
		}
		if app.Config != cfg {
			t.Errorf("Expected config %v, got %v", cfg, app.Config)
		}
		if app.FileSystem != mockFS {
			t.Error("Expected FileSystem to be the provided mock")
		}
		if app.InputReader != mockReader {
			t.Error("Expected InputReader to be the provided mock")
		}
		if app.READMEParser != mockParser {
			t.Error("Expected READMEParser to be the provided mock")
		}
		if app.ScriptGenerator != mockGenerator {
			t.Error("Expected ScriptGenerator to be the provided mock")
		}
	})
}

// TestApp_Run_ShowHelp_Normal はRun メソッドのヘルプ表示テストです
func TestApp_Run_ShowHelp_Normal(t *testing.T) {
	t.Run("Run_ShowHelp_Normal", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{
			PackageName: "",
			ShowHelp:    true,
		}
		app := NewApp(cfg)
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		// Act
		exitCode := app.Run(stdout, stderr)

		// Assert
		if exitCode != ExitCodeOK {
			t.Errorf("Expected exit code %d, got %d", ExitCodeOK, exitCode)
		}
		output := stdout.String()
		if !strings.Contains(output, "使用方法") {
			t.Errorf("Expected help message to contain '使用方法', got: %s", output)
		}
	})
}

// TestApp_Run_WithPackageName_Normal はRun メソッドのパッケージ名指定テストです
func TestApp_Run_WithPackageName_Normal(t *testing.T) {
	t.Run("Run_WithPackageName_Normal", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{
			PackageName: "test-package",
			ShowHelp:    false,
		}

		// モックを設定
		mockFS := &MockFileSystem{}
		mockReader := NewMockInputReader([]string{})
		mockParser := &MockREADMEParser{}
		mockGenerator := &MockScriptGenerator{}

		// ファイルシステムのモック設定
		mockFS.StatFunc = func(name string) (os.FileInfo, error) {
			if strings.Contains(name, "test-package") {
				return &MockFileInfo{name: "test-package", isDir: true}, nil
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

		mockFS.WriteFileFunc = func(name string, data []byte, perm os.FileMode) error {
			return nil
		}

		mockParser.ParseUsageExamplesFunc = func(content []byte) ([]string, error) {
			return []string{"echo \"  ./test_package --help\""}, nil
		}

		mockGenerator.GenerateContentFunc = func(packageName, packagePath string, usageExamples []string) string {
			return fmt.Sprintf("#!/bin/bash\necho \"Building %s...\"\n", packageName)
		}

		app := NewAppWithDependencies(cfg, mockFS, mockReader, mockParser, mockGenerator)
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		// Act
		exitCode := app.Run(stdout, stderr)

		// Assert
		if exitCode != ExitCodeOK {
			t.Errorf("Expected exit code %d, got %d", ExitCodeOK, exitCode)
		}
		output := stdout.String()
		if !strings.Contains(output, "ビルドスクリプトを生成しました") {
			t.Errorf("Expected output to contain build script message, got: %s", output)
		}
	})
}

// TestApp_getAvailablePackages_Normal は getAvailablePackages メソッドのテストです
func TestApp_getAvailablePackages_Normal(t *testing.T) {
	t.Run("getAvailablePackages_Normal", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{}
		mockFS := &MockFileSystem{}

		packages := []string{"package-a", "package-b", "package-c"}
		mockEntries := make([]os.DirEntry, len(packages))
		for i, pkg := range packages {
			mockEntries[i] = NewMockDirEntry(pkg, true)
		}

		mockFS.ReadDirFunc = func(dirname string) ([]os.DirEntry, error) {
			return mockEntries, nil
		}

		app := NewAppWithDependencies(cfg, mockFS, nil, nil, nil)

		// Act
		result, err := app.getAvailablePackages()

		// Assert
		if err != nil {
			t.Fatalf("getAvailablePackages() returned error: %v", err)
		}
		if len(result) != len(packages) {
			t.Errorf("Expected %d packages, got %d", len(packages), len(result))
		}
		for _, pkg := range packages {
			found := false
			for _, r := range result {
				if r == pkg {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected package '%s' not found in result", pkg)
			}
		}
	})
}

// TestApp_getAvailablePackages_DirectoryNotFound は getAvailablePackages メソッドのエラーテストです
func TestApp_getAvailablePackages_DirectoryNotFound(t *testing.T) {
	t.Run("getAvailablePackages_DirectoryNotFound", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{}
		mockFS := &MockFileSystem{}

		mockFS.ReadDirFunc = func(dirname string) ([]os.DirEntry, error) {
			return nil, fmt.Errorf("ディレクトリの読み取りに失敗しました")
		}

		app := NewAppWithDependencies(cfg, mockFS, nil, nil, nil)

		// Act
		result, err := app.getAvailablePackages()

		// Assert
		if err == nil {
			t.Error("Expected error when directory does not exist")
		}
		if result != nil {
			t.Errorf("Expected nil result when error occurs, got: %v", result)
		}
	})
}

// TestApp_selectPackage_Normal は selectPackage メソッドのテストです
func TestApp_selectPackage_Normal(t *testing.T) {
	t.Run("selectPackage_Normal", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{}
		mockFS := &MockFileSystem{}
		mockReader := NewMockInputReader([]string{"2"}) // 2番目のパッケージを選択

		packages := []string{"package-a", "package-b", "package-c"}
		mockEntries := make([]os.DirEntry, len(packages))
		for i, pkg := range packages {
			mockEntries[i] = NewMockDirEntry(pkg, true)
		}

		mockFS.ReadDirFunc = func(dirname string) ([]os.DirEntry, error) {
			return mockEntries, nil
		}

		app := NewAppWithDependencies(cfg, mockFS, mockReader, nil, nil)
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
		if !strings.Contains(output, "利用可能なパッケージ:") {
			t.Errorf("Expected output to contain package list, got: %s", output)
		}
	})
}

// TestApp_selectPackage_NoPackagesAvailable は selectPackage メソッドのエラーテストです
func TestApp_selectPackage_NoPackagesAvailable(t *testing.T) {
	t.Run("selectPackage_NoPackagesAvailable", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{}
		mockFS := &MockFileSystem{}
		mockReader := NewMockInputReader([]string{})

		mockFS.ReadDirFunc = func(dirname string) ([]os.DirEntry, error) {
			return []os.DirEntry{}, nil // 空のリストを返す
		}

		app := NewAppWithDependencies(cfg, mockFS, mockReader, nil, nil)
		w := &bytes.Buffer{}

		// Act
		result, err := app.selectPackage(w)

		// Assert
		if err == nil {
			t.Error("Expected error when no packages are available")
		}
		if result != "" {
			t.Errorf("Expected empty result when error occurs, got: %s", result)
		}
		if !strings.Contains(err.Error(), "利用可能なパッケージが見つかりません") {
			t.Errorf("Expected error message about no packages, got: %s", err.Error())
		}
	})
}

// TestApp_validatePackage_Normal は validatePackage メソッドのテストです
func TestApp_validatePackage_Normal(t *testing.T) {
	t.Run("validatePackage_Normal", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{}
		mockFS := &MockFileSystem{}

		mockFS.StatFunc = func(name string) (os.FileInfo, error) {
			if strings.Contains(name, "test-package") {
				return &MockFileInfo{name: "test-package", isDir: true}, nil
			}
			return nil, os.ErrNotExist
		}

		app := NewAppWithDependencies(cfg, mockFS, nil, nil, nil)

		// Act
		err := app.validatePackage("test-package")

		// Assert
		if err != nil {
			t.Errorf("validatePackage() returned unexpected error: %v", err)
		}
	})
}

// TestApp_validatePackage_PackageNotFound は validatePackage メソッドのエラーテストです
func TestApp_validatePackage_PackageNotFound(t *testing.T) {
	t.Run("validatePackage_PackageNotFound", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{}
		mockFS := &MockFileSystem{}

		mockFS.StatFunc = func(name string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		}

		app := NewAppWithDependencies(cfg, mockFS, nil, nil, nil)

		// Act
		err := app.validatePackage("non-existent-package")

		// Assert
		if err == nil {
			t.Error("Expected error when package does not exist")
		}
		if !strings.Contains(err.Error(), "が見つかりません") {
			t.Errorf("Expected error message to contain 'が見つかりません', got: %s", err.Error())
		}
	})
}

// TestApp_parseREADMEFile_Normal は parseREADMEFile メソッドのテストです
func TestApp_parseREADMEFile_Normal(t *testing.T) {
	t.Run("parseREADMEFile_Normal", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{}
		mockFS := &MockFileSystem{}
		mockParser := &MockREADMEParser{}

		readmeContent := []byte("# Test Package\n\n## 使用例\n\n```\n./test_package --help\n```\n")
		expectedUsageExamples := []string{"echo \"  ./test_package --help\""}

		mockFS.StatFunc = func(name string) (os.FileInfo, error) {
			if strings.Contains(name, "README.md") {
				return &MockFileInfo{name: "README.md", isDir: false}, nil
			}
			return nil, os.ErrNotExist
		}

		mockFS.ReadFileFunc = func(name string) ([]byte, error) {
			if strings.Contains(name, "README.md") {
				return readmeContent, nil
			}
			return nil, os.ErrNotExist
		}

		mockParser.ParseUsageExamplesFunc = func(content []byte) ([]string, error) {
			return expectedUsageExamples, nil
		}

		app := NewAppWithDependencies(cfg, mockFS, nil, mockParser, nil)
		w := &bytes.Buffer{}

		// Act
		result, err := app.parseREADMEFile("test-package", w)

		// Assert
		if err != nil {
			t.Fatalf("parseREADMEFile() returned error: %v", err)
		}
		if len(result) != len(expectedUsageExamples) {
			t.Errorf("Expected %d usage examples, got %d", len(expectedUsageExamples), len(result))
		}
		for i, expected := range expectedUsageExamples {
			if result[i] != expected {
				t.Errorf("Expected usage example '%s', got '%s'", expected, result[i])
			}
		}
	})
}

// TestApp_parseREADMEFile_NoREADME は parseREADMEFile メソッドのREADMEなしテストです
func TestApp_parseREADMEFile_NoREADME(t *testing.T) {
	t.Run("parseREADMEFile_NoREADME", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{}
		mockFS := &MockFileSystem{}

		mockFS.StatFunc = func(name string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		}

		app := NewAppWithDependencies(cfg, mockFS, nil, nil, nil)
		w := &bytes.Buffer{}

		// Act
		result, err := app.parseREADMEFile("test-package", w)

		// Assert
		if err != nil {
			t.Fatalf("parseREADMEFile() returned unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty result when README does not exist, got: %v", result)
		}
	})
}

// TestApp_writeScriptFile_Normal は writeScriptFile メソッドのテストです
func TestApp_writeScriptFile_Normal(t *testing.T) {
	t.Run("writeScriptFile_Normal", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{}
		mockFS := &MockFileSystem{}

		var writtenPath string
		var writtenContent []byte

		mockFS.MkdirAllFunc = func(path string, perm os.FileMode) error {
			return nil
		}

		mockFS.WriteFileFunc = func(name string, data []byte, perm os.FileMode) error {
			writtenPath = name
			writtenContent = data
			return nil
		}

		app := NewAppWithDependencies(cfg, mockFS, nil, nil, nil)
		w := &bytes.Buffer{}
		content := "#!/bin/bash\necho \"test script\""

		// Act
		err := app.writeScriptFile("test-package", content, w)

		// Assert
		if err != nil {
			t.Fatalf("writeScriptFile() returned error: %v", err)
		}
		if !strings.Contains(writtenPath, "build_test_package.sh") {
			t.Errorf("Expected script path to contain 'build_test_package.sh', got: %s", writtenPath)
		}
		if string(writtenContent) != content {
			t.Errorf("Expected content '%s', got '%s'", content, string(writtenContent))
		}
		output := w.String()
		if !strings.Contains(output, "ビルドスクリプトを生成しました") {
			t.Errorf("Expected output to contain success message, got: %s", output)
		}
	})
}

// TestExitCodes は終了コードの定数をテストします
func TestExitCodes(t *testing.T) {
	t.Run("ExitCodes_Values", func(t *testing.T) {
		// Assert
		if ExitCodeOK != 0 {
			t.Errorf("Expected ExitCodeOK to be 0, got %d", ExitCodeOK)
		}
		if ExitCodeError != 1 {
			t.Errorf("Expected ExitCodeError to be 1, got %d", ExitCodeError)
		}
	})
}

// TestApp_Run_ErrorCases は Run メソッドのエラーケーステストです
func TestApp_Run_ErrorCases(t *testing.T) {
	t.Run("Run_PackageSelectionError", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{
			PackageName: "",
			ShowHelp:    false,
		}
		mockFS := &MockFileSystem{}
		mockReader := NewMockInputReader([]string{})

		mockFS.ReadDirFunc = func(dirname string) ([]os.DirEntry, error) {
			return nil, fmt.Errorf("directory read error")
		}

		app := NewAppWithDependencies(cfg, mockFS, mockReader, nil, nil)
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		// Act
		exitCode := app.Run(stdout, stderr)

		// Assert
		if exitCode != ExitCodeError {
			t.Errorf("Expected exit code %d, got %d", ExitCodeError, exitCode)
		}
	})

	t.Run("Run_GenerateBuildScriptError", func(t *testing.T) {
		// Arrange
		cfg := &config.AppConfig{
			PackageName: "non-existent-package",
			ShowHelp:    false,
		}
		mockFS := &MockFileSystem{}

		mockFS.StatFunc = func(name string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		}

		app := NewAppWithDependencies(cfg, mockFS, nil, nil, nil)
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		// Act
		exitCode := app.Run(stdout, stderr)

		// Assert
		if exitCode != ExitCodeError {
			t.Errorf("Expected exit code %d, got %d", ExitCodeError, exitCode)
		}
	})
}
