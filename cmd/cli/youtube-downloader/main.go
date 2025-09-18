package main

import (
	"context"
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/youtube_downloader/config"
	domain "github.com/landmaster135/devbox/internal/youtube_downloader/domain"
	usecases "github.com/landmaster135/devbox/internal/youtube_downloader/usecases"
)

func handleDownload(cfg *config.Config) {
	// サービスを作成
	service := usecases.NewServiceWithDefaults()

	// ダウンロード要求を作成
	request := domain.DownloadRequest{
		URL:         cfg.URL,
		OutputDir:   cfg.OutputDir,
		Quality:     cfg.Quality,
		Format:      cfg.Format,
		AudioOnly:   cfg.AudioOnly,
		Playlist:    cfg.Playlist,
		MaxRoutines: cfg.MaxRoutines,
		ChunkSize:   cfg.ChunkSize,
	}

	// ダウンロード実行
	ctx := context.Background()
	result, err := service.DownloadVideo(ctx, request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を標準出力に表示
	fmt.Print(result)
}

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	if cfg.Help {
		config.PrintUsage()
		return
	}

	handleDownload(cfg)
}
