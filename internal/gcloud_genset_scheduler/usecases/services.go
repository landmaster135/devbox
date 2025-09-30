package usecases

import (
	"encoding/json"
	"fmt"
	"strings"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_scheduler/config"
)

const (
	defaultLocation          = "us-central1"
	defaultSchedule          = "0 4 * * 0-6"
	defaultStopSchedule      = "0 7 * * 0-6"
	defaultDescriptionPubSub = "Trigger Cloud Functions to start Cloud SQL instance."
	defaultDescriptionHTTP   = "Trigger Cloud Run Container."
	defaultTimeZone          = "Asia/Tokyo"
)

// Service は Cloud Scheduler 向け gcloud コマンドを生成する役割を持つ。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

const (
	discordWebhookEnvVarName = "DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD"
	discordCLIPath           = "$HOME/devbox/pkg/bin/cli/linux_amd64/discord-webhook"
	successEmbedType         = "google-cloud-scheduler-success"
	failureEmbedType         = "google-cloud-scheduler-failed"
)

// DiscordNotificationParams は Discord 通知生成に必要な情報を表す。
type DiscordNotificationParams struct {
	Operation string
}

// CreatePubSubJobParams は Pub/Sub ジョブ生成時のパラメータを表す。
type CreatePubSubJobParams struct {
	JobName     string
	ProjectID   string
	Location    string
	Schedule    string
	Description string
	TimeZone    string
	PubsubTopic string
	MessageBody string
}

// CreateHTTPJobParams は HTTP ジョブ生成時のパラメータを表す。
type CreateHTTPJobParams struct {
	JobName                 string
	ProjectID               string
	Location                string
	Schedule                string
	Description             string
	TimeZone                string
	HTTPMethod              string
	ServiceURL              string
	MessageBody             string
	Headers                 string
	OIDCServiceAccountEmail string
}

// CreateCloudSQLJobParams は Cloud SQL 向けジョブ生成に必要な値を表す。
type CreateCloudSQLJobParams struct {
	JobName           string
	ProjectID         string
	Location          string
	Schedule          string
	Description       string
	TimeZone          string
	PubsubTopic       string
	DBInstanceID      string
	Action            string
	DiscordWebhookURL string
	CloudSQLIconURL   string
}

// CreateStartStopCloudSQLJobParams は Cloud SQL 起動/停止ジョブ用の値を表す。
type CreateStartStopCloudSQLJobParams struct {
	JobName           string
	ProjectID         string
	Location          string
	Schedule          string
	TimeZone          string
	Description       string
	PubsubTopic       string
	DBInstanceID      string
	DiscordWebhookURL string
	CloudSQLIconURL   string
}

// ListJobsParams はジョブ一覧取得時のパラメータを表す。
type ListJobsParams struct {
	Location string
	Limit    string
}

// UpdateHTTPJobParams は HTTP ジョブ更新時のパラメータを表す。
type UpdateHTTPJobParams struct {
	JobName     string
	Schedule    string
	MessageBody string
	Headers     string
}

// UpdatePubSubJobParams は Pub/Sub ジョブ更新時のパラメータを表す。
type UpdatePubSubJobParams struct {
	JobName     string
	Schedule    string
	MessageBody string
}

// JobControlParams は pause/resume/delete/run で共通するパラメータを表す。
type JobControlParams struct {
	JobName  string
	Location string
}

// BuildCreatePubSubJobCommand は Pub/Sub ジョブ作成コマンドを生成する。
func (s *Service) BuildCreatePubSubJobCommand(params CreatePubSubJobParams) (string, error) {
	jobName := strings.TrimSpace(params.JobName)
	if jobName == "" {
		return "", fmt.Errorf("jobName は必須です")
	}
	projectID := strings.TrimSpace(params.ProjectID)
	if projectID == "" {
		return "", fmt.Errorf("projectID は必須です")
	}
	topic := strings.TrimSpace(params.PubsubTopic)
	if topic == "" {
		return "", fmt.Errorf("pubsubTopic は必須です")
	}

	location := defaultIfEmpty(strings.TrimSpace(params.Location), defaultLocation)
	schedule := defaultIfEmpty(strings.TrimSpace(params.Schedule), defaultSchedule)
	description := defaultIfEmpty(strings.TrimSpace(params.Description), defaultDescriptionPubSub)
	timeZone := defaultIfEmpty(strings.TrimSpace(params.TimeZone), defaultTimeZone)
	messageBody := sanitizeMessageBody(params.MessageBody)

	flags := []string{
		formatFlag("--schedule", schedule),
		formatFlag("--description", description),
		formatFlag("--project", projectID),
		formatFlag("--location", location),
		formatFlag("--time-zone", timeZone),
		formatFlag("--topic", topic),
	}
	if messageBody != "" {
		flags = append(flags, formatFlag("--message-body", messageBody))
	}

	base := []string{"gcloud", "scheduler", "jobs", "create", "pubsub", shellQuote(jobName)}
	return buildMultilineCommand(base, flags), nil
}

