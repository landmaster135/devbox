package config

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	// OperationDeployCloudRunContainer は Cloud Run コンテナのデプロイコマンドを生成する操作を表す。
	OperationDeployCloudRunContainer = "deploy-cloud-run-container"
	// OperationUpdateCloudRunContainerEnv は Cloud Run コンテナの環境変数ファイルを適用する操作を表す。
	OperationUpdateCloudRunContainerEnv = "update-cloud-run-container-env"
	// OperationDeployCloudRunFunction は HTTP トリガーの Cloud Functions (Gen2) をデプロイする操作を表す。
	OperationDeployCloudRunFunction = "deploy-cloud-run-function"
	// OperationDeployCloudRunFunctionTriggeredByPubSub は Pub/Sub トリガーの Cloud Functions (Gen2) をデプロイする操作を表す。
	OperationDeployCloudRunFunctionTriggeredByPubSub = "deploy-cloud-run-function-triggered-by-pubsub"
	// OperationUpdateCloudRunFunctionEnv は Cloud Functions (Gen2) の環境変数を更新する操作を表す。
	OperationUpdateCloudRunFunctionEnv = "update-cloud-run-function-env"
	// OperationUpdateCloudRunServiceEnv は Cloud Run サービスの環境変数をファイルから更新する操作を表す。
	OperationUpdateCloudRunServiceEnv = "update-cloud-run-service-env"
	// OperationCreateCloudPubSubTopic は Pub/Sub トピックを作成する操作を表す。
	OperationCreateCloudPubSubTopic = "create-cloud-pubsub-topic"
	// OperationListCloudPubSubTopics は Pub/Sub トピックをフィルタ付きで一覧表示する操作を表す。
	OperationListCloudPubSubTopics = "list-cloud-pubsub-topics"
	// OperationListCloudPubSubSubscriptions は Pub/Sub サブスクリプションを一覧表示する操作を表す。
	OperationListCloudPubSubSubscriptions = "list-cloud-pubsub-subscriptions"
	// OperationCreateCloudPubSubSubscription は Push 形式の Pub/Sub サブスクリプションを作成する操作を表す。
	OperationCreateCloudPubSubSubscription = "create-cloud-pubsub-subscription"
	// OperationDeleteCloudPubSubSubscriptionsAndTopics はサブスクリプションとトピックを同時に削除する操作を表す。
	OperationDeleteCloudPubSubSubscriptionsAndTopics = "delete-cloud-pubsub-subscriptions-and-topics"
	// OperationDeleteCloudPubSubSubscriptions はサブスクリプションのみを削除する操作を表す。
	OperationDeleteCloudPubSubSubscriptions = "delete-cloud-pubsub-subscriptions"
	// OperationDeleteCloudPubSubTopics はトピックのみを削除する操作を表す。
	OperationDeleteCloudPubSubTopics = "delete-cloud-pubsub-topics"
	// OperationDeleteCloudRunFunction は Cloud Functions (Gen2) が利用する Cloud Run サービスを削除する操作を表す。
	OperationDeleteCloudRunFunction = "delete-cloud-run-function"
)

var validOperations = []string{
	OperationDeployCloudRunContainer,
	OperationUpdateCloudRunContainerEnv,
	OperationDeployCloudRunFunction,
	OperationDeployCloudRunFunctionTriggeredByPubSub,
	OperationUpdateCloudRunFunctionEnv,
	OperationUpdateCloudRunServiceEnv,
	OperationCreateCloudPubSubTopic,
	OperationListCloudPubSubTopics,
	OperationListCloudPubSubSubscriptions,
	OperationCreateCloudPubSubSubscription,
	OperationDeleteCloudPubSubSubscriptionsAndTopics,
	OperationDeleteCloudPubSubSubscriptions,
	OperationDeleteCloudPubSubTopics,
	OperationDeleteCloudRunFunction,
}

