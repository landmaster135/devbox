package usecases

import (
	"fmt"
	"strings"
)

const (
	replicationPolicyAutomatic   = "automatic"
	replicationPolicyUserManaged = "user-managed"
)

// Service は Secret Manager 向け gcloud コマンドを生成する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

const (
	discordWebhookEnvVarName = "DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD"
	discordCLIPath           = "$HOME/devbox/pkg/bin/cli/linux_amd64/discord-webhook"
	successEmbedType         = "google-secret-manager-success"
	failureEmbedType         = "google-secret-manager-failed"
)

// DiscordNotificationParams はDiscord通知用コマンド生成に必要な情報を表す。
type DiscordNotificationParams struct {
	Operation  string
	SecretName string
}

// CreateSecretParams は secrets create コマンドのパラメータを表す。
type CreateSecretParams struct {
	SecretName        string
	ReplicationPolicy string
	Locations         string
}

// AddSecretVersionParams は secrets versions add コマンドのパラメータを表す。
type AddSecretVersionParams struct {
	SecretName  string
	SecretValue string
}

// CreateAndAddSecretVersionParams は secrets create と versions add をまとめて実行する際のパラメータを表す。
type CreateAndAddSecretVersionParams struct {
	SecretName        string
	SecretValue       string
	ReplicationPolicy string
	Locations         string
}

// AccessSecretVersionParams は secrets versions access コマンドのパラメータを表す。
type AccessSecretVersionParams struct {
	SecretName string
	Version    string
}

// UpdateSecretLabelsParams は secrets update --update-labels のパラメータを表す。
type UpdateSecretLabelsParams struct {
	SecretName string
	Labels     string
}

// UpdateSecretVersionAliasesParams は secrets update のバージョンエイリアス操作パラメータを表す。
type UpdateSecretVersionAliasesParams struct {
	SecretName  string
	AliasOption string
}

// BuildCreateSecretCommand は gcloud secrets create コマンドを生成する。
func (s *Service) BuildCreateSecretCommand(params CreateSecretParams) (string, error) {
	name := strings.TrimSpace(params.SecretName)
	if name == "" {
		return "", fmt.Errorf("secret-name は必須です")
	}

	policy := strings.TrimSpace(params.ReplicationPolicy)
	if policy == "" {
		policy = replicationPolicyAutomatic
	}

	if policy != replicationPolicyAutomatic && policy != replicationPolicyUserManaged {
		return "", fmt.Errorf("replication-policy には automatic または user-managed を指定してください")
	}

	command := fmt.Sprintf("gcloud secrets create %s --replication-policy=%s", shellQuote(name), shellQuote(policy))

	if policy == replicationPolicyUserManaged {
		locations := strings.TrimSpace(params.Locations)
		if locations == "" {
			return "", fmt.Errorf("replication-policy が user-managed の場合、locations は必須です")
		}
		command = fmt.Sprintf("%s --locations=%s", command, shellQuote(locations))
	}

	return command, nil
}

// BuildAddSecretVersionCommand は gcloud secrets versions add コマンドを生成する。
func (s *Service) BuildAddSecretVersionCommand(params AddSecretVersionParams) (string, error) {
	name := strings.TrimSpace(params.SecretName)
	if name == "" {
		return "", fmt.Errorf("secret-name は必須です")
	}
	if params.SecretValue == "" {
		return "", fmt.Errorf("secret-value は必須です")
	}

	command := fmt.Sprintf("echo -n %s | gcloud secrets versions add %s --data-file=-", shellQuote(params.SecretValue), shellQuote(name))
	return command, nil
}

