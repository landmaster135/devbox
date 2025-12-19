package usecases

import (
	"fmt"
	"strings"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_container/config"
)

// Service は gcloud コマンド生成を担当する。
type Service struct{}

const (
	discordWebhookEnvVarName    = "DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD"
	discordCLIPath              = "$HOME/devbox/pkg/bin/cli/linux_amd64/discord-webhook"
	runSuccessEmbedType         = "google-cloud-run-success"
	runFailureEmbedType         = "google-cloud-run-failed"
	runFunctionSuccessEmbedType = "google-cloud-run-function-success"
	runFunctionFailureEmbedType = "google-cloud-run-function-failed"
)

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// BuildCommand は設定された operation に応じて gcloud コマンドを生成する。
func (s *Service) BuildCommand(conf *cfg.Config) (string, error) {
	switch conf.Operation {
	case cfg.OperationDeployCloudRunContainer:
		return s.BuildDeployCloudRunContainerCommand(DeployCloudRunContainerParams{
			ServiceName:          conf.ServiceName,
			ProjectID:            conf.ProjectID,
			Region:               conf.Region,
			Timeout:              conf.Timeout,
			RunServiceAccount:    conf.RunServiceAccount,
			AllowUnauthenticated: conf.AllowUnauthenticated,
		})
	case cfg.OperationUpdateCloudRunContainerEnv:
		return s.BuildUpdateCloudRunContainerEnvCommand(UpdateCloudRunContainerEnvParams{
			ServiceName: conf.ServiceName,
			ProjectID:   conf.ProjectID,
			Region:      conf.Region,
			EnvFile:     conf.EnvFile,
		})
	case cfg.OperationDeployCloudRunFunction:
		return s.BuildDeployCloudRunFunctionCommand(DeployCloudRunFunctionParams{
			FunctionName: conf.FunctionName,
			Region:       conf.Region,
			EntryPoint:   conf.EntryPoint,
		})
	case cfg.OperationDeployCloudRunFunctionTriggeredByPubSub:
		return s.BuildDeployCloudRunFunctionTriggeredByPubSubCommand(DeployCloudRunFunctionTriggeredByPubSubParams{
			FunctionName:          conf.FunctionName,
			ProjectID:             conf.ProjectID,
			Region:                conf.Region,
			EntryPoint:            conf.EntryPoint,
			TriggerServiceAccount: conf.TriggerServiceAccount,
			TriggerTopic:          conf.TriggerTopic,
			APIClientID:           conf.APIClientID,
			APIClientSecret:       conf.APIClientSecret,
			APIEndpoint:           conf.APIEndpoint,
		})
	case cfg.OperationUpdateCloudRunFunctionEnv:
		return s.BuildUpdateCloudRunFunctionEnvCommand(UpdateCloudRunFunctionEnvParams{
			ServiceName: conf.ServiceName,
			Region:      conf.Region,
			EnvVars:     conf.EnvVars,
		})
	case cfg.OperationUpdateCloudRunServiceEnv:
		return s.BuildUpdateCloudRunServiceEnvCommand(UpdateCloudRunServiceEnvParams{
			ServiceName: conf.ServiceName,
			ProjectID:   conf.ProjectID,
			Region:      conf.Region,
			EnvFile:     conf.EnvFile,
		})
	case cfg.OperationCreateCloudPubSubTopic:
		return s.BuildCreateCloudPubSubTopicCommand(CreateCloudPubSubTopicParams{
			TopicName:                conf.TopicName,
			MessageRetentionDuration: conf.MessageRetentionDuration,
		})
	case cfg.OperationListCloudPubSubTopics:
		return s.BuildListCloudPubSubTopicsCommand(ListCloudPubSubTopicsParams{TopicName: conf.TopicName})
	case cfg.OperationListCloudPubSubSubscriptions:
		return s.BuildListCloudPubSubSubscriptionsCommand(ListCloudPubSubSubscriptionsParams{
			SubscriptionName: conf.SubscriptionName,
			ShowURI:          conf.ShowURI,
		})
	case cfg.OperationCreateCloudPubSubSubscription:
		return s.BuildCreateCloudPubSubSubscriptionCommand(CreateCloudPubSubSubscriptionParams{
			SubscriptionName:         conf.SubscriptionName,
			TopicName:                conf.TopicName,
			TopicProject:             conf.TopicProject,
			PushEndpoint:             conf.PushEndpoint,
			PushServiceAccount:       conf.PushServiceAccount,
			MessageRetentionDuration: conf.MessageRetentionDuration,
			ExpirationPeriod:         conf.ExpirationPeriod,
			MaxRetryDelay:            conf.MaxRetryDelay,
			MinRetryDelay:            conf.MinRetryDelay,
			AckDeadline:              conf.AckDeadline,
		})
	case cfg.OperationDeleteCloudPubSubSubscriptionsAndTopics:
		return s.BuildDeleteCloudPubSubSubscriptionsAndTopicsCommand(DeleteCloudPubSubSubscriptionsAndTopicsParams{
			SubscriptionNames: conf.SubscriptionNames,
			TopicNames:        conf.TopicNames,
		})
	case cfg.OperationDeleteCloudPubSubSubscriptions:
		return s.BuildDeleteCloudPubSubSubscriptionsCommand(DeleteCloudPubSubSubscriptionsParams{
			SubscriptionNames: conf.SubscriptionNames,
		})
	case cfg.OperationDeleteCloudPubSubTopics:
		return s.BuildDeleteCloudPubSubTopicsCommand(DeleteCloudPubSubTopicsParams{
			TopicNames: conf.TopicNames,
		})
	case cfg.OperationDeleteCloudRunFunction:
		return s.BuildDeleteCloudRunFunctionCommand(DeleteCloudRunFunctionParams{
			ServiceName: conf.ServiceName,
			Region:      conf.Region,
		})
	default:
		return "", fmt.Errorf("未対応のoperationです: %s", conf.Operation)
	}
}

