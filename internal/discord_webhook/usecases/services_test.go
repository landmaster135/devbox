package usecases

import (
	"context"
	"errors"
	"strings"
	"testing"

	discord "github.com/landmaster135/devbox/internal/discord_webhook/infrastructure/discord"
)

type payloadCall struct {
	botName string
	content string
	embeds  []*discord.Embed
	isTTS   bool
}

type webhookCall struct {
	url     string
	payload *discord.Payload
}

type createEmbedsCall struct {
	title             string
	colorInDecimal    int
	linkOnTitle       string
	footerText        string
	footerIconURL     string
	displaysTimestamp bool
}

type mockDiscordRepository struct {
	payloads          []*discord.Payload
	payloadCalls      []payloadCall
	webhookCalls      []webhookCall
	convertColorCalls []string
	createEmbedsCalls []createEmbedsCall

	availableColors []string

	convertColorErr error
	createEmbedsErr error
	payloadErr      error
	sendWebhookErr  error
}

func newMockDiscordRepository() *mockDiscordRepository {
	return &mockDiscordRepository{
		availableColors: []string{"orange", "blue"},
	}
}

func (m *mockDiscordRepository) SendWebhook(ctx context.Context, webhookURL string, payload *discord.Payload) error {
	m.webhookCalls = append(m.webhookCalls, webhookCall{url: webhookURL, payload: payload})
	if m.sendWebhookErr != nil {
		return m.sendWebhookErr
	}
	return nil
}

func (m *mockDiscordRepository) ConvertColorToDecimal(color string) (int, error) {
	m.convertColorCalls = append(m.convertColorCalls, color)
	if m.convertColorErr != nil {
		return 0, m.convertColorErr
	}
	return 123456, nil
}

func (m *mockDiscordRepository) CreateEmbed(title, description, linkOnTitle string, colorInDecimal int, fields []*discord.EmbedField, footerText, footerIconURL string, displaysTimestamp bool) (*discord.Embed, error) {
	embed := &discord.Embed{
		Title:       title,
		Description: description,
		Color:       colorInDecimal,
		URL:         linkOnTitle,
		Footer: &discord.EmbedFooter{
			Text:    footerText,
			IconURL: footerIconURL,
		},
		Fields: fields,
	}
	if displaysTimestamp {
		embed.Timestamp = "timestamp"
	}
	return embed, nil
}

func (m *mockDiscordRepository) CreateEmbeds(title string, colorInDecimal int, linkOnTitle string, footerText string, footerIconURL string, displaysTimestamp bool) ([]*discord.Embed, error) {
	m.createEmbedsCalls = append(m.createEmbedsCalls, createEmbedsCall{
		title:             title,
		colorInDecimal:    colorInDecimal,
		linkOnTitle:       linkOnTitle,
		footerText:        footerText,
		footerIconURL:     footerIconURL,
		displaysTimestamp: displaysTimestamp,
	})
	if m.createEmbedsErr != nil {
		return nil, m.createEmbedsErr
	}
	embed, _ := m.CreateEmbed(title, "", linkOnTitle, colorInDecimal, nil, footerText, footerIconURL, displaysTimestamp)
	return []*discord.Embed{embed}, nil
}

func (m *mockDiscordRepository) CreatePayload(botName string, content string, embeds []*discord.Embed, isTTS bool) (*discord.Payload, error) {
	if m.payloadErr != nil {
		return nil, m.payloadErr
	}
	payload := &discord.Payload{
		Username: botName,
		Content:  content,
		Embeds:   embeds,
		TTS:      isTTS,
	}
	m.payloads = append(m.payloads, payload)
	m.payloadCalls = append(m.payloadCalls, payloadCall{botName: botName, content: content, embeds: embeds, isTTS: isTTS})
	return payload, nil
}

func (m *mockDiscordRepository) CreateWeatherEmbeds(title, description string, colorInDecimal int, fields []*discord.EmbedField, footerText, footerIconURL string, displaysTimestamp bool) ([]*discord.Embed, error) {
	embed, _ := m.CreateEmbed(title, description, "", colorInDecimal, fields, footerText, footerIconURL, displaysTimestamp)
	return []*discord.Embed{embed}, nil
}

func (m *mockDiscordRepository) GetAvailableColors() []string {
	return m.availableColors
}

