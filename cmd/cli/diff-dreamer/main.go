package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"

	usecases "github.com/landmaster135/devbox/internal/diff_dreamer/usecases"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

//go:embed web/index.html
var indexHTML string

//go:embed web/style.css
var styleCSS string

//go:embed web/script.js
var scriptJS string

// run はdiff-dreamerの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	flagSet := flag.NewFlagSet("diff-dreamer", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	// コマンドライン引数の定義
	leftFile := flagSet.String("left", "", "左側に表示するテキストファイルのパス")
	rightFile := flagSet.String("right", "", "右側に表示するテキストファイルのパス")
	outputFile := flagSet.String("output", "diff_dreamer.html", "出力HTMLファイル名")

	// 引数の解析
	if err := flagSet.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// 統合メソッドを使用してdiff-dreamerの全処理を実行
	err := usecases.ProcessDiffDreamer(*leftFile, *rightFile, *outputFile, indexHTML, styleCSS, scriptJS)
	if err != nil {
		fmt.Fprintf(stderr, "Diff Dreamerの処理に失敗しました: %v\n", err)
		return exitCodeError
	}

	fmt.Fprintf(stdout, "Diff Dreamerが起動しました\n")
	if *leftFile != "" || *rightFile != "" {
		fmt.Fprintf(stdout, "指定されたファイルが読み込まれています\n")
	}
	fmt.Fprintf(stdout, "ブラウザでテキストを比較してください\n")
	fmt.Fprintf(stdout, "\n使用方法:\n")
	fmt.Fprintf(stdout, "  - 左右のテキストエリアにテキストを入力\n")
	fmt.Fprintf(stdout, "  - 「比較する」ボタンをクリック\n")
	fmt.Fprintf(stdout, "  - Ctrl+Enter: 比較実行\n")
	fmt.Fprintf(stdout, "  - Ctrl+Shift+C: すべてクリア\n")

	// ユーザーが確認できるように一時停止
	fmt.Fprintf(stdout, "\nEnterキーを押すとプログラムを終了します（HTMLファイルは削除されません）...\n")
	fmt.Scanln()

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
