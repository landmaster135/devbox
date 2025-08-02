package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config はArXiv CLIの設定を保持する構造体
type Config struct {
	Operation   string   // 操作タイプ (search, get_by_id)
	SearchQuery string   // 検索クエリ
	IdList      []string // arXiv IDリスト
	Start       int      // 開始位置（0ベース）
	MaxResults  int      // 最大結果数
	SortBy      string   // ソート基準 (relevance, lastUpdatedDate, submittedDate)
	SortOrder   string   // ソート順 (ascending, descending)
	Help        bool     // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(operation, searchQuery string, idList []string, start, maxResults int, sortBy, sortOrder string) (*Config, error) {
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	// 操作タイプの検証
	validOperations := []string{"search", "get_by_id"}
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

	// 操作タイプ別の検証
	switch operation {
	case "search":
		if searchQuery == "" {
			return nil, fmt.Errorf("search操作には検索クエリが必要です")
		}
	case "get_by_id":
		if len(idList) == 0 {
			return nil, fmt.Errorf("get_by_id操作にはIDリストが必要です")
		}
	}

	// パラメータの検証
	if start < 0 {
		return nil, fmt.Errorf("開始位置は0以上である必要があります")
	}
	if maxResults < 0 {
		return nil, fmt.Errorf("最大結果数は0以上である必要があります")
	}
	if maxResults > 30000 {
		return nil, fmt.Errorf("最大結果数は30000以下である必要があります")
	}

	// ソート基準の検証
	if sortBy != "" {
		validSortBy := []string{"relevance", "lastUpdatedDate", "submittedDate"}
		isValidSort := false
		for _, sort := range validSortBy {
			if sortBy == sort {
				isValidSort = true
				break
			}
		}
		if !isValidSort {
			return nil, fmt.Errorf("無効なソート基準です: %s", sortBy)
		}
	}

	// ソート順の検証
	if sortOrder != "" {
		validSortOrder := []string{"ascending", "descending"}
		isValidOrder := false
		for _, order := range validSortOrder {
			if sortOrder == order {
				isValidOrder = true
				break
			}
		}
		if !isValidOrder {
			return nil, fmt.Errorf("無効なソート順です: %s", sortOrder)
		}
	}

	return &Config{
		Operation:   operation,
		SearchQuery: searchQuery,
		IdList:      idList,
		Start:       start,
		MaxResults:  maxResults,
		SortBy:      sortBy,
		SortOrder:   sortOrder,
	}, nil
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation    = ""
		searchQuery  = ""
		idList       = ""
		startStr     = "0"
		maxResultsStr = "10"
		sortBy       = ""
		sortOrder    = ""
		help         = false
	)

	parser.StringVar(&operation, "operation", operation, "操作タイプ (search, get_by_id)")
	parser.StringVar(&operation, "o", operation, "操作タイプの短縮形")

	// 検索用のパラメータ
	parser.StringVar(&searchQuery, "query", searchQuery, "検索クエリ")
	parser.StringVar(&searchQuery, "q", searchQuery, "検索クエリの短縮形")

	// ID指定用のパラメータ
	parser.StringVar(&idList, "ids", idList, "カンマ区切りのarXiv IDリスト")
	parser.StringVar(&idList, "i", idList, "IDリストの短縮形")

	// ページング用のパラメータ
	parser.StringVar(&startStr, "start", startStr, "開始位置（0ベース）")
	parser.StringVar(&startStr, "s", startStr, "開始位置の短縮形")
	parser.StringVar(&maxResultsStr, "max_results", maxResultsStr, "最大結果数")
	parser.StringVar(&maxResultsStr, "m", maxResultsStr, "最大結果数の短縮形")

	// ソート用のパラメータ
	parser.StringVar(&sortBy, "sort_by", sortBy, "ソート基準 (relevance, lastUpdatedDate, submittedDate)")
	parser.StringVar(&sortOrder, "sort_order", sortOrder, "ソート順 (ascending, descending)")

	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	// 文字列から数値に変換
	start, err := strconv.Atoi(startStr)
	if err != nil {
		return nil, fmt.Errorf("無効な開始位置です: %s", startStr)
	}

	maxResults, err := strconv.Atoi(maxResultsStr)
	if err != nil {
		return nil, fmt.Errorf("無効な最大結果数です: %s", maxResultsStr)
	}

	// IDリストを配列に変換
	var idListArray []string
	if idList != "" {
		parts := strings.Split(idList, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				idListArray = append(idListArray, part)
			}
		}
	}

	return NewConfig(operation, searchQuery, idListArray, start, maxResults, sortBy, sortOrder)
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `ArXiv論文検索CLIツール

使用方法:
  論文検索:
    %s -operation search -query "quantum computing" -max_results 5
    %s -o search -q "au:einstein" -m 10

  ID指定取得:
    %s -operation get_by_id -ids "2301.00001,2301.00002"
    %s -o get_by_id -i "1234.5678"

  高度な検索:
    %s -operation search -query "ti:machine learning AND cat:cs.AI" -sort_by lastUpdatedDate -sort_order descending

オプション:
  -operation, -o    操作タイプ (search, get_by_id)
  -query, -q        検索クエリ (search操作用)
  -ids, -i          カンマ区切りのarXiv IDリスト (get_by_id操作用)
  -start, -s        開始位置（0ベース、デフォルト: 0）
  -max_results, -m  最大結果数（デフォルト: 10、最大: 30000）
  -sort_by          ソート基準 (relevance, lastUpdatedDate, submittedDate)
  -sort_order       ソート順 (ascending, descending)
  -help, -h         このヘルプを表示

検索クエリの例:
  all:electron                    - 全フィールドで"electron"を検索
  ti:"quantum computing"          - タイトルで"quantum computing"を検索
  au:einstein                     - 著者で"einstein"を検索
  abs:machine learning            - 要約で"machine learning"を検索
  cat:cs.AI                       - カテゴリで"cs.AI"を検索
  ti:quantum AND au:feynman       - タイトルに"quantum"かつ著者に"feynman"
  ti:physics ANDNOT cat:hep-th    - タイトルに"physics"だがカテゴリが"hep-th"でない

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
