package config

import (
	"flag"
	"fmt"
	"os"
)

// Config はCLIツールの設定を保持する構造体
type Config struct {
	InputFormat   string // json, csv, tsv (入力データの形式)
	OutputFormat  string // html, csv, tsv (出力データの形式)
	Input         string // 直接入力（JSON文字列、CSV文字列など）
	InputFilePath string // 入力ファイルパス
	Help          bool   // ヘルプ表示
}

// ParseFlags はコマンドライン引数を解析してConfigを返す
func ParseFlags() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.InputFormat, "input-format", "", "入力データの形式 (json, csv, tsv)")
	flag.StringVar(&cfg.OutputFormat, "output-format", "", "出力データの形式 (html, csv, tsv)")
	flag.StringVar(&cfg.Input, "input", "", "直接入力データ")
	flag.StringVar(&cfg.InputFilePath, "input-file-path", "", "入力ファイルのパス")
	flag.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")

	flag.Parse()

	// バリデーション
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate は設定の妥当性を検証する
func (c *Config) validate() error {
	// ヘルプが要求された場合はバリデーションをスキップ
	if c.Help {
		return nil
	}

	// 必須パラメータのチェック
	if c.InputFormat == "" {
		return fmt.Errorf("input-formatは必須です")
	}
	if c.OutputFormat == "" {
		return fmt.Errorf("output-formatは必須です")
	}

	// 排他的入力チェック
	if c.Input != "" && c.InputFilePath != "" {
		return fmt.Errorf("inputとinput-file-pathは同時に指定できません")
	}

	// 入力必須チェック
	if c.Input == "" && c.InputFilePath == "" {
		return fmt.Errorf("inputまたはinput-file-pathのいずれかを指定してください")
	}

	// 入力形式の妥当性チェック
	if !isValidInputFormat(c.InputFormat) {
		return fmt.Errorf("未対応の入力形式です: %s (対応形式: json, csv, tsv)", c.InputFormat)
	}

	// 出力形式の妥当性チェック
	if !isValidOutputFormat(c.OutputFormat) {
		return fmt.Errorf("未対応の出力形式です: %s (対応形式: html, csv, tsv)", c.OutputFormat)
	}

	return nil
}

// isValidInputFormat は入力形式が有効かどうかを判定する
func isValidInputFormat(format string) bool {
	validFormats := []string{"json", "csv", "tsv"}
	for _, valid := range validFormats {
		if format == valid {
			return true
		}
	}
	return false
}

// isValidOutputFormat は出力形式が有効かどうかを判定する
func isValidOutputFormat(format string) bool {
	validFormats := []string{"html", "csv", "tsv"}
	for _, valid := range validFormats {
		if format == valid {
			return true
		}
	}
	return false
}

// PrintUsage はヘルプメッセージを表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: data-converter [オプション]\n\n")
	fmt.Fprintf(os.Stderr, "オプション:\n")
	fmt.Fprintf(os.Stderr, "  -input-format string\n")
	fmt.Fprintf(os.Stderr, "        入力データの形式 (json, csv, tsv)\n")
	fmt.Fprintf(os.Stderr, "  -output-format string\n")
	fmt.Fprintf(os.Stderr, "        出力データの形式 (html, csv, tsv)\n")
	fmt.Fprintf(os.Stderr, "  -input string\n")
	fmt.Fprintf(os.Stderr, "        直接入力データ\n")
	fmt.Fprintf(os.Stderr, "  -input-file-path string\n")
	fmt.Fprintf(os.Stderr, "        入力ファイルのパス\n")
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        ヘルプを表示\n\n")
	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  # JSON文字列をHTMLテーブルに変換\n")
	fmt.Fprintf(os.Stderr, "  data-converter -input-format=json -output-format=html -input='[[\"A\",\"B\"],[\"1\",\"2\"]]'\n\n")
	fmt.Fprintf(os.Stderr, "  # CSVファイルをTSVに変換\n")
	fmt.Fprintf(os.Stderr, "  data-converter -input-format=csv -output-format=tsv -input-file-path=data.csv\n\n")
	fmt.Fprintf(os.Stderr, "対応する変換パターン:\n")
	fmt.Fprintf(os.Stderr, "  入力形式 → 出力形式:\n")
	fmt.Fprintf(os.Stderr, "  - json → html, csv, tsv\n")
	fmt.Fprintf(os.Stderr, "  - csv → html, json, tsv\n")
	fmt.Fprintf(os.Stderr, "  - tsv → html, json, csv\n")
}