// DeployCloudRunContainerParams は Cloud Run コンテナ デプロイ用コマンドの生成に必要な情報を保持する。
type DeployCloudRunContainerParams struct {
	ServiceName          string
	ProjectID            string
	Region               string
	Timeout              string
	RunServiceAccount    string
	AllowUnauthenticated bool
}

// BuildDeployCloudRunContainerCommand は gcloud run deploy コマンドを生成する。
func (s *Service) BuildDeployCloudRunContainerCommand(params DeployCloudRunContainerParams) (string, error) {
	if params.ServiceName == "" || params.ProjectID == "" || params.Region == "" || params.Timeout == "" {
		return "", fmt.Errorf("必須パラメータが不足しています")
	}

	command := fmt.Sprintf(
		"gcloud run deploy %s --source . --project=%s --region=%s --timeout=%s",
		shellQuote(params.ServiceName),
		shellQuote(params.ProjectID),
		shellQuote(params.Region),
		shellQuote(params.Timeout),
	)

	if params.RunServiceAccount != "" {
		command = fmt.Sprintf("%s --service-account=%s", command, shellQuote(params.RunServiceAccount))
	}

	if params.AllowUnauthenticated {
		command = fmt.Sprintf("%s --allow-unauthenticated", command)
	} else {
		command = fmt.Sprintf("%s --no-allow-unauthenticated", command)
	}

	return command, nil
}

// UpdateCloudRunContainerEnvParams は Cloud Run コンテナの環境変数更新コマンド生成に必要な情報を表す。
type UpdateCloudRunContainerEnvParams struct {
	ServiceName string
	ProjectID   string
	Region      string
	EnvFile     string
}

