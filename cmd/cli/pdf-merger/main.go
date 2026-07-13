package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	usecases "github.com/landmaster135/devbox/internal/pdf_merger/usecases"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// parseFlags はコマンドライン引数を解析してオプション構造体を返します
func parseFlags(args []string, stderr io.Writer) (usecases.PDFMergerOptions, error) {
	// フラグセットを作成
	fs := flag.NewFlagSet("pdf-merger", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// コマンドライン引数の定義
	dir := fs.String("dir", ".", "画像を検索するフォルダー")
	add := fs.String("add", "", "既存の PDF ファイルパス (指定時は既存PDFに画像を追加)")
	recursive := fs.Bool("recursive", false, "画像をサブディレクトリまで再帰探索する")

	// PDF画像抽出用のオプション
	extract := fs.String("extract", "", "PDFファイルから画像を抽出する (PDFファイルパス)")
	outputDir := fs.String("output-dir", "", "出力ディレクトリ (PDF作成/抽出時必須)")
	imageFormat := fs.String("format", "jpg", "出力画像形式 (jpg, jpeg, png, tiff, webp)")
	startPage := fs.Int("start", 0, "抽出開始ページ (1から開始、0は全ページ)")
	endPage := fs.Int("end", 0, "抽出終了ページ (0は最終ページまで)")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		return usecases.PDFMergerOptions{}, err
	}
	if *outputDir == "" {
		return usecases.PDFMergerOptions{}, fmt.Errorf("-output-dir は必須です")
	}

	// オプション構造体を作成して返す
	return usecases.PDFMergerOptions{
		Dir:         *dir,
		Add:         *add,
		Recursive:   *recursive,
		Extract:     *extract,
		OutputDir:   *outputDir,
		ImageFormat: *imageFormat,
		StartPage:   *startPage,
		EndPage:     *endPage,
	}, nil
}

// run はPDFマージャーツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグ解析
	opts, err := parseFlags(args, stderr)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// サービス実行
	logger := log.New(stdout, "", 0)
	service := usecases.NewPDFMergerServiceWithLogger(logger)
	if err := service.Process(opts); err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
