package config

import (
	"fmt"
	"os"
)

// Config はメモリCLIの設定を保持する構造体
type Config struct {
	Operation     string // 操作タイプ (create-entities, create-relations, add-observations, delete-entities, delete-observations, delete-relations, read-graph, search-nodes, open-nodes)
	Entities      string // JSON形式のエンティティ配列
	Relations     string // JSON形式のリレーション配列
	Observations  string // JSON形式の観察事項配列
	Query         string // 検索クエリ
	Names         string // カンマ区切りの名前リスト
	EntityNames   string // 削除対象エンティティ名
	Deletions     string // JSON形式の削除対象
	MemoryFile    string // メモリファイルパス
	Help          bool   // ヘルプ表示フラグ
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

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation     = ""
		entities      = ""
		relations     = ""
		observations  = ""
		query         = ""
		names         = ""
		entityNames   = ""
		deletions     = ""
		memoryFile    = "./memory.json"
		help          = false
	)

	parser.StringVar(&operation, "operation", "", "メモリ操作 (create-entities, create-relations, add-observations, delete-entities, delete-observations, delete-relations, read-graph, search-nodes, open-nodes)")
	parser.StringVar(&operation, "o", "", "操作の短縮形")

	// エンティティ関連のパラメータ
	parser.StringVar(&entities, "entities", "", "JSON形式のエンティティ配列")
	parser.StringVar(&entities, "e", "", "エンティティの短縮形")

	// リレーション関連のパラメータ
	parser.StringVar(&relations, "relations", "", "JSON形式のリレーション配列")
	parser.StringVar(&relations, "r", "", "リレーションの短縮形")

	// 観察事項関連のパラメータ
	parser.StringVar(&observations, "observations", "", "JSON形式の観察事項配列")
	parser.StringVar(&observations, "obs", "", "観察事項の短縮形")

	// 検索関連のパラメータ
	parser.StringVar(&query, "query", "", "検索クエリ")
	parser.StringVar(&query, "q", "", "クエリの短縮形")

	// ノード名関連のパラメータ
	parser.StringVar(&names, "names", "", "カンマ区切りの名前リスト")
	parser.StringVar(&names, "n", "", "名前の短縮形")

	// 削除関連のパラメータ
	parser.StringVar(&entityNames, "entity-names", "", "削除対象エンティティ名（カンマ区切り）")
	parser.StringVar(&entityNames, "en", "", "エンティティ名の短縮形")
	parser.StringVar(&deletions, "deletions", "", "JSON形式の削除対象")
	parser.StringVar(&deletions, "del", "", "削除対象の短縮形")

	// ファイル関連のパラメータ
	parser.StringVar(&memoryFile, "memory-file", "./memory.json", "メモリファイルパス")
	parser.StringVar(&memoryFile, "f", "./memory.json", "ファイルの短縮形")

	parser.BoolVar(&help, "help", false, "ヘルプを表示")
	parser.BoolVar(&help, "h", false, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	// 操作タイプが指定されていない場合のエラーチェック
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
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
  -help, -h            このヘルプを表示

エンティティのJSON形式例:
  [{"name":"John","entityType":"person","observations":["speaks English","likes coffee"]}]

リレーションのJSON形式例:
  [{"from":"John","to":"Company","relationType":"works_at"}]

観察事項のJSON形式例:
  [{"entityName":"John","contents":["new observation"]}]

削除対象のJSON形式例:
  [{"entityName":"John","observations":["old observation"]}]

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