// BuildUpdateCloudRunContainerEnvCommand は gcloud run deploy による環境変数更新コマンドを生成する。
func (s *Service) BuildUpdateCloudRunContainerEnvCommand(params UpdateCloudRunContainerEnvParams) (string, error) {
	if params.ServiceName == "" || params.ProjectID == "" || params.Region == "" || params.EnvFile == "" {
		return "", fmt.Errorf("必須パラメータが不足しています")
	}

	imagePath := fmt.Sprintf("gcr.io/%s/%s", params.ProjectID, params.ServiceName)

	command := fmt.Sprintf(
		"gcloud run deploy %s --image=%s --region=%s --env-vars-file=%s",
		shellQuote(params.ServiceName),
		shellQuote(imagePath),
		shellQuote(params.Region),
		shellQuote(params.EnvFile),
	)

	return command, nil
}

// DeployCloudRunFunctionParams は Cloud Functions (Gen2) デプロイコマンドのパラメータを表す。
type DeployCloudRunFunctionParams struct {
	FunctionName string
	Region       string
	EntryPoint   string
}

// BuildDeployCloudRunFunctionCommand は gcloud functions deploy (Gen2, HTTP トリガー) コマンドを生成する。
func (s *Service) BuildDeployCloudRunFunctionCommand(params DeployCloudRunFunctionParams) (string, error) {
	if params.FunctionName == "" || params.Region == "" || params.EntryPoint == "" {
		return "", fmt.Errorf("必須パラメータが不足しています")
	}

	command := fmt.Sprintf(
		"gcloud functions deploy %s --gen2 --runtime=go122 --region=%s --source=. --entry-point=%s --trigger-http --allow-unauthenticated --timeout=180s",
		shellQuote(params.FunctionName),
		shellQuote(params.Region),
		shellQuote(params.EntryPoint),
	)

	return command, nil
}

// DeployCloudRunFunctionTriggeredByPubSubParams は Pub/Sub トリガー用 Cloud Functions デプロイのパラメータを表す。
type DeployCloudRunFunctionTriggeredByPubSubParams struct {
	FunctionName          string
	ProjectID             string
	Region                string
	EntryPoint            string
	TriggerServiceAccount string
	TriggerTopic          string
	APIClientID           string
	APIClientSecret       string
	APIEndpoint           string
}

// BuildDeployCloudRunFunctionTriggeredByPubSubCommand は Pub/Sub トリガー向けデプロイコマンドを生成する。
func (s *Service) BuildDeployCloudRunFunctionTriggeredByPubSubCommand(params DeployCloudRunFunctionTriggeredByPubSubParams) (string, error) {
	if params.FunctionName == "" || params.ProjectID == "" || params.Region == "" || params.EntryPoint == "" || params.TriggerServiceAccount == "" || params.TriggerTopic == "" {
		return "", fmt.Errorf("必須パラメータが不足しています")
	}

	command := fmt.Sprintf(
		"gcloud functions deploy %s --gen2 --runtime=go123 --project=%s --region=%s --source=. --entry-point=%s --trigger-service-account=%s --trigger-topic=%s --allow-unauthenticated --timeout=180s",
		shellQuote(params.FunctionName),
		shellQuote(params.ProjectID),
		shellQuote(params.Region),
		shellQuote(params.EntryPoint),
		shellQuote(params.TriggerServiceAccount),
		shellQuote(params.TriggerTopic),
	)

	envVars := buildOptionalEnvVars(params.APIClientID, params.APIClientSecret, params.APIEndpoint)
	if envVars != "" {
		command = fmt.Sprintf("%s --set-env-vars=%s", command, shellQuote(envVars))
	}

	return command, nil
}

func buildOptionalEnvVars(apiClientID, apiClientSecret, apiEndpoint string) string {
	var parts []string
	if strings.TrimSpace(apiClientID) != "" {
		parts = append(parts, fmt.Sprintf("SCRIPT_MANAGER_API_CLIENT_ID=%s", apiClientID))
	}
	if strings.TrimSpace(apiClientSecret) != "" {
		parts = append(parts, fmt.Sprintf("SCRIPT_MANAGER_API_CLIENT_SECRET=%s", apiClientSecret))
	}
	if strings.TrimSpace(apiEndpoint) != "" {
		parts = append(parts, fmt.Sprintf("SCRIPT_MANAGER_API_ENDPOINT=%s", apiEndpoint))
	}
	return strings.Join(parts, ",")
}

