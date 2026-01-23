package usecases

import (
	"context"
	"fmt"
	"os"
	"strings"

	discord "github.com/landmaster135/devbox/internal/discord_webhook/infrastructure/discord"
)

// #==============================================================#
// ##       Const and Types                                      ##
// #==============================================================#
const (
	// VSCode
	botNameForVSCODE        = "VSCode生徒会長"
	footerTextInVSCODEEmbed = "VSCode"
	vsCodeIconURL           = "https://code.visualstudio.com/assets/images/code-stable.png"
	defaultVSCODEEmbedText  = "通知"
	defaultVSCODEEmbedColor = "blue"

	// Open Weather
	botNameForOpenWeather        = "お天気あゆ"
	footerTextInOpenWeatherEmbed = "OpenWeatherMap"
	openWeatherIconURL           = "https://avatars.githubusercontent.com/u/1743227?s=200&v=4"
	defaultOpenWeatherEmbedText  = "最新の天気予報"
	defaultOpenWeatherEmbedColor = "orange"

	// gcloud
	botNameForGCloud           = "クラウドウォッチャーあゆ"
	defaultGCloudSuccessColor  = "green"
	defaultGCloudFailureColor  = "red"
	defaultGCloudFooterIconURL = "https://www.vectorlogo.zone/logos/google_cloud/google_cloud-icon.svg"
)

type gcloudNotificationOutcome string

const (
	gcloudOutcomeSuccess gcloudNotificationOutcome = "success"
	gcloudOutcomeFailed  gcloudNotificationOutcome = "failed"

	gcloudSuccessSuffix = "-success"
	gcloudFailedSuffix  = "-failed"
)

type gcloudEmbedProfile struct {
	footerText     string
	displayName    string
	iconEnvVar     string
	defaultIconURL string
}

func (p gcloudEmbedProfile) iconURL() string {
	if envValue := os.Getenv(p.iconEnvVar); envValue != "" {
		return envValue
	}
	return p.defaultIconURL
}

func (p gcloudEmbedProfile) defaultMessage(outcome gcloudNotificationOutcome) string {
	switch outcome {
	case gcloudOutcomeSuccess:
		return fmt.Sprintf("%sのリクエストが成功しました", p.displayName)
	case gcloudOutcomeFailed:
		return fmt.Sprintf("%sのリクエストでエラーが発生しました", p.displayName)
	default:
		return p.displayName
	}
}

