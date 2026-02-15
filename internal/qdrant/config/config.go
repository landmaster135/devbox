package config

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	OperationCreateCollection   = "create-collection"
	OperationListCollections    = "list-collections"
	OperationUpsertTexts        = "upsert-texts"
	OperationQueryPoints        = "query-points"
	OperationDescribeCollection = "describe-collection"
	OperationDeleteCollection   = "delete-collection"
	OperationOverwritePayload   = "overwrite-payload"

	defaultDBHost         = "127.0.0.1"
	defaultDBPort         = 6334
	defaultEmbeddingHost  = "127.0.0.1"
	defaultEmbeddingPort  = 11434
	defaultCollectionSize = 4096
	defaultQueryLimit     = 5
)

var supportedOperations = []string{
	OperationCreateCollection,
	OperationListCollections,
	OperationUpsertTexts,
	OperationQueryPoints,
	OperationDescribeCollection,
	OperationDeleteCollection,
	OperationOverwritePayload,
}

// Config は CLI で受け取る設定値を保持する。
type Config struct {
	Operation      string
	DBHost         string
	DBPort         int
	CollectionName string
	Size           int

	EmbeddingHost  string
	EmbeddingPort  int
	EmbeddingModel string

	Input      string
	Payload    string
	QueryLimit int

	FilterMust           []KeyValue
	FilterMustNot        []KeyValue
	FilterShould         []KeyValue
	FilterMinShouldCount int

	Help bool
}

