package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	discordUsecases "github.com/landmaster135/devbox/internal/discord_webhook/usecases"
	config "github.com/landmaster135/devbox/internal/gcloud_wrapper_workload_identity_federation/config"
	usecases "github.com/landmaster135/devbox/internal/gcloud_wrapper_workload_identity_federation/usecases"
)

func handleWorkloadIdentitySetup(cfg *config.Config) {
	// サービスの初期化
	service := usecases.NewService()
	discordService := discordUsecases.NewDefaultDiscordWebhookService()

	// 設定の変換
	workloadConfig := &usecases.WorkloadIdentityConfig{
		ProjectID:        cfg.ProjectID,
		PoolID:           cfg.PoolID,
		ProviderID:       cfg.ProviderID,
		ServiceAccountID: cfg.ServiceAccountID,
		Location:         cfg.Location,
		PoolDescription:  cfg.PoolDescription,
		RepoOwner:        cfg.RepoOwner,
		RepoName:         cfg.RepoName,
	}

	// 通知メッセージの生成（自動分割）
	messages, err := service.GenerateNotificationMessages(workloadConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: 通知メッセージの生成に失敗しました: %v\n", err)
		os.Exit(1)
	}

	// メッセージ生成結果の表示（テスト用）
	fmt.Printf("\n=== 生成されたメッセージ (%d件) ===\n", len(messages))
	for i, message := range messages {
		fmt.Printf("\n--- メッセージ %d: %s ---\n", i+1, message.Title)
		fmt.Printf("色: %s, コード: %t\n", message.Color, message.IsCode)
		fmt.Printf("文字数: %d文字\n", len(message.Content))

		// 内容の一部を表示（1行目のみ）
		lines := strings.Split(message.Content, "\n")
		firstLine := ""
		if len(lines) > 0 {
			firstLine = lines[0]
		}
		fmt.Printf("内容（抜粋）:\n%s\n", firstLine)
	}

	// Discord通知の送信（複数回）
	ctx := context.Background()
	var successCount, failureCount int

	for i, message := range messages {
		fmt.Printf("\nDiscord通知を送信中... (%d/%d): %s\n", i+1, len(messages), message.Title)

		err = discordService.SendNotification(
			ctx,
			cfg.WebhookURL,
			message.Content,
			"vscode",
			message.Title,
			message.Color,
			"",
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "警告: Discord通知の送信に失敗しました (%d/%d): %v\n", i+1, len(messages), err)
			failureCount++
		} else {
			fmt.Printf("Discord通知を送信しました (%d/%d): %s\n", i+1, len(messages), message.Title)
			successCount++
		}
	}

	// 結果の表示
	fmt.Print("\n=== 処理結果 ===\n")
	fmt.Printf("Workload Identity Federationのセットアップスクリプトを生成しました。\n")
	fmt.Printf("プロジェクト: %s\n", cfg.ProjectID)
	fmt.Printf("リポジトリ: %s/%s\n", cfg.RepoOwner, cfg.RepoName)
	fmt.Printf("Discord通知: %d件成功, %d件失敗\n", successCount, failureCount)

	if successCount > 0 {
		fmt.Print("Discordで詳細な手順とスクリプトを確認してください。\n")
	}

	if failureCount > 0 {
		fmt.Fprintf(os.Stderr, "\n注意: 一部のDiscord通知が失敗しました。Webhook URLを確認してください。\n")
		os.Exit(1)
	}
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

	handleWorkloadIdentitySetup(cfg)
}
