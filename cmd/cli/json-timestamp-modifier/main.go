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

// run はJSON Timestamp Modifierツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	fs := flag.NewFlagSet("json-timestamp-modifier", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// コマンドライン引数の定義
	filePath := fs.String("file", "", "操作するJSONファイルのパス")
	key := fs.String("key", "timestamp", "タイムスタンプを追加するキー")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// ファイルパスが指定されていない場合はエラー
	if *filePath == "" {
		fmt.Fprintln(stderr, "エラー: JSONファイルのパスを指定してください（-file オプション）")
		fs.Usage()
		return exitCodeError
	}

	// 依存関係の注入
	fileRepo := repositories.NewFileRepository()
	jsonService := services.NewJSONService(fileRepo)
	timestampService := services.NewTimestampService(jsonService)

	// タイムスタンプを追加
	err := timestampService.AddTimestamp(*filePath, *key)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	fmt.Fprintf(stdout, "JSONファイル '%s' にキー '%s' と現在の日時のタイムスタンプを追加しました\n", *filePath, *key)

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
