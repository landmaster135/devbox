package main

import (
	"flag"
	"log"
	"os"

	"github.com/landmaster135/devbox/internal/code_analyzer/config"
	"github.com/landmaster135/devbox/internal/code_analyzer/app"
)

func main() {
	// フラグを定義
	projectPath := flag.String("path", ".", "解析対象のプロジェクトパス")
	outputFormat := flag.String("format", "text", "出力形式 (text, json, csv)")
	outputFile := flag.String("output", "", "出力ファイル（省略時は標準出力）")
	extensions := flag.String("ext", ".go", "解析対象のファイル拡張子（カンマ区切り）")
	visualReport := flag.Bool("visual", false, "視覚的なHTMLレポートを生成")
	historyPath := flag.String("history", "", "トレンド分析用の履歴データファイルのパス")
	detectClones := flag.Bool("detect-clones", false, "コードクローン検出を有効化")
	minBlockSize := flag.Int("min-block-size", 30, "クローン検出のための最小トークンブロックサイズ")
	minSimilarity := flag.Float64("min-similarity", 0.8, "クローン検出のための最小類似度閾値 (0.0-1.0)")
	verbose := flag.Bool("v", false, "詳細なログを出力")

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
		ProjectPath:    *projectPath,
		OutputFormat:   *outputFormat,
		OutputFile:     *outputFile,
		VisualReport:   *visualReport,
		HistoryPath:    *historyPath,
		DetectClones:   *detectClones,
		MinBlockSize:   *minBlockSize,
		MinSimilarity:  *minSimilarity,
		Verbose:        *verbose,
	}

	// 拡張子リストの処理
	appConfig.SetExtensions(*extensions)

	// アプリケーションを実行
	application := app.NewApp(appConfig)
	exitCode := application.Run(os.Stdout, os.Stderr)

	// 終了コードを設定
	os.Exit(exitCode)
}
