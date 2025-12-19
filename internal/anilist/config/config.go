package config

import (
	"fmt"
	"os"
	"slices"
	"strconv"
)

// Config はAniList CLIの設定を保持する構造体
type Config struct {
	Operation string // 操作タイプ (query-anime)
	Username  string // AniListユーザー名
	UserID    *int   // AniListユーザーID
	Format    string // 出力形式 (json, table)
	Limit     int    // 取得件数の制限
	Status    string // ステータスフィルタ
	OutputDir string // 出力ディレクトリ
	Help      bool   // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(operation, username string, userID *int, format string, limit int, status, outputDir string) (*Config, error) {
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	// 操作タイプの検証
	validOperations := []string{"query-anime", "query-manga"}
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

	// ユーザー名またはユーザーIDのいずれかが必要
	if username == "" && userID == nil {
		return nil, fmt.Errorf("ユーザー名またはユーザーIDのいずれかを指定してください")
	}

	// 出力形式の検証
	if format != "" {
		validFormats := []string{"json", "table"}
		isValidFormat := false
		for _, f := range validFormats {
			if format == f {
				isValidFormat = true
				break
			}
		}
		if !isValidFormat {
			return nil, fmt.Errorf("無効な出力形式です: %s", format)
		}
	}

	// 制限値の検証
	if limit < 0 {
		return nil, fmt.Errorf("制限値は0以上である必要があります")
	}

	// ステータスの検証
	if status != "" {
		validStatuses := []string{"CURRENT", "PLANNING", "COMPLETED", "DROPPED", "PAUSED", "REPEATING"}
		if !slices.Contains(validStatuses, status) {
			return nil, fmt.Errorf("無効なステータスです: %s", status)
		}
	}

	return &Config{
		Operation: operation,
		Username:  username,
		UserID:    userID,
		Format:    format,
		Limit:     limit,
		Status:    status,
		OutputDir: outputDir,
	}, nil
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation = ""
		username  = ""
		userIDStr = ""
		format    = "json"
		limitStr  = "0"
		status    = ""
		outputDir = ""
		help      = false
	)

	parser.StringVar(&operation, "operation", operation, "操作タイプ (query-anime, query-manga)")
	parser.StringVar(&operation, "o", operation, "操作タイプの短縮形")

	parser.StringVar(&username, "username", username, "AniListユーザー名")
	parser.StringVar(&username, "u", username, "ユーザー名の短縮形")

	parser.StringVar(&userIDStr, "user-id", userIDStr, "AniListユーザーID")
	parser.StringVar(&userIDStr, "i", userIDStr, "ユーザーIDの短縮形")

	parser.StringVar(&format, "format", format, "出力形式 (json, table)")
	parser.StringVar(&format, "f", format, "出力形式の短縮形")

	parser.StringVar(&limitStr, "limit", limitStr, "取得件数の制限 (0は無制限)")
	parser.StringVar(&limitStr, "l", limitStr, "制限の短縮形")

	parser.StringVar(&status, "status", status, "ステータスフィルタ (CURRENT, PLANNING, COMPLETED, DROPPED, PAUSED, REPEATING)")
	parser.StringVar(&status, "s", status, "ステータスの短縮形")

	parser.StringVar(&outputDir, "output-dir", outputDir, "出力ディレクトリ")
	parser.StringVar(&outputDir, "d", outputDir, "出力ディレクトリの短縮形")

	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	// ユーザーIDの変換
	var userID *int
	if userIDStr != "" {
		id, err := strconv.Atoi(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("無効なユーザーIDです: %s", userIDStr)
		}
		userID = &id
	}

	// 制限値の変換
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, fmt.Errorf("無効な制限値です: %s", limitStr)
	}

	// 残りの引数から操作タイプを取得（位置引数として）
	args := parser.Args()
	if len(args) >= 1 && operation == "" {
		operation = args[0]
	}

	return NewConfig(operation, username, userID, format, limit, status, outputDir)
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `AniList CLI ツール

使用方法:
  アニメ情報取得（ユーザー名指定）:
    %s -operation query-anime -username your_username
    %s -o query-anime -u your_username

  アニメ情報取得（ユーザーID指定）:
    %s -operation query-anime -user-id 123456
    %s -o query-anime -i 123456

  出力形式指定:
    %s -o query-anime -u your_username -format table
    %s -o query-anime -u your_username -f json

  ステータスフィルタ:
    %s -o query-anime -u your_username -status COMPLETED
    %s -o query-anime -u your_username -s CURRENT

  取得件数制限:
    %s -o query-anime -u your_username -limit 10
    %s -o query-anime -u your_username -l 5

  ファイル出力:
    %s -o query-anime -u your_username -output-dir ./output
    %s -o query-anime -u your_username -d ./results

  マンガ情報取得（ユーザー名指定）:
    %s -operation query-manga -username your_username
    %s -o query-manga -u your_username

  マンガ情報取得（テーブル形式）:
    %s -o query-manga -u your_username -format table
    %s -o query-manga -u your_username -f table

  マンガ情報取得（完了済みのみ）:
    %s -o query-manga -u your_username -status COMPLETED
    %s -o query-manga -u your_username -s COMPLETED

オプション:
  -operation, -o    操作タイプ (query-anime, query-manga)
  -username, -u     AniListユーザー名
  -user-id, -i      AniListユーザーID
  -format, -f       出力形式 (json, table) [デフォルト: json]
  -limit, -l        取得件数の制限 (0は無制限) [デフォルト: 0]
  -status, -s       ステータスフィルタ (CURRENT, PLANNING, COMPLETED, DROPPED, PAUSED, REPEATING)
  -output-dir, -d   出力ディレクトリ (指定時はファイルに保存)
  -help, -h         このヘルプを表示

注意:
  - ユーザー名またはユーザーIDのいずれかを指定してください
  - 環境変数は使用しません

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
