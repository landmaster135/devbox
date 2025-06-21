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

	// PDF画像抽出用のオプション
	extract := fs.String("extract", "", "PDFファイルから画像を抽出する (PDFファイルパス)")
	outputDir := fs.String("output-dir", "", "画像の出力ディレクトリ (抽出時必須)")
	imageFormat := fs.String("format", "jpg", "出力画像形式 (jpg, jpeg, png, tiff, webp)")
	startPage := fs.Int("start", 0, "抽出開始ページ (1から開始、0は全ページ)")
	endPage := fs.Int("end", 0, "抽出終了ページ (0は最終ページまで)")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// PDFからの画像抽出処理
	if *extract != "" {
		return handleImageExtraction(*extract, *outputDir, *imageFormat, *startPage, *endPage, stdout, stderr)
	}

	// 既存のPDF作成機能
	return handlePDFCreation(*dir, *out, *add, stdout, stderr)
}

// handleImageExtraction はPDFからの画像抽出を処理します
func handleImageExtraction(pdfPath, outputDir, imageFormat string, startPage, endPage int, stdout, stderr io.Writer) exitCode {
	// 出力ディレクトリが指定されていない場合はエラー
	if outputDir == "" {
		fmt.Fprintln(stderr, "エラー: 画像抽出時は --output-dir オプションが必須です")
		return exitCodeError
	}

	// PDFファイルの存在確認
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		fmt.Fprintf(stderr, "エラー: PDFファイルが見つかりません: %s\n", pdfPath)
		return exitCodeError
	}

	// 画像抽出サービスのインスタンスを作成
	imageService := usecases.NewImageExtractionService()

	// PDFのページ数を取得
	totalPages, err := imageService.GetPageCount(pdfPath)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: PDFのページ数取得に失敗しました: %v\n", err)
		return exitCodeError
	}

	// startPageとendPageのバリデーション
	if err := imageService.ValidatePageRange(startPage, endPage, totalPages); err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	// ページ範囲の表示
	var pageRangeMsg string
	if startPage > 0 && endPage > 0 {
		pageRangeMsg = fmt.Sprintf("ページ %d から %d まで", startPage, endPage)
	} else if startPage > 0 && endPage == 0 {
		pageRangeMsg = fmt.Sprintf("ページ %d から %d まで (最終ページ)", startPage, totalPages)
	} else if startPage == 0 && endPage > 0 {
		pageRangeMsg = fmt.Sprintf("ページ 1 から %d まで", endPage)
	} else {
		pageRangeMsg = fmt.Sprintf("全ページ (ページ 1 から %d まで)", totalPages)
	}

	fmt.Fprintf(stdout, "PDF画像抽出を開始します...\n")
	fmt.Fprintf(stdout, "入力PDF    : %s\n", pdfPath)
	fmt.Fprintf(stdout, "出力ディレクトリ: %s\n", outputDir)
	fmt.Fprintf(stdout, "画像形式   : %s\n", imageFormat)
	fmt.Fprintf(stdout, "ページ範囲 : %s\n", pageRangeMsg)

	// 画像抽出の実行
	err = imageService.ExtractToImages(pdfPath, outputDir, imageFormat, startPage, endPage)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	fmt.Fprintln(stdout, "画像抽出が完了しました。")
	return exitCodeOK
}

// handlePDFCreation は既存のPDF作成機能を処理します
func handlePDFCreation(dir, out, add string, stdout, stderr io.Writer) exitCode {

	// PDF作成サービスのインスタンスを作成
	service := usecases.NewPDFCreationService()

	// 画像ファイルの取得
	images, output, err := service.GetSourceImages(dir, out)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	if len(images) == 0 {
		fmt.Fprintln(stderr, "画像が見つかりませんでした。終了します。")
		return exitCodeError
	}

	fmt.Fprintf(stdout, "検出した画像: %d 枚\n", len(images))

	// 既存PDFファイルが指定されている場合は既存PDFに画像を追加
	if add != "" {
		// 既存PDFファイルの存在確認
		if _, err := os.Stat(add); os.IsNotExist(err) {
			fmt.Fprintf(stderr, "エラー: 既存PDFファイルが見つかりません: %s\n", add)
			return exitCodeError
		}

		fmt.Fprintf(stdout, "既存 PDF   : %s\n", add)

		// 既存PDFに画像を追加
		err = service.AddImagesToExistingPDF(add, images, output)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", err)
			return exitCodeError
		}

		fmt.Fprintln(stdout, "既存PDFに画像を追加しました。完了です。")
	} else {
		// 新規PDFの生成
		err = service.MergeImagesIntoPDF(images, output)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", err)
			return exitCodeError
		}

		fmt.Fprintln(stdout, "PDF を生成しました。完了です。")
	}

	fmt.Fprintf(stdout, "出力 PDF   : %s\n", output)

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
