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

// run はJSON ISO8601 Converterツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	fs := flag.NewFlagSet("json-iso8601-converter", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// コマンドライン引数の定義
	dirPath := fs.String("dir", ".", "JSONファイルを検索するディレクトリパス")
	key := fs.String("key", "", "変換対象のキー名")
	recursive := fs.Bool("recursive", false, "サブディレクトリも検索するかどうか")
	dryRun := fs.Bool("dry-run", false, "変換をシミュレーションするだけで実際には変更しない")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// キーが指定されていない場合はエラー
	if *key == "" {
		fmt.Fprintln(stderr, "エラー: 変換対象のキー名を指定してください（-key オプション）")
		fs.Usage()
		return exitCodeError
	}

	// 依存関係の注入
	fileRepo := repositories.NewFileRepository()
	iso8601Repo := repositories.NewISO8601Repository()
	jsonRepo := repositories.NewJSONRepository(fileRepo, iso8601Repo)
	iso8601Service := services.NewISO8601Service(jsonRepo)

	// ディレクトリ内のJSONファイルを検索して変換
	count, err := iso8601Service.ConvertISO8601ToTimestamp(*dirPath, *key, *recursive, *dryRun)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	if *dryRun {
		fmt.Fprintf(stdout, "ドライラン: %d 個のJSONファイルで変換対象のキー '%s' が見つかりました\n", count, *key)
	} else {
		fmt.Fprintf(stdout, "%d 個のJSONファイルのキー '%s' の値をISO8601形式からUNIXタイムスタンプに変換しました\n", count, *key)
	}

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