const (
	defaultRegion                   = "us-central1"
	defaultContainerTimeout         = "40m"
	defaultEnvFile                  = "env.yml"
	defaultMessageRetentionDuration = "1d"
	defaultExpirationPeriod         = "never"
	defaultMaxRetryDelay            = "600s"
	defaultMinRetryDelay            = "10s"
	defaultAckDeadline              = "600"
	defaultPubSubEntryPoint         = "ProcessPubSub"
)

// Config は CLI から渡されるパラメータを保持する。
type Config struct {
	Operation string
	Help      bool

	ServiceName           string
	FunctionName          string
	ProjectID             string
	Region                string
	Timeout               string
	RunServiceAccount     string
	AllowUnauthenticated  bool
	EnvFile               string
	EnvVars               string
	EntryPoint            string
	TriggerServiceAccount string
	TriggerTopic          string
	APIClientID           string
	APIClientSecret       string
	APIEndpoint           string

	TopicName        string
	TopicProject     string
	SubscriptionName string
	ShowURI          bool
	PushEndpoint     string

	MessageRetentionDuration string
	ExpirationPeriod         string
	MaxRetryDelay            string
	MinRetryDelay            string
	AckDeadline              string

	SubscriptionNames []string
	TopicNames        []string

	rawSubscriptionNames string
	rawTopicNames        string

	PushServiceAccount string
}

// ParseFlags は標準のフラグパーサーで引数を解析する。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたパーサーを用いて引数を解析する。
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{}

	parser.StringVar(&cfg.Operation, "operation", "", fmt.Sprintf("実行する操作 (%s)", strings.Join(validOperations, ", ")))
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示する")

	parser.StringVar(&cfg.ServiceName, "service-name", "", "Cloud Run サービス名")
	parser.StringVar(&cfg.FunctionName, "function-name", "", "Cloud Functions (Gen2) の関数名")
	parser.StringVar(&cfg.ProjectID, "project-id", "", "対象の GCP プロジェクトID")
	parser.StringVar(&cfg.Region, "region", "", "対象リージョン (operation により必須/任意)")
	parser.StringVar(&cfg.Timeout, "timeout", defaultContainerTimeout, "Cloud Run デプロイ時のタイムアウト (例: 40m)")
	parser.StringVar(&cfg.RunServiceAccount, "run-service-account", "", "Cloud Run デプロイで使用するサービスアカウントのメールアドレス")
	parser.BoolVar(&cfg.AllowUnauthenticated, "allow-unauthenticated", true, "認証なしアクセスを許可するか")
	parser.StringVar(&cfg.EnvFile, "env-file", defaultEnvFile, "環境変数を記述したファイルパス")
	parser.StringVar(&cfg.EnvVars, "env-vars", "", "環境変数 (KEY=VALUE,KEY2=VALUE2 形式)")

	parser.StringVar(&cfg.EntryPoint, "entry-point", "", "Cloud Functions のエントリポイント")
	parser.StringVar(&cfg.TriggerServiceAccount, "trigger-service-account", "", "Pub/Sub トリガーに使用するサービスアカウント")
	parser.StringVar(&cfg.TriggerTopic, "trigger-topic", "", "Pub/Sub トリガートピックID")
	parser.StringVar(&cfg.APIClientID, "api-client-id", "", "任意: API クライアントID")
	parser.StringVar(&cfg.APIClientSecret, "api-client-secret", "", "任意: API クライアントシークレット")
	parser.StringVar(&cfg.APIEndpoint, "api-endpoint", "", "任意: API エンドポイントURL")

	parser.StringVar(&cfg.TopicName, "topic-name", "", "Pub/Sub トピック名")
	parser.StringVar(&cfg.TopicProject, "topic-project", "", "Pub/Sub トピックが属する GCP プロジェクトID")
	parser.StringVar(&cfg.SubscriptionName, "subscription-name", "", "Pub/Sub サブスクリプション名 (フィルタ用)")
	parser.BoolVar(&cfg.ShowURI, "show-uri", false, "Pub/Sub 一覧出力に URI を含める")
	parser.StringVar(&cfg.PushEndpoint, "push-endpoint", "", "Push サブスクリプションのエンドポイント URL")

	parser.StringVar(&cfg.MessageRetentionDuration, "message-retention-duration", defaultMessageRetentionDuration, "メッセージ保持期間 (例: 1d)")
	parser.StringVar(&cfg.ExpirationPeriod, "expiration-period", defaultExpirationPeriod, "期限切れまでの期間 (例: never)")
	parser.StringVar(&cfg.MaxRetryDelay, "max-retry-delay", defaultMaxRetryDelay, "最大リトライ遅延 (例: 600s)")
	parser.StringVar(&cfg.MinRetryDelay, "min-retry-delay", defaultMinRetryDelay, "最小リトライ遅延 (例: 10s)")
	parser.StringVar(&cfg.AckDeadline, "ack-deadline", defaultAckDeadline, "ACK 期限 (秒)")

	parser.StringVar(&cfg.rawSubscriptionNames, "subscription-names", "", "削除対象のサブスクリプション名 (カンマ区切り)")
	parser.StringVar(&cfg.rawTopicNames, "topic-names", "", "削除対象のトピック名 (カンマ区切り)")
	parser.StringVar(&cfg.PushServiceAccount, "push-service-account", "", "Push サブスクリプションで使用するサービスアカウントのメールアドレス")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if cfg.Help {
		return cfg, nil
	}

	if err := parseLists(cfg); err != nil {
		return nil, err
	}

	normalizeConfig(cfg)
	applyDefaults(cfg)

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseLists(cfg *Config) error {
	subs, err := parseCommaSeparated(cfg.rawSubscriptionNames, "subscription-names")
	if err != nil {
		return err
	}
	cfg.SubscriptionNames = subs

	topics, err := parseCommaSeparated(cfg.rawTopicNames, "topic-names")
	if err != nil {
		return err
	}
	cfg.TopicNames = topics
	return nil
}

