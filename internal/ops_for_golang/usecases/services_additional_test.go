package usecases

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
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

	mockDirectoryChecker.ExistsFunc = func(path string) bool {
		return path == directory
	}
	mockDirectoryChecker.IsDirectoryFunc = func(path string) bool {
		return path == directory
	}
	mockCommandExecutor.ExecuteInDirFunc = func(dir, name string, args ...string) ([]byte, error) {
		if name == "go" && len(args) >= 2 && args[0] == "test" && args[1] == "-cover" {
			return expectedOutput, nil
		}
		return nil, fmt.Errorf("unexpected command")
	}

	// Act
	result, err := service.ExecuteTestCoverage(directory, grepPattern)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "出力のフィルタリングに失敗しました")
	assert.Contains(t, err.Error(), "無効な正規表現パターンです")
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

	mockDirectoryChecker.ExistsFunc = func(path string) bool {
		return path == directory
	}
	mockDirectoryChecker.IsDirectoryFunc = func(path string) bool {
		return path == directory
	}
	mockCommandExecutor.ExecuteInDirFunc = func(dir, name string, args ...string) ([]byte, error) {
		if name == "go" && len(args) >= 2 && args[0] == "test" && args[1] == "-cover" {
			return expectedOutput, exitError
		}
		return nil, fmt.Errorf("unexpected command")
	}

	// Act
	result, err := service.ExecuteTestCoverage(directory, grepPattern)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "テストカバレッジの実行でエラーが発生しました")
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

	mockDirectoryChecker.ExistsFunc = func(path string) bool {
		return path == directory
	}
	mockDirectoryChecker.IsDirectoryFunc = func(path string) bool {
		return path == directory
	}
	mockCommandExecutor.ExecuteInDirFunc = func(dir, name string, args ...string) ([]byte, error) {
		if name == "go" && len(args) >= 2 && args[0] == "test" && args[1] == "-coverprofile=coverage.out" {
			return step1Output, exitError
		}
		return nil, fmt.Errorf("unexpected command")
	}

	// Act
	result, err := service.ExecuteTestCoverageProject(directory)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "テストカバレッジプロファイルの生成でエラーが発生しました")
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

	mockDirectoryChecker.ExistsFunc = func(path string) bool {
		return path == directory
	}
	mockDirectoryChecker.IsDirectoryFunc = func(path string) bool {
		return path == directory
	}
	mockCommandExecutor.ExecuteInDirFunc = func(dir, name string, args ...string) ([]byte, error) {
		if name == "go" && len(args) >= 2 {
			if args[0] == "test" && args[1] == "-coverprofile=coverage.out" {
				return step1Output, nil
			} else if args[0] == "tool" && args[1] == "cover" && args[2] == "-func=coverage.out" {
				return step2Output, nil
			} else if args[0] == "tool" && args[1] == "cover" && args[2] == "-html=coverage.out" {
				return []byte(""), step3Error
			}
		}
		return nil, fmt.Errorf("unexpected command")
	}

	// Act
	result, err := service.ExecuteTestCoverageProject(directory)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "HTMLカバレッジレポートの生成に失敗しました")
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

	mockDirectoryChecker.ExistsFunc = func(path string) bool {
		return path == directory
	}
	mockDirectoryChecker.IsDirectoryFunc = func(path string) bool {
		return path == directory
	}
	mockCommandExecutor.ExecuteInDirFunc = func(dir, name string, args ...string) ([]byte, error) {
		if name == "go" && len(args) >= 2 {
			if args[0] == "test" && args[1] == "-coverprofile=coverage.out" {
				return step1Output, nil
			} else if args[0] == "tool" && args[1] == "cover" && args[2] == "-func=coverage.out" {
				return step2Output, nil
			} else if args[0] == "tool" && args[1] == "cover" && args[2] == "-html=coverage.out" {
				return step3Output, nil
			}
		}
		return nil, fmt.Errorf("unexpected command")
	}

	// Act
	result, err := service.ExecuteTestCoverageProject(directory)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, string(step3Output))
	assert.Contains(t, result, "HTMLレポートが生成されました")
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

	mockDirectoryChecker.ExistsFunc = func(path string) bool {
		if path == executionFile {
			return true
		}
		if path == rootDirectory {
			return false
		}
		return false
	}

	// Act
	result, err := service.ExecuteGoRun(executionFile, rootDirectory, parameters)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "指定されたルートディレクトリが存在しません")
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

	mockDirectoryChecker.ExistsFunc = func(path string) bool {
		return path == executionFile || path == rootDirectory
	}
	mockDirectoryChecker.IsDirectoryFunc = func(path string) bool {
		return path != rootDirectory
	}

	// Act
	result, err := service.ExecuteGoRun(executionFile, rootDirectory, parameters)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "指定されたルートパスはディレクトリではありません")
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

	mockDirectoryChecker.ExistsFunc = func(path string) bool {
		return path == coverageFile
	}
	mockCommandExecutor.ExecuteInDirFunc = func(dir, name string, args ...string) ([]byte, error) {
		if name == "go" && len(args) >= 3 && args[0] == "tool" && args[1] == "cover" && args[2] == "-func=coverage.out" {
			return expectedOutput, nil
		}
		return nil, fmt.Errorf("unexpected command")
	}

	// Act
	result, err := service.ExecuteCoverageFunc(coverageFile, grepPattern)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "出力のフィルタリングに失敗しました")
	assert.Contains(t, err.Error(), "無効な正規表現パターンです")
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

	mockDirectoryChecker.ExistsFunc = func(path string) bool {
		return path == directory
	}
	mockDirectoryChecker.IsDirectoryFunc = func(path string) bool {
		return path == directory
	}
	mockCommandExecutor.ExecuteInDirFunc = func(dir, name string, args ...string) ([]byte, error) {
		if name == "go" && len(args) >= 2 && args[0] == "test" && args[1] == "-cover" {
			return expectedOutput, nil
		}
		return nil, fmt.Errorf("unexpected command")
	}

	// Act
	result, err := service.ExecuteTestCoverage(directory, grepPattern)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "coverage: 80.0%")
	assert.Contains(t, result, "coverage: 95.5%")
	assert.NotContains(t, result, "FAIL\ttest/other")
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

	mockDirectoryChecker.ExistsFunc = func(path string) bool {
		return path == executionFile || path == rootDirectory
	}
	mockDirectoryChecker.IsDirectoryFunc = func(path string) bool {
		return path == rootDirectory
	}
	mockCommandExecutor.ExecuteInDirFunc = func(dir, name string, args ...string) ([]byte, error) {
		if name == "go" && len(args) >= 2 && args[0] == "run" && args[1] == executionFile {
			return expectedOutput, nil
		}
		return nil, fmt.Errorf("unexpected command")
	}

	// Act
	result, err := service.ExecuteGoRun(executionFile, rootDirectory, parameters)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "go runを実行中")
	assert.NotContains(t, result, "パラメータ:")
	assert.Contains(t, result, string(expectedOutput))
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

	mockDirectoryChecker.ExistsFunc = func(path string) bool {
		return path == coverageFile
	}
	mockCommandExecutor.ExecuteInDirFunc = func(dir, name string, args ...string) ([]byte, error) {
		if name == "go" && len(args) >= 3 && args[0] == "tool" && args[1] == "cover" && args[2] == "-func=coverage.out" {
			return expectedOutput, nil
		}
		return nil, fmt.Errorf("unexpected command")
	}

	// Act
	result, err := service.ExecuteCoverageFunc(coverageFile, grepPattern)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "カバレッジファイルから関数情報を取得中")
	assert.NotContains(t, result, "grepパターン:")
	assert.Contains(t, result, string(expectedOutput))
}
