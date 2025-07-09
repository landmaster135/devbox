package main

import (
	"flag"
	"log"
	"os"

	"github.com/landmaster135/devbox/internal/depends_visualizer/config"
	"github.com/landmaster135/devbox/internal/depends_visualizer/usecases"
)

func main() {
	// フラグを定義
	configPath := flag.String("config", "", "設定ファイルのパス")
	srcFile := flag.String("file", "", "解析対象のソースファイル")
	extension := flag.String("ext", "", "ファイルの拡張子 (.go, .py, .js)")
	outputPath := flag.String("out", "", "出力ファイルのパス（省略時は標準出力）")
	format := flag.String("format", "mermaid", "出力形式 (mermaid, mermaid-flowchart, plantuml, dot)")
	recursive := flag.Bool("r", false, "ディレクトリ内のファイルを再帰的に解析")
	verbose := flag.Bool("v", false, "詳細なログを出力")
	dir := flag.String("dir", ".", "解析対象のディレクトリ")

	// フラグを解析
	flag.Parse()

	// ログ設定
	if *verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
	}

	// アプリケーション設定
	appConfig := &config.AppConfig{
		ConfigPath: *configPath,
		SourceFile: *srcFile,
		Extension:  *extension,
		OutputPath: *outputPath,
		Format:     *format,
		Recursive:  *recursive,
		Verbose:    *verbose,
		Directory:  *dir,
	}

	// アプリケーションを実行
	app := usecases.NewApp(appConfig)
	exitCode := app.Run(os.Stdout, os.Stderr)

	// 終了コードを設定
	os.Exit(exitCode)
}