// UpdateCloudRunFunctionEnvParams は Cloud Functions (Gen2) の環境変数更新パラメータを表す。
type UpdateCloudRunFunctionEnvParams struct {
	ServiceName string
	Region      string
	EnvVars     string
}

// BuildUpdateCloudRunFunctionEnvCommand は gcloud run services update の環境変数更新コマンドを生成する。
func (s *Service) BuildUpdateCloudRunFunctionEnvCommand(params UpdateCloudRunFunctionEnvParams) (string, error) {
	if params.ServiceName == "" || params.Region == "" || params.EnvVars == "" {
		return "", fmt.Errorf("必須パラメータが不足しています")
	}

	return fmt.Sprintf(
		"gcloud run services update %s --region=%s --update-env-vars=%s",
		shellQuote(params.ServiceName),
		shellQuote(params.Region),
		shellQuote(params.EnvVars),
	), nil
}

// UpdateCloudRunServiceEnvParams は Cloud Run サービス環境変数更新コマンドのパラメータを表す。
type UpdateCloudRunServiceEnvParams struct {
	ServiceName string
	ProjectID   string
	Region      string
	EnvFile     string
}

// BuildUpdateCloudRunServiceEnvCommand は gcloud run services update (ファイル指定) コマンドを生成する。
func (s *Service) BuildUpdateCloudRunServiceEnvCommand(params UpdateCloudRunServiceEnvParams) (string, error) {
	if params.ServiceName == "" || params.ProjectID == "" || params.Region == "" || params.EnvFile == "" {
		return "", fmt.Errorf("必須パラメータが不足しています")
	}

	return fmt.Sprintf(
		"gcloud run services update %s --project=%s --region=%s --env-vars-file=%s",
		shellQuote(params.ServiceName),
		shellQuote(params.ProjectID),
		shellQuote(params.Region),
		shellQuote(params.EnvFile),
	), nil
}

// CreateCloudPubSubTopicParams は Pub/Sub トピック作成コマンドのパラメータを表す。
type CreateCloudPubSubTopicParams struct {
	TopicName                string
	MessageRetentionDuration string
}

// BuildCreateCloudPubSubTopicCommand は gcloud pubsub topics create コマンドを生成する。
func (s *Service) BuildCreateCloudPubSubTopicCommand(params CreateCloudPubSubTopicParams) (string, error) {
	if params.TopicName == "" || params.MessageRetentionDuration == "" {
		return "", fmt.Errorf("必須パラメータが不足しています")
	}

	return fmt.Sprintf(
		"gcloud pubsub topics create %s --message-retention-duration=%s",
		shellQuote(params.TopicName),
		shellQuote(params.MessageRetentionDuration),
	), nil
}

// ListCloudPubSubTopicsParams は Pub/Sub トピック一覧コマンドのパラメータを表す。
type ListCloudPubSubTopicsParams struct {
	TopicName string
}

// BuildListCloudPubSubTopicsCommand は gcloud pubsub topics list コマンドを生成する。
func (s *Service) BuildListCloudPubSubTopicsCommand(params ListCloudPubSubTopicsParams) (string, error) {
	if params.TopicName == "" {
		return "", fmt.Errorf("必須パラメータが不足しています")
	}

	filter := fmt.Sprintf("name.scope(topic):'%s'", params.TopicName)
	return fmt.Sprintf("gcloud pubsub topics list --filter=%s", doubleQuote(filter)), nil
}

// ListCloudPubSubSubscriptionsParams は Pub/Sub サブスクリプション一覧コマンドのパラメータを表す。
type ListCloudPubSubSubscriptionsParams struct {
	SubscriptionName string
	ShowURI          bool
}

