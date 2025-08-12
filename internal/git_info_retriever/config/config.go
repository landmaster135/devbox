package config

import (
	"flag"
	"fmt"
)

// Config はCLIツールの設定を保持する構造体
type Config struct {
	Operation              string // 必須: 操作タイプ（retrieve, archive）
	Service                string // 必須: サービスタイプ（現在は"github"のみサポート）
	Token                  string // 必須: GitHubアクセストークン
	SaveFilePath           string // オプション: 結果を保存するファイルパス
	OutputCommandFilePath  string // オプション: Bash関数出力ファイルパス
	ArchiveDir             string // オプション: アーカイブディレクトリ
	SrcFile                string // オプション: 既存ファイルからリポジトリ情報を読み込み
	Help                   bool   // ヘルプ表示フラグ
}

// ParseFlags はコマンドライン引数を解析してConfigを返す
func ParseFlags() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.Operation, "operation", "", "操作タイプ（必須）: retrieve, archive")
	flag.StringVar(&cfg.Service, "service", "", "サービスタイプ（必須）: github")
	flag.StringVar(&cfg.Token, "token", "", "GitHubアクセストークン（必須）")
	flag.StringVar(&cfg.SaveFilePath, "save-file-path", "", "結果を保存するファイルパス（オプション）")
	flag.StringVar(&cfg.OutputCommandFilePath, "output-command-file-path", "", "Bash関数出力ファイルパス（オプション）")
	flag.StringVar(&cfg.ArchiveDir, "archive-dir", "./archives", "アーカイブディレクトリ（オプション）")
	flag.StringVar(&cfg.SrcFile, "src-file", "", "既存ファイルからリポジトリ情報を読み込み（オプション）")
	flag.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")

	flag.Parse()

	// ヘルプが要求された場合は検証をスキップ
	if cfg.Help {
		return cfg, nil
	}

	// 必須パラメータの検証
	if cfg.Operation == "" {
		return nil, fmt.Errorf("operationパラメータは必須です")
	}

	if cfg.Service == "" {
		return nil, fmt.Errorf("serviceパラメータは必須です")
	}

	// archiveオペレーションでsrc-fileが指定されていない場合はtokenが必須
	if cfg.Operation == "archive" && cfg.SrcFile == "" && cfg.Token == "" {
		return nil, fmt.Errorf("archiveオペレーションでsrc-fileが未指定の場合、tokenパラメータは必須です")
	}

	// retrieveオペレーションの場合はtokenが必須
	if cfg.Operation == "retrieve" && cfg.Token == "" {
		return nil, fmt.Errorf("retrieveオペレーションの場合、tokenパラメータは必須です")
	}

	// サポートされている操作タイプの検証
	if cfg.Operation != "retrieve" && cfg.Operation != "archive" {
		return nil, fmt.Errorf("サポートされていない操作タイプです: %s（現在サポート: retrieve, archive）", cfg.Operation)
	}

	// サポートされているサービスタイプの検証
	if cfg.Service != "github" {
		return nil, fmt.Errorf("サポートされていないサービスタイプです: %s（現在サポート: github）", cfg.Service)
	}

	return cfg, nil
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Println("Git情報取得・アーカイブツール")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  git-info-retriever -operation <操作タイプ> -service <サービスタイプ> [オプション]")
	fmt.Println()
	fmt.Println("パラメータ:")
	fmt.Println("  -operation string")
	fmt.Println("        操作タイプ（必須）")
	fmt.Println("        retrieve: リポジトリ情報取得")
	fmt.Println("        archive: Bash関数生成（git clone + zip圧縮）")
	fmt.Println("  -service string")
	fmt.Println("        サービスタイプ（必須）")
	fmt.Println("        現在サポート: github")
	fmt.Println("  -token string")
	fmt.Println("        GitHubアクセストークン")
	fmt.Println("        retrieve操作では必須")
	fmt.Println("        archive操作でsrc-file未指定の場合は必須")
	fmt.Println("  -save-file-path string")
	fmt.Println("        結果を保存するファイルパス（オプション）")
	fmt.Println("        指定しない場合は標準出力に表示")
	fmt.Println("  -output-command-file-path string")
	fmt.Println("        Bash関数出力ファイルパス（オプション）")
	fmt.Println("        archive操作で使用、指定しない場合は標準出力に表示")
	fmt.Println("  -archive-dir string")
	fmt.Println("        アーカイブディレクトリ（オプション、デフォルト: ./archives）")
	fmt.Println("        archive操作で使用")
	fmt.Println("  -src-file string")
	fmt.Println("        既存ファイルからリポジトリ情報を読み込み（オプション）")
	fmt.Println("        archive操作で使用、指定時はtokenは不要")
	fmt.Println("  -help")
	fmt.Println("        このヘルプを表示")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  # リポジトリ情報取得（標準出力）")
	fmt.Println("  git-info-retriever -operation retrieve -service github -token ghp_xxxxxxxxxxxxxxxxxxxx")
	fmt.Println()
	fmt.Println("  # リポジトリ情報取得（ファイル保存）")
	fmt.Println("  git-info-retriever -operation retrieve -service github -token ghp_xxxxxxxxxxxxxxxxxxxx -save-file-path ./repos.json")
	fmt.Println()
	fmt.Println("  # Bash関数生成（GitHubから取得）")
	fmt.Println("  git-info-retriever -operation archive -service github -token ghp_xxxxxxxxxxxxxxxxxxxx -archive-dir ./my-archives")
	fmt.Println()
	fmt.Println("  # Bash関数生成（既存ファイルから）")
	fmt.Println("  git-info-retriever -operation archive -service github -src-file ./repos.json -output-command-file-path ./archive_commands.sh")
}
