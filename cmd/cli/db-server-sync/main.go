package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/db_server_sync/config"
	usecases "github.com/landmaster135/devbox/internal/db_server_sync/usecases"
)

func main() {
	// コマンドライン引数を解析
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// ヘルプが要求された場合
	if cfg.Help {
		config.PrintUsage()
		return
	}

	// データベースサーバー同期サービスを初期化
	service := usecases.NewDBServerSyncService()

	// 操作に応じて処理を実行
	switch cfg.Operation {
	case "append-anime":
		err = service.ProcessAppendAnime(cfg.InputFilePath, cfg.AdditionalInputFilePath, cfg.OutputFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("処理が完了しました。出力ファイル: %s\n", cfg.OutputFilePath)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作です: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}
