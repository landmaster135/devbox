package usecases

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// #==============================================================#
// ##          Mock Structures                                   ##
// #==============================================================#

// MockCommandExecutor はCommandExecutorのモック実装です
type MockCommandExecutor struct {
	mock.Mock
}

// Execute はコマンドを実行するモックメソッドです
func (m *MockCommandExecutor) Execute(name string, args ...string) ([]byte, error) {
	arguments := m.Called(name, args)
	return arguments.Get(0).([]byte), arguments.Error(1)
}

// ExecuteInDir は指定されたディレクトリでコマンドを実行するモックメソッドです
func (m *MockCommandExecutor) ExecuteInDir(dir, name string, args ...string) ([]byte, error) {
	arguments := m.Called(dir, name, args)
	return arguments.Get(0).([]byte), arguments.Error(1)
}

// MockDirectoryChecker はDirectoryCheckerのモック実装です
type MockDirectoryChecker struct {
	mock.Mock
}

// Exists はパスが存在するかチェックするモックメソッドです
func (m *MockDirectoryChecker) Exists(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

// IsDirectory はパスがディレクトリかチェックするモックメソッドです
func (m *MockDirectoryChecker) IsDirectory(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

// #==============================================================#
// ##          Test Helper Functions                             ##
// #==============================================================#

// createMockExitError はテスト用のExitErrorを作成します
func createMockExitError(exitCode int) *exec.ExitError {
	// 実際のExitErrorを作成するのは複雑なので、インターフェースを使用
	cmd := exec.Command("false") // 常に失敗するコマンド
	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError
	}
	return nil
}

// #==============================================================#
// ##          filterOutput Function Tests                       ##
// #==============================================================#

// TestFilterOutput_Normal はfilterOutput関数の正常系テストです
func TestFilterOutput_Normal(t *testing.T) {
	// Arrange
	input := []byte("line1\nline2 test\nline3\nline4 test")
	pattern := "test"

	// Act
	result, err := filterOutput(input, pattern)

	// Assert
	assert.NoError(t, err)
	expected := "line2 test\nline4 test"
	assert.Equal(t, expected, string(result))
}

// TestFilterOutput_EmptyPattern はfilterOutput関数の空パターンテストです
func TestFilterOutput_EmptyPattern(t *testing.T) {
	// Arrange
	input := []byte("line1\nline2\nline3")
	pattern := ""

	// Act
	result, err := filterOutput(input, pattern)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, input, result)
}

// TestFilterOutput_InvalidRegex はfilterOutput関数の無効な正規表現テストです
func TestFilterOutput_InvalidRegex(t *testing.T) {
	// Arrange
	input := []byte("line1\nline2")
	pattern := "[invalid"

	// Act
	result, err := filterOutput(input, pattern)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "無効な正規表現パターンです")
}

// TestFilterOutput_NoMatch はfilterOutput関数のマッチなしテストです
func TestFilterOutput_NoMatch(t *testing.T) {
	// Arrange
	input := []byte("line1\nline2\nline3")
	pattern := "nomatch"

	// Act
	result, err := filterOutput(input, pattern)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "", string(result))
}

// #==============================================================#
// ##          GolangOpsService Constructor Tests                ##
// #==============================================================#

// TestNewGolangOpsService_Normal はNewGolangOpsServiceの正常系テストです
func TestNewGolangOpsService_Normal(t *testing.T) {
	// Act
	service := NewGolangOpsService()

	// Assert
	assert.NotNil(t, service)
	assert.NotNil(t, service.commandExecutor)
	assert.NotNil(t, service.directoryChecker)
}

// TestNewGolangOpsServiceWithDependencies_Normal はNewGolangOpsServiceWithDependenciesの正常系テストです
func TestNewGolangOpsServiceWithDependencies_Normal(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}

	// Act
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	// Assert
	assert.NotNil(t, service)
	assert.Equal(t, mockCommandExecutor, service.commandExecutor)
	assert.Equal(t, mockDirectoryChecker, service.directoryChecker)
}

// #==============================================================#
// ##          ExecuteTestCoverage Tests                         ##
// #==============================================================#

