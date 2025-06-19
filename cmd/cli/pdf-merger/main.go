package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	usecases "github.com/landmaster135/devbox/internal/independencies/pdf_merger/usecases"
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
	add := fs.String("add", "", "既存の PDF ファイルパス (指定時は既存PDFに画像を追加)")

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

	// 既存PDFファイルが指定されている場合は既存PDFに画像を追加
	if *add != "" {
		// 既存PDFファイルの存在確認
		if _, err := os.Stat(*add); os.IsNotExist(err) {
			fmt.Fprintf(stderr, "エラー: 既存PDFファイルが見つかりません: %s\n", *add)
			return exitCodeError
		}
		
		fmt.Fprintf(stdout, "既存 PDF   : %s\n", *add)
		
		// 既存PDFに画像を追加
		err = usecases.AddImagesToExistingPDF(*add, images, output)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", err)
			return exitCodeError
		}
		
		fmt.Fprintln(stdout, "既存PDFに画像を追加しました。完了です。")
	} else {
		// 新規PDFの生成
		err = usecases.MergeImagesIntoPDF(images, output)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", err)
			return exitCodeError
		}
		
		fmt.Fprintln(stdout, "PDF を生成しました。完了です。")
	}
	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