// ParseFlags は CLI フラグを解析する。
func ParseFlags() (*Config, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cfg := &Config{
		DBHost:        defaultDBHost,
		DBPort:        defaultDBPort,
		EmbeddingHost: defaultEmbeddingHost,
		EmbeddingPort: defaultEmbeddingPort,
		Size:          defaultCollectionSize,
		QueryLimit:    defaultQueryLimit,
	}

	fs.StringVar(&cfg.Operation, "operation", cfg.Operation, fmt.Sprintf("実行する操作 (%s)", strings.Join(supportedOperations, ", ")))
	fs.StringVar(&cfg.DBHost, "db-host", cfg.DBHost, "Qdrant のホスト名")
	fs.IntVar(&cfg.DBPort, "db-port", cfg.DBPort, "Qdrant のポート番号")
	fs.StringVar(&cfg.CollectionName, "collection-name", cfg.CollectionName, "対象コレクション名")
	fs.IntVar(&cfg.Size, "size", cfg.Size, "create-collection のベクトル次元")

	fs.StringVar(&cfg.EmbeddingHost, "embedding-host", cfg.EmbeddingHost, "埋め込みサービスのホスト名")
	fs.IntVar(&cfg.EmbeddingPort, "embedding-port", cfg.EmbeddingPort, "埋め込みサービスのポート番号")
	fs.StringVar(&cfg.EmbeddingModel, "embedding-model", cfg.EmbeddingModel, "埋め込みモデル名")

	fs.StringVar(&cfg.Input, "input", cfg.Input, "埋め込み対象テキスト (単一指定)")
	fs.StringVar(&cfg.Payload, "payload", cfg.Payload, "payload 条件 (key=value) の単一指定")
	fs.IntVar(&cfg.QueryLimit, "limit", cfg.QueryLimit, "query-points で取得する最大件数")
	fs.IntVar(&cfg.FilterMinShouldCount, "filter-min-should", cfg.FilterMinShouldCount, "should 条件のうち最低いくつを満たすか (overwrite-payload 用)")

	fs.Var(&keyValueList{target: &cfg.FilterMust}, "filter-must", "overwrite-payload 用の must 条件 (key=value、複数指定可)")
	fs.Var(&keyValueList{target: &cfg.FilterMustNot}, "filter-must-not", "overwrite-payload 用の must-not 条件 (key=value、複数指定可)")
	fs.Var(&keyValueList{target: &cfg.FilterShould}, "filter-should", "overwrite-payload 用の should 条件 (key=value、複数指定可)")
	fs.BoolVar(&cfg.Help, "help", false, "ヘルプを表示する")

	if err := fs.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			return &Config{Help: true}, nil
		}
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if len(fs.Args()) > 0 {
		return nil, fmt.Errorf("未処理の位置引数があります: %v", fs.Args())
	}

	normalizeConfig(cfg)

	if cfg.Help {
		return cfg, nil
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func normalizeConfig(cfg *Config) {
	cfg.Operation = strings.ToLower(strings.TrimSpace(cfg.Operation))
	cfg.DBHost = strings.TrimSpace(cfg.DBHost)
	cfg.CollectionName = strings.TrimSpace(cfg.CollectionName)
	cfg.EmbeddingHost = strings.TrimSpace(cfg.EmbeddingHost)
	cfg.EmbeddingModel = strings.TrimSpace(cfg.EmbeddingModel)
	cfg.Input = strings.TrimSpace(cfg.Input)
	cfg.Payload = strings.TrimSpace(cfg.Payload)
}

func validateConfig(cfg *Config) error {
	if cfg.Operation == "" {
		return fmt.Errorf("operation パラメータは必須です")
	}

	if !isSupportedOperation(cfg.Operation) {
		return fmt.Errorf("未対応のoperationです: %s", cfg.Operation)
	}

	if err := validatePort(cfg.DBPort, "db-port"); err != nil {
		return err
	}

	switch cfg.Operation {
	case OperationCreateCollection:
		if cfg.CollectionName == "" {
			return fmt.Errorf("create-collection では collection-name が必要です")
		}
		if cfg.Size <= 0 {
			return fmt.Errorf("size パラメータは 1 以上を指定してください")
		}
	case OperationListCollections:
		// no-op
	case OperationDescribeCollection, OperationDeleteCollection:
		if cfg.CollectionName == "" {
			return fmt.Errorf("operation %s では collection-name が必要です", cfg.Operation)
		}
	case OperationUpsertTexts:
		if err := validateCollectionAndEmbedding(cfg); err != nil {
			return err
		}
		if cfg.Input == "" {
			return fmt.Errorf("upsert-texts では input が必要です")
		}
		if err := validatePayloadFormat(cfg.Payload); err != nil {
			return err
		}
	case OperationQueryPoints:
		if err := validateCollectionAndEmbedding(cfg); err != nil {
			return err
		}
		if cfg.Input == "" {
			return fmt.Errorf("query-points では input が必要です")
		}
		if cfg.QueryLimit <= 0 {
			return fmt.Errorf("limit パラメータは 1 以上を指定してください")
		}
		if err := validatePayloadFormat(cfg.Payload); err != nil {
			return err
		}
	case OperationOverwritePayload:
		if cfg.CollectionName == "" {
			return fmt.Errorf("overwrite-payload では collection-name が必要です")
		}
		if strings.TrimSpace(cfg.Payload) == "" {
			return fmt.Errorf("overwrite-payload では payload が必要です")
		}
		if err := validatePayloadFormat(cfg.Payload); err != nil {
			return err
		}
	}

	if cfg.FilterMinShouldCount < 0 {
		return fmt.Errorf("filter-min-should は 0 以上を指定してください")
	}
	if cfg.Operation == OperationOverwritePayload && cfg.FilterMinShouldCount > 0 && len(cfg.FilterShould) == 0 {
		return fmt.Errorf("filter-min-should を指定する場合は filter-should も少なくとも1件必要です")
	}

	return nil
}

func validateCollectionAndEmbedding(cfg *Config) error {
	if cfg.CollectionName == "" {
		return fmt.Errorf("operation %s では collection-name が必要です", cfg.Operation)
	}
	if cfg.EmbeddingModel == "" {
		return fmt.Errorf("operation %s では embedding-model が必要です", cfg.Operation)
	}
	if err := validatePort(cfg.EmbeddingPort, "embedding-port"); err != nil {
		return err
	}
	if cfg.EmbeddingHost == "" {
		return fmt.Errorf("embedding-host は必須です")
	}
	return nil
}

func validatePort(port int, name string) error {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("%s の値が不正です: %d", name, port)
	}
	return nil
}

