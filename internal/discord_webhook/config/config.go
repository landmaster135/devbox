package config

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Config はDiscord Webhook CLIの設定を保持する構造体
type Config struct {
	WebhookURL         string // webhook-url (必須)
	ContentText        string // content-text (必須)
	BotName            string // bot-name (任意)
	EmbedType          string // embed-type (none, vscode, postgres, open-weather-map, google-*-success, google-*-failed)
	EmbedText          string // embed-text (任意)
	EmbedColor         string // embed-color (任意)
	EmbedURLLinkedText string // embed-url-linked-text (任意)
	Help               bool   // ヘルプ表示フラグ
}

var validEmbedTypes = []string{
	"none",
	"vscode",
	"postgres",
	"open-weather-map",
	"google-compute-engine-success",
	"google-compute-engine-failed",
	"google-secret-manager-success",
	"google-secret-manager-failed",
	"google-cloud-storage-success",
	"google-cloud-storage-failed",
	"google-cloud-scheduler-success",
	"google-cloud-scheduler-failed",
	"google-cloud-iam-success",
	"google-cloud-iam-failed",
	"google-cloud-run-success",
	"google-cloud-run-failed",
	"google-cloud-run-function-success",
	"google-cloud-run-function-failed",
}

// NewConfig は新しいConfigを作成する
func NewConfig(embedType, webhookURL, botName, contentText, embedText, embedColor, embedURLLinkedText string) (*Config, error) {
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
	if !isValidEmbedType(embedType) {
		return nil, fmt.Errorf("無効なembed-typeです: %s (有効な値: %s)", embedType, embedTypeHelpMessage())
	}

	// webhook-urlの基本的な検証
	if !strings.HasPrefix(webhookURL, "https://discord.com/api/webhooks/") &&
		!strings.HasPrefix(webhookURL, "https://discordapp.com/api/webhooks/") {
		return nil, fmt.Errorf("無効なwebhook-urlです。Discord WebhookのURLを指定してください")
	}

	const BotNameMaxLen = 80
	if len(botName) > BotNameMaxLen {
		return nil, fmt.Errorf("bot-nameは%d文字以下である必要があります", BotNameMaxLen)
	}

	const ContentTextMaxLen = 2000
	if len(contentText) > ContentTextMaxLen {
		return nil, fmt.Errorf("content-textは%d文字以下である必要があります", ContentTextMaxLen)
	}

	// embed-textの長さ検証
	const EmbedTextMaxLen = 256
	if embedText != "" && len(embedText) > EmbedTextMaxLen {
		return nil, fmt.Errorf("embed-textは%d文字以下である必要があります", EmbedTextMaxLen)
	}

	return &Config{
		EmbedType:          embedType,
		WebhookURL:         webhookURL,
		BotName:            botName,
		ContentText:        contentText,
		EmbedText:          embedText,
		EmbedColor:         embedColor,
		EmbedURLLinkedText: embedURLLinkedText,
	}, nil
}

func isValidEmbedType(embedType string) bool {
	for _, et := range validEmbedTypes {
		if embedType == et {
			return true
		}
	}
	return false
}

func embedTypeHelpMessage() string {
	types := make([]string, len(validEmbedTypes))
	copy(types, validEmbedTypes)
	sort.Strings(types)
	return strings.Join(types, ", ")
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
		botName            = ""
		contentText        = ""
		embedText          = ""
		embedColor         = ""
		embedURLLinkedText = ""
		help               = false
	)

	embedUsage := fmt.Sprintf("Embedのタイプ (%s)", embedTypeHelpMessage())

	parser.StringVar(&embedType, "embed-type", embedType, embedUsage)
	parser.StringVar(&embedType, "et", embedType, "embed-typeの短縮形")

	parser.StringVar(&webhookURL, "webhook-url", webhookURL, "Discord WebhookのURL")
	parser.StringVar(&webhookURL, "wu", webhookURL, "webhook-urlの短縮形")

	parser.StringVar(&botName, "bot-name", botName, "ボットの名前")
	parser.StringVar(&botName, "bn", botName, "bot-nameの短縮形")

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

	return NewConfig(embedType, webhookURL, botName, contentText, embedText, embedColor, embedURLLinkedText)
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `Discord Webhook通知CLIツール

使用方法:
  基本的な通知（embedなし）:
    %s -webhook-url "https://discord.com/api/webhooks/..." -bot-name "テストボット" -content-text "Hello, Discord!" -embed-type none

  VSCode風embed付き通知:
    %s -webhook-url "https://discord.com/api/webhooks/..." -content-text "デプロイ完了" -embed-type vscode -embed-text "アプリケーションが正常にデプロイされました"

  PostgreSQLダンプ通知:
    %s -webhook-url "https://discord.com/api/webhooks/..." -content-text "ダンプが完了しました" -embed-type postgres -embed-text "最新のPostgreSQLバックアップ"

  OpenWeatherMap embed付き通知:
    %s -webhook-url "https://discord.com/api/webhooks/..." -content-text "天気予報" -embed-type open-weather-map -embed-text "今日の天気予報"

  フルオプション:
    %s -webhook-url "https://discord.com/api/webhooks/..." -content-text "通知" -embed-type vscode -embed-text "タイトル" -embed-color "green" -embed-url-linked-text "https://example.com"

必須オプション:
  -webhook-url, -wu     Discord WebhookのURL
  -content-text, -ct    メッセージの本文
  -embed-type, -et      Embedのタイプ (%s)

任意オプション:
  -bot-name, -bn            ボットの名前
  -embed-text, -et-text     Embedのタイトル
  -embed-color, -ec         Embedの色 (green, red, blue, yellow, orange, purple, pink, sky_blue, gray_blue, white, black)
  -embed-url-linked-text, -eult  EmbedタイトルのリンクURL
  -help, -h                 このヘルプを表示

embed-typeについて:
  none             : Embedを使用せず、content-textのみを送信
  vscode           : VSCode風のEmbedを使用（フッターにVSCodeアイコンを表示）
  postgres         : VSCode風レイアウトでPostgreSQL通知を送信（専用Bot名とアイコンを使用）
  open-weather-map : 天気予報向けのEmbedをOpenWeatherMap用アイコンで送信
  google-*-success / google-*-failed : Google Cloud各サービスのリクエスト結果を通知（compute-engine, secret-manager, cloud-storage, cloud-scheduler, cloud-iam, cloud-run, cloud-run-function）

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], embedTypeHelpMessage())
}
