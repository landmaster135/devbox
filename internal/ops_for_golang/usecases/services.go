package usecases

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CommandExecutor はコマンド実行のインターフェースです（テスト用）
type CommandExecutor interface {
	Execute(name string, args ...string) ([]byte, error)
	ExecuteInDir(dir, name string, args ...string) ([]byte, error)
}

// DefaultCommandExecutor は標準のexec.Commandを使用する実装です
type DefaultCommandExecutor struct{}

// Execute はコマンドを実行します
func (e *DefaultCommandExecutor) Execute(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

// ExecuteInDir は指定されたディレクトリでコマンドを実行します
func (e *DefaultCommandExecutor) ExecuteInDir(dir, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// DirectoryChecker はディレクトリ存在確認のインターフェースです（テスト用）
type DirectoryChecker interface {
	Exists(path string) bool
	IsDirectory(path string) bool
}

// DefaultDirectoryChecker は標準のos.Statを使用する実装です
type DefaultDirectoryChecker struct{}

// Exists はパスが存在するかチェックします
func (c *DefaultDirectoryChecker) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDirectory はパスがディレクトリかチェックします
func (c *DefaultDirectoryChecker) IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// #==============================================================#
// ##          GolangOpsService                                  ##
// #==============================================================#
// GolangOpsService はGo開発用操作を行うサービスです
type GolangOpsService struct {
	commandExecutor  CommandExecutor
	directoryChecker DirectoryChecker
}

// NewGolangOpsService は新しいGolangOpsServiceを作成します
func NewGolangOpsService() *GolangOpsService {
	return &GolangOpsService{
		commandExecutor:  &DefaultCommandExecutor{},
		directoryChecker: &DefaultDirectoryChecker{},
	}
}

// NewGolangOpsServiceWithDependencies はテスト用に依存性を注入できるGolangOpsServiceを作成します
func NewGolangOpsServiceWithDependencies(commandExecutor CommandExecutor, directoryChecker DirectoryChecker) *GolangOpsService {
	return &GolangOpsService{
		commandExecutor:  commandExecutor,
		directoryChecker: directoryChecker,
	}
}

// ExecuteTestCoverage はテストカバレッジを実行します
func (s *GolangOpsService) ExecuteTestCoverage(directory string) error {
	// ディレクトリの存在確認
	if !s.directoryChecker.Exists(directory) {
		return fmt.Errorf("指定されたディレクトリが存在しません: %s", directory)
	}
	if !s.directoryChecker.IsDirectory(directory) {
		return fmt.Errorf("指定されたパスはディレクトリではありません: %s", directory)
	}

	// 絶対パスに変換
	absDir, err := filepath.Abs(directory)
	if err != nil {
		return fmt.Errorf("ディレクトリパスの変換に失敗しました: %v", err)
	}

	fmt.Printf("テストカバレッジを実行中: %s\n", absDir)
	fmt.Printf("実行コマンド: go test -cover ./...\n\n")

	// go test -cover ./... を実行
	output, err := s.commandExecutor.ExecuteInDir(absDir, "go", "test", "-cover", "./...")
	if err != nil {
		return fmt.Errorf("テストカバレッジの実行に失敗しました: %v\n出力: %s", err, string(output))
	}

	fmt.Print(string(output))
	fmt.Println("\nテストカバレッジの実行が完了しました。")
	return nil
}

// ExecuteTestCoverageProject はプロジェクト全体のテストカバレッジを実行します
func (s *GolangOpsService) ExecuteTestCoverageProject(directory string) error {
	// ディレクトリの存在確認
	if !s.directoryChecker.Exists(directory) {
		return fmt.Errorf("指定されたディレクトリが存在しません: %s", directory)
	}
	if !s.directoryChecker.IsDirectory(directory) {
		return fmt.Errorf("指定されたパスはディレクトリではありません: %s", directory)
	}

	// 絶対パスに変換
	absDir, err := filepath.Abs(directory)
	if err != nil {
		return fmt.Errorf("ディレクトリパスの変換に失敗しました: %v", err)
	}

	fmt.Printf("プロジェクト全体のテストカバレッジを実行中: %s\n", absDir)
	fmt.Printf("実行コマンド: go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out && go tool cover -html=coverage.out -o coverage.html\n\n")

	// Step 1: go test -coverprofile=coverage.out ./...
	fmt.Println("Step 1: テストカバレッジプロファイルを生成中...")
	output1, err := s.commandExecutor.ExecuteInDir(absDir, "go", "test", "-coverprofile=coverage.out", "./...")
	if err != nil {
		return fmt.Errorf("テストカバレッジプロファイルの生成に失敗しました: %v\n出力: %s", err, string(output1))
	}
	fmt.Print(string(output1))

	// Step 2: go tool cover -func=coverage.out
	fmt.Println("\nStep 2: カバレッジ関数レポートを生成中...")
	output2, err := s.commandExecutor.ExecuteInDir(absDir, "go", "tool", "cover", "-func=coverage.out")
	if err != nil {
		return fmt.Errorf("カバレッジ関数レポートの生成に失敗しました: %v\n出力: %s", err, string(output2))
	}
	fmt.Print(string(output2))

	// Step 3: go tool cover -html=coverage.out -o coverage.html
	fmt.Println("\nStep 3: HTMLカバレッジレポートを生成中...")
	output3, err := s.commandExecutor.ExecuteInDir(absDir, "go", "tool", "cover", "-html=coverage.out", "-o", "coverage.html")
	if err != nil {
		return fmt.Errorf("HTMLカバレッジレポートの生成に失敗しました: %v\n出力: %s", err, string(output3))
	}
	if len(output3) > 0 {
		fmt.Print(string(output3))
	}

	fmt.Printf("\nプロジェクト全体のテストカバレッジの実行が完了しました。\n")
	fmt.Printf("HTMLレポートが生成されました: %s/coverage.html\n", absDir)
	return nil
}

// ExecuteGoRun はgo runを実行します
func (s *GolangOpsService) ExecuteGoRun(executionFile, parameters string) error {
	// 実行ファイルの存在確認
	if !s.directoryChecker.Exists(executionFile) {
		return fmt.Errorf("指定された実行ファイルが存在しません: %s", executionFile)
	}

	// コマンド引数を構築
	args := []string{"run", executionFile}

	// パラメータが指定されている場合は追加
	if parameters != "" {
		// パラメータを空白で分割して追加
		paramArgs := strings.Fields(parameters)
		args = append(args, paramArgs...)
	}

	fmt.Printf("go runを実行中: %s\n", executionFile)
	if parameters != "" {
		fmt.Printf("パラメータ: %s\n", parameters)
	}
	fmt.Printf("実行コマンド: go %s\n\n", strings.Join(args, " "))

	// go run を実行
	output, err := s.commandExecutor.Execute("go", args...)
	if err != nil {
		return fmt.Errorf("go runの実行に失敗しました: %v\n出力: %s", err, string(output))
	}

	fmt.Print(string(output))
	fmt.Println("\ngo runの実行が完了しました。")
	return nil
}

// HandleTestCoverage はテストカバレッジのハンドラーです
func (s *GolangOpsService) HandleTestCoverage(directory string) error {
	return s.ExecuteTestCoverage(directory)
}

// HandleTestCoverageProject はプロジェクト全体のテストカバレッジのハンドラーです
func (s *GolangOpsService) HandleTestCoverageProject(directory string) error {
	return s.ExecuteTestCoverageProject(directory)
}

// ExecuteCoverageFunc はカバレッジファイルから関数情報を表示します
func (s *GolangOpsService) ExecuteCoverageFunc(coverageFile string) error {
	// カバレッジファイルの存在確認
	if !s.directoryChecker.Exists(coverageFile) {
		return fmt.Errorf("指定されたカバレッジファイルが存在しません: %s", coverageFile)
	}

	fmt.Printf("カバレッジファイルから関数情報を取得中: %s\n", coverageFile)
	fmt.Printf("実行コマンド: go tool cover -func=%s\n\n", coverageFile)

	// go tool cover -func=[coverage_file] を実行
	output, err := s.commandExecutor.Execute("go", "tool", "cover", fmt.Sprintf("-func=%s", coverageFile))
	if err != nil {
		return fmt.Errorf("カバレッジ関数情報の取得に失敗しました: %v\n出力: %s", err, string(output))
	}

	fmt.Print(string(output))
	fmt.Println("\nカバレッジ関数情報の取得が完了しました。")
	return nil
}

// HandleGoRun はgo runのハンドラーです
func (s *GolangOpsService) HandleGoRun(executionFile, parameters string) error {
	return s.ExecuteGoRun(executionFile, parameters)
}

// HandleCoverageFunc はカバレッジ関数情報のハンドラーです
func (s *GolangOpsService) HandleCoverageFunc(coverageFile string) error {
	return s.ExecuteCoverageFunc(coverageFile)
}