func parseCommaSeparated(value string, fieldName string) ([]string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}

	parts := strings.Split(trimmed, ",")
	var result []string
	for _, part := range parts {
		candidate := strings.TrimSpace(part)
		if candidate != "" {
			result = append(result, candidate)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("%s の形式が不正です", fieldName)
	}

	return result, nil
}

func normalizeConfig(cfg *Config) {
	cfg.Operation = strings.TrimSpace(cfg.Operation)
	cfg.ServiceName = strings.TrimSpace(cfg.ServiceName)
	cfg.FunctionName = strings.TrimSpace(cfg.FunctionName)
	cfg.ProjectID = strings.TrimSpace(cfg.ProjectID)
	cfg.Region = strings.TrimSpace(cfg.Region)
	cfg.Timeout = strings.TrimSpace(cfg.Timeout)
	cfg.RunServiceAccount = strings.TrimSpace(cfg.RunServiceAccount)
	cfg.EnvFile = strings.TrimSpace(cfg.EnvFile)
	cfg.EnvVars = strings.TrimSpace(cfg.EnvVars)
	cfg.EntryPoint = strings.TrimSpace(cfg.EntryPoint)
	cfg.TriggerServiceAccount = strings.TrimSpace(cfg.TriggerServiceAccount)
	cfg.TriggerTopic = strings.TrimSpace(cfg.TriggerTopic)
	cfg.APIClientID = strings.TrimSpace(cfg.APIClientID)
	cfg.APIClientSecret = strings.TrimSpace(cfg.APIClientSecret)
	cfg.APIEndpoint = strings.TrimSpace(cfg.APIEndpoint)
	cfg.TopicName = strings.TrimSpace(cfg.TopicName)
	cfg.TopicProject = strings.TrimSpace(cfg.TopicProject)
	cfg.SubscriptionName = strings.TrimSpace(cfg.SubscriptionName)
	cfg.PushEndpoint = strings.TrimSpace(cfg.PushEndpoint)
	cfg.MessageRetentionDuration = strings.TrimSpace(cfg.MessageRetentionDuration)
	cfg.ExpirationPeriod = strings.TrimSpace(cfg.ExpirationPeriod)
	cfg.MaxRetryDelay = strings.TrimSpace(cfg.MaxRetryDelay)
	cfg.MinRetryDelay = strings.TrimSpace(cfg.MinRetryDelay)
	cfg.AckDeadline = strings.TrimSpace(cfg.AckDeadline)
	cfg.PushServiceAccount = strings.TrimSpace(cfg.PushServiceAccount)

	cfg.SubscriptionNames = normalizeList(cfg.SubscriptionNames)
	cfg.TopicNames = normalizeList(cfg.TopicNames)
}