// TestGolangOpsService_ExecuteTestCoverage_Normal はExecuteTestCoverageの正常系テストです
func TestGolangOpsService_ExecuteTestCoverage_Normal(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	grepPattern := ""
	expectedOutput := []byte("ok  \ttest/package\tcoverage: 80.0% of statements")

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-cover", "./..."}).Return(expectedOutput, nil)

	// Act
	result, err := service.ExecuteTestCoverage(directory, grepPattern)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "テストカバレッジを実行中")
	assert.Contains(t, result, string(expectedOutput))
	assert.Contains(t, result, "テストカバレッジの実行が完了しました")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteTestCoverage_WithGrepPattern はExecuteTestCoverageのgrepパターンありテストです
func TestGolangOpsService_ExecuteTestCoverage_WithGrepPattern(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	grepPattern := "coverage"
	expectedOutput := []byte("ok  \ttest/package\tcoverage: 80.0% of statements\nFAIL\ttest/other")

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-cover", "./..."}).Return(expectedOutput, nil)

	// Act
	result, err := service.ExecuteTestCoverage(directory, grepPattern)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "grepパターン: coverage")
	assert.Contains(t, result, "coverage: 80.0%")
	assert.NotContains(t, result, "FAIL\ttest/other")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteTestCoverage_TestFailure はExecuteTestCoverageのテスト失敗テストです