func validatePayloadFormat(payload string) error {
	if payload == "" {
		return nil
	}
	if !strings.Contains(payload, "=") {
		return fmt.Errorf("payload は key=value 形式で指定してください")
	}
	parts := strings.SplitN(payload, "=", 2)
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if key == "" || value == "" {
		return fmt.Errorf("payload は空文字を含められません")
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

// PrintUsage は CLI の利用方法を標準エラーに出力する。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Qdrant CLI\n\n")
	fmt.Fprintf(os.Stderr, "共通オプション:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n        実行する操作 (%s)\n", strings.Join(supportedOperations, ", "))
	fmt.Fprintf(os.Stderr, "  -db-host string\n        Qdrant のホスト名 (デフォルト: %s)\n", defaultDBHost)
	fmt.Fprintf(os.Stderr, "  -db-port int\n        Qdrant のポート番号 (デフォルト: %d)\n", defaultDBPort)
	fmt.Fprintf(os.Stderr, "  -collection-name string\n        対象のコレクション名\n")
	fmt.Fprintf(os.Stderr, "  -help\n        このヘルプを表示\n\n")

	fmt.Fprintf(os.Stderr, "create-collection 用:\n")
	fmt.Fprintf(os.Stderr, "  -size int\n        ベクトル次元 (デフォルト: %d)\n\n", defaultCollectionSize)

	fmt.Fprintf(os.Stderr, "describe-collection/delete-collection 用:\n")
	fmt.Fprintf(os.Stderr, "  -collection-name string\n        対象コレクション名 (必須)\n\n")

	fmt.Fprintf(os.Stderr, "upsert-texts/query-points 用:\n")
	fmt.Fprintf(os.Stderr, "  -embedding-host string (デフォルト: %s)\n", defaultEmbeddingHost)
	fmt.Fprintf(os.Stderr, "  -embedding-port int (デフォルト: %d)\n", defaultEmbeddingPort)
	fmt.Fprintf(os.Stderr, "  -embedding-model string\n        使用する埋め込みモデル\n")
	fmt.Fprintf(os.Stderr, "  -input string\n        埋め込み対象のテキスト (単一指定)\n")
	fmt.Fprintf(os.Stderr, "  -payload string\n        payload 条件 (key=value、単一指定)\n")
	fmt.Fprintf(os.Stderr, "  -limit int\n        query-points の取得件数 (デフォルト: %d)\n", defaultQueryLimit)

	fmt.Fprintf(os.Stderr, "overwrite-payload 用:\n")
	fmt.Fprintf(os.Stderr, "  -payload string\n        上書きする payload (key=value、単一指定)\n")
	fmt.Fprintf(os.Stderr, "  -filter-must key=value\n        上書き対象を絞り込む must 条件 (複数指定可)\n")
	fmt.Fprintf(os.Stderr, "  -filter-must-not key=value\n        上書き対象から除外する条件 (複数指定可)\n")
	fmt.Fprintf(os.Stderr, "  -filter-should key=value\n        should 条件 (複数指定可)\n")
	fmt.Fprintf(os.Stderr, "  -filter-min-should int\n        should 条件のうち満たすべき最小件数\n")
}

// KeyValue は key=value 形式の値を保持する。
type KeyValue struct {
	Key   string
	Value string
}

type keyValueList struct {
	target *[]KeyValue
}

func (l *keyValueList) String() string {
	if l == nil || l.target == nil {
		return ""
	}
	parts := make([]string, 0, len(*l.target))
	for _, kv := range *l.target {
		parts = append(parts, fmt.Sprintf("%s=%s", kv.Key, kv.Value))
	}
	return strings.Join(parts, ",")
}

func (l *keyValueList) Set(value string) error {
	if l == nil || l.target == nil {
		return fmt.Errorf("内部エラー: keyValueList が初期化されていません")
	}
	key, val, err := parseKeyValuePair(value)
	if err != nil {
		return err
	}
	*l.target = append(*l.target, KeyValue{Key: key, Value: val})
	return nil
}

func parseKeyValuePair(raw string) (string, string, error) {
	parts := strings.SplitN(strings.TrimSpace(raw), "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("key=value 形式で指定してください: %s", raw)
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if key == "" || value == "" {
		return "", "", fmt.Errorf("key/value に空文字は指定できません: %s", raw)
	}
	return key, value, nil
}

// SupportedOperations はサポートされている operation を返す。
func SupportedOperations() []string {
	out := make([]string, len(supportedOperations))
	copy(out, supportedOperations)
	sort.Strings(out)
	return out
}

// DumpDefaults はテスト向けに現在のデフォルト値を返す。
func DumpDefaults() map[string]string {
	return map[string]string{
		"db_host":        defaultDBHost,
		"db_port":        strconv.Itoa(defaultDBPort),
		"embedding_host": defaultEmbeddingHost,
	}
}
