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
	operation := fs.String("operation", "", "実行する処理 (merge-into-new, add-into-exist, extract-images)")
	srcDir := fs.String("src-dir", ".", "画像を検索するフォルダー")
	receivingFile := fs.String("receiving-file", "", "画像を追加する既存の PDF ファイルパス")
	recursive := fs.Bool("recursive", false, "画像をサブディレクトリまで再帰探索する")

	// PDF画像抽出用のオプション
	srcFile := fs.String("src-file", "", "PDFファイルから画像を抽出する入力PDFファイルパス")
	outputDir := fs.String("output-dir", "", "出力ディレクトリ (PDF作成/抽出時必須)")
	imageFormat := fs.String("format", "jpg", "出力画像形式 (jpg, jpeg, png, tiff, webp)")
	startPage := fs.Int("start", 0, "抽出開始ページ (1から開始、0は全ページ)")
	endPage := fs.Int("end", 0, "抽出終了ページ (0は最終ページまで)")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		return usecases.PDFMergerOptions{}, err
	}
	if *operation == "" {
		return usecases.PDFMergerOptions{}, fmt.Errorf("-operation は必須です")
	}
	if *outputDir == "" {
		return usecases.PDFMergerOptions{}, fmt.Errorf("-output-dir は必須です")
	}
	if err := validateOperationFlags(*operation, *receivingFile, *srcFile); err != nil {
		return usecases.PDFMergerOptions{}, err
	}

	// オプション構造体を作成して返す
	return usecases.PDFMergerOptions{
		Operation:     *operation,
		SrcDir:        *srcDir,
		ReceivingFile: *receivingFile,
		Recursive:     *recursive,
		SrcFile:       *srcFile,
		OutputDir:     *outputDir,
		ImageFormat:   *imageFormat,
		StartPage:     *startPage,
		EndPage:       *endPage,
	}, nil
}

func validateOperationFlags(operation, receivingFile, srcFile string) error {
	switch operation {
	case usecases.OperationMergeIntoNew:
		if receivingFile != "" {
			return fmt.Errorf("-operation=%s では -receiving-file を指定できません", operation)
		}
		if srcFile != "" {
			return fmt.Errorf("-operation=%s では -src-file を指定できません", operation)
		}
	case usecases.OperationAddIntoExist:
		if receivingFile == "" {
			return fmt.Errorf("-operation=%s では -receiving-file は必須です", operation)
		}
		if srcFile != "" {
			return fmt.Errorf("-operation=%s では -src-file を指定できません", operation)
		}
	case usecases.OperationExtractImages:
		if srcFile == "" {
			return fmt.Errorf("-operation=%s では -src-file は必須です", operation)
		}
		if receivingFile != "" {
			return fmt.Errorf("-operation=%s では -receiving-file を指定できません", operation)
		}
	default:
		return fmt.Errorf("未対応の operation です: %s", operation)
	}
	return nil
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
