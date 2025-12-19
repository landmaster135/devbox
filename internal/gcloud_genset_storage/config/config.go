package config

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	OperationUploadFiles   = "upload-files"
	OperationDownloadFiles = "download-files"
	OperationCreateBucket  = "create-bucket"
	OperationListContents  = "list-contents"
	OperationShowDetails   = "show-details"
	OperationListNames     = "list-names"
	OperationDeleteObject  = "delete-object"
	OperationGetACL        = "get-acl"
	OperationSetACL        = "set-acl"
	OperationGrantReadAll  = "grant-read-all"
	OperationRemoveReadAll = "remove-read-all"
)

// Config holds CLI parameters for gcloud genset storage operations.
type Config struct {
	Operation    string
	LocalPath    string
	BucketURL    string
	Sources      []string
	Destination  string
	StorageClass string
	Location     string
	Target       string
	ACLFile      string
	Help         bool
	rawSources   string
}

var validOperations = []string{
	OperationUploadFiles,
	OperationDownloadFiles,
	OperationCreateBucket,
	OperationListContents,
	OperationShowDetails,
	OperationListNames,
	OperationDeleteObject,
	OperationGetACL,
	OperationSetACL,
	OperationGrantReadAll,
	OperationRemoveReadAll,
}

// ParseFlags parses CLI flags using the standard parser.
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser parses CLI flags using the provided parser.
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{}

	parser.StringVar(&cfg.Operation, "operation", "", fmt.Sprintf("実行する操作 (%s)", strings.Join(validOperations, ", ")))
	parser.StringVar(&cfg.LocalPath, "local-path", "", "ローカルディレクトリへのパス")
	parser.StringVar(&cfg.BucketURL, "bucket-url", "", "対象のGCSバケットURL (gs://...)")
	parser.StringVar(&cfg.rawSources, "sources", "", "ダウンロード対象のGCSオブジェクト（カンマ区切り）")
	parser.StringVar(&cfg.Destination, "destination", "", "ダウンロード先ローカルディレクトリ")
	parser.StringVar(&cfg.StorageClass, "storage-class", "", "GCSバケットのストレージクラス")
	parser.StringVar(&cfg.Location, "location", "", "GCSバケットのロケーション")
	parser.StringVar(&cfg.Target, "target", "", "対象のパス (gs://... またはローカルパス)")
	parser.StringVar(&cfg.ACLFile, "acl-file", "", "ACL設定に使用するファイル")
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示する")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if cfg.Help {
		return cfg, nil
	}

	if err := parseSources(cfg); err != nil {
		return nil, err
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseSources(cfg *Config) error {
	if strings.TrimSpace(cfg.rawSources) == "" {
		cfg.Sources = nil
		return nil
	}
	parts := strings.Split(cfg.rawSources, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			cfg.Sources = append(cfg.Sources, trimmed)
		}
	}
	if len(cfg.Sources) == 0 {
		return fmt.Errorf("sources の形式が不正です")
	}
	return nil
}

func validateConfig(cfg *Config) error {
	if strings.TrimSpace(cfg.Operation) == "" {
		return fmt.Errorf("operation は必須です")
	}
	if !isValidOperation(cfg.Operation) {
		return fmt.Errorf("未対応のoperationです: %s", cfg.Operation)
	}

	switch cfg.Operation {
	case OperationUploadFiles:
		if strings.TrimSpace(cfg.LocalPath) == "" {
			return fmt.Errorf("local-path は必須です")
		}
		if strings.TrimSpace(cfg.BucketURL) == "" {
			return fmt.Errorf("bucket-url は必須です")
		}
	case OperationDownloadFiles:
		if len(cfg.Sources) == 0 {
			return fmt.Errorf("sources は必須です")
		}
		if strings.TrimSpace(cfg.Destination) == "" {
			return fmt.Errorf("destination は必須です")
		}
	case OperationCreateBucket:
		if strings.TrimSpace(cfg.BucketURL) == "" {
			return fmt.Errorf("bucket-url は必須です")
		}
		if strings.TrimSpace(cfg.StorageClass) == "" {
			return fmt.Errorf("storage-class は必須です")
		}
		if strings.TrimSpace(cfg.Location) == "" {
			return fmt.Errorf("location は必須です")
		}
	case OperationListContents:
		if strings.TrimSpace(cfg.Target) == "" {
			return fmt.Errorf("target は必須です")
		}
	case OperationShowDetails, OperationListNames, OperationDeleteObject, OperationGetACL, OperationGrantReadAll, OperationRemoveReadAll:
		if err := requireGCSPath(cfg.Target); err != nil {
			return err
		}
	case OperationSetACL:
		if strings.TrimSpace(cfg.ACLFile) == "" {
			return fmt.Errorf("acl-file は必須です")
		}
		if err := requireGCSPath(cfg.Target); err != nil {
			return err
		}
	}

	return nil
}

func isValidOperation(op string) bool {
	for _, candidate := range validOperations {
		if candidate == op {
			return true
		}
	}
	return false
}

func requireGCSPath(target string) error {
	if strings.TrimSpace(target) == "" {
		return fmt.Errorf("target は必須です")
	}
	if !strings.HasPrefix(target, "gs://") {
		return fmt.Errorf("target には gs:// で始まるパスを指定してください")
	}
	return nil
}

// PrintUsage displays CLI usage information.
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Google Cloud Storage 用の gcloud/gsutil コマンド生成ツール\n\n")

	fmt.Fprintf(os.Stderr, "共通パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n        実行する操作 (%s)\n", strings.Join(validOperations, ", "))
	fmt.Fprintf(os.Stderr, "  -help\n        このヘルプを表示\n\n")

	fmt.Fprintf(os.Stderr, "operation ごとの主なパラメータ:\n")
	fmt.Fprintf(os.Stderr, "  upload-files           : -local-path, -bucket-url\n")
	fmt.Fprintf(os.Stderr, "  download-files         : -sources, -destination\n")
	fmt.Fprintf(os.Stderr, "  create-bucket          : -bucket-url, -storage-class, -location\n")
	fmt.Fprintf(os.Stderr, "  list-contents          : -target\n")
	fmt.Fprintf(os.Stderr, "  show-details           : -target (gs://)\n")
	fmt.Fprintf(os.Stderr, "  list-names             : -target (gs://)\n")
	fmt.Fprintf(os.Stderr, "  delete-object          : -target (gs://)\n")
	fmt.Fprintf(os.Stderr, "  get-acl                : -target (gs://)\n")
	fmt.Fprintf(os.Stderr, "  set-acl                : -acl-file, -target (gs://)\n")
	fmt.Fprintf(os.Stderr, "  grant-read-all         : -target (gs://)\n")
	fmt.Fprintf(os.Stderr, "  remove-read-all        : -target (gs://)\n")
}

func init() {
	sort.Strings(validOperations)
}
