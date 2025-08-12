package usecases

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// #==============================================================#
// ##          Additional Coverage Tests                         ##
// #==============================================================#

// TestGolangOpsService_ExecuteTestCoverage_InvalidRegexPattern はExecuteTestCoverageの無効正規表現テストです
func TestGolangOpsService_ExecuteTestCoverage_InvalidRegexPattern(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	grepPattern := "[invalid"
	expectedOutput := []byte("ok  \ttest/package\tcoverage: 80.0% of statements")

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-cover", "./..."}).Return(expectedOutput, nil)

	// Act
	result, err := service.ExecuteTestCoverage(directory, grepPattern)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "出力のフィルタリングに失敗しました")
	assert.Contains(t, err.Error(), "無効な正規表現パターンです")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteTestCoverage_ExitCode2Error はExecuteTestCoverageのexit code 2エラーテストです
func TestGolangOpsService_ExecuteTestCoverage_ExitCode2Error(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	grepPattern := ""
	expectedOutput := []byte("go: cannot find main module")

	// exit status 2をシミュレート
	cmd := exec.Command("sh", "-c", "exit 2")
	err := cmd.Run()
	exitError, _ := err.(*exec.ExitError)

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-cover", "./..."}).Return(expectedOutput, exitError)

	// Act
	result, err := service.ExecuteTestCoverage(directory, grepPattern)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "テストカバレッジの実行でエラーが発生しました")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteTestCoverageProject_Step1ExitCode2Error はExecuteTestCoverageProjectのStep1 exit code 2エラーテストです
func TestGolangOpsService_ExecuteTestCoverageProject_Step1ExitCode2Error(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	step1Output := []byte("go: cannot find main module")

	// exit status 2をシミュレート
	cmd := exec.Command("sh", "-c", "exit 2")
	err := cmd.Run()
	exitError, _ := err.(*exec.ExitError)

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-coverprofile=coverage.out", "./..."}).Return(step1Output, exitError)

	// Act
	result, err := service.ExecuteTestCoverageProject(directory)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "テストカバレッジプロファイルの生成でエラーが発生しました")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteTestCoverageProject_Step3Error はExecuteTestCoverageProjectのStep3エラーテストです
func TestGolangOpsService_ExecuteTestCoverageProject_Step3Error(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	step1Output := []byte("ok  \ttest/package\tcoverage: 80.0% of statements")
	step2Output := []byte("test/package/file.go:10:\tfunc1\t\t\t80.0%")
	step3Error := fmt.Errorf("cannot create HTML file")

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-coverprofile=coverage.out", "./..."}).Return(step1Output, nil)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"tool", "cover", "-func=coverage.out"}).Return(step2Output, nil)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"tool", "cover", "-html=coverage.out", "-o", "coverage.html"}).Return([]byte(""), step3Error)

	// Act
	result, err := service.ExecuteTestCoverageProject(directory)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "HTMLカバレッジレポートの生成に失敗しました")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteTestCoverageProject_Step3WithOutput はExecuteTestCoverageProjectのStep3出力ありテストです
func TestGolangOpsService_ExecuteTestCoverageProject_Step3WithOutput(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	step1Output := []byte("ok  \ttest/package\tcoverage: 80.0% of statements")
	step2Output := []byte("test/package/file.go:10:\tfunc1\t\t\t80.0%")
	step3Output := []byte("HTML coverage report generated")

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-coverprofile=coverage.out", "./..."}).Return(step1Output, nil)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"tool", "cover", "-func=coverage.out"}).Return(step2Output, nil)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"tool", "cover", "-html=coverage.out", "-o", "coverage.html"}).Return(step3Output, nil)

	// Act
	result, err := service.ExecuteTestCoverageProject(directory)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, string(step3Output))
	assert.Contains(t, result, "HTMLレポートが生成されました")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteGoRun_RootDirectoryNotExists はExecuteGoRunのルートディレクトリ存在しないテストです
func TestGolangOpsService_ExecuteGoRun_RootDirectoryNotExists(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	executionFile := "/test/main.go"
	rootDirectory := "/nonexistent"
	parameters := ""

	mockDirectoryChecker.On("Exists", executionFile).Return(true)
	mockDirectoryChecker.On("Exists", rootDirectory).Return(false)

	// Act
	result, err := service.ExecuteGoRun(executionFile, rootDirectory, parameters)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "指定されたルートディレクトリが存在しません")
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteGoRun_RootDirectoryNotDirectory はExecuteGoRunのルートディレクトリがディレクトリでないテストです
func TestGolangOpsService_ExecuteGoRun_RootDirectoryNotDirectory(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	executionFile := "/test/main.go"
	rootDirectory := "/test/file.txt"
	parameters := ""

	mockDirectoryChecker.On("Exists", executionFile).Return(true)
	mockDirectoryChecker.On("Exists", rootDirectory).Return(true)
	mockDirectoryChecker.On("IsDirectory", rootDirectory).Return(false)

	// Act
	result, err := service.ExecuteGoRun(executionFile, rootDirectory, parameters)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "指定されたルートパスはディレクトリではありません")
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteCoverageFunc_InvalidRegexPattern はExecuteCoverageFuncの無効正規表現テストです
func TestGolangOpsService_ExecuteCoverageFunc_InvalidRegexPattern(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	coverageFile := "/test/dir/coverage.out"
	grepPattern := "[invalid"
	expectedOutput := []byte("test/package/file.go:10:\tfunc1\t\t\t80.0%")

	mockDirectoryChecker.On("Exists", coverageFile).Return(true)
	mockCommandExecutor.On("ExecuteInDir", "/test/dir", "go", []string{"tool", "cover", "-func=coverage.out"}).Return(expectedOutput, nil)

	// Act
	result, err := service.ExecuteCoverageFunc(coverageFile, grepPattern)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "出力のフィルタリングに失敗しました")
	assert.Contains(t, err.Error(), "無効な正規表現パターンです")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// #==============================================================#
