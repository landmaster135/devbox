package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/landmaster135/devbox/internal/claude_code_usage/app"
	"github.com/landmaster135/devbox/internal/claude_code_usage/config"
)

func main() {
	// フラグを定義
	sinceStr := flag.String("since", "", "開始日フィルター (YYYYMMDD形式)")
	untilStr := flag.String("until", "", "終了日フィルター (YYYYMMDD形式)")
	claudePath := flag.String("path", "", "Claudeデータディレクトリのカスタムパス (デフォルト: ~/.claude)")
	outputJSON := flag.Bool("json", false, "JSON形式で結果を出力")
	showHelp := flag.Bool("help", false, "ヘルプメッセージを表示")
	showVersion := flag.Bool("version", false, "バージョン情報を表示")
	verbose := flag.Bool("v", false, "詳細なログを出力")

	// カスタムUsage関数を設定
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "claude-code-usage - Claude Code Usage Analysis Tool\n\n")
		fmt.Fprintf(os.Stderr, "USAGE:\n")
		fmt.Fprintf(os.Stderr, "    claude-code-usage [COMMAND] [OPTIONS]\n\n")
		fmt.Fprintf(os.Stderr, "COMMANDS:\n")
		fmt.Fprintf(os.Stderr, "    daily    Show daily usage report (default)\n")
		fmt.Fprintf(os.Stderr, "    session  Show usage grouped by conversation sessions\n\n")
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "    claude-code-usage daily\n")
		fmt.Fprintf(os.Stderr, "    claude-code-usage session -json\n")
		fmt.Fprintf(os.Stderr, "    claude-code-usage daily -since 20250525 -until 20250530\n")
	}

	// フラグを解析
	flag.Parse()

	// ログ設定
	if *verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
	}

	// コマンドを取得（フラグ解析後の残りの引数）
	command := "daily" // デフォルト
	args := flag.Args()
	if len(args) > 0 {
		firstArg := strings.ToLower(args[0])
		if firstArg == "daily" || firstArg == "session" {
			command = firstArg
		}
	}

	// アプリケーション設定
	appConfig := &config.AppConfig{
		Command:     command,
		ClaudePath:  *claudePath,
		OutputJSON:  *outputJSON,
		ShowHelp:    *showHelp,
		ShowVersion: *showVersion,
	}

	// 日付フィルターの処理
	if *sinceStr != "" {
		since, err := app.ParseDate(*sinceStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid since date format: %v\n", err)
			fmt.Fprintf(os.Stderr, "Please use YYYYMMDD format (e.g., 20250525)\n")
			os.Exit(1)
		}
		appConfig.Since = &since
	}

	if *untilStr != "" {
		until, err := app.ParseDate(*untilStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid until date format: %v\n", err)
			fmt.Fprintf(os.Stderr, "Please use YYYYMMDD format (e.g., 20250530)\n")
			os.Exit(1)
		}
		appConfig.Until = &until
	}

	// アプリケーションを実行
	application := app.NewApp(appConfig)
	exitCode := application.Run(os.Stdout, os.Stderr)

	// 終了コードを設定
	os.Exit(exitCode)
}
