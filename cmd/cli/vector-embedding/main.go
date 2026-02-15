package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/landmaster135/devbox/internal/vector_embedding/config"
	"github.com/landmaster135/devbox/internal/vector_embedding/usecases"
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

	service, err := usecases.NewService(usecases.Options{
		Host:    cfg.Host,
		Port:    cfg.Port,
		Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
		APIKey:  cfg.APIKey,
	})
	if err != nil {
		exitWithError(err)
	}

	ctx := context.Background()
	result, err := service.Embed(ctx, cfg)
	if err != nil {
		exitWithError(err)
	}

	if err := printJSON(result); err != nil {
		exitWithError(err)
	}
}

func printJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("結果の整形に失敗しました: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
	os.Exit(1)
}
