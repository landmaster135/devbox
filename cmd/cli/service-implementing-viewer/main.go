package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/service_implementing_viewer/config"
	usecases "github.com/landmaster135/devbox/internal/service_implementing_viewer/usecases"
)

func main() {
	// コマンドライン引数を解析
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// サービスを作成
	service := usecases.NewServiceImplementingViewerService(cfg.RootDir, cfg.TargetDirs)

	// サービス実装状況を取得
	result, err := service.GetServiceImplementingStatus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Print(result)
}