// BuildCreateHTTPJobCommand は HTTP ジョブ作成コマンドを生成する。
func (s *Service) BuildCreateHTTPJobCommand(params CreateHTTPJobParams) (string, error) {
	jobName := strings.TrimSpace(params.JobName)
	if jobName == "" {
		return "", fmt.Errorf("jobName は必須です")
	}
	projectID := strings.TrimSpace(params.ProjectID)
	if projectID == "" {
		return "", fmt.Errorf("projectID は必須です")
	}
	method := strings.ToUpper(strings.TrimSpace(params.HTTPMethod))
	if method == "" {
		return "", fmt.Errorf("httpMethod は必須です")
	}
	serviceURL := strings.TrimSpace(params.ServiceURL)
	if serviceURL == "" {
		return "", fmt.Errorf("serviceURL は必須です")
	}

	location := defaultIfEmpty(strings.TrimSpace(params.Location), defaultLocation)
	schedule := defaultIfEmpty(strings.TrimSpace(params.Schedule), defaultSchedule)
	description := defaultIfEmpty(strings.TrimSpace(params.Description), defaultDescriptionHTTP)
	timeZone := defaultIfEmpty(strings.TrimSpace(params.TimeZone), defaultTimeZone)
	messageBody := sanitizeMessageBody(params.MessageBody)
	headers := strings.TrimSpace(params.Headers)
	oidc := strings.TrimSpace(params.OIDCServiceAccountEmail)

	flags := []string{
		formatFlag("--schedule", schedule),
		formatFlag("--description", description),
		formatFlag("--project", projectID),
		formatFlag("--location", location),
		formatFlag("--time-zone", timeZone),
		formatFlag("--http-method", method),
		formatFlag("--uri", serviceURL),
	}

	if headers != "" {
		flags = append(flags, formatFlag("--headers", headers))
	}
	if messageBody != "" {
		flags = append(flags, formatFlag("--message-body", messageBody))
	}
	if oidc != "" {
		flags = append(flags, formatFlag("--oidc-service-account-email", oidc))
	}

	base := []string{"gcloud", "scheduler", "jobs", "create", "http", shellQuote(jobName)}
	return buildMultilineCommand(base, flags), nil
}

// BuildCreateCloudSQLJobCommand は Cloud SQL 用の Pub/Sub ジョブ作成コマンドを生成する。
func (s *Service) BuildCreateCloudSQLJobCommand(params CreateCloudSQLJobParams) (string, error) {
	jobName := strings.TrimSpace(params.JobName)
	if jobName == "" {
		return "", fmt.Errorf("jobName は必須です")
	}
	projectID := strings.TrimSpace(params.ProjectID)
	if projectID == "" {
		return "", fmt.Errorf("projectID は必須です")
	}
	topic := strings.TrimSpace(params.PubsubTopic)
	if topic == "" {
		return "", fmt.Errorf("pubsubTopic は必須です")
	}
	dbInstance := strings.TrimSpace(params.DBInstanceID)
	if dbInstance == "" {
		return "", fmt.Errorf("dbInstanceID は必須です")
	}

	location := defaultIfEmpty(strings.TrimSpace(params.Location), defaultLocation)
	schedule := defaultIfEmpty(strings.TrimSpace(params.Schedule), defaultSchedule)
	description := defaultIfEmpty(strings.TrimSpace(params.Description), defaultDescriptionPubSub)
	timeZone := defaultIfEmpty(strings.TrimSpace(params.TimeZone), defaultTimeZone)
	action := defaultIfEmpty(strings.TrimSpace(params.Action), "start")

	payload := struct {
		Instance          string `json:"Instance"`
		Project           string `json:"Project"`
		Action            string `json:"Action"`
		DiscordWebhookURL string `json:"DiscordWebhookUrl"`
		CloudSQLIconURL   string `json:"CloudSqlIconUrl"`
	}{
		Instance:          dbInstance,
		Project:           projectID,
		Action:            action,
		DiscordWebhookURL: params.DiscordWebhookURL,
		CloudSQLIconURL:   params.CloudSQLIconURL,
	}

	messageBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("message body の生成に失敗しました: %w", err)
	}

	return s.BuildCreatePubSubJobCommand(CreatePubSubJobParams{
		JobName:     jobName,
		ProjectID:   projectID,
		Location:    location,
		Schedule:    schedule,
		Description: description,
		TimeZone:    timeZone,
		PubsubTopic: topic,
		MessageBody: string(messageBytes),
	})
}

