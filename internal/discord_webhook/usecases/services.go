package usecases

import (
	"context"
	"fmt"

	discord "github.com/landmaster135/devbox/internal/discord_webhook/infrastructure/discord"
)

const (
	botName = "VSCode生徒会長"
	footerTextInEmbed = "VSCode"
	vsCodeIconURL = "https://code.visualstudio.com/assets/images/code-stable.png"
)

// DiscordWebhookService はDiscord Webhook通知のサービス
type DiscordWebhookService struct {
	repository discord.DiscordRepository
}

// NewDiscordWebhookService は新しいDiscordWebhookServiceを作成します
func NewDiscordWebhookService(repository discord.DiscordRepository) *DiscordWebhookService {
	return &DiscordWebhookService{
		repository: repository,
	}
}

// NewDefaultDiscordWebhookService はデフォルト設定でDiscordWebhookServiceを作成します
func NewDefaultDiscordWebhookService() *DiscordWebhookService {
	repository := discord.NewDefaultDiscordRepository()
	return NewDiscordWebhookService(repository)
}

// SendNotification はDiscordに通知を送信します
func (s *DiscordWebhookService) SendNotification(ctx context.Context, webhookURL, contentText, embedType, embedText, embedColor, embedURLLinkedText string) error {
	// embed-typeに応じた処理
	switch embedType {
	case "none":
		return s.sendSimpleNotification(ctx, webhookURL, contentText)
	case "vscode":
		return s.sendVSCodeNotification(ctx, webhookURL, contentText, embedText, embedColor, embedURLLinkedText)
	default:
		return fmt.Errorf("未対応のembed-typeです: %s", embedType)
	}
}

// sendSimpleNotification はembedなしの簡単な通知を送信します
func (s *DiscordWebhookService) sendSimpleNotification(ctx context.Context, webhookURL, contentText string) error {
	// embedなしのペイロードを作成
	payload, err := s.repository.CreatePayload("Webhook Bot", contentText, nil, false)
	if err != nil {
		return fmt.Errorf("ペイロードの作成に失敗しました: %w", err)
	}

	// Webhookを送信
	if err := s.repository.SendWebhook(ctx, webhookURL, payload); err != nil {
		return fmt.Errorf("webhook送信に失敗しました: %w", err)
	}

	return nil
}

// SendWeatherNotification は天気予報専用のDiscord通知を送信します
func (s *DiscordWebhookService) SendWeatherNotification(ctx context.Context, webhookURL, title, description, embedColor string, fields []*discord.EmbedField) error {
	// 色を10進数に変換
	colorInDecimal, err := s.repository.ConvertColorToDecimal(embedColor)
	if err != nil {
		availableColors := s.repository.GetAvailableColors()
		return fmt.Errorf("色の変換に失敗しました: %w\n使用可能な色: %v", err, availableColors)
	}

	// 天気予報専用のembedsを作成
	embeds, err := s.repository.CreateWeatherEmbeds(
		title,
		description,
		colorInDecimal,
		fields,
		footerTextInEmbed,
		vsCodeIconURL,
		true, // タイムスタンプを表示
	)
	if err != nil {
		return fmt.Errorf("embedsの作成に失敗しました: %w", err)
	}

	// ペイロードを作成
	payload, err := s.repository.CreatePayload(botName, "", embeds, false)
	if err != nil {
		return fmt.Errorf("ペイロードの作成に失敗しました: %w", err)
	}

	// Webhookを送信
	if err := s.repository.SendWebhook(ctx, webhookURL, payload); err != nil {
		return fmt.Errorf("webhook送信に失敗しました: %w", err)
	}

	return nil
}

// sendVSCodeNotification はVSCode風のembed付き通知を送信します
func (s *DiscordWebhookService) sendVSCodeNotification(ctx context.Context, webhookURL, contentText, embedText, embedColor, embedURLLinkedText string) error {
	// デフォルト値の設定
	if embedText == "" {
		embedText = "通知"
	}
	if embedColor == "" {
		embedColor = "blue"
	}

	// 色を10進数に変換
	colorInDecimal, err := s.repository.ConvertColorToDecimal(embedColor)
	if err != nil {
		availableColors := s.repository.GetAvailableColors()
		return fmt.Errorf("色の変換に失敗しました: %w\n使用可能な色: %v", err, availableColors)
	}

	// embedsを作成
	embeds, err := s.repository.CreateEmbeds(
		embedText,
		colorInDecimal,
		embedURLLinkedText,
		footerTextInEmbed,
		vsCodeIconURL,
		true, // タイムスタンプを表示
	)
	if err != nil {
		return fmt.Errorf("embedsの作成に失敗しました: %w", err)
	}

	// ペイロードを作成
	payload, err := s.repository.CreatePayload(botName, contentText, embeds, false)
	if err != nil {
		return fmt.Errorf("ペイロードの作成に失敗しました: %w", err)
	}

	// Webhookを送信
	if err := s.repository.SendWebhook(ctx, webhookURL, payload); err != nil {
		return fmt.Errorf("webhook送信に失敗しました: %w", err)
	}

	return nil
}