var gcloudEmbedProfiles = map[string]gcloudEmbedProfile{
	"google-compute-engine": {
		footerText:     "GoogleComputeEngine",
		displayName:    "Google Compute Engine",
		iconEnvVar:     "GCE_ICON_URL",
		defaultIconURL: defaultGCloudFooterIconURL,
	},
	"google-secret-manager": {
		footerText:     "GoogleSecretManager",
		displayName:    "Google Secret Manager",
		iconEnvVar:     "GSM_ICON_URL",
		defaultIconURL: defaultGCloudFooterIconURL,
	},
	"google-cloud-storage": {
		footerText:     "GoogleCloudStorage",
		displayName:    "Google Cloud Storage",
		iconEnvVar:     "GCS_ICON_URL",
		defaultIconURL: defaultGCloudFooterIconURL,
	},
	"google-cloud-scheduler": {
		footerText:     "GoogleCloudScheduler",
		displayName:    "Google Cloud Scheduler",
		iconEnvVar:     "GCSCHEDULER_ICON_URL",
		defaultIconURL: defaultGCloudFooterIconURL,
	},
	"google-cloud-iam": {
		footerText:     "GoogleCloudIAM",
		displayName:    "Google Cloud IAM",
		iconEnvVar:     "GCIAM_ICON_URL",
		defaultIconURL: defaultGCloudFooterIconURL,
	},
	"google-cloud-run": {
		footerText:     "GoogleCloudRun",
		displayName:    "Google Cloud Run",
		iconEnvVar:     "GCLOUD_RUN_ICON_URL",
		defaultIconURL: defaultGCloudFooterIconURL,
	},
	"google-cloud-run-function": {
		footerText:     "GoogleCloudRunFunction",
		displayName:    "Google Cloud Run Functions",
		iconEnvVar:     "GCLOUD_RUN_FUNCTION_ICON_URL",
		defaultIconURL: defaultGCloudFooterIconURL,
	},
}

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
func (s *DiscordWebhookService) createSimplePayload(botName, contentText string) (*discord.Payload, error) {
	// ペイロードを作成
	payload, err := s.repository.CreatePayload(botName, contentText, nil, false)
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

func parseGcloudEmbedType(embedType string) (string, gcloudNotificationOutcome, error) {
	switch {
	case strings.HasSuffix(embedType, gcloudSuccessSuffix):
		base := strings.TrimSuffix(embedType, gcloudSuccessSuffix)
		if base == "" {
			return "", "", fmt.Errorf("無効なembed-typeです: %s", embedType)
		}
		return base, gcloudOutcomeSuccess, nil
	case strings.HasSuffix(embedType, gcloudFailedSuffix):
		base := strings.TrimSuffix(embedType, gcloudFailedSuffix)
		if base == "" {
			return "", "", fmt.Errorf("無効なembed-typeです: %s", embedType)
		}
		return base, gcloudOutcomeFailed, nil
	default:
		return "", "", fmt.Errorf("未対応のGoogle Cloud embed-typeです: %s", embedType)
	}
}

func isGcloudEmbedType(embedType string) bool {
	base, _, err := parseGcloudEmbedType(embedType)
	if err != nil {
		return false
	}
	_, exists := gcloudEmbedProfiles[base]
	return exists
}

func (s *DiscordWebhookService) createGcloudPayload(contentText, embedType, embedText, embedColor, embedURLLinkedText string) (*discord.Payload, error) {
	base, outcome, err := parseGcloudEmbedType(embedType)
	if err != nil {
		return nil, err
	}

	profile, exists := gcloudEmbedProfiles[base]
	if !exists {
		return nil, fmt.Errorf("未対応のGoogle Cloudサービスです: %s", base)
	}

	iconURL := profile.iconURL()
	if iconURL == "" {
		return nil, fmt.Errorf("google cloud通知用のアイコンURLが設定されていません (環境変数: %s)", profile.iconEnvVar)
	}

	resolvedEmbedText := embedText
	if resolvedEmbedText == "" {
		resolvedEmbedText = profile.defaultMessage(outcome)
	}

	resolvedColor := embedColor
	if resolvedColor == "" {
		if outcome == gcloudOutcomeSuccess {
			resolvedColor = defaultGCloudSuccessColor
		} else {
			resolvedColor = defaultGCloudFailureColor
		}
	}

	return s.createPayloadWithEmbed(
		botNameForGCloud,
		contentText,
		resolvedEmbedText,
		resolvedColor,
		embedURLLinkedText,
		profile.footerText,
		iconURL,
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
func (s *DiscordWebhookService) SendNotification(ctx context.Context, webhookURL, botName, contentText, embedType, embedText, embedColor, embedURLLinkedText string) error {
	// embed-typeに応じた処理
	var payload *discord.Payload
	var err error
	switch {
	case embedType == "none":
		payload, err = s.createSimplePayload(botName, contentText)
		if err != nil {
			return fmt.Errorf("ペイロードの作成に失敗しました: %w", err)
		}
	case embedType == "vscode":
		payload, err = s.createVSCodePayload(contentText, embedText, embedColor, embedURLLinkedText)
		if err != nil {
			return err
		}
	case embedType == "open-weather-map":
		payload, err = s.createOpenWeatherMapPayload(contentText, embedText, embedColor, embedURLLinkedText)
		if err != nil {
			return err
		}
	case isGcloudEmbedType(embedType):
		payload, err = s.createGcloudPayload(contentText, embedType, embedText, embedColor, embedURLLinkedText)
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
		footerTextInOpenWeatherEmbed,
		openWeatherIconURL,
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
