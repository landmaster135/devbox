package config

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	// OperationCreatePubSubJob は Cloud Run Functions を Pub/Sub 経由で起動するジョブ生成を表す。
	OperationCreatePubSubJob = "create-pubsub-job"
	// OperationCreateHTTPJob は HTTP 宛先向けジョブ生成を表す。
	OperationCreateHTTPJob = "create-http-job"
	// OperationCreateCloudSQLJob は Cloud SQL 用の Pub/Sub ジョブ生成を表す。
	OperationCreateCloudSQLJob = "create-cloud-sql-job"
	// OperationCreateStartCloudSQLJob は Cloud SQL インスタンス起動ジョブ生成を表す。
	OperationCreateStartCloudSQLJob = "create-start-cloud-sql-job"
	// OperationCreateStopCloudSQLJob は Cloud SQL インスタンス停止ジョブ生成を表す。
	OperationCreateStopCloudSQLJob = "create-stop-cloud-sql-job"
	// OperationListJobs はジョブ一覧取得を表す。
	OperationListJobs = "list-jobs"
	// OperationUpdateHTTPJob は HTTP ジョブの更新を表す。
	OperationUpdateHTTPJob = "update-http-job"
	// OperationUpdatePubSubJob は Pub/Sub ジョブの更新を表す。
	OperationUpdatePubSubJob = "update-pubsub-job"
	// OperationPauseJob はジョブの一時停止を表す。
	OperationPauseJob = "pause-job"
	// OperationResumeJob はジョブの再開を表す。
	OperationResumeJob = "resume-job"
	// OperationDeleteJob はジョブの削除を表す。
	OperationDeleteJob = "delete-job"
	// OperationRunJob はジョブの即時実行を表す。
	OperationRunJob = "run-job"
)

var supportedOperations = []string{
	OperationCreateCloudSQLJob,
	OperationCreateHTTPJob,
	OperationCreatePubSubJob,
	OperationCreateStartCloudSQLJob,
	OperationCreateStopCloudSQLJob,
	OperationDeleteJob,
	OperationListJobs,
	OperationPauseJob,
	OperationResumeJob,
	OperationRunJob,
	OperationUpdateHTTPJob,
	OperationUpdatePubSubJob,
}

var supportedHTTPMethods = []string{"DELETE", "GET", "HEAD", "PATCH", "POST", "PUT"}

// Config は CLI 引数から取得する値を保持する。
type Config struct {
	Operation               string
	Help                    bool
	JobName                 string
	ProjectID               string
	Location                string
	Schedule                string
	Description             string
	TimeZone                string
	PubsubTopic             string
	MessageBody             string
	HTTPMethod              string
	ServiceURL              string
	Headers                 string
	OIDCServiceAccountEmail string
	DBInstanceID            string
	Action                  string
	DiscordWebhookURL       string
	CloudSQLIconURL         string
	Limit                   string
}

