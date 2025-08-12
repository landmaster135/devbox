package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileSystem はファイルシステム操作のインターフェース
type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
}

// RealFileSystem は実際のファイルシステム操作を行う
type RealFileSystem struct{}

func (fs *RealFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (fs *RealFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// Config はメモリCLIの設定を保持する構造体
type Config struct {
	Operation    string // 操作タイプ (create-entities, create-relations, add-observations, delete-entities, delete-observations, delete-relations, read-graph, search-nodes, open-nodes)
	Entities     string // JSON形式のエンティティ配列
	Relations    string // JSON形式のリレーション配列
	Observations string // JSON形式の観察事項配列
	Query        string // 検索クエリ
	Names        string // カンマ区切りの名前リスト
	EntityNames  string // 削除対象エンティティ名
	Deletions    string // JSON形式の削除対象
	StorageType  string // ストレージタイプ (file, valkey)
	Help         bool   // ヘルプ表示フラグ

	// File
	MemoryFile string // メモリファイルパス

	// Valkey
	ValkeyHost     string // Valkeyホスト
	ValkeyPort     int    // Valkeyポート
	ValkeyPassword string // Valkeyパスワード
	ValkeyDatabase int    // データベース番号
	ValkeyKey      string // Valkeyキー
}

// NewConfig は新しいConfigを作成する
func NewConfig(operation string) (*Config, error) {
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	// 操作タイプの検証
	validOperations := []string{
		"create-entities", "create-relations", "add-observations",
		"delete-entities", "delete-observations", "delete-relations",
		"read-graph", "search-nodes", "open-nodes",
	}
	isValid := false
	for _, op := range validOperations {
		if operation == op {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, fmt.Errorf("無効な操作タイプです: %s", operation)
	}

	return &Config{
		Operation:  operation,
		MemoryFile: "./memory.json", // デフォルトパス
	}, nil
}

// BuildValkeyURL はConfigからValkey接続URLを構築する
func (c *Config) BuildValkeyURL() string {
	if c.ValkeyPassword != "" {
		// パスワードありの場合
		return fmt.Sprintf("valkey://:%s@%s:%d/%d",
			c.ValkeyPassword, c.ValkeyHost, c.ValkeyPort, c.ValkeyDatabase)
	}
	// パスワードなしの場合
	return fmt.Sprintf("valkey://%s:%d/%d",
		c.ValkeyHost, c.ValkeyPort, c.ValkeyDatabase)
}

// validateOperationParameters は操作別の必須パラメータをチェックする
func validateOperationParameters(cfg *Config) error {
	switch cfg.Operation {
	case "create-entities":
		if strings.TrimSpace(cfg.Entities) == "" {
			return fmt.Errorf("create-entities操作には-entitiesパラメータが必須です")
		}
	case "create-relations":
		if strings.TrimSpace(cfg.Relations) == "" {
			return fmt.Errorf("create-relations操作には-relationsパラメータが必須です")
		}
	case "add-observations":
		if strings.TrimSpace(cfg.Observations) == "" {
			return fmt.Errorf("add-observations操作には-observationsパラメータが必須です")
		}
	case "delete-entities":
		if strings.TrimSpace(cfg.EntityNames) == "" {
			return fmt.Errorf("delete-entities操作には-entity-namesパラメータが必須です")
		}
	case "delete-observations":
		if strings.TrimSpace(cfg.Deletions) == "" {
			return fmt.Errorf("delete-observations操作には-deletionsパラメータが必須です")
		}
	case "delete-relations":
		if strings.TrimSpace(cfg.Relations) == "" {
			return fmt.Errorf("delete-relations操作には-relationsパラメータが必須です")
		}
	case "search-nodes":
		if strings.TrimSpace(cfg.Query) == "" {
			return fmt.Errorf("search-nodes操作には-queryパラメータが必須です")
		}
	case "open-nodes":
		if strings.TrimSpace(cfg.Names) == "" {
			return fmt.Errorf("open-nodes操作には-namesパラメータが必須です")
		}
	case "read-graph":
		// read-graphは追加パラメータ不要
	default:
		return fmt.Errorf("未知の操作タイプです: %s", cfg.Operation)
	}
	return nil
}

// validateJSONParameters はJSONパラメータの形式をチェックする
func validateJSONParameters(cfg *Config) error {
	// entitiesのJSON検証
	if cfg.Entities != "" {
		var entities []any
		if err := json.Unmarshal([]byte(cfg.Entities), &entities); err != nil {
			return fmt.Errorf("entitiesパラメータが無効なJSON形式です: %v", err)
		}
		if len(entities) == 0 {
			return fmt.Errorf("entitiesパラメータは空の配列にできません")
		}
	}

	// relationsのJSON検証
	if cfg.Relations != "" {
		var relations []any
		if err := json.Unmarshal([]byte(cfg.Relations), &relations); err != nil {
			return fmt.Errorf("relationsパラメータが無効なJSON形式です: %v", err)
		}
		if len(relations) == 0 {
			return fmt.Errorf("relationsパラメータは空の配列にできません")
		}
	}

	// observationsのJSON検証
	if cfg.Observations != "" {
		var observations []any
		if err := json.Unmarshal([]byte(cfg.Observations), &observations); err != nil {
			return fmt.Errorf("observationsパラメータが無効なJSON形式です: %v", err)
		}
		if len(observations) == 0 {
			return fmt.Errorf("observationsパラメータは空の配列にできません")
		}
	}

	// deletionsのJSON検証
	if cfg.Deletions != "" {
		var deletions []any
		if err := json.Unmarshal([]byte(cfg.Deletions), &deletions); err != nil {
			return fmt.Errorf("deletionsパラメータが無効なJSON形式です: %v", err)
		}
		if len(deletions) == 0 {
			return fmt.Errorf("deletionsパラメータは空の配列にできません")
		}
	}

	return nil
}

// validateValkeyParameters はValkeyパラメータをチェックする
func validateValkeyParameters(cfg *Config) error {
	// ホストの検証
	if strings.TrimSpace(cfg.ValkeyHost) == "" {
		return fmt.Errorf("Valkeyホストが空です")
	}

	// ポートの範囲チェック
	if cfg.ValkeyPort < 1 || cfg.ValkeyPort > 65535 {
		return fmt.Errorf("Valkeyポートは1-65535の範囲で指定してください: %d", cfg.ValkeyPort)
	}

	// データベース番号の範囲チェック
	if cfg.ValkeyDatabase < 0 || cfg.ValkeyDatabase > 15 {
		return fmt.Errorf("Valkeyデータベース番号は0-15の範囲で指定してください: %d", cfg.ValkeyDatabase)
	}

	// キーの検証
	if strings.TrimSpace(cfg.ValkeyKey) == "" {
		return fmt.Errorf("Valkeyキーが空です")
	}

	return nil
}

// validateFileParametersWithFS はFileSystemインターフェースを使用してファイルパラメータをチェックする（テスト用）
func validateFileParametersWithFS(cfg *Config, fs FileSystem) error {
	// ファイルパスの検証
	if strings.TrimSpace(cfg.MemoryFile) == "" {
		return fmt.Errorf("メモリファイルパスが空です")
	}

	// ディレクトリの存在確認と作成
	dir := filepath.Dir(cfg.MemoryFile)
	if dir != "." && dir != "" {
		if _, err := fs.Stat(dir); os.IsNotExist(err) {
			if err := fs.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("ディレクトリの作成に失敗しました: %v", err)
			}
		}
	}

	return nil
}

// validateFileParameters はファイルパラメータをチェックする
func validateFileParameters(cfg *Config) error {
	return validateFileParametersWithFS(cfg, &RealFileSystem{})
}

func validateConfig(cfg *Config) error {
	// ストレージタイプが指定されていない場合のエラーチェック
	if cfg.StorageType == "" {
		return fmt.Errorf("ストレージタイプが指定されていません。-storage-type file または -storage-type valkey を指定してください")
	}

	// ストレージタイプの検証
	if cfg.StorageType != "file" && cfg.StorageType != "valkey" {
		return fmt.Errorf("無効なストレージタイプです: %s (file または valkey を指定してください)", cfg.StorageType)
	}

	// 操作別の必須パラメータチェック
	if err := validateOperationParameters(cfg); err != nil {
		return err
	}

	// JSONパラメータの形式検証
	if err := validateJSONParameters(cfg); err != nil {
		return err
	}

	// Valkeyパラメータの検証
	if cfg.StorageType == "valkey" {
		if err := validateValkeyParameters(cfg); err != nil {
			return err
		}
	}

	// ファイルパスの検証
	if cfg.StorageType == "file" {
		if err := validateFileParameters(cfg); err != nil {
			return err
		}
	}

	return nil
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation      = ""
		entities       = ""
		relations      = ""
		observations   = ""
		query          = ""
		names          = ""
		entityNames    = ""
		deletions      = ""
		memoryFile     = "./memory.json"
		storageType    = ""
		valkeyHost     = "localhost"
		valkeyPort     = 6379
		valkeyPassword = ""
		valkeyDatabase = 0
		valkeyKey      = ""
		help           = false
	)

	parser.StringVar(&operation, "operation", operation, "メモリ操作 (create-entities, create-relations, add-observations, delete-entities, delete-observations, delete-relations, read-graph, search-nodes, open-nodes)")
	parser.StringVar(&operation, "o", operation, "操作の短縮形")

	// エンティティ関連のパラメータ
	parser.StringVar(&entities, "entities", entities, "JSON形式のエンティティ配列")
	parser.StringVar(&entities, "e", entities, "エンティティの短縮形")

	// リレーション関連のパラメータ
	parser.StringVar(&relations, "relations", relations, "JSON形式のリレーション配列")
	parser.StringVar(&relations, "r", relations, "リレーションの短縮形")

	// 観察事項関連のパラメータ
	parser.StringVar(&observations, "observations", observations, "JSON形式の観察事項配列")
	parser.StringVar(&observations, "obs", observations, "観察事項の短縮形")

	// 検索関連のパラメータ
	parser.StringVar(&query, "query", query, "検索クエリ")
	parser.StringVar(&query, "q", query, "クエリの短縮形")

	// ノード名関連のパラメータ
	parser.StringVar(&names, "names", names, "カンマ区切りの名前リスト")
	parser.StringVar(&names, "n", names, "名前の短縮形")

	// 削除関連のパラメータ
	parser.StringVar(&entityNames, "entity-names", entityNames, "削除対象エンティティ名（カンマ区切り）")
	parser.StringVar(&entityNames, "en", entityNames, "エンティティ名の短縮形")
	parser.StringVar(&deletions, "deletions", deletions, "JSON形式の削除対象")
	parser.StringVar(&deletions, "del", deletions, "削除対象の短縮形")

	// ファイル関連のパラメータ
	parser.StringVar(&memoryFile, "memory-file", memoryFile, "メモリファイルパス")
	parser.StringVar(&memoryFile, "f", memoryFile, "ファイルの短縮形")

	// ストレージ関連のパラメータ
	parser.StringVar(&storageType, "storage-type", storageType, "ストレージタイプ (file, valkey)")
	parser.StringVar(&storageType, "s", storageType, "ストレージタイプの短縮形")

	// Valkey接続関連のパラメータ
	parser.StringVar(&valkeyHost, "valkey-host", valkeyHost, "Valkeyホスト")
	parser.StringVar(&valkeyHost, "vh", valkeyHost, "Valkeyホストの短縮形")
	parser.IntVar(&valkeyPort, "valkey-port", valkeyPort, "Valkeyポート")
	parser.IntVar(&valkeyPort, "vp", valkeyPort, "Valkeyポートの短縮形")
	parser.StringVar(&valkeyPassword, "valkey-password", valkeyPassword, "Valkeyパスワード")
	parser.StringVar(&valkeyPassword, "vpass", valkeyPassword, "Valkeyパスワードの短縮形")
	parser.IntVar(&valkeyDatabase, "valkey-database", valkeyDatabase, "データベース番号")
	parser.IntVar(&valkeyDatabase, "vdb", valkeyDatabase, "データベース番号の短縮形")
	parser.StringVar(&valkeyKey, "valkey-key", valkeyKey, "Valkeyキー")
	parser.StringVar(&valkeyKey, "vk", valkeyKey, "Valkeyキーの短縮形")

	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	config, err := NewConfig(operation)
	if err != nil {
		return nil, err
	}

	// パラメータを設定
	config.Entities = entities
	config.Relations = relations
	config.Observations = observations
	config.Query = query
	config.Names = names
	config.EntityNames = entityNames
	config.Deletions = deletions
	config.MemoryFile = memoryFile
	config.StorageType = storageType
	config.ValkeyHost = valkeyHost
	config.ValkeyPort = valkeyPort
	config.ValkeyPassword = valkeyPassword
	config.ValkeyDatabase = valkeyDatabase
	config.ValkeyKey = valkeyKey

	err = validateConfig(config); if err != nil {
		return nil, err
	}

	return config, nil
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `メモリCLIツール

使用方法:
  エンティティ作成:
    %s -operation create-entities -entities '[{"name":"John","entityType":"person","observations":["speaks English"]}]'

  リレーション作成:
    %s -operation create-relations -relations '[{"from":"John","to":"Company","relationType":"works_at"}]'

  観察事項追加:
    %s -operation add-observations -observations '[{"entityName":"John","contents":["likes coffee"]}]'

  グラフ読み取り:
    %s -operation read-graph

  ノード検索:
    %s -operation search-nodes -query "John"

  特定ノード取得:
    %s -operation open-nodes -names "John,Company"

  エンティティ削除:
    %s -operation delete-entities -entity-names "John,Company"

  観察事項削除:
    %s -operation delete-observations -deletions '[{"entityName":"John","observations":["old info"]}]'

  リレーション削除:
    %s -operation delete-relations -relations '[{"from":"John","to":"Company","relationType":"works_at"}]'

  短縮形:
    %s -o create-entities -e '[{"name":"John","entityType":"person","observations":["speaks English"]}]'
    %s -o search-nodes -q "John"
    %s -o open-nodes -n "John,Company"

  Valkey使用例:
    %s -operation read-graph -storage-type valkey
    %s -operation read-graph -storage-type valkey -valkey-host redis.example.com -valkey-port 6380
    %s -operation read-graph -storage-type valkey -valkey-password mypassword -valkey-database 1

オプション:
  -operation, -o       メモリ操作 (create-entities, create-relations, add-observations, delete-entities, delete-observations, delete-relations, read-graph, search-nodes, open-nodes)
  -entities, -e        JSON形式のエンティティ配列
  -relations, -r       JSON形式のリレーション配列
  -observations, -obs  JSON形式の観察事項配列
  -query, -q           検索クエリ
  -names, -n           カンマ区切りの名前リスト
  -entity-names, -en   削除対象エンティティ名（カンマ区切り）
  -deletions, -del     JSON形式の削除対象
  -memory-file, -f     メモリファイルパス (デフォルト: ./memory.json)
  -storage-type, -s    ストレージタイプ (file, valkey) (必須)
  -valkey-host, -vh    Valkeyホスト (デフォルト: localhost)
  -valkey-port, -vp    Valkeyポート (デフォルト: 6379)
  -valkey-password, -vpass  Valkeyパスワード
  -valkey-database, -vdb    データベース番号 (デフォルト: 0)
  -valkey-key, -vk     Valkeyキー (デフォルト: knowledge_graph:main)
  -help, -h            このヘルプを表示

エンティティのJSON形式例:
  [{"name":"John","entityType":"person","observations":["speaks English","likes coffee"]}]

リレーションのJSON形式例:
  [{"from":"John","to":"Company","relationType":"works_at"}]

観察事項のJSON形式例:
  [{"entityName":"John","contents":["new observation"]}]

削除対象のJSON形式例:
  [{"entityName":"John","observations":["old observation"]}]

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
