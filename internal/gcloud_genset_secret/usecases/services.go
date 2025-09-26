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
	if template.start != "" {
		lines = append(lines, template.start)
	}

	if template.failure != "" {
		lines = append(lines, fmt.Sprintf("if %s; then", gcloudCommand))
		if template.success != "" {
			lines = append(lines, fmt.Sprintf("  %s", template.success))
		}
		lines = append(lines, "else")
		lines = append(lines, fmt.Sprintf("  %s", template.failure))
		lines = append(lines, "fi")
	} else {
		lines = append(lines, gcloudCommand)
		if template.success != "" {
			lines = append(lines, template.success)
		}
	}

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

type notificationTemplate struct {
	start   string
	success string
	failure string
}

var notificationTemplates = map[string]notificationTemplate{
	"create-secret": {
		start:   "send_discord_notification \"シークレットを作るよ！\"",
		success: "send_discord_notification_about_gsm \"作ったよ！\" \"シークレットを作ったよ！\" \"green\"",
		failure: "send_discord_notification_about_gsm \"失敗…\" \"シークレットが作れなかったよ…\" \"red\"",
	},
	"add-secret-version": {
		start:   "send_discord_notification \"シークレットのバージョンを作るよ！\"",
		success: "send_discord_notification_about_gsm \"作ったよ！\" \"シークレットにバージョンを追加したよ！\" \"green\"",
		failure: "send_discord_notification_about_gsm \"失敗…\" \"シークレットにバージョンを作れなかったよ…\" \"red\"",
	},
	"create-and-add-secret-version": {
		start:   "send_discord_notification \"シークレットとバージョンを作るよ！\"",
		success: "send_discord_notification_about_gsm \"作ったよ！\" \"シークレットと新しいバージョンを追加したよ！\" \"green\"",
		failure: "send_discord_notification_about_gsm \"失敗…\" \"シークレットとバージョンの作成に失敗したよ…\" \"red\"",
	},
	"access-secret-version": {
		start:   "send_discord_notification \"シークレットの値を取得するよ！\"",
		success: "send_discord_notification_about_gsm \"取れた！\" \"シークレットを取得したよ！\" \"green\"",
		failure: "send_discord_notification_about_gsm \"失敗…\" \"シークレットを取れなかったよ…\" \"red\"",
	},
	"update-secret-labels": {
		start:   "send_discord_notification \"シークレットのラベルを更新するよ！\"",
		success: "send_discord_notification_about_gsm \"更新したよ！\" \"シークレットのラベルを更新したよ！\" \"green\"",
		failure: "send_discord_notification_about_gsm \"失敗…\" \"シークレットのラベル更新に失敗したよ…\" \"red\"",
	},
	"update-secret-version-aliases": {
		start:   "send_discord_notification \"シークレットのエイリアスを更新するよ！\"",
		success: "send_discord_notification_about_gsm \"更新したよ！\" \"シークレットのエイリアスを更新したよ！\" \"green\"",
		failure: "send_discord_notification_about_gsm \"失敗…\" \"シークレットのエイリアス更新に失敗したよ…\" \"red\"",
	},
}