// ParseFlags は標準パーサーを用いてフラグを解析する。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は注入したパーサーを利用して CLI 引数を解析する。
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{}

	parser.StringVar(&cfg.Operation, "operation", "", fmt.Sprintf("実行する操作 (%s)", strings.Join(supportedOperations, ", ")))
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示する")

	parser.StringVar(&cfg.JobName, "job-name", "", "Scheduler ジョブ名")
	parser.StringVar(&cfg.ProjectID, "project-id", "", "GCP プロジェクト ID")
	parser.StringVar(&cfg.Location, "location", "", "ジョブのロケーション (未指定時はデフォルト値を使用)")
	parser.StringVar(&cfg.Schedule, "schedule", "", "cron 形式のスケジュール文字列")
	parser.StringVar(&cfg.Description, "description", "", "ジョブの説明")
	parser.StringVar(&cfg.TimeZone, "time-zone", "", "スケジュールに利用するタイムゾーン")
	parser.StringVar(&cfg.PubsubTopic, "pubsub-topic", "", "Pub/Sub トピック名")
	parser.StringVar(&cfg.MessageBody, "message-body", "", "HTTP または Pub/Sub に送信するメッセージ本文")
	parser.StringVar(&cfg.HTTPMethod, "http-method", "", "HTTP ジョブで利用するメソッド (GET|POST|PUT|DELETE|PATCH|HEAD)")
	parser.StringVar(&cfg.ServiceURL, "service-url", "", "HTTP ジョブの送信先 URL")
	parser.StringVar(&cfg.Headers, "headers", "", "HTTP ヘッダー (例: Header1=value1,Header2=value2)")
	parser.StringVar(&cfg.OIDCServiceAccountEmail, "oidc-service-account-email", "", "OIDC トークンを生成するサービスアカウント")
	parser.StringVar(&cfg.DBInstanceID, "db-instance-id", "", "Cloud SQL インスタンス ID")
	parser.StringVar(&cfg.Action, "action", "", "Cloud SQL メッセージ本文に含めるアクション")
	parser.StringVar(&cfg.DiscordWebhookURL, "discord-webhook-url", "", "通知先の Discord Webhook URL")
	parser.StringVar(&cfg.CloudSQLIconURL, "cloud-sql-icon-url", "", "通知に利用するアイコン URL")
	parser.StringVar(&cfg.Limit, "limit", "", "一覧取得時の最大件数")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if cfg.Help {
		return cfg, nil
	}

	normalizeConfig(cfg)

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func normalizeConfig(cfg *Config) {
	cfg.Operation = strings.TrimSpace(cfg.Operation)
	cfg.JobName = strings.TrimSpace(cfg.JobName)
	cfg.ProjectID = strings.TrimSpace(cfg.ProjectID)
	cfg.Location = strings.TrimSpace(cfg.Location)
	cfg.Schedule = strings.TrimSpace(cfg.Schedule)
	cfg.Description = strings.TrimSpace(cfg.Description)
	cfg.TimeZone = strings.TrimSpace(cfg.TimeZone)
	cfg.PubsubTopic = strings.TrimSpace(cfg.PubsubTopic)
	cfg.MessageBody = strings.TrimSpace(cfg.MessageBody)
	cfg.HTTPMethod = strings.ToUpper(strings.TrimSpace(cfg.HTTPMethod))
	cfg.ServiceURL = strings.TrimSpace(cfg.ServiceURL)
	cfg.Headers = strings.TrimSpace(cfg.Headers)
	cfg.OIDCServiceAccountEmail = strings.TrimSpace(cfg.OIDCServiceAccountEmail)
	cfg.DBInstanceID = strings.TrimSpace(cfg.DBInstanceID)
	cfg.Action = strings.TrimSpace(cfg.Action)
	cfg.DiscordWebhookURL = strings.TrimSpace(cfg.DiscordWebhookURL)
	cfg.CloudSQLIconURL = strings.TrimSpace(cfg.CloudSQLIconURL)
	cfg.Limit = strings.TrimSpace(cfg.Limit)
}

func validateConfig(cfg *Config) error {
	if cfg.Operation == "" {
		return fmt.Errorf("operation パラメータは必須です")
	}
	if !isSupportedOperation(cfg.Operation) {
		return fmt.Errorf("未対応のoperationです: %s", cfg.Operation)
	}

	switch cfg.Operation {
	case OperationCreatePubSubJob:
		if cfg.JobName == "" {
			return fmt.Errorf("job-name パラメータは必須です")
		}
		if cfg.ProjectID == "" {
			return fmt.Errorf("project-id パラメータは必須です")
		}
		if cfg.PubsubTopic == "" {
			return fmt.Errorf("pubsub-topic パラメータは必須です")
		}
	case OperationCreateHTTPJob:
		if cfg.JobName == "" {
			return fmt.Errorf("job-name パラメータは必須です")
		}
		if cfg.ProjectID == "" {
			return fmt.Errorf("project-id パラメータは必須です")
		}
		if cfg.HTTPMethod == "" {
			return fmt.Errorf("http-method パラメータは必須です")
		}
		if !isSupportedHTTPMethod(cfg.HTTPMethod) {
			return fmt.Errorf("http-method には %s のいずれかを指定してください", strings.Join(supportedHTTPMethods, ", "))
		}
		if cfg.ServiceURL == "" {
			return fmt.Errorf("service-url パラメータは必須です")
		}
	case OperationCreateCloudSQLJob:
		if cfg.JobName == "" {
			return fmt.Errorf("job-name パラメータは必須です")
		}
		if err := ensureCloudSQLCommonParams(cfg); err != nil {
			return err
		}
	case OperationCreateStartCloudSQLJob, OperationCreateStopCloudSQLJob:
		if err := ensureCloudSQLCommonParams(cfg); err != nil {
			return err
		}
	case OperationUpdateHTTPJob:
		if cfg.JobName == "" {
			return fmt.Errorf("job-name パラメータは必須です")
		}
	case OperationUpdatePubSubJob:
		if cfg.JobName == "" {
			return fmt.Errorf("job-name パラメータは必須です")
		}
	case OperationPauseJob, OperationResumeJob, OperationDeleteJob, OperationRunJob:
		if cfg.JobName == "" {
			return fmt.Errorf("job-name パラメータは必須です")
		}
		if cfg.Location == "" {
			return fmt.Errorf("location パラメータは必須です")
		}
	case OperationListJobs:
		// 追加必須項目なし
	}

	if cfg.Limit != "" {
		if _, err := strconv.Atoi(cfg.Limit); err != nil {
			return fmt.Errorf("limit には数値を指定してください")
		}
	}

	return nil
}

