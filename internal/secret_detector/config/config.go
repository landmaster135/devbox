package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/landmaster135/devbox/internal/secret_detector/domain"
)

const (
	helpMsgOfVersion    = "バージョン情報を表示"
	helpMsgOfVerbose    = "詳細な出力を表示"
	helpMsgOfConfigFile = "特定の設定ファイルのみをチェック"
	helpMsgOfDryRun     = "実際のGit操作なしでテスト実行"
)

// CLIConfig はCLIツールの設定を保持する構造体
type CLIConfig struct {
	Version    bool   // バージョン表示
	Verbose    bool   // 詳細出力モード
	ConfigFile string // 特定の設定ファイルパス
	DryRun     bool   // ドライランモード
	Help       bool   // ヘルプ表示
}

// ParseFlags はコマンドライン引数を解析してCLIConfigを返す
func ParseFlags() (*CLIConfig, error) {
	cfg := &CLIConfig{}

	flag.BoolVar(&cfg.Version, "version", false, helpMsgOfVersion)
	flag.BoolVar(&cfg.Verbose, "verbose", false, helpMsgOfVerbose)
	flag.StringVar(&cfg.ConfigFile, "config-file", "", helpMsgOfConfigFile)
	flag.BoolVar(&cfg.DryRun, "dry-run", false, helpMsgOfDryRun)
	flag.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")

	flag.Parse()

	// バリデーション
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate は設定の妥当性を検証する
func (c *CLIConfig) validate() error {
	// ヘルプまたはバージョンが要求された場合はバリデーションをスキップ
	if c.Help || c.Version {
		return nil
	}

	// 特定ファイル指定時の存在チェック
	if c.ConfigFile != "" {
		if _, err := os.Stat(c.ConfigFile); os.IsNotExist(err) {
			return fmt.Errorf("指定された設定ファイルが存在しません: %s", c.ConfigFile)
		}
	}

	return nil
}

// PrintUsage はヘルプメッセージを表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: secret-detector [オプション]\n\n")
	fmt.Fprintf(os.Stderr, "説明:\n")
	fmt.Fprintf(os.Stderr, "  Git pre-commit hook用のシークレット検知ツールです。\n")
	fmt.Fprintf(os.Stderr, "  JSON設定ファイル内の機密情報を自動検知し、コミット前にブロックします。\n\n")
	fmt.Fprintf(os.Stderr, "オプション:\n")
	fmt.Fprintf(os.Stderr, "  -version\n")
	fmt.Fprintf(os.Stderr, "        "+helpMsgOfVersion+"\n")
	fmt.Fprintf(os.Stderr, "  -verbose\n")
	fmt.Fprintf(os.Stderr, "        "+helpMsgOfVerbose+"\n")
	fmt.Fprintf(os.Stderr, "  -config-file string\n")
	fmt.Fprintf(os.Stderr, "        "+helpMsgOfConfigFile+"\n")
	fmt.Fprintf(os.Stderr, "  -dry-run\n")
	fmt.Fprintf(os.Stderr, "        "+helpMsgOfDryRun+"\n")
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        ヘルプを表示\n\n")
	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  # 通常実行（Git pre-commit hookとして）\n")
	fmt.Fprintf(os.Stderr, "  secret-detector\n\n")
	fmt.Fprintf(os.Stderr, "  # 詳細出力で実行\n")
	fmt.Fprintf(os.Stderr, "  secret-detector -verbose\n\n")
	fmt.Fprintf(os.Stderr, "  # 特定ファイルのみをチェック\n")
	fmt.Fprintf(os.Stderr, "  secret-detector -config-file=config.json\n\n")
	fmt.Fprintf(os.Stderr, "  # ドライランモード（Git操作なし）\n")
	fmt.Fprintf(os.Stderr, "  secret-detector -dry-run\n\n")
	fmt.Fprintf(os.Stderr, "検知対象:\n")
	fmt.Fprintf(os.Stderr, "  - 疑わしいキー名: api_key, secret_key, access_token など\n")
	fmt.Fprintf(os.Stderr, "  - 実際のシークレットパターン: OpenAI APIキー, GitHub PAT など\n")
	fmt.Fprintf(os.Stderr, "  - 高エントロピー値: ランダムな文字列\n\n")
	fmt.Fprintf(os.Stderr, "対応ファイル形式:\n")
	fmt.Fprintf(os.Stderr, "  - *.json, *.config.js, *.config.ts\n")
	fmt.Fprintf(os.Stderr, "  - mcp_settings.json, claude_desktop_config.json, cline_mcp_settings.json\n")
}

// LoadConfig はJSONファイルから設定を読み込む
func LoadConfig(filename string) (*domain.Config, error) {
	// ファイルが存在するかチェック
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filename)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var config domain.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid JSON format in %s: %w", filename, err)
	}

	return &config, nil
}
