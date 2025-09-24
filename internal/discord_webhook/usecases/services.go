package usecases

import (
	"context"
	"fmt"

	discord "github.com/landmaster135/devbox/internal/discord_webhook/infrastructure/discord"
)

// #==============================================================#
// ##       Const and Types                                      ##
// #==============================================================#
const (
	botNameForVSCODE             = "VSCode生徒会長"
	botNameForOpenWeather        = "お天気あゆ"
	footerTextInVSCODEEmbed      = "VSCode"
	vsCodeIconURL                = "https://code.visualstudio.com/assets/images/code-stable.png"
	defaultVSCODEEmbedText       = "通知"
	defaultVSCODEEmbedColor      = "blue"
	footerTextInOpenWeatherEmbed = "OpenWeatherMap"
	openWeatherIconURL           = "https://openweathermap.org/themes/openweathermap/assets/img/logo_white_cropped.png"
	defaultOpenWeatherEmbedText  = "最新の天気予報"
	defaultOpenWeatherEmbedColor = "orange"
)

// #==============================================================#
// ##       Implementations for DiscordWebhookService            ##
// #==============================================================#
// DiscordWebhookService はDiscord Webhook通知のサービス
type DiscordWebhookService struct {
	repository discord.DiscordClientRepository
}

// NewDiscordWebhookService は新しいDiscordWebhookServiceを作成します
func NewDiscordWebhookService(repository discord.DiscordClientRepository) *DiscordWebhookService {
	return &DiscordWebhookService{
		repository: repository,
	}
}

// NewDefaultDiscordWebhookService はデフォルト設定でDiscordWebhookServiceを作成します
func NewDefaultDiscordWebhookService() *DiscordWebhookService {
	repository := discord.NewDefaultDiscordClient()
	return NewDiscordWebhookService(repository)
}

// #==============================================================#
// ##       Webhook Process                                      ##
// #==============================================================#
// createSimplePayload はembedなしの簡単な通知を送信します
func (s *DiscordWebhookService) createSimplePayload(contentText string) (*discord.Payload, error) {
	// ペイロードを作成
	payload, err := s.repository.CreatePayload(botNameForVSCODE, contentText, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ペイロードの作成に失敗しました: %w", err)
	}

	return payload, nil
}

// createPayloadWithEmbed は共通のEmbed付きペイロード作成処理
func (s *DiscordWebhookService) createPayloadWithEmbed(botName, contentText, embedText, embedColor, embedURLLinkedText, footerText, footerIconURL string) (*discord.Payload, error) {
	colorInDecimal, err := s.repository.ConvertColorToDecimal(embedColor)
	if err != nil {
		availableColors := s.repository.GetAvailableColors()
		return nil, fmt.Errorf("色の変換に失敗しました: %w\n使用可能な色: %v", err, availableColors)
	}

	embeds, err := s.repository.CreateEmbeds(
		embedText,
		colorInDecimal,
		embedURLLinkedText,
		footerText,
		footerIconURL,
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("embedsの作成に失敗しました: %w", err)
	}

	payload, err := s.repository.CreatePayload(botName, contentText, embeds, false)
	if err != nil {
		return nil, fmt.Errorf("ペイロードの作成に失敗しました: %w", err)
	}

	return payload, nil
}

// createVSCodePayload はVSCode風のembed付き通知を送信します
func (s *DiscordWebhookService) createVSCodePayload(contentText, embedText, embedColor, embedURLLinkedText string) (*discord.Payload, error) {
	if embedText == "" {
		embedText = defaultVSCODEEmbedText
	}
	if embedColor == "" {
		embedColor = defaultVSCODEEmbedColor
	}

	return s.createPayloadWithEmbed(
		botNameForVSCODE,
		contentText,
		embedText,
		embedColor,
		embedURLLinkedText,
		footerTextInVSCODEEmbed,
		vsCodeIconURL,
	)
}

func (s *DiscordWebhookService) createOpenWeatherMapPayload(contentText, embedText, embedColor, embedURLLinkedText string) (*discord.Payload, error) {
	if embedText == "" {
		embedText = defaultOpenWeatherEmbedText
	}
	if embedColor == "" {
		embedColor = defaultOpenWeatherEmbedColor
	}

	return s.createPayloadWithEmbed(
		botNameForOpenWeather,
		contentText,
		embedText,
		embedColor,
		embedURLLinkedText,
		footerTextInOpenWeatherEmbed,
		openWeatherIconURL,
	)
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
			return err
		}
	case "open-weather-map":
		payload, err = s.createOpenWeatherMapPayload(contentText, embedText, embedColor, embedURLLinkedText)
		if err != nil {
			return err
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

// #==============================================================#
// ##       Weather Notification Process                         ##
// #==============================================================#
func (s *DiscordWebhookService) CreateWeatherEmbed(title, description string, fields []*discord.EmbedField) (*discord.Embed, error) {
	// 色を10進数に変換
	embedColor := "orange"
	colorInDecimal, err := s.repository.ConvertColorToDecimal(embedColor)
	if err != nil {
		availableColors := s.repository.GetAvailableColors()
		return nil, fmt.Errorf("色の変換に失敗しました: %w\n使用可能な色: %v", err, availableColors)
	}

	embed, err := s.repository.CreateEmbed(
		title,
		description,
		"",
		colorInDecimal,
		fields,
		footerTextInVSCODEEmbed,
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
	payload, err := s.repository.CreatePayload(botNameForOpenWeather, "", embeds, false)
	if err != nil {
		return fmt.Errorf("ペイロードの作成に失敗しました: %w", err)
	}

	// Webhookを送信
	if err := s.repository.SendWebhook(ctx, webhookURL, payload); err != nil {
		return err
	}

	return nil
}