// BuildListCloudPubSubSubscriptionsCommand は gcloud pubsub subscriptions list コマンドを生成する。
func (s *Service) BuildListCloudPubSubSubscriptionsCommand(params ListCloudPubSubSubscriptionsParams) (string, error) {
	builder := strings.Builder{}
	builder.WriteString("gcloud pubsub subscriptions list")

	if params.SubscriptionName != "" {
		filter := fmt.Sprintf("name.scope(subscription):'%s'", params.SubscriptionName)
		builder.WriteString(" --filter=")
		builder.WriteString(doubleQuote(filter))
	}

	if params.ShowURI {
		builder.WriteString(" --uri")
	}

	return builder.String(), nil
}

// CreateCloudPubSubSubscriptionParams は Pub/Sub サブスクリプション作成コマンドのパラメータを表す。
type CreateCloudPubSubSubscriptionParams struct {
	SubscriptionName         string
	TopicName                string
	TopicProject             string
	PushEndpoint             string
	PushServiceAccount       string
	MessageRetentionDuration string
	ExpirationPeriod         string
	MaxRetryDelay            string
	MinRetryDelay            string
	AckDeadline              string
}

// BuildCreateCloudPubSubSubscriptionCommand は gcloud pubsub subscriptions create コマンドを生成する。
func (s *Service) BuildCreateCloudPubSubSubscriptionCommand(params CreateCloudPubSubSubscriptionParams) (string, error) {
	if params.SubscriptionName == "" || params.TopicName == "" || params.TopicProject == "" || params.PushServiceAccount == "" {
		return "", fmt.Errorf("必須パラメータが不足しています")
	}

	retention := params.MessageRetentionDuration
	if retention == "" {
		return "", fmt.Errorf("message-retention-duration パラメータは必須です")
	}
	expiration := params.ExpirationPeriod
	if expiration == "" {
		return "", fmt.Errorf("expiration-period パラメータは必須です")
	}
	maxRetry := params.MaxRetryDelay
	if maxRetry == "" {
		return "", fmt.Errorf("max-retry-delay パラメータは必須です")
	}
	minRetry := params.MinRetryDelay
	if minRetry == "" {
		return "", fmt.Errorf("min-retry-delay パラメータは必須です")
	}
	ackDeadline := params.AckDeadline
	if ackDeadline == "" {
		return "", fmt.Errorf("ack-deadline パラメータは必須です")
	}

	endpoint := params.PushEndpoint
	if strings.TrimSpace(endpoint) == "" {
		endpoint = fmt.Sprintf("https://%s.appspot.com/%s", params.TopicProject, params.SubscriptionName)
	}

	return fmt.Sprintf(
		"gcloud pubsub subscriptions create %s --topic=%s --topic-project=%s --message-retention-duration=%s --push-auth-service-account=%s --push-endpoint=%s --expiration-period=%s --max-retry-delay=%s --min-retry-delay=%s --ack-deadline=%s",
		shellQuote(params.SubscriptionName),
		shellQuote(params.TopicName),
		shellQuote(params.TopicProject),
		shellQuote(retention),
		shellQuote(params.PushServiceAccount),
		shellQuote(endpoint),
		shellQuote(expiration),
		shellQuote(maxRetry),
		shellQuote(minRetry),
		shellQuote(ackDeadline),
	), nil
}

// DeleteCloudPubSubSubscriptionsAndTopicsParams はサブスクリプションとトピックを同時に削除する際のパラメータを表す。
type DeleteCloudPubSubSubscriptionsAndTopicsParams struct {
	SubscriptionNames []string
	TopicNames        []string
}

