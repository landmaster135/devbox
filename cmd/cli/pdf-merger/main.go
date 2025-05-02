package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/landmaster135/devbox/internal/independencies/pdf_merger/usecases"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run はPDFマージャーツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	fs := flag.NewFlagSet("pdf-merger", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// コマンドライン引数の定義
	dir := fs.String("dir", ".", "画像を検索するフォルダー (再帰探索)")
	out := fs.String("out", "", "出力 PDF ファイル名 (未指定なら <dir 名>.pdf)")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// 画像ファイルの取得
	images, output, err := usecases.GetSourceImages(*dir, *out)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	if len(images) == 0 {
		fmt.Fprintln(stderr, "画像が見つかりませんでした。終了します。")
		return exitCodeError
	}

	fmt.Fprintf(stdout, "検出した画像: %d 枚\n", len(images))
	fmt.Fprintf(stdout, "出力 PDF   : %s\n", output)

	// PDFの生成
	err = usecases.MergeImagesIntoPDF(images, output)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	fmt.Fprintln(stdout, "PDF を生成しました。完了です。")
	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