func TestGolangOpsService_ExecuteTestCoverage_TestFailure(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	grepPattern := ""
	expectedOutput := []byte("FAIL\ttest/package\t0.000s")

	// 実際のExitErrorを作成（exit status 1をシミュレート）
	cmd := exec.Command("sh", "-c", "exit 1")
	err := cmd.Run()
	exitError, _ := err.(*exec.ExitError)

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-cover", "./..."}).Return(expectedOutput, exitError)

	// Act
	result, err := service.ExecuteTestCoverage(directory, grepPattern)

	// Assert
	assert.NoError(t, err) // テスト失敗はエラーとして扱わない
	assert.Contains(t, result, string(expectedOutput))
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteTestCoverage_DirectoryNotExists はExecuteTestCoverageのディレクトリ存在しないテストです
func TestGolangOpsService_ExecuteTestCoverage_DirectoryNotExists(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/nonexistent/dir"
	grepPattern := ""

	mockDirectoryChecker.On("Exists", directory).Return(false)

	// Act
	result, err := service.ExecuteTestCoverage(directory, grepPattern)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "指定されたディレクトリが存在しません")
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteTestCoverage_NotDirectory はExecuteTestCoverageのディレクトリでないテストです
func TestGolangOpsService_ExecuteTestCoverage_NotDirectory(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/file.txt"
	grepPattern := ""

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(false)

	// Act
	result, err := service.ExecuteTestCoverage(directory, grepPattern)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "指定されたパスはディレクトリではありません")
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteTestCoverage_CommandError はExecuteTestCoverageのコマンドエラーテストです
func TestGolangOpsService_ExecuteTestCoverage_CommandError(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	grepPattern := ""
	commandError := fmt.Errorf("command not found")

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-cover", "./..."}).Return([]byte(""), commandError)

	// Act
	result, err := service.ExecuteTestCoverage(directory, grepPattern)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "コマンドの実行に失敗しました")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// #==============================================================#
// ##          ExecuteTestCoverageProject Tests                  ##
// #==============================================================#

// TestGolangOpsService_ExecuteTestCoverageProject_Normal はExecuteTestCoverageProjectの正常系テストです
func TestGolangOpsService_ExecuteTestCoverageProject_Normal(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	step1Output := []byte("ok  \ttest/package\tcoverage: 80.0% of statements")
	step2Output := []byte("test/package/file.go:10:\tfunc1\t\t\t80.0%")
	step3Output := []byte("")

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-coverprofile=coverage.out", "./..."}).Return(step1Output, nil)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"tool", "cover", "-func=coverage.out"}).Return(step2Output, nil)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"tool", "cover", "-html=coverage.out", "-o", "coverage.html"}).Return(step3Output, nil)

	// Act
	result, err := service.ExecuteTestCoverageProject(directory)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "プロジェクト全体のテストカバレッジを実行中")
	assert.Contains(t, result, "Step 1: テストカバレッジプロファイルを生成中")
	assert.Contains(t, result, "Step 2: カバレッジ関数レポートを生成中")
	assert.Contains(t, result, "Step 3: HTMLカバレッジレポートを生成中")
	assert.Contains(t, result, string(step1Output))
	assert.Contains(t, result, string(step2Output))
	assert.Contains(t, result, "HTMLレポートが生成されました")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteTestCoverageProject_Step1TestFailure はExecuteTestCoverageProjectのStep1テスト失敗テストです
func TestGolangOpsService_ExecuteTestCoverageProject_Step1TestFailure(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	step1Output := []byte("FAIL\ttest/package\t0.000s")
	step2Output := []byte("test/package/file.go:10:\tfunc1\t\t\t0.0%")
	step3Output := []byte("")

	// 実際のExitErrorを作成（exit status 1をシミュレート）
	cmd := exec.Command("sh", "-c", "exit 1")
	err := cmd.Run()
	exitError, _ := err.(*exec.ExitError)

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-coverprofile=coverage.out", "./..."}).Return(step1Output, exitError)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"tool", "cover", "-func=coverage.out"}).Return(step2Output, nil)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"tool", "cover", "-html=coverage.out", "-o", "coverage.html"}).Return(step3Output, nil)

	// Act
	result, err := service.ExecuteTestCoverageProject(directory)

	// Assert
	assert.NoError(t, err) // Step1のテスト失敗はエラーとして扱わない
	assert.Contains(t, result, string(step1Output))
	assert.Contains(t, result, string(step2Output))
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteTestCoverageProject_Step2Error はExecuteTestCoverageProjectのStep2エラーテストです
func TestGolangOpsService_ExecuteTestCoverageProject_Step2Error(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	step1Output := []byte("ok  \ttest/package\tcoverage: 80.0% of statements")
	step2Error := fmt.Errorf("coverage file not found")

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-coverprofile=coverage.out", "./..."}).Return(step1Output, nil)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"tool", "cover", "-func=coverage.out"}).Return([]byte(""), step2Error)

	// Act
	result, err := service.ExecuteTestCoverageProject(directory)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "カバレッジ関数レポートの生成に失敗しました")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// #==============================================================#
// ##          ExecuteGoRun Tests                                ##
// #==============================================================#

// TestGolangOpsService_ExecuteGoRun_Normal はExecuteGoRunの正常系テストです
func TestGolangOpsService_ExecuteGoRun_Normal(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	executionFile := "/test/main.go"
	rootDirectory := "/test"
	parameters := ""
	expectedOutput := []byte("Hello, World!")

	mockDirectoryChecker.On("Exists", executionFile).Return(true)
	mockDirectoryChecker.On("Exists", rootDirectory).Return(true)
	mockDirectoryChecker.On("IsDirectory", rootDirectory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"run", executionFile}).Return(expectedOutput, nil)

	// Act
	result, err := service.ExecuteGoRun(executionFile, rootDirectory, parameters)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "go runを実行中")
	assert.Contains(t, result, "実行ディレクトリ:")
	assert.Contains(t, result, string(expectedOutput))
	assert.Contains(t, result, "go runの実行が完了しました")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteGoRun_WithParameters はExecuteGoRunのパラメータありテストです
func TestGolangOpsService_ExecuteGoRun_WithParameters(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	executionFile := "/test/main.go"
	rootDirectory := "/test"
	parameters := "-flag value arg1 arg2"
	expectedOutput := []byte("Hello with parameters!")

	mockDirectoryChecker.On("Exists", executionFile).Return(true)
	mockDirectoryChecker.On("Exists", rootDirectory).Return(true)
	mockDirectoryChecker.On("IsDirectory", rootDirectory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"run", executionFile, "-flag", "value", "arg1", "arg2"}).Return(expectedOutput, nil)

	// Act
	result, err := service.ExecuteGoRun(executionFile, rootDirectory, parameters)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "パラメータ: "+parameters)
	assert.Contains(t, result, string(expectedOutput))
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteGoRun_FileNotExists はExecuteGoRunのファイル存在しないテストです
func TestGolangOpsService_ExecuteGoRun_FileNotExists(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	executionFile := "/nonexistent/main.go"
	rootDirectory := "/test"
	parameters := ""

	mockDirectoryChecker.On("Exists", executionFile).Return(false)

	// Act
	result, err := service.ExecuteGoRun(executionFile, rootDirectory, parameters)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "指定された実行ファイルが存在しません")
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteGoRun_CommandError はExecuteGoRunのコマンドエラーテストです
func TestGolangOpsService_ExecuteGoRun_CommandError(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	executionFile := "/test/main.go"
	rootDirectory := "/test"
	parameters := ""
	commandError := fmt.Errorf("compilation error")

	mockDirectoryChecker.On("Exists", executionFile).Return(true)
	mockDirectoryChecker.On("Exists", rootDirectory).Return(true)
	mockDirectoryChecker.On("IsDirectory", rootDirectory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"run", executionFile}).Return([]byte(""), commandError)

	// Act
	result, err := service.ExecuteGoRun(executionFile, rootDirectory, parameters)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "go runの実行に失敗しました")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// #==============================================================#
