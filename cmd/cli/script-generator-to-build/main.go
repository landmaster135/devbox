package main

import (
	"flag"
	"os"

	"github.com/landmaster135/devbox/internal/script_generator_to_build/config"
	"github.com/landmaster135/devbox/internal/script_generator_to_build/usecases"
)

func main() {
	// コマンドライン引数の解析
	var showHelp bool
	flag.BoolVar(&showHelp, "h", false, "ヘルプメッセージを表示")
	flag.BoolVar(&showHelp, "help", false, "ヘルプメッセージを表示")
	flag.Parse()

	// 残りの引数があれば、最初の引数をパッケージ名として扱う
	var packageName string
	args := flag.Args()
	if len(args) > 0 {
		packageName = args[0]
	}

	// アプリケーション設定
	appConfig := &config.AppConfig{
		PackageName: packageName,
		ShowHelp:    showHelp,
	}

	// アプリケーションを実行
	app := usecases.NewApp(appConfig)
	exitCode := app.Run(os.Stdout, os.Stderr)

	// 終了コードを設定
	os.Exit(exitCode)
}