// BuildCreateStartCloudSQLJobCommand は Cloud SQL 起動ジョブ作成コマンドを生成する。
func (s *Service) BuildCreateStartCloudSQLJobCommand(params CreateStartStopCloudSQLJobParams) (string, error) {
	dbInstance := strings.TrimSpace(params.DBInstanceID)
	if dbInstance == "" {
		return "", fmt.Errorf("dbInstanceID は必須です")
	}
	jobName := strings.TrimSpace(params.JobName)
	if jobName == "" {
		jobName = fmt.Sprintf("start-%s-instance", dbInstance)
	}
	schedule := defaultIfEmpty(strings.TrimSpace(params.Schedule), defaultSchedule)

	return s.BuildCreateCloudSQLJobCommand(CreateCloudSQLJobParams{
		JobName:           jobName,
		ProjectID:         params.ProjectID,
		Location:          params.Location,
		Schedule:          schedule,
		Description:       params.Description,
		TimeZone:          params.TimeZone,
		PubsubTopic:       params.PubsubTopic,
		DBInstanceID:      dbInstance,
		Action:            "start",
		DiscordWebhookURL: params.DiscordWebhookURL,
		CloudSQLIconURL:   params.CloudSQLIconURL,
	})
}

// BuildCreateStopCloudSQLJobCommand は Cloud SQL 停止ジョブ作成コマンドを生成する。
func (s *Service) BuildCreateStopCloudSQLJobCommand(params CreateStartStopCloudSQLJobParams) (string, error) {
	dbInstance := strings.TrimSpace(params.DBInstanceID)
	if dbInstance == "" {
		return "", fmt.Errorf("dbInstanceID は必須です")
	}
	jobName := strings.TrimSpace(params.JobName)
	if jobName == "" {
		jobName = fmt.Sprintf("stop-%s-instance", dbInstance)
	}
	schedule := defaultIfEmpty(strings.TrimSpace(params.Schedule), defaultStopSchedule)

	return s.BuildCreateCloudSQLJobCommand(CreateCloudSQLJobParams{
		JobName:           jobName,
		ProjectID:         params.ProjectID,
		Location:          params.Location,
		Schedule:          schedule,
		Description:       params.Description,
		TimeZone:          params.TimeZone,
		PubsubTopic:       params.PubsubTopic,
		DBInstanceID:      dbInstance,
		Action:            "stop",
		DiscordWebhookURL: params.DiscordWebhookURL,
		CloudSQLIconURL:   params.CloudSQLIconURL,
	})
}

// BuildListJobsCommand は一覧取得コマンドを生成する。
func (s *Service) BuildListJobsCommand(params ListJobsParams) (string, error) {
	base := []string{"gcloud", "scheduler", "jobs", "list"}
	flags := make([]string, 0, 2)

	if loc := strings.TrimSpace(params.Location); loc != "" {
		flags = append(flags, formatFlag("--location", loc))
	}
	if limit := strings.TrimSpace(params.Limit); limit != "" {
		flags = append(flags, formatFlag("--limit", limit))
	}

	return buildMultilineCommand(base, flags), nil
}

