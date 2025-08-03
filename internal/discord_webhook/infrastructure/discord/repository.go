package discord

import (
	"context"
)

// DiscordRepository はDiscord通知のリポジトリインターフェース
type DiscordRepository interface {
	// SendWebhook はDiscordのwebhookにメッセージを送信します
	SendWebhook(ctx context.Context, webhookURL string, payload *Payload) error

	// ConvertColorToDecimal は文字列の色名を10進数の色コードに変換します
	ConvertColorToDecimal(color string) (int, error)

	CreateEmbed(title string, colorInDecimal int, linkOnTitle string, footerText string, footerIconURL string, displaysTimestamp bool) (*Embed, error)

	// CreateEmbeds はDiscord通知用のembedsを作成します
	CreateEmbeds(title string, colorInDecimal int, linkOnTitle string, footerText string, footerIconURL string, displaysTimestamp bool) ([]*Embed, error)

	// CreatePayload はDiscord通知用のペイロードを作成します
	CreatePayload(botName string, content string, embeds []*Embed, isTTS bool) (*Payload, error)

	// CreateWeatherEmbed は天気予報専用のDiscord Embedを作成します
	CreateWeatherEmbed(title, description string, colorInDecimal int, fields []*EmbedField, footerText, footerIconURL string, displaysTimestamp bool) (*Embed, error)

	// CreateWeatherEmbeds は天気予報専用のDiscord Embedsを作成します
	CreateWeatherEmbeds(title, description string, colorInDecimal int, fields []*EmbedField, footerText, footerIconURL string, displaysTimestamp bool) ([]*Embed, error)

	// GetAvailableColors は使用可能な色のリストを取得します
	GetAvailableColors() []string
}