func ensureCloudSQLCommonParams(cfg *Config) error {
	if cfg.ProjectID == "" {
		return fmt.Errorf("project-id パラメータは必須です")
	}
	if cfg.PubsubTopic == "" {
		return fmt.Errorf("pubsub-topic パラメータは必須です")
	}
	if cfg.DBInstanceID == "" {
		return fmt.Errorf("db-instance-id パラメータは必須です")
	}
	return nil
}

func isSupportedOperation(op string) bool {
	for _, candidate := range supportedOperations {
		if candidate == op {
			return true
		}
	}
	return false
}

func isSupportedHTTPMethod(method string) bool {
	for _, candidate := range supportedHTTPMethods {
		if candidate == method {
			return true
		}
	}
	return false
}

// PrintUsage は CLI の利用方法を標準エラーに出力する。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Cloud Scheduler 用 gcloud コマンド生成ツール\n\n")

	fmt.Fprintf(os.Stderr, "共通パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n        実行する操作 (%s)\n", strings.Join(supportedOperations, ", "))
	fmt.Fprintf(os.Stderr, "  -help\n        このヘルプを表示\n\n")

	fmt.Fprintf(os.Stderr, "主要パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -job-name string\n        対象ジョブ名 (操作により必須)\n")
	fmt.Fprintf(os.Stderr, "  -project-id string\n        GCP プロジェクト ID (ジョブ作成時に必須)\n")
	fmt.Fprintf(os.Stderr, "  -location string\n        ジョブのロケーション (未指定ならデフォルト us-central1)\n")
	fmt.Fprintf(os.Stderr, "  -schedule string\n        cron 形式のスケジュール (未指定なら操作に応じた既定値)\n")
	fmt.Fprintf(os.Stderr, "  -time-zone string\n        スケジュールのタイムゾーン (未指定なら Asia/Tokyo)\n")
	fmt.Fprintf(os.Stderr, "  -pubsub-topic string\n        Pub/Sub トピック名 (Pub/Sub 系操作で必須)\n")
	fmt.Fprintf(os.Stderr, "  -http-method string, -service-url string\n        HTTP ジョブのメソッドと送信先\n")
	fmt.Fprintf(os.Stderr, "  -db-instance-id string\n        Cloud SQL インスタンス ID (Cloud SQL 系操作で必須)\n")
	fmt.Fprintf(os.Stderr, "  -message-body string\n        リクエスト本文\n")
	fmt.Fprintf(os.Stderr, "  -headers string\n        追加ヘッダーをカンマ区切りで指定\n")
	fmt.Fprintf(os.Stderr, "  -oidc-service-account-email string\n        OIDC トークン生成用のサービスアカウント\n")
	fmt.Fprintf(os.Stderr, "  -discord-webhook-url string, -cloud-sql-icon-url string\n        Cloud SQL 用メッセージに埋め込む URL\n")
	fmt.Fprintf(os.Stderr, "  -limit string\n        list-jobs 時の取得上限\n\n")

	fmt.Fprintf(os.Stderr, "操作ごとの必須フラグ概要:\n")
	fmt.Fprintf(os.Stderr, "  %s: -job-name, -project-id, -pubsub-topic\n", OperationCreatePubSubJob)
	fmt.Fprintf(os.Stderr, "  %s: -job-name, -project-id, -http-method, -service-url\n", OperationCreateHTTPJob)
	fmt.Fprintf(os.Stderr, "  %s: -job-name, -project-id, -pubsub-topic, -db-instance-id\n", OperationCreateCloudSQLJob)
	fmt.Fprintf(os.Stderr, "  %s/%s: -project-id, -pubsub-topic, -db-instance-id (job-name は省略可)\n", OperationCreateStartCloudSQLJob, OperationCreateStopCloudSQLJob)
	fmt.Fprintf(os.Stderr, "  %s/%s/%s/%s: -job-name, -location\n", OperationPauseJob, OperationResumeJob, OperationDeleteJob, OperationRunJob)
	fmt.Fprintf(os.Stderr, "  %s/%s: -job-name\n", OperationUpdateHTTPJob, OperationUpdatePubSubJob)

	fmt.Fprintf(os.Stderr, "\n使用例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=%s -job-name=exec-cloud-run -project-id=my-project -pubsub-topic=my-topic\n", os.Args[0], OperationCreatePubSubJob)
	fmt.Fprintf(os.Stderr, "  %s -operation=%s -job-name=http-job -project-id=my-project -http-method=POST -service-url=https://example.com\n", os.Args[0], OperationCreateHTTPJob)
	fmt.Fprintf(os.Stderr, "  %s -operation=%s -job-name=start-my-db -project-id=my-project -pubsub-topic=my-topic -db-instance-id=my-db\n", os.Args[0], OperationCreateCloudSQLJob)
}

func init() {
	sort.Strings(supportedOperations)
	sort.Strings(supportedHTTPMethods)
}