func TestDiscordWebhookService_OpenWeather_Defaults(t *testing.T) {
	repo := newMockDiscordRepository()
	service := NewDiscordWebhookService(repo)

	err := service.SendNotification(
		context.Background(),
		"https://discord.com/api/webhooks/test",
		"本文",
		"open-weather-map",
		"",
		"",
		"",
	)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}

	if len(repo.convertColorCalls) != 1 || repo.convertColorCalls[0] != defaultOpenWeatherEmbedColor {
		t.Fatalf("ConvertColorToDecimalの呼び出しが不正です: %#v", repo.convertColorCalls)
	}

	if len(repo.createEmbedsCalls) != 1 {
		t.Fatalf("CreateEmbedsの呼び出し数が不正です: %d", len(repo.createEmbedsCalls))
	}
	call := repo.createEmbedsCalls[0]
	if call.title != defaultOpenWeatherEmbedText {
		t.Errorf("タイトルが期待値と異なります: got %s, want %s", call.title, defaultOpenWeatherEmbedText)
	}
	if call.footerText != footerTextInOpenWeatherEmbed {
		t.Errorf("フッター文言が期待値と異なります: got %s, want %s", call.footerText, footerTextInOpenWeatherEmbed)
	}
	if call.footerIconURL != openWeatherIconURL {
		t.Errorf("フッターアイコンが期待値と異なります: got %s, want %s", call.footerIconURL, openWeatherIconURL)
	}

	if len(repo.payloadCalls) != 1 {
		t.Fatalf("CreatePayloadの呼び出し数が不正です: %d", len(repo.payloadCalls))
	}
	payloadCall := repo.payloadCalls[0]
	if payloadCall.botName != botNameForOpenWeather {
		t.Errorf("Bot名が期待値と異なります: got %s, want %s", payloadCall.botName, botNameForOpenWeather)
	}
	if payloadCall.content != "本文" {
		t.Errorf("本文が期待値と異なります: got %s, want %s", payloadCall.content, "本文")
	}

	if len(repo.webhookCalls) != 1 {
		t.Fatalf("SendWebhookの呼び出し数が不正です: %d", len(repo.webhookCalls))
	}
	if repo.webhookCalls[0].url != "https://discord.com/api/webhooks/test" {
		t.Errorf("WebhookURLが期待値と異なります: got %s", repo.webhookCalls[0].url)
	}
}

func TestDiscordWebhookService_OpenWeather_WithOptions(t *testing.T) {
	repo := newMockDiscordRepository()
	service := NewDiscordWebhookService(repo)

	err := service.SendNotification(
		context.Background(),
		"https://discord.com/api/webhooks/test",
		"本文",
		"open-weather-map",
		"カスタムタイトル",
		"sky_blue",
		"https://example.com",
	)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}

	if len(repo.convertColorCalls) != 1 || repo.convertColorCalls[0] != "sky_blue" {
		t.Errorf("指定した色が使用されていません: %#v", repo.convertColorCalls)
	}

	call := repo.createEmbedsCalls[0]
	if call.title != "カスタムタイトル" {
		t.Errorf("タイトルが期待値と異なります: got %s", call.title)
	}
	if call.linkOnTitle != "https://example.com" {
		t.Errorf("リンクURLが期待値と異なります: got %s", call.linkOnTitle)
	}

	payload := repo.payloads[0]
	if len(payload.Embeds) == 0 || payload.Embeds[0].Title != "カスタムタイトル" {
		t.Fatalf("ペイロードのEmbedタイトルが期待値と異なります")
	}
}

func TestDiscordWebhookService_OpenWeather_ColorConversionError(t *testing.T) {
	repo := newMockDiscordRepository()
	repo.convertColorErr = errors.New("color conversion error")
	service := NewDiscordWebhookService(repo)

	err := service.SendNotification(
		context.Background(),
		"https://discord.com/api/webhooks/test",
		"本文",
		"open-weather-map",
		"",
		"",
		"",
	)
	if err == nil {
		t.Fatal("エラーが期待されましたが、nilでした")
	}
	if !strings.Contains(err.Error(), "色の変換に失敗しました") {
		t.Errorf("色変換エラーに関するメッセージが含まれていません: %v", err)
	}
}

func TestDiscordWebhookService_OpenWeather_CreateEmbedsError(t *testing.T) {
	repo := newMockDiscordRepository()
	repo.createEmbedsErr = errors.New("failed create embeds")
	service := NewDiscordWebhookService(repo)

	err := service.SendNotification(
		context.Background(),
		"https://discord.com/api/webhooks/test",
		"本文",
		"open-weather-map",
		"",
		"",
		"",
	)
	if err == nil {
		t.Fatal("エラーが期待されましたが、nilでした")
	}
	if !strings.Contains(err.Error(), "embedsの作成に失敗しました") {
		t.Errorf("embedsエラーに関するメッセージが含まれていません: %v", err)
	}
}

func TestDiscordWebhookService_OpenWeather_CreatePayloadError(t *testing.T) {
	repo := newMockDiscordRepository()
	repo.payloadErr = errors.New("failed payload")
	service := NewDiscordWebhookService(repo)

	err := service.SendNotification(
		context.Background(),
		"https://discord.com/api/webhooks/test",
		"本文",
		"open-weather-map",
		"",
		"",
		"",
	)
	if err == nil {
		t.Fatal("エラーが期待されましたが、nilでした")
	}
	if !strings.Contains(err.Error(), "ペイロードの作成に失敗しました") {
		t.Errorf("ペイロードエラーに関するメッセージが含まれていません: %v", err)
	}
}