func normalizeList(values []string) []string {
	if len(values) == 0 {
		return values
	}
	normalized := make([]string, 0, len(values))
	for _, v := range values {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}
	return normalized
}

func applyDefaults(cfg *Config) {
	if cfg.Timeout == "" {
		cfg.Timeout = defaultContainerTimeout
	}
	if cfg.EnvFile == "" {
		cfg.EnvFile = defaultEnvFile
	}
	if cfg.MessageRetentionDuration == "" {
		cfg.MessageRetentionDuration = defaultMessageRetentionDuration
	}
	if cfg.ExpirationPeriod == "" {
		cfg.ExpirationPeriod = defaultExpirationPeriod
	}
	if cfg.MaxRetryDelay == "" {
		cfg.MaxRetryDelay = defaultMaxRetryDelay
	}
	if cfg.MinRetryDelay == "" {
		cfg.MinRetryDelay = defaultMinRetryDelay
	}
	if cfg.AckDeadline == "" {
		cfg.AckDeadline = defaultAckDeadline
	}

	switch cfg.Operation {
	case OperationDeployCloudRunContainer:
		if cfg.Region == "" {
			cfg.Region = defaultRegion
		}
	case OperationDeployCloudRunFunctionTriggeredByPubSub:
		if cfg.Region == "" {
			cfg.Region = defaultRegion
		}
		if cfg.EntryPoint == "" {
			cfg.EntryPoint = defaultPubSubEntryPoint
		}
	}
}