// BuildDeleteCloudPubSubSubscriptionsAndTopicsCommand は削除用コマンドを生成する。
func (s *Service) BuildDeleteCloudPubSubSubscriptionsAndTopicsCommand(params DeleteCloudPubSubSubscriptionsAndTopicsParams) (string, error) {
	if len(params.SubscriptionNames) == 0 && len(params.TopicNames) == 0 {
		return "", fmt.Errorf("削除対象が指定されていません")
	}

	var commands []string

	if len(params.SubscriptionNames) > 0 {
		cmd, err := s.BuildDeleteCloudPubSubSubscriptionsCommand(DeleteCloudPubSubSubscriptionsParams{SubscriptionNames: params.SubscriptionNames})
		if err != nil {
			return "", err
		}
		commands = append(commands, cmd)
	}

	if len(params.TopicNames) > 0 {
		cmd, err := s.BuildDeleteCloudPubSubTopicsCommand(DeleteCloudPubSubTopicsParams{TopicNames: params.TopicNames})
		if err != nil {
			return "", err
		}
		commands = append(commands, cmd)
	}

	return strings.Join(commands, " && "), nil
}

// DeleteCloudPubSubSubscriptionsParams はサブスクリプション削除コマンドのパラメータを表す。
type DeleteCloudPubSubSubscriptionsParams struct {
	SubscriptionNames []string
}

// BuildDeleteCloudPubSubSubscriptionsCommand は gcloud pubsub subscriptions delete コマンド群を生成する。
func (s *Service) BuildDeleteCloudPubSubSubscriptionsCommand(params DeleteCloudPubSubSubscriptionsParams) (string, error) {
	if len(params.SubscriptionNames) == 0 {
		return "", fmt.Errorf("subscription-names が指定されていません")
	}

	commands := make([]string, 0, len(params.SubscriptionNames))
	for _, name := range params.SubscriptionNames {
		if strings.TrimSpace(name) == "" {
			continue
		}
		commands = append(commands, fmt.Sprintf("gcloud pubsub subscriptions delete %s", shellQuote(name)))
	}

	if len(commands) == 0 {
		return "", fmt.Errorf("有効な subscription-names が指定されていません")
	}

	return strings.Join(commands, " && "), nil
}

// DeleteCloudPubSubTopicsParams はトピック削除コマンドのパラメータを表す。
type DeleteCloudPubSubTopicsParams struct {
	TopicNames []string
}

// BuildDeleteCloudPubSubTopicsCommand は gcloud pubsub topics delete コマンド群を生成する。
func (s *Service) BuildDeleteCloudPubSubTopicsCommand(params DeleteCloudPubSubTopicsParams) (string, error) {
	if len(params.TopicNames) == 0 {
		return "", fmt.Errorf("topic-names が指定されていません")
	}

	commands := make([]string, 0, len(params.TopicNames))
	for _, name := range params.TopicNames {
		if strings.TrimSpace(name) == "" {
			continue
		}
		commands = append(commands, fmt.Sprintf("gcloud pubsub topics delete %s", shellQuote(name)))
	}

	if len(commands) == 0 {
		return "", fmt.Errorf("有効な topic-names が指定されていません")
	}

	return strings.Join(commands, " && "), nil
}

// DeleteCloudRunFunctionParams は Cloud Functions (Gen2) の削除コマンドを表す。
type DeleteCloudRunFunctionParams struct {
	ServiceName string
	Region      string
}

// BuildDeleteCloudRunFunctionCommand は gcloud run services delete コマンドを生成する。
func (s *Service) BuildDeleteCloudRunFunctionCommand(params DeleteCloudRunFunctionParams) (string, error) {
	if params.ServiceName == "" || params.Region == "" {
		return "", fmt.Errorf("必須パラメータが不足しています")
	}

	return fmt.Sprintf(
		"gcloud run services delete %s --region=%s",
		shellQuote(params.ServiceName),
		shellQuote(params.Region),
	), nil
}

// PrintHighlightedCommand は生成したコマンドを枠で囲んで出力する。
func (s *Service) PrintHighlightedCommand(command string) {
	fmt.Println()
	fmt.Println("==============================")
	fmt.Println("生成された gcloud コマンド")
	fmt.Println("==============================")
	fmt.Println(command)
	fmt.Println("==============================")
}

