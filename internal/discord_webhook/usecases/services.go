package usecases

import (
	"context"
	"fmt"

	discord "github.com/landmaster135/devbox/internal/discord_webhook/infrastructure/discord"
)

const (
	botName           = "VSCode生徒会長"
	footerTextInEmbed = "VSCode"
	vsCodeIconURL     = "https://code.visualstudio.com/assets/images/code-stable.png"
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

// createSimplePayload はembedなしの簡単な通知を送信します
func (s *DiscordWebhookService) createSimplePayload(contentText string) (*discord.Payload, error) {
	// ペイロードを作成
	payload, err := s.repository.CreatePayload(botName, contentText, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ペイロードの作成に失敗しました: %w", err)
	}

	return payload, nil
}

// createVSCodePayload はVSCode風のembed付き通知を送信します
func (s *DiscordWebhookService) createVSCodePayload(contentText, embedText, embedColor, embedURLLinkedText string) (*discord.Payload, error) {
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
		return nil, fmt.Errorf("色の変換に失敗しました: %w\n使用可能な色: %v", err, availableColors)
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
		return nil, fmt.Errorf("embedsの作成に失敗しました: %w", err)
	}

	// ペイロードを作成
	payload, err := s.repository.CreatePayload(botName, contentText, embeds, false)
	if err != nil {
		return nil, fmt.Errorf("ペイロードの作成に失敗しました: %w", err)
	}

	return payload, nil
}

func (s *DiscordWebhookService) SendWebhook(ctx context.Context, webhookURL string, payload *discord.Payload) error {
	// Webhookを送信
	if err := s.repository.SendWebhook(ctx, webhookURL, payload); err != nil {
		return fmt.Errorf("webhook送信に失敗しました: %w", err)
	}
	return nil
}

// SendNotification はDiscordに通知を送信します
func (s *DiscordWebhookService) SendNotification(ctx context.Context, webhookURL, contentText, embedType, embedText, embedColor, embedURLLinkedText string) error {
	// embed-typeに応じた処理
	var payload *discord.Payload
	var err error
	switch embedType {
	case "none":
		payload, err = s.createSimplePayload(contentText)
		if err != nil {
			return fmt.Errorf("ペイロードの作成に失敗しました: %w", err)
		}
	case "vscode":
		payload, err = s.createVSCodePayload(contentText, embedText, embedColor, embedURLLinkedText)
		if err != nil {
			return fmt.Errorf("ペイロードの作成に失敗しました: %w", err)
		}
	default:
		return fmt.Errorf("未対応のembed-typeです: %s", embedType)
	}

	// Webhookを送信
	if err := s.SendWebhook(ctx, webhookURL, payload); err != nil {
		return err
	}

	return nil
}

func (s *DiscordWebhookService) CreateWeatherEmbed(title, description string, fields []*discord.EmbedField) (*discord.Embed, error) {
	// 色を10進数に変換
	embedColor := "orange"
	colorInDecimal, err := s.repository.ConvertColorToDecimal(embedColor)
	if err != nil {
		availableColors := s.repository.GetAvailableColors()
		return nil, fmt.Errorf("色の変換に失敗しました: %w\n使用可能な色: %v", err, availableColors)
	}

	embed, err := s.repository.CreateWeatherEmbed(
		title,
		description,
		colorInDecimal,
		fields,
		footerTextInEmbed,
		vsCodeIconURL,
		true, // タイムスタンプを表示
	)
	if err != nil {
		return nil, fmt.Errorf("embedsの作成に失敗しました: %w", err)
	}
	return embed, nil
}

// SendWeatherNotification は天気予報専用のDiscord通知を送信します
func (s *DiscordWebhookService) SendWeatherNotification(ctx context.Context, webhookURL string, embeds []*discord.Embed) error {
	// ペイロードを作成
	payload, err := s.repository.CreatePayload(botName, "", embeds, false)
	if err != nil {
		return fmt.Errorf("ペイロードの作成に失敗しました: %w", err)
	}

	// Webhookを送信
	if err := s.repository.SendWebhook(ctx, webhookURL, payload); err != nil {
		return err
	}

	return nil
}
