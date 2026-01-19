package main

import (
	"flag"
	"os"

	"github.com/landmaster135/devbox/internal/script_generator_to_build/config"
	"github.com/landmaster135/devbox/internal/script_generator_to_build/usecases"
)

func main() {
	// コマンドライン引数の解析
	var (
		showHelp        bool
		baseDir         string
		cliDir          string
		scriptsDir      string
		outputDir       string
		packageName string
	)
	flag.BoolVar(&showHelp, "h", false, "ヘルプメッセージを表示")
	flag.BoolVar(&showHelp, "help", false, "ヘルプメッセージを表示")
	flag.StringVar(&baseDir, "base_dir", "", "ベースディレクトリを指定")
	flag.StringVar(&cliDir, "cli_dir", "", "CLIディレクトリを指定")
	flag.StringVar(&scriptsDir, "scripts_dir", "", "スクリプトディレクトリを指定")
	flag.StringVar(&outputDir, "output_dir", "", "出力ディレクトリを指定")
	flag.StringVar(&packageName, "package_name", "", "生成対象のパッケージ名を指定")
	flag.Parse()

	// サービス設定
	appConfig := &config.ServiceConfig{
		PackageName: packageName,
		ShowHelp:    showHelp,
		BaseDir:     baseDir,
		CLIDir:      cliDir,
		ScriptsDir:  scriptsDir,
		OutputDir:   outputDir,
	}

	// サービスを実行
	app := usecases.NewService(appConfig)
	exitCode := app.Run(os.Stdout, os.Stderr)

	// 終了コードを設定
	os.Exit(exitCode)
}