// ##          ExecuteCoverageFunc Tests                         ##
// #==============================================================#

// TestGolangOpsService_ExecuteCoverageFunc_Normal はExecuteCoverageFuncの正常系テストです
func TestGolangOpsService_ExecuteCoverageFunc_Normal(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	coverageFile := "/test/dir/coverage.out"
	grepPattern := ""
	expectedOutput := []byte("test/package/file.go:10:\tfunc1\t\t\t80.0%\ntotal:\t\t\t\t\t(statements)\t80.0%")

	mockDirectoryChecker.On("Exists", coverageFile).Return(true)
	mockCommandExecutor.On("ExecuteInDir", "/test/dir", "go", []string{"tool", "cover", "-func=coverage.out"}).Return(expectedOutput, nil)

	// Act
	result, err := service.ExecuteCoverageFunc(coverageFile, grepPattern)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "カバレッジファイルから関数情報を取得中")
	assert.Contains(t, result, string(expectedOutput))
	assert.Contains(t, result, "カバレッジ関数情報の取得が完了しました")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteCoverageFunc_WithGrepPattern はExecuteCoverageFuncのgrepパターンありテストです
func TestGolangOpsService_ExecuteCoverageFunc_WithGrepPattern(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	coverageFile := "/test/dir/coverage.out"
	grepPattern := "total"
	expectedOutput := []byte("test/package/file.go:10:\tfunc1\t\t\t80.0%\ntotal:\t\t\t\t\t(statements)\t80.0%")

	mockDirectoryChecker.On("Exists", coverageFile).Return(true)
	mockCommandExecutor.On("ExecuteInDir", "/test/dir", "go", []string{"tool", "cover", "-func=coverage.out"}).Return(expectedOutput, nil)

	// Act
	result, err := service.ExecuteCoverageFunc(coverageFile, grepPattern)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "grepパターン: total")
	assert.Contains(t, result, "total:")
	assert.NotContains(t, result, "func1")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteCoverageFunc_FileNotExists はExecuteCoverageFuncのファイル存在しないテストです
func TestGolangOpsService_ExecuteCoverageFunc_FileNotExists(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	coverageFile := "/nonexistent/coverage.out"
	grepPattern := ""

	mockDirectoryChecker.On("Exists", coverageFile).Return(false)

	// Act
	result, err := service.ExecuteCoverageFunc(coverageFile, grepPattern)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "指定されたカバレッジファイルが存在しません")
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteCoverageFunc_CommandError はExecuteCoverageFuncのコマンドエラーテストです
func TestGolangOpsService_ExecuteCoverageFunc_CommandError(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	coverageFile := "/test/dir/coverage.out"
	grepPattern := ""
	commandError := fmt.Errorf("coverage file format error")

	mockDirectoryChecker.On("Exists", coverageFile).Return(true)
	mockCommandExecutor.On("ExecuteInDir", "/test/dir", "go", []string{"tool", "cover", "-func=coverage.out"}).Return([]byte(""), commandError)

	// Act
	result, err := service.ExecuteCoverageFunc(coverageFile, grepPattern)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "カバレッジ関数情報の取得に失敗しました")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// #==============================================================#
// ##          Handler Method Tests                              ##
// #==============================================================#

