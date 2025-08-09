package config

import (
	"fmt"
	"os"
)

// Config はops-for-golang CLIの設定を保持する構造体
type Config struct {
	Operation     string // 操作タイプ (test-coverage, test-coverage-project, run, coverage-func)
	Directory     string // ディレクトリパス (test-coverage, test-coverage-project用)
	ExecutionFile string // 実行ファイル (run用)
	Parameters    string // 実行パラメータ (run用)
	CoverageFile  string // カバレッジファイル (coverage-func用)
	GrepPattern   string // 出力フィルタリング用のgrepパターン (全操作共通)
	Help          bool   // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(operation, directory, executionFile, parameters, coverageFile, grepPattern string) (*Config, error) {
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	// 操作タイプの検証
	validOperations := []string{"test-coverage", "test-coverage-project", "run", "coverage-func"}
	isValid := false
	for _, op := range validOperations {
		if operation == op {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, fmt.Errorf("無効な操作タイプです: %s", operation)
	}

	// 操作タイプ別の検証
	switch operation {
	case "test-coverage", "test-coverage-project":
		if directory == "" {
			return nil, fmt.Errorf("%s操作にはディレクトリパスが必要です", operation)
		}
	case "run":
		if executionFile == "" {
			return nil, fmt.Errorf("run操作には実行ファイルが必要です")
		}
	case "coverage-func":
		if coverageFile == "" {
			return nil, fmt.Errorf("coverage-func操作にはカバレッジファイルが必要です")
		}
	}

	return &Config{
		Operation:     operation,
		Directory:     directory,
		ExecutionFile: executionFile,
		Parameters:    parameters,
		CoverageFile:  coverageFile,
		GrepPattern:   grepPattern,
	}, nil
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation     = ""
		directory     = ""
		executionFile = ""
		parameters    = ""
		coverageFile  = ""
		grepPattern   = ""
		help          = false
	)

	parser.StringVar(&operation, "ops", operation, "操作タイプ (test-coverage, test-coverage-project, run, coverage-func)")
	parser.StringVar(&operation, "o", operation, "操作タイプの短縮形")

	// test-coverage, test-coverage-project用のパラメータ
	parser.StringVar(&directory, "directory", directory, "対象ディレクトリパス")
	parser.StringVar(&directory, "d", directory, "ディレクトリパスの短縮形")

	// run用のパラメータ
	parser.StringVar(&executionFile, "execution_file", executionFile, "実行ファイルパス (run操作用)")
	parser.StringVar(&executionFile, "e", executionFile, "実行ファイルパスの短縮形")
	parser.StringVar(&parameters, "parameters", parameters, "実行パラメータ (run操作用)")
	parser.StringVar(&parameters, "p", parameters, "実行パラメータの短縮形")

	// coverage-func用のパラメータ
	parser.StringVar(&coverageFile, "coverage_file", coverageFile, "カバレッジファイルパス (coverage-func操作用)")
	parser.StringVar(&coverageFile, "c", coverageFile, "カバレッジファイルパスの短縮形")

	// grep用のパラメータ（全操作共通）
	parser.StringVar(&grepPattern, "grep_pattern", grepPattern, "出力フィルタリング用のgrepパターン (全操作共通)")
	parser.StringVar(&grepPattern, "g", grepPattern, "grepパターンの短縮形")

	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	return NewConfig(operation, directory, executionFile, parameters, coverageFile, grepPattern)
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `Go開発用操作CLIツール

使用方法:
  テストカバレッジ取得:
    %s -ops test-coverage -directory /path/to/project
    %s -o test-coverage -d /path/to/project

  プロジェクト全体のテストカバレッジ取得:
    %s -ops test-coverage-project -directory /path/to/project
    %s -o test-coverage-project -d /path/to/project

  カバレッジファイルから関数情報取得:
    %s -ops coverage-func -coverage_file coverage.out
    %s -o coverage-func -c coverage.out

  go run実行:
    %s -ops run -execution_file ./main.go -parameters "-dry-run -token 'test_token'"
    %s -o run -e ./main.go -p "-dry-run -token 'test_token'"

  grepパターンでフィルタリング:
    %s -ops coverage-func -coverage_file coverage.out -grep_pattern "100.0%%"
    %s -o test-coverage -d /path/to/project -g "PASS"

オプション:
  -ops, -o           操作タイプ (test-coverage, test-coverage-project, run, coverage-func)
  -directory, -d     対象ディレクトリパス (test-coverage, test-coverage-project用)
  -execution_file, -e 実行ファイルパス (run用)
  -parameters, -p    実行パラメータ (run用)
  -coverage_file, -c カバレッジファイルパス (coverage-func用)
  -grep_pattern, -g  出力フィルタリング用のgrepパターン (全操作共通)
  -help, -h          このヘルプを表示

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