// ##          Edge Case Tests                                   ##
// #==============================================================#

// TestFilterOutput_EmptyInput はfilterOutput関数の空入力テストです
func TestFilterOutput_EmptyInput(t *testing.T) {
	// Arrange
	input := []byte("")
	pattern := "test"

	// Act
	result, err := filterOutput(input, pattern)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "", string(result))
}

// TestFilterOutput_SingleLineMatch はfilterOutput関数の単一行マッチテストです
func TestFilterOutput_SingleLineMatch(t *testing.T) {
	// Arrange
	input := []byte("single line with test")
	pattern := "test"

	// Act
	result, err := filterOutput(input, pattern)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "single line with test", string(result))
}

// TestFilterOutput_MultipleMatches はfilterOutput関数の複数マッチテストです
func TestFilterOutput_MultipleMatches(t *testing.T) {
	// Arrange
	input := []byte("line1 test\nline2 no match\nline3 test\nline4 no match\nline5 test")
	pattern := "test"

	// Act
	result, err := filterOutput(input, pattern)

	// Assert
	assert.NoError(t, err)
	expected := "line1 test\nline3 test\nline5 test"
	assert.Equal(t, expected, string(result))
}

// TestFilterOutput_ComplexRegex はfilterOutput関数の複雑な正規表現テストです
func TestFilterOutput_ComplexRegex(t *testing.T) {
	// Arrange
	input := []byte("coverage: 80.0% of statements\ncoverage: 90.5% of statements\nno coverage info\ncoverage: 100.0% of statements")
	pattern := `coverage: \d+\.\d+% of statements`

	// Act
	result, err := filterOutput(input, pattern)

	// Assert
	assert.NoError(t, err)
	expected := "coverage: 80.0% of statements\ncoverage: 90.5% of statements\ncoverage: 100.0% of statements"
	assert.Equal(t, expected, string(result))
}

// #==============================================================#
// ##          Integration-like Tests                            ##
// #==============================================================#

// TestGolangOpsService_ExecuteTestCoverage_WithComplexGrepPattern はExecuteTestCoverageの複雑なgrepパターンテストです
func TestGolangOpsService_ExecuteTestCoverage_WithComplexGrepPattern(t *testing.T) {
	// Arrange
	mockCommandExecutor := &MockCommandExecutor{}
	mockDirectoryChecker := &MockDirectoryChecker{}
	service := NewGolangOpsServiceWithDependencies(mockCommandExecutor, mockDirectoryChecker)

	directory := "/test/dir"
	grepPattern := `coverage: \d+\.\d+%`
	expectedOutput := []byte("ok  \ttest/package\tcoverage: 80.0% of statements\nFAIL\ttest/other\t0.000s\nok  \ttest/another\tcoverage: 95.5% of statements")

	mockDirectoryChecker.On("Exists", directory).Return(true)
	mockDirectoryChecker.On("IsDirectory", directory).Return(true)
	mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-cover", "./..."}).Return(expectedOutput, nil)

	// Act
	result, err := service.ExecuteTestCoverage(directory, grepPattern)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "coverage: 80.0%")
	assert.Contains(t, result, "coverage: 95.5%")
	assert.NotContains(t, result, "FAIL\ttest/other")
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteGoRun_EmptyParameters はExecuteGoRunの空パラメータテストです
func TestGolangOpsService_ExecuteGoRun_EmptyParameters(t *testing.T) {
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
	assert.NotContains(t, result, "パラメータ:")
	assert.Contains(t, result, string(expectedOutput))
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}

// TestGolangOpsService_ExecuteCoverageFunc_EmptyGrepPattern はExecuteCoverageFuncの空grepパターンテストです
func TestGolangOpsService_ExecuteCoverageFunc_EmptyGrepPattern(t *testing.T) {
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
	assert.NotContains(t, result, "grepパターン:")
	assert.Contains(t, result, string(expectedOutput))
	mockCommandExecutor.AssertExpectations(t)
	mockDirectoryChecker.AssertExpectations(t)
}