// TestGolangOpsService_HandleTestCoverage_Normal はHandleTestCoverageの正常系テストです
func TestGolangOpsService_HandleTestCoverage_Normal(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	grepPattern := ""
	expectedOutput := []byte("ok  \ttest/package\tcoverage: 80.0% of statements")

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-cover", "./..."}).Return(expectedOutput, nil)

	// Act
	result, err := service.HandleTestCoverage(directory, grepPattern)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "テストカバレッジを実行中")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_HandleTestCoverageProject_Normal はHandleTestCoverageProjectの正常系テストです
func TestGolangOpsService_HandleTestCoverageProject_Normal(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	step1Output := []byte("ok  \ttest/package\tcoverage: 80.0% of statements")
	step2Output := []byte("test/package/file.go:10:\tfunc1\t\t\t80.0%")
	step3Output := []byte("")

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-coverprofile=coverage.out", "./..."}).Return(step1Output, nil)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"tool", "cover", "-func=coverage.out"}).Return(step2Output, nil)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"tool", "cover", "-html=coverage.out", "-o", "coverage.html"}).Return(step3Output, nil)

	// Act
	result, err := service.HandleTestCoverageProject(directory)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "プロジェクト全体のテストカバレッジを実行中")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_HandleGoRun_Normal はHandleGoRunの正常系テストです
func TestGolangOpsService_HandleGoRun_Normal(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	executionFile := "/test/main.go"
	rootDirectory := "/test"
	parameters := ""
	expectedOutput := []byte("Hello, World!")

	mockDirectoryChecker.On("Exists", executionFile).Return(true)
	mockDirectoryChecker.On("Exists", rootDirectory).Return(true)
	mockDirectoryChecker.On("IsDirectory", rootDirectory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"run", executionFile}).Return(expectedOutput, nil)

	// Act
	result, err := service.HandleGoRun(executionFile, rootDirectory, parameters)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "go runを実行中")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_HandleCoverageFunc_Normal はHandleCoverageFuncの正常系テストです
func TestGolangOpsService_HandleCoverageFunc_Normal(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	coverageFile := "/test/dir/coverage.out"
	grepPattern := ""
	expectedOutput := []byte("test/package/file.go:10:\tfunc1\t\t\t80.0%")

	mockDirectoryChecker.On("Exists", coverageFile).Return(true)
	mockCommandExecutor.On("ExecuteInDir", "/test/dir", "go", []string{"tool", "cover", "-func=coverage.out"}).Return(expectedOutput, nil)

	// Act
	result, err := service.HandleCoverageFunc(coverageFile, grepPattern)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "カバレッジファイルから関数情報を取得中")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// #==============================================================#
// ##          Default Implementation Tests                      ##
// #==============================================================#

// TestDefaultCommandExecutor_Execute はDefaultCommandExecutorのExecuteメソッドテストです
func TestDefaultCommandExecutor_Execute_Normal(t *testing.T) {
	// Arrange
	executor := &DefaultCommandExecutor{}

	// Act
	output, err := executor.Execute("echo", "test")

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, string(output), "test")
}

// TestDefaultCommandExecutor_ExecuteInDir はDefaultCommandExecutorのExecuteInDirメソッドテストです
func TestDefaultCommandExecutor_ExecuteInDir_Normal(t *testing.T) {
	// Arrange
	executor := &DefaultCommandExecutor{}

	// Act
	output, err := executor.ExecuteInDir("/tmp", "pwd")

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, string(output), "/tmp")
}

// TestDefaultDirectoryChecker_Exists はDefaultDirectoryCheckerのExistsメソッドテストです
func TestDefaultDirectoryChecker_Exists_Normal(t *testing.T) {
	// Arrange
	checker := &DefaultDirectoryChecker{}

	// Act & Assert
	assert.True(t, checker.Exists("/tmp"))
	assert.False(t, checker.Exists("/nonexistent/path"))
}

// TestDefaultDirectoryChecker_IsDirectory はDefaultDirectoryCheckerのIsDirectoryメソッドテストです
func TestDefaultDirectoryChecker_IsDirectory_Normal(t *testing.T) {
	// Arrange
	checker := &DefaultDirectoryChecker{}

	// Act & Assert
	assert.True(t, checker.IsDirectory("/tmp"))
	assert.False(t, checker.IsDirectory("/nonexistent/path"))
}
