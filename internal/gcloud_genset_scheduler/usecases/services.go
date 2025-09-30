package usecases

import (
	"encoding/json"
	"fmt"
	"strings"
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

func shellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
}
