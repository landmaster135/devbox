package main

import (
	_ "embed"
	"flag"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	usecases "github.com/landmaster135/devbox/internal/independencies/image_trim_describer/usecases"
)

//go:embed web/index.html
var indexHTML string

//go:embed web/style.css
var styleCSS string

//go:embed web/script.js
var scriptJS string

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run は画像トリミング記述ツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	flagSet := flag.NewFlagSet("image-trim-describer", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	// コマンドライン引数の定義
	srcDir := flagSet.String("src", ".", "画像ファイルのディレクトリ")
	outDir := flagSet.String("out", "", "出力先ディレクトリ (空の場合はsrcと同じ)")
	arcDir := flagSet.String("arc", "./5_original_files", "アーカイブ先ディレクトリ")
	suffix := flagSet.String("suffix", "trimmed", "出力ファイル名に付加するサフィックス")
	move := flagSet.Bool("move", false, "元ファイルを移動する (コピーではなく)")

	// 引数の解析
	if err := flagSet.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// 出力ディレクトリが指定されていない場合は入力ディレクトリと同じに
	if *outDir == "" {
		*outDir = *srcDir
	}

	// 画像ファイル一覧を取得
	imageFiles, err := usecases.GetImageFiles(*srcDir)
	if err != nil {
		fmt.Fprintf(stderr, "画像ファイル一覧の取得に失敗しました: %v\n", err)
		return exitCodeError
	}

	// 設定情報を作成
	config := usecases.Config{
		ImageFiles: imageFiles,
		SrcDir:     *srcDir,
		OutDir:     *outDir,
		ArcDir:     *arcDir,
		Suffix:     *suffix,
		Move:       *move,
	}

	// 設定情報をJavaScript形式に変換
	configJS := usecases.GenerateConfigJS(config)

	// HTMLテンプレートのデータを作成
	data := usecases.TemplateData{
		Style:      styleCSS,
		Script:     scriptJS,
		ConfigJS:   configJS,
		ImageFiles: imageFiles,
	}

	// HTMLを生成
	html, err := usecases.GenerateHTML(indexHTML, data)
	if err != nil {
		fmt.Fprintf(stderr, "HTMLの生成に失敗しました: %v\n", err)
		return exitCodeError
	}

	// カレントディレクトリにHTMLファイルを作成
	htmlFilePath := "image_trimmer.html"
	err = usecases.OpenHTMLPage(html, htmlFilePath)
	if err != nil {
		return exitCodeError
	}

	fmt.Fprintf(stdout, "画像トリミングツールが起動しました\n")
	fmt.Fprintf(stdout, "ブラウザでトリミング座標を選択し、トリミングを実行してください\n")

	// ユーザーが確認できるように一時停止
	fmt.Fprintf(stdout, "Enterキーを押すとプログラムを終了します（HTMLファイルは削除されません）...\n")
	fmt.Scanln()

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
