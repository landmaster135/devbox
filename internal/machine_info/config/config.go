package config

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	// OperationUbuntu はUbuntu環境でマシン情報を収集するoperation名
	OperationUbuntu = "ubuntu"

	defaultNetworkInterface = "eth0"
	defaultOutputDir        = "."
)

// Config はmachine-info CLIの設定を保持する
type Config struct {
	Operation        string
	NetworkInterface string
	OutputDir        string
	Help             bool
}

// ParseFlags はコマンドライン引数からConfigを生成する
func ParseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("machine-info", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cfg := &Config{}
	fs.StringVar(&cfg.Operation, "operation", OperationUbuntu, "実行する操作 (例: ubuntu)")
	fs.StringVar(&cfg.NetworkInterface, "network-interface", defaultNetworkInterface, "ネットワークインターフェース名")
	fs.StringVar(&cfg.OutputDir, "output-dir", defaultOutputDir, "ログファイルを保存するディレクトリ")

	var help bool
	fs.BoolVar(&help, "help", false, "このヘルプを表示")
	fs.BoolVar(&help, "h", false, "このヘルプを表示 (短縮)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg.Help = help
	cfg.Operation = strings.ToLower(strings.TrimSpace(cfg.Operation))

	if cfg.Operation == "" {
		return nil, fmt.Errorf("operationフラグを指定してください")
	}

	cfg.NetworkInterface = strings.TrimSpace(cfg.NetworkInterface)
	cfg.OutputDir = strings.TrimSpace(cfg.OutputDir)
	if cfg.OutputDir == "" {
		cfg.OutputDir = defaultOutputDir
	}

	return cfg, nil
}

// SupportedOperations は利用可能なoperationリストを返す
func SupportedOperations() []string {
	return []string{OperationUbuntu}
}

// PrintUsage はCLIの使い方を標準エラーに出力する
func PrintUsage() {
	fmt.Fprintln(os.Stderr, "machine-info CLI - PC情報の収集ツール")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "使用方法:")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/cli/machine-info/main.go --operation=ubuntu --network-interface=eth0")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "フラグ:")
	fmt.Fprintln(os.Stderr, "  --operation           実行する操作 (デフォルト: ubuntu)")
	fmt.Fprintln(os.Stderr, "  --network-interface   ネットワークインターフェース名 (デフォルト: eth0)")
	fmt.Fprintln(os.Stderr, "  --output-dir          ログファイルの出力先ディレクトリ (デフォルト: カレント)")
	fmt.Fprintln(os.Stderr, "  --help, -h            このヘルプを表示")
}