func validateConfig(cfg *Config) error {
	if cfg.Operation == "" {
		return fmt.Errorf("operation パラメータは必須です")
	}
	if !isValidOperation(cfg.Operation) {
		return fmt.Errorf("未対応のoperationです: %s", cfg.Operation)
	}

	switch cfg.Operation {
	case OperationDeployCloudRunContainer:
		if cfg.ServiceName == "" {
			return fmt.Errorf("service-name パラメータは必須です")
		}
		if cfg.ProjectID == "" {
			return fmt.Errorf("project-id パラメータは必須です")
		}
		if cfg.Region == "" {
			return fmt.Errorf("region パラメータは必須です")
		}
		if cfg.Timeout == "" {
			return fmt.Errorf("timeout パラメータは必須です")
		}
	case OperationUpdateCloudRunContainerEnv:
		if cfg.ServiceName == "" {
			return fmt.Errorf("service-name パラメータは必須です")
		}
		if cfg.ProjectID == "" {
			return fmt.Errorf("project-id パラメータは必須です")
		}
		if cfg.Region == "" {
			return fmt.Errorf("region パラメータは必須です")
		}
		if cfg.EnvFile == "" {
			return fmt.Errorf("env-file パラメータは必須です")
		}
	case OperationDeployCloudRunFunction:
		if cfg.FunctionName == "" {
			return fmt.Errorf("function-name パラメータは必須です")
		}
		if cfg.Region == "" {
			return fmt.Errorf("region パラメータは必須です")
		}
		if cfg.EntryPoint == "" {
			return fmt.Errorf("entry-point パラメータは必須です")
		}
	case OperationDeployCloudRunFunctionTriggeredByPubSub:
		if cfg.FunctionName == "" {
			return fmt.Errorf("function-name パラメータは必須です")
		}
		if cfg.ProjectID == "" {
			return fmt.Errorf("project-id パラメータは必須です")
		}
		if cfg.Region == "" {
			return fmt.Errorf("region パラメータは必須です")
		}
		if cfg.TriggerServiceAccount == "" {
			return fmt.Errorf("trigger-service-account パラメータは必須です")
		}
		if cfg.TriggerTopic == "" {
			return fmt.Errorf("trigger-topic パラメータは必須です")
		}
		if cfg.EntryPoint == "" {
			return fmt.Errorf("entry-point パラメータは必須です")
		}
	case OperationUpdateCloudRunFunctionEnv:
		if cfg.ServiceName == "" {
			return fmt.Errorf("service-name パラメータは必須です")
		}
		if cfg.Region == "" {
			return fmt.Errorf("region パラメータは必須です")
		}
		if cfg.EnvVars == "" {
			return fmt.Errorf("env-vars パラメータは必須です")
		}
	case OperationUpdateCloudRunServiceEnv:
		if cfg.ServiceName == "" {
			return fmt.Errorf("service-name パラメータは必須です")
		}
		if cfg.ProjectID == "" {
			return fmt.Errorf("project-id パラメータは必須です")
		}
		if cfg.Region == "" {
			return fmt.Errorf("region パラメータは必須です")
		}
		if cfg.EnvFile == "" {
			return fmt.Errorf("env-file パラメータは必須です")
		}
	case OperationCreateCloudPubSubTopic:
		if cfg.TopicName == "" {
			return fmt.Errorf("topic-name パラメータは必須です")
		}
		if cfg.MessageRetentionDuration == "" {
			return fmt.Errorf("message-retention-duration パラメータは必須です")
		}
	case OperationListCloudPubSubTopics:
		if cfg.TopicName == "" {
			return fmt.Errorf("topic-name パラメータは必須です")
		}
	case OperationListCloudPubSubSubscriptions:
		// 特別な必須項目はなし
	case OperationCreateCloudPubSubSubscription:
		if cfg.SubscriptionName == "" {
			return fmt.Errorf("subscription-name パラメータは必須です")
		}
		if cfg.TopicName == "" {
			return fmt.Errorf("topic-name パラメータは必須です")
		}
		if cfg.TopicProject == "" {
			return fmt.Errorf("topic-project パラメータは必須です")
		}
		if cfg.PushServiceAccount == "" {
			return fmt.Errorf("push-service-account パラメータは必須です")
		}
		if cfg.MessageRetentionDuration == "" {
			return fmt.Errorf("message-retention-duration パラメータは必須です")
		}
		if cfg.ExpirationPeriod == "" {
			return fmt.Errorf("expiration-period パラメータは必須です")
		}
		if cfg.MaxRetryDelay == "" {
			return fmt.Errorf("max-retry-delay パラメータは必須です")
		}
		if cfg.MinRetryDelay == "" {
			return fmt.Errorf("min-retry-delay パラメータは必須です")
		}
		if cfg.AckDeadline == "" {
			return fmt.Errorf("ack-deadline パラメータは必須です")
		}
	case OperationDeleteCloudPubSubSubscriptionsAndTopics:
		if len(cfg.SubscriptionNames) == 0 && len(cfg.TopicNames) == 0 {
			return fmt.Errorf("subscription-names または topic-names のいずれかは必須です")
		}
	case OperationDeleteCloudPubSubSubscriptions:
		if len(cfg.SubscriptionNames) == 0 {
			return fmt.Errorf("subscription-names パラメータは必須です")
		}
	case OperationDeleteCloudPubSubTopics:
		if len(cfg.TopicNames) == 0 {
			return fmt.Errorf("topic-names パラメータは必須です")
		}
	case OperationDeleteCloudRunFunction:
		if cfg.ServiceName == "" {
			return fmt.Errorf("service-name パラメータは必須です")
		}
		if cfg.Region == "" {
			return fmt.Errorf("region パラメータは必須です")
		}
	}

	return nil
}

