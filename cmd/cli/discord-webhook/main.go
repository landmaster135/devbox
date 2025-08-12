package main

import (
	"context"
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/discord_webhook/config"
	usecases "github.com/landmaster135/devbox/internal/discord_webhook/usecases"
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

	// Discord Webhook通知を実行
	if err := handleDiscordWebhookNotification(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Discord通知が正常に送信されました")
}

// handleDiscordWebhookNotification はDiscord Webhook通知を処理する
func handleDiscordWebhookNotification(cfg *config.Config) error {
	// DiscordWebhookServiceを初期化
	service := usecases.NewDefaultDiscordWebhookService()

	// コンテキストを作成
	ctx := context.Background()

	// 通知を送信
	err := service.SendNotification(
		ctx,
		cfg.WebhookURL,
		cfg.ContentText,
		cfg.EmbedType,
		cfg.EmbedText,
		cfg.EmbedColor,
		cfg.EmbedURLLinkedText,
	)
	if err != nil {
		return fmt.Errorf("discord通知の送信に失敗しました: %w", err)
	}

	return nil
}
