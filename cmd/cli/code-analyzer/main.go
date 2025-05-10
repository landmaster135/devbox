// cmd/codemetrics/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/landmaster135/devbox/internal/independencies/code_analyzer/app"
	"github.com/landmaster135/devbox/internal/independencies/code_analyzer/config"
)

func main() {
	// 設定のインスタンス作成
	appConfig := config.NewAppConfig()

	// フラグを定義
	flag.StringVar(&appConfig.ProjectPath, "path", ".", "Path to the project to analyze")
	flag.StringVar(&appConfig.OutputFormat, "format", "text", "Output format (text, json, csv)")
	flag.StringVar(&appConfig.OutputFile, "output", "", "Output file (if not provided, writes to stdout)")

	extensions := flag.String("ext", ".go", "Comma-separated list of file extensions to analyze")

	flag.BoolVar(&appConfig.VisualReport, "visual", false, "Generate visual HTML report")
	flag.StringVar(&appConfig.HistoryPath, "history", "", "Path to historical data file for trend analysis")
	flag.BoolVar(&appConfig.DetectClones, "detect-clones", false, "Enable code clone detection")
	flag.IntVar(&appConfig.MinBlockSize, "min-block-size", 30, "Minimum token block size for clone detection")
	flag.Float64Var(&appConfig.MinSimilarity, "min-similarity", 0.8, "Minimum similarity threshold for clone detection (0.0-1.0)")
	flag.BoolVar(&appConfig.Verbose, "v", false, "Enable verbose logging")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nCode Metrics Analyzer Tool\n")
		fmt.Fprintf(os.Stderr, "Analyzes source code and generates metrics reports\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s -path /path/to/project -ext .go,.py -format json -output metrics.json -visual\n", os.Args[0])
	}

	// フラグを解析
	flag.Parse()

	// ログ設定
	if appConfig.Verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
	}

	// 拡張子リストの処理
	appConfig.SetExtensions(*extensions)

	// アプリケーションを実行
	application := app.NewApp(appConfig)
	exitCode := application.Run(os.Stdout, os.Stderr)

	// 終了コードを設定
	os.Exit(exitCode)
}
