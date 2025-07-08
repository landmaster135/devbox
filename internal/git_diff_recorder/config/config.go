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
}

// ParseFlags はコマンドライン引数を解析してConfigを返す
func ParseFlags() (*Config, error) {
	var config Config

	flag.StringVar(&config.OutputDir, "output-dir", "", "出力先ディレクトリ (必須)")
	flag.BoolVar(&config.StagedOnly, "staged-only", false, "ステージング済み差分のみ記録 (デフォルト: false)")
	flag.Parse()

	if config.OutputDir == "" {
		return nil, fmt.Errorf("--output-dir は必須パラメータです")
	}

	return &config, nil
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nオプション:\n")
	flag.PrintDefaults()
}