// BuildUpdateHTTPJobCommand は HTTP ジョブ更新コマンドを生成する。
func (s *Service) BuildUpdateHTTPJobCommand(params UpdateHTTPJobParams) (string, error) {
	jobName := strings.TrimSpace(params.JobName)
	if jobName == "" {
		return "", fmt.Errorf("jobName は必須です")
	}

	flags := make([]string, 0, 3)
	if schedule := strings.TrimSpace(params.Schedule); schedule != "" {
		flags = append(flags, formatFlag("--schedule", schedule))
	}
	if message := sanitizeMessageBody(params.MessageBody); message != "" {
		flags = append(flags, formatFlag("--message-body", message))
	}
	if headers := strings.TrimSpace(params.Headers); headers != "" {
		flags = append(flags, formatFlag("--headers", headers))
	}

	base := []string{"gcloud", "scheduler", "jobs", "update", "http", shellQuote(jobName)}
	return buildMultilineCommand(base, flags), nil
}

// BuildUpdatePubSubJobCommand は Pub/Sub ジョブ更新コマンドを生成する。
func (s *Service) BuildUpdatePubSubJobCommand(params UpdatePubSubJobParams) (string, error) {
	jobName := strings.TrimSpace(params.JobName)
	if jobName == "" {
		return "", fmt.Errorf("jobName は必須です")
	}

	flags := make([]string, 0, 2)
	if schedule := strings.TrimSpace(params.Schedule); schedule != "" {
		flags = append(flags, formatFlag("--schedule", schedule))
	}
	if message := sanitizeMessageBody(params.MessageBody); message != "" {
		flags = append(flags, formatFlag("--message-body", message))
	}

	base := []string{"gcloud", "scheduler", "jobs", "update", "pubsub", shellQuote(jobName)}
	return buildMultilineCommand(base, flags), nil
}

// BuildPauseJobCommand はジョブの一時停止コマンドを生成する。
func (s *Service) BuildPauseJobCommand(params JobControlParams) (string, error) {
	return buildJobControlCommand("pause", params)
}

// BuildResumeJobCommand はジョブの再開コマンドを生成する。
func (s *Service) BuildResumeJobCommand(params JobControlParams) (string, error) {
	return buildJobControlCommand("resume", params)
}

// BuildDeleteJobCommand はジョブの削除コマンドを生成する。
func (s *Service) BuildDeleteJobCommand(params JobControlParams) (string, error) {
	command, err := buildJobControlCommand("delete", params)
	if err != nil {
		return "", err
	}
	return command + " \\\n  --quiet", nil
}

// BuildRunJobCommand はジョブの即時実行コマンドを生成する。
func (s *Service) BuildRunJobCommand(params JobControlParams) (string, error) {
	return buildJobControlCommand("run", params)
}

// PrintHighlightedCommand は生成したコマンドを整形して標準出力に表示する。
func (s *Service) PrintHighlightedCommand(command string) {
	fmt.Println()
	fmt.Println("==============================")
	fmt.Println("生成された gcloud コマンド")
	fmt.Println("==============================")
	fmt.Println(command)
	fmt.Println("==============================")
}

// BuildNotificationWrappedCommand は通知付きシェルスクリプトを生成する。
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

	return strings.Join(lines, "\n"), true
}

// PrintNotificationScript は通知付きシェルスクリプトを整形して出力する。
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

func buildJobControlCommand(action string, params JobControlParams) (string, error) {
	jobName := strings.TrimSpace(params.JobName)
	if jobName == "" {
		return "", fmt.Errorf("jobName は必須です")
	}
	location := strings.TrimSpace(params.Location)
	if location == "" {
		return "", fmt.Errorf("location は必須です")
	}

	base := []string{"gcloud", "scheduler", "jobs", action, shellQuote(jobName)}
	flags := []string{formatFlag("--location", location)}
	return buildMultilineCommand(base, flags), nil
}

func buildMultilineCommand(base []string, flags []string) string {
	cmd := strings.Join(base, " ")
	if len(flags) == 0 {
		return cmd
	}

	var builder strings.Builder
	builder.WriteString(cmd)
	for _, flag := range flags {
		builder.WriteString(" \\\n  ")
		builder.WriteString(flag)
	}
	return builder.String()
}

