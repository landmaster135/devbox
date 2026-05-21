package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/movie_extractor/config"
	usecases "github.com/landmaster135/devbox/internal/movie_extractor/usecases"
)

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

	switch cfg.Operation {
	case "extract-frames":
		handleExtractFrames(cfg)
	case "dedup-images":
		handleDedupImages(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作です: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

func handleExtractFrames(cfg *config.Config) {
	service := usecases.NewService()
	result, err := service.HandleExtractFrames(usecases.ExtractFramesInput{
		SrcFile:       cfg.SrcFile,
		FPS:           cfg.FPS,
		Quality:       cfg.Quality,
		StartPosition: cfg.StartPosition,
		OutDir:        cfg.OutDir,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(result)
}

func handleDedupImages(cfg *config.Config) {
	service := usecases.NewService()
	var logWriter *os.File
	if cfg.Log {
		logWriter = os.Stdout
	}

	result, err := service.HandleDedupImages(usecases.DedupImagesInput{
		SrcDir:    cfg.SrcDir,
		MatchRate: cfg.MatchRate,
		Log:       cfg.Log,
		LogWriter: logWriter,
		OutDir:    cfg.OutDir,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(result)
}
