package config

import (
	"fmt"
	"os"
	"strings"
)

// Config はDiscord Webhook CLIの設定を保持する構造体
type Config struct {
	WebhookURL         string // webhook-url (必須)
	ContentText        string // content-text (必須)
	EmbedType          string // embed-type (none, vscode)
	EmbedText          string // embed-text (任意)
	EmbedColor         string // embed-color (任意)
	EmbedURLLinkedText string // embed-url-linked-text (任意)
	Help               bool   // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(embedType, webhookURL, contentText, embedText, embedColor, embedURLLinkedText string) (*Config, error) {
	// 必須パラメータの検証
	if webhookURL == "" {
		return nil, fmt.Errorf("webhook-urlが指定されていません")
	}
	if contentText == "" {
		return nil, fmt.Errorf("content-textが指定されていません")
	}
	if embedType == "" {
		return nil, fmt.Errorf("embed-typeが指定されていません")
	}

	// embed-typeの検証
	validEmbedTypes := []string{"none", "vscode"}
	isValid := false
	for _, et := range validEmbedTypes {
		if embedType == et {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, fmt.Errorf("無効なembed-typeです: %s (有効な値: none, vscode)", embedType)
	}

	// webhook-urlの基本的な検証
	if !strings.HasPrefix(webhookURL, "https://discord.com/api/webhooks/") &&
		!strings.HasPrefix(webhookURL, "https://discordapp.com/api/webhooks/") {
		return nil, fmt.Errorf("無効なwebhook-urlです。Discord WebhookのURLを指定してください")
	}

	// content-textの長さ検証
	if len(contentText) > 2000 {
		return nil, fmt.Errorf("content-textは2000文字以下である必要があります")
	}

	// embed-textの長さ検証
	if embedText != "" && len(embedText) > 256 {
		return nil, fmt.Errorf("embed-textは256文字以下である必要があります")
	}

	return &Config{
		EmbedType:          embedType,
		WebhookURL:         webhookURL,
		ContentText:        contentText,
		EmbedText:          embedText,
		EmbedColor:         embedColor,
		EmbedURLLinkedText: embedURLLinkedText,
	}, nil
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		embedType          = ""
		webhookURL         = ""
		contentText        = ""
		embedText          = ""
		embedColor         = ""
		embedURLLinkedText = ""
		help               = false
	)

	parser.StringVar(&embedType, "embed-type", embedType, "Embedのタイプ (none, vscode)")
	parser.StringVar(&embedType, "et", embedType, "embed-typeの短縮形")

	parser.StringVar(&webhookURL, "webhook-url", webhookURL, "Discord WebhookのURL")
	parser.StringVar(&webhookURL, "wu", webhookURL, "webhook-urlの短縮形")

	parser.StringVar(&contentText, "content-text", contentText, "メッセージの本文")
	parser.StringVar(&contentText, "ct", contentText, "content-textの短縮形")

	parser.StringVar(&embedText, "embed-text", embedText, "Embedのタイトル（任意）")
	parser.StringVar(&embedText, "et-text", embedText, "embed-textの短縮形")

	parser.StringVar(&embedColor, "embed-color", embedColor, "Embedの色（任意）")
	parser.StringVar(&embedColor, "ec", embedColor, "embed-colorの短縮形")

	parser.StringVar(&embedURLLinkedText, "embed-url-linked-text", embedURLLinkedText, "EmbedタイトルのリンクURL（任意）")
	parser.StringVar(&embedURLLinkedText, "eult", embedURLLinkedText, "embed-url-linked-textの短縮形")

	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	return NewConfig(embedType, webhookURL, contentText, embedText, embedColor, embedURLLinkedText)
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `Discord Webhook通知CLIツール

使用方法:
  基本的な通知（embedなし）:
    %s -webhook-url "https://discord.com/api/webhooks/..." -content-text "Hello, Discord!" -embed-type none

  VSCode風embed付き通知:
    %s -webhook-url "https://discord.com/api/webhooks/..." -content-text "デプロイ完了" -embed-type vscode -embed-text "アプリケーションが正常にデプロイされました"

  フルオプション:
    %s -webhook-url "https://discord.com/api/webhooks/..." -content-text "通知" -embed-type vscode -embed-text "タイトル" -embed-color "green" -embed-url-linked-text "https://example.com"

必須オプション:
  -webhook-url, -wu     Discord WebhookのURL
  -content-text, -ct    メッセージの本文
  -embed-type, -et      Embedのタイプ (none, vscode)

任意オプション:
  -embed-text, -et-text     Embedのタイトル
  -embed-color, -ec         Embedの色 (green, red, blue, yellow, orange, purple, pink, sky_blue, gray_blue, white, black)
  -embed-url-linked-text, -eult  EmbedタイトルのリンクURL
  -help, -h                 このヘルプを表示

embed-typeについて:
  none   : Embedを使用せず、content-textのみを送信
  vscode : VSCode風のEmbedを使用（フッターにVSCodeアイコンを表示）

`, os.Args[0], os.Args[0], os.Args[0])
}