func formatFlag(name, value string) string {
	return fmt.Sprintf("%s=%s", name, shellQuote(value))
}

func defaultIfEmpty(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func sanitizeMessageBody(body string) string {
	if body == "" {
		return ""
	}
	withoutCR := strings.ReplaceAll(body, "\r", "")
	return strings.ReplaceAll(withoutCR, "\n", "")
}

func buildSimpleNotificationCommand(content string) string {
	return buildDiscordWebhookCommand(content, "none", "")
}

func buildEmbedNotificationCommand(content, embedText, embedType string) string {
	return buildDiscordWebhookCommand(content, embedType, embedText)
}

func buildDiscordWebhookCommand(content, embedType, embedText string) string {
	lines := []string{
		fmt.Sprintf("%s \\", discordCLIPath),
		fmt.Sprintf("  -webhook-url \"$%s\" \\", discordWebhookEnvVarName),
		fmt.Sprintf("  -content-text %s \\", shellQuote(content)),
	}
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

func shellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
}

type notificationTemplate struct {
	startContent     string
	successContent   string
	successEmbedText string
	failureContent   string
	failureEmbedText string
}

var notificationTemplates = map[string]notificationTemplate{
	cfg.OperationCreatePubSubJob: {
		startContent:     "ジョブを作成するよ！",
		successContent:   "作成したよ！",
		successEmbedText: "ジョブを作成したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ジョブを作成できなかったよ…",
	},
	cfg.OperationCreateHTTPJob: {
		startContent:     "ジョブを作成するよ！",
		successContent:   "作成したよ！",
		successEmbedText: "ジョブを作成したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ジョブを作成できなかったよ…",
	},
	cfg.OperationCreateCloudSQLJob: {
		startContent:     "ジョブを作成するよ！",
		successContent:   "作成したよ！",
		successEmbedText: "ジョブを作成したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ジョブを作成できなかったよ…",
	},
	cfg.OperationCreateStartCloudSQLJob: {
		startContent:     "ジョブを作成するよ！",
		successContent:   "作成したよ！",
		successEmbedText: "ジョブを作成したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ジョブを作成できなかったよ…",
	},
	cfg.OperationCreateStopCloudSQLJob: {
		startContent:     "ジョブを作成するよ！",
		successContent:   "作成したよ！",
		successEmbedText: "ジョブを作成したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ジョブを作成できなかったよ…",
	},
	cfg.OperationListJobs: {
		startContent:     "ジョブを一覧表示するよ！",
		successContent:   "一覧表示したよ！",
		successEmbedText: "ジョブを一覧表示したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ジョブを一覧表示できなかったよ…",
	},
	cfg.OperationUpdateHTTPJob: {
		startContent:     "HTTPジョブを更新するよ！",
		successContent:   "更新したよ！",
		successEmbedText: "HTTPジョブを更新したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "HTTPジョブを更新できなかったよ…",
	},
	cfg.OperationUpdatePubSubJob: {
		startContent:     "PUB/SUBジョブを更新するよ！",
		successContent:   "更新したよ！",
		successEmbedText: "PUB/SUBジョブを更新したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "PUB/SUBジョブを更新できなかったよ…",
	},
	cfg.OperationPauseJob: {
		startContent:     "ジョブを一時停止するよ！",
		successContent:   "一時停止したよ！",
		successEmbedText: "ジョブを一時停止したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ジョブを一時停止できなかったよ…",
	},
	cfg.OperationResumeJob: {
		startContent:     "ジョブを再開するよ！",
		successContent:   "再開したよ！",
		successEmbedText: "ジョブを再開したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ジョブを再開できなかったよ…",
	},
	cfg.OperationDeleteJob: {
		startContent:     "ジョブを削除するよ！",
		successContent:   "削除したよ！",
		successEmbedText: "ジョブを削除したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ジョブを削除できなかったよ…",
	},
	cfg.OperationRunJob: {
		startContent:     "ジョブを強制実行するよ！",
		successContent:   "強制実行したよ！",
		successEmbedText: "ジョブを強制実行したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ジョブを強制実行できなかったよ…",
	},
}
