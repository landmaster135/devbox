package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Logger はロギング機能を提供するインターフェース
type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, err error, keysAndValues ...interface{})
}

// NoopLogger はログ出力を行わないロガー実装
type NoopLogger struct{}

// Info は情報ログを出力します（何も行いません）
func (l *NoopLogger) Info(msg string, keysAndValues ...interface{}) {}

// Error はエラーログを出力します（何も行いません）
func (l *NoopLogger) Error(msg string, err error, keysAndValues ...interface{}) {}

// DiscordClient はHTTPを使用してDiscord APIと通信するクライアント
type DiscordClient struct {
	client *http.Client
	logger Logger
}

// NewDiscordClient は新しいDiscordClientを作成します
func NewDiscordClient(logger Logger) *DiscordClient {
	// ロガーが指定されていない場合はNoopLoggerを使用
	if logger == nil {
		logger = &NoopLogger{}
	}

	client := &http.Client{
		Timeout: time.Second * 30, // タイムアウトを30秒に延長
	}
	return &DiscordClient{
		client: client,
		logger: logger,
	}
}

// NewDefaultDiscordClient はデフォルト設定でDiscordClientを作成します
func NewDefaultDiscordClient() *DiscordClient {
	return NewDiscordClient(nil)
}

// SendWebhook はDiscordのwebhookにメッセージを送信します
func (c *DiscordClient) SendWebhook(ctx context.Context, webhookURL string, payload *Payload) error {
	// JSONにエンコード
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// リクエストを作成
	req, err := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// ヘッダーを設定
	req.Header.Set("Content-Type", "application/json")

	// リクエストを送信
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスを確認
	if resp.StatusCode == 204 {
		c.logger.Info("notification sent successfully to Discord")
	} else {
		// エラーレスポンスを読み取り
		var respBody []byte
		if resp.Body != nil {
			respBody = make([]byte, 1024)
			resp.Body.Read(respBody)
		}
		c.logger.Error(
			"failed to send notification",
			fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(respBody)),
		)
		return fmt.Errorf("failed to send notification: %d - %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// AvailableColor は使用可能な色の情報を保持する構造体
type AvailableColor struct {
	name          string
	codeInDecimal int
}

// EmbedFooter はEmbedのフッター情報を保持する構造体
type EmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

// EmbedField はEmbedのフィールド情報を保持する構造体
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// Embed はDiscord Embedの構造体
type Embed struct {
	Title       string        `json:"title"`
	Description string        `json:"description,omitempty"`
	Color       int           `json:"color"`
	URL         string        `json:"url,omitempty"`
	Timestamp   string        `json:"timestamp,omitempty"`
	Footer      *EmbedFooter  `json:"footer,omitempty"`
	Fields      []*EmbedField `json:"fields,omitempty"`
}

// Payload はDiscord Webhookのペイロード構造体
type Payload struct {
	Username string   `json:"username"`
	Content  string   `json:"content"`
	TTS      bool     `json:"tts"`
	Embeds   []*Embed `json:"embeds,omitempty"`
}

// getAvailableColors は使用可能な色の配列を取得します
func (c *DiscordClient) getAvailableColors() []AvailableColor {
	return []AvailableColor{
		{name: "green", codeInDecimal: 4569935},     // 0x45BB4F
		{name: "red", codeInDecimal: 16711680},      // 0xFF0000
		{name: "sky_blue", codeInDecimal: 52479},    // 0xCCFF
		{name: "orange", codeInDecimal: 14177041},   // 0xD87D11
		{name: "white", codeInDecimal: 16777215},    // 0xFFFFFF
		{name: "blue", codeInDecimal: 39423},        // 0x99FF
		{name: "yellow", codeInDecimal: 16770560},   // 0xFFE000
		{name: "pink", codeInDecimal: 16711833},     // 0xFF0099
		{name: "purple", codeInDecimal: 10494192},   // 0xA020F0
		{name: "gray_blue", codeInDecimal: 9212588}, // 0x8C92AC
		{name: "black", codeInDecimal: 3355443},     // 0x333333
	}
}

// ConvertColorToDecimal は文字列の色名を10進数の色コードに変換します
func (c *DiscordClient) ConvertColorToDecimal(color string) (int, error) {
	if color == "" {
		return 0, fmt.Errorf("'color' must not be empty")
	}

	colors := c.getAvailableColors()
	for _, availableColor := range colors {
		if availableColor.name == color {
			return availableColor.codeInDecimal, nil
		}
	}

	return 0, fmt.Errorf("'color' is an unexpected value: %s", color)
}

// CreateEmbed は単一のDiscord Embedを作成します
func (c *DiscordClient) CreateEmbed(title string, colorInDecimal int, linkOnTitle string, footerText string, footerIconURL string, displaysTimestamp bool) (*Embed, error) {
	// 入力値の検証
	if title == "" {
		return nil, fmt.Errorf("'title' must not be empty")
	}
	if footerText == "" {
		return nil, fmt.Errorf("'footerText' must not be empty")
	}
	if footerIconURL == "" {
		return nil, fmt.Errorf("'footerIconURL' must not be empty")
	}

	// embedを作成
	embed := &Embed{
		Title: title,
		Color: colorInDecimal,
		Footer: &EmbedFooter{
			Text:    footerText,
			IconURL: footerIconURL,
		},
	}

	// リンクが指定されている場合は追加
	if linkOnTitle != "" {
		embed.URL = linkOnTitle
	}

	// タイムスタンプを表示する場合は追加
	if displaysTimestamp {
		// ISO 8601形式のタイムスタンプを生成
		embed.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	return embed, nil
}

// CreateEmbeds はDiscord通知用のembedsを作成します
func (c *DiscordClient) CreateEmbeds(title string, colorInDecimal int, linkOnTitle string, footerText string, footerIconURL string, displaysTimestamp bool) ([]*Embed, error) {
	embed, err := c.CreateEmbed(title, colorInDecimal, linkOnTitle, footerText, footerIconURL, displaysTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to create embed: %v", err)
	}

	// 入力値の検証
	if embed == nil {
		return nil, fmt.Errorf("embed must not be nil")
	}

	// 配列に格納して返す
	embeds := []*Embed{embed}
	return embeds, nil
}

// CreatePayload はDiscord通知用のペイロードを作成します
func (c *DiscordClient) CreatePayload(botName string, content string, embeds []*Embed, isTTS bool) (*Payload, error) {
	// 入力値の検証
	if botName == "" {
		return nil, fmt.Errorf("'botName' must not be empty")
	}
	if len(botName) > 2000 {
		return nil, fmt.Errorf("'botName' must be equals or smaller than 2000")
	}
	if len(content) > 2000 {
		return nil, fmt.Errorf("'content' must be equals or smaller than 2000")
	}
	if len(embeds) > 10 {
		return nil, fmt.Errorf("length of 'embeds' must be equals or smaller than 10")
	}

	// ペイロードを作成
	payload := &Payload{
		Username: botName,
		Content:  content,
		TTS:      isTTS,
		Embeds:   embeds,
	}

	c.logger.Info("payload created", "payload", payload)
	return payload, nil
}

// CreateWeatherEmbed は天気予報専用のDiscord Embedを作成します
func (c *DiscordClient) CreateWeatherEmbed(title, description string, colorInDecimal int, fields []*EmbedField, footerText, footerIconURL string, displaysTimestamp bool) (*Embed, error) {
	// 入力値の検証
	if title == "" {
		return nil, fmt.Errorf("'title' must not be empty")
	}
	if footerText == "" {
		return nil, fmt.Errorf("'footerText' must not be empty")
	}
	if footerIconURL == "" {
		return nil, fmt.Errorf("'footerIconURL' must not be empty")
	}
	if len(fields) > 25 {
		return nil, fmt.Errorf("fields count must be 25 or less")
	}

	// embedを作成
	embed := &Embed{
		Title:       title,
		Description: description,
		Color:       colorInDecimal,
		Fields:      fields,
		Footer: &EmbedFooter{
			Text:    footerText,
			IconURL: footerIconURL,
		},
	}

	// タイムスタンプを表示する場合は追加
	if displaysTimestamp {
		// ISO 8601形式のタイムスタンプを生成
		embed.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	return embed, nil
}

// CreateWeatherEmbeds は天気予報専用のDiscord Embedsを作成します
func (c *DiscordClient) CreateWeatherEmbeds(title, description string, colorInDecimal int, fields []*EmbedField, footerText, footerIconURL string, displaysTimestamp bool) ([]*Embed, error) {
	embed, err := c.CreateWeatherEmbed(title, description, colorInDecimal, fields, footerText, footerIconURL, displaysTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to create weather embed: %v", err)
	}

	// 入力値の検証
	if embed == nil {
		return nil, fmt.Errorf("embed must not be nil")
	}

	// 配列に格納して返す
	embeds := []*Embed{embed}
	return embeds, nil
}

// GetAvailableColors は使用可能な色のリストを取得します
func (c *DiscordClient) GetAvailableColors() []string {
	colors := c.getAvailableColors()
	colorNames := make([]string, len(colors))
	for i, color := range colors {
		colorNames[i] = color.name
	}
	return colorNames
}
