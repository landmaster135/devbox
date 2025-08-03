package discord

import (
	"context"
)

// DiscordRepositoryImpl はDiscord通知のリポジトリ実装
type DiscordRepositoryImpl struct {
	client *DiscordClient
}

// NewDiscordRepository は新しいDiscordRepositoryを作成します
func NewDiscordRepository(logger Logger) DiscordRepository {
	client := NewDiscordClient(logger)
	return &DiscordRepositoryImpl{
		client: client,
	}
}

// NewDefaultDiscordRepository はデフォルト設定でDiscordRepositoryを作成します
func NewDefaultDiscordRepository() DiscordRepository {
	return NewDiscordRepository(nil)
}

// SendWebhook はDiscordのwebhookにメッセージを送信します
func (r *DiscordRepositoryImpl) SendWebhook(ctx context.Context, webhookURL string, payload *Payload) error {
	return r.client.SendWebhook(ctx, webhookURL, payload)
}

// ConvertColorToDecimal は文字列の色名を10進数の色コードに変換します
func (r *DiscordRepositoryImpl) ConvertColorToDecimal(color string) (int, error) {
	return r.client.ConvertColorToDecimal(color)
}

func (r *DiscordRepositoryImpl) CreateEmbed(title string, colorInDecimal int, linkOnTitle string, footerText string, footerIconURL string, displaysTimestamp bool) (*Embed, error) {
	return r.client.CreateEmbed(title, colorInDecimal, linkOnTitle, footerText, footerIconURL, displaysTimestamp)
}

// CreateEmbeds はDiscord通知用のembedsを作成します
func (r *DiscordRepositoryImpl) CreateEmbeds(title string, colorInDecimal int, linkOnTitle string, footerText string, footerIconURL string, displaysTimestamp bool) ([]*Embed, error) {
	return r.client.CreateEmbeds(title, colorInDecimal, linkOnTitle, footerText, footerIconURL, displaysTimestamp)
}

// CreatePayload はDiscord通知用のペイロードを作成します
func (r *DiscordRepositoryImpl) CreatePayload(botName string, content string, embeds []*Embed, isTTS bool) (*Payload, error) {
	return r.client.CreatePayload(botName, content, embeds, isTTS)
}

// GetAvailableColors は使用可能な色のリストを取得します
func (r *DiscordRepositoryImpl) GetAvailableColors() []string {
	return r.client.GetAvailableColors()
}
