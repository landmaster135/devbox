package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/landmaster135/devbox/internal/interfaces/repositories"
	"github.com/landmaster135/devbox/internal/usecases/services"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run は環境変数ローダーの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	fs := flag.NewFlagSet("env-loader", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// コマンドライン引数の定義
	envFile := fs.String("env", "", "環境変数を読み込むYAMLファイルのパス（デフォルト: env.yml）")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// 依存関係の注入
	envRepo := repositories.NewEnvRepository()
	envService := services.NewEnvService(envRepo)

	// 環境変数ファイルのパスを解決
	envFilePath := envService.ResolveEnvFilePath(*envFile)

	// 環境変数を読み込み、設定
	if err := envService.LoadAndSetEnvFromYaml(envFilePath); err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	fmt.Fprintf(stdout, "環境変数を正常に読み込みました: %s\n", envFilePath)
	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