// BuildNotificationWrappedCommand は Discord 通知付きのシェルスクリプトを生成する。
func (s *Service) BuildNotificationWrappedCommand(operation string, gcloudCommand string) (string, bool) {
	template, ok := notificationTemplates[operation]
	if !ok {
		return "", false
	}

	var lines []string
	if template.startContent != "" {
		lines = append(lines, buildSimpleNotificationCommand(template.startContent))
	}

	successCommand := ""
	if template.successContent != "" {
		successCommand = buildEmbedNotificationCommand(template.successContent, template.successEmbedText, template.successEmbedType)
	}
	failureCommand := ""
	if template.failureContent != "" {
		failureCommand = buildEmbedNotificationCommand(template.failureContent, template.failureEmbedText, template.failureEmbedType)
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

	return strings.Join(lines, "\n"), true
}

// PrintNotificationScript は通知コマンドを枠付きで表示する。
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

func shellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\\''")
	return "'" + escaped + "'"
}

func doubleQuote(value string) string {
	escaped := strings.ReplaceAll(value, "\"", "\\\"")
	return "\"" + escaped + "\""
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
	embedLine := fmt.Sprintf("  -embed-type %s", shellQuote(embedType))
	if strings.TrimSpace(embedText) != "" {
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

type notificationTemplate struct {
	startContent     string
	successContent   string
	successEmbedText string
	successEmbedType string
	failureContent   string
	failureEmbedText string
	failureEmbedType string
}

var notificationTemplates = map[string]notificationTemplate{
	cfg.OperationDeployCloudRunContainer: {
		startContent:     "コンテナをデプロイするよ！",
		successContent:   "デプロイしたよ！",
		successEmbedText: "コンテナをデプロイしたよ！",
		successEmbedType: runSuccessEmbedType,
		failureContent:   "失敗…",
		failureEmbedText: "コンテナをデプロイできなかったよ…",
		failureEmbedType: runFailureEmbedType,
	},
	cfg.OperationUpdateCloudRunContainerEnv: {
		startContent:     "コンテナの環境変数を更新するよ！",
		successContent:   "更新したよ！",
		successEmbedText: "コンテナの環境変数を更新したよ！",
		successEmbedType: runSuccessEmbedType,
		failureContent:   "失敗…",
		failureEmbedText: "コンテナの環境変数を更新できなかったよ…",
		failureEmbedType: runFailureEmbedType,
	},
	cfg.OperationDeployCloudRunFunction: {
		startContent:     "関数をデプロイするよ！",
		successContent:   "デプロイしたよ！",
		successEmbedText: "関数をデプロイしたよ！",
		successEmbedType: runFunctionSuccessEmbedType,
		failureContent:   "失敗…",
		failureEmbedText: "関数をデプロイできなかったよ…",
		failureEmbedType: runFunctionFailureEmbedType,
	},
	cfg.OperationDeployCloudRunFunctionTriggeredByPubSub: {
		startContent:     "関数をデプロイするよ！",
		successContent:   "デプロイしたよ！",
		successEmbedText: "関数をデプロイしたよ！",
		successEmbedType: runFunctionSuccessEmbedType,
		failureContent:   "失敗…",
		failureEmbedText: "関数をデプロイできなかったよ…",
		failureEmbedType: runFunctionFailureEmbedType,
	},
	cfg.OperationUpdateCloudRunFunctionEnv: {
		startContent:     "関数の環境変数を更新するよ！",
		successContent:   "更新したよ！",
		successEmbedText: "関数の環境変数を更新したよ！",
		successEmbedType: runFunctionSuccessEmbedType,
		failureContent:   "失敗…",
		failureEmbedText: "関数の環境変数を更新できなかったよ…",
		failureEmbedType: runFunctionFailureEmbedType,
	},
	cfg.OperationUpdateCloudRunServiceEnv: {
		startContent:     "サービスの環境変数を更新するよ！",
		successContent:   "更新したよ！",
		successEmbedText: "サービスの環境変数を更新したよ！",
		successEmbedType: runSuccessEmbedType,
		failureContent:   "失敗…",
		failureEmbedText: "サービスの環境変数を更新できなかったよ…",
		failureEmbedType: runFailureEmbedType,
	},
}
