package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/anilist/config"
	usecases "github.com/landmaster135/devbox/internal/anilist/usecases"
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

	// 操作タイプに応じて処理を実行
	switch cfg.Operation {
	case "query-anime":
		handleQueryAnime(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

// handleQueryAnime はアニメ情報取得を処理する
func handleQueryAnime(cfg *config.Config) {
	// AniListServiceを初期化
	service := usecases.NewAniListService()

	// アニメ情報を取得
	result, err := service.QueryAnime(cfg.Username, cfg.UserID, cfg.Format, cfg.Limit, cfg.Status, cfg.OutputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Print(result)
}