func isValidOperation(operation string) bool {
	for _, candidate := range validOperations {
		if candidate == operation {
			return true
		}
	}
	return false
}

// PrintUsage は CLI の利用方法を標準エラーに出力する。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Cloud Run / Cloud Functions / Pub/Sub 向け gcloud コマンド生成ツール\n\n")

	fmt.Fprintf(os.Stderr, "共通パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n        実行する操作 (%s)\n", strings.Join(validOperations, ", "))
	fmt.Fprintf(os.Stderr, "  -help\n        このヘルプを表示\n\n")

	fmt.Fprintf(os.Stderr, "Cloud Run コンテナ関連:\n")
	fmt.Fprintf(os.Stderr, "  %s:\n    -service-name (必須), -project-id (必須), [-region], [-timeout], [-run-service-account], [-allow-unauthenticated]\n", OperationDeployCloudRunContainer)
	fmt.Fprintf(os.Stderr, "  %s:\n    -service-name (必須), -project-id (必須), -region (必須), [-env-file]\n\n", OperationUpdateCloudRunContainerEnv)

	fmt.Fprintf(os.Stderr, "Cloud Functions (Gen2) 関連:\n")
	fmt.Fprintf(os.Stderr, "  %s:\n    -function-name (必須), -region (必須), -entry-point (必須)\n", OperationDeployCloudRunFunction)
	fmt.Fprintf(os.Stderr, "  %s:\n    -function-name (必須), -project-id (必須), [-region], [-entry-point], -trigger-service-account (必須), -trigger-topic (必須), [-api-*]\n", OperationDeployCloudRunFunctionTriggeredByPubSub)
	fmt.Fprintf(os.Stderr, "  %s:\n    -service-name (必須), -region (必須), -env-vars (必須)\n", OperationUpdateCloudRunFunctionEnv)
	fmt.Fprintf(os.Stderr, "  %s:\n    -service-name (必須), -project-id (必須), -region (必須), [-env-file]\n", OperationUpdateCloudRunServiceEnv)
	fmt.Fprintf(os.Stderr, "  %s:\n    -service-name (必須), -region (必須)\n\n", OperationDeleteCloudRunFunction)

	fmt.Fprintf(os.Stderr, "Pub/Sub 関連:\n")
	fmt.Fprintf(os.Stderr, "  %s:\n    -topic-name (必須), [-message-retention-duration]\n", OperationCreateCloudPubSubTopic)
	fmt.Fprintf(os.Stderr, "  %s:\n    -topic-name (必須)\n", OperationListCloudPubSubTopics)
	fmt.Fprintf(os.Stderr, "  %s:\n    [-subscription-name], [-show-uri]\n", OperationListCloudPubSubSubscriptions)
	fmt.Fprintf(os.Stderr, "  %s:\n    -subscription-name (必須), -topic-name (必須), -topic-project (必須), -push-service-account (必須), [-push-endpoint], [-message-retention-duration], [-expiration-period], [-max-retry-delay], [-min-retry-delay], [-ack-deadline]\n", OperationCreateCloudPubSubSubscription)
	fmt.Fprintf(os.Stderr, "  %s:\n    [-subscription-names], [-topic-names] (少なくともどちらか必須)\n", OperationDeleteCloudPubSubSubscriptionsAndTopics)
	fmt.Fprintf(os.Stderr, "  %s:\n    -subscription-names (必須)\n", OperationDeleteCloudPubSubSubscriptions)
	fmt.Fprintf(os.Stderr, "  %s:\n    -topic-names (必須)\n", OperationDeleteCloudPubSubTopics)
}

func init() {
	sort.Strings(validOperations)
}