// BuildCreateAndAddSecretVersionCommand は secrets create と versions add を連結したコマンドを生成する。
func (s *Service) BuildCreateAndAddSecretVersionCommand(params CreateAndAddSecretVersionParams) (string, error) {
	createCmd, err := s.BuildCreateSecretCommand(CreateSecretParams{
		SecretName:        params.SecretName,
		ReplicationPolicy: params.ReplicationPolicy,
		Locations:         params.Locations,
	})
	if err != nil {
		return "", err
	}

	addCmd, err := s.BuildAddSecretVersionCommand(AddSecretVersionParams{
		SecretName:  params.SecretName,
		SecretValue: params.SecretValue,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s && %s", createCmd, addCmd), nil
}

// BuildAccessSecretVersionCommand は gcloud secrets versions access コマンドを生成する。
func (s *Service) BuildAccessSecretVersionCommand(params AccessSecretVersionParams) (string, error) {
	name := strings.TrimSpace(params.SecretName)
	if name == "" {
		return "", fmt.Errorf("secret-name は必須です")
	}

	version := strings.TrimSpace(params.Version)
	if version == "" {
		version = "latest"
	}

	command := fmt.Sprintf("gcloud secrets versions access %s --secret=%s", shellQuote(version), shellQuote(name))
	return command, nil
}

// BuildUpdateSecretLabelsCommand は gcloud secrets update --update-labels コマンドを生成する。
func (s *Service) BuildUpdateSecretLabelsCommand(params UpdateSecretLabelsParams) (string, error) {
	name := strings.TrimSpace(params.SecretName)
	if name == "" {
		return "", fmt.Errorf("secret-name は必須です")
	}

	labels := strings.TrimSpace(params.Labels)
	if labels == "" {
		return "", fmt.Errorf("labels は必須です")
	}

	command := fmt.Sprintf("gcloud secrets update %s --update-labels=%s", shellQuote(name), shellQuote(labels))
	return command, nil
}

// BuildUpdateSecretVersionAliasesCommand は gcloud secrets update のバージョンエイリアス関連コマンドを生成する。
func (s *Service) BuildUpdateSecretVersionAliasesCommand(params UpdateSecretVersionAliasesParams) (string, error) {
	name := strings.TrimSpace(params.SecretName)
	if name == "" {
		return "", fmt.Errorf("secret-name は必須です")
	}

	option := strings.TrimSpace(params.AliasOption)
	if err := validateAliasOption(option); err != nil {
		return "", err
	}

	command := fmt.Sprintf("gcloud secrets update %s %s", shellQuote(name), option)
	return command, nil
}

// BuildNotificationWrappedCommand は開始/成否通知を含んだシェルコマンドを生成する。
func (s *Service) BuildNotificationWrappedCommand(params DiscordNotificationParams, gcloudCommand string) (string, bool) {
	template, ok := notificationTemplates[params.Operation]
	if !ok {
		return "", false
	}

	var lines []string
	if template.startContent != "" {
		lines = append(lines, buildSimpleNotificationCommand(template.startContent))
	}

	successCommand := ""
	if template.successContent != "" {
		successCommand = buildEmbedNotificationCommand(template.successContent, template.successEmbedText, successEmbedType)
	}
	failureCommand := ""
	if template.failureContent != "" {
		failureCommand = buildEmbedNotificationCommand(template.failureContent, template.failureEmbedText, failureEmbedType)
	}

	lines = append(lines, fmt.Sprintf("if %s; then", gcloudCommand))
	if successCommand != "" {
		lines = append(lines, indentCommand(successCommand, "  "))
	}
	lines = append(lines, "else")
	if failureCommand != "" {
		lines = append(lines, indentCommand(failureCommand, "  "))
	}
	lines = append(lines, "fi")

	script := strings.Join(lines, "\n")
	return script, true
}

// PrintHighlightedCommand は生成したコマンドを見やすい形式で出力する。
func (s *Service) PrintHighlightedCommand(command string) {
	fmt.Println()
	fmt.Println("==============================")
	fmt.Println("生成された gcloud コマンド")
	fmt.Println("==============================")
	fmt.Println(command)
	fmt.Println("==============================")
}

// PrintNotificationScript は通知付きシェルコマンドを整形して出力する。
func (s *Service) PrintNotificationScript(script string) {
	if strings.TrimSpace(script) == "" {
		return
	}

	fmt.Println()
	fmt.Println("==============================")
	fmt.Println("通知付きシェルコマンド")
	fmt.Println("==============================")
	fmt.Println(script)
	fmt.Println("==============================")
}

func validateAliasOption(option string) error {
	if option == "" {
		return fmt.Errorf("alias-option は必須です")
	}
	if option == "--clear-version-aliases" {
		return nil
	}
	if strings.HasPrefix(option, "--remove-version-aliases=") {
		return nil
	}
	if strings.HasPrefix(option, "--update-version-aliases=") {
		return nil
	}
	return fmt.Errorf("alias-option には --clear-version-aliases もしくは --remove-version-aliases= / --update-version-aliases= の形式を指定してください")
}

func shellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
}

func buildSimpleNotificationCommand(content string) string {
	return buildDiscordWebhookCommand(content, "none", "")
}

func buildEmbedNotificationCommand(content, embedText, embedType string) string {
	return buildDiscordWebhookCommand(content, embedType, embedText)
}

func buildDiscordWebhookCommand(content, embedType, embedText string) string {
	var lines []string
	lines = append(lines, fmt.Sprintf("%s \\", discordCLIPath))
	lines = append(lines, fmt.Sprintf("  -webhook-url \"$%s\" \\", discordWebhookEnvVarName))
	lines = append(lines, fmt.Sprintf("  -content-text %s \\", shellQuote(content)))
	embedLine := fmt.Sprintf("  -embed-type %s", shelled(embedType))
	if embedText != "" {
		lines = append(lines, embedLine+" \\")
		lines = append(lines, fmt.Sprintf("  -embed-text %s", shellQuote(embedText)))
	} else {
		lines = append(lines, embedLine)
	}
	return strings.Join(lines, "\n")
}
func indentCommand(command, indent string) string {
	if command == "" {
		return ""
	}
	parts := strings.Split(command, "\n")
	for i, part := range parts {
		parts[i] = indent + part
	}
	return strings.Join(parts, "\n")
}

func shelled(value string) string {
	if value == "" {
		return "''"
	}
	return shellQuote(value)
}

type notificationTemplate struct {
	startContent     string
	successContent   string
	successEmbedText string
	failureContent   string
	failureEmbedText string
}

var notificationTemplates = map[string]notificationTemplate{
	"create-secret": {
		startContent:     "シークレットを作るよ！",
		successContent:   "作ったよ！",
		successEmbedText: "シークレットを作ったよ！",
		failureContent:   "失敗…",
		failureEmbedText: "シークレットが作れなかったよ…",
	},
	"add-secret-version": {
		startContent:     "シークレットのバージョンを作るよ！",
		successContent:   "作ったよ！",
		successEmbedText: "シークレットにバージョンを追加したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "シークレットにバージョンを作れなかったよ…",
	},
	"create-and-add-secret-version": {
		startContent:     "シークレットとバージョンを作るよ！",
		successContent:   "作ったよ！",
		successEmbedText: "シークレットと新しいバージョンを追加したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "シークレットとバージョンの作成に失敗したよ…",
	},
	"access-secret-version": {
		startContent:     "シークレットの値を取得するよ！",
		successContent:   "取れた！",
		successEmbedText: "シークレットを取得したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "シークレットを取れなかったよ…",
	},
	"update-secret-labels": {
		startContent:     "シークレットのラベルを更新するよ！",
		successContent:   "更新したよ！",
		successEmbedText: "シークレットのラベルを更新したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "シークレットのラベル更新に失敗したよ…",
	},
	"update-secret-version-aliases": {
		startContent:     "シークレットのエイリアスを更新するよ！",
		successContent:   "更新したよ！",
		successEmbedText: "シークレットのエイリアスを更新したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "シークレットのエイリアス更新に失敗したよ…",
	},
}
