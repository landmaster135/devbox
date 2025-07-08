package config

import (
	"flag"
	"fmt"
	"os"
)

// 見出し定数
const (
	HeaderGitDiffRecord      = "=== Git Diff Record ==="
	HeaderFileChangesSummary = "=== File Changes Summary ==="
	HeaderDetailedDiff       = "=== Detailed Diff ==="
	HeaderNewFiles           = "=== New Files ==="
	HeaderDeletedFiles       = "=== Deleted Files ==="
)

// Config はCLI設定を保持する構造体
type Config struct {
	OutputDir  string
	StagedOnly bool
	ReadMode   bool
	GenMode    bool
	SourceDir  string
	Repository string
	GitDir     string
}

// ParseFlags はコマンドライン引数を解析してConfigを返す
func ParseFlags() (*Config, error) {
	var config Config

	flag.StringVar(&config.OutputDir, "output-dir", "", "出力先ディレクトリ (記録モード時必須)")
	flag.BoolVar(&config.StagedOnly, "staged-only", false, "ステージング済み差分のみ記録 (デフォルト: false)")
	flag.BoolVar(&config.ReadMode, "read-mode", false, "読み取りモードを有効にする")
	flag.BoolVar(&config.GenMode, "gen-mode", false, "生成モードを有効にする")
	flag.StringVar(&config.SourceDir, "source-dir", "", "読み取り対象のディレクトリ (読み取りモード時必須)")
	flag.StringVar(&config.Repository, "repository", "", "対象リポジトリ名 (読み取りモード時必須)")
	flag.StringVar(&config.GitDir, "git-dir", "", "対象Gitディレクトリ (生成モード時必須)")
	flag.Parse()

	if config.GenMode {
		// 生成モードの場合
		if config.GitDir == "" {
			return nil, fmt.Errorf("生成モードでは --git-dir は必須パラメータです")
		}
	} else if config.ReadMode {
		// 読み取りモードの場合
		if config.SourceDir == "" {
			return nil, fmt.Errorf("読み取りモードでは --source-dir は必須パラメータです")
		}
		if config.Repository == "" {
			return nil, fmt.Errorf("読み取りモードでは --repository は必須パラメータです")
		}
	} else {
		// 記録モードの場合
		if config.OutputDir == "" {
			return nil, fmt.Errorf("記録モードでは --output-dir は必須パラメータです")
		}
	}

	return &config, nil
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nオプション:\n")
	flag.PrintDefaults()
}
