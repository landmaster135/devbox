package config

import (
	"fmt"
	"os"

	envLoader "github.com/landmaster135/devbox/internal/env_loader"
)

const envKeyGitHubToken = "GITHUB_PERSONAL_ACCESS_TOKEN"

// Config はGitHub CLIの設定を保持する構造体
type Config struct {
	Operation   string // 操作タイプ (list-issues)
	Token       string // GitHubトークン
	Owner       string // リポジトリオーナー
	Repo        string // リポジトリ名
	State       string // イシューの状態 (open, closed, all)
	Sort        string // ソート項目 (created, updated, comments)
	Direction   string // ソート方向 (asc, desc)
	PerPage     int    // ページあたりの件数
	Page        int    // ページ番号
	IssueNumber int    // イシュー番号（特定イシュー取得用）
	Help        bool   // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(operation, token, owner, repo, state, sort, direction string, perPage, page, issueNumber int) (*Config, error) {
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	// 操作タイプの検証
	validOperations := []string{"list-issues"}
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

	// 必須パラメータの検証
	if token == "" {
		return nil, fmt.Errorf("GitHubトークンが指定されていません")
	}
	if owner == "" {
		return nil, fmt.Errorf("リポジトリオーナーが指定されていません")
	}
	if repo == "" {
		return nil, fmt.Errorf("リポジトリ名が指定されていません")
	}

	// オプションパラメータのデフォルト値設定
	if perPage <= 0 {
		perPage = 30
	}
	if page <= 0 {
		page = 1
	}

	return &Config{
		Operation:   operation,
		Token:       token,
		Owner:       owner,
		Repo:        repo,
		State:       state,
		Sort:        sort,
		Direction:   direction,
		PerPage:     perPage,
		Page:        page,
		IssueNumber: issueNumber,
	}, nil
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation   = ""
		owner       = ""
		repo        = ""
		state       = ""
		sort        = ""
		direction   = ""
		perPage     = 30
		page        = 1
		issueNumber = 0
		help        = false
	)

	parser.StringVar(&operation, "operation", operation, "操作タイプ (list-issues)")
	parser.StringVar(&operation, "o", operation, "操作タイプの短縮形")

	parser.StringVar(&owner, "owner", owner, "リポジトリオーナー")
	parser.StringVar(&owner, "ow", owner, "リポジトリオーナーの短縮形")

	parser.StringVar(&repo, "repo", repo, "リポジトリ名")
	parser.StringVar(&repo, "r", repo, "リポジトリ名の短縮形")

	parser.StringVar(&state, "state", state, "イシューの状態 (open, closed, all)")
	parser.StringVar(&state, "s", state, "イシューの状態の短縮形")

	parser.StringVar(&sort, "sort", sort, "ソート項目 (created, updated, comments)")
	parser.StringVar(&sort, "so", sort, "ソート項目の短縮形")

	parser.StringVar(&direction, "direction", direction, "ソート方向 (asc, desc)")
	parser.StringVar(&direction, "d", direction, "ソート方向の短縮形")

	parser.IntVar(&perPage, "per-page", perPage, "ページあたりの件数")
	parser.IntVar(&perPage, "pp", perPage, "ページあたりの件数の短縮形")

	parser.IntVar(&page, "page", page, "ページ番号")
	parser.IntVar(&page, "p", page, "ページ番号の短縮形")

	parser.IntVar(&issueNumber, "issue-number", issueNumber, "イシュー番号（特定イシュー取得用）")
	parser.IntVar(&issueNumber, "in", issueNumber, "イシュー番号の短縮形")

	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	envValues, err := envLoader.Load([]string{envKeyGitHubToken})
	if err != nil {
		return nil, err
	}

	return NewConfig(operation, envValues[envKeyGitHubToken], owner, repo, state, sort, direction, perPage, page, issueNumber)
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `GitHub CLI ツール

使用方法:
  イシュー一覧取得:
    %s -operation list-issues -owner OWNER -repo REPO
    %s -o list-issues -ow OWNER -r REPO

  特定イシュー取得:
    %s -o list-issues -ow OWNER -r REPO -issue-number 123
    %s -o list-issues -ow OWNER -r REPO -in 123

  オプションパラメータ付き:
    %s -o list-issues -ow OWNER -r REPO -state open -sort created -direction desc
    %s -o list-issues -ow OWNER -r REPO -s closed -so updated -d asc -pp 50 -p 2

環境変数:
  GITHUB_PERSONAL_ACCESS_TOKEN GitHubトークン [必須]

オプション:
  -operation, -o     操作タイプ (list-issues) [必須]
  -owner, -ow        リポジトリオーナー [必須]
  -repo, -r          リポジトリ名 [必須]
  -state, -s         イシューの状態 (open, closed, all) [任意]
  -sort, -so         ソート項目 (created, updated, comments) [任意]
  -direction, -d     ソート方向 (asc, desc) [任意]
  -per-page, -pp     ページあたりの件数 (デフォルト: 30) [任意]
  -page, -p          ページ番号 (デフォルト: 1) [任意]
  -issue-number, -in イシュー番号（特定イシュー取得用） [任意]
  -help, -h          このヘルプを表示

例:
  # 基本的な使用方法（イシュー一覧取得）
  %s -o list-issues -ow octocat -r Hello-World

  # 特定のイシューを取得
  %s -o list-issues -ow octocat -r Hello-World -in 123

  # 状態とソートを指定
  %s -o list-issues -ow octocat -r Hello-World -s open -so created -d desc

  # ページネーション
  %s -o list-issues -ow octocat -r Hello-World -pp 10 -p 2

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
