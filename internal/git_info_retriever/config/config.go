package config

import (
	"flag"
	"fmt"
)

// Config はCLIツールの設定を保持する構造体
type Config struct {
	Service      string // 必須: サービスタイプ（現在は"github"のみサポート）
	Token        string // 必須: GitHubアクセストークン
	SaveFilePath string // オプション: 結果を保存するファイルパス
	Help         bool   // ヘルプ表示フラグ
}

// ParseFlags はコマンドライン引数を解析してConfigを返す
func ParseFlags() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.Service, "service", "", "サービスタイプ（必須）: github")
	flag.StringVar(&cfg.Token, "token", "", "GitHubアクセストークン（必須）")
	flag.StringVar(&cfg.SaveFilePath, "save-file-path", "", "結果を保存するファイルパス（オプション）")
	flag.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")

	flag.Parse()

	// ヘルプが要求された場合は検証をスキップ
	if cfg.Help {
		return cfg, nil
	}

	// 必須パラメータの検証
	if cfg.Service == "" {
		return nil, fmt.Errorf("serviceパラメータは必須です")
	}

	if cfg.Token == "" {
		return nil, fmt.Errorf("tokenパラメータは必須です")
	}

	// サポートされているサービスタイプの検証
	if cfg.Service != "github" {
		return nil, fmt.Errorf("サポートされていないサービスタイプです: %s（現在サポート: github）", cfg.Service)
	}

	return cfg, nil
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Println("Git情報取得ツール")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  git-info-retriever -service <サービスタイプ> -token <アクセストークン> [-save-file <ファイルパス>]")
	fmt.Println()
	fmt.Println("パラメータ:")
	fmt.Println("  -service string")
	fmt.Println("        サービスタイプ（必須）")
	fmt.Println("        現在サポート: github")
	fmt.Println("  -token string")
	fmt.Println("        GitHubアクセストークン（必須）")
	fmt.Println("  -save-file string")
	fmt.Println("        結果を保存するファイルパス（オプション）")
	fmt.Println("        指定しない場合は標準出力に表示")
	fmt.Println("  -help")
	fmt.Println("        このヘルプを表示")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  # 標準出力に表示")
	fmt.Println("  git-info-retriever -service github -token ghp_xxxxxxxxxxxxxxxxxxxx")
	fmt.Println()
	fmt.Println("  # ファイルに保存")
	fmt.Println("  git-info-retriever -service github -token ghp_xxxxxxxxxxxxxxxxxxxx -save-file ./output.json")
}
